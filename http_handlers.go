package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"sort"
	"strconv"
	"time"
)

// API

func contextListRest(JWTConig jwtConfig, itemChan ItemsChannel, operations GroupedOperations) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query := parseURLParameters(r)
		fmt.Println("request", r.URL, "items", len(ITEMS))
		items := runQuery(ITEMS, query, operations)

		headerData := getHeaderData(items, query)

		if !query.EarlyExit() {
			items = sortLimit(items, query)
		}

		w.Header().Set("Content-Type", "application/json")
		for key, val := range headerData {
			w.Header().Set(key, val)
		}

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
}

func ItemChanWorker(itemChan ItemsChannel) {
	for items := range itemChan {
		for _, item := range items {
			ITEMS = append(ITEMS, item)
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
		}
		fmt.Println("adding", len(items))
		itemChan <- items
		w.WriteHeader(204)
	}
}

func rmRest(w http.ResponseWriter, r *http.Request) {
	ITEMS = make(Items, 0, 100*1000)
	go func() {
		time.Sleep(1 * time.Second)
		runtime.GC()
	}()
	w.WriteHeader(204)
}

func loadRest(w http.ResponseWriter, r *http.Request) {
	filename := "ITEMS.txt.gz"
	fi, err := os.Open(filename)
	if err != nil {
		return
	}
	defer fi.Close()

	fz, err := gzip.NewReader(fi)
	if err != nil {
		return
	}
	defer fz.Close()

	s, err := ioutil.ReadAll(fz)
	if err != nil {
		return
	}

	// TODO find out why wihout GC memory keeps growing
	runtime.GC()
	ITEMS = make(Items, 0, 100*100)
	runtime.GC()
	json.Unmarshal(s, &ITEMS)
	go func() {
		time.Sleep(1 * time.Second)
		runtime.GC()
	}()
	return
}

func saveRest(w http.ResponseWriter, r *http.Request) {
	filename := "ITEMS.txt.gz"
	var b bytes.Buffer
	writer := gzip.NewWriter(&b)
	itemJSON, _ := json.Marshal(ITEMS)
	writer.Write(itemJSON)
	writer.Flush()
	writer.Close()
	err := ioutil.WriteFile(filename, b.Bytes(), 0666)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(204)
}

func helpRest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := make(map[string][]string)
	registeredFilters := []string{}
	for k := range RegisterFuncMap {
		registeredFilters = append(registeredFilters, k)
	}

	registeredGroupbys := []string{}
	for k := range RegisterGroupBy {
		registeredGroupbys = append(registeredGroupbys, k)
	}

	_, registeredSortings := sortBy(ITEMS, []string{})

	sort.Strings(registeredFilters)
	sort.Strings(registeredGroupbys)
	sort.Strings(registeredSortings)
	response["filters"] = registeredFilters
	response["groupby"] = registeredGroupbys
	response["sortby"] = registeredSortings
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
		fmt.Sprintf("filtering: http://%s/list/?search=10&ontains-case=144&contains-case=10&page=1&pagesize=1", host),
	}
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}
