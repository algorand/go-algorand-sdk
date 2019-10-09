package logic

//go:generate ./bundle_langspec_json.sh

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/algorand/go-algorand-sdk/types"
)

type langSpec struct {
	EvalMaxVersion  int
	LogicSigVersion int
	Ops             []operation
}

type operation struct {
	Opcode        int
	Name          string
	Cost          int
	Size          int
	Returns       string
	ArgEnum       []string
	ArgEnumTypes  string
	Doc           string
	ImmediateNote string
	Group         []string
}

var spec *langSpec
var opcodes []operation

// CheckProgram performs basic program validation: instruction count and program cost
func CheckProgram(program []byte, args [][]byte) error {
	const intcblockOpcode = 32
	const bytecblockOpcode = 38
	if program == nil || len(program) == 0 {
		return fmt.Errorf("empty program")
	}

	if spec == nil {
		spec = new(langSpec)
		if err := json.Unmarshal(langSpecJson, spec); err != nil {
			return err
		}
	}
	version, vlen := binary.Uvarint(program)
	if vlen <= 0 {
		return fmt.Errorf("version parsing error")
	}
	if int(version) > spec.EvalMaxVersion {
		return fmt.Errorf("unsupported version")
	}

	cost := 0
	length := len(program)
	for _, arg := range args {
		length += len(arg)
	}
	if length > types.LogicSigMaxSize {
		return fmt.Errorf("program too long")
	}

	if opcodes == nil {
		opcodes = make([]operation, 256)
		for _, op := range spec.Ops {
			opcodes[op.Opcode] = op
		}
	}

	for pc := vlen; pc < len(program); {
		op := opcodes[program[pc]]
		if op.Name == "" {
			return fmt.Errorf("invalid instruction")
		}

		cost = cost + op.Cost
		size := op.Size
		var err error
		if size == 0 {
			switch op.Opcode {
			case intcblockOpcode:
				size, err = checkIntConstBlock(program, pc)
				if err != nil {
					return err
				}
			case bytecblockOpcode:
				size, err = checkByteConstBlock(program, pc)
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("invalid instruction")
			}
		}
		pc = pc + size
	}

	if cost > types.LogicSigMaxCost {
		return fmt.Errorf("program too costly to run")
	}

	return nil
}

func checkIntConstBlock(program []byte, pc int) (size int, err error) {
	size = 1
	numInts, bytesUsed := binary.Uvarint(program[pc+size:])
	if bytesUsed <= 0 {
		err = fmt.Errorf("could not decode int const block size at pc=%d", pc+size)
		return
	}

	size += bytesUsed
	for i := uint64(0); i < numInts; i++ {
		if pc+size >= len(program) {
			err = fmt.Errorf("intcblock ran past end of program")
			return
		}
		_, bytesUsed = binary.Uvarint(program[pc+size:])
		if bytesUsed <= 0 {
			err = fmt.Errorf("could not decode int const[%d] at pc=%d", i, pc+size)
			return
		}
		size += bytesUsed
	}
	return
}

func checkByteConstBlock(program []byte, pc int) (size int, err error) {
	size = 1
	numInts, bytesUsed := binary.Uvarint(program[pc+size:])
	if bytesUsed <= 0 {
		err = fmt.Errorf("could not decode []byte const block size at pc=%d", pc+size)
		return
	}

	size += bytesUsed
	for i := uint64(0); i < numInts; i++ {
		if pc+size >= len(program) {
			err = fmt.Errorf("bytecblock ran past end of program")
			return
		}
		itemLen, bytesUsed := binary.Uvarint(program[pc+size:])
		if bytesUsed <= 0 {
			err = fmt.Errorf("could not decode []byte const[%d] at pc=%d", i, pc+size)
			return
		}
		size += bytesUsed
		if pc+size >= len(program) {
			err = fmt.Errorf("bytecblock ran past end of program")
			return
		}
		size += int(itemLen)
	}
	return
}
