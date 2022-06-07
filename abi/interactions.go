package abi

import (
	"crypto/sha512"
	"fmt"
	"strings"

	"github.com/algorand/go-algorand-sdk/types"
)

// AnyTransactionType is the ABI argument type string for a nonspecific
// transaction argument
const AnyTransactionType = "txn"

var transactionArgTypes = map[string]interface{}{
	AnyTransactionType:              nil, // denotes a placeholder for any of the size types below
	string(types.PaymentTx):         nil,
	string(types.KeyRegistrationTx): nil,
	string(types.AssetConfigTx):     nil,
	string(types.AssetTransferTx):   nil,
	string(types.AssetFreezeTx):     nil,
	string(types.ApplicationCallTx): nil,
}

// IsTransactionType checks if a type string represents a transaction type
// argument, such as "txn", "pay", "keyreg", etc.
func IsTransactionType(typeStr string) bool {
	_, ok := transactionArgTypes[typeStr]
	return ok
}

// AccountReferenceType is the ABI argument type string for account references
const AccountReferenceType = "account"

// AssetReferenceType is the ABI argument type string for asset references
const AssetReferenceType = "asset"

// ApplicationReferenceType is the ABI argument type string for application
// references
const ApplicationReferenceType = "application"

var referenceArgTypes = map[string]interface{}{
	AccountReferenceType:     nil,
	AssetReferenceType:       nil,
	ApplicationReferenceType: nil,
}

// IsReferenceType checks if a type string represents a reference type argument,
// such as "account", "asset", or "application".
func IsReferenceType(typeStr string) bool {
	_, ok := referenceArgTypes[typeStr]
	return ok
}

// Arg represents an ABI Method argument
type Arg struct {
	// Optional, user-friendly name for the argument
	Name string `json:"name,omitempty"`
	// The type of the argument as a string. See the method GetTypeObject to
	// obtain the ABI type object
	Type string `json:"type"`
	// A hidden type object cache that holds the parsed type object
	typeObject *Type `json:"-"`
	// Optional, user-friendly description for the argument
	Desc string `json:"desc,omitempty"`
}

// IsTransactionArg checks if this argument's type is a transaction type
func (a Arg) IsTransactionArg() bool {
	return IsTransactionType(a.Type)
}

// IsReferenceArg checks if this argument's type is a reference type
func (a Arg) IsReferenceArg() bool {
	return IsReferenceType(a.Type)
}

// GetTypeObject parses and returns the ABI type object for this argument's
// type. An error will be returned if this argument's type is a transaction or
// reference type
func (a *Arg) GetTypeObject() (Type, error) {
	if a.IsTransactionArg() {
		return Type{}, fmt.Errorf("Invalid operation on transaction type %s", a.Type)
	}
	if a.IsReferenceArg() {
		return Type{}, fmt.Errorf("Invalid operation on reference type %s", a.Type)
	}
	if a.typeObject != nil {
		return *a.typeObject, nil
	}
	typeObject, err := TypeOf(a.Type)
	if err == nil {
		a.typeObject = &typeObject
	}
	return typeObject, err
}

// VoidReturnType is the ABI return type string for a method that does not have
// a return value
const VoidReturnType = "void"

// Return represents an ABI method return value
type Return struct {
	// The type of the return value as a string. See the method GetTypeObject to
	// obtain the ABI type object
	Type string `json:"type"`
	// A hidden type object cache that holds the parsed type object
	typeObject *Type `json:"-"`
	// Optional, user-friendly description for the return value
	Desc string `json:"desc,omitempty"`
}

// IsVoid checks if this return type is void, meaning the method does not have
// any return value
func (r Return) IsVoid() bool {
	return r.Type == VoidReturnType
}

// GetTypeObject parses and returns the ABI type object for this return type.
// An error will be returned if this is a void return type.
func (r *Return) GetTypeObject() (Type, error) {
	if r.IsVoid() {
		return Type{}, fmt.Errorf("Invalid operation on void return type")
	}
	if r.typeObject != nil {
		return *r.typeObject, nil
	}
	typeObject, err := TypeOf(r.Type)
	if err == nil {
		r.typeObject = &typeObject
	}
	return typeObject, err
}

// Method represents an individual ABI method definition
type Method struct {
	// The name of the method
	Name string `json:"name"`
	// Optional, user-friendly description for the method
	Desc string `json:"desc,omitempty"`
	// The arguments of the method, in order
	Args []Arg `json:"args"`
	// Information about the method's return value
	Returns Return `json:"returns"`
}

// parseMethodArgs parses the arguments from a method signature string.
// trMethod is the complete method signature and startIdx is the index of the
// opening parenthesis of the arguments list. This function returns a list of
// the argument types from the method signature and the index of the closing
// parenthesis of the arguments list.
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
			return nil, -1, fmt.Errorf("method signature parentheses mismatch")
		} else if parenCnt > 1 {
			continue
		}

		if strMethod[curPos] == ',' || parenCnt == 0 {
			strArg := strMethod[prevPos:curPos]
			argTypes = append(argTypes, strArg)
			prevPos = curPos + 1
		}

		if parenCnt == 0 {
			closeIdx = curPos
			break
		}
	}

	if closeIdx == -1 {
		return nil, -1, fmt.Errorf("method signature parentheses mismatch")
	}

	return argTypes, closeIdx, nil
}

// MethodFromSignature decoded a method signature string into a Method object.
func MethodFromSignature(methodStr string) (Method, error) {
	openIdx := strings.Index(methodStr, "(")
	if openIdx == -1 {
		return Method{}, fmt.Errorf("method signature is missing an open parenthesis")
	}

	name := methodStr[:openIdx]
	if name == "" {
		return Method{}, fmt.Errorf("method must have a non empty name")
	}

	argTypes, closeIdx, err := parseMethodArgs(methodStr, openIdx)
	if err != nil {
		return Method{}, err
	}

	returnType := Return{Type: methodStr[closeIdx+1:]}
	if !returnType.IsVoid() {
		// fill type object cache and catch any errors
		_, err := returnType.GetTypeObject()
		if err != nil {
			return Method{}, err
		}
	}

	args := make([]Arg, len(argTypes))
	for i, argType := range argTypes {
		args[i].Type = argType

		if IsTransactionType(argType) || IsReferenceType(argType) {
			continue
		}

		// fill type object cache and catch any errors
		_, err := args[i].GetTypeObject()
		if err != nil {
			return Method{}, err
		}
	}

	return Method{
		Name:    name,
		Args:    args,
		Returns: returnType,
	}, nil
}

// GetSignature calculates and returns the signature of the method
func (method *Method) GetSignature() string {
	var methodSignature string
	methodSignature += method.Name + "("

	var strTypes []string
	for _, arg := range method.Args {
		strTypes = append(strTypes, arg.Type)
	}

	methodSignature += strings.Join(strTypes, ",")
	methodSignature += ")"
	methodSignature += method.Returns.Type

	return methodSignature
}

// GetSelector calculates and returns the 4-byte selector of the method
func (method *Method) GetSelector() []byte {
	sig := method.GetSignature()
	sigHash := sha512.Sum512_256([]byte(sig))
	return sigHash[:4]
}

// GetTxCount returns the number of transactions required to invoke this method
func (method *Method) GetTxCount() int {
	cnt := 1
	for _, arg := range method.Args {
		if arg.IsTransactionArg() {
			cnt++
		}
	}
	return cnt
}

// Interface represents an ABI interface, which is a logically grouped
// collection of methods
type Interface struct {
	// A user-friendly name for the interface
	Name string `json:"name"`
	// Optional, user-friendly description for the interface
	Desc string `json:"desc,omitempty"`
	// The methods that the interface contains
	Methods []Method `json:"methods"`
}

func (i *Interface) GetMethodByName(name string) (Method, error) {
	var methods []Method
	for _, method := range i.Methods {
		if method.Name == name {
			methods = append(methods, method)
		}
	}

	if len(methods) > 1 {
		var sigs []string
		for _, method := range methods {
			sigs = append(sigs, method.GetSignature())
		}

		return Method{}, fmt.Errorf("found %d methods with the same name %s", len(methods), strings.Join(sigs, ","))
	}

	if len(methods) == 0 {
		return Method{}, fmt.Errorf("found 0 methods with the name %s", name)
	}

	return methods[0], nil
}

// ContractNetworkInfo contains network-specific information about the contract
type ContractNetworkInfo struct {
	// The application ID of the contract for this network
	AppID uint64 `json:"appID"`
}

// Contract represents an ABI contract, which is a concrete set of methods
// implemented by a single app
type Contract struct {
	// A user-friendly name for the contract
	Name string `json:"name"`
	// Optional, user-friendly description for the contract
	Desc string `json:"desc,omitempty"`
	// Optional information about the contract's instances across different
	// networks
	Networks map[string]ContractNetworkInfo `json:"networks,omitempty"`
	// The methods that the contract implements
	Methods []Method `json:"methods"`
}

func (c *Contract) GetMethodByName(name string) (Method, error) {
	var methods []Method
	for _, method := range c.Methods {
		if method.Name == name {
			methods = append(methods, method)
		}
	}

	if len(methods) > 1 {
		var sigs []string
		for _, method := range methods {
			sigs = append(sigs, method.GetSignature())
		}

		return Method{}, fmt.Errorf("found %d methods with the same name %s", len(methods), strings.Join(sigs, ","))
	}

	if len(methods) == 0 {
		return Method{}, fmt.Errorf("found 0 methods with the name %s", name)
	}

	return methods[0], nil
}
