package miscellaneous

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func FriendCircle(queries [][]int) []int {
	friendship := map[int]map[int]bool{}
	max := 0
	result := make([]int, len(queries))
	for i, q := range queries {
		// Is one of these nodes already in the map?
		left, right := q[0], q[1]
		// Map length is stored, so this is cheap.
		if len(friendship[left]) < len(friendship[right]) {
			left, right = right, left
		}
		// If left is new, right must also be new.
		if _, l := friendship[left]; !l {
			// This turns out to be more efficient, since most of these
			// maps are combined into larger ones. Allocating larger
			// ones turns out to waste time.
			m := map[int]bool{ left: true, right: true}
			friendship[left] = m
 			friendship[right] = m
		} else if _, r := friendship[right]; !r {
			// right is new; add it to left's friends.
			friendship[left][right] = true
			friendship[right] = friendship[left]
		} else {
			// Copy right's friends to left.
			for friend := range friendship[right] {
				friendship[left][friend] = true
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

func readEdges(f *os.File) [][]int {
	scanner := bufio.NewScanner(f)
	scanner.Scan()
	size, _ := strconv.Atoi(scanner.Text())
	edges := make([][]int, 0, size)
	for scanner.Scan() {
		edge := strings.Split(scanner.Text(), " ")
		right, _ := strconv.Atoi(edge[0])
		left, _ := strconv.Atoi(edge[0])
		edges = append(edges, []int{left, right})
	}

	if len(edges) != size {
		panic(fmt.Sprintf("Expected inputs of size %d; got %d", size, len(edges)))
	}

	return edges
}

func BenchmarkFriendCircle(b *testing.B) {
	inputFileName := "friend-circle/input10.txt"
	b.Logf("Opening %s", inputFileName)
	inputFile, err := os.Open(inputFileName)
	if err != nil {
		b.Fatal(err)
	}
	defer inputFile.Close()
	edges := readEdges(inputFile)

	b.Run("Test Case 10", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			FriendCircle(edges)
		}
	})
}

func TestFriendCircle(t *testing.T) {

	type Table struct {
		input    [][]int
		expected []int
	}
	table := []Table{
		{
			[][]int{
				{1, 2},
				{1, 3},
			},
			[]int{2, 3},
		},
		{
			[][]int{
				{1, 2},
				{3, 4},
				{2, 3},
			},
			[]int{2, 2, 4},
		},
		{
			[][]int{
				{6, 4},
				{5, 9},
				{8, 5},
				{4, 1},
				{1, 5},
				{7, 2},
				{4, 2},
				{7, 6},
			},
			[]int{2, 2, 3, 3, 6, 6, 8, 8},
		},
	}

	for i, row := range table {
		output := FriendCircle(row.input)
		if !slices.Equal(output, row.expected) {
			t.Errorf("Test %d expected %v, got %v", i, row.expected, output)
		} else {
			t.Logf("Test %d returns %v", i, output)
		}
	}
}