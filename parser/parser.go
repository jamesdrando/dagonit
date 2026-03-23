package parser

// Language represents a supported programming language.
type Language string

const (
	LangGo         Language = "go"
	LangTypeScript Language = "typescript"
	LangJavaScript Language = "javascript"
	LangPython     Language = "python"
	LangUnknown    Language = "unknown"
)

// SymbolKind represents the type of symbol (function, class, interface, etc.).
type SymbolKind string

const (
	KindFunction  SymbolKind = "function"
	KindClass     SymbolKind = "class"
	KindInterface SymbolKind = "interface"
	KindType      SymbolKind = "type"
	KindVariable  SymbolKind = "variable"
	KindMethod    SymbolKind = "method"
	KindModule    SymbolKind = "module"
)

// Span represents a range of lines in a file.
type Span struct {
	Start int
	End   int
}

// Symbol represents a code symbol (function, class, etc.).
type Symbol struct {
	Name      string
	Kind      SymbolKind
	Signature string
	Span      Span
}

// FileMetadata contains extracted information about a source file.
type FileMetadata struct {
	Path     string
	Language Language
	Imports  []string
	Symbols  []Symbol
}

// Parser defines the interface for language-specific parsers.
type Parser interface {
	Parse(path string, content string) (*FileMetadata, error)
}
