package future

import (
	"crypto/sha512"
	"errors"
	"strconv"
	"strings"

	"github.com/algorand/go-algorand-sdk/abi"
)

type Arg struct {
	Name    string `json:"name"`
	AbiType string `json:"type"`
	Desc    string `json:"desc,omitempty"`
}

type Return struct {
	AbiType string `json:"type"`
	Desc    string `json:"desc,omitempty"`
}

type Method struct {
	Name    string `json:"name"`
	Desc    string `json:"desc,omitempty"`
	Args    []Arg  `json:"args"`
	Returns Return `json:"returns"`
}

func MethodFromSignature(methodStr string) (Method, error) {
	openCnt := strings.Count(methodStr, "(")
	if openCnt == 0 || openCnt > 1 {
		return Method{}, errors.New("Method signature has invalid format")
	}

	closeCnt := strings.Count(methodStr, ")")
	if closeCnt == 0 || closeCnt > 1 {
		return Method{}, errors.New("Method signature has invalid format")
	}

	openIdx := strings.Index(methodStr, "(")
	closeIdx := strings.Index(methodStr, ")")
	if openIdx > closeIdx {
		return Method{}, errors.New("Method signature has invalid format")
	}

	// signature must include method name
	if openIdx == 0 {
		return Method{}, errors.New("Method signature has invalid format")
	}

	name := methodStr[:openIdx]
	stringArgs := strings.Split(methodStr[openIdx+1:closeIdx], ",")
	returnType := methodStr[closeIdx+1:]

	for _, argType := range stringArgs {
		_, err := abi.TypeOf(argType)
		if err != nil {
			return Method{}, err
		}
	}

	_, err := abi.TypeOf(returnType)
	if err != nil {
		return Method{}, err
	}

	args := make([]Arg, len(stringArgs))
	for i, argType := range stringArgs {
		args[i] = Arg{
			Name:    strconv.Itoa(i),
			AbiType: argType,
		}
	}

	return Method{
		Name:    name,
		Args:    args,
		Returns: Return{AbiType: returnType},
	}, nil
}

func (method Method) GetSignature() string {
	var methodSignature string
	methodSignature += method.Name + "("

	var strTypes []string
	for _, arg := range method.Args {
		strTypes = append(strTypes, arg.AbiType)
	}

	methodSignature += strings.Join(strTypes, ",")
	methodSignature += ")"
	methodSignature += method.Returns.AbiType

	return methodSignature
}

func (method Method) GetSelector() []byte {
	sig := method.GetSignature()
	sigHash := sha512.Sum512_256([]byte(sig))
	return sigHash[:4]
}

type Interface struct {
	Name    string   `json:"name"`
	Methods []Method `json:"methods"`
}

type Contract struct {
	Name    string   `json:"name"`
	AppId   uint64   `json:"appId"`
	Methods []Method `json:"methods"`
}
