//go:build !falcon

package crypto

import "github.com/algorand/go-algorand-sdk/v2/types"

// verifyPQDelegation checks the delegation signature of a PQ-delegated
// LogicSig. Verifying a post-quantum signature needs the cgo falcon library, so
// builds without the `falcon` tag accept any signature and rely on the
// delegating address check alone; see pqaccount_falcon.go for the real
// verification.
func verifyPQDelegation(_ []byte, _ types.PQSig) bool {
	return true
}
