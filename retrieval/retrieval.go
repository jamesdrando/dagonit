package retrieval

import (
	"dagonit/graph"
)

// Engine handles subgraph retrieval from the code graph.
type Engine struct {
	Graph *graph.Graph
}

func NewEngine(g *graph.Graph) *Engine {
	return &Engine{Graph: g}
}

// FindNodesByTerm finds nodes by exact name or file path match.
func (e *Engine) FindNodesByTerm(term string) []graph.Node {
	var results []graph.Node
	nodes := e.Graph.GetNodes()
	for _, n := range nodes {
		if n.ID == term {
			results = append(results, n)
			continue
		}
		if name, ok := n.Metadata["name"].(string); ok && name == term {
			results = append(results, n)
		}
	}
	return results
}

// GetSubgraph retrieves a subgraph around a seed node.
func (e *Engine) GetSubgraph(seedID string, depth int) (*graph.Graph, error) {
	sub := graph.NewGraph()
	seed, ok := e.Graph.GetNode(seedID)
	if !ok {
		return nil, nil
	}
	sub.AddNode(seed)

	visited := make(map[string]bool)
	e.expand(seedID, depth, sub, visited)

	return sub, nil
}

func (e *Engine) expand(id string, depth int, sub *graph.Graph, visited map[string]bool) {
	if depth <= 0 || visited[id] {
		return
	}
	visited[id] = true

	neighbors := e.Graph.GetNeighbors(id)
	for _, n := range neighbors {
		if node, ok := e.Graph.GetNode(n.ID); ok {
			sub.AddNode(node)
			sub.AddEdge(graph.Edge{From: id, To: n.ID, Type: n.Type})
			e.expand(n.ID, depth-1, sub, visited)
		}
	}

	inNeighbors := e.Graph.GetInNeighbors(id)
	for _, n := range inNeighbors {
		if node, ok := e.Graph.GetNode(n.ID); ok {
			sub.AddNode(node)
			sub.AddEdge(graph.Edge{From: n.ID, To: id, Type: n.Type})
			e.expand(n.ID, depth-1, sub, visited)
		}
	}
}
