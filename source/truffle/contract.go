package truffle

import (
	"fmt"
	"github.com/tenderly/tenderly-trace/ethereum/core/types"
	"github.com/tenderly/tenderly-trace/source"
	"time"
)

type ContractCompiler struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ContractNetwork struct {
	Events          interface{} `json:"events"`
	Links           interface{} `json:"links"`
	Address         string      `json:"address"`
	TransactionHash string      `json:"transactionHash"`
}

type ContractAst struct {
	AbsolutePath    string           `json:"absolutePath"`
	ExportedSymbols map[string][]int `json:"exportedSymbols"`
	Id              int              `json:"id"`
	NodeType        string           `json:"nodeType"`
	Nodes           []Node           `json:"nodes"`
	Src             string           `json:"src"`
}

type Contract struct {
	Name                    string      `json:"contractName"`
	Abi                     interface{} `json:"abi"`
	Bytecode                string      `json:"bytecode"`
	DeployedBytecode        string      `json:"deployedBytecode"`
	SourceMap               string      `json:"sourceMap"`
	DeployedSourceMap       string      `json:"deployedSourceMap"`
	ParsedDeployedSourceMap source.SourceMap
	Source                  string      `json:"source"`
	SourcePath              string      `json:"sourcePath"`
	Ast                     ContractAst `json:"ast"`
	ParsedStateVariable     []*types.Node
	ParsedAst               types.Ast
	Compiler                ContractCompiler           `json:"compiler"`
	Networks                map[string]ContractNetwork `json:"networks"`

	SchemaVersion string    `json:"schemaVersion"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type Expression struct {
	Id               int
	IsConstant       bool
	IsLValue         bool
	IsPure           bool
	LValueRequested  bool
	leftHandSide     *Expression
	NodeType         string
	Operator         string
	RightHandSide    *Expression
	Src              string
	TypeDescriptions TypeDescriptions
}

type Statement struct {
	Assignments  []int       `json:"assignments"`
	Declarations []Node      `json:"declarations"`
	Expression   *Expression `json:"expression"`
	Id           int         `json:"id"`
	InitialValue Expression  `json:"initialValue"`
	NodeType     string      `json:"nodeType"`
	Src          string      `json:"src"`
}

type Body struct {
	Id         int         `json:"id"`
	NodeType   string      `json:"nodeType"`
	Src        string      `json:"src"`
	Statements []Statement `json:"statements"`
}
type TypeDescriptions struct {
	TypeIdentifier string `json:"typeIdentifier"`
	TypeString     string `json:"typeString"`
}

type TypeName struct {
	ContractScope         string           `json:"contractScope"`
	Id                    int              `json:"id"`
	Name                  string           `json:"name"`
	NodeType              string           `json:"nodeType"`
	ReferencedDeclaration int              `json:"referencedDeclaration"`
	Src                   string           `json:"src"`
	TypeDescription       TypeDescriptions `json:"typeDescription"`
}

type Parameters struct {
	Id         int    `json:"id"`
	NodeType   string `json:"nodeType"`
	Parameters []Node `json:"parameters"`
	Src        string `json:"src"`
}

type Node struct {
	Body                 Body             `json:"body"`
	BaseContracts        []string         `json:"baseContracts"`
	ContractDependencies []string         `json:"contractDependencies"`
	ContractKind         string           `json:"contractKind"`
	Documentation        interface{}      `json:"documentation"`
	FullyImplemented     bool             `json:"fullyImplemented"`
	Id                   int              `json:"id"`
	Literals             []string         `json:"literals"`
	Implemented          bool             `json:"implemented"`
	IsConstructor        bool             `json:"isConstructor"`
	IsDeclaredConst      bool             `json:"isDeclaredConst"`
	Modifiers            []interface{}    `json:"modifiers"`
	Constant             bool             `json:"constant"`
	Name                 string           `json:"name"`
	Nodes                []Node           `json:"nodes"`
	NodeType             string           `json:"nodeType"`
	Parameters           Parameters       `json:"parameters"`
	Payable              bool             `json:"payable"`
	ReturnParameters     Parameters       `json:"returnParameters"`
	Scope                int              `json:"scope"`
	Src                  string           `json:"src"`
	StateVariable        bool             `json:"stateVariable"`
	StorageLocation      string           `json:"storageLocation"`
	TypeDescriptions     TypeDescriptions `json:"typeDescriptions"`
	TypeName             TypeName         `json:"typeName"`
	Value                *string          `json:"value"`
	StateMutability      string           `json:"stateMutability"`
	SuperFunction        *string          `json:"superFunction"`
	Visibility           string           `json:"visibility"`
}

func (c *Contract) GetContractAst(sourceMap source.SourceMap) types.Ast {
	if c.ParsedAst != nil {
		return c.ParsedAst
	}

	return ParseAst(c, sourceMap)
}

func (c *Contract) GetContractStateVariables() []*types.Node {
	if c.ParsedStateVariable != nil {
		return c.ParsedStateVariable
	}

	return ParseStateVariables(c)
}

func (c *Contract) GetContractSourceMap() (source.SourceMap, error) {
	if c.ParsedDeployedSourceMap != nil {
		return c.ParsedDeployedSourceMap, nil
	}

	sourceMap, err := ParseSourcecode(c)
	if err != nil {
		return nil, fmt.Errorf("unable to parse source map, err %s\n", err)
	}

	return sourceMap, nil
}
