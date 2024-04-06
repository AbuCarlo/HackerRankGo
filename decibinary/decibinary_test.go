package main

import (
	"fmt"
	"testing"
)

func printDecibinaryNumeral(n []int)  string {
	s := ""
	for _, d := range n {
		s = fmt.Sprint(d) + s
	}
	return s
}

func TestEnumeration(t *testing.T) {
	n := 10
	numeral := lowestDecibinaryNumeral(n)
	t.Logf("10 => %s", printDecibinaryNumeral(numeral))
	for i := 1; i < int(counts[n]); i++ {
		// Optional: set capacity to longest numeral?
		copy := append(make([]int, 0, len(numeral)), numeral...)
		carries := i
		// Transformations...
		for j := 0; carries > 0; j++ {
			if copy[j] == 0 || copy[j] == 1 {
				continue
			}
			// You *have* to carry these.
			if copy[j] > 9 {
				// You can only carry forward an even number.
				reduction := 8
				if copy[j] % 2 == 0 {
					reduction = 9
				}
				remainder := copy[j] - reduction
				copy[j] = reduction
				if j + 1 == len(copy) {
					copy = append(copy, remainder)
				} else {
					copy[j + 1] += remainder >> 1
				}
			}
			// *Now* perform additional transformations.
			k := min(carries, copy[j] / 2)
			copy[j] -= k * 2
			if len(copy) == j + 1 {
				copy = append(copy, k)
			} else {
				copy[j + 1] += k
			}
			carries -= k
		}
		t.Logf("10 => %s", printDecibinaryNumeral(copy))
	}
}