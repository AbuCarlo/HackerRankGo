package graphs

/*
	https://www.hackerrank.com/challenges/ctci-connected-cell-in-a-grid/problem

	The problem is to find the largest disjoint set in a graph, i.e.
	to implement the "disjoint sets" algorithm with path reduction.
*/

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"pgregory.net/rapid"
)

type _Cell struct {
	row, column int
}

func maxRegion(grid [][]int32) int32 {
	adjacency := make(map[_Cell][]_Cell)

	for i, r := range grid {
		for j, value := range r {
			if value == 0 {
				continue
			}
			cell := _Cell{i, j}
			// Treat a cell as its own neighbor, so we can handle 1-size grids.
			adjacency[cell] = []_Cell{cell}
			if j > 0 {
				// to the left
				if grid[i][j-1] == 1 {
					adjacency[cell] = append(adjacency[cell], _Cell{i, j - 1})
				}
				// up, left
				if i > 0 && grid[i-1][j-1] == 1 {
					adjacency[cell] = append(adjacency[cell], _Cell{i - 1, j - 1})
				}
			}
			// up
			if i > 0 {
				if grid[i-1][j] == 1 {
					adjacency[cell] = append(adjacency[cell], _Cell{i - 1, j})
				}

				// up, right
				if j < len(r) - 1 && grid[i-1][j+1] == 1 {
					adjacency[cell] = append(adjacency[cell], _Cell{i - 1, j + 1})
				}
			}
		}
	}

	parents := make(map[_Cell]_Cell)

	disjoints := make(map[_Cell][]_Cell)

	for cell := range adjacency {
		parents[cell] = cell
		// Initially all singletons.
		disjoints[cell] = []_Cell{cell}
	}

	_findRoot := func(cell _Cell) _Cell {
		parent := parents[cell]
		// Follow the path to the root, compressing all the while.
		// See https://en.wikipedia.org/wiki/Disjoint-set_data_structure#Finding_set_representatives
		for parent != parents[parent] {
			parent, parents[parent] = parents[parent], parents[parents[parent]]
		}

		return parent
	}

	// Now give each subgraph its adjacency matrix.
	for u, s := range adjacency {
		for _, v := range s {
			// This can't be lifted out of the loop: the root may have changed.
			x := _findRoot(u)
			y := _findRoot(v)
			if x == y {
				continue
			}
			// Minor optimization: combine the smaller subgraph into the larger one.
			if len(disjoints[x]) < len(disjoints[y]) {
				x, y = y, x
			}
			parents[y] = x
			disjoints[x] = append(disjoints[x], disjoints[y]...)
			delete(disjoints, y)
		}
	}

	result := 0
	for _, d := range disjoints {
		if len(d) > result {
			result = len(d)
		}
	}

	return int32(result)
}

func loadGrid(reader *bufio.Reader) [][]int32 {
	var l string
	l, _ = reader.ReadString('\n')
	n, _ := strconv.ParseInt(strings.TrimSpace(l), 10, 32)
	l, _ = reader.ReadString('\n')
	m, _ := strconv.ParseInt(strings.TrimSpace(l), 10, 32)
	grid := make([][]int32, n)

	for i := range n {
		grid[i] = make([]int32, m)
		row, _ := reader.ReadString('\n')
		a := strings.Split(strings.TrimSpace(row), " ")
		for j, s := range a {
			v, _ := strconv.ParseInt(s, 10, 32)
			grid[i][j] = int32(v)
		}
	}

	return grid
}

func TestAllOnes(t *testing.T) {
	f := func(t *rapid.T) {
		height := rapid.Int32Range(0, 64).Draw(t, "height")
		width := rapid.Int32Range(0, 64).Draw(t, "width")

		row := make([]int32, width)
		for i := range(row) {
			row[i] = 1
		}
		grid := make([][]int32, height)
		for i := range grid {
			grid[i] = row
		}

		actual := maxRegion(grid)
		if actual != height * width {
			t.Errorf("%d by %d grid should produce answer %d", height, width, height * width)
		}
	}

	rapid.Check(t, f)
}

func TestGridSamples(t *testing.T) {
	testCases := []struct {
		input    string
		expected int32
	}{
		{`4
4
1 1 0 0
0 1 1 0
0 0 1 0
1 0 0 0`, 5},
		{
			`5
4
0 0 1 1
0 0 1 0
0 1 1 0
0 1 0 0
1 1 0 0`, 8,
		},
		{
			`5
5
1 0 1 1 0
1 1 0 0 1
0 1 1 1 0
0 0 0 0 1
1 1 1 0 0`, 10,
		},
	}

	for i, test := range testCases {
		t.Run(fmt.Sprintf("Sample_%d", i), func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(test.input))
			grid := loadGrid(reader)
			actual := maxRegion(grid)
			if actual != test.expected {
			t.Errorf("Test %d expected %d; got %d", i, test.expected, actual)
			}
		})
	}

}
