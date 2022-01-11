package test

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/algorand/go-algorand-sdk/types"

	"github.com/cucumber/godog"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
)

func feeFieldIsInTxn() error {
	var txn map[string]interface{}
	err := msgpack.Decode(stx, &txn)
	if err != nil {
		return fmt.Errorf("Error while decoding txn. %v", err)
	}
	if _, ok := txn["txn"].(map[interface{}]interface{})["fee"]; !ok {
		return fmt.Errorf("fee field missing. %v", err)
	}
	return nil
}

func feeFieldNotInTxn() error {
	var txn map[string]interface{}
	err := msgpack.Decode(stx, &txn)
	if err != nil {
		return fmt.Errorf("Error while decoding txn. %v", err)
	}
	if _, ok := txn["txn"].(map[interface{}]interface{})["fee"]; ok {
		return fmt.Errorf("fee field found but it should have been omitted. %v", err)
	}
	return nil
}

func getSuggestedParams(
	fee, fv, lv uint64,
	gen, ghb64 string,
	flat bool) (types.SuggestedParams, error) {
	gh, err := base64.StdEncoding.DecodeString(ghb64)
	if err != nil {
		return types.SuggestedParams{}, err
	}
	return types.SuggestedParams{
		Fee:             types.MicroAlgos(fee),
		GenesisID:       gen,
		GenesisHash:     gh,
		FirstRoundValid: types.Round(fv),
		LastRoundValid:  types.Round(lv),
		FlatFee:         flat,
	}, err
}

func weMakeAGetAssetByIDCall(assetID int) error {
	clt, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	clt.GetAssetByID(uint64(assetID)).Do(context.Background())
	return nil
}

func weMakeAGetApplicationByIDCall(applicationID int) error {
	clt, err := algod.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	clt.GetApplicationByID(uint64(applicationID)).Do(context.Background())
	return nil
}

func weMakeASearchForApplicationsCall(applicationID int) error {
	clt, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	clt.SearchForApplications().ApplicationId(uint64(applicationID)).Do(context.Background())
	return nil
}

func weMakeALookupApplicationsCall(applicationID int) error {
	clt, err := indexer.MakeClient(mockServer.URL, "")
	if err != nil {
		return err
	}
	clt.LookupApplicationByID(uint64(applicationID)).Do(context.Background())
	return nil
}

func ApplicationsUnitContext(s *godog.Suite) {
	// @unit.transactions
	s.Step(`^fee field is in txn$`, feeFieldIsInTxn)
	s.Step(`^fee field not in txn$`, feeFieldNotInTxn)

	//@unit.applications
	s.Step(`^we make a GetAssetByID call for assetID (\d+)$`, weMakeAGetAssetByIDCall)
	s.Step(`^we make a GetApplicationByID call for applicationID (\d+)$`, weMakeAGetApplicationByIDCall)
	s.Step(`^we make a SearchForApplications call with applicationID (\d+)$`, weMakeASearchForApplicationsCall)
	s.Step(`^we make a LookupApplications call with applicationID (\d+)$`, weMakeALookupApplicationsCall)
}
