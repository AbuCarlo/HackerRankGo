package trees

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"

	"pgregory.net/rapid"
)

type Node struct {
	Id       int
	Value    int
	Parent   *Node
	Children []*Node
	// TODO: Make this a long.
	Sum int
}

func wire(node *Node) {
	sum := node.Value
	for i := 0; i < len(node.Children); i++ {
		child := node.Children[i]
		wire(child)
		sum += child.Sum
	}
	node.Sum = sum
}

func Disconnect(node *Node) *Node {
	var result *Node
	for parent := node.Parent; parent != nil; {
		result = &Node{parent.Id, parent.Value, nil, nil, parent.Sum - node.Sum}
		// Copy the children, with one exception.
		for _, n := range parent.Children {
			if n == node {
				result.Children = append(result.Children, result)
				n.Parent = result
			} else {
				result.Children = append(result.Children, n)
			}
		}
		tmp := parent
		parent = parent.Parent
		node = tmp
	}
	return result
}

func TestTreeGeneration(t *testing.T) {
	f := func(t *rapid.T) {
		size := rapid.IntRange(1, 100).Draw(t, "size")

		valueGenerator := rapid.IntRange(1, 100)
		blahGenerator := rapid.IntRange(1, 3)

		id := 1
		c := valueGenerator.Draw(t, "c")
		root := Node{id, c, nil, nil, 0}
		nodes := []*Node{&root}
		id++
		for id < size {
			node := nodes[0]
			nodes = nodes[1:]
			blah := blahGenerator.Draw(t, "blah")
			for j := 1; j <= blah && id <= size; j++ {
				c := valueGenerator.Draw(t, "c")
				child := Node{id, c, node, nil, c}
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

func balancedForest(_ *Node) int {
	return 0
}

func TestSamples(t *testing.T) {
	type Test struct {
		path     string
		expected int
	}

	tests := []Test{
		{"input00-1.txt", 1},
		{"input00-2.txt", 1},
		{"input06.txt", 1},
		{"input07.txt", 1},
	}

	for _, test := range tests {
		tree := read("./balanced-forest-inputs" + "/" + test.path)
		actual := balancedForest(tree)
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

	nodes := make([]*Node, n + 1)

	cTemp := strings.Split(strings.TrimSpace(readLine(reader)), " ")

	for i := range n {
		cItemTemp, err := strconv.ParseInt(cTemp[i], 10, 64)
		checkError(err)
		cItem := int(cItemTemp)
		id := i + 1
		nodes[id] = &Node{id, cItem, nil, nil, 0}
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
