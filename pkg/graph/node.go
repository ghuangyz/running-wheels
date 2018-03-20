package graph

import (
	"fmt"
	"github.com/ghuangyz/running-wheels/pkg/types"
	"sync/atomic"
)

const (
	NullNodeError      = "NullNodeError"
	DuplicateEdgeError = "DuplicateEdgeError"
	InvalidEdgeError   = "InvalidEdgeError"
)

// Global graph node count, also graph node indices are assigned by this count
var nodeCount int64 = 0

func nextNodeId() int64 {
	prev := atomic.LoadInt64(&nodeCount)
	atomic.AddInt64(&nodeCount, 1)
	return prev
}

// class Node
type Node struct {
	id        int64
	neighbors map[int64]*Node
	Value     interface{}
}

func NewNode(value interface{}) *Node {
	node := new(Node)
	node.neighbors = make(map[int64]*Node)
	node.Value = value
	node.id = nextNodeId()
	return node
}

func (node *Node) Id() int64 {
	return node.id
}

func (node *Node) Neighbors() []*Node {
	var neighbors []*Node
	for _, v := range node.neighbors {
		neighbors = append(neighbors, v)
	}
	return neighbors
}

func (node *Node) HasNeighbor() bool {
	return len(node.neighbors) > 0
}

func (from *Node) ConnectTo(to *Node) error {
	if from == nil {
		return types.NewError(NullNodeError, "")
	}

	errMsg := fmt.Sprintf("Node %d->Node %d", from.id, to.id)
	if from == to {
		return types.NewError(DuplicateEdgeError, errMsg)
	}

	if _, exist := from.neighbors[to.id]; exist {
		return types.NewError(InvalidEdgeError, errMsg)
	}

	from.neighbors[to.id] = to
	return nil
}
