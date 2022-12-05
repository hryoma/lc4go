package machine

const CODE_SIZE = 65536
const MEM_SIZE = 65536
const NUM_REGS = 8

const (
	OpNOP		= iota
	OpBRp		= iota
	OpBRz		= iota
	OpBRzp		= iota
	OpBRn		= iota
	OpBRnp		= iota
	OpBRnz		= iota
	OpBRnzp		= iota
	OpADD		= iota
	OpMUL		= iota
	OpSUB		= iota
	OpDIV		= iota
	OpADDI		= iota
	OpCMP		= iota
	OpCMPU		= iota
	OpCMPI		= iota
	OpCMPIU		= iota
	OpJSRR		= iota
	OpJSR		= iota
	OpAND		= iota
	OpNOT		= iota
	OpOR		= iota
	OpXOR		= iota
	OpANDI		= iota
	OpLDR		= iota
	OpSTR		= iota
	OpRTI		= iota
	OpCONST		= iota
	OpSLL		= iota
	OpSRA		= iota
	OpSRL		= iota
	OpMOD		= iota
	OpJMPR		= iota
	OpJMP		= iota
	OpHICONST	= iota
	OpTRAP		= iota
	// psuedo instructions
	OpRET
	OpLEA
	OpLC
)

type Insn struct {
	Data		uint16
	Op			int
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

