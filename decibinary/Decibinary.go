package main

import (
	"fmt"
	"sort"
)

const MaximumIndex = 1e16

var numbers = map[int64][]int64{0: {0}, 1: {1}}
var counts = map[int64]int64{0: 1, 1: 1}

func countDecibinaryNumbers(x int64) int64 {
	if x % 2 == 1 {
		return countDecibinaryNumbers(x - 1)
	}
	result, ok := counts[x]
	if ok {
		return result
	}

	for least := x % 2; least < 10 && least <= x; least += 2 {
		most := (x - least) >> 1
		result += countDecibinaryNumbers(most)
	}
	counts[x] = result
	return result
}

func getDecibinaryNumbers(x int64) []int64 {
	result := numbers[x]
	if result != nil {
		return result
	}
	// An empty list will be the one we just created.
	for least := int64(x % 2); least < 10 && least <= x; least += 2 {
		most := (x - least) >> 1
		prefixes := getDecibinaryNumbers(most)
		for _, prefix := range prefixes {
			result = append(result, prefix*int64(10)+least)
		}
	}
	numbers[x] = result
	return result
}

// NthDecibinaryNumber yadda yadda.
func NthDecibinaryNumber(x int64) int64 {
	count := int64(0)
	var numerals []int64
	for m := int64(0); m < x; m++ {
		c := countDecibinaryNumbers(m)
		if count+c < x {
			count += c
		} else {
			numerals = getDecibinaryNumbers(m)
			break
		}
	}
	sort.Slice(numerals, func(i, j int) bool { return numerals[i] < numerals[j] })

	return numerals[x-count-1]
}


func main() {
	var sum int64
	for x := int64(0); sum < MaximumIndex; x++ {
		result := countDecibinaryNumbers(x)
		sum += result
		fmt.Printf("x = %d, count = %d, sum = %d\n", x, result, sum)
	}

	fmt.Printf("Stopped at sum %d\n", sum)
}
