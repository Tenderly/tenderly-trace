package core

import (
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tenderly/tenderly-trace/ethereum/signer/accounts/abi"
	"reflect"
	"strings"
)

type DecodedCallData struct {
	Signature string
	Name      string
	Inputs    []DecodedArgument
}

type DecodedArgument struct {
	Soltype abi.Argument
	Value   interface{}
}

type Argument struct {
	Name    string
	Type    Type
	Indexed bool // indexed is only used by events
}

type Type struct {
	Elem *Type

	Kind reflect.Kind
	Type reflect.Type
	Size int
	T    byte // Our own type checking

	stringKind string // holds the unparsed string for deriving signatures
}

type ABI struct {
	Constructor Method
	Methods     map[string]Method
	Events      map[string]Event
}

type Method struct {
	Name    string
	Const   bool
	Inputs  Arguments
	Outputs Arguments
}

type Arguments []Argument

type Event struct {
	Name      string
	Anonymous bool
	Inputs    Arguments
}

type ArgumentsUint struct {
	Position int `json:"position"`
	Value    int `json:"value"`
}

type ArgumentsString struct {
	Position int    `json:"position"`
	Value    string `json:"value"`
}

func ParseCallData(calldata []byte, abidata string) (*DecodedCallData, error) {

	if len(calldata) < 4 {
		return nil, fmt.Errorf("Invalid ABI-data, incomplete method signature of (%d bytes)", len(calldata))
	}

	sigdata, argdata := calldata[:4], calldata[4:]
	if len(argdata)%32 != 0 {
		return nil, fmt.Errorf("Not ABI-encoded data; length should be a multiple of 32 (was %d)", len(argdata))
	}

	abispec, err := abi.JSON(strings.NewReader(abidata))
	if err != nil {
		return nil, fmt.Errorf("Failed parsing JSON ABI: %v, abidata: %v", err, abidata)
	}

	method, err := abispec.MethodById(sigdata)
	if err != nil {
		return nil, err
	}

	v, err := method.Inputs.UnpackValues(argdata)
	if err != nil {
		return nil, err
	}

	decoded := DecodedCallData{Signature: method.Sig(), Name: method.Name}

	for n, transaction := range method.Inputs {
		if err != nil {
			return nil, fmt.Errorf("Failed to decode transaction %d (signature %v): %v", n, method.Sig(), err)
		}
		decodedArg := DecodedArgument{
			Soltype: transaction,
			Value:   v[n],
		}
		decoded.Inputs = append(decoded.Inputs, decodedArg)
	}

	// We're finished decoding the data. At this point, we encode the decoded data to see if it matches with the
	// original data. If we didn't do that, it would e.g. be possible to stuff extra data into the transactions, which
	// is not detected by merely decoding the data.

	var (
		encoded []byte
	)
	encoded, err = method.Inputs.PackValues(v)

	if err != nil {
		return nil, err
	}

	if !bytes.Equal(encoded, argdata) {
		was := common.Bytes2Hex(encoded)
		exp := common.Bytes2Hex(argdata)
		return nil, fmt.Errorf("WARNING: Supplied data is stuffed with extra data. \nWant %s\nHave %s\nfor method %v", exp, was, method.Sig())
	}
	return &decoded, nil
}
