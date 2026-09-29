package crypto

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ed25519"

	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

func TestGenerateAddressFromSK(t *testing.T) {
	acct := GenerateAccount()

	addr, err := GenerateAddressFromSK(acct.PrivateKey)
	require.NoError(t, err)
	require.Equal(t, acct.Address, addr)
}

func TestSKToInMemorySigner(t *testing.T) {
	acct := GenerateAccount()

	sgnr, err := SKToInMemorySigner(acct.PrivateKey)
	require.NoError(t, err)
	require.Equal(t, Ed25519PublicKey(acct.Address), sgnr.Ed25519PublicKey())

	message := []byte("test message")
	sig, err := sgnr.Ed25519Sign(message)
	require.NoError(t, err)
	require.True(t, ed25519.Verify(acct.PublicKey, message, sig))

	// Invalid keys
	_, err = SKToInMemorySigner(nil)
	require.ErrorIs(t, err, errInvalidPrivateKey)

	_, err = SKToInMemorySigner(ed25519.PrivateKey(make([]byte, 32)))
	require.ErrorIs(t, err, errInvalidPrivateKey)

	_, err = SKToInMemorySigner(ed25519.PrivateKey(make([]byte, 65)))
	require.ErrorIs(t, err, errInvalidPrivateKey)
}

func TestSignBid(t *testing.T) {
	bidder := GenerateAccount()
	auction := GenerateAccount()
	bid := types.Bid{
		BidderKey:   bidder.Address,
		BidCurrency: 1000,
		MaxPrice:    10,
		BidID:       1,
		AuctionKey:  auction.Address,
		AuctionID:   2,
	}

	signedBid, err := SignBid(bidder.PrivateKey, bid)
	require.NoError(t, err)

	var nf types.NoteField
	require.NoError(t, msgpack.Decode(signedBid, &nf))
	require.Equal(t, types.NoteBid, nf.Type)
	require.Equal(t, bid, nf.SignedBid.Bid)

	toBeSigned := append([]byte("aB"), msgpack.Encode(bid)...)
	require.True(t, ed25519.Verify(bidder.PublicKey, toBeSigned, nf.SignedBid.Sig[:]))

	_, err = SignBid(nil, bid)
	require.ErrorIs(t, err, errInvalidPrivateKey)
}
