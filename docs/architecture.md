# Architecture Notes

## Overview

Xavi now has a working source-to-execution pipeline:

1. load a `.xavi` entry file
2. resolve imports
3. tokenize source text
4. parse an AST
5. compile AST functions into bytecode
6. execute bytecode in the VM

The architecture is still small and prototype-oriented, but it is no longer just a handwritten bytecode demo. The main runtime path is implemented and exercised by the example programs under `examples/`.

## End-To-End Flow

### 1. Entrypoint

[`main.go`](../main.go) acts as the current CLI entrypoint.

Responsibilities:

- require an input path such as `examples/hello.xavi`
- load the program graph through `vm/loader`
- compile it through `compiler/bytecode`
- execute `main` through `vm/exec`

If no argument is provided, the process exits with:

```text
Usage: go run . <file.xavi>
```

### 2. Loading And Import Resolution

[`vm/loader/loader.go`](../vm/loader/loader.go) is the bridge from source files into the compiler pipeline.

It currently:

- resolves the entry path to an absolute path
- reads source from disk
- tokenizes and parses each file
- follows imports relative to the importing file
- prevents repeated loading with a `visited` map
- merges imported functions ahead of the current file's functions

Supported import shapes:

- `import "./sum.xavi"`
- `import sum < "./sum.xavi"`

The named-import form selects a single function by name from the imported file. The unnamed form includes all imported functions.

### 3. Lexing

[`compiler/lexer/lexer.go`](../compiler/lexer/lexer.go) contains a real tokenizer rather than just token definitions.

Implemented token categories include:

- keywords: `fn`, `import`, `return`, `let`, `record`, `agent`, `on`, `event`
- punctuation/operators: `:`, `<`, `->`, `(`, `)`, `,`, `=`, `+`, `-`, `*`, `/`
- literals: identifiers, numbers, strings
- layout tokens: `INDENT`, `DEDENT`, `NEWLINE`

Current lexer behavior:

- normalizes CRLF to LF
- emits explicit newline tokens
- tracks indentation width across lines
- treats tabs as width 4 for indentation counting
- parses decimal numbers into token text later converted by the parser

This makes the language block-structured and indentation-sensitive.

### 4. Parsing

[`compiler/parser/parser.go`](../compiler/parser/parser.go) is a recursive-descent parser over the lexer token stream.

The parser currently builds:

- `ast.Program`
- `ast.Import`
- `ast.Function`
- `ast.LetStmt`
- `ast.ReturnStmt`
- `ast.ExprStmt`
- `ast.Ident`
- `ast.NumberLiteral`
- `ast.StringLiteral`
- `ast.BinaryExpr`
- `ast.CallExpr`

Grammar features supported in practice:

- top-level `import` and `fn`
- optional parameter type annotations like `x: Number`
- optional function return type syntax `-> Type`
- colon-plus-indented blocks
- arithmetic precedence for `+`, `-`, `*`, `/`
- parenthesized expressions
- function and builtin calls

Error handling is still panic-based. Unsupported top-level forms, invalid call targets, and unexpected tokens currently fail fast instead of returning structured diagnostics.

## Compiler Layer

### `compiler/ast`

[`compiler/ast/nodes.go`](../compiler/ast/nodes.go) defines the current tree shape.

Key design points:

- `Node` exposes `Pos() int`
- statements and expressions are modeled through `Stmt` and `Expr` marker interfaces
- `Program` separates `Imports` from `Functions`
- function bodies are stored as `[]Stmt`

Compared with earlier revisions, the AST is now substantially more complete:

- position methods exist on concrete nodes
- identifiers, literals, binary expressions, and calls are represented explicitly
- import declarations are first-class nodes

### `compiler/bytecode`

[`compiler/bytecode/gen.go`](../compiler/bytecode/gen.go) is the main lowering stage used by the project today.

Core types:

- `CompiledProgram`
- `CompiledFunction`
- `Generator`

Compilation behavior:

- indexes functions by name before code generation
- resets local state per function
- assigns parameter slots first
- allocates local slots on first `let`
- deduplicates constants within a function
- emits bytecode for expressions and statements

Statement lowering:

- `let` evaluates the expression and stores it in a local slot
- `return` evaluates and emits `RETURN`
- bare expression statements evaluate and emit `POP`

Expression lowering:

- literals load from the function constant pool
- identifiers load from frame slots
- arithmetic maps to VM opcodes
- calls dispatch either to user-defined functions or builtin indexes

This package is now the real compiler backend for the repository. By contrast, [`compiler/ir/ir.go`](../compiler/ir/ir.go) remains a disconnected minimal type and is not part of the active pipeline.

### Placeholder Compiler Packages

The following directories still exist without active implementation files:

- `compiler/emitter`
- `compiler/sema`
- `compiler/toon`

These appear reserved for future compilation stages or experiments, but they do not participate in the current build.

## VM Layer

### `vm/opcode`

[`vm/opcode/opcode.go`](../vm/opcode/opcode.go) defines the interpreter instruction set:

- `LOAD_CONST`
- `LOAD_VAR`
- `STORE_VAR`
- `ADD`
- `SUB`
- `MUL`
- `DIV`
- `CALL`
- `CALL_BUILTIN`
- `POP`
- `RETURN`

It also maps builtin names to builtin indexes:

- `pt`
- `ptln`

### `vm/runtime`

[`vm/runtime/stack.go`](../vm/runtime/stack.go) and [`vm/runtime/frame.go`](../vm/runtime/frame.go) provide the mutable runtime containers.

`Stack`:

- backed by `[]interface{}`
- supports `Push` and `Pop`
- returns `nil` on underflow

`Frame`:

- backed by `[]interface{}`
- stores local values by slot index
- is sized from the compiled function's `FrameSize`

The model is intentionally simple and works well for the current single-frame-per-call execution design.

### `vm/exec`

[`vm/exec/interpreter.go`](../vm/exec/interpreter.go) executes compiled bytecode function-by-function.

Interpreter behavior:

- resolves the entry function by name from `CompiledProgram.FunctionIndex`
- creates a new frame for each call
- loads incoming arguments into frame slots
- uses a fresh operand stack per function invocation
- dispatches opcodes in a single loop

Implemented execution features:

- arithmetic on `float64`
- local loads and stores
- user-defined function calls
- builtin calls through `vm/builtin`
- returning the top value from the stack

There is no structured runtime error system yet. Invalid indexes, wrong value shapes, or unsupported states tend to surface as panics or raw Go runtime failures.

### `vm/builtin`

[`vm/builtin/print.go`](../vm/builtin/print.go) currently exposes two builtins:

- `pt`: print arguments without a newline
- `ptln`: print arguments and then add a newline

Both builtins accept variadic `interface{}` values and write directly to standard output.

### `vm/agents`

[`vm/agents/agents.go`](../vm/agents/agents.go) remains separate from the source language and VM execution flow.

It currently provides:

- an `Agent` struct with `Name` and `Run`
- a `Start()` method that runs `Run` repeatedly in a goroutine every 10 milliseconds

This is still a runtime concept stub rather than a language feature wired into parsing, code generation, or execution.

## Current Example Path

[`examples/hello.xavi`](../examples/hello.xavi) is the best reference for the implemented system.

It demonstrates:

- named imports from sibling files
- nested user-defined function calls
- builtin printing
- arithmetic evaluation

Supporting files:

- [`examples/sum.xavi`](../examples/sum.xavi)
- [`examples/operation/mult.xavi`](../examples/operation/mult.xavi)

Running:

```bash
go run . examples/hello.xavi
```

Currently produces:

```text
Hello world20 hello new line
Result: 30
```

## Build Status

As of August 17, 2026:

- `go test ./...` succeeds
- there are no Go test files yet, so the build status currently reflects compile-time validation rather than behavioral test coverage

This is an important change from earlier repository states where parser incompleteness prevented a clean build.

## Architectural Gaps

Even with the active end-to-end path, several architectural gaps remain:

- no semantic analysis stage
- no type checker despite typed syntax hooks in the parser
- no namespace or module object model for imports
- no explicit diagnostic/reporting abstraction
- no control flow beyond straight-line statements and function calls
- no records, agents, events, or `on` handlers in the executable pipeline yet
- no dedicated CLI under `cmd/xavic`
- no tests around lexer, parser, loader, bytecode generation, or VM execution

## Likely Next Steps

1. Add tests that lock down the current example-driven workflow.
2. Replace parser and generator panics with structured errors.
3. Introduce semantic analysis for symbol checks and duplicate definition handling.
4. Decide whether `compiler/ir` should become a real stage or be removed.
5. Expand the bytecode and VM model before adding richer language syntax.
