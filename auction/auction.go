package auction

import (
	"fmt"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/types"
)

// MakeBid constructs a bid using the passed parameters. `bidderAddress` and
// `auctionAddress` should be checksummed, human-readable addresses
func MakeBid(bidderAddress string, bidAmount, maxPrice, bidID int64, auctionAddress string, auctionID int64) (encodedBid []byte, err error) {

	// Sanity check for int64
	if bidAmount < 0 ||
		maxPrice < 0 ||
		bidID < 0 ||
		auctionID < 0 {
		err = fmt.Errorf("all numbers must not be negative")
		return
	}

	// Decode from address
	bidderAddr, err := types.DecodeAddress(bidderAddress)
	if err != nil {
		return
	}

	// Decode to address
	auctionAddr, err := types.DecodeAddress(auctionAddress)
	if err != nil {
		return
	}

	bid := types.Bid{
		BidderKey:   bidderAddr,
		BidCurrency: uint64(bidAmount),
		MaxPrice:    uint64(maxPrice),
		BidID:       uint64(bidID),
		AuctionKey:  auctionAddr,
		AuctionID:   uint64(auctionID),
	}

	encodedBid = msgpack.Encode(bid)
	return
}
