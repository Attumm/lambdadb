/*
	convert geojson to s2 cover

	derived from code found @
	github.com/akhenakh/oureadb
*/

package main

import (
	"github.com/go-spatial/geom"
	//"github.com/go-spatial/geom/encoding/geojson"
	"fmt"
	"github.com/golang/geo/s2"
	"github.com/pkg/errors"
)

//geoDataCoverCellUnion given geometry create an s2 cover for it
func geoDataCoverCellUnion(g geom.Geometry, coverer *s2.RegionCoverer, interior bool) (s2.CellUnion, error) {

	if geom.IsEmpty(g) {
		fmt.Println("empty?")
		return nil, errors.New("invalid geometry")
	}
	var cu s2.CellUnion

	switch gg := g.(type) {
	case geom.Point:
		points, _ := geom.GetCoordinates(gg)
		c := s2.CellIDFromLatLng(
			s2.LatLngFromDegrees(points[0][1], points[0][0]),
		)
		cu = append(cu, c.Parent(coverer.MinLevel))

	case geom.Polygon:
		points, _ := geom.GetCoordinates(gg)
		cup, err := coverPolygon(points, coverer, interior)
		if err != nil {
			return nil, errors.Wrap(err, "can't cover polygon")
		}
		cu = append(cu, cup...)

	case geom.MultiPolygon:
		for _, p := range gg.Polygons() {
			points, _ := geom.GetCoordinates(p)
			cup, err := coverPolygon(points, coverer, interior)
			if err != nil {
				return nil, errors.Wrap(err, "can't cover multipolygon")
			}

			cu = append(cu, cup...)
		}

	case geom.LineString:
		points, _ := geom.GetCoordinates(gg)
		if len(points)%2 != 0 {
			return nil, errors.New("invalid coordinates count for line")
		}

		pl := make(s2.Polyline, len(points))
		for i := 0; i < len(points); i += 1 {
			ll := s2.LatLngFromDegrees(points[i][1], points[i][0])
			pl[i] = s2.PointFromLatLng(ll)
		}

		var cupl s2.CellUnion
		if interior {
			cupl = coverer.InteriorCellUnion(&pl)
		} else {
			cupl = coverer.CellUnion(&pl)
		}
		cu = append(cu, cupl...)

	default:
		fmt.Println(gg)
		return nil, errors.New("unsupported geojson data type")
	}

	return cu, nil
}

func CoverDefault(g geom.Geometry) s2.CellUnion {

	coverer := &s2.RegionCoverer{MinLevel: minLevel, MaxLevel: maxLevel, MaxCells: maxCells}
	cu, err := Cover(g, coverer)

	// no cover for this geo object this is probably an error
	if len(cu) == 0 || err != nil {
		fmt.Println("geo object can't be indexed, empty cover")
		fmt.Println(err)
	}
	return cu
}

// Cover generates an s2 cover for GeoData gd
func Cover(g geom.Geometry, coverer *s2.RegionCoverer) (s2.CellUnion, error) {
	return geoDataCoverCellUnion(g, coverer, false)
}

// returns an s2 cover from a list of lng, lat forming a closed polygon
func coverPolygon(p []geom.Point, coverer *s2.RegionCoverer, interior bool) (s2.CellUnion, error) {
	if len(p) < 3 {
		return nil, errors.New("invalid polygons not enough coordinates for a closed polygon")
	}
	if len(p)%2 != 0 {
		return nil, errors.New("invalid polygons odd coordinates number")
	}

	l := LoopFromCoordinatesAndCCW(p, true)
	if l.IsEmpty() || l.IsFull() {
		return nil, errors.New("invalid polygons")
	}

	// super hacky try reverse if ContainsOrigin
	if l.ContainsOrigin() {
		// reversing the slice
		for i := len(p)/2 - 1; i >= 0; i-- {
			opp := len(p) - 1 - i
			p[i], p[opp] = p[opp], p[i]
		}
	}

	if interior {
		return coverer.InteriorCovering(l), nil
	}
	return coverer.Covering(l), nil
}

// LoopFromCoordinatesAndCCW creates a LoopFence from a list of lng lat
// if checkCCW is true also try to fix CCW
func LoopFromCoordinatesAndCCW(p []geom.Point, checkCCW bool) *s2.Loop {
	if len(p)%2 != 0 || len(p) < 3 {
		return nil
	}
	points := make([]s2.Point, len(p))

	for i := 0; i < len(p); i += 1 {
		points[i] = s2.PointFromLatLng(s2.LatLngFromDegrees(p[i][1], p[i][0]))
	}

	if checkCCW && s2.RobustSign(points[0], points[1], points[2]) != s2.CounterClockwise {
		// reversing the slice
		for i := len(points)/2 - 1; i >= 0; i-- {
			opp := len(points) - 1 - i
			points[i], points[opp] = points[opp], points[i]
		}
	}

	if points[0] == points[len(points)-1] {
		// remove last item if same as 1st
		points = append(points[:len(points)-1], points[len(points)-1+1:]...)
	}

	loop := s2.LoopFromPoints(points)
	return loop
}
