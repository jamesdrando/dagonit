package graph

import (
	"testing"
)

func TestGraph(t *testing.T) {
	g := NewGraph()

	node1 := Node{ID: "node1", Type: NodeFile, Metadata: map[string]interface{}{"path": "file1.go"}}
	node2 := Node{ID: "node2", Type: NodeSymbol, Metadata: map[string]interface{}{"name": "Function1"}}
	node3 := Node{ID: "node3", Type: NodeSymbol, Metadata: map[string]interface{}{"name": "Function2"}}

	g.AddNode(node1)
	g.AddNode(node2)
	g.AddNode(node3)

	err := g.AddEdge(Edge{From: "node1", To: "node2", Type: EdgeContains})
	if err != nil {
		t.Fatalf("Failed to add edge: %v", err)
	}

	err = g.AddEdge(Edge{From: "node2", To: "node3", Type: EdgeCalls})
	if err != nil {
		t.Fatalf("Failed to add edge: %v", err)
	}

	// Test GetNode
	retrievedNode, ok := g.GetNode("node1")
	if !ok || retrievedNode.ID != "node1" {
		t.Errorf("Expected node1, got %v", retrievedNode)
	}

	// Test GetNeighbors
	neighbors := g.GetNeighbors("node1")
	if len(neighbors) != 1 || neighbors[0].ID != "node2" {
		t.Errorf("Expected node2 as neighbor of node1, got %v", neighbors)
	}

	// Test GetInNeighbors
	inNeighbors := g.GetInNeighbors("node3")
	if len(inNeighbors) != 1 || inNeighbors[0].ID != "node2" {
		t.Errorf("Expected node2 as in-neighbor of node3, got %v", inNeighbors)
	}

	// Test error on missing node
	err = g.AddEdge(Edge{From: "node1", To: "missing_node", Type: EdgeCalls})
	if err == nil {
		t.Error("Expected error when adding edge to non-existent node")
	}
}
