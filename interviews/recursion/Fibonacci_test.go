package recursion

// https://www.hackerrank.com/challenges/ctci-fibonacci-numbers/problem

import (
	"fmt"
	"testing"
)

var m = map[int]int{
	0: 0,
	1: 1,
}

func Fibonacci(n int) int {
	if f, has := m[n]; has {
		return f
	} else {
		f = Fibonacci(n-1) + Fibonacci(n-2)
		m[n] = f
		return f
	}
}

func TestFibonacci(t *testing.T) {
	for i := 0; i <= 20; i++ {
		fmt.Printf("Fibonacci(%d) = %d\n", i, Fibonacci(i))
	}
}
