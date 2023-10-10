package main

import (
	"encoding/json"
	"fmt"
	. "github.com/Attumm/settingo/settingo"
	"index/suffixarray"
	"net/http"
	"net/url"
	"runtime"
	"sort"
	"strconv"
	"strings"
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
func parseURLParameters(r *http.Request) (Query, bool) {
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
	indexGiven = indexGiven && SETTINGS.GetBool("indexed")
	indexUsed := indexGiven && len(indexL[0]) > 2
	if indexUsed {
		index = strings.ToLower(indexL[0])
	}
	validQuery := true
	if SETTINGS.GetBool("JWTENABLED") {
		token, err := getJWT(r, SETTINGS.Get("JWTSECRET"), SETTINGS.Get("JWTHEADER"))
		if err != nil {
			validQuery = false
			fmt.Println("jwt token had a issue:", err)
			return Query{}, validQuery
		}
		column := SETTINGS.Get("JWTCOLUMN")
		columnValues := getColumnValues(token.Groups, SETTINGS.GetMap("JWTGROUPSTOVALUES"))
		if !containsWildCard(columnValues) {
			overrideAnyFilter(anyMap, column, columnValues)
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

		ReturnFormat: format,
	}, validQuery
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
		// empty keys crashes frontend.
		// should be fixed in frontend then below can go.
		// NOTE: add a special field so we can filter on 'nil' / empty values.
		if len(k) > 0 {
			results = append(results, k)
		}
	}
	return results
}

func runQuery(items Items, query Query, operations GroupedOperations) (Items, int64) {
	start := time.Now()
	var newItems Items

	//TODO this still needs a cleanup, but it's currently the solution to solve column and the indexes

	//if query.IndexGiven && len(STR_INDEX) > 0 {
	//	items = make(Items, 0)
	//	indices := INDEX.Lookup([]byte(query.IndexQuery), -1)
	//	seen := make(map[string]bool)
	//	for _, idx := range indices {
	//		key := getStringFromIndex(STR_INDEX, idx)
	//		if !seen[key] {
	//			seen[key] = true
	//			for _, item := range LOOKUP[key] {
	//				items = append(items, item)
	//			}
	//		}

	//	}
	//}
	if query.IndexGiven {
		items = runIndexQuery(query)
	}

	if query.EarlyExit() {
		newItems = filteredEarlyExit(items, operations, query)
	} else {
		newItems = filtered(items, operations, query)
	}
	diff := time.Since(start)
	return newItems, int64(diff) / int64(1000000)
}

func runTypeAheadQuery(items Items, column string, query Query, operations GroupedOperations) ([]string, int64) {
	start := time.Now()
	if query.IndexGiven {
		items = runIndexQuery(query)
	}
	results := filteredEarlyExitSingle(items, column, operations, query)
	diff := time.Since(start)
	return results, int64(diff) / int64(1000000)
}

func runIndexQuery(query Query) Items {
	items := make(Items, 0)
	indices := INDEX.Lookup([]byte(query.IndexQuery), -1)
	seen := make(map[string]bool)
	added := make(map[int]bool)
	for _, idx := range indices {
		key := getStringFromIndex(STR_INDEX, idx)
		seen[key] = true
	}
	for key := range seen {
		for _, index := range LOOKUPINDEX[key] {
			added[index] = true
		}
	}
	for index := range added {
		items = append(items, ITEMS[index])
	}
	return items
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
	// https://play.golang.org/p/GxhbBGNaXwL
	if len(items) < query.Limit {
		return items
	}
	return items[:query.Limit]
}

var LOOKUPINDEX map[string][]int
var INDEX *suffixarray.Index
var STR_INDEX []byte

const FILENAME = "./files/name"

func getStringFromIndex(data []byte, index int) string {
	var start, end int
	for i := index - 1; i >= 0; i-- {
		if data[i] == 0 {
			start = i + 1
			break
		}
	}
	for i := index + 1; i < len(data); i++ {
		if data[i] == 0 {
			end = i
			break
		}
	}
	return string(data[start:end])
}

//make an index on column values in dataset
var F_INDEX = "files/STR_INDEX"
var F_LOOKUP = "files/LOOKUPINDEX"

func makeIndex() {
	gcCount := SETTINGS.GetInt("INDEXEDGC")

	sort.Slice(ITEMS, func(i, j int) bool {
		return ITEMS[i].GetIndex() < ITEMS[j].GetIndex()
	})
	LOOKUPINDEX = make(map[string][]int)

	if SETTINGS.GetBool("INDEXSTORED") {
		fmt.Println("index is stored, trying to load")
		LOOKUPINDEX = DecodeMapStrSInt(ReadFromFile(F_LOOKUP))
		runtime.GC()
		STR_INDEX = ReadFromFile(F_INDEX)
		if len(LOOKUPINDEX) != 0 && len(STR_INDEX) != 0 {
			runtime.GC()
			INDEX = suffixarray.New(STR_INDEX)
			runtime.GC()
			fmt.Println("loading index from file is done")
			return
		}
		fmt.Println("failed to set indexes from files us")
	}

	//TODO this still needs a cleanup, but it's currently the solution to solve column and the indexes
	//for _, item := range ITEMS {
	//	key := strings.ToLower(item.GetIndex())
	//	kSet[key] = true
	//	LOOKUP[key] = append(LOOKUP[key], item)
	//}

	kSet := make(map[string]bool)
	for index, item := range ITEMS {
		if gcCount != 0 && index%gcCount == 0 {
			runtime.GC()
		}
		for _, v := range item.Row() {
			v := strings.ToLower(v)
			kSet[v] = true
			//LOOKUP[v] = append(LOOKUP[v], item)
			LOOKUPINDEX[v] = append(LOOKUPINDEX[v], index)
		}
	}

	runtime.GC()
	//make a list of all used keys
	keys := []string{}
	for key := range kSet {
		keys = append(keys, key)
	}

	runtime.GC()
	//join all keys together
	STR_INDEX = []byte("\x00" + strings.Join(keys, "\x00") + "\x00")

	WriteToFile(STR_INDEX, F_INDEX)
	WriteToFile(EncodeMapStrSInt(LOOKUPINDEX), F_LOOKUP)

	runtime.GC()
	INDEX = suffixarray.New(STR_INDEX)

	runtime.GC()
}
