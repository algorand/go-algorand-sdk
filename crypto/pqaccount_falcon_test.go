//go:build falcon

package crypto

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

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
	accountAddr, err := account.Address()
	require.NoError(t, err)
	require.Equal(t, expectedAddress, accountAddr)
}

func TestGenerateFalcon1024Account(t *testing.T) {
	pqa := GenerateFalcon1024Account()
	addr, err := pqa.Address()
	require.NoError(t, err)
	require.NotEqual(t, types.Address{}, addr)
}

func TestFalcon1024AccountFromPQSeed(t *testing.T) {
	// Same seed must always yield the same account (deterministic keygen).
	pqa1 := makeTestFalcon1024Account(t)
	pqa2 := makeTestFalcon1024Account(t)
	require.Equal(t, pqa1, pqa2)

	// A valid PQ address must not double as a valid ed25519 point.
	addr, err := pqa1.Address()
	require.NoError(t, err)
	require.False(t, IsEdwards25519Point(addr[:]))
}

func TestSignFalcon1024AccountSigner(t *testing.T) {
	pqa := makeTestFalcon1024Account(t)
	sgnr := pqa.AsSigner()

	// The signer signs on behalf of the account it was derived from.
	addr, err := PQSignerAddress(sgnr)
	require.NoError(t, err)
	pqaAddr, err := pqa.Address()
	require.NoError(t, err)
	require.Equal(t, pqaAddr, addr)
	require.Equal(t, types.PQSchemeFalcon1024, sgnr.PQScheme())

	toBeSigned := TransactionBytesToSign(makeTestPaymentTxn(t, addr))
	pqsig, err := signWith(sgnr, toBeSigned)
	require.NoError(t, err)
	require.True(t, VerifyPQSig(toBeSigned, pqsig))

	// A tampered message must not verify against the signature.
	require.False(t, VerifyPQSig(append(toBeSigned, 0), pqsig))
}

// signWith signs the given bytes and returns the resulting PQSig envelope
func signWith(sgnr PQSigner, toBeSigned []byte) (types.PQSig, error) {
	signature, err := sgnr.PQSign(toBeSigned)
	if err != nil {
		return types.PQSig{}, err
	}
	pqsig, _, err := PQSigFor(sgnr, signature)
	return pqsig, err
}

type customFalconSigner struct {
	pqa Falcon1024Account
}

// PQSign signs the given bytes with a pq signature
func (sgnr customFalconSigner) PQSign(toBeSigned []byte) ([]byte, error) {
	return nil, fmt.Errorf("Unimplemented")
}

// PQPublicKey returns the public key that should be used to verify the
// signatures performed by this signer
func (sgnr customFalconSigner) PQPublicKey() []byte {
	return sgnr.pqa.PublicKey[:]
}

// PQScheme returns the identifier for the post-quantum scheme used by this
// signer
func (sgnr customFalconSigner) PQScheme() types.PQScheme {
	return types.PQSchemeFalcon1024
}

// Any signer for the same key material selects the same account: the salt is
// derived from the (scheme, public key) pair, not carried by the signer.
func TestEverySignerGetsCanonicalSalt(t *testing.T) {
	pqa := makeTestFalcon1024Account(t)

	sgnr := customFalconSigner{pqa: pqa}
	salt, err := SaltForPQSigner(sgnr)
	require.NoError(t, err)

	defaultSgnr := pqa.AsSigner()
	defaultSalt, err := SaltForPQSigner(defaultSgnr)
	require.NoError(t, err)

	require.Equal(t, defaultSalt, salt)

	// That salt is the canonical one, and the addresses agree on it.
	canonical, err := canonicalSaltForPQPK(pqa.PublicKey[:], types.PQSchemeFalcon1024)
	require.NoError(t, err)
	require.Equal(t, canonical, salt)

	addr, err := PQSignerAddress(sgnr)
	require.NoError(t, err)
	defaultAddr, err := PQSignerAddress(defaultSgnr)
	require.NoError(t, err)
	require.Equal(t, defaultAddr, addr)
	require.Equal(t, pqAddressWithSalt(pqa.PublicKey[:], types.PQSchemeFalcon1024, canonical), addr)

	// Signatures made by that signer carry the canonical salt too.
	toBeSigned := TransactionBytesToSign(makeTestPaymentTxn(t, addr))
	pqsig, err := signWith(defaultSgnr, toBeSigned)
	require.NoError(t, err)
	require.Equal(t, canonical, pqsig.Salt)
	require.Equal(t, addr, PQAddressFromSig(pqsig))
	require.True(t, VerifyPQSig(toBeSigned, pqsig))
}

func TestMakeLogicSigAccountDelegatedFalcon1024(t *testing.T) {
	pqa := makeTestFalcon1024Account(t)
	program := []byte{1, 32, 1, 1, 34}
	args := [][]byte{{0x01}, {0x02, 0x03}}

	lsa, err := MakeLogicSigAccountDelegatedPQ(program, args, pqa.AsSigner())
	require.NoError(t, err)
	require.True(t, lsa.IsDelegated())
	require.False(t, lsa.Lsig.PQsig.Blank())

	// A delegated PQ lsig's address is the delegating PQ account.
	addr, err := lsa.Address()
	require.NoError(t, err)
	pqaAddr, err := pqa.Address()
	require.NoError(t, err)
	require.Equal(t, pqaAddr, addr)

	toBeSigned := pqsigProgramToSign(addr, lsa.Lsig.Logic)
	require.True(t, VerifyPQSig(toBeSigned, lsa.Lsig.PQsig))

	// Tampering with the program must break verification.
	tampered := lsa.Lsig
	tampered.Logic = append([]byte{}, program...)
	tampered.Logic[3] = 2
	tamperedToBeSigned := pqsigProgramToSign(addr, tampered.Logic)
	require.False(t, VerifyPQSig(tamperedToBeSigned, tampered.PQsig))

	// VerifyLogicSig checks that the delegating singleSigner matches the PQ signature
	require.True(t, VerifyLogicSig(lsa.Lsig, addr))
	wrongAddr := types.Address{1, 2, 3}
	require.False(t, VerifyLogicSig(lsa.Lsig, wrongAddr))
	require.False(t, VerifyLogicSig(lsa.Lsig, types.Address{}))

	// VerifyPQSig rejects a mismatched scheme, even though the signature is
	// otherwise valid for toBeSigned.
	wrongSchemeSig := lsa.Lsig.PQsig
	wrongSchemeSig.Scheme = types.PQScheme{'x', 'x'}
	require.False(t, VerifyPQSig(toBeSigned, wrongSchemeSig))

	// Likewise for a public key that is not falcon1024-sized: without the
	// length check, converting it to a falcon.PublicKey array would panic.
	wrongLenSig := lsa.Lsig.PQsig
	wrongLenSig.PublicKey = make([]byte, 32)
	require.False(t, VerifyPQSig(toBeSigned, wrongLenSig))
}

func TestPQAccountNilSignerChecks(t *testing.T) {
	_, err := SaltForPQSigner(nil)
	require.Error(t, err)

	_, err = PQSignerAddress(nil)
	require.Error(t, err)

	_, err = MakeLogicSigAccountDelegatedPQ([]byte{1, 2, 3}, nil, nil)
	require.Error(t, err)
}
