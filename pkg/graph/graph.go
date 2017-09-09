package graph

type Graph struct {
	nodes []*Node
}

type Path []*Node

func (graph *Graph) MakeNode(value interface{}) *Node {
	node := NewNode(value)
	graph.nodes = append(graph.nodes, node)
	return node
}

func (graph *Graph) AddEdge(from, to *Node) error {
	return from.ConnectTo(to)
}

func (graph *Graph) Size() int {
	return len(graph.nodes)
}

func (graph *Graph) Nodes() []*Node {
	return graph.nodes
}

const (
	notVisited = 0
	visiting   = 1
	visited    = 2
)

func (graph *Graph) HasCycle() (bool, Path) {
	states := make(map[*Node]int)
	for _, node := range graph.nodes {
		states[node] = notVisited
	}

	for _, node := range graph.nodes {
		if states[node] == notVisited {
			var path Path
			hasCycle := graph.hasCycleHelper(node, states, &path)
			if hasCycle {
				return hasCycle, path
			}
		}
	}
	return false, nil
}

func (graph *Graph) hasCycleHelper(node *Node, states map[*Node]int, path *Path) bool {
	states[node] = visiting
	*path = append(*path, node)
	for _, neighbor := range node.Neighbors() {
		switch state := states[neighbor]; state {
		case notVisited:
			hasCycle := graph.hasCycleHelper(neighbor, states, path)
			if hasCycle == true {
				return hasCycle
			}
		case visiting:
			*path = append(*path, neighbor)
			return true
		case visited:
		default:
			//do nothing
		}
	}
	states[node] = visited
	*path = (*path)[:len(*path)-1]
	return false
}
