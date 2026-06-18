package main

import (
	"fmt"
	"os"

	"github.com/Hanningtone03/build-your-own-bundler/internal/bundler"
	"github.com/Hanningtone03/build-your-own-bundler/internal/minifier"
	"github.com/Hanningtone03/build-your-own-bundler/internal/resolver"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("Usage: go run main.go <entry-file> [output-file]")
		os.Exit(1)
	}

	entryPath := args[0]
	outputPath := "bundle.js"
	if len(args) >= 2 {
		outputPath = args[1]
	}

	fmt.Printf("\n  Bundling from: %s\n", entryPath)

	modules, order, err := resolver.Resolve(entryPath)
	if err != nil {
		fmt.Printf("  Error: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("  Found %d module(s):\n", len(order))
	for i, path := range order {
		fmt.Printf("    %d. %s\n", i+1, path)
	}

	bundled := bundler.Bundle(modules, order, modules[order[len(order)-1]].Path)
	minified := minifier.Minify(bundled)

	origSize, minSize, reduction := minifier.Stats(bundled, minified)

	err = os.WriteFile(outputPath, []byte(minified), 0644)
	if err != nil {
		fmt.Printf("  Error writing output: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n  Bundle written to: %s\n", outputPath)
	fmt.Printf("  Original size: %d bytes\n", origSize)
	fmt.Printf("  Minified size: %d bytes\n", minSize)
	fmt.Printf("  Reduction: %.1f%%\n\n", reduction)
}
