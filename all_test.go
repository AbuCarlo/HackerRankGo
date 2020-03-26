package hackerrank

import (
	"fmt"
	"testing"
	"github.com/abucarlo/hackerrank/dictionaries"
)

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
			ans := dictionaries.TwoStrings(test.s, test.t)
			if ans != test.common {
				t.Errorf("got %t, want %t", ans, test.common)
			}
		})
	}
}
