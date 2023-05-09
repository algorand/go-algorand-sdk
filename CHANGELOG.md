# 2.1.0

## What's Changed

Supports new devmode block timestamp offset endpoints.

### Bugfixes
* bugfix: Fix wrong response error type. by @winder in https://github.com/algorand/go-algorand-sdk/pull/461
* debug: Remove debug output. by @winder in https://github.com/algorand/go-algorand-sdk/pull/465
* bug: Fix extractError parsing by @Eric-Warehime in https://github.com/algorand/go-algorand-sdk/pull/492
### Enhancements
* enhancement: add genesis type by @shiqizng in https://github.com/algorand/go-algorand-sdk/pull/443
* docs: Update README.md by @algochoi in https://github.com/algorand/go-algorand-sdk/pull/460
* tests: Add disassembly test for Go SDK by @algochoi in https://github.com/algorand/go-algorand-sdk/pull/462
* marshaling: Lenient Address and BlockHash go-codec unmarshalers. by @winder in https://github.com/algorand/go-algorand-sdk/pull/464
* API: Add types for ledgercore.StateDelta. by @winder in https://github.com/algorand/go-algorand-sdk/pull/467
* docs: Create runnable examples to be pulled into docs by @barnjamin in https://github.com/algorand/go-algorand-sdk/pull/480
* Docs: Examples by @barnjamin in https://github.com/algorand/go-algorand-sdk/pull/491
* api: Regenerate client interfaces for timestamp, ready, and simulate endpoints by @algochoi in https://github.com/algorand/go-algorand-sdk/pull/513
* Performance: Add MakeClientWithTransport client override that allows the user to pass a custom http transport by @pbennett in https://github.com/algorand/go-algorand-sdk/pull/520
* DevOps: Add CODEOWNERS to restrict workflow editing by @onetechnical in https://github.com/algorand/go-algorand-sdk/pull/524
* Performance: Add custom http transport to MakeClientWithTransport  by @algochoi in https://github.com/algorand/go-algorand-sdk/pull/523
### Other
* Regenerate code with the latest specification file (c90fd645) by @github-actions in https://github.com/algorand/go-algorand-sdk/pull/522

## New Contributors
* @pbennett made their first contribution in https://github.com/algorand/go-algorand-sdk/pull/520

**Full Changelog**: https://github.com/algorand/go-algorand-sdk/compare/v2.0.0...v2.1.0

# 2.0.0

## What's Changed
### Breaking Changes
* Remove `future` package.  Move package contents to `transaction`.
* Remove `MakeLogicSigAccount` and replace with `MakeLogicSigAccountEscrow`. Mark `MakeLogicSig` as a private function as well, only intended for internal use.
* Rename `SignLogicsigTransaction` to `SignLogicSigTransaction`.
* Remove logicsig templates, `logic/langspec.json`, and all methods depending on it.
* Remove `DryrunTxnResult.Cost` in favor of 2 fields: `BudgetAdded` and `BudgetConsumed`. `cost` can be derived by `BudgetConsumed - BudgetAdded`.
* Remove v1 algod API (client/algod) due to API end-of-life (2022-12-01). Instead, use v2 algod API (client/v2/algod).
* Remove unused generated types:  `CatchpointAbortResponse`, `CatchpointStartResponse`.

**Full Changelog**: https://github.com/algorand/go-algorand-sdk/compare/v1.24.0...v2.0.0

# 1.24.0
## What's Changed
### Bugfixes
* BugFix: Fix disassemble endpoint by @zyablitsev in https://github.com/algorand/go-algorand-sdk/pull/436
### Enhancements
* Tests: Support for new cucumber app call txn decoding test by @jasonpaulos in https://github.com/algorand/go-algorand-sdk/pull/433
* Tests: Migrate v1 algod dependencies to v2 in cucumber tests by @algochoi in https://github.com/algorand/go-algorand-sdk/pull/434
* REST API:  Add KV counts to NodeStatusResponse by @michaeldiamant in https://github.com/algorand/go-algorand-sdk/pull/437
* Enhancement: allowing zero length static array by @ahangsu in https://github.com/algorand/go-algorand-sdk/pull/438
* Enhancement: revert generic StateProof txn field by @shiqizng in https://github.com/algorand/go-algorand-sdk/pull/439
* Refactoring: Move old transaction dependencies to future.transaction by @algochoi in https://github.com/algorand/go-algorand-sdk/pull/435

## New Contributors
* @zyablitsev made their first contribution in https://github.com/algorand/go-algorand-sdk/pull/436

**Full Changelog**: https://github.com/algorand/go-algorand-sdk/compare/v1.23.0...v1.24.0

# 1.23.0
## What's Changed
### New Features
* Boxes: Add support for Boxes by @michaeldiamant in https://github.com/algorand/go-algorand-sdk/pull/341

**Full Changelog**: https://github.com/algorand/go-algorand-sdk/compare/v1.22.0...v1.23.0

# 1.22.0
## What's Changed
## Enhancements
* REST API: Add algod block hash endpoint, add indexer block header-only param. ([#421](https://github.com/algorand/go-algorand-sdk/pull/421))

# 1.21.0
## What's Changed
### Enhancements
* Deprecation: Add deprecated tags to v1 `algod1` API ([#392](https://github.com/algorand/go-algorand-sdk/pull/392))
* Enhancement: update block model ([#401](https://github.com/algorand/go-algorand-sdk/pull/401))
### Bugfixes
* Bugfix: Fix dryrun parser ([#400](https://github.com/algorand/go-algorand-sdk/pull/400))

# 1.20.0
## What's Changed
### Bugfixes
* Bug-Fix: passthru verbosity by @tzaffi in https://github.com/algorand/go-algorand-sdk/pull/371
* BugFix: Src map type assert fix by @barnjamin in https://github.com/algorand/go-algorand-sdk/pull/370
### New Features
* StateProof: State proof support by @Eric-Warehime in https://github.com/algorand/go-algorand-sdk/pull/374
* StateProof: State Proof Verification additions by @almog-t in https://github.com/algorand/go-algorand-sdk/pull/377
* State Proofs: added compute leaf function for light block header to sdk by @almog-t in https://github.com/algorand/go-algorand-sdk/pull/382
* State Proofs: renamed light block header hash func by @almog-t in https://github.com/algorand/go-algorand-sdk/pull/383
### Enhancements
* Enhancement: Use Sandbox for Testing by @tzaffi in https://github.com/algorand/go-algorand-sdk/pull/360
* Enhancement: Deprecating use of langspec by @ahangsu in https://github.com/algorand/go-algorand-sdk/pull/366
* State Proofs: Use generic type for StateProof txn field. by @winder in https://github.com/algorand/go-algorand-sdk/pull/378
* Improvement: Better SourceMap decoding by @barnjamin in https://github.com/algorand/go-algorand-sdk/pull/380
* tests: Enable stpf cucumber unit tests by @Eric-Warehime in https://github.com/algorand/go-algorand-sdk/pull/386

## New Contributors
* @tzaffi made their first contribution in https://github.com/algorand/go-algorand-sdk/pull/360
* @almog-t made their first contribution in https://github.com/algorand/go-algorand-sdk/pull/377

**Full Changelog**: https://github.com/algorand/go-algorand-sdk/compare/v1.19.0...v1.20.0

# 1.19.0
## Enhancements
* AVM: Consolidate TEAL and AVM versions ([#345](https://github.com/algorand/go-algorand-sdk/pull/345))
* Testing: Use Dev mode network for cucumber tests ([#349](https://github.com/algorand/go-algorand-sdk/pull/349))
* AVM: Use avm-abi repo ([#352](https://github.com/algorand/go-algorand-sdk/pull/352))

# 1.18.0

## What's Changed

### New Features
* Dev Tools: Source map decoder by @barnjamin in https://github.com/algorand/go-algorand-sdk/pull/335

### Enhancements
* Github-Actions: Adding pr title and label checks by @algojack in https://github.com/algorand/go-algorand-sdk/pull/336

**Full Changelog**: https://github.com/algorand/go-algorand-sdk/compare/v1.17.0...v1.18.0


# 1.17.0

## What's Changed
* Added GetMethodByName on Interface and Contract ([#330](https://github.com/algorand/go-algorand-sdk/pull/330))
* Regenerated code with the latest specification file (d012c9f5) ([#332](https://github.com/algorand/go-algorand-sdk/pull/332))
* Added helper method for formatting the algod API path ([#331](https://github.com/algorand/go-algorand-sdk/pull/331))
* Added method in ABI results object ([#329](https://github.com/algorand/go-algorand-sdk/pull/329))

# 1.16.0

## Important Note
This release includes an upgrade to golang 1.17.

## What's Changed
* Adding `Foreign*` args to AddMethodCallParams by @barnjamin in https://github.com/algorand/go-algorand-sdk/pull/318
* build: Bump golang to 1.17 by @Eric-Warehime in https://github.com/algorand/go-algorand-sdk/pull/314
* Update generated files by @Eric-Warehime in https://github.com/algorand/go-algorand-sdk/pull/321
* Copy foreign arrays before modifying by @algoidurovic in https://github.com/algorand/go-algorand-sdk/pull/323
* Build: Sdk code generation automation by @Eric-Warehime in https://github.com/algorand/go-algorand-sdk/pull/324
* Update codegen.yml by @Eric-Warehime in https://github.com/algorand/go-algorand-sdk/pull/325
* Generate updated API client code by @algoidurovic in https://github.com/algorand/go-algorand-sdk/pull/316

## New Contributors
* @Eric-Warehime made their first contribution in https://github.com/algorand/go-algorand-sdk/pull/314

# 1.15.0
* adding foreign app addr to dryrun creator ([#312](https://github.com/algorand/go-algorand-sdk/pull/312))
* adding dryrun stack printer ([#289](https://github.com/algorand/go-algorand-sdk/pull/289))
* Readme updates ([#296](https://github.com/algorand/go-algorand-sdk/pull/296))
# 1.14.1
- Avoid client response failure on unknown field (#307)
- Add ParticipationUpdates to BlockHeader (#306)
# 1.14.0
- Unlimited assets changes (#294)
- Update abi impl from go-algorand (#303)
- Update go to version 1.16 (#301)
# 1.14.0-beta.1
- Unlimited assets changes (#294)
# 1.13.0
## Added
- Add app creator to dryrun request (#283)
- Stateproof keys APIs changes (#284)
- adding status code checker to msgpack decode (#286)
## Changed:
- Update to use v2 client (#270)
- Implement C2C tests (#282)
- Add circleci job (#287)
- Update langspec for TEAL 6 (#291)
# 1.12.0
## Added
- Add stateproof to keyreg transaction (#278)
- Create response object for "AtomicTransactionComposer.Execute" (#276)
- Support ABI reference types and other improvements (#273)
- Add CreateDryrun function (#265)
- Add EncodeAddress (#264)
- ABI Interaction (#258)
- Add ABI-encoding feature (#247)
- Implemented WaitForConfirmation function (#232)
## Changed
- Update abi exported interface (#255)
- Update ApplyData and EvalDelta (#249)
# 1.12.0-beta.2
## Added
- Support ABI reference types and other improvements (#273)
## Changed
- Fix wait for confirmation function (#267)
# 1.12.0-beta.1
## Added
- EncodeAddress
- ABI Interaction
- WaitForConfirmation function
## Changed
- ABI Interface
# 1.11.0
## Added
- add TealVerify function (#242)
- Support AVM 1.0 (#248)
- Test with go 1.17 in ci (#237)
## Changed
- Mark contract binary template code as Deprecated (#241)
# 1.10.0
## Added
- New github Issue template
- Signing support for rekeying to LogicSig/MultiSig account
- Asset Base64 Fields
## BugFix
- Use correct go version in CI
# 1.9.2
## Bug Fix
- Update FromBase64String() to correctly return the signed transaction
- Make MakeApplicationCreateTxWithExtraPages() and revert MakeApplicationCallTx() to make non-API-breaking
# 1.9.1
## Bugfix
- Allow asset URLs to be up to 96 bytes
# 1.9.0 - API COMPATIBILITY CHANGE
## Added
- Support for TEAL 4 programs
- Support for creating application with extra program pages
- Support for setting a transaction fee below the network minimum, for use with fee pooling
## Bugfix
- Algod and Indexer responses will now produce JSON that matches the other SDKs
# 1.8.0
## Added
- V2: Add MakeClientWithHeaders wrapper functions
## Bugfix
- Fix FlatFee computation
# 1.7.0
## Bugfix
- Fix GetGenesis endpoint.
# 1.6.0
## Added
- Code generation for more of the http client
- Add TEAL 3 support
- template UX tweaks
## Bugfix
- Make limitorder.GetSwapAssetsTransaction behave the same as other SDKs
# 1.5.1
## Added
- Add `BlockRaw` method to algod API V2 client.
# 1.5.0
## Added
- Support for Applications
# 1.4.2
## Bugfix
- Fix incorrect `SendRawTransaction` path in API V2 client.
# 1.4.1
## Bugfix
-  Fix go get, test package names needed to be renamed.
# 1.4.0
## Added
-  Clients for Indexer V2 and algod API V2
# 1.3.0
## Added
-  additional Algorand Smart Contracts (ASC)
    -  support for Dynamic Fee contract
    -  support for Limit Order contract
    -  support for Periodic Payment contract
- support for SuggestedParams
- support for RawBlock request
- Missing transaction types
# 1.2.1
## Added
- Added asset decimals field.
# 1.2.0
## Added
- Added support for Algorand Standardized Assets (ASA)
- Added support for Algorand Smart Contracts (ASC)
    - Added support for Hashed Time Lock Contract (HTLC)
    - Added support for Split contract
- Added support for Group Transactions
- Added support for leases
# 1.1.3
## Added
- Signing and verifying arbitrary bytes
- Deleting multisigs
- Support for flat fees in transactions
## Changed
- Add note parameter to key registration transaction constructors
# 1.1.2
## Added
- Support for GenesisHash
- Updated API Models.
# 1.1.1
## Added
- Indexer support
# 1.1.0
## Added
- Multisignature support
# 1.0.6
## Added
- Support in new SuggestedFee functionality
# 1.0.5
## Added
- Added helper functions for preventing overflow
# 1.0.4
## Added
- Better auction support
# 1.0.3
## Added
- Additional mnemonic support
# 1.0.2
## Added
- Support for "genesis ID" field in transactions
- Support for "close remainder to" field in transactions
# 1.0.0
## Added
- SDK released
