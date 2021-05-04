package main

import (
	"testing"
)

func TestBytesSaving(t *testing.T) {

	size := len(ITEMS)

	if size != 10 {
		t.Errorf("expected 10 ITEMS got %d", size)
	}

}

func TestBytes(t *testing.T) {

	saveAsBytes("testdata/testbytes")

	RegisteredColumns = make(ColumnRegister)
	ITEMS = Items{}

	clearBitArrays()
	clearGeoIndex()

	loadAsBytes("testdata/testbytes")

	if len(ITEMS) != 10 {
		t.Error("bytes save / load failed")
	}

	saveAsBytes("testdata/testbytesz")
	ITEMS = Items{}
	loadAsBytes("testdata/testbytesz")
	if len(ITEMS) != 10 {
		t.Error("bytes compressed save / load failed")
	}

	if len(BitArrays) == 0 {
		t.Error("bitarrays are not restored")
	}

	if len(S2CELLS) == 0 {
		t.Error("geoindex is not restored")
	}

	if len(RegisteredColumns) == 0 {
		t.Error("colom register is not restored")
	}

}

func TestJson(t *testing.T) {

	saveAsJsonZipped("testdata/test.json")
	ITEMS = Items{} // Clear ITEMS
	loadAsJsonZipped("testdata/test.json")
	if len(ITEMS) != 10 {
		t.Error("bytes compressed save / load failed")
	}
}
