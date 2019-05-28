package transaction

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/mnemonic"
)

func byteFromBase64(s string) []byte {
	b, _ := base64.StdEncoding.DecodeString(s)
	return b
}

func TestMakePaymentTxn(t *testing.T) {
	const fromAddress = "47YPQTIGQEO7T4Y4RWDYWEKV6RTR2UNBQXBABEEGM72ESWDQNCQ52OPASU"
	const toAddress = "PNWOET7LLOWMBMLE4KOCELCX6X3D3Q4H2Q4QJASYIEOF7YIPPQBG3YQ5YI"
	const referenceTxID = "5FJDJD5LMZC3EHUYYJNH5I23U4X6H2KXABNDGPIL557ZMJ33GZHQ"
	const mn = "advice pudding treat near rule blouse same whisper inner electric quit surface sunny dismiss leader blood seat clown cost exist hospital century reform able sponsor"
	const golden = "gqNzaWfEQPhUAZ3xkDDcc8FvOVo6UinzmKBCqs0woYSfodlmBMfQvGbeUx3Srxy3dyJDzv7rLm26BRv9FnL2/AuT7NYfiAWjdHhui6NhbXTNA+ilY2xvc2XEIEDpNJKIJWTLzpxZpptnVCaJ6aHDoqnqW2Wm6KRCH/xXo2ZlZc0EmKJmds0wsqNnZW6sZGV2bmV0LXYzMy4womdoxCAmCyAJoJOohot5WHIvpeVG7eftF+TYXEx4r7BFJpDt0qJsds00mqRub3RlxAjqABVHQ2y/lqNyY3bEIHts4k/rW6zAsWTinCIsV/X2PcOH1DkEglhBHF/hD3wCo3NuZMQg5/D4TQaBHfnzHI2HixFV9GcdUaGFwgCQhmf0SVhwaKGkdHlwZaNwYXk="
	gh := byteFromBase64("JgsgCaCTqIaLeVhyL6XlRu3n7Rfk2FxMeK+wRSaQ7dI=")

	txn, err := MakePaymentTxn(fromAddress, toAddress, 4, 1000, 12466, 13466, byteFromBase64("6gAVR0Nsv5Y="), "IDUTJEUIEVSMXTU4LGTJWZ2UE2E6TIODUKU6UW3FU3UKIQQ77RLUBBBFLA", "devnet-v33.0", gh)
	require.NoError(t, err)

	key, err := mnemonic.ToPrivateKey(mn)
	require.NoError(t, err)

	id, bytes, err := crypto.SignTransaction(key, txn)

	stxBytes := byteFromBase64(golden)
	require.Equal(t, stxBytes, bytes)

	require.Equal(t, referenceTxID, id)
}

// should fail on a lack of GenesisHash
func TestMakePaymentTxn2(t *testing.T) {
	const fromAddress = "47YPQTIGQEO7T4Y4RWDYWEKV6RTR2UNBQXBABEEGM72ESWDQNCQ52OPASU"
	const toAddress = "PNWOET7LLOWMBMLE4KOCELCX6X3D3Q4H2Q4QJASYIEOF7YIPPQBG3YQ5YI"

	_, err := MakePaymentTxn(fromAddress, toAddress, 4, 1000, 12466, 13466, byteFromBase64("6gAVR0Nsv5Y="), "IDUTJEUIEVSMXTU4LGTJWZ2UE2E6TIODUKU6UW3FU3UKIQQ77RLUBBBFLA", "devnet-v33.0", []byte{})
	require.Error(t, err)

}
