package vm

import (
	"fmt"
)

// OpCode is an EVM opcode
type OpCode byte

func (op OpCode) IsPush() bool {
	switch op {
	case PUSH1, PUSH2, PUSH3, PUSH4, PUSH5, PUSH6, PUSH7, PUSH8, PUSH9, PUSH10, PUSH11, PUSH12, PUSH13, PUSH14, PUSH15, PUSH16, PUSH17, PUSH18, PUSH19, PUSH20, PUSH21, PUSH22, PUSH23, PUSH24, PUSH25, PUSH26, PUSH27, PUSH28, PUSH29, PUSH30, PUSH31, PUSH32:
		return true
	}
	return false
}

func (op OpCode) IsStaticJump() bool {
	return op == JUMP
}

const (
	// 0x0 range - arithmetic ops
	STOP OpCode = iota
	ADD
	MUL
	SUB
	DIV
	SDIV
	MOD
	SMOD
	ADDMOD
	MULMOD
	EXP
	SIGNEXTEND
)

const (
	LT OpCode = iota + 0x10
	GT
	SLT
	SGT
	EQ
	ISZERO
	AND
	OR
	XOR
	NOT
	BYTE
	SHL
	SHR
	SAR

	SHA3 = 0x20
)

const (
	// 0x30 range - closure state
	ADDRESS OpCode = 0x30 + iota
	BALANCE
	ORIGIN
	CALLER
	CALLVALUE
	CALLDATALOAD
	CALLDATASIZE
	CALLDATACOPY
	CODESIZE
	CODECOPY
	GASPRICE
	EXTCODESIZE
	EXTCODECOPY
	RETURNDATASIZE
	RETURNDATACOPY
)

const (
	// 0x40 range - block operations
	BLOCKHASH OpCode = 0x40 + iota
	COINBASE
	TIMESTAMP
	NUMBER
	DIFFICULTY
	GASLIMIT
)

const (
	// 0x50 range - 'storage' and execution
	POP OpCode = 0x50 + iota
	MLOAD
	MSTORE
	MSTORE8
	SLOAD
	SSTORE
	JUMP
	JUMPI
	PC
	MSIZE
	GAS
	JUMPDEST
)

const (
	// 0x60 range
	PUSH1 OpCode = 0x60 + iota
	PUSH2
	PUSH3
	PUSH4
	PUSH5
	PUSH6
	PUSH7
	PUSH8
	PUSH9
	PUSH10
	PUSH11
	PUSH12
	PUSH13
	PUSH14
	PUSH15
	PUSH16
	PUSH17
	PUSH18
	PUSH19
	PUSH20
	PUSH21
	PUSH22
	PUSH23
	PUSH24
	PUSH25
	PUSH26
	PUSH27
	PUSH28
	PUSH29
	PUSH30
	PUSH31
	PUSH32
	DUP1
	DUP2
	DUP3
	DUP4
	DUP5
	DUP6
	DUP7
	DUP8
	DUP9
	DUP10
	DUP11
	DUP12
	DUP13
	DUP14
	DUP15
	DUP16
	SWAP1
	SWAP2
	SWAP3
	SWAP4
	SWAP5
	SWAP6
	SWAP7
	SWAP8
	SWAP9
	SWAP10
	SWAP11
	SWAP12
	SWAP13
	SWAP14
	SWAP15
	SWAP16
)

const (
	LOG0 OpCode = 0xa0 + iota
	LOG1
	LOG2
	LOG3
	LOG4
)

// unofficial opcodes used for parsing
const (
	PUSH OpCode = 0xb0 + iota
	DUP
	SWAP
)

const (
	// 0xf0 range - closures
	CREATE OpCode = 0xf0 + iota
	CALL
	CALLCODE
	RETURN
	DELEGATECALL
	STATICCALL = 0xfa

	REVERT       = 0xfd
	SELFDESTRUCT = 0xff
)

// Since the opcodes aren't all in order we can't use a regular slice
var opCodeToString = map[OpCode]string{
	// 0x0 range - arithmetic ops
	STOP:       "STOP",
	ADD:        "ADD",
	MUL:        "MUL",
	SUB:        "SUB",
	DIV:        "DIV",
	SDIV:       "SDIV",
	MOD:        "MOD",
	SMOD:       "SMOD",
	EXP:        "EXP",
	NOT:        "NOT",
	LT:         "LT",
	GT:         "GT",
	SLT:        "SLT",
	SGT:        "SGT",
	EQ:         "EQ",
	ISZERO:     "ISZERO",
	SIGNEXTEND: "SIGNEXTEND",

	// 0x10 range - bit ops
	AND:    "AND",
	OR:     "OR",
	XOR:    "XOR",
	BYTE:   "BYTE",
	SHL:    "SHL",
	SHR:    "SHR",
	SAR:    "SAR",
	ADDMOD: "ADDMOD",
	MULMOD: "MULMOD",

	// 0x20 range - crypto
	SHA3: "SHA3",

	// 0x30 range - closure state
	ADDRESS:        "ADDRESS",
	BALANCE:        "BALANCE",
	ORIGIN:         "ORIGIN",
	CALLER:         "CALLER",
	CALLVALUE:      "CALLVALUE",
	CALLDATALOAD:   "CALLDATALOAD",
	CALLDATASIZE:   "CALLDATASIZE",
	CALLDATACOPY:   "CALLDATACOPY",
	CODESIZE:       "CODESIZE",
	CODECOPY:       "CODECOPY",
	GASPRICE:       "GASPRICE",
	EXTCODESIZE:    "EXTCODESIZE",
	EXTCODECOPY:    "EXTCODECOPY",
	RETURNDATASIZE: "RETURNDATASIZE",
	RETURNDATACOPY: "RETURNDATACOPY",

	// 0x40 range - block operations
	BLOCKHASH:  "BLOCKHASH",
	COINBASE:   "COINBASE",
	TIMESTAMP:  "TIMESTAMP",
	NUMBER:     "NUMBER",
	DIFFICULTY: "DIFFICULTY",
	GASLIMIT:   "GASLIMIT",

	// 0x50 range - 'storage' and execution
	POP: "POP",
	//DUP:     "DUP",
	//SWAP:    "SWAP",
	MLOAD:    "MLOAD",
	MSTORE:   "MSTORE",
	MSTORE8:  "MSTORE8",
	SLOAD:    "SLOAD",
	SSTORE:   "SSTORE",
	JUMP:     "JUMP",
	JUMPI:    "JUMPI",
	PC:       "PC",
	MSIZE:    "MSIZE",
	GAS:      "GAS",
	JUMPDEST: "JUMPDEST",

	// 0x60 range - push
	PUSH1:  "PUSH1",
	PUSH2:  "PUSH2",
	PUSH3:  "PUSH3",
	PUSH4:  "PUSH4",
	PUSH5:  "PUSH5",
	PUSH6:  "PUSH6",
	PUSH7:  "PUSH7",
	PUSH8:  "PUSH8",
	PUSH9:  "PUSH9",
	PUSH10: "PUSH10",
	PUSH11: "PUSH11",
	PUSH12: "PUSH12",
	PUSH13: "PUSH13",
	PUSH14: "PUSH14",
	PUSH15: "PUSH15",
	PUSH16: "PUSH16",
	PUSH17: "PUSH17",
	PUSH18: "PUSH18",
	PUSH19: "PUSH19",
	PUSH20: "PUSH20",
	PUSH21: "PUSH21",
	PUSH22: "PUSH22",
	PUSH23: "PUSH23",
	PUSH24: "PUSH24",
	PUSH25: "PUSH25",
	PUSH26: "PUSH26",
	PUSH27: "PUSH27",
	PUSH28: "PUSH28",
	PUSH29: "PUSH29",
	PUSH30: "PUSH30",
	PUSH31: "PUSH31",
	PUSH32: "PUSH32",

	DUP1:  "DUP1",
	DUP2:  "DUP2",
	DUP3:  "DUP3",
	DUP4:  "DUP4",
	DUP5:  "DUP5",
	DUP6:  "DUP6",
	DUP7:  "DUP7",
	DUP8:  "DUP8",
	DUP9:  "DUP9",
	DUP10: "DUP10",
	DUP11: "DUP11",
	DUP12: "DUP12",
	DUP13: "DUP13",
	DUP14: "DUP14",
	DUP15: "DUP15",
	DUP16: "DUP16",

	SWAP1:  "SWAP1",
	SWAP2:  "SWAP2",
	SWAP3:  "SWAP3",
	SWAP4:  "SWAP4",
	SWAP5:  "SWAP5",
	SWAP6:  "SWAP6",
	SWAP7:  "SWAP7",
	SWAP8:  "SWAP8",
	SWAP9:  "SWAP9",
	SWAP10: "SWAP10",
	SWAP11: "SWAP11",
	SWAP12: "SWAP12",
	SWAP13: "SWAP13",
	SWAP14: "SWAP14",
	SWAP15: "SWAP15",
	SWAP16: "SWAP16",
	LOG0:   "LOG0",
	LOG1:   "LOG1",
	LOG2:   "LOG2",
	LOG3:   "LOG3",
	LOG4:   "LOG4",

	// 0xf0 range
	CREATE:       "CREATE",
	CALL:         "CALL",
	RETURN:       "RETURN",
	CALLCODE:     "CALLCODE",
	DELEGATECALL: "DELEGATECALL",
	STATICCALL:   "STATICCALL",
	REVERT:       "REVERT",
	SELFDESTRUCT: "SELFDESTRUCT",

	PUSH: "PUSH",
	DUP:  "DUP",
	SWAP: "SWAP",
}

func (op OpCode) String() string {
	str := opCodeToString[op]
	if len(str) == 0 {
		return fmt.Sprintf("Missing opcode 0x%x", int(op))
	}

	return str
}

var stringToOp = map[string]OpCode{
	"STOP":           STOP,
	"ADD":            ADD,
	"MUL":            MUL,
	"SUB":            SUB,
	"DIV":            DIV,
	"SDIV":           SDIV,
	"MOD":            MOD,
	"SMOD":           SMOD,
	"EXP":            EXP,
	"NOT":            NOT,
	"LT":             LT,
	"GT":             GT,
	"SLT":            SLT,
	"SGT":            SGT,
	"EQ":             EQ,
	"ISZERO":         ISZERO,
	"SIGNEXTEND":     SIGNEXTEND,
	"AND":            AND,
	"OR":             OR,
	"XOR":            XOR,
	"BYTE":           BYTE,
	"SHL":            SHL,
	"SHR":            SHR,
	"SAR":            SAR,
	"ADDMOD":         ADDMOD,
	"MULMOD":         MULMOD,
	"SHA3":           SHA3,
	"ADDRESS":        ADDRESS,
	"BALANCE":        BALANCE,
	"ORIGIN":         ORIGIN,
	"CALLER":         CALLER,
	"CALLVALUE":      CALLVALUE,
	"CALLDATALOAD":   CALLDATALOAD,
	"CALLDATASIZE":   CALLDATASIZE,
	"CALLDATACOPY":   CALLDATACOPY,
	"DELEGATECALL":   DELEGATECALL,
	"STATICCALL":     STATICCALL,
	"CODESIZE":       CODESIZE,
	"CODECOPY":       CODECOPY,
	"GASPRICE":       GASPRICE,
	"EXTCODESIZE":    EXTCODESIZE,
	"EXTCODECOPY":    EXTCODECOPY,
	"RETURNDATASIZE": RETURNDATASIZE,
	"RETURNDATACOPY": RETURNDATACOPY,
	"BLOCKHASH":      BLOCKHASH,
	"COINBASE":       COINBASE,
	"TIMESTAMP":      TIMESTAMP,
	"NUMBER":         NUMBER,
	"DIFFICULTY":     DIFFICULTY,
	"GASLIMIT":       GASLIMIT,
	"POP":            POP,
	"MLOAD":          MLOAD,
	"MSTORE":         MSTORE,
	"MSTORE8":        MSTORE8,
	"SLOAD":          SLOAD,
	"SSTORE":         SSTORE,
	"JUMP":           JUMP,
	"JUMPI":          JUMPI,
	"PC":             PC,
	"MSIZE":          MSIZE,
	"GAS":            GAS,
	"JUMPDEST":       JUMPDEST,
	"PUSH1":          PUSH1,
	"PUSH2":          PUSH2,
	"PUSH3":          PUSH3,
	"PUSH4":          PUSH4,
	"PUSH5":          PUSH5,
	"PUSH6":          PUSH6,
	"PUSH7":          PUSH7,
	"PUSH8":          PUSH8,
	"PUSH9":          PUSH9,
	"PUSH10":         PUSH10,
	"PUSH11":         PUSH11,
	"PUSH12":         PUSH12,
	"PUSH13":         PUSH13,
	"PUSH14":         PUSH14,
	"PUSH15":         PUSH15,
	"PUSH16":         PUSH16,
	"PUSH17":         PUSH17,
	"PUSH18":         PUSH18,
	"PUSH19":         PUSH19,
	"PUSH20":         PUSH20,
	"PUSH21":         PUSH21,
	"PUSH22":         PUSH22,
	"PUSH23":         PUSH23,
	"PUSH24":         PUSH24,
	"PUSH25":         PUSH25,
	"PUSH26":         PUSH26,
	"PUSH27":         PUSH27,
	"PUSH28":         PUSH28,
	"PUSH29":         PUSH29,
	"PUSH30":         PUSH30,
	"PUSH31":         PUSH31,
	"PUSH32":         PUSH32,
	"DUP1":           DUP1,
	"DUP2":           DUP2,
	"DUP3":           DUP3,
	"DUP4":           DUP4,
	"DUP5":           DUP5,
	"DUP6":           DUP6,
	"DUP7":           DUP7,
	"DUP8":           DUP8,
	"DUP9":           DUP9,
	"DUP10":          DUP10,
	"DUP11":          DUP11,
	"DUP12":          DUP12,
	"DUP13":          DUP13,
	"DUP14":          DUP14,
	"DUP15":          DUP15,
	"DUP16":          DUP16,
	"SWAP1":          SWAP1,
	"SWAP2":          SWAP2,
	"SWAP3":          SWAP3,
	"SWAP4":          SWAP4,
	"SWAP5":          SWAP5,
	"SWAP6":          SWAP6,
	"SWAP7":          SWAP7,
	"SWAP8":          SWAP8,
	"SWAP9":          SWAP9,
	"SWAP10":         SWAP10,
	"SWAP11":         SWAP11,
	"SWAP12":         SWAP12,
	"SWAP13":         SWAP13,
	"SWAP14":         SWAP14,
	"SWAP15":         SWAP15,
	"SWAP16":         SWAP16,
	"LOG0":           LOG0,
	"LOG1":           LOG1,
	"LOG2":           LOG2,
	"LOG3":           LOG3,
	"LOG4":           LOG4,
	"CREATE":         CREATE,
	"CALL":           CALL,
	"RETURN":         RETURN,
	"CALLCODE":       CALLCODE,
	"REVERT":         REVERT,
	"SELFDESTRUCT":   SELFDESTRUCT,
}

func (op OpCode) OpCodeStackChange(stackPosition int) int {
	switch op {
	case SHL:
		return 1
	case SHR:
		return 1
	case SAR:
		return 1
	case STATICCALL:
		return 5
	case RETURNDATASIZE:
		return -1
	case RETURNDATACOPY:
		return 3
	case REVERT:
		return 2
	case STOP:
		return 0
	case ADD:
		return 1
	case MUL:
		return 1
	case SUB:
		return 1
	case DIV:
		return 1
	case SDIV:
		return 1
	case MOD:
		return 1
	case SMOD:
		return 1
	case ADDMOD:
		return 2
	case MULMOD:
		return 2
	case EXP:
		return 1
	case SIGNEXTEND:
		return 1
	case LT:
		return 1
	case GT:
		return 1
	case SLT:
		return 1
	case SGT:
		return 1
	case EQ:
		return 1
	case ISZERO:
		return 0
	case AND:
		return 1
	case XOR:
		return 1
	case OR:
		return 1
	case NOT:
		return 0
	case BYTE:
		return 1
	case SHA3:
		return 1
	case ADDRESS:
		return -1
	case BALANCE:
		return 0
	case ORIGIN:
		return -1
	case CALLER:
		return -1
	case CALLVALUE:
		return -1
	case CALLDATALOAD:
		return 0
	case CALLDATASIZE:
		return -1
	case CALLDATACOPY:
		return 3
	case CODESIZE:
		return -1
	case CODECOPY:
		return 3
	case GASPRICE:
		return -1
	case EXTCODESIZE:
		return 0
	case EXTCODECOPY:
		return 4
	case BLOCKHASH:
		return 0
	case COINBASE:
		return -1
	case TIMESTAMP:
		return -1
	case NUMBER:
		return -1
	case DIFFICULTY:
		return -1
	case GASLIMIT:
		return -1
	case POP:
		return 1
	case MLOAD:
		return 0
	case MSTORE:
		return 2
	case MSTORE8:
		return 2
	case SLOAD:
		return 0
	case SSTORE:
		return 2
	case JUMP:
		return 1
	case JUMPI:
		return 2
	case PC:
		return -1
	case MSIZE:
		return -1
	case GAS:
		return -1
	case JUMPDEST:
		return 0
	case PUSH1:
		return -1
	case PUSH2:
		return -1
	case PUSH3:
		return -1
	case PUSH4:
		return -1
	case PUSH5:
		return -1
	case PUSH6:
		return -1
	case PUSH7:
		return -1
	case PUSH8:
		return -1
	case PUSH9:
		return -1
	case PUSH10:
		return -1
	case PUSH11:
		return -1
	case PUSH12:
		return -1
	case PUSH13:
		return -1
	case PUSH14:
		return -1
	case PUSH15:
		return -1
	case PUSH16:
		return -1
	case PUSH17:
		return -1
	case PUSH18:
		return -1
	case PUSH19:
		return -1
	case PUSH20:
		return -1
	case PUSH21:
		return -1
	case PUSH22:
		return -1
	case PUSH23:
		return -1
	case PUSH24:
		return -1
	case PUSH25:
		return -1
	case PUSH26:
		return -1
	case PUSH27:
		return -1
	case PUSH28:
		return -1
	case PUSH29:
		return -1
	case PUSH30:
		return -1
	case PUSH31:
		return -1
	case PUSH32:
		return -1
	case DUP1:
		return -1
	case DUP2:
		return -1
	case DUP3:
		return -1
	case DUP4:
		return -1
	case DUP5:
		return -1
	case DUP6:
		return -1
	case DUP7:
		return -1
	case DUP8:
		return -1
	case DUP9:
		return -1
	case DUP10:
		return -1
	case DUP11:
		return -1
	case DUP12:
		return -1
	case DUP13:
		return -1
	case DUP14:
		return -1
	case DUP15:
		return -1
	case DUP16:
		return -1
	case SWAP1:
		if stackPosition == 0 {
			return -1
		}
		return 0
	case SWAP2:
		if stackPosition == 0 {
			return -2
		}
		return 0
	case SWAP3:
		if stackPosition == 0 {
			return -3
		}
		return 0
	case SWAP4:
		if stackPosition == 0 {
			return -4
		}
		return 0
	case SWAP5:
		if stackPosition == 0 {
			return -5
		}
		return 0
	case SWAP6:
		if stackPosition == 0 {
			return -6
		}
		return 0
	case SWAP7:
		if stackPosition == 0 {
			return -7
		}
		return 0
	case SWAP8:
		if stackPosition == 0 {
			return -8
		}
		return 0
	case SWAP9:
		if stackPosition == 0 {
			return -9
		}
		return 0
	case SWAP10:
		if stackPosition == 0 {
			return -10
		}
		return 0
	case SWAP11:
		if stackPosition == 0 {
			return -11
		}
		return 0
	case SWAP12:
		if stackPosition == 0 {
			return -12
		}
		return 0
	case SWAP13:
		if stackPosition == 0 {
			return -13
		}
		return 0
	case SWAP14:
		if stackPosition == 0 {
			return -14
		}
		return 0
	case SWAP15:
		if stackPosition == 0 {
			return -15
		}
		return 0
	case SWAP16:
		if stackPosition == 0 {
			return -16
		}
		return 0
	case LOG0:
		return 2
	case LOG1:
		return 3
	case LOG2:
		return 4
	case LOG3:
		return 5
	case LOG4:
		return 6
	case CREATE:
		return 2
	case CALL:
		return 6
	case CALLCODE:
		return 6
	case RETURN:
		return 2
	case SELFDESTRUCT:
		return 1
	default:
		return 0
	}
}

func StringToOp(str string) OpCode {
	return stringToOp[str]
}
