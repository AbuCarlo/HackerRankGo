package search

import "sort"

// FindPair returns the indices of the values
// adding up to "money"
func FindPair(cost []int32, money int32) (int, int) {
	// Index the values by original location.
	// WARN This will not deal with duplicate values!
	var index = map[int32]int{}
	for i, value := range cost {
		index[value] = i
	}

	// Sort descending, because of how binary search in Go works.
	sort.Slice(cost, func(i, j int) bool { return cost[i] > cost[j] })
	var high, low int32
	var max int32 = 0

	// TODO Start with money.
	for i, hi := range cost {
		var j = sort.Search(len(cost), func(j int) bool { return j > i && hi+cost[j] <= money })
		if j == len(cost) {
			continue
		}
		if hi+cost[j] > max {
			high = hi
			low = cost[j]
			max = high + low
		}
	}
	return index[low], index[high]
}

func whatFlavors(cost []int32, money int32) {

}
