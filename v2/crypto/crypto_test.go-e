package crypto

import (
	"bytes"
	"encoding/base64"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ed25519"

	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/algorand/go-algorand-sdk/types"
)

func makeTestMultisigAccount(t *testing.T) (MultisigAccount, ed25519.PrivateKey, ed25519.PrivateKey, ed25519.PrivateKey) {
	addr1, err := types.DecodeAddress("DN7MBMCL5JQ3PFUQS7TMX5AH4EEKOBJVDUF4TCV6WERATKFLQF4MQUPZTA")
	require.NoError(t, err)
	addr2, err := types.DecodeAddress("BFRTECKTOOE7A5LHCF3TTEOH2A7BW46IYT2SX5VP6ANKEXHZYJY77SJTVM")
	require.NoError(t, err)
	addr3, err := types.DecodeAddress("47YPQTIGQEO7T4Y4RWDYWEKV6RTR2UNBQXBABEEGM72ESWDQNCQ52OPASU")
	require.NoError(t, err)
	ma, err := MultisigAccountWithParams(1, 2, []types.Address{
		addr1,
		addr2,
		addr3,
	})
	require.NoError(t, err)
	mn1 := "auction inquiry lava second expand liberty glass involve ginger illness length room item discover ahead table doctor term tackle cement bonus profit right above catch"
	sk1, err := mnemonic.ToPrivateKey(mn1)
	require.NoError(t, err)
	mn2 := "since during average anxiety protect cherry club long lawsuit loan expand embark forum theory winter park twenty ball kangaroo cram burst board host ability left"
	sk2, err := mnemonic.ToPrivateKey(mn2)
	require.NoError(t, err)
	mn3 := "advice pudding treat near rule blouse same whisper inner electric quit surface sunny dismiss leader blood seat clown cost exist hospital century reform able sponsor"
	sk3, err := mnemonic.ToPrivateKey(mn3)
	return ma, sk1, sk2, sk3
}

func TestSignMultisigTransaction(t *testing.T) {
	ma, sk1, _, _ := makeTestMultisigAccount(t)
	fromAddr, err := ma.Address()
	require.NoError(t, err)
	toAddr, err := types.DecodeAddress("DN7MBMCL5JQ3PFUQS7TMX5AH4EEKOBJVDUF4TCV6WERATKFLQF4MQUPZTA")
	require.NoError(t, err)
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
	txid, txBytes, err := SignMultisigTransaction(sk1, ma, tx)
	require.NoError(t, err)
	expectedBytes := []byte{130, 164, 109, 115, 105, 103, 131, 166, 115, 117, 98, 115, 105, 103, 147, 130, 162, 112, 107, 196, 32, 27, 126, 192, 176, 75, 234, 97, 183, 150, 144, 151, 230, 203, 244, 7, 225, 8, 167, 5, 53, 29, 11, 201, 138, 190, 177, 34, 9, 168, 171, 129, 120, 161, 115, 196, 64, 118, 246, 119, 203, 209, 172, 34, 112, 79, 186, 215, 112, 41, 206, 201, 203, 230, 167, 215, 112, 156, 141, 37, 117, 149, 203, 209, 1, 132, 10, 96, 236, 87, 193, 248, 19, 228, 31, 230, 43, 94, 17, 231, 187, 158, 96, 148, 216, 202, 128, 206, 243, 48, 88, 234, 68, 38, 5, 169, 86, 146, 111, 121, 0, 129, 162, 112, 107, 196, 32, 9, 99, 50, 9, 83, 115, 137, 240, 117, 103, 17, 119, 57, 145, 199, 208, 62, 27, 115, 200, 196, 245, 43, 246, 175, 240, 26, 162, 92, 249, 194, 113, 129, 162, 112, 107, 196, 32, 231, 240, 248, 77, 6, 129, 29, 249, 243, 28, 141, 135, 139, 17, 85, 244, 103, 29, 81, 161, 133, 194, 0, 144, 134, 103, 244, 73, 88, 112, 104, 161, 163, 116, 104, 114, 2, 161, 118, 1, 163, 116, 120, 110, 137, 163, 97, 109, 116, 205, 19, 136, 163, 102, 101, 101, 206, 0, 3, 79, 168, 162, 102, 118, 206, 0, 14, 214, 220, 163, 103, 101, 110, 173, 116, 101, 115, 116, 110, 101, 116, 45, 118, 51, 49, 46, 48, 162, 108, 118, 206, 0, 14, 218, 196, 164, 110, 111, 116, 101, 196, 8, 180, 81, 121, 57, 252, 250, 210, 113, 163, 114, 99, 118, 196, 32, 27, 126, 192, 176, 75, 234, 97, 183, 150, 144, 151, 230, 203, 244, 7, 225, 8, 167, 5, 53, 29, 11, 201, 138, 190, 177, 34, 9, 168, 171, 129, 120, 163, 115, 110, 100, 196, 32, 141, 146, 180, 137, 144, 1, 115, 160, 77, 250, 67, 89, 163, 102, 106, 106, 252, 234, 44, 66, 160, 93, 217, 193, 247, 62, 235, 165, 71, 128, 55, 233, 164, 116, 121, 112, 101, 163, 112, 97, 121}
	require.EqualValues(t, expectedBytes, txBytes)
	expectedTxid := "KY6I7NQXQDAMDUCAVATI45BAODW6NRYQKFH4KIHLH2HQI4DO4XBA"
	require.Equal(t, expectedTxid, txid)

	// decode and verify
	var stx types.SignedTxn
	err = msgpack.Decode(txBytes, &stx)
	require.NoError(t, err)
	bytesToSign := rawTransactionBytesToSign(stx.Txn)

	verified := VerifyMultisig(fromAddr, bytesToSign, stx.Msig)
	require.False(t, verified) // not enough signatures
}

func TestAppendMultisigTransaction(t *testing.T) {
	uniSigTxBytes := []byte{130, 164, 109, 115, 105, 103, 131, 166, 115, 117, 98, 115, 105, 103, 147, 130, 162, 112, 107, 196, 32, 27, 126, 192, 176, 75, 234, 97, 183, 150, 144, 151, 230, 203, 244, 7, 225, 8, 167, 5, 53, 29, 11, 201, 138, 190, 177, 34, 9, 168, 171, 129, 120, 161, 115, 196, 64, 118, 246, 119, 203, 209, 172, 34, 112, 79, 186, 215, 112, 41, 206, 201, 203, 230, 167, 215, 112, 156, 141, 37, 117, 149, 203, 209, 1, 132, 10, 96, 236, 87, 193, 248, 19, 228, 31, 230, 43, 94, 17, 231, 187, 158, 96, 148, 216, 202, 128, 206, 243, 48, 88, 234, 68, 38, 5, 169, 86, 146, 111, 121, 0, 129, 162, 112, 107, 196, 32, 9, 99, 50, 9, 83, 115, 137, 240, 117, 103, 17, 119, 57, 145, 199, 208, 62, 27, 115, 200, 196, 245, 43, 246, 175, 240, 26, 162, 92, 249, 194, 113, 129, 162, 112, 107, 196, 32, 231, 240, 248, 77, 6, 129, 29, 249, 243, 28, 141, 135, 139, 17, 85, 244, 103, 29, 81, 161, 133, 194, 0, 144, 134, 103, 244, 73, 88, 112, 104, 161, 163, 116, 104, 114, 2, 161, 118, 1, 163, 116, 120, 110, 137, 163, 97, 109, 116, 205, 19, 136, 163, 102, 101, 101, 206, 0, 3, 79, 168, 162, 102, 118, 206, 0, 14, 214, 220, 163, 103, 101, 110, 173, 116, 101, 115, 116, 110, 101, 116, 45, 118, 51, 49, 46, 48, 162, 108, 118, 206, 0, 14, 218, 196, 164, 110, 111, 116, 101, 196, 8, 180, 81, 121, 57, 252, 250, 210, 113, 163, 114, 99, 118, 196, 32, 27, 126, 192, 176, 75, 234, 97, 183, 150, 144, 151, 230, 203, 244, 7, 225, 8, 167, 5, 53, 29, 11, 201, 138, 190, 177, 34, 9, 168, 171, 129, 120, 163, 115, 110, 100, 196, 32, 141, 146, 180, 137, 144, 1, 115, 160, 77, 250, 67, 89, 163, 102, 106, 106, 252, 234, 44, 66, 160, 93, 217, 193, 247, 62, 235, 165, 71, 128, 55, 233, 164, 116, 121, 112, 101, 163, 112, 97, 121}
	ma, _, sk2, _ := makeTestMultisigAccount(t)
	txid, txBytes, err := AppendMultisigTransaction(sk2, ma, uniSigTxBytes)
	require.NoError(t, err)
	expectedBytes := []byte{130, 164, 109, 115, 105, 103, 131, 166, 115, 117, 98, 115, 105, 103, 147, 130, 162, 112, 107, 196, 32, 27, 126, 192, 176, 75, 234, 97, 183, 150, 144, 151, 230, 203, 244, 7, 225, 8, 167, 5, 53, 29, 11, 201, 138, 190, 177, 34, 9, 168, 171, 129, 120, 161, 115, 196, 64, 118, 246, 119, 203, 209, 172, 34, 112, 79, 186, 215, 112, 41, 206, 201, 203, 230, 167, 215, 112, 156, 141, 37, 117, 149, 203, 209, 1, 132, 10, 96, 236, 87, 193, 248, 19, 228, 31, 230, 43, 94, 17, 231, 187, 158, 96, 148, 216, 202, 128, 206, 243, 48, 88, 234, 68, 38, 5, 169, 86, 146, 111, 121, 0, 130, 162, 112, 107, 196, 32, 9, 99, 50, 9, 83, 115, 137, 240, 117, 103, 17, 119, 57, 145, 199, 208, 62, 27, 115, 200, 196, 245, 43, 246, 175, 240, 26, 162, 92, 249, 194, 113, 161, 115, 196, 64, 78, 28, 117, 80, 233, 161, 90, 21, 86, 137, 23, 68, 108, 250, 59, 209, 183, 58, 57, 99, 119, 233, 29, 233, 221, 128, 106, 21, 203, 60, 96, 207, 181, 63, 208, 3, 208, 200, 214, 176, 120, 195, 199, 4, 13, 60, 187, 129, 222, 79, 105, 217, 39, 228, 70, 66, 218, 152, 243, 26, 29, 12, 112, 7, 129, 162, 112, 107, 196, 32, 231, 240, 248, 77, 6, 129, 29, 249, 243, 28, 141, 135, 139, 17, 85, 244, 103, 29, 81, 161, 133, 194, 0, 144, 134, 103, 244, 73, 88, 112, 104, 161, 163, 116, 104, 114, 2, 161, 118, 1, 163, 116, 120, 110, 137, 163, 97, 109, 116, 205, 19, 136, 163, 102, 101, 101, 206, 0, 3, 79, 168, 162, 102, 118, 206, 0, 14, 214, 220, 163, 103, 101, 110, 173, 116, 101, 115, 116, 110, 101, 116, 45, 118, 51, 49, 46, 48, 162, 108, 118, 206, 0, 14, 218, 196, 164, 110, 111, 116, 101, 196, 8, 180, 81, 121, 57, 252, 250, 210, 113, 163, 114, 99, 118, 196, 32, 27, 126, 192, 176, 75, 234, 97, 183, 150, 144, 151, 230, 203, 244, 7, 225, 8, 167, 5, 53, 29, 11, 201, 138, 190, 177, 34, 9, 168, 171, 129, 120, 163, 115, 110, 100, 196, 32, 141, 146, 180, 137, 144, 1, 115, 160, 77, 250, 67, 89, 163, 102, 106, 106, 252, 234, 44, 66, 160, 93, 217, 193, 247, 62, 235, 165, 71, 128, 55, 233, 164, 116, 121, 112, 101, 163, 112, 97, 121}
	require.EqualValues(t, expectedBytes, txBytes)
	expectedTxid := "KY6I7NQXQDAMDUCAVATI45BAODW6NRYQKFH4KIHLH2HQI4DO4XBA"
	require.Equal(t, expectedTxid, txid)

	// decode and verify
	var stx types.SignedTxn
	err = msgpack.Decode(txBytes, &stx)
	require.NoError(t, err)
	bytesToSign := rawTransactionBytesToSign(stx.Txn)

	fromAddr, err := ma.Address()
	require.NoError(t, err)

	verified := VerifyMultisig(fromAddr, bytesToSign, stx.Msig)
	require.True(t, verified)

	// now change fee and ensure signature verification fails
	idx := bytes.IndexAny(bytesToSign, "fee")
	require.NotEqual(t, -1, idx)
	bytesToSign[idx+4] = 1
	verified = VerifyMultisig(fromAddr, bytesToSign, stx.Msig)
	require.False(t, verified)
}

func TestSignMultisigTransactionKeyReg(t *testing.T) {
	rawKeyRegTx := []byte{129, 163, 116, 120, 110, 137, 163, 102, 101, 101, 206, 0, 3, 200, 192, 162, 102, 118, 206, 0, 14, 249, 218, 162, 108, 118, 206, 0, 14, 253, 194, 166, 115, 101, 108, 107, 101, 121, 196, 32, 50, 18, 43, 43, 214, 61, 220, 83, 49, 150, 23, 165, 170, 83, 196, 177, 194, 111, 227, 220, 202, 242, 141, 54, 34, 181, 105, 119, 161, 64, 92, 134, 163, 115, 110, 100, 196, 32, 141, 146, 180, 137, 144, 1, 115, 160, 77, 250, 67, 89, 163, 102, 106, 106, 252, 234, 44, 66, 160, 93, 217, 193, 247, 62, 235, 165, 71, 128, 55, 233, 164, 116, 121, 112, 101, 166, 107, 101, 121, 114, 101, 103, 166, 118, 111, 116, 101, 107, 100, 205, 39, 16, 167, 118, 111, 116, 101, 107, 101, 121, 196, 32, 112, 27, 215, 251, 145, 43, 7, 179, 8, 17, 255, 40, 29, 159, 238, 149, 99, 229, 128, 46, 32, 38, 137, 35, 25, 37, 143, 119, 250, 147, 30, 136, 167, 118, 111, 116, 101, 108, 115, 116, 206, 0, 15, 66, 64}
	var tx types.SignedTxn
	err := msgpack.Decode(rawKeyRegTx, &tx)
	require.NoError(t, err)

	ma, sk1, _, _ := makeTestMultisigAccount(t)
	txid, txBytes, err := SignMultisigTransaction(sk1, ma, tx.Txn)
	require.NoError(t, err)
	expectedBytes := []byte{130, 164, 109, 115, 105, 103, 131, 166, 115, 117, 98, 115, 105, 103, 147, 130, 162, 112, 107, 196, 32, 27, 126, 192, 176, 75, 234, 97, 183, 150, 144, 151, 230, 203, 244, 7, 225, 8, 167, 5, 53, 29, 11, 201, 138, 190, 177, 34, 9, 168, 171, 129, 120, 161, 115, 196, 64, 186, 52, 94, 163, 20, 123, 21, 228, 212, 78, 168, 14, 159, 234, 210, 219, 69, 206, 23, 113, 13, 3, 226, 107, 74, 6, 121, 202, 250, 195, 62, 13, 205, 64, 12, 208, 205, 69, 221, 116, 29, 15, 86, 243, 209, 159, 143, 116, 161, 84, 144, 104, 113, 8, 99, 78, 68, 12, 149, 213, 4, 83, 201, 15, 129, 162, 112, 107, 196, 32, 9, 99, 50, 9, 83, 115, 137, 240, 117, 103, 17, 119, 57, 145, 199, 208, 62, 27, 115, 200, 196, 245, 43, 246, 175, 240, 26, 162, 92, 249, 194, 113, 129, 162, 112, 107, 196, 32, 231, 240, 248, 77, 6, 129, 29, 249, 243, 28, 141, 135, 139, 17, 85, 244, 103, 29, 81, 161, 133, 194, 0, 144, 134, 103, 244, 73, 88, 112, 104, 161, 163, 116, 104, 114, 2, 161, 118, 1, 163, 116, 120, 110, 137, 163, 102, 101, 101, 206, 0, 3, 200, 192, 162, 102, 118, 206, 0, 14, 249, 218, 162, 108, 118, 206, 0, 14, 253, 194, 166, 115, 101, 108, 107, 101, 121, 196, 32, 50, 18, 43, 43, 214, 61, 220, 83, 49, 150, 23, 165, 170, 83, 196, 177, 194, 111, 227, 220, 202, 242, 141, 54, 34, 181, 105, 119, 161, 64, 92, 134, 163, 115, 110, 100, 196, 32, 141, 146, 180, 137, 144, 1, 115, 160, 77, 250, 67, 89, 163, 102, 106, 106, 252, 234, 44, 66, 160, 93, 217, 193, 247, 62, 235, 165, 71, 128, 55, 233, 164, 116, 121, 112, 101, 166, 107, 101, 121, 114, 101, 103, 166, 118, 111, 116, 101, 107, 100, 205, 39, 16, 167, 118, 111, 116, 101, 107, 101, 121, 196, 32, 112, 27, 215, 251, 145, 43, 7, 179, 8, 17, 255, 40, 29, 159, 238, 149, 99, 229, 128, 46, 32, 38, 137, 35, 25, 37, 143, 119, 250, 147, 30, 136, 167, 118, 111, 116, 101, 108, 115, 116, 206, 0, 15, 66, 64}
	require.EqualValues(t, expectedBytes, txBytes)
	expectedTxid := "2U2QKCYSYA3DUJNVP5T2KWBLWVUMX4PPIGN43IPPVGVZKTNFKJFQ"
	require.Equal(t, expectedTxid, txid)
}

func TestAppendMultisigTransactionKeyReg(t *testing.T) {
	uniSigTxBytes := []byte{130, 164, 109, 115, 105, 103, 131, 166, 115, 117, 98, 115, 105, 103, 147, 130, 162, 112, 107, 196, 32, 27, 126, 192, 176, 75, 234, 97, 183, 150, 144, 151, 230, 203, 244, 7, 225, 8, 167, 5, 53, 29, 11, 201, 138, 190, 177, 34, 9, 168, 171, 129, 120, 161, 115, 196, 64, 186, 52, 94, 163, 20, 123, 21, 228, 212, 78, 168, 14, 159, 234, 210, 219, 69, 206, 23, 113, 13, 3, 226, 107, 74, 6, 121, 202, 250, 195, 62, 13, 205, 64, 12, 208, 205, 69, 221, 116, 29, 15, 86, 243, 209, 159, 143, 116, 161, 84, 144, 104, 113, 8, 99, 78, 68, 12, 149, 213, 4, 83, 201, 15, 129, 162, 112, 107, 196, 32, 9, 99, 50, 9, 83, 115, 137, 240, 117, 103, 17, 119, 57, 145, 199, 208, 62, 27, 115, 200, 196, 245, 43, 246, 175, 240, 26, 162, 92, 249, 194, 113, 129, 162, 112, 107, 196, 32, 231, 240, 248, 77, 6, 129, 29, 249, 243, 28, 141, 135, 139, 17, 85, 244, 103, 29, 81, 161, 133, 194, 0, 144, 134, 103, 244, 73, 88, 112, 104, 161, 163, 116, 104, 114, 2, 161, 118, 1, 163, 116, 120, 110, 137, 163, 102, 101, 101, 206, 0, 3, 200, 192, 162, 102, 118, 206, 0, 14, 249, 218, 162, 108, 118, 206, 0, 14, 253, 194, 166, 115, 101, 108, 107, 101, 121, 196, 32, 50, 18, 43, 43, 214, 61, 220, 83, 49, 150, 23, 165, 170, 83, 196, 177, 194, 111, 227, 220, 202, 242, 141, 54, 34, 181, 105, 119, 161, 64, 92, 134, 163, 115, 110, 100, 196, 32, 141, 146, 180, 137, 144, 1, 115, 160, 77, 250, 67, 89, 163, 102, 106, 106, 252, 234, 44, 66, 160, 93, 217, 193, 247, 62, 235, 165, 71, 128, 55, 233, 164, 116, 121, 112, 101, 166, 107, 101, 121, 114, 101, 103, 166, 118, 111, 116, 101, 107, 100, 205, 39, 16, 167, 118, 111, 116, 101, 107, 101, 121, 196, 32, 112, 27, 215, 251, 145, 43, 7, 179, 8, 17, 255, 40, 29, 159, 238, 149, 99, 229, 128, 46, 32, 38, 137, 35, 25, 37, 143, 119, 250, 147, 30, 136, 167, 118, 111, 116, 101, 108, 115, 116, 206, 0, 15, 66, 64}
	ma, _, _, sk3 := makeTestMultisigAccount(t)
	txid, txBytes, err := AppendMultisigTransaction(sk3, ma, uniSigTxBytes)
	require.NoError(t, err)
	expectedBytes := []byte{130, 164, 109, 115, 105, 103, 131, 166, 115, 117, 98, 115, 105, 103, 147, 130, 162, 112, 107, 196, 32, 27, 126, 192, 176, 75, 234, 97, 183, 150, 144, 151, 230, 203, 244, 7, 225, 8, 167, 5, 53, 29, 11, 201, 138, 190, 177, 34, 9, 168, 171, 129, 120, 161, 115, 196, 64, 186, 52, 94, 163, 20, 123, 21, 228, 212, 78, 168, 14, 159, 234, 210, 219, 69, 206, 23, 113, 13, 3, 226, 107, 74, 6, 121, 202, 250, 195, 62, 13, 205, 64, 12, 208, 205, 69, 221, 116, 29, 15, 86, 243, 209, 159, 143, 116, 161, 84, 144, 104, 113, 8, 99, 78, 68, 12, 149, 213, 4, 83, 201, 15, 129, 162, 112, 107, 196, 32, 9, 99, 50, 9, 83, 115, 137, 240, 117, 103, 17, 119, 57, 145, 199, 208, 62, 27, 115, 200, 196, 245, 43, 246, 175, 240, 26, 162, 92, 249, 194, 113, 130, 162, 112, 107, 196, 32, 231, 240, 248, 77, 6, 129, 29, 249, 243, 28, 141, 135, 139, 17, 85, 244, 103, 29, 81, 161, 133, 194, 0, 144, 134, 103, 244, 73, 88, 112, 104, 161, 161, 115, 196, 64, 172, 133, 89, 89, 172, 158, 161, 188, 202, 74, 255, 179, 164, 146, 102, 110, 184, 236, 130, 86, 57, 39, 79, 127, 212, 165, 55, 237, 62, 92, 74, 94, 125, 230, 99, 40, 182, 163, 187, 107, 97, 230, 207, 69, 218, 71, 26, 18, 234, 149, 97, 177, 205, 152, 74, 67, 34, 83, 246, 33, 28, 144, 156, 3, 163, 116, 104, 114, 2, 161, 118, 1, 163, 116, 120, 110, 137, 163, 102, 101, 101, 206, 0, 3, 200, 192, 162, 102, 118, 206, 0, 14, 249, 218, 162, 108, 118, 206, 0, 14, 253, 194, 166, 115, 101, 108, 107, 101, 121, 196, 32, 50, 18, 43, 43, 214, 61, 220, 83, 49, 150, 23, 165, 170, 83, 196, 177, 194, 111, 227, 220, 202, 242, 141, 54, 34, 181, 105, 119, 161, 64, 92, 134, 163, 115, 110, 100, 196, 32, 141, 146, 180, 137, 144, 1, 115, 160, 77, 250, 67, 89, 163, 102, 106, 106, 252, 234, 44, 66, 160, 93, 217, 193, 247, 62, 235, 165, 71, 128, 55, 233, 164, 116, 121, 112, 101, 166, 107, 101, 121, 114, 101, 103, 166, 118, 111, 116, 101, 107, 100, 205, 39, 16, 167, 118, 111, 116, 101, 107, 101, 121, 196, 32, 112, 27, 215, 251, 145, 43, 7, 179, 8, 17, 255, 40, 29, 159, 238, 149, 99, 229, 128, 46, 32, 38, 137, 35, 25, 37, 143, 119, 250, 147, 30, 136, 167, 118, 111, 116, 101, 108, 115, 116, 206, 0, 15, 66, 64}
	require.EqualValues(t, expectedBytes, txBytes)
	expectedTxid := "2U2QKCYSYA3DUJNVP5T2KWBLWVUMX4PPIGN43IPPVGVZKTNFKJFQ"
	require.Equal(t, expectedTxid, txid)
}

func TestMergeMultisigTransactions(t *testing.T) {
	oneThreeSigTxBytes := []byte{130, 164, 109, 115, 105, 103, 131, 166, 115, 117, 98, 115, 105, 103, 147, 130, 162, 112, 107, 196, 32, 27, 126, 192, 176, 75, 234, 97, 183, 150, 144, 151, 230, 203, 244, 7, 225, 8, 167, 5, 53, 29, 11, 201, 138, 190, 177, 34, 9, 168, 171, 129, 120, 161, 115, 196, 64, 186, 52, 94, 163, 20, 123, 21, 228, 212, 78, 168, 14, 159, 234, 210, 219, 69, 206, 23, 113, 13, 3, 226, 107, 74, 6, 121, 202, 250, 195, 62, 13, 205, 64, 12, 208, 205, 69, 221, 116, 29, 15, 86, 243, 209, 159, 143, 116, 161, 84, 144, 104, 113, 8, 99, 78, 68, 12, 149, 213, 4, 83, 201, 15, 129, 162, 112, 107, 196, 32, 9, 99, 50, 9, 83, 115, 137, 240, 117, 103, 17, 119, 57, 145, 199, 208, 62, 27, 115, 200, 196, 245, 43, 246, 175, 240, 26, 162, 92, 249, 194, 113, 130, 162, 112, 107, 196, 32, 231, 240, 248, 77, 6, 129, 29, 249, 243, 28, 141, 135, 139, 17, 85, 244, 103, 29, 81, 161, 133, 194, 0, 144, 134, 103, 244, 73, 88, 112, 104, 161, 161, 115, 196, 64, 172, 133, 89, 89, 172, 158, 161, 188, 202, 74, 255, 179, 164, 146, 102, 110, 184, 236, 130, 86, 57, 39, 79, 127, 212, 165, 55, 237, 62, 92, 74, 94, 125, 230, 99, 40, 182, 163, 187, 107, 97, 230, 207, 69, 218, 71, 26, 18, 234, 149, 97, 177, 205, 152, 74, 67, 34, 83, 246, 33, 28, 144, 156, 3, 163, 116, 104, 114, 2, 161, 118, 1, 163, 116, 120, 110, 137, 163, 102, 101, 101, 206, 0, 3, 200, 192, 162, 102, 118, 206, 0, 14, 249, 218, 162, 108, 118, 206, 0, 14, 253, 194, 166, 115, 101, 108, 107, 101, 121, 196, 32, 50, 18, 43, 43, 214, 61, 220, 83, 49, 150, 23, 165, 170, 83, 196, 177, 194, 111, 227, 220, 202, 242, 141, 54, 34, 181, 105, 119, 161, 64, 92, 134, 163, 115, 110, 100, 196, 32, 141, 146, 180, 137, 144, 1, 115, 160, 77, 250, 67, 89, 163, 102, 106, 106, 252, 234, 44, 66, 160, 93, 217, 193, 247, 62, 235, 165, 71, 128, 55, 233, 164, 116, 121, 112, 101, 166, 107, 101, 121, 114, 101, 103, 166, 118, 111, 116, 101, 107, 100, 205, 39, 16, 167, 118, 111, 116, 101, 107, 101, 121, 196, 32, 112, 27, 215, 251, 145, 43, 7, 179, 8, 17, 255, 40, 29, 159, 238, 149, 99, 229, 128, 46, 32, 38, 137, 35, 25, 37, 143, 119, 250, 147, 30, 136, 167, 118, 111, 116, 101, 108, 115, 116, 206, 0, 15, 66, 64}
	twoThreeSigTxBytes := []byte{130, 164, 109, 115, 105, 103, 131, 166, 115, 117, 98, 115, 105, 103, 147, 129, 162, 112, 107, 196, 32, 27, 126, 192, 176, 75, 234, 97, 183, 150, 144, 151, 230, 203, 244, 7, 225, 8, 167, 5, 53, 29, 11, 201, 138, 190, 177, 34, 9, 168, 171, 129, 120, 130, 162, 112, 107, 196, 32, 9, 99, 50, 9, 83, 115, 137, 240, 117, 103, 17, 119, 57, 145, 199, 208, 62, 27, 115, 200, 196, 245, 43, 246, 175, 240, 26, 162, 92, 249, 194, 113, 161, 115, 196, 64, 191, 142, 166, 135, 208, 59, 232, 220, 86, 180, 101, 85, 236, 64, 3, 252, 51, 149, 11, 247, 226, 113, 205, 104, 169, 14, 112, 53, 194, 96, 41, 170, 89, 114, 185, 145, 228, 100, 220, 6, 209, 228, 152, 248, 176, 202, 48, 26, 1, 217, 102, 152, 112, 147, 86, 202, 146, 98, 226, 93, 95, 233, 162, 15, 130, 162, 112, 107, 196, 32, 231, 240, 248, 77, 6, 129, 29, 249, 243, 28, 141, 135, 139, 17, 85, 244, 103, 29, 81, 161, 133, 194, 0, 144, 134, 103, 244, 73, 88, 112, 104, 161, 161, 115, 196, 64, 172, 133, 89, 89, 172, 158, 161, 188, 202, 74, 255, 179, 164, 146, 102, 110, 184, 236, 130, 86, 57, 39, 79, 127, 212, 165, 55, 237, 62, 92, 74, 94, 125, 230, 99, 40, 182, 163, 187, 107, 97, 230, 207, 69, 218, 71, 26, 18, 234, 149, 97, 177, 205, 152, 74, 67, 34, 83, 246, 33, 28, 144, 156, 3, 163, 116, 104, 114, 2, 161, 118, 1, 163, 116, 120, 110, 137, 163, 102, 101, 101, 206, 0, 3, 200, 192, 162, 102, 118, 206, 0, 14, 249, 218, 162, 108, 118, 206, 0, 14, 253, 194, 166, 115, 101, 108, 107, 101, 121, 196, 32, 50, 18, 43, 43, 214, 61, 220, 83, 49, 150, 23, 165, 170, 83, 196, 177, 194, 111, 227, 220, 202, 242, 141, 54, 34, 181, 105, 119, 161, 64, 92, 134, 163, 115, 110, 100, 196, 32, 141, 146, 180, 137, 144, 1, 115, 160, 77, 250, 67, 89, 163, 102, 106, 106, 252, 234, 44, 66, 160, 93, 217, 193, 247, 62, 235, 165, 71, 128, 55, 233, 164, 116, 121, 112, 101, 166, 107, 101, 121, 114, 101, 103, 166, 118, 111, 116, 101, 107, 100, 205, 39, 16, 167, 118, 111, 116, 101, 107, 101, 121, 196, 32, 112, 27, 215, 251, 145, 43, 7, 179, 8, 17, 255, 40, 29, 159, 238, 149, 99, 229, 128, 46, 32, 38, 137, 35, 25, 37, 143, 119, 250, 147, 30, 136, 167, 118, 111, 116, 101, 108, 115, 116, 206, 0, 15, 66, 64}
	expectedBytes := []byte{130, 164, 109, 115, 105, 103, 131, 166, 115, 117, 98, 115, 105, 103, 147, 130, 162, 112, 107, 196, 32, 27, 126, 192, 176, 75, 234, 97, 183, 150, 144, 151, 230, 203, 244, 7, 225, 8, 167, 5, 53, 29, 11, 201, 138, 190, 177, 34, 9, 168, 171, 129, 120, 161, 115, 196, 64, 186, 52, 94, 163, 20, 123, 21, 228, 212, 78, 168, 14, 159, 234, 210, 219, 69, 206, 23, 113, 13, 3, 226, 107, 74, 6, 121, 202, 250, 195, 62, 13, 205, 64, 12, 208, 205, 69, 221, 116, 29, 15, 86, 243, 209, 159, 143, 116, 161, 84, 144, 104, 113, 8, 99, 78, 68, 12, 149, 213, 4, 83, 201, 15, 130, 162, 112, 107, 196, 32, 9, 99, 50, 9, 83, 115, 137, 240, 117, 103, 17, 119, 57, 145, 199, 208, 62, 27, 115, 200, 196, 245, 43, 246, 175, 240, 26, 162, 92, 249, 194, 113, 161, 115, 196, 64, 191, 142, 166, 135, 208, 59, 232, 220, 86, 180, 101, 85, 236, 64, 3, 252, 51, 149, 11, 247, 226, 113, 205, 104, 169, 14, 112, 53, 194, 96, 41, 170, 89, 114, 185, 145, 228, 100, 220, 6, 209, 228, 152, 248, 176, 202, 48, 26, 1, 217, 102, 152, 112, 147, 86, 202, 146, 98, 226, 93, 95, 233, 162, 15, 130, 162, 112, 107, 196, 32, 231, 240, 248, 77, 6, 129, 29, 249, 243, 28, 141, 135, 139, 17, 85, 244, 103, 29, 81, 161, 133, 194, 0, 144, 134, 103, 244, 73, 88, 112, 104, 161, 161, 115, 196, 64, 172, 133, 89, 89, 172, 158, 161, 188, 202, 74, 255, 179, 164, 146, 102, 110, 184, 236, 130, 86, 57, 39, 79, 127, 212, 165, 55, 237, 62, 92, 74, 94, 125, 230, 99, 40, 182, 163, 187, 107, 97, 230, 207, 69, 218, 71, 26, 18, 234, 149, 97, 177, 205, 152, 74, 67, 34, 83, 246, 33, 28, 144, 156, 3, 163, 116, 104, 114, 2, 161, 118, 1, 163, 116, 120, 110, 137, 163, 102, 101, 101, 206, 0, 3, 200, 192, 162, 102, 118, 206, 0, 14, 249, 218, 162, 108, 118, 206, 0, 14, 253, 194, 166, 115, 101, 108, 107, 101, 121, 196, 32, 50, 18, 43, 43, 214, 61, 220, 83, 49, 150, 23, 165, 170, 83, 196, 177, 194, 111, 227, 220, 202, 242, 141, 54, 34, 181, 105, 119, 161, 64, 92, 134, 163, 115, 110, 100, 196, 32, 141, 146, 180, 137, 144, 1, 115, 160, 77, 250, 67, 89, 163, 102, 106, 106, 252, 234, 44, 66, 160, 93, 217, 193, 247, 62, 235, 165, 71, 128, 55, 233, 164, 116, 121, 112, 101, 166, 107, 101, 121, 114, 101, 103, 166, 118, 111, 116, 101, 107, 100, 205, 39, 16, 167, 118, 111, 116, 101, 107, 101, 121, 196, 32, 112, 27, 215, 251, 145, 43, 7, 179, 8, 17, 255, 40, 29, 159, 238, 149, 99, 229, 128, 46, 32, 38, 137, 35, 25, 37, 143, 119, 250, 147, 30, 136, 167, 118, 111, 116, 101, 108, 115, 116, 206, 0, 15, 66, 64}
	_, txBytesSym, err := MergeMultisigTransactions(twoThreeSigTxBytes, oneThreeSigTxBytes)
	require.NoError(t, err)
	txid, txBytesSym2, err := MergeMultisigTransactions(oneThreeSigTxBytes, twoThreeSigTxBytes)
	require.NoError(t, err)
	// make sure merging is symmetric
	require.EqualValues(t, txBytesSym2, txBytesSym)
	// make sure they match expected output
	require.EqualValues(t, expectedBytes, txBytesSym)
	expectedTxid := "2U2QKCYSYA3DUJNVP5T2KWBLWVUMX4PPIGN43IPPVGVZKTNFKJFQ"
	require.Equal(t, expectedTxid, txid)
}

func TestSignBytes(t *testing.T) {
	account := GenerateAccount()
	message := make([]byte, 15)
	rand.Read(message)
	signature, err := SignBytes(account.PrivateKey, message)
	require.NoError(t, err)
	require.True(t, VerifyBytes(account.PublicKey, message, signature))
	if message[0] == 255 {
		message[0] = 0
	} else {
		message[0] = message[0] + 1
	}
	require.False(t, VerifyBytes(account.PublicKey, message, signature))
}

func TestMakeLogicSigBasic(t *testing.T) {
	// basic checks and contracts without delegation
	var program []byte
	var args [][]byte
	var sk ed25519.PrivateKey
	var pk MultisigAccount

	// check empty LogicSig
	lsig, err := MakeLogicSig(program, args, sk, pk)
	require.Error(t, err)
	require.Equal(t, types.LogicSig{}, lsig)

	program = []byte{1, 32, 1, 1, 34}
	programHash := "6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY"
	contractSender, err := types.DecodeAddress(programHash)
	require.NoError(t, err)

	lsig, err = MakeLogicSig(program, args, sk, pk)
	require.NoError(t, err)
	require.Equal(t, program, lsig.Logic)
	require.Equal(t, args, lsig.Args)
	require.Equal(t, types.Signature{}, lsig.Sig)
	require.True(t, lsig.Msig.Blank())
	verified := VerifyLogicSig(lsig, contractSender)
	require.True(t, verified)
	require.Equal(t, LogicSigAddress(lsig), contractSender)

	// check arguments
	args = make([][]byte, 2)
	args[0] = []byte{1, 2, 3}
	args[1] = []byte{4, 5, 6}
	lsig, err = MakeLogicSig(program, args, sk, pk)
	require.NoError(t, err)
	require.Equal(t, program, lsig.Logic)
	require.Equal(t, args, lsig.Args)
	require.Equal(t, types.Signature{}, lsig.Sig)
	require.True(t, lsig.Msig.Blank())
	verified = VerifyLogicSig(lsig, contractSender)
	require.True(t, verified)

	// check serialization
	var lsig1 types.LogicSig
	encoded := msgpack.Encode(lsig)
	err = msgpack.Decode(encoded, &lsig1)
	require.NoError(t, err)
	require.Equal(t, lsig, lsig1)

	// check verification of modified program fails
	programMod := make([]byte, len(program))
	copy(programMod[:], program)
	programMod[3] = 2
	lsig, err = MakeLogicSig(programMod, args, sk, pk)
	require.NoError(t, err)
	verified = VerifyLogicSig(lsig, contractSender)
	require.False(t, verified)

	// check invalid program fails
	copy(programMod[:], program)
	programMod[0] = 128
	lsig, err = MakeLogicSig(programMod, args, sk, pk)
	require.Error(t, err)
}

func TestMakeLogicSigSingle(t *testing.T) {
	var program []byte
	var args [][]byte
	var sk ed25519.PrivateKey
	var pk MultisigAccount

	acc := GenerateAccount()
	program = []byte{1, 32, 1, 1, 34}
	sk = acc.PrivateKey
	lsig, err := MakeLogicSig(program, args, sk, pk)
	require.NoError(t, err)
	require.NotEqual(t, types.Signature{}, lsig.Sig)
	require.True(t, lsig.Msig.Blank())

	var sender types.Address
	copy(sender[:], acc.PublicKey)

	verified := VerifyLogicSig(lsig, sender)
	require.True(t, verified)

	// check serialization
	var lsig1 types.LogicSig
	encoded := msgpack.Encode(lsig)
	err = msgpack.Decode(encoded, &lsig1)
	require.NoError(t, err)
	require.Equal(t, lsig, lsig1)
}

func TestMakeLogicSigMulti(t *testing.T) {
	var program []byte
	var args [][]byte
	var sk ed25519.PrivateKey
	var pk MultisigAccount

	ma, sk1, sk2, _ := makeTestMultisigAccount(t)
	program = []byte{1, 32, 1, 1, 34}
	sender, err := ma.Address()
	require.NoError(t, err)
	acc := GenerateAccount()
	sk = acc.PrivateKey

	lsig, err := MakeLogicSig(program, args, sk1, ma)
	require.NoError(t, err)
	require.Equal(t, program, lsig.Logic)
	require.Equal(t, args, lsig.Args)
	require.Equal(t, types.Signature{}, lsig.Sig)
	require.False(t, lsig.Msig.Blank())

	verified := VerifyLogicSig(lsig, sender)
	require.False(t, verified) // not enough signatures

	err = AppendMultisigToLogicSig(&lsig, sk)
	require.Error(t, err) // sk not part of multisig

	err = AppendMultisigToLogicSig(&lsig, sk2)
	require.NoError(t, err)

	verified = VerifyLogicSig(lsig, sender)
	require.True(t, verified)

	// combine sig and multisig, ensure it fails
	lsigf, err := MakeLogicSig(program, args, sk, pk)
	require.NoError(t, err)
	lsig.Sig = lsigf.Sig

	verified = VerifyLogicSig(lsig, sender)
	require.False(t, verified) // sig + msig

	// remove sig and ensure things are good
	lsig.Sig = types.Signature{}
	verified = VerifyLogicSig(lsig, sender)
	require.True(t, verified)

	// check serialization
	var lsig1 types.LogicSig
	encoded := msgpack.Encode(lsig)
	err = msgpack.Decode(encoded, &lsig1)
	require.NoError(t, err)
	require.Equal(t, lsig, lsig1)
}

func TestTealSign(t *testing.T) {
	data, err := base64.StdEncoding.DecodeString("Ux8jntyBJQarjKGF8A==")
	require.NoError(t, err)

	seed, err := base64.StdEncoding.DecodeString("5Pf7eGMA52qfMT4R4/vYCt7con/7U3yejkdXkrcb26Q=")
	require.NoError(t, err)
	sk := ed25519.NewKeyFromSeed(seed)

	addr, err := types.DecodeAddress("6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY")
	require.NoError(t, err)

	prog, err := base64.StdEncoding.DecodeString("ASABASI=")
	require.NoError(t, err)

	sig1, err := TealSign(sk, data, addr)
	require.NoError(t, err)

	sig2, err := TealSignFromProgram(sk, data, prog)
	require.NoError(t, err)

	require.Equal(t, sig1, sig2)

	pk := sk.Public().(ed25519.PublicKey)
	msg := bytes.Join([][]byte{programDataPrefix, addr[:], data}, nil)
	verified := ed25519.Verify(pk, msg, sig1[:])
	require.True(t, verified)
}
