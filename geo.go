/*

  Determine S2 cells involved in geometries.

  inspired by
        "github.com/akhenakh/oureadb/index/geodata"
	"github.com/akhenakh/oureadb/store"

*/
package main

import (
	"fmt"
	"github.com/go-spatial/geom"
	"github.com/go-spatial/geom/encoding/wkt"
	"github.com/golang/geo/s2"
	"strings"
)

var minLevel int
var maxLevel int
var maxCells int

func init() {
	minLevel = 19
	maxLevel = 19
	maxCells = 1
}

var sidx = s2.NewShapeIndex()

func BuildGeoIndex() {

	for i, v := range ITEMS {
		v.GeoIndex(i)
	}

}

func (i Item) GeoIndex(idx int) error {
	sreader := strings.NewReader(i.Point)
	g, err := wkt.Decode(sreader)
	if err != nil {
		fmt.Printf("error encountered with %s", i.Point)
	}
	p, err := geom.GetCoordinates(g)

	if err != nil {
		fmt.Printf("error encountered with %s", i.Point)
	}

	fmt.Println(p)

	x := p[0][0]
	y := p[0][1]
	center := s2.PointFromLatLng(s2.LatLngFromDegrees(x, y))
	cap := s2.CapFromCenterArea(center, s2RadialAreaMeters(2))

	coverer := &s2.RegionCoverer{MinLevel: minLevel, MaxLevel: maxLevel, MaxCells: maxCells}
	cu := coverer.Covering(cap)

	// no cover for this geo object this is probably an error
	if len(cu) == 0 {
		fmt.Printf("geo object can't be indexed, empty cover")
	}
	return nil

}

//CalculateCover calculate S2 covering from given user polygon.
func CalculateCover(geom string) {

}

// GeoIdsAtCells returns all GeoData keys contained in the cells, without duplicates
func (idx *S2FlatIdx) GeoIdsAtCells(cells []s2.CellID) ([]GeoID, error) {
	m := make(map[string]struct{})

	for _, c := range cells {
		ids, err := idx.GeoIdsAtCell(c)
		if err != nil {
			return nil, errors.Wrap(err, "fetching geo ids from cells failed")
		}
		for _, id := range ids {
			m[string(id)] = struct{}{}
		}
	}

	res := make([]GeoID, len(m))
	var i int
	for k := range m {
		res[i] = []byte(k)
		i++
	}

	return res, nil
}
