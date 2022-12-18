package emulator

import (
	"fmt"
	"github.com/hryoma/lc4go/machine"
	"github.com/hryoma/lc4go/tokenizer"
)

func InitLc4() {
	machine.Lc4.Code = [machine.CODE_SIZE]machine.Insn{}
	machine.Lc4.Mem = [machine.MEM_SIZE]uint16{}
	machine.Lc4.Reg = [machine.NUM_REGS]uint16{}
	machine.Lc4.Nzp = 0
	machine.Lc4.Pc = 0
	machine.Lc4.Psr = 1
	machine.Lc4.Labels = map[string]uint16{}
}

func Breakpoint() {
	fmt.Println("breakpoint / b")
}

func Continue() {
	fmt.Println("continue / c")
	for {
		if machine.Lc4.Code[machine.Lc4.Pc].Breakpoint {
			// TODO print breakpoint message
			return
		}

		if ok := Step(); !ok {
			return
		}
	}
}

func Load(fileName string) {
	fmt.Println("load / l")
	fmt.Printf("obj file: %s\n", fileName)
	
	InitLc4()
	tokenizer.TokenizeObj(fileName)
}

func Next() {
	fmt.Println("next / n")
	nextPc := machine.Lc4.Pc + 1
	// like continue, but loop over step until pc = pc_curr + 1
	for {
		if machine.Lc4.Pc == nextPc {
			return
		}
		if machine.Lc4.Code[machine.Lc4.Pc].Breakpoint {
			// TODO print breakpoint message
			return
		}

		if ok := Step(); !ok {
			return
		}
	}
}

func Print() {
	fmt.Println("print / p")
	// print current line of code, nzp, pc, reg
}

func PrintCode() {
	fmt.Println("print / p -c")
	// print code at address
}

func PrintMem() {
	fmt.Println("print / p -m")
	// print memory at address
}

func PrintNzp() {
	fmt.Println("print / p -n")
	// print nzp bits
}

func PrintReg() {
	fmt.Println("print / p -r")
	// check if -r has a value
	// if it doesn't, print all regs
	// if it does, check that it falls between 0-7
	// otherwise, error, invalid reg num
}

func Run() {
	fmt.Println("run / r")
	// TODO reset machine
	Continue()
}

func Reset() {
	fmt.Println("reset")
	// call init again, and then load data
}

func Step() (ok bool) {
	fmt.Println("step / s")
	// TODO add check to see if pc in valid address
	// TODO add check to see if pc is at end
	if machine.Lc4.Pc == 0x80FF {
		return false
	}
	// then execute once

	return true
}

