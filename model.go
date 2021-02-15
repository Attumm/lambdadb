/*
	model.go define the 'items' to store.
	All columns with getters and setters are defined here.

	ItemIn, represent rows from the Input data
	Item, the compact item stored in memmory
	ItemOut, defines how and which fields are exported out
	of the API. It is possible to ignore input columns

	Repeated values are stored in maps with int numbers
	as keys.  Optionally bitarrays are created for reapeated
	column values to do fast bit-wise filtering.

	A S2 geo index in created for lat, lon values.

	Unique values are stored as-is.

	The generated codes leaves room to create custom
	index functions yourself to create an API with an
	< 1 ms response time for your specific needs.

	This codebase solves: I need to have an API on this
	tabular dataset fast!
*/

package main

import (
	"errors"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/Workiva/go-datastructures/bitarray"
)

type registerCustomGroupByFunc map[string]func(*Item, ItemsGroupedBy)
type registerGroupByFunc map[string]func(*Item) string
type registerGettersMap map[string]func(*Item) string
type registerReduce map[string]func(Items) map[string]string

type registerBitArray map[string]func(s string) (bitarray.BitArray, error)

type fieldIdxMap map[string]uint16
type fieldMapIdx map[uint16]string
type fieldItemsMap map[uint16]bitarray.BitArray

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

func init() {

	WoningTypeTracker = 0
	WoningTypeIdxMap = make(fieldIdxMap)
	WoningType = make(fieldMapIdx)

	WoningTypeItems = make(fieldItemsMap)

	LabelscoreVoorlopigTracker = 0
	LabelscoreVoorlopigIdxMap = make(fieldIdxMap)
	LabelscoreVoorlopig = make(fieldMapIdx)

	LabelscoreVoorlopigItems = make(fieldItemsMap)

	LabelscoreDefinitiefTracker = 0
	LabelscoreDefinitiefIdxMap = make(fieldIdxMap)
	LabelscoreDefinitief = make(fieldMapIdx)

	LabelscoreDefinitiefItems = make(fieldItemsMap)

	GemeentecodeTracker = 0
	GemeentecodeIdxMap = make(fieldIdxMap)
	Gemeentecode = make(fieldMapIdx)

	GemeentecodeItems = make(fieldItemsMap)

	GemeentenaamTracker = 0
	GemeentenaamIdxMap = make(fieldIdxMap)
	Gemeentenaam = make(fieldMapIdx)

	BuurtcodeTracker = 0
	BuurtcodeIdxMap = make(fieldIdxMap)
	Buurtcode = make(fieldMapIdx)

	BuurtcodeItems = make(fieldItemsMap)

	BuurtnaamTracker = 0
	BuurtnaamIdxMap = make(fieldIdxMap)
	Buurtnaam = make(fieldMapIdx)

	WijkcodeTracker = 0
	WijkcodeIdxMap = make(fieldIdxMap)
	Wijkcode = make(fieldMapIdx)

	WijkcodeItems = make(fieldItemsMap)

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

	GebruiksdoelenTracker = 0
	GebruiksdoelenIdxMap = make(fieldIdxMap)
	Gebruiksdoelen = make(fieldMapIdx)

	/*
		labelscoredefinitiefTracker = 0
		labelscoredefinitiefIdxMap = make(fieldIdxMap)
		labelscoredefinitief = make(fieldMapIdx)
	*/
}

type ItemIn struct {
	Pid                     string `json:"pid"`
	Vid                     string `json:"vid"`
	Numid                   string `json:"numid"`
	Postcode                string `json:"postcode"`
	Oppervlakte             string `json:"oppervlakte"`
	Woningequivalent        string `json:"woningequivalent"`
	Ekey                    string `json:"ekey"`
	WoningType              string `json:"woning_type"`
	LabelscoreVoorlopig     string `json:"labelscore_voorlopig"`
	LabelscoreDefinitief    string `json:"labelscore_definitief"`
	Gemeentecode            string `json:"gemeentecode"`
	Gemeentenaam            string `json:"gemeentenaam"`
	Buurtcode               string `json:"buurtcode"`
	Buurtnaam               string `json:"buurtnaam"`
	Wijkcode                string `json:"wijkcode"`
	Wijknaam                string `json:"wijknaam"`
	Provinciecode           string `json:"provinciecode"`
	Provincienaam           string `json:"provincienaam"`
	Point                   string `json:"point"`
	PandGasEanAansluitingen string `json:"pand_gas_ean_aansluitingen"`
	GroupId2020             string `json:"group_id_2020"`
	GasAansluitingen2020    string `json:"gas_aansluitingen_2020"`
	Gasm32020               string `json:"gasm3_2020"`
	Kwh2020                 string `json:"kwh_2020"`
	Gebruiksdoelen          string `json:"gebruiksdoelen"`
}

type ItemOut struct {
	Pid                     string `json:"pid"`
	Vid                     string `json:"vid"`
	Numid                   string `json:"numid"`
	Postcode                string `json:"postcode"`
	Oppervlakte             string `json:"oppervlakte"`
	Woningequivalent        string `json:"woningequivalent"`
	Ekey                    string `json:"ekey"`
	WoningType              string `json:"woning_type"`
	LabelscoreVoorlopig     string `json:"labelscore_voorlopig"`
	LabelscoreDefinitief    string `json:"labelscore_definitief"`
	Gemeentecode            string `json:"gemeentecode"`
	Gemeentenaam            string `json:"gemeentenaam"`
	Buurtcode               string `json:"buurtcode"`
	Buurtnaam               string `json:"buurtnaam"`
	Wijkcode                string `json:"wijkcode"`
	Wijknaam                string `json:"wijknaam"`
	Provinciecode           string `json:"provinciecode"`
	Provincienaam           string `json:"provincienaam"`
	Point                   string `json:"point"`
	PandGasEanAansluitingen string `json:"pand_gas_ean_aansluitingen"`
	GroupId2020             string `json:"group_id_2020"`
	GasAansluitingen2020    string `json:"gas_aansluitingen_2020"`
	Gasm32020               string `json:"gasm3_2020"`
	Kwh2020                 string `json:"kwh_2020"`
	Gebruiksdoelen          string `json:"gebruiksdoelen"`
}

type Item struct {
	Label                   int // internal index in ITEMS
	Pid                     string
	Vid                     string
	Numid                   string
	Postcode                string
	Oppervlakte             string
	Woningequivalent        string
	Ekey                    string
	WoningType              uint16
	LabelscoreVoorlopig     uint16
	LabelscoreDefinitief    uint16
	Gemeentecode            uint16
	Gemeentenaam            uint16
	Buurtcode               uint16
	Buurtnaam               uint16
	Wijkcode                uint16
	Wijknaam                uint16
	Provinciecode           uint16
	Provincienaam           uint16
	Point                   string
	PandGasEanAansluitingen string
	GroupId2020             string
	GasAansluitingen2020    string
	Gasm32020               string
	Kwh2020                 string
	Gebruiksdoelen          []uint16
}

// Shrink create smaller Item using uint16
func (i ItemIn) Shrink(label int) Item {

	lock.Lock()
	defer lock.Unlock()

	//check if column value is already present
	//else store new key
	if _, ok := WoningTypeIdxMap[i.WoningType]; !ok {
		// store WoningType in map at current index of tracker
		WoningType[WoningTypeTracker] = i.WoningType
		// store key - idx
		WoningTypeIdxMap[i.WoningType] = WoningTypeTracker
		// increase tracker
		WoningTypeTracker += 1
	}

	//check if column value is already present
	//else store new key
	if _, ok := LabelscoreVoorlopigIdxMap[i.LabelscoreVoorlopig]; !ok {
		// store LabelscoreVoorlopig in map at current index of tracker
		LabelscoreVoorlopig[LabelscoreVoorlopigTracker] = i.LabelscoreVoorlopig
		// store key - idx
		LabelscoreVoorlopigIdxMap[i.LabelscoreVoorlopig] = LabelscoreVoorlopigTracker
		// increase tracker
		LabelscoreVoorlopigTracker += 1
	}

	//check if column value is already present
	//else store new key
	if _, ok := LabelscoreDefinitiefIdxMap[i.LabelscoreDefinitief]; !ok {
		// store LabelscoreDefinitief in map at current index of tracker
		LabelscoreDefinitief[LabelscoreDefinitiefTracker] = i.LabelscoreDefinitief
		// store key - idx
		LabelscoreDefinitiefIdxMap[i.LabelscoreDefinitief] = LabelscoreDefinitiefTracker
		// increase tracker
		LabelscoreDefinitiefTracker += 1
	}

	//check if column value is already present
	//else store new key
	if _, ok := GemeentecodeIdxMap[i.Gemeentecode]; !ok {
		// store Gemeentecode in map at current index of tracker
		Gemeentecode[GemeentecodeTracker] = i.Gemeentecode
		// store key - idx
		GemeentecodeIdxMap[i.Gemeentecode] = GemeentecodeTracker
		// increase tracker
		GemeentecodeTracker += 1
	}

	//check if column value is already present
	//else store new key
	if _, ok := GemeentenaamIdxMap[i.Gemeentenaam]; !ok {
		// store Gemeentenaam in map at current index of tracker
		Gemeentenaam[GemeentenaamTracker] = i.Gemeentenaam
		// store key - idx
		GemeentenaamIdxMap[i.Gemeentenaam] = GemeentenaamTracker
		// increase tracker
		GemeentenaamTracker += 1
	}

	//check if column value is already present
	//else store new key
	if _, ok := BuurtcodeIdxMap[i.Buurtcode]; !ok {
		// store Buurtcode in map at current index of tracker
		Buurtcode[BuurtcodeTracker] = i.Buurtcode
		// store key - idx
		BuurtcodeIdxMap[i.Buurtcode] = BuurtcodeTracker
		// increase tracker
		BuurtcodeTracker += 1
	}

	//check if column value is already present
	//else store new key
	if _, ok := BuurtnaamIdxMap[i.Buurtnaam]; !ok {
		// store Buurtnaam in map at current index of tracker
		Buurtnaam[BuurtnaamTracker] = i.Buurtnaam
		// store key - idx
		BuurtnaamIdxMap[i.Buurtnaam] = BuurtnaamTracker
		// increase tracker
		BuurtnaamTracker += 1
	}

	//check if column value is already present
	//else store new key
	if _, ok := WijkcodeIdxMap[i.Wijkcode]; !ok {
		// store Wijkcode in map at current index of tracker
		Wijkcode[WijkcodeTracker] = i.Wijkcode
		// store key - idx
		WijkcodeIdxMap[i.Wijkcode] = WijkcodeTracker
		// increase tracker
		WijkcodeTracker += 1
	}

	//check if column value is already present
	//else store new key
	if _, ok := WijknaamIdxMap[i.Wijknaam]; !ok {
		// store Wijknaam in map at current index of tracker
		Wijknaam[WijknaamTracker] = i.Wijknaam
		// store key - idx
		WijknaamIdxMap[i.Wijknaam] = WijknaamTracker
		// increase tracker
		WijknaamTracker += 1
	}

	//check if column value is already present
	//else store new key
	if _, ok := ProvinciecodeIdxMap[i.Provinciecode]; !ok {
		// store Provinciecode in map at current index of tracker
		Provinciecode[ProvinciecodeTracker] = i.Provinciecode
		// store key - idx
		ProvinciecodeIdxMap[i.Provinciecode] = ProvinciecodeTracker
		// increase tracker
		ProvinciecodeTracker += 1
	}

	//check if column value is already present
	//else store new key
	if _, ok := ProvincienaamIdxMap[i.Provincienaam]; !ok {
		// store Provincienaam in map at current index of tracker
		Provincienaam[ProvincienaamTracker] = i.Provincienaam
		// store key - idx
		ProvincienaamIdxMap[i.Provincienaam] = ProvincienaamTracker
		// increase tracker
		ProvincienaamTracker += 1
	}

	//check if column value is already present
	//else store new key
	if _, ok := GebruiksdoelenIdxMap[i.Gebruiksdoelen]; !ok {
		// store Gebruiksdoelen in map at current index of tracker
		Gebruiksdoelen[GebruiksdoelenTracker] = i.Gebruiksdoelen
		// store key - idx
		GebruiksdoelenIdxMap[i.Gebruiksdoelen] = GebruiksdoelenTracker
		// increase tracker
		GebruiksdoelenTracker += 1
	}

	//check if column value is already present
	//else store new key
	doelen := make([]uint16, 0)

	// parsing {a, b} array values
	// string should be at least 2 example "{}" == size 2
	if len(i.Gebruiksdoelen) > 2 {

		gebruiksdoelen, err := ParsePGArray(i.Gebruiksdoelen)
		if err != nil {
			log.Fatal(err, "error parsing array ")
		}

		for _, gd := range gebruiksdoelen {
			if _, ok := GebruiksdoelenIdxMap[gd]; !ok {
				// store Gebruiksdoelen in map at current index of tracker
				Gebruiksdoelen[GebruiksdoelenTracker] = gd
				// store key - idx
				GebruiksdoelenIdxMap[gd] = GebruiksdoelenTracker
				// increase tracker
				GebruiksdoelenTracker += 1
			}
		}

		for _, v := range gebruiksdoelen {
			doelen = append(doelen, GebruiksdoelenIdxMap[v])
		}
	}

	return Item{

		label,

		i.Pid,
		i.Vid,
		i.Numid,
		i.Postcode,
		i.Oppervlakte,
		i.Woningequivalent,
		i.Ekey,
		WoningTypeIdxMap[i.WoningType],
		LabelscoreVoorlopigIdxMap[i.LabelscoreVoorlopig],
		LabelscoreDefinitiefIdxMap[i.LabelscoreDefinitief],
		GemeentecodeIdxMap[i.Gemeentecode],
		GemeentenaamIdxMap[i.Gemeentenaam],
		BuurtcodeIdxMap[i.Buurtcode],
		BuurtnaamIdxMap[i.Buurtnaam],
		WijkcodeIdxMap[i.Wijkcode],
		WijknaamIdxMap[i.Wijknaam],
		ProvinciecodeIdxMap[i.Provinciecode],
		ProvincienaamIdxMap[i.Provincienaam],
		i.Point,
		i.PandGasEanAansluitingen,
		i.GroupId2020,
		i.GasAansluitingen2020,
		i.Gasm32020,
		i.Kwh2020,
		doelen,
	}
}

// Store selected columns in seperate map[columnvalue]bitarray
// for gast item lookup
func (i Item) StoreBitArrayColumns() {

	balock.Lock()
	defer balock.Unlock()

	lock.RLock()
	defer lock.RUnlock()

	var ba bitarray.BitArray
	var ok bool

	// Column WoningType has byte arrays for
	ba, ok = WoningTypeItems[i.WoningType]
	if !ok {
		ba = bitarray.NewSparseBitArray()
		WoningTypeItems[i.WoningType] = ba
	}

	ba.SetBit(uint64(i.Label))
	// Column LabelscoreVoorlopig has byte arrays for
	ba, ok = LabelscoreVoorlopigItems[i.LabelscoreVoorlopig]
	if !ok {
		ba = bitarray.NewSparseBitArray()
		LabelscoreVoorlopigItems[i.LabelscoreVoorlopig] = ba
	}

	ba.SetBit(uint64(i.Label))
	// Column LabelscoreDefinitief has byte arrays for
	ba, ok = LabelscoreDefinitiefItems[i.LabelscoreDefinitief]
	if !ok {
		ba = bitarray.NewSparseBitArray()
		LabelscoreDefinitiefItems[i.LabelscoreDefinitief] = ba
	}

	ba.SetBit(uint64(i.Label))
	// Column Gemeentecode has byte arrays for
	ba, ok = GemeentecodeItems[i.Gemeentecode]
	if !ok {
		ba = bitarray.NewSparseBitArray()
		GemeentecodeItems[i.Gemeentecode] = ba
	}

	ba.SetBit(uint64(i.Label))
	// Column Buurtcode has byte arrays for
	ba, ok = BuurtcodeItems[i.Buurtcode]
	if !ok {
		ba = bitarray.NewSparseBitArray()
		BuurtcodeItems[i.Buurtcode] = ba
	}

	ba.SetBit(uint64(i.Label))
	// Column Wijkcode has byte arrays for
	ba, ok = WijkcodeItems[i.Wijkcode]
	if !ok {
		ba = bitarray.NewSparseBitArray()
		WijkcodeItems[i.Wijkcode] = ba
	}

	ba.SetBit(uint64(i.Label))
	// Column Provinciecode has byte arrays for
	ba, ok = ProvinciecodeItems[i.Provinciecode]
	if !ok {
		ba = bitarray.NewSparseBitArray()
		ProvinciecodeItems[i.Provinciecode] = ba
	}

	ba.SetBit(uint64(i.Label))

	/*
		// Column Buurtcode has byte arrays for
		ba, ok = BuurtcodeItems[i.Buurtcode]
		if !ok {
			ba = bitarray.NewSparseBitArray()
			BuurtcodeItems[i.Buurtcode] = ba
		}
		ba.SetBit(uint64(i.Label))
	*/

}

func (i Item) Serialize() ItemOut {

	lock.RLock()
	defer lock.RUnlock()

	return ItemOut{

		i.Pid,
		i.Vid,
		i.Numid,
		i.Postcode,
		i.Oppervlakte,
		i.Woningequivalent,
		i.Ekey,
		WoningType[i.WoningType],
		LabelscoreVoorlopig[i.LabelscoreVoorlopig],
		LabelscoreDefinitief[i.LabelscoreDefinitief],
		Gemeentecode[i.Gemeentecode],
		Gemeentenaam[i.Gemeentenaam],
		Buurtcode[i.Buurtcode],
		Buurtnaam[i.Buurtnaam],
		Wijkcode[i.Wijkcode],
		Wijknaam[i.Wijknaam],
		Provinciecode[i.Provinciecode],
		Provincienaam[i.Provincienaam],
		i.Point,
		i.PandGasEanAansluitingen,
		i.GroupId2020,
		i.GasAansluitingen2020,
		i.Gasm32020,
		i.Kwh2020,
		GettersGebruiksdoelen(&i),
	}
}

func (i ItemIn) Columns() []string {
	return []string{

		"pid",
		"vid",
		"numid",
		"postcode",
		"oppervlakte",
		"woningequivalent",
		"ekey",
		"woning_type",
		"labelscore_voorlopig",
		"labelscore_definitief",
		"gemeentecode",
		"gemeentenaam",
		"buurtcode",
		"buurtnaam",
		"wijkcode",
		"wijknaam",
		"provinciecode",
		"provincienaam",
		"point",
		"pand_gas_ean_aansluitingen",
		"group_id_2020",
		"gas_aansluitingen_2020",
		"gasm3_2020",
		"kwh_2020",
		"gebruiksdoelen",
	}
}

func (i ItemOut) Columns() []string {
	return []string{

		"pid",
		"vid",
		"numid",
		"postcode",
		"oppervlakte",
		"woningequivalent",
		"ekey",
		"woning_type",
		"labelscore_voorlopig",
		"labelscore_definitief",
		"gemeentecode",
		"gemeentenaam",
		"buurtcode",
		"buurtnaam",
		"wijkcode",
		"wijknaam",
		"provinciecode",
		"provincienaam",
		"point",
		"pand_gas_ean_aansluitingen",
		"group_id_2020",
		"gas_aansluitingen_2020",
		"gasm3_2020",
		"kwh_2020",
		"gebruiksdoelen",
	}
}

func (i Item) Row() []string {

	lock.RLock()
	defer lock.RUnlock()

	return []string{

		i.Pid,
		i.Vid,
		i.Numid,
		i.Postcode,
		i.Oppervlakte,
		i.Woningequivalent,
		i.Ekey,
		WoningType[i.WoningType],
		LabelscoreVoorlopig[i.LabelscoreVoorlopig],
		LabelscoreDefinitief[i.LabelscoreDefinitief],
		Gemeentecode[i.Gemeentecode],
		Gemeentenaam[i.Gemeentenaam],
		Buurtcode[i.Buurtcode],
		Buurtnaam[i.Buurtnaam],
		Wijkcode[i.Wijkcode],
		Wijknaam[i.Wijknaam],
		Provinciecode[i.Provinciecode],
		Provincienaam[i.Provincienaam],
		i.Point,
		i.PandGasEanAansluitingen,
		i.GroupId2020,
		i.GasAansluitingen2020,
		i.Gasm32020,
		i.Kwh2020,
		GettersGebruiksdoelen(&i),
	}
}

func (i Item) GetIndex() string {
	return GettersEkey(&i)
}

func (i Item) GetGeometry() string {
	return GettersPoint(&i)
}

// contain filter Pid
func FilterPidContains(i *Item, s string) bool {
	return strings.Contains(i.Pid, s)
}

// startswith filter Pid
func FilterPidStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Pid, s)
}

// match filters Pid
func FilterPidMatch(i *Item, s string) bool {
	return i.Pid == s
}

// getter Pid
func GettersPid(i *Item) string {
	return i.Pid
}

// contain filter Vid
func FilterVidContains(i *Item, s string) bool {
	return strings.Contains(i.Vid, s)
}

// startswith filter Vid
func FilterVidStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Vid, s)
}

// match filters Vid
func FilterVidMatch(i *Item, s string) bool {
	return i.Vid == s
}

// getter Vid
func GettersVid(i *Item) string {
	return i.Vid
}

// contain filter Numid
func FilterNumidContains(i *Item, s string) bool {
	return strings.Contains(i.Numid, s)
}

// startswith filter Numid
func FilterNumidStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Numid, s)
}

// match filters Numid
func FilterNumidMatch(i *Item, s string) bool {
	return i.Numid == s
}

// getter Numid
func GettersNumid(i *Item) string {
	return i.Numid
}

// contain filter Postcode
func FilterPostcodeContains(i *Item, s string) bool {
	return strings.Contains(i.Postcode, s)
}

// startswith filter Postcode
func FilterPostcodeStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Postcode, s)
}

// match filters Postcode
func FilterPostcodeMatch(i *Item, s string) bool {
	return i.Postcode == s
}

// getter Postcode
func GettersPostcode(i *Item) string {
	return i.Postcode
}

// contain filter Oppervlakte
func FilterOppervlakteContains(i *Item, s string) bool {
	return strings.Contains(i.Oppervlakte, s)
}

// startswith filter Oppervlakte
func FilterOppervlakteStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Oppervlakte, s)
}

// match filters Oppervlakte
func FilterOppervlakteMatch(i *Item, s string) bool {
	return i.Oppervlakte == s
}

// getter Oppervlakte
func GettersOppervlakte(i *Item) string {
	return i.Oppervlakte
}

// contain filter Woningequivalent
func FilterWoningequivalentContains(i *Item, s string) bool {
	return strings.Contains(i.Woningequivalent, s)
}

// startswith filter Woningequivalent
func FilterWoningequivalentStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Woningequivalent, s)
}

// match filters Woningequivalent
func FilterWoningequivalentMatch(i *Item, s string) bool {
	return i.Woningequivalent == s
}

// getter Woningequivalent
func GettersWoningequivalent(i *Item) string {
	return i.Woningequivalent
}

// contain filter Ekey
func FilterEkeyContains(i *Item, s string) bool {
	return strings.Contains(i.Ekey, s)
}

// startswith filter Ekey
func FilterEkeyStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Ekey, s)
}

// match filters Ekey
func FilterEkeyMatch(i *Item, s string) bool {
	return i.Ekey == s
}

// getter Ekey
func GettersEkey(i *Item) string {
	return i.Ekey
}

// contain filter WoningType
func FilterWoningTypeContains(i *Item, s string) bool {
	return strings.Contains(WoningType[i.WoningType], s)
}

// startswith filter WoningType
func FilterWoningTypeStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(WoningType[i.WoningType], s)
}

// match filters WoningType
func FilterWoningTypeMatch(i *Item, s string) bool {
	return WoningType[i.WoningType] == s
}

// getter WoningType
func GettersWoningType(i *Item) string {
	return WoningType[i.WoningType]
}

// contain filter LabelscoreVoorlopig
func FilterLabelscoreVoorlopigContains(i *Item, s string) bool {
	return strings.Contains(LabelscoreVoorlopig[i.LabelscoreVoorlopig], s)
}

// startswith filter LabelscoreVoorlopig
func FilterLabelscoreVoorlopigStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(LabelscoreVoorlopig[i.LabelscoreVoorlopig], s)
}

// match filters LabelscoreVoorlopig
func FilterLabelscoreVoorlopigMatch(i *Item, s string) bool {
	return LabelscoreVoorlopig[i.LabelscoreVoorlopig] == s
}

// getter LabelscoreVoorlopig
func GettersLabelscoreVoorlopig(i *Item) string {
	return LabelscoreVoorlopig[i.LabelscoreVoorlopig]
}

// contain filter LabelscoreDefinitief
func FilterLabelscoreDefinitiefContains(i *Item, s string) bool {
	return strings.Contains(LabelscoreDefinitief[i.LabelscoreDefinitief], s)
}

// startswith filter LabelscoreDefinitief
func FilterLabelscoreDefinitiefStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(LabelscoreDefinitief[i.LabelscoreDefinitief], s)
}

// match filters LabelscoreDefinitief
func FilterLabelscoreDefinitiefMatch(i *Item, s string) bool {
	return LabelscoreDefinitief[i.LabelscoreDefinitief] == s
}

// getter LabelscoreDefinitief
func GettersLabelscoreDefinitief(i *Item) string {
	return LabelscoreDefinitief[i.LabelscoreDefinitief]
}

// contain filter Gemeentecode
func FilterGemeentecodeContains(i *Item, s string) bool {
	return strings.Contains(Gemeentecode[i.Gemeentecode], s)
}

// startswith filter Gemeentecode
func FilterGemeentecodeStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(Gemeentecode[i.Gemeentecode], s)
}

// match filters Gemeentecode
func FilterGemeentecodeMatch(i *Item, s string) bool {
	return Gemeentecode[i.Gemeentecode] == s
}

// getter Gemeentecode
func GettersGemeentecode(i *Item) string {
	return Gemeentecode[i.Gemeentecode]
}

// contain filter Gemeentenaam
func FilterGemeentenaamContains(i *Item, s string) bool {
	return strings.Contains(Gemeentenaam[i.Gemeentenaam], s)
}

// startswith filter Gemeentenaam
func FilterGemeentenaamStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(Gemeentenaam[i.Gemeentenaam], s)
}

// match filters Gemeentenaam
func FilterGemeentenaamMatch(i *Item, s string) bool {
	return Gemeentenaam[i.Gemeentenaam] == s
}

// getter Gemeentenaam
func GettersGemeentenaam(i *Item) string {
	return Gemeentenaam[i.Gemeentenaam]
}

// contain filter Buurtcode
func FilterBuurtcodeContains(i *Item, s string) bool {
	return strings.Contains(Buurtcode[i.Buurtcode], s)
}

// startswith filter Buurtcode
func FilterBuurtcodeStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(Buurtcode[i.Buurtcode], s)
}

// match filters Buurtcode
func FilterBuurtcodeMatch(i *Item, s string) bool {
	return Buurtcode[i.Buurtcode] == s
}

// getter Buurtcode
func GettersBuurtcode(i *Item) string {
	return Buurtcode[i.Buurtcode]
}

// contain filter Buurtnaam
func FilterBuurtnaamContains(i *Item, s string) bool {
	return strings.Contains(Buurtnaam[i.Buurtnaam], s)
}

// startswith filter Buurtnaam
func FilterBuurtnaamStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(Buurtnaam[i.Buurtnaam], s)
}

// match filters Buurtnaam
func FilterBuurtnaamMatch(i *Item, s string) bool {
	return Buurtnaam[i.Buurtnaam] == s
}

// getter Buurtnaam
func GettersBuurtnaam(i *Item) string {
	return Buurtnaam[i.Buurtnaam]
}

// contain filter Wijkcode
func FilterWijkcodeContains(i *Item, s string) bool {
	return strings.Contains(Wijkcode[i.Wijkcode], s)
}

// startswith filter Wijkcode
func FilterWijkcodeStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(Wijkcode[i.Wijkcode], s)
}

// match filters Wijkcode
func FilterWijkcodeMatch(i *Item, s string) bool {
	return Wijkcode[i.Wijkcode] == s
}

// getter Wijkcode
func GettersWijkcode(i *Item) string {
	return Wijkcode[i.Wijkcode]
}

// contain filter Wijknaam
func FilterWijknaamContains(i *Item, s string) bool {
	return strings.Contains(Wijknaam[i.Wijknaam], s)
}

// startswith filter Wijknaam
func FilterWijknaamStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(Wijknaam[i.Wijknaam], s)
}

// match filters Wijknaam
func FilterWijknaamMatch(i *Item, s string) bool {
	return Wijknaam[i.Wijknaam] == s
}

// getter Wijknaam
func GettersWijknaam(i *Item) string {
	return Wijknaam[i.Wijknaam]
}

// contain filter Provinciecode
func FilterProvinciecodeContains(i *Item, s string) bool {
	return strings.Contains(Provinciecode[i.Provinciecode], s)
}

// startswith filter Provinciecode
func FilterProvinciecodeStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(Provinciecode[i.Provinciecode], s)
}

// match filters Provinciecode
func FilterProvinciecodeMatch(i *Item, s string) bool {
	return Provinciecode[i.Provinciecode] == s
}

// getter Provinciecode
func GettersProvinciecode(i *Item) string {
	return Provinciecode[i.Provinciecode]
}

// contain filter Provincienaam
func FilterProvincienaamContains(i *Item, s string) bool {
	return strings.Contains(Provincienaam[i.Provincienaam], s)
}

// startswith filter Provincienaam
func FilterProvincienaamStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(Provincienaam[i.Provincienaam], s)
}

// match filters Provincienaam
func FilterProvincienaamMatch(i *Item, s string) bool {
	return Provincienaam[i.Provincienaam] == s
}

// getter Provincienaam
func GettersProvincienaam(i *Item) string {
	return Provincienaam[i.Provincienaam]
}

// contain filter Point
func FilterPointContains(i *Item, s string) bool {
	return strings.Contains(i.Point, s)
}

// startswith filter Point
func FilterPointStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Point, s)
}

// match filters Point
func FilterPointMatch(i *Item, s string) bool {
	return i.Point == s
}

// getter Point
func GettersPoint(i *Item) string {
	return i.Point
}

// contain filter PandGasEanAansluitingen
func FilterPandGasEanAansluitingenContains(i *Item, s string) bool {
	return strings.Contains(i.PandGasEanAansluitingen, s)
}

// startswith filter PandGasEanAansluitingen
func FilterPandGasEanAansluitingenStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.PandGasEanAansluitingen, s)
}

// match filters PandGasEanAansluitingen
func FilterPandGasEanAansluitingenMatch(i *Item, s string) bool {
	return i.PandGasEanAansluitingen == s
}

// getter PandGasEanAansluitingen
func GettersPandGasEanAansluitingen(i *Item) string {
	return i.PandGasEanAansluitingen
}

// contain filter GroupId2020
func FilterGroupId2020Contains(i *Item, s string) bool {
	return strings.Contains(i.GroupId2020, s)
}

// startswith filter GroupId2020
func FilterGroupId2020StartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.GroupId2020, s)
}

// match filters GroupId2020
func FilterGroupId2020Match(i *Item, s string) bool {
	return i.GroupId2020 == s
}

// getter GroupId2020
func GettersGroupId2020(i *Item) string {
	return i.GroupId2020
}

// contain filter GasAansluitingen2020
func FilterGasAansluitingen2020Contains(i *Item, s string) bool {
	return strings.Contains(i.GasAansluitingen2020, s)
}

// startswith filter GasAansluitingen2020
func FilterGasAansluitingen2020StartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.GasAansluitingen2020, s)
}

// match filters GasAansluitingen2020
func FilterGasAansluitingen2020Match(i *Item, s string) bool {
	return i.GasAansluitingen2020 == s
}

// getter GasAansluitingen2020
func GettersGasAansluitingen2020(i *Item) string {
	return i.GasAansluitingen2020
}

// contain filter Gasm32020
func FilterGasm32020Contains(i *Item, s string) bool {
	return strings.Contains(i.Gasm32020, s)
}

// startswith filter Gasm32020
func FilterGasm32020StartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Gasm32020, s)
}

// match filters Gasm32020
func FilterGasm32020Match(i *Item, s string) bool {
	return i.Gasm32020 == s
}

// getter Gasm32020
func GettersGasm32020(i *Item) string {
	return i.Gasm32020
}

// contain filter Kwh2020
func FilterKwh2020Contains(i *Item, s string) bool {
	return strings.Contains(i.Kwh2020, s)
}

// startswith filter Kwh2020
func FilterKwh2020StartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Kwh2020, s)
}

// match filters Kwh2020
func FilterKwh2020Match(i *Item, s string) bool {
	return i.Kwh2020 == s
}

// getter Kwh2020
func GettersKwh2020(i *Item) string {
	return i.Kwh2020
}

// contain filter Gebruiksdoelen
func FilterGebruiksdoelenContains(i *Item, s string) bool {
	for _, v := range i.Gebruiksdoelen {
		vs := Gebruiksdoelen[v]
		if strings.Contains(vs, s) {
			return true
		}
	}
	return false
}

// startswith filter Gebruiksdoelen
func FilterGebruiksdoelenStartsWith(i *Item, s string) bool {
	for _, v := range i.Gebruiksdoelen {
		vs := Gebruiksdoelen[v]
		if strings.HasPrefix(vs, s) {
			return true
		}
	}
	return false

}

// match filters Gebruiksdoelen
func FilterGebruiksdoelenMatch(i *Item, s string) bool {
	for _, v := range i.Gebruiksdoelen {
		vs := Gebruiksdoelen[v]
		if vs == s {
			return true
		}
	}
	return false
}

// getter Gebruiksdoelen
func GettersGebruiksdoelen(i *Item) string {
	doelen := make([]string, 0)
	for _, v := range i.Gebruiksdoelen {
		vs := Gebruiksdoelen[v]
		doelen = append(doelen, vs)
	}
	return strings.Join(doelen, ", ")
}

// getter Gebruiksdoelen
func GroupByGettersGebruiksdoelen(item *Item, grouping ItemsGroupedBy) {

	for i := range item.Gebruiksdoelen {
		groupkey := Gebruiksdoelen[item.Gebruiksdoelen[i]]
		grouping[groupkey] = append(grouping[groupkey], item)
	}
}

/*
// contain filters
func FilterEkeyContains(i *Item, s string) bool {
	return strings.Contains(i.Ekey, s)
}


// startswith filters
func FilterEkeyStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Ekey, s)
}


// match filters
func FilterEkeyMatch(i *Item, s string) bool {
	return i.Ekey == s
}

// getters
func GettersEkey(i *Item) string {
	return i.Ekey
}
*/

// reduce functions
func reduceCount(items Items) map[string]string {
	result := make(map[string]string)
	result["count"] = strconv.Itoa(len(items))
	return result
}

type GroupedOperations struct {
	Funcs     registerFuncType
	GroupBy   registerGroupByFunc
	Getters   registerGettersMap
	Reduce    registerReduce
	BitArrays registerBitArray
}

var Operations GroupedOperations

var RegisterFuncMap registerFuncType
var RegisterGroupBy registerGroupByFunc
var RegisterGroupByCustom registerCustomGroupByFunc
var RegisterGetters registerGettersMap
var RegisterReduce registerReduce
var RegisterBitArray registerBitArray

// ValidateRegsiters validate exposed columns do match filter names
func validateRegisters() error {
	var i = ItemOut{}
	var filters = []string{"match", "contains", "startswith"}
	for _, c := range i.Columns() {
		for _, f := range filters {
			if _, ok := RegisterFuncMap[f+"-"+c]; !ok {
				return errors.New(c + " is missing in RegisterMap")
			}
		}
	}
	return nil
}

// GetBitArrayWoningType for given v string see if there is
// a bitarray created.
func GetBitArrayWoningType(v string) (bitarray.BitArray, error) {

	bpi, ok := WoningTypeIdxMap[v]

	if !ok {
		return nil, errors.New("no bitarray filter found for column value WoningType")
	}

	ba, ok := WoningTypeItems[bpi]

	if !ok {
		return nil, errors.New("no bitarray filter found for column idx value WoningType")
	}

	return ba, nil
}

// GetBitArrayLabelscoreVoorlopig for given v string see if there is
// a bitarray created.
func GetBitArrayLabelscoreVoorlopig(v string) (bitarray.BitArray, error) {

	bpi, ok := LabelscoreVoorlopigIdxMap[v]

	if !ok {
		return nil, errors.New("no bitarray filter found for column value LabelscoreVoorlopig")
	}

	ba, ok := LabelscoreVoorlopigItems[bpi]

	if !ok {
		return nil, errors.New("no bitarray filter found for column idx value LabelscoreVoorlopig")
	}

	return ba, nil
}

// GetBitArrayLabelscoreDefinitief for given v string see if there is
// a bitarray created.
func GetBitArrayLabelscoreDefinitief(v string) (bitarray.BitArray, error) {

	bpi, ok := LabelscoreDefinitiefIdxMap[v]

	if !ok {
		return nil, errors.New("no bitarray filter found for column value LabelscoreDefinitief")
	}

	ba, ok := LabelscoreDefinitiefItems[bpi]

	if !ok {
		return nil, errors.New("no bitarray filter found for column idx value LabelscoreDefinitief")
	}

	return ba, nil
}

// GetBitArrayGemeentecode for given v string see if there is
// a bitarray created.
func GetBitArrayGemeentecode(v string) (bitarray.BitArray, error) {

	bpi, ok := GemeentecodeIdxMap[v]

	if !ok {
		return nil, errors.New("no bitarray filter found for column value Gemeentecode")
	}

	ba, ok := GemeentecodeItems[bpi]

	if !ok {
		return nil, errors.New("no bitarray filter found for column idx value Gemeentecode")
	}

	return ba, nil
}

// GetBitArrayBuurtcode for given v string see if there is
// a bitarray created.
func GetBitArrayBuurtcode(v string) (bitarray.BitArray, error) {

	bpi, ok := BuurtcodeIdxMap[v]

	if !ok {
		return nil, errors.New("no bitarray filter found for column value Buurtcode")
	}

	ba, ok := BuurtcodeItems[bpi]

	if !ok {
		return nil, errors.New("no bitarray filter found for column idx value Buurtcode")
	}

	return ba, nil
}

// GetBitArrayWijkcode for given v string see if there is
// a bitarray created.
func GetBitArrayWijkcode(v string) (bitarray.BitArray, error) {

	bpi, ok := WijkcodeIdxMap[v]

	if !ok {
		return nil, errors.New("no bitarray filter found for column value Wijkcode")
	}

	ba, ok := WijkcodeItems[bpi]

	if !ok {
		return nil, errors.New("no bitarray filter found for column idx value Wijkcode")
	}

	return ba, nil
}

// GetBitArrayProvinciecode for given v string see if there is
// a bitarray created.
func GetBitArrayProvinciecode(v string) (bitarray.BitArray, error) {

	bpi, ok := ProvinciecodeIdxMap[v]

	if !ok {
		return nil, errors.New("no bitarray filter found for column value Provinciecode")
	}

	ba, ok := ProvinciecodeItems[bpi]

	if !ok {
		return nil, errors.New("no bitarray filter found for column idx value Provinciecode")
	}

	return ba, nil
}

func init() {

	RegisterFuncMap = make(registerFuncType)
	RegisterGroupBy = make(registerGroupByFunc)
	RegisterGetters = make(registerGettersMap)
	RegisterReduce = make(registerReduce)
	RegisterBitArray = make(registerBitArray)
	RegisterGroupByCustom = make(registerCustomGroupByFunc)

	// register search filter.
	//RegisterFuncMap["search"] = 'EDITYOURSELF'
	// example RegisterFuncMap["search"] = FilterEkeyStartsWith

	//RegisterFuncMap["value"] = 'EDITYOURSELF'
	RegisterGetters["value"] = GettersEkey

	// register filters

	//register filters for Pid
	RegisterFuncMap["match-pid"] = FilterPidMatch
	RegisterFuncMap["contains-pid"] = FilterPidContains
	RegisterFuncMap["startswith-pid"] = FilterPidStartsWith
	RegisterGetters["pid"] = GettersPid
	RegisterGroupBy["pid"] = GettersPid

	//register filters for Vid
	RegisterFuncMap["match-vid"] = FilterVidMatch
	RegisterFuncMap["contains-vid"] = FilterVidContains
	RegisterFuncMap["startswith-vid"] = FilterVidStartsWith
	RegisterGetters["vid"] = GettersVid
	RegisterGroupBy["vid"] = GettersVid

	//register filters for Numid
	RegisterFuncMap["match-numid"] = FilterNumidMatch
	RegisterFuncMap["contains-numid"] = FilterNumidContains
	RegisterFuncMap["startswith-numid"] = FilterNumidStartsWith
	RegisterGetters["numid"] = GettersNumid
	RegisterGroupBy["numid"] = GettersNumid

	//register filters for Postcode
	RegisterFuncMap["match-postcode"] = FilterPostcodeMatch
	RegisterFuncMap["contains-postcode"] = FilterPostcodeContains
	RegisterFuncMap["startswith-postcode"] = FilterPostcodeStartsWith
	RegisterGetters["postcode"] = GettersPostcode
	RegisterGroupBy["postcode"] = GettersPostcode

	//register filters for Oppervlakte
	RegisterFuncMap["match-oppervlakte"] = FilterOppervlakteMatch
	RegisterFuncMap["contains-oppervlakte"] = FilterOppervlakteContains
	RegisterFuncMap["startswith-oppervlakte"] = FilterOppervlakteStartsWith
	RegisterGetters["oppervlakte"] = GettersOppervlakte
	RegisterGroupBy["oppervlakte"] = GettersOppervlakte

	//register filters for Woningequivalent
	RegisterFuncMap["match-woningequivalent"] = FilterWoningequivalentMatch
	RegisterFuncMap["contains-woningequivalent"] = FilterWoningequivalentContains
	RegisterFuncMap["startswith-woningequivalent"] = FilterWoningequivalentStartsWith
	RegisterGetters["woningequivalent"] = GettersWoningequivalent
	RegisterGroupBy["woningequivalent"] = GettersWoningequivalent

	//register filters for Ekey
	RegisterFuncMap["match-ekey"] = FilterEkeyMatch
	RegisterFuncMap["contains-ekey"] = FilterEkeyContains
	RegisterFuncMap["startswith-ekey"] = FilterEkeyStartsWith
	RegisterGetters["ekey"] = GettersEkey
	RegisterGroupBy["ekey"] = GettersEkey

	//register filters for WoningType
	RegisterFuncMap["match-woning_type"] = FilterWoningTypeMatch
	RegisterFuncMap["contains-woning_type"] = FilterWoningTypeContains
	RegisterFuncMap["startswith-woning_type"] = FilterWoningTypeStartsWith
	RegisterGetters["woning_type"] = GettersWoningType
	RegisterGroupBy["woning_type"] = GettersWoningType

	RegisterBitArray["woning_type"] = GetBitArrayWoningType

	//register filters for LabelscoreVoorlopig
	RegisterFuncMap["match-labelscore_voorlopig"] = FilterLabelscoreVoorlopigMatch
	RegisterFuncMap["contains-labelscore_voorlopig"] = FilterLabelscoreVoorlopigContains
	RegisterFuncMap["startswith-labelscore_voorlopig"] = FilterLabelscoreVoorlopigStartsWith
	RegisterGetters["labelscore_voorlopig"] = GettersLabelscoreVoorlopig
	RegisterGroupBy["labelscore_voorlopig"] = GettersLabelscoreVoorlopig

	RegisterBitArray["labelscore_voorlopig"] = GetBitArrayLabelscoreVoorlopig

	//register filters for LabelscoreDefinitief
	RegisterFuncMap["match-labelscore_definitief"] = FilterLabelscoreDefinitiefMatch
	RegisterFuncMap["contains-labelscore_definitief"] = FilterLabelscoreDefinitiefContains
	RegisterFuncMap["startswith-labelscore_definitief"] = FilterLabelscoreDefinitiefStartsWith
	RegisterGetters["labelscore_definitief"] = GettersLabelscoreDefinitief
	RegisterGroupBy["labelscore_definitief"] = GettersLabelscoreDefinitief

	RegisterBitArray["labelscore_definitief"] = GetBitArrayLabelscoreDefinitief

	//register filters for Gemeentecode
	RegisterFuncMap["match-gemeentecode"] = FilterGemeentecodeMatch
	RegisterFuncMap["contains-gemeentecode"] = FilterGemeentecodeContains
	RegisterFuncMap["startswith-gemeentecode"] = FilterGemeentecodeStartsWith
	RegisterGetters["gemeentecode"] = GettersGemeentecode
	RegisterGroupBy["gemeentecode"] = GettersGemeentecode

	RegisterBitArray["gemeentecode"] = GetBitArrayGemeentecode

	//register filters for Gemeentenaam
	RegisterFuncMap["match-gemeentenaam"] = FilterGemeentenaamMatch
	RegisterFuncMap["contains-gemeentenaam"] = FilterGemeentenaamContains
	RegisterFuncMap["startswith-gemeentenaam"] = FilterGemeentenaamStartsWith
	RegisterGetters["gemeentenaam"] = GettersGemeentenaam
	RegisterGroupBy["gemeentenaam"] = GettersGemeentenaam

	//register filters for Buurtcode
	RegisterFuncMap["match-buurtcode"] = FilterBuurtcodeMatch
	RegisterFuncMap["contains-buurtcode"] = FilterBuurtcodeContains
	RegisterFuncMap["startswith-buurtcode"] = FilterBuurtcodeStartsWith
	RegisterGetters["buurtcode"] = GettersBuurtcode
	RegisterGroupBy["buurtcode"] = GettersBuurtcode

	RegisterBitArray["buurtcode"] = GetBitArrayBuurtcode

	//register filters for Buurtnaam
	RegisterFuncMap["match-buurtnaam"] = FilterBuurtnaamMatch
	RegisterFuncMap["contains-buurtnaam"] = FilterBuurtnaamContains
	RegisterFuncMap["startswith-buurtnaam"] = FilterBuurtnaamStartsWith
	RegisterGetters["buurtnaam"] = GettersBuurtnaam
	RegisterGroupBy["buurtnaam"] = GettersBuurtnaam

	//register filters for Wijkcode
	RegisterFuncMap["match-wijkcode"] = FilterWijkcodeMatch
	RegisterFuncMap["contains-wijkcode"] = FilterWijkcodeContains
	RegisterFuncMap["startswith-wijkcode"] = FilterWijkcodeStartsWith
	RegisterGetters["wijkcode"] = GettersWijkcode
	RegisterGroupBy["wijkcode"] = GettersWijkcode

	RegisterBitArray["wijkcode"] = GetBitArrayWijkcode

	//register filters for Wijknaam
	RegisterFuncMap["match-wijknaam"] = FilterWijknaamMatch
	RegisterFuncMap["contains-wijknaam"] = FilterWijknaamContains
	RegisterFuncMap["startswith-wijknaam"] = FilterWijknaamStartsWith
	RegisterGetters["wijknaam"] = GettersWijknaam
	RegisterGroupBy["wijknaam"] = GettersWijknaam

	//register filters for Provinciecode
	RegisterFuncMap["match-provinciecode"] = FilterProvinciecodeMatch
	RegisterFuncMap["contains-provinciecode"] = FilterProvinciecodeContains
	RegisterFuncMap["startswith-provinciecode"] = FilterProvinciecodeStartsWith
	RegisterGetters["provinciecode"] = GettersProvinciecode
	RegisterGroupBy["provinciecode"] = GettersProvinciecode

	RegisterBitArray["provinciecode"] = GetBitArrayProvinciecode

	//register filters for Provincienaam
	RegisterFuncMap["match-provincienaam"] = FilterProvincienaamMatch
	RegisterFuncMap["contains-provincienaam"] = FilterProvincienaamContains
	RegisterFuncMap["startswith-provincienaam"] = FilterProvincienaamStartsWith
	RegisterGetters["provincienaam"] = GettersProvincienaam
	RegisterGroupBy["provincienaam"] = GettersProvincienaam

	//register filters for Point
	RegisterFuncMap["match-point"] = FilterPointMatch
	RegisterFuncMap["contains-point"] = FilterPointContains
	RegisterFuncMap["startswith-point"] = FilterPointStartsWith
	RegisterGetters["point"] = GettersPoint
	RegisterGroupBy["point"] = GettersPoint

	//register filters for PandGasEanAansluitingen
	RegisterFuncMap["match-pand_gas_ean_aansluitingen"] = FilterPandGasEanAansluitingenMatch
	RegisterFuncMap["contains-pand_gas_ean_aansluitingen"] = FilterPandGasEanAansluitingenContains
	RegisterFuncMap["startswith-pand_gas_ean_aansluitingen"] = FilterPandGasEanAansluitingenStartsWith
	RegisterGetters["pand_gas_ean_aansluitingen"] = GettersPandGasEanAansluitingen
	RegisterGroupBy["pand_gas_ean_aansluitingen"] = GettersPandGasEanAansluitingen

	//register filters for GroupId2020
	RegisterFuncMap["match-group_id_2020"] = FilterGroupId2020Match
	RegisterFuncMap["contains-group_id_2020"] = FilterGroupId2020Contains
	RegisterFuncMap["startswith-group_id_2020"] = FilterGroupId2020StartsWith
	RegisterGetters["group_id_2020"] = GettersGroupId2020
	RegisterGroupBy["group_id_2020"] = GettersGroupId2020

	//register filters for GasAansluitingen2020
	RegisterFuncMap["match-gas_aansluitingen_2020"] = FilterGasAansluitingen2020Match
	RegisterFuncMap["contains-gas_aansluitingen_2020"] = FilterGasAansluitingen2020Contains
	RegisterFuncMap["startswith-gas_aansluitingen_2020"] = FilterGasAansluitingen2020StartsWith
	RegisterGetters["gas_aansluitingen_2020"] = GettersGasAansluitingen2020
	RegisterGroupBy["gas_aansluitingen_2020"] = GettersGasAansluitingen2020

	//register filters for Gasm32020
	RegisterFuncMap["match-gasm3_2020"] = FilterGasm32020Match
	RegisterFuncMap["contains-gasm3_2020"] = FilterGasm32020Contains
	RegisterFuncMap["startswith-gasm3_2020"] = FilterGasm32020StartsWith
	RegisterGetters["gasm3_2020"] = GettersGasm32020
	RegisterGroupBy["gasm3_2020"] = GettersGasm32020

	//register filters for Kwh2020
	RegisterFuncMap["match-kwh_2020"] = FilterKwh2020Match
	RegisterFuncMap["contains-kwh_2020"] = FilterKwh2020Contains
	RegisterFuncMap["startswith-kwh_2020"] = FilterKwh2020StartsWith
	RegisterGetters["kwh_2020"] = GettersKwh2020
	RegisterGroupBy["kwh_2020"] = GettersKwh2020

	//register filters for Gebruiksdoelen
	RegisterFuncMap["match-gebruiksdoelen"] = FilterGebruiksdoelenMatch
	RegisterFuncMap["contains-gebruiksdoelen"] = FilterGebruiksdoelenContains
	RegisterFuncMap["startswith-gebruiksdoelen"] = FilterGebruiksdoelenStartsWith
	RegisterGetters["gebruiksdoelen"] = GettersGebruiksdoelen
	RegisterGroupBy["gebruiksdoelen"] = GettersGebruiksdoelen

	RegisterGroupByCustom["gebruiksdoelen-mixed"] = GroupByGettersGebruiksdoelen

	err := validateRegisters()
	if err != nil {
		log.Fatal(err)
	}

	/*
		RegisterFuncMap["match-ekey"] = FilterEkeyMatch
		RegisterFuncMap["contains-ekey"] = FilterEkeyContains
		// register startswith filters
		RegisterFuncMap["startswith-ekey"] = FilterEkeyStartsWith
		// register getters
		RegisterGetters["ekey"] = GettersEkey
		// register groupby
		RegisterGroupBy["ekey"] = GettersEkey

	*/

	// register reduce functions
	RegisterReduce["count"] = reduceCount
}

type sortLookup map[string]func(int, int) bool

func createSort(items Items) sortLookup {

	sortFuncs := sortLookup{

		"pid":  func(i, j int) bool { return items[i].Pid < items[j].Pid },
		"-pid": func(i, j int) bool { return items[i].Pid > items[j].Pid },

		"vid":  func(i, j int) bool { return items[i].Vid < items[j].Vid },
		"-vid": func(i, j int) bool { return items[i].Vid > items[j].Vid },

		"numid":  func(i, j int) bool { return items[i].Numid < items[j].Numid },
		"-numid": func(i, j int) bool { return items[i].Numid > items[j].Numid },

		"postcode":  func(i, j int) bool { return items[i].Postcode < items[j].Postcode },
		"-postcode": func(i, j int) bool { return items[i].Postcode > items[j].Postcode },

		"oppervlakte":  func(i, j int) bool { return items[i].Oppervlakte < items[j].Oppervlakte },
		"-oppervlakte": func(i, j int) bool { return items[i].Oppervlakte > items[j].Oppervlakte },

		"woningequivalent":  func(i, j int) bool { return items[i].Woningequivalent < items[j].Woningequivalent },
		"-woningequivalent": func(i, j int) bool { return items[i].Woningequivalent > items[j].Woningequivalent },

		"ekey":  func(i, j int) bool { return items[i].Ekey < items[j].Ekey },
		"-ekey": func(i, j int) bool { return items[i].Ekey > items[j].Ekey },

		"woning_type":  func(i, j int) bool { return WoningType[items[i].WoningType] < WoningType[items[j].WoningType] },
		"-woning_type": func(i, j int) bool { return WoningType[items[i].WoningType] > WoningType[items[j].WoningType] },

		"labelscore_voorlopig": func(i, j int) bool {
			return LabelscoreVoorlopig[items[i].LabelscoreVoorlopig] < LabelscoreVoorlopig[items[j].LabelscoreVoorlopig]
		},
		"-labelscore_voorlopig": func(i, j int) bool {
			return LabelscoreVoorlopig[items[i].LabelscoreVoorlopig] > LabelscoreVoorlopig[items[j].LabelscoreVoorlopig]
		},

		"labelscore_definitief": func(i, j int) bool {
			return LabelscoreDefinitief[items[i].LabelscoreDefinitief] < LabelscoreDefinitief[items[j].LabelscoreDefinitief]
		},
		"-labelscore_definitief": func(i, j int) bool {
			return LabelscoreDefinitief[items[i].LabelscoreDefinitief] > LabelscoreDefinitief[items[j].LabelscoreDefinitief]
		},

		"gemeentecode":  func(i, j int) bool { return Gemeentecode[items[i].Gemeentecode] < Gemeentecode[items[j].Gemeentecode] },
		"-gemeentecode": func(i, j int) bool { return Gemeentecode[items[i].Gemeentecode] > Gemeentecode[items[j].Gemeentecode] },

		"gemeentenaam":  func(i, j int) bool { return Gemeentenaam[items[i].Gemeentenaam] < Gemeentenaam[items[j].Gemeentenaam] },
		"-gemeentenaam": func(i, j int) bool { return Gemeentenaam[items[i].Gemeentenaam] > Gemeentenaam[items[j].Gemeentenaam] },

		"buurtcode":  func(i, j int) bool { return Buurtcode[items[i].Buurtcode] < Buurtcode[items[j].Buurtcode] },
		"-buurtcode": func(i, j int) bool { return Buurtcode[items[i].Buurtcode] > Buurtcode[items[j].Buurtcode] },

		"buurtnaam":  func(i, j int) bool { return Buurtnaam[items[i].Buurtnaam] < Buurtnaam[items[j].Buurtnaam] },
		"-buurtnaam": func(i, j int) bool { return Buurtnaam[items[i].Buurtnaam] > Buurtnaam[items[j].Buurtnaam] },

		"wijkcode":  func(i, j int) bool { return Wijkcode[items[i].Wijkcode] < Wijkcode[items[j].Wijkcode] },
		"-wijkcode": func(i, j int) bool { return Wijkcode[items[i].Wijkcode] > Wijkcode[items[j].Wijkcode] },

		"wijknaam":  func(i, j int) bool { return Wijknaam[items[i].Wijknaam] < Wijknaam[items[j].Wijknaam] },
		"-wijknaam": func(i, j int) bool { return Wijknaam[items[i].Wijknaam] > Wijknaam[items[j].Wijknaam] },

		"provinciecode": func(i, j int) bool {
			return Provinciecode[items[i].Provinciecode] < Provinciecode[items[j].Provinciecode]
		},
		"-provinciecode": func(i, j int) bool {
			return Provinciecode[items[i].Provinciecode] > Provinciecode[items[j].Provinciecode]
		},

		"provincienaam": func(i, j int) bool {
			return Provincienaam[items[i].Provincienaam] < Provincienaam[items[j].Provincienaam]
		},
		"-provincienaam": func(i, j int) bool {
			return Provincienaam[items[i].Provincienaam] > Provincienaam[items[j].Provincienaam]
		},

		"point":  func(i, j int) bool { return items[i].Point < items[j].Point },
		"-point": func(i, j int) bool { return items[i].Point > items[j].Point },

		"pand_gas_ean_aansluitingen":  func(i, j int) bool { return items[i].PandGasEanAansluitingen < items[j].PandGasEanAansluitingen },
		"-pand_gas_ean_aansluitingen": func(i, j int) bool { return items[i].PandGasEanAansluitingen > items[j].PandGasEanAansluitingen },

		"group_id_2020":  func(i, j int) bool { return items[i].GroupId2020 < items[j].GroupId2020 },
		"-group_id_2020": func(i, j int) bool { return items[i].GroupId2020 > items[j].GroupId2020 },

		"gas_aansluitingen_2020":  func(i, j int) bool { return items[i].GasAansluitingen2020 < items[j].GasAansluitingen2020 },
		"-gas_aansluitingen_2020": func(i, j int) bool { return items[i].GasAansluitingen2020 > items[j].GasAansluitingen2020 },

		"gasm3_2020":  func(i, j int) bool { return items[i].Gasm32020 < items[j].Gasm32020 },
		"-gasm3_2020": func(i, j int) bool { return items[i].Gasm32020 > items[j].Gasm32020 },

		"kwh_2020":  func(i, j int) bool { return items[i].Kwh2020 < items[j].Kwh2020 },
		"-kwh_2020": func(i, j int) bool { return items[i].Kwh2020 > items[j].Kwh2020 },

		"gebruiksdoelen": func(i, j int) bool {
			return GettersGebruiksdoelen(items[i]) < GettersGebruiksdoelen(items[j])
		},
		"-gebruiksdoelen": func(i, j int) bool {
			return GettersGebruiksdoelen(items[i]) > GettersGebruiksdoelen(items[j])
		},

		/*
			"ekey":  func(i, j int) bool { return items[i].Ekey < items[j].Ekey },
			"-ekey": func(i, j int) bool { return items[i].Ekey > items[j].Ekey },
		*/
	}
	return sortFuncs
}

func sortBy(items Items, sortingL []string) (Items, []string) {

	lock.Lock()
	defer lock.Unlock()

	sortFuncs := createSort(items)

	for _, sortFuncName := range sortingL {
		sortFunc := sortFuncs[sortFuncName]
		sort.Slice(items, sortFunc)
	}

	// TODO must be nicer way
	keys := []string{}
	for key := range sortFuncs {
		keys = append(keys, key)
	}

	return items, keys
}
