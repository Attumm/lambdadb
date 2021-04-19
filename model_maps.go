/*

	When transforming ItemsIn to Items and back again to ItemsOut

	maps are needed to store lookup values.

	those are generated here.

*/

package main

import (
	"sync"
)

type ModelMaps struct {
	WoningTypeTracker uint16
	WoningTypeIdxMap  fieldIdxMap
	WoningType        fieldMapIdx

	LabelscoreVoorlopigTracker uint16
	LabelscoreVoorlopigIdxMap  fieldIdxMap
	LabelscoreVoorlopig        fieldMapIdx

	// LabelscoreVoorlopigItems fieldItemsMap

	LabelscoreDefinitiefTracker uint16
	LabelscoreDefinitiefIdxMap  fieldIdxMap
	LabelscoreDefinitief        fieldMapIdx

	// LabelscoreDefinitiefItems fieldItemsMap

	GemeentecodeTracker uint16
	GemeentecodeIdxMap  fieldIdxMap
	Gemeentecode        fieldMapIdx

	// GemeentecodeItems fieldItemsMap

	GemeentenaamTracker uint16
	GemeentenaamIdxMap  fieldIdxMap
	Gemeentenaam        fieldMapIdx

	BuurtcodeTracker uint16
	BuurtcodeIdxMap  fieldIdxMap
	Buurtcode        fieldMapIdx

	// BuurtcodeItems fieldItemsMap

	BuurtnaamTracker uint16
	BuurtnaamIdxMap  fieldIdxMap
	Buurtnaam        fieldMapIdx

	WijkcodeTracker uint16
	WijkcodeIdxMap  fieldIdxMap
	Wijkcode        fieldMapIdx

	// WijkcodeItems fieldItemsMap

	WijknaamTracker uint16
	WijknaamIdxMap  fieldIdxMap
	Wijknaam        fieldMapIdx

	ProvinciecodeTracker uint16
	ProvinciecodeIdxMap  fieldIdxMap
	Provinciecode        fieldMapIdx

	// ProvinciecodeItems fieldItemsMap

	ProvincienaamTracker uint16
	ProvincienaamIdxMap  fieldIdxMap
	Provincienaam        fieldMapIdx

	PandGasEanAansluitingenTracker uint16
	PandGasEanAansluitingenIdxMap  fieldIdxMap
	PandGasEanAansluitingen        fieldMapIdx

	P6GasAansluitingen2020Tracker uint16
	P6GasAansluitingen2020IdxMap  fieldIdxMap
	P6GasAansluitingen2020        fieldMapIdx

	P6Gasm32020Tracker uint16
	P6Gasm32020IdxMap  fieldIdxMap
	P6Gasm32020        fieldMapIdx

	P6Kwh2020Tracker uint16
	P6Kwh2020IdxMap  fieldIdxMap
	P6Kwh2020        fieldMapIdx

	PandBouwjaarTracker uint16
	PandBouwjaarIdxMap  fieldIdxMap
	PandBouwjaar        fieldMapIdx

	PandGasAansluitingenTracker uint16
	PandGasAansluitingenIdxMap  fieldIdxMap
	PandGasAansluitingen        fieldMapIdx

	GebruiksdoelenTracker uint16
	GebruiksdoelenIdxMap  fieldIdxMap
	Gebruiksdoelen        fieldMapIdx
}

// Column maps.
// Store for each non distinct/repeated column
// unit16 -> string map and
// string -> unit16 map
// track count of distinct values

var WoningTypeTracker uint16
var WoningTypeIdxMap fieldIdxMap
var WoningType fieldMapIdx

var WoningTypeItems fieldItemsMap

var LabelscoreVoorlopigTracker uint16
var LabelscoreVoorlopigIdxMap fieldIdxMap
var LabelscoreVoorlopig fieldMapIdx

var LabelscoreVoorlopigItems fieldItemsMap

var LabelscoreDefinitiefTracker uint16
var LabelscoreDefinitiefIdxMap fieldIdxMap
var LabelscoreDefinitief fieldMapIdx

var LabelscoreDefinitiefItems fieldItemsMap

var GemeentecodeTracker uint16
var GemeentecodeIdxMap fieldIdxMap
var Gemeentecode fieldMapIdx

var GemeentecodeItems fieldItemsMap

var GemeentenaamTracker uint16
var GemeentenaamIdxMap fieldIdxMap
var Gemeentenaam fieldMapIdx

var BuurtcodeTracker uint16
var BuurtcodeIdxMap fieldIdxMap
var Buurtcode fieldMapIdx

var BuurtcodeItems fieldItemsMap

var BuurtnaamTracker uint16
var BuurtnaamIdxMap fieldIdxMap
var Buurtnaam fieldMapIdx

var WijkcodeTracker uint16
var WijkcodeIdxMap fieldIdxMap
var Wijkcode fieldMapIdx

var WijkcodeItems fieldItemsMap

var WijknaamTracker uint16
var WijknaamIdxMap fieldIdxMap
var Wijknaam fieldMapIdx

var ProvinciecodeTracker uint16
var ProvinciecodeIdxMap fieldIdxMap
var Provinciecode fieldMapIdx

var ProvinciecodeItems fieldItemsMap

var ProvincienaamTracker uint16
var ProvincienaamIdxMap fieldIdxMap
var Provincienaam fieldMapIdx

var PandGasEanAansluitingenTracker uint16
var PandGasEanAansluitingenIdxMap fieldIdxMap
var PandGasEanAansluitingen fieldMapIdx

var P6GasAansluitingen2020Tracker uint16
var P6GasAansluitingen2020IdxMap fieldIdxMap
var P6GasAansluitingen2020 fieldMapIdx

var P6Gasm32020Tracker uint16
var P6Gasm32020IdxMap fieldIdxMap
var P6Gasm32020 fieldMapIdx

var P6Kwh2020Tracker uint16
var P6Kwh2020IdxMap fieldIdxMap
var P6Kwh2020 fieldMapIdx

var PandBouwjaarTracker uint16
var PandBouwjaarIdxMap fieldIdxMap
var PandBouwjaar fieldMapIdx

var PandGasAansluitingenTracker uint16
var PandGasAansluitingenIdxMap fieldIdxMap
var PandGasAansluitingen fieldMapIdx

var GebruiksdoelenTracker uint16
var GebruiksdoelenIdxMap fieldIdxMap
var Gebruiksdoelen fieldMapIdx

/*
var {columnname}Tracker uint16
var {columnname}IdxMap fieldIdxMap
var {columnname} fieldMapIdx
var {columnname}Items fieldItemmap
*/

// item map lock
var lock = sync.RWMutex{}

// bitArray Lock
var balock = sync.RWMutex{}

func initBitarrays() {

	WoningTypeItems = make(fieldItemsMap)
	LabelscoreVoorlopigItems = make(fieldItemsMap)
	LabelscoreDefinitiefItems = make(fieldItemsMap)
	GemeentecodeItems = make(fieldItemsMap)
	BuurtcodeItems = make(fieldItemsMap)
	WijkcodeItems = make(fieldItemsMap)
}

func setUpMaps() {
	initBitarrays()
	WoningTypeTracker = 0
	WoningTypeIdxMap = make(fieldIdxMap)
	WoningType = make(fieldMapIdx)

	LabelscoreVoorlopigTracker = 0
	LabelscoreVoorlopigIdxMap = make(fieldIdxMap)
	LabelscoreVoorlopig = make(fieldMapIdx)

	LabelscoreDefinitiefTracker = 0
	LabelscoreDefinitiefIdxMap = make(fieldIdxMap)
	LabelscoreDefinitief = make(fieldMapIdx)

	GemeentecodeTracker = 0
	GemeentecodeIdxMap = make(fieldIdxMap)
	Gemeentecode = make(fieldMapIdx)

	GemeentenaamTracker = 0
	GemeentenaamIdxMap = make(fieldIdxMap)
	Gemeentenaam = make(fieldMapIdx)

	BuurtcodeTracker = 0
	BuurtcodeIdxMap = make(fieldIdxMap)
	Buurtcode = make(fieldMapIdx)

	BuurtnaamTracker = 0
	BuurtnaamIdxMap = make(fieldIdxMap)
	Buurtnaam = make(fieldMapIdx)

	WijkcodeTracker = 0
	WijkcodeIdxMap = make(fieldIdxMap)
	Wijkcode = make(fieldMapIdx)

	WijknaamTracker = 0
	WijknaamIdxMap = make(fieldIdxMap)
	Wijknaam = make(fieldMapIdx)

	ProvinciecodeTracker = 0
	ProvinciecodeIdxMap = make(fieldIdxMap)
	Provinciecode = make(fieldMapIdx)

	ProvinciecodeItems = make(fieldItemsMap)

	ProvincienaamTracker = 0
	ProvincienaamIdxMap = make(fieldIdxMap)
	Provincienaam = make(fieldMapIdx)

	PandGasEanAansluitingenTracker = 0
	PandGasEanAansluitingenIdxMap = make(fieldIdxMap)
	PandGasEanAansluitingen = make(fieldMapIdx)

	P6GasAansluitingen2020Tracker = 0
	P6GasAansluitingen2020IdxMap = make(fieldIdxMap)
	P6GasAansluitingen2020 = make(fieldMapIdx)

	P6Gasm32020Tracker = 0
	P6Gasm32020IdxMap = make(fieldIdxMap)
	P6Gasm32020 = make(fieldMapIdx)

	P6Kwh2020Tracker = 0
	P6Kwh2020IdxMap = make(fieldIdxMap)
	P6Kwh2020 = make(fieldMapIdx)

	PandBouwjaarTracker = 0
	PandBouwjaarIdxMap = make(fieldIdxMap)
	PandBouwjaar = make(fieldMapIdx)

	PandGasAansluitingenTracker = 0
	PandGasAansluitingenIdxMap = make(fieldIdxMap)
	PandGasAansluitingen = make(fieldMapIdx)

	GebruiksdoelenTracker = 0
	GebruiksdoelenIdxMap = make(fieldIdxMap)
	Gebruiksdoelen = make(fieldMapIdx)
}

func CreateMapstore() ModelMaps {
	return ModelMaps{
		WoningTypeTracker,
		WoningTypeIdxMap,
		WoningType,

		LabelscoreVoorlopigTracker,
		LabelscoreVoorlopigIdxMap,
		LabelscoreVoorlopig,

		LabelscoreDefinitiefTracker,
		LabelscoreDefinitiefIdxMap,
		LabelscoreDefinitief,

		GemeentecodeTracker,
		GemeentecodeIdxMap,
		Gemeentecode,

		GemeentenaamTracker,
		GemeentenaamIdxMap,
		Gemeentenaam,

		BuurtcodeTracker,
		BuurtcodeIdxMap,
		Buurtcode,

		BuurtnaamTracker,
		BuurtnaamIdxMap,
		Buurtnaam,

		WijkcodeTracker,
		WijkcodeIdxMap,
		Wijkcode,

		WijknaamTracker,
		WijknaamIdxMap,
		Wijknaam,

		ProvinciecodeTracker,
		ProvinciecodeIdxMap,
		Provinciecode,

		ProvincienaamTracker,
		ProvincienaamIdxMap,
		Provincienaam,

		PandGasEanAansluitingenTracker,
		PandGasEanAansluitingenIdxMap,
		PandGasEanAansluitingen,

		P6GasAansluitingen2020Tracker,
		P6GasAansluitingen2020IdxMap,
		P6GasAansluitingen2020,

		P6Gasm32020Tracker,
		P6Gasm32020IdxMap,
		P6Gasm32020,

		P6Kwh2020Tracker,
		P6Kwh2020IdxMap,
		P6Kwh2020,

		PandBouwjaarTracker,
		PandBouwjaarIdxMap,
		PandBouwjaar,

		PandGasAansluitingenTracker,
		PandGasAansluitingenIdxMap,
		PandGasAansluitingen,

		GebruiksdoelenTracker,
		GebruiksdoelenIdxMap,
		Gebruiksdoelen,
	}
}

func LoadMapstore(m ModelMaps) {

	WoningTypeTracker = m.WoningTypeTracker
	WoningTypeIdxMap = m.WoningTypeIdxMap
	WoningType = m.WoningType

	LabelscoreVoorlopigTracker = m.LabelscoreVoorlopigTracker
	LabelscoreVoorlopigIdxMap = m.LabelscoreVoorlopigIdxMap
	LabelscoreVoorlopig = m.LabelscoreVoorlopig

	LabelscoreDefinitiefTracker = m.LabelscoreDefinitiefTracker
	LabelscoreDefinitiefIdxMap = m.LabelscoreDefinitiefIdxMap
	LabelscoreDefinitief = m.LabelscoreDefinitief

	GemeentecodeTracker = m.GemeentecodeTracker
	GemeentecodeIdxMap = m.GemeentecodeIdxMap
	Gemeentecode = m.Gemeentecode

	GemeentenaamTracker = m.GemeentenaamTracker
	GemeentenaamIdxMap = m.GemeentenaamIdxMap
	Gemeentenaam = m.Gemeentenaam

	BuurtcodeTracker = m.BuurtcodeTracker
	BuurtcodeIdxMap = m.BuurtcodeIdxMap
	Buurtcode = m.Buurtcode

	BuurtnaamTracker = m.BuurtnaamTracker
	BuurtnaamIdxMap = m.BuurtnaamIdxMap
	Buurtnaam = m.Buurtnaam

	WijkcodeTracker = m.WijkcodeTracker
	WijkcodeIdxMap = m.WijkcodeIdxMap
	Wijkcode = m.Wijkcode

	WijknaamTracker = m.WijknaamTracker
	WijknaamIdxMap = m.WijknaamIdxMap
	Wijknaam = m.Wijknaam

	ProvinciecodeTracker = m.ProvinciecodeTracker
	ProvinciecodeIdxMap = m.ProvinciecodeIdxMap
	Provinciecode = m.Provinciecode

	ProvincienaamTracker = m.ProvincienaamTracker
	ProvincienaamIdxMap = m.ProvincienaamIdxMap
	Provincienaam = m.Provincienaam

	PandGasEanAansluitingenTracker = m.PandGasEanAansluitingenTracker
	PandGasEanAansluitingenIdxMap = m.PandGasEanAansluitingenIdxMap
	PandGasEanAansluitingen = m.PandGasEanAansluitingen

	P6GasAansluitingen2020Tracker = m.P6GasAansluitingen2020Tracker
	P6GasAansluitingen2020IdxMap = m.P6GasAansluitingen2020IdxMap
	P6GasAansluitingen2020 = m.P6GasAansluitingen2020

	P6Gasm32020Tracker = m.P6Gasm32020Tracker
	P6Gasm32020IdxMap = m.P6Gasm32020IdxMap
	P6Gasm32020 = m.P6Gasm32020

	P6Kwh2020Tracker = m.P6Kwh2020Tracker
	P6Kwh2020IdxMap = m.P6Kwh2020IdxMap
	P6Kwh2020 = m.P6Kwh2020

	PandBouwjaarTracker = m.PandBouwjaarTracker
	PandBouwjaarIdxMap = m.PandBouwjaarIdxMap
	PandBouwjaar = m.PandBouwjaar

	PandGasAansluitingenTracker = m.PandGasAansluitingenTracker
	PandGasAansluitingenIdxMap = m.PandGasEanAansluitingenIdxMap
	PandGasAansluitingen = m.PandGasAansluitingen

	GebruiksdoelenTracker = m.GebruiksdoelenTracker
	GebruiksdoelenIdxMap = m.GebruiksdoelenIdxMap
	Gebruiksdoelen = m.Gebruiksdoelen
}
