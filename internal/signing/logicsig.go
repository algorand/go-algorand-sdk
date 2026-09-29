package signing

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/types"
)

// ErrLsigEmptyMsig is returned when appending a multisig signature to a
// LogicSig that is not delegated to a multisig account.
var ErrLsigEmptyMsig = errors.New("empty multisig in logicsig")

func isASCIIPrintableByte(symbol byte) bool {
	isBreakLine := symbol == '\n'
	isStdPrintable := symbol >= ' ' && symbol <= '~'
	return isBreakLine || isStdPrintable
}

func isASCIIPrintable(program []byte) bool {
	for _, b := range program {
		if !isASCIIPrintableByte(b) {
			return false
		}
	}
	return true
}

// SanityCheckProgram performs heuristic program validation:
// check if passed in bytes are Algorand address or is B64 encoded, rather than Teal bytes
func SanityCheckProgram(program []byte) error {
	if len(program) == 0 {
		return fmt.Errorf("empty program")
	}
	if isASCIIPrintable(program) {
		if _, err := types.DecodeAddress(string(program)); err == nil {
			return fmt.Errorf("requesting program bytes, get Algorand address")
		}
		if _, err := base64.StdEncoding.DecodeString(string(program)); err == nil {
			return fmt.Errorf("program should not be b64 encoded")
		}
		return fmt.Errorf("program bytes are all ASCII printable characters, not looking like Teal byte code")
	}
	return nil
}

// EscrowLogicSig returns a LogicSig that is not delegated to any account.
func EscrowLogicSig(program []byte, args [][]byte) (lsig types.LogicSig, err error) {
	if err = SanityCheckProgram(program); err != nil {
		return
	}
	return types.LogicSig{Logic: program, Args: args}, nil
}

// Ed25519DelegatedLogicSig returns a LogicSig delegated to the single ed25519
// account of sgnr.
func Ed25519DelegatedLogicSig(program []byte, args [][]byte, sgnr Ed25519Signer) (lsig types.LogicSig, err error) {
	if err = SanityCheckProgram(program); err != nil {
		return
	}

	sig, err := Ed25519RawSignature(sgnr, ProgramToSign(program))
	if err != nil {
		return
	}
	return types.LogicSig{Logic: program, Args: args, Sig: sig}, nil
}

// Ed25519MultisigDelegatedLogicSig returns a LogicSig delegated to the multisig
// account ma, carrying only the signature of its member sgnr. Additional
// signatures can be added with Ed25519AppendMultisigToLogicSig.
func Ed25519MultisigDelegatedLogicSig(program []byte, args [][]byte, ma MultisigAccount, sgnr Ed25519Signer) (lsig types.LogicSig, err error) {
	if err = SanityCheckProgram(program); err != nil {
		return
	}

	multisigAddr, err := ma.Address()
	if err != nil {
		return
	}

	msig, _, err := ed25519MultisigSig(sgnr, ma, MsigProgramToSign(multisigAddr, program))
	if err != nil {
		return
	}
	return types.LogicSig{Logic: program, Args: args, LMsig: msig}, nil
}

// Ed25519AppendMultisigToLogicSig adds the signature of sgnr to a LogicSig
// delegated to a multisig account.
func Ed25519AppendMultisigToLogicSig(lsig *types.LogicSig, sgnr Ed25519Signer) error {
	if lsig.LMsig.Blank() {
		return ErrLsigEmptyMsig
	}

	ma, err := MultisigAccountFromSig(lsig.LMsig)
	if err != nil {
		return err
	}

	multisigAddr, err := ma.Address()
	if err != nil {
		return err
	}

	msig, idx, err := ed25519MultisigSig(sgnr, ma, MsigProgramToSign(multisigAddr, lsig.Logic))
	if err != nil {
		return err
	}

	lsig.LMsig.Subsigs[idx] = msig.Subsigs[idx]

	return nil
}

// PQDelegatedLogicSig returns a LogicSig delegated to the PQ account of sgnr.
func PQDelegatedLogicSig(program []byte, args [][]byte, sgnr PQSigner) (lsig types.LogicSig, err error) {
	if err = SanityCheckProgram(program); err != nil {
		return
	}

	// the delegation signature commits to the address it delegates from, so the
	// envelope is assembled first and its signature filled in afterwards
	pqsig, addr, err := PQSigFor(sgnr, nil)
	if err != nil {
		return
	}

	pqsig.Signature, err = sgnr.PQSign(PQProgramToSign(addr, program))
	if err != nil {
		return
	}

	return types.LogicSig{Logic: program, Args: args, PQsig: pqsig}, nil
}
