package runtime

type Frame struct {
	Locals []interface{}
}

func NewFrame(size int) *Frame {
	return &Frame{
		Locals: make([]interface{}, size),
	}
}

func (f *Frame) Store(index int, val interface{}) {
	f.Locals[index] = val
}

func (f *Frame) Load(index int) interface{} {
	return f.Locals[index]
}
