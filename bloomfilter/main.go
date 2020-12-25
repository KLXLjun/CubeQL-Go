package bloomfilter

import "github.com/willf/bitset"

const DEFAULT_SIZE = 2 << 24

var seeds = []uint{7, 11, 13, 31, 37, 61}

//BloomFilter 类型
type BloomFilter struct {
	Set   *bitset.BitSet
	Funcs [6]SimpleHash
}

//NewBloomFilter  新的布隆过滤器
func NewBloomFilter() *BloomFilter {
	bf := new(BloomFilter)
	for i := 0; i < len(bf.Funcs); i++ {
		bf.Funcs[i] = SimpleHash{DEFAULT_SIZE, seeds[i]}
	}
	bf.Set = bitset.New(DEFAULT_SIZE)
	return bf
}

//Add 添加
func (bf BloomFilter) Add(value string) {
	for _, f := range bf.Funcs {
		bf.Set.Set(f.hash(value))
	}
}

//Contains 是否存在
func (bf BloomFilter) Contains(value string) bool {
	if value == "" {
		return false
	}
	ret := true
	for _, f := range bf.Funcs {
		ret = ret && bf.Set.Test(f.hash(value))
	}
	return ret
}

type SimpleHash struct {
	Cap  uint
	Seed uint
}

func (s SimpleHash) hash(value string) uint {
	var result uint = 0
	for i := 0; i < len(value); i++ {
		result = result*s.Seed + uint(value[i])
	}
	return (s.Cap - 1) & result
}
