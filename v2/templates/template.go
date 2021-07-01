package templates

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

// ContractTemplate template representation
type ContractTemplate struct {
	address string
	program []byte
}

// GetAddress returns the contract address
func (contract ContractTemplate) GetAddress() string {
	return contract.address
}

// GetProgram returns the program bytes
func (contract ContractTemplate) GetProgram() []byte {
	return contract.program
}

func replace(buf, newBytes []byte, offset, placeholderLength uint64) []byte {
	firstChunk := make([]byte, len(buf[:offset]))
	copy(firstChunk, buf[:offset])
	firstChunkAmended := append(firstChunk, newBytes...)
	secondChunk := make([]byte, len(buf[(offset+placeholderLength):]))
	copy(secondChunk, buf[(offset+placeholderLength):])
	return append(firstChunkAmended, secondChunk...)
}

func inject(original []byte, offsets []uint64, values []interface{}) (result []byte, err error) {
	result = original
	if len(offsets) != len(values) {
		err = fmt.Errorf("length of offsets %v does not match length of replacement values %v", len(offsets), len(values))
		return
	}

	for i, value := range values {
		decodedLength := 0
		if valueAsUint, ok := value.(uint64); ok {
			// make the exact minimum buffer needed and no larger
			// because otherwise there will be extra bytes inserted
			sizingBuffer := make([]byte, binary.MaxVarintLen64)
			decodedLength = binary.PutUvarint(sizingBuffer, valueAsUint)
			fillingBuffer := make([]byte, decodedLength)
			decodedLength = binary.PutUvarint(fillingBuffer, valueAsUint)
			result = replace(result, fillingBuffer, offsets[i], uint64(1))
		} else if address, ok := value.(types.Address); ok {
			addressLen := uint64(32)
			addressBytes := make([]byte, addressLen)
			copy(addressBytes, address[:])
			result = replace(result, addressBytes, offsets[i], addressLen)
		} else if b64string, ok := value.(string); ok {
			decodeBytes, decodeErr := base64.StdEncoding.DecodeString(b64string)
			if decodeErr != nil {
				err = decodeErr
				return
			}
			// do the same thing as in the uint64 case to trim empty bytes:
			// first fill one buffer to figure out the number of bytes to be written,
			// then fill a second buffer of exactly the right size
			sizingBuffer := make([]byte, binary.MaxVarintLen64)
			numBytesWritten := binary.PutUvarint(sizingBuffer, uint64(len(decodeBytes)))
			fillingBuffer := make([]byte, numBytesWritten)
			binary.PutUvarint(fillingBuffer, uint64(len(decodeBytes))) // indicate length of b64 bytes
			// want to write [length of b64 bytes, b64 bytes]
			decodeBytes = append(fillingBuffer, decodeBytes...)
			result = replace(result, decodeBytes, offsets[i], uint64(33))
		}

		if decodedLength != 0 {
			for j := range offsets {
				offsets[j] = offsets[j] + uint64(decodedLength) - 1
			}
		}
	}
	return
}
