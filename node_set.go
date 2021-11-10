package it

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"sort"
)

type NodeSet struct {
	ID     int64
	Nodes  []Node
	MinKey Key
	MaxKey Key
}

func NewNodeSet(nodes []Node) *NodeSet {
	set := &NodeSet{
		ID:    rand.Int63(),
		Nodes: nodes,
	}
	if len(nodes) > 0 {
		set.MinKey, _ = nodes[0].KeyRange()
		_, set.MaxKey = nodes[len(nodes)-1].KeyRange()
	}
	return set
}

var _ Node = new(NodeSet)

func (n *NodeSet) Equal(n2 Node) bool {
	switch n2 := n2.(type) {
	case *NodeSet:
		return n.ID == n2.ID
	}
	panic("bad type")
}

func (n NodeSet) KeyRange() (Key, Key) {
	return n.MinKey, n.MaxKey
}

func (n *NodeSet) Mutate(
	ctx Scope,
	path KeyPath,
	fn func(Node) (Node, error),
) (
	retNode Node,
	err error,
) {

	if len(path) == 0 {
		return nil, we(ErrNotFound)
	}

	var nodes []Node
	if n != nil {
		nodes = n.Nodes
	}
	key := path[0]

	// search
	i := sort.Search(len(nodes), func(i int) bool {
		min, _ := nodes[i].KeyRange()
		return Compare(min, key) >= 0
	})

	if i == len(nodes) {
		// not found
		if len(path) > 1 {
			return nil, we(ErrNotFound)
		}
		newNode, err := fn(nil)
		if err != nil {
			return nil, we(err)
		}
		if newNode != nil {
			// append
			newNodes := make([]Node, len(nodes), len(nodes)+1)
			copy(newNodes, nodes)
			newNodes = append(newNodes, newNode)
			if !sort.SliceIsSorted(newNodes, func(i, j int) bool {
				min1, _ := newNodes[i].KeyRange()
				min2, _ := newNodes[i].KeyRange()
				return Compare(min1, min2) < 0
			}) {
				return nil, we(ErrBadOrder)
			}
			return NewNodeSet(newNodes), nil
		}
		// not changed
		return n, nil
	}

	min, max := nodes[i].KeyRange()
	if Compare(key, min) < 0 {
		// not found
		if len(path) > 1 {
			return nil, we(ErrNotFound)
		}
		newNode, err := fn(nil)
		if err != nil {
			return nil, we(err)
		}
		if newNode != nil {
			// insert
			newNodes := make([]Node, 0, len(nodes)+1)
			newNodes = append(newNodes, nodes[:i]...)
			newNodes = append(newNodes, newNode)
			newNodes = append(newNodes, nodes[i:]...)
			if !sort.SliceIsSorted(newNodes, func(i, j int) bool {
				min1, _ := newNodes[i].KeyRange()
				min2, _ := newNodes[i].KeyRange()
				return Compare(min1, min2) < 0
			}) {
				return nil, we(ErrBadOrder)
			}
			return NewNodeSet(newNodes), nil
		}
		// not changed
		return n, nil

	} else if Compare(min, key) <= 0 && Compare(key, max) <= 0 {
		// in range
		var newNode Node
		if min == key && max == key {
			// exactly
			newNode, err = nodes[i].Mutate(ctx, path[1:], fn)
		} else {
			// descend
			newNode, err = nodes[i].Mutate(ctx, path, fn)
		}
		if err != nil {
			return nil, we(err)
		}

		if newNode == nil {
			// delete
			newNodes := make([]Node, 0, len(nodes)-1)
			newNodes = append(newNodes, nodes[:i]...)
			newNodes = append(newNodes, nodes[i+1:]...)
			return NewNodeSet(newNodes), nil

		} else if !newNode.Equal(nodes[i]) {
			// replace
			newMin, newMax := newNode.KeyRange()
			if newMin != min || newMax != max {
				return nil, we(ErrInvalidName)
			}
			newNodes := make([]Node, 0, len(nodes))
			newNodes = append(newNodes, nodes[:i]...)
			newNodes = append(newNodes, newNode)
			newNodes = append(newNodes, nodes[i+1:]...)
			if !sort.SliceIsSorted(newNodes, func(i, j int) bool {
				min1, _ := newNodes[i].KeyRange()
				min2, _ := newNodes[i].KeyRange()
				return Compare(min1, min2) < 0
			}) {
				return nil, we(ErrBadOrder)
			}
			return NewNodeSet(newNodes), nil
		}

		// not changed
		return n, nil
	}

	panic("impossible")
}

func (n NodeSet) Dump(w io.Writer, level int) {
	if w == nil {
		w = os.Stdout
	}
	for _, node := range n.Nodes {
		node.Dump(w, level)
	}
}

func (n NodeSet) Range(cont Src) Src {
	i := 0
	var src Src
	src = func() (any, Src, error) {
		if i == len(n.Nodes) {
			return nil, cont, nil
		}
		i++
		switch node := n.Nodes[i-1].(type) {
		case *NodeSet:
			return node.Range(src)()
		default:
			return node, src, nil
		}
	}
	return src
}

func (n *NodeSet) Merge(ctx Scope, node2 Node) (Node, error) {
	if node2 == nil {
		return n, nil
	}
	if n == nil {
		return node2, nil
	}
	n2, ok := node2.(*NodeSet)
	if !ok {
		panic(fmt.Errorf("bad merge type: %T", node2))
	}
	if n2.Equal(n) {
		return n, nil
	}

	var nodes []Node
	nodes1 := n.Nodes
	nodes2 := n2.Nodes

	for {
		if len(nodes1) == 0 {
			nodes = append(nodes, nodes2...)
			break
		}
		if len(nodes2) == 0 {
			nodes = append(nodes, nodes1...)
			break
		}
		node1 := nodes1[0]
		node2 := nodes2[0]
		min1, _ := node1.KeyRange()
		min2, _ := node2.KeyRange()
		res := Compare(min1, min2)
		if res < 0 {
			nodes = append(nodes, node1)
			nodes1 = nodes1[1:]
			continue
		}
		if res > 0 {
			nodes = append(nodes, node2)
			nodes2 = nodes2[1:]
			continue
		}
		node, err := node1.Merge(ctx, node2)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
		nodes1 = nodes1[1:]
		nodes2 = nodes2[1:]
	}

	return NewNodeSet(nodes), nil
}
