package dynamicprogramming

import "sort"

var numbers = map[int][]int{0: {0}, 1: {1}}
var counts = map[int]int{0: 1, 1: 1}

func countDecibinaryNumbers(x int) int {
	result := counts[x]
	if result > 0 {
		return result
	}
	for least := x % 2; least < 10 && least <= x; least += 2 {
		most := (x - least) >> 1
		result += countDecibinaryNumbers(most)
	}
	counts[x] = result
	return result
}

func getDecibinaryNumbers(x int) []int {
	result := numbers[x]
	if result != nil {
		return result
	}
	// An empty list will be the one we just created.
	for least := x % 2; least < 10 && least <= x; least += 2 {
		most := (x - least) >> 1
		prefixes := getDecibinaryNumbers(most)
		for _, prefix := range prefixes {
			result = append(result, prefix*10+least)
		}
	}
	return result
}

// NthDecibinaryNumber yadda yadda.
func NthDecibinaryNumber(n int) int {
	count := 0
	var numerals []int
	for m := 0; m < n; m++ {
		c := countDecibinaryNumbers(m)
		if count+c < n {
			count += c
		} else {
			numerals = getDecibinaryNumbers(m)
			break
		}
	}
	sort.Ints(numerals)
	return numerals[n-count-1]
}

func decibinaryNumbers(x int64) int {
	return NthDecibinaryNumber(int(x))
}
