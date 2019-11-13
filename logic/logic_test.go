package logic

import (
	"github.com/algorand/go-algorand-sdk/logic/templates"
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

func TestSplit(t *testing.T) {
	// Inputs
	owner := "WO3QIJ6T4DZHBX5PWJH26JLHFSRT7W7M2DJOULPXDTUS6TUX7ZRIO4KDFY"
	receivers := [2]string{"W6UUUSEAOGLBHT7VFT4H2SDATKKSG6ZBUIJXTZMSLW36YS44FRP5NVAU7U", "XCIBIN7RT4ZXGBMVAMU3QS6L5EKB7XGROC5EPCNHHYXUIBAA5Q6C5Y7NEU"}
	ratn, ratd := uint64(30), uint64(100)
	expiryRound := uint64(123456)
	minPay := uint64(100000)
	maxFee := uint64(1000)
	c, err := templates.MakeSplit(owner, receivers[0], receivers[1], ratn, ratd, expiryRound, minPay, maxFee)
	require.NoError(t, err)
	// outputs
	goldenAddress := "45L2LBNYS2WH3VHDHCGBL65QB7XMFU43RBW2IJYIB3NLNTHVLI5AP6TZVM"
	goldenProgram := "dHhuIFR5cGVFbnVtCmludCAxCj09CnR4biBGZWUKaW50IDEwMDAKPAomJgpnbG9iYWwgR3JvdXBTaXplCmludCAyCj09CmJueiBsYWJlbDAKdHhuIENsb3NlUmVtYWluZGVyVG8KYWRkciBXTzNRSUo2VDREWkhCWDVQV0pIMjZKTEhGU1JUN1c3TTJESk9VTFBYRFRVUzZUVVg3WlJJTzRLREZZCj09CnR4biBSZWNlaXZlcgpnbG9iYWwgWmVyb0FkZHJlc3MKPT0KdHhuIEFtb3VudAppbnQgMAo9PQp0eG4gRmlyc3RWYWxpZAppbnQgMTIzNDU2Cj4KJiYKJiYKJiYKaW50IDEKYm56IGxhYmVsMQpwb3AgLy8gY2hlY2sgYnVnIHdvcmthcm91bmQKbGFiZWwwOgpndHhuIDAgU2VuZGVyCmd0eG4gMSBTZW5kZXIKPT0KdHhuIENsb3NlUmVtYWluZGVyVG8KZ2xvYmFsIFplcm9BZGRyZXNzCj09Cmd0eG4gMCBSZWNlaXZlcgphZGRyIFc2VVVVU0VBT0dMQkhUN1ZGVDRIMlNEQVRLS1NHNlpCVUlKWFRaTVNMVzM2WVM0NEZSUDVOVkFVN1UKPT0KZ3R4biAxIFJlY2VpdmVyCmFkZHIgWENJQklON1JUNFpYR0JNVkFNVTNRUzZMNUVLQjdYR1JPQzVFUENOSEhZWFVJQkFBNVE2QzVZN05FVQo9PQpndHhuIDAgQW1vdW50Cmd0eG4gMCBBbW91bnQKZ3R4biAxIEFtb3VudAorCmludCAzMAoqCmludCAxMDAKLwo9PQpndHhuIDAgQW1vdW50CmludCAxMDAwMDAKPgomJgomJgomJgomJgomJgpsYWJlbDE6CiYmCg=="
	require.Equal(t, goldenAddress, c.GetAddress())
	require.Equal(t, goldenProgram, c.GetProgram())
}
