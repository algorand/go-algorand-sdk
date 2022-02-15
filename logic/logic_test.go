package logic

import (
	"encoding/hex"
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

	// check 800x keccak256 fail for v3 and below
	versions := []byte{0x1, 0x2, 0x3}
	program = append(program, []byte(strings.Repeat("\x02", 800))...) // append 800x keccak256
	for _, v := range versions {
		program[0] = v
		err = CheckProgram(program, args)
		require.EqualError(t, err, "program too costly for Teal version < 4. consider using v4.")
	}

	// check 800x keccak256 ok for v4 and above
	versions = []byte{0x4}
	for _, v := range versions {
		program[0] = v
		err = CheckProgram(program, args)
		require.NoError(t, err)
	}
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
	// check TEAL v3 opcodes
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

func TestCheckProgramV4(t *testing.T) {
	// check TEAL v4 opcodes
	require.True(t, spec.EvalMaxVersion >= 4)

	args := make([][]byte, 0)

	// divmodw
	program := []byte{0x04, 0x20, 0x03, 0x01, 0x00, 0x02, 0x22, 0x81, 0xd0, 0x0f, 0x23, 0x24, 0x1f} // int 1; pushint 2000; int 0; int 2; divmodw
	err := CheckProgram(program, args)
	require.NoError(t, err)

	// gloads i
	program = []byte{0x04, 0x20, 0x01, 0x00, 0x22, 0x3b, 0x00} // int 0; gloads 0
	err = CheckProgram(program, args)
	require.NoError(t, err)

	// callsub
	program = []byte{0x04, 0x20, 0x02, 0x01, 0x02, 0x22, 0x88, 0x00, 0x02, 0x23, 0x12, 0x49} // int 1; callsub double; int 2; ==; double: dup;
	err = CheckProgram(program, args)
	require.NoError(t, err)

	// b>=
	program = []byte{0x04, 0x26, 0x02, 0x01, 0x11, 0x01, 0x10, 0x28, 0x29, 0xa7} // byte 0x11; byte 0x10; b>=
	err = CheckProgram(program, args)
	require.NoError(t, err)

	// b^
	program = []byte{0x04, 0x26, 0x03, 0x01, 0x11, 0x01, 0x10, 0x01, 0x01, 0x28, 0x29, 0xad, 0x2a, 0x12} // byte 0x11; byte 0x10; b^; byte 0x01; ==
	err = CheckProgram(program, args)
	require.NoError(t, err)

	// callsub, retsub.
	program = []byte{0x04, 0x20, 0x02, 0x01, 0x02, 0x22, 0x88, 0x00, 0x03, 0x23, 0x12, 0x43, 0x49, 0x08, 0x89} // int 1; callsub double; int 2; ==; return; double: dup; +; retsub;
	err = CheckProgram(program, args)
	require.NoError(t, err)

	// loop
	program = []byte{0x04, 0x20, 0x04, 0x01, 0x02, 0x0a, 0x10, 0x22, 0x23, 0x0b, 0x49, 0x24, 0x0c, 0x40, 0xff, 0xf8, 0x25, 0x12} // int 1; loop: int 2; *; dup; int 10; <; bnz loop; int 16; ==
	err = CheckProgram(program, args)
	require.NoError(t, err)
}

func TestCheckProgramV5(t *testing.T) {
	// check TEAL v5 opcodes
	require.True(t, spec.EvalMaxVersion >= 5)

	args := make([][]byte, 0)

	// itxn ops
	program, err := hex.DecodeString("052001c0843db18101b21022b2083100b207b3b4082212")
	// itxn_begin; int pay; itxn_field TypeEnum; int 1000000; itxn_field Amount; txn Sender; itxn_field Receiver; itxn_submit; itxn Amount; int 1000000; ==
	require.NoError(t, err)
	err = CheckProgram(program, args)
	require.NoError(t, err)

	// ECDSA ops
	program, err = hex.DecodeString("058008746573746461746103802079bfa8245aeac0e714b7bd2b3252d03979e5e7a43cb039715a5f8109a7dd9ba180200753d317e54350d1d102289afbde3002add4529f10b9f7d3d223843985de62e0802103abfb5e6e331fb871e423f354e2bd78a384ef7cb07ac8bbf27d2dd1eca00e73c106000500")
	// byte "testdata"; sha512_256; byte 0x79bfa8245aeac0e714b7bd2b3252d03979e5e7a43cb039715a5f8109a7dd9ba1; byte 0x0753d317e54350d1d102289afbde3002add4529f10b9f7d3d223843985de62e0; byte 0x03abfb5e6e331fb871e423f354e2bd78a384ef7cb07ac8bbf27d2dd1eca00e73c1; ecdsa_pk_decompress Secp256k1; ecdsa_verify Secp256k1
	require.NoError(t, err)
	err = CheckProgram(program, args)
	require.NoError(t, err)

	// cover, uncover, log
	program, err = hex.DecodeString("058001618001628001634e024f025050b08101")
	// byte "a"; byte "b"; byte "c"; cover 2; uncover 2; concat; concat; log; int 1
	require.NoError(t, err)
	err = CheckProgram(program, args)
	require.NoError(t, err)
}

func TestCheckProgramV6(t *testing.T) {
	// check TEAL v6 opcodes
	require.True(t, spec.EvalMaxVersion >= 6)

	args := make([][]byte, 0)

	// bsqrt
	program, err := hex.DecodeString("068001909680010c12")
	// byte 0x90; bsqrt; byte 0x0c; ==
	require.NoError(t, err)
	err = CheckProgram(program, args)
	require.NoError(t, err)

	// divw
	program, err = hex.DecodeString("06810981ecffffffffffffffff01810a9781feffffffffffffffff0112")
	// int 9; int 18446744073709551596; int 10; divw; int 18446744073709551614; ==
	require.NoError(t, err)
	err = CheckProgram(program, args)
	require.NoError(t, err)

	// txn fields
	program, err = hex.DecodeString("06313f1581401233003e15810a1210")
	// txn StateProofPK; len; int 64; ==; gtxn 0 LastLog; len; int 10; ==; &&
	require.NoError(t, err)
	err = CheckProgram(program, args)
	require.NoError(t, err)
}
