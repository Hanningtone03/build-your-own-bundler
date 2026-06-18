package bundler

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Hanningtone03/build-your-own-bundler/internal/resolver"
)

func Bundle(modules map[string]*resolver.Module, order []string, entryPath string) string {
	var output strings.Builder

	output.WriteString("(function() {\n")
	output.WriteString("  const __modules = {};\n")
	output.WriteString("  const __cache = {};\n\n")
	output.WriteString("  function __require(id) {\n")
	output.WriteString("    if (__cache[id]) return __cache[id].exports;\n")
	output.WriteString("    const module = { exports: {} };\n")
	output.WriteString("    __cache[id] = module;\n")
	output.WriteString("    __modules[id](module, module.exports, __require);\n")
	output.WriteString("    return module.exports;\n")
	output.WriteString("  }\n\n")

	for _, path := range order {
		module := modules[path]
		id := moduleId(path)

		output.WriteString(fmt.Sprintf("  __modules[%q] = function(module, exports, require) {\n", id))

		source := module.Source
		for _, dep := range module.Deps {
			depId := moduleId(dep)
			source = replaceImportWithRequire(source, dep, depId)
		}
		source = stripExportKeywords(source)
		source = addModuleExports(source)

		indented := indentLines(source, "    ")
		output.WriteString(indented)
		output.WriteString("\n  };\n\n")
	}

	entryId := moduleId(entryPath)
	output.WriteString(fmt.Sprintf("  __require(%q);\n", entryId))
	output.WriteString("})();\n")

	return output.String()
}

func moduleId(path string) string {
	path = strings.ReplaceAll(path, "\\", "/")
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		return path
	}
	return parts[len(parts)-1]
}

var importLineRegex = regexp.MustCompile(`^\s*import\s+\{([^}]+)\}\s+from\s+['"]([^'"]+)['"]`)

func replaceImportWithRequire(source, originalPath, moduleId string) string {
	lines := strings.Split(source, "\n")
	for i, line := range lines {
		m := importLineRegex.FindStringSubmatch(line)
		if m != nil && strings.HasSuffix(originalPath, strings.TrimPrefix(m[2], "./")) {
			names := strings.TrimSpace(m[1])
			lines[i] = fmt.Sprintf("const { %s } = require(%q);", names, moduleId)
		}
	}
	return strings.Join(lines, "\n")
}

func stripExportKeywords(source string) string {
	lines := strings.Split(source, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "export function") {
			lines[i] = strings.Replace(line, "export function", "function", 1)
		} else if strings.HasPrefix(trimmed, "export default function") {
			lines[i] = strings.Replace(line, "export default function", "function", 1)
		} else if strings.HasPrefix(trimmed, "export const") {
			lines[i] = strings.Replace(line, "export const", "const", 1)
		} else if strings.HasPrefix(trimmed, "export class") {
			lines[i] = strings.Replace(line, "export class", "class", 1)
		}
	}
	return strings.Join(lines, "\n")
}

func addModuleExports(source string) string {
	exportRegex := regexp.MustCompile(`(?:^|\n)(?:function|const|class)\s+(\w+)`)
	matches := exportRegex.FindAllStringSubmatch(source, -1)

	var exportsCode strings.Builder
	for _, m := range matches {
		exportsCode.WriteString(fmt.Sprintf("\n  module.exports.%s = %s;", m[1], m[1]))
	}

	return source + exportsCode.String()
}

func sanitize(s string) string {
	s = strings.ReplaceAll(s, ".", "_")
	s = strings.ReplaceAll(s, "-", "_")
	return s
}

func indentLines(source, indent string) string {
	lines := strings.Split(source, "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) != "" {
			lines[i] = indent + line
		}
	}
	return strings.Join(lines, "\n")
}
