package hackerrank

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/abucarlo/hackerrank/dictionaries"
	"github.com/abucarlo/hackerrank/dynamicprogramming"
	"github.com/abucarlo/hackerrank/recursion"
	"github.com/abucarlo/hackerrank/search"
	"github.com/golang/glog"
)

func TestIceCream(t *testing.T) {
	var sample00 = []int32{1, 4, 5, 3, 2}
	var x, y = search.FindPair(sample00, 4)
	println(x, y)
	var sample01 = []int32{2, 2, 4, 3}
	var l, r = search.FindPair(sample01, 4)
	println(l, r)
}

func TestTwoStrings(t *testing.T) {
	tests := []struct {
		s      string
		t      string
		common bool
	}{
		{"HELLO", "WORLD", true},
		{"hi", "world", false},
		{"TONY", "NASSAR", true},
	}

	for _, test := range tests {

		testname := fmt.Sprintf("%s,%s", test.s, test.t)
		t.Run(testname, func(t *testing.T) {
			ans := dictionaries.TwoStrings(test.s, test.t)
			if ans != test.common {
				t.Errorf("got %t, want %t", ans, test.common)
			}
		})
	}
}

func TestGreedyFlorist(t *testing.T) {
	result00 := dynamicprogramming.Optimize(3, []int32{2, 5, 6})
	if result00 != 13 {
		t.Errorf("got %d, want %d", result00, 5)
	}
	result01 := dynamicprogramming.Optimize(2, []int32{2, 5, 6})
	if result01 != 15 {
		t.Errorf("got %d, want %d", result00, 5)
	}
}

func TestAbbreviation(t *testing.T) {
	tests := []struct {
		source string
		target string
		expect bool
	}{
		{"", "", true},
		{"abc", "abc", true},
		{"ABC", "ABC", true},
		{"abCde", "C", true},
		{"abcdE", "E", true},
	}
	for _, test := range tests {
		if test.expect {
			if !dynamicprogramming.Abbreviate(test.source, test.target) {
				t.Errorf("%s should match %s", test.source, test.target)
			}
		} else {
			if dynamicprogramming.Abbreviate(test.source, test.target) {
				t.Errorf("%s should not match %s", test.source, test.target)
			}
		}
	}
}

func TestAbbreviationTestCases(t *testing.T) {
	name := "dynamicprogramming/abbreviation-input12.txt"
	path, _ := filepath.Abs(name)
	file, e := os.Open(path)
	if e != nil {
		glog.Error(e)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	te := scanner.Text()
	n, _ := strconv.Atoi(te)
	for i := 0; i < n; i++ {
		scanner.Scan()
		source := scanner.Text()
		scanner.Scan()
		target := scanner.Text()
		result := dynamicprogramming.Abbreviate(source, target)
		fmt.Printf("Result of %s... / %s...: %t\n", source[0:10], target[0:10], result)
	}
}

func TestSuperDigit(t *testing.T) {
	result := recursion.SuperDigit("148", 3)
	fmt.Printf("Result: %d\n", result)
}

func TestDavisStaircase(t *testing.T) {
	result := recursion.Climb(7)
	fmt.Printf("Result: %v\n", result)
}

func TestFibonacci(t *testing.T) {
	for i := 0; i <= 20; i++ {
		fmt.Printf("Fibonacci(%d) = %d\n", i, recursion.Fibonacci(i))
	}
}

func TestDecibinary(t *testing.T) {
	tests := []struct {
		index    int64
		expected int64
	}{
		{1, 0},
		{2, 1},
		{3, 2},
		{4, 10},
		{10, 100},
		{8, 12},
		{23, 23},
		{19, 102},
		{16, 14},
		{26, 111},
		{7, 4},
		{6, 11},
	}

	for _, test := range tests {
		actual := dynamicprogramming.NthDecibinaryNumber(test.index)
		if actual != test.expected {
			t.Errorf("%d-th number should be %d, not %d\n", test.index, test.expected, actual)
		}
	}
}

func TestDecibinaryTestCases(t *testing.T) {
	name := "dynamicprogramming/decibinary-input07.txt"
	path, _ := filepath.Abs(name)
	fmt.Printf("Opening %s\n", name)
	file, e := os.Open(path)
	if e != nil {
		glog.Error(e)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()
	n, _ := strconv.Atoi(line)
	for i := 0; i < n; i++ {
		scanner.Scan()
		line = scanner.Text()
		x, _ := strconv.ParseInt(line, 10, 64)
		actual := dynamicprogramming.NthDecibinaryNumber(x)
		fmt.Printf("%d-th debinary number is %d\n", x, actual)
	}
}
