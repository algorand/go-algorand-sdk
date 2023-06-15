package test

import (
	"bytes"
	"context"
	"encoding/base32"
	"encoding/base64"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/algorand/go-algorand-sdk/v2/transaction"

	"golang.org/x/crypto/ed25519"

	"github.com/algorand/go-algorand-sdk/v2/abi"
	"github.com/algorand/go-algorand-sdk/v2/auction"
	"github.com/algorand/go-algorand-sdk/v2/client/kmd"
	algodV2 "github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	commonV2 "github.com/algorand/go-algorand-sdk/v2/client/v2/common"
	modelsV2 "github.com/algorand/go-algorand-sdk/v2/client/v2/common/models"
	indexerV2 "github.com/algorand/go-algorand-sdk/v2/client/v2/indexer"
	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/v2/logic"
	"github.com/algorand/go-algorand-sdk/v2/mnemonic"
	"github.com/algorand/go-algorand-sdk/v2/types"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

var txn types.Transaction
var stx []byte
var stxKmd []byte
var stxObj types.SignedTxn
var txid string
var account crypto.Account
var note []byte
var fee uint64
var fv uint64
var lv uint64
var to string
var gh []byte
var close string
var amt uint64
var gen string
var a types.Address
var msig crypto.MultisigAccount
var msigsig types.MultisigSig
var kcl kmd.Client
var aclv2 *algodV2.Client
var iclv2 *indexerV2.Client
var walletName string
var walletPswd string
var walletID string
var handle string
var versions []string
var msigExp kmd.ExportMultisigResponse
var pk string
var rekey string
var accounts []string
var e bool
var lastRound uint64
var sugParams types.SuggestedParams
var bid types.Bid
var sbid types.NoteField
var oldBid types.NoteField
var oldPk string
var newMn string
var mdk types.MasterDerivationKey
var microalgos types.MicroAlgos
var bytetxs [][]byte
var votekey string
var selkey string
var stateProofPK string
var votefst uint64
var votelst uint64
var votekd uint64
var nonpart bool
var data []byte
var sig types.Signature
var abiMethod abi.Method
var abiMethods []abi.Method
var abiJsonString string
var abiInterface abi.Interface
var abiContract abi.Contract
var txComposer transaction.AtomicTransactionComposer
var accountTxSigner transaction.BasicAccountTransactionSigner
var methodArgs []interface{}
var sigTxs [][]byte
var accountTxAndSigner transaction.TransactionWithSigner
var txTrace transaction.DryrunTxnResult
var trace string
var sourceMap logic.SourceMap
var srcMapping map[string]interface{}
var seeminglyProgram []byte
var sanityCheckError error

var assetTestFixture struct {
	Creator               string
	AssetIndex            uint64
	AssetName             string
	AssetUnitName         string
	AssetURL              string
	AssetMetadataHash     string
	ExpectedParams        modelsV2.AssetParams
	QueriedParams         modelsV2.AssetParams
	LastTransactionIssued types.Transaction
}

var tealCompleResult struct {
	status   int
	response modelsV2.CompileResponse
}

var tealDryrunResult struct {
	status   int
	response modelsV2.DryrunResponse
}

var opt = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "progress", // can define default values
}

// Dev mode helper functions
const devModeInitialAmount = 10_000_000

/**
 * waitForAlgodInDevMode is a Dev mode helper method
 * to wait for blocks to be finalized.
 * Since Dev mode produces blocks on a per transaction basis, it's possible
 * algod generates a block _before_ the corresponding SDK call to wait for a block.
 * Without _any_ wait, it's possible the SDK looks for the transaction before algod completes processing.
 * So, the method performs a local sleep to simulate waiting for a block.
 */
func waitForAlgodInDevMode() {
	time.Sleep(500 * time.Millisecond)
}

func initializeAccount(accountAddress string) error {
	params, err := aclv2.SuggestedParams().Do(context.Background())
	if err != nil {
		return err
	}

	txn, err = transaction.MakePaymentTxn(accounts[0], accountAddress, devModeInitialAmount, []byte{}, "", params)
	if err != nil {
		return err
	}

	res, err := kcl.SignTransaction(handle, walletPswd, txn)
	if err != nil {
		return err
	}

	_, err = aclv2.SendRawTransaction(res.SignedTransaction).Do(context.Background())
	if err != nil {
		return err
	}
	waitForAlgodInDevMode()
	return err
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opt)
}

func TestMain(m *testing.M) {
	flag.Parse()
	opt.Paths = flag.Args()

	status := godog.RunWithOptions("godogs", func(s *godog.Suite) {
		FeatureContext(s)
		AlgodClientV2Context(s)
		IndexerUnitTestContext(s)
		TransactionsUnitContext(s)
		ApplicationsContext(s)
		ApplicationsUnitContext(s)
		ResponsesContext(s)
	}, opt)

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

func FeatureContext(s *godog.Suite) {
	s.Step("I create a wallet", createWallet)
	s.Step("the wallet should exist", walletExist)
	s.Step("I get the wallet handle", getHandle)
	s.Step("I can get the master derivation key", getMdk)
	s.Step("I rename the wallet", renameWallet)
	s.Step("I can still get the wallet information with the same handle", getWalletInfo)
	s.Step("I renew the wallet handle", renewHandle)
	s.Step("I release the wallet handle", releaseHandle)
	s.Step("the wallet handle should not work", tryHandle)
	s.Step(`payment transaction parameters (\d+) (\d+) (\d+) "([^"]*)" "([^"]*)" "([^"]*)" (\d+) "([^"]*)" "([^"]*)"`, txnParams)
	s.Step(`mnemonic for private key "([^"]*)"`, mnForSk)
	s.Step(`multisig addresses "([^"]*)"`, msigAddresses)
	s.Step("I create the multisig payment transaction$", createMsigTxn)
	s.Step("I create the multisig payment transaction with zero fee", createMsigTxnZeroFee)
	s.Step("I sign the multisig transaction with the private key", signMsigTxn)
	s.Step("I sign the transaction with the private key", signTxn)
	s.Step(`^I add a rekeyTo field with address "([^"]*)"$`, iAddARekeyToFieldWithAddress)
	s.Step(`^I add a rekeyTo field with the private key algorand address$`, iAddARekeyToFieldWithThePrivateKeyAlgorandAddress)
	s.Step(`^I set the from address to "([^"]*)"$`, iSetTheFromAddressTo)
	s.Step(`the signed transaction should equal the golden "([^"]*)"`, equalGolden)
	s.Step(`the multisig transaction should equal the golden "([^"]*)"`, equalMsigGolden)
	s.Step(`the multisig address should equal the golden "([^"]*)"`, equalMsigAddrGolden)
	s.Step("I get versions with algod", aclV)
	s.Step("v1 should be in the versions", v1InVersions)
	s.Step("v2 should be in the versions", v2InVersions)
	s.Step("I get versions with kmd", kclV)
	s.Step("I import the multisig", importMsig)
	s.Step("the multisig should be in the wallet", msigInWallet)
	s.Step("I export the multisig", expMsig)
	s.Step("the multisig should equal the exported multisig", msigEq)
	s.Step("I delete the multisig", deleteMsig)
	s.Step("the multisig should not be in the wallet", msigNotInWallet)
	s.Step("^I generate a key using kmd for rekeying and fund it$", genRekeyKmd)
	s.Step("^I generate a key using kmd$", genKeyKmd)
	s.Step("the key should be in the wallet", keyInWallet)
	s.Step("I delete the key", deleteKey)
	s.Step("the key should not be in the wallet", keyNotInWallet)
	s.Step("I generate a key", genKey)
	s.Step("I import the key", importKey)
	s.Step("the private key should be equal to the exported private key", skEqExport)
	s.Step("a kmd client", kmdClient)
	s.Step("wallet information", walletInfo)
	s.Step(`default transaction with parameters (\d+) "([^"]*)"`, defaultTxn)
	s.Step(`default multisig transaction with parameters (\d+) "([^"]*)"`, defaultMsigTxn)
	s.Step("I get the private key", getSk)
	s.Step("I send the transaction", sendTxn)
	s.Step("I send the kmd-signed transaction", sendTxnKmd)
	s.Step("I send the bogus kmd-signed transaction", sendTxnKmdFailureExpected)
	s.Step("I send the multisig transaction", sendMsigTxn)
	s.Step("the transaction should not go through", txnFail)
	s.Step("I sign the transaction with kmd", signKmd)
	s.Step("the signed transaction should equal the kmd signed transaction", signBothEqual)
	s.Step("I sign the multisig transaction with kmd", signMsigKmd)
	s.Step("the multisig transaction should equal the kmd signed multisig transaction", signMsigBothEqual)
	s.Step(`^the node should be healthy`, nodeHealth)
	s.Step(`^I get the ledger supply`, ledger)
	s.Step(`^I create a bid`, createBid)
	s.Step(`^I encode and decode the bid`, encDecBid)
	s.Step(`^the bid should still be the same`, checkBid)
	s.Step(`^I decode the address`, decAddr)
	s.Step(`^I encode the address`, encAddr)
	s.Step(`^the address should still be the same`, checkAddr)
	s.Step(`^I convert the private key back to a mnemonic`, skToMn)
	s.Step(`^the mnemonic should still be the same as "([^"]*)"`, checkMn)
	s.Step(`^mnemonic for master derivation key "([^"]*)"`, mnToMdk)
	s.Step(`^I convert the master derivation key back to a mnemonic`, mdkToMn)
	s.Step(`^I create the flat fee payment transaction`, createTxnFlat)
	s.Step(`^encoded multisig transaction "([^"]*)"`, encMsigTxn)
	s.Step(`^I append a signature to the multisig transaction`, appendMsig)
	s.Step(`^encoded multisig transactions "([^"]*)"`, encMtxs)
	s.Step(`^I merge the multisig transactions`, mergeMsig)
	s.Step(`^I convert (\d+) microalgos to algos and back`, microToAlgos)
	s.Step(`^it should still be the same amount of microalgos (\d+)`, checkAlgos)
	s.Step("I sign the bid", signBid)
	s.Step(`default V2 key registration transaction "([^"]*)"`, createKeyregWithStateProof)
	s.Step(`^I can get account information`, newAccInfo)
	s.Step("asset test fixture", createAssetTestFixture)
	s.Step(`^default asset creation transaction with total issuance (\d+)$`, defaultAssetCreateTxn)
	s.Step(`^I update the asset index$`, getAssetIndex)
	s.Step(`^I get the asset info$`, getAssetInfo)
	s.Step(`^I should be unable to get the asset info`, failToGetAssetInfo)
	s.Step(`^the asset info should match the expected asset info$`, checkExpectedVsActualAssetParams)
	s.Step(`^I create a no-managers asset reconfigure transaction$`, createNoManagerAssetReconfigure)
	s.Step(`^I create an asset destroy transaction$`, createAssetDestroy)
	s.Step(`^I create a transaction for a second account, signalling asset acceptance$`, createAssetAcceptanceForSecondAccount)
	s.Step(`^I create a transaction transferring (\d+) assets from creator to a second account$`, createAssetTransferTransactionToSecondAccount)
	s.Step(`^the creator should have (\d+) assets remaining$`, theCreatorShouldHaveAssetsRemaining)
	s.Step(`^I create a freeze transaction targeting the second account$`, createFreezeTransactionTargetingSecondAccount)
	s.Step(`^I create a transaction transferring (\d+) assets from a second account to creator$`, createAssetTransferTransactionFromSecondAccountToCreator)
	s.Step(`^I create an un-freeze transaction targeting the second account$`, createUnfreezeTransactionTargetingSecondAccount)
	s.Step(`^default-frozen asset creation transaction with total issuance (\d+)$`, defaultAssetCreateTxnWithDefaultFrozen)
	s.Step(`^I create a transaction revoking (\d+) assets from a second account to creator$`, createRevocationTransaction)
	s.Step(`^I create a transaction transferring <amount> assets from creator to a second account$`, iCreateATransactionTransferringAmountAssetsFromCreatorToASecondAccount) // provide handler for when godog misreads
	s.Step(`^base64 encoded data to sign "([^"]*)"$`, baseEncodedDataToSign)
	s.Step(`^program hash "([^"]*)"$`, programHash)
	s.Step(`^I perform tealsign$`, iPerformTealsign)
	s.Step(`^the signature should be equal to "([^"]*)"$`, theSignatureShouldBeEqualTo)
	s.Step(`^base64 encoded program "([^"]*)"$`, baseEncodedProgram)
	s.Step(`^base64 encoded private key "([^"]*)"$`, baseEncodedPrivateKey)
	s.Step("an algod v2 client$", algodClientV2)
	s.Step("an indexer v2 client$", indexerClientV2)
	s.Step(`^I compile a teal program "([^"]*)"$`, tealCompile)
	s.Step(`^it is compiled with (\d+) and "([^"]*)" and "([^"]*)"$`, tealCheckCompile)
	s.Step(`^base64 decoding the response is the same as the binary "([^"]*)"$`, tealCheckCompileAgainstFile)
	s.Step(`^I dryrun a "([^"]*)" program "([^"]*)"$`, tealDryrun)
	s.Step(`^I get execution result "([^"]*)"$`, tealCheckDryrun)
	s.Step(`^I create the Method object from method signature "([^"]*)"$`, createMethodObjectFromSignature)
	s.Step(`^I serialize the Method object into json$`, serializeMethodObjectIntoJson)
	s.Step(`^the produced json should equal "([^"]*)" loaded from "([^"]*)"$`, checkSerializedMethodObject)
	s.Step(`^I create the Method object with name "([^"]*)" first argument type "([^"]*)" second argument type "([^"]*)" and return type "([^"]*)"$`, createMethodObjectFromProperties)
	s.Step(`^I create the Method object with name "([^"]*)" first argument name "([^"]*)" first argument type "([^"]*)" second argument name "([^"]*)" second argument type "([^"]*)" and return type "([^"]*)"$`, createMethodObjectWithArgNames)
	s.Step(`^I create the Method object with name "([^"]*)" method description "([^"]*)" first argument type "([^"]*)" first argument description "([^"]*)" second argument type "([^"]*)" second argument description "([^"]*)" and return type "([^"]*)"$`, createMethodObjectWithDescription)
	s.Step(`^the txn count should be (\d+)$`, checkTxnCount)
	s.Step(`^the method selector should be "([^"]*)"$`, checkMethodSelector)
	s.Step(`^I create an Interface object from the Method object with name "([^"]*)" and description "([^"]*)"$`, createInterfaceObject)
	s.Step(`^I serialize the Interface object into json$`, serializeInterfaceObjectIntoJson)
	s.Step(`^I create a Contract object from the Method object with name "([^"]*)" and description "([^"]*)"$`, createContractObject)
	s.Step(`^I set the Contract\'s appID to (\d+) for the network "([^"]*)"$`, iSetTheContractsAppIDToForTheNetwork)
	s.Step(`^I serialize the Contract object into json$`, serializeContractObjectIntoJson)
	s.Step(`^the deserialized json should equal the original Method object`, deserializeMethodJson)
	s.Step(`^the deserialized json should equal the original Interface object`, deserializeInterfaceJson)
	s.Step(`^the deserialized json should equal the original Contract object`, deserializeContractJson)
	s.Step(`^a new AtomicTransactionComposer$`, aNewAtomicTransactionComposer)
	s.Step(`^suggested transaction parameters fee (\d+), flat-fee "([^"]*)", first-valid (\d+), last-valid (\d+), genesis-hash "([^"]*)", genesis-id "([^"]*)"$`, suggestedTransactionParameters)
	s.Step(`^an application id (\d+)$`, anApplicationId)
	s.Step(`^I make a transaction signer for the ([^"]*) account\.$`, iMakeATransactionSignerForTheAccount)
	s.Step(`^I create a new method arguments array\.$`, iCreateANewMethodArgumentsArray)
	s.Step(`^I append the encoded arguments "([^"]*)" to the method arguments array\.$`, iAppendTheEncodedArgumentsToTheMethodArgumentsArray)
	s.Step(`^I add a method call with the ([^"]*) account, the current application, suggested params, on complete "([^"]*)", current transaction signer, current method arguments\.$`, addMethodCall)
	s.Step(`^I add a method call with the ([^"]*) account, the current application, suggested params, on complete "([^"]*)", current transaction signer, current method arguments, approval-program "([^"]*)", clear-program "([^"]*)"\.$`, addMethodCallForUpdate)
	s.Step(`^I add a method call with the ([^"]*) account, the current application, suggested params, on complete "([^"]*)", current transaction signer, current method arguments, approval-program "([^"]*)", clear-program "([^"]*)", global-bytes (\d+), global-ints (\d+), local-bytes (\d+), local-ints (\d+), extra-pages (\d+)\.$`, addMethodCallForCreate)
	s.Step(`^I add a nonced method call with the ([^"]*) account, the current application, suggested params, on complete "([^"]*)", current transaction signer, current method arguments\.$`, addMethodCallWithNonce)
	s.Step(`^I add the nonce "([^"]*)"$`, iAddTheNonce)
	s.Step(`^I build the transaction group with the composer\. If there is an error it is "([^"]*)"\.$`, buildTheTransactionGroupWithTheComposer)
	s.Step(`^The composer should have a status of "([^"]*)"\.$`, theComposerShouldHaveAStatusOf)
	s.Step(`^I gather signatures with the composer\.$`, iGatherSignaturesWithTheComposer)
	s.Step(`^the base64 encoded signed transactions should equal "([^"]*)"$`, theBaseEncodedSignedTransactionsShouldEqual)
	s.Step(`^I build a payment transaction with sender "([^"]*)", receiver "([^"]*)", amount (\d+), close remainder to "([^"]*)"$`, iBuildAPaymentTransactionWithSenderReceiverAmountCloseRemainderTo)
	s.Step(`^I create a transaction with signer with the current transaction\.$`, iCreateATransactionWithSignerWithTheCurrentTransaction)
	s.Step(`^I append the current transaction with signer to the method arguments array\.$`, iAppendTheCurrentTransactionWithSignerToTheMethodArgumentsArray)
	s.Step(`^a dryrun response file "([^"]*)" and a transaction at index "([^"]*)"$`, aDryrunResponseFileAndATransactionAtIndex)
	s.Step(`^calling app trace produces "([^"]*)"$`, callingAppTraceProduces)
	s.Step(`^I append to my Method objects list in the case of a non-empty signature "([^"]*)"$`, iAppendToMyMethodObjectsListInTheCaseOfANonemptySignature)
	s.Step(`^I create an Interface object from my Method objects list$`, iCreateAnInterfaceObjectFromMyMethodObjectsList)
	s.Step(`^I create a Contract object from my Method objects list$`, iCreateAContractObjectFromMyMethodObjectsList)
	s.Step(`^I get the method from the Interface by name "([^"]*)"$`, iGetTheMethodFromTheInterfaceByName)
	s.Step(`^I get the method from the Contract by name "([^"]*)"$`, iGetTheMethodFromTheContractByName)
	s.Step(`^the produced method signature should equal "([^"]*)"\. If there is an error it begins with "([^"]*)"$`, theProducedMethodSignatureShouldEqualIfThereIsAnErrorItBeginsWith)
	s.Step(`^a source map json file "([^"]*)"$`, aSourceMapJsonFile)
	s.Step(`^the string composed of pc:line number equals "([^"]*)"$`, theStringComposedOfPclineNumberEquals)
	s.Step(`^I compile a teal program "([^"]*)" with mapping enabled$`, iCompileATealProgramWithMappingEnabled)
	s.Step(`^the resulting source map is the same as the json "([^"]*)"$`, theResultingSourceMapIsTheSameAsTheJson)
	s.Step(`^getting the line associated with a pc "([^"]*)" equals "([^"]*)"$`, gettingTheLineAssociatedWithAPcEquals)
	s.Step(`^getting the last pc associated with a line "([^"]*)" equals "([^"]*)"$`, gettingTheLastPcAssociatedWithALineEquals)
	s.Step(`^a base64 encoded program bytes for heuristic sanity check "([^"]*)"$`, takeB64encodedBytes)
	s.Step(`^I start heuristic sanity check over the bytes$`, heuristicCheckOverBytes)
	s.Step(`^if the heuristic sanity check throws an error, the error contains "([^"]*)"$`, checkErrorIfMatching)
	s.Step(`^disassembly of "([^"]*)" matches "([^"]*)"$`, disassemblyMatches)

	s.BeforeScenario(func(interface{}) {
		stxObj = types.SignedTxn{}
		abiMethods = nil
		kcl.RenewWalletHandle(handle)
	})
}

func createWallet() error {
	walletName = "Walletgo"
	walletPswd = ""
	resp, err := kcl.CreateWallet(walletName, walletPswd, "sqlite", types.MasterDerivationKey{})
	if err != nil {
		return err
	}
	walletID = resp.Wallet.ID
	return nil
}

func walletExist() error {
	wallets, err := kcl.ListWallets()
	if err != nil {
		return err
	}
	for _, w := range wallets.Wallets {
		if w.Name == walletName {
			return nil
		}
	}
	return fmt.Errorf("Wallet not found")
}

func getHandle() error {
	h, err := kcl.InitWalletHandle(walletID, walletPswd)
	if err != nil {
		return err
	}
	handle = h.WalletHandleToken
	return nil
}

func getMdk() error {
	_, err := kcl.ExportMasterDerivationKey(handle, walletPswd)
	return err
}

func renameWallet() error {
	walletName = "Walletgo_new"
	_, err := kcl.RenameWallet(walletID, walletPswd, walletName)
	return err
}

func getWalletInfo() error {
	resp, err := kcl.GetWallet(handle)
	if resp.WalletHandle.Wallet.Name != walletName {
		return fmt.Errorf("Wallet name not equal")
	}
	return err
}

func renewHandle() error {
	_, err := kcl.RenewWalletHandle(handle)
	return err
}

func releaseHandle() error {
	_, err := kcl.ReleaseWalletHandle(handle)
	return err
}

func tryHandle() error {
	_, err := kcl.RenewWalletHandle(handle)
	if err == nil {
		return fmt.Errorf("should be an error; handle was released")
	}
	return nil
}

func iAddARekeyToFieldWithThePrivateKeyAlgorandAddress() error {
	pk, err := crypto.GenerateAddressFromSK(account.PrivateKey)
	if err != nil {
		return err
	}

	err = txn.Rekey(pk.String())

	if err != nil {
		return err
	}
	return nil
}

func txnParams(ifee, ifv, ilv int, igh, ito, iclose string, iamt int, igen, inote string) error {
	var err error
	if inote != "none" {
		note, err = base64.StdEncoding.DecodeString(inote)
		if err != nil {
			return err
		}
	} else {
		note, err = base64.StdEncoding.DecodeString("")
		if err != nil {
			return err
		}
	}
	gh, err = base64.StdEncoding.DecodeString(igh)
	if err != nil {
		return err
	}
	to = ito
	fee = uint64(ifee)
	fv = uint64(ifv)
	lv = uint64(ilv)
	if iclose != "none" {
		close = iclose
	} else {
		close = ""
	}
	amt = uint64(iamt)
	if igen != "none" {
		gen = igen
	} else {
		gen = ""
	}
	if err != nil {
		return err
	}
	return nil
}

func mnForSk(mn string) error {
	sk, err := mnemonic.ToPrivateKey(mn)
	if err != nil {
		return err
	}
	account.PrivateKey = sk
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(sk.Public())
	if err != nil {
		return err
	}
	addr := buf.Bytes()[4:]

	n := copy(a[:], addr)
	if n != 32 {
		return fmt.Errorf("wrong address bytes length")
	}
	return err
}
func msigAddresses(addresses string) error {
	var err error
	addrlist := strings.Fields(addresses)

	var addrStructs []types.Address
	for _, a := range addrlist {
		addr, err := types.DecodeAddress(a)
		if err != nil {
			return err
		}

		addrStructs = append(addrStructs, addr)
	}
	msig, err = crypto.MultisigAccountWithParams(1, 2, addrStructs)

	return err
}

func iSetTheFromAddressTo(address string) error {
	addr, err := types.DecodeAddress(address)
	if err != nil {
		return err
	}
	txn.Sender = addr
	return nil
}

func createMsigTxn() error {
	var err error
	paramsToUse := types.SuggestedParams{
		Fee:             types.MicroAlgos(fee),
		GenesisID:       gen,
		GenesisHash:     gh,
		FirstRoundValid: types.Round(fv),
		LastRoundValid:  types.Round(lv),
		FlatFee:         false,
	}
	msigaddr, _ := msig.Address()
	txn, err = transaction.MakePaymentTxn(msigaddr.String(), to, amt, note, close, paramsToUse)
	if err != nil {
		return err
	}
	return err

}

func createMsigTxnZeroFee() error {
	var err error
	paramsToUse := types.SuggestedParams{
		Fee:             types.MicroAlgos(fee),
		GenesisID:       gen,
		GenesisHash:     gh,
		FirstRoundValid: types.Round(fv),
		LastRoundValid:  types.Round(lv),
		FlatFee:         true,
	}
	msigaddr, _ := msig.Address()
	txn, err = transaction.MakePaymentTxn(msigaddr.String(), to, amt, note, close, paramsToUse)
	if err != nil {
		return err
	}
	return err

}

func signMsigTxn() error {
	var err error
	txid, stx, err = crypto.SignMultisigTransaction(account.PrivateKey, msig, txn)

	return err
}

func signTxn() error {
	var err error
	txid, stx, err = crypto.SignTransaction(account.PrivateKey, txn)
	if err != nil {
		return err
	}
	return nil
}

func iAddARekeyToFieldWithAddress(address string) error {
	err := txn.Rekey(address)
	if err != nil {
		return err
	}
	return nil
}

func equalGolden(golden string) error {
	goldenDecoded, err := base64.StdEncoding.DecodeString(golden)
	if err != nil {
		return err
	}

	if !bytes.Equal(goldenDecoded, stx) {
		return fmt.Errorf(base64.StdEncoding.EncodeToString(stx))
	}
	return nil
}

func equalMsigAddrGolden(golden string) error {
	msigAddr, err := msig.Address()
	if err != nil {
		return err
	}
	if golden != msigAddr.String() {
		return fmt.Errorf("NOT EQUAL")
	}
	return nil
}

func equalMsigGolden(golden string) error {
	goldenDecoded, err := base64.StdEncoding.DecodeString(golden)
	if err != nil {
		return err
	}
	if !bytes.Equal(goldenDecoded, stx) {
		return fmt.Errorf("NOT EQUAL")
	}
	return nil
}

func aclV() error {
	v, err := aclv2.Versions().Do(context.Background())
	if err != nil {
		return err
	}
	versions = v.Versions
	return nil
}

func v1InVersions() error {
	for _, b := range versions {
		if b == "v1" {
			return nil
		}
	}
	return fmt.Errorf("v1 not found")
}

func v2InVersions() error {
	for _, b := range versions {
		if b == "v2" {
			return nil
		}
	}
	return fmt.Errorf("v2 not found")
}

func kclV() error {
	v, err := kcl.Version()
	versions = v.Versions
	return err
}

func importMsig() error {
	_, err := kcl.ImportMultisig(handle, msig.Version, msig.Threshold, msig.Pks)
	return err
}

func msigInWallet() error {
	msigs, err := kcl.ListMultisig(handle)
	if err != nil {
		return err
	}
	addrs := msigs.Addresses
	for _, a := range addrs {
		addr, err := msig.Address()
		if err != nil {
			return err
		}
		if a == addr.String() {
			return nil
		}
	}
	return fmt.Errorf("msig not found")

}

func expMsig() error {
	addr, err := msig.Address()
	if err != nil {
		return err
	}
	msigExp, err = kcl.ExportMultisig(handle, walletPswd, addr.String())

	return err
}

func msigEq() error {
	eq := true

	if (msig.Pks == nil) != (msigExp.PKs == nil) {
		eq = false
	}

	if len(msig.Pks) != len(msigExp.PKs) {
		eq = false
	}

	for i := range msig.Pks {

		if !bytes.Equal(msig.Pks[i], msigExp.PKs[i]) {
			eq = false
		}
	}

	if !eq {
		return fmt.Errorf("exported msig not equal to original msig")
	}
	return nil
}

func deleteMsig() error {
	addr, err := msig.Address()
	kcl.DeleteMultisig(handle, walletPswd, addr.String())
	return err
}

func msigNotInWallet() error {
	msigs, err := kcl.ListMultisig(handle)
	if err != nil {
		return err
	}
	addrs := msigs.Addresses
	for _, a := range addrs {
		addr, err := msig.Address()
		if err != nil {
			return err
		}
		if a == addr.String() {
			return fmt.Errorf("msig found unexpectedly; should have been deleted")
		}
	}
	return nil

}

func genKeyKmd() error {
	p, err := kcl.GenerateKey(handle)
	if err != nil {
		return err
	}
	pk = p.Address
	return nil
}

func genRekeyKmd() error {
	p, err := kcl.GenerateKey(handle)
	if err != nil {
		return err
	}
	rekey = p.Address
	initializeAccount(rekey)
	return nil
}

func keyInWallet() error {
	resp, err := kcl.ListKeys(handle)
	if err != nil {
		return err
	}
	for _, a := range resp.Addresses {
		if pk == a {
			return nil
		}
	}
	return fmt.Errorf("key not found")
}

func deleteKey() error {
	_, err := kcl.DeleteKey(handle, walletPswd, pk)
	return err
}

func keyNotInWallet() error {
	resp, err := kcl.ListKeys(handle)
	if err != nil {
		return err
	}
	for _, a := range resp.Addresses {
		if pk == a {
			return fmt.Errorf("key found unexpectedly; should have been deleted")
		}
	}
	return nil
}

func genKey() error {
	account = crypto.GenerateAccount()
	a = account.Address
	pk = a.String()
	return nil
}

func importKey() error {
	_, err := kcl.ImportKey(handle, account.PrivateKey)
	return err
}

func skEqExport() error {
	exp, err := kcl.ExportKey(handle, walletPswd, a.String())
	if err != nil {
		return err
	}
	kcl.DeleteKey(handle, walletPswd, a.String())
	if bytes.Equal(exp.PrivateKey.Seed(), account.PrivateKey.Seed()) {
		return nil
	}
	return fmt.Errorf("private keys not equal")
}

func kmdClient() error {
	kmdToken := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	kmdAddress := "http://localhost:" + "60001"
	var err error
	kcl, err = kmd.MakeClient(kmdAddress, kmdToken)
	return err
}

func algodClientV2() error {
	algodToken := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	algodAddress := "http://localhost:" + "60000"
	var err error
	aclv2, err = algodV2.MakeClient(algodAddress, algodToken)
	algodV2client = aclv2
	if err != nil {
		return err
	}
	_, err = aclv2.StatusAfterBlock(1).Do(context.Background())
	return err
}

func indexerClientV2() error {
	indexerAddress := "http://localhost:" + "59999"
	var err error
	iclv2, err = indexerV2.MakeClient(indexerAddress, "")
	indexerV2client = iclv2
	return err
}

func walletInfo() error {
	walletName = "unencrypted-default-wallet"
	walletPswd = ""
	wallets, err := kcl.ListWallets()
	if err != nil {
		return err
	}
	for _, w := range wallets.Wallets {
		if w.Name == walletName {
			walletID = w.ID
		}
	}
	h, err := kcl.InitWalletHandle(walletID, walletPswd)
	if err != nil {
		return err
	}
	handle = h.WalletHandleToken
	accs, err := kcl.ListKeys(handle)
	accounts = accs.Addresses
	return err
}

// Helper function for making default transactions.
func defaultTxnWithAddress(iamt int, inote string, senderAddress string) error {
	var err error
	if inote != "none" {
		note, err = base64.StdEncoding.DecodeString(inote)
		if err != nil {
			return err
		}
	} else {
		note, err = base64.StdEncoding.DecodeString("")
		if err != nil {
			return err
		}
	}

	amt = uint64(iamt)
	pk = senderAddress
	params, err := aclv2.SuggestedParams().Do(context.Background())
	if err != nil {
		return err
	}
	lastRound = uint64(params.FirstRoundValid)
	txn, err = transaction.MakePaymentTxn(senderAddress, accounts[1], amt, note, "", params)
	return err
}

func defaultTxn(iamt int, inote string) error {
	return defaultTxnWithAddress(iamt, inote, accounts[0])
}

func defaultMsigTxn(iamt int, inote string) error {
	var err error
	if inote != "none" {
		note, err = base64.StdEncoding.DecodeString(inote)
		if err != nil {
			return err
		}
	} else {
		note, err = base64.StdEncoding.DecodeString("")
		if err != nil {
			return err
		}
	}

	amt = uint64(iamt)
	pk = accounts[0]

	var addrStructs []types.Address
	for _, a := range accounts {
		addr, err := types.DecodeAddress(a)
		if err != nil {
			return err
		}

		addrStructs = append(addrStructs, addr)
	}

	msig, err = crypto.MultisigAccountWithParams(1, 1, addrStructs)
	if err != nil {
		return err
	}
	params, err := aclv2.SuggestedParams().Do(context.Background())
	if err != nil {
		return err
	}
	lastRound = uint64(params.FirstRoundValid)
	addr, err := msig.Address()
	if err != nil {
		return err
	}
	txn, err = transaction.MakePaymentTxn(addr.String(), accounts[1], amt, note, "", params)
	if err != nil {
		return err
	}
	return nil
}

func getSk() error {
	sk, err := kcl.ExportKey(handle, walletPswd, pk)
	if err != nil {
		return err
	}
	account.PrivateKey = sk.PrivateKey
	return nil
}

func sendTxn() error {
	tx, err := aclv2.SendRawTransaction(stx).Do(context.Background())
	if err != nil {
		return err
	}
	txid = tx
	return nil
}

func sendTxnKmd() error {
	var err error
	txid, err = aclv2.SendRawTransaction(stxKmd).Do(context.Background())
	return err
}

func sendTxnKmdFailureExpected() error {
	tx, err := aclv2.SendRawTransaction(stxKmd).Do(context.Background())
	if err == nil {
		e = false
		return fmt.Errorf("expected an error when sending kmd-signed transaction but no error occurred")
	}
	e = true
	txid = tx
	return nil
}

func sendMsigTxn() error {
	_, err := aclv2.SendRawTransaction(stx).Do(context.Background())
	if err != nil {
		e = true
	}

	return nil
}

func txnFail() error {
	if e {
		return nil
	}
	return fmt.Errorf("sending the transaction should have failed")
}

func signKmd() error {
	s, err := kcl.SignTransaction(handle, walletPswd, txn)
	if err != nil {
		return err
	}
	stxKmd = s.SignedTransaction
	return nil
}

func signBothEqual() error {
	if bytes.Equal(stx, stxKmd) {
		return nil
	}
	return fmt.Errorf("signed transactions not equal")
}

func signMsigKmd() error {
	kcl.ImportMultisig(handle, msig.Version, msig.Threshold, msig.Pks)
	decoded, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(pk)
	if err != nil {
		return fmt.Errorf("signMsigKmd: %w", err)
	}
	s, err := kcl.MultisigSignTransaction(handle, walletPswd, txn, decoded[:32], types.MultisigSig{})
	if err != nil {
		return err
	}
	msgpack.Decode(s.Multisig, &msigsig)
	stxObj.Msig = msigsig
	stxObj.Sig = types.Signature{}
	stxObj.Txn = txn
	stxKmd = msgpack.Encode(stxObj)
	return nil
}

func signMsigBothEqual() error {
	addr, err := msig.Address()
	if err != nil {
		return err
	}
	kcl.DeleteMultisig(handle, walletPswd, addr.String())
	if bytes.Equal(stx, stxKmd) {
		return nil
	}
	return fmt.Errorf("signed transactions not equal")

}

func nodeHealth() error {
	err := aclv2.HealthCheck().Do(context.Background())
	return err
}

func ledger() error {
	_, err := aclv2.Supply().Do(context.Background())
	return err
}

func createBid() error {
	var err error
	account = crypto.GenerateAccount()
	bid, err = auction.MakeBid(account.Address.String(), 1, 2, 3, account.Address.String(), 4)
	return err
}

func encDecBid() error {
	temp := msgpack.Encode(sbid)
	err := msgpack.Decode(temp, &sbid)
	return err
}

func signBid() error {
	signedBytes, err := crypto.SignBid(account.PrivateKey, bid)
	if err != nil {
		return err
	}
	err = msgpack.Decode(signedBytes, &sbid)
	if err != nil {
		return err
	}
	err = msgpack.Decode(signedBytes, &oldBid)
	return err
}

func checkBid() error {
	if sbid != oldBid {
		return fmt.Errorf("bid should still be the same")
	}
	return nil
}

func decAddr() error {
	var err error
	oldPk = pk
	a, err = types.DecodeAddress(pk)
	return err
}

func encAddr() error {
	pk = a.String()
	return nil
}

func checkAddr() error {
	if pk != oldPk {
		return fmt.Errorf("A decoded and encoded address should equal the original address")
	}
	return nil
}

func skToMn() error {
	var err error
	newMn, err = mnemonic.FromPrivateKey(account.PrivateKey)
	return err
}

func checkMn(mn string) error {
	if mn != newMn {
		return fmt.Errorf("the mnemonic should equal the original mnemonic")
	}
	return nil
}

func mnToMdk(mn string) error {
	var err error
	mdk, err = mnemonic.ToMasterDerivationKey(mn)
	return err
}

func mdkToMn() error {
	var err error
	newMn, err = mnemonic.FromMasterDerivationKey(mdk)
	return err
}

func createTxnFlat() error {
	var err error
	paramsToUse := types.SuggestedParams{
		Fee:             types.MicroAlgos(fee),
		GenesisID:       gen,
		GenesisHash:     gh,
		FirstRoundValid: types.Round(fv),
		LastRoundValid:  types.Round(lv),
		FlatFee:         true,
	}
	txn, err = transaction.MakePaymentTxn(a.String(), to, amt, note, close, paramsToUse)
	if err != nil {
		return err
	}
	return err
}

func encMsigTxn(encoded string) error {
	var err error
	stx, err = base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return err
	}
	err = msgpack.Decode(stx, &stxObj)
	return err
}

func appendMsig() error {
	var err error
	msig, err = crypto.MultisigAccountFromSig(stxObj.Msig)
	if err != nil {
		return err
	}
	_, stx, err = crypto.AppendMultisigTransaction(account.PrivateKey, msig, stx)
	return err
}

func encMtxs(txs string) error {
	var err error
	enctxs := strings.Split(txs, " ")
	bytetxs = make([][]byte, len(enctxs))
	for i := range enctxs {
		bytetxs[i], err = base64.StdEncoding.DecodeString(enctxs[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func mergeMsig() (err error) {
	_, stx, err = crypto.MergeMultisigTransactions(bytetxs...)
	return
}

func microToAlgos(ma int) error {
	microalgos = types.MicroAlgos(ma)
	microalgos = types.ToMicroAlgos(microalgos.ToAlgos())
	return nil
}

func checkAlgos(ma int) error {
	if types.MicroAlgos(ma) != microalgos {
		return fmt.Errorf("Converting to and from algos should not change the value")
	}
	return nil
}

func newAccInfo() error {
	_, err := aclv2.AccountInformation(pk).Do(context.Background())
	_, _ = kcl.DeleteKey(handle, walletPswd, pk)
	return err
}

func createKeyregWithStateProof(keyregType string) (err error) {
	params, err := aclv2.SuggestedParams().Do(context.Background())
	if err != nil {
		return err
	}
	lastRound = uint64(params.LastRoundValid)
	pk = accounts[0]
	if keyregType == "online" {
		nonpart = false
		votekey = "9mr13Ri8rFepxN3ghIUrZNui6LqqM5hEzB45Rri5lkU="
		selkey = "dx717L3uOIIb/jr9OIyls1l5Ei00NFgRa380w7TnPr4="
		votefst = uint64(0)
		votelst = uint64(30001)
		votekd = uint64(10000)
		stateProofPK = "mYR0GVEObMTSNdsKM6RwYywHYPqVDqg3E4JFzxZOreH9NU8B+tKzUanyY8AQ144hETgSMX7fXWwjBdHz6AWk9w=="
	} else if keyregType == "nonparticipation" {
		nonpart = true
		votekey = ""
		selkey = ""
		votefst = 0
		votelst = 0
		votekd = 0
		stateProofPK = ""
	} else if keyregType == "offline" {
		nonpart = false
		votekey = ""
		selkey = ""
		votefst = 0
		votelst = 0
		votekd = 0
		stateProofPK = ""
	}

	txn, err = transaction.MakeKeyRegTxnWithStateProofKey(accounts[0], note, params, votekey, selkey, stateProofPK, votefst, votelst, votekd, nonpart)
	if err != nil {
		return err
	}

	return err
}

func createAssetTestFixture() error {
	assetTestFixture.Creator = ""
	assetTestFixture.AssetIndex = 1
	assetTestFixture.AssetName = "testcoin"
	assetTestFixture.AssetUnitName = "coins"
	assetTestFixture.AssetURL = "http://test"
	assetTestFixture.AssetMetadataHash = "fACPO4nRgO55j1ndAK3W6Sgc4APkcyFh"
	assetTestFixture.ExpectedParams = modelsV2.AssetParams{}
	assetTestFixture.QueriedParams = modelsV2.AssetParams{}
	assetTestFixture.LastTransactionIssued = types.Transaction{}
	return nil
}

func convertTransactionAssetParamsToModelsAssetParam(input types.AssetParams) modelsV2.AssetParams {
	result := modelsV2.AssetParams{
		Total:         input.Total,
		Decimals:      uint64(input.Decimals),
		DefaultFrozen: input.DefaultFrozen,
		Manager:       input.Manager.String(),
		Reserve:       input.Reserve.String(),
		Freeze:        input.Freeze.String(),
		Clawback:      input.Clawback.String(),
		UnitName:      input.UnitName,
		Name:          input.AssetName,
		Url:           input.URL,
		MetadataHash:  input.MetadataHash[:],
	}
	// input doesn't have Creator so that will remain empty
	return result
}

func assetCreateTxnHelper(issuance int, frozenState bool) error {
	accountToUse := accounts[0]
	assetTestFixture.Creator = accountToUse
	creator := assetTestFixture.Creator
	params, err := aclv2.SuggestedParams().Do(context.Background())
	if err != nil {
		return err
	}
	lastRound = uint64(params.FirstRoundValid)
	assetNote := []byte(nil)
	assetIssuance := uint64(issuance)
	manager := creator
	reserve := creator
	freeze := creator
	clawback := creator
	unitName := assetTestFixture.AssetUnitName
	assetName := assetTestFixture.AssetName
	url := assetTestFixture.AssetURL
	metadataHash := assetTestFixture.AssetMetadataHash
	assetCreateTxn, err := transaction.MakeAssetCreateTxn(creator, assetNote, params, assetIssuance, 0, frozenState, manager, reserve, freeze, clawback, unitName, assetName, url, metadataHash)
	assetTestFixture.LastTransactionIssued = assetCreateTxn
	txn = assetCreateTxn
	assetTestFixture.ExpectedParams = convertTransactionAssetParamsToModelsAssetParam(assetCreateTxn.AssetParams)
	//convertTransactionAssetParamsToModelsAssetParam leaves creator blank, repopulate
	assetTestFixture.ExpectedParams.Creator = creator
	return err
}

func defaultAssetCreateTxn(issuance int) error {
	return assetCreateTxnHelper(issuance, false)
}

func defaultAssetCreateTxnWithDefaultFrozen(issuance int) error {
	return assetCreateTxnHelper(issuance, true)
}

func createNoManagerAssetReconfigure() error {
	creator := assetTestFixture.Creator
	params, err := aclv2.SuggestedParams().Do(context.Background())
	if err != nil {
		return err
	}
	lastRound = uint64(params.FirstRoundValid)
	assetNote := []byte(nil)
	reserve := ""
	freeze := ""
	clawback := ""
	manager := creator // if this were "" as well, this wouldn't be a reconfigure txn, it would be a destroy txn
	assetReconfigureTxn, err := transaction.MakeAssetConfigTxn(creator, assetNote, params, assetTestFixture.AssetIndex, manager, reserve, freeze, clawback, false)
	assetTestFixture.LastTransactionIssued = assetReconfigureTxn
	txn = assetReconfigureTxn
	// update expected params
	assetTestFixture.ExpectedParams.Reserve = reserve
	assetTestFixture.ExpectedParams.Freeze = freeze
	assetTestFixture.ExpectedParams.Clawback = clawback
	return err
}

func createAssetDestroy() error {
	creator := assetTestFixture.Creator
	params, err := aclv2.SuggestedParams().Do(context.Background())
	if err != nil {
		return err
	}
	lastRound = uint64(params.FirstRoundValid)
	assetNote := []byte(nil)
	assetDestroyTxn, err := transaction.MakeAssetDestroyTxn(creator, assetNote, params, assetTestFixture.AssetIndex)
	assetTestFixture.LastTransactionIssued = assetDestroyTxn
	txn = assetDestroyTxn
	// update expected params
	assetTestFixture.ExpectedParams.Reserve = ""
	assetTestFixture.ExpectedParams.Freeze = ""
	assetTestFixture.ExpectedParams.Clawback = ""
	assetTestFixture.ExpectedParams.Manager = ""
	return err
}

// used in getAssetIndex and similar to get the index of the most recently operated on asset
func getMaxKey(numbers []modelsV2.Asset) uint64 {
	var maxNumber uint64 = 0
	for _, asset := range numbers {
		idx := asset.Index
		if idx > maxNumber {
			maxNumber = idx
		}
	}
	return maxNumber
}

func getAssetIndex() error {
	accountResp, err := aclv2.AccountInformation(assetTestFixture.Creator).Do(context.Background())
	if err != nil {
		return err
	}
	// get most recent asset index
	assetTestFixture.AssetIndex = getMaxKey(accountResp.CreatedAssets)
	return nil
}

func getAssetInfo() error {
	response, err := aclv2.GetAssetByID(assetTestFixture.AssetIndex).Do(context.Background())
	assetTestFixture.QueriedParams = response.Params
	return err
}

func failToGetAssetInfo() error {
	_, err := aclv2.GetAssetByID(assetTestFixture.AssetIndex).Do(context.Background())
	if err != nil {
		return nil
	}
	return fmt.Errorf("expected an error getting asset with index %v and creator %v, but no error was returned",
		assetTestFixture.AssetIndex, assetTestFixture.Creator)
}

func checkExpectedVsActualAssetParams() error {
	expectedParams := assetTestFixture.ExpectedParams
	actualParams := assetTestFixture.QueriedParams
	nameMatch := expectedParams.Name == actualParams.Name
	if !nameMatch {
		return fmt.Errorf("expected asset name was %v but actual asset name was %v",
			expectedParams.Name, actualParams.Name)
	}
	unitMatch := expectedParams.UnitName == actualParams.UnitName
	if !unitMatch {
		return fmt.Errorf("expected unit name was %v but actual unit name was %v",
			expectedParams.UnitName, actualParams.UnitName)
	}
	urlMatch := expectedParams.Url == actualParams.Url
	if !urlMatch {
		return fmt.Errorf("expected URL was %v but actual URL was %v",
			expectedParams.Url, actualParams.Url)
	}
	hashMatch := reflect.DeepEqual(expectedParams.MetadataHash, actualParams.MetadataHash)
	if !hashMatch {
		return fmt.Errorf("expected MetadataHash was %v but actual MetadataHash was %v",
			expectedParams.MetadataHash, actualParams.MetadataHash)
	}
	issuanceMatch := expectedParams.Total == actualParams.Total
	if !issuanceMatch {
		return fmt.Errorf("expected total issuance was %v but actual issuance was %v",
			expectedParams.Total, actualParams.Total)
	}
	defaultFrozenMatch := expectedParams.DefaultFrozen == actualParams.DefaultFrozen
	if !defaultFrozenMatch {
		return fmt.Errorf("expected default frozen state %v but actual default frozen state was %v",
			expectedParams.DefaultFrozen, actualParams.DefaultFrozen)
	}
	managerMatch := expectedParams.Manager == actualParams.Manager
	if !managerMatch {
		return fmt.Errorf("expected asset manager was %v but actual asset manager was %v",
			expectedParams.Manager, actualParams.Manager)
	}
	reserveMatch := expectedParams.Reserve == actualParams.Reserve
	if !reserveMatch {
		return fmt.Errorf("expected asset reserve was %v but actual asset reserve was %v",
			expectedParams.Reserve, actualParams.Reserve)
	}
	freezeMatch := expectedParams.Freeze == actualParams.Freeze
	if !freezeMatch {
		return fmt.Errorf("expected freeze manager was %v but actual freeze manager was %v",
			expectedParams.Freeze, actualParams.Freeze)
	}
	clawbackMatch := expectedParams.Clawback == actualParams.Clawback
	if !clawbackMatch {
		return fmt.Errorf("expected revocation (clawback) manager was %v but actual revocation manager was %v",
			expectedParams.Clawback, actualParams.Clawback)
	}
	return nil
}

func findAssetID(assets []modelsV2.AssetHolding, assetID uint64) (foundAsset modelsV2.AssetHolding, err error) {
	for _, asset := range assets {
		if asset.AssetId == assetID {
			return asset, nil
		}
	}
	return modelsV2.AssetHolding{}, fmt.Errorf("asset ID %d was not found", assetID)
}

func theCreatorShouldHaveAssetsRemaining(expectedBal int) error {
	expectedBalance := uint64(expectedBal)
	accountResp, err := aclv2.AccountInformation(assetTestFixture.Creator).Do(context.Background())
	if err != nil {
		return err
	}
	// Find asset ID
	holding, err := findAssetID(accountResp.Assets, assetTestFixture.AssetIndex)
	if err != nil {
		return fmt.Errorf("attempted to get balance of account %v for creator %v and index %v, but no balance was found for that index", assetTestFixture.Creator, assetTestFixture.Creator, assetTestFixture.AssetIndex)
	}
	if holding.Amount != expectedBalance {
		return fmt.Errorf("actual balance %v differed from expected balance %v", holding.Amount, expectedBalance)
	}
	return nil
}

func createAssetAcceptanceForSecondAccount() error {
	accountToUse := accounts[1]
	params, err := aclv2.SuggestedParams().Do(context.Background())
	if err != nil {
		return err
	}
	lastRound = uint64(params.FirstRoundValid)
	assetNote := []byte(nil)
	assetAcceptanceTxn, err := transaction.MakeAssetAcceptanceTxn(accountToUse, assetNote, params, assetTestFixture.AssetIndex)
	assetTestFixture.LastTransactionIssued = assetAcceptanceTxn
	txn = assetAcceptanceTxn
	return err
}

func createAssetTransferTransactionToSecondAccount(amount int) error {
	recipient := accounts[1]
	creator := assetTestFixture.Creator
	params, err := aclv2.SuggestedParams().Do(context.Background())
	if err != nil {
		return err
	}
	sendAmount := uint64(amount)
	closeAssetsTo := ""
	lastRound = uint64(params.FirstRoundValid)
	assetNote := []byte(nil)
	assetAcceptanceTxn, err := transaction.MakeAssetTransferTxn(creator, recipient, sendAmount, assetNote, params, closeAssetsTo, assetTestFixture.AssetIndex)
	assetTestFixture.LastTransactionIssued = assetAcceptanceTxn
	txn = assetAcceptanceTxn
	return err
}

func createAssetTransferTransactionFromSecondAccountToCreator(amount int) error {
	recipient := assetTestFixture.Creator
	sender := accounts[1]
	params, err := aclv2.SuggestedParams().Do(context.Background())
	if err != nil {
		return err
	}
	sendAmount := uint64(amount)
	closeAssetsTo := ""
	lastRound = uint64(params.FirstRoundValid)
	assetNote := []byte(nil)
	assetAcceptanceTxn, err := transaction.MakeAssetTransferTxn(sender, recipient, sendAmount, assetNote, params, closeAssetsTo, assetTestFixture.AssetIndex)
	assetTestFixture.LastTransactionIssued = assetAcceptanceTxn
	txn = assetAcceptanceTxn
	return err
}

// sets up a freeze transaction, with freeze state `setting` against target account `target`
// assumes creator is asset freeze manager
func freezeTransactionHelper(target string, setting bool) error {
	params, err := aclv2.SuggestedParams().Do(context.Background())
	if err != nil {
		return err
	}
	lastRound = uint64(params.FirstRoundValid)
	assetNote := []byte(nil)
	assetFreezeOrUnfreezeTxn, err := transaction.MakeAssetFreezeTxn(assetTestFixture.Creator, assetNote, params, assetTestFixture.AssetIndex, target, setting)
	assetTestFixture.LastTransactionIssued = assetFreezeOrUnfreezeTxn
	txn = assetFreezeOrUnfreezeTxn
	return err
}

func createFreezeTransactionTargetingSecondAccount() error {
	return freezeTransactionHelper(accounts[1], true)
}

func createUnfreezeTransactionTargetingSecondAccount() error {
	return freezeTransactionHelper(accounts[1], false)
}

func createRevocationTransaction(amount int) error {
	params, err := aclv2.SuggestedParams().Do(context.Background())
	if err != nil {
		return err
	}
	lastRound = uint64(params.FirstRoundValid)
	revocationAmount := uint64(amount)
	assetNote := []byte(nil)
	assetRevokeTxn, err := transaction.MakeAssetRevocationTxn(assetTestFixture.Creator, accounts[1], revocationAmount, assetTestFixture.Creator, assetNote, params, assetTestFixture.AssetIndex)
	assetTestFixture.LastTransactionIssued = assetRevokeTxn
	txn = assetRevokeTxn
	return err
}

// godog misreads the step for this function, so provide a handler for when it does so
func iCreateATransactionTransferringAmountAssetsFromCreatorToASecondAccount() error {
	return createAssetTransferTransactionToSecondAccount(500000)
}

func baseEncodedDataToSign(dataEnc string) (err error) {
	data, err = base64.StdEncoding.DecodeString(dataEnc)
	return
}

func programHash(addr string) (err error) {
	account.Address, err = types.DecodeAddress(addr)
	return
}

func iPerformTealsign() (err error) {
	sig, err = crypto.TealSign(account.PrivateKey, data, account.Address)
	return
}

func theSignatureShouldBeEqualTo(sigEnc string) error {
	expected, err := base64.StdEncoding.DecodeString(sigEnc)
	if err != nil {
		return err
	}
	if !bytes.Equal(expected, sig[:]) {
		return fmt.Errorf("%v != %v", expected, sig[:])
	}
	return nil
}

func baseEncodedProgram(programEnc string) error {
	program, err := base64.StdEncoding.DecodeString(programEnc)
	if err != nil {
		return err
	}
	account.Address = crypto.AddressFromProgram(program)
	return nil
}

func baseEncodedPrivateKey(skEnc string) error {
	seed, err := base64.StdEncoding.DecodeString(skEnc)
	if err != nil {
		return err
	}
	account.PrivateKey = ed25519.NewKeyFromSeed(seed)
	return nil
}

func tealCompile(filename string) (err error) {
	if len(filename) == 0 {
		return fmt.Errorf("empty teal program file name")
	}
	tealProgram, err := loadResource(filename)
	if err != nil {
		return err
	}
	result, err := aclv2.TealCompile(tealProgram).Do(context.Background())
	if err == nil {
		tealCompleResult.status = 200
		tealCompleResult.response = result
		return
	}
	if _, ok := err.(commonV2.BadRequest); ok {
		tealCompleResult.status = 400
		tealCompleResult.response.Hash = ""
		tealCompleResult.response.Result = ""
		return nil
	}

	return
}

func tealCheckCompile(status int, result string, hash string) error {
	if status != tealCompleResult.status {
		return fmt.Errorf("status: %d != %d", status, tealCompleResult.status)
	}
	if result != tealCompleResult.response.Result {
		return fmt.Errorf("result: %s != %s", result, tealCompleResult.response.Result)
	}

	if hash != tealCompleResult.response.Hash {
		return fmt.Errorf("hash: %s != %s", hash, tealCompleResult.response.Hash)
	}
	return nil
}

func tealCheckCompileAgainstFile(expectedFile string) error {
	if len(expectedFile) == 0 {
		return fmt.Errorf("empty teal program file name")
	}

	expectedTeal, err := loadResource(expectedFile)
	if err != nil {
		return err
	}

	actualTeal, err := base64.StdEncoding.DecodeString(tealCompleResult.response.Result)
	if err != nil {
		return err
	}

	if !bytes.Equal(actualTeal, expectedTeal) {
		return fmt.Errorf("Actual program does not match expected")
	}

	return nil
}

func tealDryrun(kind string, filename string) (err error) {
	if len(filename) == 0 {
		return fmt.Errorf("empty teal program file name")
	}
	tealProgram, err := loadResource(filename)
	if err != nil {
		return err
	}

	txns := []types.SignedTxn{{}}
	sources := []modelsV2.DryrunSource{}
	switch kind {
	case "compiled":
		txns[0].Lsig.Logic = tealProgram
	case "source":
		sources = append(sources, modelsV2.DryrunSource{
			FieldName: "lsig",
			Source:    string(tealProgram),
			TxnIndex:  0,
		})
	default:
		return fmt.Errorf("kind %s not in (source, compiled)", kind)
	}

	ddr := modelsV2.DryrunRequest{
		Txns:    txns,
		Sources: sources,
	}

	result, err := aclv2.TealDryrun(ddr).Do(context.Background())
	if err != nil {
		return
	}

	tealDryrunResult.response = result
	return
}

func tealCheckDryrun(result string) error {
	txnResult := tealDryrunResult.response.Txns[0]
	var msgs []string
	if txnResult.AppCallMessages != nil && len(txnResult.AppCallMessages) > 0 {
		msgs = txnResult.AppCallMessages
	} else if txnResult.LogicSigMessages != nil && len(txnResult.LogicSigMessages) > 0 {
		msgs = txnResult.LogicSigMessages
	}
	if len(msgs) == 0 {
		return fmt.Errorf("received no messages")
	}

	if msgs[len(msgs)-1] != result {
		return fmt.Errorf("dryrun status %s != %s", result, msgs[len(msgs)-1])
	}
	return nil
}

func createMethodObjectFromSignature(methodSig string) error {
	abiMethodLocal, err := abi.MethodFromSignature(methodSig)
	abiMethod = abiMethodLocal
	return err
}

func serializeMethodObjectIntoJson() error {
	abiMethodJson, err := json.Marshal(abiMethod)
	if err != nil {
		return err
	}

	abiJsonString = string(abiMethodJson)
	return nil
}

func checkSerializedMethodObject(jsonFile, loadedFrom string) error {
	directory := path.Join("./features/unit/", loadedFrom)
	jsons, err := loadMockJsons(jsonFile, directory)
	if err != nil {
		return err
	}
	correctJson := string(jsons[0])

	var actualJson interface{}
	err = json.Unmarshal([]byte(abiJsonString), &actualJson)
	if err != nil {
		return err
	}

	var expectedJson interface{}
	err = json.Unmarshal([]byte(correctJson), &expectedJson)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(actualJson, expectedJson) {
		return fmt.Errorf("json strings %s != %s", correctJson, abiJsonString)
	}

	return nil
}

func createMethodObjectFromProperties(name, firstArgType, secondArgType, returnType string) error {
	args := []abi.Arg{
		{Name: "", Type: firstArgType, Desc: ""},
		{Name: "", Type: secondArgType, Desc: ""},
	}
	abiMethod = abi.Method{
		Name:    name,
		Desc:    "",
		Args:    args,
		Returns: abi.Return{Type: returnType, Desc: ""},
	}
	return nil
}

func createMethodObjectWithArgNames(name, firstArgName, firstArgType, secondArgName, secondArgType, returnType string) error {
	args := []abi.Arg{
		{Name: firstArgName, Type: firstArgType, Desc: ""},
		{Name: secondArgName, Type: secondArgType, Desc: ""},
	}
	abiMethod = abi.Method{
		Name:    name,
		Desc:    "",
		Args:    args,
		Returns: abi.Return{Type: returnType, Desc: ""},
	}
	return nil
}

func createMethodObjectWithDescription(name, nameDesc, firstArgType, firstDesc, secondArgType, secondDesc, returnType string) error {
	args := []abi.Arg{
		{Name: "", Type: firstArgType, Desc: firstDesc},
		{Name: "", Type: secondArgType, Desc: secondDesc},
	}
	abiMethod = abi.Method{
		Name:    name,
		Desc:    nameDesc,
		Args:    args,
		Returns: abi.Return{Type: returnType, Desc: ""},
	}
	return nil
}

func checkTxnCount(givenTxnCount int) error {
	correctTxnCount := abiMethod.GetTxCount()
	if correctTxnCount != givenTxnCount {
		return fmt.Errorf("txn count %d != %d", givenTxnCount, correctTxnCount)
	}
	return nil
}

func checkMethodSelector(givenMethodSelector string) error {
	correctMethodSelector := hex.EncodeToString(abiMethod.GetSelector())
	if correctMethodSelector != givenMethodSelector {
		return fmt.Errorf("method selector %s != %s", givenMethodSelector, correctMethodSelector)
	}
	return nil
}

func createInterfaceObject(name string, desc string) error {
	abiInterface = abi.Interface{
		Name:    name,
		Desc:    desc,
		Methods: []abi.Method{abiMethod},
	}
	return nil
}

func serializeInterfaceObjectIntoJson() error {
	abiInterfaceJson, err := json.Marshal(abiInterface)
	if err != nil {
		return err
	}

	abiJsonString = string(abiInterfaceJson)
	return nil
}

func createContractObject(name string, desc string) error {
	abiContract = abi.Contract{
		Name:     name,
		Desc:     desc,
		Networks: make(map[string]abi.ContractNetworkInfo),
		Methods:  []abi.Method{abiMethod},
	}
	return nil
}

func iSetTheContractsAppIDToForTheNetwork(appID int, network string) error {
	if appID < 0 {
		return fmt.Errorf("App ID must not be negative. Got: %d", appID)
	}
	abiContract.Networks[network] = abi.ContractNetworkInfo{AppID: uint64(appID)}
	return nil
}

func serializeContractObjectIntoJson() error {
	abiContractJson, err := json.Marshal(abiContract)
	if err != nil {
		return err
	}

	abiJsonString = string(abiContractJson)
	return nil
}

func iAppendToMyMethodObjectsListInTheCaseOfANonemptySignature(arg1 string) error {
	if arg1 == "" {
		return nil
	}

	meth, err := abi.MethodFromSignature(arg1)
	abiMethods = append(abiMethods, meth)
	return err
}

func iCreateAnInterfaceObjectFromMyMethodObjectsList() error {
	abiInterface = abi.Interface{
		Name:    "",
		Methods: abiMethods,
	}
	return nil
}

func iGetTheMethodFromTheInterfaceByName(arg1 string) error {
	abiMethod, globalErrForExamination = abiInterface.GetMethodByName(arg1)
	return nil
}

func iCreateAContractObjectFromMyMethodObjectsList() error {
	abiContract = abi.Contract{
		Name:    "",
		Methods: abiMethods,
	}
	return nil
}

func iGetTheMethodFromTheContractByName(arg1 string) error {
	abiMethod, globalErrForExamination = abiContract.GetMethodByName(arg1)
	return nil
}

func theProducedMethodSignatureShouldEqualIfThereIsAnErrorItBeginsWith(arg1, arg2 string) error {
	if abiMethod.Name != "" {
		if arg2 != "" {
			return fmt.Errorf("expected error condition but got a method")
		}
		if arg1 != abiMethod.GetSignature() {
			return fmt.Errorf("signature mismatch: %s != %s", arg1, abiMethod.GetSignature())
		}
	} else if globalErrForExamination != nil {
		if arg2 == "" {
			return fmt.Errorf("got error %s, expected no error", globalErrForExamination)
		}

		if !strings.Contains(globalErrForExamination.Error(), arg2) {
			return fmt.Errorf("produced error does not match expected: %q does not contain %q", globalErrForExamination.Error(), arg2)
		}

	} else {
		return fmt.Errorf("both abi method and error string are empty")
	}

	return nil
}

// equality helper methods
func checkEqualMethods(method1, method2 abi.Method) bool {
	if method1.Name != method2.Name || method1.Desc != method2.Desc {
		return false
	}

	if method1.Returns.Type != method2.Returns.Type || method1.Returns.Desc != method2.Returns.Desc {
		return false
	}

	if len(method1.Args) != len(method2.Args) {
		return false
	}

	for i, arg1 := range method1.Args {
		arg2 := method2.Args[i]
		if arg1.Name != arg2.Name || arg1.Type != arg2.Type || arg1.Desc != arg2.Desc {
			return false
		}
	}
	return true
}

func checkEqualInterfaces(interface1, interface2 abi.Interface) bool {
	if interface1.Name != interface2.Name || interface1.Desc != interface2.Desc {
		return false
	}

	if len(interface1.Methods) != len(interface2.Methods) {
		return false
	}

	for i, method := range interface1.Methods {
		if !checkEqualMethods(method, interface2.Methods[i]) {
			return false
		}
	}
	return true
}

func checkEqualContracts(contract1, contract2 abi.Contract) bool {
	if contract1.Name != contract2.Name || contract1.Desc != contract2.Desc {
		return false
	}

	if len(contract1.Networks) != len(contract2.Networks) {
		return false
	}

	for network, info1 := range contract1.Networks {
		info2, ok := contract2.Networks[network]
		if !ok || info1 != info2 {
			return false
		}
	}

	if len(contract1.Methods) != len(contract2.Methods) {
		return false
	}

	for i, method := range contract1.Methods {
		if !checkEqualMethods(method, contract2.Methods[i]) {
			return false
		}
	}
	return true
}

func deserializeMethodJson() error {
	var deserializedMethod abi.Method
	err := json.Unmarshal([]byte(abiJsonString), &deserializedMethod)
	if err != nil {
		return err
	}

	if !checkEqualMethods(deserializedMethod, abiMethod) {
		return fmt.Errorf("Deserialized method does not match original method")
	}
	return nil
}

func deserializeInterfaceJson() error {
	var deserializedInterface abi.Interface
	err := json.Unmarshal([]byte(abiJsonString), &deserializedInterface)
	if err != nil {
		return err
	}
	if !checkEqualInterfaces(deserializedInterface, abiInterface) {
		return fmt.Errorf("Deserialized interface does not match original interface")
	}
	return nil
}

func deserializeContractJson() error {
	var deserializedContract abi.Contract
	err := json.Unmarshal([]byte(abiJsonString), &deserializedContract)
	if err != nil {
		return err
	}
	if !checkEqualContracts(deserializedContract, abiContract) {
		return fmt.Errorf("Deserialized contract does not match original contract")
	}
	return nil
}

func aNewAtomicTransactionComposer() error {
	txComposer = transaction.AtomicTransactionComposer{}
	return nil
}

func suggestedTransactionParameters(fee int, flatFee string, firstValid, LastValid int, genesisHash, genesisId string) error {
	if flatFee != "true" && flatFee != "false" {
		return fmt.Errorf("flatFee must be either 'true' or 'false'")
	}

	genHash, err := base64.StdEncoding.DecodeString(genesisHash)
	if err != nil {
		return err
	}

	sugParams = types.SuggestedParams{
		Fee:             types.MicroAlgos(fee),
		GenesisID:       genesisId,
		GenesisHash:     genHash,
		FirstRoundValid: types.Round(firstValid),
		LastRoundValid:  types.Round(LastValid),
		FlatFee:         flatFee == "true",
	}

	return nil
}

func anApplicationId(id int) error {
	if id < 0 {
		return fmt.Errorf("app id must be positive integer")
	}

	applicationId = uint64(id)
	return nil
}

func iMakeATransactionSignerForTheAccount(accountType string) error {
	if accountType == "signing" {
		accountTxSigner = transaction.BasicAccountTransactionSigner{
			Account: account,
		}
	} else if accountType == "transient" {
		accountTxSigner = transaction.BasicAccountTransactionSigner{
			Account: transientAccount,
		}
	}

	return nil
}

func iCreateANewMethodArgumentsArray() error {
	methodArgs = make([]interface{}, 0)
	return nil
}

func iAppendTheEncodedArgumentsToTheMethodArgumentsArray(commaSeparatedB64Args string) error {
	if len(commaSeparatedB64Args) == 0 {
		return nil
	}

	b64Args := strings.Split(commaSeparatedB64Args, ",")
	for _, b64Arg := range b64Args {
		if strings.Contains(b64Arg, ":") {
			// special case for inserting existing application ID
			parts := strings.Split(b64Arg, ":")
			if len(parts) != 2 || parts[0] != "ctxAppIdx" {
				return fmt.Errorf("Cannot process argument: %s", b64Arg)
			}
			parsedIndex, err := strconv.ParseUint(parts[1], 10, 64)
			if err != nil {
				return err
			}
			if parsedIndex >= uint64(len(applicationIds)) {
				return fmt.Errorf("Application index out of bounds: %d, number of app IDs is %d", parsedIndex, len(applicationIds))
			}
			abiUint64, err := abi.TypeOf("uint64")
			if err != nil {
				return err
			}
			encodedUint64, err := abiUint64.Encode(applicationIds[parsedIndex])
			if err != nil {
				return err
			}
			methodArgs = append(methodArgs, encodedUint64)
			continue
		}
		decodedArg, err := base64.StdEncoding.DecodeString(b64Arg)
		if err != nil {
			return err
		}
		methodArgs = append(methodArgs, decodedArg)
	}

	return nil
}

func addMethodCall(accountType, strOnComplete string) error {
	return addMethodCallHelper(accountType, strOnComplete, "", "", 0, 0, 0, 0, 0, false)
}

func addMethodCallForUpdate(accountType, strOnComplete, approvalProgram, clearProgram string) error {
	return addMethodCallHelper(accountType, strOnComplete, approvalProgram, clearProgram, 0, 0, 0, 0, 0, false)
}

func addMethodCallForCreate(accountType, strOnComplete, approvalProgram, clearProgram string, globalBytes, globalInts, localBytes, localInts, extraPages int) error {
	return addMethodCallHelper(accountType, strOnComplete, approvalProgram, clearProgram, globalBytes, globalInts, localBytes, localInts, extraPages, false)
}

func addMethodCallWithNonce(accountType, strOnComplete string) error {
	return addMethodCallHelper(accountType, strOnComplete, "", "", 0, 0, 0, 0, 0, true)
}

func addMethodCallHelper(accountType, strOnComplete, approvalProgram, clearProgram string, globalBytes, globalInts, localBytes, localInts, extraPages int, useNonce bool) error {
	var onComplete types.OnCompletion
	switch strOnComplete {
	case "create":
		onComplete = types.NoOpOC
	case "noop":
		onComplete = types.NoOpOC
	case "update":
		onComplete = types.UpdateApplicationOC
	case "call":
		onComplete = types.NoOpOC
	case "optin":
		onComplete = types.OptInOC
	case "clear":
		onComplete = types.ClearStateOC
	case "closeout":
		onComplete = types.CloseOutOC
	case "delete":
		onComplete = types.DeleteApplicationOC
	default:
		return fmt.Errorf("invalid onComplete value")
	}

	var useAccount crypto.Account
	if accountType == "signing" {
		useAccount = account
	} else if accountType == "transient" {
		useAccount = transientAccount
	}

	var approvalProgramBytes []byte
	var clearProgramBytes []byte
	var err error

	if approvalProgram != "" {
		approvalProgramBytes, err = readTealProgram(approvalProgram)
		if err != nil {
			return err
		}
	}

	if clearProgram != "" {
		clearProgramBytes, err = readTealProgram(clearProgram)
		if err != nil {
			return err
		}
	}

	if globalInts < 0 || globalBytes < 0 || localInts < 0 || localBytes < 0 || extraPages < 0 {
		return fmt.Errorf("Values for globalInts, globalBytes, localInts, localBytes, and extraPages cannot be negative")
	}

	// populate args from methodArgs
	if len(methodArgs) != len(abiMethod.Args) {
		return fmt.Errorf("Provided argument count is incorrect. Expected %d, got %d", len(abiMethod.Args), len(methodArgs))
	}

	var preparedArgs []interface{}
	for i, argSpec := range abiMethod.Args {
		if argSpec.IsTransactionArg() {
			// encodedArg is already a TransactionWithSigner
			preparedArgs = append(preparedArgs, methodArgs[i])
			continue
		}

		encodedArg, ok := methodArgs[i].([]byte)
		if !ok {
			return fmt.Errorf("Argument should be a byte slice")
		}

		var typeToDecode abi.Type
		var err error

		if argSpec.IsReferenceArg() {
			switch argSpec.Type {
			case abi.AccountReferenceType:
				typeToDecode, err = abi.TypeOf("address")
			case abi.ApplicationReferenceType, abi.AssetReferenceType:
				typeToDecode, err = abi.TypeOf("uint64")
			default:
				return fmt.Errorf("Unknown reference type: %s", argSpec.Type)
			}
		} else {
			typeToDecode, err = argSpec.GetTypeObject()
		}
		if err != nil {
			return err
		}

		decodedArg, err := typeToDecode.Decode(encodedArg)
		if err != nil {
			return err
		}

		preparedArgs = append(preparedArgs, decodedArg)
	}

	methodCallParams := transaction.AddMethodCallParams{
		AppID:           applicationId,
		Method:          abiMethod,
		MethodArgs:      preparedArgs,
		Sender:          useAccount.Address,
		SuggestedParams: sugParams,
		OnComplete:      onComplete,
		ApprovalProgram: approvalProgramBytes,
		ClearProgram:    clearProgramBytes,
		GlobalSchema: types.StateSchema{
			NumUint:      uint64(globalInts),
			NumByteSlice: uint64(globalBytes),
		},
		LocalSchema: types.StateSchema{
			NumUint:      uint64(localInts),
			NumByteSlice: uint64(localBytes),
		},
		ExtraPages: uint32(extraPages),
		Signer:     accountTxSigner,
	}

	if useNonce {
		methodCallParams.Note = note
	}

	return txComposer.AddMethodCall(methodCallParams)
}

func iAddTheNonce(nonce string) error {
	note = []byte("I should be unique thanks to this nonce: " + nonce)
	return nil
}

func buildTheTransactionGroupWithTheComposer(errorType string) error {
	_, err := txComposer.BuildGroup()

	switch errorType {
	case "":
		// no error expected
		return err
	case "zero group size error":
		if err == nil || err.Error() != "attempting to build group with zero transactions" {
			return fmt.Errorf("Expected error, but got: %v", err)
		}
		return nil
	default:
		return fmt.Errorf("Unknown error type: %s", errorType)
	}
}

func theComposerShouldHaveAStatusOf(strStatus string) error {
	var status transaction.AtomicTransactionComposerStatus
	switch strStatus {
	case "BUILDING":
		status = transaction.BUILDING
	case "BUILT":
		status = transaction.BUILT
	case "SIGNED":
		status = transaction.SIGNED
	case "SUBMITTED":
		status = transaction.SUBMITTED
	case "COMMITTED":
		status = transaction.COMMITTED
	default:
		return fmt.Errorf("invalid status provided")
	}

	if status != txComposer.GetStatus() {
		return fmt.Errorf("status does not match")
	}

	return nil
}

func iGatherSignaturesWithTheComposer() error {
	signedTxs, err := txComposer.GatherSignatures()
	sigTxs = signedTxs
	return err
}

func theBaseEncodedSignedTransactionsShouldEqual(encodedTxsStr string) error {
	encodedTxs := strings.Split(encodedTxsStr, ",")
	if len(encodedTxs) != len(sigTxs) {
		return fmt.Errorf("Actual and expected number of signed transactions don't match")
	}

	for i, encodedTx := range encodedTxs {
		gold, err := base64.StdEncoding.DecodeString(encodedTx)
		if err != nil {
			return err
		}
		stxStr := base64.StdEncoding.EncodeToString(sigTxs[i])
		if !bytes.Equal(gold, sigTxs[i]) {
			return fmt.Errorf("Application signed transaction does not match the golden: %s != %s", stxStr, encodedTx)
		}
	}

	return nil
}

func iBuildAPaymentTransactionWithSenderReceiverAmountCloseRemainderTo(sender, receiver string, amount int, closeTo string) error {
	if amount < 0 {
		return fmt.Errorf("amount must be a positive integer")
	}

	if sender == "transient" {
		sender = transientAccount.Address.String()
	}

	if receiver == "transient" {
		receiver = transientAccount.Address.String()
	}

	var err error
	txn, err = transaction.MakePaymentTxn(sender, receiver, uint64(amount), nil, closeTo, sugParams)
	tx = txn
	return err
}

func iCreateATransactionWithSignerWithTheCurrentTransaction() error {
	accountTxAndSigner = transaction.TransactionWithSigner{
		Signer: accountTxSigner,
		Txn:    txn,
	}
	return nil
}

func iAppendTheCurrentTransactionWithSignerToTheMethodArgumentsArray() error {
	methodArgs = append(methodArgs, accountTxAndSigner)
	return nil
}

func theDecodedTransactionShouldEqualTheOriginal() error {
	var decodedTx types.SignedTxn
	err := msgpack.Decode(stx, &decodedTx)
	if err != nil {
		return err
	}

	// This test isn't perfect as it's sensitive to non-meaningful changes (e.g. nil slice vs 0
	// length slice), but it's good enough for now. We may want a Transaction.Equals method in the
	// future.
	if !reflect.DeepEqual(tx, decodedTx.Txn) {
		return fmt.Errorf("Transactions unequal: %#v != %#v", tx, decodedTx.Txn)
	}

	return nil
}

func aDryrunResponseFileAndATransactionAtIndex(arg1, arg2 string) error {
	data, err := loadResource(arg1)
	if err != nil {
		return err
	}
	dr, err := transaction.NewDryrunResponseFromJSON(data)
	if err != nil {
		return err
	}
	idx, err := strconv.Atoi(arg2)
	if err != nil {
		return err
	}
	txTrace = dr.Txns[idx]
	return nil
}

func callingAppTraceProduces(arg1 string) error {
	cfg := transaction.DefaultStackPrinterConfig()
	cfg.TopOfStackFirst = false
	trace = txTrace.GetAppCallTrace(cfg)

	data, err := loadResource(arg1)
	if err != nil {
		return err
	}
	if string(data) != trace {
		return fmt.Errorf("No matching trace: \n'%s'\nvs\n'%s'\n", string(data), trace)
	}
	return nil
}

func aSourceMapJsonFile(srcMapJsonPath string) error {
	b, err := loadResource(srcMapJsonPath)
	if err != nil {
		return err
	}

	ism := map[string]interface{}{}
	if err := json.Unmarshal(b, &ism); err != nil {
		return err
	}

	sourceMap, err = logic.DecodeSourceMap(ism)

	return err
}

func theStringComposedOfPclineNumberEquals(expectedPcToLineString string) error {
	var buff []string
	for pc := 0; pc < len(sourceMap.PcToLine); pc++ {
		line := sourceMap.PcToLine[pc]
		buff = append(buff, fmt.Sprintf("%d:%d", pc, line))
	}
	actualStr := strings.Join(buff, ";")
	if expectedPcToLineString != actualStr {
		return fmt.Errorf("Expected %s got %s", expectedPcToLineString, actualStr)
	}
	return nil
}

func gettingTheLineAssociatedWithAPcEquals(strPc, strLine string) error {
	pc, _ := strconv.Atoi(strPc)
	expectedLine, _ := strconv.Atoi(strLine)

	actualLine, ok := sourceMap.GetLineForPc(pc)
	if !ok {
		return fmt.Errorf("expected valid line, got !ok")
	}

	if actualLine != expectedLine {
		return fmt.Errorf("expected %d got %d", expectedLine, actualLine)
	}

	return nil
}

func gettingTheLastPcAssociatedWithALineEquals(strLine, strPc string) error {
	expectedPc, _ := strconv.Atoi(strPc)
	line, _ := strconv.Atoi(strLine)

	pcs := sourceMap.GetPcsForLine(line)
	actualPc := pcs[len(pcs)-1]

	if actualPc != expectedPc {
		return fmt.Errorf("expected %d got %d", expectedPc, actualPc)
	}

	return nil
}

func iCompileATealProgramWithMappingEnabled(programPath string) error {
	fileContents, err := loadResource(programPath)
	if err != nil {
		return err
	}

	result, err := aclv2.TealCompile(fileContents).Sourcemap(true).Do(context.Background())
	if err != nil {
		return err
	}

	if result.Sourcemap == nil {
		return fmt.Errorf("No source map returned")
	}

	srcMapping = *result.Sourcemap
	return nil
}

func theResultingSourceMapIsTheSameAsTheJson(expectedJsonPath string) error {

	expectedJson, err := loadResource(expectedJsonPath)
	if err != nil {
		return err
	}

	// Marshal the map to json then unmarshal it so we get alphabetic ordering
	expectedMap := map[string]interface{}{}
	err = json.Unmarshal(expectedJson, &expectedMap)
	if err != nil {
		return err
	}

	expectedJson, err = json.Marshal(expectedMap)
	if err != nil {
		return err
	}

	// Turn it back into a string
	actualJson, err := json.Marshal(srcMapping)
	if err != nil {
		return nil
	}

	if !bytes.Equal(expectedJson, actualJson) {
		return fmt.Errorf("expected %s got %s", expectedJson, actualJson)
	}

	return nil
}

func takeB64encodedBytes(b64encodedBytes string) error {
	var err error
	seeminglyProgram, err = base64.StdEncoding.DecodeString(b64encodedBytes)
	if err != nil {
		return err
	}
	return nil
}

func heuristicCheckOverBytes() error {
	_, sanityCheckError = crypto.MakeLogicSigAccountEscrowChecked(seeminglyProgram, nil)
	return nil
}

func checkErrorIfMatching(errMsg string) error {
	if len(errMsg) == 0 {
		if sanityCheckError != nil {
			return fmt.Errorf("expected err message to be empty, but sanity check says %w", sanityCheckError)
		}
	} else {
		if sanityCheckError == nil || !strings.Contains(sanityCheckError.Error(), errMsg) {
			return fmt.Errorf("expected err to contain %s, but sanity check error not matching: %w", errMsg, sanityCheckError)
		}
	}
	return nil
}

func disassemblyMatches(bytecodeFilename, sourceFilename string) error {
	disassembledBytes, err := loadResource(bytecodeFilename)
	if err != nil {
		return err
	}
	actualResult, err := aclv2.TealDisassemble(disassembledBytes).Do(context.Background())
	if err != nil {
		return err
	}
	expectedBytes, err := loadResource(sourceFilename)
	if err != nil {
		return err
	}
	expectedResult := string(expectedBytes)
	if actualResult.Result != expectedResult {
		return fmt.Errorf("Actual program does not match expected: %s != %s", actualResult.Result, expectedResult)
	}
	return nil
}
