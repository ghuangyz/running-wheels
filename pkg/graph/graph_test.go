package graph

import (
	"fmt"
	"github.com/ghuangyz/running-wheels/pkg/types"
	"strings"
	"testing"
)

func makeGraphFromEdges(edges []string) *Graph {
	graph := Graph{}
	nodesTable := make(map[string]*Node)
	for _, edge := range edges {
		node_names := strings.Split(edge, "->")
		from := node_names[0]
		to := node_names[1]
		if _, exist := nodesTable[from]; !exist {
			nodesTable[from] = graph.MakeNode(from)
		}
		if _, exist := nodesTable[to]; !exist {
			nodesTable[to] = graph.MakeNode(to)
		}

		err := graph.AddEdge(nodesTable[from], nodesTable[to])
		if err != nil {
			fmt.Println(types.ErrorStackTrace(err))
		}
	}
	return &graph
}

func TestGraphHasCycle(t *testing.T) {
	testGraph1 := makeGraphFromEdges([]string{
		"a->c",
		"b->c",
		"c->d",
	})
	hasCycle, path := testGraph1.HasCycle()

	if hasCycle {
		fmt.Printf("Graph1 has cyclic!\n")
		fmt.Printf("%+v\n", path)
	} else {
		fmt.Printf("Graph1 is acyclic!\n")
	}

	testGraph2 := makeGraphFromEdges([]string{
		"a->c",
		"a->b",
		"c->d",
		"b->d",
		"d->a",
	})
	hasCycle, path = testGraph2.HasCycle()
	if hasCycle {
		fmt.Printf("Graph2 has cyclic!\n")
		fmt.Printf("Cycle Path: ")
		for index, element := range path {
			fmt.Printf("%s", element.Value)
			if index+1 < len(path) {
				fmt.Printf("->")
			}
		}
		fmt.Printf("\n")
	} else {
		fmt.Printf("Graph2 is acyclic!\n")
	}

	testGraph3 := makeGraphFromEdges([]string{
		"a->b",
		"a->b",
	})
	testGraph3.HasCycle()
}
