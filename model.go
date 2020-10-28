package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

type fieldIdxMap map[string]uint16
type fieldMapIdx map[uint16]string
type fieldItemmap map[uint16][]*Item

// Column maps.
// Store for each non distinct/repeated column
// unit16 -> string map and
// string -> unit16 map
// track count of distinct values

var PostcodeTracker uint16
var PostcodeIdxMap fieldIdxMap
var Postcode fieldMapIdx

var HuisnummerTracker uint16
var HuisnummerIdxMap fieldIdxMap
var Huisnummer fieldMapIdx

var WoningTypeTracker uint16
var WoningTypeIdxMap fieldIdxMap
var WoningType fieldMapIdx

var LabelscoreVoorlopigTracker uint16
var LabelscoreVoorlopigIdxMap fieldIdxMap
var LabelscoreVoorlopig fieldMapIdx

var LabelscoreDefinitiefTracker uint16
var LabelscoreDefinitiefIdxMap fieldIdxMap
var LabelscoreDefinitief fieldMapIdx

var GemeentecodeTracker uint16
var GemeentecodeIdxMap fieldIdxMap
var Gemeentecode fieldMapIdx

var GemeentenaamTracker uint16
var GemeentenaamIdxMap fieldIdxMap
var Gemeentenaam fieldMapIdx

var BuurtcodeTracker uint16
var BuurtcodeIdxMap fieldIdxMap
var Buurtcode fieldMapIdx

var WijkcodeTracker uint16
var WijkcodeIdxMap fieldIdxMap
var Wijkcode fieldMapIdx

var ProvinciecodeTracker uint16
var ProvinciecodeIdxMap fieldIdxMap
var Provinciecode fieldMapIdx

var ProvincienaamTracker uint16
var ProvincienaamIdxMap fieldIdxMap
var Provincienaam fieldMapIdx

/*
var {columnname}Tracker uint16
var {columnname}IdxMap fieldIdxMap
var {columnname} fieldMapIdx
var {columnname}Items fieldItemmap
*/

var lock = sync.RWMutex{}

func init() {

	PostcodeTracker = 0
	PostcodeIdxMap = make(fieldIdxMap)
	Postcode = make(fieldMapIdx)

	HuisnummerTracker = 0
	HuisnummerIdxMap = make(fieldIdxMap)
	Huisnummer = make(fieldMapIdx)

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

	WijkcodeTracker = 0
	WijkcodeIdxMap = make(fieldIdxMap)
	Wijkcode = make(fieldMapIdx)

	ProvinciecodeTracker = 0
	ProvinciecodeIdxMap = make(fieldIdxMap)
	Provinciecode = make(fieldMapIdx)

	ProvincienaamTracker = 0
	ProvincienaamIdxMap = make(fieldIdxMap)
	Provincienaam = make(fieldMapIdx)

	/*
		labelscoredefinitiefTracker = 0
		labelscoredefinitiefIdxMap = make(fieldIdxMap)
		labelscoredefinitief = make(fieldMapIdx)
	*/
}

type ItemFull struct {
	Pid                  string `json:"pid"`
	Vid                  string `json:"vid"`
	Numid                string `json:"numid"`
	Postcode             string `json:"postcode"`
	Huisnummer           string `json:"huisnummer"`
	Ekey                 string `json:"ekey"`
	WoningType           string `json:"woning_type"`
	LabelscoreVoorlopig  string `json:"labelscore_voorlopig"`
	LabelscoreDefinitief string `json:"labelscore_definitief"`
	Identificatie        string `json:"identificatie"`
	Gemeentecode         string `json:"gemeentecode"`
	Gemeentenaam         string `json:"gemeentenaam"`
	Buurtcode            string `json:"buurtcode"`
	Wijkcode             string `json:"wijkcode"`
	Provinciecode        string `json:"provinciecode"`
	Provincienaam        string `json:"provincienaam"`
}

type Item struct {
	Pid                  string
	Vid                  string
	Numid                string
	Postcode             uint16
	Huisnummer           uint16
	Ekey                 string
	WoningType           uint16
	LabelscoreVoorlopig  uint16
	LabelscoreDefinitief uint16
	Identificatie        string
	Gemeentecode         uint16
	Gemeentenaam         uint16
	Buurtcode            uint16
	Wijkcode             uint16
	Provinciecode        uint16
	Provincienaam        uint16
}

func (i Item) Columns() []string {

	return []string{
		"pid",
		"vid",
		"numid",
		"postcode",
		"huisnummer",
		"ekey",
		"woning_type",
		"labelscore_voorlopig",
		"labelscore_definitief",
		"identificatie",
		"gemeentecode",
		"gemeentenaam",
		"buurtcode",
		"wijkcode",
		"provinciecode",
		"provincienaam",
	}
}

// Shrink create smaller Item using uint16
func (i ItemFull) Shrink() Item {

	lock.Lock()
	defer lock.Unlock()

	//check if column value is already present
	//else store new key
	if _, ok := PostcodeIdxMap[i.Postcode]; !ok {
		// store Postcode in map at current index of tracker
		Postcode[PostcodeTracker] = i.Postcode
		// store key - idx
		PostcodeIdxMap[i.Postcode] = PostcodeTracker
		// increase tracker
		PostcodeTracker += 1
	}

	//check if column value is already present
	//else store new key
	if _, ok := HuisnummerIdxMap[i.Huisnummer]; !ok {
		// store Huisnummer in map at current index of tracker
		Huisnummer[HuisnummerTracker] = i.Huisnummer
		// store key - idx
		HuisnummerIdxMap[i.Huisnummer] = HuisnummerTracker
		// increase tracker
		HuisnummerTracker += 1
	}

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

	return Item{

		i.Pid,
		i.Vid,
		i.Numid,
		PostcodeIdxMap[i.Postcode],
		HuisnummerIdxMap[i.Huisnummer],
		i.Ekey,
		WoningTypeIdxMap[i.WoningType],
		LabelscoreVoorlopigIdxMap[i.LabelscoreVoorlopig],
		LabelscoreDefinitiefIdxMap[i.LabelscoreDefinitief],
		i.Identificatie,
		GemeentecodeIdxMap[i.Gemeentecode],
		GemeentenaamIdxMap[i.Gemeentenaam],
		BuurtcodeIdxMap[i.Buurtcode],
		WijkcodeIdxMap[i.Wijkcode],
		ProvinciecodeIdxMap[i.Provinciecode],
		ProvincienaamIdxMap[i.Provincienaam],
	}
}

func (i Item) Serialize() ItemFull {

	lock.RLock()
	defer lock.RUnlock()

	return ItemFull{

		i.Pid,
		i.Vid,
		i.Numid,
		Postcode[i.Postcode],
		Huisnummer[i.Huisnummer],
		i.Ekey,
		WoningType[i.WoningType],
		LabelscoreVoorlopig[i.LabelscoreVoorlopig],
		LabelscoreDefinitief[i.LabelscoreDefinitief],
		i.Identificatie,
		Gemeentecode[i.Gemeentecode],
		Gemeentenaam[i.Gemeentenaam],
		Buurtcode[i.Buurtcode],
		Wijkcode[i.Wijkcode],
		Provinciecode[i.Provinciecode],
		Provincienaam[i.Provincienaam],
	}
}

func (i ItemFull) Columns() []string {
	return []string{

		"pid",
		"vid",
		"numid",
		"postcode",
		"huisnummer",
		"ekey",
		"woning_type",
		"labelscore_voorlopig",
		"labelscore_definitief",
		"identificatie",
		"gemeentecode",
		"gemeentenaam",
		"buurtcode",
		"wijkcode",
		"provinciecode",
		"provincienaam",
	}
}

func (i Item) Row() []string {

	lock.RLock()
	defer lock.RUnlock()

	return []string{

		i.Pid,
		i.Vid,
		i.Numid,
		Postcode[i.Postcode],
		Huisnummer[i.Huisnummer],
		i.Ekey,
		WoningType[i.WoningType],
		LabelscoreVoorlopig[i.LabelscoreVoorlopig],
		LabelscoreDefinitief[i.LabelscoreDefinitief],
		i.Identificatie,
		Gemeentecode[i.Gemeentecode],
		Gemeentenaam[i.Gemeentenaam],
		Buurtcode[i.Buurtcode],
		Wijkcode[i.Wijkcode],
		Provinciecode[i.Provinciecode],
		Provincienaam[i.Provincienaam],
	}
}

func (i Item) GetIndex() string {
	return GettersPid(&i)
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
	return strings.Contains(Postcode[i.Postcode], s)
}

// startswith filter Postcode
func FilterPostcodeStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(Postcode[i.Postcode], s)
}

// match filters Postcode
func FilterPostcodeMatch(i *Item, s string) bool {
	return Postcode[i.Postcode] == s
}

// getter Postcode
func GettersPostcode(i *Item) string {
	return Postcode[i.Postcode]
}

// contain filter Huisnummer
func FilterHuisnummerContains(i *Item, s string) bool {
	return strings.Contains(Huisnummer[i.Huisnummer], s)
}

// startswith filter Huisnummer
func FilterHuisnummerStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(Huisnummer[i.Huisnummer], s)
}

// match filters Huisnummer
func FilterHuisnummerMatch(i *Item, s string) bool {
	return Huisnummer[i.Huisnummer] == s
}

// getter Huisnummer
func GettersHuisnummer(i *Item) string {
	return Huisnummer[i.Huisnummer]
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

// contain filter Identificatie
func FilterIdentificatieContains(i *Item, s string) bool {
	return strings.Contains(i.Identificatie, s)
}

// startswith filter Identificatie
func FilterIdentificatieStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Identificatie, s)
}

// match filters Identificatie
func FilterIdentificatieMatch(i *Item, s string) bool {
	return i.Identificatie == s
}

// getter Identificatie
func GettersIdentificatie(i *Item) string {
	return i.Identificatie
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
	Funcs   registerFuncType
	GroupBy registerGroupByFunc
	Getters registerGettersMap
	Reduce  registerReduce
}

var Operations GroupedOperations

var RegisterFuncMap registerFuncType
var RegisterGroupBy registerGroupByFunc
var RegisterGetters registerGettersMap
var RegisterReduce registerReduce

func init() {

	RegisterFuncMap = make(registerFuncType)
	RegisterGroupBy = make(registerGroupByFunc)
	RegisterGetters = make(registerGettersMap)
	RegisterReduce = make(registerReduce)

	// register search filter.
	//RegisterFuncMap["search"] = 'EDITYOURSELF'

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

	//register filters for Huisnummer
	RegisterFuncMap["match-huisnummer"] = FilterHuisnummerMatch
	RegisterFuncMap["contains-huisnummer"] = FilterHuisnummerContains
	RegisterFuncMap["startswith-huisnummer"] = FilterHuisnummerStartsWith
	RegisterGetters["huisnummer"] = GettersHuisnummer
	RegisterGroupBy["huisnummer"] = GettersHuisnummer

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

	//register filters for LabelscoreVoorlopig
	RegisterFuncMap["match-labelscore_voorlopig"] = FilterLabelscoreVoorlopigMatch
	RegisterFuncMap["contains-labelscore_voorlopig"] = FilterLabelscoreVoorlopigContains
	RegisterFuncMap["startswith-labelscore_voorlopig"] = FilterLabelscoreVoorlopigStartsWith
	RegisterGetters["labelscore_voorlopig"] = GettersLabelscoreVoorlopig
	RegisterGroupBy["labelscore_voorlopig"] = GettersLabelscoreVoorlopig

	//register filters for LabelscoreDefinitief
	RegisterFuncMap["match-labelscore_definitief"] = FilterLabelscoreDefinitiefMatch
	RegisterFuncMap["contains-labelscore_definitief"] = FilterLabelscoreDefinitiefContains
	RegisterFuncMap["startswith-labelscore_definitief"] = FilterLabelscoreDefinitiefStartsWith
	RegisterGetters["labelscore_definitief"] = GettersLabelscoreDefinitief
	RegisterGroupBy["labelscore_definitief"] = GettersLabelscoreDefinitief

	//register filters for Identificatie
	RegisterFuncMap["match-identificatie"] = FilterIdentificatieMatch
	RegisterFuncMap["contains-identificatie"] = FilterIdentificatieContains
	RegisterFuncMap["startswith-identificatie"] = FilterIdentificatieStartsWith
	RegisterGetters["identificatie"] = GettersIdentificatie
	RegisterGroupBy["identificatie"] = GettersIdentificatie

	//register filters for Gemeentecode
	RegisterFuncMap["match-gemeentecode"] = FilterGemeentecodeMatch
	RegisterFuncMap["contains-gemeentecode"] = FilterGemeentecodeContains
	RegisterFuncMap["startswith-gemeentecode"] = FilterGemeentecodeStartsWith
	RegisterGetters["gemeentecode"] = GettersGemeentecode
	RegisterGroupBy["gemeentecode"] = GettersGemeentecode

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

	//register filters for Wijkcode
	RegisterFuncMap["match-wijkcode"] = FilterWijkcodeMatch
	RegisterFuncMap["contains-wijkcode"] = FilterWijkcodeContains
	RegisterFuncMap["startswith-wijkcode"] = FilterWijkcodeStartsWith
	RegisterGetters["wijkcode"] = GettersWijkcode
	RegisterGroupBy["wijkcode"] = GettersWijkcode

	//register filters for Provinciecode
	RegisterFuncMap["match-provinciecode"] = FilterProvinciecodeMatch
	RegisterFuncMap["contains-provinciecode"] = FilterProvinciecodeContains
	RegisterFuncMap["startswith-provinciecode"] = FilterProvinciecodeStartsWith
	RegisterGetters["provinciecode"] = GettersProvinciecode
	RegisterGroupBy["provinciecode"] = GettersProvinciecode

	//register filters for Provincienaam
	RegisterFuncMap["match-provincienaam"] = FilterProvincienaamMatch
	RegisterFuncMap["contains-provincienaam"] = FilterProvincienaamContains
	RegisterFuncMap["startswith-provincienaam"] = FilterProvincienaamStartsWith
	RegisterGetters["provincienaam"] = GettersProvincienaam
	RegisterGroupBy["provincienaam"] = GettersProvincienaam

	/*
		RegisterFuncMap["match-ekey"] = FilterEkeyMatch
		RegisterFuncMap["contains-ekey"] = FilterEkeyContains
		// register startswith filters
		RegisterFuncMap["startswith-ekey"] = FilterEkeyStartsWith
		// register getters
		RegisterGetters["ekey"] = GettersEkey
		// register groupby
		RegisterGroupBy["ekey"] = GettersEkey

		// register reduce functions
		RegisterReduce["count"] = reduceCount
	*/
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

		"postcode":  func(i, j int) bool { return Postcode[items[i].Postcode] < Postcode[items[j].Postcode] },
		"-postcode": func(i, j int) bool { return Postcode[items[i].Postcode] > Postcode[items[j].Postcode] },

		"huisnummer":  func(i, j int) bool { return Huisnummer[items[i].Huisnummer] < Huisnummer[items[j].Huisnummer] },
		"-huisnummer": func(i, j int) bool { return Huisnummer[items[i].Huisnummer] > Huisnummer[items[j].Huisnummer] },

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

		"identificatie":  func(i, j int) bool { return items[i].Identificatie < items[j].Identificatie },
		"-identificatie": func(i, j int) bool { return items[i].Identificatie > items[j].Identificatie },

		"gemeentecode":  func(i, j int) bool { return Gemeentecode[items[i].Gemeentecode] < Gemeentecode[items[j].Gemeentecode] },
		"-gemeentecode": func(i, j int) bool { return Gemeentecode[items[i].Gemeentecode] > Gemeentecode[items[j].Gemeentecode] },

		"gemeentenaam":  func(i, j int) bool { return Gemeentenaam[items[i].Gemeentenaam] < Gemeentenaam[items[j].Gemeentenaam] },
		"-gemeentenaam": func(i, j int) bool { return Gemeentenaam[items[i].Gemeentenaam] > Gemeentenaam[items[j].Gemeentenaam] },

		"buurtcode":  func(i, j int) bool { return Buurtcode[items[i].Buurtcode] < Buurtcode[items[j].Buurtcode] },
		"-buurtcode": func(i, j int) bool { return Buurtcode[items[i].Buurtcode] > Buurtcode[items[j].Buurtcode] },

		"wijkcode":  func(i, j int) bool { return Wijkcode[items[i].Wijkcode] < Wijkcode[items[j].Wijkcode] },
		"-wijkcode": func(i, j int) bool { return Wijkcode[items[i].Wijkcode] > Wijkcode[items[j].Wijkcode] },

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
