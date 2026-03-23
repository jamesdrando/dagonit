package parser

import (
	"testing"
)

func TestGoParser(t *testing.T) {
	content := `
package test
import "fmt"
import "os"

type MyStruct struct {}

func (m *MyStruct) MyMethod() {}

func MyFunction() {
	fmt.Println("hello")
}
`
	p := NewGoParser()
	meta, err := p.Parse("test.go", content)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(meta.Imports) != 2 {
		t.Errorf("Expected 2 imports, got %d", len(meta.Imports))
	}

	expectedSymbols := map[string]SymbolKind{
		"MyStruct":   KindClass,
		"MyMethod":   KindMethod,
		"MyFunction": KindFunction,
	}

	if len(meta.Symbols) != 3 {
		t.Errorf("Expected 3 symbols, got %d", len(meta.Symbols))
	}

	for _, s := range meta.Symbols {
		if kind, ok := expectedSymbols[s.Name]; !ok || kind != s.Kind {
			t.Errorf("Unexpected symbol: %s (%s)", s.Name, s.Kind)
		}
	}
}
