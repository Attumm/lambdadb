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
	ID    s2.CellID
	Label int
}

type s2CellIndex []cellIndexNode
type s2CellMap map[int]s2.CellID

// Implement Sort interface for s2CellIndex
func (c s2CellIndex) Len() int           { return len(c) }
func (c s2CellIndex) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c s2CellIndex) Less(i, j int) bool { return c[i].ID < c[j].ID }

var S2CELLS s2CellIndex
var S2CELLMAP s2CellMap

func init() {
	minLevel = 2
	maxLevel = 21
	maxCells = 450

	S2CELLS = make(s2CellIndex, 0)
	S2CELLMAP = s2CellMap{}
}

func BuildGeoIndex() {
	for i, v := range ITEMS {
		v.GeoIndex(i)
	}
}

func (c cellIndexNode) IsEmpty() bool {
	return c.ID == 0
}

// GeoIndex for each items determine S2Cell and store it.
func (i Item) GeoIndex(label int) error {

	lock.Lock()
	defer lock.Unlock()

	if i.GetGeometry() == "" {
		return fmt.Errorf("missing wkt geometry")
	}
	sreader := strings.NewReader(i.GetGeometry())
	g, err := wkt.Decode(sreader)

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(i.GetGeometry())
		return fmt.Errorf("wkt error encountered with %s", i.GetGeometry())
	}

	p, err := geom.GetCoordinates(g)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(i.GetGeometry())
		fmt.Printf("geom error encountered with %s", i.GetGeometry())
		return fmt.Errorf("geom error")
	}

	s2Lock.Lock()
	defer s2Lock.Unlock()

	y := p[0][0]
	x := p[0][1]
	ll := s2.LatLngFromDegrees(x, y)

	if !ll.IsValid() {
		fmt.Println(i.GetGeometry())
		fmt.Printf("ll geom error encountered with %f %f", x, y)
		return fmt.Errorf("geom error")
	}

	center := s2.PointFromLatLng(ll)
	cell := s2.CellFromPoint(center)

	cnode := cellIndexNode{ID: cell.ID(), Label: label}
	S2CELLS = append(S2CELLS, cnode)
	S2CELLMAP[label] = cell.ID()

	// Update index while loading data so queries already work
	if label%100000 == 0 {
		S2CELLS.Sort()
	}

	return nil

}

// Simple search algo
func SearchOverlapItems(items *Items, cu s2.CellUnion) Items {

	s2Lock.RLock()
	defer s2Lock.RUnlock()

	newItems := Items{}

	for k, i := range *items {
		if cu.ContainsCellID(S2CELLMAP[k]) {
			newItems[k] = i
		}
	}
	return newItems
}

// Given only a cell Union return Items
func SearchGeoItems(cu s2.CellUnion) Items {

	newItems := Items{}

	cu.Normalize()

	//for i, c := range cu {
	//	fmt.Printf("%d %s \n", i, c)
	//}

	min := S2CELLS.Seek(cu[0].ChildBegin())
	max := S2CELLS.Seek(cu[len(cu)-1].ChildEnd())

	// ITEMS read lock
	lock.RLock()
	defer lock.RUnlock()

	for _, i := range S2CELLS[min : max+1] {
		if cu.ContainsCellID(i.ID) {
			newItems = append(newItems, ITEMS[i.Label])
		}
	}

	return newItems

}

// Seek position in index which is close to target
func (ci s2CellIndex) Seek(target s2.CellID) int {

	pos := sort.Search(len(ci), func(i int) bool {
		return ci[i].ID > target
	}) - 1

	// Ensure we don't go beyond the beginning.
	if pos < 0 {
		pos = 0
	}
	return pos
}

// Sort CellIndex so Binary search can work.
func (ci s2CellIndex) Sort() {
	sort.Sort(ci)
}
