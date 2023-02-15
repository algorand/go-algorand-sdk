package mnemonic

import (
	"github.com/algorand/go-algorand-sdk/v2/test"
	"github.com/stretchr/testify/require"
)

var expected = `
{
  "Accts": {
    "Accts": [
      {
        "Addr": "7777777777777777777777777777777777777777777777777774MSJUVU",
        "AuthAddr": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAY5HFKQ",
        "MicroAlgos": 2051818321,
        "RewardedMicroAlgos": 0,
        "RewardsBase": 0,
        "SelectionID": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
        "StateProofID": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==",
        "Status": 2,
        "TotalAppLocalStates": 0,
        "TotalAppParams": 0,
        "TotalAppSchema": {},
        "TotalAssetParams": 0,
        "TotalAssets": 0,
        "TotalBoxBytes": 0,
        "TotalBoxes": 0,
        "TotalExtraAppPages": 0,
        "VoteFirstValid": 0,
        "VoteID": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
        "VoteKeyDilution": 0,
        "VoteLastValid": 0
      },
      {
        "Addr": "UKSBQH6FOHPZMEQ7YPDSD3AV54XSXRFG744W6WF4QUFEF5SLXVIL7SQ4GM",
        "AuthAddr": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAY5HFKQ",
        "MicroAlgos": 295200,
        "RewardedMicroAlgos": 0,
        "RewardsBase": 27521,
        "SelectionID": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
        "StateProofID": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==",
        "Status": 0,
        "TotalAppLocalStates": 0,
        "TotalAppParams": 0,
        "TotalAppSchema": {},
        "TotalAssetParams": 0,
        "TotalAssets": 0,
        "TotalBoxBytes": 0,
        "TotalBoxes": 0,
        "TotalExtraAppPages": 0,
        "VoteFirstValid": 0,
        "VoteID": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
        "VoteKeyDilution": 0,
        "VoteLastValid": 0
      },
      {
        "Addr": "A7NMWS3NT3IUDMLVO26ULGXGIIOUQ3ND2TXSER6EBGRZNOBOUIQXHIBGDE",
        "AuthAddr": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAY5HFKQ",
        "MicroAlgos": 342106549750,
        "RewardedMicroAlgos": 0,
        "RewardsBase": 0,
        "SelectionID": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
        "StateProofID": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==",
        "Status": 2,
        "TotalAppLocalStates": 0,
        "TotalAppParams": 0,
        "TotalAppSchema": {},
        "TotalAssetParams": 0,
        "TotalAssets": 0,
        "TotalBoxBytes": 0,
        "TotalBoxes": 0,
        "TotalExtraAppPages": 0,
        "VoteFirstValid": 0,
        "VoteID": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
        "VoteKeyDilution": 0,
        "VoteLastValid": 0
      }
    ],
    "AppResources": null,
    "AssetResources": null
  },
  "Creatables": null,
  "Hdr": {
    "earn": 27521,
    "fees": "A7NMWS3NT3IUDMLVO26ULGXGIIOUQ3ND2TXSER6EBGRZNOBOUIQXHIBGDE",
    "frac": 2047917471,
    "gen": "testnet-v1.0",
    "gh": "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI=",
    "prev": "blk-5TF3YYGCEP3F5Y67LS4QND6GO2J3LROT3JYVCJBULDWQC3KIPATA",
    "proto": "https://github.com/algorandfoundation/specs/tree/44fa607d6051730f5264526bf3c108d51f0eadb6",
    "rate": 42,
    "rnd": 26910003,
    "rwcalr": 27000000,
    "rwd": "7777777777777777777777777777777777777777777777777774MSJUVU",
    "seed": "cTBVGCaNSy6lpFHEZC/wogPmDeNM7PCUliLQl+HEjH4=",
    "spt": {
      "0": {
        "n": 26909952
      }
    },
    "tc": 153348366,
    "ts": 1673397878,
    "txn": "7t4+E0Vhj47kfCrg1PUvDeBffSPhp86wUuI4kBsCjls=",
    "txn256": "2EoaVSWHtOTfTf/5g88WppeKdTPtIcMx7wc+nHm1OyM="
  },
  "KvMods": null,
  "PrevTimestamp": 1673397875,
  "StateProofNext": 0,
  "Totals": {
    "notpart": {
      "mon": 347207146836,
      "rwd": 346468
    },
    "offline": {
      "mon": 2544590110560384,
      "rwd": 2526267564
    },
    "online": {
      "mon": 7580062682392780,
      "rwd": 7580062232
    },
    "rwdlvl": 27521
  },
  "Txids": {
    "582J5rzaCICaxgGnKz0WiCIB5kpB8krGHZX/W6Bd8dY=": {
      "Intra": 4,
      "LastValid": 26910007
    },
    "DXMq5UdOubNb2b/CtSv+3i2Sk+/lFY/Kmg0E9P+wlEc=": {
      "Intra": 1,
      "LastValid": 26910007
    },
    "VH47fPqBWEnB3q0U0/kTcJCpC3VlkCTYer39o8xwlvE=": {
      "Intra": 2,
      "LastValid": 26910007
    },
    "tLiR0xUwx2TpaWCP8q6ycBZqtiGWoUUFdFSmVETieYs=": {
      "Intra": 3,
      "LastValid": 26910007
    },
    "vVLDNK8yxToHiFKuS4CqscfA59+jyppclzD9zX9+sIc=": {
      "Intra": 0,
      "LastValid": 26910007
    }
  },
  "Txleases": null
}`
var actual = `{
  "Accts": {
    "Accts": [
      {
        "Addr": "//////////////////////////////////////////8=",
        "AuthAddr": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
        "MicroAlgos": 2051818321,
        "RewardedMicroAlgos": 0,
        "RewardsBase": 0,
        "SelectionID": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
        "StateProofID": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==",
        "Status": 2,
        "TotalAppLocalStates": 0,
        "TotalAppParams": 0,
        "TotalAppSchema": {},
        "TotalAssetParams": 0,
        "TotalAssets": 0,
        "TotalBoxBytes": 0,
        "TotalBoxes": 0,
        "TotalExtraAppPages": 0,
        "VoteFirstValid": 0,
        "VoteID": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
        "VoteKeyDilution": 0,
        "VoteLastValid": 0
      },
      {
        "Addr": "oqQYH8Vx35YSH8PHIewV7y8rxKb/OW9YvIUKQvZLvVA=",
        "AuthAddr": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
        "MicroAlgos": 295200,
        "RewardedMicroAlgos": 0,
        "RewardsBase": 27521,
        "SelectionID": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
        "StateProofID": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==",
        "Status": 0,
        "TotalAppLocalStates": 0,
        "TotalAppParams": 0,
        "TotalAppSchema": {},
        "TotalAssetParams": 0,
        "TotalAssets": 0,
        "TotalBoxBytes": 0,
        "TotalBoxes": 0,
        "TotalExtraAppPages": 0,
        "VoteFirstValid": 0,
        "VoteID": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
        "VoteKeyDilution": 0,
        "VoteLastValid": 0
      },
      {
        "Addr": "B9rLS22e0UGxdXa9RZrmQh1IbaPU7yJHxAmjlrguoiE=",
        "AuthAddr": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
        "MicroAlgos": 342106549750,
        "RewardedMicroAlgos": 0,
        "RewardsBase": 0,
        "SelectionID": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
        "StateProofID": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==",
        "Status": 2,
        "TotalAppLocalStates": 0,
        "TotalAppParams": 0,
        "TotalAppSchema": {},
        "TotalAssetParams": 0,
        "TotalAssets": 0,
        "TotalBoxBytes": 0,
        "TotalBoxes": 0,
        "TotalExtraAppPages": 0,
        "VoteFirstValid": 0,
        "VoteID": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
        "VoteKeyDilution": 0,
        "VoteLastValid": 0
      }
    ],
    "AppResources": null,
    "AssetResources": null
  },
  "Creatables": null,
  "Hdr": {
    "earn": 27521,
    "fees": "B9rLS22e0UGxdXa9RZrmQh1IbaPU7yJHxAmjlrguoiE=",
    "frac": 2047917471,
    "gen": "testnet-v1.0",
    "gh": "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI=",
    "prev": "7Mu8YMIj9l7j31y5Bo/GdpO1xdPacVEkNFjtAW1IeCY=",
    "proto": "https://github.com/algorandfoundation/specs/tree/44fa607d6051730f5264526bf3c108d51f0eadb6",
    "rate": 42,
    "rnd": 26910003,
    "rwcalr": 27000000,
    "rwd": "//////////////////////////////////////////8=",
    "seed": "cTBVGCaNSy6lpFHEZC/wogPmDeNM7PCUliLQl+HEjH4=",
    "spt": {
      0: {
        "n": 26909952
      }
    },
    "tc": 153348366,
    "ts": 1673397878,
    "txn": "7t4+E0Vhj47kfCrg1PUvDeBffSPhp86wUuI4kBsCjls=",
    "txn256": "2EoaVSWHtOTfTf/5g88WppeKdTPtIcMx7wc+nHm1OyM="
  },
  "KvMods": null,
  "PrevTimestamp": 1673397875,
  "StateProofNext": 0,
  "Totals": {
    "notpart": {
      "mon": 347207146836,
      "rwd": 346468
    },
    "offline": {
      "mon": 2544590110560384,
      "rwd": 2526267564
    },
    "online": {
      "mon": 7580062682392780,
      "rwd": 7580062232
    },
    "rwdlvl": 27521
  },
  "Txids": {
    "582J5rzaCICaxgGnKz0WiCIB5kpB8krGHZX/W6Bd8dY=": {
      "Intra": 4,
      "LastValid": 26910007
    },
    "DXMq5UdOubNb2b/CtSv+3i2Sk+/lFY/Kmg0E9P+wlEc=": {
      "Intra": 1,
      "LastValid": 26910007
    },
    "VH47fPqBWEnB3q0U0/kTcJCpC3VlkCTYer39o8xwlvE=": {
      "Intra": 2,
      "LastValid": 26910007
    },
    "tLiR0xUwx2TpaWCP8q6ycBZqtiGWoUUFdFSmVETieYs=": {
      "Intra": 3,
      "LastValid": 26910007
    },
    "vVLDNK8yxToHiFKuS4CqscfA59+jyppclzD9zX9+sIc=": {
      "Intra": 0,
      "LastValid": 26910007
    }
  },
  "Txleases": null
}`

func TestIt(t *testing.T) {
	fmt.Println("???")
	require.NoError(t, test.EqualJson2(expected, actual))
}
