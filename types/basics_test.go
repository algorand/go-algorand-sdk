package types

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
)

func TestSignedTxnFromBase64String(t *testing.T) {
	note := []byte("Testing SignedTxnFromBase64String()")

	// Create a signed transaction with a note and base64 encode it.
	var st SignedTxn
	st.Txn.Note = note
	b64data := base64.StdEncoding.EncodeToString([]byte(string(msgpack.Encode(st))))

	// Verify that FromBase64String() decodes the txn.
	var vst SignedTxn
	err := vst.FromBase64String(b64data)
	require.NoError(t, err)
	require.Equal(t, vst.Txn.Note, note)
}

func TestBlockFromBase64String(t *testing.T) {
	protocol := "Testing BlockFromBase64String()"

	// Create a block with a protocol string and base64 encode it.
	var bl Block
	bl.CurrentProtocol = protocol
	b64data := base64.StdEncoding.EncodeToString([]byte(string(msgpack.Encode(bl))))

	// Verify that FromBase64String() decodes the txn.
	var vbl Block
	err := vbl.FromBase64String(b64data)
	require.NoError(t, err)
	require.Equal(t, vbl.CurrentProtocol, protocol)
}

func TestLogicSigBlank(t *testing.T) {
	require.True(t, (LogicSig{}).Blank())

	require.False(t, (LogicSig{Args: [][]byte{{1}}}).Blank())
	require.False(t, (LogicSig{Logic: []byte{1}}).Blank())
	require.False(t, (LogicSig{Sig: Signature{1}}).Blank())
	require.False(t, (LogicSig{Msig: MultisigSig{Version: 1}}).Blank())
	require.False(t, (LogicSig{LMsig: MultisigSig{Version: 1}}).Blank())
	require.False(t, (LogicSig{PQsig: PQSig{Scheme: PQSchemeFalcon1024}}).Blank())
	require.False(t, (LogicSig{PQsig: PQSig{Signature: []byte{1}}}).Blank())
}
