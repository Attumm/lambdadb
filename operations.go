package main

import (
	"encoding/json"
	"net/http"
	"strconv"
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
	return !q.SortByGiven && q.LimitGiven
}

// util for api
func parseURLParameters(r *http.Request) Query {
	filterMap := make(filterType)
	excludeMap := make(filterType)
	anyMap := make(filterType)
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

func include_filter(items Items, filters filterType, excludes filterType, registerFuncs registerFuncType) {
	for _, item := range items {
		if !all(item, filters, registerFuncs) {
			continue
		}
		items = append(items, item)
	}
}

func exclude_filter(items Items, filters filterType, excludes filterType, registerFuncs registerFuncType) {
	for _, item := range items {
		if !all(item, filters, registerFuncs) {
			continue
		}
		if !exclude(item, excludes, registerFuncs) {
			continue
		}
		items = append(items, item)
	}
}

func filteredEarlyExit(items Items, filters filterType, operations GroupedOperations, query Query) Items {
	registerFuncs := operations.Funcs
	filteredItems := make(Items, 0, len(items))
	limit := query.Limit
	//excludes := query.Excludes
	for _, item := range items {
		//if !any(item, filters, registerFuncs) {
		//	continue
		//}
		if !all(item, filters, registerFuncs) {
			continue
		}
		//if !exclude(item, excludes, registerFuncs) {
		//	continue
		//}
		filteredItems = append(filteredItems, item)
		if len(filteredItems) == limit {
			break
		}
	}
	return filteredItems
}

func run_query(items Items, query Query, operations GroupedOperations) Items {

	filters := query.Filters
	if query.EarlyExit() {
		return filteredEarlyExit(items, filters, operations, query)
	}
	new_items := filtered(items, filters, operations, query)
	return new_items

}

func filtered(items Items, filters filterType, operations GroupedOperations, query Query) Items {
	registerFuncs := operations.Funcs
	filteredItems := make(Items, 0, len(items))
	//excludes := query.Excludes
	for _, item := range items {
		//if !any(item, filters, registerFuncs) {
		//	continue
		//}
		if !all(item, filters, registerFuncs) {
			continue
		}
		//if !exclude(item, excludes, registerFuncs) {
		//	continue
		//}
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

func getHeaderData(items Items, query Query) HeaderData {
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
	bQuery, _ := json.Marshal(query)
	headerData["query"] = string(bQuery)

	return headerData
}

func sortLimit(items Items, query Query) Items {

	if !query.LimitGiven && !query.PageGiven {
		return items
	}
	if query.SortByGiven {
		items, _ = sortBy(items, query.SortBy)
	}
	//TODO fix below
	start := (query.Page - 1) * query.PageSize
	end := start + query.PageSize
	if end > len(items) {
		end = len(items)
	}
	if query.LimitGiven && len(items) <= query.Limit {
		return items[start:end]
	}
	items = items[start:end]
	if len(items) < query.Limit {
		return items
	}
	return items[:query.Limit]
}
