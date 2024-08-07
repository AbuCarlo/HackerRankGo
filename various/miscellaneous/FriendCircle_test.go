package miscellaneous

import (
	"testing"
)

func FriendCircle(queries [][]int) []int {
	friends := map[int]*map[int]bool{}
	max := 0
	result := make([]int, len(queries))
	for i, q := range queries {
		// Is one of these nodes already in the map?
		left, right := q[0], q[1]
		// Each node is its own friend.
		if _, ok := friends[left]; !ok {
			friends[left] = &map[int]bool{ left: true }
		}
		if _, ok := friends[right]; !ok {
			friends[right] = &map[int]bool{ right: true }
		}
		// Are the nodes already friends?
		if _, ok := (*friends[left])[right]; ok {
			// The current friend group's size hasn't changed.
			result[i] = max
			continue
		}
		// Copy right's friends to left.
		for friend := range *friends[right] {
			(*friends[left])[friend] = true
		}
		// Every friend shares the same map.
		for friend := range *friends[left] {
			friends[friend] = friends[left]
		}
		if len(*friends[left]) > max {
			max = len(*friends[left])
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
