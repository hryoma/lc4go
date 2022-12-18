package machine

import ("fmt")

const MEM_SIZE = 65536
const NUM_REGS = 8

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
	// psuedo instructions
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
	Data		uint16
	OpName		Op
	Rd			uint8
	Rs			uint8
	Rt			uint8
	Imm			uint16
	Name		string
	Breakpoint	bool
}

func (insn Insn) String() string {
	// data
	// OpName: Rd, Rs, Rt, Imm
	return fmt.Sprintf("%016b\n%s: R%d, R%d, R%d, %d", insn.Data, insn.OpName, insn.Rd, insn.Rs, insn.Rt, insn.Imm)
}

type Machine struct {
	Mem		[MEM_SIZE]uint16
	Reg		[NUM_REGS]uint16
	Nzp		int8
	Psr		uint16
	Pc		uint16
	Labels	map[string]uint16
}

var Lc4 Machine

func wordToInsn(addr uint16) (insn Insn) {
	word := Lc4.Mem[addr]
	opCode := word >> 12

	var op Op
	var rd uint8
	var rs uint8
	var rt uint8
	var imm uint16
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

		imm = Lc4.Mem[addr] & 0x01FF
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
			imm = Lc4.Mem[addr] & 0x001F
			break parse_opcode
		}

		rd = uint8(Lc4.Mem[addr] >> 9) & 0b0111
		rs = uint8(Lc4.Mem[addr] >> 6) & 0b0111
		rt = uint8(Lc4.Mem[addr]) & 0b0111
	case 0b1010:
		// MOD or shift instructions
		rd = uint8(Lc4.Mem[addr] >> 9) & 0b0111
		rs = uint8(Lc4.Mem[addr] >> 6) & 0b0111

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

		imm = Lc4.Mem[addr] & 0x000F
	case 0b0101:
		// boolean instructions
		rd = uint8(Lc4.Mem[addr] >> 9) & 0b0111
		rs = uint8(Lc4.Mem[addr] >> 6) & 0b0111

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
			imm = Lc4.Mem[addr] & 0x001F
			break parse_opcode
		}

		rt = uint8(Lc4.Mem[addr]) & 0b0111
	case 0b0110:
		// LDR
		op = OpLDR
		rd = uint8(Lc4.Mem[addr] >> 9) & 0b0111
		rs = uint8(Lc4.Mem[addr] >> 6) & 0b0111
		imm = Lc4.Mem[addr] & 0x003F
	case 0b0111:
		// STR
		op = OpSTR
		rt = uint8(Lc4.Mem[addr] >> 9) & 0b0111
		rs = uint8(Lc4.Mem[addr] >> 6) & 0b0111
		imm = Lc4.Mem[addr] & 0x003F
	case 0b1001:
		// CONST
		rd = uint8(Lc4.Mem[addr] >> 9) & 0b0111
		imm = Lc4.Mem[addr] & 0x01FF
	case 0b1101:
		// HICONST
		rd = uint8(Lc4.Mem[addr] >> 9) & 0b0111
		imm = Lc4.Mem[addr] & 0x00FF
	case 0b0010:
		// comparison instructions
		rs = uint8(Lc4.Mem[addr] >> 9) & 0b0111

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
			imm = Lc4.Mem[addr] & 0x00FF
		case 0b11:
			op = OpCMPIU
			imm = Lc4.Mem[addr] & 0x00FF
		}
	case 0b0100:
		// JSRR, JSR
		subOpCode := (word >> 11) & 0b0001
		switch subOpCode {
		case 0b0:
			op = OpJSRR
			rs = uint8(Lc4.Mem[addr] >> 6) & 0b0111
		case 0b1:
			op = OpJSR
			imm = Lc4.Mem[addr] & 0x07FF
		}
	case 0b1100:
		// JMPR, JMP
		subOpCode := (word >> 11) & 0b0001
		switch subOpCode {
		case 0b0:
			op = OpJSRR
			rs = uint8(Lc4.Mem[addr] >> 6) & 0b0111
		case 0b1:
			op = OpJSR
			imm = Lc4.Mem[addr] & 0x07FF
		}
	case 0b1111:
		// TRAP
		imm = Lc4.Mem[addr] & 0x00FF
	case 0b1000:
		// RTI
	}

	return Insn{
		Data: Lc4.Mem[addr],
		OpName: op,
		Rd: rd,
		Rs: rs,
		Rt: rt,
		Imm: imm,
		Name: name,
	}
}

func signExtN(data uint16, nBits uint16) (signExtData uint16) {
	// get the sign and generate a mask
	var sign uint16 = data & (1 << (nBits - 1))
	var mask uint16 = (0xFF << nBits)

	// sign extend it
	if sign == 0 {
		signExtData = data & ^mask
	} else {
		signExtData = data | mask
	}

	return
}

func setNzp(testVal int16) {
	// reset nzp bits to 0's
	Lc4.Psr &= 0xFFF8

	if testVal < 0 {
		Lc4.Psr |= 0b0100
	} else if testVal == 0 {
		Lc4.Psr |= 0b0010
	} else {
		Lc4.Psr |= 0b0001
	}
}

func Execute() (err int) {
	// check if PC is at the end
	if Lc4.Pc == 0x80FF {
		return 0
	} else if Lc4.Pc < 0x0000 {
		// TODO - update these with the correct values
		// TODO - check for data/code region, privilege bit, etc.
		// PC is out of bounds
		return -1
	} else if Lc4.Pc > 0xFFFF {
		// PC is out of bounds
		return -1
	}

	insn := wordToInsn(Lc4.Pc)
	if insn.OpName != OpNOP {
		fmt.Printf("%s\n", insn)
	}

	switch insn.OpName {
		// branch instructions
		case OpNOP:
			// PC = PC + 1
			Lc4.Pc += 1
		case OpBRp:
			// if P, PC = PC + 1 + sext(IMM9)
			Lc4.Pc += 1
			if Lc4.Nzp > 0 {
				Lc4.Pc = uint16(int32(Lc4.Pc) + int32(signExtN(insn.Data, 9)))
			}
		case OpBRz:
			// if Z, PC = PC + 1 + sext(IMM9)
			Lc4.Pc += 1
			if Lc4.Nzp == 0 {
				Lc4.Pc = uint16(int32(Lc4.Pc) + int32(signExtN(insn.Data, 9)))
			}
		case OpBRzp:
			// if Z/P, PC = PC + 1 + sext(IMM9)
			Lc4.Pc += 1
			if Lc4.Nzp >= 0 {
				Lc4.Pc = uint16(int32(Lc4.Pc) + int32(signExtN(insn.Data, 9)))
			}
		case OpBRn:
			// if N, PC = PC + 1 + sext(IMM9)
			Lc4.Pc += 1
			if Lc4.Nzp < 0 {
				Lc4.Pc = uint16(int32(Lc4.Pc) + int32(signExtN(insn.Data, 9)))
			}
		case OpBRnp:
			// if NP, PC = PC + 1 + sext(IMM9)
			Lc4.Pc += 1
			if Lc4.Nzp != 0 {
				Lc4.Pc = uint16(int32(Lc4.Pc) + int32(signExtN(insn.Data, 9)))
			}
		case OpBRnz:
			// if NZ, PC = PC + 1 + sext(IMM9)
			Lc4.Pc += 1
			if Lc4.Nzp <= 0 {
				Lc4.Pc = uint16(int32(Lc4.Pc) + int32(signExtN(insn.Data, 9)))
			}
		case OpBRnzp:
			// if NZP, PC = PC + 1 + sext(IMM9)
			Lc4.Pc += 1
			Lc4.Pc = uint16(int32(Lc4.Pc) + int32(signExtN(insn.Data, 9)))
		// arithmetic instructions
		case OpADD:
			// Rd = Rs + Rt
			Lc4.Reg[insn.Rd] = uint16(int32(Lc4.Reg[insn.Rs]) + int32(Lc4.Reg[insn.Rt]))
			setNzp(int16(Lc4.Reg[insn.Rd]))
			Lc4.Pc += 1
		case OpMUL:
			// Rd = Rs * Rt
			Lc4.Reg[insn.Rd] = uint16(int32(Lc4.Reg[insn.Rs]) * int32(Lc4.Reg[insn.Rt]))
			setNzp(int16(Lc4.Reg[insn.Rd]))
			Lc4.Pc += 1
		case OpSUB:
			// Rd = Rs - Rt
			Lc4.Reg[insn.Rd] = uint16(int32(Lc4.Reg[insn.Rs]) - int32(Lc4.Reg[insn.Rt]))
			setNzp(int16(Lc4.Reg[insn.Rd]))
			Lc4.Pc += 1
		case OpDIV:
			// Rd = Rs / Rt
			Lc4.Reg[insn.Rd] = uint16(int32(Lc4.Reg[insn.Rs]) / int32(Lc4.Reg[insn.Rt]))
			setNzp(int16(Lc4.Reg[insn.Rd]))
			Lc4.Pc += 1
			// TODO error handling for div by 0
		case OpADDI:
			// Rd = Rs + sext(IMM5)
			imm := signExtN(insn.Data, 5)
			Lc4.Reg[insn.Rd] = uint16(int32(Lc4.Reg[insn.Rs]) + int32(imm))
			setNzp(int16(Lc4.Reg[insn.Rd]))
			Lc4.Pc += 1
		case OpMOD:
			// Rd = Rs % Rt
			Lc4.Reg[insn.Rd] = uint16(int32(Lc4.Reg[insn.Rs]) % int32(Lc4.Reg[insn.Rt]))
			setNzp(int16(Lc4.Reg[insn.Rd]))
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
			imm := signExtN(insn.Data, 5)
			Lc4.Reg[insn.Rd] = uint16(int32(Lc4.Reg[insn.Rs]) & int32(imm))
			setNzp(int16(Lc4.Reg[insn.Rd]))
			Lc4.Pc += 1
		// mem instructions
		case OpLDR:
			// Rd = dmem[Rs + sext(IMM6)]
			imm := signExtN(insn.Data, 6)
			Lc4.Reg[insn.Rd] = Lc4.Mem[uint16(int32(insn.Rs) + int32(imm))]
			setNzp(int16(Lc4.Reg[insn.Rd]))
			Lc4.Pc += 1
		case OpSTR:
			// dmem[Rs + sext(IMM6)] = Rt
			imm := signExtN(insn.Data, 6)
			Lc4.Mem[uint16(int32(insn.Rs) + int32(imm))] = Lc4.Reg[insn.Rt]
			Lc4.Pc += 1
		case OpCONST:
			// Rd = sext(IMM9)
			imm := signExtN(insn.Data, 9)
			Lc4.Reg[insn.Rd] = uint16(imm)
			setNzp(int16(Lc4.Reg[insn.Rd]))
			Lc4.Pc += 1
		case OpHICONST:
			// Rd = (Rd & 0xFF) | (UIMM8 << 8)
			imm := signExtN(insn.Data, 8)
			Lc4.Reg[insn.Rd] = (Lc4.Reg[insn.Rd] & 0xFF) | (uint16(imm) << 8)
			setNzp(int16(Lc4.Reg[insn.Rd]))
			Lc4.Pc += 1
		// comparison instructions
		case OpCMP:
			// NZP = sign(Rs - Rt)
			diff := int32(Lc4.Reg[insn.Rs]) - int32(Lc4.Reg[insn.Rt])
			setNzp(int16(diff))
			Lc4.Pc += 1
		case OpCMPU:
			// NZP = sign(uRs - uRt)
			diff := int32(Lc4.Reg[insn.Rs] - Lc4.Reg[insn.Rt])
			setNzp(int16(diff))
			Lc4.Pc += 1
		case OpCMPI:
			// NZP = sign(Rs - IMM7)
			imm := signExtN(insn.Data, 7)
			diff := int32(Lc4.Reg[insn.Rs]) - int32(imm)
			setNzp(int16(diff))
			Lc4.Pc += 1
		case OpCMPIU:
			// NZP = sign(uRs - UIMM7)
			imm := signExtN(insn.Data, 7)
			diff := int32(Lc4.Reg[insn.Rs] - uint16(imm))
			setNzp(int16(diff))
			Lc4.Pc += 1
		// shift instructions
		case OpSLL:
			// Rd = Rs << UIMM4
			imm := 0x0F & insn.Data
			Lc4.Reg[insn.Rd] = Lc4.Reg[insn.Rs] << imm
			setNzp(int16(Lc4.Reg[insn.Rd]))
			Lc4.Pc += 1
		case OpSRA:
			// Rd = Rs >>> UIMM4
			imm := 0x0F & insn.Data
			Lc4.Reg[insn.Rd] = uint16(int32(Lc4.Reg[insn.Rs]) << 16 >> 16 >> imm)
			setNzp(int16(Lc4.Reg[insn.Rd]))
			Lc4.Pc += 1
		case OpSRL:
			// Rd = Rs >> UIMM4
			imm := 0x0F & insn.Data
			Lc4.Reg[insn.Rd] = Lc4.Reg[insn.Rs] >> imm
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
			Lc4.Pc = (Lc4.Pc & 0x8000) | uint16(signExtN(insn.Data, 11) << 4)
		case OpJMPR:
			// PC = Rs
			Lc4.Pc = Lc4.Reg[insn.Rs]
		case OpJMP:
			// PC = PC + 1 + sext(IMM11)
			imm := signExtN(insn.Data, 11)
			Lc4.Pc = uint16(int32(Lc4.Pc) + 1 + int32(imm))
		// privilege instructions
		case OpTRAP:
			// R7 = PC + 1; PC = (0x8000 | IMM8); PSR[15] = 1
			imm := signExtN(insn.Data, 8)
			Lc4.Reg[7] = Lc4.Pc + 1
			setNzp(int16(Lc4.Reg[7]))
			Lc4.Pc = (0x8000 | imm)
			Lc4.Psr |= 0x8000
		case OpRTI:
			// PC = R7; PSR[15] = 0
			Lc4.Pc = Lc4.Reg[7]
			Lc4.Psr = Lc4.Psr & 0x7FFF
		// psuedo instructions:
		case OpRET:
			// return to R7
		case OpLEA:
			// store address of <Label> in Rd
		case OpLC:
			// store value of constant <Label> in Rd
		default:
			fmt.Println("Unexpected op code")
	}

	return 0
}

