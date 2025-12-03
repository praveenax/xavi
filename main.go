package main

import (
    "fmt"
    "xavi/vm/exec"
)

func main() {
    // Bytecode for:
    // return 10 + 20
    bytecode := []uint8{
        exec.OP_LOAD_CONST, 0,   // push const[0] = 10
        exec.OP_LOAD_CONST, 1,   // push const[1] = 20
        exec.OP_ADD,             // add
        exec.OP_RETURN,          // return result
    }

    constPool := []float64{10, 20}

    vm := exec.NewInterpreter(bytecode, constPool)
    result := vm.Run()

    fmt.Println("Result:", result)
}
