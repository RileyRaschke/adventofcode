package main

import (
    "os"
    "bufio"
    "io"
    "fmt"
    "strings"
    "strconv"
)

type InstructionOp struct {
    operation string
    arg int
}

type Instruction struct {
    op InstructionOp
    evalCount int
}

type Parser struct {
    code []Instruction
}
func NewParser() *Parser { return &Parser{ []Instruction{} } }

func (self *Parser) AddStringOp(operation string) {
    opts := strings.Split(operation, " ")
    instr := opts[0]
    arg, err := strconv.Atoi(opts[1])
    if err != nil { panic("Invalid input") }
    self.code = append(self.code, Instruction{ InstructionOp{ instr, arg }, 0 } )
}

func (self *Parser) Parse() (valid bool, accumulation int) {
    valid = true
    var i = 0
    for {
        self.code[i].evalCount++
        if self.code[i].evalCount > 1 {
            valid = false
            break
        }
        switch op := self.code[i].op.operation; op {
            case "nop":
                i++
                break
            case "jmp":
                i = i + self.code[i].op.arg
                break
            case "acc":
                accumulation = accumulation + self.code[i].op.arg
                i++
                break
        }
    }
    return
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    parser := NewParser()
    for {
        str, _, err := reader.ReadLine()
        if err == io.EOF {
            break;
        }
        line := strings.TrimSpace(string(str))
        parser.AddStringOp(line)
    }

    if ok, accum := parser.Parse(); !ok {
        fmt.Printf("Found loop with accumulation at: %v\n", accum )
    }
}

