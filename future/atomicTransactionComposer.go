package future

import (
	"context"
	"encoding/json"
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
type TransactionSigner interface {
	SignTransactions(txGroup []types.Transaction, indexesToSign []int) ([][]byte, []string, error)
	Equals(other TransactionSigner) bool
}

/**
 * TransactionSigner that can sign transactions for the provided basic Account.
 */
type BasicAccountTransactionSigner struct {
	Account crypto.Account
}

func (txSigner BasicAccountTransactionSigner) SignTransactions(txGroup []types.Transaction, indexesToSign []int) ([][]byte, []string, error) {
	stxs := make([][]byte, len(indexesToSign))
	txids := make([]string, len(indexesToSign))
	for i, pos := range indexesToSign {
		txid, stxBytes, err := crypto.SignTransaction(txSigner.Account.PrivateKey, txGroup[pos])
		if err != nil {
			return nil, nil, err
		}

		stxs[i] = stxBytes
		txids[i] = txid
	}

	return stxs, txids, nil
}

func (txSigner BasicAccountTransactionSigner) Equals(other TransactionSigner) bool {
	if castedSigner, ok := other.(BasicAccountTransactionSigner); ok {
		otherJson, err := json.Marshal(castedSigner)
		if err != nil {
			return false
		}

		selfJson, err := json.Marshal(txSigner)
		if err != nil {
			return false
		}

		return string(otherJson) == string(selfJson)
	}
	return false
}

/**
 * TransactionSigner that can sign transactions for the provided LogicSigAccount.
 */
type LogicSigAccountTransactionSigner struct {
	LogicSigAccount crypto.LogicSigAccount
}

func (txSigner LogicSigAccountTransactionSigner) SignTransactions(txGroup []types.Transaction, indexesToSign []int) ([][]byte, []string, error) {
	stxs := make([][]byte, len(indexesToSign))
	txids := make([]string, len(indexesToSign))
	for i, pos := range indexesToSign {
		txid, stxBytes, err := crypto.SignLogicSigAccountTransaction(txSigner.LogicSigAccount, txGroup[pos])
		if err != nil {
			return nil, nil, err
		}

		stxs[i] = stxBytes
		txids[i] = txid
	}

	return stxs, txids, nil
}

func (txSigner LogicSigAccountTransactionSigner) Equals(other TransactionSigner) bool {
	if castedSigner, ok := other.(LogicSigAccountTransactionSigner); ok {
		otherJson, err := json.Marshal(castedSigner)
		if err != nil {
			return false
		}

		selfJson, err := json.Marshal(txSigner)
		if err != nil {
			return false
		}

		return string(otherJson) == string(selfJson)
	}
	return false
}

/**
 * TransactionSigner that can sign transactions for the provided MultiSig Account
 */
type MultiSigAccountTransactionSigner struct {
	Msig crypto.MultisigAccount
	Sks  [][]byte
}

func (txSigner MultiSigAccountTransactionSigner) SignTransactions(txGroup []types.Transaction, indexesToSign []int) ([][]byte, []string, error) {
	stxs := make([][]byte, len(indexesToSign))
	txids := make([]string, len(indexesToSign))
	for i, pos := range indexesToSign {
		var unmergedStxs [][]byte
		var txid string
		for _, sk := range txSigner.Sks {
			tempTxid, unmergedStxBytes, err := crypto.SignMultisigTransaction(sk, txSigner.Msig, txGroup[pos])
			txid = tempTxid
			if err != nil {
				return nil, nil, err
			}

			unmergedStxs = append(unmergedStxs, unmergedStxBytes)
		}

		if len(txSigner.Sks) > 1 {
			tempTxid, stxBytes, err := crypto.MergeMultisigTransactions(unmergedStxs...)
			txid = tempTxid
			if err != nil {
				return nil, nil, err
			}

			stxs[i] = stxBytes
			txids[i] = txid
		} else {
			stxs[i] = unmergedStxs[0]
			txids[i] = txid
		}
	}

	return stxs, txids, nil
}

func (txSigner MultiSigAccountTransactionSigner) Equals(other TransactionSigner) bool {
	if castedSigner, ok := other.(MultiSigAccountTransactionSigner); ok {
		otherJson, err := json.Marshal(castedSigner)
		if err != nil {
			return false
		}

		selfJson, err := json.Marshal(txSigner)
		if err != nil {
			return false
		}

		return string(otherJson) == string(selfJson)
	}
	return false
}

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
	ReturnValue abi.Value
	/** If the SDK was unable to decode a return value, the error will be here. */
	DecodeError error
}

type AddMethoCallParams struct {
	/** The ID of the smart contract to call */
	AppID uint64
	/** The method to call on the smart contract */
	AbiMethod Method
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
	methodCalls map[int]Method

	/** The raw signed transactions populated after invocation of GatherSignatures(). */
	signedTxs [][]byte

	/** Txids of the transactions in this group. Populated after invocation to GatherSignatures() and used
	in Execute() to gather transaction information */
	txids map[int]string
}

func MakeAtomicTransactionComposer() AtomicTransactionComposer {
	return AtomicTransactionComposer{
		status:       BUILDING,
		transactions: []TransactionWithSigner{},
		methodCalls:  map[int]Method{},
		signedTxs:    [][]byte{},
		txids:        map[int]string{},
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
	newTxs := make([]TransactionWithSigner, len(atc.transactions))
	copy(newTxs, atc.transactions)
	for i := range newTxs {
		newTxs[i].Txn.Group = types.Digest{}
	}

	newMethodCalls := make(map[int]Method, len(atc.methodCalls))
	for k, v := range atc.methodCalls {
		newMethodCalls[k] = v
	}

	return AtomicTransactionComposer{
		status:       BUILDING,
		transactions: newTxs,
		methodCalls:  newMethodCalls,
		signedTxs:    [][]byte{},
		txids:        map[int]string{},
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
	params AddMethoCallParams,
	/** A transaction signer that can authorize this application call from sender */
	signer TransactionSigner) error {

	if atc.status != BUILDING {
		return errors.New("status must be BUILDING in order to add transactions")
	}

	if len(params.MethodArgs) != len(params.AbiMethod.Args) {
		return fmt.Errorf("the incorrect number of arguments were provided: %d != %d", len(params.MethodArgs), len(params.AbiMethod.Args))
	}

	if atc.Count()+params.AbiMethod.GetTxCountFromMethod() > MAX_GROUP_SIZE {
		return fmt.Errorf("reached max group size: %d", MAX_GROUP_SIZE)
	}

	selectorValue := params.AbiMethod.GetSelector()
	encodedAbiArgs := [][]byte{selectorValue}
	var txsToAdd []TransactionWithSigner
	var abiArgs []abi.Value
	for i, arg := range params.AbiMethod.Args {
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

			abiArgs = append(abiArgs, abi.Value{AbiType: abiType, RawValue: params.MethodArgs[i]})
		}
	}

	// Up to 16 args can be passed to app call. First is reserved for method selector and last is
	// reserved for remaining args packaged into one tuple.
	if len(abiArgs) > 14 {
		var tupleTypes []abi.Type
		var tupleValues []interface{}
		for _, abiValue := range abiArgs[14:] {
			tupleTypes = append(tupleTypes, abiValue.AbiType)
			tupleValues = append(tupleValues, abiValue.RawValue)
		}

		tupleType, err := abi.MakeTupleType(tupleTypes)
		if err != nil {
			return err
		}

		tupleValue := abi.Value{
			AbiType:  tupleType,
			RawValue: tupleValues,
		}

		abiArgs = abiArgs[:14]
		abiArgs = append(abiArgs, tupleValue)
	}

	for _, abiArg := range abiArgs {
		encodedArg, err := abiArg.AbiType.Encode(abiArg.RawValue)
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

	atc.transactions = append(atc.transactions, txsToAdd...)
	atc.methodCalls[len(atc.transactions)] = params.AbiMethod
	atc.transactions = append(atc.transactions, txAndSigner)
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

	for i := range atc.transactions {
		atc.transactions[i].Txn.Group = gid
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
	txsWithSigners, err := atc.BuildGroup()
	if err != nil {
		return nil, err
	}

	var txs []types.Transaction
	for _, txWithSigner := range txsWithSigners {
		txs = append(txs, txWithSigner.Txn)
	}

	visited := make([]bool, len(txs))
	txids := make(map[int]string)
	var signedTxs [][]byte
	for i, txAndSigner := range txsWithSigners {
		if visited[i] {
			continue
		}

		var indexesToSign []int
		for j, other := range txsWithSigners {
			if !visited[j] && txAndSigner.Signer.Equals(other.Signer) {
				indexesToSign = append(indexesToSign, j)
				visited[j] = true
			}
		}

		sigStxs, curTxids, err := txAndSigner.Signer.SignTransactions(txs, indexesToSign)
		if err != nil {
			return nil, err
		}

		signedTxs = append(signedTxs, sigStxs...)
		for i, index := range indexesToSign {
			txids[index] = curTxids[i]
		}
	}

	atc.signedTxs = signedTxs
	for index, txid := range txids {
		atc.txids[index] = txid
	}
	atc.status = SIGNED
	return atc.signedTxs, nil
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

	txids := make([]string, len(atc.txids))
	for _, v := range atc.txids {
		txids = append(txids, v)
	}

	atc.status = SUBMITTED
	return txids, nil
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

	if atc.status < SUBMITTED {
		_, err := atc.Submit(client, ctx)
		if err != nil {
			return 0, nil, nil, err
		}
	}

	txinfo, err := WaitForConfirmation(client, atc.txids[0], waitRounds, context.Background())
	if err != nil {
		return 0, nil, nil, err
	}
	atc.status = COMMITTED

	var returnValues []ABIResult
	var txids []string
	for i, txid := range atc.txids {
		if atc.transactions[i].Txn.Type != types.ApplicationCallTx {
			continue
		}

		// Verify method call is available. This may not be the case if the App Call Tx wasn't created
		// by AddMethodCall().
		if _, ok := atc.methodCalls[i]; !ok {
			continue
		}

		txids = append(txids, txid)
		txinfo, _, err := client.PendingTransactionInformation(txid).Do(context.Background())
		if err != nil {
			return 0, nil, nil, err
		}

		if atc.methodCalls[i].Returns.AbiType == "void" {
			returnValues = append(returnValues, ABIResult{TxID: txid, RawReturnValue: []byte{}, ReturnValue: abi.Value{}, DecodeError: nil})
			continue
		}

		failedToFindResult := true
		for j := len(txinfo.Logs) - 1; j >= 0; j-- {
			log := txinfo.Logs[j]
			if string(log[:4]) == string(ABI_RETURN_HASH) {
				failedToFindResult = false
				abiType, err := abi.TypeOf(atc.methodCalls[i].Returns.AbiType)
				if err != nil {
					returnValues = append(returnValues, ABIResult{TxID: txid, RawReturnValue: log[4:], ReturnValue: abi.Value{}, DecodeError: err})
					break
				}

				rawAbiValue, err := abiType.Decode(log[4:])
				if err != nil {
					returnValues = append(returnValues, ABIResult{TxID: txid, RawReturnValue: log[4:], ReturnValue: abi.Value{}, DecodeError: err})
					break
				}

				abiValue := abi.Value{
					AbiType:  abiType,
					RawValue: rawAbiValue,
				}
				returnValues = append(returnValues, ABIResult{TxID: txid, RawReturnValue: log[4:], ReturnValue: abiValue, DecodeError: nil})
				break
			}
		}

		if failedToFindResult {
			returnValues = append(returnValues, ABIResult{
				TxID:           txid,
				RawReturnValue: nil,
				ReturnValue:    abi.Value{},
				DecodeError:    fmt.Errorf("no return value found when one was expected"),
			})
		}
	}

	return txinfo.ConfirmedRound, txids, returnValues, err
}
