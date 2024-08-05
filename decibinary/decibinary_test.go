package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"testing"
)

func readInt64(scanner *bufio.Scanner) []int64 {
	inputs := make([]int64, 0)
	for scanner.Scan() {
		i, _ := strconv.ParseInt(scanner.Text(), 10, 64)
		inputs = append(inputs, i)
	}
	return inputs
}

func BenchmarkIntToArray(b *testing.B) {
	inputs := []int64{100, 1000, 74383, 35700000, 1000000000000}

	for _, input := range inputs {
		b.Run(fmt.Sprintf("input_size_%d", input), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				decibinaryToArray(input)
			}
		})
	}

	for n := 0; n < b.N; n++ {
		decibinaryToArray(100000000)
	}
}

func FuzzRoundTripThroughArray(f *testing.F) {
	f.Fuzz(func(t *testing.T, input int) {
		if input < 1 || input > MaximumDecimalNumber {
			t.Skip()
		}
		highest := highestDecibinaryNumeral(input)
		highestArray := decibinaryToArray(highest)
		x := decibinaryArrayToInt(highestArray)
		if x != input {
			t.Errorf("Input %d has minimum decibinary numeral %d; this round-tripped as %d", input, highest, x)
		}

		lowest := lowestDecibinaryNumeral(input)
		lowestArray := decibinaryToArray(lowest)
		y := decibinaryArrayToInt(lowestArray)
		if y != input {
			t.Errorf("Input %d has maximum decibinary numeral %d; this round-tripped as %d", input, highest, x)

		}
	})
}

func FuzzHighestDecibinaryNumeral(f *testing.F) {
	f.Fuzz(func(t *testing.T, input int) {
		if input < 1 || input > MaximumDecimalNumber {
			t.Skip()
		}
		t.Logf("Input: %d", input)
		highest := highestDecibinaryNumeral(input)
		actual := fmt.Sprintf("%b", highest)
		expected, _ := strconv.ParseInt(actual, 2, 32)
		if input != int(expected) {
			t.Errorf("Input was %d; round-tripped as %d instead of %d", input, highest, expected)
		}
	})
}

func FuzzTranslation(f *testing.F) {
	f.Fuzz(func(t *testing.T, input int) {
		if input < 1 || input > MaximumDecimalNumber {
			t.Skip()
		}
		highest := highestDecibinaryNumeral(input)
		roundTrip := decibinaryToInt(highest)
		if roundTrip != input {
			t.Errorf("Input %d has maximum decibinary numeral %d; this round-tripped as %d", input, highest, roundTrip)
		}

		lowest := lowestDecibinaryNumeral(input)
		roundTrip = decibinaryToInt(lowest)
		if roundTrip != input {
			t.Errorf("Input %d has minimum decibinary numeral %d; this round-tripped as %d", input, lowest, roundTrip)
		}
	})
}

func TestBoundaries(t *testing.T) {
	inputFile, err := os.Open("input03.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	scanner.Scan()
	size, _ := strconv.ParseInt(scanner.Text(), 10, 32)
	inputs := readInt64(scanner)
	if size != int64(len(inputs)) {
		t.Fatalf("Expected %d inputs; got %d", size, len(inputs))
	}

	outputFile, err := os.Open("output03.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer outputFile.Close()
	outputs := readInt64(bufio.NewScanner(outputFile))
	if size != int64(len(outputs)) {
		t.Fatalf("Expected %d outputs; got %d", size, len(outputs))
	}

	for i := 0; i < int(size); i++ {
		// This is the rank of a decibinary number: the "query".
		rank := inputs[i]
		// This is the decibinary numeral having that rank.
		expected := outputs[i]
		// This is its decimal representation.
		d := decibinaryToInt(expected)
		native := rankToNative(rank)

		if d != native {
			t.Errorf("Expected output %d does not match actual native integer %d at rank %d", expected, native, rank)
		}

		actual := locate(rank)
		if actual != expected {
			t.Errorf("Expected output %d does not match actual output %d", expected, actual)

		}
	}
}

func TestAlgorithm(t *testing.T) {
	type Table struct {
		query    int64
		response int64
	}

	table := []Table{
		// {query: 1, response: 0},
		// {query: 2, response: 1},
		// {query: 3, response: 2},
		// {query: 4, response: 10},
		// {query: 5, response: 3},
		// {query: 6, response: 11},
		// {query: 7, response: 4},
		// {query: 8, response: 12},
		// {query: 9, response: 20},
		{query: 10, response: 100},
		{query: 11, response: 5},
		{query: 14, response: 101},
		{query: 15, response: 6},
	}

	// Actually, we could generate this entire test from our arrays.
	for _, row := range table {
		actual := locate(row.query)
		if actual != row.response {
			t.Errorf("For %v got %d", row, actual)
		}
	}
}
