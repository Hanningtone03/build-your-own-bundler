![CI](https://github.com/Hanningtone03/build-your-own-bundler/actions/workflows/ci.yml/badge.svg)

# Build Your Own Bundler

A JavaScript module bundler in Go; resolves imports, bundles modules in dependency order, and minifies the output.

## How it works

Starting from an entry file, the resolver walks every import statement, finds the matching file, and builds a dependency graph. The bundler wraps each module in a function and stitches them together with a small runtime that mimics CommonJS `require`. A minifier strips comments and collapses whitespace before writing the final file.

## Project structure

```
main.go
internal/
├── resolver/
│   └── resolver.go
├── parser/
│   └── parser.go
├── bundler/
│   └── bundler.go
└── minifier/
    └── minifier.go
```

## Running locally

```bash
go run main.go examples/main.js bundle.js
node bundle.js
```

## Example

```javascript
// math.js
export function add(a, b) {
  return a + b;
}

// main.js
import { add } from './math.js';
console.log(add(2, 3));
```

Bundles into a single self-contained file with no external loader needed.

## Tech

- Go
- `regexp` for import/export detection
- No external dependencies
