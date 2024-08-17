package graphs

type DisjointSet struct {
	root int
	size int
}

type UndirectedGraph struct {
	roots map[int]int
	sets map[int]DisjointSet
}

func (g *UndirectedGraph) MakeSet(v int) {
	if _, ok := g.sets[v]; !ok {
		g.sets[v] = DisjointSet{v, 1}
	}
}

func (g *UndirectedGraph) Merge(u, v) {
	
}

func (g *UndirectedGraph) Insert(u, v int) {
	// TODO: I have no idea if this is correct.
	g.roots[u] = g.Find(v)
	g.roots[v] = g.Find(u)
}

func (*UndirectedGraph g) Find(v int) {
	u, ok := g.roots[v]
	if !ok {
		g.roots[v] = v
		return v
	}

	for u != g.roots[u] {
		u, g.roots[u] = g.roots[u], g[g.roots[u]]
	}

	return u
}
