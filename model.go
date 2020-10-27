package main

import (
	"sort"
	"strconv"
	"strings"
)

type Item struct {
	Pid                   string `json:"pid"`
	Vid                   string `json:"vid"`
	Numid                 string `json:"numid"`
	Postcode              string `json:"postcode"`
	Huisnummer            string `json:"huisnummer"`
	Ekey                  string `json:"ekey"`
	Woning_type           string `json:"woning_type"`
	Labelscore_voorlopig  string `json:"labelscore_voorlopig"`
	Labelscore_definitief string `json:"labelscore_definitief"`
	Identificatie         string `json:"identificatie"`
	Gemeentecode          string `json:"gemeentecode"`
	Gemeentenaam          string `json:"gemeentenaam"`
	Buurtcode             string `json:"buurtcode"`
	Wijkcode              string `json:"wijkcode"`
	Provinciecode         string `json:"provinciecode"`
	Provincienaam         string `json:"provincienaam"`
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

func (i Item) Row() []string {
	return []string{
		i.Pid,
		i.Vid,
		i.Numid,
		i.Postcode,
		i.Huisnummer,
		i.Ekey,
		i.Woning_type,
		i.Labelscore_voorlopig,
		i.Labelscore_definitief,
		i.Identificatie,
		i.Gemeentecode,
		i.Gemeentenaam,
		i.Buurtcode,
		i.Wijkcode,
		i.Provinciecode,
		i.Provincienaam,
	}
}

func (i Item) GetIndex() string {
	return i.Buurtcode
}

func FilterSearch(i *Item, s string) bool {

	return (strings.Contains(i.Postcode, s) ||
		strings.Contains(i.Buurtcode, s) ||
		strings.Contains(i.Provincienaam, s) ||
		strings.Contains(i.Gemeentenaam, s))
}

// contain filters
func FilterPidContains(i *Item, s string) bool {
	return strings.Contains(i.Pid, s)
}
func FilterVidContains(i *Item, s string) bool {
	return strings.Contains(i.Vid, s)
}
func FilterNumidContains(i *Item, s string) bool {
	return strings.Contains(i.Numid, s)
}
func FilterPostcodeContains(i *Item, s string) bool {
	return strings.Contains(i.Postcode, s)
}
func FilterHuisnummerContains(i *Item, s string) bool {
	return strings.Contains(i.Huisnummer, s)
}
func FilterEkeyContains(i *Item, s string) bool {
	return strings.Contains(i.Ekey, s)
}
func FilterWoning_typeContains(i *Item, s string) bool {
	return strings.Contains(i.Woning_type, s)
}
func FilterLabelscore_voorlopigContains(i *Item, s string) bool {
	return strings.Contains(i.Labelscore_voorlopig, s)
}
func FilterLabelscore_definitiefContains(i *Item, s string) bool {
	return strings.Contains(i.Labelscore_definitief, s)
}
func FilterIdentificatieContains(i *Item, s string) bool {
	return strings.Contains(i.Identificatie, s)
}
func FilterGemeentecodeContains(i *Item, s string) bool {
	return strings.Contains(i.Gemeentecode, s)
}
func FilterGemeentenaamContains(i *Item, s string) bool {
	return strings.Contains(i.Gemeentenaam, s)
}
func FilterBuurtcodeContains(i *Item, s string) bool {
	return strings.Contains(i.Buurtcode, s)
}
func FilterWijkcodeContains(i *Item, s string) bool {
	return strings.Contains(i.Wijkcode, s)
}
func FilterProvinciecodeContains(i *Item, s string) bool {
	return strings.Contains(i.Provinciecode, s)
}
func FilterProvincienaamContains(i *Item, s string) bool {
	return strings.Contains(i.Provincienaam, s)
}

// startswith filters
func FilterPidStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Pid, s)
}
func FilterVidStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Vid, s)
}
func FilterNumidStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Numid, s)
}
func FilterPostcodeStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Postcode, s)
}
func FilterHuisnummerStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Huisnummer, s)
}
func FilterEkeyStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Ekey, s)
}
func FilterWoning_typeStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Woning_type, s)
}
func FilterLabelscore_voorlopigStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Labelscore_voorlopig, s)
}
func FilterLabelscore_definitiefStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Labelscore_definitief, s)
}
func FilterIdentificatieStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Identificatie, s)
}
func FilterGemeentecodeStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Gemeentecode, s)
}
func FilterGemeentenaamStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Gemeentenaam, s)
}
func FilterBuurtcodeStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Buurtcode, s)
}
func FilterWijkcodeStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Wijkcode, s)
}
func FilterProvinciecodeStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Provinciecode, s)
}
func FilterProvincienaamStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Provincienaam, s)
}

// match filters
func FilterPidMatch(i *Item, s string) bool {
	return i.Pid == s
}
func FilterVidMatch(i *Item, s string) bool {
	return i.Vid == s
}
func FilterNumidMatch(i *Item, s string) bool {
	return i.Numid == s
}
func FilterPostcodeMatch(i *Item, s string) bool {
	return i.Postcode == s
}
func FilterHuisnummerMatch(i *Item, s string) bool {
	return i.Huisnummer == s
}
func FilterEkeyMatch(i *Item, s string) bool {
	return i.Ekey == s
}
func FilterWoning_typeMatch(i *Item, s string) bool {
	return i.Woning_type == s
}
func FilterLabelscore_voorlopigMatch(i *Item, s string) bool {
	return i.Labelscore_voorlopig == s
}
func FilterLabelscore_definitiefMatch(i *Item, s string) bool {
	return i.Labelscore_definitief == s
}

func FilterIdentificatieMatch(i *Item, s string) bool {
	return i.Identificatie == s
}
func FilterGemeentecodeMatch(i *Item, s string) bool {
	return i.Gemeentecode == s
}
func FilterGemeentenaamMatch(i *Item, s string) bool {
	return i.Gemeentenaam == s
}
func FilterBuurtcodeMatch(i *Item, s string) bool {
	return i.Buurtcode == s
}
func FilterWijkcodeMatch(i *Item, s string) bool {
	return i.Wijkcode == s
}
func FilterProvinciecodeMatch(i *Item, s string) bool {
	return i.Provinciecode == s
}
func FilterProvincienaamMatch(i *Item, s string) bool {
	return i.Provincienaam == s
}

// reduce functions

func reduceCount(items Items) map[string]string {
	result := make(map[string]string)
	result["count"] = strconv.Itoa(len(items))
	return result
}

// getters
func GettersPid(i *Item) string {
	return i.Pid
}
func GettersVid(i *Item) string {
	return i.Vid
}
func GettersNumid(i *Item) string {
	return i.Numid
}
func GettersPostcode(i *Item) string {
	return i.Postcode
}
func GettersHuisnummer(i *Item) string {
	return i.Huisnummer
}
func GettersEkey(i *Item) string {
	return i.Ekey
}
func GettersWoning_type(i *Item) string {
	return i.Woning_type
}
func GettersLabelscore_voorlopig(i *Item) string {
	return i.Labelscore_voorlopig
}
func GettersLabelscore_definitief(i *Item) string {
	return i.Labelscore_definitief
}
func GettersIdentificatie(i *Item) string {
	return i.Identificatie
}
func GettersGemeentecode(i *Item) string {
	return i.Gemeentecode
}
func GettersGemeentenaam(i *Item) string {
	return i.Gemeentenaam
}
func GettersBuurtcode(i *Item) string {
	return i.Buurtcode
}
func GettersWijkcode(i *Item) string {
	return i.Wijkcode
}
func GettersProvinciecode(i *Item) string {
	return i.Provinciecode
}
func GettersProvincienaam(i *Item) string {
	return i.Provincienaam
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

	RegisterFuncMap["search"] = FilterSearch

	// register match filters
	RegisterFuncMap["match-pid"] = FilterPidMatch
	RegisterFuncMap["match-vid"] = FilterVidMatch
	RegisterFuncMap["match-numid"] = FilterNumidMatch
	RegisterFuncMap["match-postcode"] = FilterPostcodeMatch
	RegisterFuncMap["match-huisnummer"] = FilterHuisnummerMatch
	RegisterFuncMap["match-ekey"] = FilterEkeyMatch
	RegisterFuncMap["match-woning_type"] = FilterWoning_typeMatch
	RegisterFuncMap["match-labelscore_voorlopig"] = FilterLabelscore_voorlopigMatch
	RegisterFuncMap["match-labelscore_definitief"] = FilterLabelscore_definitiefMatch
	RegisterFuncMap["match-identificatie"] = FilterIdentificatieMatch
	RegisterFuncMap["match-gemeentecode"] = FilterGemeentecodeMatch
	RegisterFuncMap["match-gemeentenaam"] = FilterGemeentenaamMatch
	RegisterFuncMap["match-buurtcode"] = FilterBuurtcodeMatch
	RegisterFuncMap["match-wijkcode"] = FilterWijkcodeMatch
	RegisterFuncMap["match-provinciecode"] = FilterProvinciecodeMatch
	RegisterFuncMap["match-provincienaam"] = FilterProvincienaamMatch

	// register contains filters
	RegisterFuncMap["contains-pid"] = FilterPidContains
	RegisterFuncMap["contains-vid"] = FilterVidContains
	RegisterFuncMap["contains-numid"] = FilterNumidContains
	RegisterFuncMap["contains-postcode"] = FilterPostcodeContains
	RegisterFuncMap["contains-huisnummer"] = FilterHuisnummerContains
	RegisterFuncMap["contains-ekey"] = FilterEkeyContains
	RegisterFuncMap["contains-woning_type"] = FilterWoning_typeContains
	RegisterFuncMap["contains-labelscore_voorlopig"] = FilterLabelscore_voorlopigContains
	RegisterFuncMap["contains-labelscore_definitief"] = FilterLabelscore_definitiefContains
	RegisterFuncMap["contains-identificatie"] = FilterIdentificatieContains
	RegisterFuncMap["contains-gemeentecode"] = FilterGemeentecodeContains
	RegisterFuncMap["contains-gemeentenaam"] = FilterGemeentenaamContains
	RegisterFuncMap["contains-buurtcode"] = FilterBuurtcodeContains
	RegisterFuncMap["contains-wijkcode"] = FilterWijkcodeContains
	RegisterFuncMap["contains-provinciecode"] = FilterProvinciecodeContains
	RegisterFuncMap["contains-provincienaam"] = FilterProvincienaamContains

	// register startswith filters
	RegisterFuncMap["startswith-pid"] = FilterPidStartsWith
	RegisterFuncMap["startswith-vid"] = FilterVidStartsWith
	RegisterFuncMap["startswith-numid"] = FilterNumidStartsWith
	RegisterFuncMap["startswith-postcode"] = FilterPostcodeStartsWith
	RegisterFuncMap["startswith-huisnummer"] = FilterHuisnummerStartsWith
	RegisterFuncMap["startswith-ekey"] = FilterEkeyStartsWith
	RegisterFuncMap["startswith-woning_type"] = FilterWoning_typeStartsWith
	RegisterFuncMap["startswith-labelscore_voorlopig"] = FilterLabelscore_voorlopigStartsWith
	RegisterFuncMap["startswith-labelscore_definitief"] = FilterLabelscore_definitiefStartsWith
	RegisterFuncMap["startswith-identificatie"] = FilterIdentificatieStartsWith
	RegisterFuncMap["startswith-gemeentecode"] = FilterGemeentecodeStartsWith
	RegisterFuncMap["startswith-gemeentenaam"] = FilterGemeentenaamStartsWith
	RegisterFuncMap["startswith-buurtcode"] = FilterBuurtcodeStartsWith
	RegisterFuncMap["startswith-wijkcode"] = FilterWijkcodeStartsWith
	RegisterFuncMap["startswith-provinciecode"] = FilterProvinciecodeStartsWith
	RegisterFuncMap["startswith-provincienaam"] = FilterProvincienaamStartsWith

	// register getters
	RegisterGetters["pid"] = GettersPid
	RegisterGetters["vid"] = GettersVid
	RegisterGetters["numid"] = GettersNumid
	RegisterGetters["postcode"] = GettersPostcode
	RegisterGetters["huisnummer"] = GettersHuisnummer
	RegisterGetters["ekey"] = GettersEkey
	RegisterGetters["woning_type"] = GettersWoning_type
	RegisterGetters["labelscore_voorlopig"] = GettersLabelscore_voorlopig
	RegisterGetters["labelscore_definitief"] = GettersLabelscore_definitief
	RegisterGetters["identificatie"] = GettersIdentificatie
	RegisterGetters["gemeentecode"] = GettersGemeentecode
	RegisterGetters["gemeentenaam"] = GettersGemeentenaam
	RegisterGetters["buurtcode"] = GettersBuurtcode
	RegisterGetters["wijkcode"] = GettersWijkcode
	RegisterGetters["provinciecode"] = GettersProvinciecode
	RegisterGetters["provincienaam"] = GettersProvincienaam

	// register groupby
	RegisterGroupBy["pid"] = GettersPid
	RegisterGroupBy["vid"] = GettersVid
	RegisterGroupBy["numid"] = GettersNumid
	RegisterGroupBy["postcode"] = GettersPostcode
	RegisterGroupBy["huisnummer"] = GettersHuisnummer
	RegisterGroupBy["ekey"] = GettersEkey
	RegisterGroupBy["woning_type"] = GettersWoning_type
	RegisterGroupBy["labelscore_voorlopig"] = GettersLabelscore_voorlopig
	RegisterGroupBy["labelscore_definitief"] = GettersLabelscore_definitief
	RegisterGroupBy["identificatie"] = GettersIdentificatie
	RegisterGroupBy["gemeentecode"] = GettersGemeentecode
	RegisterGroupBy["gemeentenaam"] = GettersGemeentenaam
	RegisterGroupBy["buurtcode"] = GettersBuurtcode
	RegisterGroupBy["wijkcode"] = GettersWijkcode
	RegisterGroupBy["provinciecode"] = GettersProvinciecode
	RegisterGroupBy["provincienaam"] = GettersProvincienaam

	// register reduce functions
	RegisterReduce["count"] = reduceCount
}
func sortBy(items Items, sortingL []string) (Items, []string) {
	sortFuncs := map[string]func(int, int) bool{"pid": func(i, j int) bool { return items[i].Pid < items[j].Pid },
		"-pid": func(i, j int) bool { return items[i].Pid > items[j].Pid },

		"vid":  func(i, j int) bool { return items[i].Vid < items[j].Vid },
		"-vid": func(i, j int) bool { return items[i].Vid > items[j].Vid },

		"numid":  func(i, j int) bool { return items[i].Numid < items[j].Numid },
		"-numid": func(i, j int) bool { return items[i].Numid > items[j].Numid },

		"postcode":  func(i, j int) bool { return items[i].Postcode < items[j].Postcode },
		"-postcode": func(i, j int) bool { return items[i].Postcode > items[j].Postcode },

		"huisnummer":  func(i, j int) bool { return items[i].Huisnummer < items[j].Huisnummer },
		"-huisnummer": func(i, j int) bool { return items[i].Huisnummer > items[j].Huisnummer },

		"ekey":  func(i, j int) bool { return items[i].Ekey < items[j].Ekey },
		"-ekey": func(i, j int) bool { return items[i].Ekey > items[j].Ekey },

		"woning_type":  func(i, j int) bool { return items[i].Woning_type < items[j].Woning_type },
		"-woning_type": func(i, j int) bool { return items[i].Woning_type > items[j].Woning_type },

		"labelscore_voorlopig":  func(i, j int) bool { return items[i].Labelscore_voorlopig < items[j].Labelscore_voorlopig },
		"-labelscore_voorlopig": func(i, j int) bool { return items[i].Labelscore_voorlopig > items[j].Labelscore_voorlopig },

		"labelscore_definitief":  func(i, j int) bool { return items[i].Labelscore_definitief < items[j].Labelscore_definitief },
		"-labelscore_definitief": func(i, j int) bool { return items[i].Labelscore_definitief > items[j].Labelscore_definitief },

		"identificatie":  func(i, j int) bool { return items[i].Identificatie < items[j].Identificatie },
		"-identificatie": func(i, j int) bool { return items[i].Identificatie > items[j].Identificatie },

		"gemeentecode":  func(i, j int) bool { return items[i].Gemeentecode < items[j].Gemeentecode },
		"-gemeentecode": func(i, j int) bool { return items[i].Gemeentecode > items[j].Gemeentecode },

		"gemeentenaam":  func(i, j int) bool { return items[i].Gemeentenaam < items[j].Gemeentenaam },
		"-gemeentenaam": func(i, j int) bool { return items[i].Gemeentenaam > items[j].Gemeentenaam },

		"buurtcode":  func(i, j int) bool { return items[i].Buurtcode < items[j].Buurtcode },
		"-buurtcode": func(i, j int) bool { return items[i].Buurtcode > items[j].Buurtcode },

		"wijkcode":  func(i, j int) bool { return items[i].Wijkcode < items[j].Wijkcode },
		"-wijkcode": func(i, j int) bool { return items[i].Wijkcode > items[j].Wijkcode },

		"provinciecode":  func(i, j int) bool { return items[i].Provinciecode < items[j].Provinciecode },
		"-provinciecode": func(i, j int) bool { return items[i].Provinciecode > items[j].Provinciecode },

		"provincienaam":  func(i, j int) bool { return items[i].Provincienaam < items[j].Provincienaam },
		"-provincienaam": func(i, j int) bool { return items[i].Provincienaam > items[j].Provincienaam },
	}
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
