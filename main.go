package main

import (
	"fmt"
	"log"
	"net/http" //	"runtime/debug" "github.com/pkg/profile")
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

type jwtConfig struct {
	Enabled      bool
	SharedSecret string
}

func init() {
	SETTINGS.Set("http_db_host", "0.0.0.0:8000", "host with port")
	SETTINGS.Set("SHAREDSECRET", "", "jwt shared secret")
	SETTINGS.Set("JWTENABLED", "yes", "JWT enabled")
	SETTINGS.Parse()

	ITEMS = make(Items, 0, 100*1000)

	fmt.Println(Operations)
}

func main() {
	Operations = GroupedOperations{Funcs: RegisterFuncMap, GroupBy: RegisterGroupBy}
	itemChan := make(ItemsChannel, 1000)

	go ItemChanWorker(itemChan)
	JWTConfig := jwtConfig{
		Enabled:      SETTINGS.Get("JWTENABLED") == "yes",
		SharedSecret: SETTINGS.Get("SHAREDSECRET"),
	}

	listRest := contextListRest(JWTConfig, itemChan, Operations)
	addRest := contextAddRest(JWTConfig, itemChan, Operations)

	ipPort := SETTINGS.Get("http_db_host")
	http.HandleFunc("/", listRest)
	http.HandleFunc("/help/", helpRest)
	http.HandleFunc("/add/", addRest)
	http.HandleFunc("/rm/", rmRest)
	fmt.Println("starting server", ipPort, " with:", len(ITEMS), "items", "jwt enabled: ", JWTConfig.Enabled)
	log.Fatal(http.ListenAndServe(ipPort, nil))
}
