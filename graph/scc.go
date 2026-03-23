package graph

import (
	"fmt"
	"sort"
)

// ComputeSCCs finds the Strongly Connected Components (SCCs) in the graph using Tarjan's algorithm.
// It returns a slice of slices, where each inner slice contains the IDs of nodes in an SCC.
func (g *Graph) ComputeSCCs() [][]string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	var (
		index   int
		stack   []string
		onStack = make(map[string]bool)
		indices = make(map[string]int)
		lowlink = make(map[string]int)
		sccs    [][]string
	)

	// We need to iterate over all nodes to ensure we cover disconnected components.
	// To make it deterministic, we sort the node IDs.
	nodeIDs := make([]string, 0, len(g.nodes))
	for id := range g.nodes {
		nodeIDs = append(nodeIDs, id)
	}
	sort.Strings(nodeIDs)

	var strongConnect func(v string)
	strongConnect = func(v string) {
		indices[v] = index
		lowlink[v] = index
		index++
		stack = append(stack, v)
		onStack[v] = true

		for _, w := range g.outNeighbors[v] {
			if _, ok := indices[w.ID]; !ok {
				strongConnect(w.ID)
				if lowlink[w.ID] < lowlink[v] {
					lowlink[v] = lowlink[w.ID]
				}
			} else if onStack[w.ID] {
				if indices[w.ID] < lowlink[v] {
					lowlink[v] = indices[w.ID]
				}
			}
		}

		if lowlink[v] == indices[v] {
			var scc []string
			for {
				w := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				onStack[w] = false
				scc = append(scc, w)
				if w == v {
					break
				}
			}
			// Sort the SCC to make results deterministic
			sort.Strings(scc)
			sccs = append(sccs, scc)
		}
	}

	for _, nodeID := range nodeIDs {
		if _, ok := indices[nodeID]; !ok {
			strongConnect(nodeID)
		}
	}

	// Sort SCCs by their first element to make the overall result deterministic
	sort.Slice(sccs, func(i, j int) bool {
		if len(sccs[i]) == 0 {
			return true
		}
		if len(sccs[j]) == 0 {
			return false
		}
		return sccs[i][0] < sccs[j][0]
	})

	return sccs
}

// ToDAG collapses all Strongly Connected Components into a new Directed Acyclic Graph (DAG).
// Each node in the new DAG represents an SCC from the original graph.
func (g *Graph) ToDAG() *Graph {
	sccs := g.ComputeSCCs()
	dag := NewGraph()

	nodeToSCCIdx := make(map[string]int)
	for i, scc := range sccs {
		sccID := fmt.Sprintf("scc-%d", i)
		dag.AddNode(Node{
			ID:   sccID,
			Type: NodeSCC, // Representing SCC as a specialized SCC node
			Metadata: map[string]interface{}{
				"nodes": scc,
			},
		})
		for _, nodeID := range scc {
			nodeToSCCIdx[nodeID] = i
		}
	}

	g.mu.RLock()
	defer g.mu.RUnlock()

	// Use a map to track added edges to avoid duplicates in the DAG
	addedEdges := make(map[string]bool)

	for from, neighbors := range g.outNeighbors {
		fromSCCIdx := nodeToSCCIdx[from]
		for _, neighbor := range neighbors {
			toSCCIdx := nodeToSCCIdx[neighbor.ID]
			if fromSCCIdx != toSCCIdx {
				fromSCCID := fmt.Sprintf("scc-%d", fromSCCIdx)
				toSCCID := fmt.Sprintf("scc-%d", toSCCIdx)
				edgeKey := fmt.Sprintf("%s->%s", fromSCCID, toSCCID)
				if !addedEdges[edgeKey] {
					dag.AddEdge(Edge{
						From: fromSCCID,
						To:   toSCCID,
						Type: neighbor.Type,
					})
					addedEdges[edgeKey] = true
				}
			}
		}
	}

	return dag
}
