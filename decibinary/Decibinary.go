package main

import (
	"fmt"
	"slices"
	"sort"
)

// See https://www.hackerrank.com/challenges/decibinary-numbers/problem

const MaximumIndex = 1e16

// We know from experience that this is the largest decimal
// number allowed by the problem definition, i.e. the total number of
// decibinary numerals needed past this point will exceed
// MaximumIndex.
const MaximumDecimalNumber = 285112

var countsBySize = make(map[int]map[int]int)

func countNumerals() []int64 {
	// Allow one more element for the odd number 285113
	a := make([]int64, MaximumDecimalNumber+2)
	a[0] = 1
	a[1] = 1

	countsBySize[0] = map[int]int{ 0: 1}
	countsBySize[1] = map[int]int{ 1: 1}

	for n := 2; n <= MaximumDecimalNumber; n += 2 {
		countsBySize[n] = make(map[int]int)
		var count int64
		// Populate the least-significant "decibinary" digit. How 
		// many decibinary numerals correspond to the remaining 
		// value? 
		for least := n % 2; least < 10 && least <= n; least += 2 {
			most := (n - least) >> 1
			count += a[most]

			for prefixSize, prefixCount := range countsBySize[most] {
				countsBySize[n][prefixSize+1] += prefixCount
			}
		}
		a[n] = count
		// For any even decimal number, the final digit will be no more than
		// 8. It's always possible to add 1 to final digit, giving you the next
		// decimal integer, which is perforce odd. For any digit, we can subtract
		// 2, halve it, and shift it leftward, to the next-higher order digit,
		// but an odd digit cannot be reduced below 1.
		a[n+1] = count
	}

	return a
}

// The numbers do start at 0 (see the problem definition).
var counts = countNumerals()

// This array allows a binary search for a decimal number based on the rank
// of the decibinary numeral. In other words, if we want the xth decibinary
// numeral, we can look for the lowest value in this array > than x. Its index
// will be the decimal number 1 greater than the one we want.
var partialSums = func() []int64 {
	var sum int64
	a := make([]int64, len(counts))
	for i, c := range counts {
		sum += c
		a[i] = sum
	}
	return a
}()

// Taken from https://stackoverflow.com/a/11398748/476942

var _ uint64

// Using []int8 does not seem to affect the speed.
var _tab = []int{
	0, 9, 1, 10, 13, 21, 2, 29,
	11, 14, 16, 18, 22, 25, 3, 30,
	8, 12, 20, 28, 15, 17, 24, 7,
	19, 27, 23, 6, 26, 5, 4, 31}

// The value must be unsigned in order to be right-shifted properly.
func log2_32(value uint32) int {
	value |= value >> 1
	value |= value >> 2
	value |= value >> 4
	value |= value >> 8
	value |= value >> 16
	return int(_tab[(value*0x07C4ACDD)>>27])
}

func maximumDecibinaryDigits(n int) int {
	return log2_32(uint32(n)) + 1
}

func minimumDecibinaryDigits(n int) int {
	// Does 0 require 1 digit or 0?
	result := 0
	for n > 0 {
		if n%2 == 1 {
			n--
		}
		n -= min(n, 8)
		n /= 2
		result++
	}

	return result
}

func decibinaryToArray(d int64) []int {
	if d == 0 {
		return []int{0}
	}
	// It's much faster to request an initial capacity.
	// The reallocation, if this number is 16, is expensive,
	// whereas an allocation of 32 is cheap. Since most of
	// inputs are large, let's just go with 32.
	a := make([]int, 0, 32)
	for ; d > 0; d /= 10 {
		digit := int(d % 10)
		a = append(a, digit)
	}
	slices.Reverse(a)
	return a
}

func decibinaryArrayToDecimal(d []int) int {
	result := 0
	for _, digit := range d {
		result *= 2
		result += digit
	}
	return result
}

// TODO This is poorly named.
func decibinaryToBinary(d int64) int {
	significance := 1
	result := 0
	for d > 0 {
		result += significance * int(d%10)
		d /= 10
		significance *= 2
	}
	return result
}

func lowestDecibinaryNumeral(n int) int64 {
	result := int64(0)
	place := int64(1)
	// Go from lower-order to higher.
	for n > 0 {
		var digit int
		if n < 10 {
			digit = n
		} else if n%2 == 0 {
			digit = 8
		} else {
			digit = 9
		}
		result += int64(digit) * place
		n -= digit
		n /= 2
		place *= 10
	}
	return result
}

func highestDecibinaryNumeral(n int) int64 {
	result := int64(0)
	if n == 0 {
		return result
	}
	bit := 1 << 30
	for n&bit == 0 {
		bit >>= 1
	}
	for bit > 0 {
		if n&bit != 0 {
			result = result*10 + 1
		} else {
			result = result * 10
		}
		bit >>= 1
	}
	return result
}

func locate(rank int64) int64 {
	if rank < 1 {
		panic("Queries are 1-based.")
	}
	result := make([]int, 0, 16)
	nativeValue := rankToNative(rank)
	// How many numerals do we want to skip over, proceeding backward?
	target := partialSums[nativeValue] - rank
	countForPrefix := int64(0)
	highest := highestDecibinaryNumeral(nativeValue)
	// These array elements go from *most* significant to least.
	suffix := decibinaryToArray(highest)
	for len(suffix) > 1 {
		if suffix[0] == 0 {
			// All the digits to the right are 1 or 0. So none
			// of them could have been transposed leftward.
			// Thus the current digit doesn't affect the
			// ceiling.
			// TODO Improve this comment.
			result = append(result, 0)
			suffix = suffix[1:]
			continue
		}
		decimalValueOfSuffix := decibinaryArrayToDecimal(suffix[1:])
		countForSuffix := counts[decimalValueOfSuffix]
		// If there are enough numerals with this prefix, we leave
		// it, and just reduce the suffix further.
		if target < countForSuffix+countForPrefix {
			result = append(result, suffix[0])
			suffix = suffix[1:]
			countForPrefix += countForSuffix
		} else {
			// Shift a value rightward. On the next iteration,
			// the suffix will have more bits to permute.
			suffix[0]--
			suffix[1] += 2
			// This step will probably be taken care of above.
			if suffix[0] == 0 {
				suffix = suffix[1:]
			}
			// countForPrefix += countForSuffix
		}
	}
	result = append(result, suffix[0])
	var blah int64 = 0
	for _, digit := range result {
		blah *= 10
		blah += int64(digit)
	}

	// Actually, we want to return a *printable* decibinary number.
	// We could just return a string, I suppose.
	return blah
}

// The lower bound in the array of partial sums should be the rank of the minimal
// decibinary representation of d.
func rankToNative(rank int64) int {
	at := sort.Search(len(partialSums), func(ix int) bool { return rank <= partialSums[ix] })
	return at
}

func main() {
	var lastQuery int64 = 0
	for i := 0; i < 10; i++ {
		result := []int64{}
		for j := lastQuery + 1; j <= lastQuery + counts[i]; j++ {
			result = append(result, locate(j))
		}
		fmt.Printf("%d: %v\n", i, result)
		lastQuery += counts[i]
	}
}
