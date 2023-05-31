package abi

import (
	avm_abi "github.com/algorand/avm-abi/abi"
)

// Type is an ABI type
type Type = avm_abi.Type

// TypeOf returns the ABI type of a string representation of a type.
func TypeOf(str string) (Type, error) {
	return avm_abi.TypeOf(str)
}

// MakeTupleType makes tuple ABI type by taking an array of tuple element types as argument.
func MakeTupleType(argumentTypes []Type) (Type, error) {
	return avm_abi.MakeTupleType(argumentTypes)
}

// AnyTransactionType is the ABI argument type string for a nonspecific transaction argument
const AnyTransactionType = avm_abi.AnyTransactionType

// PaymentTransactionType is the ABI argument type string for a payment transaction argument
const PaymentTransactionType = avm_abi.PaymentTransactionType

// KeyRegistrationTransactionType is the ABI argument type string for a key registration transaction
// argument
const KeyRegistrationTransactionType = avm_abi.KeyRegistrationTransactionType

// AssetConfigTransactionType is the ABI argument type string for an asset configuration transaction
// argument
const AssetConfigTransactionType = avm_abi.AssetConfigTransactionType

// AssetTransferTransactionType is the ABI argument type string for an asset transfer transaction
// argument
const AssetTransferTransactionType = avm_abi.AssetTransferTransactionType

// AssetFreezeTransactionType is the ABI argument type string for an asset freeze transaction
// argument
const AssetFreezeTransactionType = avm_abi.AssetFreezeTransactionType

// ApplicationCallTransactionType is the ABI argument type string for an application call
// transaction argument
const ApplicationCallTransactionType = avm_abi.ApplicationCallTransactionType

// IsTransactionType checks if a type string represents a transaction type
// argument, such as "txn", "pay", "keyreg", etc.
func IsTransactionType(typeStr string) bool {
	return avm_abi.IsTransactionType(typeStr)
}

// AccountReferenceType is the ABI argument type string for account references
const AccountReferenceType = avm_abi.AccountReferenceType

// AssetReferenceType is the ABI argument type string for asset references
const AssetReferenceType = avm_abi.AssetReferenceType

// ApplicationReferenceType is the ABI argument type string for application references
const ApplicationReferenceType = avm_abi.ApplicationReferenceType

// IsReferenceType checks if a type string represents a reference type argument,
// such as "account", "asset", or "application".
func IsReferenceType(typeStr string) bool {
	return avm_abi.IsReferenceType(typeStr)
}

// VoidReturnType is the ABI return type string for a method that does not return any value
const VoidReturnType = avm_abi.VoidReturnType
