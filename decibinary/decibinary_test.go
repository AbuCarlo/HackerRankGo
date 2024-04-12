package main

import (
	"bufio"
	"os"
	"sort"
	"strconv"
	"testing"
)

func printDecibinaryNumeral(n []int)  string {
	s := ""
	for _, d := range n {
		s = strconv.Itoa(d) + s
	}
	return s
}

/*
func XTestEnumeration(t *testing.T) {
	n := 10
	numeral := lowestDecibinaryNumeral(n)
	t.Logf("10 => %s", printDecibinaryNumeral(numeral))
	for i := 1; i < int(counts[n]); i++ {
		// Optional: set capacity to longest numeral?
		copy := append(make([]int, 0, len(numeral)), numeral...)
		carries := i
		// Transformations...
		for j := 0; carries > 0; j++ {
			if copy[j] == 0 || copy[j] == 1 {
				continue
			}
			// You *have* to carry these.
			if copy[j] > 9 {
				// You can only carry forward an even number.
				reduction := 8
				if copy[j] % 2 == 0 {
					reduction = 9
				}
				remainder := copy[j] - reduction
				copy[j] = reduction
				if j + 1 == len(copy) {
					copy = append(copy, remainder)
				} else {
					copy[j + 1] += remainder >> 1
				}
			}
			// *Now* perform additional transformations.
			k := min(carries, copy[j] / 2)
			copy[j] -= k * 2
			if len(copy) == j + 1 {
				copy = append(copy, k)
			} else {
				copy[j + 1] += k
			}
			carries -= k
		}
		t.Logf("10 => %s", printDecibinaryNumeral(copy))
	}
}

*/

func readInt64(scanner *bufio.Scanner) []int64 {
	inputs := make([]int64, 0)
    for scanner.Scan() {
		i, _ := strconv.ParseInt(scanner.Text(), 10, 64)
        inputs = append(inputs, i)
    }
	return inputs
}

func FuzzTranslation(f *testing.F) {
	f.Fuzz(func(t *testing.T, input int) {
		if input < 1 || input > MaximumDecimalNumber {
			t.Skip()
		}
        highest := highestDecibinaryNumeral(int64(input))
		roundTrip := decibinaryToBinary(highest)
		if (roundTrip != input) {
			t.Errorf("Input %d has maximum decibinary numeral %d; this round-tripped as %d", input, highest, roundTrip)
		}

		lowest := highestDecibinaryNumeral(int64(input))
		roundTrip = decibinaryToBinary(lowest)
		if (roundTrip != input) {
			t.Errorf("Input %d has minimum decibinary numeral %d; this round-tripped as %d", input, lowest, roundTrip)
		}
    })
}
 
func TestBoundaries(t *testing.T) {
	inputFile, err := os.Open("decibinary-input07.txt")
    if err != nil {
        t.Fatal(err)
    }
    defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	scanner.Scan()
	size, _ := strconv.ParseInt(scanner.Text(), 10, 32)
	inputs := readInt64(scanner)
	if (size != int64(len(inputs))) {
		t.Fatalf("Expected %d inputs; got %d", size, len(inputs))
	}

	outputFile, err := os.Open("decibinary-output07.txt")
    if err != nil {
        t.Fatal(err)
    }
	defer outputFile.Close()
	outputs := readInt64(bufio.NewScanner(outputFile))
	if (size != int64(len(outputs))) {
		t.Fatalf("Expected %d inputs; got %d", size, len(outputs))
	}

	for i := 0; i < int(size); i++ {
		// This is the rank of a decibinary number: the "query".
		rank := inputs[i]
		// This is the decibinary number of that rank.
		output := outputs[i]
		// This is its decimal value.
		d := decibinaryToBinary(output)
		// The lower bound in the array of partial sums should be the rank of the minimal
		// decibinary representation of d.
		lowerBound := sort.Search(len(partialSums), func(ix int) bool { return partialSums[ix] >= rank }) - 1
		t.Logf("Decimal %d ranked between %d and %d; looking for %d", d, partialSums[lowerBound], partialSums[lowerBound + 1], rank)
		if (!(rank >= partialSums[lowerBound] && rank < partialSums[lowerBound + 1])) {
			t.Errorf("Input %d expected to be between %d and %d (decimal value %d)", rank, partialSums[lowerBound], partialSums[lowerBound + 1], d)
		}
	}
}
