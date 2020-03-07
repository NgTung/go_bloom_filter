package filters

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash"
	"math"
)

const P = 0.000001 // Probability of false positive: 1 in 1,000,000

//BloomFilters ...
type BloomFilters interface {
	Add(value interface{})

	MightContain(value interface{}) bool
}

type bloomFilterImpl struct {
	bitSet       []int
	hashFuncList []hash.Hash
	size         int64
}

// Factory of Bloom Filter Implementation
// - elementNumber - The estimated item size
// - hashList - The list of hash functions
func InitBloomFilter(elementNumber int32, hashList []hash.Hash) BloomFilters {
	bitSetSize := getSizeOfBitSet(elementNumber) // Number of bits in the filter
	return &bloomFilterImpl{
		bitSet:       make([]int, bitSetSize),
		hashFuncList: hashList,
		size:         bitSetSize, // Re-use purpose, prevent calculate frequently the size
	}
}

// Add a value to the bit set
func (b *bloomFilterImpl) Add(value interface{}) {
	for _, h := range b.hashFuncList {
		hashPosition := getPosition(hashing(h, value), b.size)
		b.bitSet[hashPosition] = 1
	}
}

// MightContain check if an value is contained in the bit set
func (b *bloomFilterImpl) MightContain(value interface{}) bool {
	for _, h := range b.hashFuncList {
		hashPosition := getPosition(hashing(h, value), b.size)
		if b.bitSet[hashPosition] == 0 {
			return false
		}
	}
	return true
}

// Return number of bits in the filter calculated by following formula:
//    bitSetSize = (elementNum * log(P)) / log(1 / 2^log(2))
//
//    where:
//    - elementNum: Number of items
//    - P: False Positive probability
//
// For further formula info, refer to
//  - https://en.wikipedia.org/wiki/Bloom_filter
//  - https://hur.st/bloomfilter/ (for visualize)
//
func getSizeOfBitSet(elementNum int32) int64 {
	bitSetSize := math.Ceil((float64(elementNum) * math.Log(P)) / math.Log(1/math.Pow(2, math.Log(2))))
	return int64(bitSetSize)
}

func getPosition(hashValue int64, bitSetSize int64) int64 {
	hashPosition := hashValue % bitSetSize
	return int64(math.Abs(float64(hashPosition)))
}

func hashing(h hash.Hash, value interface{}) int64 {
	h.Write([]byte(toString(value)))
	bits := h.Sum(nil)
	buffer := bytes.NewBuffer(bits)
	result, _ := binary.ReadVarint(buffer)
	h.Reset()
	return result
}

func toString(value interface{}) string {
	return fmt.Sprintf("%v", value)
}
