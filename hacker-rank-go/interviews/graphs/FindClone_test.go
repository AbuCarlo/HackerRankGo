package graphs

import (
	"bufio"
	"math"
	"os"
	"strconv"
	"testing"
)

type ColoredGraph struct {
	adjacency map[int32]*Set[int32]
	// This seems insane.
	colors []int32
}

type _DisjointSets map[int32]*Set[int32]

func NewColoredGraph() *ColoredGraph {
	g := ColoredGraph{make(map[int32]*Set[int32]), make([]int32, 1 << 10)}
	return &g
}

func (g *ColoredGraph) Order() int32 {
	return int32(len(g.adjacency))
}

func (g *ColoredGraph) AddEdge(u, v int32) {
	if _, ok := g.adjacency[u]; !ok {
		g.adjacency[u] = NewSet[int32]()
	}
	g.adjacency[u].Add(v)
	if _, ok := g.adjacency[v]; !ok {
		g.adjacency[v] = NewSet[int32]()
	}
	g.adjacency[v].Add(u)
}

func (g *ColoredGraph) SetColor(v int32, color int32) {
	if int32(len(g.colors)) < v+1 {
		l := len(g.colors)
		for l < int(v + 1) {
			l *= 2
		}
		tmp := make([]int32, l)
		copy(tmp, g.colors)
		g.colors = tmp
	}
	g.colors[int(v)] = color
}

func (g *ColoredGraph) FindDisconnected() _DisjointSets {
	parents := make(map[int32]int32)
	disjoints := make(_DisjointSets)

	_findRoot := func (v int32) int32 {
		parent := parents[v]
		// Follow the path to the root, compressing all the while.
		// See https://en.wikipedia.org/wiki/Disjoint-set_data_structure#Finding_set_representatives
		for parent != parents[parent] {
			parent, parents[parent] = parents[parent], parents[parents[parent]]
		}
	
		return parent
	}

	for v := range g.adjacency {
		parents[v] = v
		disjoints[v] = NewSet[int32]()
		disjoints[v].Add(v)
	}

	// Now give each subgraph its adjacency matrix.
	for u, s := range g.adjacency {
		for _, v := range s.Items() {
			x := _findRoot(u)
			y := _findRoot(v)
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

func (g *ColoredGraph) SolveSubgraph(color int32) int32 {
	countColored := 0
	for u := range g.adjacency {
		if g.colors[u] == color {
			countColored++
		}
	}
	if countColored < 2 {
		return -1
	}

	closestClone := int32(math.MaxInt32)
	// Test how many colored nodes there are.
	// https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm#Pseudocode
	for source := range g.adjacency {
		if g.colors[source] != color {
			continue
		}

		visited := NewSet[int32]()
		q := []int32{source}

		distances := make(map[int32]int)
		distances[source] = 0

		u := source

		for {

			visited.Add(u)

			q = q[1:]

			for _, v := range g.adjacency[u].Items() {
				if visited.Has(v) {
					continue
				}
				q = append(q, v)
				alt := distances[u] + 1
				if d, ok := distances[v]; ok {
					if alt < d {
						distances[v] = alt
					}
				} else {
					distances[v] = alt
				}
			}

			if len(q) == 0 {
				break
			}
			u = q[0]
		}

		for target, distance := range distances {
			if target == source {
				continue
			}
			if g.colors[target] == color {
				closestClone = min(closestClone, int32(distance))
			}
		}
	}

	if closestClone == math.MaxInt32 {
		return -1
	}
	return closestClone
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

		solution = min(solution, sub.SolveSubgraph(color))
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
		g.SetColor(int32(i+1), int32(color))
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
		actual := g.SolveSubgraph(test.clone)
		if actual != test.expected {
			t.Errorf("Test %d expected %d, found %d", i, test.expected, actual)
		} else {
			t.Logf("Test %d expected %d", i, test.expected)
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
	scanner.Split(bufio.ScanWords)

	scanner.Scan()
	order, _ := strconv.ParseInt(scanner.Text(), 10, 32)
	scanner.Scan()
	size, _ := strconv.ParseInt(scanner.Text(), 10, 32)

	g := NewColoredGraph()

	for range int(size) {
		scanner.Scan()
		u, _ := strconv.ParseInt(scanner.Text(), 10, 32)
		scanner.Scan()
		v, _ := strconv.ParseInt(scanner.Text(), 10, 32)
		g.AddEdge(int32(u), int32(v))
	}

	for v := int32(1); v <= int32(order); v++ {
		scanner.Scan()
		color, _ := strconv.ParseInt(scanner.Text(), 10, 32)
		g.SetColor(v, int32(color))
	}

	scanner.Scan()
	value, _ := strconv.ParseInt(scanner.Text(), 10, 32)

	return g, int32(value)
}

var directory = "./find-clone-inputs"

func TestFindCloneFiles(t *testing.T) {
	testCases := []struct {
		file     string
		expected int32
	}{
		{"input02.txt", -1},
		{"input04.txt", -1},
		{"input05.txt", -1},
		{"input06.txt", -1},
		{"input07.txt", -1},
	}
	for _, test := range testCases {
		t.Logf("Test %s expecting %d", test.file, test.expected)
		g, color := loadTestCase(directory + "/" + test.file)
		actual := g.SolveDijkstra(color)
		if actual != test.expected {
			t.Errorf("Test %s expected %d, found %d", test.file, test.expected, actual)
		}
	}
}

func BenchmarkFindClone(b *testing.B) {
	g, color := loadTestCase(directory + "/" + "input05.txt")

	for i := 0; i < b.N; i++ {
		g.SolveSubgraph(color)
	}
}
