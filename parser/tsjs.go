package parser

import (
	"regexp"
	"strings"
)

type TSJSParser struct{}

func NewTSJSParser() *TSJSParser {
	return &TSJSParser{}
}

var (
	tsjsImportRegex   = regexp.MustCompile(`import\s+.*\s+from\s+['"](.+)['"]`)
	tsjsFunctionRegex = regexp.MustCompile(`function\s+(\w+)`)
	tsjsClassRegex    = regexp.MustCompile(`class\s+(\w+)`)
	tsjsConstFnRegex  = regexp.MustCompile(`const\s+(\w+)\s*=\s*(?:\([^)]*\)|[\w]+)\s*=>`)
)

func (p *TSJSParser) Parse(path string, content string) (*FileMetadata, error) {
	meta := &FileMetadata{
		Path:     path,
		Language: LangTypeScript, // Or JavaScript, based on extension
		Imports:  []string{},
		Symbols:  []Symbol{},
	}

	lines := strings.Split(content, "\n")

	for i, line := range lines {
		// Imports
		if matches := tsjsImportRegex.FindStringSubmatch(line); len(matches) > 1 {
			meta.Imports = append(meta.Imports, matches[1])
		}

		// Functions
		if matches := tsjsFunctionRegex.FindStringSubmatch(line); len(matches) > 1 {
			meta.Symbols = append(meta.Symbols, Symbol{
				Name: matches[1],
				Kind: KindFunction,
				Span: Span{Start: i + 1, End: i + 1}, // Simplification for v1
			})
		} else if matches := tsjsConstFnRegex.FindStringSubmatch(line); len(matches) > 1 {
			meta.Symbols = append(meta.Symbols, Symbol{
				Name: matches[1],
				Kind: KindFunction,
				Span: Span{Start: i + 1, End: i + 1},
			})
		}

		// Classes
		if matches := tsjsClassRegex.FindStringSubmatch(line); len(matches) > 1 {
			meta.Symbols = append(meta.Symbols, Symbol{
				Name: matches[1],
				Kind: KindClass,
				Span: Span{Start: i + 1, End: i + 1},
			})
		}
	}

	return meta, nil
}
