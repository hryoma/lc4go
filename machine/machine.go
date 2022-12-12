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

