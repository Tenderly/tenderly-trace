package tenderly

import (
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	core2 "github.com/tenderly/tenderly-trace/ethereum/core"
	"github.com/tenderly/tenderly-trace/ethereum/core/state"
	"github.com/tenderly/tenderly-trace/ethereum/core/vm"
	"github.com/tenderly/tenderly-trace/ethereum/eth/tracers"
	"github.com/tenderly/tenderly-trace/source"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/tenderly/tenderly-trace/ethereum"

	"github.com/tenderly/tenderly-trace/ethereum/signer/accounts/signer/core"
)

type Stack struct {
	traces []*Trace
}

func NewStack(trace *Trace) *Stack {
	return &Stack{
		traces: []*Trace{trace},
	}
}

func (s *Stack) Push(trace *Trace) {
	s.traces = append(s.traces, trace)
}

func (s *Stack) Pop() {
	if len(s.traces) == 1 {
		return
	}

	s.traces = s.traces[:len(s.traces)-1]
}

func (s *Stack) Get() *Trace {
	return s.traces[len(s.traces)-1]
}

type Trace struct {
	CallType            ethereum.OpCode
	From                *common.Address
	To                  *common.Address
	Value               *hexutil.Big
	Gas                 *hexutil.Big
	GasUsed             *hexutil.Big
	GasPrice            *hexutil.Big
	State               hexutil.Bytes
	DecodedState        *[]core.DecodedArgument
	Input               hexutil.Bytes
	DecodedInput        *[]core.DecodedArgument
	Output              hexutil.Bytes
	DecodedOutput       *[]core.DecodedArgument
	Locals              []hexutil.Bytes
	DecodedLocals       *[]core.DecodedArgument
	ParentLocals        []hexutil.Bytes
	DecodedParentLocals *[]core.DecodedArgument
	Error               ethereum.OpCode
	ErrorLine           *int64
	Trace               []Trace
}

func (t Tenderly) Trace(txHash string, cs source.Source) (*Trace, error) {
	tx, err := t.client.GetTransaction(txHash)
	if err != nil {
		return nil, fmt.Errorf("failed fetching transaction %s, err: %s\n", txHash, err)
	}

	if tx.BlockNumber() == nil {
		return nil, fmt.Errorf("transaction is in pending status")
	}

	blockHeader, err := t.client.GetBlockByHash(tx.BlockHash().String())
	if err != nil {
		return nil, fmt.Errorf("failed fetcing block %s, err: %s\n", tx.BlockNumber().String(), err)
	}

	tracer, err := tracers.New("callTracerFinal")

	context := buildContext(tx, blockHeader)
	stateDB := state.New(t.client, blockHeader.Number().Value(), cs.GetSource())
	chainConfig := params.TestChainConfig
	vmConfig := vm.Config{Debug: true, Tracer: tracer}

	env := vm.NewEVM(context, stateDB, chainConfig, vmConfig)
	env.Call(vm.AccountRef(*tx.From()), *tx.To(), tx.Input(), tx.Gas().ToInt().Uint64(), tx.Value().ToInt())

	results, err := tracer.GetResult()
	fmt.Println(string(results))

	return nil, nil
}

func buildContext(tx ethereum.Transaction, blockHeader ethereum.BlockHeader) vm.Context {
	message := types.NewMessage(*tx.From(), tx.To(), 0, tx.Value().ToInt(), tx.Gas().ToInt().Uint64(),
		tx.GasPrice().ToInt(), tx.Input(), false)
	header := types.Header{
		Number:     big.NewInt(blockHeader.Number().Value()),
		ParentHash: *blockHeader.ParentHash(),
		Time:       blockHeader.Time().ToInt(),
		Difficulty: blockHeader.Difficulty().ToInt(),
		GasLimit:   blockHeader.GasLimit().ToInt().Uint64(),
	}
	chain := core2.NewChain(&header)

	return core2.NewEVMContext(message, &header, chain, blockHeader.Coinbase())
}
