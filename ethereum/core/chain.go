package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/types"
)

type Chain struct {
	Header *types.Header
}

func NewChain(header *types.Header) *Chain {
	return &Chain{
		Header: header,
	}
}

func (c *Chain) Engine() consensus.Engine {
	panic("implement me")
}

func (c *Chain) GetHeader(common.Hash, uint64) *types.Header {
	return c.Header
}
