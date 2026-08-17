# Architecture Notes

## Overview

Xavi is structured like a language implementation with three broad layers:

1. source language front-end under `compiler/`
2. execution runtime under `vm/`
3. application or tooling entrypoints under `main.go` and `cmd/`

Today, only part of that end-to-end story is implemented. The architecture is still useful because it shows the intended separation of concerns.

## Compiler Layer

### `compiler/lexer`

[`compiler/lexer/lexer.go`](../compiler/lexer/lexer.go) defines token categories but does not yet include a lexer implementation.

Interesting tokens:

- keywords: `FN`, `RECORD`, `RETURN`, `LET`, `AGENT`, `ON`, `EVENT`
- punctuation: `COLON`, `ARROW`, `LPAREN`, `RPAREN`
- layout-aware tokens: `INDENT`, `DEDENT`, `NEWLINE`

This suggests the language is intended to support indentation-sensitive parsing, similar to Python-style blocks.

### `compiler/ast`

[`compiler/ast/nodes.go`](../compiler/ast/nodes.go) currently defines:

- `Node`
- `Function`
- `Param`
- `Return`
- `BinaryOp`

Observations:

- `Node` requires a `Pos() int` method, but the concrete structs shown in this file do not implement it yet.
- expression coverage is still incomplete: there are no dedicated number, string, identifier, or call nodes.
- `Function.Body` is represented as `[]Node`, which is a reasonable base for a statement list.

### `compiler/parser`

[`compiler/parser/parser.go`](../compiler/parser/parser.go) is a partial recursive-descent parser.

Implemented flow:

1. expect `fn`
2. read function name
3. parse parameter list
4. read return arrow
5. read return type
6. parse block body
7. return an AST function node

Missing pieces:

- token lookahead utilities
- `expect(...)`
- `parseParams()`
- `parseBlock()`
- error handling and recovery

Because of those missing methods, the parser package does not currently compile.

### `compiler/ir`

[`compiler/ir/ir.go`](../compiler/ir/ir.go) defines a minimal instruction shape:

```go
type Instr struct {
    Op   string
    Arg1 string
    Arg2 string
}
```

This can serve as a transitional layer between AST and VM bytecode, but it is not connected to the rest of the pipeline yet.

### Placeholder Compiler Packages

These directories currently exist without implementation files:

- `compiler/bytecode`
- `compiler/emitter`
- `compiler/sema`
- `compiler/toon`

Probable roles:

- `bytecode`: VM-oriented compiled output structures
- `emitter`: lowering AST or IR into bytecode
- `sema`: type checking and symbol validation
- `toon`: likely experimental or domain-specific compilation work

## VM Layer

### `vm/runtime`

This package provides the mutable runtime containers used by the interpreter.

`Stack`:

- backed by `[]interface{}`
- supports `Push` and `Pop`
- returns `nil` on underflow rather than panicking

`Frame`:

- backed by `[]interface{}`
- represents local variable slots
- supports indexed `Store` and `Load`

The runtime keeps the model intentionally simple, which is a good fit for a prototype interpreter.

### `vm/exec`

This is the core execution engine.

Opcode set:

- `OP_LOAD_CONST`
- `OP_LOAD_VAR`
- `OP_STORE_VAR`
- `OP_ADD`
- `OP_SUB`
- `OP_MUL`
- `OP_DIV`
- `OP_RETURN`

Interpreter responsibilities:

- hold bytecode and instruction pointer
- own the operand stack
- own the constant pool
- own the current frame
- dispatch instructions in a loop

Data assumptions:

- arithmetic opcodes cast stack operands to `float64`
- local variables are stored as `interface{}`
- bytecode operands are single-byte indexes

Those choices keep the implementation small, but they also create natural next steps:

- bounds checks for locals and constant indexes
- structured runtime errors
- richer value tagging or typed value containers
- support for calls and multiple frames

### `vm/agents`

The agent package is separate from the interpreter and models repeated background work.

Current behavior:

- `Start()` launches a goroutine
- the goroutine runs forever
- `Run()` executes every 10 milliseconds

What is missing:

- cancellation / stop control
- error handling
- scheduler integration
- event source wiring

## Entrypoints

### `main.go`

`main.go` currently bypasses the compiler and feeds handcrafted bytecode directly into the VM.

This is useful because it validates the interpreter model independently from the unfinished compiler. It also serves as a concrete example of the expected low-level instruction format.

### `cmd/xavic`

The `cmd/xavic` directory exists but is empty. This is the natural location for a future CLI such as:

- `xavic run examples/hello.xavi`
- `xavic build file.xavi`
- `xavic ast file.xavi`
- `xavic disasm file.xavi`

## Gaps In The Current End-To-End Flow

The intended language pipeline appears to be:

1. source text
2. lexer tokens
3. AST
4. IR or direct lowering
5. bytecode
6. VM execution

Right now, the implemented path is effectively:

1. handcrafted bytecode
2. VM execution

That means the documentation and future planning should treat the compiler and runtime as partially decoupled tracks rather than a complete single feature.

## Recommended Next Engineering Steps

1. Make the AST compile cleanly by implementing `Pos()` on node types or relaxing the `Node` interface temporarily.
2. Complete parser support for function signatures and blocks.
3. Add a real lexer that emits the token types already defined.
4. Decide whether IR is required for v1 or whether AST can lower directly to bytecode.
5. Add tests around `vm/exec`, especially arithmetic and local variable behavior.
6. Create one real `.xavi` example and wire it into an executable CLI or test.
