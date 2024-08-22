package graphs

import (
	"bufio"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
)

type ColoredGraph struct {
	order int32
	adjacency map[int32]*Set[int32]
	// This seems insane.
	colors []int64
}

func NewColoredGraph(order int32) *ColoredGraph {
	g := ColoredGraph{order, make(map[int32]*Set[int32]), make([]int64, order + 1)}
	return &g
}

func (g *ColoredGraph) AddEdge(u, v int32) {
	if u < v {
		u, v = v, u
	}
	if _, ok := g.adjacency[u]; !ok {
		g.adjacency[u] = NewSet[int32]()
	}
	g.adjacency[u].Add(v)
}

func (g *ColoredGraph) Color(v int32, color int64) {
	g.colors[int(v)] = color
}

func (g *ColoredGraph) FindClone(color int64) int64 {
	distances := g.FloydWarshall()
	result := int64(math.MaxInt64)
	for u := int32(1); u <= g.order; u++ {
		if g.colors[u] != color {
			continue
		}
		for v := int32(1); v < u; v++ {
			if g.colors[v] != color {
				continue
			}
			result = min(result, distances[u][v])
		}
	}
	if result == math.MaxInt64 {
		return -1
	}
	return result
}

func (g *ColoredGraph) FloydWarshall() [][]int64 {
	distances := make([][]int64, g.order + 1)
	// There is no Fill() or Repeat() function yet.
	pattern := make([]int64, g.order + 1)
	for i := 1; i <= int(g.order); i++ {
		pattern[i] = math.MaxInt32
	}
	for i := 1; i <= int(g.order); i++ {
		distances[i] = slices.Clone(pattern[:i + 1])
	}
	for u, a := range g.adjacency {
		distances[u][u] = 0
		if a == nil {
			continue
		}
		for _, v := range a.Items() {
			distances[u][v] = 1
		}
	}
	// Cut this in half.
	for k := int32(1); k <= g.order; k++ {
		for i := int32(1); i <= g.order; i++ {
			for j := int32(1); j < i; j++ {
				var choice int64
				if i > k {
					choice = distances[i][k]
				} else {
					choice = distances[k][i]
				}
				if k > j {
					choice += distances[k][j]
				} else {
					choice += distances[j][k]
				}
				distances[i][j] = min(choice, distances[i][j])
			}
		}
	}

	return distances
}

func ConstructTestCase(order int32, from []int32, to []int32, colors []int64) *ColoredGraph {
	g := NewColoredGraph(order)
	for j, u := range from {
		g.AddEdge(u, to[j])
	}
	for i, color := range colors {
		g.Color(int32(i + 1), color)
	}
	return g
}

func findShortest(graphNodes int32, graphFrom []int32, graphTo []int32, ids []int64, val int32) int32 {
    g := ConstructTestCase(graphNodes, graphFrom, graphTo, ids)
	return int32(g.FindClone(int64(val)))
}

func TestFindCloneSamples(t *testing.T) {
	testCases := []struct {
		order int32
		from []int32
		to []int32
		colors []int64
		clone int32
		expected int32
	}{
		// Sample 0, Test Case 0
		{ 4, []int32{1, 1, 2}, []int32{2, 3, 4}, []int64{1, 2, 1, 1 }, 1, 1 },
		// Sample 1, Test Case 1
		{ 4, []int32{1, 1, 4}, []int32{2, 3, 2}, []int64{1, 2, 3, 4}, 2, -1 },
		// Sample 2
		{ 5, []int32{1, 1, 2, 3}, []int32{2, 3, 4, 5}, []int64{1, 2, 3, 3, 2}, 2, 3 },
		// Test Case 2
	}

	for i, test := range testCases {
		g := ConstructTestCase(test.order, test.from, test.to, test.colors)
		actual := g.FindClone(int64(test.clone))
		if actual != int64(test.expected) {
			t.Errorf("Test %d expected %d, found %d", i, test.expected, actual)
		} else {
			t.Logf("Test %d expected %d, found %d", i, test.expected, actual)
		}
	}
}

func loadTestCase(file string) (*ColoredGraph, int64) {
	inputFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	scanner.Scan()
	
    graphNodesEdges := strings.Split(scanner.Text(), " ")
    order, _ := strconv.ParseInt(graphNodesEdges[0], 10, 32)
    size, _ := strconv.ParseInt(graphNodesEdges[1], 10, 32)

	g := NewColoredGraph(int32(order))

	for range int(size) {
		scanner.Scan()
		edge := strings.Split(scanner.Text(), " ")
		u, _ := strconv.ParseInt(edge[0], 10, 32)
		v, _ := strconv.ParseInt(edge[1], 10, 32)
		g.AddEdge(int32(u), int32(v))
	}

	scanner.Scan()
	for i, s := range strings.Split(scanner.Text(), " ") {
		color, _ := strconv.ParseInt(s, 10, 32)
		g.Color(int32(i + 1), color)
	}

	scanner.Scan()
	value, _ := strconv.ParseInt(scanner.Text(), 10, 32)

	return g, value
}

var directory = "./find-clone-inputs"

func TestFindCloneFiles(t *testing.T) {
	// Benchmark?
	
	testCases := []struct{ file string; expected int64 }{
		{ "input02.txt", -1 },
		// { "input04.txt", -1 },
		// { "input05.txt", -1 },
	}
	for _, test := range testCases {
		g, color := loadTestCase(directory + "/" + test.file)
		actual := g.FindClone(color)
		if actual != test.expected {
			t.Errorf("Test %s expected %d, found %d", test.file, test.expected, actual)
		} else {
				t.Logf("Test %s expected %d, found %d", test.file, test.expected, actual)
		}
	}
}

func BenchmarkFindClone(b *testing.B) {
	g, color := loadTestCase(directory + "/" + "input04.txt")

	for i := 0; i < b.N; i++ {
        g.FindClone(color)
    }
}



