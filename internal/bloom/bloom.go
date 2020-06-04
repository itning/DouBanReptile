package bloom

import (
	"github.com/willf/bitset"
)

const DefaultSize = 2 << 24

var seeds = []uint{7, 11, 13, 31, 37, 61}

type SimpleHash struct {
	cap  uint
	seed uint
}

type Filter struct {
	set   *bitset.BitSet
	funcs [6]SimpleHash
}

func NewBloomFilter() *Filter {
	bf := new(Filter)
	for i := 0; i < len(bf.funcs); i++ {
		bf.funcs[i] = SimpleHash{DefaultSize, seeds[i]}
	}
	bf.set = bitset.New(DefaultSize)
	return bf
}

func (bf Filter) Add(value string) {
	for _, f := range bf.funcs {
		bf.set.Set(f.hash(value))
	}
}

func (bf Filter) Contains(value string) bool {
	if value == "" {
		return false
	}
	ret := true
	for _, f := range bf.funcs {
		ret = ret && bf.set.Test(f.hash(value))
	}
	return ret
}

func (s SimpleHash) hash(value string) uint {
	var result uint = 0
	for i := 0; i < len(value); i++ {
		result = result*s.seed + uint(value[i])
	}
	return (s.cap - 1) & result
}
