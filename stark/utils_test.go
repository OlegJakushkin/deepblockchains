package stark

import (
	"fmt"
	"math/big"
	"runtime"
	"testing"
	"time"
)

func compute_with_cores(ch_column chan IndexedInt, xs []*big.Int, i0 int) {
	MODULUS := big.NewInt(351)
	SEVEN := big.NewInt(7)
	t := new(big.Int)
	for _i, x := range xs {
		ch_column <- IndexedInt{idx: i0 + _i, data: new(big.Int).Mod(t.Mul(x, SEVEN), MODULUS)}
	}
}

// Tiny Cores test
func TestCores(t *testing.T) {
	runtime.GOMAXPROCS(NUM_CORES)

	nev := 65536
	inp := make([]*big.Int, nev)
	t0 := new(big.Int)
	MODULUS := big.NewInt(351)
	for i := 0; i < nev; i++ {
		inp[i] = new(big.Int).Mod(t0.SetInt64(int64(i)), MODULUS)
	}
	column := make([]*big.Int, nev)
	ch_column := make(chan IndexedInt, nev)
	nj := min_iterations(nev, 4)
	start := time.Now()
	for i := 0; i < nev; i += nj {
		i1 := i + nj
		if i1 > nev {
			i1 = nev
		}
		go compute_with_cores(ch_column, inp[i:i1], i)
	}

	for i := 0; i < nev; i++ {
		d := <-ch_column
		column[d.idx] = d.data
	}
	fmt.Printf("Cores: %s\n", time.Since(start))
}
