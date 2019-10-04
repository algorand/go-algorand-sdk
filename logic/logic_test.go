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
	program[4] = 128
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
