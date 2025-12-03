package ast

type Node interface {
    Pos() int
}

type Function struct {
    Name       string
    Params     []Param
    ReturnType string
    Body       []Node
}

type Param struct {
    Name string
    Type string
}

type Return struct {
    Value Node
}

type BinaryOp struct {
    Op    string
    Left  Node
    Right Node
}
