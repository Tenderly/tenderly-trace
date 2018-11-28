// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package state provides a caching layer atop the Ethereum state trie.
package state

import (
	"encoding/hex"
	"github.com/tenderly/tenderly-trace/ethereum"
	"github.com/tenderly/tenderly-trace/ethereum/client"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type revision struct {
	id           int
	journalIndex int
}

var (
	// emptyState is the known hash of an empty state trie entry.
	emptyState = crypto.Keccak256Hash(nil)

	// emptyCode is the known hash of the empty EVM bytecode.
	emptyCode = crypto.Keccak256Hash(nil)
)

type Cache struct {
	balance map[common.Address]*big.Int
	code    map[common.Address]*[]byte
	state   map[common.Address][]byte
}

func NewCache() *Cache {
	return &Cache{
		balance: make(map[common.Address]*big.Int),
		code:    make(map[common.Address]*[]byte),
		state:   make(map[common.Address][]byte),
	}
}

// StateDBs within the ethereum protocol are used to store anything
// within the merkle trie. StateDBs take care of caching and storing
// nested states. It's the general query interface to retrieve:
// * Contracts
// * Accounts
type StateDB struct {
	client      client.Client
	blockNumber int64
	cache       *Cache

	// This map holds 'live' objects, which will get modified while processing a state transition.
	stateObjectsDirty map[common.Address]struct{}

	// DB error.
	// State objects are used by the consensus core and VM which are
	// unable to deal with database-level errors. Any error that occurs
	// during a database read is memoized here and will eventually be returned
	// by StateDB.Commit.
	dbErr error

	// The refund counter, also used by state transitioning.
	refund uint64

	thash, bhash common.Hash
	txIndex      int
	logs         map[common.Hash][]*types.Log
	logSize      uint

	preimages map[common.Hash][]byte

	// Journal of state modifications. This is the backbone of
	// Snapshot and RevertToSnapshot.
	validRevisions []revision
	nextRevisionId int

	lock sync.Mutex
}

// Create a new state from a given trie.
func New(client client.Client, blockNumber int64) *StateDB {
	c := NewCache()
	return &StateDB{
		client:            client,
		blockNumber:       blockNumber,
		cache:             c,
		stateObjectsDirty: make(map[common.Address]struct{}),
		logs:              make(map[common.Hash][]*types.Log),
		preimages:         make(map[common.Hash][]byte),
	}
}

func (self *StateDB) AddLog(log *types.Log) {
	return
}

// AddPreimage records a SHA3 preimage seen by the VM.
func (self *StateDB) AddPreimage(hash common.Hash, preimage []byte) {
	return
}

func (self *StateDB) AddRefund(gas uint64) {
}

// Exist reports whether the given account address exists in the state.
// Notably this also returns true for suicided accounts.
func (self *StateDB) Exist(addr common.Address) bool {
	return true
}

// Empty returns whether the state object is either non-existent
// or empty according to the EIP161 specification (balance = nonce = code = 0)
func (self *StateDB) Empty(addr common.Address) bool {
	return false
}

// Retrieve the balance from the given address or 0 if object not found
func (self *StateDB) GetBalance(addr common.Address) *big.Int {
	balance, err := self.client.GetBalance(addr.String(), ethereum.Number(self.blockNumber))
	if err != nil || balance == nil {
		self.cache.balance[addr] = big.NewInt(0)
		return big.NewInt(0)
	}
	self.cache.balance[addr] = balance
	return balance
}

func (self *StateDB) GetNonce(addr common.Address) uint64 {
	return 0
}

func (self *StateDB) GetCode(addr common.Address) []byte {
	codeCache := self.cache.code[addr]
	if codeCache != nil {
		return *codeCache
	}
	code, err := self.client.GetCode(addr.String(), ethereum.Number(self.blockNumber))
	if err != nil {
		return []byte{}
	}
	raw := *code
	if strings.HasPrefix(raw, "0x") {
		raw = raw[2:]
	}
	bin, err := hex.DecodeString(raw)
	if err != nil {
		return []byte{}
	}
	return bin
}

func (self *StateDB) GetCodeSize(addr common.Address) int {
	codeCache := self.cache.code[addr]
	if codeCache != nil {
		return len(*codeCache)
	}
	code, err := self.client.GetCode(addr.String(), ethereum.Number(self.blockNumber))
	if err != nil {
		return 0
	}
	return len(*code)
}

func (self *StateDB) GetCodeHash(addr common.Address) common.Hash {
	codeCache := self.cache.code[addr]
	if codeCache != nil {
		return common.BytesToHash(crypto.Keccak256([]byte(*codeCache)))
	}
	code, err := self.client.GetCode(addr.String(), ethereum.Number(self.blockNumber))
	if err != nil {
		return common.Hash{}
	}
	return common.BytesToHash(crypto.Keccak256([]byte(*code)))
}

func (self *StateDB) GetState(addr common.Address, bhash common.Hash) common.Hash {
	data, _ := self.client.GetStorageAt(addr.String(), bhash, ethereum.Number(self.blockNumber))
	if data != nil {
		return *data
	}
	return common.Hash{}
}

func (self *StateDB) HasSuicided(addr common.Address) bool {
	return false
}

/*
 * SETTERS
 */

// AddBalance adds amount to the account associated with addr.
func (self *StateDB) AddBalance(addr common.Address, amount *big.Int) {
	balanceCache := self.cache.balance[addr]
	if balanceCache == nil {
		balanceCache = self.GetBalance(addr)
	}
	balanceCache.Add(balanceCache, amount)
}

// SubBalance subtracts amount from the account associated with addr.
func (self *StateDB) SubBalance(addr common.Address, amount *big.Int) {
	balanceCache := self.cache.balance[addr]
	if balanceCache == nil {
		balanceCache = self.GetBalance(addr)
	}
	balanceCache.Sub(balanceCache, amount)
}

func (self *StateDB) SetNonce(addr common.Address, nonce uint64) {
}

func (self *StateDB) SetCode(addr common.Address, code []byte) {
	self.cache.code[addr] = &code
}

func (self *StateDB) SetState(addr common.Address, key, value common.Hash) {
}

// Suicide marks the given account as suicided.
// This clears the account balance.
//
// The account's state object is still available until the state is committed,
// getStateObject will return a non-nil account after Suicide.
func (self *StateDB) Suicide(addr common.Address) bool {
	return false
}

// CreateAccount explicitly creates a state object. If a state object with the address
// already exists the balance is carried over to the new account.
//
// CreateAccount is called during the EVM CREATE operation. The situation might arise that
// a contract does the following:
//
//   1. sends funds to sha(account ++ (nonce + 1))
//   2. tx_create(sha(account ++ nonce)) (note that this gets the address of 1)
//
// Carrying over the balance ensures that Ether doesn't disappear.
func (self *StateDB) CreateAccount(addr common.Address) {
}

func (db *StateDB) ForEachStorage(addr common.Address, cb func(key, value common.Hash) bool) {
}

// Snapshot returns an identifier for the current revision of the state.
func (self *StateDB) Snapshot() int {
	return 0
}

// RevertToSnapshot reverts all state changes made since the given revision.
func (self *StateDB) RevertToSnapshot(revid int) {
}

// GetRefund returns the current value of the refund counter.
func (self *StateDB) GetRefund() uint64 {
	return 0
}
