package bloom

import (
	"hash"
	"math/big"
	"math/rand"
)

type BloomFilter struct {
	i       *big.Int // use a bignum to simulate a bitset (it's bounded by size below)
	size    uint32   // how many bits?
	k       int      // how many hash functions?
	seeds   []uint32 // one seed per hash function
	hfs     []hash.Hash32
	records int64 // number of added records
}

func (bf *BloomFilter) getPosition(s string, hf hash.Hash32) int {

	// hash, reuse
	hf.Write([]byte(s))
	h := hf.Sum32()
	hf.Reset()

	return int(h % bf.size)
}

func New(numElements int, maxFalsePositiveProb float32) *BloomFilter {

	bucketsPerElement, k := ComputeOptimal(maxFalsePositiveProb)

	bf := &BloomFilter{
		i:    big.NewInt(0),
		k:    k,
		size: uint32(numElements*bucketsPerElement + 20),
	}

	// generate hash functions
	for i := 0; i < bf.k; i++ {
		seed := rand.Uint32()
		bf.seeds = append(bf.seeds, seed)
		bf.hfs = append(bf.hfs, New32(seed))
	}

	return bf
}

func (bf *BloomFilter) Add(s string) {

	for _, hf := range bf.hfs {
		bf.i.SetBit(bf.i, bf.getPosition(s, hf), 1)
	}

	bf.records += 1
}

func (bf *BloomFilter) Contains(s string) bool {

	for _, hf := range bf.hfs {

		if bf.i.Bit(bf.getPosition(s, hf)) == 0 {
			return false
		}
	}

	return true
}

func (bf *BloomFilter) Size() int64 {
	return bf.records
}
