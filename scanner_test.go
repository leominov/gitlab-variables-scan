package main

import "testing"

func TestIsVariableContainsSensitiveData(t *testing.T) {
	s := &Scanner{
		c: &Config{
			PairsRaw:  []string{"FOO=BAR"},
			ValuesRaw: []string{"[0-9]"},
			KeysRaw:   []string{"TEST"},
		},
	}
	s.c.parseRawData()
	tests := []struct {
		v     *Variable
		match bool
	}{
		{
			v: &Variable{
				Key:   "FOO",
				Value: "FOO",
			},
			match: false,
		},
		{
			v: &Variable{
				Key:   "FOO",
				Value: "BAR",
			},
			match: true,
		},
		{
			v: &Variable{
				Key:   "TEST",
				Value: "ABC",
			},
			match: true,
		},
		{
			v: &Variable{
				Key:   "TEST",
				Value: "DEF",
			},
			match: true,
		},
		{
			v: &Variable{
				Key:   "FOO",
				Value: "1234567890",
			},
			match: true,
		},
	}
	for id, test := range tests {
		_, contains := s.IsVariableContainsSensitiveData(test.v)
		if test.match != contains {
			t.Errorf("%d. Must be %v, but got %v", id, test.match, contains)
		}
	}
}
