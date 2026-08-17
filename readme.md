# Xavi

Xavi is a small Go-based language project with a working end-to-end path from `.xavi` source files to bytecode execution.

Today the repository can:

- tokenize indentation-based source files
- parse imports, functions, statements, arithmetic, and function calls
- load imported `.xavi` files recursively
- compile AST nodes into bytecode
- execute that bytecode in a stack-based virtual machine
- call simple builtins such as `pt` and `ptln`

## Repository Layout

```text
xavi/
|-- cmd/
|   `-- xavic/                # Reserved for a future dedicated CLI
|-- compiler/
|   |-- ast/                  # Program, statement, and expression nodes
|   |-- bytecode/             # AST -> bytecode generator
|   |-- emitter/              # Placeholder
|   |-- ir/                   # Minimal IR type
|   |-- lexer/                # Tokenizer with indent/dedent support
|   |-- parser/               # Recursive-descent parser
|   |-- sema/                 # Placeholder
|   `-- toon/                 # Placeholder
|-- docs/
|   `-- architecture.md
|-- examples/
|   |-- hello.xavi            # Example entry program with imports
|   |-- sum.xavi
|   `-- operation/
|       `-- mult.xavi
|-- vm/
|   |-- agents/               # Minimal background agent loop abstraction
|   |-- builtin/              # Builtin function dispatch
|   |-- exec/                 # Bytecode interpreter
|   |-- loader/               # Source loading + import resolution
|   |-- opcode/               # Opcode and builtin indexes
|   `-- runtime/              # Stack and frame storage
|-- go.mod
|-- LICENSE
|-- main.go                   # Current CLI entrypoint
`-- readme.md
```

## Current Execution Flow

The current entrypoint in [`main.go`](./main.go) expects a Xavi source file:

```bash
go run . examples/hello.xavi
```

That flow is:

1. `vm/loader` reads the entry file and resolves imports.
2. `compiler/lexer` tokenizes the source.
3. `compiler/parser` builds an AST program.
4. `compiler/bytecode` compiles functions into VM bytecode.
5. `vm/exec` runs the compiled `main` function.

## Language Features Implemented Today

### Top-Level Forms

- `import "./file.xavi"`
- `import name < "./file.xavi"` to import one named function
- `fn main():`

### Statements

- `let name = expr`
- `return expr`
- expression statements such as builtin calls

### Expressions

- identifiers
- number literals
- string literals
- binary arithmetic: `+`, `-`, `*`, `/`
- grouped expressions with parentheses
- function calls

### Runtime Features

- local variables stored in indexed frame slots
- per-function constant pools
- user-defined function calls
- builtins `pt(...)` and `ptln(...)`

## Example Program

[`examples/hello.xavi`](./examples/hello.xavi) demonstrates:

- recursive imports
- nested function calls
- arithmetic
- printing through builtins

Running it currently prints:

```text
Hello world20 hello new line
Result: 30
```

The formatting is expected with the current builtins:

- `pt(...)` prints arguments without a trailing newline
- `ptln(...)` prints arguments and then appends a newline

## Package Status

### Implemented

- [`compiler/lexer/lexer.go`](./compiler/lexer/lexer.go): tokenization, keywords, strings, numbers, indentation tokens
- [`compiler/parser/parser.go`](./compiler/parser/parser.go): program, import, function, block, statement, and expression parsing
- [`compiler/ast/nodes.go`](./compiler/ast/nodes.go): program, statement, and expression node definitions
- [`compiler/bytecode/gen.go`](./compiler/bytecode/gen.go): AST lowering to compiled functions
- [`vm/loader/loader.go`](./vm/loader/loader.go): import-aware source loading
- [`vm/exec/interpreter.go`](./vm/exec/interpreter.go): bytecode interpreter with function and builtin calls
- [`vm/runtime`](./vm/runtime): stack and frame storage
- [`vm/opcode/opcode.go`](./vm/opcode/opcode.go): opcode definitions and builtin indexes
- [`vm/builtin/print.go`](./vm/builtin/print.go): console-print builtins

### Partial Or Placeholder

- `cmd/xavic/` exists but does not yet provide a dedicated CLI binary
- `compiler/ir/` defines a minimal instruction struct but is not wired into compilation
- `compiler/emitter/`, `compiler/sema/`, and `compiler/toon/` are present but currently empty
- `vm/agents/` contains a minimal repeated callback loop and is not integrated with the language pipeline

## Build And Verification

From the project root:

```bash
go test ./...
go run . examples/hello.xavi
```

As of August 17, 2026:

- `go test ./...` succeeds
- there are currently no package test files, so the command only verifies that the code builds cleanly

## Current Limitations

- no semantic analysis pass yet
- no type checking
- parser errors still panic instead of returning structured diagnostics
- bytecode generation panics on unsupported or undefined references
- interpreter arithmetic assumes `float64` values
- runtime containers do not perform defensive bounds checks
- imports merge functions by name only and do not provide module namespaces
- builtin support is limited to printing

## Suggested Next Steps

1. Add parser, loader, bytecode, and VM tests around the example workflows.
2. Introduce structured diagnostics instead of `panic`-driven failures.
3. Add semantic analysis for symbol resolution, duplicate definitions, and type validation.
4. Decide whether `compiler/ir` will remain a real intermediate stage or be removed.
5. Build a dedicated `cmd/xavic` CLI for run/build/check commands.

## Additional Documentation

For a package-by-package walkthrough, see [docs/architecture.md](./docs/architecture.md).
