package builtin

import "fmt"

func Call(index uint8, args []interface{}) interface{} {
	switch index {
	case 0:
		printArgs(args...)
		return nil
	case 1:
		printArgs(args...)
		fmt.Println()
		return nil
	default:
		panic(fmt.Sprintf("unknown builtin function: %d", index))
	}
}

func printArgs(args ...interface{}) {
	for _, arg := range args {
		fmt.Print(arg)
	}
}
