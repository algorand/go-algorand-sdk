package types

import (
	"bytes"
	"crypto/sha512"
	"encoding/base32"
	"encoding/base64"
)

const (
	checksumLenBytes = 4
	hashLenBytes     = sha512.Size256
)

// Address represents an Algorand address.
type Address [hashLenBytes]byte

// String grabs a human-readable representation of the address. This
// representation includes a 4-byte checksum.
func (a Address) String() string {
	// Compute the checksum
	checksumHash := sha512.Sum512_256(a[:])
	checksumLenBytes := checksumHash[hashLenBytes-checksumLenBytes:]

	// Append the checksum and encode as base32
	checksumAddress := append(a[:], checksumLenBytes...)
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(checksumAddress)
}

// ZeroAddress is Address with all zero bytes. For handy == != comparisons.
var ZeroAddress Address = [hashLenBytes]byte{}

// IsZero returs true if the Address is all zero bytes.
func (a Address) IsZero() bool {
	return a == ZeroAddress
}

// MarshalText returns the address string as an array of bytes
func (addr *Address) MarshalText() ([]byte, error) {
	result := base64.StdEncoding.EncodeToString(addr[:])
	return []byte(result), nil
}

// UnmarshalText initializes the Address from an array of bytes.
// The bytes may be in the base32 checksum format, or the raw bytes base64 encoded.
func (addr *Address) UnmarshalText(text []byte) error {
	address, err := DecodeAddress(string(text))
	if err == nil {
		*addr = address
		return nil
	}
	// ignore the DecodeAddress error because it isn't the native MarshalText format.

	// Check if its b64 encoded
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err == nil {
		if len(data) != len(addr[:]) {
			return errWrongAddressLen
		}
		copy(addr[:], data[:])
		return nil
	}
	return err
}

// DecodeAddress turns a checksum address string into an Address object. It
// checks that the checksum is correct, and returns an error if it's not.
func DecodeAddress(addr string) (a Address, err error) {
	// Interpret the address as base32
	decoded, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(addr)
	if err != nil {
		return
	}

	// Ensure the decoded address is the correct length
	if len(decoded) != len(a)+checksumLenBytes {
		err = errWrongAddressLen
		return
	}

	// Split into address + checksum
	addressBytes := decoded[:len(a)]
	checksumBytes := decoded[len(a):]

	// Compute the expected checksum
	checksumHash := sha512.Sum512_256(addressBytes)
	expectedChecksumBytes := checksumHash[hashLenBytes-checksumLenBytes:]

	// Check the checksum
	if !bytes.Equal(expectedChecksumBytes, checksumBytes) {
		err = errWrongChecksum
		return
	}

	// Checksum is good, copy address bytes into output
	copy(a[:], addressBytes)
	return a, nil
}

// EncodeAddress turns a byte slice into the human readable representation of the address.
// This representation includes a 4-byte checksum
func EncodeAddress(addr []byte) (a string, err error) {
	if len(addr) != hashLenBytes {
		err = errWrongAddressByteLen
		return
	}

	var address Address
	copy(address[:], addr)
	a = address.String()

	return
}
