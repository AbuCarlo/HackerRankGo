package graphs

import (
	"bufio"
	"fmt"
	"math"
	"math/bits"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
)

type ColoredGraph struct {
	adjacency map[int32]*Set[int32]
	// This seems insane.
	colors []int32
}

type _DisjointSets map[int32]*Set[int32]

func NewColoredGraph() *ColoredGraph {
	g := ColoredGraph{make(map[int32]*Set[int32]), []int32{}}
	return &g
}

func (g *ColoredGraph) Order() int32 {
	return int32(len(g.adjacency))
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

func (g *ColoredGraph) Color(v int32, color int32) {
	if int32(len(g.colors)) < v+1 {
		newCapacity := 1 << (32 - bits.LeadingZeros32(uint32(v)) + 1)
		g.colors = slices.Grow(g.colors, newCapacity-len(g.colors)+1)
		g.colors = g.colors[:cap(g.colors)]
	}
	g.colors[int(v)] = color
}

func findRoot(parents map[int32]int32, v int32) int32 {
	parent := parents[v]
	// Follow the path to the root, compressing all the while.
	// See https://en.wikipedia.org/wiki/Disjoint-set_data_structure#Finding_set_representatives
	for parent != parents[parent] {
		parent, parents[parent] = parents[parent], parents[parents[parent]]
	}

	return parent
}

func (g *ColoredGraph) FindDisconnected() _DisjointSets {
	parents := make(map[int32]int32)
	disjoints := make(_DisjointSets)

	for v := int32(1); v <= g.Order(); v++ {
		parents[v] = v
		disjoints[v] = NewSet[int32]()
		disjoints[v].Add(v)
	}

	// Now give each subgraph its adjacency matrix.
	for u, s := range g.adjacency {
		for _, v := range s.Items() {
			x := findRoot(parents, u)
			y := findRoot(parents, v)
			if x == y {
				continue
			}
			if disjoints[x].Size() < disjoints[y].Size() {
				x, y = y, x
			}
			parents[y] = x
			disjoints[x].Union(disjoints[y])
			delete(disjoints, y)
		}
	}
	return disjoints
}

func (g *ColoredGraph) IsConnected(color int32) bool {
	fmt.Printf("Graph has order %d\n", g.Order())

	disjoints := g.FindDisconnected()

	fmt.Printf("%d disjoint sets\n", len(disjoints))
	for _, s := range disjoints {
		fmt.Printf("Disjoint set of size %d\n", s.Size())
		count := 0
		for _, v := range s.Items() {
			if g.colors[v] == color {
				count++
				if count == 2 {
					return true
				}
			}
		}
	}

	return false
}

func (g *ColoredGraph) Solve(color int32) int32 {
	// https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm#Pseudocode
	for source := range g.adjacency {
		if g.colors[source] != color {
			continue
		}
		unvisited := NewSet[int32]()
		unvisited.Add(source)

		distances := make(map[int32]int)
		distances[source] = 0

		for !unvisited.Empty() {
			u := unvisited.First()
			unvisited.Remove(u)

			for _, v := range g.adjacency[u].Items() {

				alt := distances[u] + 1
				if d, ok := distances[v]; ok {
					if alt < d {
						distances[v] = alt
					} else {
						distances[v] = alt
					}
				}
			}
		}

		closestClone := int32(math.MaxInt32)
		for target, distance := range distances {
			if target == source {
				continue
			}
			if g.colors[target] == color {
				closestClone = min(closestClone, int32(distance))
			}
		}
	}

	return 0
}

func (g *ColoredGraph) SolveDijkstra(color int32) int32 {

	disjoints := g.FindDisconnected()

	solution := int32(math.MaxInt32)

	for _, h := range disjoints {
		if h.Size() < 2 {
			continue
		}
		sub := ColoredGraph{make(map[int32]*Set[int32]), g.colors}
		for _, u := range h.Items() {
			sub.adjacency[u] = g.adjacency[u]
		}

		solution = min(solution, sub.Solve(color))
	}

	if solution == math.MaxInt32 {
		return -1
	}
	return solution
}

func ConstructTestCase(from []int32, to []int32, colors []int64) *ColoredGraph {
	g := NewColoredGraph()
	for j, u := range from {
		g.AddEdge(u, to[j])
	}
	for i, color := range colors {
		g.Color(int32(i+1), int32(color))
	}
	return g
}

func findShortest(_ int32, graphFrom []int32, graphTo []int32, ids []int64, val int32) int32 {
	g := ConstructTestCase(graphFrom, graphTo, ids)
	return int32(g.SolveDijkstra(val))
}

func TestFindCloneSamples(t *testing.T) {
	testCases := []struct {
		order    int32
		from     []int32
		to       []int32
		colors   []int64
		clone    int32
		expected int32
	}{
		// Sample 0, Test Case 0
		{4, []int32{1, 1, 2}, []int32{2, 3, 4}, []int64{1, 2, 1, 1}, 1, 1},
		// Sample 1, Test Case 1
		{4, []int32{1, 1, 4}, []int32{2, 3, 2}, []int64{1, 2, 3, 4}, 2, -1},
		// Sample 2
		{5, []int32{1, 1, 2, 3}, []int32{2, 3, 4, 5}, []int64{1, 2, 3, 3, 2}, 2, 3},
		// Test Case 2
	}

	for i, test := range testCases {
		g := ConstructTestCase(test.from, test.to, test.colors)
		actual := g.Solve(test.clone)
		if actual != test.expected {
			t.Errorf("Test %d expected %d, found %d", i, test.expected, actual)
		} else {
			t.Logf("Test %d expected %d, found %d", i, test.expected, actual)
		}
	}
}

func loadTestCase(file string) (*ColoredGraph, int32) {
	inputFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	scanner.Scan()

	graphNodesEdges := strings.Split(scanner.Text(), " ")
	// order, _ := strconv.ParseInt(graphNodesEdges[0], 10, 32)
	size, _ := strconv.ParseInt(graphNodesEdges[1], 10, 32)

	g := NewColoredGraph()

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
		g.Color(int32(i+1), int32(color))
	}

	scanner.Scan()
	value, _ := strconv.ParseInt(scanner.Text(), 10, 32)

	return g, int32(value)
}

var directory = "./find-clone-inputs"

func TestFindCloneFiles(t *testing.T) {
	// Benchmark?

	testCases := []struct {
		file     string
		expected int32
	}{
		{"input02.txt", -1},
		//{"input04.txt", -1},
		//{"input05.txt", -1},
	}
	for _, test := range testCases {
		t.Logf("Test %s expecting %d", test.file, test.expected)
		g, color := loadTestCase(directory + "/" + test.file)
		actual := g.Solve(color)
		if actual != test.expected {
			t.Errorf("Test %s expected %d, found %d", test.file, test.expected, actual)
		}
	}
}

func BenchmarkFindClone(b *testing.B) {
	g, color := loadTestCase(directory + "/" + "input04.txt")

	for i := 0; i < b.N; i++ {
		g.Solve(color)
	}
}
