package trees

import (
	"bufio"
	"io"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"
)

type Problem struct {
	Values []int32
	Edges  [][]int32
}

type Node struct {
	Id       int32
	Value    int32
	Subtotal int64
	Parent   *Node
	Children []*Node
}

func wire(node *Node) {
	node.Subtotal = int64(node.Value)
	for i := 0; i < len(node.Children); i++ {
		child := node.Children[i]
		wire(child)
		node.Subtotal += child.Subtotal
	}
}

// Disjoint determines if one node is the descendant of another.
func Disjoint(m, n *Node) bool {
	// Edge case: a node is not disjoint with itself.
	if m == n {
		return false
	}
	// Since every node has a value of at least 1, every node
	// must have a subtotal greater than any of its descendants.
	// Therefore two nodes with the same subtotal must be disjoint.
	if m.Subtotal == n.Subtotal {
		return true
	}
	// The node with the lower subtotal cannot be the ancestor.
	if m.Subtotal < n.Subtotal {
		m, n = n, m
	}
	// Now follow the path to the root.
	for ; n != nil; n = n.Parent {
		if n.Parent == m {
			return false
		}
	}
	return true
}

func mkArray(n *Node, sorted []*Node) []*Node {
	sorted = append(sorted, n)
	for _, child := range n.Children {
		sorted = mkArray(child, sorted)
	}

	return sorted
}
 
func Solve(root *Node) int64 {
	wire(root)
	sortedBySubtotal := mkArray(root, nil)
	sort.Slice(sortedBySubtotal, func(i, j int) bool { return sortedBySubtotal[i].Subtotal < sortedBySubtotal[j].Subtotal })

	// First option: two disjoint subtrees have the same total value. Detach them
	// and add a balancing node to the remaining tree. Since every node has a value
	// of at least one, two with the same total value must be disjoint (i.e. one
	// cannot be the ancestor of another without having a higher total value).
	lowerBound := (root.Subtotal + 2) / 3
	// It's not clear from the problem statement, but yes, we are allowed to synthesize
	// an entirely new node to balance the tree. So the highest value to try is half
	// the total value of the tree.
	upperBound := root.Subtotal / 2
	for value := lowerBound; value <= upperBound; value++ {
		// A subtree with this subtotal will have to be balanced.
		remainder := root.Subtotal - 2*value
		// See https://pkg.go.dev/sort#Search
		targetIndex := sort.Search(len(sortedBySubtotal), func(i int) bool { return sortedBySubtotal[i].Subtotal >= value })

		// Are there at least 2 subtrees with this subtotal? They must be disjoint.
		if sortedBySubtotal[targetIndex].Subtotal == value && sortedBySubtotal[targetIndex+1].Subtotal == value {
			return value - remainder		
		}
		// The slight optimization of cutting off the search at targetIndex has no effect!
		remainderIndex := sort.Search(targetIndex, func(i int) bool { return sortedBySubtotal[i].Subtotal >= remainder })

		// Second option: There are two disjoint subtrees such that if they're both removed from the
		// tree, the remaining value will have the same subtotal as one of them. The lesser subtree
		// can then be balanced.

		for i := remainderIndex; sortedBySubtotal[i].Subtotal == remainder; i++ {
			for j := targetIndex; sortedBySubtotal[j].Subtotal == value; j++ {
				if Disjoint(sortedBySubtotal[j], sortedBySubtotal[i]) {
					return value - remainder
				}
			}
		}

		// Third option: walk up the tree from one of the selection.
		for i := targetIndex; sortedBySubtotal[i].Subtotal == value; i++ {
			candidate := sortedBySubtotal[i]
			for p := candidate.Parent; p != nil; p = p.Parent {
				if p.Subtotal-value == remainder || p.Subtotal-value == value {
					return value - remainder				
				}
			}
		}

		for i := remainderIndex; sortedBySubtotal[i].Subtotal == remainder; i++ {
			candidate := sortedBySubtotal[i]
			for p := candidate.Parent; p != nil; p = p.Parent {
				if p.Subtotal-remainder == value {
					return value - remainder
				}
			}
		}
	}

	return -1
}

func mkNode(node *Node, nodes []*Node, adjacency [][]int32) {
	for _, id := range adjacency[node.Id] {
		child := nodes[id]
		if child == node.Parent {
			continue
		}
		child.Parent = node
		node.Children = append(node.Children, child)

		mkNode(child, nodes, adjacency)
	}
}

func mkTree(c []int32, edges [][]int32) *Node {
	// The first value is 0: there is no node 0.
	adjacency := make([][]int32, len(c) + 1)

	for _, edge := range edges {
		u, v := edge[0], edge[1]
		adjacency[u] = append(adjacency[u], v)
		adjacency[v] = append(adjacency[v], u)
	}

	nodes := make([]*Node, len(c) + 1)
	for i, cost := range c {
		nodes[i + 1] = &Node{int32(i + 1), cost, 0, nil, nil}
	}

	r := rand.Int31n(int32(len(c)-1)) + 1
	root := nodes[r]
	mkNode(root, nodes, adjacency)

	return root
}

func balancedForest(c []int32, edges [][]int32) int64 {
	tree := mkTree(c, edges)
	return Solve(tree)
}

func TestSamples(t *testing.T) {
	type Test struct {
		path     string
		expected []int64
	}

	tests := []Test{
		{"input00.txt", []int64{2, -1}},
		{"input01.txt", []int64{-1, 10, 13, 5, 297}},
		{"input02.txt", []int64{1112, 2041, 959, -1, -1}},
		// {"input03.txt", []int64{1714, 5016, 759000000000, -1, 6}},
		// {"input04.txt", []int64{1357940809, 397705399909, 439044899265, 104805614260, -1}},
		// {"input05.txt", []int64{24999687487500, 16217607772, 4, 0, -1}},
		{"input06.txt", []int64{19}},
		{"input07.txt", []int64{4}},
	}

	for _, test := range tests {
		problems := read("./balanced-forest-inputs" + "/" + test.path)
		for i, problem := range problems {
			actual := balancedForest(problem.Values, problem.Edges)
			if actual != test.expected[i] {
				t.Errorf("Test of %s[%d] expected %d; was %d", test.path, i, test.expected[i], actual)
			}
		}
	}
}

func BenchmarkBalancedForest(b *testing.B) {
	problems := read("./balanced-forest-inputs" + "/" + "input02.txt")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, problem := range problems {
			balancedForest(problem.Values, problem.Edges)
		}
	}
}

func read(path string) []Problem {
	// This is basically the code from HackerRank.
	f, err := os.Open(path)
	checkError(err)
	reader := bufio.NewReaderSize(f, 16*1024*1024)

	qTemp, err := strconv.ParseInt(strings.TrimSpace(readLine(reader)), 10, 64)
	checkError(err)
	q := int32(qTemp)

	problems := []Problem{}

	for range q {

		nTemp, err := strconv.ParseInt(strings.TrimSpace(readLine(reader)), 10, 64)
		checkError(err)
		n := int(nTemp)
		// Nodes are 1-indexed.
		c := make([]int32, n)

		cTemp := strings.Split(strings.TrimSpace(readLine(reader)), " ")

		for i := range n {
			cItemTemp, err := strconv.ParseInt(cTemp[i], 10, 64)
			checkError(err)
			cItem := int32(cItemTemp)
			c[i] = cItem
		}

		edges := make([][]int32, n-1)
		for i := range n - 1 {
			a := strings.Split(strings.TrimRight(readLine(reader), " \t\r\n"), " ")
			// Assume the input is valid; no error-checking.
			parent, _ := strconv.ParseInt(a[0], 10, 32)
			child, _ := strconv.ParseInt(a[1], 10, 32)
			edges[i] = []int32{int32(parent), int32(child)}
		}

		problems = append(problems, Problem{c, edges})
	}
	return problems
}

func readLine(reader *bufio.Reader) string {
	str, _, err := reader.ReadLine()
	if err == io.EOF {
		return ""
	}

	return strings.TrimRight(string(str), "\r\n")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
