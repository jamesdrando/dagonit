package main

import (
	"dagonit/graph"
	"dagonit/indexer"
	"dagonit/planner"
	"dagonit/retrieval"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dagonit",
	Short: "Dagonit is a codebase dependency analyzer and planner",
}

var indexCmd = &cobra.Command{
	Use:   "index <path>",
	Short: "Scans, parses, and builds the graph",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		idx := indexer.NewIndexer(path)
		fmt.Printf("Indexing %s...\n", path)
		if err := idx.Index(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		data, err := idx.Graph.ExportToJSON()
		if err != nil {
			fmt.Printf("Error exporting: %v\n", err)
			return
		}
		if err := os.WriteFile("dagonit.json", data, 0644); err != nil {
			fmt.Printf("Error writing file: %v\n", err)
			return
		}
		fmt.Println("Graph saved to dagonit.json")
	},
}

func loadGraph() (*graph.Graph, error) {
	data, err := os.ReadFile("dagonit.json")
	if err != nil {
		return nil, err
	}
	g := graph.NewGraph()
	if err := g.ImportFromJSON(data); err != nil {
		return nil, err
	}
	return g, nil
}

var queryCmd = &cobra.Command{
	Use:   "query <term>",
	Short: "Finds nodes by name/ID and shows their neighbors",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		term := args[0]
		g, err := loadGraph()
		if err != nil {
			fmt.Printf("Error loading graph: %v\n", err)
			return
		}
		engine := retrieval.NewEngine(g)
		nodes := engine.FindNodesByTerm(term)
		if len(nodes) == 0 {
			fmt.Println("No nodes found.")
			return
		}
		for _, n := range nodes {
			fmt.Printf("Node: %s [%s]\n", n.ID, n.Type)
			
			// Get all nodes to reconstruct edge types for the output
			// (Since the current Graph doesn't store EdgeType in neighbors maps)
			// Wait, the AddEdge stores them in maps but without type.
			// Let's just iterate through nodes and show relationships.
			
			neighbors := g.GetNeighbors(n.ID)
			if len(neighbors) > 0 {
				fmt.Println("  Out-Neighbors:")
				for _, m := range neighbors {
					fmt.Printf("    -> %s [%s]\n", m.ID, m.Type)
				}
			}
			inNeighbors := g.GetInNeighbors(n.ID)
			if len(inNeighbors) > 0 {
				fmt.Println("  In-Neighbors:")
				for _, m := range inNeighbors {
					fmt.Printf("    <- %s [%s]\n", m.ID, m.Type)
				}
			}
		}
	},
}

var planCmd = &cobra.Command{
	Use:   "plan <seeds...>",
	Short: "Generates a plan from seed nodes and prints it",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		g, err := loadGraph()
		if err != nil {
			fmt.Printf("Error loading graph: %v\n", err)
			return
		}
		plan, err := planner.GeneratePlan(g, args)
		if err != nil {
			fmt.Printf("Error generating plan: %v\n", err)
			return
		}
		fmt.Println("Execution Plan:")
		for i, step := range plan.OrderedSteps {
			fmt.Printf("Step %d: %s\n", i+1, strings.Join(step, ", "))
		}
		if len(plan.Risks) > 0 {
			fmt.Println("\nRisks:")
			for _, r := range plan.Risks {
				fmt.Printf("- %s\n", r)
			}
		}
	},
}

var visualizeCmd = &cobra.Command{
	Use:   "visualize [mode]",
	Short: "Exports the graph to DOT or Mermaid format",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mode := "full"
		if len(args) > 0 {
			mode = args[0]
		}

		g, err := loadGraph()
		if err != nil {
			fmt.Printf("Error loading graph: %v\n", err)
			return
		}

		var format string
		format, _ = cmd.Flags().GetString("format")

		var outputGraph *graph.Graph
		switch mode {
		case "full":
			outputGraph = g
		case "dag":
			outputGraph = g.ToDAG()
		case "plan":
			seeds, _ := cmd.Flags().GetStringSlice("seeds")
			if len(seeds) == 0 {
				fmt.Println("Error: --seeds required for plan mode")
				return
			}
			p, err := planner.GeneratePlan(g, seeds)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			// Highlight nodes in the plan
			for _, nodeID := range p.AffectedNodes {
				if n, ok := g.GetNode(nodeID); ok {
					if n.Metadata == nil {
						n.Metadata = make(map[string]interface{})
					}
					n.Metadata["highlighted"] = true
					g.AddNode(n)
				}
			}
			outputGraph = g
		default:
			fmt.Printf("Unknown mode: %s\n", mode)
			return
		}

		var output string
		if format == "mermaid" {
			output = graph.ToMermaid(outputGraph)
		} else {
			output = graph.ToDOT(outputGraph)
		}
		fmt.Println(output)
	},
}

func main() {
	visualizeCmd.Flags().String("format", "dot", "Output format: dot or mermaid")
	visualizeCmd.Flags().StringSlice("seeds", []string{}, "Seeds for plan mode")

	rootCmd.AddCommand(indexCmd, queryCmd, planCmd, visualizeCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
