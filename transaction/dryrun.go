package transaction

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

const (
	defaultAppID uint64 = 1380011588

	rejectMsg       = "REJECT"
	defaultMaxWidth = 30
)

// CreateDryrun creates a DryrunRequest object from a client and slice of SignedTxn objects and a default configuration
// Passed in as a pointer to a DryrunRequest object to use for extra parameters
func CreateDryrun(client *algod.Client, txns []types.SignedTxn, dr *models.DryrunRequest, ctx context.Context) (drr models.DryrunRequest, err error) { //nolint:revive // Ignore Context order for backwards compatibility
	var (
		apps   []types.AppIndex
		assets []types.AssetIndex
		accts  []types.Address
	)

	drr.Txns = txns

	if dr != nil {
		drr.LatestTimestamp = dr.LatestTimestamp
		drr.Round = dr.Round
		drr.ProtocolVersion = dr.ProtocolVersion
		drr.Sources = dr.Sources
	}

	for _, t := range txns {
		if t.Txn.Type != types.ApplicationCallTx {
			continue
		}

		accts = append(accts, t.Txn.Sender)

		accts = append(accts, t.Txn.Accounts...)
		apps = append(apps, t.Txn.ForeignApps...)

		for _, aidx := range t.Txn.ForeignApps {
			accts = append(accts, crypto.GetApplicationAddress(uint64(aidx)))
		}

		assets = append(assets, t.Txn.ForeignAssets...)

		if t.Txn.ApplicationID == 0 {
			drr.Apps = append(drr.Apps, models.Application{
				Id: defaultAppID,
				Params: models.ApplicationParams{
					Creator:           t.Txn.Sender.String(),
					ApprovalProgram:   t.Txn.ApprovalProgram,
					ClearStateProgram: t.Txn.ClearStateProgram,
					LocalStateSchema: models.ApplicationStateSchema{
						NumByteSlice: t.Txn.LocalStateSchema.NumByteSlice,
						NumUint:      t.Txn.LocalStateSchema.NumUint,
					},
					GlobalStateSchema: models.ApplicationStateSchema{
						NumByteSlice: t.Txn.GlobalStateSchema.NumByteSlice,
						NumUint:      t.Txn.GlobalStateSchema.NumUint,
					},
				},
			})
		} else {
			apps = append(apps, t.Txn.ApplicationID)
			accts = append(accts, crypto.GetApplicationAddress(uint64(t.Txn.ApplicationID)))
		}
	}

	seenAssets := map[types.AssetIndex]bool{}
	for _, assetID := range assets {
		if _, ok := seenAssets[assetID]; ok {
			continue
		}

		assetInfo, err := client.GetAssetByID(uint64(assetID)).Do(ctx)
		if err != nil {
			return drr, fmt.Errorf("failed to get asset %d: %+v", assetID, err)
		}

		addr, err := types.DecodeAddress(assetInfo.Params.Creator)
		if err != nil {
			return drr, fmt.Errorf("failed to decode creator adddress %s: %+v", assetInfo.Params.Creator, err)
		}

		accts = append(accts, addr)
		seenAssets[assetID] = true
	}

	seenApps := map[types.AppIndex]bool{}
	for _, appID := range apps {
		if _, ok := seenApps[appID]; ok {
			continue
		}

		appInfo, err := client.GetApplicationByID(uint64(appID)).Do(ctx)
		if err != nil {
			return drr, fmt.Errorf("failed to get application %d: %+v", appID, err)
		}
		drr.Apps = append(drr.Apps, appInfo)

		creator, err := types.DecodeAddress(appInfo.Params.Creator)
		if err != nil {
			return drr, fmt.Errorf("failed to decode creator address %s: %+v", appInfo.Params.Creator, err)
		}
		accts = append(accts, creator)

		seenApps[appID] = true
	}

	seenAccts := map[types.Address]bool{}
	for _, acct := range accts {
		if _, ok := seenAccts[acct]; ok {
			continue
		}
		acctInfo, err := client.AccountInformation(acct.String()).Do(ctx)
		if err != nil {
			return drr, fmt.Errorf("failed to get application %s: %+v", acct, err)
		}
		drr.Accounts = append(drr.Accounts, acctInfo)
		seenAccts[acct] = true
	}

	return
}

// StackPrinterConfig holds some configuration parameters for how to print a DryrunResponse stack trace
type StackPrinterConfig struct {
	MaxValueWidth   int  // Set the max width of the column, 0 is no max
	TopOfStackFirst bool // Set the order of the stack values printed, true is top of stack (last pushed) first
}

// DefaultStackPrinterConfig returns a new StackPrinterConfig with reasonable defaults
func DefaultStackPrinterConfig() StackPrinterConfig {
	return StackPrinterConfig{MaxValueWidth: defaultMaxWidth, TopOfStackFirst: true}
}

// DryrunResponse represents the response from a dryrun call
type DryrunResponse struct {
	Error           string            `json:"error"`
	ProtocolVersion string            `json:"protocol-version"`
	Txns            []DryrunTxnResult `json:"txns"`
}

// NewDryrunResponse creates a new DryrunResponse from a models.DryrunResponse
func NewDryrunResponse(d models.DryrunResponse) (DryrunResponse, error) {
	// Marshal and unmarshal to fix integer types.
	b, err := json.Marshal(d)
	if err != nil {
		return DryrunResponse{}, err
	}
	return NewDryrunResponseFromJSON(b)
}

// NewDryrunResponseFromJSON creates a new DryrunResponse from a JSON byte array
func NewDryrunResponseFromJSON(js []byte) (DryrunResponse, error) {
	dr := DryrunResponse{}
	err := json.Unmarshal(js, &dr)
	return dr, err
}

// DryrunTxnResult is a wrapper around models.DryrunTxnResult
type DryrunTxnResult struct {
	models.DryrunTxnResult
}

// AppCallRejected returns true if the Application Call was rejected
// for this transaction
func (d *DryrunTxnResult) AppCallRejected() bool {
	for _, m := range d.AppCallMessages {
		if m == rejectMsg {
			return true
		}
	}
	return false
}

// LogicSigRejected returns true if the LogicSig was rejected
// for this transaction
func (d *DryrunTxnResult) LogicSigRejected() bool {
	for _, m := range d.LogicSigMessages {
		if m == rejectMsg {
			return true
		}
	}
	return false
}

type indexedScratchValue struct {
	Idx int
	Val models.TealValue
}

func scratchToString(prevScratch, currScratch []models.TealValue) string {
	if len(currScratch) == 0 {
		return ""
	}

	var newScratchVal *indexedScratchValue

	prevScratchMap := map[indexedScratchValue]bool{}
	for idx, psv := range prevScratch {
		isv := indexedScratchValue{Idx: idx, Val: psv}
		prevScratchMap[isv] = true
	}

	for idx, csv := range currScratch {
		isv := indexedScratchValue{Idx: idx, Val: csv}
		if _, ok := prevScratchMap[isv]; !ok {
			newScratchVal = &isv
		}
	}

	if newScratchVal == nil {
		return ""
	}

	if len(newScratchVal.Val.Bytes) > 0 {
		decoded, _ := base64.StdEncoding.DecodeString(newScratchVal.Val.Bytes)
		return fmt.Sprintf("%d = %#x", newScratchVal.Idx, decoded)
	}

	return fmt.Sprintf("%d = %d", newScratchVal.Idx, newScratchVal.Val.Uint)
}

func stackToString(reverse bool, stack []models.TealValue) string {
	elems := len(stack)
	svs := make([]string, elems)
	for idx, s := range stack {

		svidx := idx
		if reverse {
			svidx = (elems - 1) - idx
		}

		if s.Type == 1 {
			// Just returns empty string if there is an error, use it
			decoded, _ := base64.StdEncoding.DecodeString(s.Bytes)
			svs[svidx] = fmt.Sprintf("%#x", decoded)
		} else {
			svs[svidx] = fmt.Sprintf("%d", s.Uint)
		}
	}

	return fmt.Sprintf("[%s]", strings.Join(svs, ", "))
}

func (d *DryrunTxnResult) trace(state []models.DryrunState, disassemmbly []string, spc StackPrinterConfig) string {
	buff := bytes.NewBuffer(nil)
	w := tabwriter.NewWriter(buff, 0, 0, 1, ' ', tabwriter.Debug)

	fmt.Fprintln(w, "pc#\tln#\tsource\tscratch\tstack")
	for idx, s := range state {

		prevScratch := []models.TealValue{}
		if idx > 0 {
			prevScratch = state[idx-1].Scratch
		}

		src := disassemmbly[s.Line]
		if s.Error != "" {
			src = fmt.Sprintf("!! %s !!", s.Error)
		}

		srcLine := fmt.Sprintf("%d\t%d\t%s\t%s\t%s",
			s.Pc, s.Line,
			truncate(src, spc.MaxValueWidth),
			truncate(scratchToString(prevScratch, s.Scratch), spc.MaxValueWidth),
			truncate(stackToString(spc.TopOfStackFirst, s.Stack), spc.MaxValueWidth))

		fmt.Fprintln(w, srcLine)
	}
	w.Flush()

	return buff.String()
}

// GetAppCallTrace returns a string representing a stack trace for this transactions
// application logic evaluation
func (d *DryrunTxnResult) GetAppCallTrace(spc StackPrinterConfig) string {
	return d.trace(d.AppCallTrace, d.Disassembly, spc)
}

// GetLogicSigTrace returns a string representing a stack trace for this transactions
// logic signature evaluation
func (d *DryrunTxnResult) GetLogicSigTrace(spc StackPrinterConfig) string {
	return d.trace(d.LogicSigTrace, d.LogicSigDisassembly, spc)
}

func truncate(str string, width int) string {
	if len(str) > width && width > 0 {
		return str[:width] + "..."
	}
	return str
}
