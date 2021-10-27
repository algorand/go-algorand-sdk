package future

import (
	"context"
	"errors"

	"github.com/algorand/go-algorand-sdk/abi"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/types"
)

/**
 * This type represents a function which can sign transactions from an atomic transaction group.
 * @param txnGroup - The atomic group containing transactions to be signed
 * @param indexesToSign - An array of indexes in the atomic transaction group that should be signed
 * @returns A promise which resolves an array of encoded signed transactions. The length of the
 *   array will be the same as the length of indexesToSign, and each index i in the array
 *   corresponds to the signed transaction from txnGroup[indexesToSign[i]]
 */
type TransactionSigner = func(txGroup []types.Transaction, indexesToSign []int) ([][]byte, []string, error)

/**
 * Create a TransactionSigner that can sign transactions for the provided basic Account.
 */
func makeBasicAccountTransactionSigner(account crypto.Account) TransactionSigner {

	return func(txGroup []types.Transaction, indexesToSign []int) ([][]byte, []string, error) {
		var stxs [][]byte
		var txids []string
		for _, pos := range indexesToSign {
			txid, stxBytes, err := crypto.SignTransaction(account.PrivateKey, txGroup[pos])
			if err != nil {
				return nil, nil, err
			}

			stxs = append(stxs, stxBytes)
			txids = append(txids, txid)
		}

		return stxs, txids, nil
	}
}

/**
 * Create a TransactionSigner that can sign transactions for the provided LogicSigAccount.
 */
func makeLogicSigAccountTransactionSigner(logicSigAccount crypto.LogicSigAccount) TransactionSigner {

	return func(txGroup []types.Transaction, indexesToSign []int) ([][]byte, []string, error) {
		var stxs [][]byte
		var txids []string
		for _, pos := range indexesToSign {
			txid, stxBytes, err := crypto.SignLogicSigAccountTransaction(logicSigAccount, txGroup[pos])
			if err != nil {
				return nil, nil, err
			}

			stxs = append(stxs, stxBytes)
			txids = append(txids, txid)
		}

		return stxs, txids, nil
	}
}

/**
 * Create a TransactionSigner that can sign transactions for the provided Multisig account.
 * @param msig - The Multisig account metadata
 * @param sks - An array of private keys belonging to the msig which should sign the transactions.
 */
func makeMultiSigAccountTransactionSigner(msig crypto.MultisigAccount, sks [][]byte) TransactionSigner {

	return func(txGroup []types.Transaction, indexesToSign []int) ([][]byte, []string, error) {
		var stxs [][]byte
		var txids []string
		for _, pos := range indexesToSign {
			var unmergedStxs [][]byte
			for _, sk := range sks {
				_, unmergedStxBytes, err := crypto.SignMultisigTransaction(sk, msig, txGroup[pos])
				if err != nil {
					return nil, nil, err
				}

				unmergedStxs = append(unmergedStxs, unmergedStxBytes)
			}

			txid, stxBytes, err := crypto.MergeMultisigTransactions(unmergedStxs...)
			if err != nil {
				return nil, nil, err
			}

			stxs = append(stxs, stxBytes)
			txids = append(txids, txid)
		}

		return stxs, txids, nil
	}
}

/** Represents an unsigned transactions and a signer that can authorize that transaction. */
type TransactionWithSigner struct {
	/** An unsigned transaction */
	Txn types.Transaction
	/** A transaction signer that can authorize txn */
	Signer TransactionSigner
}

type MethodArgument interface {
	IsMethodArgument()
}

type ABIValue struct {
	abi.Value
}

func (abiValue ABIValue) IsMethodArgument()              {}
func (txSigner TransactionWithSigner) IsMethodArgument() {}

type AtomicTransactionComposerStatus = int

const (
	/** The atomic group is still under construction. */
	BUILDING AtomicTransactionComposerStatus = 1

	/** The atomic group has been finalized, but not yet signed. */
	BUILT = 2

	/** The atomic group has been finalized and signed, but not yet submitted to the network. */
	SIGNED = 3

	/** The atomic gorup has been finalized, signed, and submitted to the network. */
	SUBMITTED = 4

	/** The atomic group has been finalized, signed, submitted, and successfully committed to a block. */
	COMMITTED = 5
)

/** A class used to construct and execute atomic transaction groups */
type AtomicTransactionComposer struct {

	/** The maximum size of an atomic transaction group. */
	MAX_GROUP_SIZE int

	status AtomicTransactionComposerStatus

	transactions []TransactionWithSigner

	signedTxs [][]byte

	txids []string
}

/**
* Get the status of this composer's transaction group.
 */
func (atc *AtomicTransactionComposer) getStatus() AtomicTransactionComposerStatus {
	return atc.status
}

/**
* Get the number of transactions currently in this atomic group.
 */
func (atc *AtomicTransactionComposer) count() int {
	return len(atc.transactions)
}

/**
* Create a new composer with the same underlying transactions. The new composer's status will be
* BUILDING, so additional transactions may be added to it.
 */
func (atc *AtomicTransactionComposer) clone() AtomicTransactionComposer {

}

/**
* Add a transaction to this atomic group.
*
* An error will be thrown if the composer's status is not BUILDING, or if adding this transaction
* causes the current group to exceed MAX_GROUP_SIZE.
 */
func (atc *AtomicTransactionComposer) addTransaction(txnAndSigner TransactionWithSigner) error {
	if atc.status != BUILDING {
		return errors.New("Transactions cannot be added when the status of the transaction composer is not BUILDING")
	}

	if atc.count() == atc.MAX_GROUP_SIZE {
		return errors.New("The transaction composer has reached its MAX_GROUP_SIZE: %d", atc.MAX_GROUP_SIZE)
	}

	atc.transactions = append(atc.transactions, txnAndSigner)
	return nil
}

/**
* Add a smart contract method call to this atomic group.
*
* An error will be thrown if the composer's status is not BUILDING, if adding this transaction
* causes the current group to exceed MAX_GROUP_SIZE, or if the provided arguments are invalid
* for the given method.
 */
func (atc *AtomicTransactionComposer) addMethodCall(
	/** The ID of the smart contract to call */
	appID uint64,
	/** The method to call on the smart contract */
	method abi.Method,
	/** The arguments to include in the method call. If omitted, no arguments will be passed to the method. */
	methodArgs []MethodArgument,
	/** The address of the sender of this application call */
	sender types.Address,
	/** Transactions params to use for this application call */
	suggestedParams types.SuggestedParams,
	/** The OnComplete action to take for this application call. If omitted, OnApplicationComplete.NoOpOC will be used. */
	onComplete types.OnCompletion,
	/** The note value for this application call */
	note []byte,
	/** The lease value for this application call */
	lease [32]byte,
	/** If provided, the address that the sender will be rekeyed to at the conclusion of this application call */
	rekeyTo types.Address,
	/** A transaction signer that can authorize this application call from sender */
	signer TransactionSigner) error {

	if atc.status != BUILDING {
		return errors.New("Transactions cannot be added when the status of the transaction composer is not BUILDING")
	}

	var abiArgs []ABIValue
	for _, methodArg := range methodArgs {
		switch v := methodArg.(type) {
		case ABIValue:
			abiArgs = append(abiArgs, v)
		case TransactionWithSigner:
			err := atc.addTransaction(v)
			if err != nil {
				return err
			}
		default:
			return errors.New("MethodArg has type %s, only ABIValues and TransactionSigners are allowed")
		}
	}

	if atc.count() == atc.MAX_GROUP_SIZE {
		return errors.New("The transaction composer has reached its MAX_GROUP_SIZE: %d", atc.MAX_GROUP_SIZE)
	}

	var encodedAbiArgs [][]byte
	for _, abiArg := range abiArgs {
		encodedArg, err := abiArg.Encode()
		if err != nil {
			return err
		}

		encodedAbiArgs = append(encodedAbiArgs, encodedArg)
	}

	tx, err := MakeApplicationNoOpTx(
		appID,
		encodedAbiArgs,
		nil,
		nil,
		nil,
		suggestedParams,
		sender,
		note,
		types.Digest{},
		lease,
		rekeyTo)

	if err != nil {
		return err
	}

	txSigner := TransactionWithSigner{
		Txn:    tx,
		Signer: signer,
	}

	atc.transactions = append(atc.transactions, txSigner)
	return nil
}

/**
* Finalize the transaction group and returned the finalized transactions.
*
* The composer's status will be at least BUILT after executing this method.
 */
func (atc *AtomicTransactionComposer) buildGroup() ([]TransactionWithSigner, error) {
	if atc.status > BUILDING {
		return atc.transactions, nil
	}

	var txns []types.Transaction
	for _, txSigner := range atc.transactions {
		txns = append(txns, txSigner.Txn)
	}

	gid, err := crypto.ComputeGroupID(txns)
	if err != nil {
		return nil, err
	}

	for _, txSigner := range atc.transactions {
		txSigner.Txn.Group = gid
	}

	if atc.status == BUILDING {
		atc.status = BUILT
	}

	return atc.transactions, nil
}

/**
* Obtain signatures for each transaction in this group. If signatures have already been obtained,
* this method will return cached versions of the signatures.
*
* The composer's status will be at least SIGNED after executing this method.
*
* An error will be thrown if signing any of the transactions fails.
*
* @returns A promise that resolves to an array of signed transactions.
 */
func (atc *AtomicTransactionComposer) gatherSignatures() ([][]byte, error) {
	// if status is at least signed then return cached signed transactions
	if atc.status >= SIGNED {
		return atc.signedTxs, nil
	}

	// retrieve built transactions and verify status is BUILT
	txs, err := atc.buildGroup()
	if err != nil {
		return nil, err
	}

	var sigTxs [][]byte
	var txids []string
	for _, txSigner := range txs {
		txGroup := []types.Transaction{txSigner.Txn}
		indexesToSign := []int{0}
		sigStx, txid, err := txSigner.Signer(txGroup, indexesToSign)
		if err != nil {
			return nil, err
		}

		sigTxs = append(sigTxs, sigStx[0])
		txids = append(txids, txid[0])
	}

	atc.signedTxs = sigTxs
	atc.txids = txids
	atc.status = SIGNED
	return sigTxs, nil
}

/**
* Send the transaction group to the network, but don't wait for it to be committed to a block. An
* error will be thrown if submission fails.
*
* The composer's status must be SUBMITTED or lower before calling this method. If submission is
* successful, this composer's status will update to SUBMITTED.
*
* Note: a group can only be submitted again if it fails.
*
* @returns A promise that, upon success, resolves to a list of TxIDs of the submitted transactions.
 */
func (atc *AtomicTransactionComposer) submit(client *algod.Client) ([]string, error) {
	if atc.status > SUBMITTED {
		return nil, errors.New("Transaction composer status must be SUBMITTED or lower")
	}

	if atc.status == SUBMITTED {
		return atc.txids, nil
	}

	stxs, err := atc.gatherSignatures()
	if err != nil {
		return nil, err
	}

	var serializedStxs []byte
	for _, stx := range stxs {
		serializedStxs = append(serializedStxs, stx...)
	}

	_, err = client.SendRawTransaction(serializedStxs).Do(context.Background())
	if err != nil {
		return nil, err
	}

	return atc.txids, nil
}

/**
* Send the transaction group to the network and wait until it's committed to a block. An error
* will be thrown if submission or execution fails.
*
* The composer's status must be SUBMITTED or lower before calling this method, since execution is
* only allowed once. If submission is successful, this composer's status will update to SUBMITTED.
* If the execution is also successful, this composer's status will update to COMMITTED.
*
* Note: a group can only be submitted again if it fails.
*
* @returns A promise that, upon success, resolves to an object containing the confirmed round for
*   this transaction, the txIDs of the submitted transactions, and an array of results containing
*   one element for each method call transaction in this group. If a method has no return value
*   (void), then the method results array will contain null for that method's return value.
 */
func (atc *AtomicTransactionComposer) execute(client *algod.Client) (int, []string, []ABIValue, error) {
	if atc.status > SUBMITTED {
		return 0, nil, nil, errors.New("Transaction composer status is already COMMITTED")
	}

	_, err := atc.submit(client)
	if err != nil {
		return 0, nil, nil, err
	}

	return 0, nil, nil, err
}
