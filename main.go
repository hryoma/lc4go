package main

import (
	"fmt"
	"github.com/hryoma/lc4go/emulator"
	"github.com/hryoma/lc4go/tokenizer"
)

func main() {
	fmt.Println("LC4 ISA Emulator using Go")
	tokenizer.Tokenize()
	emulator.Emulate()
}
