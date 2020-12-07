/*

  Determine S2 cells involved in geometries. Provide a fast way to lookup
  data from based on a geojson query.

  inspired by
        "github.com/akhenakh/oureadb/index/geodata"
	"github.com/akhenakh/oureadb/store"

	s2 cell index code.

  With S2 CillIDs we can find which items are contained in given
  filter geometry (S2 cell union).

*/

package main

import (
	"fmt"
	"github.com/go-spatial/geom"
	"github.com/go-spatial/geom/encoding/wkt"
	"github.com/golang/geo/s2"
	"sort"
	"strings"
	"sync"
)

var minLevel int
var maxLevel int
var maxCells int

var s2Lock = sync.RWMutex{}

type cellIndexNode struct {
	Cell  s2.Cell
	Label int
}

type s2CellIndex []cellIndexNode
type s2CellMap map[int]s2.CellID

var S2CELLS s2CellIndex
var S2CELLMAP s2CellMap

func init() {
	minLevel = 7
	maxLevel = 20
	maxCells = 50

	S2CELLS = make(s2CellIndex, 100000)
	S2CELLMAP = s2CellMap{}
}

func BuildGeoIndex() {
	for i, v := range ITEMS {
		v.GeoIndex(i)
	}
}

//GeoIndex for each items determine S2Cell and store it.
func (i Item) GeoIndex(label int) error {
	if i.GetGeometry() == "" {
		return fmt.Errorf("missing wkt geometry")
	}
	sreader := strings.NewReader(i.GetGeometry())
	g, err := wkt.Decode(sreader)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(i.GetGeometry())
		fmt.Println(i.Ekey)
		return fmt.Errorf("wkt error encountered with %s", i.Point)
	}

	p, err := geom.GetCoordinates(g)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(i.Ekey)
		fmt.Println(i.GetGeometry())
		fmt.Printf("geom error encountered with %s", i.Point)
		return fmt.Errorf("geom error")
	}
	s2Lock.Lock()
	defer s2Lock.Unlock()

	y := p[0][0]
	x := p[0][1]
	center := s2.PointFromLatLng(s2.LatLngFromDegrees(x, y))
	cell := s2.CellFromPoint(center)

	cnode := cellIndexNode{Cell: cell, Label: label}
	S2CELLS = append(S2CELLS, cnode)
	S2CELLMAP[label] = cell.ID()

	return nil

}

//CalculateCover calculate S2 covering from given user polygon.
func CalculateCover(geom string) {

}

// Simple search algo
func SearchOverlapItems(items *labeledItems, cu s2.CellUnion) labeledItems {

	cellUnion := make([]s2.Cell, 0)

	// Create S2cells from cell id.
	for _, c := range cu {
		cell := s2.CellFromCellID(c)
		cellUnion = append(cellUnion, cell)
	}

	newItems := labeledItems{}

	for k, i := range *items {
		if cu.ContainsCellID(S2CELLMAP[k]) {
			newItems[k] = i
		}
	}

	return newItems
}

// Given only a cell Union return items
/*
func SearchRelevantItems(cu s2.CellUnion) Items {

	cellUnion := make([]s2.Cell, 0)

	//Create S2cells from cell id.
	for _, c := range cu {
		cell := s2.CellFromCellID(c)
		cellUnion = append(cellUnion, cell)
	}

	newItems := make(Items, 0)

	min = S2CellS.Seek(cu.RectBound().

	for idx, i := range S2CellS {
		if SearchOverlap(idx, cellUnion) {
			newItems = append(newItems, i)
		}
	}

	return newItems

}
*/

// SearchOverlap check if any cell of celluntion contains item points
func SearchOverlap(i int, cu []s2.Cell) bool {

	s2Lock.RLock()
	defer s2Lock.RUnlock()

	for _, c := range cu {
		if c.ContainsCell(S2CELLS[i].Cell) {
			return true
		}
	}
	return false
}

// GeoIdsAtCells returns all GeoData keys contained in the cells, without duplicates
/*
func (idx *s2.S2FlatIdx) GeoIdsAtCells(cells []s2.CellID) ([]GeoID, error) {
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
*/

//func (ca *s2Cell)

// Seek position in index which is close to target
func (ci s2CellIndex) Seek(target s2.CellID) int {
	pos := sort.Search(len(ci), func(i int) bool {
		return ci[i].Cell.ID() > target
	}) - 1

	// Ensure we don't go beyond the beginning.
	if pos < 0 {
		pos = 0
	}
	return pos
}
