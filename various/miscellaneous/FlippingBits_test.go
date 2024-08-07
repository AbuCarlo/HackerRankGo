package miscellaneous

import (
	"testing"
)

func flippingBits(n int64) int64 {
	u := uint64(n)
	u = -u - 1
	return int64(u)
}

func TestFlippingBits(t *testing.T) {
	inputs := []int64{2147483647, 1, 0}
	expected := []int64{2147483648, 4294967294, 4294967295}

	for i, n := range inputs {
		if actual := flippingBits(n); actual != expected[i] {
			t.Errorf("Test case %d, %d: expected %d, actual %d", i, n, expected[i], actual)
		}
	}
}
