package main

import (
	"fmt"
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

// add me
