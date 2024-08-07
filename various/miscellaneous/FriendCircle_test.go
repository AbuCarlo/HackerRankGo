package miscellaneous

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

func FriendCircle(queries [][]int) []int {
	friendship := map[int]*[]int{}
	max := 0
	result := make([]int, len(queries))
	for i, q := range queries {
		// Is one of these nodes already in the map?
		left, right := q[0], q[1]
		if _, l := friendship[left]; !l {
			// left is new.
			if _, r := friendship[right]; r {
				// Add it to an existing friend group.
				s := append(*friendship[right], left)
				for _, friend := range s {
					friendship[friend] = &s
				}
			} else {
				s := make([]int, 0, 32)
				s = append(s, left, right)
				friendship[right] = &s
				friendship[left] = friendship[right]
			}
		} else if _, r := friendship[right]; !r {
			// right is new; add it to left's friends.
			s := append(*friendship[left], right)
			for _, friend := range s {
				friendship[friend] = &s
			}
		} else if friendship[left] != friendship[right] {

			s := append(*friendship[left], *friendship[right]...)
			// Every friend shares the same map.
			for _, friend := range s {
				friendship[friend] = &s
			}
		}
		if len(*friendship[left]) > max {
			max = len(*friendship[left])
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
		input [][]int;
		expected []int
	}
	table := []Table{
		{
			[][]int{
				{1, 2},
				{1, 3},
			},
			[]int{},
		},
		{
			[][]int{
				{1, 2},
				{3, 4},
				{2, 3},
			},
			[]int{},
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
			[]int{},		
		},
	}

	for i, row := range table {
		output := FriendCircle(row.input)
		t.Logf("Test %d returns %v", i, output)
	}
}
