package indexer

import (
	"dagonit/graph"
	"dagonit/parser"
	"fmt"
	"os"
	"strings"
)

// Indexer orchestrates the scanning, parsing, and graph construction.
type Indexer struct {
	Scanner *Scanner
	Graph   *graph.Graph
	Parsers map[parser.Language]parser.Parser
}

func NewIndexer(root string) *Indexer {
	return &Indexer{
		Scanner: NewScanner(root),
		Graph:   graph.NewGraph(),
		Parsers: map[parser.Language]parser.Parser{
			parser.LangGo:         parser.NewGoParser(),
			parser.LangTypeScript: parser.NewTSJSParser(),
			parser.LangJavaScript: parser.NewTSJSParser(),
			parser.LangPython:     parser.NewPythonParser(),
		},
	}
}

func (idx *Indexer) Index() error {
	files, err := idx.Scanner.Scan()
	if err != nil {
		return err
	}

	fileNodes := make(map[string]graph.Node)
	symbolNodes := make(map[string]graph.Node)

	// Phase 1: Create File, Symbol, and Build Nodes
	for _, f := range files {
		if f.Type == FileBuild {
			buildNode := graph.Node{
				ID:   f.Path,
				Type: graph.NodeBuildTarget,
				Metadata: map[string]interface{}{
					"path": f.Path,
				},
			}
			idx.Graph.AddNode(buildNode)
			continue
		}

		content, err := os.ReadFile(f.Path)
		if err != nil {
			fmt.Printf("Warning: could not read file %s: %v\n", f.Path, err)
			continue
		}

		p, ok := idx.Parsers[f.Language]
		if !ok {
			continue
		}

		meta, err := p.Parse(f.Path, string(content))
		if err != nil {
			fmt.Printf("Warning: could not parse file %s: %v\n", f.Path, err)
			continue
		}

		nodeType := graph.NodeFile
		if f.IsTest {
			nodeType = graph.NodeTest
		}

		fileNode := graph.Node{
			ID:   f.Path,
			Type: nodeType,
			Metadata: map[string]interface{}{
				"language": f.Language,
				"imports":  meta.Imports,
				"isTest":   f.IsTest,
			},
		}
		idx.Graph.AddNode(fileNode)
		fileNodes[f.Path] = fileNode

		for _, s := range meta.Symbols {
			symbolID := fmt.Sprintf("%s#%s", f.Path, s.Name)
			symType := graph.NodeType(s.Kind)
			if f.IsTest {
				symType = graph.NodeTest
			}

			symbolNode := graph.Node{
				ID:   symbolID,
				Type: symType,
				Metadata: map[string]interface{}{
					"name":  s.Name,
					"kind":  s.Kind,
					"span":  s.Span,
					"file":  f.Path,
				},
			}
			idx.Graph.AddNode(symbolNode)
			symbolNodes[symbolID] = symbolNode

			// Containment edge
			idx.Graph.AddEdge(graph.Edge{
				From: fileNode.ID,
				To:   symbolNode.ID,
				Type: graph.EdgeContains,
			})
		}
	}

	// Phase 2: Create Import Edges and Validation Mappings
	for _, fNode := range fileNodes {
		imports := fNode.Metadata["imports"].([]string)
		for _, imp := range imports {
			if target, ok := fileNodes[imp]; ok {
				idx.Graph.AddEdge(graph.Edge{
					From: fNode.ID,
					To:   target.ID,
					Type: graph.EdgeImports,
				})
			}
		}

		// Heuristic for tests: Link test file to source file
		if fNode.Type == graph.NodeTest {
			// e.g., scanner_test.go -> scanner.go
			possibleSource := strings.Replace(fNode.ID, "_test.go", ".go", 1)
			possibleSource = strings.Replace(possibleSource, ".test.ts", ".ts", 1)
			possibleSource = strings.Replace(possibleSource, ".spec.ts", ".ts", 1)
			possibleSource = strings.Replace(possibleSource, ".test.js", ".js", 1)
			possibleSource = strings.Replace(possibleSource, ".spec.js", ".js", 1)

			if _, ok := fileNodes[possibleSource]; ok {
				idx.Graph.AddEdge(graph.Edge{
					From: possibleSource,
					To:   fNode.ID,
					Type: graph.EdgeTestedBy,
				})
			}
		}
	}

	// Heuristic for symbol tests: Link TestFuncName to FuncName
	for _, sNode := range symbolNodes {
		if sNode.Type == graph.NodeTest {
			name := sNode.Metadata["name"].(string)
			if strings.HasPrefix(name, "Test") {
				targetName := strings.TrimPrefix(name, "Test")
				// Look for this name in other symbols
				for _, other := range symbolNodes {
					if other.Type != graph.NodeTest && other.Metadata["name"] == targetName {
						idx.Graph.AddEdge(graph.Edge{
							From: other.ID,
							To:   sNode.ID,
							Type: graph.EdgeTestedBy,
						})
					}
				}
			}
		}
	}

	return nil
}
