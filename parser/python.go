package parser

import (
	"regexp"
	"strings"
)

type PythonParser struct{}

func NewPythonParser() *PythonParser {
	return &PythonParser{}
}

var (
	pyImportRegex  = regexp.MustCompile(`(?:import\s+(\w+)|from\s+(\w+)\s+import)`)
	pyDefRegex     = regexp.MustCompile(`def\s+(\w+)\(`)
	pyClassRegex   = regexp.MustCompile(`class\s+(\w+)(?:\(.*\))?:`)
)

func (p *PythonParser) Parse(path string, content string) (*FileMetadata, error) {
	meta := &FileMetadata{
		Path:     path,
		Language: LangPython,
		Imports:  []string{},
		Symbols:  []Symbol{},
	}

	lines := strings.Split(content, "\n")

	for i, line := range lines {
		// Imports
		if matches := pyImportRegex.FindStringSubmatch(line); len(matches) > 0 {
			for _, m := range matches[1:] {
				if m != "" {
					meta.Imports = append(meta.Imports, m)
					break
				}
			}
		}

		// Functions
		if matches := pyDefRegex.FindStringSubmatch(line); len(matches) > 1 {
			meta.Symbols = append(meta.Symbols, Symbol{
				Name: matches[1],
				Kind: KindFunction,
				Span: Span{Start: i + 1, End: i + 1},
			})
		}

		// Classes
		if matches := pyClassRegex.FindStringSubmatch(line); len(matches) > 1 {
			meta.Symbols = append(meta.Symbols, Symbol{
				Name: matches[1],
				Kind: KindClass,
				Span: Span{Start: i + 1, End: i + 1},
			})
		}
	}

	return meta, nil
}
