package merklesignature

// firstRoundInKeyLifetime calculates the round of the valid key for a given round by lowering to the closest KeyLiftime divisor.
// It is implicitly assumed that round is larger than keyLifetime, as an MSS key for round 0 is not valid.
// A key lifetime of 0 is invalid.
func firstRoundInKeyLifetime(round, keyLifetime uint64) uint64 {
	return round - (round % keyLifetime)
}
