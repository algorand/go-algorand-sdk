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

	require.Equal(t, types.Address{}, stx.AuthAddr)
	require.Equal(t, tx, stx.Txn)

	bytesToSign := rawTransactionBytesToSign(stx.Txn)
	verified := VerifyMultisig(fromAddr, bytesToSign, stx.Msig)
	require.False(t, verified) // not enough signatures
}

func TestSignMultisigTransactionWithAuthAddr(t *testing.T) {
	ma, sk1, _, _ := makeTestMultisigAccount(t)
	multisigAddr, err := ma.Address()
	require.NoError(t, err)
	fromAddr, err := types.DecodeAddress("47YPQTIGQEO7T4Y4RWDYWEKV6RTR2UNBQXBABEEGM72ESWDQNCQ52OPASU")
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
	expectedBytes := []byte{0x83, 0xa4, 0x6d, 0x73, 0x69, 0x67, 0x83, 0xa6, 0x73, 0x75, 0x62, 0x73, 0x69, 0x67, 0x93, 0x82, 0xa2, 0x70, 0x6b, 0xc4, 0x20, 0x1b, 0x7e, 0xc0, 0xb0, 0x4b, 0xea, 0x61, 0xb7, 0x96, 0x90, 0x97, 0xe6, 0xcb, 0xf4, 0x7, 0xe1, 0x8, 0xa7, 0x5, 0x35, 0x1d, 0xb, 0xc9, 0x8a, 0xbe, 0xb1, 0x22, 0x9, 0xa8, 0xab, 0x81, 0x78, 0xa1, 0x73, 0xc4, 0x40, 0x3e, 0xd3, 0xf2, 0x25, 0x45, 0x61, 0x38, 0x4a, 0x72, 0xca, 0x59, 0x2b, 0xd7, 0xca, 0xfa, 0x7, 0xe, 0x1, 0xb1, 0xce, 0x3c, 0x2e, 0xd8, 0x8c, 0x70, 0x80, 0xa7, 0x9d, 0xe2, 0x88, 0xf7, 0xa2, 0xa3, 0xc9, 0xee, 0x88, 0x0, 0xb, 0x6c, 0x2c, 0x6a, 0x7, 0xca, 0x53, 0x22, 0xcb, 0xc, 0x78, 0x82, 0xcd, 0x42, 0xc9, 0x76, 0x6a, 0x72, 0xed, 0x3f, 0x7c, 0x8a, 0x47, 0x47, 0xaf, 0xb7, 0x0, 0x81, 0xa2, 0x70, 0x6b, 0xc4, 0x20, 0x9, 0x63, 0x32, 0x9, 0x53, 0x73, 0x89, 0xf0, 0x75, 0x67, 0x11, 0x77, 0x39, 0x91, 0xc7, 0xd0, 0x3e, 0x1b, 0x73, 0xc8, 0xc4, 0xf5, 0x2b, 0xf6, 0xaf, 0xf0, 0x1a, 0xa2, 0x5c, 0xf9, 0xc2, 0x71, 0x81, 0xa2, 0x70, 0x6b, 0xc4, 0x20, 0xe7, 0xf0, 0xf8, 0x4d, 0x6, 0x81, 0x1d, 0xf9, 0xf3, 0x1c, 0x8d, 0x87, 0x8b, 0x11, 0x55, 0xf4, 0x67, 0x1d, 0x51, 0xa1, 0x85, 0xc2, 0x0, 0x90, 0x86, 0x67, 0xf4, 0x49, 0x58, 0x70, 0x68, 0xa1, 0xa3, 0x74, 0x68, 0x72, 0x2, 0xa1, 0x76, 0x1, 0xa4, 0x73, 0x67, 0x6e, 0x72, 0xc4, 0x20, 0x8d, 0x92, 0xb4, 0x89, 0x90, 0x1, 0x73, 0xa0, 0x4d, 0xfa, 0x43, 0x59, 0xa3, 0x66, 0x6a, 0x6a, 0xfc, 0xea, 0x2c, 0x42, 0xa0, 0x5d, 0xd9, 0xc1, 0xf7, 0x3e, 0xeb, 0xa5, 0x47, 0x80, 0x37, 0xe9, 0xa3, 0x74, 0x78, 0x6e, 0x89, 0xa3, 0x61, 0x6d, 0x74, 0xcd, 0x13, 0x88, 0xa3, 0x66, 0x65, 0x65, 0xce, 0x0, 0x3, 0x4f, 0xa8, 0xa2, 0x66, 0x76, 0xce, 0x0, 0xe, 0xd6, 0xdc, 0xa3, 0x67, 0x65, 0x6e, 0xad, 0x74, 0x65, 0x73, 0x74, 0x6e, 0x65, 0x74, 0x2d, 0x76, 0x33, 0x31, 0x2e, 0x30, 0xa2, 0x6c, 0x76, 0xce, 0x0, 0xe, 0xda, 0xc4, 0xa4, 0x6e, 0x6f, 0x74, 0x65, 0xc4, 0x8, 0xb4, 0x51, 0x79, 0x39, 0xfc, 0xfa, 0xd2, 0x71, 0xa3, 0x72, 0x63, 0x76, 0xc4, 0x20, 0x1b, 0x7e, 0xc0, 0xb0, 0x4b, 0xea, 0x61, 0xb7, 0x96, 0x90, 0x97, 0xe6, 0xcb, 0xf4, 0x7, 0xe1, 0x8, 0xa7, 0x5, 0x35, 0x1d, 0xb, 0xc9, 0x8a, 0xbe, 0xb1, 0x22, 0x9, 0xa8, 0xab, 0x81, 0x78, 0xa3, 0x73, 0x6e, 0x64, 0xc4, 0x20, 0xe7, 0xf0, 0xf8, 0x4d, 0x6, 0x81, 0x1d, 0xf9, 0xf3, 0x1c, 0x8d, 0x87, 0x8b, 0x11, 0x55, 0xf4, 0x67, 0x1d, 0x51, 0xa1, 0x85, 0xc2, 0x0, 0x90, 0x86, 0x67, 0xf4, 0x49, 0x58, 0x70, 0x68, 0xa1, 0xa4, 0x74, 0x79, 0x70, 0x65, 0xa3, 0x70, 0x61, 0x79}
	require.EqualValues(t, expectedBytes, txBytes)
	expectedTxid := "HXNXFLSFZQXGOEN33IIQU4OAVZTGXIYHE4HPWFG3RDKKA4OM64JQ"
	require.Equal(t, expectedTxid, txid)

	// decode and verify
	var stx types.SignedTxn
	err = msgpack.Decode(txBytes, &stx)
	require.NoError(t, err)

	require.Equal(t, multisigAddr, stx.AuthAddr)
	require.Equal(t, tx, stx.Txn)

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
	require.True(t, lsig.Blank())

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

	// check invalid program fails
	programMod := make([]byte, len(program))
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

	acc, err := AccountFromPrivateKey(ed25519.PrivateKey{0xd2, 0xdc, 0x4c, 0xcc, 0xe9, 0x98, 0x62, 0xff, 0xcf, 0x8c, 0xeb, 0x93, 0x6, 0xc4, 0x8d, 0xa6, 0x80, 0x50, 0x82, 0xa, 0xbb, 0x29, 0x95, 0x7a, 0xac, 0x82, 0x68, 0x9a, 0x8c, 0x49, 0x5a, 0x38, 0x5e, 0x67, 0x4f, 0x1c, 0xa, 0xee, 0xec, 0x37, 0x71, 0x89, 0x8f, 0x61, 0xc7, 0x6f, 0xf5, 0xd2, 0x4a, 0x19, 0x79, 0x3e, 0x2c, 0x91, 0xfa, 0x8, 0x51, 0x62, 0x63, 0xe3, 0x85, 0x73, 0xea, 0x42})
	require.NoError(t, err)
	program = []byte{1, 32, 1, 1, 34}
	sk = acc.PrivateKey
	lsig, err := MakeLogicSig(program, args, sk, pk)
	require.NoError(t, err)
	expectedSig := types.Signature{0x3e, 0x5, 0x3d, 0x39, 0x4d, 0xfb, 0x12, 0xbc, 0x65, 0x79, 0x9f, 0xea, 0x31, 0x8a, 0x7b, 0x8e, 0xa2, 0x51, 0x8b, 0x55, 0x2c, 0x8a, 0xbe, 0x6c, 0xd7, 0xa7, 0x65, 0x2d, 0xd8, 0xb0, 0x18, 0x7e, 0x21, 0x5, 0x2d, 0xb9, 0x24, 0x62, 0x89, 0x16, 0xe5, 0x61, 0x74, 0xcd, 0xf, 0x19, 0xac, 0xb9, 0x6c, 0x45, 0xa4, 0x29, 0x91, 0x99, 0x11, 0x1d, 0xe4, 0x7c, 0xe4, 0xfc, 0x12, 0xec, 0xce, 0x2}
	require.Equal(t, expectedSig, lsig.Sig)
	require.True(t, lsig.Msig.Blank())

	verified := VerifyLogicSig(lsig, acc.Address)
	require.True(t, verified)

	// check that a modified program fails verification
	modProgram := make([]byte, len(program))
	copy(modProgram, program)
	lsigModified, err := MakeLogicSig(modProgram, args, sk, pk)
	require.NoError(t, err)
	modProgram[3] = 2
	verified = VerifyLogicSig(lsigModified, acc.Address)
	require.False(t, verified)

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

	// check that a modified program fails verification
	modProgram := make([]byte, len(program))
	copy(modProgram, program)
	lsigModified, err := MakeLogicSig(modProgram, args, sk1, ma)
	require.NoError(t, err)
	modProgram[3] = 2
	verified = VerifyLogicSig(lsigModified, sender)
	require.False(t, verified)

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

func TestSignLogicsigTransaction(t *testing.T) {
	program := []byte{1, 32, 1, 1, 34}
	args := [][]byte{
		{0x01},
		{0x02, 0x03},
	}

	otherAddrStr := "WTDCE2FEYM2VB5MKNXKLRSRDTSPR2EFTIGVH4GRW4PHGD6747GFJTBGT2A"
	otherAddr, err := types.DecodeAddress(otherAddrStr)
	require.NoError(t, err)

	testSign := func(t *testing.T, lsig types.LogicSig, sender types.Address, expectedBytes []byte, expectedTxid string, expectedAuthAddr types.Address) {
		txn := types.Transaction{
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
				Receiver: otherAddr,
				Amount:   5000,
			},
		}

		txid, stxnBytes, err := SignLogicsigTransaction(lsig, txn)
		require.NoError(t, err)
		require.EqualValues(t, expectedBytes, stxnBytes)
		require.Equal(t, expectedTxid, txid)

		// decode and verify
		var stx types.SignedTxn
		err = msgpack.Decode(stxnBytes, &stx)
		require.NoError(t, err)

		require.Equal(t, expectedAuthAddr, stx.AuthAddr)

		require.Equal(t, lsig, stx.Lsig)
		require.Equal(t, types.Signature{}, stx.Sig)
		require.True(t, stx.Msig.Blank())

		require.Equal(t, txn, stx.Txn)
	}

	t.Run("no sig", func(t *testing.T) {
		var sk ed25519.PrivateKey
		var ma MultisigAccount
		lsig, err := MakeLogicSig(program, args, sk, ma)
		require.NoError(t, err)

		programHash := "6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY"
		programAddr, err := types.DecodeAddress(programHash)
		require.NoError(t, err)

		t.Run("sender is contract addr", func(t *testing.T) {
			expectedBytes := []byte{0x82, 0xa4, 0x6c, 0x73, 0x69, 0x67, 0x82, 0xa3, 0x61, 0x72, 0x67, 0x92, 0xc4, 0x1, 0x1, 0xc4, 0x2, 0x2, 0x3, 0xa1, 0x6c, 0xc4, 0x5, 0x1, 0x20, 0x1, 0x1, 0x22, 0xa3, 0x74, 0x78, 0x6e, 0x89, 0xa3, 0x61, 0x6d, 0x74, 0xcd, 0x13, 0x88, 0xa3, 0x66, 0x65, 0x65, 0xce, 0x0, 0x3, 0x4f, 0xa8, 0xa2, 0x66, 0x76, 0xce, 0x0, 0xe, 0xd6, 0xdc, 0xa3, 0x67, 0x65, 0x6e, 0xad, 0x74, 0x65, 0x73, 0x74, 0x6e, 0x65, 0x74, 0x2d, 0x76, 0x33, 0x31, 0x2e, 0x30, 0xa2, 0x6c, 0x76, 0xce, 0x0, 0xe, 0xda, 0xc4, 0xa4, 0x6e, 0x6f, 0x74, 0x65, 0xc4, 0x8, 0xb4, 0x51, 0x79, 0x39, 0xfc, 0xfa, 0xd2, 0x71, 0xa3, 0x72, 0x63, 0x76, 0xc4, 0x20, 0xb4, 0xc6, 0x22, 0x68, 0xa4, 0xc3, 0x35, 0x50, 0xf5, 0x8a, 0x6d, 0xd4, 0xb8, 0xca, 0x23, 0x9c, 0x9f, 0x1d, 0x10, 0xb3, 0x41, 0xaa, 0x7e, 0x1a, 0x36, 0xe3, 0xce, 0x61, 0xfb, 0xfc, 0xf9, 0x8a, 0xa3, 0x73, 0x6e, 0x64, 0xc4, 0x20, 0xf6, 0x76, 0x2d, 0xac, 0x75, 0xb1, 0x99, 0x7d, 0x6c, 0x2c, 0x96, 0x18, 0x6, 0x80, 0x50, 0x74, 0x90, 0xd7, 0x95, 0x11, 0x2f, 0xfe, 0x7f, 0xb7, 0x60, 0xb2, 0x73, 0x8a, 0xf9, 0xc7, 0xf1, 0xad, 0xa4, 0x74, 0x79, 0x70, 0x65, 0xa3, 0x70, 0x61, 0x79}
			expectedTxid := "IL5UCKXGWBA2MQ4YYFQKYC3BFCWO2KHZSNZWVDIXOOZS3AWVIQDA"
			expectedAuthAddr := types.Address{}
			sender := programAddr
			testSign(t, lsig, sender, expectedBytes, expectedTxid, expectedAuthAddr)
		})

		t.Run("sender is not contract addr", func(t *testing.T) {
			expectedBytes := []byte{0x83, 0xa4, 0x6c, 0x73, 0x69, 0x67, 0x82, 0xa3, 0x61, 0x72, 0x67, 0x92, 0xc4, 0x1, 0x1, 0xc4, 0x2, 0x2, 0x3, 0xa1, 0x6c, 0xc4, 0x5, 0x1, 0x20, 0x1, 0x1, 0x22, 0xa4, 0x73, 0x67, 0x6e, 0x72, 0xc4, 0x20, 0xf6, 0x76, 0x2d, 0xac, 0x75, 0xb1, 0x99, 0x7d, 0x6c, 0x2c, 0x96, 0x18, 0x6, 0x80, 0x50, 0x74, 0x90, 0xd7, 0x95, 0x11, 0x2f, 0xfe, 0x7f, 0xb7, 0x60, 0xb2, 0x73, 0x8a, 0xf9, 0xc7, 0xf1, 0xad, 0xa3, 0x74, 0x78, 0x6e, 0x89, 0xa3, 0x61, 0x6d, 0x74, 0xcd, 0x13, 0x88, 0xa3, 0x66, 0x65, 0x65, 0xce, 0x0, 0x3, 0x4f, 0xa8, 0xa2, 0x66, 0x76, 0xce, 0x0, 0xe, 0xd6, 0xdc, 0xa3, 0x67, 0x65, 0x6e, 0xad, 0x74, 0x65, 0x73, 0x74, 0x6e, 0x65, 0x74, 0x2d, 0x76, 0x33, 0x31, 0x2e, 0x30, 0xa2, 0x6c, 0x76, 0xce, 0x0, 0xe, 0xda, 0xc4, 0xa4, 0x6e, 0x6f, 0x74, 0x65, 0xc4, 0x8, 0xb4, 0x51, 0x79, 0x39, 0xfc, 0xfa, 0xd2, 0x71, 0xa3, 0x72, 0x63, 0x76, 0xc4, 0x20, 0xb4, 0xc6, 0x22, 0x68, 0xa4, 0xc3, 0x35, 0x50, 0xf5, 0x8a, 0x6d, 0xd4, 0xb8, 0xca, 0x23, 0x9c, 0x9f, 0x1d, 0x10, 0xb3, 0x41, 0xaa, 0x7e, 0x1a, 0x36, 0xe3, 0xce, 0x61, 0xfb, 0xfc, 0xf9, 0x8a, 0xa3, 0x73, 0x6e, 0x64, 0xc4, 0x20, 0xb4, 0xc6, 0x22, 0x68, 0xa4, 0xc3, 0x35, 0x50, 0xf5, 0x8a, 0x6d, 0xd4, 0xb8, 0xca, 0x23, 0x9c, 0x9f, 0x1d, 0x10, 0xb3, 0x41, 0xaa, 0x7e, 0x1a, 0x36, 0xe3, 0xce, 0x61, 0xfb, 0xfc, 0xf9, 0x8a, 0xa4, 0x74, 0x79, 0x70, 0x65, 0xa3, 0x70, 0x61, 0x79}
			expectedTxid := "U4X24Q45MCZ6JSL343QNR3RC6AJO2NUDQXFCTNONIW3SU2AYQJOA"
			expectedAuthAddr := programAddr
			sender := otherAddr
			testSign(t, lsig, sender, expectedBytes, expectedTxid, expectedAuthAddr)
		})
	})

	t.Run("single sig", func(t *testing.T) {
		var ma MultisigAccount
		acc, err := AccountFromPrivateKey(ed25519.PrivateKey{0xd2, 0xdc, 0x4c, 0xcc, 0xe9, 0x98, 0x62, 0xff, 0xcf, 0x8c, 0xeb, 0x93, 0x6, 0xc4, 0x8d, 0xa6, 0x80, 0x50, 0x82, 0xa, 0xbb, 0x29, 0x95, 0x7a, 0xac, 0x82, 0x68, 0x9a, 0x8c, 0x49, 0x5a, 0x38, 0x5e, 0x67, 0x4f, 0x1c, 0xa, 0xee, 0xec, 0x37, 0x71, 0x89, 0x8f, 0x61, 0xc7, 0x6f, 0xf5, 0xd2, 0x4a, 0x19, 0x79, 0x3e, 0x2c, 0x91, 0xfa, 0x8, 0x51, 0x62, 0x63, 0xe3, 0x85, 0x73, 0xea, 0x42})
		require.NoError(t, err)
		lsig, err := MakeLogicSig(program, args, acc.PrivateKey, ma)
		require.NoError(t, err)

		t.Run("sender is contract addr", func(t *testing.T) {
			expectedBytes := []byte{0x82, 0xa4, 0x6c, 0x73, 0x69, 0x67, 0x83, 0xa3, 0x61, 0x72, 0x67, 0x92, 0xc4, 0x1, 0x1, 0xc4, 0x2, 0x2, 0x3, 0xa1, 0x6c, 0xc4, 0x5, 0x1, 0x20, 0x1, 0x1, 0x22, 0xa3, 0x73, 0x69, 0x67, 0xc4, 0x40, 0x3e, 0x5, 0x3d, 0x39, 0x4d, 0xfb, 0x12, 0xbc, 0x65, 0x79, 0x9f, 0xea, 0x31, 0x8a, 0x7b, 0x8e, 0xa2, 0x51, 0x8b, 0x55, 0x2c, 0x8a, 0xbe, 0x6c, 0xd7, 0xa7, 0x65, 0x2d, 0xd8, 0xb0, 0x18, 0x7e, 0x21, 0x5, 0x2d, 0xb9, 0x24, 0x62, 0x89, 0x16, 0xe5, 0x61, 0x74, 0xcd, 0xf, 0x19, 0xac, 0xb9, 0x6c, 0x45, 0xa4, 0x29, 0x91, 0x99, 0x11, 0x1d, 0xe4, 0x7c, 0xe4, 0xfc, 0x12, 0xec, 0xce, 0x2, 0xa3, 0x74, 0x78, 0x6e, 0x89, 0xa3, 0x61, 0x6d, 0x74, 0xcd, 0x13, 0x88, 0xa3, 0x66, 0x65, 0x65, 0xce, 0x0, 0x3, 0x4f, 0xa8, 0xa2, 0x66, 0x76, 0xce, 0x0, 0xe, 0xd6, 0xdc, 0xa3, 0x67, 0x65, 0x6e, 0xad, 0x74, 0x65, 0x73, 0x74, 0x6e, 0x65, 0x74, 0x2d, 0x76, 0x33, 0x31, 0x2e, 0x30, 0xa2, 0x6c, 0x76, 0xce, 0x0, 0xe, 0xda, 0xc4, 0xa4, 0x6e, 0x6f, 0x74, 0x65, 0xc4, 0x8, 0xb4, 0x51, 0x79, 0x39, 0xfc, 0xfa, 0xd2, 0x71, 0xa3, 0x72, 0x63, 0x76, 0xc4, 0x20, 0xb4, 0xc6, 0x22, 0x68, 0xa4, 0xc3, 0x35, 0x50, 0xf5, 0x8a, 0x6d, 0xd4, 0xb8, 0xca, 0x23, 0x9c, 0x9f, 0x1d, 0x10, 0xb3, 0x41, 0xaa, 0x7e, 0x1a, 0x36, 0xe3, 0xce, 0x61, 0xfb, 0xfc, 0xf9, 0x8a, 0xa3, 0x73, 0x6e, 0x64, 0xc4, 0x20, 0x5e, 0x67, 0x4f, 0x1c, 0xa, 0xee, 0xec, 0x37, 0x71, 0x89, 0x8f, 0x61, 0xc7, 0x6f, 0xf5, 0xd2, 0x4a, 0x19, 0x79, 0x3e, 0x2c, 0x91, 0xfa, 0x8, 0x51, 0x62, 0x63, 0xe3, 0x85, 0x73, 0xea, 0x42, 0xa4, 0x74, 0x79, 0x70, 0x65, 0xa3, 0x70, 0x61, 0x79}
			expectedTxid := "XPFTYDOV5RCL7K5EGTC32PSTRKFO3ITZIL64CRTYIJEPHXOTCCXA"
			expectedAuthAddr := types.Address{}
			sender := acc.Address
			testSign(t, lsig, sender, expectedBytes, expectedTxid, expectedAuthAddr)
		})

		t.Run("sender is not contract addr", func(t *testing.T) {
			txn := types.Transaction{
				Type: types.PaymentTx,
				Header: types.Header{
					Sender:     otherAddr,
					Fee:        217000,
					FirstValid: 972508,
					LastValid:  973508,
					Note:       []byte{180, 81, 121, 57, 252, 250, 210, 113},
					GenesisID:  "testnet-v31.0",
				},
				PaymentTxnFields: types.PaymentTxnFields{
					Receiver: otherAddr,
					Amount:   5000,
				},
			}

			_, _, err := SignLogicsigTransaction(lsig, txn)
			require.Error(t, err, errLsigInvalidSignature)
		})
	})

	t.Run("multi sig", func(t *testing.T) {
		ma, sk1, sk2, _ := makeTestMultisigAccount(t)
		maAddr, err := ma.Address()
		require.NoError(t, err)

		lsig, err := MakeLogicSig(program, args, sk1, ma)
		require.NoError(t, err)

		err = AppendMultisigToLogicSig(&lsig, sk2)
		require.NoError(t, err)

		t.Run("sender is contract addr", func(t *testing.T) {
			expectedBytes := []byte{0x82, 0xa4, 0x6c, 0x73, 0x69, 0x67, 0x83, 0xa3, 0x61, 0x72, 0x67, 0x92, 0xc4, 0x1, 0x1, 0xc4, 0x2, 0x2, 0x3, 0xa1, 0x6c, 0xc4, 0x5, 0x1, 0x20, 0x1, 0x1, 0x22, 0xa4, 0x6d, 0x73, 0x69, 0x67, 0x83, 0xa6, 0x73, 0x75, 0x62, 0x73, 0x69, 0x67, 0x93, 0x82, 0xa2, 0x70, 0x6b, 0xc4, 0x20, 0x1b, 0x7e, 0xc0, 0xb0, 0x4b, 0xea, 0x61, 0xb7, 0x96, 0x90, 0x97, 0xe6, 0xcb, 0xf4, 0x7, 0xe1, 0x8, 0xa7, 0x5, 0x35, 0x1d, 0xb, 0xc9, 0x8a, 0xbe, 0xb1, 0x22, 0x9, 0xa8, 0xab, 0x81, 0x78, 0xa1, 0x73, 0xc4, 0x40, 0x49, 0x13, 0xb8, 0x5, 0xd1, 0x9e, 0x7f, 0x2c, 0x10, 0x80, 0xf6, 0x33, 0x7e, 0x18, 0x54, 0xa7, 0xce, 0xea, 0xee, 0x10, 0xdd, 0xbd, 0x13, 0x65, 0x84, 0xbf, 0x93, 0xb7, 0x5f, 0x30, 0x63, 0x15, 0x91, 0xca, 0x23, 0xc, 0xed, 0xef, 0x23, 0xd1, 0x74, 0x1b, 0x52, 0x9d, 0xb0, 0xff, 0xef, 0x37, 0x54, 0xd6, 0x46, 0xf4, 0xb5, 0x61, 0xfc, 0x8b, 0xbc, 0x2d, 0x7b, 0x4e, 0x63, 0x5c, 0xbd, 0x2, 0x82, 0xa2, 0x70, 0x6b, 0xc4, 0x20, 0x9, 0x63, 0x32, 0x9, 0x53, 0x73, 0x89, 0xf0, 0x75, 0x67, 0x11, 0x77, 0x39, 0x91, 0xc7, 0xd0, 0x3e, 0x1b, 0x73, 0xc8, 0xc4, 0xf5, 0x2b, 0xf6, 0xaf, 0xf0, 0x1a, 0xa2, 0x5c, 0xf9, 0xc2, 0x71, 0xa1, 0x73, 0xc4, 0x40, 0x64, 0xbc, 0x55, 0xdb, 0xed, 0x91, 0xa2, 0x41, 0xd4, 0x2a, 0xb6, 0x60, 0xf7, 0xe1, 0x4a, 0xb9, 0x99, 0x9a, 0x52, 0xb3, 0xb1, 0x71, 0x58, 0xce, 0xfc, 0x3f, 0x4f, 0xe7, 0xcb, 0x22, 0x41, 0x14, 0xad, 0xa9, 0x3d, 0x5e, 0x84, 0x5, 0x2, 0xa, 0x17, 0xa6, 0x69, 0x83, 0x3, 0x22, 0x4e, 0x86, 0xa3, 0x8b, 0x6a, 0x36, 0xc5, 0x54, 0xbe, 0x20, 0x50, 0xff, 0xd3, 0xee, 0xa8, 0xb3, 0x4, 0x9, 0x81, 0xa2, 0x70, 0x6b, 0xc4, 0x20, 0xe7, 0xf0, 0xf8, 0x4d, 0x6, 0x81, 0x1d, 0xf9, 0xf3, 0x1c, 0x8d, 0x87, 0x8b, 0x11, 0x55, 0xf4, 0x67, 0x1d, 0x51, 0xa1, 0x85, 0xc2, 0x0, 0x90, 0x86, 0x67, 0xf4, 0x49, 0x58, 0x70, 0x68, 0xa1, 0xa3, 0x74, 0x68, 0x72, 0x2, 0xa1, 0x76, 0x1, 0xa3, 0x74, 0x78, 0x6e, 0x89, 0xa3, 0x61, 0x6d, 0x74, 0xcd, 0x13, 0x88, 0xa3, 0x66, 0x65, 0x65, 0xce, 0x0, 0x3, 0x4f, 0xa8, 0xa2, 0x66, 0x76, 0xce, 0x0, 0xe, 0xd6, 0xdc, 0xa3, 0x67, 0x65, 0x6e, 0xad, 0x74, 0x65, 0x73, 0x74, 0x6e, 0x65, 0x74, 0x2d, 0x76, 0x33, 0x31, 0x2e, 0x30, 0xa2, 0x6c, 0x76, 0xce, 0x0, 0xe, 0xda, 0xc4, 0xa4, 0x6e, 0x6f, 0x74, 0x65, 0xc4, 0x8, 0xb4, 0x51, 0x79, 0x39, 0xfc, 0xfa, 0xd2, 0x71, 0xa3, 0x72, 0x63, 0x76, 0xc4, 0x20, 0xb4, 0xc6, 0x22, 0x68, 0xa4, 0xc3, 0x35, 0x50, 0xf5, 0x8a, 0x6d, 0xd4, 0xb8, 0xca, 0x23, 0x9c, 0x9f, 0x1d, 0x10, 0xb3, 0x41, 0xaa, 0x7e, 0x1a, 0x36, 0xe3, 0xce, 0x61, 0xfb, 0xfc, 0xf9, 0x8a, 0xa3, 0x73, 0x6e, 0x64, 0xc4, 0x20, 0x8d, 0x92, 0xb4, 0x89, 0x90, 0x1, 0x73, 0xa0, 0x4d, 0xfa, 0x43, 0x59, 0xa3, 0x66, 0x6a, 0x6a, 0xfc, 0xea, 0x2c, 0x42, 0xa0, 0x5d, 0xd9, 0xc1, 0xf7, 0x3e, 0xeb, 0xa5, 0x47, 0x80, 0x37, 0xe9, 0xa4, 0x74, 0x79, 0x70, 0x65, 0xa3, 0x70, 0x61, 0x79}
			expectedTxid := "2I2QT3MGXNMB5IOTNWEQZWUJCHLJ5QFBYI264X5NWMUDKZXPZ5RA"
			expectedAuthAddr := types.Address{}
			sender := maAddr
			testSign(t, lsig, sender, expectedBytes, expectedTxid, expectedAuthAddr)
		})

		t.Run("sender is not contract addr", func(t *testing.T) {
			expectedBytes := []byte{0x83, 0xa4, 0x6c, 0x73, 0x69, 0x67, 0x83, 0xa3, 0x61, 0x72, 0x67, 0x92, 0xc4, 0x1, 0x1, 0xc4, 0x2, 0x2, 0x3, 0xa1, 0x6c, 0xc4, 0x5, 0x1, 0x20, 0x1, 0x1, 0x22, 0xa4, 0x6d, 0x73, 0x69, 0x67, 0x83, 0xa6, 0x73, 0x75, 0x62, 0x73, 0x69, 0x67, 0x93, 0x82, 0xa2, 0x70, 0x6b, 0xc4, 0x20, 0x1b, 0x7e, 0xc0, 0xb0, 0x4b, 0xea, 0x61, 0xb7, 0x96, 0x90, 0x97, 0xe6, 0xcb, 0xf4, 0x7, 0xe1, 0x8, 0xa7, 0x5, 0x35, 0x1d, 0xb, 0xc9, 0x8a, 0xbe, 0xb1, 0x22, 0x9, 0xa8, 0xab, 0x81, 0x78, 0xa1, 0x73, 0xc4, 0x40, 0x49, 0x13, 0xb8, 0x5, 0xd1, 0x9e, 0x7f, 0x2c, 0x10, 0x80, 0xf6, 0x33, 0x7e, 0x18, 0x54, 0xa7, 0xce, 0xea, 0xee, 0x10, 0xdd, 0xbd, 0x13, 0x65, 0x84, 0xbf, 0x93, 0xb7, 0x5f, 0x30, 0x63, 0x15, 0x91, 0xca, 0x23, 0xc, 0xed, 0xef, 0x23, 0xd1, 0x74, 0x1b, 0x52, 0x9d, 0xb0, 0xff, 0xef, 0x37, 0x54, 0xd6, 0x46, 0xf4, 0xb5, 0x61, 0xfc, 0x8b, 0xbc, 0x2d, 0x7b, 0x4e, 0x63, 0x5c, 0xbd, 0x2, 0x82, 0xa2, 0x70, 0x6b, 0xc4, 0x20, 0x9, 0x63, 0x32, 0x9, 0x53, 0x73, 0x89, 0xf0, 0x75, 0x67, 0x11, 0x77, 0x39, 0x91, 0xc7, 0xd0, 0x3e, 0x1b, 0x73, 0xc8, 0xc4, 0xf5, 0x2b, 0xf6, 0xaf, 0xf0, 0x1a, 0xa2, 0x5c, 0xf9, 0xc2, 0x71, 0xa1, 0x73, 0xc4, 0x40, 0x64, 0xbc, 0x55, 0xdb, 0xed, 0x91, 0xa2, 0x41, 0xd4, 0x2a, 0xb6, 0x60, 0xf7, 0xe1, 0x4a, 0xb9, 0x99, 0x9a, 0x52, 0xb3, 0xb1, 0x71, 0x58, 0xce, 0xfc, 0x3f, 0x4f, 0xe7, 0xcb, 0x22, 0x41, 0x14, 0xad, 0xa9, 0x3d, 0x5e, 0x84, 0x5, 0x2, 0xa, 0x17, 0xa6, 0x69, 0x83, 0x3, 0x22, 0x4e, 0x86, 0xa3, 0x8b, 0x6a, 0x36, 0xc5, 0x54, 0xbe, 0x20, 0x50, 0xff, 0xd3, 0xee, 0xa8, 0xb3, 0x4, 0x9, 0x81, 0xa2, 0x70, 0x6b, 0xc4, 0x20, 0xe7, 0xf0, 0xf8, 0x4d, 0x6, 0x81, 0x1d, 0xf9, 0xf3, 0x1c, 0x8d, 0x87, 0x8b, 0x11, 0x55, 0xf4, 0x67, 0x1d, 0x51, 0xa1, 0x85, 0xc2, 0x0, 0x90, 0x86, 0x67, 0xf4, 0x49, 0x58, 0x70, 0x68, 0xa1, 0xa3, 0x74, 0x68, 0x72, 0x2, 0xa1, 0x76, 0x1, 0xa4, 0x73, 0x67, 0x6e, 0x72, 0xc4, 0x20, 0x8d, 0x92, 0xb4, 0x89, 0x90, 0x1, 0x73, 0xa0, 0x4d, 0xfa, 0x43, 0x59, 0xa3, 0x66, 0x6a, 0x6a, 0xfc, 0xea, 0x2c, 0x42, 0xa0, 0x5d, 0xd9, 0xc1, 0xf7, 0x3e, 0xeb, 0xa5, 0x47, 0x80, 0x37, 0xe9, 0xa3, 0x74, 0x78, 0x6e, 0x89, 0xa3, 0x61, 0x6d, 0x74, 0xcd, 0x13, 0x88, 0xa3, 0x66, 0x65, 0x65, 0xce, 0x0, 0x3, 0x4f, 0xa8, 0xa2, 0x66, 0x76, 0xce, 0x0, 0xe, 0xd6, 0xdc, 0xa3, 0x67, 0x65, 0x6e, 0xad, 0x74, 0x65, 0x73, 0x74, 0x6e, 0x65, 0x74, 0x2d, 0x76, 0x33, 0x31, 0x2e, 0x30, 0xa2, 0x6c, 0x76, 0xce, 0x0, 0xe, 0xda, 0xc4, 0xa4, 0x6e, 0x6f, 0x74, 0x65, 0xc4, 0x8, 0xb4, 0x51, 0x79, 0x39, 0xfc, 0xfa, 0xd2, 0x71, 0xa3, 0x72, 0x63, 0x76, 0xc4, 0x20, 0xb4, 0xc6, 0x22, 0x68, 0xa4, 0xc3, 0x35, 0x50, 0xf5, 0x8a, 0x6d, 0xd4, 0xb8, 0xca, 0x23, 0x9c, 0x9f, 0x1d, 0x10, 0xb3, 0x41, 0xaa, 0x7e, 0x1a, 0x36, 0xe3, 0xce, 0x61, 0xfb, 0xfc, 0xf9, 0x8a, 0xa3, 0x73, 0x6e, 0x64, 0xc4, 0x20, 0xb4, 0xc6, 0x22, 0x68, 0xa4, 0xc3, 0x35, 0x50, 0xf5, 0x8a, 0x6d, 0xd4, 0xb8, 0xca, 0x23, 0x9c, 0x9f, 0x1d, 0x10, 0xb3, 0x41, 0xaa, 0x7e, 0x1a, 0x36, 0xe3, 0xce, 0x61, 0xfb, 0xfc, 0xf9, 0x8a, 0xa4, 0x74, 0x79, 0x70, 0x65, 0xa3, 0x70, 0x61, 0x79}
			expectedTxid := "U4X24Q45MCZ6JSL343QNR3RC6AJO2NUDQXFCTNONIW3SU2AYQJOA"
			expectedAuthAddr := maAddr
			sender := otherAddr
			testSign(t, lsig, sender, expectedBytes, expectedTxid, expectedAuthAddr)
		})
	})
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
