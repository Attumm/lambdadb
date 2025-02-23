package main

import (
	"sort"
	"strconv"
	"strings"
)

// Item struct definition.
type Item struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	Score      string `json:"score"`
	Created_at string `json:"created_at"`
}

// Columns returns the column names.
func (i Item) Columns() []string {
	return []string{
		"id",
		"name",
		"status",
		"score",
		"created_at",
	}
}

// Row returns the values of the item as a string slice.
func (i Item) Row() []string {
	return []string{
		i.Id,
		i.Name,
		i.Status,
		i.Score,
		i.Created_at,
	}
}

// GetIndex returns the index for the item.
// Here, the first column is used as the index.
func (i Item) GetIndex() string {
	return i.Id
}

// ---- Getter Functions ----
func GettersId(i *Item) string {
	return i.Id
}
func GettersName(i *Item) string {
	return i.Name
}
func GettersStatus(i *Item) string {
	return i.Status
}
func GettersScore(i *Item) string {
	return i.Score
}
func GettersCreated_at(i *Item) string {
	return i.Created_at
}

// ---- Standard Filter Functions (case-sensitive) ----
func FilterIdContains(i *Item, s string) bool {
	return strings.Contains(i.Id, s)
}

func FilterIdStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Id, s)
}

func FilterIdMatch(i *Item, s string) bool {
	return i.Id == s
}
func FilterNameContains(i *Item, s string) bool {
	return strings.Contains(i.Name, s)
}

func FilterNameStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Name, s)
}

func FilterNameMatch(i *Item, s string) bool {
	return i.Name == s
}
func FilterStatusContains(i *Item, s string) bool {
	return strings.Contains(i.Status, s)
}

func FilterStatusStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Status, s)
}

func FilterStatusMatch(i *Item, s string) bool {
	return i.Status == s
}
func FilterScoreContains(i *Item, s string) bool {
	return strings.Contains(i.Score, s)
}

func FilterScoreStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Score, s)
}

func FilterScoreMatch(i *Item, s string) bool {
	return i.Score == s
}
func FilterCreated_atContains(i *Item, s string) bool {
	return strings.Contains(i.Created_at, s)
}

func FilterCreated_atStartsWith(i *Item, s string) bool {
	return strings.HasPrefix(i.Created_at, s)
}

func FilterCreated_atMatch(i *Item, s string) bool {
	return i.Created_at == s
}

// ---- Case-Insensitive Helper Functions ----
func IContains(s, substr string) bool {
	return strings.Contains(strings.ToLower(substr), strings.ToLower(s))
}

func IMatch(s, match string) bool {
	return strings.EqualFold(s, match)
}

func hasPrefixCaseInsensitive(s, prefix string) bool {
	if len(prefix) > len(s) {
		return false
	}
	return strings.EqualFold(s[:len(s)], prefix)
}

// ---- Case-Insensitive Filter Functions ----
func FilterIdIContains(i *Item, s string) bool {
	return strings.Contains(strings.ToLower(i.Id), strings.ToLower(s))
}

func FilterIdIMatch(i *Item, s string) bool {
	return IMatch(i.Id, s)
}

func FilterIdIPrefix(i *Item, s string) bool {
	return hasPrefixCaseInsensitive(i.Id, s)
}
func FilterNameIContains(i *Item, s string) bool {
	return strings.Contains(strings.ToLower(i.Name), strings.ToLower(s))
}

func FilterNameIMatch(i *Item, s string) bool {
	return IMatch(i.Name, s)
}

func FilterNameIPrefix(i *Item, s string) bool {
	return hasPrefixCaseInsensitive(i.Name, s)
}
func FilterStatusIContains(i *Item, s string) bool {
	return strings.Contains(strings.ToLower(i.Status), strings.ToLower(s))
}

func FilterStatusIMatch(i *Item, s string) bool {
	return IMatch(i.Status, s)
}

func FilterStatusIPrefix(i *Item, s string) bool {
	return hasPrefixCaseInsensitive(i.Status, s)
}
func FilterScoreIContains(i *Item, s string) bool {
	return strings.Contains(strings.ToLower(i.Score), strings.ToLower(s))
}

func FilterScoreIMatch(i *Item, s string) bool {
	return IMatch(i.Score, s)
}

func FilterScoreIPrefix(i *Item, s string) bool {
	return hasPrefixCaseInsensitive(i.Score, s)
}
func FilterCreated_atIContains(i *Item, s string) bool {
	return strings.Contains(strings.ToLower(i.Created_at), strings.ToLower(s))
}

func FilterCreated_atIMatch(i *Item, s string) bool {
	return IMatch(i.Created_at, s)
}

func FilterCreated_atIPrefix(i *Item, s string) bool {
	return hasPrefixCaseInsensitive(i.Created_at, s)
}

// ---- Reduce Functions ----
func reduceCount(items Items) map[string]string {
	result := make(map[string]string)
	result["count"] = strconv.Itoa(len(items))
	return result
}

// ---- Type Definitions and Registration ----
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

	// Register standard match filters
	RegisterFuncMap["match-id"] = FilterIdMatch
	RegisterFuncMap["match-name"] = FilterNameMatch
	RegisterFuncMap["match-status"] = FilterStatusMatch
	RegisterFuncMap["match-score"] = FilterScoreMatch
	RegisterFuncMap["match-created_at"] = FilterCreated_atMatch

	// Register standard contains filters
	RegisterFuncMap["contains-id"] = FilterIdContains
	RegisterFuncMap["contains-name"] = FilterNameContains
	RegisterFuncMap["contains-status"] = FilterStatusContains
	RegisterFuncMap["contains-score"] = FilterScoreContains
	RegisterFuncMap["contains-created_at"] = FilterCreated_atContains

	// Register standard startswith filters
	RegisterFuncMap["startswith-id"] = FilterIdStartsWith
	RegisterFuncMap["startswith-name"] = FilterNameStartsWith
	RegisterFuncMap["startswith-status"] = FilterStatusStartsWith
	RegisterFuncMap["startswith-score"] = FilterScoreStartsWith
	RegisterFuncMap["startswith-created_at"] = FilterCreated_atStartsWith

	// Register case-insensitive contains filters
	RegisterFuncMap["icontains-id"] = FilterIdIContains
	RegisterFuncMap["icontains-name"] = FilterNameIContains
	RegisterFuncMap["icontains-status"] = FilterStatusIContains
	RegisterFuncMap["icontains-score"] = FilterScoreIContains
	RegisterFuncMap["icontains-created_at"] = FilterCreated_atIContains

	// Register case-insensitive match filters
	RegisterFuncMap["imatch-id"] = FilterIdIMatch
	RegisterFuncMap["imatch-name"] = FilterNameIMatch
	RegisterFuncMap["imatch-status"] = FilterStatusIMatch
	RegisterFuncMap["imatch-score"] = FilterScoreIMatch
	RegisterFuncMap["imatch-created_at"] = FilterCreated_atIMatch

	// Register case-insensitive prefix filters
	RegisterFuncMap["iprefix-id"] = FilterIdIPrefix
	RegisterFuncMap["iprefix-name"] = FilterNameIPrefix
	RegisterFuncMap["iprefix-status"] = FilterStatusIPrefix
	RegisterFuncMap["iprefix-score"] = FilterScoreIPrefix
	RegisterFuncMap["iprefix-created_at"] = FilterCreated_atIPrefix

	// Register getters
	RegisterGetters["id"] = GettersId
	RegisterGetters["name"] = GettersName
	RegisterGetters["status"] = GettersStatus
	RegisterGetters["score"] = GettersScore
	RegisterGetters["created_at"] = GettersCreated_at

	// Register groupby functions
	RegisterGroupBy["id"] = GettersId
	RegisterGroupBy["name"] = GettersName
	RegisterGroupBy["status"] = GettersStatus
	RegisterGroupBy["score"] = GettersScore
	RegisterGroupBy["created_at"] = GettersCreated_at

	// Register reduce functions
	RegisterReduce["count"] = reduceCount
}

// ---- Sorting Functions ----
func sortBy(items Items, sortingL []string) (Items, []string) {
	sortFuncs := map[string]func(i, j int) bool{
		"id":          func(i, j int) bool { return items[i].Id < items[j].Id },
		"-id":         func(i, j int) bool { return items[i].Id > items[j].Id },
		"name":        func(i, j int) bool { return items[i].Name < items[j].Name },
		"-name":       func(i, j int) bool { return items[i].Name > items[j].Name },
		"status":      func(i, j int) bool { return items[i].Status < items[j].Status },
		"-status":     func(i, j int) bool { return items[i].Status > items[j].Status },
		"score":       func(i, j int) bool { return items[i].Score < items[j].Score },
		"-score":      func(i, j int) bool { return items[i].Score > items[j].Score },
		"created_at":  func(i, j int) bool { return items[i].Created_at < items[j].Created_at },
		"-created_at": func(i, j int) bool { return items[i].Created_at > items[j].Created_at },
	}

	for _, sortFuncName := range sortingL {
		sortFunc := sortFuncs[sortFuncName]
		sort.Slice(items, sortFunc)
	}

	// Gather available sort keys
	keys := []string{}
	for key := range sortFuncs {
		keys = append(keys, key)
	}
	return items, keys
}
