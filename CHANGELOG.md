# 1.14.0-beta.1
## Added:
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
