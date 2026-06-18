package minifier

import (
	"regexp"
	"strings"
)

var singleLineComment = regexp.MustCompile(`//[^\n]*`)
var blockComment = regexp.MustCompile(`/\*[\s\S]*?\*/`)
var multipleSpaces = regexp.MustCompile(`[ \t]+`)
var multipleNewlines = regexp.MustCompile(`\n+`)

func Minify(source string) string {
	source = blockComment.ReplaceAllString(source, "")
	source = singleLineComment.ReplaceAllString(source, "")

	lines := strings.Split(source, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			trimmed = multipleSpaces.ReplaceAllString(trimmed, " ")
			result = append(result, trimmed)
		}
	}

	return strings.Join(result, "\n")
}

func Stats(original, minified string) (int, int, float64) {
	origSize := len(original)
	minSize := len(minified)
	reduction := 0.0
	if origSize > 0 {
		reduction = (1 - float64(minSize)/float64(origSize)) * 100
	}
	return origSize, minSize, reduction
}
