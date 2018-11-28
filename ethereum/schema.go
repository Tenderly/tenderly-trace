package ethereum

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/tenderly/tenderly-trace/jsonrpc2"
)

type Schema interface {
	Eth() EthSchema
	Net() NetSchema
	Trace() TraceSchema
	PubSub() PubSubSchema
}

// Eth

type EthSchema interface {
	BlockNumber() (*jsonrpc2.Request, *Number)
	GetBlockByNumber(num Number) (*jsonrpc2.Request, Block)
	GetBlockByHash(hash string) (*jsonrpc2.Request, BlockHeader)
	GetTransaction(hash string) (*jsonrpc2.Request, Transaction)
	GetTransactionReceipt(hash string) (*jsonrpc2.Request, TransactionReceipt)
	GetBalance(address string, block Number) (*jsonrpc2.Request, *hexutil.Big)
	GetCode(address string, block Number) (*jsonrpc2.Request, *string)
	GetStorage(address string, offset common.Hash, block Number) (*jsonrpc2.Request, *common.Hash)
}

// Net

type NetSchema interface {
	Version() (*jsonrpc2.Request, *string)
}

// States

type TraceSchema interface {
	VMTrace(hash string) (*jsonrpc2.Request, TransactionStates)
	CallTrace(hash string) (*jsonrpc2.Request, CallTraces)
}

// PubSub

type PubSubSchema interface {
	Subscribe() (*jsonrpc2.Request, *SubscriptionID)
	Unsubscribe(id SubscriptionID) (*jsonrpc2.Request, *UnsubscribeSuccess)
}

type pubSubSchema struct {
}

func (pubSubSchema) Subscribe() (*jsonrpc2.Request, *SubscriptionID) {
	id := NewNilSubscriptionID()

	return jsonrpc2.NewRequest("eth_subscribe", "newHeads"), &id
}

func (pubSubSchema) Unsubscribe(id SubscriptionID) (*jsonrpc2.Request, *UnsubscribeSuccess) {
	var success UnsubscribeSuccess

	return jsonrpc2.NewRequest("eth_unsubscribe", id.String()), &success
}
