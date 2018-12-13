package truffle

import (
	"encoding/hex"
	"fmt"
	"github.com/tenderly/tenderly-trace/ethereum/core/vm"
	"github.com/tenderly/tenderly-trace/source"
	"strconv"
	"strings"
)

func ParseSourcecode(contract *Contract) (source.SourceMap, error) {
	rawSrcMap := contract.DeployedSourceMap
	instructionSrcMap, err := parseInstructionSourceMap(rawSrcMap)
	if err != nil {
		return nil, fmt.Errorf("sourcemap.Parse: %s", err)
	}

	rawBytecode := contract.DeployedBytecode

	memSrcMap, err := convertToMemoryMap(instructionSrcMap, rawBytecode)
	if err != nil {
		return nil, fmt.Errorf("sourcemap.Parse: %s", err)
	}

	rawSrc := contract.Source

	for _, instruction := range memSrcMap {
		if instruction == nil {
			continue
		}

		i := 0
		line := 1
		column := 1

		for i < instruction.Start && i < len(rawSrc) {
			if rawSrc[i] == '\n' {
				line++
				column = 0
			}

			column++
			i++
		}

		instruction.Line = line
		instruction.Column = column
	}

	return memSrcMap, nil
}

func parseInstructionSourceMap(rawSrcMap string) (source.SourceMap, error) {
	instructionSrcMap := make(source.SourceMap)

	var err error
	var s int64
	var l int64
	var f int64
	var j string
	for index, mapping := range strings.Split(rawSrcMap, ";") {
		if mapping == "" {
			instructionSrcMap[index] = instructionSrcMap[index-1]
			continue
		}

		info := strings.Split(mapping, ":")
		if len(info) > 0 && info[0] != "" {
			s, err = strconv.ParseInt(info[0], 0, 64)
			if err != nil {
				return nil, fmt.Errorf("failed parsing integer: %s", err)
			}
		}
		if len(info) > 1 && info[1] != "" {
			l, err = strconv.ParseInt(info[1], 0, 64)
			if err != nil {
				return nil, fmt.Errorf("failed parsing integer: %s", err)
			}
		} else {
			l = int64(instructionSrcMap[index-1].Length)
		}

		if len(info) > 2 && info[2] != "" {
			f, err = strconv.ParseInt(info[2], 0, 64)
			if err != nil {
				return nil, fmt.Errorf("failed parsing integer: %s", err)
			}
		} else {
			f = int64(instructionSrcMap[index-1].FileIndex)
		}

		if len(info) > 3 && info[3] != "" {
			j = info[3]
		}

		instructionSrcMap[index] = &source.InstructionMapping{
			Src:       strconv.Itoa(int(s)) + ":" + strconv.Itoa(int(l)) + ":" + strconv.Itoa(int(f)),
			Index:     index,
			Start:     int(s),
			Length:    int(l),
			FileIndex: int(f),
			Jump:      j,
		}
	}

	return instructionSrcMap, nil
}

func convertToMemoryMap(sourceMap source.SourceMap, binData string) (source.SourceMap, error) {
	memSrcMap := make(source.SourceMap)

	bin, err := hex.DecodeString(binData[2:])
	if err != nil {
		return nil, fmt.Errorf("failed decoding runtime binary: %s", err)
	}

	instruction := 0
	for i := 0; i < len(bin); i++ {

		op := vm.OpCode(bin[i])
		extraPush := 0
		if op.IsPush() {
			// Skip more here
			extraPush = int(op - vm.PUSH1 + 1)
		}

		memSrcMap[i] = sourceMap[instruction]

		instruction++
		i += extraPush
	}

	return memSrcMap, nil
}
