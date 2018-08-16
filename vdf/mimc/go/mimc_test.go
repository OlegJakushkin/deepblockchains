package stark

import (
	"fmt"
	"math/big"
	"testing"
	"time"
)

/*
$ go test -run MiMC
Forward MiMC1: 7.074633ms
Reverse MiMC1: 279.084117ms 3
MATCH rtrace=3
Forward MiMC2: 70.682971ms
Reverse MiMC2: 2.172787695s 3
MATCH rtrace[0]=3
PASS
*/

func DefaultModulus() *big.Int {
	t0 := new(big.Int)
	t1 := new(big.Int)
	t2 := new(big.Int)
	return new(big.Int).Add(t1.Sub(t0.Exp(big.NewInt(2), big.NewInt(256), nil), t2.Mul(t1.Exp(big.NewInt(2), big.NewInt(32), nil), big.NewInt(351))), big.NewInt(1))
}

func TestMiMC(t *testing.T) {
	t0 := new(big.Int)
	t1 := new(big.Int)
	t2 := new(big.Int)
	modulus := DefaultModulus()

	// 64 constants
	SEVEN := big.NewInt(7)
	FORTYTWO := big.NewInt(42)
	round_constants := make([]*big.Int, 64)
	for i := int64(0); i < 64; i++ {
		round_constants[i] = new(big.Int).Xor(t2.Exp(big.NewInt(i), SEVEN, nil), FORTYTWO)
	}

	// Forward MiMC
	nsteps := 8192
	input := big.NewInt(3)
	trace := new(big.Int).Set(input)
	start := time.Now()
	for i := 1; i < nsteps; i++ {
		trace.Mod(t2.Add(t1.Mul(trace, t0.Mul(trace, trace)), round_constants[i%len(round_constants)]), modulus)
	}
	output := new(big.Int).Set(trace)
	fmt.Printf("forward-mimc: %s\n", time.Since(start))

	// Reverse MiMC
	start = time.Now()
	rtrace := new(big.Int).Set(output)
	little_fermat_expt := new(big.Int).Div(t2.Sub(t1.Mul(big.NewInt(2), modulus), big.NewInt(1)), big.NewInt(3))
	for i := nsteps - 1; i > 0; i-- {
		rtrace.Exp(t2.Sub(rtrace, round_constants[i%len(round_constants)]), little_fermat_expt, modulus)
	}
	fmt.Printf("reverse-mimc: %s\n", time.Since(start))
	if rtrace.Cmp(input) == 0 {
		fmt.Printf("PASS\n")
	} else {
		t.Fatalf("FAIL\n")
	}
}
