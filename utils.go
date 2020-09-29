package main

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func runPrintMem() {
	for {
		PrintMemUsage()
		time.Sleep(1 * time.Second)
	}
}

// Parse int from string
// if parsed int or error return default value
func intMoreDefault(s string, defaultN int) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	if n < defaultN {
		return defaultN
	}
	return n
}

// add me
const ITEM_DELIMITER = ";"
const VAL_SEP = ","
const KEY_SEP = ":"

// should refactor to use error.
// Now I gotta go fast
func parseKeyValue(s string) (string, []string, bool) {
	items := strings.Split(s, KEY_SEP)
	if len(items) != 2 {
		return "", nil, true
	}
	values := strings.Split(items[1], VAL_SEP)
	return items[0], values, false
}

func ParseLineToMap(s string) map[string][]string {
	parsed := make(map[string][]string)
	items := strings.Split(s, ITEM_DELIMITER)
	for _, item := range items {
		key, values, err := parseKeyValue(item)
		if err {
			fmt.Println("Unable to parse line, discarded:", item)
			continue
		}
		parsed[key] = values
	}
	return parsed
}

func FlattenMapStrSlice(ss map[string][]string) []string {
	uniqItems := make(map[string]bool)
	for _, values := range ss {
		for _, val := range values {
			uniqItems[val] = true
		}
	}
	items := []string{}
	for item := range uniqItems {
		items = append(items, item)
	}
	return items
}
