package main

import (
	"sort"
	"strconv"
	"strings"
)

type Item struct {
	Tconst         string `json:"tconst"`
	Titletype      string `json:"titletype"`
	Primarytitle   string `json:"primarytitle"`
	Originaltitle  string `json:"originaltitle"`
	Isadult        string `json:"isadult"`
	Startyear      string `json:"startyear"`
	Endyear        string `json:"endyear"`
	Runtimeminutes string `json:"runtimeminutes"`
	Genres         string `json:"genres"`
}

func (i Item) Columns() []string {
	return []string{
		"tconst",
		"titletype",
		"primarytitle",
		"originaltitle",
		"isadult",
		"startyear",
		"endyear",
		"runtimeminutes",
		"genres",
	}
}

func (i Item) Row() []string {
	return []string{
		i.Tconst,
		i.Titletype,
		i.Primarytitle,
		i.Originaltitle,
		i.Isadult,
		i.Startyear,
		i.Endyear,
		i.Runtimeminutes,
		i.Genres,
	}
}

func (i Item) GetIndex() string {
	return i.Tconst
}

// contain filters
func FilterTconstContains(i *Item, s string) bool {
	return strings.Contains(i.Tconst, s)
}
func FilterTitletypeContains(i *Item, s string) bool {
	return strings.Contains(i.Titletype, s)
}
func FilterPrimarytitleContains(i *Item, s string) bool {
	return strings.Contains(i.Primarytitle, s)
}
func FilterOriginaltitleContains(i *Item, s string) bool {
	return strings.Contains(i.Originaltitle, s)
}
func FilterIsadultContains(i *Item, s string) bool {
	return strings.Contains(i.Isadult, s)
}
func FilterStartyearContains(i *Item, s string) bool {
	return strings.Contains(i.Startyear, s)
}
func FilterEndyearContains(i *Item, s string) bool {
	return strings.Contains(i.Endyear, s)
}
func FilterRuntimeminutesContains(i *Item, s string) bool {
	return strings.Contains(i.Runtimeminutes, s)
}
func FilterGenresContains(i *Item, s string) bool {
	return strings.Contains(i.Genres, s)
}

// startswith filters
func FilterTconstStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Tconst, s)
}
func FilterTitletypeStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Titletype, s)
}
func FilterPrimarytitleStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Primarytitle, s)
}
func FilterOriginaltitleStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Originaltitle, s)
}
func FilterIsadultStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Isadult, s)
}
func FilterStartyearStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Startyear, s)
}
func FilterEndyearStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Endyear, s)
}
func FilterRuntimeminutesStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Runtimeminutes, s)
}
func FilterGenresStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Genres, s)
}

// match filters
func FilterTconstMatch(i *Item, s string) bool {
	return i.Tconst == s
}
func FilterTitletypeMatch(i *Item, s string) bool {
	return i.Titletype == s
}
func FilterPrimarytitleMatch(i *Item, s string) bool {
	return i.Primarytitle == s
}
func FilterOriginaltitleMatch(i *Item, s string) bool {
	return i.Originaltitle == s
}
func FilterIsadultMatch(i *Item, s string) bool {
	return i.Isadult == s
}
func FilterStartyearMatch(i *Item, s string) bool {
	return i.Startyear == s
}
func FilterEndyearMatch(i *Item, s string) bool {
	return i.Endyear == s
}
func FilterRuntimeminutesMatch(i *Item, s string) bool {
	return i.Runtimeminutes == s
}
func FilterGenresMatch(i *Item, s string) bool {
	return i.Genres == s
}

// reduce functions

func reduceCount(items Items) map[string]string {
	result := make(map[string]string)
	result["count"] = strconv.Itoa(len(items))
	return result
}

// getters
func GettersTconst(i *Item) string {
	return i.Tconst
}
func GettersTitletype(i *Item) string {
	return i.Titletype
}
func GettersPrimarytitle(i *Item) string {
	return i.Primarytitle
}
func GettersOriginaltitle(i *Item) string {
	return i.Originaltitle
}
func GettersIsadult(i *Item) string {
	return i.Isadult
}
func GettersStartyear(i *Item) string {
	return i.Startyear
}
func GettersEndyear(i *Item) string {
	return i.Endyear
}
func GettersRuntimeminutes(i *Item) string {
	return i.Runtimeminutes
}
func GettersGenres(i *Item) string {
	return i.Genres
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

	RegisterFuncMap["match-tconst"] = FilterTconstMatch
	RegisterFuncMap["match-titletype"] = FilterTitletypeMatch
	RegisterFuncMap["match-primarytitle"] = FilterPrimarytitleMatch
	RegisterFuncMap["match-originaltitle"] = FilterOriginaltitleMatch
	RegisterFuncMap["match-isadult"] = FilterIsadultMatch
	RegisterFuncMap["match-startyear"] = FilterStartyearMatch
	RegisterFuncMap["match-endyear"] = FilterEndyearMatch
	RegisterFuncMap["match-runtimeminutes"] = FilterRuntimeminutesMatch
	RegisterFuncMap["match-genres"] = FilterGenresMatch

	// register contains filters
	RegisterFuncMap["contains-tconst"] = FilterTconstContains
	RegisterFuncMap["contains-titletype"] = FilterTitletypeContains
	RegisterFuncMap["contains-primarytitle"] = FilterPrimarytitleContains
	RegisterFuncMap["contains-originaltitle"] = FilterOriginaltitleContains
	RegisterFuncMap["contains-isadult"] = FilterIsadultContains
	RegisterFuncMap["contains-startyear"] = FilterStartyearContains
	RegisterFuncMap["contains-endyear"] = FilterEndyearContains
	RegisterFuncMap["contains-runtimeminutes"] = FilterRuntimeminutesContains
	RegisterFuncMap["contains-genres"] = FilterGenresContains

	// register startswith filters
	RegisterFuncMap["startswith-tconst"] = FilterTconstStartsWith
	RegisterFuncMap["startswith-titletype"] = FilterTitletypeStartsWith
	RegisterFuncMap["startswith-primarytitle"] = FilterPrimarytitleStartsWith
	RegisterFuncMap["startswith-originaltitle"] = FilterOriginaltitleStartsWith
	RegisterFuncMap["startswith-isadult"] = FilterIsadultStartsWith
	RegisterFuncMap["startswith-startyear"] = FilterStartyearStartsWith
	RegisterFuncMap["startswith-endyear"] = FilterEndyearStartsWith
	RegisterFuncMap["startswith-runtimeminutes"] = FilterRuntimeminutesStartsWith
	RegisterFuncMap["startswith-genres"] = FilterGenresStartsWith

	// register getters
	RegisterGetters["tconst"] = GettersTconst
	RegisterGetters["titletype"] = GettersTitletype
	RegisterGetters["primarytitle"] = GettersPrimarytitle
	RegisterGetters["originaltitle"] = GettersOriginaltitle
	RegisterGetters["isadult"] = GettersIsadult
	RegisterGetters["startyear"] = GettersStartyear
	RegisterGetters["endyear"] = GettersEndyear
	RegisterGetters["runtimeminutes"] = GettersRuntimeminutes
	RegisterGetters["genres"] = GettersGenres

	// register groupby
	RegisterGroupBy["tconst"] = GettersTconst
	RegisterGroupBy["titletype"] = GettersTitletype
	RegisterGroupBy["primarytitle"] = GettersPrimarytitle
	RegisterGroupBy["originaltitle"] = GettersOriginaltitle
	RegisterGroupBy["isadult"] = GettersIsadult
	RegisterGroupBy["startyear"] = GettersStartyear
	RegisterGroupBy["endyear"] = GettersEndyear
	RegisterGroupBy["runtimeminutes"] = GettersRuntimeminutes
	RegisterGroupBy["genres"] = GettersGenres

	// register reduce functions
	RegisterReduce["count"] = reduceCount
}
func sortBy(items Items, sortingL []string) (Items, []string) {
	sortFuncs := map[string]func(int, int) bool{"tconst": func(i, j int) bool { return items[i].Tconst < items[j].Tconst },
		"-tconst": func(i, j int) bool { return items[i].Tconst > items[j].Tconst },

		"titletype":  func(i, j int) bool { return items[i].Titletype < items[j].Titletype },
		"-titletype": func(i, j int) bool { return items[i].Titletype > items[j].Titletype },

		"primarytitle":  func(i, j int) bool { return items[i].Primarytitle < items[j].Primarytitle },
		"-primarytitle": func(i, j int) bool { return items[i].Primarytitle > items[j].Primarytitle },

		"originaltitle":  func(i, j int) bool { return items[i].Originaltitle < items[j].Originaltitle },
		"-originaltitle": func(i, j int) bool { return items[i].Originaltitle > items[j].Originaltitle },

		"isadult":  func(i, j int) bool { return items[i].Isadult < items[j].Isadult },
		"-isadult": func(i, j int) bool { return items[i].Isadult > items[j].Isadult },

		"startyear":  func(i, j int) bool { return items[i].Startyear < items[j].Startyear },
		"-startyear": func(i, j int) bool { return items[i].Startyear > items[j].Startyear },

		"endyear":  func(i, j int) bool { return items[i].Endyear < items[j].Endyear },
		"-endyear": func(i, j int) bool { return items[i].Endyear > items[j].Endyear },

		"runtimeminutes":  func(i, j int) bool { return items[i].Runtimeminutes < items[j].Runtimeminutes },
		"-runtimeminutes": func(i, j int) bool { return items[i].Runtimeminutes > items[j].Runtimeminutes },

		"genres":  func(i, j int) bool { return items[i].Genres < items[j].Genres },
		"-genres": func(i, j int) bool { return items[i].Genres > items[j].Genres },
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
