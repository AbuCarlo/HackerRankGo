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
        highest := highestDecibinaryNumeral(input)
		roundTrip := decibinaryToBinary(highest)
		if (roundTrip != input) {
			t.Errorf("Input %d has maximum decibinary numeral %d; this round-tripped as %d", input, highest, roundTrip)
		}

		lowest := highestDecibinaryNumeral(input)
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
