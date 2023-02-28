package json

import (
	"io"

	"github.com/algorand/go-codec/codec"
)

// CodecHandle is used to instantiate JSON encoders and decoders
// with our settings (canonical, paranoid about decoding errors)
var CodecHandle *codec.JsonHandle

// LenientCodecHandle is used to instantiate msgpack encoders for the REST API.
var LenientCodecHandle *codec.JsonHandle

// JSONStrictHandle is the same as CodecHandle but with MapKeyAsString=true
// for correct maps[int]interface{} encoding
var JSONStrictHandle *codec.JsonHandle

// init configures our json encoder and decoder
func init() {
	CodecHandle = new(codec.JsonHandle)
	CodecHandle.ErrorIfNoField = true
	CodecHandle.ErrorIfNoArrayExpand = true
	CodecHandle.Canonical = true
	CodecHandle.RecursiveEmptyCheck = true
	CodecHandle.Indent = 2
	CodecHandle.HTMLCharsAsIs = true

	LenientCodecHandle = new(codec.JsonHandle)
	LenientCodecHandle.ErrorIfNoField = false
	LenientCodecHandle.ErrorIfNoArrayExpand = true
	LenientCodecHandle.Canonical = true
	LenientCodecHandle.RecursiveEmptyCheck = true
	LenientCodecHandle.Indent = 2
	LenientCodecHandle.HTMLCharsAsIs = true

	JSONStrictHandle = new(codec.JsonHandle)
	JSONStrictHandle.ErrorIfNoField = CodecHandle.ErrorIfNoField
	JSONStrictHandle.ErrorIfNoArrayExpand = CodecHandle.ErrorIfNoArrayExpand
	JSONStrictHandle.Canonical = CodecHandle.Canonical
	JSONStrictHandle.RecursiveEmptyCheck = CodecHandle.RecursiveEmptyCheck
	JSONStrictHandle.Indent = CodecHandle.Indent
	JSONStrictHandle.HTMLCharsAsIs = CodecHandle.HTMLCharsAsIs
	JSONStrictHandle.MapKeyAsString = true
}

// Encode returns a json-encoded byte buffer for a given object
func Encode(obj interface{}) []byte {
	var b []byte
	enc := codec.NewEncoderBytes(&b, CodecHandle)
	enc.MustEncode(obj)
	return b
}

// EncodeStrict returns a JSON-encoded byte buffer for a given object
// It is the same Encode but encodes map's int keys as strings
func EncodeStrict(obj interface{}) []byte {
	var b []byte
	enc := codec.NewEncoderBytes(&b, JSONStrictHandle)
	enc.MustEncode(obj)
	return b
}

// Decode attempts to decode a json-encoded byte buffer into an
// object instance pointed to by objptr
func Decode(b []byte, objptr interface{}) error {
	dec := codec.NewDecoderBytes(b, CodecHandle)
	err := dec.Decode(objptr)
	if err != nil {
		return err
	}
	return nil
}

// LenientDecode attempts to decode a json-encoded byte buffer into an
// object instance pointed to by objptr
func LenientDecode(b []byte, objptr interface{}) error {
	dec := codec.NewDecoderBytes(b, LenientCodecHandle)
	err := dec.Decode(objptr)
	if err != nil {
		return err
	}
	return nil
}

// NewDecoder returns a json decoder
func NewDecoder(r io.Reader) *codec.Decoder {
	return codec.NewDecoder(r, CodecHandle)
}

// NewLenientDecoder returns a json decoder
func NewLenientDecoder(r io.Reader) *codec.Decoder {
	return codec.NewDecoder(r, LenientCodecHandle)
}
