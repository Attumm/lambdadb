package main

import (
	"fmt"
	"log"
	"net/http"
	//	"runtime/debug"
	//	"github.com/pkg/profile"
)

type filterFuncc func(*Item, string) bool
type registerFuncType map[string]filterFuncc
type registerGroupByFunc map[string]func(*Item) string
type filterType map[string][]string
type formatRespFunc func(w http.ResponseWriter, r *http.Request, items Items)
type registerFormatMap map[string]formatRespFunc

//Items as Example
type Items []*Item

type ItemsGroupedBy map[string]Items
type ItemsChannel chan Items

var ITEMS Items

func init() {
	SETTINGS.Set("http_db_host", "0.0.0.0:8000", "host with port")
	SETTINGS.Parse()

	ITEMS = make(Items, 0, 100*1000)

	registerFormat = make(registerFormatMap)
	registerFormat["json"] = formatResponseJSON
	registerFormat["csv"] = formatResponseCSV

	fmt.Println(Operations)
}

func main() {
	Operations = GroupedOperations{Funcs: RegisterFuncMap, GroupBy: RegisterGroupBy}
	itemChan := make(ItemsChannel, 1000)

	go ItemChanWorker(itemChan)

	listRest := contextListRest(itemChan, Operations)
	addRest := contextAddRest(itemChan, Operations)

	ipPort := SETTINGS.Get("http_db_host")
	http.HandleFunc("/", listRest)
	http.HandleFunc("/help/", helpRest)
	http.HandleFunc("/add/", addRest)
	http.HandleFunc("/rm/", rmRest)
	fmt.Println("starting server", ipPort, " with:", len(ITEMS), "items")
	log.Fatal(http.ListenAndServe(ipPort, nil))
}
