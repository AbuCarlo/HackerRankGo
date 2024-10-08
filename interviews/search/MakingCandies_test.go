package search

// https://www.hackerrank.com/challenges/making-candies/problem

import (
	"math"
	"testing"
)

func minimumPasses(machines, workers, price, target int) int {
	if target <= price {
		return int(math.Ceil(float64(target) / float64(machines*workers)))
	}

	candies := 0
	iterations := 0
	current := math.MaxInt

	for candies < target {
		if candies < price {
			i := int(math.Ceil(float64(price-candies) / float64(machines*workers)))
			iterations += i
			candies += machines * workers * i
			continue
		}

		purchased := candies/price
		candies = candies%price
		assets := machines + workers + purchased
		half := assets / 2

		if machines > workers {
			machines = max(machines, half)
			workers = assets - machines
		} else {
			workers = max(workers, half)
			machines = assets - workers
		}

		iterations += 1
		candies += machines * workers
		current = min(current, iterations+int(math.Ceil(float64(target-candies)/float64(machines*workers))))
	}

	return min(current, iterations)
}

func TestMinimumPasses(t *testing.T) {
	tests := []struct {
		inputs   []int
		expected int
	}{
		{[]int{3, 1, 2, 12}, 3},
		{[]int{1, 1, 6, 45}, 16},
		{[]int{184889632, 5184889632, 20, 10000}, 1},
	}
	for i, test := range tests {
		m, w, p, n := test.inputs[0], test.inputs[1], test.inputs[2], test.inputs[3]
		actual := minimumPasses(m, w, p, n)
		if actual != test.expected {
			t.Errorf("Test %d on %v expected %d; got %d", i, test.inputs, test.expected, actual)
		}
	}
}
