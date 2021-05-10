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
	"encoding/json"
	"fmt"
	// "io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

var handler http.Handler

/* load some data 19 records*/
func TestMain(m *testing.M) {

	defaultSettings()

	SETTINGS.Set(
		"csv", "./testdata/dataselectie_vbo_energie_20210505.head.csv",
		"test dataset")

	SETTINGS.Set("channelwait", "0.01s", "timeout for channel loading")

	loadcsv(itemChan)
	close(itemChan)
	ItemChanWorker(itemChan)

	handler = setupHandler()

	// Run the test
	m.Run()
}

func TestCsvLoading(t *testing.T) {

	size := len(ITEMS)

	if size != 10 {
		t.Errorf("expected 10 ITEMS got %d", size)
	}
}

func TestBasicHandlers(t *testing.T) {

	if len(ITEMS) < 10 {
		t.Error("no items")
	}

	type testCase struct {
		url      string
		expected string
	}

	tests := []testCase{
		testCase{"/list/?search=1", "10"},
		testCase{"/typeahead/huisnummer/?search=1", "3"},
		testCase{"/typeahead/pid/?search=1", "2"},
		testCase{"/help/", ""},
	}

	for i := range tests {
		req := httptest.NewRequest("GET", tests[i].url, nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		resp := w.Result()
		if resp.StatusCode != 200 {
			t.Errorf("request to %s failed", tests[i].url)
			t.Error(resp)
		}

		if tests[i].expected == "" {
			continue
		}

		if resp.Header.Get("Total-Items") != tests[i].expected {
			t.Errorf("total hits mismatch from %s %s != %s",
				tests[i].url,
				tests[i].expected,
				resp.Header.Get("Total-Items"),
			)
			t.Error(resp)
		}
	}
}

// Test geojson queries combined with groupby and reduce.
func TestGeoQuery(t *testing.T) {

	BuildGeoIndex()

	if len(ITEMS) < 10 {
		t.Error("no items")
	}

	if len(S2CELLS) == 0 {
		t.Error("geo indexing failed")
	}

	if len(S2CELLMAP) == 0 {
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
		    [4.905321, 52.377706],
		    [4.90527, 52.377706],
		    [4.90527, 52.377869],
		    [4.905321, 52.377869],
		    [4.905321, 52.377706]
		]
	]
}
	`)
	data.Set("geojson", geojson)

	params := strings.NewReader(data.Encode())

	req := httptest.NewRequest("POST", "/list/", params)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != 200 {
		t.Errorf("geo request to %s failed statuscode", req.URL)
		t.Error(resp)
	}

	headerQuery := resp.Header.Get("Query")
	query := Query{}
	json.Unmarshal([]byte(headerQuery), &query)

	if query.GeometryGiven != true {
		t.Errorf("geo request to %s failed ", req.URL)
		t.Error(resp.Header.Get("Query"))
		// t.Error(resp.Header.Get("GeometryGiven"))
		t.Error(resp.Body)
	}

	if resp.Header.Get("Total-Items") != "7" {
		t.Error("geo request count is not 7")
	}

	// parse json GroupBy response
	defer resp.Body.Close()
	j := GroupByResult{}
	err := json.NewDecoder(resp.Body).Decode(&j)

	if err != nil {
		t.Error(err)
	}

	if j["1011AB"]["count"] != "7" {
		t.Error("geo request json response count is not 7")
	}
}
