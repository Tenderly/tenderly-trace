package tenderly

import (
	"github.com/tenderly/tenderly-trace/source"
	"time"
)

type Compiler struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func (c Compiler) GetName() string {
	return c.Name
}

func (c Compiler) GetVersion() string {
	return c.Version
}

type Network struct {
	Events          interface{} `json:"events"`
	Links           interface{} `json:"links"`
	Address         string      `json:"address"`
	TransactionHash string      `json:"transactionHash"`
}

func (n Network) GetEvents() interface{} {
	return n.Events
}

func (n Network) GetLinks() interface{} {
	return n.Links
}

func (n Network) GetAddress() string {
	return n.Address
}

func (n Network) GetTransactionHash() string {
	return n.TransactionHash
}

type Networks struct {
	Networks map[string]Network
}

func (n Networks) GetNetwork(address string) source.Network {
	return n.Networks[address]
}

type Ast struct {
	AbsolutePath    string           `json:"absolutePath"`
	ExportedSymbols map[string][]int `json:"exportedSymbols"`
	Id              int              `json:"id"`
	NodeType        string           `json:"nodeType"`
	Nodes           interface{}      `json:"nodes"`
	Src             string           `json:"src"`
}

func (a Ast) GetAbsolutePath() string {
	return a.AbsolutePath
}

func (a Ast) GetExportedSymbols() map[string][]int {
	return a.ExportedSymbols
}

func (a Ast) GetId() int {
	return a.Id
}

func (a Ast) GetNodeType() string {
	return a.NodeType
}

func (a Ast) GetNodes() interface{} {
	return a.Nodes
}

func (a Ast) GetSrc() string {
	return a.Src
}

type Contract struct {
	Name              string      `json:"contractName"`
	Abi               interface{} `json:"abi"`
	Bytecode          string      `json:"bytecode"`
	DeployedBytecode  string      `json:"deployedBytecode"`
	SourceMap         string      `json:"sourceMap"`
	DeployedSourceMap string      `json:"deployedSourceMap"`
	Source            string      `json:"source"`
	SourcePath        string      `json:"sourcePath"`
	Ast               Ast         `json:"legacyAST"`
	Compiler          Compiler    `json:"compiler"`
	Networks          Networks    `json:"networks"`

	SchemaVersion string    `json:"schemaVersion"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func (c Contract) GetName() string {
	return c.Name
}

func (c Contract) GetAbi() interface{} {
	return c.Abi
}

func (c Contract) GetBytecode() string {
	return c.Bytecode
}

func (c Contract) GetDeployedBytecode() string {
	return c.DeployedBytecode
}

func (c Contract) GetSourceMap() string {
	return c.SourcePath
}

func (c Contract) GetDeployedSourceMap() string {
	return c.DeployedSourceMap
}

func (c Contract) GetSource() string {
	return c.Source
}

func (c Contract) GetSourcePath() string {
	return c.SourcePath
}

func (c Contract) GetAst() source.Ast {
	return c.Ast
}

func (c Contract) GetCompiler() source.Compiler {
	return c.Compiler
}

func (c Contract) GetNetworks() source.Networks {
	return c.Networks
}

func (c Contract) GetSchemaVersion() string {
	return c.SchemaVersion
}

func (c Contract) GetUpdatedAt() time.Time {
	return c.UpdatedAt
}
