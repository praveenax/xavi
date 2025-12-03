package runtime

type Stack struct {
    data []interface{}
}

func NewStack() *Stack {
    return &Stack{data: []interface{}{}}
}

func (s *Stack) Push(v interface{}) {
    s.data = append(s.data, v)
}

func (s *Stack) Pop() interface{} {
    if len(s.data) == 0 {
        return nil
    }
    v := s.data[len(s.data)-1]
    s.data = s.data[:len(s.data)-1]
    return v
}
