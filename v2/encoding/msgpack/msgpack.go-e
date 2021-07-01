package msgpack

import (
	"io"

	"github.com/algorand/go-codec/codec"
)

// CodecHandle is used to instantiate msgpack encoders and decoders
// with our settings (canonical, paranoid about decoding errors)
var CodecHandle *codec.MsgpackHandle

// init configures our msgpack encoder and decoder
func init() {
	CodecHandle = new(codec.MsgpackHandle)
	CodecHandle.ErrorIfNoField = true
	CodecHandle.ErrorIfNoArrayExpand = true
	CodecHandle.Canonical = true
	CodecHandle.RecursiveEmptyCheck = true
	CodecHandle.WriteExt = true
	CodecHandle.PositiveIntUnsigned = true
}

// Encode returns a msgpack-encoded byte buffer for a given object
func Encode(obj interface{}) []byte {
	var b []byte
	enc := codec.NewEncoderBytes(&b, CodecHandle)
	enc.MustEncode(obj)
	return b
}

// Decode attempts to decode a msgpack-encoded byte buffer into an
// object instance pointed to by objptr
func Decode(b []byte, objptr interface{}) error {
	dec := codec.NewDecoderBytes(b, CodecHandle)
	err := dec.Decode(objptr)
	if err != nil {
		return err
	}
	return nil
}

// NewDecoder returns a msgpack decoder
func NewDecoder(r io.Reader) *codec.Decoder {
	return codec.NewDecoder(r, CodecHandle)
}
