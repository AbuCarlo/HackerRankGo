package miscellaneous

import (
	"slices"
	"testing"

	"golang.org/x/exp/constraints"
)

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
        arr []int
        queries []int
        expected []int
    }{
        {
            []int{ 0, 1, 2 },
            []int{3, 7, 2},
            []int{ 3, 7, 3 },
        },
        {
            []int{ 5, 1, 7, 4, 3 },
            []int{2, 0},
            []int{ 7, 7 }, 
        },
        {
            []int{1, 3, 5, 7 },
            []int{17, 6},
            []int{22, 7 },
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