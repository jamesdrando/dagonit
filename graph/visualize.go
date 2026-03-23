package graph

import (
	"fmt"
	"strings"
)

// ToDOT exports the graph to DOT format for Graphviz.
func ToDOT(g *Graph) string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	var sb strings.Builder
	sb.WriteString("digraph G {\n")
	sb.WriteString("  rankdir=LR;\n")
	sb.WriteString("  node [shape=box, style=filled, fillcolor=lightgrey];\n")

	// Collect and sort node IDs for deterministic output
	nodeIDs := make([]string, 0, len(g.nodes))
	for id := range g.nodes {
		nodeIDs = append(nodeIDs, id)
	}
	// Sort nodeIDs or use a deterministic traversal if needed, but for visualization it's less critical
	// though good for testing.

	for _, id := range nodeIDs {
		n := g.nodes[id]
		label := id
		fillcolor := "lightgrey"
		if n.Type == NodeSCC {
			if nodes, ok := n.Metadata["nodes"].([]string); ok {
				label = fmt.Sprintf("SCC: %v", nodes)
			}
			fillcolor = "lightblue"
		} else if n.Type == NodeFile {
			fillcolor = "lightyellow"
		} else if n.Type == NodeSymbol {
			fillcolor = "white"
		}

		if highlighted, ok := n.Metadata["highlighted"].(bool); ok && highlighted {
			fillcolor = "orange"
		}

		sb.WriteString(fmt.Sprintf("  \"%s\" [label=\"%s\", fillcolor=\"%s\"];\n", id, label, fillcolor))
	}

	for from, neighbors := range g.outNeighbors {
		for _, to := range neighbors {
			sb.WriteString(fmt.Sprintf("  \"%s\" -> \"%s\" [label=\"%s\"];\n", from, to.ID, to.Type))
		}
	}

	sb.WriteString("}\n")
	return sb.String()
}

// ToMermaid exports the graph to Mermaid format.
func ToMermaid(g *Graph) string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	var sb strings.Builder
	sb.WriteString("graph TD\n")

	// Collect node IDs
	nodeIDs := make([]string, 0, len(g.nodes))
	for id := range g.nodes {
		nodeIDs = append(nodeIDs, id)
	}

	for _, id := range nodeIDs {
		n := g.nodes[id]
		label := id
		if n.Type == NodeSCC {
			if nodes, ok := n.Metadata["nodes"].([]string); ok {
				label = fmt.Sprintf("SCC: %v", nodes)
			}
		}

		// Mermaid uses different syntax for labels and styles
		// We'll use [ ] for square, ( ) for rounded, etc.
		// For simplicity, just use [ ] and apply classes later if needed.
		// But Mermaid labels can't easily contain special characters without quoting.
		safeID := strings.ReplaceAll(id, ".", "_")
		safeID = strings.ReplaceAll(safeID, "/", "_")
		safeID = strings.ReplaceAll(safeID, ":", "_")
		safeID = strings.ReplaceAll(safeID, "#", "_")
		safeID = strings.ReplaceAll(safeID, "-", "_")

		sb.WriteString(fmt.Sprintf("  %s[\"%s\"]\n", safeID, label))

		if highlighted, ok := n.Metadata["highlighted"].(bool); ok && highlighted {
			sb.WriteString(fmt.Sprintf("  style %s fill:#f96\n", safeID))
		} else if n.Type == NodeSCC {
			sb.WriteString(fmt.Sprintf("  style %s fill:#add8e6\n", safeID))
		} else if n.Type == NodeFile {
			sb.WriteString(fmt.Sprintf("  style %s fill:#ffffe0\n", safeID))
		}
	}

	for from, neighbors := range g.outNeighbors {
		safeFrom := strings.ReplaceAll(from, ".", "_")
		safeFrom = strings.ReplaceAll(safeFrom, "/", "_")
		safeFrom = strings.ReplaceAll(safeFrom, ":", "_")
		safeFrom = strings.ReplaceAll(safeFrom, "#", "_")
		safeFrom = strings.ReplaceAll(safeFrom, "-", "_")

		for _, to := range neighbors {
			safeTo := strings.ReplaceAll(to.ID, ".", "_")
			safeTo = strings.ReplaceAll(safeTo, "/", "_")
			safeTo = strings.ReplaceAll(safeTo, ":", "_")
			safeTo = strings.ReplaceAll(safeTo, "#", "_")
			safeTo = strings.ReplaceAll(safeTo, "-", "_")
			sb.WriteString(fmt.Sprintf("  %s -- %s --> %s\n", safeFrom, string(to.Type), safeTo))
		}
	}

	return sb.String()
}
