package future

import (
	"crypto/sha512"
	"errors"
	"regexp"
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

func parseMethodArgs(strMethod string, startIdx int) ([]string, int, error) {
	// handle no args
	if startIdx < len(strMethod)-1 && strMethod[startIdx+1] == ')' {
		return []string{}, startIdx + 1, nil
	}

	var argTypes []string
	parenCnt := 1
	prevPos := startIdx + 1
	closeIdx := -1
	for curPos := prevPos; curPos < len(strMethod); curPos++ {
		if strMethod[curPos] == '(' {
			parenCnt++
		} else if strMethod[curPos] == ')' {
			parenCnt--
		}

		if parenCnt < 0 {
			return nil, -1, errors.New("Method Signature has invalid format")
		} else if parenCnt > 1 {
			continue
		}

		if strMethod[curPos] == ',' || parenCnt == 0 {
			strArg := strMethod[prevPos:curPos]
			_, err := abi.TypeOf(strArg)
			if err != nil {
				return nil, -1, err
			}

			argTypes = append(argTypes, strArg)
			prevPos = curPos + 1
		}

		if parenCnt == 0 {
			closeIdx = curPos
			break
		}
	}

	if closeIdx == -1 {
		return nil, -1, errors.New("Method Signature has invalid format")
	}

	return argTypes, closeIdx, nil
}

func MethodFromSignature(methodStr string) (Method, error) {
	openCnt := strings.Count(methodStr, "(")
	if openCnt == 0 {
		return Method{}, errors.New("Method signature has invalid format")
	}

	openIdx := strings.Index(methodStr, "(")
	name := methodStr[:openIdx]
	match, err := regexp.MatchString(`[_a-zA-Z][_a-zA-Z0-9]+`, name)
	if err != nil {
		return Method{}, err
	}
	if !match {
		return Method{}, errors.New("Method signature has invalid format")
	}

	argTypes, closeIdx, err := parseMethodArgs(methodStr, openIdx)
	if err != nil {
		return Method{}, err
	}

	returnType := methodStr[closeIdx+1:]
	if returnType != "void" {
		_, err := abi.TypeOf(returnType)
		if err != nil {
			return Method{}, err
		}
	}

	args := make([]Arg, len(argTypes))
	for i, argType := range argTypes {
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
