package miscellaneous

import (
	"testing"
)

func FriendCircle(queries [][]int) []int {
	friendship := map[int]map[int]bool{}
	max := 0
	result := make([]int, len(queries))
	for i, q := range queries {
		// Is one of these nodes already in the map?
		left, right := q[0], q[1]
		if _, l := friendship[left]; !l {
			// left is new.
			if friends, r := friendship[right]; r {
				// Add it to an existing friend group.
				friendship[left] = friends
				friendship[left][left] = true
			} else {
				friendship[left] = map[int]bool{left: true, right: true}
				friendship[right] = friendship[left]
			}
		} else if _, r := friendship[right]; !r {
			// right is new; add it to left's friends.
			friendship[left][right] = true
			friendship[right] = friendship[left]
		} else {
			// Copy right's friends to left.
			for friend := range friendship[right] {
				(friendship[left])[friend] = true
			}
			// Every friend shares the same map.
			for friend := range friendship[left] {
				friendship[friend] = friendship[left]
			}
		}
		if len(friendship[left]) > max {
			max = len(friendship[left])
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
