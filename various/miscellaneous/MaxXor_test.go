package miscellaneous

import (
	"reflect"
	"slices"
	"testing"

	"golang.org/x/exp/constraints"
)

func maxXorArray[N constraints.Integer](a []N) N {
	max := N(0)
	mask := N(0)

	sizeOfN := reflect.TypeOf(N(0)).Bits()
    // Start with the highest-order bit, and end with 0b1.
	for position := sizeOfN - 1; position >= 0; position -= 1 {
		bit := N(1) << position
        // Add the next-highest-order bit to the mask.
		mask |= bit

		set := map[N]bool{}
        // Find all possible prefixes with respect the the current mask.
		for _, num := range a {
			left := num & mask
			set[left] = true
		}
        // Try to find a value better than the current maximum,
        // i.e. with one more bit. Any value with a 0 bit in the 
        // current position is, as far as we know, no better than
        // max, since we can't peek ahead at bits to the right.
		greed := max | bit

		for prefix := range set {
			if set[greed^prefix] {
				max = greed
				break
			}
		}
	}
	return max
}

func TestMaxXorArray(t *testing.T) {
	// https://stackoverflow.com/a/66822115/476942
	sample := []int{3, 10, 5, 25, 2, 8}
	actual := maxXorArray(sample)
	if actual != 28 {
		t.Errorf("Expected %d for %v; actual %v", 28, sample, actual)

	}
}

func maxXor[N constraints.Integer](arr []N, queries []N) []N {
	result := make([]N, len(queries))
	for i, q := range queries {
		result[i] = N(0)
		for _, a := range arr {
			x := q ^ a
			if x > result[i] {
				result[i] = x
			}
		}
	}

	return result
}

func TestSamples(t *testing.T) {
	tests := []struct {
		arr      []int
		queries  []int
		expected []int
	}{
		{
			[]int{0, 1, 2},
			[]int{3, 7, 2},
			[]int{3, 7, 3},
		},
		{
			[]int{5, 1, 7, 4, 3},
			[]int{2, 0},
			[]int{7, 7},
		},
		{
			[]int{1, 3, 5, 7},
			[]int{17, 6},
			[]int{22, 7},
		},
	}

	for i, test := range tests {
		actual := maxXor(test.arr, test.queries)
		if !slices.Equal(actual, test.expected) {
			t.Errorf("Test # %d expected %v; actual %v", i, test.expected, actual)
		} else {
			t.Logf("Result %02d: %v", i, actual)
		}
	}
}
