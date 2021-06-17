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
	_, _, err := ReadProgram(program, args)
	return err
}

// ReadProgram is used to validate a program as well as extract found variables
func ReadProgram(program []byte, args [][]byte) (ints []uint64, byteArrays [][]byte, err error) {
	const intcblockOpcode = 32
	const bytecblockOpcode = 38
	const pushbytesOpcode = 128
	const pushintOpcode = 129
	if program == nil || len(program) == 0 {
		err = fmt.Errorf("empty program")
		return
	}

	if spec == nil {
		spec = new(langSpec)
		if err = json.Unmarshal(langSpecJson, spec); err != nil {
			return
		}
	}
	version, vlen := binary.Uvarint(program)
	if vlen <= 0 {
		err = fmt.Errorf("version parsing error")
		return
	}
	if int(version) > spec.EvalMaxVersion {
		err = fmt.Errorf("unsupported version")
		return
	}

	cost := 0
	length := len(program)
	for _, arg := range args {
		length += len(arg)
	}

	if length > types.LogicSigMaxSize {
		err = fmt.Errorf("program too long")
		return
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
			err = fmt.Errorf("invalid instruction")
			return
		}

		cost = cost + op.Cost
		size := op.Size
		if size == 0 {
			switch op.Opcode {
			case intcblockOpcode:
				var foundInts []uint64
				size, foundInts, err = readIntConstBlock(program, pc)
				ints = append(ints, foundInts...)
				if err != nil {
					return
				}
			case bytecblockOpcode:
				var foundByteArrays [][]byte
				size, foundByteArrays, err = readByteConstBlock(program, pc)
				byteArrays = append(byteArrays, foundByteArrays...)
				if err != nil {
					return
				}
			case pushintOpcode:
				var foundInt uint64
				size, foundInt, err = readPushIntOp(program, pc)
				ints = append(ints, foundInt)
				if err != nil {
					return
				}
			case pushbytesOpcode:
				var foundByteArray []byte
				size, foundByteArray, err = readPushByteOp(program, pc)
				byteArrays = append(byteArrays, foundByteArray)
				if err != nil {
					return
				}
			default:
				err = fmt.Errorf("invalid instruction")
				return
			}
		}
		pc = pc + size
	}

	// costs calculated dynamically starting in v4
	if version < 4 && cost > types.LogicSigMaxCost {
		err = fmt.Errorf("program too costly for Teal version < 4. consider using v4.")
	}

	return
}

func readIntConstBlock(program []byte, pc int) (size int, ints []uint64, err error) {
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
		num, bytesUsed := binary.Uvarint(program[pc+size:])
		if bytesUsed <= 0 {
			err = fmt.Errorf("could not decode int const[%d] at pc=%d", i, pc+size)
			return
		}
		ints = append(ints, num)
		size += bytesUsed
	}
	return
}

func readByteConstBlock(program []byte, pc int) (size int, byteArrays [][]byte, err error) {
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
		scanTarget := program[pc+size:]
		itemLen, bytesUsed := binary.Uvarint(scanTarget)
		if bytesUsed <= 0 {
			err = fmt.Errorf("could not decode []byte const[%d] at pc=%d", i, pc+size)
			return
		}
		size += bytesUsed
		if pc+size+int(itemLen) > len(program) {
			err = fmt.Errorf("bytecblock ran past end of program")
			return
		}
		byteArray := program[pc+size : pc+size+int(itemLen)]
		byteArrays = append(byteArrays, byteArray)
		size += int(itemLen)
	}
	return
}

func readPushIntOp(program []byte, pc int) (size int, foundInt uint64, err error) {
	size = 1
	foundInt, bytesUsed := binary.Uvarint(program[pc+size:])
	if bytesUsed <= 0 {
		err = fmt.Errorf("could not decode push int const at pc=%d", pc+size)
		return
	}

	size += bytesUsed
	return
}

func readPushByteOp(program []byte, pc int) (size int, byteArray []byte, err error) {
	size = 1
	itemLen, bytesUsed := binary.Uvarint(program[pc+size:])
	if bytesUsed <= 0 {
		err = fmt.Errorf("could not decode push []byte const size at pc=%d", pc+size)
		return
	}

	size += bytesUsed
	if pc+size+int(itemLen) > len(program) {
		err = fmt.Errorf("pushbytes ran past end of program")
		return
	}
	byteArray = program[pc+size : pc+size+int(itemLen)]
	size += int(itemLen)
	return
}
