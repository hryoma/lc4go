package emulator

import (
	"fmt"
	"github.com/hryoma/lc4go/machine"
	"github.com/hryoma/lc4go/tokenizer"
)

func InitLc4(lc4 *machine.Machine) {
	lc4.Code = [machine.CODE_SIZE]machine.Insn{}
	lc4.Mem = [machine.MEM_SIZE]uint16{}
	lc4.Reg = [machine.NUM_REGS]uint16{}
	lc4.Nzp = 0
	lc4.Pc = 0
	lc4.Psr = 1
	lc4.Labels = map[string]uint16{}
}

func Breakpoint(lc4 *machine.Machine) {
	fmt.Println("breakpoint / b")
}

func Continue(lc4 *machine.Machine) {
	fmt.Println("continue / c")
	for {
		if lc4.Code[lc4.Pc].Breakpoint {
			// TODO print breakpoint message
			return
		}

		if ok := Step(lc4); !ok {
			return
		}
	}
}

func Load(lc4 *machine.Machine) {
	fmt.Println("load / l")

	// make sure everything is reset by calling init

	// read file, make sure it exists

	// load input file
	tokenizer.Tokenize()

	// set machine state

}

func Next(lc4 *machine.Machine) {
	fmt.Println("next / n")
	nextPc := lc4.Pc + 1
	// like continue, but loop over step until pc = pc_curr + 1
	for {
		if lc4.Pc == nextPc {
			return
		}
		if lc4.Code[lc4.Pc].Breakpoint {
			// TODO print breakpoint message
			return
		}

		if ok := Step(lc4); !ok {
			return
		}
	}
}

func Print(lc4 *machine.Machine) {
	fmt.Println("print / p")
	// print current line of code, nzp, pc, reg
}

func PrintCode(lc4 *machine.Machine) {
	fmt.Println("print / p -c")
	// print code at address
}

func PrintMem(lc4 *machine.Machine) {
	fmt.Println("print / p -m")
	// print memory at address
}

func PrintNzp(lc4 *machine.Machine) {
	fmt.Println("print / p -n")
	// print nzp bits
}

func PrintReg(lc4 *machine.Machine) {
	fmt.Println("print / p -r")
	// check if -r has a value
	// if it doesn't, print all regs
	// if it does, check that it falls between 0-7
	// otherwise, error, invalid reg num
}

func Run(lc4 *machine.Machine) {
	fmt.Println("run / r")
	// TODO reset machine
	Continue(lc4)
}

func Reset(lc4 *machine.Machine) {
	fmt.Println("reset")
	// call init again, and then load data
}

func Step(lc4 *machine.Machine) (ok bool) {
	fmt.Println("step / s")
	// TODO add check to see if pc in valid address
	// TODO add check to see if pc is at end
	if lc4.Pc == 0x80FF {
		return false
	}
	// then execute once
	lc4.Pc++
	return true
}

