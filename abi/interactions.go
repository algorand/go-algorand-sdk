package abi

import (
	"strconv"
	"strings"

	"github.com/algorand/go-algorand/data/abi"
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

func MethodFromString(methodStr string) (Method, bool) {
	openCnt := strings.Count(methodStr, "(")
	if openCnt == 0 || openCnt > 1 {
		return Method{}, false
	}

	closeCnt := strings.Count(methodStr, ")")
	if closeCnt == 0 || closeCnt > 1 {
		return Method{}, false
	}

	openIdx := strings.Index(methodStr, "(")
	closeIdx := strings.Index(methodStr, ")")
	if openIdx > closeIdx {
		return Method{}, false
	}

	name := methodStr[:openIdx]
	stringArgs := strings.Split(methodStr[openIdx+1:closeIdx], ",")
	returnType := methodStr[closeIdx+1:]

	for _, argType := range stringArgs {
		_, err := abi.TypeFromString(argType)
		if err != nil {
			return Method{}, false
		}
	}

	_, err := abi.TypeFromString(returnType)
	if err != nil {
		return Method{}, false
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
	}, true
}

type Interface struct {
	Name    string   `json:"name"`
	Methods []Method `json:"methods"`
}

type Contract struct {
	Name    string   `json:"name"`
	AppId   string   `json:"appId"`
	Methods []Method `json:"methods"`
}
