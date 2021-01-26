package main

import (
	"bytes"
	"errors"
)

func ParsePGArray(array string) ([]string, error) {
	var out []string
	var arrayOpened, quoteOpened, escapeOpened bool
	item := &bytes.Buffer{}
	for _, r := range array {
		switch {
		case !arrayOpened:
			if r != '{' {
				return nil, errors.New("Doesn't appear to be a postgres array.  Doesn't start with an opening curly brace.")
			}
			arrayOpened = true
		case escapeOpened:
			item.WriteRune(r)
			escapeOpened = false
		case quoteOpened:
			switch r {
			case '\\':
				escapeOpened = true
			case '"':
				quoteOpened = false
				if item.String() == "NULL" {
					item.Reset()
				}
			default:
				item.WriteRune(r)
			}
		case r == '}':
			// done
			out = append(out, item.String())
			return out, nil
		case r == '"':
			quoteOpened = true
		case r == ',':
			// end of item
			out = append(out, item.String())
			item.Reset()
		default:
			item.WriteRune(r)
		}
	}
	return nil, errors.New("Doesn't appear to be a postgres array.  Premature end of string.")
}
