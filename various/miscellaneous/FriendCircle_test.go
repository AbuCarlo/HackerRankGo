package miscellaneous

import (
	"testing"
)

func FriendCircle(queries [][]int) []int {
	friends := map[int]map[int]bool{}
	max := 0
	result := make([]int, len(queries))
	for i, q := range queries {
		left, right := q[0], q[1]
		_, ok := friends[left]
		if !ok {
			friends[left] = map[int]bool{ left: true }
		}
		_, ok = friends[right]
		if !ok {
			friends[right] = map[int]bool{ right: true }
		}
		for friend := range friends[left] {
			friends[right][friend] = true
		}
		for friend := range friends[right] {
			friends[friend] = friends[right]
		}
		if len(friends[right]) > max {
			max = len(friends[right])
		}
		result[i] = max
	}
	return result
}

func TestFriendCircle(t *testing.T) {
	input := [][]int{
		{1, 2},
		{3, 4},
		{2, 3},
	}

	output := FriendCircle(input)

	t.Logf("%v", output)
}
