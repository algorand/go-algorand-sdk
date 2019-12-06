package templates

import (
	"encoding/base64"
	"github.com/stretchr/testify/require"
	"testing"
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
	goldenProgram := "ASAIAcCWsQICAMDEBx5kkE4mAyCztwQn0+DycN+vsk+vJWcsoz/b7NDS6i33HOkvTpf+YiC3qUpIgHGWE8/1LPh9SGCalSN7IaITeeWSXbfsS5wsXyC4kBQ38Z8zcwWVAym4S8vpFB/c0XC6R4mnPi9EBADsPDEQIhIxASMMEDIEJBJAABkxCSgSMQcyAxIQMQglEhAxAiEEDRAiQAAuMwAAMwEAEjEJMgMSEDMABykSEDMBByoSEDMACCEFCzMBCCEGCxIQMwAIIQcPEBA="
	require.Equal(t, goldenProgram, base64.StdEncoding.EncodeToString(c.GetProgram()))
	goldenAddress := "KPYGWKTV7CKMPMTLQRNGMEQRSYTYDHUOFNV4UDSBDLC44CLIJPQWRTCPBU"
	require.Equal(t, goldenAddress, c.GetAddress())
}

func TestHTLC(t *testing.T) {
	// Inputs
	owner := "ç"
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
