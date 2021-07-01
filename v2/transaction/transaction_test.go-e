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

func TestMakePaymentTxnWithLease(t *testing.T) {
	const fromAddress = "47YPQTIGQEO7T4Y4RWDYWEKV6RTR2UNBQXBABEEGM72ESWDQNCQ52OPASU"
	const toAddress = "PNWOET7LLOWMBMLE4KOCELCX6X3D3Q4H2Q4QJASYIEOF7YIPPQBG3YQ5YI"
	const referenceTxID = "7BG6COBZKF6I6W5XY72ZE4HXV6LLZ6ENSR6DASEGSTXYXR4XJOOQ"
	const mn = "advice pudding treat near rule blouse same whisper inner electric quit surface sunny dismiss leader blood seat clown cost exist hospital century reform able sponsor"
	const golden = "gqNzaWfEQOMmFSIKsZvpW0txwzhmbgQjxv6IyN7BbV5sZ2aNgFbVcrWUnqPpQQxfPhV/wdu9jzEPUU1jAujYtcNCxJ7ONgejdHhujKNhbXTNA+ilY2xvc2XEIEDpNJKIJWTLzpxZpptnVCaJ6aHDoqnqW2Wm6KRCH/xXo2ZlZc0FLKJmds0wsqNnZW6sZGV2bmV0LXYzMy4womdoxCAmCyAJoJOohot5WHIvpeVG7eftF+TYXEx4r7BFJpDt0qJsds00mqJseMQgAQIDBAECAwQBAgMEAQIDBAECAwQBAgMEAQIDBAECAwSkbm90ZcQI6gAVR0Nsv5ajcmN2xCB7bOJP61uswLFk4pwiLFf19j3Dh9Q5BIJYQRxf4Q98AqNzbmTEIOfw+E0GgR358xyNh4sRVfRnHVGhhcIAkIZn9ElYcGihpHR5cGWjcGF5"
	gh := byteFromBase64("JgsgCaCTqIaLeVhyL6XlRu3n7Rfk2FxMeK+wRSaQ7dI=")

	lease := [32]byte{1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4}
	txn, err := MakePaymentTxn(fromAddress, toAddress, 4, 1000, 12466, 13466, byteFromBase64("6gAVR0Nsv5Y="), "IDUTJEUIEVSMXTU4LGTJWZ2UE2E6TIODUKU6UW3FU3UKIQQ77RLUBBBFLA", "devnet-v33.0", gh)
	require.NoError(t, err)
	txn.AddLease(lease, 4)
	require.NoError(t, err)

	key, err := mnemonic.ToPrivateKey(mn)
	require.NoError(t, err)

	id, stxBytes, err := crypto.SignTransaction(key, txn)

	goldenBytes := byteFromBase64(golden)
	require.Equal(t, goldenBytes, stxBytes)
	require.Equal(t, referenceTxID, id)
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
	const testURL = "websitewebsitewebsitewebsitewebsitewebsitewebsitewebsitewebsitewebsitewebsitewebsitewebsitewebsi" //96 characters
	const metadataHash = "fACPO4nRgO55j1ndAK3W6Sgc4APkcyFh"
	tx, err := MakeAssetCreateTxn(addr, 10, 322575, 323575, nil, "", genesisHash, total, 0, defaultFrozen, addr, reserve, freeze, clawback, unitName, assetName, testURL, metadataHash)
	require.NoError(t, err)

	a, err := types.DecodeAddress(addr)
	require.NoError(t, err)
	expectedAssetCreationTxn := types.Transaction{
		Type: types.AssetConfigTx,
		Header: types.Header{
			Sender:      a,
			Fee:         4920,
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
		UnitName:      unitName,
		AssetName:     assetName,
		URL:           testURL,
	}
	copy(expectedAssetCreationTxn.AssetParams.MetadataHash[:], []byte(metadataHash))
	require.Equal(t, expectedAssetCreationTxn, tx)

	const addrSK = "awful drop leaf tennis indoor begin mandate discover uncle seven only coil atom any hospital uncover make any climb actor armed measure need above hundred"
	private, err := mnemonic.ToPrivateKey(addrSK)
	require.NoError(t, err)
	_, newStxBytes, err := crypto.SignTransaction(private, tx)
	signedGolden := "gqNzaWfEQCPNXED15sdMwiSxl9K4LALy5brxBpnOtBurPK2TAtjXpxKQqN9flVtjxP0GGYVUjNx/Sez0W4v1CSU9FhYSUgOjdHhuh6RhcGFyiaJhbcQgZkFDUE80blJnTzU1ajFuZEFLM1c2U2djNEFQa2N5RmiiYW6odGVzdGNvaW6iYXXZYHdlYnNpdGV3ZWJzaXRld2Vic2l0ZXdlYnNpdGV3ZWJzaXRld2Vic2l0ZXdlYnNpdGV3ZWJzaXRld2Vic2l0ZXdlYnNpdGV3ZWJzaXRld2Vic2l0ZXdlYnNpdGV3ZWJzaaFjxCAJ+9J2LAj4bFrmv23Xp6kB3mZ111Dgfoxcdphkfbbh/aFmxCAJ+9J2LAj4bFrmv23Xp6kB3mZ111Dgfoxcdphkfbbh/aFtxCAJ+9J2LAj4bFrmv23Xp6kB3mZ111Dgfoxcdphkfbbh/aFyxCAJ+9J2LAj4bFrmv23Xp6kB3mZ111Dgfoxcdphkfbbh/aF0ZKJ1bqN0c3SjZmVlzRM4omZ2zgAE7A+iZ2jEIEhjtRiks8hOyBDyLU8QgcsPcfBZp6wg3sYvf3DlCToiomx2zgAE7/ejc25kxCAJ+9J2LAj4bFrmv23Xp6kB3mZ111Dgfoxcdphkfbbh/aR0eXBlpGFjZmc="
	require.EqualValues(t, newStxBytes, byteFromBase64(signedGolden))
}

func TestMakeAssetCreateTxnWithDecimals(t *testing.T) {
	const addr = "BH55E5RMBD4GYWXGX5W5PJ5JAHPGM5OXKDQH5DC4O2MGI7NW4H6VOE4CP4"
	const defaultFrozen = false
	const genesisHash = "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI="
	const total = 100
	const decimals = 1
	const reserve = addr
	const freeze = addr
	const clawback = addr
	const unitName = "tst"
	const assetName = "testcoin"
	const testURL = "website"
	const metadataHash = "fACPO4nRgO55j1ndAK3W6Sgc4APkcyFh"
	tx, err := MakeAssetCreateTxn(addr, 10, 322575, 323575, nil, "", genesisHash, total, decimals, defaultFrozen, addr, reserve, freeze, clawback, unitName, assetName, testURL, metadataHash)
	require.NoError(t, err)

	a, err := types.DecodeAddress(addr)
	require.NoError(t, err)
	expectedAssetCreationTxn := types.Transaction{
		Type: types.AssetConfigTx,
		Header: types.Header{
			Sender:      a,
			Fee:         4060,
			FirstValid:  322575,
			LastValid:   323575,
			GenesisHash: byte32ArrayFromBase64(genesisHash),
			GenesisID:   "",
		},
	}
	expectedAssetCreationTxn.AssetParams = types.AssetParams{
		Total:         total,
		Decimals:      decimals,
		DefaultFrozen: defaultFrozen,
		Manager:       a,
		Reserve:       a,
		Freeze:        a,
		Clawback:      a,
		UnitName:      unitName,
		AssetName:     assetName,
		URL:           testURL,
	}
	copy(expectedAssetCreationTxn.AssetParams.MetadataHash[:], []byte(metadataHash))
	require.Equal(t, expectedAssetCreationTxn, tx)

	const addrSK = "awful drop leaf tennis indoor begin mandate discover uncle seven only coil atom any hospital uncover make any climb actor armed measure need above hundred"
	private, err := mnemonic.ToPrivateKey(addrSK)
	require.NoError(t, err)
	_, newStxBytes, err := crypto.SignTransaction(private, tx)
	signedGolden := "gqNzaWfEQCj5xLqNozR5ahB+LNBlTG+d0gl0vWBrGdAXj1ibsCkvAwOsXs5KHZK1YdLgkdJecQiWm4oiZ+pm5Yg0m3KFqgqjdHhuh6RhcGFyiqJhbcQgZkFDUE80blJnTzU1ajFuZEFLM1c2U2djNEFQa2N5RmiiYW6odGVzdGNvaW6iYXWnd2Vic2l0ZaFjxCAJ+9J2LAj4bFrmv23Xp6kB3mZ111Dgfoxcdphkfbbh/aJkYwGhZsQgCfvSdiwI+Gxa5r9t16epAd5mdddQ4H6MXHaYZH224f2hbcQgCfvSdiwI+Gxa5r9t16epAd5mdddQ4H6MXHaYZH224f2hcsQgCfvSdiwI+Gxa5r9t16epAd5mdddQ4H6MXHaYZH224f2hdGSidW6jdHN0o2ZlZc0P3KJmds4ABOwPomdoxCBIY7UYpLPITsgQ8i1PEIHLD3HwWaesIN7GL39w5Qk6IqJsds4ABO/3o3NuZMQgCfvSdiwI+Gxa5r9t16epAd5mdddQ4H6MXHaYZH224f2kdHlwZaRhY2Zn"
	require.EqualValues(t, newStxBytes, byteFromBase64(signedGolden))
}

func TestMakeAssetConfigTxn(t *testing.T) {
	const addr = "BH55E5RMBD4GYWXGX5W5PJ5JAHPGM5OXKDQH5DC4O2MGI7NW4H6VOE4CP4"
	const genesisHash = "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI="
	const manager = addr
	const reserve = addr
	const freeze = addr
	const clawback = addr
	const assetIndex = 1234
	tx, err := MakeAssetConfigTxn(addr, 10, 322575, 323575, nil, "", genesisHash,
		assetIndex, manager, reserve, freeze, clawback, false)
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
	require.Equal(t, expectedAssetConfigTxn, tx)

	const addrSK = "awful drop leaf tennis indoor begin mandate discover uncle seven only coil atom any hospital uncover make any climb actor armed measure need above hundred"
	private, err := mnemonic.ToPrivateKey(addrSK)
	require.NoError(t, err)
	_, newStxBytes, err := crypto.SignTransaction(private, tx)
	signedGolden := "gqNzaWfEQBBkfw5n6UevuIMDo2lHyU4dS80JCCQ/vTRUcTx5m0ivX68zTKyuVRrHaTbxbRRc3YpJ4zeVEnC9Fiw3Wf4REwejdHhuiKRhcGFyhKFjxCAJ+9J2LAj4bFrmv23Xp6kB3mZ111Dgfoxcdphkfbbh/aFmxCAJ+9J2LAj4bFrmv23Xp6kB3mZ111Dgfoxcdphkfbbh/aFtxCAJ+9J2LAj4bFrmv23Xp6kB3mZ111Dgfoxcdphkfbbh/aFyxCAJ+9J2LAj4bFrmv23Xp6kB3mZ111Dgfoxcdphkfbbh/aRjYWlkzQTSo2ZlZc0NSKJmds4ABOwPomdoxCBIY7UYpLPITsgQ8i1PEIHLD3HwWaesIN7GL39w5Qk6IqJsds4ABO/3o3NuZMQgCfvSdiwI+Gxa5r9t16epAd5mdddQ4H6MXHaYZH224f2kdHlwZaRhY2Zn"
	require.EqualValues(t, newStxBytes, byteFromBase64(signedGolden))
}

func TestMakeAssetConfigTxnStrictChecking(t *testing.T) {
	const addr = "BH55E5RMBD4GYWXGX5W5PJ5JAHPGM5OXKDQH5DC4O2MGI7NW4H6VOE4CP4"
	const genesisHash = "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI="
	const manager = addr
	const reserve = addr
	const freeze = ""
	const clawback = addr
	const assetIndex = 1234
	_, err := MakeAssetConfigTxn(addr, 10, 322575, 323575, nil, "", genesisHash,
		assetIndex, manager, reserve, freeze, clawback, true)
	require.Error(t, err)
}

func TestMakeAssetDestroyTxn(t *testing.T) {
	const addr = "BH55E5RMBD4GYWXGX5W5PJ5JAHPGM5OXKDQH5DC4O2MGI7NW4H6VOE4CP4"
	const genesisHash = "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI="
	const creator = addr
	const assetIndex = 1
	const firstValidRound = 322575
	const lastValidRound = 323575
	tx, err := MakeAssetDestroyTxn(creator, 10, firstValidRound, lastValidRound, nil, "", genesisHash, assetIndex)
	require.NoError(t, err)

	a, err := types.DecodeAddress(creator)
	require.NoError(t, err)

	expectedAssetDestroyTxn := types.Transaction{
		Type: types.AssetConfigTx,
		Header: types.Header{
			Sender:      a,
			Fee:         1880,
			FirstValid:  firstValidRound,
			LastValid:   lastValidRound,
			GenesisHash: byte32ArrayFromBase64(genesisHash),
			GenesisID:   "",
		},
	}
	expectedAssetDestroyTxn.AssetParams = types.AssetParams{}
	expectedAssetDestroyTxn.ConfigAsset = types.AssetIndex(assetIndex)
	require.Equal(t, expectedAssetDestroyTxn, tx)

	const addrSK = "awful drop leaf tennis indoor begin mandate discover uncle seven only coil atom any hospital uncover make any climb actor armed measure need above hundred"
	private, err := mnemonic.ToPrivateKey(addrSK)
	require.NoError(t, err)
	_, newStxBytes, err := crypto.SignTransaction(private, tx)
	signedGolden := "gqNzaWfEQBSP7HtzD/Lvn4aVvaNpeR4T93dQgo4LvywEwcZgDEoc/WVl3aKsZGcZkcRFoiWk8AidhfOZzZYutckkccB8RgGjdHhuh6RjYWlkAaNmZWXNB1iiZnbOAATsD6JnaMQgSGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiKibHbOAATv96NzbmTEIAn70nYsCPhsWua/bdenqQHeZnXXUOB+jFx2mGR9tuH9pHR5cGWkYWNmZw=="
	require.EqualValues(t, newStxBytes, byteFromBase64(signedGolden))
}

func TestMakeAssetFreezeTxn(t *testing.T) {
	const addr = "BH55E5RMBD4GYWXGX5W5PJ5JAHPGM5OXKDQH5DC4O2MGI7NW4H6VOE4CP4"
	const genesisHash = "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI="
	const assetIndex = 1
	const firstValidRound = 322575
	const lastValidRound = 323576
	const freezeSetting = true
	const target = addr
	tx, err := MakeAssetFreezeTxn(addr, 10, firstValidRound, lastValidRound, nil, "", genesisHash, assetIndex, target, freezeSetting)
	require.NoError(t, err)

	a, err := types.DecodeAddress(addr)
	require.NoError(t, err)

	expectedAssetFreezeTxn := types.Transaction{
		Type: types.AssetFreezeTx,
		Header: types.Header{
			Sender:      a,
			Fee:         2330,
			FirstValid:  firstValidRound,
			LastValid:   lastValidRound,
			GenesisHash: byte32ArrayFromBase64(genesisHash),
			GenesisID:   "",
		},
	}
	expectedAssetFreezeTxn.FreezeAsset = types.AssetIndex(assetIndex)
	expectedAssetFreezeTxn.AssetFrozen = freezeSetting
	expectedAssetFreezeTxn.FreezeAccount = a
	require.Equal(t, expectedAssetFreezeTxn, tx)

	const addrSK = "awful drop leaf tennis indoor begin mandate discover uncle seven only coil atom any hospital uncover make any climb actor armed measure need above hundred"
	private, err := mnemonic.ToPrivateKey(addrSK)
	require.NoError(t, err)
	_, newStxBytes, err := crypto.SignTransaction(private, tx)
	signedGolden := "gqNzaWfEQAhru5V2Xvr19s4pGnI0aslqwY4lA2skzpYtDTAN9DKSH5+qsfQQhm4oq+9VHVj7e1rQC49S28vQZmzDTVnYDQGjdHhuiaRhZnJ6w6RmYWRkxCAJ+9J2LAj4bFrmv23Xp6kB3mZ111Dgfoxcdphkfbbh/aRmYWlkAaNmZWXNCRqiZnbOAATsD6JnaMQgSGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiKibHbOAATv+KNzbmTEIAn70nYsCPhsWua/bdenqQHeZnXXUOB+jFx2mGR9tuH9pHR5cGWkYWZyeg=="
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

	tx, err := MakeAssetTransferTxn(sender, recipient, closeAssetsTo, amountToSend, 10, firstValidRound,
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

	require.Equal(t, expectedAssetTransferTxn, tx)

	// now compare tx against a golden
	const signedGolden = "gqNzaWfEQNkEs3WdfFq6IQKJdF1n0/hbV9waLsvojy9pM1T4fvwfMNdjGQDy+LeesuQUfQVTneJD4VfMP7zKx4OUlItbrwSjdHhuiqRhYW10AaZhY2xvc2XEIAn70nYsCPhsWua/bdenqQHeZnXXUOB+jFx2mGR9tuH9pGFyY3bEIAn70nYsCPhsWua/bdenqQHeZnXXUOB+jFx2mGR9tuH9o2ZlZc0KvqJmds4ABOwPomdoxCBIY7UYpLPITsgQ8i1PEIHLD3HwWaesIN7GL39w5Qk6IqJsds4ABO/4o3NuZMQgCfvSdiwI+Gxa5r9t16epAd5mdddQ4H6MXHaYZH224f2kdHlwZaVheGZlcqR4YWlkAQ=="
	_, newStxBytes, err := crypto.SignTransaction(private, tx)
	require.NoError(t, err)
	require.EqualValues(t, newStxBytes, byteFromBase64(signedGolden))
}

func TestMakeAssetAcceptanceTxn(t *testing.T) {
	const sender = "BH55E5RMBD4GYWXGX5W5PJ5JAHPGM5OXKDQH5DC4O2MGI7NW4H6VOE4CP4"
	const genesisHash = "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI="
	const assetIndex = 1
	const firstValidRound = 322575
	const lastValidRound = 323575

	tx, err := MakeAssetAcceptanceTxn(sender, 10, firstValidRound,
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

	require.Equal(t, expectedAssetAcceptanceTxn, tx)

	const addrSK = "awful drop leaf tennis indoor begin mandate discover uncle seven only coil atom any hospital uncover make any climb actor armed measure need above hundred"
	private, err := mnemonic.ToPrivateKey(addrSK)
	require.NoError(t, err)
	_, newStxBytes, err := crypto.SignTransaction(private, tx)
	signedGolden := "gqNzaWfEQJ7q2rOT8Sb/wB0F87ld+1zMprxVlYqbUbe+oz0WM63FctIi+K9eYFSqT26XBZ4Rr3+VTJpBE+JLKs8nctl9hgijdHhuiKRhcmN2xCAJ+9J2LAj4bFrmv23Xp6kB3mZ111Dgfoxcdphkfbbh/aNmZWXNCOiiZnbOAATsD6JnaMQgSGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiKibHbOAATv96NzbmTEIAn70nYsCPhsWua/bdenqQHeZnXXUOB+jFx2mGR9tuH9pHR5cGWlYXhmZXKkeGFpZAE="
	require.EqualValues(t, newStxBytes, byteFromBase64(signedGolden))
}

func TestMakeAssetRevocationTransaction(t *testing.T) {
	const addr = "BH55E5RMBD4GYWXGX5W5PJ5JAHPGM5OXKDQH5DC4O2MGI7NW4H6VOE4CP4"
	const genesisHash = "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI="
	const revoker, recipient, revoked = addr, addr, addr
	const assetIndex = 1
	const firstValidRound = 322575
	const lastValidRound = 323575
	const amountToSend = 1

	tx, err := MakeAssetRevocationTxn(revoker, revoked, recipient, amountToSend, 10, firstValidRound,
		lastValidRound, nil, "", genesisHash, assetIndex)
	require.NoError(t, err)

	sendAddr, err := types.DecodeAddress(revoker)
	require.NoError(t, err)

	expectedAssetRevocationTxn := types.Transaction{
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

	expectedAssetID := types.AssetIndex(assetIndex)
	expectedAssetRevocationTxn.XferAsset = expectedAssetID

	receiveAddr, err := types.DecodeAddress(recipient)
	require.NoError(t, err)
	expectedAssetRevocationTxn.AssetReceiver = receiveAddr

	expectedAssetRevocationTxn.AssetAmount = amountToSend

	targetAddr, err := types.DecodeAddress(revoked)
	require.NoError(t, err)
	expectedAssetRevocationTxn.AssetSender = targetAddr

	require.Equal(t, expectedAssetRevocationTxn, tx)

	const addrSK = "awful drop leaf tennis indoor begin mandate discover uncle seven only coil atom any hospital uncover make any climb actor armed measure need above hundred"
	private, err := mnemonic.ToPrivateKey(addrSK)
	require.NoError(t, err)
	_, newStxBytes, err := crypto.SignTransaction(private, tx)
	signedGolden := "gqNzaWfEQHsgfEAmEHUxLLLR9s+Y/yq5WeoGo/jAArCbany+7ZYwExMySzAhmV7M7S8+LBtJalB4EhzEUMKmt3kNKk6+vAWjdHhuiqRhYW10AaRhcmN2xCAJ+9J2LAj4bFrmv23Xp6kB3mZ111Dgfoxcdphkfbbh/aRhc25kxCAJ+9J2LAj4bFrmv23Xp6kB3mZ111Dgfoxcdphkfbbh/aNmZWXNCqqiZnbOAATsD6JnaMQgSGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiKibHbOAATv96NzbmTEIAn70nYsCPhsWua/bdenqQHeZnXXUOB+jFx2mGR9tuH9pHR5cGWlYXhmZXKkeGFpZAE="
	require.EqualValues(t, newStxBytes, byteFromBase64(signedGolden))
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

func TestLogicSig(t *testing.T) {
	// validate LogicSig signed transaction against goal
	const fromAddress = "47YPQTIGQEO7T4Y4RWDYWEKV6RTR2UNBQXBABEEGM72ESWDQNCQ52OPASU"
	const toAddress = "PNWOET7LLOWMBMLE4KOCELCX6X3D3Q4H2Q4QJASYIEOF7YIPPQBG3YQ5YI"
	const referenceTxID = "5FJDJD5LMZC3EHUYYJNH5I23U4X6H2KXABNDGPIL557ZMJ33GZHQ"
	const mn = "advice pudding treat near rule blouse same whisper inner electric quit surface sunny dismiss leader blood seat clown cost exist hospital century reform able sponsor"
	const fee = 1000
	const amount = 2000
	const firstRound = 2063137
	const genesisID = "devnet-v1.0"
	genesisHash := byteFromBase64("sC3P7e2SdbqKJK0tbiCdK9tdSpbe6XeCGKdoNzmlj0E=")
	note := byteFromBase64("8xMCTuLQ810=")

	tx, err := MakePaymentTxnWithFlatFee(
		fromAddress, toAddress, fee, amount, firstRound, firstRound+1000,
		note, "", genesisID, genesisHash,
	)
	require.NoError(t, err)

	// goal clerk send -o tx3 -a 2000 --fee 1000 -d ~/.algorand -w test -L sig.lsig --argb64 MTIz --argb64 NDU2 \
	// -f 47YPQTIGQEO7T4Y4RWDYWEKV6RTR2UNBQXBABEEGM72ESWDQNCQ52OPASU \
	// -t PNWOET7LLOWMBMLE4KOCELCX6X3D3Q4H2Q4QJASYIEOF7YIPPQBG3YQ5YI
	const golden = "gqRsc2lng6NhcmeSxAMxMjPEAzQ1NqFsxAUBIAEBIqNzaWfEQE6HXaI5K0lcq50o/y3bWOYsyw9TLi/oorZB4xaNdn1Z14351u2f6JTON478fl+JhIP4HNRRAIh/I8EWXBPpJQ2jdHhuiqNhbXTNB9CjZmVlzQPoomZ2zgAfeyGjZ2Vuq2Rldm5ldC12MS4womdoxCCwLc/t7ZJ1uookrS1uIJ0r211Klt7pd4IYp2g3OaWPQaJsds4AH38JpG5vdGXECPMTAk7i0PNdo3JjdsQge2ziT+tbrMCxZOKcIixX9fY9w4fUOQSCWEEcX+EPfAKjc25kxCDn8PhNBoEd+fMcjYeLEVX0Zx1RoYXCAJCGZ/RJWHBooaR0eXBlo3BheQ=="

	program := []byte{1, 32, 1, 1, 34}
	args := make([][]byte, 2)
	args[0] = []byte("123")
	args[1] = []byte("456")
	key, err := mnemonic.ToPrivateKey(mn)
	var pk crypto.MultisigAccount
	require.NoError(t, err)
	lsig, err := crypto.MakeLogicSig(program, args, key, pk)
	require.NoError(t, err)

	_, stxBytes, err := crypto.SignLogicsigTransaction(lsig, tx)
	require.NoError(t, err)

	require.Equal(t, byteFromBase64(golden), stxBytes)

	sender, err := types.DecodeAddress(fromAddress)
	require.NoError(t, err)

	verified := crypto.VerifyLogicSig(lsig, sender)
	require.True(t, verified)

}
