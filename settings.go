package main

import (
	"flag"
	"os"
	"strconv"
)

type Settings struct {
	msg        map[string]string
	VarString  map[string]string
	VarInt     map[string]int
	Parsers    map[string]func(string) string
	ParsersInt map[string]func(int) int
}

func (s *Settings) Set(flagName, defaultVar, message string) {
	s.msg[flagName] = message
	s.VarString[flagName] = defaultVar
}

func (s *Settings) SetString(flagName, defaultVar, message string) {
	s.Set(flagName, defaultVar, message)
}

func (s *Settings) SetInt(flagName string, defaultVar int, message string) {
	s.msg[flagName] = message
	s.VarInt[flagName] = defaultVar
}

func (s *Settings) SetParsed(flagName, defaultVar, message string, parserFunc func(string) string) {
	s.msg[flagName] = message
	s.VarString[flagName] = defaultVar
	s.Parsers[flagName] = parserFunc
}

func (s *Settings) SetParsedInt(flagName, defaultVar, message string, parserFunc func(int) int) {
	s.msg[flagName] = message
	s.VarString[flagName] = defaultVar
	s.ParsersInt[flagName] = parserFunc
}

func (s Settings) Get(flagName string) string {
	return s.VarString[flagName]
}

func (s Settings) GetInt(flagName string) int {
	return s.VarInt[flagName]
}

func (s *Settings) HandleCMDLineInput() {
	parsedString := make(map[string]*string)
	for key, val := range s.VarString {
		var newV = flag.String(key, val, s.msg[key])
		parsedString[key] = newV
	}
	parsedInt := make(map[string]*int)
	for key, val := range s.VarInt {
		var newV = flag.Int(key, val, s.msg[key])
		parsedInt[key] = newV
	}
	flag.Parse()

	for key, val := range parsedString {
		if parseFunc, found := s.Parsers[key]; found {
			s.VarString[key] = parseFunc(*val)
		} else {
			s.VarString[key] = *val
		}
	}
	for key, val := range parsedInt {
		if parseFunc, found := s.ParsersInt[key]; found {
			s.VarInt[key] = parseFunc(*val)
		} else {
			s.VarInt[key] = *val
		}
	}
}

func (s *Settings) HandleOSInput() {
	for key, _ := range s.VarString {
		varEnv, found := os.LookupEnv(key)
		if found {
			s.VarString[key] = varEnv
		}
	}
	for key, _ := range s.VarInt {
		varEnv, found := os.LookupEnv(key)
		if found {
			if num, err := strconv.Atoi(varEnv); err == nil {
				s.VarInt[key] = num
			}
		}
	}
}

func (s *Settings) Parse() {
	s.HandleOSInput()
	s.HandleCMDLineInput()
}

var SETTINGS = Settings{
	msg:        make(map[string]string),
	VarString:  make(map[string]string),
	VarInt:     make(map[string]int),
	Parsers:    make(map[string]func(string) string),
	ParsersInt: make(map[string]func(int) int),
}
