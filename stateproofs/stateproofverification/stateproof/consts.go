package stateproof

const (
	precisionBits       = uint8(16)     // number of bits used for log approximation. This should not exceed 63
	ln2IntApproximation = uint64(45427) // the value of the ln(2) with 16 bits of precision (i.e  ln2IntApproximation = ceil( 2^precisionBits * ln(2) ))
	MaxReveals          = 640           // MaxReveals is a bound on allocation and on numReveals to limit log computation
	// VersionForCoinGenerator is used as part of the seed for Fiat-Shamir. We would change this
	// value if the state proof verifier algorithm changes. This will allow us to make different coins for different state proof verification algorithms
	VersionForCoinGenerator = byte(0)
	// MaxTreeDepth defines the maximal size of a merkle tree depth the state proof allows.
	MaxTreeDepth = 20
)
