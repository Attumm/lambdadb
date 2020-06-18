package main

import (
	"sort"
	"strings"
)

type Item struct {
	ID        string   `json:"id"`
	Value     string   `json:"value"`
	Type      string   `json:"type"`
	Name      string   `json:"name"`
	Vendor    string   `json:"vendor"`
	IP        string   `json:"ip"`
	Dn        []string `json:"dn"`
	ValueType string   `json:"value_type"`
	Country   string   `json:"country"`
}

// All the function underneath should become generated from the above Item struct

func (i Item) Columns() []string {
	return []string{
		"id",
		"name",
		"value",
		"vendor",
		"type",
		"ip",
		"dn",
		"value_type",
		"country",
	}
}

func (i Item) Row() []string {
	return []string{
		i.ID,
		i.Value,
		i.Type,
		i.Name,
		i.Vendor,
		i.IP,
		strings.Join(i.Dn, "."),
		i.ValueType,
		i.Country,
	}
}

// Filter Functions
func FilterValueContains(i *Item, s string) bool {
	return strings.Contains(i.Value, s)
}

func FilterValueStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Value, s)
}

func FilterValueMatch(i *Item, s string) bool {
	return i.Value == s
}

func FilterNameMatch(i *Item, s string) bool {
	return i.Name == s
}

func FilterIPMatch(i *Item, s string) bool {
	return i.IP == s
}

func FilterCountryMatch(i *Item, s string) bool {
	return i.Country == s
}

func FilterDNMatch(i *Item, s string) bool {
	return strings.Join(i.Dn, ".") == s
}

func FilterValueContainsCase(i *Item, s string) bool {
	return strings.Contains(strings.ToLower(i.Value), strings.ToLower(s))
}

// match
func FilterMatchID(i *Item, s string) bool {
	return i.ID == s
}

func FilterMatchType(i *Item, s string) bool {
	return i.Type == s
}

func FilterMatchName(i *Item, s string) bool {
	return i.Name == s
}

func FilterMatchVendor(i *Item, s string) bool {
	return i.Vendor == s
}

func FilterMatchIP(i *Item, s string) bool {
	return i.IP == s
}

func FilterMatchDN(i *Item, s string) bool {
	return strings.Join(i.Dn, ".") == s
}

func FilterMatchValueType(i *Item, s string) bool {
	return i.ValueType == s
}

func FilterMatchCountry(i *Item, s string) bool {
	return i.Country == s
}

// Contains
func FilterContainsID(i *Item, s string) bool {
	return strings.Contains(i.ID, s)
}

func FilterContainsType(i *Item, s string) bool {
	return strings.Contains(i.Type, s)
}

func FilterContainsName(i *Item, s string) bool {
	return strings.Contains(i.Name, s)
}

func FilterContainsVendor(i *Item, s string) bool {
	return strings.Contains(i.Vendor, s)
}

func FilterContainsIP(i *Item, s string) bool {
	return strings.Contains(i.IP, s)
}

func FilterContainsDN(i *Item, s string) bool {
	return strings.Contains(strings.Join(i.Dn, "."), s)
}

func FilterContainsValueType(i *Item, s string) bool {
	return strings.Contains(i.ValueType, s)
}

func FilterContainsCountry(i *Item, s string) bool {
	return strings.Contains(i.Country, s)
}

// StartsWith
func FilterStartsWithID(i *Item, s string) bool {
	return strings.HasPrefix(i.ID, s)
}

func FilterStartsWithType(i *Item, s string) bool {
	return strings.HasPrefix(i.Type, s)
}

func FilterStartsWithName(i *Item, s string) bool {
	return strings.HasPrefix(i.Name, s)
}

func FilterStartsWithVendor(i *Item, s string) bool {
	return strings.HasPrefix(i.Vendor, s)
}

func FilterStartsWithDN(i *Item, s string) bool {
	return strings.HasPrefix(strings.Join(i.Dn, "."), s)
}

func FilterStartsWithIP(i *Item, s string) bool {
	return strings.HasPrefix(i.IP, s)
}

func FilterStartsWithValueType(i *Item, s string) bool {
	return strings.HasPrefix(i.ValueType, s)
}

func FilterStartsWithCountry(i *Item, s string) bool {
	return strings.HasPrefix(i.Country, s)
}

// Getter functions
func GettersID(i *Item) string {
	return i.ID
}

func GettersValue(i *Item) string {
	return i.Value
}

func GettersType(i *Item) string {
	return i.Type
}

func GettersName(i *Item) string {
	return i.Name
}

func GettersVendor(i *Item) string {
	return i.Vendor
}

func GettersDN(i *Item) string {
	return strings.Join(i.Dn, ".")
}

func GettersValueType(i *Item) string {
	return i.ValueType
}

func GettersIP(i *Item) string {
	return i.IP
}

func GettersCountry(i *Item) string {
	return i.Country
}

type GroupedOperations struct {
	Funcs   registerFuncType
	GroupBy registerGroupByFunc
	Getters registerGettersMap
}

var Operations GroupedOperations

var RegisterFuncMap registerFuncType
var RegisterGroupBy registerGroupByFunc
var RegisterGetters registerGettersMap

func init() {

	RegisterFuncMap = make(registerFuncType)
	RegisterFuncMap["contains"] = FilterValueContains
	RegisterFuncMap["contains-case"] = FilterValueContainsCase
	RegisterFuncMap["match"] = FilterValueMatch
	RegisterFuncMap["starts-with"] = FilterValueStartsWith

	// operations
	RegisterFuncMap["typeahead"] = FilterValueStartsWith
	RegisterFuncMap["search"] = FilterValueContainsCase
	RegisterFuncMap["q"] = FilterValueContainsCase

	// match
	RegisterFuncMap["match-id"] = FilterMatchID
	RegisterFuncMap["match-type"] = FilterMatchType
	RegisterFuncMap["match-name"] = FilterMatchName
	RegisterFuncMap["match-vendor"] = FilterMatchVendor
	RegisterFuncMap["match-ip"] = FilterMatchIP
	RegisterFuncMap["match-dn"] = FilterMatchDN
	RegisterFuncMap["match-valuetype"] = FilterMatchValueType
	RegisterFuncMap["match-country"] = FilterMatchCountry
	RegisterFuncMap["match-value"] = FilterValueMatch

	// contains
	RegisterFuncMap["contains-id"] = FilterContainsID
	RegisterFuncMap["contains-type"] = FilterContainsType
	RegisterFuncMap["contains-name"] = FilterContainsName
	RegisterFuncMap["contains-vendor"] = FilterContainsVendor
	RegisterFuncMap["contains-ip"] = FilterContainsIP
	RegisterFuncMap["contains-dn"] = FilterContainsDN
	RegisterFuncMap["contains-valuetype"] = FilterContainsValueType
	RegisterFuncMap["contains-country"] = FilterContainsCountry
	RegisterFuncMap["contains-value"] = FilterValueContainsCase

	// startwith
	RegisterFuncMap["startswith-id"] = FilterStartsWithID
	RegisterFuncMap["startswith-type"] = FilterStartsWithType
	RegisterFuncMap["startswith-name"] = FilterStartsWithName
	RegisterFuncMap["startswith-vendor"] = FilterStartsWithVendor
	RegisterFuncMap["startswith-ip"] = FilterStartsWithIP
	RegisterFuncMap["startswith-dn"] = FilterStartsWithDN
	RegisterFuncMap["startswith-valuetype"] = FilterStartsWithValueType
	RegisterFuncMap["startswith-country"] = FilterStartsWithCountry
	RegisterFuncMap["startswith-value"] = FilterValueStartsWith

	RegisterGroupBy = make(registerGroupByFunc)
	RegisterGroupBy["id"] = GettersID
	RegisterGroupBy["value"] = GettersValue
	RegisterGroupBy["type"] = GettersType
	RegisterGroupBy["name"] = GettersName
	RegisterGroupBy["vendor"] = GettersVendor
	RegisterGroupBy["ip"] = GettersIP
	RegisterGroupBy["dn"] = GettersDN
	RegisterGroupBy["valuetype"] = GettersValueType
	RegisterGroupBy["country"] = GettersCountry

	RegisterGetters = make(registerGettersMap)
	RegisterGetters["id"] = GettersID
	RegisterGetters["value"] = GettersValue
	RegisterGetters["type"] = GettersType
	RegisterGetters["name"] = GettersName
	RegisterGetters["vendor"] = GettersVendor
	RegisterGetters["ip"] = GettersIP
	RegisterGetters["dn"] = GettersDN
	RegisterGetters["valuetype"] = GettersValueType
	RegisterGetters["country"] = GettersCountry

}

func sortBy(items Items, sortingL []string) (Items, []string) {
	sortFuncs := map[string]func(int, int) bool{
		"id":  func(i, j int) bool { return items[i].ID < items[j].ID },
		"-id": func(i, j int) bool { return items[i].ID > items[j].ID },

		"name":  func(i, j int) bool { return items[i].Name < items[j].Name },
		"-name": func(i, j int) bool { return items[i].Name > items[j].Name },

		"value":  func(i, j int) bool { return items[i].Value < items[j].Value },
		"-value": func(i, j int) bool { return items[i].Value > items[j].Value },

		"type":  func(i, j int) bool { return items[i].Type < items[j].Type },
		"-type": func(i, j int) bool { return items[i].Type > items[j].Type },

		"vendor":  func(i, j int) bool { return items[i].Vendor < items[j].Vendor },
		"-vendor": func(i, j int) bool { return items[i].Vendor > items[j].Vendor },

		"ip":  func(i, j int) bool { return items[i].IP < items[j].IP },
		"-ip": func(i, j int) bool { return items[i].IP > items[j].IP },

		"valueType":  func(i, j int) bool { return items[i].ValueType < items[j].ValueType },
		"-valueType": func(i, j int) bool { return items[i].ValueType > items[j].ValueType },

		"country":  func(i, j int) bool { return items[i].Country < items[j].Country },
		"-country": func(i, j int) bool { return items[i].Country > items[j].Country },
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
