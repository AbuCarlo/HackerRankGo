package miscellaneous

import (
    "testing"
)

func maxXor(arr []int32, queries []int32) []int32 {
    result := make([]int32, len(queries))
    for i, q := range queries {
        result[i] = 0
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
	arr00 := []int32{ 0, 1, 2 }
    queries00 := []int32{ 3, 7, 2 }
    result00 := maxXor(arr00, queries00)

    t.Logf("Result #02: %v", result00)

    arr02 := []int32{ 1, 3, 5, 7 }
    queries02 := []int32{ 17, 6 }
    result02 := maxXor(arr02, queries02)

    t.Logf("Result #02: %v", result02)
}