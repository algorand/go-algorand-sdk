//go:build falcon

package crypto

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/algorand/go-algorand-sdk/v2/types"
)

// getFalcon1024GoldenAccount returns the fixed PQ account shared by all
// fixtures, alongside its address.
func getFalcon1024GoldenAccount(t *testing.T) (Falcon1024Account, types.Address) {
	var falcon1024GoldenSeed = [32]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}
	pqa, err := Falcon1024AccountFromPQSeed(falcon1024GoldenSeed[:])
	require.NoError(t, err)
	addr, err := pqa.Address()
	require.NoError(t, err)
	return pqa, addr
}

// With falcon available, VerifyLogicSig cryptographically verifies the PQ
// delegation signature rather than only checking the delegating address.
func TestVerifyLogicSigFalcon1024Delegation(t *testing.T) {
	pqa, pqaAddr := getFalcon1024GoldenAccount(t)

	program := []byte{1, 32, 1, 1, 34}
	lsa := pqDelegatedLogicSigAccount(t, pqa, program, nil)
	require.True(t, VerifyLogicSig(lsa.Lsig, pqaAddr))

	// A delegation signature over a different program must be rejected, even
	// though the envelope still names the right delegating account.
	otherProgram := []byte{1, 32, 1, 2, 34}
	otherLsa := pqDelegatedLogicSigAccount(t, pqa, otherProgram, nil)

	swapped := lsa.Lsig
	swapped.PQsig.Signature = otherLsa.Lsig.PQsig.Signature
	require.False(t, VerifyLogicSig(swapped, pqaAddr))

	// So must a corrupted one.
	tampered := lsa.Lsig
	tampered.PQsig.Signature = append([]byte(nil), lsa.Lsig.PQsig.Signature...)
	tampered.PQsig.Signature[len(tampered.PQsig.Signature)-1] ^= 0xff
	require.False(t, VerifyLogicSig(tampered, pqaAddr))

	// Signing gates on VerifyLogicSig, so a bad delegation cannot be sent.
	_, _, err := SignLogicSigAccountTransaction(
		LogicSigAccount{Lsig: tampered},
		falcon1024GoldenTxn(t, pqaAddr),
	)
	require.ErrorIs(t, err, errLsigInvalidSignature)
}

// pqDelegatedLogicSigAccount builds a PQ-delegated LogicSigAccount by hand: the
// delegating account signs ("PQProgram" || address || program), and the
// signature travels in the envelope naming that account.
func pqDelegatedLogicSigAccount(t *testing.T, pqa Falcon1024Account, program []byte, args [][]byte) LogicSigAccount {
	sgnr := pqa.AsSigner()
	addr, err := pqa.Address()
	require.NoError(t, err)
	salt, err := CanonicalSaltForPQAddress(sgnr.PQPublicKey(), sgnr.PQScheme())
	require.NoError(t, err)

	signature, err := sgnr.PQSign(bytes.Join([][]byte{[]byte("PQProgram"), addr[:], program}, nil))
	require.NoError(t, err)

	return LogicSigAccount{Lsig: types.LogicSig{
		Logic: program,
		Args:  args,
		PQsig: types.PQSig{
			Scheme:    sgnr.PQScheme(),
			Salt:      salt,
			PublicKey: sgnr.PQPublicKey(),
			Signature: signature,
		},
	}}
}
