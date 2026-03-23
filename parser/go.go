package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type GoParser struct{}

func NewGoParser() *GoParser {
	return &GoParser{}
}

func (p *GoParser) Parse(path string, content string) (*FileMetadata, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, content, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	meta := &FileMetadata{
		Path:     path,
		Language: LangGo,
		Imports:  []string{},
		Symbols:  []Symbol{},
	}

	for _, imp := range f.Imports {
		meta.Imports = append(meta.Imports, strings.Trim(imp.Path.Value, `"`))
	}

	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			kind := KindFunction
			if x.Recv != nil {
				kind = KindMethod
			}
			meta.Symbols = append(meta.Symbols, Symbol{
				Name: x.Name.Name,
				Kind: kind,
				Span: Span{
					Start: fset.Position(x.Pos()).Line,
					End:   fset.Position(x.End()).Line,
				},
			})
		case *ast.GenDecl:
			for _, spec := range x.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					kind := KindType
					if _, ok := s.Type.(*ast.StructType); ok {
						kind = KindClass // Mapping struct to Class for now
					} else if _, ok := s.Type.(*ast.InterfaceType); ok {
						kind = KindInterface
					}
					meta.Symbols = append(meta.Symbols, Symbol{
						Name: s.Name.Name,
						Kind: kind,
						Span: Span{
							Start: fset.Position(s.Pos()).Line,
							End:   fset.Position(s.End()).Line,
						},
					})
				}
			}
		}
		return true
	})

	return meta, nil
}
