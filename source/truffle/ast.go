package truffle

import (
	"github.com/tenderly/tenderly-trace/ethereum/core/types"
	"github.com/tenderly/tenderly-trace/source"
)

type TempAst map[string]*types.Node

func ParseAst(contract *Contract, sourceMap source.SourceMap) types.Ast {
	tast := make(TempAst)
	recursiveParseAst(tast, contract.Ast.Nodes)

	ast := make(types.Ast)

	for k, v := range sourceMap {
		if v != nil {
			ast[uint(k)] = tast[v.Src]
		}
	}

	return ast
}

func ParseStateVariables(contract *Contract) []*types.Node {
	tast := make(TempAst)
	recursiveParseAst(tast, contract.Ast.Nodes)

	var sv []*types.Node

	for _, node := range tast {
		if node.StateVariable {
			sv = append(sv, node)
		}
	}

	return sv
}

func recursiveParseAst(tast TempAst, nodes []Node) TempAst {
	for _, node := range nodes {
		tast[node.Src] = convert(node)

		if node.Nodes != nil {
			recursiveParseAst(tast, node.Nodes)
		}

		if node.Parameters.Parameters != nil {
			recursiveParseAst(tast, node.Parameters.Parameters)
		}

		if node.ReturnParameters.Parameters != nil {
			recursiveParseAst(tast, node.ReturnParameters.Parameters)
		}

		if node.Body.Statements != nil {

			for _, statement := range node.Body.Statements {
				recursiveParseAst(tast, statement.Declarations)
			}
		}
	}

	return tast
}

func convert(node Node) *types.Node {
	return &types.Node{
		Name:          node.Name,
		NodeType:      node.NodeType,
		StateVariable: node.StateVariable,
		TypeName: types.TypeName{
			ContractScope:         node.TypeName.ContractScope,
			Id:                    node.TypeName.Id,
			Name:                  node.TypeName.Name,
			NodeType:              node.TypeName.NodeType,
			ReferencedDeclaration: node.TypeName.ReferencedDeclaration,
			Src: node.TypeName.Src,
			TypeDescription: types.TypeDescriptions{
				TypeIdentifier: node.TypeDescriptions.TypeIdentifier,
				TypeString:     node.TypeDescriptions.TypeString,
			},
		},
		Parameters: types.Parameters{
			Id:         node.Parameters.Id,
			NodeType:   node.Parameters.NodeType,
			Parameters: convertArray(node.Parameters.Parameters),
			Src:        node.Parameters.Src,
		},
		ReturnParameters: types.Parameters{
			Id:         node.Parameters.Id,
			NodeType:   node.Parameters.NodeType,
			Parameters: convertArray(node.ReturnParameters.Parameters),
			Src:        node.Parameters.Src,
		},
	}
}

func convertArray(nodes []Node) []types.Node {
	var snodes []types.Node
	for _, node := range nodes {
		snodes = append(snodes, *convert(node))
	}

	return snodes
}
