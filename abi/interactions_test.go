package abi

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMethodFromSignature(t *testing.T) {
	Uint32, err := TypeOf("uint32")
	require.NoError(t, err)

	expectedArgs := []Arg{
		{Name: "", Type: "uint32", typeObject: &Uint32, Desc: ""},
		{Name: "", Type: "uint32", typeObject: &Uint32, Desc: ""},
	}
	expected := Method{
		Name:    "add",
		Desc:    "",
		Args:    expectedArgs,
		Returns: Return{Type: "uint32", typeObject: &Uint32, Desc: ""},
	}

	methodSig := "add(uint32,uint32)uint32"
	result, err := MethodFromSignature(methodSig)

	require.NoError(t, err)
	require.Equal(t, result, expected)
}

func TestMethodFromSignatureWithTuple(t *testing.T) {
	Uint32, err := TypeOf("uint32")
	require.NoError(t, err)

	Uint32Tuple, err := TypeOf("(uint32,uint32)")
	require.NoError(t, err)

	Uint32TupleTuple, err := TypeOf("(uint32,(uint32,uint32))")
	require.NoError(t, err)

	expectedArgs := []Arg{
		{Name: "", Type: "(uint32,(uint32,uint32))", typeObject: &Uint32TupleTuple, Desc: ""},
		{Name: "", Type: "uint32", typeObject: &Uint32, Desc: ""},
	}
	expected := Method{
		Name:    "add",
		Desc:    "",
		Args:    expectedArgs,
		Returns: Return{Type: "(uint32,uint32)", typeObject: &Uint32Tuple, Desc: ""},
	}

	methodSig := "add((uint32,(uint32,uint32)),uint32)(uint32,uint32)"
	result, err := MethodFromSignature(methodSig)

	require.NoError(t, err)
	require.Equal(t, expected, result)
}

func TestMethodFromSignatureWithVoidReturn(t *testing.T) {
	Uint32, err := TypeOf("uint32")
	require.NoError(t, err)

	expectedArgs := []Arg{
		{Name: "", Type: "uint32", typeObject: &Uint32, Desc: ""},
		{Name: "", Type: "uint32", typeObject: &Uint32, Desc: ""},
	}
	expected := Method{
		Name:    "add",
		Desc:    "",
		Args:    expectedArgs,
		Returns: Return{Type: "void", Desc: ""},
	}

	methodSig := "add(uint32,uint32)void"
	result, err := MethodFromSignature(methodSig)

	require.NoError(t, err)
	require.Equal(t, expected, result)
}

func TestMethodFromSignatureWithNoArgs(t *testing.T) {
	expectedArgs := []Arg{}
	expected := Method{
		Name:    "add",
		Desc:    "",
		Args:    expectedArgs,
		Returns: Return{Type: "void", Desc: ""},
	}

	methodSig := "add()void"
	result, err := MethodFromSignature(methodSig)

	require.NoError(t, err)
	require.Equal(t, expected, result)
}

func TestMethodFromSignatureInvalidFormat(t *testing.T) {
	methodSig := "add)uint32,uint32)uint32"
	_, err := MethodFromSignature(methodSig)
	require.Error(t, err)

	methodSig = "add(uint32, uint32)uint32"
	_, err = MethodFromSignature(methodSig)
	require.Error(t, err)

	methodSig = "(uint32,uint32)uint32"
	_, err = MethodFromSignature(methodSig)
	require.Error(t, err)

	methodSig = "add((uint32, uint32)uint32"
	_, err = MethodFromSignature(methodSig)
	require.Error(t, err)
}

func TestMethodFromSignatureInvalidAbiType(t *testing.T) {
	methodSig := "add(uint32,uint32)int32"
	_, err := MethodFromSignature(methodSig)
	require.Error(t, err)
}

func TestGetSignature(t *testing.T) {
	expectedArgs := []Arg{
		{Name: "", Type: "uint32", Desc: ""},
		{Name: "", Type: "uint32", Desc: ""},
	}

	method := Method{
		Name:    "add",
		Desc:    "",
		Args:    expectedArgs,
		Returns: Return{Type: "uint32", Desc: ""},
	}

	expected := "add(uint32,uint32)uint32"
	require.Equal(t, method.GetSignature(), expected)
}

func TestGetSelector(t *testing.T) {
	args := []Arg{
		{Name: "", Type: "uint32", Desc: ""},
		{Name: "", Type: "uint32", Desc: ""},
	}

	method := Method{
		Name:    "add",
		Desc:    "",
		Args:    args,
		Returns: Return{Type: "uint32", Desc: ""},
	}

	expected := []byte{0x3e, 0x1e, 0x52, 0xbd}
	require.Equal(t, method.GetSelector(), expected)
}

func TestEncodeJsonMethod(t *testing.T) {
	args := []Arg{
		{Name: "0", Type: "uint32", Desc: ""},
		{Name: "1", Type: "uint32", Desc: ""},
	}

	method := Method{
		Name:    "add",
		Desc:    "",
		Args:    args,
		Returns: Return{Type: "uint32", Desc: ""},
	}

	expected := `{"name":"add","args":[{"name":"0","type":"uint32"},{"name":"1","type":"uint32"}],"returns":{"type":"uint32"}}`

	jsonMethod, err := json.Marshal(method)
	require.NoError(t, err)
	require.Equal(t, expected, string(jsonMethod))
}

func TestEncodeJsonMethodWithDescription(t *testing.T) {
	args := []Arg{
		{Name: "0", Type: "uint32", Desc: "description"},
		{Name: "1", Type: "uint32", Desc: "description"},
	}

	method := Method{
		Name:    "add",
		Desc:    "description",
		Args:    args,
		Returns: Return{Type: "uint32", Desc: "description"},
	}

	expected := `{"name":"add","desc":"description","args":[{"name":"0","type":"uint32","desc":"description"},{"name":"1","type":"uint32","desc":"description"}],"returns":{"type":"uint32","desc":"description"}}`

	jsonMethod, err := json.Marshal(method)
	require.NoError(t, err)
	require.Equal(t, expected, string(jsonMethod))
}

func TestEncodeJsonInterface(t *testing.T) {
	args := []Arg{
		{Name: "0", Type: "uint32", Desc: ""},
		{Name: "1", Type: "uint32", Desc: ""},
	}

	method := Method{
		Name:    "add",
		Desc:    "",
		Args:    args,
		Returns: Return{Type: "uint32", Desc: ""},
	}

	interfaceObject := Interface{
		Name:    "interface",
		Methods: []Method{method},
	}

	expected := `{"name":"interface","methods":[{"name":"add","args":[{"name":"0","type":"uint32"},{"name":"1","type":"uint32"}],"returns":{"type":"uint32"}}]}`

	jsonInterface, err := json.Marshal(interfaceObject)
	require.NoError(t, err)
	require.Equal(t, expected, string(jsonInterface))
}

func TestEncodeJsonInterfaceWithDescription(t *testing.T) {
	args := []Arg{
		{Name: "0", Type: "uint32", Desc: "description"},
		{Name: "1", Type: "uint32", Desc: "description"},
	}

	method := Method{
		Name:    "add",
		Desc:    "description",
		Args:    args,
		Returns: Return{Type: "uint32", Desc: "description"},
	}

	interfaceObject := Interface{
		Name:    "interface",
		Methods: []Method{method},
	}

	expected := `{"name":"interface","methods":[{"name":"add","desc":"description","args":[{"name":"0","type":"uint32","desc":"description"},{"name":"1","type":"uint32","desc":"description"}],"returns":{"type":"uint32","desc":"description"}}]}`

	jsonInterface, err := json.Marshal(interfaceObject)
	require.NoError(t, err)
	require.Equal(t, expected, string(jsonInterface))
}

func TestEncodeJsonContract(t *testing.T) {
	args := []Arg{
		{Name: "0", Type: "uint32", Desc: ""},
		{Name: "1", Type: "uint32", Desc: ""},
	}

	method := Method{
		Name:    "add",
		Desc:    "",
		Args:    args,
		Returns: Return{Type: "uint32", Desc: ""},
	}

	network := ContractNetworkInfo{AppID: 123}

	contract := Contract{
		Name:     "contract",
		Networks: map[string]ContractNetworkInfo{"genesis hash": network},
		Methods:  []Method{method},
	}

	expected := `{"name":"contract","networks":{"genesis hash":{"appID":123}},"methods":[{"name":"add","args":[{"name":"0","type":"uint32"},{"name":"1","type":"uint32"}],"returns":{"type":"uint32"}}]}`

	jsonContract, err := json.Marshal(contract)
	require.NoError(t, err)
	require.Equal(t, expected, string(jsonContract))
}

func TestEncodeJsonContractWithDescription(t *testing.T) {
	args := []Arg{
		{Name: "0", Type: "uint32", Desc: "description"},
		{Name: "1", Type: "uint32", Desc: "description"},
	}

	method := Method{
		Name:    "add",
		Desc:    "description",
		Args:    args,
		Returns: Return{Type: "uint32", Desc: "description"},
	}

	network := ContractNetworkInfo{AppID: 123}

	contract := Contract{
		Name:     "contract",
		Desc:     "description for contract",
		Networks: map[string]ContractNetworkInfo{"genesis hash": network},
		Methods:  []Method{method},
	}

	expected := `{"name":"contract","desc":"description for contract","networks":{"genesis hash":{"appID":123}},"methods":[{"name":"add","desc":"description","args":[{"name":"0","type":"uint32","desc":"description"},{"name":"1","type":"uint32","desc":"description"}],"returns":{"type":"uint32","desc":"description"}}]}`

	jsonContract, err := json.Marshal(contract)
	require.NoError(t, err)
	require.Equal(t, expected, string(jsonContract))
}
