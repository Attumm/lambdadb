package main

import (
	"sort"
	"strconv"
	"strings"
	"fmt"
)

type Item struct {
	Gid                        string `json:"gid"`
	Identificatie              string `json:"identificatie"`
	Gemeentecode               string `json:"gemeentecode"`
	Gemeentenaam               string `json:"gemeentenaam"`
	Buurtcode                  string `json:"buurtcode"`
	Wijkcode                   string `json:"wijkcode"`
	Provinciecode              string `json:"provinciecode"`
	Provincienaam              string `json:"provincienaam"`
	Bouwjaar                   string `json:"bouwjaar"`
	Hoogte                     string `json:"hoogte"`
	Gas_aansluitingen_2020     string `json:"gas_aansluitingen_2020"`
	Ean_code_count             string `json:"ean_code_count"`
	Elabel_definitief          string `json:"elabel_definitief"`
	Elabel_voorlopig           string `json:"elabel_voorlopig"`
	Gasm3_per_m2               string `json:"gasm3_per_m2"`
	Gasm3_per_m3               string `json:"gasm3_per_m3"`
	Gasm3_2017                 string `json:"gasm3_2017"`
	Gasm3_2018                 string `json:"gasm3_2018"`
	Gasm3_2019                 string `json:"gasm3_2019"`
	Gasm3_2020                 string `json:"gasm3_2020"`
	Kwh_2020                   string `json:"kwh_2020"`
	Kwh_2019                   string `json:"kwh_2019"`
	Kwh_2018                   string `json:"kwh_2018"`
	Kwh_2017                   string `json:"kwh_2017"`
	Kwh_leveringsrichting_2020 string `json:"kwh_leveringsrichting_2020"`
	Kwh_leveringsrichting_2019 string `json:"kwh_leveringsrichting_2019"`
	Kwh_leveringsrichting_2018 string `json:"kwh_leveringsrichting_2018"`
	Group_id_2020              string `json:"group_id_2020"`
	Group_id_2019              string `json:"group_id_2019"`
	Group_id_2018              string `json:"group_id_2018"`
	Pandcount_2020             string `json:"pandcount_2020"`
	Pandcount_2019             string `json:"pandcount_2019"`
	Pandcount_2018             string `json:"pandcount_2018"`
	M2                         string `json:"m2"`
	Totaal_oppervlak_m2        string `json:"totaal_oppervlak_m2"`
	Totaal_volume_m3           string `json:"totaal_volume_m3"`
	Totaal_verbruik_m3         string `json:"totaal_verbruik_m3"`
	Geovlak                    string `json:"geovlak"`
}

func (i Item) Columns() []string {
	return []string{
		"gid",
		"identificatie",
		"gemeentecode",
		"gemeentenaam",
		"buurtcode",
		"wijkcode",
		"provinciecode",
		"provincienaam",
		"bouwjaar",
		"hoogte",
		"gas_aansluitingen_2020",
		"ean_code_count",
		"elabel_definitief",
		"elabel_voorlopig",
		"gasm3_per_m2",
		"gasm3_per_m3",
		"gasm3_2017",
		"gasm3_2018",
		"gasm3_2019",
		"gasm3_2020",
		"kwh_2020",
		"kwh_2019",
		"kwh_2018",
		"kwh_2017",
		"kwh_leveringsrichting_2020",
		"kwh_leveringsrichting_2019",
		"kwh_leveringsrichting_2018",
		"group_id_2020",
		"group_id_2019",
		"group_id_2018",
		"pandcount_2020",
		"pandcount_2019",
		"pandcount_2018",
		"m2",
		"totaal_oppervlak_m2",
		"totaal_volume_m3",
		"totaal_verbruik_m3",
		"geovlak",
	}
}

func (i Item) Row() []string {
	return []string{
		i.Gid,
		i.Identificatie,
		i.Gemeentecode,
		i.Gemeentenaam,
		i.Buurtcode,
		i.Wijkcode,
		i.Provinciecode,
		i.Provincienaam,
		i.Bouwjaar,
		i.Hoogte,
		i.Gas_aansluitingen_2020,
		i.Ean_code_count,
		i.Elabel_definitief,
		i.Elabel_voorlopig,
		i.Gasm3_per_m2,
		i.Gasm3_per_m3,
		i.Gasm3_2017,
		i.Gasm3_2018,
		i.Gasm3_2019,
		i.Gasm3_2020,
		i.Kwh_2020,
		i.Kwh_2019,
		i.Kwh_2018,
		i.Kwh_2017,
		i.Kwh_leveringsrichting_2020,
		i.Kwh_leveringsrichting_2019,
		i.Kwh_leveringsrichting_2018,
		i.Group_id_2020,
		i.Group_id_2019,
		i.Group_id_2018,
		i.Pandcount_2020,
		i.Pandcount_2019,
		i.Pandcount_2018,
		i.M2,
		i.Totaal_oppervlak_m2,
		i.Totaal_volume_m3,
		i.Totaal_verbruik_m3,
		i.Geovlak,
	}
}

// contain filters
func FilterGidContains(i *Item, s string) bool {
	return strings.Contains(i.Gid, s)
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
func FilterBouwjaarContains(i *Item, s string) bool {
	return strings.Contains(i.Bouwjaar, s)
}
func FilterHoogteContains(i *Item, s string) bool {
	return strings.Contains(i.Hoogte, s)
}
func FilterGas_aansluitingen_2020Contains(i *Item, s string) bool {
	return strings.Contains(i.Gas_aansluitingen_2020, s)
}
func FilterEan_code_countContains(i *Item, s string) bool {
	return strings.Contains(i.Ean_code_count, s)
}
func FilterElabel_definitiefContains(i *Item, s string) bool {
	return strings.Contains(i.Elabel_definitief, s)
}
func FilterElabel_voorlopigContains(i *Item, s string) bool {
	return strings.Contains(i.Elabel_voorlopig, s)
}
func FilterGasm3_per_m2Contains(i *Item, s string) bool {
	return strings.Contains(i.Gasm3_per_m2, s)
}
func FilterGasm3_per_m3Contains(i *Item, s string) bool {
	return strings.Contains(i.Gasm3_per_m3, s)
}
func FilterGasm3_2017Contains(i *Item, s string) bool {
	return strings.Contains(i.Gasm3_2017, s)
}
func FilterGasm3_2018Contains(i *Item, s string) bool {
	return strings.Contains(i.Gasm3_2018, s)
}
func FilterGasm3_2019Contains(i *Item, s string) bool {
	return strings.Contains(i.Gasm3_2019, s)
}
func FilterGasm3_2020Contains(i *Item, s string) bool {
	return strings.Contains(i.Gasm3_2020, s)
}
func FilterKwh_2020Contains(i *Item, s string) bool {
	return strings.Contains(i.Kwh_2020, s)
}
func FilterKwh_2019Contains(i *Item, s string) bool {
	return strings.Contains(i.Kwh_2019, s)
}
func FilterKwh_2018Contains(i *Item, s string) bool {
	return strings.Contains(i.Kwh_2018, s)
}
func FilterKwh_2017Contains(i *Item, s string) bool {
	return strings.Contains(i.Kwh_2017, s)
}
func FilterKwh_leveringsrichting_2020Contains(i *Item, s string) bool {
	return strings.Contains(i.Kwh_leveringsrichting_2020, s)
}
func FilterKwh_leveringsrichting_2019Contains(i *Item, s string) bool {
	return strings.Contains(i.Kwh_leveringsrichting_2019, s)
}
func FilterKwh_leveringsrichting_2018Contains(i *Item, s string) bool {
	return strings.Contains(i.Kwh_leveringsrichting_2018, s)
}
func FilterGroup_id_2020Contains(i *Item, s string) bool {
	return strings.Contains(i.Group_id_2020, s)
}
func FilterGroup_id_2019Contains(i *Item, s string) bool {
	return strings.Contains(i.Group_id_2019, s)
}
func FilterGroup_id_2018Contains(i *Item, s string) bool {
	return strings.Contains(i.Group_id_2018, s)
}
func FilterPandcount_2020Contains(i *Item, s string) bool {
	return strings.Contains(i.Pandcount_2020, s)
}
func FilterPandcount_2019Contains(i *Item, s string) bool {
	return strings.Contains(i.Pandcount_2019, s)
}
func FilterPandcount_2018Contains(i *Item, s string) bool {
	return strings.Contains(i.Pandcount_2018, s)
}
func FilterM2Contains(i *Item, s string) bool {
	return strings.Contains(i.M2, s)
}
func FilterTotaal_oppervlak_m2Contains(i *Item, s string) bool {
	return strings.Contains(i.Totaal_oppervlak_m2, s)
}
func FilterTotaal_volume_m3Contains(i *Item, s string) bool {
	return strings.Contains(i.Totaal_volume_m3, s)
}
func FilterTotaal_verbruik_m3Contains(i *Item, s string) bool {
	return strings.Contains(i.Totaal_verbruik_m3, s)
}
func FilterGeovlakContains(i *Item, s string) bool {
	return strings.Contains(i.Geovlak, s)
}

// startswith filters
func FilterGidStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Gid, s)
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
func FilterBouwjaarStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Bouwjaar, s)
}
func FilterHoogteStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Hoogte, s)
}
func FilterGas_aansluitingen_2020StartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Gas_aansluitingen_2020, s)
}
func FilterEan_code_countStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Ean_code_count, s)
}
func FilterElabel_definitiefStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Elabel_definitief, s)
}
func FilterElabel_voorlopigStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Elabel_voorlopig, s)
}
func FilterGasm3_per_m2StartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Gasm3_per_m2, s)
}
func FilterGasm3_per_m3StartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Gasm3_per_m3, s)
}
func FilterGasm3_2017StartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Gasm3_2017, s)
}
func FilterGasm3_2018StartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Gasm3_2018, s)
}
func FilterGasm3_2019StartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Gasm3_2019, s)
}
func FilterGasm3_2020StartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Gasm3_2020, s)
}
func FilterKwh_2020StartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Kwh_2020, s)
}
func FilterKwh_2019StartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Kwh_2019, s)
}
func FilterKwh_2018StartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Kwh_2018, s)
}
func FilterKwh_2017StartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Kwh_2017, s)
}
func FilterKwh_leveringsrichting_2020StartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Kwh_leveringsrichting_2020, s)
}
func FilterKwh_leveringsrichting_2019StartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Kwh_leveringsrichting_2019, s)
}
func FilterKwh_leveringsrichting_2018StartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Kwh_leveringsrichting_2018, s)
}
func FilterGroup_id_2020StartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Group_id_2020, s)
}
func FilterGroup_id_2019StartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Group_id_2019, s)
}
func FilterGroup_id_2018StartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Group_id_2018, s)
}
func FilterPandcount_2020StartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Pandcount_2020, s)
}
func FilterPandcount_2019StartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Pandcount_2019, s)
}
func FilterPandcount_2018StartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Pandcount_2018, s)
}
func FilterM2StartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.M2, s)
}
func FilterTotaal_oppervlak_m2StartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Totaal_oppervlak_m2, s)
}
func FilterTotaal_volume_m3StartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Totaal_volume_m3, s)
}
func FilterTotaal_verbruik_m3StartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Totaal_verbruik_m3, s)
}
func FilterGeovlakStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Geovlak, s)
}

// match filters
func FilterGidMatch(i *Item, s string) bool {
	return i.Gid == s
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
func FilterBouwjaarMatch(i *Item, s string) bool {
	return i.Bouwjaar == s
}
func FilterHoogteMatch(i *Item, s string) bool {
	return i.Hoogte == s
}
func FilterGas_aansluitingen_2020Match(i *Item, s string) bool {
	return i.Gas_aansluitingen_2020 == s
}
func FilterEan_code_countMatch(i *Item, s string) bool {
	return i.Ean_code_count == s
}
func FilterElabel_definitiefMatch(i *Item, s string) bool {
	return i.Elabel_definitief == s
}
func FilterElabel_voorlopigMatch(i *Item, s string) bool {
	return i.Elabel_voorlopig == s
}
func FilterGasm3_per_m2Match(i *Item, s string) bool {
	return i.Gasm3_per_m2 == s
}
func FilterGasm3_per_m3Match(i *Item, s string) bool {
	return i.Gasm3_per_m3 == s
}
func FilterGasm3_2017Match(i *Item, s string) bool {
	return i.Gasm3_2017 == s
}
func FilterGasm3_2018Match(i *Item, s string) bool {
	return i.Gasm3_2018 == s
}
func FilterGasm3_2019Match(i *Item, s string) bool {
	return i.Gasm3_2019 == s
}
func FilterGasm3_2020Match(i *Item, s string) bool {
	return i.Gasm3_2020 == s
}
func FilterKwh_2020Match(i *Item, s string) bool {
	return i.Kwh_2020 == s
}
func FilterKwh_2019Match(i *Item, s string) bool {
	return i.Kwh_2019 == s
}
func FilterKwh_2018Match(i *Item, s string) bool {
	return i.Kwh_2018 == s
}
func FilterKwh_2017Match(i *Item, s string) bool {
	return i.Kwh_2017 == s
}
func FilterKwh_leveringsrichting_2020Match(i *Item, s string) bool {
	return i.Kwh_leveringsrichting_2020 == s
}
func FilterKwh_leveringsrichting_2019Match(i *Item, s string) bool {
	return i.Kwh_leveringsrichting_2019 == s
}
func FilterKwh_leveringsrichting_2018Match(i *Item, s string) bool {
	return i.Kwh_leveringsrichting_2018 == s
}
func FilterGroup_id_2020Match(i *Item, s string) bool {
	return i.Group_id_2020 == s
}
func FilterGroup_id_2019Match(i *Item, s string) bool {
	return i.Group_id_2019 == s
}
func FilterGroup_id_2018Match(i *Item, s string) bool {
	return i.Group_id_2018 == s
}
func FilterPandcount_2020Match(i *Item, s string) bool {
	return i.Pandcount_2020 == s
}
func FilterPandcount_2019Match(i *Item, s string) bool {
	return i.Pandcount_2019 == s
}
func FilterPandcount_2018Match(i *Item, s string) bool {
	return i.Pandcount_2018 == s
}
func FilterM2Match(i *Item, s string) bool {
	return i.M2 == s
}
func FilterTotaal_oppervlak_m2Match(i *Item, s string) bool {
	return i.Totaal_oppervlak_m2 == s
}
func FilterTotaal_volume_m3Match(i *Item, s string) bool {
	return i.Totaal_volume_m3 == s
}
func FilterTotaal_verbruik_m3Match(i *Item, s string) bool {
	return i.Totaal_verbruik_m3 == s
}
func FilterGeovlakMatch(i *Item, s string) bool {
	return i.Geovlak == s
}

// reduce functions

func reduceCount(items Items) map[string]string {
	result := make(map[string]string)
	result["count"] = strconv.Itoa(len(items))
	return result
}

// getters
func GettersGid(i *Item) string {
	return i.Gid
}


func Bla(i *Item) string {
	return fmt.Sprint("%s %s", i.Gemeentecode, i.Gemeentenaam)
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
func GettersBouwjaar(i *Item) string {
	return i.Bouwjaar
}
func GettersHoogte(i *Item) string {
	return i.Hoogte
}
func GettersGas_aansluitingen_2020(i *Item) string {
	return i.Gas_aansluitingen_2020
}
func GettersEan_code_count(i *Item) string {
	return i.Ean_code_count
}
func GettersElabel_definitief(i *Item) string {
	return i.Elabel_definitief
}
func GettersElabel_voorlopig(i *Item) string {
	return i.Elabel_voorlopig
}
func GettersGasm3_per_m2(i *Item) string {
	return i.Gasm3_per_m2
}
func GettersGasm3_per_m3(i *Item) string {
	return i.Gasm3_per_m3
}
func GettersGasm3_2017(i *Item) string {
	return i.Gasm3_2017
}
func GettersGasm3_2018(i *Item) string {
	return i.Gasm3_2018
}
func GettersGasm3_2019(i *Item) string {
	return i.Gasm3_2019
}
func GettersGasm3_2020(i *Item) string {
	return i.Gasm3_2020
}
func GettersKwh_2020(i *Item) string {
	return i.Kwh_2020
}
func GettersKwh_2019(i *Item) string {
	return i.Kwh_2019
}
func GettersKwh_2018(i *Item) string {
	return i.Kwh_2018
}
func GettersKwh_2017(i *Item) string {
	return i.Kwh_2017
}
func GettersKwh_leveringsrichting_2020(i *Item) string {
	return i.Kwh_leveringsrichting_2020
}
func GettersKwh_leveringsrichting_2019(i *Item) string {
	return i.Kwh_leveringsrichting_2019
}
func GettersKwh_leveringsrichting_2018(i *Item) string {
	return i.Kwh_leveringsrichting_2018
}
func GettersGroup_id_2020(i *Item) string {
	return i.Group_id_2020
}
func GettersGroup_id_2019(i *Item) string {
	return i.Group_id_2019
}
func GettersGroup_id_2018(i *Item) string {
	return i.Group_id_2018
}
func GettersPandcount_2020(i *Item) string {
	return i.Pandcount_2020
}
func GettersPandcount_2019(i *Item) string {
	return i.Pandcount_2019
}
func GettersPandcount_2018(i *Item) string {
	return i.Pandcount_2018
}
func GettersM2(i *Item) string {
	return i.M2
}
func GettersTotaal_oppervlak_m2(i *Item) string {
	return i.Totaal_oppervlak_m2
}
func GettersTotaal_volume_m3(i *Item) string {
	return i.Totaal_volume_m3
}
func GettersTotaal_verbruik_m3(i *Item) string {
	return i.Totaal_verbruik_m3
}
func GettersGeovlak(i *Item) string {
	return i.Geovlak
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

	// register match filters
	RegisterFuncMap["match-gid"] = FilterGidMatch
	RegisterFuncMap["match-identificatie"] = FilterIdentificatieMatch
	RegisterFuncMap["match-gemeentecode"] = FilterGemeentecodeMatch
	RegisterFuncMap["match-gemeentenaam"] = FilterGemeentenaamMatch
	RegisterFuncMap["match-buurtcode"] = FilterBuurtcodeMatch
	RegisterFuncMap["match-wijkcode"] = FilterWijkcodeMatch
	RegisterFuncMap["match-provinciecode"] = FilterProvinciecodeMatch
	RegisterFuncMap["match-provincienaam"] = FilterProvincienaamMatch
	RegisterFuncMap["match-bouwjaar"] = FilterBouwjaarMatch
	RegisterFuncMap["match-hoogte"] = FilterHoogteMatch
	RegisterFuncMap["match-gas_aansluitingen_2020"] = FilterGas_aansluitingen_2020Match
	RegisterFuncMap["match-ean_code_count"] = FilterEan_code_countMatch
	RegisterFuncMap["match-elabel_definitief"] = FilterElabel_definitiefMatch
	RegisterFuncMap["match-elabel_voorlopig"] = FilterElabel_voorlopigMatch
	RegisterFuncMap["match-gasm3_per_m2"] = FilterGasm3_per_m2Match
	RegisterFuncMap["match-gasm3_per_m3"] = FilterGasm3_per_m3Match
	RegisterFuncMap["match-gasm3_2017"] = FilterGasm3_2017Match
	RegisterFuncMap["match-gasm3_2018"] = FilterGasm3_2018Match
	RegisterFuncMap["match-gasm3_2019"] = FilterGasm3_2019Match
	RegisterFuncMap["match-gasm3_2020"] = FilterGasm3_2020Match
	RegisterFuncMap["match-kwh_2020"] = FilterKwh_2020Match
	RegisterFuncMap["match-kwh_2019"] = FilterKwh_2019Match
	RegisterFuncMap["match-kwh_2018"] = FilterKwh_2018Match
	RegisterFuncMap["match-kwh_2017"] = FilterKwh_2017Match
	RegisterFuncMap["match-kwh_leveringsrichting_2020"] = FilterKwh_leveringsrichting_2020Match
	RegisterFuncMap["match-kwh_leveringsrichting_2019"] = FilterKwh_leveringsrichting_2019Match
	RegisterFuncMap["match-kwh_leveringsrichting_2018"] = FilterKwh_leveringsrichting_2018Match
	RegisterFuncMap["match-group_id_2020"] = FilterGroup_id_2020Match
	RegisterFuncMap["match-group_id_2019"] = FilterGroup_id_2019Match
	RegisterFuncMap["match-group_id_2018"] = FilterGroup_id_2018Match
	RegisterFuncMap["match-pandcount_2020"] = FilterPandcount_2020Match
	RegisterFuncMap["match-pandcount_2019"] = FilterPandcount_2019Match
	RegisterFuncMap["match-pandcount_2018"] = FilterPandcount_2018Match
	RegisterFuncMap["match-m2"] = FilterM2Match
	RegisterFuncMap["match-totaal_oppervlak_m2"] = FilterTotaal_oppervlak_m2Match
	RegisterFuncMap["match-totaal_volume_m3"] = FilterTotaal_volume_m3Match
	RegisterFuncMap["match-totaal_verbruik_m3"] = FilterTotaal_verbruik_m3Match
	RegisterFuncMap["match-geovlak"] = FilterGeovlakMatch

	// register contains filters
	RegisterFuncMap["contains-gid"] = FilterGidContains
	RegisterFuncMap["contains-identificatie"] = FilterIdentificatieContains
	RegisterFuncMap["contains-gemeentecode"] = FilterGemeentecodeContains
	RegisterFuncMap["contains-gemeentenaam"] = FilterGemeentenaamContains
	RegisterFuncMap["contains-buurtcode"] = FilterBuurtcodeContains
	RegisterFuncMap["contains-wijkcode"] = FilterWijkcodeContains
	RegisterFuncMap["contains-provinciecode"] = FilterProvinciecodeContains
	RegisterFuncMap["contains-provincienaam"] = FilterProvincienaamContains
	RegisterFuncMap["contains-bouwjaar"] = FilterBouwjaarContains
	RegisterFuncMap["contains-hoogte"] = FilterHoogteContains
	RegisterFuncMap["contains-gas_aansluitingen_2020"] = FilterGas_aansluitingen_2020Contains
	RegisterFuncMap["contains-ean_code_count"] = FilterEan_code_countContains
	RegisterFuncMap["contains-elabel_definitief"] = FilterElabel_definitiefContains
	RegisterFuncMap["contains-elabel_voorlopig"] = FilterElabel_voorlopigContains
	RegisterFuncMap["contains-gasm3_per_m2"] = FilterGasm3_per_m2Contains
	RegisterFuncMap["contains-gasm3_per_m3"] = FilterGasm3_per_m3Contains
	RegisterFuncMap["contains-gasm3_2017"] = FilterGasm3_2017Contains
	RegisterFuncMap["contains-gasm3_2018"] = FilterGasm3_2018Contains
	RegisterFuncMap["contains-gasm3_2019"] = FilterGasm3_2019Contains
	RegisterFuncMap["contains-gasm3_2020"] = FilterGasm3_2020Contains
	RegisterFuncMap["contains-kwh_2020"] = FilterKwh_2020Contains
	RegisterFuncMap["contains-kwh_2019"] = FilterKwh_2019Contains
	RegisterFuncMap["contains-kwh_2018"] = FilterKwh_2018Contains
	RegisterFuncMap["contains-kwh_2017"] = FilterKwh_2017Contains
	RegisterFuncMap["contains-kwh_leveringsrichting_2020"] = FilterKwh_leveringsrichting_2020Contains
	RegisterFuncMap["contains-kwh_leveringsrichting_2019"] = FilterKwh_leveringsrichting_2019Contains
	RegisterFuncMap["contains-kwh_leveringsrichting_2018"] = FilterKwh_leveringsrichting_2018Contains
	RegisterFuncMap["contains-group_id_2020"] = FilterGroup_id_2020Contains
	RegisterFuncMap["contains-group_id_2019"] = FilterGroup_id_2019Contains
	RegisterFuncMap["contains-group_id_2018"] = FilterGroup_id_2018Contains
	RegisterFuncMap["contains-pandcount_2020"] = FilterPandcount_2020Contains
	RegisterFuncMap["contains-pandcount_2019"] = FilterPandcount_2019Contains
	RegisterFuncMap["contains-pandcount_2018"] = FilterPandcount_2018Contains
	RegisterFuncMap["contains-m2"] = FilterM2Contains
	RegisterFuncMap["contains-totaal_oppervlak_m2"] = FilterTotaal_oppervlak_m2Contains
	RegisterFuncMap["contains-totaal_volume_m3"] = FilterTotaal_volume_m3Contains
	RegisterFuncMap["contains-totaal_verbruik_m3"] = FilterTotaal_verbruik_m3Contains
	RegisterFuncMap["contains-geovlak"] = FilterGeovlakContains

	// register startswith filters
	RegisterFuncMap["startswith-gid"] = FilterGidStartsWith
	RegisterFuncMap["startswith-identificatie"] = FilterIdentificatieStartsWith
	RegisterFuncMap["startswith-gemeentecode"] = FilterGemeentecodeStartsWith
	RegisterFuncMap["startswith-gemeentenaam"] = FilterGemeentenaamStartsWith
	RegisterFuncMap["startswith-buurtcode"] = FilterBuurtcodeStartsWith
	RegisterFuncMap["startswith-wijkcode"] = FilterWijkcodeStartsWith
	RegisterFuncMap["startswith-provinciecode"] = FilterProvinciecodeStartsWith
	RegisterFuncMap["startswith-provincienaam"] = FilterProvincienaamStartsWith
	RegisterFuncMap["startswith-bouwjaar"] = FilterBouwjaarStartsWith
	RegisterFuncMap["startswith-hoogte"] = FilterHoogteStartsWith
	RegisterFuncMap["startswith-gas_aansluitingen_2020"] = FilterGas_aansluitingen_2020StartsWith
	RegisterFuncMap["startswith-ean_code_count"] = FilterEan_code_countStartsWith
	RegisterFuncMap["startswith-elabel_definitief"] = FilterElabel_definitiefStartsWith
	RegisterFuncMap["startswith-elabel_voorlopig"] = FilterElabel_voorlopigStartsWith
	RegisterFuncMap["startswith-gasm3_per_m2"] = FilterGasm3_per_m2StartsWith
	RegisterFuncMap["startswith-gasm3_per_m3"] = FilterGasm3_per_m3StartsWith
	RegisterFuncMap["startswith-gasm3_2017"] = FilterGasm3_2017StartsWith
	RegisterFuncMap["startswith-gasm3_2018"] = FilterGasm3_2018StartsWith
	RegisterFuncMap["startswith-gasm3_2019"] = FilterGasm3_2019StartsWith
	RegisterFuncMap["startswith-gasm3_2020"] = FilterGasm3_2020StartsWith
	RegisterFuncMap["startswith-kwh_2020"] = FilterKwh_2020StartsWith
	RegisterFuncMap["startswith-kwh_2019"] = FilterKwh_2019StartsWith
	RegisterFuncMap["startswith-kwh_2018"] = FilterKwh_2018StartsWith
	RegisterFuncMap["startswith-kwh_2017"] = FilterKwh_2017StartsWith
	RegisterFuncMap["startswith-kwh_leveringsrichting_2020"] = FilterKwh_leveringsrichting_2020StartsWith
	RegisterFuncMap["startswith-kwh_leveringsrichting_2019"] = FilterKwh_leveringsrichting_2019StartsWith
	RegisterFuncMap["startswith-kwh_leveringsrichting_2018"] = FilterKwh_leveringsrichting_2018StartsWith
	RegisterFuncMap["startswith-group_id_2020"] = FilterGroup_id_2020StartsWith
	RegisterFuncMap["startswith-group_id_2019"] = FilterGroup_id_2019StartsWith
	RegisterFuncMap["startswith-group_id_2018"] = FilterGroup_id_2018StartsWith
	RegisterFuncMap["startswith-pandcount_2020"] = FilterPandcount_2020StartsWith
	RegisterFuncMap["startswith-pandcount_2019"] = FilterPandcount_2019StartsWith
	RegisterFuncMap["startswith-pandcount_2018"] = FilterPandcount_2018StartsWith
	RegisterFuncMap["startswith-m2"] = FilterM2StartsWith
	RegisterFuncMap["startswith-totaal_oppervlak_m2"] = FilterTotaal_oppervlak_m2StartsWith
	RegisterFuncMap["startswith-totaal_volume_m3"] = FilterTotaal_volume_m3StartsWith
	RegisterFuncMap["startswith-totaal_verbruik_m3"] = FilterTotaal_verbruik_m3StartsWith
	RegisterFuncMap["startswith-geovlak"] = FilterGeovlakStartsWith

	// register getters
	RegisterGetters["gid"] = GettersGid
	RegisterGetters["identificatie"] = GettersIdentificatie
	RegisterGetters["gemeentecode"] = GettersGemeentecode
	RegisterGetters["gemeentenaam"] = GettersGemeentenaam
	RegisterGetters["buurtcode"] = GettersBuurtcode
	RegisterGetters["wijkcode"] = GettersWijkcode
	RegisterGetters["provinciecode"] = GettersProvinciecode
	RegisterGetters["provincienaam"] = GettersProvincienaam
	RegisterGetters["bouwjaar"] = GettersBouwjaar
	RegisterGetters["hoogte"] = GettersHoogte
	RegisterGetters["gas_aansluitingen_2020"] = GettersGas_aansluitingen_2020
	RegisterGetters["ean_code_count"] = GettersEan_code_count
	RegisterGetters["elabel_definitief"] = GettersElabel_definitief
	RegisterGetters["elabel_voorlopig"] = GettersElabel_voorlopig
	RegisterGetters["gasm3_per_m2"] = GettersGasm3_per_m2
	RegisterGetters["gasm3_per_m3"] = GettersGasm3_per_m3
	RegisterGetters["gasm3_2017"] = GettersGasm3_2017
	RegisterGetters["gasm3_2018"] = GettersGasm3_2018
	RegisterGetters["gasm3_2019"] = GettersGasm3_2019
	RegisterGetters["gasm3_2020"] = GettersGasm3_2020
	RegisterGetters["kwh_2020"] = GettersKwh_2020
	RegisterGetters["kwh_2019"] = GettersKwh_2019
	RegisterGetters["kwh_2018"] = GettersKwh_2018
	RegisterGetters["kwh_2017"] = GettersKwh_2017
	RegisterGetters["kwh_leveringsrichting_2020"] = GettersKwh_leveringsrichting_2020
	RegisterGetters["kwh_leveringsrichting_2019"] = GettersKwh_leveringsrichting_2019
	RegisterGetters["kwh_leveringsrichting_2018"] = GettersKwh_leveringsrichting_2018
	RegisterGetters["group_id_2020"] = GettersGroup_id_2020
	RegisterGetters["group_id_2019"] = GettersGroup_id_2019
	RegisterGetters["group_id_2018"] = GettersGroup_id_2018
	RegisterGetters["pandcount_2020"] = GettersPandcount_2020
	RegisterGetters["pandcount_2019"] = GettersPandcount_2019
	RegisterGetters["pandcount_2018"] = GettersPandcount_2018
	RegisterGetters["m2"] = GettersM2
	RegisterGetters["totaal_oppervlak_m2"] = GettersTotaal_oppervlak_m2
	RegisterGetters["totaal_volume_m3"] = GettersTotaal_volume_m3
	RegisterGetters["totaal_verbruik_m3"] = GettersTotaal_verbruik_m3
	RegisterGetters["geovlak"] = GettersGeovlak

	RegisterGetters["foobar"] = Bla

	// register groupby
	RegisterGroupBy["gid"] = GettersGid
	RegisterGroupBy["foobar"] = Bla
	RegisterGroupBy["identificatie"] = GettersIdentificatie
	RegisterGroupBy["gemeentecode"] = GettersGemeentecode
	RegisterGroupBy["gemeentenaam"] = GettersGemeentenaam
	RegisterGroupBy["buurtcode"] = GettersBuurtcode
	RegisterGroupBy["wijkcode"] = GettersWijkcode
	RegisterGroupBy["provinciecode"] = GettersProvinciecode
	RegisterGroupBy["provincienaam"] = GettersProvincienaam
	RegisterGroupBy["bouwjaar"] = GettersBouwjaar
	RegisterGroupBy["hoogte"] = GettersHoogte
	RegisterGroupBy["gas_aansluitingen_2020"] = GettersGas_aansluitingen_2020
	RegisterGroupBy["ean_code_count"] = GettersEan_code_count
	RegisterGroupBy["elabel_definitief"] = GettersElabel_definitief
	RegisterGroupBy["elabel_voorlopig"] = GettersElabel_voorlopig
	RegisterGroupBy["gasm3_per_m2"] = GettersGasm3_per_m2
	RegisterGroupBy["gasm3_per_m3"] = GettersGasm3_per_m3
	RegisterGroupBy["gasm3_2017"] = GettersGasm3_2017
	RegisterGroupBy["gasm3_2018"] = GettersGasm3_2018
	RegisterGroupBy["gasm3_2019"] = GettersGasm3_2019
	RegisterGroupBy["gasm3_2020"] = GettersGasm3_2020
	RegisterGroupBy["kwh_2020"] = GettersKwh_2020
	RegisterGroupBy["kwh_2019"] = GettersKwh_2019
	RegisterGroupBy["kwh_2018"] = GettersKwh_2018
	RegisterGroupBy["kwh_2017"] = GettersKwh_2017
	RegisterGroupBy["kwh_leveringsrichting_2020"] = GettersKwh_leveringsrichting_2020
	RegisterGroupBy["kwh_leveringsrichting_2019"] = GettersKwh_leveringsrichting_2019
	RegisterGroupBy["kwh_leveringsrichting_2018"] = GettersKwh_leveringsrichting_2018
	RegisterGroupBy["group_id_2020"] = GettersGroup_id_2020
	RegisterGroupBy["group_id_2019"] = GettersGroup_id_2019
	RegisterGroupBy["group_id_2018"] = GettersGroup_id_2018
	RegisterGroupBy["pandcount_2020"] = GettersPandcount_2020
	RegisterGroupBy["pandcount_2019"] = GettersPandcount_2019
	RegisterGroupBy["pandcount_2018"] = GettersPandcount_2018
	RegisterGroupBy["m2"] = GettersM2
	RegisterGroupBy["totaal_oppervlak_m2"] = GettersTotaal_oppervlak_m2
	RegisterGroupBy["totaal_volume_m3"] = GettersTotaal_volume_m3
	RegisterGroupBy["totaal_verbruik_m3"] = GettersTotaal_verbruik_m3
	RegisterGroupBy["geovlak"] = GettersGeovlak

	// register reduce functions
	RegisterReduce["count"] = reduceCount
}
func sortBy(items Items, sortingL []string) (Items, []string) {
	sortFuncs := map[string]func(int, int) bool{"gid": func(i, j int) bool { return items[i].Gid < items[j].Gid },
		"-gid": func(i, j int) bool { return items[i].Gid > items[j].Gid },

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

		"bouwjaar":  func(i, j int) bool { return items[i].Bouwjaar < items[j].Bouwjaar },
		"-bouwjaar": func(i, j int) bool { return items[i].Bouwjaar > items[j].Bouwjaar },

		"hoogte":  func(i, j int) bool { return items[i].Hoogte < items[j].Hoogte },
		"-hoogte": func(i, j int) bool { return items[i].Hoogte > items[j].Hoogte },

		"gas_aansluitingen_2020":  func(i, j int) bool { return items[i].Gas_aansluitingen_2020 < items[j].Gas_aansluitingen_2020 },
		"-gas_aansluitingen_2020": func(i, j int) bool { return items[i].Gas_aansluitingen_2020 > items[j].Gas_aansluitingen_2020 },

		"ean_code_count":  func(i, j int) bool { return items[i].Ean_code_count < items[j].Ean_code_count },
		"-ean_code_count": func(i, j int) bool { return items[i].Ean_code_count > items[j].Ean_code_count },

		"elabel_definitief":  func(i, j int) bool { return items[i].Elabel_definitief < items[j].Elabel_definitief },
		"-elabel_definitief": func(i, j int) bool { return items[i].Elabel_definitief > items[j].Elabel_definitief },

		"elabel_voorlopig":  func(i, j int) bool { return items[i].Elabel_voorlopig < items[j].Elabel_voorlopig },
		"-elabel_voorlopig": func(i, j int) bool { return items[i].Elabel_voorlopig > items[j].Elabel_voorlopig },

		"gasm3_per_m2":  func(i, j int) bool { return items[i].Gasm3_per_m2 < items[j].Gasm3_per_m2 },
		"-gasm3_per_m2": func(i, j int) bool { return items[i].Gasm3_per_m2 > items[j].Gasm3_per_m2 },

		"gasm3_per_m3":  func(i, j int) bool { return items[i].Gasm3_per_m3 < items[j].Gasm3_per_m3 },
		"-gasm3_per_m3": func(i, j int) bool { return items[i].Gasm3_per_m3 > items[j].Gasm3_per_m3 },

		"gasm3_2017":  func(i, j int) bool { return items[i].Gasm3_2017 < items[j].Gasm3_2017 },
		"-gasm3_2017": func(i, j int) bool { return items[i].Gasm3_2017 > items[j].Gasm3_2017 },

		"gasm3_2018":  func(i, j int) bool { return items[i].Gasm3_2018 < items[j].Gasm3_2018 },
		"-gasm3_2018": func(i, j int) bool { return items[i].Gasm3_2018 > items[j].Gasm3_2018 },

		"gasm3_2019":  func(i, j int) bool { return items[i].Gasm3_2019 < items[j].Gasm3_2019 },
		"-gasm3_2019": func(i, j int) bool { return items[i].Gasm3_2019 > items[j].Gasm3_2019 },

		"gasm3_2020":  func(i, j int) bool { return items[i].Gasm3_2020 < items[j].Gasm3_2020 },
		"-gasm3_2020": func(i, j int) bool { return items[i].Gasm3_2020 > items[j].Gasm3_2020 },

		"kwh_2020":  func(i, j int) bool { return items[i].Kwh_2020 < items[j].Kwh_2020 },
		"-kwh_2020": func(i, j int) bool { return items[i].Kwh_2020 > items[j].Kwh_2020 },

		"kwh_2019":  func(i, j int) bool { return items[i].Kwh_2019 < items[j].Kwh_2019 },
		"-kwh_2019": func(i, j int) bool { return items[i].Kwh_2019 > items[j].Kwh_2019 },

		"kwh_2018":  func(i, j int) bool { return items[i].Kwh_2018 < items[j].Kwh_2018 },
		"-kwh_2018": func(i, j int) bool { return items[i].Kwh_2018 > items[j].Kwh_2018 },

		"kwh_2017":  func(i, j int) bool { return items[i].Kwh_2017 < items[j].Kwh_2017 },
		"-kwh_2017": func(i, j int) bool { return items[i].Kwh_2017 > items[j].Kwh_2017 },

		"kwh_leveringsrichting_2020":  func(i, j int) bool { return items[i].Kwh_leveringsrichting_2020 < items[j].Kwh_leveringsrichting_2020 },
		"-kwh_leveringsrichting_2020": func(i, j int) bool { return items[i].Kwh_leveringsrichting_2020 > items[j].Kwh_leveringsrichting_2020 },

		"kwh_leveringsrichting_2019":  func(i, j int) bool { return items[i].Kwh_leveringsrichting_2019 < items[j].Kwh_leveringsrichting_2019 },
		"-kwh_leveringsrichting_2019": func(i, j int) bool { return items[i].Kwh_leveringsrichting_2019 > items[j].Kwh_leveringsrichting_2019 },

		"kwh_leveringsrichting_2018":  func(i, j int) bool { return items[i].Kwh_leveringsrichting_2018 < items[j].Kwh_leveringsrichting_2018 },
		"-kwh_leveringsrichting_2018": func(i, j int) bool { return items[i].Kwh_leveringsrichting_2018 > items[j].Kwh_leveringsrichting_2018 },

		"group_id_2020":  func(i, j int) bool { return items[i].Group_id_2020 < items[j].Group_id_2020 },
		"-group_id_2020": func(i, j int) bool { return items[i].Group_id_2020 > items[j].Group_id_2020 },

		"group_id_2019":  func(i, j int) bool { return items[i].Group_id_2019 < items[j].Group_id_2019 },
		"-group_id_2019": func(i, j int) bool { return items[i].Group_id_2019 > items[j].Group_id_2019 },

		"group_id_2018":  func(i, j int) bool { return items[i].Group_id_2018 < items[j].Group_id_2018 },
		"-group_id_2018": func(i, j int) bool { return items[i].Group_id_2018 > items[j].Group_id_2018 },

		"pandcount_2020":  func(i, j int) bool { return items[i].Pandcount_2020 < items[j].Pandcount_2020 },
		"-pandcount_2020": func(i, j int) bool { return items[i].Pandcount_2020 > items[j].Pandcount_2020 },

		"pandcount_2019":  func(i, j int) bool { return items[i].Pandcount_2019 < items[j].Pandcount_2019 },
		"-pandcount_2019": func(i, j int) bool { return items[i].Pandcount_2019 > items[j].Pandcount_2019 },

		"pandcount_2018":  func(i, j int) bool { return items[i].Pandcount_2018 < items[j].Pandcount_2018 },
		"-pandcount_2018": func(i, j int) bool { return items[i].Pandcount_2018 > items[j].Pandcount_2018 },

		"m2":  func(i, j int) bool { return items[i].M2 < items[j].M2 },
		"-m2": func(i, j int) bool { return items[i].M2 > items[j].M2 },

		"totaal_oppervlak_m2":  func(i, j int) bool { return items[i].Totaal_oppervlak_m2 < items[j].Totaal_oppervlak_m2 },
		"-totaal_oppervlak_m2": func(i, j int) bool { return items[i].Totaal_oppervlak_m2 > items[j].Totaal_oppervlak_m2 },

		"totaal_volume_m3":  func(i, j int) bool { return items[i].Totaal_volume_m3 < items[j].Totaal_volume_m3 },
		"-totaal_volume_m3": func(i, j int) bool { return items[i].Totaal_volume_m3 > items[j].Totaal_volume_m3 },

		"totaal_verbruik_m3":  func(i, j int) bool { return items[i].Totaal_verbruik_m3 < items[j].Totaal_verbruik_m3 },
		"-totaal_verbruik_m3": func(i, j int) bool { return items[i].Totaal_verbruik_m3 > items[j].Totaal_verbruik_m3 },

		"geovlak":  func(i, j int) bool { return items[i].Geovlak < items[j].Geovlak },
		"-geovlak": func(i, j int) bool { return items[i].Geovlak > items[j].Geovlak },
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
