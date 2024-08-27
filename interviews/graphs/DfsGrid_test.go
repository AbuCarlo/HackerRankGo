package graphs

import (
	"bufio"
	"strconv"
	"strings"
	"testing"
)

type _Cell struct {
	row, column int
}

func maxRegion(grid [][]int32) int32 {
	adjacency := make(map[_Cell][]_Cell)
	disjoints := make(map[_Cell][]_Cell)

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

	for cell := range adjacency {
		parents[cell] = cell
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
			x := _findRoot(u)
			y := _findRoot(v)
			if x == y {
				continue
			}
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
		result = max(result, len(d))
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
		reader := bufio.NewReader(strings.NewReader(test.input))
		grid := loadGrid(reader)
		actual := maxRegion(grid)
		if actual != test.expected {
			t.Errorf("Test %d expected %d; got %d", i, test.expected, actual)
		}
	}

}
