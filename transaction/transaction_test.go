package transaction

import (
	"encoding/base64"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ed25519"
	"testing"
)

func byteFromBase64(s string) []byte {
	b, _ := base64.StdEncoding.DecodeString(s)
	return b
}

func TestSignTransaction(t *testing.T) {
	// Reference transaction + keys generated using goal
	const fromAddress = "47YPQTIGQEO7T4Y4RWDYWEKV6RTR2UNBQXBABEEGM72ESWDQNCQ52OPASU"
	const toAddress = "PNWOET7LLOWMBMLE4KOCELCX6X3D3Q4H2Q4QJASYIEOF7YIPPQBG3YQ5YI"
	const referenceTxID = "5FJDJD5LMZC3EHUYYJNH5I23U4X6H2KXABNDGPIL557ZMJ33GZHQ"
	const mn = "advice pudding treat near rule blouse same whisper inner electric quit surface sunny dismiss leader blood seat clown cost exist hospital century reform able sponsor"
	const golden = "gqNzaWfEQPhUAZ3xkDDcc8FvOVo6UinzmKBCqs0woYSfodlmBMfQvGbeUx3Srxy3dyJDzv7rLm26BRv9FnL2/AuT7NYfiAWjdHhui6NhbXTNA+ilY2xvc2XEIEDpNJKIJWTLzpxZpptnVCaJ6aHDoqnqW2Wm6KRCH/xXo2ZlZc0EmKJmds0wsqNnZW6sZGV2bmV0LXYzMy4womdoxCAmCyAJoJOohot5WHIvpeVG7eftF+TYXEx4r7BFJpDt0qJsds00mqRub3RlxAjqABVHQ2y/lqNyY3bEIHts4k/rW6zAsWTinCIsV/X2PcOH1DkEglhBHF/hD3wCo3NuZMQg5/D4TQaBHfnzHI2HixFV9GcdUaGFwgCQhmf0SVhwaKGkdHlwZaNwYXk="
	gh := byteFromBase64("JgsgCaCTqIaLeVhyL6XlRu3n7Rfk2FxMeK+wRSaQ7dI=")

	// Build the unsigned transaction
	tx, err := MakePaymentTxn(fromAddress, toAddress, 4, 1000, 12466, 13466, byteFromBase64("6gAVR0Nsv5Y="), "IDUTJEUIEVSMXTU4LGTJWZ2UE2E6TIODUKU6UW3FU3UKIQQ77RLUBBBFLA", "devnet-v33.0", gh)
	require.NoError(t, err)

	// Decode the secret key for the sender
	seed, err := mnemonic.ToKey(mn)
	require.NoError(t, err)
	sk := ed25519.NewKeyFromSeed(seed)

	// Check that we have the correct pk
	derivedFromPK := sk.Public()
	var derivedFromAddr types.Address
	copy(derivedFromAddr[:], derivedFromPK.(ed25519.PublicKey))
	require.Equal(t, fromAddress, derivedFromAddr.String())

	// Sign the transaction
	stxBytes, err := crypto.SignTransaction(sk, tx)
	require.NoError(t, err)

	txid := crypto.GetTxID(tx)
	require.Equal(t, txid, referenceTxID)

	require.Equal(t, byteFromBase64(golden), stxBytes)
}
