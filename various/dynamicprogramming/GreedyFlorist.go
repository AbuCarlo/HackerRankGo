package dynamicprogramming

import "sort"

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