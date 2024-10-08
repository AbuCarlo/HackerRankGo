package dictionaries

import (
	"fmt"
	"testing"
)

// TwoStrings checks for the simplest common substring,
// i.e. one common character.
func TwoStrings(s string, t string) bool {
	d := map[rune]bool{}
	for _, c := range s {
		d[c] = true
	}
	for _, c := range t {
		if _, ok := d[c]; ok {
			return true
		}
	}
	return false
}

func TestTwoStrings(t *testing.T) {
	tests := []struct {
		s      string
		t      string
		common bool
	}{
		{"HELLO", "WORLD", true},
		{"hi", "world", false},
		{"TONY", "NASSAR", true},
	}

	for _, test := range tests {

		testname := fmt.Sprintf("%s,%s", test.s, test.t)
		t.Run(testname, func(t *testing.T) {
			ans := TwoStrings(test.s, test.t)
			if ans != test.common {
				t.Errorf("got %t, want %t", ans, test.common)
			}
		})
	}
}
