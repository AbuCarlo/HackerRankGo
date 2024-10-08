package greedy

// https://www.hackerrank.com/challenges/greedy-florist/problem

import (
	"sort"
	"testing"
)

// Optimize returns the minimum cost for k customers
// to buy all the flowers whose prices are given in c.
func Optimize(k int32, costs[]int32) int32 {
	sort.Slice(costs, func(i int, j int) bool { return costs[i] > costs[j] })
	purchases := make([]int32, k)
	var result int32 = 0
	var j int32 = 0
	for _, cost := range costs {
		result += (purchases[j] + 1) * cost
		purchases[j]++
		j++
		if (j == k) {
			j = 0
		}
	}

	return result
}

func TestGreedyFlorist(t *testing.T) {
	result00 := Optimize(3, []int32{2, 5, 6})
	if result00 != 13 {
		t.Errorf("got %d, want %d", result00, 5)
	}
	result01 := Optimize(2, []int32{2, 5, 6})
	if result01 != 15 {
		t.Errorf("got %d, want %d", result00, 5)
	}
}