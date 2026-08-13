//go:build falcon

package crypto

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/v2/mnemonic"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

// makeTestFalcon1024Account returns a deterministic Falcon-1024 account derived
// from a fixed mnemonic, so the tests are reproducible.
func makeTestFalcon1024Account(t *testing.T) Falcon1024Account {
	mn := "auction inquiry lava second expand liberty glass involve ginger illness length room item discover ahead table doctor term tackle cement bonus profit right above catch"
	seed, err := mnemonic.ToPQSeed(mn, types.PQSchemeFalcon1024)
	require.NoError(t, err)
	pqa, err := Falcon1024AccountFromPQSeed(seed)
	require.NoError(t, err)
	return pqa
}

func makeTestPaymentTxn(t *testing.T, sender types.Address) types.Transaction {
	toAddr, err := types.DecodeAddress("DN7MBMCL5JQ3PFUQS7TMX5AH4EEKOBJVDUF4TCV6WERATKFLQF4MQUPZTA")
	require.NoError(t, err)
	return types.Transaction{
		Type: types.PaymentTx,
		Header: types.Header{
			Sender:     sender,
			Fee:        217000,
			FirstValid: 972508,
			LastValid:  973508,
			Note:       []byte{180, 81, 121, 57, 252, 250, 210, 113},
			GenesisID:  "testnet-v31.0",
		},
		PaymentTxnFields: types.PaymentTxnFields{
			Receiver: toAddr,
			Amount:   5000,
		},
	}
}

func TestAddress(t *testing.T) {
	seed, err := base64.StdEncoding.DecodeString("EI+JCEv/+Kyqo5yvW6O2A/u0KKtLp5wWIjAvS5sT488=")
	require.NoError(t, err)
	account, err := Falcon1024AccountFromPQSeed(seed)
	require.NoError(t, err)
	expectedAddress, err := types.DecodeAddress("UGEDBJQD4LZF6OMFQDQ3BLY6CRX36Y75AZPDKJ3TTRU4TOGJ36EL34CWRI")
	require.NoError(t, err)
	require.Equal(t, expectedAddress, account.Address())
}

func TestGenerateFalcon1024Account(t *testing.T) {
	pqa := GenerateFalcon1024Account()
	require.NoError(t, pqa.Validate())
	require.NotEqual(t, types.Address{}, pqa.Address())
}

func TestFalcon1024AccountFromPQSeed(t *testing.T) {
	// Same seed must always yield the same account (deterministic keygen + salt).
	pqa1 := makeTestFalcon1024Account(t)
	pqa2 := makeTestFalcon1024Account(t)
	require.Equal(t, pqa1, pqa2)
	require.NoError(t, pqa1.Validate())

	// A valid PQ address must not double as a valid ed25519 point.
	addr := pqa1.Address()
	require.False(t, IsEdwards25519Point(addr[:]))
}

func TestSignFalcon1024AccountTransaction(t *testing.T) {
	pqa := makeTestFalcon1024Account(t)
	fromAddr := pqa.Address()
	tx := makeTestPaymentTxn(t, fromAddr)

	txid, txBytes, err := SignFalcon1024AccountTransaction(pqa.AsSigner(), tx)
	require.NoError(t, err)
	require.NotEmpty(t, txid)

	var stx types.SignedTxn
	require.NoError(t, msgpack.Decode(txBytes, &stx))

	// Sender == signer, so no AuthAddr is set.
	require.Equal(t, types.Address{}, stx.AuthAddr)
	require.Equal(t, tx, stx.Txn)
	require.Equal(t, types.PQSchemeFalcon1024, stx.PQsig.Scheme)

	bytesToSign := rawTransactionBytesToSign(stx.Txn)
	require.True(t, VerifyPQSig(bytesToSign, stx.PQsig))

	// A tampered transaction must not verify against the signature.
	stx.Txn.Amount++
	require.False(t, VerifyPQSig(rawTransactionBytesToSign(stx.Txn), stx.PQsig))
}

type customFalconSigner struct {
	pqa Falcon1024Account
}

// Falcon1024 signs the given bytes with a falcon1024 signature
func (sgnr customFalconSigner) Falcon1024Sign(toBeSigned []byte) ([]byte, error) {
	return nil, fmt.Errorf("Unimplemented")
}

// Falcon1024PublicKey returns the public key that should be used to verify the
// signatures performed by this signer
func (sgnr customFalconSigner) Falcon1024PublicKey() Falcon1024PublicKey {
	return sgnr.pqa.PublicKey
}

func TestBasicSignerGetsCanonicalSalt(t *testing.T) {
	pqa := makeTestFalcon1024Account(t)

	sgnr := customFalconSigner{pqa: pqa}
	salt, err := SaltForFalcon1024Signer(sgnr)
	require.NoError(t, err)

	defaultSgnr := pqa.AsSigner()
	defaultSalt, err := SaltForFalcon1024Signer(defaultSgnr)
	require.NoError(t, err)

	require.Equal(t, defaultSalt, salt)
}

func TestSaltedSignerOnlyDiffersInSaltAndAddress(t *testing.T) {
	pqa := makeTestFalcon1024Account(t)
	defaultSgnr := pqa.AsSigner()
	saltedSgnr := SaltedFalcon1024Signer{
		Signer: defaultSgnr,
		Salt:   types.PQAddressSalt(99),
	}
	fromAddr := pqa.Address()
	tx := makeTestPaymentTxn(t, fromAddr)

	_, txBytes, err := SignFalcon1024AccountTransaction(saltedSgnr, tx)
	require.NoError(t, err)

	var stx types.SignedTxn
	require.NoError(t, msgpack.Decode(txBytes, &stx))

	// We modified the salt, this means a different account made the signature
	// and therefore AuthAddr is set.
	require.NotEqual(t, types.Address{}, stx.AuthAddr)
	require.Equal(t, types.PQAddressSalt(99), stx.PQsig.Salt)

	bytesToSign := rawTransactionBytesToSign(stx.Txn)
	require.True(t, VerifyPQSig(bytesToSign, stx.PQsig))
}

func TestSignFalcon1024AccountTransactionWithAuthAddr(t *testing.T) {
	pqa := makeTestFalcon1024Account(t)
	authAddr := pqa.Address()

	// Sender differs from the signer: the account has been rekeyed to the PQ key.
	fromAddr, err := types.DecodeAddress("DN7MBMCL5JQ3PFUQS7TMX5AH4EEKOBJVDUF4TCV6WERATKFLQF4MQUPZTA")
	require.NoError(t, err)
	tx := makeTestPaymentTxn(t, fromAddr)

	_, txBytes, err := SignFalcon1024AccountTransaction(pqa.AsSigner(), tx)
	require.NoError(t, err)

	var stx types.SignedTxn
	require.NoError(t, msgpack.Decode(txBytes, &stx))
	require.Equal(t, authAddr, stx.AuthAddr)
}

func TestMakeLogicSigAccountDelegatedFalcon1024(t *testing.T) {
	pqa := makeTestFalcon1024Account(t)
	program := []byte{1, 32, 1, 1, 34}
	args := [][]byte{{0x01}, {0x02, 0x03}}

	lsa, err := MakeLogicSigAccountDelegatedFalcon1024(program, args, pqa.AsSigner())
	require.NoError(t, err)
	require.True(t, lsa.IsDelegated())
	require.False(t, lsa.Lsig.PQsig.Blank())

	// A delegated PQ lsig's address is the delegating PQ account.
	addr, err := lsa.Address()
	require.NoError(t, err)
	require.Equal(t, pqa.Address(), addr)

	require.True(t, VerifyLogicSig(lsa.Lsig, addr))

	// Tampering with the program must break verification.
	tampered := lsa.Lsig
	tampered.Logic = append([]byte{}, program...)
	tampered.Logic[3] = 2
	require.False(t, VerifyLogicSig(tampered, addr))
}
