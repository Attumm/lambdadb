package main

import (
	"fmt"
	"net/http"
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

func (i Item) Columns() []string {
	return []string{
		"ID",
		"Value",
		"Type",
		"Name",
		"Vendor",
		"IP",
		"Dn",
		"ValueType",
		"Country",
	}
}

//"type", "vendor", "country"

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

// Group by functions

func GroupByName(i *Item) string {
	return i.Name
}

func GroupByIP(i *Item) string {
	return i.IP
}

func GroupByID(i *Item) string {
	return i.ID
}

func GroupByCountry(i *Item) string {
	return i.Country
}

type GroupedOperations struct {
	Funcs   registerFuncType
	GroupBy registerGroupByFunc
}

var Operations GroupedOperations

var RegisterFuncMap registerFuncType
var RegisterGroupBy registerGroupByFunc
var registerFormat registerFormatMap

func init() {

	RegisterFuncMap = make(registerFuncType)
	RegisterFuncMap["contains"] = FilterValueContains
	RegisterFuncMap["contains-case"] = FilterValueContainsCase
	RegisterFuncMap["match"] = FilterValueMatch
	RegisterFuncMap["starts-with"] = FilterValueStartsWith

	RegisterGroupBy = make(registerGroupByFunc)
	RegisterGroupBy["name"] = GroupByName
	RegisterGroupBy["ip"] = GroupByIP
	RegisterGroupBy["id"] = GroupByID
	RegisterGroupBy["country"] = GroupByCountry

	fmt.Println(RegisterFuncMap)
}

func sortBy(items Items, r *http.Request) Items {
	sortFuncs := map[string]func(int, int) bool{
		"id":     func(i, j int) bool { return items[i].ID < items[j].ID },
		"-id":    func(i, j int) bool { return items[i].ID > items[j].ID },
		"name":   func(i, j int) bool { return items[i].Name < items[j].Name },
		"-name":  func(i, j int) bool { return items[i].Name > items[j].Name },
		"value":  func(i, j int) bool { return items[i].Value < items[j].Value },
		"-value": func(i, j int) bool { return items[i].Value > items[j].Value },

		"type":  func(i, j int) bool { return items[i].Type < items[j].Type },
		"-type": func(i, j int) bool { return items[i].Type > items[j].Type },

		"vendor":  func(i, j int) bool { return items[i].Vendor < items[j].Vendor },
		"-vendor": func(i, j int) bool { return items[i].Vendor > items[j].Vendor },

		"ip":  func(i, j int) bool { return items[i].IP < items[j].IP },
		"-ip": func(i, j int) bool { return items[i].IP > items[j].IP },

		//"dn":  func(i, j int) bool { return items[i].Dn < items[j].Dn },
		//"-dn": func(i, j int) bool { return items[i].Dn > items[j].Dn },

		"valueType":  func(i, j int) bool { return items[i].ValueType < items[j].ValueType },
		"-valueType": func(i, j int) bool { return items[i].ValueType > items[j].ValueType },

		"country":  func(i, j int) bool { return items[i].Country < items[j].Country },
		"-country": func(i, j int) bool { return items[i].Country > items[j].Country },
	}
	sortingL, sortingGiven := r.URL.Query()["sorting"]
	if sortingGiven {
		for _, sortFuncName := range sortingL {
			sortFunc := sortFuncs[sortFuncName]
			sort.Slice(items, sortFunc)
		}
	}
	return items
}
