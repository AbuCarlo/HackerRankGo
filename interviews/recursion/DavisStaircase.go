package recursion

const Modulus int64 = 10000000007

func Climb(n int32) int32 {
	strides := []int32{1, 2, 3}
	ways := make([]int64, n+1)
	// There is 1 way to get to the 0th step, i.e. do nothing.
	ways[0] = 1
	for i := int32(1); i < n+1; i++ {
		for _, s := range strides {
			if i-s >= 0 {
				ways[i] += ways[i-s]
			}
		}
		ways[i] %= Modulus
	}
	return int32(ways[n])
}
