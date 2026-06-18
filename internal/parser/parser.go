package parser

import (
	"regexp"
)

type ImportInfo struct {
	Path  string
	Start int
	End   int
}

var importRegex = regexp.MustCompile(`import\s+(?:[\w{},\s*]+\s+from\s+)?['"]([^'"]+)['"]`)
var requireRegex = regexp.MustCompile(`require\(['"]([^'"]+)['"]\)`)
var exportRegex = regexp.MustCompile(`export\s+(default\s+)?(function|class|const|let|var)\s+(\w+)`)

func FindImports(source string) []ImportInfo {
	var imports []ImportInfo

	matches := importRegex.FindAllStringSubmatchIndex(source, -1)
	for _, m := range matches {
		imports = append(imports, ImportInfo{
			Path:  source[m[2]:m[3]],
			Start: m[0],
			End:   m[1],
		})
	}

	requireMatches := requireRegex.FindAllStringSubmatchIndex(source, -1)
	for _, m := range requireMatches {
		imports = append(imports, ImportInfo{
			Path:  source[m[2]:m[3]],
			Start: m[0],
			End:   m[1],
		})
	}

	return imports
}

func FindExports(source string) []string {
	var exports []string
	matches := exportRegex.FindAllStringSubmatch(source, -1)
	for _, m := range matches {
		exports = append(exports, m[3])
	}
	return exports
}

func StripImports(source string) string {
	source = importRegex.ReplaceAllString(source, "")
	return source
}
