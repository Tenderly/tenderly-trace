package truffle

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/tenderly/tenderly-trace/source"
)

type ContractSource struct {
	contracts map[string]*Contract
}

func (cs ContractSource) GetSource() source.ContractSource {
	cast := make(map[string]source.Contract)
	for k, v := range cs.contracts {
		cast[k] = v
	}

	return source.ContractSource{Contracts: cast}
}

// NewContractSource builds the Contract Source from the provided config, and scoped to the provided network.
func NewContractSource(absBuildDir string) (*ContractSource, error) {
	truffleContracts, err := loadTruffleContracts(absBuildDir)
	if err != nil {
		return nil, err
	}

	return mapTruffleContracts(truffleContracts), nil
}

func loadTruffleContracts(absBuildDir string) ([]*Contract, error) {
	files, err := ioutil.ReadDir(absBuildDir)
	if err != nil {
		return nil, fmt.Errorf("failed listing truffle build files: %s", err)
	}

	var contracts []*Contract
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		data, err := ioutil.ReadFile(filepath.Join(absBuildDir, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed reading truffle build files: %s", err)
		}

		var contract Contract
		err = json.Unmarshal(data, &contract)
		if err != nil {
			return nil, fmt.Errorf("failed parsing truffle build files: %s", err)
		}

		contracts = append(contracts, &contract)
	}

	return contracts, nil
}

func mapTruffleContracts(truffleContracts []*Contract) *ContractSource {
	contracts := make(map[string]*Contract)

	for _, truffleContract := range truffleContracts {
		contracts[truffleContract.DeployedBytecode] = truffleContract
	}

	return &ContractSource{contracts: contracts}
}
