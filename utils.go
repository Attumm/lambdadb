package main

import (
	"fmt"
	"net/http"
	"runtime"
	"strconv"
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

func combineSlices(sss ...[]string) []string {
	set := make(map[string]bool)
	for _, ss := range sss {
		for _, s := range ss {
			set[s] = true
		}
	}
	l := []string{}
	for k, _ := range set {
		l = append(l, k)
	}
	return l
}

const JWT_WILDCARD string = "all"

func containsWildCard(ss []string) bool {
	for _, s := range ss {
		if s == JWT_WILDCARD {
			return true
		}
	}
	return false
}

func getColumnValues(groups []string, groupToValues map[string][]string) []string {
	values := []string{}
	for _, group := range groups {
		values = append(values, groupToValues[group]...)
	}
	return combineSlices(values)
}

func overrideAnyFilter(anyMap filterType, column string, columnValues []string) {
	key := fmt.Sprintf("match-%s", column)
	anyMap[key] = columnValues
}

func getJWT(r *http.Request, jwtSecret, headerName string) (Claims, error) {
	tokenString := r.Header.Get(headerName)
	return handleJWT(tokenString, jwtSecret)
}
