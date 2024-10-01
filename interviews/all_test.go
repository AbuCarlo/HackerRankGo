package hackerrank

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/abucarlo/hackerrank/interviews/dictionaries"
	"github.com/abucarlo/hackerrank/interviews/recursion"
	"github.com/abucarlo/hackerrank/interviews/search"
	"github.com/abucarlo/hackerrank/interviews/dynamicprogramming"

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
