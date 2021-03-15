package logic

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckProgram(t *testing.T) {
	program := []byte{1, 32, 1, 1, 34} // int 1
	var args [][]byte
	err := CheckProgram(program, args)
	require.NoError(t, err)

	args = make([][]byte, 1)
	args[0] = []byte(strings.Repeat("a", 10))
	err = CheckProgram(program, args)
	require.NoError(t, err)

	// too long arg
	args[0] = []byte(strings.Repeat("a", 1000))
	err = CheckProgram(program, args)
	require.Error(t, err)

	program = append(program, []byte(strings.Repeat("\x22", 10))...)
	args[0] = []byte(strings.Repeat("a", 10))
	err = CheckProgram(program, args)
	require.NoError(t, err)

	// too long program
	program = append(program, []byte(strings.Repeat("\x22", 1000))...)
	args[0] = []byte(strings.Repeat("a", 10))
	err = CheckProgram(program, args)
	require.EqualError(t, err, "program too long")

	// invalid opcode
	program = []byte{1, 32, 1, 1, 34} // int 1
	program[4] = 255
	args[0] = []byte(strings.Repeat("a", 10))
	err = CheckProgram(program, args)
	require.EqualError(t, err, "invalid instruction")

	// check single keccak256 and 10x keccak256 work
	program = []byte{0x01, 0x26, 0x01, 0x01, 0x01, 0x28, 0x02} // byte 0x01 + keccak256
	err = CheckProgram(program, args)
	require.NoError(t, err)

	program = append(program, []byte(strings.Repeat("\x02", 10))...) // append 10x keccak256
	err = CheckProgram(program, args)
	require.NoError(t, err)

	// check 800x keccak256 fail
	program = append(program, []byte(strings.Repeat("\x02", 800))...) // append 800x keccak256
	err = CheckProgram(program, args)
	require.EqualError(t, err, "program too costly to run")
}

func TestCheckProgramV2(t *testing.T) {
	// check TEAL v2 opcodes
	require.True(t, spec.EvalMaxVersion >= 2)
	require.True(t, spec.LogicSigVersion >= 2)

	args := make([][]byte, 0)

	// balance
	program := []byte{0x02, 0x20, 0x01, 0x00, 0x22, 0x60} // int 0; balance
	err := CheckProgram(program, args)
	require.NoError(t, err)

	// app_opted_in
	program = []byte{0x02, 0x20, 0x01, 0x00, 0x22, 0x22, 0x61} // int 0; int 0; app_opted_in
	err = CheckProgram(program, args)
	require.NoError(t, err)

	// asset_holding_get
	program = []byte{0x02, 0x20, 0x01, 0x00, 0x22, 0x70, 0x00} // int 0; int 0; asset_holding_get Balance
	err = CheckProgram(program, args)
	require.NoError(t, err)
}

func TestCheckProgramV3(t *testing.T) {
	// check TEAL v2 opcodes
	require.True(t, spec.EvalMaxVersion >= 3)
	require.True(t, spec.LogicSigVersion >= 3)

	args := make([][]byte, 0)

	// min_balance
	program := []byte{0x03, 0x20, 0x01, 0x00, 0x22, 0x78} // int 0; min_balance
	err := CheckProgram(program, args)
	require.NoError(t, err)

	// pushbytes
	program = []byte{0x03, 0x20, 0x01, 0x00, 0x22, 0x80, 0x02, 0x68, 0x69, 0x48} // int 0; pushbytes "hi"; pop
	err = CheckProgram(program, args)
	require.NoError(t, err)

	// pushint
	program = []byte{0x03, 0x20, 0x01, 0x00, 0x22, 0x81, 0x01, 0x48} // int 0; pushint 1; pop
	err = CheckProgram(program, args)
	require.NoError(t, err)

	// swap
	program = []byte{0x03, 0x20, 0x02, 0x00, 0x01, 0x22, 0x23, 0x4c, 0x48} // int 0; int 1; swap; pop
	err = CheckProgram(program, args)
	require.NoError(t, err)
}
