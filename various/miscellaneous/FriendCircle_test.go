package miscellaneous

import (
	"testing"
)

func FriendCircle(queries [][]int) []int {
	friends := map[int][]int{}
	max := 0
	result := make([]int, len(queries))
	for i, q := range queries {
		// Is one of these nodes already in the map?
		left, right := q[0], q[1]
		_, ok := friends[left]
		if !ok {
			friends[left] = []int{left}
		}
		_, ok = friends[right]
		if !ok {
			friends[right] = []int{right}
		}
		// If left and right are already friends, they
		// point to equivalent slices. We can't check
		// reference equality in Go: a slice is a value 
		// object.
		friends[left] = append(friends[left], friends[right]...)
		for _, friend := range friends[left] {
			if friend != left {
				friends[friend] = friends[left]
			}
		}
		if len(friends[left]) > max {
			max = len(friends[left])
		}
		result[i] = max
	}
	return result
}

func TestFriendCircle(t *testing.T) {
	input00 := [][]int{
		{1, 2},
		{1, 3},
	}

	output00 := FriendCircle(input00)
	t.Logf("%v", output00)

	input := [][]int{
		{1, 2},
		{3, 4},
		{2, 3},
	}

	output := FriendCircle(input)

	t.Logf("%v", output)

	input02 := [][]int{
		{6, 4},
		{5, 9},
		{8, 5},
		{4, 1},
		{1, 5},
		{7, 2},
		{4, 2},
		{7, 6},
	}

	output02 := FriendCircle(input02)

	t.Logf("%v", output02)

}
