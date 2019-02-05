package msgpack

import (
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
