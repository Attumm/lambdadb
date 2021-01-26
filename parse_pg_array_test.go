package main

import (
	"reflect"
	"testing"
)

func TestParseArray(t *testing.T) {
	scanTests := []struct {
		in  string
		out []string
	}{
		{"{one,two}", []string{"one", "two"}},
		{`{"one, sdf",two}`, []string{"one, sdf", "two"}},
		{`{"\"one\"",two}`, []string{`"one"`, "two"}},
		{`{"\\one\\",two}`, []string{`\one\`, "two"}},
		{`{"{one}",two}`, []string{`{one}`, "two"}},
		{`{"one two"}`, []string{`one two`}},
		{`{"one,two"}`, []string{`one,two`}},
		{`{abcdef:83bf98cc-fec9-4e77-b4cf-99f9fb6655fa-0NH:zxcvzxc:wers:vxdfw-asdf-asdf}`, []string{"abcdef:83bf98cc-fec9-4e77-b4cf-99f9fb6655fa-0NH:zxcvzxc:wers:vxdfw-asdf-asdf"}},
		{`{"",two}`, []string{"", "two"}},
		{`{" ","NULL"}`, []string{" ", ""}},
	}

	for tcNumber, testcase := range scanTests {
		result, err := ParsePGArray(testcase.in)
		if err != nil {
			t.Error("testcase", tcNumber, "gave error")
		}
		if len(result) == 0 {
			t.Error("testcase", tcNumber, "expected", "found", "!=", "not found")
		}
		if !reflect.DeepEqual(result, testcase.out) {
			t.Error("testcase", tcNumber, "expected", testcase.out, "!=", result)
		}
	}
}
