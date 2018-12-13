package types

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
	Expression *Expression `json:"expression"`
	Id         int         `json:"id"`
	NodeType   string      `json:"node_type"`
	Src        string      `json:"src"`
}

type Body struct {
	Id        int       `json:"id"`
	NodeType  string    `json:"node_type"`
	Src       string    `json:"src"`
	Statement Statement `json:"statement"`
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
	Name             string     `json:"name"`
	NodeType         string     `json:"nodeType"`
	TypeName         TypeName   `json:"typeName"`
	StateVariable    bool       `json:"stateVariable"`
	Parameters       Parameters `json:"parameters"`
	ReturnParameters Parameters `json:"returnParameters"`
}

type Ast map[uint]*Node
