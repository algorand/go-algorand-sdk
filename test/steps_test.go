package test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base32"
	"encoding/base64"
	"encoding/gob"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"path/filepath"

	"golang.org/x/crypto/ed25519"

	"github.com/algorand/go-algorand-sdk/auction"
	"github.com/algorand/go-algorand-sdk/client/algod"
	"github.com/algorand/go-algorand-sdk/client/algod/models"
	"github.com/algorand/go-algorand-sdk/client/kmd"
	algodV2 "github.com/algorand/go-algorand-sdk/client/v2/algod"
	commonV2 "github.com/algorand/go-algorand-sdk/client/v2/common"
	modelsV2 "github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/algorand/go-algorand-sdk/templates"
	"github.com/algorand/go-algorand-sdk/types"
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
var acl algod.Client
var aclv2 *algodV2.Client
var walletName string
var walletPswd string
var walletID string
var handle string
var versions []string
var status models.NodeStatus
var statusAfter models.NodeStatus
var msigExp kmd.ExportMultisigResponse
var pk string
var accounts []string
var e bool
var lastRound uint64
var sugParams types.SuggestedParams
var sugFee models.TransactionFee
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
var votefst uint64
var votelst uint64
var votekd uint64
var num string
var backupTxnSender string
var groupTxnBytes []byte
var data []byte
var sig types.Signature

var assetTestFixture struct {
	Creator               string
	AssetIndex            uint64
	AssetName             string
	AssetUnitName         string
	AssetURL              string
	AssetMetadataHash     string
	ExpectedParams        models.AssetParams
	QueriedParams         models.AssetParams
	LastTransactionIssued types.Transaction
}

var contractTestFixture struct {
	activeAddress      string
	contractFundAmount uint64
	split              templates.Split
	splitN             uint64
	splitD             uint64
	splitMin           uint64
	htlc               templates.HTLC
	htlcPreImage       string
	periodicPay        templates.PeriodicPayment
	periodicPayPeriod  uint64
	limitOrder         templates.LimitOrder
	limitOrderN        uint64
	limitOrderD        uint64
	limitOrderMin      uint64
	dynamicFee         templates.DynamicFee
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
		IndexerIntegrationTestContext(s)
		ApplicationsContext(s)
		ApplicationsUnitContext(s)
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
	s.Step("I create the payment transaction", createTxn)
	s.Step(`multisig addresses "([^"]*)"`, msigAddresses)
	s.Step("I create the multisig payment transaction", createMsigTxn)
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
	s.Step("I get versions with kmd", kclV)
	s.Step("I get the status", getStatus)
	s.Step(`^I get status after this block`, statusAfterBlock)
	s.Step("I can get the block info", block)
	s.Step("I import the multisig", importMsig)
	s.Step("the multisig should be in the wallet", msigInWallet)
	s.Step("I export the multisig", expMsig)
	s.Step("the multisig should equal the exported multisig", msigEq)
	s.Step("I delete the multisig", deleteMsig)
	s.Step("the multisig should not be in the wallet", msigNotInWallet)
	s.Step("I generate a key using kmd", genKeyKmd)
	s.Step("the key should be in the wallet", keyInWallet)
	s.Step("I delete the key", deleteKey)
	s.Step("the key should not be in the wallet", keyNotInWallet)
	s.Step("I generate a key", genKey)
	s.Step("I import the key", importKey)
	s.Step("the private key should be equal to the exported private key", skEqExport)
	s.Step("a kmd client", kmdClient)
	s.Step("an algod client", algodClient)
	s.Step("wallet information", walletInfo)
	s.Step(`default transaction with parameters (\d+) "([^"]*)"`, defaultTxn)
	s.Step(`default multisig transaction with parameters (\d+) "([^"]*)"`, defaultMsigTxn)
	s.Step("I get the private key", getSk)
	s.Step("I send the transaction", sendTxn)
	s.Step("I send the kmd-signed transaction", sendTxnKmd)
	s.Step("I send the bogus kmd-signed transaction", sendTxnKmdFailureExpected)
	s.Step("I send the multisig transaction", sendMsigTxn)
	s.Step("the transaction should go through", checkTxn)
	s.Step("the transaction should not go through", txnFail)
	s.Step("I sign the transaction with kmd", signKmd)
	s.Step("the signed transaction should equal the kmd signed transaction", signBothEqual)
	s.Step("I sign the multisig transaction with kmd", signMsigKmd)
	s.Step("the multisig transaction should equal the kmd signed multisig transaction", signMsigBothEqual)
	s.Step(`I read a transaction "([^"]*)" from file "([^"]*)"`, readTxn)
	s.Step("I write the transaction to file", writeTxn)
	s.Step("the transaction should still be the same", checkEnc)
	s.Step("I do my part", createSaveTxn)
	s.Step(`^the node should be healthy`, nodeHealth)
	s.Step(`^I get the ledger supply`, ledger)
	s.Step(`^I get transactions by address and round`, txnsByAddrRound)
	s.Step(`^I get pending transactions`, txnsPending)
	s.Step(`^I get the suggested params`, suggestedParams)
	s.Step(`^I get the suggested fee`, suggestedFee)
	s.Step(`^the fee in the suggested params should equal the suggested fee`, checkSuggested)
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
	s.Step(`I get account information`, accInfo)
	s.Step("I sign the bid", signBid)
	s.Step("I get transactions by address only", txnsByAddrOnly)
	s.Step("I get transactions by address and date", txnsByAddrDate)
	s.Step(`key registration transaction parameters (\d+) (\d+) (\d+) "([^"]*)" "([^"]*)" "([^"]*)" (\d+) (\d+) (\d+) "([^"]*)" "([^"]*)`, keyregTxnParams)
	s.Step("I create the key registration transaction", createKeyregTxn)
	s.Step(`^I get recent transactions, limited by (\d+) transactions$`, getTxnsByCount)
	s.Step(`^I can get account information`, newAccInfo)
	s.Step(`^I can get the transaction by ID$`, txnbyID)
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
	s.Step(`^a split contract with ratio (\d+) to (\d+) and minimum payment (\d+)$`, aSplitContractWithRatioToAndMinimumPayment)
	s.Step(`^I send the split transactions$`, iSendTheSplitTransactions)
	s.Step(`^an HTLC contract with hash preimage "([^"]*)"$`, anHTLCContractWithHashPreimage)
	s.Step(`^I fund the contract account$`, iFundTheContractAccount)
	s.Step(`^I claim the algos$`, iClaimTheAlgosHTLC)
	s.Step(`^a periodic payment contract with withdrawing window (\d+) and period (\d+)$`, aPeriodicPaymentContractWithWithdrawingWindowAndPeriod)
	s.Step(`^I claim the periodic payment$`, iClaimThePeriodicPayment)
	s.Step(`^a limit order contract with parameters (\d+) (\d+) (\d+)$`, aLimitOrderContractWithParameters)
	s.Step(`^I swap assets for algos$`, iSwapAssetsForAlgos)
	s.Step(`^a dynamic fee contract with amount (\d+)$`, aDynamicFeeContractWithAmount)
	s.Step(`^I send the dynamic fee transactions$`, iSendTheDynamicFeeTransaction)
	s.Step("contract test fixture", createContractTestFixture)
	s.Step(`^I create a transaction transferring <amount> assets from creator to a second account$`, iCreateATransactionTransferringAmountAssetsFromCreatorToASecondAccount) // provide handler for when godog misreads
	s.Step(`^base64 encoded data to sign "([^"]*)"$`, baseEncodedDataToSign)
	s.Step(`^program hash "([^"]*)"$`, programHash)
	s.Step(`^I perform tealsign$`, iPerformTealsign)
	s.Step(`^the signature should be equal to "([^"]*)"$`, theSignatureShouldBeEqualTo)
	s.Step(`^base64 encoded program "([^"]*)"$`, baseEncodedProgram)
	s.Step(`^base64 encoded private key "([^"]*)"$`, baseEncodedPrivateKey)
	s.Step("an algod v2 client", algodClientV2)
	s.Step(`^I compile a teal program "([^"]*)"$`, tealCompile)
	s.Step(`^it is compiled with (\d+) and "([^"]*)" and "([^"]*)"$`, tealCheckCompile)
	s.Step(`^I dryrun a "([^"]*)" program "([^"]*)"$`, tealDryrun)
	s.Step(`^I get execution result "([^"]*)"$`, tealCheckDryrun)

	s.BeforeScenario(func(interface{}) {
		stxObj = types.SignedTxn{}
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

func createTxn() error {
	var err error
	paramsToUse := types.SuggestedParams{
		Fee:             types.MicroAlgos(fee),
		GenesisID:       gen,
		GenesisHash:     gh,
		FirstRoundValid: types.Round(fv),
		LastRoundValid:  types.Round(lv),
		FlatFee:         false,
	}
	txn, err = future.MakePaymentTxn(a.String(), to, amt, note, close, paramsToUse)
	if err != nil {
		return err
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
	txn, err = future.MakePaymentTxn(msigaddr.String(), to, amt, note, close, paramsToUse)
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
	v, err := acl.Versions()
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

func kclV() error {
	v, err := kcl.Version()
	versions = v.Versions
	return err
}

func getStatus() error {
	var err error
	status, err = acl.Status()
	lastRound = status.LastRound
	return err
}

func statusAfterBlock() error {
	var err error
	statusAfter, err = acl.StatusAfterBlock(lastRound)
	if err != nil {
		return err
	}
	return nil
}

func block() error {
	_, err := acl.Block(status.LastRound)
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

func algodClient() error {
	algodToken := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	algodAddress := "http://localhost:" + "60000"
	var err error
	acl, err = algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		return err
	}
	_, err = acl.StatusAfterBlock(1)
	return err
}

func algodClientV2() error {
	algodToken := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	algodAddress := "http://localhost:" + "60000"
	var err error
	aclv2, err = algodV2.MakeClient(algodAddress, algodToken)
	if err != nil {
		return err
	}
	_, err = aclv2.StatusAfterBlock(1).Do(context.Background())
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

func defaultTxn(iamt int, inote string) error {
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
	params, err := acl.BuildSuggestedParams()
	if err != nil {
		return err
	}
	lastRound = uint64(params.FirstRoundValid)
	txn, err = future.MakePaymentTxn(accounts[0], accounts[1], amt, note, "", params)
	return err
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
	params, err := acl.BuildSuggestedParams()
	if err != nil {
		return err
	}
	lastRound = uint64(params.FirstRoundValid)
	addr, err := msig.Address()
	if err != nil {
		return err
	}
	txn, err = future.MakePaymentTxn(addr.String(), accounts[1], amt, note, "", params)
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
	tx, err := acl.SendRawTransaction(stx)
	if err != nil {
		return err
	}
	txid = tx.TxID
	return nil
}

func sendTxnKmd() error {
	tx, err := acl.SendRawTransaction(stxKmd)
	if err != nil {
		e = true
	}
	txid = tx.TxID
	return nil
}

func sendTxnKmdFailureExpected() error {
	tx, err := acl.SendRawTransaction(stxKmd)
	if err == nil {
		e = false
		return fmt.Errorf("expected an error when sending kmd-signed transaction but no error occurred")
	}
	e = true
	txid = tx.TxID
	return nil
}

func sendMsigTxn() error {
	_, err := acl.SendRawTransaction(stx)

	if err != nil {
		e = true
	}

	return nil
}

func checkTxn() error {
	_, err := acl.PendingTransactionInformation(txid)
	if err != nil {
		return err
	}
	_, err = acl.StatusAfterBlock(lastRound + 2)
	if err != nil {
		return err
	}
	if txn.Sender.String() != "" && txn.Sender.String() != "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAY5HFKQ" {
		_, err = acl.TransactionInformation(txn.Sender.String(), txid)
	} else {
		_, err = acl.TransactionInformation(backupTxnSender, txid)
	}
	if err != nil {
		return err
	}
	_, err = acl.TransactionByID(txid)
	return err
}

func txnbyID() error {
	var err error
	_, err = acl.StatusAfterBlock(lastRound + 2)
	if err != nil {
		return err
	}
	_, err = acl.TransactionByID(txid)
	return err
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

func readTxn(encodedTxn string, inum string) error {
	encodedBytes, err := base64.StdEncoding.DecodeString(encodedTxn)
	if err != nil {
		return err
	}
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	num = inum
	path = filepath.Dir(filepath.Dir(path)) + "/temp/old" + num + ".tx"
	err = ioutil.WriteFile(path, encodedBytes, 0644)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = msgpack.Decode(data, &stxObj)
	return err
}

func writeTxn() error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	path = filepath.Dir(filepath.Dir(path)) + "/temp/raw" + num + ".tx"
	data := msgpack.Encode(stxObj)
	err = ioutil.WriteFile(path, data, 0644)
	return err
}

func checkEnc() error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	pathold := filepath.Dir(filepath.Dir(path)) + "/temp/old" + num + ".tx"
	dataold, err := ioutil.ReadFile(pathold)

	pathnew := filepath.Dir(filepath.Dir(path)) + "/temp/raw" + num + ".tx"
	datanew, err := ioutil.ReadFile(pathnew)

	if bytes.Equal(dataold, datanew) {
		return nil
	}
	return fmt.Errorf("should be equal")
}

func createSaveTxn() error {
	var err error

	amt = 100000
	pk = accounts[0]
	params, err := acl.BuildSuggestedParams()
	if err != nil {
		return err
	}
	lastRound = uint64(params.FirstRoundValid)
	txn, err = future.MakePaymentTxn(accounts[0], accounts[1], amt, note, "", params)
	if err != nil {
		return err
	}

	path, err := os.Getwd()
	if err != nil {
		return err
	}
	path = filepath.Dir(filepath.Dir(path)) + "/temp/txn.tx"
	data := msgpack.Encode(txn)
	err = ioutil.WriteFile(path, data, 0644)
	return err
}

func nodeHealth() error {
	err := acl.HealthCheck()
	return err
}

func ledger() error {
	_, err := acl.LedgerSupply()
	return err
}

func txnsByAddrRound() error {
	lr, err := acl.Status()
	if err != nil {
		return err
	}
	_, err = acl.TransactionsByAddr(accounts[0], 1, lr.LastRound)
	return err
}

func txnsByAddrOnly() error {
	_, err := acl.TransactionsByAddrLimit(accounts[0], 10)
	return err
}

func txnsByAddrDate() error {
	fromDate := time.Now().Format("2006-01-02")
	_, err := acl.TransactionsByAddrForDate(accounts[0], fromDate, fromDate)
	return err
}

func txnsPending() error {
	_, err := acl.GetPendingTransactions(10)
	return err
}

func suggestedParams() error {
	var err error
	sugParams, err = acl.BuildSuggestedParams()
	return err
}

func suggestedFee() error {
	var err error
	sugFee, err = acl.SuggestedFee()
	return err
}

func checkSuggested() error {
	if uint64(sugParams.Fee) != sugFee.Fee {
		return fmt.Errorf("suggested fee from params should be equal to suggested fee")
	}
	return nil
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
	txn, err = future.MakePaymentTxn(a.String(), to, amt, note, close, paramsToUse)
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

func accInfo() error {
	_, err := acl.AccountInformation(accounts[0])
	return err
}

func newAccInfo() error {
	_, err := acl.AccountInformation(pk)
	_, _ = kcl.DeleteKey(handle, walletPswd, pk)
	return err
}

func keyregTxnParams(ifee, ifv, ilv int, igh, ivotekey, iselkey string, ivotefst, ivotelst, ivotekd int, igen, inote string) error {
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
	votekey = ivotekey
	selkey = iselkey
	fee = uint64(ifee)
	fv = uint64(ifv)
	lv = uint64(ilv)
	votefst = uint64(ivotefst)
	votelst = uint64(ivotelst)
	votekd = uint64(ivotekd)
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

func createKeyregTxn() (err error) {
	paramsToUse := types.SuggestedParams{
		Fee:             types.MicroAlgos(fee),
		GenesisID:       gen,
		GenesisHash:     gh,
		FirstRoundValid: types.Round(fv),
		LastRoundValid:  types.Round(lv),
		FlatFee:         false,
	}
	txn, err = future.MakeKeyRegTxn(a.String(), note, paramsToUse, votekey, selkey, votefst, votelst, votekd)
	if err != nil {
		return err
	}
	return err
}

func getTxnsByCount(cnt int) error {
	_, err := acl.TransactionsByAddrLimit(accounts[0], uint64(cnt))
	return err
}

func createAssetTestFixture() error {
	assetTestFixture.Creator = ""
	assetTestFixture.AssetIndex = 1
	assetTestFixture.AssetName = "testcoin"
	assetTestFixture.AssetUnitName = "coins"
	assetTestFixture.AssetURL = "http://test"
	assetTestFixture.AssetMetadataHash = "fACPO4nRgO55j1ndAK3W6Sgc4APkcyFh"
	assetTestFixture.ExpectedParams = models.AssetParams{}
	assetTestFixture.QueriedParams = models.AssetParams{}
	assetTestFixture.LastTransactionIssued = types.Transaction{}
	return nil
}

func convertTransactionAssetParamsToModelsAssetParam(input types.AssetParams) models.AssetParams {
	result := models.AssetParams{
		Total:         input.Total,
		Decimals:      input.Decimals,
		DefaultFrozen: input.DefaultFrozen,
		ManagerAddr:   input.Manager.String(),
		ReserveAddr:   input.Reserve.String(),
		FreezeAddr:    input.Freeze.String(),
		ClawbackAddr:  input.Clawback.String(),
		UnitName:      input.UnitName,
		AssetName:     input.AssetName,
		URL:           input.URL,
		MetadataHash:  input.MetadataHash[:],
	}
	// input doesn't have Creator so that will remain empty
	return result
}

func assetCreateTxnHelper(issuance int, frozenState bool) error {
	accountToUse := accounts[0]
	assetTestFixture.Creator = accountToUse
	creator := assetTestFixture.Creator
	params, err := acl.BuildSuggestedParams()
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
	assetCreateTxn, err := future.MakeAssetCreateTxn(creator, assetNote, params, assetIssuance, 0, frozenState, manager, reserve, freeze, clawback, unitName, assetName, url, metadataHash)
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
	params, err := acl.BuildSuggestedParams()
	if err != nil {
		return err
	}
	lastRound = uint64(params.FirstRoundValid)
	assetNote := []byte(nil)
	reserve := ""
	freeze := ""
	clawback := ""
	manager := creator // if this were "" as well, this wouldn't be a reconfigure txn, it would be a destroy txn
	assetReconfigureTxn, err := future.MakeAssetConfigTxn(creator, assetNote, params, assetTestFixture.AssetIndex, manager, reserve, freeze, clawback, false)
	assetTestFixture.LastTransactionIssued = assetReconfigureTxn
	txn = assetReconfigureTxn
	// update expected params
	assetTestFixture.ExpectedParams.ReserveAddr = reserve
	assetTestFixture.ExpectedParams.FreezeAddr = freeze
	assetTestFixture.ExpectedParams.ClawbackAddr = clawback
	return err
}

func createAssetDestroy() error {
	creator := assetTestFixture.Creator
	params, err := acl.BuildSuggestedParams()
	if err != nil {
		return err
	}
	lastRound = uint64(params.FirstRoundValid)
	assetNote := []byte(nil)
	assetDestroyTxn, err := future.MakeAssetDestroyTxn(creator, assetNote, params, assetTestFixture.AssetIndex)
	assetTestFixture.LastTransactionIssued = assetDestroyTxn
	txn = assetDestroyTxn
	// update expected params
	assetTestFixture.ExpectedParams.ReserveAddr = ""
	assetTestFixture.ExpectedParams.FreezeAddr = ""
	assetTestFixture.ExpectedParams.ClawbackAddr = ""
	assetTestFixture.ExpectedParams.ManagerAddr = ""
	return err
}

// used in getAssetIndex and similar to get the index of the most recently operated on asset
func getMaxKey(numbers map[uint64]models.AssetParams) uint64 {
	var maxNumber uint64
	for n := range numbers {
		maxNumber = n
		break
	}
	for n := range numbers {
		if n > maxNumber {
			maxNumber = n
		}
	}
	return maxNumber
}

func getAssetIndex() error {
	accountResp, err := acl.AccountInformation(assetTestFixture.Creator)
	if err != nil {
		return err
	}
	// get most recent asset index
	assetTestFixture.AssetIndex = getMaxKey(accountResp.AssetParams)
	return nil
}

func getAssetInfo() error {
	response, err := acl.AssetInformation(assetTestFixture.AssetIndex)
	assetTestFixture.QueriedParams = response
	return err
}

func failToGetAssetInfo() error {
	_, err := acl.AssetInformation(assetTestFixture.AssetIndex)
	if err != nil {
		return nil
	}
	return fmt.Errorf("expected an error getting asset with index %v and creator %v, but no error was returned",
		assetTestFixture.AssetIndex, assetTestFixture.Creator)
}

func checkExpectedVsActualAssetParams() error {
	expectedParams := assetTestFixture.ExpectedParams
	actualParams := assetTestFixture.QueriedParams
	nameMatch := expectedParams.AssetName == actualParams.AssetName
	if !nameMatch {
		return fmt.Errorf("expected asset name was %v but actual asset name was %v",
			expectedParams.AssetName, actualParams.AssetName)
	}
	unitMatch := expectedParams.UnitName == actualParams.UnitName
	if !unitMatch {
		return fmt.Errorf("expected unit name was %v but actual unit name was %v",
			expectedParams.UnitName, actualParams.UnitName)
	}
	urlMatch := expectedParams.URL == actualParams.URL
	if !urlMatch {
		return fmt.Errorf("expected URL was %v but actual URL was %v",
			expectedParams.URL, actualParams.URL)
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
	managerMatch := expectedParams.ManagerAddr == actualParams.ManagerAddr
	if !managerMatch {
		return fmt.Errorf("expected asset manager was %v but actual asset manager was %v",
			expectedParams.ManagerAddr, actualParams.ManagerAddr)
	}
	reserveMatch := expectedParams.ReserveAddr == actualParams.ReserveAddr
	if !reserveMatch {
		return fmt.Errorf("expected asset reserve was %v but actual asset reserve was %v",
			expectedParams.ReserveAddr, actualParams.ReserveAddr)
	}
	freezeMatch := expectedParams.FreezeAddr == actualParams.FreezeAddr
	if !freezeMatch {
		return fmt.Errorf("expected freeze manager was %v but actual freeze manager was %v",
			expectedParams.FreezeAddr, actualParams.FreezeAddr)
	}
	clawbackMatch := expectedParams.ClawbackAddr == actualParams.ClawbackAddr
	if !clawbackMatch {
		return fmt.Errorf("expected revocation (clawback) manager was %v but actual revocation manager was %v",
			expectedParams.ClawbackAddr, actualParams.ClawbackAddr)
	}
	return nil
}

func theCreatorShouldHaveAssetsRemaining(expectedBal int) error {
	expectedBalance := uint64(expectedBal)
	accountResp, err := acl.AccountInformation(assetTestFixture.Creator)
	if err != nil {
		return err
	}
	holding, ok := accountResp.Assets[assetTestFixture.AssetIndex]
	if !ok {
		return fmt.Errorf("attempted to get balance of account %v for creator %v and index %v, but no balance was found for that index", assetTestFixture.Creator, assetTestFixture.Creator, assetTestFixture.AssetIndex)
	}
	if holding.Amount != expectedBalance {
		return fmt.Errorf("actual balance %v differed from expected balance %v", holding.Amount, expectedBalance)
	}
	return nil
}

func createAssetAcceptanceForSecondAccount() error {
	accountToUse := accounts[1]
	params, err := acl.BuildSuggestedParams()
	if err != nil {
		return err
	}
	lastRound = uint64(params.FirstRoundValid)
	assetNote := []byte(nil)
	assetAcceptanceTxn, err := future.MakeAssetAcceptanceTxn(accountToUse, assetNote, params, assetTestFixture.AssetIndex)
	assetTestFixture.LastTransactionIssued = assetAcceptanceTxn
	txn = assetAcceptanceTxn
	return err
}

func createAssetTransferTransactionToSecondAccount(amount int) error {
	recipient := accounts[1]
	creator := assetTestFixture.Creator
	params, err := acl.BuildSuggestedParams()
	if err != nil {
		return err
	}
	sendAmount := uint64(amount)
	closeAssetsTo := ""
	lastRound = uint64(params.FirstRoundValid)
	assetNote := []byte(nil)
	assetAcceptanceTxn, err := future.MakeAssetTransferTxn(creator, recipient, sendAmount, assetNote, params, closeAssetsTo, assetTestFixture.AssetIndex)
	assetTestFixture.LastTransactionIssued = assetAcceptanceTxn
	txn = assetAcceptanceTxn
	return err
}

func createAssetTransferTransactionFromSecondAccountToCreator(amount int) error {
	recipient := assetTestFixture.Creator
	sender := accounts[1]
	params, err := acl.BuildSuggestedParams()
	if err != nil {
		return err
	}
	sendAmount := uint64(amount)
	closeAssetsTo := ""
	lastRound = uint64(params.FirstRoundValid)
	assetNote := []byte(nil)
	assetAcceptanceTxn, err := future.MakeAssetTransferTxn(sender, recipient, sendAmount, assetNote, params, closeAssetsTo, assetTestFixture.AssetIndex)
	assetTestFixture.LastTransactionIssued = assetAcceptanceTxn
	txn = assetAcceptanceTxn
	return err
}

// sets up a freeze transaction, with freeze state `setting` against target account `target`
// assumes creator is asset freeze manager
func freezeTransactionHelper(target string, setting bool) error {
	params, err := acl.BuildSuggestedParams()
	if err != nil {
		return err
	}
	lastRound = uint64(params.FirstRoundValid)
	assetNote := []byte(nil)
	assetFreezeOrUnfreezeTxn, err := future.MakeAssetFreezeTxn(assetTestFixture.Creator, assetNote, params, assetTestFixture.AssetIndex, target, setting)
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
	params, err := acl.BuildSuggestedParams()
	if err != nil {
		return err
	}
	lastRound = uint64(params.FirstRoundValid)
	revocationAmount := uint64(amount)
	assetNote := []byte(nil)
	assetRevokeTxn, err := future.MakeAssetRevocationTxn(assetTestFixture.Creator, accounts[1], revocationAmount, assetTestFixture.Creator, assetNote, params, assetTestFixture.AssetIndex)
	assetTestFixture.LastTransactionIssued = assetRevokeTxn
	txn = assetRevokeTxn
	return err
}

func createContractTestFixture() error {
	contractTestFixture.split = templates.Split{}
	contractTestFixture.htlc = templates.HTLC{}
	contractTestFixture.periodicPay = templates.PeriodicPayment{}
	contractTestFixture.limitOrder = templates.LimitOrder{}
	contractTestFixture.dynamicFee = templates.DynamicFee{}
	contractTestFixture.activeAddress = ""
	contractTestFixture.htlcPreImage = ""
	contractTestFixture.limitOrderN = 0
	contractTestFixture.limitOrderD = 0
	contractTestFixture.limitOrderMin = 0
	contractTestFixture.splitN = 0
	contractTestFixture.splitD = 0
	contractTestFixture.splitMin = 0
	contractTestFixture.contractFundAmount = 0
	contractTestFixture.periodicPayPeriod = 0
	return nil
}

func aSplitContractWithRatioToAndMinimumPayment(ratn, ratd, minPay int) error {
	owner := accounts[0]
	receivers := [2]string{accounts[0], accounts[1]}
	expiryRound := uint64(100)
	maxFee := uint64(5000000)
	contractTestFixture.splitN = uint64(ratn)
	contractTestFixture.splitD = uint64(ratd)
	contractTestFixture.splitMin = uint64(minPay)
	c, err := templates.MakeSplit(owner, receivers[0], receivers[1], uint64(ratn), uint64(ratd), expiryRound, uint64(minPay), maxFee)
	contractTestFixture.split = c
	contractTestFixture.activeAddress = c.GetAddress()
	// add much more than enough to evenly split
	contractTestFixture.contractFundAmount = uint64(minPay*(ratd+ratn)) * 10
	return err
}

func iSendTheSplitTransactions() error {
	amount := contractTestFixture.splitMin * (contractTestFixture.splitN + contractTestFixture.splitD)
	params, err := acl.BuildSuggestedParams()
	if err != nil {
		return err
	}
	lastRound = uint64(params.FirstRoundValid)
	txnBytes, err := templates.GetSplitFundsTransaction(contractTestFixture.split.GetProgram(), amount, params)
	if err != nil {
		return err
	}
	txIdResponse, err := acl.SendRawTransaction(txnBytes)
	txid = txIdResponse.TxID
	// hack to make checkTxn work
	backupTxnSender = contractTestFixture.split.GetAddress()
	return err
}

func anHTLCContractWithHashPreimage(preImage string) error {
	hashImage := sha256.Sum256([]byte(preImage))
	owner := accounts[0]
	receiver := accounts[1]
	hashFn := "sha256"
	expiryRound := uint64(100)
	maxFee := uint64(1000000)
	hashB64 := base64.StdEncoding.EncodeToString(hashImage[:])
	c, err := templates.MakeHTLC(owner, receiver, hashFn, hashB64, expiryRound, maxFee)
	contractTestFixture.htlc = c
	contractTestFixture.htlcPreImage = preImage
	contractTestFixture.activeAddress = c.GetAddress()
	contractTestFixture.contractFundAmount = 100000000
	return err
}

func iFundTheContractAccount() error {
	// send the requested money to c.address
	amount := contractTestFixture.contractFundAmount
	params, err := acl.BuildSuggestedParams()
	if err != nil {
		return err
	}
	lastRound = uint64(params.FirstRoundValid)
	txn, err = future.MakePaymentTxn(accounts[0], contractTestFixture.activeAddress, amount, nil, "", params)
	if err != nil {
		return err
	}
	err = signKmd()
	if err != nil {
		return err
	}
	err = sendTxnKmd()
	if err != nil {
		return err
	}
	return checkTxn()
}

// used in HTLC
func iClaimTheAlgosHTLC() error {
	preImage := contractTestFixture.htlcPreImage
	preImageAsArgument := []byte(preImage)
	args := make([][]byte, 1)
	args[0] = preImageAsArgument
	receiver := accounts[1] // was set as receiver in setup step
	var blankMultisig crypto.MultisigAccount
	lsig, err := crypto.MakeLogicSig(contractTestFixture.htlc.GetProgram(), args, nil, blankMultisig)
	if err != nil {
		return err
	}
	params, err := acl.BuildSuggestedParams()
	if err != nil {
		return err
	}
	lastRound = uint64(params.FirstRoundValid)
	txn, err = future.MakePaymentTxn(contractTestFixture.activeAddress, receiver, 0, nil, receiver, params)
	if err != nil {
		return err
	}
	txn.Receiver = types.Address{} //txn must have no receiver but MakePayment disallows this.
	txid, stx, err = crypto.SignLogicsigTransaction(lsig, txn)
	if err != nil {
		return err
	}
	return sendTxn()
}

func aPeriodicPaymentContractWithWithdrawingWindowAndPeriod(withdrawWindow, period int) error {
	receiver := accounts[0]
	amount := uint64(10000000)
	// add far more than enough to withdraw
	contractTestFixture.contractFundAmount = amount * 10
	expiryRound := uint64(100)
	maxFee := uint64(1000000000000)
	contract, err := templates.MakePeriodicPayment(receiver, amount, uint64(withdrawWindow), uint64(period), expiryRound, maxFee)
	contractTestFixture.activeAddress = contract.GetAddress()
	contractTestFixture.periodicPay = contract
	contractTestFixture.periodicPayPeriod = uint64(period)
	return err
}

func iClaimThePeriodicPayment() error {
	params, err := acl.BuildSuggestedParams()
	if err != nil {
		return err
	}
	txnFirstValid := uint64(params.FirstRoundValid)
	remainder := txnFirstValid % contractTestFixture.periodicPayPeriod
	txnFirstValid += remainder
	stx, err = templates.GetPeriodicPaymentWithdrawalTransaction(contractTestFixture.periodicPay.GetProgram(), txnFirstValid, uint64(params.Fee), params.GenesisHash)
	if err != nil {
		return err
	}
	lastRound = uint64(params.FirstRoundValid) // used in send/checkTxn
	return sendTxn()
}

func aLimitOrderContractWithParameters(ratn, ratd, minTrade int) error {
	maxFee := uint64(100000)
	expiryRound := uint64(100)
	contractTestFixture.limitOrderN = uint64(ratn)
	contractTestFixture.limitOrderD = uint64(ratd)
	contractTestFixture.limitOrderMin = uint64(minTrade)
	contractTestFixture.contractFundAmount = 2 * uint64(minTrade)
	if contractTestFixture.contractFundAmount < 1000000 {
		contractTestFixture.contractFundAmount = 1000000
	}
	contract, err := templates.MakeLimitOrder(accounts[0], assetTestFixture.AssetIndex, uint64(ratn), uint64(ratd), expiryRound, uint64(minTrade), maxFee)
	contractTestFixture.activeAddress = contract.GetAddress()
	contractTestFixture.limitOrder = contract
	return err
}

// godog misreads the step for this function, so provide a handler for when it does so
func iCreateATransactionTransferringAmountAssetsFromCreatorToASecondAccount() error {
	return createAssetTransferTransactionToSecondAccount(500000)
}

func iSwapAssetsForAlgos() error {
	exp, err := kcl.ExportKey(handle, walletPswd, accounts[1])
	if err != nil {
		return err
	}
	secretKey := exp.PrivateKey
	params, err := acl.BuildSuggestedParams()
	if err != nil {
		return err
	}
	lastRound = uint64(params.FirstRoundValid)
	contract := contractTestFixture.limitOrder.GetProgram()
	microAlgoAmount := contractTestFixture.limitOrderMin + 1 // just over the minimum
	assetAmount := microAlgoAmount * contractTestFixture.limitOrderN / contractTestFixture.limitOrderD
	assetAmount += 1 // assetAmount initialized to absolute minimum, will fail greater-than check, so increment by one for a better deal
	stx, err = contractTestFixture.limitOrder.GetSwapAssetsTransaction(assetAmount, microAlgoAmount, contract, secretKey, params)
	if err != nil {
		return err
	}
	// hack to make checktxn work
	txn = types.Transaction{}
	backupTxnSender = contractTestFixture.limitOrder.GetAddress() // used in checktxn
	return sendTxn()
}

func aDynamicFeeContractWithAmount(amount int) error {
	params, err := acl.BuildSuggestedParams()
	if err != nil {
		return err
	}
	lastRound = uint64(params.FirstRoundValid)
	txnFirstValid := lastRound
	txnLastValid := txnFirstValid + 10
	contractTestFixture.contractFundAmount = uint64(10 * amount)
	contract, err := templates.MakeDynamicFee(accounts[1], "", uint64(amount), txnFirstValid, txnLastValid)

	contractTestFixture.dynamicFee = contract
	contractTestFixture.activeAddress = contract.GetAddress()
	return err
}

func iSendTheDynamicFeeTransaction() error {
	params, err := acl.BuildSuggestedParams()
	if err != nil {
		return err
	}
	lastRound = uint64(params.FirstRoundValid)
	exp, err := kcl.ExportKey(handle, walletPswd, accounts[0])
	if err != nil {
		return err
	}
	secretKeyOne := exp.PrivateKey
	initialTxn, lsig, err := templates.SignDynamicFee(contractTestFixture.dynamicFee.GetProgram(), secretKeyOne, params.GenesisHash)
	if err != nil {
		return err
	}
	exp, err = kcl.ExportKey(handle, walletPswd, accounts[1])
	if err != nil {
		return err
	}
	secretKeyTwo := exp.PrivateKey
	groupTxnBytes, err := templates.GetDynamicFeeTransactions(initialTxn, lsig, secretKeyTwo, uint64(params.Fee))
	// hack to make checkTxn work
	txn = initialTxn
	// end hack to make checkTxn work
	response, err := acl.SendRawTransaction(groupTxnBytes)
	txid = response.TxID
	return err
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
	data := msgpack.Encode(&ddr)

	result, err := aclv2.TealDryrun(data).Do(context.Background())
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
