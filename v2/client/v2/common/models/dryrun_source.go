package models

// DryrunSource dryrunSource is TEAL source text that gets uploaded, compiled, and
// inserted into transactions or application state.
type DryrunSource struct {
	// AppIndex
	AppIndex uint64 `json:"app-index"`

	// FieldName fieldName is what kind of sources this is. If lsig then it goes into
	// the transactions[this.TxnIndex].LogicSig. If approv or clearp it goes into the
	// Approval Program or Clear State Program of application[this.AppIndex].
	FieldName string `json:"field-name"`

	// Source
	Source string `json:"source"`

	// TxnIndex
	TxnIndex uint64 `json:"txn-index"`
}
