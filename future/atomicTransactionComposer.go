package future

import (
	"context"
	"errors"
	"fmt"

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
func MakeBasicAccountTransactionSigner(account crypto.Account) TransactionSigner {

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
func MakeLogicSigAccountTransactionSigner(logicSigAccount crypto.LogicSigAccount) TransactionSigner {

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
func MakeMultiSigAccountTransactionSigner(msig crypto.MultisigAccount, sks [][]byte) TransactionSigner {

	return func(txGroup []types.Transaction, indexesToSign []int) ([][]byte, []string, error) {
		var stxs [][]byte
		var txids []string
		for _, pos := range indexesToSign {
			var unmergedStxs [][]byte
			var txid string
			for _, sk := range sks {
				tempTxid, unmergedStxBytes, err := crypto.SignMultisigTransaction(sk, msig, txGroup[pos])
				txid = tempTxid
				if err != nil {
					return nil, nil, err
				}

				unmergedStxs = append(unmergedStxs, unmergedStxBytes)
			}

			if len(sks) > 1 {
				tempTxid, stxBytes, err := crypto.MergeMultisigTransactions(unmergedStxs...)
				txid = tempTxid
				if err != nil {
					return nil, nil, err
				}

				stxs = append(stxs, stxBytes)
				txids = append(txids, txid)
			} else {
				stxs = append(stxs, unmergedStxs[0])
				txids = append(txids, txid)
			}
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

type ABIValue abi.Value

func (abiValue ABIValue) IsMethodArgument()                 {}
func (txAndSigner TransactionWithSigner) IsMethodArgument() {}

type AtomicTransactionComposerStatus = int

const (
	/** The atomic group is still under construction. */
	BUILDING AtomicTransactionComposerStatus = iota

	/** The atomic group has been finalized, but not yet signed. */
	BUILT

	/** The atomic group has been finalized and signed, but not yet submitted to the network. */
	SIGNED

	/** The atomic gorup has been finalized, signed, and submitted to the network. */
	SUBMITTED

	/** The atomic group has been finalized, signed, submitted, and successfully committed to a block. */
	COMMITTED
)

/** The maximum size of an atomic transaction group. */
const MAX_GROUP_SIZE = 16

/** A class used to construct and execute atomic transaction groups */
type AtomicTransactionComposer struct {
	/** The current status of the composer. The status increases monotonically. */
	status AtomicTransactionComposerStatus

	/** The transactions in the group with their respective signers. If status is greater then BUILDING
	then this slice cannot change. */
	transactions []TransactionWithSigner

	/** Method calls constructed from information passed into AddMethodCall(). Used in Execute() to extract
	return values. */
	methodCalls []Method

	/** The raw signed transactions populated after invocation of GatherSignatures(). */
	signedTxs [][]byte

	/** Txids of the transactions in this group. Populated after invocation to GatherSignatures() and used
	in Execute() to gather transaction information */
	txids []string
}

func MakeAtomicTransactionComposer() AtomicTransactionComposer {
	return AtomicTransactionComposer{
		status:       BUILDING,
		transactions: []TransactionWithSigner{},
		methodCalls:  []Method{},
		signedTxs:    [][]byte{},
		txids:        []string{},
	}
}

/**
* Get the status of this composer's transaction group.
 */
func (atc *AtomicTransactionComposer) GetStatus() AtomicTransactionComposerStatus {
	return atc.status
}

/**
* Get the number of transactions currently in this atomic group.
 */
func (atc *AtomicTransactionComposer) Count() int {
	return len(atc.transactions)
}

/**
* Create a new composer with the same underlying transactions. The new composer's status will be
* BUILDING, so additional transactions may be added to it.
 */
func (atc *AtomicTransactionComposer) Clone() AtomicTransactionComposer {
	newTxs := atc.transactions
	for _, txAndSigner := range newTxs {
		txAndSigner.Txn.Group = types.Digest{}
	}

	return AtomicTransactionComposer{
		status:       BUILDING,
		transactions: newTxs,
		methodCalls:  atc.methodCalls,
		signedTxs:    [][]byte{},
		txids:        []string{},
	}
}

/**
* Add a transaction to this atomic group.
*
* An error will be thrown if the composer's status is not BUILDING, or if adding this transaction
* causes the current group to exceed MAX_GROUP_SIZE.
 */
func (atc *AtomicTransactionComposer) AddTransaction(txnAndSigner TransactionWithSigner) error {
	if atc.status != BUILDING {
		return errors.New("STATUS MUST BE BUILDING IN ORDER TO ADD TRANSACTION")
	}

	if atc.Count() == MAX_GROUP_SIZE {
		return fmt.Errorf("REACHED MAX_GROUP_SIZE: %d", MAX_GROUP_SIZE)
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
func (atc *AtomicTransactionComposer) AddMethodCall(
	/** The ID of the smart contract to call */
	appID uint64,
	/** The method to call on the smart contract */
	method Method,
	/** The arguments to include in the method call. If omitted, no arguments will be passed to the method. */
	methodArgs []MethodArgument,
	/** The address of the sender of this application call */
	sender types.Address,
	/** Transactions params to use for this application call */
	suggestedParams types.SuggestedParams,
	/** The note value for this application call */
	note []byte,
	/** The lease value for this application call */
	lease [32]byte,
	/** If provided, the address that the sender will be rekeyed to at the conclusion of this application call */
	rekeyTo types.Address,
	/** A transaction signer that can authorize this application call from sender */
	signer TransactionSigner) error {

	if atc.status != BUILDING {
		return errors.New("STATUS MUST BE BUILDING IN ORDER TO ADD TRANSACTION")
	}

	selectorValue := method.GetSelector()
	encodedAbiArgs := [][]byte{selectorValue}

	// insert selector as first param
	var abiArgs []ABIValue
	for _, methodArg := range methodArgs {
		switch v := methodArg.(type) {
		case ABIValue:
			abiArgs = append(abiArgs, v)
		case TransactionWithSigner:
			err := atc.AddTransaction(v)
			if err != nil {
				return err
			}
		default:
			return errors.New("MethodArg must be either ABIValue or TransactionSigner")
		}
	}

	if atc.Count() == MAX_GROUP_SIZE {
		return fmt.Errorf("REACHED MAX_GROUP_SIZE: %d", MAX_GROUP_SIZE)
	}

	for _, abiArg := range abiArgs {
		encodedArg, err := abiArg.AbiType.Encode(abiArg.RawValue)
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

	txAndSigner := TransactionWithSigner{
		Txn:    tx,
		Signer: signer,
	}

	atc.transactions = append(atc.transactions, txAndSigner)
	atc.methodCalls = append(atc.methodCalls, method)
	return nil
}

/**
* Finalize the transaction group and returned the finalized transactions.
*
* The composer's status will be at least BUILT after executing this method.
 */
func (atc *AtomicTransactionComposer) BuildGroup() ([]TransactionWithSigner, error) {
	if atc.status > BUILDING {
		return atc.transactions, nil
	}

	var txns []types.Transaction
	for _, txAndSigner := range atc.transactions {
		txns = append(txns, txAndSigner.Txn)
	}

	gid, err := crypto.ComputeGroupID(txns)
	if err != nil {
		return nil, err
	}

	for _, txAndSigner := range atc.transactions {
		txAndSigner.Txn.Group = gid
	}

	atc.status = BUILT
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
func (atc *AtomicTransactionComposer) GatherSignatures() ([][]byte, error) {
	// if status is at least signed then return cached signed transactions
	if atc.status >= SIGNED {
		return atc.signedTxs, nil
	}

	// retrieve built transactions and verify status is BUILT
	txs, err := atc.BuildGroup()
	if err != nil {
		return nil, err
	}

	var sigTxs [][]byte
	var txids []string
	for _, txAndSigner := range txs {
		txGroup := []types.Transaction{txAndSigner.Txn}
		indexesToSign := []int{0}
		sigStx, txid, err := txAndSigner.Signer(txGroup, indexesToSign)
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
func (atc *AtomicTransactionComposer) Submit(client *algod.Client) ([]string, error) {
	if atc.status > SUBMITTED {
		return nil, errors.New("STATUS MUST BE SUBMITTED OR LOWER")
	}

	if atc.status == SUBMITTED {
		return atc.txids, nil
	}

	stxs, err := atc.GatherSignatures()
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
func (atc *AtomicTransactionComposer) Execute(client *algod.Client) (int, []string, []ABIValue, error) {
	if atc.status > SUBMITTED {
		return 0, nil, nil, errors.New("STATUS IS ALREADY COMMITTED")
	}

	_, err := atc.Submit(client)
	if err != nil {
		return 0, nil, nil, err
	}

	txinfo, err := WaitForConfirmation(client, atc.txids[0], 0, context.Background())
	if err != nil {
		return 0, nil, nil, err
	}

	var returnValues []ABIValue
	for i, txid := range atc.txids {
		if atc.transactions[i].Txn.Type != types.ApplicationCallTx {
			continue
		}

		txinfo, _, err := client.PendingTransactionInformation(txid).Do(context.Background())
		if err != nil {
			return 0, nil, nil, err
		}

		for _, log := range txinfo.Logs {
			isReturnValue := true
			for _, b := range log[:4] {
				isReturnValue = isReturnValue && (b == 0)
			}

			if isReturnValue {
				abiType, err := abi.TypeOf(atc.methodCalls[i].Returns.AbiType)
				if err != nil {
					return 0, nil, nil, err
				}

				rawAbiValue, err := abiType.Decode(log[4:])
				if err != nil {
					return 0, nil, nil, err
				}

				abiValue := ABIValue{
					AbiType:  abiType,
					RawValue: rawAbiValue,
				}
				returnValues = append(returnValues, abiValue)

				break
			}
		}
	}

	return int(txinfo.ConfirmedRound), atc.txids, returnValues, err
}
