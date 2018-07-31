package stark

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/blake2s"
)

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
	one := big.NewInt(1)
	if x.Cmp(one) == 0 {
		return true
	}
	two := big.NewInt(2)
	lowbit := new(big.Int).Mod(x, two)
	if lowbit.Cmp(big.NewInt(0)) > 0 {
		return false
	}
	y := new(big.Int).Div(x, two)
	return is_a_power_of_2(y)
}

// Get the set of powers of R, until but not including when the powers loop back to 1
func get_power_cycle(r *big.Int, modulus *big.Int) []*big.Int {
	o := make([]*big.Int, 0)
	o = append(o, new(big.Int).SetUint64(1))
	o = append(o, r)
	for o[len(o)-1].Cmp(big.NewInt(1)) != 0 {
		t := new(big.Int).Mul(o[len(o)-1], r)
		t.Mod(t, modulus)
		o = append(o, t)
	}
	return o[:len(o)-1]
}
