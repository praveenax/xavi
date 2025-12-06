package exec

import "xavi/vm/runtime"

type Interpreter struct {
	BC     []uint8
	IP     int
	Stack  *runtime.Stack
	Consts []float64
	Frame  *runtime.Frame
}

func NewInterpreter(bc []uint8, consts []float64, frameSize int) *Interpreter {
	return &Interpreter{
		BC:     bc,
		IP:     0,
		Stack:  runtime.NewStack(),
		Consts: consts,
		Frame:  runtime.NewFrame(frameSize),
	}
}

func (vm *Interpreter) Run() interface{} {
	for vm.IP < len(vm.BC) {
		op := vm.BC[vm.IP]
		vm.IP++

		switch op {

		case OP_LOAD_CONST:
			idx := vm.BC[vm.IP]
			vm.IP++
			vm.Stack.Push(vm.Consts[idx])

		case OP_ADD:
			b := vm.Stack.Pop().(float64)
			a := vm.Stack.Pop().(float64)
			vm.Stack.Push(a + b)

		case OP_MUL:
			b := vm.Stack.Pop().(float64)
			a := vm.Stack.Pop().(float64)
			vm.Stack.Push(a * b)

		case OP_SUB:
			b := vm.Stack.Pop().(float64)
			a := vm.Stack.Pop().(float64)
			vm.Stack.Push(a - b)

		case OP_DIV:
			b := vm.Stack.Pop().(float64)
			a := vm.Stack.Pop().(float64)
			vm.Stack.Push(a / b)

		case OP_LOAD_VAR:
			idx := vm.BC[vm.IP]
			vm.IP++
			vm.Stack.Push(vm.Frame.Load(int(idx)))

		case OP_STORE_VAR:
			idx := vm.BC[vm.IP]
			vm.IP++
			val := vm.Stack.Pop()
			vm.Frame.Store(int(idx), val)

		case OP_RETURN:
			return vm.Stack.Pop()
		}
	}

	return nil
}
