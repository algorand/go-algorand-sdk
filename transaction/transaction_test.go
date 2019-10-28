package transaction

import (
	"encoding/base64"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
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

func byte32ArrayFromBase64(s string) (out [32]byte) {
	slice := byteFromBase64(s)
	if len(slice) != 32 {
		panic("wrong length: input slice not 32 bytes")
	}
	copy(out[:], slice)
	return
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

func TestMakeAssetConfigTxn(t *testing.T) {
	const addr = "BH55E5RMBD4GYWXGX5W5PJ5JAHPGM5OXKDQH5DC4O2MGI7NW4H6VOE4CP4"
	const genesisHash = "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI="
	const manager = addr
	const reserve = addr
	const freeze = addr
	const clawback = addr
	const assetIndex = 1234
	encoded, err := MakeAssetConfigTxn(addr, 10, 322575, 323575, nil, "", genesisHash,
		assetIndex, manager, reserve, freeze, clawback)
	require.NoError(t, err)

	a, err := types.DecodeAddress(addr)
	require.NoError(t, err)
	expectedAssetConfigTxn := types.Transaction{
		Type: types.AssetConfigTx,
		Header: types.Header{
			Sender:      a,
			Fee:         3400,
			FirstValid:  322575,
			LastValid:   323575,
			GenesisHash: byte32ArrayFromBase64(genesisHash),
			GenesisID:   "",
		},
	}

	expectedAssetConfigTxn.AssetParams = types.AssetParams{
		Manager:  a,
		Reserve:  a,
		Freeze:   a,
		Clawback: a,
	}
	expectedAssetConfigTxn.ConfigAsset = types.AssetIndex(assetIndex)

	var decoded types.Transaction
	err = msgpack.Decode(encoded, &decoded)
	require.NoError(t, err)

	require.Equal(t, expectedAssetConfigTxn, decoded)

	const addrSK = "awful drop leaf tennis indoor begin mandate discover uncle seven only coil atom any hospital uncover make any climb actor armed measure need above hundred"
	private, err := mnemonic.ToPrivateKey(addrSK)
	require.NoError(t, err)
	newStxBytes, err := crypto.SignTransaction(private, encoded)
	signedGolden := "gqNzaWfEQBBkfw5n6UevuIMDo2lHyU4dS80JCCQ/vTRUcTx5m0ivX68zTKyuVRrHaTbxbRRc3YpJ4zeVEnC9Fiw3Wf4REwejdHhuiKRhcGFyhKFjxCAJ+9J2LAj4bFrmv23Xp6kB3mZ111Dgfoxcdphkfbbh/aFmxCAJ+9J2LAj4bFrmv23Xp6kB3mZ111Dgfoxcdphkfbbh/aFtxCAJ+9J2LAj4bFrmv23Xp6kB3mZ111Dgfoxcdphkfbbh/aFyxCAJ+9J2LAj4bFrmv23Xp6kB3mZ111Dgfoxcdphkfbbh/aRjYWlkzQTSo2ZlZc0NSKJmds4ABOwPomdoxCBIY7UYpLPITsgQ8i1PEIHLD3HwWaesIN7GL39w5Qk6IqJsds4ABO/3o3NuZMQgCfvSdiwI+Gxa5r9t16epAd5mdddQ4H6MXHaYZH224f2kdHlwZaRhY2Zn"
	require.EqualValues(t, newStxBytes, byteFromBase64(signedGolden))
}

func TestMakeAssetTransferTxn(t *testing.T) {
	const addrSK = "awful drop leaf tennis indoor begin mandate discover uncle seven only coil atom any hospital uncover make any climb actor armed measure need above hundred"
	private, err := mnemonic.ToPrivateKey(addrSK)
	require.NoError(t, err)

	const addr = "BH55E5RMBD4GYWXGX5W5PJ5JAHPGM5OXKDQH5DC4O2MGI7NW4H6VOE4CP4"
	const genesisHash = "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI="
	const sender, recipient, closeAssetsTo = addr, addr, addr
	const assetIndex = 1
	const firstValidRound = 322575
	const lastValidRound = 323576
	const amountToSend = 1

	encoded, err := MakeAssetTransferTxn(sender, recipient, closeAssetsTo, amountToSend, 10, firstValidRound,
		lastValidRound, nil, "", genesisHash, assetIndex)
	require.NoError(t, err)

	sendAddr, err := types.DecodeAddress(sender)
	require.NoError(t, err)

	expectedAssetTransferTxn := types.Transaction{
		Type: types.AssetTransferTx,
		Header: types.Header{
			Sender:      sendAddr,
			Fee:         2750,
			FirstValid:  firstValidRound,
			LastValid:   lastValidRound,
			GenesisHash: byte32ArrayFromBase64(genesisHash),
			GenesisID:   "",
		},
	}

	expectedAssetID := types.AssetIndex(assetIndex)
	expectedAssetTransferTxn.XferAsset = expectedAssetID

	receiveAddr, err := types.DecodeAddress(recipient)
	require.NoError(t, err)
	expectedAssetTransferTxn.AssetReceiver = receiveAddr

	closeAddr, err := types.DecodeAddress(closeAssetsTo)
	require.NoError(t, err)
	expectedAssetTransferTxn.AssetCloseTo = closeAddr

	expectedAssetTransferTxn.AssetAmount = amountToSend

	var decoded types.Transaction
	err = msgpack.Decode(encoded, &decoded)
	require.NoError(t, err)

	require.Equal(t, expectedAssetTransferTxn, decoded)

	// now compare tx against a golden
	const signedGolden = "gqNzaWfEQNkEs3WdfFq6IQKJdF1n0/hbV9waLsvojy9pM1T4fvwfMNdjGQDy+LeesuQUfQVTneJD4VfMP7zKx4OUlItbrwSjdHhuiqRhYW10AaZhY2xvc2XEIAn70nYsCPhsWua/bdenqQHeZnXXUOB+jFx2mGR9tuH9pGFyY3bEIAn70nYsCPhsWua/bdenqQHeZnXXUOB+jFx2mGR9tuH9o2ZlZc0KvqJmds4ABOwPomdoxCBIY7UYpLPITsgQ8i1PEIHLD3HwWaesIN7GL39w5Qk6IqJsds4ABO/4o3NuZMQgCfvSdiwI+Gxa5r9t16epAd5mdddQ4H6MXHaYZH224f2kdHlwZaVheGZlcqR4YWlkAQ=="
	newStxBytes, err := crypto.SignTransaction(private, encoded)
	require.NoError(t, err)
	require.EqualValues(t, newStxBytes, byteFromBase64(signedGolden))
}

func TestMakeAssetAcceptanceTxn(t *testing.T) {
	const sender = "BH55E5RMBD4GYWXGX5W5PJ5JAHPGM5OXKDQH5DC4O2MGI7NW4H6VOE4CP4"
	const genesisHash = "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI="
	const assetIndex = 1
	const firstValidRound = 322575
	const lastValidRound = 323575

	encoded, err := MakeAssetAcceptanceTxn(sender, 10, firstValidRound,
		lastValidRound, nil, "", genesisHash, assetIndex)
	require.NoError(t, err)

	sendAddr, err := types.DecodeAddress(sender)
	require.NoError(t, err)

	expectedAssetAcceptanceTxn := types.Transaction{
		Type: types.AssetTransferTx,
		Header: types.Header{
			Sender:      sendAddr,
			Fee:         2280,
			FirstValid:  firstValidRound,
			LastValid:   lastValidRound,
			GenesisHash: byte32ArrayFromBase64(genesisHash),
			GenesisID:   "",
		},
	}

	expectedAssetID := types.AssetIndex(assetIndex)
	expectedAssetAcceptanceTxn.XferAsset = expectedAssetID
	expectedAssetAcceptanceTxn.AssetReceiver = sendAddr
	expectedAssetAcceptanceTxn.AssetAmount = 0

	var decoded types.Transaction
	err = msgpack.Decode(encoded, &decoded)
	require.NoError(t, err)

	require.Equal(t, expectedAssetAcceptanceTxn, decoded)

	const addrSK = "awful drop leaf tennis indoor begin mandate discover uncle seven only coil atom any hospital uncover make any climb actor armed measure need above hundred"
	private, err := mnemonic.ToPrivateKey(addrSK)
	require.NoError(t, err)
	newStxBytes, err := crypto.SignTransaction(private, encoded)
	signedGolden := "gqNzaWfEQJ7q2rOT8Sb/wB0F87ld+1zMprxVlYqbUbe+oz0WM63FctIi+K9eYFSqT26XBZ4Rr3+VTJpBE+JLKs8nctl9hgijdHhuiKRhcmN2xCAJ+9J2LAj4bFrmv23Xp6kB3mZ111Dgfoxcdphkfbbh/aNmZWXNCOiiZnbOAATsD6JnaMQgSGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiKibHbOAATv96NzbmTEIAn70nYsCPhsWua/bdenqQHeZnXXUOB+jFx2mGR9tuH9pHR5cGWlYXhmZXKkeGFpZAE="
	require.EqualValues(t, newStxBytes, byteFromBase64(signedGolden))
}
