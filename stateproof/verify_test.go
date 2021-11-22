package stateproof

import (
	"fmt"
	"testing"

	crypto "github.com/algorand/go-algorand-sdk/algod-crypto"
	"github.com/algorand/go-algorand-sdk/algod-crypto/merklearray"
	"github.com/algorand/go-algorand/protocol"
	"github.com/algorand/go-algorand/test/partitiontest"
	"github.com/stretchr/testify/require"

)

func TestMerkle(t *testing.T) {
	partitiontest.PartitionTest(t)

	for i := uint64(0); i < 1024; i++ {
		testMerkle(t, crypto.Sha512_256, i)
	}

	if !testing.Short() {
		for i := uint64(0); i < 10; i++ {
			testMerkle(t, crypto.Sumhash, i)
		}
	}

}

type TestData crypto.Digest
func (d TestData) ToBeHashed() (protocol.HashID, []byte) {
	return protocol.Message, d[:]
}

type TestArray []TestData
func (a TestArray) Length() uint64 {
	return uint64(len(a))
}
func hashRep(h crypto.Hashable) []byte {
	hashid, data := h.ToBeHashed()
	return append([]byte(hashid), data...)
}
func (a TestArray) Marshal(pos uint64) ([]byte, error) {
	if pos >= uint64(len(a)) {
		return nil, fmt.Errorf("pos %d larger than length %d", pos, len(a))
	}

	return hashRep(a[pos]), nil
}

func testMerkle(t *testing.T, hashtype crypto.HashType, size uint64) {
	var junk TestData
	crypto.RandBytes(junk[:])

	a := make(TestArray, size)
	for i := uint64(0); i < size; i++ {
		crypto.RandBytes(a[i][:])
	}

	tree, err := merklearray.Build(a, crypto.HashFactory{HashType: hashtype})
	require.NoError(t, err)

	root := tree.Root()

	var allpos []uint64
	allmap := make(map[uint64]crypto.Hashable)

	for i := uint64(0); i < size; i++ {
		proof, err := tree.Prove([]uint64{i})
		require.NoError(t, err)

		err = Verify(root, map[uint64]crypto.Hashable{i: a[i]}, proof)
		require.NoError(t, err)

		err = Verify(root, map[uint64]crypto.Hashable{i: junk}, proof)
		require.Error(t, err, "no error when verifying junk")

		allpos = append(allpos, i)
		allmap[i] = a[i]
	}

	proof, err := tree.Prove(allpos)
	require.NoError(t, err)

	err = Verify(root, allmap, proof)
	require.NoError(t, err)

	err = Verify(root, map[uint64]crypto.Hashable{0: junk}, proof)
	require.Error(t, err, "no error when verifying junk batch")

	err = Verify(root, map[uint64]crypto.Hashable{0: junk}, nil)
	require.Error(t, err, "no error when verifying junk batch")

	_, err = tree.Prove([]uint64{size})
	require.Error(t, err, "no error when proving past the end")

	err = Verify(root, map[uint64]crypto.Hashable{size: junk}, nil)
	require.Error(t, err, "no error when verifying past the end")

	if size > 0 {
		var somepos []uint64
		somemap := make(map[uint64]crypto.Hashable)
		for i := 0; i < 10; i++ {
			pos := crypto.RandUint64() % size
			somepos = append(somepos, pos)
			somemap[pos] = a[pos]
		}

		proof, err = tree.Prove(somepos)
		require.NoError(t, err)

		err = Verify(root, somemap, proof)
		require.NoError(t, err)
	}
}
