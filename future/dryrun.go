package future

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/types"
)

const (
	defaultAppId uint64 = 1380011588
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
			return drr, fmt.Errorf("faiiled to decode creator address %s: %+v", appInfo.Params.Creator, err)
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
