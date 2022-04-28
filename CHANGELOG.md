# 1.15.0
- Update submodule to 3.5.1 [929](https://github.com/algorand/indexer/pull/929)
- Fix dev mode by addressing off by one error. [920](https://github.com/algorand/indexer/pull/920)
- Enable validation and acceptance of new parameter configuration [912](https://github.com/algorand/indexer/pull/912)
- Revert imported_tx_per_block change and add new imported_txns gauge. [913](https://github.com/algorand/indexer/pull/913)
- Update e2e_subs docs to revise artifact upload process [910](https://github.com/algorand/indexer/pull/910)
- Document instructions for updating indexer e2e test input [906](https://github.com/algorand/indexer/pull/906)
- Allow viewing of disabled params [902](https://github.com/algorand/indexer/pull/902)
- Fix failing e2e test [899](https://github.com/algorand/indexer/pull/899)
- Monitoring dashboard [876](https://github.com/algorand/indexer/pull/876)
- Disable parameters in rest api [892](https://github.com/algorand/indexer/pull/892)
- Ensure query is canceled after account rewind. [893](https://github.com/algorand/indexer/pull/893)
- Temporarily revert to inline glob search [889](https://github.com/algorand/indexer/pull/889)
- Configurable query parameters runtime data structures. [873](https://github.com/algorand/indexer/pull/873)
- BugfixReturn all inner transactions are returned for logs endpoint. [915](https://github.com/algorand/indexer/pull/915)
- Default to including python goAlgorand e2e tests in buildtestdata.sh [888](https://github.com/algorand/indexer/pull/888)
- Do not mergeIndexer 2.10.0 [923](https://github.com/algorand/indexer/pull/923)
- FeatureAdd support for unlimited assets. [900](https://github.com/algorand/indexer/pull/900)
- Include block-generator and validator as algorandIndexer subcommands. [891](https://github.com/algorand/indexer/pull/891)
- Release preparationFeature flag to disable configurable api parameters. [917](https://github.com/algorand/indexer/pull/917)
- ToolsAdd cli doc generation command. [919](https://github.com/algorand/indexer/pull/919)
- Use fetch/reset instead of pull for goAlgorand in nightly test. [885](https://github.com/algorand/indexer/pull/885)
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
