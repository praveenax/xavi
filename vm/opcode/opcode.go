package opcode

const (
	LOAD_CONST = iota
	LOAD_VAR
	STORE_VAR
	ADD
	SUB
	MUL
	DIV
	CALL
	CALL_BUILTIN
	POP
	RETURN
)

const (
	BuiltinPt uint8 = iota
	BuiltinPtln
)

func BuiltinIndex(name string) (uint8, bool) {
	switch name {
	case "pt":
		return BuiltinPt, true
	case "ptln":
		return BuiltinPtln, true
	default:
		return 0, false
	}
}
