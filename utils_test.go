package main

import (
	"reflect"
	"sort"
	"testing"
)

func TestCombineSlices(t *testing.T) {
	testcases := []struct {
		input    [][]string
		expected []string
	}{
		{[][]string{[]string{"1", "2"}, []string{"2", "3"}}, []string{"1", "2", "3"}},
		{[][]string{[]string{}, []string{"2", "3"}}, []string{"2", "3"}},
		{[][]string{[]string{"1", "2"}, []string{}}, []string{"1", "2"}},
		{[][]string{[]string{"1"}, []string{"2"}, []string{"3"}}, []string{"1", "2", "3"}},
		{[][]string{}, []string{}},
	}

	for tcNumber, testcase := range testcases {
		result := combineSlices(testcase.input...)
		sort.Strings(result)
		sort.Strings(testcase.expected)
		if !reflect.DeepEqual(result, testcase.expected) {
			t.Error("testcase", tcNumber, "expected", testcase.expected, "!=", result)
		}
	}

}
