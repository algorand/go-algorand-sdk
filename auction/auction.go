package auction

import (
	"github.com/algorand/go-algorand-sdk/types"
)

// MakeBid constructs a bid using the passed parameters. `bidderAddress` and
// `auctionAddress` should be checksummed, human-readable addresses
func MakeBid(bidderAddress string, bidAmount, maxPrice, bidID uint64, auctionAddress string, auctionID uint64) (bid types.Bid, err error) {
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

	bid = types.Bid{
		BidderKey:   bidderAddr,
		BidCurrency: bidAmount,
		MaxPrice:    maxPrice,
		BidID:       bidID,
		AuctionKey:  auctionAddr,
		AuctionID:   auctionID,
	}
	return
}
