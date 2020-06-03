package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

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

func formatResponseCSV(w http.ResponseWriter, r *http.Request, items Items) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=output.csv")
	wr := csv.NewWriter(w)
	if err := wr.Write(items[0].Columns()); err != nil {
		log.Fatal(err)
	}
	for _, item := range items { // make a loop for 100 rows just for testing purposes
		if err := wr.Write(item.Row()); err != nil {
			log.Fatal(err)
		}
	}
	wr.Flush() // writes the csv writer data to  the buffered data io writer(b(bytes.buffer))
}

func FormatAndSend(w http.ResponseWriter, r *http.Request, items Items) {
	respFormatSlice, respFormatFound := r.URL.Query()["format"]
	respFormat := ""
	if respFormatFound {
		respFormat = respFormatSlice[0]
	}

	w.Header().Set("Total-Items", strconv.Itoa(len(items)))
	w.WriteHeader(http.StatusOK)

	respFormatFunc, found := registerFormat[respFormat]
	if !found {
		respFormatFunc = registerFormat["json"]
	}
	respFormatFunc(w, r, items)
}

func formatResponseJSON(w http.ResponseWriter, r *http.Request, items Items) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// API

func contextListRest(JWTConig jwtConfig, itemChan ItemsChannel, operations GroupedOperations) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		filterMap, excludeMap := parseURLParameters(r)
		fmt.Println("request", r.URL, "items", len(ITEMS))
		items := filtered(ITEMS, filterMap, excludeMap, operations)

		items = sortLimit(w, r, items)

		w.Header().Set("Content-Type", "application/json")
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
	response["filters"] = registeredFilters
	response["groupby"] = registeredGroupbys
	totalItems := strconv.Itoa(len(ITEMS))
	response["total-items"] = []string{"total size", totalItems}
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}
