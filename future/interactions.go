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

func MethodFromSignature(methodStr string) (Method, error) {
	openCnt := strings.Count(methodStr, "(")
	if openCnt == 0 {
		return Method{}, errors.New("Method signature has invalid format")
	}

	closeCnt := strings.Count(methodStr, ")")
	if closeCnt != openCnt {
		return Method{}, errors.New("Method signature has invalid format")
	}

	openIdx := strings.Index(methodStr, "(")
	parenCnt := 0
	closeIdx := -1
	for pos := openIdx; pos < len(methodStr); pos++ {
		if methodStr[pos] == '(' {
			parenCnt++
		} else if methodStr[pos] == ')' {
			parenCnt--
		}

		if parenCnt == 0 {
			closeIdx = pos
			break
		}
	}

	if closeIdx == -1 {
		return Method{}, errors.New("Method signature has invalid format")
	}

	// signature must include method name
	if openIdx == 0 {
		return Method{}, errors.New("Method signature has invalid format")
	}

	name := methodStr[:openIdx]
	match, err := regexp.MatchString(`[_a-zA-Z][_a-zA-Z0-9]+`, name)
	if err != nil {
		return Method{}, err
	}
	if !match {
		return Method{}, errors.New("Method signature has invalid format")
	}

	stringArgs := methodStr[openIdx+1 : closeIdx]
	var argTypes []string
	parenCnt = 0
	prevPos := 0
	for curPos, c := range stringArgs {
		if c == '(' {
			parenCnt++
		} else if c == ')' {
			parenCnt--
		}

		if parenCnt < 0 {
			return Method{}, errors.New("Method Signature has invalid format")
		} else if parenCnt > 0 {
			continue
		}

		// parenCnt == 0 indicates top level comma
		if c == ',' {
			strArg := stringArgs[prevPos:curPos]
			_, err := abi.TypeOf(strArg)
			if err != nil {
				return Method{}, err
			}

			argTypes = append(argTypes, strArg)
			prevPos = curPos + 1
		}
	}

	// must process final arg
	strArg := stringArgs[prevPos:]
	_, err = abi.TypeOf(strArg)
	if err != nil {
		return Method{}, err
	}
	argTypes = append(argTypes, strArg)

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
