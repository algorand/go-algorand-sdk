package main

import (
	"context"
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/types"
)

// Missing APP_SCHEMA in GOSDK examples (in ../docs/get-details/dapps/smart-contracts/frontend/apps.md:163)
// Missing APP_SOURCE in GOSDK examples (in ../docs/get-details/dapps/smart-contracts/frontend/apps.md:223)
// Missing APP_COMPILE in GOSDK examples (in ../docs/get-details/dapps/smart-contracts/frontend/apps.md:336)
// Missing APP_CREATE in GOSDK examples (in ../docs/get-details/dapps/smart-contracts/frontend/apps.md:431)
// Missing APP_OPTIN in GOSDK examples (in ../docs/get-details/dapps/smart-contracts/frontend/apps.md:566)
// Missing APP_NOOP in GOSDK examples (in ../docs/get-details/dapps/smart-contracts/frontend/apps.md:641)
// Missing APP_READ_STATE in GOSDK examples (in ../docs/get-details/dapps/smart-contracts/frontend/apps.md:707)
// Missing APP_UPDATE in GOSDK examples (in ../docs/get-details/dapps/smart-contracts/frontend/apps.md:814)
// Missing APP_CALL in GOSDK examples (in ../docs/get-details/dapps/smart-contracts/frontend/apps.md:920)
// Missing APP_CLOSEOUT in GOSDK examples (in ../docs/get-details/dapps/smart-contracts/frontend/apps.md:994)
// Missing APP_DELETE in GOSDK examples (in ../docs/get-details/dapps/smart-contracts/frontend/apps.md:1057)
// Missing APP_CLEAR in GOSDK examples (in ../docs/get-details/dapps/smart-contracts/frontend/apps.md:1116)

func main() {

	algodClient := getAlgodClient()

	nodeStatus, err := algodClient.Status().Do(context.Background())
	if err != nil {
		fmt.Printf("Failed to get status: %s\n", err)
		return
	}

	fmt.Printf("Last Round: %d\n", nodeStatus.LastRound)

	// example: GO_STATE_SCHEMA
	// declare application state storage (immutable)
	const localInts = 1
	const localBytes = 1
	const globalInts = 1
	const globalBytes = 0

	// define schema
	globalSchema := types.StateSchema{NumUint: uint64(globalInts), NumByteSlice: uint64(globalBytes)}
	localSchema := types.StateSchema{NumUint: uint64(localInts), NumByteSlice: uint64(localBytes)}
	// example: GO_STATE_SCHEMA
}
