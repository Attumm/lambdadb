package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sort"
	"strconv"
	"time"
)

type HeaderData map[string]string

func sortLimit(r *http.Request, items Items) (Items, HeaderData) {

	pageStr, pageGiven := r.URL.Query()["page"]
	pageSizeStr, pageSizeGiven := r.URL.Query()["pagesize"]
	limitStr, limitGiven := r.URL.Query()["limit"]
	headerData := make(HeaderData)

	limit := len(items)
	if limitGiven {
		limit = intMoreDefault(limitStr[0], 1)
		headerData["Limit"] = strconv.Itoa(limit)
	}

	pageSize := 10
	if pageSizeGiven {
		pageSize = intMoreDefault(pageSizeStr[0], 1)

	}

	page := 1
	if pageGiven {
		page = intMoreDefault(pageStr[0], 1)
		headerData["Page"] = strconv.Itoa(page)
		headerData["Page-Size"] = strconv.Itoa(pageSize)
		headerData["Total-Pages"] = strconv.Itoa((len(items) / pageSize) + 1)
	}

	headerData["Total-Items"] = strconv.Itoa(len(items))

	if !limitGiven && !pageGiven {
		return items, headerData
	}
	sortingL, sortingGiven := r.URL.Query()["sortby"]
	if sortingGiven {
		items, _ = sortBy(items, sortingL)
	}
	//TODO fix below
	start := (page - 1) * pageSize
	end := start + pageSize
	if end > len(items) {
		end = len(items)
	}
	if len(items) <= limit {
		return items[start:end], headerData
	}
	items = items[start:end]
	if len(items) < limit {
		return items, headerData
	}
	return items[:limit], headerData
}

// API

func contextListRest(JWTConig jwtConfig, itemChan ItemsChannel, operations GroupedOperations) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		filterMap, excludeMap := parseURLParameters(r)
		fmt.Println("request", r.URL, "items", len(ITEMS))
		items := filtered(ITEMS, filterMap, excludeMap, operations)

		items, headerData := sortLimit(r, items)

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
	response["total-items"] = []string{totalItems}
	response["settings"] = []string{
		fmt.Sprintf("host: %s", SETTINGS.Get("http_db_host")),
		fmt.Sprintf("JWT: %s", SETTINGS.Get("JWTENABLED")),
	}
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}
