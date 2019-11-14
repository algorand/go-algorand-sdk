package templates

import (
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
	require.NoError(t, err)
	// Outputs
	goldenProgram := "ASAIAcCWsQICAMDEBx5kkE4mAyCztwQn0+DycN+vsk+vJWcsoz/b7NDS6i33HOkvTpf+YiC3qUpIgHGWE8/1LPh9SGCalSN7IaITeeWSXbfsS5wsXyC4kBQ38Z8zcwWVAym4S8vpFB/c0XC6R4mnPi9EBADsPDEQIhIxASMMEDIEJBJAABkxCSgSMQcyAxIQMQglEhAxAiEEDRAiQAAuMwAAMwEAEjEJMgMSEDMABykSEDMBByoSEDMACCEFCzMBCCEGCxIQMwAIIQcPEBA="
	require.Equal(t, goldenProgram, c.GetProgram())
	goldenAddress := "KPYGWKTV7CKMPMTLQRNGMEQRSYTYDHUOFNV4UDSBDLC44CLIJPQWRTCPBU"
	require.Equal(t, goldenAddress, c.GetAddress())
}

func TestDelegatedFee(t *testing.T) {
	require.Equal(t, 0, 1)
}

func TestPeriodicPayment(t *testing.T) {
	require.Equal(t, 0, 1)
}

func TestLimitOrder(t *testing.T) {
	require.Equal(t, 0, 1)
}
