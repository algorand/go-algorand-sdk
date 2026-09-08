package transaction

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/mnemonic"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

func TestMakeEd25519AccountTransactionSigner(t *testing.T) {
	account := crypto.GenerateAccount()
	txSigner := Ed25519AccountTransactionSigner{Signer: account.AsSigner()}

	addr, err := types.DecodeAddress("DN7MBMCL5JQ3PFUQS7TMX5AH4EEKOBJVDUF4TCV6WERATKFLQF4MQUPZTA")
	require.NoError(t, err)
	tx := types.Transaction{
		Type: types.PaymentTx,
		Header: types.Header{
			Sender:     addr,
			Fee:        217000,
			FirstValid: 972508,
			LastValid:  973508,
			Note:       []byte{180, 81, 121, 57, 252, 250, 210, 113},
			GenesisID:  "testnet-v31.0",
		},
		PaymentTxnFields: types.PaymentTxnFields{
			Receiver: addr,
			Amount:   5000,
		},
	}

	sigs, err := txSigner.SignTransactions([]types.Transaction{tx}, []int{0})
	require.NoError(t, err)

	expectedSig, err := ed25519SignTransaction(account.AsSigner(), tx)
	require.NoError(t, err)
	require.Len(t, sigs, 1)
	require.Equal(t, sigs[0], expectedSig)
}

func TestMakeLogicSigAccountTransactionSigner(t *testing.T) {
	program := []byte{1, 32, 1, 1, 34}
	args := [][]byte{
		{0x01},
		{0x02, 0x03},
	}
	account := crypto.GenerateAccount()
	lsig, err := crypto.Ed25519MakeLogicSigAccountDelegated(program, args, account.AsSigner())
	require.NoError(t, err)

	programHash := "6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY"
	programAddr, err := types.DecodeAddress(programHash)
	require.NoError(t, err)

	txSigner := LogicSigAccountTransactionSigner{LogicSigAccount: lsig}

	require.NoError(t, err)
	tx := types.Transaction{
		Type: types.PaymentTx,
		Header: types.Header{
			Sender:     programAddr,
			Fee:        217000,
			FirstValid: 972508,
			LastValid:  973508,
			Note:       []byte{180, 81, 121, 57, 252, 250, 210, 113},
			GenesisID:  "testnet-v31.0",
		},
		PaymentTxnFields: types.PaymentTxnFields{
			Receiver: programAddr,
			Amount:   5000,
		},
	}

	sigs, err := txSigner.SignTransactions([]types.Transaction{tx}, []int{0})
	require.NoError(t, err)

	_, expectedSig, err := crypto.SignLogicSigAccountTransaction(lsig, tx)
	require.NoError(t, err)
	require.Equal(t, sigs[0], expectedSig)
}

func makeTestMultisigAccount(t *testing.T) (crypto.MultisigAccount, crypto.Ed25519Signer, crypto.Ed25519Signer, crypto.Ed25519Signer) {
	addr1, err := types.DecodeAddress("DN7MBMCL5JQ3PFUQS7TMX5AH4EEKOBJVDUF4TCV6WERATKFLQF4MQUPZTA")
	require.NoError(t, err)
	addr2, err := types.DecodeAddress("BFRTECKTOOE7A5LHCF3TTEOH2A7BW46IYT2SX5VP6ANKEXHZYJY77SJTVM")
	require.NoError(t, err)
	addr3, err := types.DecodeAddress("47YPQTIGQEO7T4Y4RWDYWEKV6RTR2UNBQXBABEEGM72ESWDQNCQ52OPASU")
	require.NoError(t, err)
	ma, err := crypto.MultisigAccountWithParams(1, 2, []types.Address{
		addr1,
		addr2,
		addr3,
	})
	require.NoError(t, err)
	mn1 := "auction inquiry lava second expand liberty glass involve ginger illness length room item discover ahead table doctor term tackle cement bonus profit right above catch"
	sk1, err := mnemonic.ToPrivateKey(mn1)
	require.NoError(t, err)
	sgnr1, err := crypto.SKToInMemorySigner(sk1)
	require.NoError(t, err)
	mn2 := "since during average anxiety protect cherry club long lawsuit loan expand embark forum theory winter park twenty ball kangaroo cram burst board host ability left"
	sk2, err := mnemonic.ToPrivateKey(mn2)
	require.NoError(t, err)
	sgnr2, err := crypto.SKToInMemorySigner(sk2)
	require.NoError(t, err)
	mn3 := "advice pudding treat near rule blouse same whisper inner electric quit surface sunny dismiss leader blood seat clown cost exist hospital century reform able sponsor"
	sk3, err := mnemonic.ToPrivateKey(mn3)
	sgnr3, err := crypto.SKToInMemorySigner(sk3)
	require.NoError(t, err)
	return ma, sgnr1, sgnr2, sgnr3
}

func TestMakeMultiSigEd25519AccountTransactionSigner(t *testing.T) {
	ma, sgnr1, _, _ := makeTestMultisigAccount(t)
	fromAddr, err := ma.Address()
	require.NoError(t, err)
	toAddr, err := types.DecodeAddress("DN7MBMCL5JQ3PFUQS7TMX5AH4EEKOBJVDUF4TCV6WERATKFLQF4MQUPZTA")
	require.NoError(t, err)

	txSigner := MultiSigEd25519AccountTransactionSigner{Msig: ma, Signers: []crypto.Ed25519Signer{sgnr1}}
	tx := types.Transaction{
		Type: types.PaymentTx,
		Header: types.Header{
			Sender:     fromAddr,
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

	sigs, err := txSigner.SignTransactions([]types.Transaction{tx}, []int{0})
	require.NoError(t, err)

	expectedSig, err := ed25519SignMultisigTransaction(sgnr1, ma, tx)
	require.NoError(t, err)
	require.Equal(t, sigs[0], expectedSig)
}

func TestMultiSigEd25519AccountTransactionSignerEmptySigners(t *testing.T) {
	ma, _, _, _ := makeTestMultisigAccount(t)
	txSigner := MultiSigEd25519AccountTransactionSigner{Msig: ma, Signers: nil}

	tx := types.Transaction{}
	_, err := txSigner.SignTransactions([]types.Transaction{tx}, []int{0})
	require.Error(t, err)
	require.Contains(t, err.Error(), "multisig signer has no signing keys")

	_, err = txSigner.SignDelegationTo([]byte{1, 2, 3}, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "multisig signer has no signing keys")
}

type mockBasicPQSigner struct {
	scheme    types.PQScheme
	publicKey []byte
}

func (m mockBasicPQSigner) PQSign(toBeSigned []byte) ([]byte, error) {
	return []byte("sig"), nil
}

func (m mockBasicPQSigner) PQPublicKey() []byte {
	return m.publicKey
}

func (m mockBasicPQSigner) PQScheme() types.PQScheme {
	return m.scheme
}

type mockSaltedPQSigner struct {
	mockBasicPQSigner
	salt types.PQAddressSalt
}

func (m mockSaltedPQSigner) PQSalt() types.PQAddressSalt {
	return m.salt
}

func TestPQAccountTransactionSignerEquals(t *testing.T) {
	pk := []byte("12345678901234567890123456789012")
	scheme1 := types.PQScheme{'f', '1'}
	scheme2 := types.PQScheme{'m', '2'}

	s1 := PQAccountTransactionSigner{Signer: mockSaltedPQSigner{mockBasicPQSigner: mockBasicPQSigner{scheme: scheme1, publicKey: pk}, salt: 0}}
	s1Same := PQAccountTransactionSigner{Signer: mockSaltedPQSigner{mockBasicPQSigner: mockBasicPQSigner{scheme: scheme1, publicKey: pk}, salt: 0}}
	sDiffScheme := PQAccountTransactionSigner{Signer: mockSaltedPQSigner{mockBasicPQSigner: mockBasicPQSigner{scheme: scheme2, publicKey: pk}, salt: 0}}
	sDiffPK := PQAccountTransactionSigner{Signer: mockSaltedPQSigner{mockBasicPQSigner: mockBasicPQSigner{scheme: scheme1, publicKey: []byte("other-pk-1234567890123456789012")}, salt: 0}}
	sDiffSalt := PQAccountTransactionSigner{Signer: mockSaltedPQSigner{mockBasicPQSigner: mockBasicPQSigner{scheme: scheme1, publicKey: pk}, salt: 1}}

	require.True(t, s1.Equals(s1Same))
	require.False(t, s1.Equals(sDiffScheme))
	require.False(t, s1.Equals(sDiffPK))
	require.False(t, s1.Equals(sDiffSalt))
	require.False(t, s1.Equals(EmptyTransactionSigner{}))
	require.False(t, s1.Equals(PQAccountTransactionSigner{Signer: nil}))
	require.True(t, (PQAccountTransactionSigner{Signer: nil}).Equals(PQAccountTransactionSigner{Signer: nil}))
}

func TestEd25519AccountTransactionSignerEqualsNilSigner(t *testing.T) {
	account := crypto.GenerateAccount()
	s1 := Ed25519AccountTransactionSigner{Signer: account.AsSigner()}

	require.True(t, (Ed25519AccountTransactionSigner{Signer: nil}).Equals(Ed25519AccountTransactionSigner{Signer: nil}))
	require.False(t, s1.Equals(Ed25519AccountTransactionSigner{Signer: nil}))
	require.False(t, (Ed25519AccountTransactionSigner{Signer: nil}).Equals(s1))
}

func TestMultiSigEd25519AccountTransactionSignerEqualsNilSigners(t *testing.T) {
	ma, sgnr1, _, _ := makeTestMultisigAccount(t)
	s1 := MultiSigEd25519AccountTransactionSigner{Msig: ma, Signers: []crypto.Ed25519Signer{sgnr1, nil}}
	s1Same := MultiSigEd25519AccountTransactionSigner{Msig: ma, Signers: []crypto.Ed25519Signer{sgnr1, nil}}
	s1Diff := MultiSigEd25519AccountTransactionSigner{Msig: ma, Signers: []crypto.Ed25519Signer{sgnr1, sgnr1}}

	require.True(t, s1.Equals(s1Same))
	require.False(t, s1.Equals(s1Diff))
}

type mockEd25519SignerWithLen struct {
	sigLen int
}

func (m mockEd25519SignerWithLen) Ed25519Sign(message []byte) ([]byte, error) {
	return make([]byte, m.sigLen), nil
}

func (m mockEd25519SignerWithLen) Ed25519PublicKey() crypto.Ed25519PublicKey {
	return crypto.Ed25519PublicKey{}
}

func TestEd25519SignatureLengthValidation(t *testing.T) {
	tx := types.Transaction{}

	// Oversized signature (65 bytes) must be rejected
	signerOversized := Ed25519AccountTransactionSigner{Signer: mockEd25519SignerWithLen{sigLen: 65}}
	_, _, err := SignTransaction(signerOversized, tx)
	require.ErrorIs(t, err, errInvalidSignatureReturned)

	// Undersized signature (63 bytes) must be rejected
	signerUndersized := Ed25519AccountTransactionSigner{Signer: mockEd25519SignerWithLen{sigLen: 63}}
	_, _, err = SignTransaction(signerUndersized, tx)
	require.ErrorIs(t, err, errInvalidSignatureReturned)
}
