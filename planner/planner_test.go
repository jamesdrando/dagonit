package planner

import (
	"dagonit/graph"
	"reflect"
	"sort"
	"testing"
)

func TestGeneratePlan(t *testing.T) {
	g := graph.NewGraph()

	// Setup a graph:
	// A -> B -> C
	// A -> D
	// E -> F
	// G -> H -> G (Cycle)
	// I (High fan-in target)
	// J, K, L, M, N, O -> I

	nodes := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O"}
	for _, id := range nodes {
		g.AddNode(graph.Node{ID: id, Type: graph.NodeSymbol})
	}

	g.AddEdge(graph.Edge{From: "A", To: "B", Type: graph.EdgeImports})
	g.AddEdge(graph.Edge{From: "B", To: "C", Type: graph.EdgeImports})
	g.AddEdge(graph.Edge{From: "A", To: "D", Type: graph.EdgeImports})
	g.AddEdge(graph.Edge{From: "E", To: "F", Type: graph.EdgeImports})
	g.AddEdge(graph.Edge{From: "G", To: "H", Type: graph.EdgeImports})
	g.AddEdge(graph.Edge{From: "H", To: "G", Type: graph.EdgeImports})
	g.AddEdge(graph.Edge{From: "J", To: "I", Type: graph.EdgeImports})
	g.AddEdge(graph.Edge{From: "K", To: "I", Type: graph.EdgeImports})
	g.AddEdge(graph.Edge{From: "L", To: "I", Type: graph.EdgeImports})
	g.AddEdge(graph.Edge{From: "M", To: "I", Type: graph.EdgeImports})
	g.AddEdge(graph.Edge{From: "N", To: "I", Type: graph.EdgeImports})
	g.AddEdge(graph.Edge{From: "O", To: "I", Type: graph.EdgeImports})

	t.Run("Plan for A", func(t *testing.T) {
		plan, err := GeneratePlan(g, []string{"A"})
		if err != nil {
			t.Fatalf("Failed to generate plan: %v", err)
		}

		// A depends on B and D. B depends on C.
		// Expected subgraph nodes: A, B, C, D
		expectedAffected := []string{"A", "B", "C", "D"}
		sort.Strings(plan.AffectedNodes)
		if !reflect.DeepEqual(plan.AffectedNodes, expectedAffected) {
			t.Errorf("Expected affected nodes %v, got %v", expectedAffected, plan.AffectedNodes)
		}

		// Expected order:
		// Step 0: C, D (no dependencies)
		// Step 1: B (depends on C)
		// Step 2: A (depends on B and D)
		if len(plan.OrderedSteps) != 3 {
			t.Errorf("Expected 3 steps, got %d", len(plan.OrderedSteps))
		}
		
		// Note: Steps are grouped by levels.
		// Level 0: C, D
		// Level 1: B
		// Level 2: A
		
		step0 := plan.OrderedSteps[0]
		sort.Strings(step0)
		if !reflect.DeepEqual(step0, []string{"C", "D"}) {
			t.Errorf("Step 0 mismatch: expected [C D], got %v", step0)
		}
		
		if !reflect.DeepEqual(plan.OrderedSteps[1], []string{"B"}) {
			t.Errorf("Step 1 mismatch: expected [B], got %v", plan.OrderedSteps[1])
		}
		
		if !reflect.DeepEqual(plan.OrderedSteps[2], []string{"A"}) {
			t.Errorf("Step 2 mismatch: expected [A], got %v", plan.OrderedSteps[2])
		}
	})

	t.Run("Plan for G (Cycle)", func(t *testing.T) {
		plan, err := GeneratePlan(g, []string{"G"})
		if err != nil {
			t.Fatalf("Failed to generate plan: %v", err)
		}

		// G -> H -> G. Subgraph nodes: G, H. They form an SCC.
		expectedAffected := []string{"G", "H"}
		sort.Strings(plan.AffectedNodes)
		if !reflect.DeepEqual(plan.AffectedNodes, expectedAffected) {
			t.Errorf("Expected affected nodes %v, got %v", expectedAffected, plan.AffectedNodes)
		}

		// SCC should be one step.
		if len(plan.OrderedSteps) != 1 {
			t.Errorf("Expected 1 step for cycle, got %d", len(plan.OrderedSteps))
		}
		
		step0 := plan.OrderedSteps[0]
		sort.Strings(step0)
		if !reflect.DeepEqual(step0, []string{"G", "H"}) {
			t.Errorf("Step 0 mismatch: expected [G H], got %v", step0)
		}

		foundCycleRisk := false
		for _, risk := range plan.Risks {
			if risk == "Cycle detected: nodes [G H] form an SCC" {
				foundCycleRisk = true
				break
			}
		}
		if !foundCycleRisk {
			t.Errorf("Expected cycle risk not found in %v", plan.Risks)
		}
	})

	t.Run("Plan for J (High Fan-in)", func(t *testing.T) {
		plan, err := GeneratePlan(g, []string{"J"})
		if err != nil {
			t.Fatalf("Failed to generate plan: %v", err)
		}

		// J depends on I.
		expectedAffected := []string{"I", "J"}
		sort.Strings(plan.AffectedNodes)
		if !reflect.DeepEqual(plan.AffectedNodes, expectedAffected) {
			t.Errorf("Expected affected nodes %v, got %v", expectedAffected, plan.AffectedNodes)
		}

		foundHighFanInRisk := false
		for _, risk := range plan.Risks {
			if risk == "High risk: node I has high fan-in (6)" {
				foundHighFanInRisk = true
				break
			}
		}
		if !foundHighFanInRisk {
			t.Errorf("Expected high fan-in risk not found in %v", plan.Risks)
		}
	})
}
