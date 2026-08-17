# Xavi

Xavi is an early-stage Go project for a small language and virtual machine. The repository already contains the main building blocks of a language toolchain:

- a compiler-facing package layout with lexer, parser, AST, and IR packages
- a bytecode interpreter with a stack and local frame model
- a lightweight agent runtime concept for repeated background execution

At the moment, the VM layer is the most concrete part of the codebase. The compiler side is still partial, with some packages acting as placeholders for future work.

## Repository Layout

```text
xavi/
├── cmd/
│   └── xavic/                # Reserved for a future CLI entrypoint
├── compiler/
│   ├── ast/                  # AST node types
│   ├── bytecode/             # Placeholder
│   ├── emitter/              # Placeholder
│   ├── ir/                   # Intermediate representation types
│   ├── lexer/                # Token and token-type definitions
│   ├── parser/               # Partial parser implementation
│   ├── sema/                 # Placeholder
│   └── toon/                 # Placeholder
├── examples/
│   └── hello.xavi            # Currently empty
├── vm/
│   ├── agents/               # Minimal agent loop abstraction
│   ├── builtin/              # Placeholder
│   ├── exec/                 # Bytecode interpreter and opcodes
│   ├── loader/               # Placeholder
│   └── runtime/              # Stack and frame storage
├── go.mod
├── LICENSE
└── main.go                   # Demo program using handwritten bytecode
```

## What Exists Today

### 1. Demo Entrypoint

[`main.go`](./main.go) manually constructs bytecode for a tiny arithmetic program equivalent to:

```text
return 10 + 20
```

The program:

- loads two numeric constants
- stores them in local frame slots
- reloads them
- adds them
- returns the result from the stack

This file is best read as a VM demonstration rather than a full compiler pipeline.

### 2. Virtual Machine

The interpreter in [`vm/exec/interpreter.go`](./vm/exec/interpreter.go) executes a byte slice instruction-by-instruction.

Supported opcodes in [`vm/exec/opcodes.go`](./vm/exec/opcodes.go):

- `OP_LOAD_CONST`
- `OP_LOAD_VAR`
- `OP_STORE_VAR`
- `OP_ADD`
- `OP_SUB`
- `OP_MUL`
- `OP_DIV`
- `OP_RETURN`

Runtime support:

- [`vm/runtime/stack.go`](./vm/runtime/stack.go) implements a simple LIFO stack over `[]interface{}`
- [`vm/runtime/frame.go`](./vm/runtime/frame.go) stores locals in indexed slots

Execution model:

1. Read opcode from `BC`
2. Advance instruction pointer
3. Perform stack or frame operation
4. Continue until `OP_RETURN` or end of bytecode

The interpreter currently assumes numeric arithmetic with `float64` values for arithmetic operations.

### 3. Compiler Skeleton

The `compiler` tree outlines the intended front-end pipeline:

- [`compiler/lexer/lexer.go`](./compiler/lexer/lexer.go) defines token types and token metadata
- [`compiler/ast/nodes.go`](./compiler/ast/nodes.go) defines core AST structures like `Function`, `Param`, `Return`, and `BinaryOp`
- [`compiler/ir/ir.go`](./compiler/ir/ir.go) defines a minimal IR instruction type
- [`compiler/parser/parser.go`](./compiler/parser/parser.go) begins a function parser

The parser currently expects a function syntax shaped like:

```text
fn <name>(<params>) -> <returnType>
```

It then attempts to parse a block body, but helper methods such as `expect`, `parseParams`, and `parseBlock` are not implemented yet.

### 4. Agent Runtime Concept

[`vm/agents/agents.go`](./vm/agents/agents.go) defines:

- an `Agent` with a `Name`
- a `Run` callback
- a `Start()` method that repeatedly executes `Run` in a goroutine every 10ms

This is a minimal building block for evented or autonomous runtime behavior. It is not yet integrated with the compiler or interpreter.

## Current Project Status

This repository is best understood as a prototype / scaffold.

Implemented:

- bytecode execution for simple arithmetic
- local variable storage via frame slots
- basic AST and IR type definitions
- token type definitions for a higher-level language design

Incomplete or placeholder areas:

- actual lexing logic
- parser helper methods and block parsing
- semantic analysis
- bytecode generation from AST/IR
- CLI compiler command under `cmd/xavic`
- builtins, loader, and example programs

## Build And Run

From the project root:

```bash
go run .
```

Expected behavior from the current demo:

```text
Result: 30
```

To run package checks:

```bash
go test ./...
```

## Known Limitations

As of August 17, 2026, the repository does not build cleanly as a full project because `compiler/parser/parser.go` references methods that do not exist yet:

- `expect`
- `parseParams`
- `parseBlock`

That means:

- `go test ./...` fails at compile time in the parser package
- the current entrypoint is a handwritten VM example rather than the output of a compiler

## Suggested Roadmap

### Near-Term

1. Implement parser helpers (`expect`, token navigation, parameter parsing, block parsing).
2. Add real lexer logic to turn source text into tokens.
3. Define literal and identifier AST nodes so expressions can be represented cleanly.
4. Add a bytecode emitter from AST or IR to VM opcodes.
5. Populate `examples/hello.xavi` with a minimal source program.

### Mid-Term

1. Add a CLI under `cmd/xavic` for compile/run workflows.
2. Introduce semantic checks for symbols, return types, and arity.
3. Add built-in functions and loader support.
4. Define how `agent`, `on`, and `event` syntax should lower into runtime constructs.

## Design Notes

The token set suggests the language may be aiming for:

- function declarations with explicit return types
- records / structured data
- variable bindings
- agent/event-oriented constructs
- indentation-aware parsing (`INDENT` / `DEDENT`)

That design direction is not fully implemented yet, but the package layout already reflects it.

## Additional Documentation

For a package-by-package walkthrough, see [docs/architecture.md](./docs/architecture.md).
