package exec

import (
	"fmt"

	"xavi/compiler/bytecode"
	"xavi/vm/opcode"
	"xavi/vm/runtime"
)

type Interpreter struct {
	Program *bytecode.CompiledProgram
}

func NewInterpreter(program *bytecode.CompiledProgram) *Interpreter {
	return &Interpreter{Program: program}
}

func (vm *Interpreter) Run(functionName string) interface{} {
	functionIndex, exists := vm.Program.FunctionIndex[functionName]
	if !exists {
		panic(fmt.Sprintf("entry function not found: %s", functionName))
	}

	return vm.runFunction(vm.Program.Functions[functionIndex], nil)
}

func (vm *Interpreter) runFunction(fn *bytecode.CompiledFunction, args []interface{}) interface{} {
	frame := runtime.NewFrame(fn.FrameSize)
	for index, arg := range args {
		frame.Store(index, arg)
	}

	stack := runtime.NewStack()
	ip := 0

	for ip < len(fn.Bytecode) {
		op := fn.Bytecode[ip]
		ip++

		switch op {
		case opcode.LOAD_CONST:
			idx := fn.Bytecode[ip]
			ip++
			stack.Push(fn.Consts[idx])

		case opcode.ADD:
			b := stack.Pop().(float64)
			a := stack.Pop().(float64)
			stack.Push(a + b)

		case opcode.MUL:
			b := stack.Pop().(float64)
			a := stack.Pop().(float64)
			stack.Push(a * b)

		case opcode.SUB:
			b := stack.Pop().(float64)
			a := stack.Pop().(float64)
			stack.Push(a - b)

		case opcode.DIV:
			b := stack.Pop().(float64)
			a := stack.Pop().(float64)
			stack.Push(a / b)

		case opcode.LOAD_VAR:
			idx := fn.Bytecode[ip]
			ip++
			stack.Push(frame.Load(int(idx)))

		case opcode.STORE_VAR:
			idx := fn.Bytecode[ip]
			ip++
			val := stack.Pop()
			frame.Store(int(idx), val)

		case opcode.CALL:
			functionIndex := fn.Bytecode[ip]
			ip++
			argCount := int(fn.Bytecode[ip])
			ip++

			callArgs := make([]interface{}, argCount)
			for i := argCount - 1; i >= 0; i-- {
				callArgs[i] = stack.Pop()
			}

			result := vm.runFunction(vm.Program.Functions[functionIndex], callArgs)
			stack.Push(result)

		case opcode.RETURN:
			return stack.Pop()
		}
	}

	return nil
}
