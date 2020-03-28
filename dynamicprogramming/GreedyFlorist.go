package dynamicprogramming

import "math"

type history struct {
	// Numbers of purchases mapped to number of customers.
	// The map is initialized to { 0: k }. After the first purchase,
	// its value will be { 1: 1, 0: k - 1}.
	purchases map[int32]int32
	// Costs of the remaining flowers.
	costs []int32
	spent int32
}

func deleteFromSlice(a []int32, i int) []int32 {
	a[i] = a[len(a)-1] 
	return a[:len(a)-1]
}

// Optimize returns the minimum cost for k customers
// to buy all the flowers whose prices are given in c.
func Optimize(k int32, costs[]int32) int32 {
	initialState := history { map[int32]int32{ 0: k}, costs, 0}
	var optimum int32 = math.MaxInt32
	// For any given (multi)set of remaining flowers, what was the cheapest 
	// way to get there?
	states := []history{ initialState }
	for {
		nextStates := []history{}
		var localOptimum int32 = math.MaxInt32
		// What is the cheapest way to buy the next flower?
		for _, state := range states {
			var previousCost int32 = -1
			for i, cost := range state.costs {
				if (cost == previousCost) {
					continue
				}
				previousCost = cost
				for k, v := range state.purchases {
					thisCost := previousCost + cost * (1 + k)

				}
			}
		}
		states = nextStates
	}

	return optimum
}