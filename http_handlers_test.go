/*

# test some basic request handling.

 - typeahead: http://%s/list/?typeahead=ams&limit=10", host),
 - search: http://%s/list/?search=ams&page=1&pagesize=1", host),
 - search with limit: http://%s/list/?search=10&page=1&pagesize=10&limit=5", host),
 - sorting: http://%s/list/?search=100&page=10&pagesize=100&sortby=-country", host),
 - filtering: http://%s/list/?search=10&ontains=144&contains-case=10&page=1&pagesize=1", host),
 - groupby: http://%s/list/?search=10&contains-case=10&groupby=country", host),
 - aggregation: http://%s/list/?search=10&contains-case=10&groupby=country&reduce=count", host),
 - chain the same filters: http://%s/list/?search=10&contains-case=127&contains-case=0&contains-case=1", host),
 - typeahead use the name of the column in this case IP: http://%s/typeahead/ip/?starts-with=127&limit=15", host),


*/
package main

import (
	"fmt"
	// "io"
	// "net/http"
	"net/http/httptest"
	"testing"
)

/* load some data 19 records*/
func TestMain(m *testing.M) {

	SETTINGS.Set(
		"csv", "./testdata/dataselectie_vbo_energie_20210217.head.csv.gz",
		"test dataset")

	SETTINGS.Set("channelwait", "0.001s", "timeout for channel loading")
	itemChan := make(ItemsChannel, 1)
	loadcsv(itemChan)
	close(itemChan)
	ItemChanWorker(itemChan)
	// Run the test
	m.Run()
}

func TestCsvLoading(t *testing.T) {

	fmt.Println(len(ITEMS))
	size := len(ITEMS)
	if size != 19 {
		t.Errorf("expected 19 ITEMS got %d", size)
	}
}

func TestBasicHandlers(t *testing.T) {

	handler := setupHandler()

	urls := []string{
		"/list/",
		"/typeahead/pid/?search=1",
		"/help/",
	}

	for i := range urls {
		req := httptest.NewRequest("GET", urls[i], nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		resp := w.Result()
		if resp.StatusCode != 200 {
			t.Errorf("request to %s failed", urls[i])
			t.Error(resp)
		}
	}
}
