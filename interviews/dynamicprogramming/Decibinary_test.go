package dynamicprogramming

// https://www.hackerrank.com/challenges/decibinary-numbers/problem

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
	"strconv"
	"testing"
)

const MaximumIndex = 1e16

// We know from experience that this is the largest decimal
// number allowed by the problem definition, i.e. the total number of
// decibinary numerals needed past this point will exceed
// MaximumIndex.
const MaximumDecimalNumber = 285112

// The numbers do start at 0 (see the problem definition).
// Allow one more element for the odd number 285113
var counts = make([]int64, MaximumDecimalNumber+2)
var partialSums = make([]int64, MaximumDecimalNumber+2)
var countsBySize = make(map[int]map[int]int)

func countNumerals() {

	counts[0] = 1
	counts[1] = 1

	countsBySize[1] = map[int]int{1: 1}

	for n := 2; n <= MaximumDecimalNumber; n += 2 {
		countsBySize[n] = make(map[int]int)
		var count int64
		// Populate the least-significant "decibinary" digit. How
		// many decibinary numerals correspond to the remaining
		// value?
		for least := n % 2; least < 10 && least <= n; least += 2 {
			most := (n - least) >> 1
			count += counts[most]

			if most == 0 {
				countsBySize[n][1] = 1
				// There are no more significant digits.
				// This numeral has 1-digit representation.
				continue
			}

			key := most
			if key > 1 && key%2 == 1 {
				key--
			}

			for prefixSize, prefixCount := range countsBySize[key] {
				countsBySize[n][prefixSize+1] += prefixCount
			}
		}
		counts[n] = count
		// For any even decimal number, the final digit will be no more than
		// 8. It's always possible to add 1 to final digit, giving you the next
		// decimal integer, which is perforce odd. For any digit, we can subtract
		// 2, halve it, and shift it leftward, to the next-higher order digit,
		// but an odd digit cannot be reduced below 1.
		counts[n+1] = count
	}
}

// This array allows a binary search for a decimal number based on the rank
// of the decibinary numeral. In other words, if we want the xth decibinary
// numeral, we can look for the lowest value in this array > than x. Its index
// will be the decimal number 1 greater than the one we want.
func calculatePartialSums() {
	var sum int64
	for i, c := range counts {
		sum += c
		partialSums[i] = sum
	}
}

func init() {
	countNumerals()
	calculatePartialSums()
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
	// We can't yet use Go 1.22
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}

	return a
}

func decibinaryArrayToDecibinary(a []int) int64 {
	var result int64
	for _, digit := range a {
		result *= 10
		result += int64(digit)
	}
	return result
}

func decibinaryArrayToInt(d []int) int {
	result := 0
	for _, digit := range d {
		result *= 2
		result += digit
	}
	return result
}

func decibinaryToInt(d int64) int {
	significance := 1
	result := 0
	for d > 0 {
		result += significance * int(d%10)
		d /= 10
		significance *= 2
	}
	return result
}

func highestDecibinaryNumeral(n int) int64 {
	if n == 0 {
		return 0
	}
	result := int64(0)

	for bit := 1 << (bits.Len(uint(n)) - 1); bit > 0; bit >>= 1 {
		result *= 10
		if n&bit != 0 {
			result++
		}
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

// TODO Cache these?
func countSuffixes(value int, size int) int {
	if value == 0 {
		return 1
	}
	count := 0
	if value > 1 && value%2 == 1 {
		value--
	}
	for s, f := range countsBySize[value] {
		if s <= size {
			count += f
		}
	}
	return count
}

func locateSized(a []int, operations int64) {
	if operations == 0 {
		return
	}

	if len(a) == 0 {
		if operations != 0 {
			panic("Expected 0")
		}
		return
	}

	// Do I have to roll the leading digit?
	for {
		tailValue := decibinaryArrayToInt(a[1:])
		tailVersions := countSuffixes(tailValue, len(a)-1)
		// If enough operations are possible on the suffix,
		// leave the highest digit in place and recurse.
		if int(operations) < tailVersions {
			locateSized(a[1:], operations)
			return
		}

		// TODO Explain this.
		operations -= int64(tailVersions)
		a[0] -= 1
		a[1] += 2
	
		if a[0] == 0 {
			a = a[1:]
		}
	
		for i, d := range a {
			if d > 9 {
				// d could be 11 or 10.
				a[i] -= d - 9
				a[i+1] += 2 * (d - 9)
			}
		}
	}
}

func locate(rank int64) int64 {
	if rank < 1 {
		panic("Queries are 1-based.")
	}

	nativeValue := rankToNative(rank)
	array := decibinaryToArray(highestDecibinaryNumeral(nativeValue))
	operations := partialSums[nativeValue] - rank

	locateSized(array, operations)

	return int64(decibinaryArrayToDecibinary(array))
}

// The lower bound in the array of partial sums should be the rank of the minimal
// decibinary representation of d.
func rankToNative(rank int64) int {
	at := sort.Search(len(partialSums), func(ix int) bool { return rank <= partialSums[ix] })
	return at
}

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

func readTestFiles(t *testing.T, n int) ([]int64, []int64) {
	inputFileName := fmt.Sprintf("decibinary/input%02d.txt", n)
	outputFileName := fmt.Sprintf("decibinary/output%02d.txt", n)

	t.Logf("Opening %s", inputFileName)
	inputFile, err := os.Open(inputFileName)
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

	outputFile, err := os.Open(outputFileName)
	if err != nil {
		t.Fatal(err)
	}
	defer outputFile.Close()
	outputs := readInt64(bufio.NewScanner(outputFile))
	if size != int64(len(outputs)) {
		t.Fatalf("Expected %d outputs; got %d", size, len(outputs))
	}

	return inputs, outputs
}

func TestBoundaries(t *testing.T) {

	testFileNumbers := []int{ 0, 3, 7, 9, 10 }

	for _, n := range testFileNumbers {

		inputs, outputs := readTestFiles(t, n)

		for i, rank := range inputs {
			// This is the rank of a decibinary number: the "query".
			// This is the decibinary numeral having that rank.
			expected := outputs[i]
			// This is its decimal representation.
			d := decibinaryToInt(expected)
			native := rankToNative(rank)
	
			if d != native {
				t.Errorf("Expected %d-th output %d does not match actual native integer %d at rank %d", i, expected, native, rank)
			}
	
			actual := locate(rank)
			if actual != expected {
				t.Errorf("Expected %d-th output %d does not match actual output %d for input %d", i, expected, actual, rank)
	
			}
		}
	}
}

func TestAlgorithm(t *testing.T) {
	type Table struct {
		query    int64
		response int64
	}

	table := []Table{
		{query: 1, response: 0},
		{query: 2, response: 1},
		{query: 3, response: 2},
		{query: 4, response: 10},
		{query: 5, response: 3},
		{query: 6, response: 11},
		{query: 7, response: 4},
		{query: 8, response: 12},
		{query: 9, response: 20},
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

func _main() {

	for i := 43; i <= 43; i++ {
		first := partialSums[i] - counts[i] + 1
		result := []int64{}
		for j := first; j <= partialSums[i]; j++ {
			d := locate(j)
			result = append(result, d)
			if i != decibinaryToInt(d) {
				fmt.Printf("Decimal value %d as decibinary %d fails round-trip: %d\n", i, d, decibinaryToInt(d))
			}
		}
		fmt.Printf("%d: %v\n", i, result)
	}

	fmt.Printf("Input %d; actual output %d; expected %d\n", 2714, locate(2714), 755)
}
