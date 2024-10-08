package recursion

// https://www.hackerrank.com/challenges/recursive-digit-sum/problem

import (
	"fmt"
	"testing"
)

// SuperDigit really does not need recursion.
func SuperDigit(n string, k int32) int32 {
	var d int32 = 0
	for _, c := range n {
		d += int32(c) - int32('0')
		if d > 9 {
			d = d/10 + d%10
		}
	}
	var result int32 = 0
	for i := int32(0); i < k; i++ {
		result += d
		if result > 9 {
			result = result/10 + result%10
		}
	}
	return result
}

func TestSuperDigit(t *testing.T) {
	result := SuperDigit("148", 3)
	fmt.Printf("Result: %d\n", result)
}
