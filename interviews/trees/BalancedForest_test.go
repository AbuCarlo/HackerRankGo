package trees

import (
	"bufio"
	"fmt"
	"io"
	"math"
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

func (n *Node) String() string {
	s := fmt.Sprintf("{Id: %d, Value: %d, Subtotal: %d, Children: ", n.Id, n.Value, n.Subtotal)
	children := []int32{}
	for _, child := range n.Children {
		children = append(children, child.Id)
	}
	s += fmt.Sprintf("%v}", children)
	return s
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
	for ; n != nil && n.Subtotal <= m.Subtotal; n = n.Parent {
		if n.Parent == m {
			return false
		}
	}
	return true
}

func mkMap(nodes []*Node) map[int64][]*Node {	
	m := make(map[int64][]*Node)
	for _, n := range nodes {
		m[n.Subtotal] = append(m[n.Subtotal], n)
	}
	return m
}

func mkNode(node *Node, nodes []*Node, adjacency [][]int32) {
	node.Children = make([]*Node, 0, len(adjacency[node.Id]))
	for _, id := range adjacency[node.Id] {
		if node.Parent != nil && id == node.Parent.Id {
			continue
		}
		child := nodes[id]
		child.Parent = node
		node.Children = append(node.Children, child)

		mkNode(child, nodes, adjacency)
	}
}

func mkTree(c []int32, edges [][]int32) ([]*Node, *Node) {
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

	r := rand.Intn(len(nodes) - 1) + 1
	// TODO: Is 1 always the root?
	root := nodes[r]
	mkNode(root, nodes, adjacency)

	return nodes[1:], root
}

func balancedForest(c []int32, edges [][]int32) int64 {
	nodes, root := mkTree(c, edges)
	wire(root)
	// TODO: Keep only counts.
	sort.Slice(nodes, func (i, j int) bool { return nodes[i].Subtotal <= nodes[j].Subtotal })
	// TODO: Get rid of children.
	countsBySubtotal := mkMap(nodes)

	// First option: two disjoint subtrees have the same total value. Detach them
	// and add a balancing node to the remaining tree. Since every node has a value
	// of at least one, two with the same total value must be disjoint (i.e. one
	// cannot be the ancestor of another without having a higher total value).
	lowerBound := (root.Subtotal + 2) / 3
	// It's not clear from the problem statement, but yes, we are allowed to synthesize
	// an entirely new node to balance the tree. So the highest value to try is half
	// the total value of the tree.
	upperBound := root.Subtotal / 2

	current := int64(math.MaxInt64)

	for _, subtree := range nodes {
		// This is the minimum possible result.
		if current == 3 - (root.Subtotal % 3) {
			break
		}
		if subtree.Subtotal > upperBound {
			break
		}

		if subtree.Subtotal < lowerBound {
			target := (root.Subtotal - subtree.Subtotal) / 2

			mutated := make(map[int64]*Node)
			blah := make(map[*Node]struct{})
			for p := subtree.Parent; p != nil; p = p.Parent {
				mutated[p.Subtotal - subtree.Subtotal] = p
				blah[p] = struct{}{}
			}

			if _, ok := mutated[target]; ok {
				current = min(current, target - subtree.Subtotal)
			} else if nodes, ok := countsBySubtotal[target]; ok {
				for _, n := range nodes {
					if _, nok := blah[n]; !nok {
						current = min(current, target - subtree.Subtotal)
						break
					}
				}				
			}

		} else {
			target := subtree.Subtotal
			remainder := root.Subtotal - 2 * target

			if len(countsBySubtotal[target]) > 1 {
				current = min(current, target - remainder)
				continue
			}

			for p := subtree.Parent; p != nil; p = p.Parent {
				if p.Subtotal == 2 * target || p.Subtotal == target + remainder {
					current = min(current, target - remainder)
					break
				}
			}
		}
	}

	if current == int64(math.MaxInt64) {
		return -1
	}

	return current
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
		// {"input06.txt", []int64{19}},
		// {"input07.txt", []int64{4}},
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
