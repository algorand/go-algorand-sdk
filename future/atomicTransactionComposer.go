package future

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/algorand/go-algorand-sdk/abi"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/types"
)

var ABI_RETURN_HASH = []byte{0x15, 0x1f, 0x7c, 0x75}

/** Represents an unsigned transactions and a signer that can authorize that transaction. */
type TransactionWithSigner struct {
	/** An unsigned transaction */
	Txn types.Transaction
	/** A transaction signer that can authorize txn */
	Signer TransactionSigner
}

/** Represents the output from a successful ABI method call. */
type ABIResult struct {
	/** The TxID of the transaction that invoked the ABI method call. */
	TxID string
	/**
	 * The raw bytes of the return value from the ABI method call. This will be empty if the method
	 * does not return a value (return type "void").
	 */
	RawReturnValue []byte
	/**
	 * The return value from the ABI method call. This will be undefined if the method does not return
	 * a value (return type "void"), or if the SDK was unable to decode the returned value.
	 */
	ReturnValue interface{}
	/** If the SDK was unable to decode a return value, the error will be here. */
	DecodeError error
}

type AddMethodCallParams struct {
	/** The ID of the smart contract to call */
	AppID uint64
	/** The method to call on the smart contract */
	Method Method
	/** The arguments to include in the method call. If omitted, no arguments will be passed to the method. */
	MethodArgs []interface{}
	/** The address of the sender of this application call */
	Sender types.Address
	/** Transactions params to use for this application call */
	SuggestedParams types.SuggestedParams
	/** The OnComplete action to take for this application call. If omitted, OnApplicationComplete.NoOpOC will be used. */
	OnComplete types.OnCompletion
	/** The note value for this application call */
	Note []byte
	/** The lease value for this application call */
	Lease [32]byte
	/** If provided, the address that the sender will be rekeyed to at the conclusion of this application call */
	RekeyTo types.Address
}

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

type transactionContext struct {
	// The main transaction.
	txn types.Transaction

	// The corresponding signer responsible for producing the signed transaction.
	signer TransactionSigner

	// The corresponding Method constructed from information passed into atc.AddMethodCall().
	method *Method

	// The raw signed transaction populated after invocation of atc.GatherSignatures().
	stxBytes []byte
}

func (txContext *transactionContext) txID() string {
	return crypto.GetTxID(txContext.txn)
}

func (txContext *transactionContext) isMethodCallTx() bool {
	return txContext.method != nil
}

/** The maximum size of an atomic transaction group. */
const MAX_GROUP_SIZE = 16

/** A class used to construct and execute atomic transaction groups */
type AtomicTransactionComposer struct {
	/** The current status of the composer. The status increases monotonically. */
	status AtomicTransactionComposerStatus

	/** The transaction contexts in the group with their respective signers. If status is greater then
	 *  BUILDING then this slice cannot change.
	 */
	txContexts []transactionContext

	// caches that will be populated when respective contents are calculated
	cachedTxWithSigners []TransactionWithSigner
	cachedStxs          [][]byte
	cachedTxIDs         []string
}

func MakeAtomicTransactionComposer() AtomicTransactionComposer {
	return AtomicTransactionComposer{
		status:              BUILDING,
		txContexts:          []transactionContext{},
		cachedTxWithSigners: nil,
		cachedStxs:          nil,
		cachedTxIDs:         nil,
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
	return len(atc.txContexts)
}

/**
* Create a new composer with the same underlying transactions. The new composer's status will be
* BUILDING, so additional transactions may be added to it.
 */
func (atc *AtomicTransactionComposer) Clone() AtomicTransactionComposer {
	newTxContexts := make([]transactionContext, len(atc.txContexts))
	copy(newTxContexts, atc.txContexts)
	for i := range newTxContexts {
		newTxContexts[i].txn.Group = types.Digest{}
	}

	return AtomicTransactionComposer{
		status:     BUILDING,
		txContexts: newTxContexts,
	}
}

func (atc *AtomicTransactionComposer) validateTransaction(txn types.Transaction) error {
	emtpyGroup := types.Digest{}
	if txn.Group != emtpyGroup {
		return fmt.Errorf("expected empty group id")
	}

	return nil
}

/**
* Add a transaction to this atomic group.
*
* An error will be thrown if the composer's status is not BUILDING, or if adding this transaction
* causes the current group to exceed MAX_GROUP_SIZE.
 */
func (atc *AtomicTransactionComposer) AddTransaction(txnAndSigner TransactionWithSigner) error {
	if atc.status != BUILDING {
		return errors.New("status must be BUILDING in order to add tranactions")
	}

	if atc.Count() == MAX_GROUP_SIZE {
		return fmt.Errorf("reached max group size: %d", MAX_GROUP_SIZE)
	}

	err := atc.validateTransaction(txnAndSigner.Txn)
	if err != nil {
		return err
	}

	txContext := transactionContext{
		txn:    txnAndSigner.Txn,
		signer: txnAndSigner.Signer,
	}
	atc.txContexts = append(atc.txContexts, txContext)
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
	params AddMethodCallParams,
	/** A transaction signer that can authorize this application call from sender */
	signer TransactionSigner) error {

	if atc.status != BUILDING {
		return errors.New("status must be BUILDING in order to add transactions")
	}

	if len(params.MethodArgs) != len(params.Method.Args) {
		return fmt.Errorf("the incorrect number of arguments were provided: %d != %d", len(params.MethodArgs), len(params.Method.Args))
	}

	if atc.Count()+params.Method.GetTxCountFromMethod() > MAX_GROUP_SIZE {
		return fmt.Errorf("reached max group size: %d", MAX_GROUP_SIZE)
	}

	selectorValue := params.Method.GetSelector()
	encodedAbiArgs := [][]byte{selectorValue}
	var txsToAdd []TransactionWithSigner
	var abiArgs []interface{}
	var abiTypes []abi.Type
	for i, arg := range params.Method.Args {
		if _, ok := TransactionArgTypes[arg.AbiType]; ok {
			txnAndSigner, ok := params.MethodArgs[i].(TransactionWithSigner)
			if !ok {
				return fmt.Errorf("invalid arg type, expected transaction")
			}

			err := atc.validateTransaction(txnAndSigner.Txn)
			if err != nil {
				return err
			}
			txsToAdd = append(txsToAdd, txnAndSigner)
		} else {
			abiType, err := abi.TypeOf(arg.AbiType)
			if err != nil {
				return err
			}

			abiArgs = append(abiArgs, params.MethodArgs[i])
			abiTypes = append(abiTypes, abiType)
		}
	}

	// Up to 16 args can be passed to app call. First is reserved for method selector and last is
	// reserved for remaining args packaged into one tuple.
	if len(abiArgs) > 14 {
		var tupleTypes []abi.Type
		var tupleValues []interface{}
		for i, abiValue := range abiArgs[14:] {
			tupleTypes = append(tupleTypes, abiTypes[i])
			tupleValues = append(tupleValues, abiValue)
		}

		tupleType, err := abi.MakeTupleType(tupleTypes)
		if err != nil {
			return err
		}

		abiArgs = abiArgs[:14]
		abiArgs = append(abiArgs, tupleValues)
		abiTypes = abiTypes[:14]
		abiTypes = append(abiTypes, tupleType)
	}

	for i, abiArg := range abiArgs {
		encodedArg, err := abiTypes[i].Encode(abiArg)
		if err != nil {
			return err
		}

		encodedAbiArgs = append(encodedAbiArgs, encodedArg)
	}

	tx, err := MakeApplicationCallTx(
		params.AppID,
		encodedAbiArgs,
		nil,
		nil,
		nil,
		params.OnComplete,
		nil,
		nil,
		types.StateSchema{},
		types.StateSchema{},
		params.SuggestedParams,
		params.Sender,
		params.Note,
		types.Digest{},
		params.Lease,
		params.RekeyTo)

	if err != nil {
		return err
	}

	txAndSigner := TransactionWithSigner{
		Txn:    tx,
		Signer: signer,
	}

	for _, txAndSigner := range txsToAdd {
		txContext := transactionContext{
			txn:    txAndSigner.Txn,
			signer: txAndSigner.Signer,
		}
		atc.txContexts = append(atc.txContexts, txContext)
	}

	methodCallTxContext := transactionContext{
		txn:    txAndSigner.Txn,
		signer: txAndSigner.Signer,
		method: &params.Method,
	}
	atc.txContexts = append(atc.txContexts, methodCallTxContext)
	return nil
}

func (atc *AtomicTransactionComposer) getFinalizedTxWithSigners() []TransactionWithSigner {
	if atc.cachedTxWithSigners != nil {
		return atc.cachedTxWithSigners
	}

	atc.cachedTxWithSigners = make([]TransactionWithSigner, len(atc.txContexts))
	for i, txContext := range atc.txContexts {
		atc.cachedTxWithSigners[i] = TransactionWithSigner{
			Txn:    txContext.txn,
			Signer: txContext.signer,
		}
	}
	return atc.cachedTxWithSigners
}

/**
* Finalize the transaction group and returned the finalized transactions.
*
* The composer's status will be at least BUILT after executing this method.
 */
func (atc *AtomicTransactionComposer) BuildGroup() ([]TransactionWithSigner, error) {
	if atc.status > BUILDING {
		return atc.getFinalizedTxWithSigners(), nil
	}

	if atc.Count() == 0 {
		return nil, fmt.Errorf("attempting to build group with zero transactions")
	}

	var txns []types.Transaction
	for _, txContext := range atc.txContexts {
		txns = append(txns, txContext.txn)
	}

	gid, err := crypto.ComputeGroupID(txns)
	if err != nil {
		return nil, err
	}

	for i := range atc.txContexts {
		atc.txContexts[i].txn.Group = gid
	}

	atc.status = BUILT
	return atc.getFinalizedTxWithSigners(), nil
}

func (atc *AtomicTransactionComposer) getRawSignedTxs() [][]byte {
	if atc.cachedStxs != nil {
		return atc.cachedStxs
	}

	atc.cachedStxs = make([][]byte, len(atc.txContexts))
	for i, txContext := range atc.txContexts {
		atc.cachedStxs[i] = txContext.stxBytes
	}
	return atc.cachedStxs
}

/**
* Obtain signatures for each transaction in this group. If signatures have already been obtained,
* this method will return cached versions of the signatures.
*
* The composer's status will be at least SIGNED after executing this method.
*
* An error will be thrown if signing any of the transactions fails.
*
* @returns An array of signed transactions.
 */
func (atc *AtomicTransactionComposer) GatherSignatures() ([][]byte, error) {
	// if status is at least signed then return cached signed transactions
	if atc.status >= SIGNED {
		return atc.getRawSignedTxs(), nil
	}

	// retrieve built transactions and verify status is BUILT
	txsWithSigners, err := atc.BuildGroup()
	if err != nil {
		return nil, err
	}

	var txs []types.Transaction
	for _, txWithSigner := range txsWithSigners {
		txs = append(txs, txWithSigner.Txn)
	}

	visited := make([]bool, len(txs))
	rawSignedTxs := make([][]byte, len(txs))
	for i, txWithSigner := range txsWithSigners {
		if visited[i] {
			continue
		}

		var indexesToSign []int
		for j, other := range txsWithSigners {
			if !visited[j] && txWithSigner.Signer.Equals(other.Signer) {
				indexesToSign = append(indexesToSign, j)
				visited[j] = true
			}
		}

		sigStxs, err := txWithSigner.Signer.SignTransactions(txs, indexesToSign)
		if err != nil {
			return nil, err
		}

		for i, index := range indexesToSign {
			rawSignedTxs[index] = sigStxs[i]
		}
	}

	for i, stxBytes := range rawSignedTxs {
		atc.txContexts[i].stxBytes = stxBytes
	}
	atc.status = SIGNED
	return rawSignedTxs, nil
}

func (atc *AtomicTransactionComposer) getTxIDs() []string {
	if atc.cachedTxIDs != nil {
		return atc.cachedTxIDs
	}

	atc.cachedTxIDs = make([]string, len(atc.txContexts))
	for i, txContext := range atc.txContexts {
		atc.cachedTxIDs[i] = txContext.txID()
	}
	return atc.cachedTxIDs
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
* @returns a list of TxIDs of the submitted transactions.
 */
func (atc *AtomicTransactionComposer) Submit(client *algod.Client, ctx context.Context) ([]string, error) {
	if atc.status > SUBMITTED {
		return nil, errors.New("status must be SUBMITTED or lower in order to call Submit()")
	}

	stxs, err := atc.GatherSignatures()
	if err != nil {
		return nil, err
	}

	var serializedStxs []byte
	for _, stx := range stxs {
		serializedStxs = append(serializedStxs, stx...)
	}

	_, err = client.SendRawTransaction(serializedStxs).Do(ctx)
	if err != nil {
		return nil, err
	}

	atc.status = SUBMITTED
	return atc.getTxIDs(), nil
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
* @returns an object containing the confirmed round for this transaction, the txIDs of the submitted
*   transactions, and an array of results containing one element for each method call transaction in
*   this group. If a method has no return value (void), then the method results array will contain
*   null for that method's return value.
 */
func (atc *AtomicTransactionComposer) Execute(client *algod.Client, ctx context.Context, waitRounds uint64) (uint64, []string, []ABIResult, error) {
	if atc.status == COMMITTED {
		return 0, nil, nil, errors.New("status is already committed")
	}

	_, err := atc.Submit(client, ctx)
	if err != nil {
		return 0, nil, nil, err
	}

	txinfo, err := WaitForConfirmation(client, atc.txContexts[0].txID(), waitRounds, ctx)
	if err != nil {
		return 0, nil, nil, err
	}
	atc.status = COMMITTED

	var returnValues []ABIResult
	for _, txContext := range atc.txContexts {
		txid := txContext.txID()

		// Verify method call is available. This may not be the case if the App Call Tx wasn't created
		// by AddMethodCall().
		if txContext.isMethodCallTx() {
			continue
		}

		txinfo, _, err := client.PendingTransactionInformation(txid).Do(ctx)
		if err != nil {
			returnValues = append(returnValues, ABIResult{TxID: txid, RawReturnValue: []byte{}, ReturnValue: nil, DecodeError: err})
			continue
		}

		if txContext.method.Returns.AbiType == "void" {
			returnValues = append(returnValues, ABIResult{TxID: txid, RawReturnValue: []byte{}, ReturnValue: nil, DecodeError: nil})
			continue
		}

		failedToFindResult := true
		for j := len(txinfo.Logs) - 1; j >= 0; j-- {
			log := txinfo.Logs[j]
			if len(log) >= 4 && bytes.Equal(log[:4], ABI_RETURN_HASH) {
				failedToFindResult = false
				abiType, err := abi.TypeOf(txContext.method.Returns.AbiType)
				if err != nil {
					returnValues = append(returnValues, ABIResult{TxID: txid, RawReturnValue: log[4:], ReturnValue: nil, DecodeError: err})
					break
				}

				abiValue, err := abiType.Decode(log[4:])
				if err != nil {
					returnValues = append(returnValues, ABIResult{TxID: txid, RawReturnValue: log[4:], ReturnValue: nil, DecodeError: err})
					break
				}

				returnValues = append(returnValues, ABIResult{TxID: txid, RawReturnValue: log[4:], ReturnValue: abiValue, DecodeError: nil})
				break
			}
		}

		if failedToFindResult {
			returnValues = append(returnValues, ABIResult{
				TxID:           txid,
				RawReturnValue: nil,
				ReturnValue:    nil,
				DecodeError:    fmt.Errorf("no return value found when one was expected"),
			})
		}
	}

	return txinfo.ConfirmedRound, atc.getTxIDs(), returnValues, err
}
