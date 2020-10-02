package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

	IndexQuery string
	IndexGiven bool
}

func (q Query) EarlyExit() bool {
	return q.LimitGiven && !q.PageGiven && !q.SortByGiven
}

func decodeUrl(s string) string {
	newS, err := url.QueryUnescape(s)
	if err != nil {
		fmt.Println("oh no error", err)
		return s
	}
	fmt.Println("decoded", s, newS)
	return newS
}

// util for api
func parseURLParameters(r *http.Request) Query {
	filterMap := make(filterType)
	excludeMap := make(filterType)
	anyMap := make(filterType)
	//TODO change query to be based on input.

	urlItems := r.URL.Query()

	for k := range RegisterFuncMap {
		parameter, parameterFound := urlItems[k]
		if parameterFound && parameter[0] != "" {
			newSl := make([]string, len(parameter))
			for i, v := range parameter {
				newSl[i] = decodeUrl(v)
			}
			//filterMap[k] = parameter
			filterMap[k] = newSl
		}
		parameter, parameterFound = urlItems["!"+k]
		if parameterFound && parameter[0] != "" {
			excludeMap[k] = parameter
		}
		parameter, parameterFound = urlItems["any_"+k]
		if parameterFound && parameter[0] != "" {
			anyMap[k] = parameter
		}
	}

	// TODO there must be better way
	page := 1
	pageStr, pageGiven := urlItems["page"]
	if pageGiven {
		page = intMoreDefault(pageStr[0], 1)
	}

	pageSize := 10
	pageSizeStr, pageSizeGiven := urlItems["pagesize"]
	if pageSizeGiven {
		pageSize = intMoreDefault(pageSizeStr[0], 1)
	}

	limit := 0
	limitStr, limitGiven := urlItems["limit"]
	if limitGiven {
		limit = intMoreDefault(limitStr[0], 1)
	}

	sortingL, sortingGiven := urlItems["sortby"]

	index := ""
	indexL, indexGiven := urlItems["search-uiot"]
	if indexGiven {
		index = indexL[0]
	}
	return Query{
		Filters:  filterMap,
		Excludes: excludeMap,
		Anys:     anyMap,

		Limit:      limit,
		LimitGiven: limitGiven,

		Page:          page,
		PageGiven:     pageGiven,
		PageSize:      pageSize,
		PageSizeGiven: pageSizeGiven,

		SortBy:      sortingL,
		SortByGiven: sortingGiven,

		IndexQuery: index,
		IndexGiven: indexGiven,
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

func filteredEarlyExitSingle(items Items, column string, operations GroupedOperations, query Query) []string {
	registerFuncs := operations.Funcs
	filteredItemsSet := make(map[string]bool)
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
		single := operations.Getters[column](item)
		filteredItemsSet[single] = true

		if len(filteredItemsSet) == stop {
			break
		}
	}
	results := []string{}
	for k := range filteredItemsSet {
		results = append(results, k)
	}
	return results
}

func runQuery(items Items, query Query, operations GroupedOperations) (Items, int64) {
	start := time.Now()
	var newItems Items

	if query.IndexGiven && len(STR_INDEX) > 0 {
		items = make(Items, 0)
		indices := INDEX.Lookup([]byte(query.IndexQuery), -1)
		seen := make(map[string]bool)
		for _, idx := range indices {
			key := getStringFromIndex(STR_INDEX, idx)
			if !seen[key] {
				seen[key] = true
				for _, item := range LOOKUP[key] {
					items = append(items, item)
				}
			}

		}
	}

	if query.EarlyExit() {
		newItems = filteredEarlyExit(items, operations, query)
	} else {
		newItems = filtered(items, operations, query)
	}
	diff := time.Now().Sub(start)
	return newItems, int64(diff) / int64(1000000)
}

func runTypeAheadQuery(items Items, column string, query Query, operations GroupedOperations) ([]string, int64) {
	start := time.Now()
	results := filteredEarlyExitSingle(items, column, operations, query)
	diff := time.Now().Sub(start)
	return results, int64(diff) / int64(1000000)
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

func getHeaderDataSlice(items []string, query Query, queryDuration int64) HeaderData {
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
