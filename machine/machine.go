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
		case OpNOP:
			Lc4.Pc += 1
		case OpBRp:
			Lc4.Pc += 1
			if Lc4.Nzp > 0 {
				Lc4.Pc = uint16(int32(Lc4.Pc) + int32(getSignExtN(insn.Data, 9)))
			}
		case OpBRz:
			Lc4.Pc += 1
			if Lc4.Nzp == 0 {
				Lc4.Pc = uint16(int32(Lc4.Pc) + int32(getSignExtN(insn.Data, 9)))
			}
		case OpBRzp:
			Lc4.Pc += 1
			if Lc4.Nzp >= 0 {
				Lc4.Pc = uint16(int32(Lc4.Pc) + int32(getSignExtN(insn.Data, 9)))
			}
		case OpBRn:
			Lc4.Pc += 1
			if Lc4.Nzp < 0 {
				Lc4.Pc = uint16(int32(Lc4.Pc) + int32(getSignExtN(insn.Data, 9)))
			}
		case OpBRnp:
			Lc4.Pc += 1
			if Lc4.Nzp != 0 {
				Lc4.Pc = uint16(int32(Lc4.Pc) + int32(getSignExtN(insn.Data, 9)))
			}
		case OpBRnz:
			Lc4.Pc += 1
			if Lc4.Nzp <= 0 {
				Lc4.Pc = uint16(int32(Lc4.Pc) + int32(getSignExtN(insn.Data, 9)))
			}
		case OpBRnzp:
			Lc4.Pc += 1
			Lc4.Pc = uint16(int32(Lc4.Pc) + int32(getSignExtN(insn.Data, 9)))
		case OpADD:
			rd := insn.Data & (0b0111 << 9)
			rs := insn.Data & (0b0111 << 6)
			rt := insn.Data & (0b0111)
			Lc4.Reg[rd] = Lc4.Reg[rs] + Lc4.Reg[rt]
		case OpMUL:
			rd := insn.Data & (0b0111 << 9)
			rs := insn.Data & (0b0111 << 6)
			rt := insn.Data & (0b0111)
			Lc4.Reg[rd] = Lc4.Reg[rs] * Lc4.Reg[rt]
		case OpSUB:
			rd := insn.Data & (0b0111 << 9)
			rs := insn.Data & (0b0111 << 6)
			rt := insn.Data & (0b0111)
			Lc4.Reg[rd] = Lc4.Reg[rs] - Lc4.Reg[rt]
		case OpDIV:
			rd := insn.Data & (0b0111 << 9)
			rs := insn.Data & (0b0111 << 6)
			rt := insn.Data & (0b0111)
			Lc4.Reg[rd] = Lc4.Reg[rs] / Lc4.Reg[rt]
			// TODO error handling for div by 0
		case OpADDI:
			rd := insn.Data & (0b0111 << 9)
			rs := insn.Data & (0b0111 << 6)
			imm := getSignExtN(insn.Data, 5)
			Lc4.Reg[rd] = int16(int32(Lc4.Reg[rs]) + int32(imm))
		case OpCMP:
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
		case OpJSRR:
			rsVal := Lc4.Reg[insn.Data & (0b0111 << 6)]
			Lc4.Reg[7] = int16(Lc4.Pc) + 1
			Lc4.Pc = uint16(rsVal)
		case OpJSR:
			Lc4.Reg[7] = int16(Lc4.Pc) + 1
			Lc4.Pc = (Lc4.Pc & 0x8000) | uint16(getSignExtN(insn.Data, 11) << 4)
		case OpAND:
		case OpNOT:
		case OpOR:
		case OpXOR:
		case OpANDI:
		case OpLDR:
		case OpSTR:
		case OpRTI:
		case OpCONST:
		case OpSLL:
		case OpSRA:
		case OpSRL:
		case OpMOD:
		case OpJMPR:
		case OpJMP:
		case OpHICONST:
		case OpTRAP:
		// psuedo instructions:
		case OpRET:
		case OpLEA:
		case OpLC:
	}

	return 0
}

