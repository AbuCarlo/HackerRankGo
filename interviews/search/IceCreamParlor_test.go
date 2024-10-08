package search

// https://www.hackerrank.com/challenges/ctci-ice-cream-parlor/problem

import (
	"sort"
	"testing"
)

// FindPair returns the indices of the values
// adding up to "money"
func FindPair(cost []int32, money int32) (int, int) {
	type element struct {
		value int32
		index int
	}
	// Index the values by original location.
	var indexed = []element{}

	for i, value := range cost {
		indexed = append(indexed, element{value, i})
	}

	// Sort descending, because of how binary search in Go works.
	sort.Slice(indexed, func(i int, j int) bool { return indexed[i].value > indexed[j].value })
	var high, low int
	var max int32 = 0

	// TODO Start with money.
	for i, hi := range indexed {
		j := sort.Search(len(indexed), func(j int) bool { return j > i && hi.value+indexed[j].value <= money })
		if j == len(indexed) {
			continue
		}
		lo := indexed[j]
		if hi.value+indexed[j].value > max {
			high = hi.index
			low = lo.index
			max = hi.value + lo.value
		}
	}
	if low < high {
		return low, high
	} // else
	return high, low
}

func TestIceCream(t *testing.T) {
	var sample00 = []int32{1, 4, 5, 3, 2}
	var x, y = FindPair(sample00, 4)
	println(x, y)
	var sample01 = []int32{2, 2, 4, 3}
	var l, r = FindPair(sample01, 4)
	println(l, r)
}
