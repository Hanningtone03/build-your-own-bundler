package resolver

import (
	"fmt"
	"os"
	"path/filepath"
)

type Module struct {
	Path   string
	Source string
	Deps   []string
}

func Resolve(entryPath string) (map[string]*Module, []string, error) {
	modules := make(map[string]*Module)
	order := []string{}

	var visit func(path string) error
	visit = func(path string) error {
		absPath, err := resolveExtension(path)
		if err != nil {
			return err
		}

		if _, exists := modules[absPath]; exists {
			return nil
		}

		content, err := os.ReadFile(absPath)
		if err != nil {
			return fmt.Errorf("could not read %s: %w", absPath, err)
		}

		module := &Module{
			Path:   absPath,
			Source: string(content),
		}
		modules[absPath] = module

		imports := extractImportPaths(string(content))
		dir := filepath.Dir(absPath)

		for _, imp := range imports {
			if isRelative(imp) {
				depPath := filepath.Join(dir, imp)
				resolvedDep, err := resolveExtension(depPath)
				if err != nil {
					return err
				}
				module.Deps = append(module.Deps, resolvedDep)
				if err := visit(resolvedDep); err != nil {
					return err
				}
			}
		}

		order = append(order, absPath)
		return nil
	}

	if err := visit(entryPath); err != nil {
		return nil, nil, err
	}

	return modules, order, nil
}

func resolveExtension(path string) (string, error) {
	if fileExists(path) {
		return filepath.Abs(path)
	}
	for _, ext := range []string{".js", ".mjs"} {
		candidate := path + ext
		if fileExists(candidate) {
			return filepath.Abs(candidate)
		}
	}
	return "", fmt.Errorf("module not found: %s", path)
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func isRelative(path string) bool {
	return len(path) > 0 && (path[0] == '.' || path[0] == '/')
}

func extractImportPaths(source string) []string {
	var paths []string
	for i := 0; i < len(source); i++ {
		if matchKeyword(source, i, "import") || matchKeyword(source, i, "require") {
			start := -1
			for j := i; j < len(source) && j < i+200; j++ {
				if source[j] == '\'' || source[j] == '"' {
					if start == -1 {
						start = j + 1
					} else {
						paths = append(paths, source[start:j])
						break
					}
				}
			}
		}
	}
	return paths
}

func matchKeyword(source string, pos int, keyword string) bool {
	if pos+len(keyword) > len(source) {
		return false
	}
	return source[pos:pos+len(keyword)] == keyword
}
