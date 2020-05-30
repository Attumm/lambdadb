package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"runtime"
	//	"runtime/debug"
	//	"github.com/pkg/profile"
)

type filterFuncc func(*Item, string) bool
type registerFuncType map[string]filterFuncc
type registerGroupByFunc map[string]func(*Item) string
type filterType map[string][]string

//Item as Example

type Item struct {
	ID        string   `json:"id"`
	Value     string   `json:"value"`
	Type      string   `json:"type"`
	Name      string   `json:"name"`
	Vendor    string   `json:"vendor"`
	IP        string   `json:"ip"`
	Dn        []string `json:"dn"`
	ValueType string   `json:"value_type"`
	Country   string   `json:"country"`
}

//Items as Example
type Items []*Item

type ItemsGroupedBy map[string]Items

// Filter Functions
func filterValueContains(i *Item, s string) bool {
	return strings.Contains(i.Value, s)
}

func filterValueStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Value, s)
}

func filterValueMatch(i *Item, s string) bool {
	return i.Value == s
}

func filterNameMatch(i *Item, s string) bool {
	return i.Name == s
}

func filterIPMatch(i *Item, s string) bool {
	return i.IP == s
}

func filterCountryMatch(i *Item, s string) bool {
	return i.Country == s
}

func filterDNMatch(i *Item, s string) bool {
	return strings.Join(i.Dn, ".") == s
}

func filterValueContainsCase(i *Item, s string) bool {
	return strings.Contains(strings.ToLower(i.Value), strings.ToLower(s))
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

func filtered(items Items, filters filterType, excludes filterType, registerFuncs registerFuncType) Items {
	filteredItems := make(Items, 0, 100000)
	for _, item := range items {
		if !all(item, filters, registerFuncs) {
			continue
		}
                //speed
                //if !exclude(item, excludes, registerFuncs) {
                //    #	continue
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

var registerFuncMap registerFuncType
var registerGroupBy registerGroupByFunc
var itemsUnfilterd Items

func init() {
        SETTINGS.Set("http_db_host", "0.0.0.0:8000", "host with port")
        SETTINGS.Parse()

	registerFuncMap = make(registerFuncType)
	registerFuncMap["contains"] = filterValueContains
	registerFuncMap["contains-case"] = filterValueContainsCase
	registerFuncMap["match"] = filterValueMatch
	registerFuncMap["starts-with"] = filterValueStartsWith

	registerGroupBy = make(registerGroupByFunc)
	registerGroupBy["name"] = groupByName
	registerGroupBy["ip"] = groupByIP
	registerGroupBy["id"] = groupByID
	registerGroupBy["country"] = groupByCountry

	itemsUnfilterd = make(Items, 0, 100*1000)

}

// API

func listRest(w http.ResponseWriter, r *http.Request) {
	filterMap, excludeMap := parseURLParameters(r)

	items := filtered(itemsUnfilterd, filterMap, excludeMap, registerFuncMap)
	items = sortLimit(w, r, items)

	w.Header().Set("Content-Type", "application/json")
        //w.Header().Set("Total-Items", strconv.Itoa(len(items)))
	w.WriteHeader(http.StatusOK)

	groupByS, groupByFound := r.URL.Query()["groupby"]
	if !groupByFound {
		json.NewEncoder(w).Encode(items)
		return
	}

	groupByItems := groupByRunner(items, groupByS[0])
	json.NewEncoder(w).Encode(groupByItems)
	go func() {
		time.Sleep(2 * time.Second)
		runtime.GC()
	}()
}

func addRest(w http.ResponseWriter, r *http.Request) {
	jsonDecoder := json.NewDecoder(r.Body)
	var items Items
	err := jsonDecoder.Decode(&items)
	if err != nil {
		fmt.Println(err)
	}
	for _, item := range items {
		itemsUnfilterd = append(itemsUnfilterd, item)
	}
	w.WriteHeader(204)
}

func rmRest(w http.ResponseWriter, r *http.Request) {
	itemsUnfilterd = make(Items, 0, 100*1000)
	go func() {
		time.Sleep(1 * time.Second)
		runtime.GC()
	}()
	w.WriteHeader(204)
}

func helpRest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := make(map[string][]string)
	registeredFilters := []string{}
	for k := range registerFuncMap {
		registeredFilters = append(registeredFilters, k)
	}
	registeredGroupbys := []string{}
	for k := range registerGroupBy {
		registeredGroupbys = append(registeredGroupbys, k)
	}
	response["filters"] = registeredFilters
	response["groupby"] = registeredGroupbys
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

// util for api
func parseURLParameters(r *http.Request) (filterType, filterType) {
	filterMap := make(filterType)
	excludeMap := make(filterType)
	for k := range registerFuncMap {
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

// Group by functions

func groupByName(i *Item) string {
	return i.Name
}

func groupByIP(i *Item) string {
	return i.IP
}

func groupByID(i *Item) string {
	return i.ID
}

func groupByCountry(i *Item) string {
	return i.Country
}

func groupByRunner(items Items, groubByParameter string) ItemsGroupedBy {
	grouping := make(ItemsGroupedBy)
	groupingFunc := registerGroupBy[groubByParameter]
	if groupingFunc == nil {
		return grouping
	}
	for _, item := range items {
		GroupingKey := groupingFunc(item)
		grouping[GroupingKey] = append(grouping[GroupingKey], item)
	}
	return grouping
}

func main() {
	//defer profile.Start(profile.MemProfile).Stop()
	//go runPrintMem()
	//go func() {
	//	for {
	//		time.Sleep(5 * time.Second)
	//		runtime.GC()
	//	}
	//}()
	//defer profile.Start().Stop()
	//amount := 125000000	//100M
	//amount := 10000	// 10M
	//amount := 1000000	// 1M
	//itemsUnfilterd = createItems(amount)
	runserver := true
	groupedByExample := false
	runScript := false
	runProblem := false

	if runserver {
		ip_port := SETTINGS.Get("http_db_host")
		http.HandleFunc("/", listRest)
		http.HandleFunc("/help/", helpRest)
		http.HandleFunc("/add/", addRest)
		http.HandleFunc("/rm/", rmRest)
		fmt.Println("starting server", ip_port, " with:", len(itemsUnfilterd), "items")
		log.Fatal(http.ListenAndServe(ip_port, nil))
	}
	if runScript {

		//Examples how to use filtered function

		filterMap := make(filterType)
		//filterMap["ideven"] = []string{""}
		//filterMap["namecontains"] = []string{"", "8"}
		//filterMap["before"] = []string{"2010-06-01"}
		filterMap["after"] = []string{"0001-03-01"}
		//filterMap["match"] = []string{"2840-03-28"}
		//filterMap["FInName"] = []string{"a", "9"}

		exludeMap := make(filterType)
		items := filtered(itemsUnfilterd, filterMap, exludeMap, registerFuncMap)
		var itemJSON []byte
		//groupByItems := groupByRunner(items, groupByS[0])
		if !groupedByExample {
			itemJSON, _ = json.Marshal(items)
		} else {
			groupByItems := groupByRunner(items, "year")
			itemJSON, _ = json.Marshal(groupByItems)
		}
		fmt.Println(string(itemJSON), "amount:", len(items))
	}

	runtime.GC()
	fmt.Println("start first")
	runtime.GC()
	if runProblem {
		time.Sleep(1 * time.Second)
		for i := 0; i < 1; i++ {
			time.Sleep(100 * time.Millisecond)
			go problem(i)
		}
	}
	time.Sleep(20 * time.Second)
	//debug.FreeOSMemory()
	fmt.Println("start second")
	if runProblem {
		time.Sleep(10 * time.Second)
		for i := 0; i < 1; i++ {
			time.Sleep(1000 * time.Millisecond)
			//runtime.GC()
			go problem(i)
		}
	}
	time.Sleep(20 * time.Second)
	fmt.Println("Done!!!")
	//debug.FreeOSMemory()
	time.Sleep(30 * time.Second)
}

func problem(workerID int) {
	filterMap := make(filterType)
	exludeMap := make(filterType)
	items := filtered(itemsUnfilterd, filterMap, exludeMap, registerFuncMap)
	fmt.Println(workerID, len(items))
	//for i, _ := range items {
	//	items[i] = nil
	//}
	//items = nil
	runtime.GC()
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func runPrintMem() {
	for {
		PrintMemUsage()
		time.Sleep(1 * time.Second)
	}
}

func sortLimit(w http.ResponseWriter, r *http.Request, items Items) Items {

	pageStr, pageGiven := r.URL.Query()["page"]
	pageSizeStr, pageSizeGiven := r.URL.Query()["pagesize"]
	limitStr, limitGiven := r.URL.Query()["limit"]

	if !limitGiven && !pageGiven {
		return items
	}

	limit := len(items)
	if limitGiven {
		limit = intMoreDefault(limitStr[0], 1)
	}

	pageSize := 10
	if pageSizeGiven {
		pageSize = intMoreDefault(pageSizeStr[0], 1)
	}

	page := 1
	if pageGiven {
		page = intMoreDefault(pageStr[0], 1)
	}

	w.Header().Set("Page", strconv.Itoa(page))
	w.Header().Set("Page-Size", strconv.Itoa(pageSize))
	w.Header().Set("Total-Items", strconv.Itoa(len(items)))
        w.Header().Set("Total-Pages", strconv.Itoa(len(items)/pageSize))
	//TODO fix below
	start := (page - 1) * pageSize
	end := start + pageSize
	if end > len(items) {
		end = len(items)
	}
	if len(items) <= limit {
		return items[start:end]
	}
	items = items[start:end]
	if len(items) < limit {
		return items
	}
	return items[:limit]
}

func sortBy(items Items, r *http.Request) Items {
	sortFuncs := map[string]func(int, int) bool{
		"id":     func(i, j int) bool { return items[i].ID < items[j].ID },
		"-id":    func(i, j int) bool { return items[i].ID > items[j].ID },
		"name":   func(i, j int) bool { return items[i].Name < items[j].Name },
		"-name":  func(i, j int) bool { return items[i].Name > items[j].Name },
		"value":  func(i, j int) bool { return items[i].Value < items[j].Value },
		"-value": func(i, j int) bool { return items[i].Value > items[j].Value },

		"type":  func(i, j int) bool { return items[i].Type < items[j].Type },
		"-type": func(i, j int) bool { return items[i].Type > items[j].Type },

		"vendor":  func(i, j int) bool { return items[i].Vendor < items[j].Vendor },
		"-vendor": func(i, j int) bool { return items[i].Vendor > items[j].Vendor },

		"ip":  func(i, j int) bool { return items[i].IP < items[j].IP },
		"-ip": func(i, j int) bool { return items[i].IP > items[j].IP },

		//"dn":  func(i, j int) bool { return items[i].Dn < items[j].Dn },
		//"-dn": func(i, j int) bool { return items[i].Dn > items[j].Dn },

		"valueType":  func(i, j int) bool { return items[i].ValueType < items[j].ValueType },
		"-valueType": func(i, j int) bool { return items[i].ValueType > items[j].ValueType },

		"country":  func(i, j int) bool { return items[i].Country < items[j].Country },
		"-country": func(i, j int) bool { return items[i].Country > items[j].Country },
	}
	sortingL, sortingGiven := r.URL.Query()["sorting"]
	if sortingGiven {
		for _, sortFuncName := range sortingL {
			sortFunc := sortFuncs[sortFuncName]
			sort.Slice(items, sortFunc)
		}
	}
	return items
}

// Parse int from string
// if parsed int or error return default value
func intMoreDefault(s string, defaultN int) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	if n < defaultN {
		return defaultN
	}
	return n
}
