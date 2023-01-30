package types

import (
	"bytes"
	"errors"
	"fmt"
)

func (b BlockHash) String() string {
	return fmt.Sprintf("blk-%v", Digest(b))
}

// MarshalText returns the BlockHash string as an array of bytes
func (b BlockHash) MarshalText() ([]byte, error) {
	return []byte(b.String()), nil
}

// UnmarshalText initializes the BlockHash from an array of bytes.
func (b *BlockHash) UnmarshalText(text []byte) error {
	if len(text) < 4 || !bytes.Equal(text[0:4], []byte("blk-")) {
		return errors.New("unrecognized blockhash format")
	}
	d, err := DigestFromString(string(text[4:]))
	*b = BlockHash(d)
	return err
}
