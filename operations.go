package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-spatial/geom"
	"github.com/go-spatial/geom/encoding/geojson"
	"net/http"
	"net/url"
	//"reflect"
	"errors"
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

	Geometry      geom.Geometry
	GeometryGiven bool

	ReturnFormat string
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
	return newS
}

// util for api
func parseURLParameters(r *http.Request) (Query, error) {
	filterMap := make(filterType)
	excludeMap := make(filterType)
	anyMap := make(filterType)

	//TODO change query to be based on input.

	urlItems := r.URL.Query()

	// we can post gejson data
	r.ParseForm()

	if SETTINGS.Get("debug") == "yes" {

		for key, value := range r.Form {
			fmt.Printf("F %s = %s\n", key, value)
		}
		for key, value := range urlItems {

			fmt.Printf("P %s = %s\n", key, value)
		}
	}

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

	format := "json"
	formatStr, formatGiven := urlItems["format"]

	if formatGiven {
		if formatStr[0] == "csv" {
			format = "csv"
		}
	}

	sortingL, sortingGiven := urlItems["sortby"]

	index := ""
	indexL, indexGiven := urlItems["search"]
	indexUsed := indexGiven && indexL[0] != ""

	if indexUsed {
		index = indexL[0]
	}

	// check for geojson geometry stuff.
	geometryS, geometryGiven := r.Form["geojson"]
	var geoinput geojson.Geometry
	if geometryGiven && geometryS[0] != "" {
		err := json.Unmarshal([]byte(geometryS[0]), &geoinput)
		if err != nil {
			fmt.Println("parsing geojson error")
			fmt.Println(err)
			geometryGiven = false
			return Query{}, errors.New("failed to parse geojson")
		}
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
		IndexGiven: indexUsed,

		Geometry: geoinput.Geometry,

		GeometryGiven: geometryGiven,

		ReturnFormat: format,
	}, nil
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

func filteredEarlyExit(items *labeledItems, operations GroupedOperations, query Query) Items {
	registerFuncs := operations.Funcs
	filteredItems := make(Items, 0, len(*items)/4)
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

	lock.RLock()
	defer lock.RUnlock()

	for _, item := range *items {
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

func filteredEarlyExitSingle(items *labeledItems, column string, operations GroupedOperations, query Query) []string {
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

	lock.RLock()
	defer lock.RUnlock()

	// TODO candidate for speedup
	for _, item := range *items {
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
		// empty keys crashes frontend.
		// should be fixed in frontend then below can go.
		// NOTE: add a special field so we can filter on 'nil' / empty values.
		if len(k) > 0 {
			results = append(results, k)
		}
	}
	return results
}

func runQuery(items *labeledItems, query Query, operations GroupedOperations) (Items, int64) {
	start := time.Now()
	var newItems Items

	if query.GeometryGiven {
		cu := CoverDefault(query.Geometry)
		if len(cu) == 0 {
			fmt.Println("covering cell union not created")
		} else {
			geoitems := SearchGeoItems(cu)
			items = &geoitems
			fmt.Println(len(geoitems))
		}
	}

	if query.EarlyExit() {
		newItems = filteredEarlyExit(items, operations, query)
	} else {
		newItems = filtered(items, operations, query)
	}

	diff := time.Since(start)
	return newItems, int64(diff) / int64(1000000)
}

func runTypeAheadQuery(
	items *labeledItems, column string, query Query,
	operations GroupedOperations) ([]string, int64) {

	start := time.Now()
	results := filteredEarlyExitSingle(items, column, operations, query)
	diff := time.Since(start)
	return results, int64(diff) / int64(1000000)
}

func filtered(items *labeledItems, operations GroupedOperations, query Query) Items {
	registerFuncs := operations.Funcs
	filteredItems := make(Items, 0)
	excludes := query.Excludes
	filters := query.Filters
	anys := query.Anys

	lock.RLock()
	defer lock.RUnlock()

	for _, item := range *items {
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

//getHeaderDataSlice extract from header information with data slice we want
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
