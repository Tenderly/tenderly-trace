package tenderly

import (
	"fmt"
	"github.com/tenderly/tenderly-trace/ethereum/client"
)

type Tenderly struct {
	client client.Client
}

//TODO: change contracts parameter to project folder root parameter and implement contract loading and framework detection
func NewTenderly(target string) (*Tenderly, error) {
	rpcClient, err := client.Dial(target)
	if err != nil {
		return nil, fmt.Errorf("failed dialing ethereum rpc: %s", err)
	}

	return &Tenderly{
		client: *rpcClient,
	}, nil
}
