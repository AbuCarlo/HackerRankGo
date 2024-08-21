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
	colors map[int32]int64
}

func NewColoredGraph(order int32) *ColoredGraph {
	g := ColoredGraph{order, make(map[int32]*Set[int32]), make(map[int32]int64)}
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
	g := NewColoredGraph(4)
	g.AddEdge(1, 2)
	g.AddEdge(1, 3)
	g.AddEdge(2, 4)

	distances := g.FloydWarshall()
	t.Logf("Distances: %v", distances)
}


