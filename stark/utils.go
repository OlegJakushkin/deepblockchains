package stark

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/blake2s"
)

const (
	NUM_CORES = 8
)

// structures used for channel

type IndexedBytes struct {
	idx  int
	data []byte
}

type IndexedBranchSet struct {
	idx  int
	data [][][]byte
}

type IndexedInt struct {
	idx  int
	data *big.Int
}

type IndexedIntArray struct {
	idx  int
	data []*big.Int
}

type IndexedIntArrays struct {
	idx  int
	data [][]*big.Int
}

func blake(b []byte) *big.Int {
	h, _ := blake2s.New256(nil)
	h.Write(b)
	a := h.Sum(nil)
	return common.BytesToHash(a).Big()
}

func blakeb(b []byte) []byte {
	h, _ := blake2s.New256(nil)
	h.Write(b)
	return h.Sum(nil)
}

func is_a_power_of_2(x *big.Int) bool {
	if x.Cmp(ONE) == 0 {
		return true
	}
	lowbit := new(big.Int).Mod(x, TWO)
	if lowbit.Cmp(ZERO) > 0 {
		return false
	}
	y := new(big.Int).Div(x, TWO)
	return is_a_power_of_2(y)
}

func BigToBytes(i *big.Int) []byte {
	b := i.Bytes()
	out := make([]byte, 32)
	copy(out[(32-len(b)):32], b[0:len(b)])
	return out
}

func BytesToBig(i []byte) *big.Int {
	return new(big.Int).SetBytes(i)
}

// Get the set of powers of R, until but not including when the powers loop back to 1
func get_power_cycle(r *big.Int, modulus *big.Int) []*big.Int {
	o := []*big.Int{ONE, r}
	t1 := new(big.Int)
	last := new(big.Int).Set(r)
	for last.Cmp(ONE) != 0 {
		last = new(big.Int).Mod(t1.Mul(last, r), modulus)
		o = append(o, last)
	}
	return o[:len(o)-1]
}

func min_iterations(total_iterations int, min_size_level int) int {
	min_per_iteration := total_iterations / NUM_CORES
	m := 512
	switch min_size_level {
	case -1:
		m = 8
	case 0:
		m = 512
	case 1:
		m = 512
	case 2:
		m = 512
	case 3:
		m = 512
	case 4:
		m = 512
	case 5:
		m = 1024
	case 6:
		m = 2048
	}
	if min_per_iteration < m {
		min_per_iteration = m
	}
	return min_per_iteration
}
