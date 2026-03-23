package graph

import (
	"reflect"
	"sort"
	"testing"
)

func TestComputeSCCs(t *testing.T) {
	g := NewGraph()
	// Create a graph with two SCCs and one isolated node
	// SCC 1: A -> B -> C -> A
	// SCC 2: D -> E -> D
	// Isolated: F
	// Connection: C -> D, E -> F

	nodes := []Node{
		{ID: "A", Type: NodeSymbol},
		{ID: "B", Type: NodeSymbol},
		{ID: "C", Type: NodeSymbol},
		{ID: "D", Type: NodeSymbol},
		{ID: "E", Type: NodeSymbol},
		{ID: "F", Type: NodeSymbol},
	}

	for _, n := range nodes {
		g.AddNode(n)
	}

	edges := []Edge{
		{From: "A", To: "B", Type: EdgeCalls},
		{From: "B", To: "C", Type: EdgeCalls},
		{From: "C", To: "A", Type: EdgeCalls},
		{From: "C", To: "D", Type: EdgeCalls},
		{From: "D", To: "E", Type: EdgeCalls},
		{From: "E", To: "D", Type: EdgeCalls},
		{From: "E", To: "F", Type: EdgeCalls},
	}

	for _, e := range edges {
		if err := g.AddEdge(e); err != nil {
			t.Fatalf("Failed to add edge: %v", err)
		}
	}

	sccs := g.ComputeSCCs()

	expectedSCCs := [][]string{
		{"A", "B", "C"},
		{"D", "E"},
		{"F"},
	}

	if !reflect.DeepEqual(sccs, expectedSCCs) {
		t.Errorf("Expected SCCs %v, but got %v", expectedSCCs, sccs)
	}
}

func TestToDAG(t *testing.T) {
	g := NewGraph()
	// SCC 1: A -> B -> A
	// SCC 2: C -> D -> C
	// Connection: B -> C

	g.AddNode(Node{ID: "A", Type: NodeSymbol})
	g.AddNode(Node{ID: "B", Type: NodeSymbol})
	g.AddNode(Node{ID: "C", Type: NodeSymbol})
	g.AddNode(Node{ID: "D", Type: NodeSymbol})

	g.AddEdge(Edge{From: "A", To: "B", Type: EdgeCalls})
	g.AddEdge(Edge{From: "B", To: "A", Type: EdgeCalls})
	g.AddEdge(Edge{From: "B", To: "C", Type: EdgeCalls})
	g.AddEdge(Edge{From: "C", To: "D", Type: EdgeCalls})
	g.AddEdge(Edge{From: "D", To: "C", Type: EdgeCalls})

	dag := g.ToDAG()

	// There should be 2 nodes in the DAG
	if len(dag.nodes) != 2 {
		t.Errorf("Expected 2 nodes in DAG, but got %d", len(dag.nodes))
	}

	// Identify SCCs based on their content
	var scc1ID, scc2ID string
	for id, node := range dag.nodes {
		nodesInSCC := node.Metadata["nodes"].([]string)
		sort.Strings(nodesInSCC)
		if reflect.DeepEqual(nodesInSCC, []string{"A", "B"}) {
			scc1ID = id
		} else if reflect.DeepEqual(nodesInSCC, []string{"C", "D"}) {
			scc2ID = id
		}
	}

	if scc1ID == "" || scc2ID == "" {
		t.Fatal("Could not identify SCC nodes in DAG")
	}

	// There should be an edge from SCC1 to SCC2
	neighbors := dag.GetNeighbors(scc1ID)
	found := false
	for _, neighbor := range neighbors {
		if neighbor.ID == scc2ID {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected edge from %s to %s in DAG, but not found", scc1ID, scc2ID)
	}
}
