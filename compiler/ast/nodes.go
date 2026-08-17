package ast

type Node interface {
	Pos() int
}

type Stmt interface {
	Node
	stmtNode()
}

type Expr interface {
	Node
	exprNode()
}

type Program struct {
	Imports   []*Import
	Functions []*Function
}

type Import struct {
	Name     string
	Path     string
	position int
}

func (i *Import) Pos() int { return i.position }

type Function struct {
	Name       string
	Params     []Param
	ReturnType string
	Body       []Stmt
	position   int
}

func (f *Function) Pos() int { return f.position }

type Param struct {
	Name string
	Type string
}

type LetStmt struct {
	Name     string
	Value    Expr
	position int
}

func (s *LetStmt) Pos() int  { return s.position }
func (s *LetStmt) stmtNode() {}

type ReturnStmt struct {
	Value    Expr
	position int
}

func (s *ReturnStmt) Pos() int  { return s.position }
func (s *ReturnStmt) stmtNode() {}

type ExprStmt struct {
	Value    Expr
	position int
}

func (s *ExprStmt) Pos() int  { return s.position }
func (s *ExprStmt) stmtNode() {}

type Ident struct {
	Name     string
	position int
}

func (e *Ident) Pos() int  { return e.position }
func (e *Ident) exprNode() {}

type NumberLiteral struct {
	Value    float64
	position int
}

func (e *NumberLiteral) Pos() int  { return e.position }
func (e *NumberLiteral) exprNode() {}

type StringLiteral struct {
	Value    string
	position int
}

func (e *StringLiteral) Pos() int  { return e.position }
func (e *StringLiteral) exprNode() {}

type BinaryExpr struct {
	Op       string
	Left     Expr
	Right    Expr
	position int
}

func (e *BinaryExpr) Pos() int  { return e.position }
func (e *BinaryExpr) exprNode() {}

type CallExpr struct {
	Callee   string
	Args     []Expr
	position int
}

func (e *CallExpr) Pos() int  { return e.position }
func (e *CallExpr) exprNode() {}
