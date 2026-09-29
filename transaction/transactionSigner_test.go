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

	_, expectedSig, err := crypto.SignTransaction(account.PrivateKey, tx)
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
	lsig, err := Ed25519AccountTransactionSigner{Signer: account.AsSigner()}.SignDelegationTo(program, args)
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

	// same transaction and keys as crypto's TestSignMultisigTransaction
	expectedBytes := []byte{130, 164, 109, 115, 105, 103, 131, 166, 115, 117, 98, 115, 105, 103, 147, 130, 162, 112, 107, 196, 32, 27, 126, 192, 176, 75, 234, 97, 183, 150, 144, 151, 230, 203, 244, 7, 225, 8, 167, 5, 53, 29, 11, 201, 138, 190, 177, 34, 9, 168, 171, 129, 120, 161, 115, 196, 64, 118, 246, 119, 203, 209, 172, 34, 112, 79, 186, 215, 112, 41, 206, 201, 203, 230, 167, 215, 112, 156, 141, 37, 117, 149, 203, 209, 1, 132, 10, 96, 236, 87, 193, 248, 19, 228, 31, 230, 43, 94, 17, 231, 187, 158, 96, 148, 216, 202, 128, 206, 243, 48, 88, 234, 68, 38, 5, 169, 86, 146, 111, 121, 0, 129, 162, 112, 107, 196, 32, 9, 99, 50, 9, 83, 115, 137, 240, 117, 103, 17, 119, 57, 145, 199, 208, 62, 27, 115, 200, 196, 245, 43, 246, 175, 240, 26, 162, 92, 249, 194, 113, 129, 162, 112, 107, 196, 32, 231, 240, 248, 77, 6, 129, 29, 249, 243, 28, 141, 135, 139, 17, 85, 244, 103, 29, 81, 161, 133, 194, 0, 144, 134, 103, 244, 73, 88, 112, 104, 161, 163, 116, 104, 114, 2, 161, 118, 1, 163, 116, 120, 110, 137, 163, 97, 109, 116, 205, 19, 136, 163, 102, 101, 101, 206, 0, 3, 79, 168, 162, 102, 118, 206, 0, 14, 214, 220, 163, 103, 101, 110, 173, 116, 101, 115, 116, 110, 101, 116, 45, 118, 51, 49, 46, 48, 162, 108, 118, 206, 0, 14, 218, 196, 164, 110, 111, 116, 101, 196, 8, 180, 81, 121, 57, 252, 250, 210, 113, 163, 114, 99, 118, 196, 32, 27, 126, 192, 176, 75, 234, 97, 183, 150, 144, 151, 230, 203, 244, 7, 225, 8, 167, 5, 53, 29, 11, 201, 138, 190, 177, 34, 9, 168, 171, 129, 120, 163, 115, 110, 100, 196, 32, 141, 146, 180, 137, 144, 1, 115, 160, 77, 250, 67, 89, 163, 102, 106, 106, 252, 234, 44, 66, 160, 93, 217, 193, 247, 62, 235, 165, 71, 128, 55, 233, 164, 116, 121, 112, 101, 163, 112, 97, 121}
	require.Equal(t, expectedBytes, sigs[0])
}

func TestMultiSigEd25519AccountTransactionSignerEmptySigners(t *testing.T) {
	ma, _, _, _ := makeTestMultisigAccount(t)
	txSigner := MultiSigEd25519AccountTransactionSigner{Msig: ma, Signers: nil}

	// A signer with no keys is misconfigured, so every entry point rejects it
	// eagerly rather than only once there is a transaction to sign.
	_, err := txSigner.SignTransactions(nil, nil)
	require.ErrorIs(t, err, errNoMultisigSigners)

	tx := types.Transaction{}
	_, err = txSigner.SignTransactions([]types.Transaction{tx}, []int{0})
	require.ErrorIs(t, err, errNoMultisigSigners)

	_, err = txSigner.SignDelegationTo([]byte{1, 2, 3}, nil)
	require.ErrorIs(t, err, errNoMultisigSigners)
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

func TestPQAccountTransactionSignerEquals(t *testing.T) {
	pk := []byte("12345678901234567890123456789012")
	scheme1 := types.PQScheme{'f', '1'}
	scheme2 := types.PQScheme{'m', '2'}
	signer := &mockBasicPQSigner{scheme: scheme1, publicKey: pk}

	s1 := PQAccountTransactionSigner{Signer: signer}
	s1Same := PQAccountTransactionSigner{Signer: signer}
	sSameAccount := PQAccountTransactionSigner{Signer: &mockBasicPQSigner{scheme: scheme1, publicKey: pk}}
	sDiffScheme := PQAccountTransactionSigner{Signer: &mockBasicPQSigner{scheme: scheme2, publicKey: pk}}
	sDiffPK := PQAccountTransactionSigner{Signer: &mockBasicPQSigner{scheme: scheme1, publicKey: []byte("other-pk-1234567890123456789012")}}

	require.True(t, s1.Equals(s1Same))
	require.False(t, s1.Equals(sSameAccount))
	require.False(t, s1.Equals(sDiffScheme))
	require.False(t, s1.Equals(sDiffPK))
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
	require.ErrorIs(t, err, crypto.ErrInvalidSignatureReturned)

	// Undersized signature (63 bytes) must be rejected
	signerUndersized := Ed25519AccountTransactionSigner{Signer: mockEd25519SignerWithLen{sigLen: 63}}
	_, _, err = SignTransaction(signerUndersized, tx)
	require.ErrorIs(t, err, crypto.ErrInvalidSignatureReturned)
}
