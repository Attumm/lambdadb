package main

import (
	"reflect"
	"sort"
	"testing"
)

func TestParseEnvStringToMap(t *testing.T) {
	testcases := []struct {
		input          string
		expected_key   string
		expected_value []string
	}{
		{"foo:bla", "foo", []string{"bla"}},
		{"foo:bar,bla", "foo", []string{"bar", "bla"}},
		{"foo:bar,bla;foo2:bar2,bla2", "foo", []string{"bar", "bla"}},
		{"foo:bar,bla;foo2:bar2,bla2", "foo2", []string{"bar2", "bla2"}},
	}

	for tcNumber, testcase := range testcases {
		mapResult := ParseLineToMap(testcase.input)
		result, found := mapResult[testcase.expected_key]
		if !found {
			t.Error("testcase", tcNumber, "expected", "found", "!=", "not found")
		}
		if !reflect.DeepEqual(result, testcase.expected_value) {
			t.Error("testcase", tcNumber, "expected", testcase.expected_value, "!=", result)
		}
	}
}

func TestParseEnvStringToFlatten(t *testing.T) {
	testcases := []struct {
		input    string
		expected []string
	}{
		{"foo:bla", []string{"bla"}},
		{"foo:bar,bla", []string{"bar", "bla"}},
		{"foo:bar,bla;foo2:bar2,bla2", []string{"bar", "bla", "bar2", "bla2"}},
		{"foo:bar,bla;foo2:bar2,bla2", []string{"bar", "bla", "bar2", "bla2"}},
	}

	for tcNumber, testcase := range testcases {
		mapResult := ParseLineToMap(testcase.input)
		result := FlattenMapStrSlice(mapResult)
		sort.Strings(result)
		sort.Strings(testcase.expected)
		if !reflect.DeepEqual(result, testcase.expected) {
			t.Error("testcase", tcNumber, "expected", testcase.expected, "!=", result)
		}
	}
}
