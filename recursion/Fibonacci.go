package recursion

var m = map[int]int{
	0: 0,
	1: 1,
}

func Fibonacci(n int) int {
	f, has := m[n]
	if has {
		return f
	} else {
		f = Fibonacci(n-1) + Fibonacci(n-2)
		m[n] = f
		return f
	}
}
