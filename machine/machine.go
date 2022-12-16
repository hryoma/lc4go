package machine

const CODE_SIZE = 65536
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

type Insn struct {
	Data		uint16
	Op			Op
	String		string
	Breakpoint	bool
}

type Machine struct {
	Code	[CODE_SIZE]Insn
	Mem		[MEM_SIZE]uint16
	Reg		[NUM_REGS]int16
	Nzp		int8
	Psr		uint16
	Pc		uint16
	Labels	map[string]uint16
}

var Lc4 Machine

func getSignExtN(data uint16, nBits uint16) (signExtData uint16) {
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

	insn := Lc4.Code[Lc4.Pc]
	switch insn.Op {
		// branch instructions
		case OpNOP:
			// PC = PC + 1
			Lc4.Pc += 1
		case OpBRp:
			// if P, PC = PC + 1 + sext(IMM9)
			if Lc4.Nzp > 0 {
				Lc4.Pc += 1
				Lc4.Pc = uint16(int32(Lc4.Pc) + int32(getSignExtN(insn.Data, 9)))
			}
		case OpBRz:
			// if Z, PC = PC + 1 + sext(IMM9)
			if Lc4.Nzp == 0 {
				Lc4.Pc += 1
				Lc4.Pc = uint16(int32(Lc4.Pc) + int32(getSignExtN(insn.Data, 9)))
			}
		case OpBRzp:
			// if Z/P, PC = PC + 1 + sext(IMM9)
			if Lc4.Nzp >= 0 {
				Lc4.Pc += 1
				Lc4.Pc = uint16(int32(Lc4.Pc) + int32(getSignExtN(insn.Data, 9)))
			}
		case OpBRn:
			// if N, PC = PC + 1 + sext(IMM9)
			if Lc4.Nzp < 0 {
				Lc4.Pc += 1
				Lc4.Pc = uint16(int32(Lc4.Pc) + int32(getSignExtN(insn.Data, 9)))
			}
		case OpBRnp:
			// if NP, PC = PC + 1 + sext(IMM9)
			if Lc4.Nzp != 0 {
				Lc4.Pc += 1
				Lc4.Pc = uint16(int32(Lc4.Pc) + int32(getSignExtN(insn.Data, 9)))
			}
		case OpBRnz:
			// if NZ, PC = PC + 1 + sext(IMM9)
			if Lc4.Nzp <= 0 {
				Lc4.Pc += 1
				Lc4.Pc = uint16(int32(Lc4.Pc) + int32(getSignExtN(insn.Data, 9)))
			}
		case OpBRnzp:
			// if NZP, PC = PC + 1 + sext(IMM9)
			Lc4.Pc += 1
			Lc4.Pc = uint16(int32(Lc4.Pc) + int32(getSignExtN(insn.Data, 9)))
		// arithmetic instructions
		case OpADD:
			// Rd = Rs + Rt
			rd := insn.Data & (0b0111 << 9)
			rs := insn.Data & (0b0111 << 6)
			rt := insn.Data & (0b0111)
			Lc4.Reg[rd] = Lc4.Reg[rs] + Lc4.Reg[rt]
		case OpMUL:
			// Rd = Rs * Rt
			rd := insn.Data & (0b0111 << 9)
			rs := insn.Data & (0b0111 << 6)
			rt := insn.Data & (0b0111)
			Lc4.Reg[rd] = Lc4.Reg[rs] * Lc4.Reg[rt]
		case OpSUB:
			// Rd = Rs - Rt
			rd := insn.Data & (0b0111 << 9)
			rs := insn.Data & (0b0111 << 6)
			rt := insn.Data & (0b0111)
			Lc4.Reg[rd] = Lc4.Reg[rs] - Lc4.Reg[rt]
		case OpDIV:
			// Rd = Rs / Rt
			rd := insn.Data & (0b0111 << 9)
			rs := insn.Data & (0b0111 << 6)
			rt := insn.Data & (0b0111)
			Lc4.Reg[rd] = Lc4.Reg[rs] / Lc4.Reg[rt]
			// TODO error handling for div by 0
		case OpADDI:
			// Rd = Rs + sext(IMM5)
			rd := insn.Data & (0b0111 << 9)
			rs := insn.Data & (0b0111 << 6)
			imm := getSignExtN(insn.Data, 5)
			Lc4.Reg[rd] = int16(int32(Lc4.Reg[rs]) + int32(imm))
		case OpMOD:
			// Rd = Rs % Rt
		// logical instructions
		case OpAND:
			// Rd = Rs & Rt
		case OpNOT:
			// Rd = ~Rs
		case OpOR:
			// Rd = Rs | Rt
		case OpXOR:
			// Rd = Rs ^ Rt
		case OpANDI:
			// Rd = Rs & sext(IMM5)
		// mem instructions
		case OpLDR:
			// Rd = dmem[Rs + sext(IMM6)]
		case OpSTR:
			// dmem[Rs + sext(IMM6)] = Rt
		case OpCONST:
			// Rd = sext(IMM9)
		case OpHICONST:
			// Rd = (Rd & 0xFF) | (UIMM8 << 8)
		// comparison instructions
		case OpCMP:
			// NZP = sign(Rs - Rt)
			rs := insn.Data & (0b0111 << 9)
			rt := insn.Data & (0b0111)
			diff := Lc4.Reg[rs] - Lc4.Reg[rt]
			if diff < 0 {
				Lc4.Nzp = -1
			} else if diff == 0 {
				Lc4.Nzp = 0
			} else {
				Lc4.Nzp = 1
			}
		case OpCMPU:
			// NZP = sign(uRs - uRt)
			rs := insn.Data & (0b0111 << 9)
			rt := insn.Data & (0b0111)
			diff := int16(uint16(Lc4.Reg[rs]) - uint16(Lc4.Reg[rt]))
			if diff < 0 {
				Lc4.Nzp = -1
			} else if diff == 0 {
				Lc4.Nzp = 0
			} else {
				Lc4.Nzp = 1
			}
		case OpCMPI:
			// NZP = sign(Rs - IMM7)
			rs := insn.Data & (0b0111 << 9)
			imm := getSignExtN(insn.Data, 7)
			diff := Lc4.Reg[rs] - int16(imm)
			if diff < 0 {
				Lc4.Nzp = -1
			} else if diff == 0 {
				Lc4.Nzp = 0
			} else {
				Lc4.Nzp = 1
			}
		case OpCMPIU:
			// NZP = sign(uRs - UIMM7)
			rs := insn.Data & (0b0111 << 9)
			imm := getSignExtN(insn.Data, 7)
			diff := int16(uint16(Lc4.Reg[rs]) - imm)
			if diff < 0 {
				Lc4.Nzp = -1
			} else if diff == 0 {
				Lc4.Nzp = 0
			} else {
				Lc4.Nzp = 1
			}
		// shift instructions
		case OpSLL:
			// Rd = Rs << UIMM4
		case OpSRA:
			// Rd = Rs >>> UIMM4
		case OpSRL:
			// Rd = Rs >> UIMM4
		// jjump instructions
		case OpJSRR:
			// R7 = PC + 1; PC = Rs
			rsVal := Lc4.Reg[insn.Data & (0b0111 << 6)]
			Lc4.Reg[7] = int16(Lc4.Pc) + 1
			Lc4.Pc = uint16(rsVal)
		case OpJSR:
			// R7 = PC + 1; PC = (PC & 0x8000) | (IMM11 << 4)
			Lc4.Reg[7] = int16(Lc4.Pc) + 1
			Lc4.Pc = (Lc4.Pc & 0x8000) | uint16(getSignExtN(insn.Data, 11) << 4)
		case OpJMPR:
			// PC = Rs
		case OpJMP:
			// PC = PC + 1 + sext(IMM11)
		// privilege instructions
		case OpTRAP:
			// R7 = PC + 1; PC = (0x8000 | IMM8); PSR[15] = 1
		case OpRTI:
			// PC = R7; PSR[15] = 0
		// psuedo instructions:
		case OpRET:
			// return to R7
		case OpLEA:
			// store address of <Label> in Rd
		case OpLC:
			// store value of constant <Label> in Rd
	}

	return 0
}

