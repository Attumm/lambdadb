package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type Query struct {
	Filters  filterType
	Excludes filterType
	Anys     filterType

	Limit         int
	LimitGiven    bool
	Page          int
	PageGiven     bool
	PageSize      int
	PageSizeGiven bool

	SortBy      []string
	SortByGiven bool
}

func (q Query) EarlyExit() bool {
	return q.LimitGiven && !q.PageGiven && !q.SortByGiven
}

// util for api
func parseURLParameters(r *http.Request, JWTConfig jwtConfig) Query {
	filterMap := make(filterType)
	excludeMap := make(filterType)
	anyMap := make(filterType)
	//TODO speedup gains are present
	for k := range RegisterFuncMap {
		parameter, parameterFound := r.URL.Query()[k]
		if parameterFound {
			filterMap[k] = parameter
		}
		parameter, parameterFound = r.URL.Query()["!"+k]
		if parameterFound {
			excludeMap[k] = parameter
		}
		parameter, parameterFound = r.URL.Query()["any_"+k]
		if parameterFound {
			anyMap[k] = parameter
		}
	}

	// TODO there must be better way
	page := 1
	pageStr, pageGiven := r.URL.Query()["page"]
	if pageGiven {
		page = intMoreDefault(pageStr[0], 1)
	}

	pageSize := 10
	pageSizeStr, pageSizeGiven := r.URL.Query()["pagesize"]
	if pageSizeGiven {
		pageSize = intMoreDefault(pageSizeStr[0], 1)
	}

	limit := 0
	limitStr, limitGiven := r.URL.Query()["limit"]
	if limitGiven {
		limit = intMoreDefault(limitStr[0], 1)
	}

	sortingL, sortingGiven := r.URL.Query()["sortby"]
	return Query{
		Filters:  filterMap,
		Excludes: excludeMap,
		Anys:     anyMap,

		Limit:         limit,
		LimitGiven:    limitGiven,
		Page:          page,
		PageGiven:     pageGiven,
		PageSize:      pageSize,
		PageSizeGiven: pageSizeGiven,

		SortBy:      sortingL,
		SortByGiven: sortingGiven,
	}
}

func groupByRunner(items Items, groubByParameter string) ItemsGroupedBy {
	grouping := make(ItemsGroupedBy)
	groupingFunc := RegisterGroupBy[groubByParameter]
	if groupingFunc == nil {
		return grouping
	}
	for _, item := range items {
		GroupingKey := groupingFunc(item)
		grouping[GroupingKey] = append(grouping[GroupingKey], item)
	}
	return grouping
}

//Runner of filter functions, Item Should pass all
func all(item *Item, filters filterType, registerFuncs registerFuncType) bool {
	for funcName, args := range filters {
		filterFunc := registerFuncs[funcName]
		if filterFunc == nil {
			continue
		}
		for _, arg := range args {
			if !filterFunc(item, arg) {
				return false
			}
		}
	}
	return true
}

//Runner of filter functions, Item Should pass all
func any(item *Item, filters filterType, registerFuncs registerFuncType) bool {
	if len(filters) == 0 {
		return true
	}
	for funcName, args := range filters {
		filterFunc := registerFuncs[funcName]
		if filterFunc == nil {
			continue
		}
		for _, arg := range args {
			if filterFunc(item, arg) {
				return true
			}
		}
	}
	return false
}

//Runner of exlude functions, Item Should pass all
func exclude(item *Item, excludes filterType, registerFuncs registerFuncType) bool {
	for funcName, args := range excludes {
		excludeFunc := registerFuncs[funcName]
		if excludeFunc == nil {
			continue
		}
		for _, arg := range args {
			if excludeFunc(item, arg) {
				return false
			}
		}
	}
	return true
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func filteredEarlyExit(items Items, operations GroupedOperations, query Query) Items {
	registerFuncs := operations.Funcs
	filteredItems := make(Items, 0, len(items)/4)
	excludes := query.Excludes
	filters := query.Filters
	anys := query.Anys

	limit := query.Limit
	start := (query.Page - 1) * query.PageSize
	end := start + query.PageSize
	stop := end
	if query.LimitGiven {
		stop = limit
	}

	//TODO candidate for speedup
	for _, item := range items {
		if !any(item, anys, registerFuncs) {
			continue
		}
		if !all(item, filters, registerFuncs) {
			continue
		}
		if !exclude(item, excludes, registerFuncs) {
			continue
		}
		filteredItems = append(filteredItems, item)

		if len(filteredItems) == stop {
			break
		}
	}
	return filteredItems
}

func runQuery(items Items, query Query, operations GroupedOperations) (Items, int64) {
	start := time.Now()
	var newItems Items
	if query.EarlyExit() {
		newItems = filteredEarlyExit(items, operations, query)
	} else {
		newItems = filtered(items, operations, query)
	}
	diff := time.Now().Sub(start)
	return newItems, int64(diff) / int64(1000000)
}

func filtered(items Items, operations GroupedOperations, query Query) Items {
	registerFuncs := operations.Funcs
	excludes := query.Excludes
	filters := query.Filters
	anys := query.Anys

	filteredItems := make(Items, 0)
	for _, item := range items {
		if !any(item, anys, registerFuncs) {
			continue
		}
		if !all(item, filters, registerFuncs) {
			continue
		}
		if !exclude(item, excludes, registerFuncs) {
			continue
		}
		filteredItems = append(filteredItems, item)
	}
	return filteredItems
}

func mapIndex(items Items, indexes []int) Items {
	o := Items{}
	for _, index := range indexes {
		o = append(o, items[index])
	}
	return o
}

type HeaderData map[string]string

func getHeaderData(items Items, query Query, queryDuration int64) HeaderData {
	headerData := make(HeaderData)

	if query.LimitGiven {
		headerData["Limit"] = strconv.Itoa(query.Limit)
	}

	if query.PageGiven {
		headerData["Page"] = strconv.Itoa(query.Page)
		headerData["Page-Size"] = strconv.Itoa(query.PageSize)
		headerData["Total-Pages"] = strconv.Itoa((len(items) / query.PageSize) + 1)
	}

	headerData["Total-Items"] = strconv.Itoa(len(items))
	headerData["Query-Duration"] = strconv.FormatInt(queryDuration, 10) + "ms"
	bytesQuery, _ := json.Marshal(query)
	headerData["query"] = string(bytesQuery)

	return headerData
}

func sortLimit(items Items, query Query) Items {
	count := len(items)
	if count == 0 {
		return items
	}

	if query.SortByGiven {
		items, _ = sortBy(items, query.SortBy)
	}

	if !query.LimitGiven && !query.PageGiven {
		return items
	}

	//TODO there should be nicer way
	start := (query.Page - 1) * query.PageSize
	end := start + query.PageSize

	items = items[min(start, count):min(end, count)]
	if !query.LimitGiven {
		return items
	}

	// Note the slice built on array, slicing a slice larger then the the slice adds array items
	if len(items) < query.Limit {
		return items
	}
	return items[:query.Limit]
}
