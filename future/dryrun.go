package future

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/types"
)

// CreateDryrun creates a DryrunRequest object from a client and slice of SignedTxn objects and a default configuration
// Passed in as a pointer to a DryrunRequest object to use for extra parameters
func CreateDryrun(client algod.Client, txns []types.SignedTxn, dr *models.DryrunRequest) (drr models.DryrunRequest, err error) {
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
			appId := 1 //TODO: make the magic number?
			drr.Apps = append(drr.Apps, models.Application{
				Id: uint64(appId),
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
		}
	}

	for _, assetId := range assets {
		assetInfo, err := client.GetAssetByID(uint64(assetId)).Do(context.Background())
		if err != nil {
			return drr, fmt.Errorf("failed to get asset %d: %+v", assetId, err)
		}

		addr, err := types.DecodeAddress(assetInfo.Params.Creator)
		if err != nil {
			return drr, fmt.Errorf("failed to decode creator adddress %s: %+v", assetInfo.Params.Creator, err)
		}
		accts = append(accts, addr)
	}

	for _, appId := range apps {
		appInfo, err := client.GetApplicationByID(uint64(appId)).Do(context.Background())
		if err != nil {
			return drr, fmt.Errorf("failed to get application %d: %+v", appId, err)
		}

		drr.Apps = append(drr.Apps, appInfo)
	}

	for _, acct := range accts {
		acctInfo, err := client.AccountInformation(acct.String()).Do(context.Background())
		if err != nil {
			return drr, fmt.Errorf("failed to get application %s: %+v", acct, err)
		}

		drr.Accounts = append(drr.Accounts, acctInfo)
	}

	return
}
