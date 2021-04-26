/*
	Transforming ItemsIn -> Items -> ItemsOut
	Where Items has column values ar integers to save memmory
	maps are needed to restore integers back to the actual values.
	those are generated and stored here.
*/
package main

import ()

type ModelMaps struct {
	WoningType              MappedColumn
	LabelscoreVoorlopig     MappedColumn
	LabelscoreDefinitief    MappedColumn
	Gemeentecode            MappedColumn
	Gemeentenaam            MappedColumn
	Buurtcode               MappedColumn
	Buurtnaam               MappedColumn
	Wijkcode                MappedColumn
	Wijknaam                MappedColumn
	Provinciecode           MappedColumn
	Provincienaam           MappedColumn
	PandGasEanAansluitingen MappedColumn
	P6GasAansluitingen2020  MappedColumn
	P6Gasm32020             MappedColumn
	P6Kwh2020               MappedColumn
	PandBouwjaar            MappedColumn
	PandGasAansluitingen    MappedColumn
	Gebruiksdoelen          MappedColumn
}

var modelmaps2 map[string]MappedColumn

// Column maps.
// Store for each non distinct/repeated column

var BitArrays map[string]fieldBitarrayMap

var WoningType MappedColumn
var LabelscoreVoorlopig MappedColumn
var Gemeentecode MappedColumn
var LabelscoreDefinitief MappedColumn
var Gemeentenaam MappedColumn
var Buurtcode MappedColumn
var Buurtnaam MappedColumn
var Provinciecode MappedColumn
var Wijkcode MappedColumn
var Wijknaam MappedColumn
var Provincienaam MappedColumn
var PandGasEanAansluitingen MappedColumn
var P6GasAansluitingen2020 MappedColumn
var P6Gasm32020 MappedColumn
var P6Kwh2020 MappedColumn
var PandBouwjaar MappedColumn
var PandGasAansluitingen MappedColumn
var Gebruiksdoelen MappedColumn

func clearBitArrays() {
	BitArrays = make(map[string]fieldBitarrayMap)
}

func init() {
	clearBitArrays()
}

func setUpRepeatedColumns() {
	WoningType = NewReapeatedColumn("woning_type")
	LabelscoreVoorlopig = NewReapeatedColumn("labelscore_voorlopig")
	LabelscoreDefinitief = NewReapeatedColumn("labelscore_definitief")
	Gemeentecode = NewReapeatedColumn("gemeentecode")
	Gemeentenaam = NewReapeatedColumn("gemeentenaam")
	Buurtcode = NewReapeatedColumn("buurtcode")
	Buurtnaam = NewReapeatedColumn("buurtnaam")
	Wijkcode = NewReapeatedColumn("wijkcode")
	Wijknaam = NewReapeatedColumn("wijknaam")
	Provinciecode = NewReapeatedColumn("provinciecode")
	Provincienaam = NewReapeatedColumn("provincienaam")
	PandGasEanAansluitingen = NewReapeatedColumn("pand_gas_ean_aansluitingen")
	P6GasAansluitingen2020 = NewReapeatedColumn("p6_gas_aansluitingen_2020")
	P6Gasm32020 = NewReapeatedColumn("p6_gasm3_2020")
	P6Kwh2020 = NewReapeatedColumn("p6_kwh_2020")
	PandBouwjaar = NewReapeatedColumn("pand_bouwjaar")
	PandGasAansluitingen = NewReapeatedColumn("pand_gas_aansluitingen")
	Gebruiksdoelen = NewReapeatedColumn("gebruiksdoelen")
}

func CreateMapstore() ModelMaps {
	return ModelMaps{
		WoningType,
		LabelscoreVoorlopig,
		LabelscoreDefinitief,
		Gemeentecode,
		Gemeentenaam,
		Buurtcode,
		Buurtnaam,
		Wijkcode,
		Wijknaam,
		Provinciecode,
		Provincienaam,
		PandGasEanAansluitingen,
		P6GasAansluitingen2020,
		P6Gasm32020,
		P6Kwh2020,
		PandBouwjaar,
		PandGasAansluitingen,
		Gebruiksdoelen,
	}
}

func LoadMapstore(m ModelMaps) {
	WoningType = m.WoningType
	LabelscoreVoorlopig = m.LabelscoreVoorlopig
	LabelscoreDefinitief = m.LabelscoreDefinitief
	Gemeentecode = m.Gemeentecode
	Gemeentenaam = m.Gemeentenaam
	Buurtcode = m.Buurtcode
	Buurtnaam = m.Buurtnaam
	Wijkcode = m.Wijkcode
	Wijknaam = m.Wijknaam
	Provinciecode = m.Provinciecode
	Provincienaam = m.Provincienaam
	PandGasEanAansluitingen = m.PandGasEanAansluitingen
	P6GasAansluitingen2020 = m.P6GasAansluitingen2020
	P6Gasm32020 = m.P6Gasm32020
	P6Kwh2020 = m.P6Kwh2020
	PandBouwjaar = m.PandBouwjaar
	PandGasAansluitingen = m.PandGasAansluitingen
	Gebruiksdoelen = m.Gebruiksdoelen
}
