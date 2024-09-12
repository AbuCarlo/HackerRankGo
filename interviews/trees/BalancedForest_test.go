package trees

import (
	"pgregory.net/rapid"
	"testing"
)

type Node struct {
	Id int
	Value int
	Parent *Node
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

func TestTreeGeneration(t *testing.T) {
	f := func(t *rapid.T) {
		size := rapid.IntRange(1, 100).Draw(t, "size")

		valueGenerator := rapid.IntRange(1, 100);
		blahGenerator := rapid.IntRange(1, 3);

		id := 1
		c := valueGenerator.Draw(t, "c")
		root := Node{id, c, nil, nil, c}
		nodes := []*Node{&root}
		id++
		// TODO Clean up the sums afterward, and the parent links.
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
				sum += node.Children[i].Sum
			}

			if sum != node.Sum {
				t.Errorf("Node %d has a sum of %d, but its children sum to %d", node.Id, node.Sum, sum);
			}

		}
	}
	rapid.Check(t, f)
}

