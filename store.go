package main

import (
	"fmt"
	"log"
	"time"
)

//Items
type Items []*Item
type ItemsIn []*ItemIn
type ItemsOut []*ItemOut

type ItemsGroupedBy map[string]Items
type ItemsChannel chan ItemsIn

var ITEMS Items

var itemChan ItemsChannel

func init() {
	ITEMS = Items{}
}

func ItemChanWorker(itemChan ItemsChannel) {
	label := 0

	for items := range itemChan {
		for _, itm := range items {
			if itm != nil {
				smallItem := itm.Shrink(label)
				smallItem.StoreBitArrayColumns()
				ITEMS = append(ITEMS, &smallItem)
				// ITEMS[label] = &smallItem
				if ITEMS[label] != &smallItem {
					log.Fatal("storing item index off")
				}
				smallItem.GeoIndex(label)
				label++
			}
		}
	}
}

func (items Items) FillIndexes() {

	start := time.Now()

	clearGeoIndex()
	initBitarrays()

	for i := range items {
		ITEMS[i].StoreBitArrayColumns()
		ITEMS[i].GeoIndex(ITEMS[i].Label)
	}

	diff := time.Since(start)
	msg := fmt.Sprint("Index set time: ", diff)
	fmt.Printf(WarningColorN, msg)
}
