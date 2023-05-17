package abi

import (
	"crypto/sha512"
	"fmt"
	"strings"

	avm_abi "github.com/algorand/avm-abi/abi"
)

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

// MethodFromSignature decoded a method signature string into a Method object.
func MethodFromSignature(methodStr string) (Method, error) {
	name, argTypesStr, retTypeStr, err := avm_abi.ParseMethodSignature(methodStr)
	if err != nil {
		return Method{}, err
	}

	returnType := Return{Type: retTypeStr}
	if !returnType.IsVoid() {
		// fill type object cache and catch any errors
		_, err := returnType.GetTypeObject()
		if err != nil {
			return Method{}, fmt.Errorf("Could not parse method return type: %w", err)
		}
	}

	args := make([]Arg, len(argTypesStr))
	for i, argTypeStr := range argTypesStr {
		args[i].Type = argTypeStr

		if IsTransactionType(argTypeStr) || IsReferenceType(argTypeStr) {
			continue
		}

		// fill type object cache and catch any errors
		_, err := args[i].GetTypeObject()
		if err != nil {
			return Method{}, fmt.Errorf("Could not parse argument type at index %d: %w", i, err)
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

// GetMethodByName returns the method with the given name from the given list.
// Returns an error if there are multiple or no methods with the same name.
func GetMethodByName(methods []Method, name string) (Method, error) {
	var filteredMethods []Method
	for _, method := range methods {
		if method.Name == name {
			filteredMethods = append(filteredMethods, method)
		}
	}

	if len(filteredMethods) > 1 {
		var sigs []string
		for _, method := range filteredMethods {
			sigs = append(sigs, method.GetSignature())
		}

		return Method{}, fmt.Errorf("found %d methods with the same name %s", len(filteredMethods), strings.Join(sigs, ","))
	}

	if len(filteredMethods) == 0 {
		return Method{}, fmt.Errorf("found 0 methods with the name %s", name)
	}

	return filteredMethods[0], nil
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

// GetMethodByName returns the method with the given name
func (i *Interface) GetMethodByName(name string) (Method, error) {
	return GetMethodByName(i.Methods, name)
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

// GetMethodByName returns the method with the given name
func (c *Contract) GetMethodByName(name string) (Method, error) {
	return GetMethodByName(c.Methods, name)
}
