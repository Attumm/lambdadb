package main

import (
	"errors"
	"github.com/Workiva/go-datastructures/bitarray"
	"log"
	"strings"
)

type fieldIdxMap map[string]uint32
type IdxFieldMap map[uint32]string

type MappedColumn struct {
	Idx        fieldIdxMap
	Field      IdxFieldMap
	IdxTracker uint32
}

type ColumnRegister map[string]MappedColumn

var RepeatedColumns ColumnRegister

func init() {
	RepeatedColumns = make(ColumnRegister)
}

func NewReapeatedColumn(column string) MappedColumn {
	m := MappedColumn{
		make(fieldIdxMap),
		make(IdxFieldMap),
		0,
	}
	RepeatedColumns[column] = m
	return m
}

// Store field name as idx value and idx as field value
func (m *MappedColumn) Store(field string) {

	if _, ok := m.Idx[field]; !ok {
		m.Idx[field] = m.IdxTracker
		m.Field[m.IdxTracker] = field
		m.IdxTracker++
	}
}

// Store Array field (postgres Array).
func (m *MappedColumn) StoreArray(field string) []uint32 {

	fieldsArray := make([]uint32, 0)

	// parsing {a, b} array values
	// string should be at least 2 example "{}" == size 2
	if len(field) > 2 {
		fields, err := ParsePGArray(field)

		if err != nil {
			log.Fatal(err, "error parsing array ")
		}

		for _, gd := range fields {
			m.Store(gd)
		}

		for _, v := range fields {
			fieldsArray = append(fieldsArray, Gebruiksdoelen.GetIndex(v))
		}
	}
	return fieldsArray
}

func (m *MappedColumn) GetValue(idx uint32) string {
	return m.Field[idx]
}

func (m *MappedColumn) GetArrayValue(idxs []uint32) string {

	result := make([]string, 0)
	for _, v := range idxs {
		vs := m.GetValue(v)
		result = append(result, vs)
	}
	return strings.Join(result, ", ")
}

func (m *MappedColumn) GetIndex(s string) uint32 {
	return m.Idx[s]
}

// SetBitArray WIP
func SetBitArray(column string, i uint32, label int) {

	var ba bitarray.BitArray
	var ok bool

	// check if map of bitmaps is present for column
	var map_ba fieldBitarrayMap
	if _, ok = BitArrays[column]; !ok {
		map_ba := make(fieldBitarrayMap)
		BitArrays[column] = map_ba
	}

	map_ba = BitArrays[column]

	// check for existing bitarray for i value
	ba, ok = map_ba[i]
	if !ok {
		ba = bitarray.NewSparseBitArray()
		map_ba[i] = ba
	}
	// set bit for item label.
	ba.SetBit(uint64(label))
}

func GetBitArray(column, value string) (bitarray.BitArray, error) {

	var ok bool

	if _, ok = BitArrays[column]; !ok {
		return nil, errors.New("no bitarray filter found for column " + column)
	}

	bpi, ok := RepeatedColumns[column].Idx[value]

	if !ok {
		return nil, errors.New("no bitarray filter found for column value WoningType")
	}

	ba, ok := BitArrays[column][bpi]

	if !ok {
		return nil, errors.New("no bitarray filter found for column idx value WoningType")
	}

	return ba, nil
}
