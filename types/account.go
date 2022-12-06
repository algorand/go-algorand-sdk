package types

// AccountData contains the data associated with a given address.
// This includes the account balance, cryptographic public keys,
// consensus delegation status, asset data, and application data.
type Account struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	Status             string                      `codec:"onl"`
	MicroAlgos         uint64                      `codec:"algo"`
	RewardsBase        uint64                      `codec:"ebase"`
	RewardedMicroAlgos uint64                      `codec:"ern"`
	VoteID             [32]byte                    `codec:"vote"`
	SelectionID        [32]byte                    `codec:"sel"`
	StateProofID       [64]byte                    `codec:"stprf"`
	VoteFirstValid     uint64                      `codec:"voteFst"`
	VoteLastValid      uint64                      `codec:"voteLst"`
	VoteKeyDilution    uint64                      `codec:"voteKD"`
	AssetParams        map[AssetIndex]AssetParams  `codec:"apar,allocbound=encodedMaxAssetsPerAccount"`
	Assets             map[AssetIndex]AssetHolding `codec:"asset,allocbound=encodedMaxAssetsPerAccount"`
	AuthAddr           Address                     `codec:"spend"`
	AppLocalStates     map[AppIndex]AppLocalState  `codec:"appl,allocbound=EncodedMaxAppLocalStates"`
	AppParams          map[AppIndex]AppParams      `codec:"appp,allocbound=EncodedMaxAppParams"`
	TotalAppSchema     StateSchema                 `codec:"tsch"`
	TotalExtraAppPages uint32                      `codec:"teap"`
	TotalBoxes         uint64                      `codec:"tbx"`
	TotalBoxBytes      uint64                      `codec:"tbxb"`
}
