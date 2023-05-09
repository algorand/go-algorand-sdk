package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/common"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/indexer"
)

func main() {
	// example: CREATE_INDEXER_CLIENT
	// Create a new indexer client, configured to connect to out local sandbox
	var indexerAddress = "http://localhost:8980"
	var indexerToken = strings.Repeat("a", 64)
	indexerClient, err := indexer.MakeClient(
		indexerAddress,
		indexerToken,
	)

	// Or, if necessary, pass alternate headers

	var indexerHeader common.Header
	indexerHeader.Key = "X-API-Key"
	indexerHeader.Value = indexerToken
	indexerClientWithHeaders, err := indexer.MakeClientWithHeaders(
		indexerAddress,
		indexerToken,
		[]*common.Header{&indexerHeader},
	)
	// example: CREATE_INDEXER_CLIENT

	// Suppress `indexerClientWithHeaders declared but not used`
	_ = indexerClientWithHeaders

	if err != nil {
		fmt.Printf("failed to make indexer client: %s\n", err)
		return
	}

	indexerHealth, err := indexerClient.HealthCheck().Do(context.Background())
	if err != nil {
		fmt.Printf("Failed to get status: %s\n", err)
		return
	}

	fmt.Printf("Indexer Round: %d\n", indexerHealth.Round)

	// example: INDEXER_LOOKUP_ASSET
	// query parameters
	var assetId uint64 = 2044572
	var minBalance uint64 = 50

	// Lookup accounts with minimum balance of asset
	assetResult, err := indexerClient.
		LookupAssetBalances(assetId).
		CurrencyGreaterThan(minBalance).
		Do(context.Background())

	// Print the results
	assetJson, err := json.MarshalIndent(assetResult, "", "\t")
	fmt.Printf(string(assetJson) + "\n")
	// example: INDEXER_LOOKUP_ASSET

	assetJson = nil

	// example: INDEXER_SEARCH_MIN_AMOUNT
	// query parameters
	var transactionMinAmount uint64 = 10

	// Query
	transactionResult, err := indexerClient.
		SearchForTransactions().
		CurrencyGreaterThan(transactionMinAmount).
		Do(context.Background())

	// Print results
	transactionJson, err := json.MarshalIndent(transactionResult, "", "\t")
	fmt.Printf(string(transactionJson) + "\n")
	// example: INDEXER_SEARCH_MIN_AMOUNT

	// example: INDEXER_PAGINATE_RESULTS
	var nextToken = ""
	var numTx = 1
	var numPages = 1
	var pagedMinAmount uint64 = 10
	var limit uint64 = 1

	for numTx > 0 {
		// Query
		pagedResults, err := indexerClient.
			SearchForTransactions().
			CurrencyGreaterThan(pagedMinAmount).
			Limit(limit).
			NextToken(nextToken).
			Do(context.Background())
		if err != nil {
			return
		}
		pagedTransactions := pagedResults.Transactions
		numTx = len(pagedTransactions)
		nextToken = pagedResults.NextToken

		if numTx > 0 {
			// Print results
			pagedJson, err := json.MarshalIndent(pagedTransactions, "", "\t")
			if err != nil {
				return
			}
			fmt.Printf(string(pagedJson) + "\n")
			fmt.Println("End of page : ", numPages)
			fmt.Println("Transaction printed : ", len(pagedTransactions))
			fmt.Println("Next Token : ", nextToken)
			numPages++
		}
	}
	// example: INDEXER_PAGINATE_RESULTS

	// example: INDEXER_PREFIX_SEARCH
	// Parameters
	var notePrefix = "showing prefix"

	// Query
	prefixResult, err := indexerClient.
		SearchForTransactions().
		NotePrefix([]byte(notePrefix)).
		Do(context.Background())

	// Print results
	prefixJson, err := json.MarshalIndent(prefixResult, "", "\t")
	fmt.Printf(string(prefixJson) + "\n")
	// example: INDEXER_PREFIX_SEARCH
}
