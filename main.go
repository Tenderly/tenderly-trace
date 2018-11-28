package main

import (
	"github.com/tenderly/tenderly-trace/tenderly"
	"log"
)

func main() {
	tenderly, err := tenderly.NewTenderly("http://localhost:8545")
	if err != nil {
		log.Fatalf("Unable to connect to Ethereum RPC server")
	}

	tenderly.Trace("0x9b351a31c4fdbce74bbda251d195625c78a76f085b93d8a15d7bc9fa2421d829")
}
