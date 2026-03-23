package planner

import (
	"dagonit/graph"
	"fmt"
	"sort"
)

// Plan holds the details of an execution strategy.
type Plan struct {
	// OrderedSteps is a list of parallel groups of node IDs.
	// Each inner slice contains nodes that can be executed in parallel.
	OrderedSteps [][]string
	// AffectedNodes is the list of all nodes involved in the plan.
	AffectedNodes []string
	// Risks identifies potential issues in the plan, such as high fan-in nodes or cycles.
	Risks []string
}

// GeneratePlan creates an execution plan based on seed node IDs by expanding to dependencies.
func GeneratePlan(g *graph.Graph, seeds []string) (*Plan, error) {
	// 1. Expand subgraph to include dependencies (transitive out-neighbors)
	visited := make(map[string]bool)
	for _, seed := range seeds {
		expand(g, seed, visited)
	}

	// 2. Create the subgraph
	sub := graph.NewGraph()
	var affectedNodes []string
	for id := range visited {
		node, ok := g.GetNode(id)
		if ok {
			sub.AddNode(node)
			affectedNodes = append(affectedNodes, id)
		}
	}
	sort.Strings(affectedNodes)

	// Add edges in the subgraph
	for id := range visited {
		neighbors := g.GetNeighbors(id)
		for _, neighbor := range neighbors {
			if visited[neighbor.ID] {
				sub.AddEdge(graph.Edge{
					From: id,
					To:   neighbor.ID,
					Type: neighbor.Type,
				})
			}
		}
	}

	// 3. Collapse to DAG using SCC logic
	dag := sub.ToDAG()

	// 4. Calculate Levels for Topological Sort
	// Level 0: nodes with NO dependencies (out-degree 0 in the dependency graph)
	// Level N: nodes whose dependencies are all in Level N-1 or below.
	levels := make(map[string]int)
	nodes := dag.GetNodes()

	// Sort node IDs for deterministic traversal
	nodeIDs := make([]string, 0, len(nodes))
	for _, node := range nodes {
		nodeIDs = append(nodeIDs, node.ID)
	}
	sort.Strings(nodeIDs)

	var getLevel func(id string) (int, error)
	stack := make(map[string]bool)

	getLevel = func(id string) (int, error) {
		if l, ok := levels[id]; ok {
			return l, nil
		}
		if stack[id] {
			return 0, fmt.Errorf("cycle detected in DAG for node %s", id)
		}
		stack[id] = true
		defer func() { delete(stack, id) }()

		neighbors := dag.GetNeighbors(id)
		if len(neighbors) == 0 {
			levels[id] = 0
			return 0, nil
		}

		maxL := -1
		for _, neighbor := range neighbors {
			l, err := getLevel(neighbor.ID)
			if err != nil {
				return 0, err
			}
			if l > maxL {
				maxL = l
			}
		}
		levels[id] = maxL + 1
		return maxL + 1, nil
	}

	for _, id := range nodeIDs {
		if _, err := getLevel(id); err != nil {
			return nil, err
		}
	}

	// Group SCCs by level
	levelGroups := make(map[int][]string)
	maxLevel := -1
	for id, l := range levels {
		levelGroups[l] = append(levelGroups[l], id)
		if l > maxLevel {
			maxLevel = l
		}
	}

	// 5. Build OrderedSteps and Identify Risks
	orderedSteps := make([][]string, maxLevel+1)
	var risks []string

	for i := 0; i <= maxLevel; i++ {
		sccIDs := levelGroups[i]
		sort.Strings(sccIDs)

		var step []string
		for _, sccID := range sccIDs {
			sccNode, _ := dag.GetNode(sccID)
			nodesInSCC, ok := sccNode.Metadata["nodes"].([]string)
			if ok {
				if len(nodesInSCC) > 1 {
					risks = append(risks, fmt.Sprintf("Cycle detected: nodes %v form an SCC", nodesInSCC))
				}
				sort.Strings(nodesInSCC)
				step = append(step, nodesInSCC...)
			}
		}
		orderedSteps[i] = step
	}

	// Additional Risk: High fan-in in original graph
	for _, id := range affectedNodes {
		inNeighbors := g.GetInNeighbors(id)
		if len(inNeighbors) > 5 {
			risks = append(risks, fmt.Sprintf("High risk: node %s has high fan-in (%d)", id, len(inNeighbors)))
		}
	}
	sort.Strings(risks)

	return &Plan{
		OrderedSteps:  orderedSteps,
		AffectedNodes: affectedNodes,
		Risks:         risks,
	}, nil
}

// expand performs a depth-first search to find all reachable nodes following dependencies.
func expand(g *graph.Graph, id string, visited map[string]bool) {
	if visited[id] {
		return
	}
	visited[id] = true
	neighbors := g.GetNeighbors(id)
	for _, neighbor := range neighbors {
		expand(g, neighbor.ID, visited)
	}
}
