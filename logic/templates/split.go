package templates

import "encoding/base64"

type Split struct {
	address string
	program string
}

const referenceProgram = "ASAIAQUCBgcICQomAyCztwQn0+DycN+vsk+vJWcsoz/b7NDS6i33HOkvTpf+YiC3qUpIgHGWE8/1LPh9SGCalSN7IaITeeWSXbfsS5wsXyC4kBQ38Z8zcwWVAym4S8vpFB/c0XC6R4mnPi9EBADsPDEQIhIxASMMEDIEJBJAABoxCSgSMQcyAxIxCCUSMQIhBA0QEBAiQAAzSDMAADMBABIxCTIDEjMABykSMwEHKhIzAAgzAAgzAQgIIQULIQYKEjMACCEHDRAQEBAQEA=="

var referenceOffsets = []uint64{1, 2, 3, 4, 5, 6, 7} // TODO values

func (contract Split) getAddress() string {
	return contract.address
}

func (contract Split) getProgram() string {
	return contract.program
}

func MakeSplit(owner, receiverOne, receiverTwo string, ratn, ratd, expiryRound, minPay, maxFee uint64) (Split, error) {
	referenceAsBytes, err := base64.StdEncoding.DecodeString(referenceProgram)
	if err != nil {
		return Split{}, err
	}
	injectionVector := []interface{}{owner, receiverOne, receiverTwo, ratn, ratd, expiryRound, minPay, maxFee} // TODO ordering
	injectedBytes, err := inject(referenceAsBytes, referenceOffsets, injectionVector)
	if err != nil {
		return Split{}, err
	}
	injectedProgram := base64.StdEncoding.EncodeToString(injectedBytes)
	return Split{}, err
}
