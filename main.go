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

	tenderly.Trace("0x2c0c228f4d182185b9f11cc3e85a08baf60ef44fa19a948f06826584a1449f8e")
}
