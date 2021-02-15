package main

import (
	"fmt"
	"log"
	"net/http" //	"runtime/debug" "github.com/pkg/profile")
)

type filterFuncc func(*Item, string) bool
type registerFuncType map[string]filterFuncc
type registerGroupByFunc map[string]func(*Item) string
type registerGettersMap map[string]func(*Item) string
type registerReduce map[string]func(Items) map[string]string
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

}

func loadcsv(itemChan ItemsChannel) {
	log.Print("loading given csv")
	err := importCSV(SETTINGS.Get("csv"), itemChan,
		true, true,
		SETTINGS.Get("delimiter"),
		SETTINGS.Get("null-delimiter"))
	if err != nil {
		log.Fatal(err)
	}
	makeIndex()
}

func main() {
	SETTINGS.Set("http_db_host", "0.0.0.0:8000", "host with port")
	SETTINGS.Set("SHAREDSECRET", "", "jwt shared secret")
	SETTINGS.Set("JWTENABLED", "yes", "JWT enabled")

	SETTINGS.Set("csv", "", "load a gzipped csv file on starup")
	SETTINGS.Set("null-delimiter", "\\N", "null delimiter")
	SETTINGS.Set("delimiter", ",", "delimiter")

	SETTINGS.Set("readonly", "yes", "only allow read only funcions")

	SETTINGS.Set("indexed", "no", "is the data made for the suffix tree index?")
	SETTINGS.Parse()

	ITEMS = make(Items, 0, 100*1000)

	Operations = GroupedOperations{Funcs: RegisterFuncMap, GroupBy: RegisterGroupBy, Getters: RegisterGetters, Reduce: RegisterReduce}
	itemChan := make(ItemsChannel, 1000)

	go ItemChanWorker(itemChan)

	if SETTINGS.Get("csv") != "" {
		go loadcsv(itemChan)
	}

	JWTConfig := jwtConfig{
		Enabled:      SETTINGS.Get("JWTENABLED") == "yes",
		SharedSecret: SETTINGS.Get("SHAREDSECRET"),
	}

	listRest := contextListRest(JWTConfig, itemChan, Operations)
	addRest := contextAddRest(JWTConfig, itemChan, Operations)

	searchRest := contextSearchRest(JWTConfig, itemChan, Operations)
	typeAheadRest := contextTypeAheadRest(JWTConfig, itemChan, Operations)

	ipPort := SETTINGS.Get("http_db_host")

	mux := http.NewServeMux()

	mux.HandleFunc("/search/", searchRest)
	mux.HandleFunc("/typeahead/", typeAheadRest)
	mux.HandleFunc("/list/", listRest)
	mux.HandleFunc("/help/", helpRest)

	if SETTINGS.Get("readonly") != "yes" {
		mux.HandleFunc("/add/", addRest)
		mux.HandleFunc("/rm/", rmRest)
		mux.HandleFunc("/save/", saveRest)
		mux.HandleFunc("/load/", loadRest)

		mux.Handle("/", http.FileServer(http.Dir("./www")))
		mux.Handle("/dsm-search", http.FileServer(http.Dir("./www")))
	}

	msg := fmt.Sprint("starting server\nhost: ", ipPort, " with:", len(ITEMS), "items ", "jwt enabled: ", JWTConfig.Enabled)
	fmt.Printf(InfoColorN, msg)

	log.Fatal(http.ListenAndServe(ipPort, CORS(mux)))
}
