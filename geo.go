package main

import (
	"github.com/go-spatial/geom/encoding/wkt"
	"github.com/golang/geo/s2"
	"strings"
)

var sidx = s2.NewShapeIndex()

func buildGeoIndex() {

	for _, v := range ITEMS {
		addItem(v)
	}

}

func addItem(i *Item) {
	sreader := strings.NewReader(i.geom)
	geometry, err := wkt.Decode(sreader)
}
