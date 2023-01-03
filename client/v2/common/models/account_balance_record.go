package models

// AccountBalanceRecord account and its address
type AccountBalanceRecord struct {
	// AccountData updated account data.
	AccountData Account `json:"account-data"`

	// Address address of the updated account.
	Address string `json:"address"`
}
