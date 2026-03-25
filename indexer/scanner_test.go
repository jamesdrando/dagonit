package indexer

import (
	"dagonit/parser"
	"os"
	"path/filepath"
	"testing"
)

func TestScanner_Scan(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "scanner_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	filesToCreate := []string{
		"main.go",
		"src/app.ts",
		"src/util.js",
		"python/script.py",
		"node_modules/library.js",
		".git/config",
		"build/main.o",
		"dist/bundle.js",
		"other.txt",
	}

	for _, f := range filesToCreate {
		path := filepath.Join(tempDir, f)
		err := os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory for %s: %v", f, err)
		}
		err = os.WriteFile(path, []byte(""), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", f, err)
		}
	}

	scanner := NewScanner(tempDir)
	files, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	expectedFiles := map[string]parser.Language{
		"main.go":          parser.LangGo,
		"src/app.ts":       parser.LangTypeScript,
		"src/util.js":      parser.LangJavaScript,
		"python/script.py": parser.LangPython,
	}

	if len(files) != len(expectedFiles) {
		t.Errorf("Expected %d files, got %d", len(expectedFiles), len(files))
	}

	for _, f := range files {
		lang, ok := expectedFiles[f.Path]
		if !ok {
			t.Errorf("Unexpected file discovered: %s", f.Path)
			continue
		}
		if f.Language != lang {
			t.Errorf("Expected language %s for file %s, got %s", lang, f.Path, f.Language)
		}
	}
}
