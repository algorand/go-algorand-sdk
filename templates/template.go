package templates

import (
	"encoding/binary"
	"fmt"
	"github.com/algorand/go-algorand-sdk/types"
)

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
			sizingBuffer := make([]byte, 32)
			decodedLength = binary.PutUvarint(sizingBuffer, valueAsUint)
			fillingBuffer := make([]byte, decodedLength)
			decodedLength = binary.PutUvarint(fillingBuffer, valueAsUint)
			result = replace(result, fillingBuffer, offsets[i], uint64(1))
		} else if addressString, ok := value.(string); ok {
			address, decodeErr := types.DecodeAddress(addressString)
			if decodeErr != nil {
				err = decodeErr // fix "err is shadowed during return" error
				return
			}
			addressLen := uint64(32)
			addressBytes := make([]byte, addressLen)
			copy(addressBytes, address[:])
			result = replace(result, addressBytes, offsets[i], addressLen)
		}

		if decodedLength != 0 {
			for j, _ := range offsets {
				offsets[j] = offsets[j] + uint64(decodedLength) - 1
			}
		}
	}
	return
}
