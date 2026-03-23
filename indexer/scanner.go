package indexer

import (
	"dagonit/parser"
	"os"
	"path/filepath"
	"strings"
)

// FileType indicates if a file is a source file, build file, or config file.
type FileType string

const (
	FileSource FileType = "source"
	FileBuild  FileType = "build"
	FileConfig FileType = "config"
)

// FileInfo contains basic information about a discovered file.
type FileInfo struct {
	Path     string
	Language parser.Language
	Type     FileType
	IsTest   bool
}

// Scanner handles recursive file discovery with ignore rules and language detection.
type Scanner struct {
	Root           string
	IgnoreDirs     map[string]bool
	Extensions     map[string]parser.Language
	BuildFileNames map[string]bool
	TestPatterns   []string
}

// NewScanner creates a new Scanner with default ignore rules and supported extensions.
func NewScanner(root string) *Scanner {
	return &Scanner{
		Root: root,
		IgnoreDirs: map[string]bool{
			".git":         true,
			"node_modules": true,
			"dist":         true,
			"build":        true,
			"vendor":       true,
			"venv":         true,
			"__pycache__":  true,
		},
		Extensions: map[string]parser.Language{
			".go":   parser.LangGo,
			".ts":   parser.LangTypeScript,
			".tsx":  parser.LangTypeScript,
			".js":   parser.LangJavaScript,
			".jsx":  parser.LangJavaScript,
			".py":   parser.LangPython,
		},
		BuildFileNames: map[string]bool{
			"go.mod":           true,
			"package.json":     true,
			"requirements.txt": true,
			"Makefile":         true,
			"Dockerfile":       true,
		},
		TestPatterns: []string{
			"_test.go",
			".test.ts",
			".spec.ts",
			".test.js",
			".spec.js",
			"test_", // Common for python
		},
	}
}

// Scan recursively walks the root directory and returns a list of discovered files.
func (s *Scanner) Scan() ([]FileInfo, error) {
	var files []FileInfo

	err := filepath.Walk(s.Root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if s.IgnoreDirs[info.Name()] {
				return filepath.SkipDir
			}
			return nil
		}

		relPath, err := filepath.Rel(s.Root, path)
		if err != nil {
			relPath = path
		}

		// Check for build files
		if s.BuildFileNames[info.Name()] {
			files = append(files, FileInfo{
				Path:     relPath,
				Language: parser.LangUnknown,
				Type:     FileBuild,
			})
			return nil
		}

		// Check for source files
		ext := filepath.Ext(path)
		if lang, ok := s.Extensions[ext]; ok {
			isTest := false
			for _, pattern := range s.TestPatterns {
				if strings.Contains(info.Name(), pattern) {
					isTest = true
					break
				}
			}

			files = append(files, FileInfo{
				Path:     relPath,
				Language: lang,
				Type:     FileSource,
				IsTest:   isTest,
			})
		}

		return nil
	})

	return files, err
}

