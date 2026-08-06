//go:build !falcon

package crypto

import "github.com/algorand/go-algorand-sdk/v2/types"

// verifyPQSig is the unexported version of VerifyPQSig used by VerifyLogicSig.
// On falcon-enabled environment the implementation will use the falcon impl to validate the signature.
// See pqaccount_falcon.go for more info.
func verifyPQSig(toBeSigned []byte, pqsig types.PQSig) bool {
	return false
}
