package source

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/tenderly/tenderly-trace/ethereum/client"
	"github.com/tenderly/tenderly-trace/ethereum/core/accounts/signer/core"
	"github.com/tenderly/tenderly-trace/source/tenderly"
	"time"
)

type Ast interface {
	GetAbsolutePath() string
	GetExportedSymbols() map[string][]int
	GetId() int
	GetNodeType() string
	GetNodes() interface{}
	GetSrc() string
}

type Compiler interface {
	GetName() string
	GetVersion() string
}

type Network interface {
	GetEvents() interface{}
	GetLinks() interface{}
	GetAddress() string
	GetTransactionHash() string
}

type Networks interface {
	GetNetwork(string) Network
}

type Contract interface {
	GetName() string
	GetAbi() interface{}
	GetBytecode() string
	GetDeployedBytecode() string
	GetSourceMap() string
	GetDeployedSourceMap() string
	GetSource() string
	GetSourcePath() string
	GetAst() Ast
	GetCompiler() Compiler
	GetNetworks() Networks

	GetSchemaVersion() string
	GetUpdatedAt() time.Time
}

type ContractSource struct {
	contracts map[string]Contract //mapping address => contract interface
	client    client.Client
}

func NewContractSource(contracts map[string]Contract, client client.Client) *ContractSource {
	return &ContractSource{
		contracts: contracts,
		client:    client,
	}
}

func (cs ContractSource) GetContract(address string) Contract {
	if contract, ok := cs.contracts[address]; ok {
		return contract
	}

	deployedBytecode, err := cs.client.GetCode(address)
	if err != nil {
		// TODO: Log
		return tenderly.Contract{
			DeployedBytecode: *deployedBytecode,
		}
	}

	for _, c := range cs.contracts {
		if c.GetDeployedBytecode() == *deployedBytecode {
			return c
		}
	}

	// TODO: Log
	return tenderly.Contract{
		DeployedBytecode: *deployedBytecode,
	}
}

func (cs ContractSource) DecodeInput(address string, rawInput hexutil.Bytes) *[]core.DecodedArgument {
	c := cs.GetContract(address)
	if c.GetAbi() != nil {
		if decodedInput, err := core.ParseCallData([]byte(rawInput), c.GetAbi().(string)); err != nil {
			return nil
		} else {
			return &decodedInput.Inputs
		}
	}

	return nil
}
