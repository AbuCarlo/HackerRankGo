package main

import (
	"fmt"
)

// See https://www.hackerrank.com/challenges/decibinary-numbers/problem

const MaximumIndex = 1e16
// We know from experience that this is the largest decimal
// number allowed by the problem definition, i.e. the total number of
// decibinary numerals needed past this point will exceed 
// MaximumIndex.
const MaximumDecimalNumber = 40 // 285112

// The numbers do start at 0 (see the problem definition).
var counts = func() []int64 {
	// Allow one more element for the odd number 285113
	a := make([]int64, MaximumDecimalNumber + 2)
	a[0] = 1
	a[1] = 1

	
	for n := 2; n <= MaximumDecimalNumber; n += 2 {
		var count int64

		for least := n % 2; least < 10 && least <= n; least += 2 {
			most := (n - least) >> 1
			count += a[most]
		}
		a[n] = count
		// For any even decimal number, the final digit will be no more than
		// 8. It's always possible to add 1 to final digit, giving you the next
		// decimal integer, which is perforce odd. Any multiple of 2 can be 
		// halved and shifted leftward, but this additional 1 cannot. 
		a[n + 1] = count
	}

	return a
}()

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

var _tab = []int {
	0,  9,  1, 10, 13, 21,  2, 29,
   11, 14, 16, 18, 22, 25,  3, 30,
	8, 12, 20, 28, 15, 17, 24,  7,
   19, 27, 23,  6, 26,  5,  4, 31};

func log2_32(value uint32) int {
   value |= value >> 1;
   value |= value >> 2;
   value |= value >> 4;
   value |= value >> 8;
   value |= value >> 16;
   return _tab[(value*0x07C4ACDD) >> 27];
}

func maximumDecibinaryDigits(n int) int {
	return log2_32(uint32(n)) + 1
}

func minimumDecibinaryDigits(n int) int {
	// Does 0 require 1 digit or 0?
	result := 0
	for n > 0 {
		if n % 2 == 1 {
			n--
		}
		n -= min(n, 8)
		n /= 2
		result++
	}

	return result
}

func main() {
	for i, c := range counts {
		fmt.Printf("n = %d, maximum size = %d, minimum size = %d, count = %d, sum = %d\n", i, maximumDecibinaryDigits(i), minimumDecibinaryDigits(i), c, partialSums[i]);
	}
}
