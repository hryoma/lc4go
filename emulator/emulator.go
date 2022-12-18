package emulator

import (
	"fmt"
	"github.com/hryoma/lc4go/machine"
	"github.com/hryoma/lc4go/tokenizer"
)

func InitLc4() {
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
	for {
		// TODO store bp's in its own bit mask
		// if machine.Lc4.Code[machine.Lc4.Pc].Breakpoint {
			// TODO print breakpoint message
		//	return
		//}

		if ok := Step(); !ok {
			return
		}
	}
}

func Load(fileName string) {
	InitLc4()
	tokenizer.TokenizeObj(fileName)
}

func Next() {
	nextPc := machine.Lc4.Pc + 1
	// like continue, but loop over step until pc = pc_curr + 1
	for {
		if machine.Lc4.Pc == nextPc {
			return
		}
		// TODO store bp's in it's own bitmask
		//if machine.Lc4.Code[machine.Lc4.Pc].Breakpoint {
			//// TODO print breakpoint message
			//return
		//}

		if ok := Step(); !ok {
			return
		}
	}
}

func Print() {
	PrintCode()
	PrintNzp()
	PrintReg()
}

func PrintCode() {
	pc := machine.Lc4.Pc
	data := machine.Lc4.Mem[pc]
	fmt.Printf("%d:\t0b%016b / 0x%04X\n", pc, data, data)
}

func PrintMem() {
	fmt.Println("print / p -m")
	// print memory at address
}

func PrintNzp() {
	var n, z, p uint8

	if machine.Lc4.Psr & 0b100 != 0 {
		n = 1
	}
	if machine.Lc4.Psr & 0b010 != 0 {
		z = 1
	}
	if machine.Lc4.Psr & 0b001 != 0 {
		p = 1
	}

	fmt.Printf("nzp:\t%01b/%01b/%01b\n", n, z, p)
}

func PrintReg() {
	var regMask uint16 = 0x00FF

	for i := uint16(0); i < 8; i++ {
		if (regMask >> i) & 1 == 1 {
			regVal := machine.Lc4.Reg[i]
			fmt.Printf("\tR%d: %016b / 0x%04X\n", i, regVal, regVal)
		}
	}
}

func Run() {
	// TODO reset machine
	Continue()
}

func Reset() {
	// call init again, and then load data
}

func Step() (ok bool) {
	// TODO add check to see if pc in valid address
	// TODO add check to see if pc is at end
	if machine.Lc4.Pc == 0x80FF {
		return false
	}
	// then execute once
	err := machine.Execute()
	if err != 0 {
		fmt.Println("execute error")
	}

	return true
}

