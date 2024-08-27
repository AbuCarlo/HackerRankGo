package graphs

import (
	"bufio"
	"strconv"
	"strings"
	"testing"
)

func maxRegion(grid [][]int32) int32 {
    type _Cell struct {
		row, column int
	}

	adjacency := make(map[_Cell][]_Cell)

	for i, r := range grid {
		for j, value := range r {
			if value == 0 {
				continue
			}
			cell := _Cell{i, j}
			if j > 0 {
				// to the left
				if grid[i][j - 1] == 1 {
					adjacency[cell] = append(adjacency[cell], _Cell{i, j - 1})
				}
				// up, left
				if i > 0 && grid[i - 1][j - 1] == 1 {
					adjacency[cell] = append(adjacency[cell], _Cell{i - 1, j - 1})
				}
			}
			// up
			if i > 0 {
				if grid[i - 1][j] == 1 {
					adjacency[cell] = append(adjacency[cell], _Cell{i - 1, j})
				}

				// up, right
				if j < len(r) && grid[i - 1][j + 1] == 1 {
					adjacency[cell] = append(adjacency[cell], _Cell{i - 1, j + 1})
				}
			}
		}
	}

	return 0
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
		input string
		expected int32
	}{
		{`4
4
1 1 0 0
0 1 1 0
0 0 1 0
1 0 0 0`, 5},
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