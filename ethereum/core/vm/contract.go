package vm

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tenderly/tenderly-trace/ethereum/core/types"
	"math/big"
)

type ContractRef interface {
	Address() common.Address
}

type AccountRef common.Address

// Address casts AccountRef to a Address
func (ar AccountRef) Address() common.Address { return (common.Address)(ar) }

type Contract struct {
	// CallerAddress is the result of the caller which initialised this
	// contract. However when the "call method" is delegated this value
	// needs to be initialised to that of the caller's caller.
	CallerAddress common.Address
	caller        ContractRef
	self          ContractRef

	jumpdests destinations // result of JUMPDEST analysis.

	Ast            types.Ast
	StateVariables []*types.Node
	Code           []byte
	CodeHash       common.Hash
	CodeAddr       *common.Address
	Input          []byte

	Gas   uint64
	value *big.Int

	Args []byte

	DelegateCall bool
}

func NewContract(caller ContractRef, object ContractRef, value *big.Int, gas uint64) *Contract {
	c := &Contract{CallerAddress: caller.Address(), caller: caller, self: object, Args: nil}

	if parent, ok := caller.(*Contract); ok {
		// Reuse JUMPDEST analysis from parent context if available.
		c.jumpdests = parent.jumpdests
	} else {
		c.jumpdests = make(destinations)
	}

	// Gas should be a pointer so it can safely be reduced through the run
	// This pointer will be off the state transition
	c.Gas = gas
	// ensures a value is set
	c.value = value

	return c
}

// AsDelegate sets the contract to be a delegate call and returns the current
// contract (for chaining calls)
func (c *Contract) AsDelegate() *Contract {
	c.DelegateCall = true
	// NOTE: caller must, at all times be a contract. It should never happen
	// that caller is something other than a Contract.
	parent := c.caller.(*Contract)
	c.CallerAddress = parent.CallerAddress
	c.value = parent.value

	return c
}

// GetOp returns the n'th element in the contract's byte array
func (c *Contract) GetOp(n uint64) OpCode {
	return OpCode(c.GetByte(n))
}

// GetByte returns the n'th byte in the contract's byte array
func (c *Contract) GetByte(n uint64) byte {
	if n < uint64(len(c.Code)) {
		return c.Code[n]
	}

	return 0
}

// Caller returns the caller of the contract.
//
// Caller will recursively call caller when the contract is a delegate
// call, including that of caller's caller.
func (c *Contract) Caller() common.Address {
	return c.CallerAddress
}

// UseGas attempts the use gas and subtracts it and returns true on success
func (c *Contract) UseGas(gas uint64) (ok bool) {
	if c.Gas < gas {
		return false
	}
	c.Gas -= gas
	return true
}

// Address returns the contracts address
func (c *Contract) Address() common.Address {
	return c.self.Address()
}

// Value returns the contracts value (sent to it from it's caller)
func (c *Contract) Value() *big.Int {
	return c.value
}

// ContractAst returns the contracts ast
func (c *Contract) ContractAst(pc uint) []byte {
	jast, err := json.Marshal(c.Ast[pc])
	if err != nil {
		return []byte{}
	}

	return jast
}

// ContractStateVariables returns the contracts state variables
func (c *Contract) ContractStateVariables() []byte {
	jast, err := json.Marshal(c.StateVariables)
	if err != nil {
		return []byte{}
	}

	return jast
}

// SetCode sets the code to the contract
func (c *Contract) SetCode(hash common.Hash, code []byte) {
	c.Code = code
	c.CodeHash = hash
}

// SetCallCode sets the code of the contract and address of the backing data
// object
func (c *Contract) SetCallCode(addr *common.Address, hash common.Hash, code []byte, ast types.Ast, stateVariables []*types.Node) {
	c.Code = code
	c.CodeHash = hash
	c.CodeAddr = addr
	c.Ast = ast
	c.StateVariables = stateVariables
}
