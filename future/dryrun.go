package future

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/types"
)

const (
	defaultAppId uint64 = 1380011588

	rejectMsg      = "REJECT"
	defaultSpacing = 20
)

// CreateDryrun creates a DryrunRequest object from a client and slice of SignedTxn objects and a default configuration
// Passed in as a pointer to a DryrunRequest object to use for extra parameters
func CreateDryrun(client *algod.Client, txns []types.SignedTxn, dr *models.DryrunRequest, ctx context.Context) (drr models.DryrunRequest, err error) {
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
		assets = append(assets, t.Txn.ForeignAssets...)

		if t.Txn.ApplicationID == 0 {
			drr.Apps = append(drr.Apps, models.Application{
				Id: defaultAppId,
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
	for _, assetId := range assets {
		if _, ok := seenAssets[assetId]; ok {
			continue
		}

		assetInfo, err := client.GetAssetByID(uint64(assetId)).Do(ctx)
		if err != nil {
			return drr, fmt.Errorf("failed to get asset %d: %+v", assetId, err)
		}

		addr, err := types.DecodeAddress(assetInfo.Params.Creator)
		if err != nil {
			return drr, fmt.Errorf("failed to decode creator adddress %s: %+v", assetInfo.Params.Creator, err)
		}

		accts = append(accts, addr)
		seenAssets[assetId] = true
	}

	seenApps := map[types.AppIndex]bool{}
	for _, appId := range apps {
		if _, ok := seenApps[appId]; ok {
			continue
		}

		appInfo, err := client.GetApplicationByID(uint64(appId)).Do(ctx)
		if err != nil {
			return drr, fmt.Errorf("failed to get application %d: %+v", appId, err)
		}
		drr.Apps = append(drr.Apps, appInfo)

		creator, err := types.DecodeAddress(appInfo.Params.Creator)
		if err != nil {
			return drr, fmt.Errorf("failed to decode creator address %s: %+v", appInfo.Params.Creator, err)
		}
		accts = append(accts, creator)

		seenApps[appId] = true
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

type DryrunResponse struct {
	Error           string            `json:"error"`
	ProtocolVersion string            `json:"protocol-version"`
	Txns            []DryrunTxnResult `json:"txns"`
}

func NewDryrunResponseFromJson(js []byte) (DryrunResponse, error) {
	dr := DryrunResponse{}
	err := json.Unmarshal(js, &dr)
	return dr, err
}

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

func stackToString(stack []models.TealValue) string {
	svs := []string{}
	for _, s := range stack {
		if s.Type == 1 {
			// Just returns empty string if there is an error, use it
			decoded, _ := base64.StdEncoding.DecodeString(s.Bytes)
			svs = append(svs, fmt.Sprintf("0x%x", decoded))
		} else {
			svs = append(svs, fmt.Sprintf("%d", s.Uint))
		}
	}
	return fmt.Sprintf("[%s]", strings.Join(svs, ", "))
}

func (d *DryrunTxnResult) trace(state []models.DryrunState, disassemmbly []string, spaces ...int) string {

	var padSpacing int
	if len(spaces) == 0 {
		padSpacing = defaultSpacing
	} else {
		padSpacing = spaces[0]
	}

	if padSpacing == 0 {
		for _, l := range disassemmbly {
			if len(l) > padSpacing {
				padSpacing = len(l)
			}
		}
		//  4 for line number + 1 for space between number and line, *2 for pc
		padSpacing += 10
	}

	traceLines := []string{"pc# line# source" + strings.Repeat(" ", padSpacing-16) + "stack"}

	for _, s := range state {
		lineNumberPadding := strings.Repeat(" ", 4-len(strconv.Itoa(int(s.Line))))
		pcNumberPadding := strings.Repeat(" ", 4-len(strconv.Itoa(int(s.Pc))))

		srcLine := fmt.Sprintf("%s%d %s%d %s", pcNumberPadding, s.Pc, lineNumberPadding, s.Line, disassemmbly[s.Line])

		stack := stackToString(s.Stack)
		padding := ""
		if repeat := padSpacing - len(srcLine); repeat > 0 {
			padding = strings.Repeat(" ", repeat)
		}
		traceLines = append(traceLines, fmt.Sprintf("%s%s%s", srcLine, padding, stack))
	}

	return strings.Join(traceLines, "\n")
}

// GetAppCallTrace returns a string representing a stack trace for this transactions
// application logic evaluation
func (d *DryrunTxnResult) GetAppCallTrace(spaces ...int) string {
	return d.trace(d.AppCallTrace, d.Disassembly, spaces...)
}

// GetLogicSigTrace returns a string representing a stack trace for this transactions
// logic signature evaluation
func (d *DryrunTxnResult) GetLogicSigTrace(spaces ...int) string {
	// TODO: regen structs with LsigDisassembly
	return d.trace(d.LogicSigTrace, d.Disassembly, spaces...)
}
