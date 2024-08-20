package graphs

// See https://en.wikipedia.org/wiki/Disjoint-set_data_structure
// See 

import (
	"testing"

	"github.com/abucarlo/hackerrank/interviews/graphs/sets"
)

type DisjointSets map[int]*sets.Set[int]

type UndirectedGraph struct {
	// TODO: Rename "parents"
	roots map[int]int
	adjacency map[int]*sets.Set[int]
}

func NewUndirectedGraph() *UndirectedGraph {
	g := UndirectedGraph{ make(map[int]int), make(map[int]*sets.Set[int]) }
	return &g
}

func (g *UndirectedGraph) FindDisjoint() []UndirectedGraph {
	// To get distinct values in a slice:
	// https://stackoverflow.com/a/76471309/476942
	// To get a map's keys as a slice, see
	// https://stackoverflow.com/a/69889828/476942
	// I don't want to exploit any of these options here,
	// in case HackerRank doesn't let me.
	// See mainly https://en.wikipedia.org/wiki/Kruskal%27s_algorithm

	disjoints := make(DisjointSets)
	// Now give each subgraph its adjacency matrix.
	for v, _ := range g.roots {
		root := g.FindRoot(v)
		s, has := disjoints[root]
		if !has {
			s = sets.New[int]()
			disjoints[root] = s
		}
		s.Add(v)
		s.Add(root)
	}

	var trees []UndirectedGraph 
	for _, s := range disjoints {
		// Construct a new subgraph
		// in which every adjacency list
		// has a size of 1 or 2.
		mst := NewUndirectedGraph()
		seen := sets.New[int]()
		for _, u := range s.Items() {
			if seen.Has(u) {
				continue
			}
			for _, v := range g.adjacency[u].Items() {
				if seen.Has(v) {
					continue
				}
				mst.Insert(u, v)
			}
		}
		
		trees = append(trees, *mst)
	}

	return trees
}

func (g *UndirectedGraph) Insert(u, v int) {
	// Since it's an undirected graph, let's just
	// decide to make the the lower vertex number
	// the root.
	if u > v {
		u, v = v, u
	}
	g.roots[v] = g.FindRoot(u)
	if _, ok := g.adjacency[u]; !ok {
		g.adjacency[u] = sets.New[int]()
	}
	if _, ok := g.adjacency[v]; !ok {
		g.adjacency[v] = sets.New[int]()
	}
	g.adjacency[u].Add(v)
	g.adjacency[v].Add(u)
}

func (g *UndirectedGraph) FindRoot(v int) int {
	u, ok := g.roots[v]
	if !ok {
		g.roots[v] = v
		return v
	}
	// Follow the path to the root, compressing all the while.
	// See https://en.wikipedia.org/wiki/Disjoint-set_data_structure#Finding_set_representatives
	for u != g.roots[u] {
		u, g.roots[u] = g.roots[u], g.roots[g.roots[u]]
	}

	return u
}

func TestStarGraph(t *testing.T) {
	g := NewUndirectedGraph()
	for i := 2; i < 10; i++ {
		g.Insert(1, i)
	}

	d := g.FindDisjoint()
	t.Logf("%v", d)
}

// TODO: Star graphs, paths.
