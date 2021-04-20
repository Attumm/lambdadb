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
	"net/url"
	"strings"
	"testing"
)

/* load some data 19 records*/
func TestMain(m *testing.M) {

	defaultSettings()

	SETTINGS.Set(
		"csv", "./testdata/dataselectie_vbo_energie_20210217.head.csv.gz",
		"test dataset")

	SETTINGS.Set("channelwait", "0.01s", "timeout for channel loading")

	loadcsv(itemChan)
	close(itemChan)
	ItemChanWorker(itemChan)

	// Run the test
	m.Run()
}

func TestCsvLoading(t *testing.T) {

	fmt.Println(len(ITEMS))

	size := len(ITEMS)

	if size != 9 {
		t.Errorf("expected 9 ITEMS got %d", size)
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

func TestGeoQuery(t *testing.T) {

	BuildGeoIndex()

	if len(S2CELLS) == 0 {
		t.Error("geo indexing failed")
	}

	data := url.Values{}
	data.Set("groupby", "postcode")
	data.Set("reduce", "count")

	geojson := fmt.Sprint(`
{
	"type": "Polygon",
	"coordinates": [
		[
		    [4.902321, 52.428306],
		    [4.90127, 52.427024],
		    [4.905281, 52.426069],
		    [4.906782, 52.426226],
		    [4.906418, 52.427469],
		    [4.902321, 52.428306]
		]
	]
}
	`)
	data.Set("geojson", geojson)

	params := strings.NewReader(data.Encode())

	handler := setupHandler()
	req := httptest.NewRequest("POST", "/list/", params)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != 201 {
		t.Errorf("request to %s failed", req.URL)
		t.Error(resp)
	}
}
