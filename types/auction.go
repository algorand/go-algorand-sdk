package types

// Bid represents a bid by a user as part of an auction.
type Bid struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// BidderKey identifies the bidder placing this bid.
	BidderKey Address `codec:"bidder"`

	// BidCurrency specifies how much external currency the bidder
	// is putting in with this bid.
	BidCurrency uint64 `codec:"cur"`

	// MaxPrice specifies the maximum price, in units of external
	// currency per Algo, that the bidder is willing to pay.
	// This must be at least as high as the current price of the
	// auction in the block in which this bid appears.
	MaxPrice uint64 `codec:"price"`

	// BidID identifies this bid.  The first bid by a bidder (identified
	// by BidderKey) with a particular BidID on the blockchain will be
	// considered, preventing replay of bids.  Specifying a different
	// BidID allows the bidder to place multiple bids in an auction.
	BidID uint64 `codec:"id"`

	// AuctionKey specifies the auction for this bid.
	AuctionKey Address `codec:"auc"`

	// AuctionID identifies the auction for which this bid is intended.
	AuctionID uint64 `codec:"aid"`
}

// SignedBid represents a signed bid by a bidder.
type SignedBid struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// Bid contains information about the bid.
	Bid Bid `codec:"bid"`

	// Sig is a signature by the bidder, as identified in the bid
	// (Bid.BidderKey) over the hash of the Bid.
	Sig Signature `codec:"sig"`
}

// NoteFieldType indicates a type of auction message encoded into a
// transaction's Note field.
type NoteFieldType string

const (
	// NoteDeposit indicates a SignedDeposit message.
	NoteDeposit NoteFieldType = "d"

	// NoteBid indicates a SignedBid message.
	NoteBid NoteFieldType = "b"

	// NoteSettlement indicates a SignedSettlement message.
	NoteSettlement NoteFieldType = "s"

	// NoteParams indicates a SignedParams message.
	NoteParams NoteFieldType = "p"
)

// Settlement describes the outcome of an auction.
type Settlement struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// AuctionKey is the public key for the series of auctions.
	AuctionKey Digest `codec:"auc"`

	// AuctionID identifies the auction being settled.
	AuctionID uint64 `codec:"aid"`

	// Cleared indicates whether the auction fully cleared.
	// It is the same as in the BidOutcomes.
	Cleared bool `codec:"cleared"`

	// OutcomesHash is a hash of the BidOutcomes for this auction.
	// The pre-image (the actual BidOutcomes struct) should be published
	// out-of-band (e.g., on the web site of the Algorand company).
	OutcomesHash Digest `codec:"outhash"`

	// Canceled indicates that the auction was canceled.
	// When Canceled is true, clear and OutcomeHash are false and empty, respectively.
	Canceled bool `codec:"canceled"`
}

// SignedSettlement is a settlement signed by the auction operator
// (e.g., the Algorand company).
type SignedSettlement struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// Settlement describes the outcome of an auction.
	Settlement Settlement `codec:"settle"`

	// Sig is a signature by the auction operator on the hash
	// of the Settlement struct above.
	Sig Signature `codec:"sig"`
}

// Params describes the parameters for a particular auction.
type Params struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// AuctionKey is the public key of the auction operator.
	// This is somewhat superfluous, because the Params struct
	// gets signed by the corresponding private key, so in order
	// to verify a SignedParams, the caller must already know
	// the correct AuctionKey to use.  However, having this field
	// here is useful to allow code to keep track of all auction
	// parameters (including the public key) using a Params struct.
	AuctionKey Digest `codec:"auc"`

	// AuctionID identifies the auction whose parameters are being
	// specified.
	AuctionID uint64 `codec:"aid"`

	// BankKey specifies the key of the external bank that will
	// be signing deposits.
	BankKey Digest `codec:"bank"`

	// DispensingKey specifies the public key of the account from
	// which auction winnings will be dispensed.
	DispensingKey Digest `codec:"dispense"`

	// LastPrice specifies the price at the end of the auction
	// (i.e., in the last chunk), in units of external currency
	// per Algo.  This is called ``reserve price'' in the design doc.
	LastPrice uint64 `codec:"lastprice"`

	// DepositRound specifies the first block in which deposits
	// will be considered.  This can be less than FirstRound to
	// allow the external bank (e.g., CoinList) to place deposits
	// for an auction before bidding begins.
	DepositRound uint64 `codec:"depositrnd"`

	// FirstRound specifies the first block in which bids will be
	// considered.
	FirstRound uint64 `codec:"firstrnd"`

	// PriceChunkRounds specifies the number of blocks for which
	// the price remains the same.  The auction proceeds in chunks
	// of PriceChunkRounds at a time, starting from FirstRound.
	PriceChunkRounds uint64 `codec:"chunkrnds"`

	// NumChunks specifies the number of PriceChunkRounds-sized
	// chunks for which the auction will run.  This means that
	// the last block in which a bid can be placed will be
	// (FirstRound + PriceChunkRounds*NumChunnks - 1).
	NumChunks uint64 `codec:"numchunks"`

	// MaxPriceMultiple defines the ratio between MaxPrice (the
	// starting price of the auction) and LastPrice.  Expect this
	// is on the order of 100.
	MaxPriceMultiple uint64 `codec:"maxmult"`

	// NumAlgos specifies the maximum number of MicroAlgos that will be
	// sold in this auction.
	NumAlgos uint64 `codec:"maxalgos"`

	// MinBidAlgos specifies the minimum amount of a bid, in terms
	// of the number of MicroAlgos at the bid's maximum price.  This
	// should not be less than MinBalance, otherwise the transaction
	// that dispenses winnings might be rejected.
	MinBidAlgos uint64 `codec:"minbidalgos"`
}

// SignedParams is a signed statement by the auction operator attesting
// to the start of an auction.
type SignedParams struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// Params specifies the parameters for the auction.
	Params Params `codec:"param"`

	// Sig is a signature over Params by the operator's key.
	Sig Signature `codec:"sig"`
}

// NoteField is the struct that represents an auction message.
type NoteField struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// Type indicates which type of a message this is
	Type NoteFieldType `codec:"t"`

	// SignedBid, for NoteBid type
	SignedBid SignedBid `codec:"b"`

	// SignedSettlement, for NoteSettlement type
	SignedSettlement SignedSettlement `codec:"s"`

	// SignedParams, for NoteParams type
	SignedParams SignedParams `codec:"p"`
}
