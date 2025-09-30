# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is the official Algorand Go SDK, providing HTTP clients for algod (agreement) and kmd (key management) APIs, plus standalone functionality for interacting with the Algorand protocol including transaction signing, message encoding, and crypto utilities.

## Common Commands

### Building and Testing
```bash
# Build the project (includes code generation)
make build

# Run all tests (unit + go tests)
make test

# Run unit tests only (fast, for development)
make unit

# Run integration tests (requires test harness)
make integration

# Run tests in Docker with full test environment
make docker-test
```

### Linting and Formatting
```bash
# Run linters (required before PR submission)
make lint

# Format code
make fmt

# Generate code (TEAL langspec and other generated files)
make generate
```

### Test Harness Management
```bash
# Set up test environment (clones algorand-sdk-testing and stands up sandbox)
make harness

# Tear down test environment
make harness-down

# Alternative: Use test-harness.sh directly
./test-harness.sh up
./test-harness.sh down
```

### Running Specific Tests
```bash
# Run a specific test by name
go test -run TestName ./path/to/package

# Run cucumber tests with specific tags
cd test && go test -timeout 0s --godog.tags="@unit.algod" --test.v .

# Run smoke tests for examples
make smoke-test-examples
```

## Architecture

### Package Structure

- **`client/`**: REST API clients
  - `kmd/`: Key Management Daemon v1 client for wallet and key management
  - `v2/algod/`: Algod v2 client for blockchain operations (querying blocks, submitting transactions)
  - `v2/indexer/`: Indexer v2 client for historical blockchain data queries
  - `v2/common/models/`: Auto-generated data models from OpenAPI specs

- **`types/`**: Core Algorand data structures (Address, Transaction, MultisigSig, etc.)

- **`transaction/`**: Transaction building utilities
  - `transaction.go`: Core transaction construction functions
  - `atomicTransactionComposer.go`: ATC for atomic transaction groups and ABI calls
  - `transactionSigner.go`: Transaction signing interfaces
  - `dryrun.go`: Dryrun request utilities

- **`encoding/`**: Serialization utilities
  - `json/`: JSON encoding/decoding
  - `msgpack/`: MessagePack encoding with canonical field ordering

- **`crypto/`**: Cryptographic operations (ED25519, multisig, LogicSig)

- **`mnemonic/`**: BIP39 mnemonic phrase utilities for key backup/recovery

- **`abi/`**: ABI (Application Binary Interface) encoding/decoding per ARC-004

- **`logic/`**: TEAL program utilities including source maps

- **`protocol/`**: Protocol-level constants and configuration

- **`examples/`**: Working examples demonstrating SDK usage

### Code Generation

The `client/v2/common/models/` package is auto-generated from OpenAPI specifications. A GitHub Action runs nightly to regenerate code from the latest specs. Do not manually edit generated files in `client/v2/`.

### Cucumber Testing

This SDK uses Cucumber/Gherkin feature files for cross-SDK test standardization. Feature files are cloned from [algorand-sdk-testing](https://github.com/algorand/algorand-sdk-testing) during test setup.

- `test/unit.tags`: Lists all implemented unit test tags
- `test/integration.tags`: Lists all implemented integration test tags
- `test/features/`: Cucumber feature files (auto-populated by test-harness.sh)
- `test/*.go`: Step definitions that implement the Cucumber tests

When adding support for new test scenarios, add the corresponding tag to the appropriate `.tags` file.

### Client Architecture

Both algod and indexer clients follow a builder pattern:
1. Create client with `MakeClient(address, apiToken)`
2. Call method to create request builder (e.g., `client.Block(round)`)
3. Chain optional parameters as needed
4. Call `.Do(ctx)` to execute the request

Example:
```go
client, _ := algod.MakeClient(address, token)
block, _ := client.Block(round).Do(context.Background())
```

### Transaction Types

All transaction types embed common fields from `types.Header` and include type-specific fields:
- **Payment**: `PaymentTxnFields` (receiver, amount, close remainder)
- **Key Registration**: `KeyregTxnFields` (vote keys, participation)
- **Asset Config**: `AssetConfigTxnFields` (create/reconfigure/destroy assets)
- **Asset Transfer**: `AssetTransferTxnFields` (transfer/opt-in/clawback)
- **Asset Freeze**: `AssetFreezeTxnFields` (freeze/unfreeze asset holdings)
- **Application Call**: `ApplicationFields` (smart contract interactions)
- **State Proof**: `StateProofTxnFields` (consensus proofs)
- **Heartbeat**: `HeartbeatTxnFields` (participation heartbeats)

### Encoding Standards

Algorand uses MessagePack with specific requirements:
- Fields with default/zero values are omitted from encoding
- Field names must be alphabetically ordered for consistent hashing
- The `_struct` tag controls omitempty behavior for all fields

## Development Workflow

### Making Changes

1. Make code changes
2. Run `make fmt` to format code
3. Run `make generate` if modifying logic package
4. Run `make lint` to check for issues
5. Run `make unit` for quick feedback
6. Run `make integration` for full validation (requires harness)
7. Commit changes

### Working with Tests

- Unit tests run quickly without external dependencies
- Integration tests require the test harness (sandbox with algod, kmd, indexer, postgres)
- The test harness can be left running during development for faster iteration
- Cucumber tests are stateful; restart harness if tests behave unexpectedly

### Linter Configuration

The project uses golangci-lint with custom configuration in `.golangci.yml`:
- Generated code in `client/v2/` is excluded from linting
- Specific test utilities in `test/helpers.go` and `test/utilities.go` have relaxed unused rules
- Fields in `types/` package may appear unused but are used via reflection for encoding

## Important Notes

- Go version: 1.23.0+ (see go.mod)
- CI runs on Go 1.23.9 (see .github/workflows/ci-pr.yml)
- Main branch for PRs: `main` (not `master`)
- When adding new transaction types or fields, ensure MessagePack encoding follows canonical ordering
- The `_struct` field with codec tags is essential for proper encoding/decoding
- Test harness uses specific ports and tokens configured in algorand-sdk-testing's .env