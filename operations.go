package main

import (
	"net/http"
)

type Query struct {
	Filters  filterType
	Excludes filterType
}

// util for api
func parseURLParameters(r *http.Request) (filterType, filterType) {
	filterMap := make(filterType)
	excludeMap := make(filterType)
	for k := range RegisterFuncMap {
		parameter, parameterFound := r.URL.Query()[k]
		if parameterFound {
			filterMap[k] = parameter
		}
		parameter, parameterFound = r.URL.Query()["!"+k]
		if parameterFound {
			excludeMap[k] = parameter
		}
	}
	return filterMap, excludeMap
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

func earlyExitFilter(items Items, filters filterType, excludes filterType, registerFuncs registerFuncType) {
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

func filtered(items Items, filters filterType, excludes filterType, operations GroupedOperations) Items {
	registerFuncs := operations.Funcs
	filteredItems := make(Items, 0, len(items))
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
