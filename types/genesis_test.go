package types

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/algorand/go-algorand-sdk/encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeDecodeGenesis(t *testing.T) {
	genesisStr := "{\n  \"alloc\": [\n    {\n      \"addr\": \"7777777777777777777777777777777777777777777777777774MSJUVU\",\n      \"comment\": \"RewardsPool\",\n      \"state\": {\n        \"algo\": 125000000000000,\n        \"onl\": 2\n      }\n    },\n    {\n      \"addr\": \"LX6VRFCL4723LRO3VEX5XNOXHFXFQKATKDLDBT5HNT5DFDC3SQG4EI2LSI\",\n      \"comment\": \"Wallet1\",\n      \"state\": {\n        \"algo\": 1000000000000000,\n        \"onl\": 1,\n        \"sel\": \"XGIljovdG8IjSQJRcvc/vi6UXXTIddUgjZACWZPW23M=\",\n        \"vote\": \"zccp6OB5EOh21WMV7A7iVg+CitiZeOao+VYbIkNdPgM=\",\n        \"voteKD\": 10000,\n        \"voteLst\": 3000000\n      }\n    }\n  ],\n  \"fees\": \"A7NMWS3NT3IUDMLVO26ULGXGIIOUQ3ND2TXSER6EBGRZNOBOUIQXHIBGDE\",\n  \"id\": \"v1\",\n  \"network\": \"sandnet\",\n  \"proto\": \"future\",\n  \"rwd\": \"7777777777777777777777777777777777777777777777777774MSJUVU\"\n}"
	var genesis Genesis
	json.Decode([]byte(genesisStr), &genesis)
	assert.Equal(t, base64.StdEncoding.EncodeToString(json.Encode(genesis)), base64.StdEncoding.EncodeToString([]byte(genesisStr)))
}

func TestGenesis_Balances(t *testing.T) {
	containsErrorFunc := func(str string) assert.ErrorAssertionFunc {
		return func(_ assert.TestingT, err error, i ...interface{}) bool {
			require.ErrorContains(t, err, str)
			return true
		}
	}
	mustAddr := func(addr string) Address {
		address, err := DecodeAddress(addr)
		require.NoError(t, err)
		return address
	}
	makeAddr := func(addr uint64) Address {
		var address Address
		address[0] = byte(addr)
		return address
	}
	acctWith := func(algos uint64, addr string) GenesisAllocation {
		return GenesisAllocation{
			_struct: struct{}{},
			Address: addr,
			Comment: "",
			State: Account{
				MicroAlgos: algos,
			},
		}
	}
	goodAddr := makeAddr(100)
	allocation1 := acctWith(1000, makeAddr(1).String())
	allocation2 := acctWith(2000, makeAddr(2).String())
	badAllocation := acctWith(1234, "El Toro Loco")
	type fields struct {
		Allocation  []GenesisAllocation
		FeeSink     string
		RewardsPool string
	}
	tests := []struct {
		name    string
		fields  fields
		want    GenesisBalances
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "basic test",
			fields: fields{
				Allocation:  []GenesisAllocation{allocation1},
				FeeSink:     goodAddr.String(),
				RewardsPool: goodAddr.String(),
			},
			want: GenesisBalances{
				Balances: map[Address]Account{
					mustAddr(allocation1.Address): allocation1.State,
				},
				FeeSink:     goodAddr,
				RewardsPool: goodAddr,
				Timestamp:   0,
			},
			wantErr: assert.NoError,
		},
		{
			name: "two test",
			fields: fields{
				Allocation:  []GenesisAllocation{allocation1, allocation2},
				FeeSink:     goodAddr.String(),
				RewardsPool: goodAddr.String(),
			},
			want: GenesisBalances{
				Balances: map[Address]Account{
					mustAddr(allocation1.Address): allocation1.State,
					mustAddr(allocation2.Address): allocation2.State,
				},
				FeeSink:     goodAddr,
				RewardsPool: goodAddr,
				Timestamp:   0,
			},
			wantErr: assert.NoError,
		},
		{
			name: "bad fee sink",
			fields: fields{
				Allocation:  []GenesisAllocation{allocation1, allocation2},
				RewardsPool: goodAddr.String(),
			},
			wantErr: containsErrorFunc("cannot parse fee sink addr"),
		},
		{
			name: "bad rewards pool",
			fields: fields{
				Allocation: []GenesisAllocation{allocation1, allocation2},
				FeeSink:    goodAddr.String(),
			},
			wantErr: containsErrorFunc("cannot parse rewards pool addr"),
		},
		{
			name: "bad genesis addr",
			fields: fields{
				Allocation:  []GenesisAllocation{badAllocation},
				FeeSink:     goodAddr.String(),
				RewardsPool: goodAddr.String(),
			},
			wantErr: containsErrorFunc("cannot parse genesis addr"),
		},
		{
			name: "repeat address",
			fields: fields{
				Allocation:  []GenesisAllocation{allocation1, allocation1},
				FeeSink:     goodAddr.String(),
				RewardsPool: goodAddr.String(),
			},
			wantErr: containsErrorFunc("repeated allocation to"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			genesis := Genesis{
				Allocation:  tt.fields.Allocation,
				FeeSink:     tt.fields.FeeSink,
				RewardsPool: tt.fields.RewardsPool,
			}
			got, err := genesis.Balances()
			if tt.wantErr(t, err, fmt.Sprintf("Balances()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "Balances()")
		})
	}
}
