package graph

import (
	"encoding/json"
	"errors"
	"sync"
)

// NodeType represents the type of a node in the graph.
type NodeType string

const (
	NodeRepository  NodeType = "repository"
	NodeModule      NodeType = "module"
	NodeFile        NodeType = "file"
	NodeSymbol      NodeType = "symbol"
	NodeSCC         NodeType = "scc"
	NodeTest        NodeType = "test"
	NodeBuildTarget NodeType = "build_target"
	NodeConfig      NodeType = "config"
)

// Node represents a single entity in the graph.
type Node struct {
	ID       string                 `json:"id"`
	Type     NodeType               `json:"type"`
	Metadata map[string]interface{} `json:"metadata"`
}

// EdgeType represents the type of connection between nodes.
type EdgeType string

const (
	EdgeContains   EdgeType = "contains"
	EdgeImports    EdgeType = "imports"
	EdgeCalls      EdgeType = "calls"
	EdgeTestedBy   EdgeType = "tested_by"
	EdgeBuilds     EdgeType = "builds"
	EdgeConfigures EdgeType = "configures"
)

// Edge represents a directed connection between two nodes.
type Edge struct {
	From string   `json:"from"`
	To   string   `json:"to"`
	Type EdgeType `json:"type"`
}

// Neighbor represents a connected node and the type of edge connecting to it.
type Neighbor struct {
	ID   string
	Type EdgeType
}

// Graph manages nodes and edges and provides methods for graph traversal.
type Graph struct {
	mu           sync.RWMutex
	nodes        map[string]Node
	outNeighbors map[string][]Neighbor
	inNeighbors  map[string][]Neighbor
}

// NewGraph creates a new instance of Graph.
func NewGraph() *Graph {
	return &Graph{
		nodes:        make(map[string]Node),
		outNeighbors: make(map[string][]Neighbor),
		inNeighbors:  make(map[string][]Neighbor),
	}
}

// AddNode adds a node to the graph. If a node with the same ID already exists, it is updated.
func (g *Graph) AddNode(node Node) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.nodes[node.ID] = node
}

// AddEdge adds a directed edge to the graph.
func (g *Graph) AddEdge(edge Edge) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, ok := g.nodes[edge.From]; !ok {
		return errors.New("source node does not exist: " + edge.From)
	}
	if _, ok := g.nodes[edge.To]; !ok {
		return errors.New("target node does not exist: " + edge.To)
	}

	g.outNeighbors[edge.From] = append(g.outNeighbors[edge.From], Neighbor{ID: edge.To, Type: edge.Type})
	g.inNeighbors[edge.To] = append(g.inNeighbors[edge.To], Neighbor{ID: edge.From, Type: edge.Type})
	return nil
}

// GetNode retrieves a node by its ID.
func (g *Graph) GetNode(id string) (Node, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	node, ok := g.nodes[id]
	return node, ok
}

// GetNodes returns all nodes in the graph.
func (g *Graph) GetNodes() []Node {
	g.mu.RLock()
	defer g.mu.RUnlock()

	nodes := make([]Node, 0, len(g.nodes))
	for _, node := range g.nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

// GetNeighbors returns all neighbors that the given node points to with edge types.
func (g *Graph) GetNeighbors(id string) []Neighbor {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.outNeighbors[id]
}

// GetInNeighbors returns all neighbors that point to the given node with edge types.
func (g *Graph) GetInNeighbors(id string) []Neighbor {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.inNeighbors[id]
}

// ExportToJSON serializes the graph's nodes and edges into a JSON format.
func (g *Graph) ExportToJSON() ([]byte, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	data := struct {
		Nodes []Node `json:"nodes"`
		Edges []Edge `json:"edges"`
	}{
		Nodes: make([]Node, 0, len(g.nodes)),
		Edges: []Edge{},
	}

	for _, node := range g.nodes {
		data.Nodes = append(data.Nodes, node)
	}

	for from, neighbors := range g.outNeighbors {
		for _, n := range neighbors {
			data.Edges = append(data.Edges, Edge{
				From: from,
				To:   n.ID,
				Type: n.Type,
			})
		}
	}

	return json.MarshalIndent(data, "", "  ")
}

// ImportFromJSON deserializes the graph's nodes and edges from a JSON format.
func (g *Graph) ImportFromJSON(data []byte) error {
	var input struct {
		Nodes []Node `json:"nodes"`
		Edges []Edge `json:"edges"`
	}

	if err := json.Unmarshal(data, &input); err != nil {
		return err
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	g.nodes = make(map[string]Node)
	g.outNeighbors = make(map[string][]Neighbor)
	g.inNeighbors = make(map[string][]Neighbor)

	for _, n := range input.Nodes {
		g.nodes[n.ID] = n
	}

	for _, e := range input.Edges {
		g.outNeighbors[e.From] = append(g.outNeighbors[e.From], Neighbor{ID: e.To, Type: e.Type})
		g.inNeighbors[e.To] = append(g.inNeighbors[e.To], Neighbor{ID: e.From, Type: e.Type})
	}

	return nil
}
