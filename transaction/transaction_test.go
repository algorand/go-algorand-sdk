package transaction

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/algorand/go-algorand-sdk/types"
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

func TestKeyRegTxn(t *testing.T) {
	// preKeyRegTxn is an unsigned signed keyreg txn with zero Sender
	const addr = "BH55E5RMBD4GYWXGX5W5PJ5JAHPGM5OXKDQH5DC4O2MGI7NW4H6VOE4CP4"
	a, err := types.DecodeAddress(addr)
	require.NoError(t, err)
	const addrSK = "awful drop leaf tennis indoor begin mandate discover uncle seven only coil atom any hospital uncover make any climb actor armed measure need above hundred"
	expKeyRegTxn := types.Transaction{
		Type: types.KeyRegistrationTx,
		Header: types.Header{
			Sender:      a,
			Fee:         1000,
			FirstValid:  322575,
			LastValid:   323575,
			GenesisHash: byte32ArrayFromBase64("SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI="),
			GenesisID:   "",
		},
		KeyregTxnFields: types.KeyregTxnFields{
			VotePK:          byte32ArrayFromBase64("Kv7QI7chi1y6axoy+t7wzAVpePqRq/rkjzWh/RMYyLo="),
			SelectionPK:     byte32ArrayFromBase64("bPgrv4YogPcdaUAxrt1QysYZTVyRAuUMD4zQmCu9llc="),
			VoteFirst:       10000,
			VoteLast:        10111,
			VoteKeyDilution: 11,
		},
	}
	const signedGolden = "gqNzaWfEQEA8ANbrvTRxU9c8v6WERcEPw7D/HacRgg4vICa61vEof60Wwtx6KJKDyvBuvViFeacLlngPY6vYCVP0DktTwQ2jdHhui6NmZWXNA+iiZnbOAATsD6JnaMQgSGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiKibHbOAATv96ZzZWxrZXnEIGz4K7+GKID3HWlAMa7dUMrGGU1ckQLlDA+M0JgrvZZXo3NuZMQgCfvSdiwI+Gxa5r9t16epAd5mdddQ4H6MXHaYZH224f2kdHlwZaZrZXlyZWendm90ZWZzdM0nEKZ2b3Rla2QLp3ZvdGVrZXnEICr+0CO3IYtcumsaMvre8MwFaXj6kav65I81of0TGMi6p3ZvdGVsc3TNJ38="
	// now, sign
	private, err := mnemonic.ToPrivateKey(addrSK)
	require.NoError(t, err)
	txid, newStxBytes, err := crypto.SignTransaction(private, expKeyRegTxn)
	require.NoError(t, err)
	require.Equal(t, "MDRIUVH5AW4Z3GMOB67WP44LYLEVM2MP3ZEPKFHUB5J47A2J6TUQ", txid)
	require.EqualValues(t, newStxBytes, byteFromBase64(signedGolden))
}

func TestMakeKeyRegTxn(t *testing.T) {
	const addr = "BH55E5RMBD4GYWXGX5W5PJ5JAHPGM5OXKDQH5DC4O2MGI7NW4H6VOE4CP4"
	tx, err := MakeKeyRegTxn(addr, 10, 322575, 323575, []byte{45, 67}, "", "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI=",
		"Kv7QI7chi1y6axoy+t7wzAVpePqRq/rkjzWh/RMYyLo=", "bPgrv4YogPcdaUAxrt1QysYZTVyRAuUMD4zQmCu9llc=", 10000, 10111, 11)
	require.NoError(t, err)

	a, err := types.DecodeAddress(addr)
	require.NoError(t, err)
	expKeyRegTxn := types.Transaction{
		Type: types.KeyRegistrationTx,
		Header: types.Header{
			Sender:      a,
			Fee:         3060,
			FirstValid:  322575,
			LastValid:   323575,
			Note:        []byte{45, 67},
			GenesisHash: byte32ArrayFromBase64("SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI="),
			GenesisID:   "",
		},
		KeyregTxnFields: types.KeyregTxnFields{
			VotePK:          byte32ArrayFromBase64("Kv7QI7chi1y6axoy+t7wzAVpePqRq/rkjzWh/RMYyLo="),
			SelectionPK:     byte32ArrayFromBase64("bPgrv4YogPcdaUAxrt1QysYZTVyRAuUMD4zQmCu9llc="),
			VoteFirst:       10000,
			VoteLast:        10111,
			VoteKeyDilution: 11,
		},
	}
	require.Equal(t, expKeyRegTxn, tx)
}

func TestMakeAssetCreateTxn(t *testing.T) {
	const addr = "BH55E5RMBD4GYWXGX5W5PJ5JAHPGM5OXKDQH5DC4O2MGI7NW4H6VOE4CP4"
	const defaultFrozen = false
	const genesisHash = "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI="
	const total = 100
	const reserve = addr
	const freeze = addr
	const clawback = addr
	const unitName = "tst"
	const assetName = "testcoin"
	tx, err := MakeAssetCreateTxn(addr, 10, 322575, 323575, nil, "", genesisHash,
		total, defaultFrozen, addr, reserve, freeze, clawback, unitName, assetName)
	require.NoError(t, err)

	a, err := types.DecodeAddress(addr)
	require.NoError(t, err)
	expectedAssetCreationTxn := types.Transaction{
		Type: types.AssetConfigTx,
		Header: types.Header{
			Sender:      a,
			Fee:         3890,
			FirstValid:  322575,
			LastValid:   323575,
			GenesisHash: byte32ArrayFromBase64(genesisHash),
			GenesisID:   "",
		},
	}
	expectedAssetCreationTxn.AssetParams = types.AssetParams{
		Total:         total,
		DefaultFrozen: defaultFrozen,
		Manager:       a,
		Reserve:       a,
		Freeze:        a,
		Clawback:      a,
	}
	copy(expectedAssetCreationTxn.AssetParams.UnitName[:], []byte(unitName))
	copy(expectedAssetCreationTxn.AssetParams.AssetName[:], []byte(assetName))
	require.Equal(t, expectedAssetCreationTxn, tx)
}

func TestMakeAssetConfigTxn(t *testing.T) {
	const addr = "BH55E5RMBD4GYWXGX5W5PJ5JAHPGM5OXKDQH5DC4O2MGI7NW4H6VOE4CP4"
	const genesisHash = "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI="
	const creator = addr
	const manager = addr
	const reserve = addr
	const freeze = addr
	const clawback = addr
	tx, err := MakeAssetConfigTxn(addr, 10, 322575, 323575, nil, "", genesisHash,
		creator, 1234, manager, reserve, freeze, clawback)
	require.NoError(t, err)

	a, err := types.DecodeAddress(creator)
	require.NoError(t, err)
	expectedAssetConfigTxn := types.Transaction{
		Type: types.AssetConfigTx,
		Header: types.Header{
			Sender:      a,
			Fee:         3790,
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
	expectedAssetConfigTxn.ConfigAsset = types.AssetID{
		Creator: a,
		Index:   1234,
	}
	require.Equal(t, expectedAssetConfigTxn, tx)
}

func TestMakeAssetDestroyTxn(t *testing.T) {
	const addr = "BH55E5RMBD4GYWXGX5W5PJ5JAHPGM5OXKDQH5DC4O2MGI7NW4H6VOE4CP4"
	const genesisHash = "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI="
	const creator = addr
	const assetIndex = 1
	const firstValidRound = 322575
	const lastValidRound = 323575
	tx, err := MakeAssetDestroyTxn(creator, 10, firstValidRound, lastValidRound, nil, "", genesisHash, creator, assetIndex)
	require.NoError(t, err)

	a, err := types.DecodeAddress(creator)
	require.NoError(t, err)

	expectedAssetConfigTxn := types.Transaction{
		Type: types.AssetConfigTx,
		Header: types.Header{
			Sender:      a,
			Fee:         2290,
			FirstValid:  firstValidRound,
			LastValid:   lastValidRound,
			GenesisHash: byte32ArrayFromBase64(genesisHash),
			GenesisID:   "",
		},
	}
	expectedAssetConfigTxn.AssetParams = types.AssetParams{}
	expectedAssetConfigTxn.ConfigAsset = types.AssetID{
		Creator: a,
		Index:   assetIndex,
	}
	require.Equal(t, expectedAssetConfigTxn, tx)
}

func TestMakeAssetFreezeTxn(t *testing.T) {
	const addr = "BH55E5RMBD4GYWXGX5W5PJ5JAHPGM5OXKDQH5DC4O2MGI7NW4H6VOE4CP4"
	const genesisHash = "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI="
	const creator = addr
	const assetIndex = 1
	const firstValidRound = 322575
	const lastValidRound = 323575
	const freezeSetting = true
	const target = creator
	tx, err := MakeAssetFreezeTxn(creator, 10, firstValidRound, lastValidRound, nil, "", genesisHash, creator, assetIndex, target, freezeSetting)
	require.NoError(t, err)

	a, err := types.DecodeAddress(creator)
	require.NoError(t, err)

	expectedAssetFreezeTxn := types.Transaction{
		Type: types.AssetFreezeTx,
		Header: types.Header{
			Sender:      a,
			Fee:         2720,
			FirstValid:  firstValidRound,
			LastValid:   lastValidRound,
			GenesisHash: byte32ArrayFromBase64(genesisHash),
			GenesisID:   "",
		},
	}
	expectedAssetFreezeTxn.FreezeAsset.Creator = a
	expectedAssetFreezeTxn.FreezeAsset.Index = assetIndex
	expectedAssetFreezeTxn.AssetFrozen = freezeSetting
	expectedAssetFreezeTxn.FreezeAccount = a
	require.Equal(t, expectedAssetFreezeTxn, tx)
}

func TestMakeAssetTransferTxn(t *testing.T) {
	const addr = "BH55E5RMBD4GYWXGX5W5PJ5JAHPGM5OXKDQH5DC4O2MGI7NW4H6VOE4CP4"
	const genesisHash = "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI="
	const sender, recipient, creator, closeAssetsTo = addr, addr, addr, addr
	const assetIndex = 1
	const firstValidRound = 322575
	const lastValidRound = 323576
	const amountToSend = 1

	tx, err := MakeAssetTransferTxn(sender, recipient, closeAssetsTo, amountToSend, 10, firstValidRound,
		lastValidRound, nil, "", genesisHash, creator, assetIndex)
	require.NoError(t, err)

	sendAddr, err := types.DecodeAddress(sender)
	require.NoError(t, err)

	expectedAssetTransferTxn := types.Transaction{
		Type: types.AssetTransferTx,
		Header: types.Header{
			Sender:      sendAddr,
			Fee:         3140,
			FirstValid:  firstValidRound,
			LastValid:   lastValidRound,
			GenesisHash: byte32ArrayFromBase64(genesisHash),
			GenesisID:   "",
		},
	}

	creatorAddr, err := types.DecodeAddress(creator)
	require.NoError(t, err)

	expectedAssetID := types.AssetID{
		Creator: creatorAddr,
		Index:   assetIndex,
	}
	expectedAssetTransferTxn.XferAsset = expectedAssetID

	receiveAddr, err := types.DecodeAddress(recipient)
	require.NoError(t, err)
	expectedAssetTransferTxn.AssetReceiver = receiveAddr

	closeAddr, err := types.DecodeAddress(closeAssetsTo)
	require.NoError(t, err)
	expectedAssetTransferTxn.AssetCloseTo = closeAddr

	expectedAssetTransferTxn.AssetAmount = amountToSend

	require.Equal(t, expectedAssetTransferTxn, tx)

	// now compare tx against a golden
	const addrSK = "awful drop leaf tennis indoor begin mandate discover uncle seven only coil atom any hospital uncover make any climb actor armed measure need above hundred"
	const signedGolden = "gqNzaWfEQGkk9CtvOKnn4nU59xmPGoZvYv+6TCu5B95PgwQ/YytwE9dr199ehEqAnSS0C2SaO4YhEBAk+JVOiwZiRq/w1gijdHhuiqRhYW10AaZhY2xvc2XEIAn70nYsCPhsWua/bdenqQHeZnXXUOB+jFx2mGR9tuH9pGFyY3bEIAn70nYsCPhsWua/bdenqQHeZnXXUOB+jFx2mGR9tuH9o2ZlZc0MRKJmds4ABOwPomdoxCBIY7UYpLPITsgQ8i1PEIHLD3HwWaesIN7GL39w5Qk6IqJsds4ABO/4o3NuZMQgCfvSdiwI+Gxa5r9t16epAd5mdddQ4H6MXHaYZH224f2kdHlwZaVheGZlcqR4YWlkgqFjxCAJ+9J2LAj4bFrmv23Xp6kB3mZ111Dgfoxcdphkfbbh/aFpAQ=="
	private, err := mnemonic.ToPrivateKey(addrSK)
	require.NoError(t, err)
	_, newStxBytes, err := crypto.SignTransaction(private, tx)
	require.NoError(t, err)
	require.EqualValues(t, newStxBytes, byteFromBase64(signedGolden))
}

func TestMakeAssetAcceptanceTxn(t *testing.T) {
	const addr = "BH55E5RMBD4GYWXGX5W5PJ5JAHPGM5OXKDQH5DC4O2MGI7NW4H6VOE4CP4"
	const genesisHash = "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI="
	const sender, creator = addr, addr
	const assetIndex = 1
	const firstValidRound = 322575
	const lastValidRound = 323575
	const amountToSend = 1

	tx, err := MakeAssetAcceptanceTxn(sender, 10, firstValidRound,
		lastValidRound, nil, "", genesisHash, creator, assetIndex)
	require.NoError(t, err)

	sendAddr, err := types.DecodeAddress(sender)
	require.NoError(t, err)

	expectedAssetAcceptanceTxn := types.Transaction{
		Type: types.AssetTransferTx,
		Header: types.Header{
			Sender:      sendAddr,
			Fee:         2670,
			FirstValid:  firstValidRound,
			LastValid:   lastValidRound,
			GenesisHash: byte32ArrayFromBase64(genesisHash),
			GenesisID:   "",
		},
	}

	creatorAddr, err := types.DecodeAddress(creator)
	require.NoError(t, err)

	expectedAssetID := types.AssetID{
		Creator: creatorAddr,
		Index:   assetIndex,
	}
	expectedAssetAcceptanceTxn.XferAsset = expectedAssetID
	expectedAssetAcceptanceTxn.AssetReceiver = sendAddr
	expectedAssetAcceptanceTxn.AssetAmount = 0

	require.Equal(t, expectedAssetAcceptanceTxn, tx)
}

func TestMakeAssetRevocationTransaction(t *testing.T) {
	const addr = "BH55E5RMBD4GYWXGX5W5PJ5JAHPGM5OXKDQH5DC4O2MGI7NW4H6VOE4CP4"
	const genesisHash = "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI="
	const revoker, recipient, revoked, creator = addr, addr, addr, addr
	const assetIndex = 1
	const firstValidRound = 322575
	const lastValidRound = 323575
	const amountToSend = 1

	tx, err := MakeAssetRevocationTxn(revoker, revoked, recipient, amountToSend, 10, firstValidRound,
		lastValidRound, nil, "", genesisHash, creator, assetIndex)
	require.NoError(t, err)

	sendAddr, err := types.DecodeAddress(revoker)
	require.NoError(t, err)

	expectedAssetRevocationTxn := types.Transaction{
		Type: types.AssetTransferTx,
		Header: types.Header{
			Sender:      sendAddr,
			Fee:         3140,
			FirstValid:  firstValidRound,
			LastValid:   lastValidRound,
			GenesisHash: byte32ArrayFromBase64(genesisHash),
			GenesisID:   "",
		},
	}

	creatorAddr, err := types.DecodeAddress(creator)
	require.NoError(t, err)

	expectedAssetID := types.AssetID{
		Creator: creatorAddr,
		Index:   assetIndex,
	}
	expectedAssetRevocationTxn.XferAsset = expectedAssetID

	receiveAddr, err := types.DecodeAddress(recipient)
	require.NoError(t, err)
	expectedAssetRevocationTxn.AssetReceiver = receiveAddr

	expectedAssetRevocationTxn.AssetAmount = amountToSend

	targetAddr, err := types.DecodeAddress(revoked)
	require.NoError(t, err)
	expectedAssetRevocationTxn.AssetSender = targetAddr

	require.Equal(t, expectedAssetRevocationTxn, tx)
}

func TestMakeAssetMintTransaction(t *testing.T) {
	const addr = "BH55E5RMBD4GYWXGX5W5PJ5JAHPGM5OXKDQH5DC4O2MGI7NW4H6VOE4CP4"
	const genesisHash = "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI="
	const reserve, recipient, creator = addr, addr, addr
	const assetIndex = 1
	const firstValidRound = 322575
	const lastValidRound = 323575
	const amountToSend = 1

	tx, err := MakeAssetMintTxn(reserve, recipient, amountToSend, 10, firstValidRound,
		lastValidRound, nil, "", genesisHash, creator, assetIndex)
	require.NoError(t, err)

	sendAddr, err := types.DecodeAddress(reserve)
	require.NoError(t, err)

	expectedAssetMintTxn := types.Transaction{
		Type: types.AssetTransferTx,
		Header: types.Header{
			Sender:      sendAddr,
			Fee:         2730,
			FirstValid:  firstValidRound,
			LastValid:   lastValidRound,
			GenesisHash: byte32ArrayFromBase64(genesisHash),
			GenesisID:   "",
		},
	}

	creatorAddr, err := types.DecodeAddress(creator)
	require.NoError(t, err)

	expectedAssetID := types.AssetID{
		Creator: creatorAddr,
		Index:   assetIndex,
	}
	expectedAssetMintTxn.XferAsset = expectedAssetID

	receiveAddr, err := types.DecodeAddress(recipient)
	require.NoError(t, err)
	expectedAssetMintTxn.AssetReceiver = receiveAddr

	expectedAssetMintTxn.AssetAmount = amountToSend

	require.Equal(t, expectedAssetMintTxn, tx)
}

func TestMakeAssetBurnTransaction(t *testing.T) {
	const addr = "BH55E5RMBD4GYWXGX5W5PJ5JAHPGM5OXKDQH5DC4O2MGI7NW4H6VOE4CP4"
	const genesisHash = "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI="
	const reserve, sender, creator = addr, addr, addr
	const assetIndex = 1
	const firstValidRound = 322575
	const lastValidRound = 323575
	const amountToSend = 1

	tx, err := MakeAssetBurnTxn(sender, reserve, amountToSend, 10, firstValidRound,
		lastValidRound, nil, "", genesisHash, creator, assetIndex)
	require.NoError(t, err)

	sendAddr, err := types.DecodeAddress(sender)
	require.NoError(t, err)

	expectedAssetBurnTxn := types.Transaction{
		Type: types.AssetTransferTx,
		Header: types.Header{
			Sender:      sendAddr,
			Fee:         2730,
			FirstValid:  firstValidRound,
			LastValid:   lastValidRound,
			GenesisHash: byte32ArrayFromBase64(genesisHash),
			GenesisID:   "",
		},
	}

	creatorAddr, err := types.DecodeAddress(creator)
	require.NoError(t, err)

	expectedAssetID := types.AssetID{
		Creator: creatorAddr,
		Index:   assetIndex,
	}
	expectedAssetBurnTxn.XferAsset = expectedAssetID

	receiveAddr, err := types.DecodeAddress(reserve)
	require.NoError(t, err)
	expectedAssetBurnTxn.AssetReceiver = receiveAddr

	expectedAssetBurnTxn.AssetAmount = amountToSend

	require.Equal(t, expectedAssetBurnTxn, tx)
}

func TestComputeGroupID(t *testing.T) {
	// compare regular transactions created in SDK with 'goal clerk send' result
	// compare transaction group created in SDK with 'goal clerk group' result
	const address = "UPYAFLHSIPMJOHVXU2MPLQ46GXJKSDCEMZ6RLCQ7GWB5PRDKJUWKKXECXI"
	const fromAddress, toAddress = address, address
	const fee = 1000
	const amount = 2000
	const genesisID = "devnet-v1.0"
	genesisHash := byteFromBase64("sC3P7e2SdbqKJK0tbiCdK9tdSpbe6XeCGKdoNzmlj0E=")

	const firstRound1 = 710399
	note1 := byteFromBase64("wRKw5cJ0CMo=")
	tx1, err := MakePaymentTxnWithFlatFee(
		fromAddress, toAddress, fee, amount, firstRound1, firstRound1+1000,
		note1, "", genesisID, genesisHash,
	)
	require.NoError(t, err)

	const firstRound2 = 710515
	note2 := byteFromBase64("dBlHI6BdrIg=")
	tx2, err := MakePaymentTxnWithFlatFee(
		fromAddress, toAddress, fee, amount, firstRound2, firstRound2+1000,
		note2, "", genesisID, genesisHash,
	)
	require.NoError(t, err)

	const goldenTx1 = "gaN0eG6Ko2FtdM0H0KNmZWXNA+iiZnbOAArW/6NnZW6rZGV2bmV0LXYxLjCiZ2jEILAtz+3tknW6iiStLW4gnSvbXUqW3ul3ghinaDc5pY9Bomx2zgAK2uekbm90ZcQIwRKw5cJ0CMqjcmN2xCCj8AKs8kPYlx63ppj1w5410qkMRGZ9FYofNYPXxGpNLKNzbmTEIKPwAqzyQ9iXHremmPXDnjXSqQxEZn0Vih81g9fEak0spHR5cGWjcGF5"
	const goldenTx2 = "gaN0eG6Ko2FtdM0H0KNmZWXNA+iiZnbOAArXc6NnZW6rZGV2bmV0LXYxLjCiZ2jEILAtz+3tknW6iiStLW4gnSvbXUqW3ul3ghinaDc5pY9Bomx2zgAK21ukbm90ZcQIdBlHI6BdrIijcmN2xCCj8AKs8kPYlx63ppj1w5410qkMRGZ9FYofNYPXxGpNLKNzbmTEIKPwAqzyQ9iXHremmPXDnjXSqQxEZn0Vih81g9fEak0spHR5cGWjcGF5"

	// goal clerk send dumps unsigned transaction as signed with empty signature in order to save tx type
	stx1 := types.SignedTxn{Sig: types.Signature{}, Msig: types.MultisigSig{}, Txn: tx1}
	stx2 := types.SignedTxn{Sig: types.Signature{}, Msig: types.MultisigSig{}, Txn: tx2}
	require.Equal(t, byteFromBase64(goldenTx1), msgpack.Encode(stx1))
	require.Equal(t, byteFromBase64(goldenTx2), msgpack.Encode(stx2))

	gid, err := crypto.ComputeGroupID([]types.Transaction{tx1, tx2})

	// goal clerk group sets Group to every transaction and concatenate them in output file
	// simulating that behavior here
	stx1.Txn.Group = gid
	stx2.Txn.Group = gid

	var txg []byte
	txg = append(txg, msgpack.Encode(stx1)...)
	txg = append(txg, msgpack.Encode(stx2)...)

	const goldenTxg = "gaN0eG6Lo2FtdM0H0KNmZWXNA+iiZnbOAArW/6NnZW6rZGV2bmV0LXYxLjCiZ2jEILAtz+3tknW6iiStLW4gnSvbXUqW3ul3ghinaDc5pY9Bo2dycMQgLiQ9OBup9H/bZLSfQUH2S6iHUM6FQ3PLuv9FNKyt09SibHbOAAra56Rub3RlxAjBErDlwnQIyqNyY3bEIKPwAqzyQ9iXHremmPXDnjXSqQxEZn0Vih81g9fEak0so3NuZMQgo/ACrPJD2Jcet6aY9cOeNdKpDERmfRWKHzWD18RqTSykdHlwZaNwYXmBo3R4boujYW10zQfQo2ZlZc0D6KJmds4ACtdzo2dlbqtkZXZuZXQtdjEuMKJnaMQgsC3P7e2SdbqKJK0tbiCdK9tdSpbe6XeCGKdoNzmlj0GjZ3JwxCAuJD04G6n0f9tktJ9BQfZLqIdQzoVDc8u6/0U0rK3T1KJsds4ACttbpG5vdGXECHQZRyOgXayIo3JjdsQgo/ACrPJD2Jcet6aY9cOeNdKpDERmfRWKHzWD18RqTSyjc25kxCCj8AKs8kPYlx63ppj1w5410qkMRGZ9FYofNYPXxGpNLKR0eXBlo3BheQ=="

	require.Equal(t, byteFromBase64(goldenTxg), txg)

	// check AssignGroupID, do not validate correctness of Group field calculation
	result, err := AssignGroupID([]types.Transaction{tx1, tx2}, "BH55E5RMBD4GYWXGX5W5PJ5JAHPGM5OXKDQH5DC4O2MGI7NW4H6VOE4CP4")
	require.NoError(t, err)
	require.Equal(t, 0, len(result))

	result, err = AssignGroupID([]types.Transaction{tx1, tx2}, address)
	require.NoError(t, err)
	require.Equal(t, 2, len(result))

	result, err = AssignGroupID([]types.Transaction{tx1, tx2}, "")
	require.NoError(t, err)
	require.Equal(t, 2, len(result))
}
