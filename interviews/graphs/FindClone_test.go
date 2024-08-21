package graphs

import (
	"math"
	"slices"
	"testing"
)

type ColoredGraph struct {
	order int32
	adjacency map[int32]*Set[int32]
	// This seems insane.
	colors map[int32]int32
}

func NewColoredGraph(order int32) *ColoredGraph {
	g := ColoredGraph{order, make(map[int32]*Set[int32]), make(map[int32]int32)}
	return &g
}

func (g *ColoredGraph) AddEdge(u, v int32) {
	if _, ok := g.adjacency[u]; !ok {
		g.adjacency[u] = NewSet[int32]()
	}
	if _, ok := g.adjacency[v]; !ok {
		g.adjacency[v] = NewSet[int32]()
	}
	g.adjacency[u].Add(v)
	g.adjacency[v].Add(u)
}

func (g *ColoredGraph) Color(colors... int32) {
	for i, color := range colors {
		g.colors[int32(i + 1)] = color
	}
}

func (g *ColoredGraph) FindClone(color int32) int64 {
	distances := g.FloydWarshall()
	result := int64(math.MaxInt64)
	for v := int32(1); v <= g.order; v++ {
		if g.colors[v] != color {
			continue
		}
		for u := v + 1; u <= g.order; u++ {
			if g.colors[u] != color {
				continue
			}
			result = min(result, distances[v][u])
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
	for i := range pattern {
		pattern[i] = math.MaxInt32
	}
	for i := range distances {
		distances[i] = slices.Clone(pattern)
	}
	for v, a := range g.adjacency {
		distances[v][v] = 0
		if a == nil {
			continue
		}
		for _, u := range a.Items() {
			distances[v][u] = 1
			distances[u][v] = 1
		}
	}
	// Cut this in half.
	for k := int32(1); k <= g.order; k++ {
		for j := int32(1); j <= g.order; j++ {
			for i := int32(1); i <= g.order; i++ {
				choice := distances[i][k] + distances[k][j]
				if choice < distances[i][j] {
					distances[i][j] = choice
					// Eliminate these.
					distances[j][i] = choice
				}
			}
		}
	}

	return distances
}

func TestFindCloneSamples(t *testing.T) {
	testCases := []struct {
		order int32
		from []int32
		to []int32
		colors []int32
		clone int32
		expected int32
	}{
		{ 4, []int32{1, 1, 2}, []int32{2, 3, 4}, []int32{1, 2, 1, 1 }, 1, 1 },
		{ 4, []int32{1, 1, 4}, []int32{2, 3, 2}, []int32{1, 2, 3, 4}, 2, -1 },
	}

	for i, test := range testCases {
		g := NewColoredGraph(test.order)
		for j, u := range test.from {
			g.AddEdge(u, test.to[j])
		}
		g.Color(test.colors...)
		actual := g.FindClone(test.clone)
		if actual != int64(test.expected) {
			t.Errorf("Test %d expected %d, found %d", i, test.expected, actual)
		} else {
			t.Logf("Test %d expected %d, found %d", i, test.expected, actual)
		}
	}
}


