package types

import (
	"encoding/base64"
	"strings"
)

// MarshalText returns the BlockHash string as an array of bytes
func (b *BlockHash) MarshalText() ([]byte, error) {
	result := base64.StdEncoding.EncodeToString(b[:])
	return []byte(result), nil
}

// UnmarshalText initializes the BlockHash from an array of bytes.
func (b *BlockHash) UnmarshalText(text []byte) error {
	// Remove the blk- prefix if it is present to allow decoding either format.
	if strings.HasPrefix(string(text), "blk-") {
		text = text[4:]
	}

	// Attempt to decode base32 format
	d, err := DigestFromString(string(text))
	if err == nil {
		*b = BlockHash(d)
		return nil
	}
	// ignore the DigestFromString error because it isn't the native MarshalText format.

	// Attempt to decode base64 format
	var data BlockHash
	n, err := base64.StdEncoding.Decode(data[:], text)
	if err == nil {
		if n != len(b[:]) {
			return errWrongBlockHashLen
		}
		copy(b[:], data[:])
		return nil
	}
	return err
}
