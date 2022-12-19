package machine

import (
	"fmt"
)

const MEM_SIZE = 65536
const NUM_REGS = 8

const USER_CODE_START = 0x0000
const USER_CODE_END = 0x1FFF
const USER_DATA_START = 0x2000
const USER_DATA_END = 0x7FFF
const OS_CODE_START = 0x8000
const OS_CODE_END = 0x9FFF
const OS_DATA_START = 0xA000
const OS_DATA_END = 0xFFFF

type Op int

const (
	OpNOP Op = iota
	OpBRp
	OpBRz
	OpBRzp
	OpBRn
	OpBRnp
	OpBRnz
	OpBRnzp
	OpADD
	OpMUL
	OpSUB
	OpDIV
	OpADDI
	OpCMP
	OpCMPU
	OpCMPI
	OpCMPIU
	OpJSRR
	OpJSR
	OpAND
	OpNOT
	OpOR
	OpXOR
	OpANDI
	OpLDR
	OpSTR
	OpRTI
	OpCONST
	OpSLL
	OpSRA
	OpSRL
	OpMOD
	OpJMPR
	OpJMP
	OpHICONST
	OpTRAP
	// pseudo instructions
	OpRET
	OpLEA
	OpLC
)

func (op Op) String() string {
	return [...]string{
		"NOP",
		"BRp",
		"BRz",
		"BRzp",
		"BRn",
		"BRnp",
		"BRnz",
		"BRnzp",
		"ADD",
		"MUL",
		"SUB",
		"DIV",
		"ADDI",
		"CMP",
		"CMPU",
		"CMPI",
		"CMPIU",
		"JSRR",
		"JSR",
		"AND",
		"NOT",
		"OR",
		"XOR",
		"ANDI",
		"LDR",
		"STR",
		"RTI",
		"CONST",
		"SLL",
		"SRA",
		"SRL",
		"MOD",
		"JMPR",
		"JMP",
		"HICONST",
		"TRAP",
		"RET",
		"LEA",
		"LC",
	}[op]
}

type Insn struct {
	Data   uint16
	OpName Op
	Rd     uint8
	Rs     uint8
	Rt     uint8
	Imm    int16
	Name   string
}

func (insn Insn) String() string {
	// data
	// OpName: Rd, Rs, Rt, Imm
	return fmt.Sprintf("%016b\n%s: R%d, R%d, R%d, %d", insn.Data, insn.OpName, insn.Rd, insn.Rs, insn.Rt, insn.Imm)
}

type MemMetadata struct {
	Label      string
	Breakpoint bool
}

type Machine struct {
	Mem    [MEM_SIZE]uint16
	Reg    [NUM_REGS]uint16
	Nzp    int8
	Psr    uint16
	Pc     uint16
	Labels map[string]uint16
	Meta   map[uint16]MemMetadata
}

var Lc4 Machine

func wordToInsn(addr uint16) (insn Insn) {
	word := Lc4.Mem[addr]
	opCode := word >> 12

	var op Op
	var rd uint8
	var rs uint8
	var rt uint8
	var imm int16
	var name string

parse_opcode:
	switch opCode {
	case 0b0000:
		// branch instructions
		subOpCode := (word >> 9) & 0b111

		switch subOpCode {
		case 0b000:
			op = OpNOP
			break parse_opcode
		case 0b001:
			op = OpBRp
		case 0b010:
			op = OpBRz
		case 0b011:
			op = OpBRzp
		case 0b100:
			op = OpBRn
		case 0b101:
			op = OpBRnp
		case 0b110:
			op = OpBRnz
		case 0b111:
			op = OpBRnzp
		}

		imm = signExtN(Lc4.Mem[addr]&0x01FF, 9)
	case 0b0001:
		// arithmetic instructions
		subOpCode := (word >> 3) & 0b111
		switch subOpCode {
		case 0b000:
			op = OpADD
		case 0b001:
			op = OpMUL
		case 0b010:
			op = OpSUB
		case 0b011:
			op = OpDIV
		default:
			op = OpADDI
			imm = signExtN(Lc4.Mem[addr]&0x001F, 5)
			break parse_opcode
		}

		rd = uint8(Lc4.Mem[addr]>>9) & 0b0111
		rs = uint8(Lc4.Mem[addr]>>6) & 0b0111
		rt = uint8(Lc4.Mem[addr]) & 0b0111
	case 0b1010:
		// MOD or shift instructions
		rd = uint8(Lc4.Mem[addr]>>9) & 0b0111
		rs = uint8(Lc4.Mem[addr]>>6) & 0b0111

		subOpCode := (word >> 4) & 0b11
		switch subOpCode {
		case 0b00:
			op = OpSLL
		case 0b01:
			op = OpSRA
		case 0b10:
			op = OpSRL
		case 0b11:
			op = OpMOD
			rt = uint8(Lc4.Mem[addr]) & 0b0111
			break parse_opcode
		}

		imm = int16(Lc4.Mem[addr]) & 0x000F
	case 0b0101:
		// boolean instructions
		rd = uint8(Lc4.Mem[addr]>>9) & 0b0111
		rs = uint8(Lc4.Mem[addr]>>6) & 0b0111

		subOpCode := (word >> 3) & 0b111
		switch subOpCode {
		case 0b000:
			op = OpAND
		case 0b001:
			op = OpNOT
			break parse_opcode
		case 0b010:
			op = OpOR
		case 0b011:
			op = OpXOR
		default:
			op = OpANDI
			imm = signExtN(Lc4.Mem[addr]&0x001F, 5)
			break parse_opcode
		}

		rt = uint8(Lc4.Mem[addr]) & 0b0111
	case 0b0110:
		// LDR
		op = OpLDR
		rd = uint8(Lc4.Mem[addr]>>9) & 0b0111
		rs = uint8(Lc4.Mem[addr]>>6) & 0b0111
		imm = signExtN(Lc4.Mem[addr]&0x003F, 6)
	case 0b0111:
		// STR
		op = OpSTR
		rt = uint8(Lc4.Mem[addr]>>9) & 0b0111
		rs = uint8(Lc4.Mem[addr]>>6) & 0b0111
		imm = signExtN(Lc4.Mem[addr]&0x003F, 6)
	case 0b1001:
		// CONST
		op = OpCONST
		rd = uint8(Lc4.Mem[addr]>>9) & 0b0111
		imm = signExtN(Lc4.Mem[addr]&0x003F, 6)
	case 0b1101:
		// HICONST
		op = OpHICONST
		rd = uint8(Lc4.Mem[addr]>>9) & 0b0111
		imm = int16(Lc4.Mem[addr]) & 0x003F
	case 0b0010:
		// comparison instructions
		rs = uint8(Lc4.Mem[addr]>>9) & 0b0111

		subOpCode := (word >> 7) & 0b0011
		switch subOpCode {
		case 0b00:
			op = OpCMP
			rt = uint8(Lc4.Mem[addr]) & 0b0111
		case 0b01:
			op = OpCMPU
			rt = uint8(Lc4.Mem[addr]) & 0b0111
		case 0b10:
			op = OpCMPI
			imm = signExtN(Lc4.Mem[addr]&0x00FF, 8)
		case 0b11:
			op = OpCMPIU
			imm = int16(Lc4.Mem[addr]) & 0x00FF
		}
	case 0b0100:
		// JSRR, JSR
		subOpCode := (word >> 11) & 0b0001
		switch subOpCode {
		case 0b0:
			op = OpJSRR
			rs = uint8(Lc4.Mem[addr]>>6) & 0b0111
		case 0b1:
			op = OpJSR
			imm = signExtN(Lc4.Mem[addr]&0x07FF, 11)
		}
	case 0b1100:
		// JMPR, JMP
		subOpCode := (word >> 11) & 0b0001
		switch subOpCode {
		case 0b0:
			op = OpJSRR
			rs = uint8(Lc4.Mem[addr]>>6) & 0b0111
		case 0b1:
			op = OpJSR
			imm = signExtN(Lc4.Mem[addr]&0x07FF, 11)
		}
	case 0b1111:
		// TRAP
		op = OpTRAP
		imm = int16(Lc4.Mem[addr]) & 0x00FF
	case 0b1000:
		// RTI
		op = OpRTI
		imm = signExtN(Lc4.Mem[addr]&0x00FF, 8)
	}

	return Insn{
		Data:   Lc4.Mem[addr],
		OpName: op,
		Rd:     rd,
		Rs:     rs,
		Rt:     rt,
		Imm:    imm,
		Name:   name,
	}
}

func signExtN(data uint16, nBits uint16) int16 {
	// get the sign and generate a mask
	var sign uint16 = data & (1 << (nBits - 1))
	var mask uint16 = (0xFF << nBits)

	// sign extend it
	if sign == 0 {
		return int16(data & ^mask)
	} else {
		return int16(data | mask)
	}
}

func setNzp(testVal int16) {
	// reset nzp bits to 0's
	Lc4.Psr &= 0xFFF8

	if testVal < 0 {
		Lc4.Psr |= 0b100
	} else if testVal == 0 {
		Lc4.Psr |= 0b010
	} else {
		Lc4.Psr |= 0b001
	}
}

func uintPlusInt(val uint16, offset int16) (res uint16, err int) {
	var tempRes int32 = int32(val) + int32(offset)
	if 0 <= tempRes && tempRes <= 0xFFFF {
		res = uint16(tempRes)
	} else {
		err = -1
	}
	return
}

func Execute() (err int) {
	if (USER_DATA_START <= Lc4.Pc) && (Lc4.Pc <= USER_DATA_END) {
		return -1
	} else if (OS_DATA_START <= Lc4.Pc) && (Lc4.Pc <= OS_DATA_END) {
		return -1
	} else if (OS_CODE_START <= Lc4.Pc) && (Lc4.Pc <= OS_CODE_END) {
		if (Lc4.Psr & 0x8000) == 0 {
			// os code section, ran with insufficient privilege
			return -1
		}
	}

	insn := wordToInsn(Lc4.Pc)
	switch insn.OpName {
	// branch instructions
	case OpNOP:
		// PC = PC + 1
		Lc4.Pc += 1
	case OpBRp:
		// if P, PC = PC + 1 + sext(IMM9)
		Lc4.Pc += 1
		if Lc4.Nzp > 0 {
			Lc4.Pc, err = uintPlusInt(Lc4.Pc, insn.Imm)
		}
	case OpBRz:
		// if Z, PC = PC + 1 + sext(IMM9)
		Lc4.Pc += 1
		if Lc4.Nzp == 0 {
			Lc4.Pc, err = uintPlusInt(Lc4.Pc, insn.Imm)
		}
	case OpBRzp:
		// if Z/P, PC = PC + 1 + sext(IMM9)
		Lc4.Pc += 1
		if Lc4.Nzp >= 0 {
			Lc4.Pc, err = uintPlusInt(Lc4.Pc, insn.Imm)
		}
	case OpBRn:
		// if N, PC = PC + 1 + sext(IMM9)
		Lc4.Pc += 1
		if Lc4.Nzp < 0 {
			Lc4.Pc, err = uintPlusInt(Lc4.Pc, insn.Imm)
		}
	case OpBRnp:
		// if NP, PC = PC + 1 + sext(IMM9)
		Lc4.Pc += 1
		if Lc4.Nzp != 0 {
			Lc4.Pc, err = uintPlusInt(Lc4.Pc, insn.Imm)
		}
	case OpBRnz:
		// if NZ, PC = PC + 1 + sext(IMM9)
		Lc4.Pc += 1
		if Lc4.Nzp <= 0 {
			Lc4.Pc, err = uintPlusInt(Lc4.Pc, insn.Imm)
		}
	case OpBRnzp:
		// if NZP, PC = PC + 1 + sext(IMM9)
		Lc4.Pc += 1
		Lc4.Pc, err = uintPlusInt(Lc4.Pc, insn.Imm)
	// arithmetic instructions
	case OpADD:
		// Rd = Rs + Rt
		calc := int16(Lc4.Reg[insn.Rs]) + int16(Lc4.Reg[insn.Rt])
		Lc4.Reg[insn.Rd] = uint16(calc)
		setNzp(calc)
		Lc4.Pc += 1
	case OpMUL:
		// Rd = Rs * Rt
		calc := int16(Lc4.Reg[insn.Rs]) * int16(Lc4.Reg[insn.Rt])
		Lc4.Reg[insn.Rd] = uint16(calc)
		setNzp(calc)
		Lc4.Pc += 1
	case OpSUB:
		// Rd = Rs - Rt
		calc := int16(Lc4.Reg[insn.Rs]) - int16(Lc4.Reg[insn.Rt])
		Lc4.Reg[insn.Rd] = uint16(calc)
		setNzp(calc)
		Lc4.Pc += 1
	case OpDIV:
		// Rd = Rs / Rt
		if Lc4.Reg[insn.Rt] == 0 {
			Lc4.Reg[insn.Rd] = 0
			setNzp(0)
		} else {
			calc := int16(Lc4.Reg[insn.Rs]) / int16(Lc4.Reg[insn.Rt])
			Lc4.Reg[insn.Rd] = uint16(calc)
			setNzp(calc)
		}
		Lc4.Pc += 1
	case OpADDI:
		// Rd = Rs + sext(IMM5)
		Lc4.Reg[insn.Rd] = uint16(int16(Lc4.Reg[insn.Rs]) + insn.Imm)
		setNzp(int16(Lc4.Reg[insn.Rd]))
		Lc4.Pc += 1
	case OpMOD:
		// Rd = Rs % Rt
		if Lc4.Reg[insn.Rt] == 0 {
			Lc4.Reg[insn.Rd] = 0
			setNzp(0)
		} else {
			calc := int16(Lc4.Reg[insn.Rs]) % int16(Lc4.Reg[insn.Rt])
			Lc4.Reg[insn.Rd] = uint16(calc)
			setNzp(calc)
		}
		Lc4.Pc += 1
		// TODO error handling for mod by 0
	// logical instructions
	case OpAND:
		// Rd = Rs & Rt
		Lc4.Reg[insn.Rd] = Lc4.Reg[insn.Rs] & Lc4.Reg[insn.Rt]
		setNzp(int16(Lc4.Reg[insn.Rd]))
		Lc4.Pc += 1
	case OpNOT:
		// Rd = ~Rs
		Lc4.Reg[insn.Rd] = 0xFFFF ^ Lc4.Reg[insn.Rs]
		setNzp(int16(Lc4.Reg[insn.Rd]))
		Lc4.Pc += 1
	case OpOR:
		// Rd = Rs | Rt
		Lc4.Reg[insn.Rd] = Lc4.Reg[insn.Rs] | Lc4.Reg[insn.Rt]
		setNzp(int16(Lc4.Reg[insn.Rd]))
		Lc4.Pc += 1
	case OpXOR:
		// Rd = Rs ^ Rt
		Lc4.Reg[insn.Rd] = Lc4.Reg[insn.Rs] ^ Lc4.Reg[insn.Rt]
		setNzp(int16(Lc4.Reg[insn.Rd]))
		Lc4.Pc += 1
	case OpANDI:
		// Rd = Rs & sext(IMM5)
		Lc4.Reg[insn.Rd] = uint16(int16(Lc4.Reg[insn.Rs]) & insn.Imm)
		setNzp(int16(Lc4.Reg[insn.Rd]))
		Lc4.Pc += 1
	// mem instructions
	case OpLDR:
		// Rd = dmem[Rs + sext(IMM6)]
		Lc4.Reg[insn.Rd] = Lc4.Mem[uint16(int16(insn.Rs)+insn.Imm)]
		setNzp(int16(Lc4.Reg[insn.Rd]))
		Lc4.Pc += 1
	case OpSTR:
		// dmem[Rs + sext(IMM6)] = Rt
		dmemAddr := uint16(int16(insn.Rs) + insn.Imm)

		// only write to data sections
		if (USER_DATA_START <= dmemAddr) && (dmemAddr <= USER_DATA_END) {
			Lc4.Mem[dmemAddr] = Lc4.Reg[insn.Rt]
		} else if (OS_DATA_END <= dmemAddr) && (dmemAddr <= OS_DATA_END) {
			// only write to os data if enough privilege
			if (Lc4.Psr & 0x8000) == 0 {
				Lc4.Mem[dmemAddr] = Lc4.Reg[insn.Rt]
			}
		}

		Lc4.Pc += 1
	case OpCONST:
		// Rd = sext(IMM9)
		Lc4.Reg[insn.Rd] = uint16(insn.Imm)
		setNzp(int16(Lc4.Reg[insn.Rd]))
		Lc4.Pc += 1
	case OpHICONST:
		// Rd = (Rd & 0xFF) | (UIMM8 << 8)
		Lc4.Reg[insn.Rd] = (Lc4.Reg[insn.Rd] & 0xFF) | (uint16(insn.Imm) << 8)
		setNzp(int16(Lc4.Reg[insn.Rd]))
		Lc4.Pc += 1
	// comparison instructions
	case OpCMP:
		// NZP = sign(Rs - Rt)
		setNzp(int16(Lc4.Reg[insn.Rs]) - int16(Lc4.Reg[insn.Rt]))
		Lc4.Pc += 1
	case OpCMPU:
		// NZP = sign(uRs - uRt)
		setNzp(int16(Lc4.Reg[insn.Rs] - Lc4.Reg[insn.Rt]))
		Lc4.Pc += 1
	case OpCMPI:
		// NZP = sign(Rs - IMM7)
		setNzp(int16(Lc4.Reg[insn.Rs]) - insn.Imm)
		Lc4.Pc += 1
	case OpCMPIU:
		// NZP = sign(uRs - UIMM7)
		setNzp(int16(Lc4.Reg[insn.Rs] - uint16(insn.Imm)))
		Lc4.Pc += 1
	// shift instructions
	case OpSLL:
		// Rd = Rs << UIMM4
		Lc4.Reg[insn.Rd] = Lc4.Reg[insn.Rs] << insn.Imm
		setNzp(int16(Lc4.Reg[insn.Rd]))
		Lc4.Pc += 1
	case OpSRA:
		// Rd = Rs >>> UIMM4
		Lc4.Reg[insn.Rd] = uint16(int32(Lc4.Reg[insn.Rs]) << 16 >> 16 >> insn.Imm)
		setNzp(int16(Lc4.Reg[insn.Rd]))
		Lc4.Pc += 1
	case OpSRL:
		// Rd = Rs >> UIMM4
		Lc4.Reg[insn.Rd] = Lc4.Reg[insn.Rs] >> insn.Imm
		setNzp(int16(Lc4.Reg[insn.Rd]))
		Lc4.Pc += 1
	// jjump instructions
	case OpJSRR:
		// R7 = PC + 1; PC = Rs
		temp_rs := Lc4.Reg[insn.Rs]
		Lc4.Reg[7] = Lc4.Pc + 1
		setNzp(int16(Lc4.Reg[7]))
		Lc4.Pc = uint16(temp_rs)
	case OpJSR:
		// R7 = PC + 1; PC = (PC & 0x8000) | (IMM11 << 4)
		Lc4.Reg[7] = Lc4.Pc + 1
		setNzp(int16(Lc4.Reg[7]))
		Lc4.Pc = (Lc4.Pc & 0x8000) | (uint16(insn.Imm) << 4)
	case OpJMPR:
		// PC = Rs
		Lc4.Pc = Lc4.Reg[insn.Rs]
	case OpJMP:
		// PC = PC + 1 + sext(IMM11)
		Lc4.Pc += 1
		Lc4.Pc, err = uintPlusInt(Lc4.Pc, insn.Imm)
	// privilege instructions
	case OpTRAP:
		// R7 = PC + 1; PC = (0x8000 | IMM8); PSR[15] = 1
		Lc4.Reg[7] = Lc4.Pc + 1
		setNzp(int16(Lc4.Reg[7]))
		Lc4.Pc = 0x8000 | uint16(insn.Imm)
		Lc4.Psr |= 0x8000
	case OpRTI:
		// PC = R7; PSR[15] = 0
		Lc4.Pc = Lc4.Reg[7]
		Lc4.Psr = Lc4.Psr & 0x7FFF
	// pseudo instructions:
	case OpRET:
		// return to R7
	case OpLEA:
		// store address of <Label> in Rd
	case OpLC:
		// store value of constant <Label> in Rd
	default:
		fmt.Println("Unexpected op code")
	}

	if err == -1 {
		fmt.Println("Arithmetic error: uintPlusInt")
		return -1
	}

	return 0
}
