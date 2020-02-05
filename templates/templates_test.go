package templates

import (
	"encoding/base64"
	"testing"

	"github.com/algorand/go-algorand-sdk/encoding/msgpack"

	"github.com/stretchr/testify/require"
)

func TestSplit(t *testing.T) {
	// Inputs
	owner := "WO3QIJ6T4DZHBX5PWJH26JLHFSRT7W7M2DJOULPXDTUS6TUX7ZRIO4KDFY"
	receivers := [2]string{"W6UUUSEAOGLBHT7VFT4H2SDATKKSG6ZBUIJXTZMSLW36YS44FRP5NVAU7U", "XCIBIN7RT4ZXGBMVAMU3QS6L5EKB7XGROC5EPCNHHYXUIBAA5Q6C5Y7NEU"}
	ratn, ratd := uint64(30), uint64(100) // receiverOne gets 30/100 of whatever is sent to contract address
	expiryRound := uint64(123456)
	minPay := uint64(10000)
	maxFee := uint64(5000000)
	c, err := MakeSplit(owner, receivers[0], receivers[1], ratn, ratd, expiryRound, minPay, maxFee)
	// Outputs
	require.NoError(t, err)
	goldenProgram := "ASAIAcCWsQICAMDEB2QekE4mAyCztwQn0+DycN+vsk+vJWcsoz/b7NDS6i33HOkvTpf+YiC3qUpIgHGWE8/1LPh9SGCalSN7IaITeeWSXbfsS5wsXyC4kBQ38Z8zcwWVAym4S8vpFB/c0XC6R4mnPi9EBADsPDEQIhIxASMMEDIEJBJAABkxCSgSMQcyAxIQMQglEhAxAiEEDRAiQAAuMwAAMwEAEjEJMgMSEDMABykSEDMBByoSEDMACCEFCzMBCCEGCxIQMwAIIQcPEBA="
	require.Equal(t, goldenProgram, base64.StdEncoding.EncodeToString(c.GetProgram()))
	goldenAddress := "HDY7A4VHBWQWQZJBEMASFOUZKBNGWBMJEMUXAGZ4SPIRQ6C24MJHUZKFGY"
	require.Equal(t, goldenAddress, c.GetAddress())
	goldenGenesisHash := "f4OxZX/x/FO5LcGBSKHWXfwtSx+j1ncoSt3SABJtkGk="
	genesisBytes, _ := base64.StdEncoding.DecodeString(goldenGenesisHash)
	goldenStx := "gqRsc2lngaFsxM4BIAgBwJaxAgIAwMQHZB6QTiYDILO3BCfT4PJw36+yT68lZyyjP9vs0NLqLfcc6S9Ol/5iILepSkiAcZYTz/Us+H1IYJqVI3shohN55ZJdt+xLnCxfILiQFDfxnzNzBZUDKbhLy+kUH9zRcLpHiac+L0QEAOw8MRAiEjEBIwwQMgQkEkAAGTEJKBIxBzIDEhAxCCUSEDECIQQNECJAAC4zAAAzAQASMQkyAxIQMwAHKRIQMwEHKhIQMwAIIQULMwEIIQYLEhAzAAghBw8QEKN0eG6Jo2FtdM4AB9M0o2ZlZc4AId/gomZ2AaJnaMQgf4OxZX/x/FO5LcGBSKHWXfwtSx+j1ncoSt3SABJtkGmjZ3JwxCBAATLqJB2EUR98mHEDe8EVRvziz3fPALAaBRe6xhOqZKJsdmSjcmN2xCC3qUpIgHGWE8/1LPh9SGCalSN7IaITeeWSXbfsS5wsX6NzbmTEIDjx8HKnDaFoZSEjASK6mVBaawWJIylwGzyT0Rh4WuMSpHR5cGWjcGF5gqRsc2lngaFsxM4BIAgBwJaxAgIAwMQHZB6QTiYDILO3BCfT4PJw36+yT68lZyyjP9vs0NLqLfcc6S9Ol/5iILepSkiAcZYTz/Us+H1IYJqVI3shohN55ZJdt+xLnCxfILiQFDfxnzNzBZUDKbhLy+kUH9zRcLpHiac+L0QEAOw8MRAiEjEBIwwQMgQkEkAAGTEJKBIxBzIDEhAxCCUSEDECIQQNECJAAC4zAAAzAQASMQkyAxIQMwAHKRIQMwEHKhIQMwAIIQULMwEIIQYLEhAzAAghBw8QEKN0eG6Jo2FtdM4AAlj2o2ZlZc4AId/gomZ2AaJnaMQgf4OxZX/x/FO5LcGBSKHWXfwtSx+j1ncoSt3SABJtkGmjZ3JwxCBAATLqJB2EUR98mHEDe8EVRvziz3fPALAaBRe6xhOqZKJsdmSjcmN2xCC4kBQ38Z8zcwWVAym4S8vpFB/c0XC6R4mnPi9EBADsPKNzbmTEIDjx8HKnDaFoZSEjASK6mVBaawWJIylwGzyT0Rh4WuMSpHR5cGWjcGF5"
	stx, err := GetSplitFundsTransaction(c.GetProgram(), 2*100000*ratd/ratn, 1, 100, 10000, genesisBytes)
	require.NoError(t, err)
	require.Equal(t, goldenStx, base64.StdEncoding.EncodeToString(stx))
}

func TestHTLC(t *testing.T) {
	// Inputs
	owner := "726KBOYUJJNE5J5UHCSGQGWIBZWKCBN4WYD7YVSTEXEVNFPWUIJ7TAEOPM"
	receiver := "42NJMHTPFVPXVSDGA6JGKUV6TARV5UZTMPFIREMLXHETRKIVW34QFSDFRE"
	hashFn := "sha256"
	hashImg := "f4OxZX/x/FO5LcGBSKHWXfwtSx+j1ncoSt3SABJtkGk="
	expiryRound := uint64(600000)
	maxFee := uint64(1000)
	c, err := MakeHTLC(owner, receiver, hashFn, hashImg, expiryRound, maxFee)
	// Outputs
	require.NoError(t, err)
	goldenProgram := "ASAE6AcBAMDPJCYDIOaalh5vLV96yGYHkmVSvpgjXtMzY8qIkYu5yTipFbb5IH+DsWV/8fxTuS3BgUih1l38LUsfo9Z3KErd0gASbZBpIP68oLsUSlpOp7Q4pGgayA5soQW8tgf8VlMlyVaV9qITMQEiDjEQIxIQMQcyAxIQMQgkEhAxCSgSLQEpEhAxCSoSMQIlDRAREA=="
	require.Equal(t, goldenProgram, base64.StdEncoding.EncodeToString(c.GetProgram()))
	goldenAddress := "KNBD7ATNUVQ4NTLOI72EEUWBVMBNKMPHWVBCETERV2W7T2YO6CVMLJRBM4"
	require.Equal(t, goldenAddress, c.GetAddress())
}

func TestPeriodicPayment(t *testing.T) {
	// Inputs
	receiver := "SKXZDBHECM6AS73GVPGJHMIRDMJKEAN5TUGMUPSKJCQ44E6M6TC2H2UJ3I"
	artificialLease := "AQIDBAUGBwgBAgMEBQYHCAECAwQFBgcIAQIDBAUGBwg="
	amount := uint64(500000)
	withdrawalWindow := uint64(95)
	period := uint64(100)
	maxFee := uint64(1000)
	expiryRound := uint64(2445756)
	c, err := makePeriodicPaymentWithLease(receiver, artificialLease, amount, withdrawalWindow, period, expiryRound, maxFee)
	// Outputs
	require.NoError(t, err)
	goldenProgram := "ASAHAegHZABfoMIevKOVASYCIAECAwQFBgcIAQIDBAUGBwgBAgMEBQYHCAECAwQFBgcIIJKvkYTkEzwJf2arzJOxERsSogG9nQzKPkpIoc4TzPTFMRAiEjEBIw4QMQIkGCUSEDEEIQQxAggSEDEGKBIQMQkyAxIxBykSEDEIIQUSEDEJKRIxBzIDEhAxAiEGDRAxCCUSEBEQ"
	contractBytes := c.GetProgram()
	require.Equal(t, goldenProgram, base64.StdEncoding.EncodeToString(contractBytes))
	goldenAddress := "JMS3K4LSHPULANJIVQBTEDP5PZK6HHMDQS4OKHIMHUZZ6OILYO3FVQW7IY"
	require.Equal(t, goldenAddress, c.GetAddress())
	goldenGenesisHash := "f4OxZX/x/FO5LcGBSKHWXfwtSx+j1ncoSt3SABJtkGk="
	genesisBytes, err := base64.StdEncoding.DecodeString(goldenGenesisHash)
	require.NoError(t, err)
	stx, err := GetPeriodicPaymentWithdrawalTransaction(contractBytes, 1200, maxFee, genesisBytes)
	require.NoError(t, err)
	goldenStx := "gqRsc2lngaFsxJkBIAcB6AdkAF+gwh68o5UBJgIgAQIDBAUGBwgBAgMEBQYHCAECAwQFBgcIAQIDBAUGBwggkq+RhOQTPAl/ZqvMk7ERGxKiAb2dDMo+SkihzhPM9MUxECISMQEjDhAxAiQYJRIQMQQhBDECCBIQMQYoEhAxCTIDEjEHKRIQMQghBRIQMQkpEjEHMgMSEDECIQYNEDEIJRIQERCjdHhuiaNhbXTOAAehIKNmZWXOAAQDWKJmds0EsKJnaMQgf4OxZX/x/FO5LcGBSKHWXfwtSx+j1ncoSt3SABJtkGmibHbNBQ+ibHjEIAECAwQFBgcIAQIDBAUGBwgBAgMEBQYHCAECAwQFBgcIo3JjdsQgkq+RhOQTPAl/ZqvMk7ERGxKiAb2dDMo+SkihzhPM9MWjc25kxCBLJbVxcjvosDUorAMyDf1+VeOdg4S45R0MPTOfOQvDtqR0eXBlo3BheQ=="
	require.Equal(t, goldenStx, base64.StdEncoding.EncodeToString(stx))
}

func TestDynamicFee(t *testing.T) {
	// Inputs
	receiver := "726KBOYUJJNE5J5UHCSGQGWIBZWKCBN4WYD7YVSTEXEVNFPWUIJ7TAEOPM"
	amount := uint64(5000)
	firstValid := uint64(12345)
	lastValid := uint64(12346)
	closeRemainder := "42NJMHTPFVPXVSDGA6JGKUV6TARV5UZTMPFIREMLXHETRKIVW34QFSDFRE"
	artificialLease := "f4OxZX/x/FO5LcGBSKHWXfwtSx+j1ncoSt3SABJtkGk="
	c, err := makeDynamicFeeWithLease(receiver, closeRemainder, artificialLease, amount, firstValid, lastValid)
	require.NoError(t, err)
	goldenGenesisHash := "f4OxZX/x/FO5LcGBSKHWXfwtSx+j1ncoSt3SABJtkGk="
	genesisBytes, err := base64.StdEncoding.DecodeString(goldenGenesisHash)
	require.NoError(t, err)
	contractBytes := c.GetProgram()
	require.NoError(t, err)
	txn, lsig, err := SignDynamicFee(contractBytes, genesisBytes)
	require.NoError(t, err)
	goldenLsig := "gaFsxLEBIAUCAYgnuWC6YCYDIP68oLsUSlpOp7Q4pGgayA5soQW8tgf8VlMlyVaV9qITIOaalh5vLV96yGYHkmVSvpgjXtMzY8qIkYu5yTipFbb5IH+DsWV/8fxTuS3BgUih1l38LUsfo9Z3KErd0gASbZBpMgQiEjMAECMSEDMABzEAEhAzAAgxARIQMRYjEhAxECMSEDEHKBIQMQkpEhAxCCQSEDECJRIQMQQhBBIQMQYqEhA="
	require.Equal(t, goldenLsig, base64.StdEncoding.EncodeToString(msgpack.Encode(lsig)))
	goldenTxn := "iqNhbXTNE4ilY2xvc2XEIOaalh5vLV96yGYHkmVSvpgjXtMzY8qIkYu5yTipFbb5o2ZlZc0D6KJmds0wOaJnaMQgf4OxZX/x/FO5LcGBSKHWXfwtSx+j1ncoSt3SABJtkGmibHbNMDqibHjEIH+DsWV/8fxTuS3BgUih1l38LUsfo9Z3KErd0gASbZBpo3JjdsQg/ryguxRKWk6ntDikaBrIDmyhBby2B/xWUyXJVpX2ohOjc25kxCAwkctYaLUKCdXuhv+nZwbuiUW+mVI/jfGXVvKVIDm5x6R0eXBlo3BheQ=="
	require.Equal(t, goldenTxn, base64.StdEncoding.EncodeToString(msgpack.Encode(txn)))
	privateKeyB64 := "2qjz96Vj9M6YOqtNlfJUOKac13EHCXyDty94ozCjuwwriI+jzFgStFx9E6kEk1l4+lFsW4Te2PY1KV8kNcccRg=="
	privateKey, err := base64.StdEncoding.DecodeString(privateKeyB64)
	stxns, err := GetDynamicFeeTransactions(txn, lsig, privateKey, 1234)
	require.NoError(t, err)
	// Outputs
	goldenProgram := "ASAFAgGIJ7lgumAmAyD+vKC7FEpaTqe0OKRoGsgObKEFvLYH/FZTJclWlfaiEyDmmpYeby1feshmB5JlUr6YI17TM2PKiJGLuck4qRW2+SB/g7Flf/H8U7ktwYFIodZd/C1LH6PWdyhK3dIAEm2QaTIEIhIzABAjEhAzAAcxABIQMwAIMQESEDEWIxIQMRAjEhAxBygSEDEJKRIQMQgkEhAxAiUSEDEEIQQSEDEGKhIQ"
	require.Equal(t, goldenProgram, base64.StdEncoding.EncodeToString(contractBytes))
	goldenAddress := "GCI4WWDIWUFATVPOQ372OZYG52EULPUZKI7Y34MXK3ZJKIBZXHD2H5C5TI"
	require.Equal(t, goldenAddress, c.GetAddress())
	goldenStxns := "gqNzaWfEQGeQXRB3Jud/EGWv+cvIJa+xGdz79sA8xf76fCI16KX1YbUUMXbu2lXksXXMsejkw05sYawl7o0VkstARzxu9gajdHhuiqNhbXTOAAWq6qNmZWXOAATzvqJmds0wOaJnaMQgf4OxZX/x/FO5LcGBSKHWXfwtSx+j1ncoSt3SABJtkGmjZ3JwxCCg8Vspk8YbZ5yeSiz//h8ExJ0NundK91p6hOILxCKehaJsds0wOqJseMQgf4OxZX/x/FO5LcGBSKHWXfwtSx+j1ncoSt3SABJtkGmjcmN2xCAwkctYaLUKCdXuhv+nZwbuiUW+mVI/jfGXVvKVIDm5x6NzbmTEICuIj6PMWBK0XH0TqQSTWXj6UWxbhN7Y9jUpXyQ1xxxGpHR5cGWjcGF5gqRsc2lngaFsxLEBIAUCAYgnuWC6YCYDIP68oLsUSlpOp7Q4pGgayA5soQW8tgf8VlMlyVaV9qITIOaalh5vLV96yGYHkmVSvpgjXtMzY8qIkYu5yTipFbb5IH+DsWV/8fxTuS3BgUih1l38LUsfo9Z3KErd0gASbZBpMgQiEjMAECMSEDMABzEAEhAzAAgxARIQMRYjEhAxECMSEDEHKBIQMQkpEhAxCCQSEDECJRIQMQQhBBIQMQYqEhCjdHhui6NhbXTNE4ilY2xvc2XEIOaalh5vLV96yGYHkmVSvpgjXtMzY8qIkYu5yTipFbb5o2ZlZc4ABarqomZ2zTA5omdoxCB/g7Flf/H8U7ktwYFIodZd/C1LH6PWdyhK3dIAEm2QaaNncnDEIKDxWymTxhtnnJ5KLP/+HwTEnQ26d0r3WnqE4gvEIp6Fomx2zTA6omx4xCB/g7Flf/H8U7ktwYFIodZd/C1LH6PWdyhK3dIAEm2QaaNyY3bEIP68oLsUSlpOp7Q4pGgayA5soQW8tgf8VlMlyVaV9qITo3NuZMQgMJHLWGi1CgnV7ob/p2cG7olFvplSP43xl1bylSA5ucekdHlwZaNwYXk="
	require.Equal(t, goldenStxns, base64.StdEncoding.EncodeToString(stxns))
}

func TestLimitOrder(t *testing.T) {
	// Inputs
	owner := "726KBOYUJJNE5J5UHCSGQGWIBZWKCBN4WYD7YVSTEXEVNFPWUIJ7TAEOPM"
	assetid := uint64(12345)
	ratn, ratd := uint64(30), uint64(100)
	expiryRound := uint64(123456)
	minTrade := uint64(10000)
	maxFee := uint64(5000000)
	c, err := MakeLimitOrder(owner, assetid, ratn, ratd, expiryRound, minTrade, maxFee)
	// Outputs
	require.NoError(t, err)
	goldenProgram := "ASAKAAHAlrECApBOBLlgZB7AxAcmASD+vKC7FEpaTqe0OKRoGsgObKEFvLYH/FZTJclWlfaiEzEWIhIxECMSEDEBJA4QMgQjEkAAVTIEJRIxCCEEDRAxCTIDEhAzARAhBRIQMwERIQYSEDMBFCgSEDMBEzIDEhAzARIhBx01AjUBMQghCB01BDUDNAE0Aw1AACQ0ATQDEjQCNAQPEEAAFgAxCSgSMQIhCQ0QMQcyAxIQMQgiEhAQ"
	require.Equal(t, goldenProgram, base64.StdEncoding.EncodeToString(c.GetProgram()))
	goldenAddress := "LXQWT2XLIVNFS54VTLR63UY5K6AMIEWI7YTVE6LB4RWZDBZKH22ZO3S36I"
	require.Equal(t, goldenAddress, c.GetAddress())
}
