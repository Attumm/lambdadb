package main

import (
	"encoding/csv"
	"fmt"
	"index/suffixarray"
	"log"
	"net/http"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

// API list headers
func setHeader(items Items, w http.ResponseWriter, query Query, queryTime int64) {

	headerData := getHeaderData(items, query, queryTime)

	for key, val := range headerData {
		w.Header().Set(key, val)
	}

	if query.ReturnFormat == "csv" {
		w.Header().Set("Content-Disposition", "attachment; filename=\"items.csv\"")
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	} else {
		w.Header().Set("Content-Type", "application/json")
	}

	if len(items) > 0 {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func contextListRest(JWTConig jwtConfig, itemChan ItemsChannel, operations GroupedOperations) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query := parseURLParameters(r)

		items, queryTime := runQuery(&ITEMS, query, operations)

		msg := fmt.Sprint("total: ", len(ITEMS), " hits: ", len(items), " time: ", queryTime, "ms ", "url: ", r.URL)
		fmt.Printf(NoticeColorN, msg)

		if !query.EarlyExit() {
			items = sortLimit(items, query)
		}

		setHeader(items, w, query, queryTime)

		groupByS, groupByFound := r.URL.Query()["groupby"]

		if !groupByFound {
			if query.ReturnFormat == "csv" {
				writeCSV(items, w)
			} else {
				json.NewEncoder(w).Encode(items)
			}
			// force empty helps garbage collection
			items = nil
			return
		}

		groupByItems := groupByRunner(items, groupByS[0])
		items = nil

		reduceName, reduceFound := r.URL.Query()["reduce"]

		if reduceFound {
			result := make(map[string]map[string]string)
			reduceFunc, reduceFuncFound := operations.Reduce[reduceName[0]]
			if !reduceFuncFound {
				json.NewEncoder(w).Encode(result)
				return
			}
			for key, val := range groupByItems {
				result[key] = reduceFunc(val)
			}
			groupByItems = nil

			if len(result) == 0 {
				return
			}

			json.NewEncoder(w).Encode(result)
			return
		}

		json.NewEncoder(w).Encode(groupByItems)
	}
}

func ItemChanWorker(itemChan ItemsChannel) {
	label := 0
	for items := range itemChan {
		for _, itm := range items {
			ITEMS[label] = itm
			itm.GeoIndex(label)
			label += 1
		}
	}
}

func contextAddRest(JWTConig jwtConfig, itemChan ItemsChannel, operations GroupedOperations) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		jsonDecoder := json.NewDecoder(r.Body)
		var items Items
		err := jsonDecoder.Decode(&items)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}
		msg := fmt.Sprint("adding ", len(items))
		fmt.Printf(WarningColorN, msg)

		strictMode := SETTINGS.Get("strict-mode") == "y"
		for n, item := range items {
			if (*item == Item{}) {
				fmt.Printf("unable to process item %d of batch\n", n)
				if strictMode {
					fmt.Printf("strict mode stopping ingestion of batch\n")
					w.WriteHeader(406)
					return
				}
			}

		}
		itemChan <- items
		w.WriteHeader(204)
	}
}

func rmRest(w http.ResponseWriter, r *http.Request) {
	ITEMS = make(Items, 0, 100*1000)
	msg := fmt.Sprint("removed items from database")
	fmt.Printf(WarningColorN, msg)
	ITEMS = labeledItems{}

	go func() {
		time.Sleep(1 * time.Second)
		runtime.GC()
	}()
	w.WriteHeader(204)
}

var LOOKUP map[string]Items
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
func makeIndex() {

	sort.Slice(ITEMS, func(i, j int) bool {
		return ITEMS[i].GetIndex() < ITEMS[j].GetIndex()
	})

	LOOKUP = make(map[string]Items)
	LOOKUPINDEX = make(map[string][]int)
	kSet := make(map[string]bool)

	//TODO this still needs a cleanup, but it's currently the solution to solve column and the indexes
	//for _, item := range ITEMS {
	//	key := strings.ToLower(item.GetIndex())
	//	kSet[key] = true
	//	LOOKUP[key] = append(LOOKUP[key], item)
	//}

	for index, item := range ITEMS {
		for _, v := range item.Row() {
			v := strings.ToLower(v)
			kSet[v] = true
			//LOOKUP[v] = append(LOOKUP[v], item)
			LOOKUPINDEX[v] = append(LOOKUPINDEX[v], index)
		}
	}

	//make a list of all used keys
	keys := []string{}
	for key := range kSet {
		keys = append(keys, key)
	}

	//join all keys together
	STR_INDEX = []byte("\x00" + strings.Join(keys, "\x00") + "\x00")
	INDEX = suffixarray.New(STR_INDEX)
}

func writeCSV(items Items, w http.ResponseWriter) {
	writer := csv.NewWriter(w)
	for i := range items {
		writer.Write(items[i].Row())
		writer.Flush()
	}
}

func loadRest(w http.ResponseWriter, r *http.Request) {
	storagename, _, retrievefunc, filename := handleInputStorage(r)

	msg := fmt.Sprintf("retrieving with: %s, with filename: %s", storagename, filename)
	fmt.Printf(WarningColorN, msg)
	itemsAdded, err := retrievefunc(ITEMS, filename)
	if err != nil {
		log.Printf("could not open %s reason %s", filename, err)
		w.Write([]byte("500 - could not load data"))
	}

	msg = fmt.Sprint("Loaded new items in memory amount: ", itemsAdded)
	fmt.Printf(WarningColorN, msg)

	if SETTINGS.Get("indexed") == "y" {
		msg := fmt.Sprint("Creating index")
		fmt.Printf(WarningColorN, msg)
		makeIndex()
		msg = fmt.Sprint("Index set")
		fmt.Printf(WarningColorN, msg)
	}
}

func handleInputStorage(r *http.Request) (string, storageFunc, retrieveFunc, string) {
	urlPath := r.URL.Path
	storagename := SETTINGS.Get("STORAGEMETHOD")

	if len(urlPath) > len("/mgmt/save/") {
		storagename = urlPath[len("/mgmt/save/"):]
	}
	storagefunc, found := STORAGEFUNCS[storagename]
	if !found {
		storagename := SETTINGS.Get("STORAGEMETHOD")
		storagefunc = STORAGEFUNCS[storagename]
	}

	s, err := ioutil.ReadAll(fz)
	if err != nil {
		return
	}

	retrievefunc, found := RETRIEVEFUNCS[storagename]
	if !found {
		storagename := SETTINGS.Get("STORAGEMETHOD")
		retrievefunc = RETRIEVEFUNCS[storagename]
	}

	// empty exising ITEMS
	ITEMS = labeledItems{}
	json.Unmarshal(s, &ITEMS)

	filename := fmt.Sprintf("%s.%s", FILENAME, storagename)

	msg := fmt.Sprint("Loaded new items in memory amount: ", len(ITEMS))
	fmt.Printf(WarningColorN, msg)
	//makeIndex()
	BuildGeoIndex()
	return storagename, storagefunc, retrievefunc, filename
}

func saveRest(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("storing items %d", len(ITEMS))
	fmt.Printf(WarningColor, msg)

	fmt.Println("full", r.URL.Path)

	storagename, storagefunc, _, filename := handleInputStorage(r)
	msg = fmt.Sprintf("storage method: %s filename: %s\n", storagename, filename)
	fmt.Printf(WarningColor, msg)

	size, err := storagefunc(ITEMS, filename)
	if err != nil {
		fmt.Println("unable to write file reason:", err)
		w.WriteHeader(500)
		return

	}
	msg = fmt.Sprintf("filname %s, filesize: %d mb\n", filename, size/1024/1025)
	fmt.Printf(WarningColor, msg)

	w.WriteHeader(204)
}

func validColumn(column string, columns []string) bool {
	for _, item := range columns {
		if column == item {
			return true
		}
	}
	return false
}

// Other wise also known in mathematics as set but in http name it would be confused with the verb set.
//func UniqueValuesInColumn(w http.ResponseWriter, r *http.Request) {
//	column := r.URL.Path[1:]
//	response := make(map[string]string)
//	if len(ITEMS) == 0 {
//		response["message"] = fmt.Sprint("invalid input: ", column)
//		w.WriteHeader(400)
//		json.NewEncoder(w).Encode(response)
//		return
//
//	}
//	validColumns := ITEMS[0].Columns()
//
//	if !validColumn(column, validColumns) {
//		w.WriteHeader(400)
//
//		response["message"] = fmt.Sprint("invalid input: ", column)
//		response["input"] = column
//		response["valid input"] = strings.Join(validColumns, ", ")
//		json.NewEncoder(w).Encode(response)
//		return
//	}
//	set := make(map[string]bool)
//	for item := range ITEMS {
//		r := reflect.ValueOf(item)
//		value := reflect.Indirect(r).FieldByName(column)
//		valu
//		set[value.Str()] = true
//	}
//
//}
type ShowItem struct {
	IsShow bool   `json:"isShow"`
	Label  string `json:"label"`
	Name   string `json:"name"`
}

type Meta struct {
	Fields []ShowItem `json:"fields"`
	View   string     `json:"view"`
}

type searchResponse struct {
	Count int       `json:"count"`
	Data  ItemsFull `json:"data"`
	MMeta *Meta     `json:"meta"`
}

func makeResp(items Items) searchResponse {

	itemsfull := make(ItemsFull, 0, len(items))

	for _, oneitem := range items {
		orgItem := oneitem.Serialize()
		itemsfull = append(itemsfull, &orgItem)
	}

	fields := []ShowItem{}
	columns := ItemFull{}.Columns()
	for _, column := range columns {
		fields = append(fields, ShowItem{IsShow: true, Name: column, Label: column})
	}

	return searchResponse{
		Count: len(items),
		Data:  itemsfull,
		MMeta: &Meta{Fields: fields, View: "table"},
	}
}

func corsEnabled(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Expose-Headers", "*")

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST")
			w.Header().Set("Access-Control-Allow-Headers", "Page, Page-Size, Total-Pages, query, Total-Items, Query-Duration, Content-Type, X-CSRF-Token, Authorization")
			return
		} else {
			h.ServeHTTP(w, r)
		}
	})

}
func passThrough(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)

	})
}

func MIDDLEWARE(cors bool) func(http.Handler) http.Handler {
	if cors {
		return corsEnabled
	}

	return passThrough
}

func contextSearchRest(JWTConig jwtConfig, itemChan ItemsChannel, operations GroupedOperations) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query := parseURLParameters(r)

		items, queryTime := runQuery(&ITEMS, query, operations)
		if len(items) == 0 {
			w.WriteHeader(404)
			return
		}
		msg := fmt.Sprint("total: ", len(ITEMS), " hits: ", len(items), " time: ", queryTime, "ms ", "url: ", r.URL)
		fmt.Printf(NoticeColorN, msg)
		headerData := getHeaderData(items, query, queryTime)

		if !query.EarlyExit() {
			items = sortLimit(items, query)
		}

		w.Header().Set("Content-Type", "application/json")
		for key, val := range headerData {
			w.Header().Set(key, val)
		}

		w.WriteHeader(http.StatusOK)

		response := makeResp(items)
		json.NewEncoder(w).Encode(response)
		items = nil
	}
}

func contextTypeAheadRest(JWTConig jwtConfig, itemChan ItemsChannel, operations GroupedOperations) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query := parseURLParameters(r)
		column := r.URL.Path[len("/typeahead/"):]
		if column[len(column)-1] == '/' {
			column = column[:len(column)-1]
		}
		if _, ok := operations.Getters[column]; !ok {
			w.WriteHeader(404)
			w.Write([]byte("column is not found"))
			return
		}

		results, queryTime := runTypeAheadQuery(&ITEMS, column, query, operations)
		if len(results) == 0 {
			w.WriteHeader(404)
			return
		}
		msg := fmt.Sprint("total: ", len(ITEMS), " hits: ", len(results), " time: ", queryTime, "ms ", "url: ", r.URL)
		fmt.Printf(NoticeColorN, msg)
		headerData := getHeaderDataSlice(results, query, queryTime)

		w.Header().Set("Content-Type", "application/json")
		for key, val := range headerData {
			w.Header().Set(key, val)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(results)
		results = nil
	}
}

func helpRest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := make(map[string][]string)
	registeredFilters := []string{}
	for k := range RegisterFuncMap {
		registeredFilters = append(registeredFilters, k)
	}

	registeredExcludes := []string{}
	for k := range RegisterFuncMap {
		registeredExcludes = append(registeredExcludes, "!"+k)
	}

	registeredAnys := []string{}
	// TODO create const for the "any_" or exclude prefix
	for k := range RegisterFuncMap {
		registeredAnys = append(registeredAnys, "any_"+k)
	}

	registeredGroupbys := []string{}
	for k := range RegisterGroupBy {
		registeredGroupbys = append(registeredGroupbys, k)
	}

	registerReduces := []string{}
	for k := range RegisterReduce {
		registerReduces = append(registerReduces, k)
	}

	newItems := make(Items, 10)
	for i := 0; i < 10; i++ {
		newItems = append(newItems, ITEMS[i])
	}

	_, registeredSortings := sortBy(newItems, []string{})

	sort.Strings(registeredFilters)
	sort.Strings(registeredExcludes)
	sort.Strings(registeredAnys)
	sort.Strings(registeredGroupbys)
	sort.Strings(registeredSortings)
	sort.Strings(registerReduces)

	response["filters"] = registeredFilters
	response["exclude_filters"] = registeredExcludes
	response["anys_filters"] = registeredAnys
	response["groupby"] = registeredGroupbys
	response["sortby"] = registeredSortings
	response["reduce"] = registerReduces

	totalItems := strconv.Itoa(len(ITEMS))

	host := SETTINGS.Get("http_db_host")
	response["total-items"] = []string{totalItems}
	response["settings"] = []string{
		fmt.Sprintf("host: %s", host),
		fmt.Sprintf("JWT: %s", SETTINGS.Get("JWTENABLED")),
	}
	response["examples"] = []string{
		fmt.Sprintf("typeahead: http://%s/list/?typeahead=ams&limit=10", host),
		fmt.Sprintf("search: http://%s/list/?search=ams&page=1&pagesize=1", host),
		fmt.Sprintf("search with limit: http://%s/list/?search=10&page=1&pagesize=10&limit=5", host),
		fmt.Sprintf("sorting: http://%s/list/?search=100&page=10&pagesize=100&sortby=-country", host),
		fmt.Sprintf("filtering: http://%s/list/?search=10&ontains=144&contains-case=10&page=1&pagesize=1", host),
		fmt.Sprintf("groupby: http://%s/list/?search=10&contains-case=10&groupby=country", host),
		fmt.Sprintf("aggregation: http://%s/list/?search=10&contains-case=10&groupby=country&reduce=count", host),
		fmt.Sprintf("chain the same filters: http://%s/list/?search=10&contains-case=127&contains-case=0&contains-case=1", host),
		fmt.Sprintf("typeahead use the name of the column in this case IP: http://%s/typeahead/ip/?starts-with=127&limit=15", host),
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
