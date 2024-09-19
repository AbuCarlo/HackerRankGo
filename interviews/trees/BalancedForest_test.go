package trees

import (
	"bufio"
	"io"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
	"testing"

	"pgregory.net/rapid"
)

type Node struct {
	Id       int
	Value    int
	Subtotal int
	Parent   *Node
	Children []*Node
	// TODO: Make this a long.
}

func wire(node *Node) {
	node.Subtotal = node.Value
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

func Solve(root *Node) int {
	wire(root)
	sums := mkArray(root, nil)
	slices.SortFunc(sums, func(m *Node, n *Node) int { return m.Subtotal - n.Subtotal })

	// First option: two disjoint subtrees have the same total value. Detach them
	// and add a balancing node to the remaining tree. Since every node has a value
	// of at least one, two with the same total value must be disjoint (i.e. one
	// cannot be the ancestor of another without having a higher total value).
	resultForBlah := -1
	lowerBound := (root.Subtotal + 2) / 3
	// Any subtree must have a subtotal of at least 1; we're not going to
	// synthesize on from a null subtree.
	upperBound := (root.Subtotal - 1) / 2
	for v := lowerBound; v <= upperBound; v++ {
		// See https://pkg.go.dev/sort#Search
		index := sort.Search(len(sums), func(i int) bool { return sums[i].Subtotal >= v })
		if sums[index].Subtotal != v {
			// v = sums[index].Sum
			continue
		}
		// Are there at least 2 subtrees with this subtotal?
		if sums[index].Subtotal != sums[index+1].Subtotal {
			// We start at a node with no more than half the value of the entire tree,
			// so index + 1 will not be out of bounds.
			// v = sums[index+1].Sum
			continue
		}
		resultForBlah = v - (root.Subtotal - 2*v)
		break
	}

	// Second option: There are two disjoint subtrees such that if they're both removed from the
	// tree, the remaining value will have the same subtotal as one of them. The lesser subtree
	// can then be balanced.
	resultForPoo := -1
	for v := lowerBound; v <= upperBound; v++ {
		index := sort.Search(len(sums), func(i int) bool { return sums[i].Subtotal >= v })
		// We could just count down here.
		if sums[index].Subtotal != v {
			// TODO Raise v here.
			continue
		}
		target := root.Subtotal - 2*v
		blah := sort.Search(len(sums), func(i int) bool { return sums[i].Subtotal >= target })
		if sums[blah].Subtotal != target {
			continue
		}
		// Filter out descendants.
		for i := blah; sums[i].Subtotal == target; i++ {
			if Disjoint(sums[index], sums[i]) {
				resultForPoo = v - sums[i].Subtotal
				if resultForBlah == -1 || resultForPoo < resultForBlah {
					return resultForPoo
				}
			}
		}
		break
	}

	if resultForBlah != -1 {
		return resultForBlah
	}

	return -1
}

func TestTreeGeneration(t *testing.T) {
	f := func(t *rapid.T) {
		size := rapid.IntRange(1, 100).Draw(t, "size")

		valueGenerator := rapid.IntRange(1, 100)
		blahGenerator := rapid.IntRange(1, 3)

		id := 1
		c := valueGenerator.Draw(t, "c")
		root := Node{id, c, 0, nil, nil}
		nodes := []*Node{&root}
		id++
		for id < size {
			node := nodes[0]
			nodes = nodes[1:]
			blah := blahGenerator.Draw(t, "blah")
			for j := 1; j <= blah && id <= size; j++ {
				c := valueGenerator.Draw(t, "c")
				child := Node{id, c, c, node, nil}
				node.Children = append(node.Children, &child)
				nodes = append(nodes, &child)
				child.Parent = node
				node.Subtotal += child.Value

				id++
			}
		}
		wire(&root)

		// Invariants: 1, any node's children should point back to it; 2, a node's Sum should == the Sum of its children plus its own value.

		traversal := []*Node{&root}
		for len(traversal) > 0 {
			node := traversal[0]
			sum := node.Value
			traversal = traversal[1:]
			for i := 0; i < len(node.Children); i++ {
				child := node.Children[i]
				sum += child.Subtotal
				if child.Parent != node {
					t.Errorf("Node %d has child %d, but child does not point to its parent.", node.Id, child.Id)
				}
			}

			if sum != node.Subtotal {
				t.Errorf("Node %d has a sum of %d, but its children sum to %d", node.Id, node.Subtotal, sum)
			}

		}
	}
	rapid.Check(t, f)
}

func TestSamples(t *testing.T) {
	type Test struct {
		path     string
		expected []int
	}

	tests := []Test{
		{"input00.txt", []int{2, -1}},
		{"input01.txt", []int{-1, 10, 13, 5, 297}},
		{"input02.txt", []int{1112, 2041, 959, -1, -1}},
		{"input03.txt", []int{1714, 5016, 759000000000, -1, 6}},
		{"input04.txt", []int{1357940809, 397705399909, 439044899265, 104805614260, -1}},
		{"input05.txt", []int{24999687487500, 16217607772, 4, 0, -1}},
		{"input06.txt", []int{19}},
		{"input07.txt", []int{4}},
	}

	for _, test := range tests {
		trees := read("./balanced-forest-inputs" + "/" + test.path)
		for i, tree := range trees {
			actual := Solve(tree)
			if actual != test.expected[i] {
				t.Errorf("Test of %s[%d] expected %d; was %d", test.path, i, test.expected[i], actual)
			}
		}

	}
}

func BenchmarkBalancedForest(b *testing.B) {
	trees := read("./balanced-forest-inputs" + "/" + "input07.txt")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, tree := range trees {
			Solve(tree)
		}
	}
}

func read(path string) []*Node {
	// This is basically the code from HackerRank.
	f, err := os.Open(path)
	checkError(err)
	reader := bufio.NewReaderSize(f, 16*1024*1024)

	qTemp, err := strconv.ParseInt(strings.TrimSpace(readLine(reader)), 10, 64)
	checkError(err)
	q := int32(qTemp)

	trees := []*Node{}

	for range q {

		nTemp, err := strconv.ParseInt(strings.TrimSpace(readLine(reader)), 10, 64)
		checkError(err)
		n := int(nTemp)
		// Nodes are 1-indexed.
		nodes := make([]*Node, n+1)

		cTemp := strings.Split(strings.TrimSpace(readLine(reader)), " ")

		for i := range n {
			cItemTemp, err := strconv.ParseInt(cTemp[i], 10, 64)
			checkError(err)
			cItem := int(cItemTemp)
			id := i + 1
			nodes[id] = &Node{id, cItem, 0, nil, nil}
		}

		for i := 0; i < int(n)-1; i++ {
			a := strings.Split(strings.TrimRight(readLine(reader), " \t\r\n"), " ")
			// Assume the input is valid; no error-checking.
			parent, _ := strconv.ParseInt(a[0], 10, 32)
			child, _ := strconv.ParseInt(a[1], 10, 32)
			// TODO: Can I cheat?
			if parent > child {
				parent, child = child, parent
			}
			nodes[child].Parent = nodes[parent]
			nodes[parent].Children = append(nodes[parent].Children, nodes[child])
		}
		trees = append(trees, nodes[1])
	}
	return trees
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
