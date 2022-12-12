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

const (
	nzpNONE	= 0b00000000
	nzpN	= 0b00000100
	nzpZ	= 0b00000010
	nzpP	= 0b00000001
	nzpNZ	= 0b00000110
	nzpNP	= 0b00000101
	nzpZP	= 0b00000011
	nzpNZP	= 0b00000111
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
	Reg		[NUM_REGS]uint16
	Nzp		int8
	Psr		uint16
	Pc		uint16
	Labels	map[string]uint16
}

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

func Execute(lc4 *Machine) (err int) {
	// check if PC is at the end
	if lc4.Pc == 0x80FF {
		return 0
	} else if lc4.Pc < 0x0000 {
		// TODO - update these with the correct values
		// TODO - check for data/code region, privilege bit, etc.
		// PC is out of bounds
		return -1
	} else if lc4.Pc > 0xFFFF {
		// PC is out of bounds
		return -1
	}

	insn := lc4.Code[lc4.Pc]
	switch insn.Op {
		case OpNOP:
			lc4.Pc += 1
		case OpBRp:
			lc4.Pc += 1
			if lc4.Nzp == nzpP {
				lc4.Pc += getSignExtN(insn.Data, 9)
			}
		case OpBRz:
		case OpBRzp:
		case OpBRn:
		case OpBRnp:
		case OpBRnz:
		case OpBRnzp:
		case OpADD:
		case OpMUL:
		case OpSUB:
		case OpDIV:
		case OpADDI:
		case OpCMP:
		case OpCMPU:
		case OpCMPI:
		case OpCMPIU:
		case OpJSRR:
		case OpJSR:
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

