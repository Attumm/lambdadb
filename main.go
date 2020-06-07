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

// Colors are fun, and can be used to note that this is joyfull and fun project.
const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"

	InfoColorN    = "\033[1;34m%s\033[0m\n"
	NoticeColorN  = "\033[1;36m%s\033[0m\n"
	WarningColorN = "\033[1;33m%s\033[0m\n"
	ErrorColorN   = "\033[1;31m%s\033[0m\n"
	DebugColorN   = "\033[0;36m%s\033[0m\n"
)

func init() {
	SETTINGS.Set("http_db_host", "0.0.0.0:8000", "host with port")
	SETTINGS.Set("SHAREDSECRET", "", "jwt shared secret")
	SETTINGS.Set("JWTENABLED", "yes", "JWT enabled")
	SETTINGS.Parse()

	ITEMS = make(Items, 0, 100*1000)
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
	http.HandleFunc("/save/", saveRest)
	http.HandleFunc("/load/", loadRest)
	msg := fmt.Sprint("starting server\nhost: ", ipPort, " with:", len(ITEMS), "items ", "jwt enabled: ", JWTConfig.Enabled)
	fmt.Printf(InfoColorN, msg)
	log.Fatal(http.ListenAndServe(ipPort, nil))
}
