package source

import (
	"github.com/tenderly/tenderly-trace/ethereum/core/types"
)

type InstructionMapping struct {
	Src   string
	Index int

	Start  int
	Length int

	Line   int
	Column int

	FileIndex int
	Jump      string
}

type Ast map[int]*types.Node
type SourceMap map[int]*InstructionMapping

type Contract interface {
	GetContractAst(sourceMap SourceMap) types.Ast
	GetContractStateVariables() []*types.Node
	GetContractSourceMap() (SourceMap, error)
}

type ContractSource struct {
	Contracts map[string]Contract //mapping code => contract interface
}

func (cs ContractSource) GetAst(code string) types.Ast {
	contractCode, ok := cs.Contracts[code]
	if !ok {
		return nil
	}

	sourceMap, err := contractCode.GetContractSourceMap()
	if err != nil {
		return nil
	}

	return contractCode.GetContractAst(sourceMap)
}

func (cs ContractSource) GetStateVariables(code string) []*types.Node {
	contractCode, ok := cs.Contracts[code]
	if !ok {
		return nil
	}

	return contractCode.GetContractStateVariables()
}

//func (cs ContractSource) GetAst(code string) *state.Variables {
//	contract := cs.Contracts[code]
//	if contract == nil {
//		return nil
//	}
//
//	ast := contract.GetContractAst()
//
//	return parseAst(&ast)
//}
//
//func parseAst(ast *Ast) *state.Variables {
//	if ast == nil {
//		return nil
//	}
//
//	variables := state.NewVariables()
//
//	for _, node := range ast.Nodes {
//		explore(node, variables)
//	}
//
//	return variables
//}
//
//func explore(node Node, variables *state.Variables) {
//	if node.Nodes != nil {
//		for _, childNode := range node.Nodes {
//			explore(childNode, variables)
//		}
//	}
//
//	if node.StateVariable {
//		variables.IsGlobal[node.Name] = true
//	}
//}

type Source interface {
	GetSource() ContractSource
}

//func NewContractSource(contracts map[string]Contract, client client.Client) *ContractSource {
//	return &ContractSource{
//		contracts: contracts,
//		client:    client,
//	}
//}
//
//func (cs ContractSource) GetContract(address string) Contract {
//	if contract, ok := cs.contracts[address]; ok {
//		return contract
//	}
//
//	deployedBytecode, err := cs.client.GetCode(address)
//	if err != nil {
//		// TODO: Log
//		return tenderly.Contract{
//			DeployedBytecode: *deployedBytecode,
//		}
//	}
//
//	for _, c := range cs.contracts {
//		if c.GetDeployedBytecode() == *deployedBytecode {
//			return c
//		}
//	}
//
//	// TODO: Log
//	return tenderly.Contract{
//		DeployedBytecode: *deployedBytecode,
//	}
//}

//func (cs ContractSource) DecodeInput(address string, rawInput hexutil.Bytes) *[]core.DecodedArgument {
//	c := cs.GetContract(address)
//	if c.GetAbi() != nil {
//		if decodedInput, err := core.ParseCallData([]byte(rawInput), c.GetAbi().(string)); err != nil {
//			return nil
//		} else {
//			return &decodedInput.Inputs
//		}
//	}
//
//	return nil
//}
