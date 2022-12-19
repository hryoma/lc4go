package emulator

import (
	"fmt"
	"github.com/hryoma/lc4go/machine"
	"github.com/hryoma/lc4go/tokenizer"
	"strconv"
)

const PC_INIT_VAL = 0x8200
const PSR_INIT_VAL = 0x8002
const PC_TERM = 0x80FF

func Breakpoint(strAddr string) {
	if addr, err := strconv.ParseUint(strAddr, 0, 16); err == nil {

		if meta, exists := machine.Lc4.Meta[uint16(addr)]; exists {
			meta.Breakpoint = true
		} else {
			machine.Lc4.Meta[uint16(addr)] = machine.MemMetadata{
				Breakpoint: true,
			}
		}
		fmt.Printf("Breakpoint set at 0x%04X\n", addr)
	} else {
		fmt.Println("Invalid address:", strAddr)
	}
}

func Clear() {
	machine.Lc4.Mem = [machine.MEM_SIZE]uint16{}
	machine.Lc4.Meta = map[uint16]machine.MemMetadata{}
	Reset()
}

func Continue() {
	for {
		if ok := Step(); !ok {
			return
		}

		// stop if breakpoint is hit
		addr := machine.Lc4.Pc
		if meta, exists := machine.Lc4.Meta[addr]; exists && meta.Breakpoint {
			fmt.Printf("Hit breakpoint at 0x%04X\n", addr)
			return
		}
	}
}

func Load(fileName string) {
	tokenizer.TokenizeObj(fileName)
}

func Next() {
	nextPc := machine.Lc4.Pc + 1
	for {
		if ok := Step(); !ok {
			return
		}

		if machine.Lc4.Pc == nextPc {
			return
		}

		// stop if breakpoint is hit
		addr := machine.Lc4.Pc
		if meta, exists := machine.Lc4.Meta[addr]; exists && meta.Breakpoint {
			fmt.Printf("Hit breakpoint at 0x%04X\n", addr)
			return
		}
	}
}

func Print() {
	PrintCode()
	PrintPsr()
	PrintReg()
}

func PrintCode() {
	pc := machine.Lc4.Pc
	data := machine.Lc4.Mem[pc]
	fmt.Printf("0x%04X:\t0b%016b / 0x%04X\n", pc, data, data)
}

func PrintMem(strAddr string) {
	if addr, err := strconv.ParseUint(strAddr, 0, 16); err == nil {
		data := machine.Lc4.Mem[addr]
		fmt.Printf("0x%04X:\t0b%016b / 0x%04X\n", addr, data, data)
	} else {
		fmt.Println("Invalid address:", strAddr)
	}
}

func PrintPsr() {
	var n, z, p uint8

	if machine.Lc4.Psr&0b100 != 0 {
		n = 1
	}
	if machine.Lc4.Psr&0b010 != 0 {
		z = 1
	}
	if machine.Lc4.Psr&0b001 != 0 {
		p = 1
	}
	priv := machine.Lc4.Psr&0x8000 != 0

	fmt.Printf("psr:\t%01b/%01b/%01b (privilege %t)\n", n, z, p, priv)
}

func PrintReg() {
	var regMask uint16 = 0x00FF

	for i := uint16(0); i < 8; i++ {
		if (regMask>>i)&1 == 1 {
			regVal := machine.Lc4.Reg[i]
			fmt.Printf("\tR%d: %016b / 0x%04X\n", i, regVal, regVal)
		}
	}
}

func Run() {
	Reset()
	Continue()
}

func Reset() {
	machine.Lc4.Reg = [machine.NUM_REGS]uint16{}
	machine.Lc4.Pc = PC_INIT_VAL
	machine.Lc4.Psr = PSR_INIT_VAL
}

func Step() (ok bool) {
	if machine.Lc4.Pc == PC_TERM {
		return false
	}

	if err := machine.Execute(); err != 0 {
		fmt.Println("Execution error")
		return false
	}

	return true
}
