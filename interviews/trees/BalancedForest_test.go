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
	Sum      int
	Parent   *Node
	Children []*Node
	// TODO: Make this a long.
}

func wire(node *Node) {
	node.Sum = node.Value
	for i := 0; i < len(node.Children); i++ {
		child := node.Children[i]
		wire(child)
		node.Sum += child.Sum
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
	if m.Sum == n.Sum {
		return true
	}
	if m.Sum < n.Sum {
		m, n = n, m
	}

	for ; n != nil; n = n.Parent {
		if n.Parent == m {
			return false
		}
	}
	return true
}

func MkArray(n *Node, sorted []*Node) []*Node {
	sorted = append(sorted, n)
	for _, child := range n.Children {
		sorted = MkArray(child, sorted)
	}

	return sorted
}

func Solve(root *Node) int {
	sums := MkArray(root, nil)
	slices.SortFunc(sums, func(m *Node, n *Node) int { return m.Sum - n.Sum })

	// First option: two disjoint subtrees have the same total value. Detach them
	// and add a balancing node to the remaining tree. Since every node has a value
	// of at least one, two with the same total value must be disjoint (i.e. one
	// cannot be the ancestor of another without having a higher total value).
	resultForBlah := -1
	lowerBound := (root.Sum + 2) / 3
	// Any subtree must have a subtotal of at least 1; we're not going to
	// synthesize on from a null subtree. 
	upperBound := (root.Sum - 1) / 2
	for v := lowerBound; v <= upperBound; v++ {
		// See https://pkg.go.dev/sort#Search
		index := sort.Search(len(sums), func(i int) bool { return sums[i].Sum >= v })
		if sums[index].Sum != v {
			continue
		}
		// Are there at least 2 subtrees with this subtotal?
		if sums[index].Sum != sums[index+1].Sum {
			// We start at a node with no more than half the value of the entire tree,
			// so index + 1 will not be out of bounds.
			continue
		}
		resultForBlah = v - (root.Sum - 2*v)
		break
	}

	// Second option: There are two disjoint subtrees such that if they're both removed from the
	// tree, the remaining value will have the same subtotal as one of them. The lesser subtree
	// can then be balanced.
	resultForPoo := -1
	for v := lowerBound; v <= upperBound; v++ {
		index := sort.Search(len(sums), func(i int) bool { return sums[i].Sum >= v })
		// We could just count down here.
		if sums[index].Sum != v {
			// TODO Raise v here.
			continue
		}
		target := root.Sum - 2 * v
		blah := sort.Search(len(sums), func(i int) bool { return sums[i].Sum >= target })
		if sums[blah].Sum != target {
			continue
		}
		// Filter out descendants.
		for i := blah; sums[i].Sum == target; i++ {
			if Disjoint(sums[index], sums[i]) {
				resultForPoo = v - sums[i].Sum
				if resultForBlah == -1 || resultForPoo < resultForBlah  {
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
				node.Sum += child.Value

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
				sum += child.Sum
				if child.Parent != node {
					t.Errorf("Node %d has child %d, but child does not point to its parent.", node.Id, child.Id)
				}
			}

			if sum != node.Sum {
				t.Errorf("Node %d has a sum of %d, but its children sum to %d", node.Id, node.Sum, sum)
			}

		}
	}
	rapid.Check(t, f)
}

func TestSamples(t *testing.T) {
	type Test struct {
		path     string
		expected int
	}

	tests := []Test{
		{"sample00-1.txt", 2},
		{"sample00-2.txt", -1},
		{"input06.txt", 19},
		{"input07.txt", 4},
	}

	for _, test := range tests {
		tree := read("./balanced-forest-inputs" + "/" + test.path)
		wire(tree)
		actual := Solve(tree)
		if actual != test.expected {
			t.Errorf("Test of %s expected %d; was %d", test.path, test.expected, actual)
		}
	}
}

func read(path string) *Node {
	// This is basically the code from HackerRank.
	f, err := os.Open(path)
	checkError(err)
	reader := bufio.NewReaderSize(f, 16*1024*1024)

	nTemp, err := strconv.ParseInt(strings.TrimSpace(readLine(reader)), 10, 64)
	checkError(err)
	n := int(nTemp)

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
		parent, _ := strconv.ParseInt(a[0], 10, 32)
		child, _ := strconv.ParseInt(a[1], 10, 32)
		nodes[child].Parent = nodes[parent]
		nodes[parent].Children = append(nodes[parent].Children, nodes[child])
	}

	return nodes[1]
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
