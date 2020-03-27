package hackerrank

import (
	"fmt"
	"testing"
	"github.com/abucarlo/hackerrank/dictionaries"
	"github.com/abucarlo/hackerrank/search"
)

func TestIceCream(t *testing.T) {
	var sample00 = []int32 { 1, 4, 5, 3, 2 }
	var x, y = search.FindPair(sample00, 4)
	println(x, y)
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
			ans := dictionaries.TwoStrings(test.s, test.t)
			if ans != test.common {
				t.Errorf("got %t, want %t", ans, test.common)
			}
		})
	}
}
