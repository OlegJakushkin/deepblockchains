package stark

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/rlp"
)

// Full STARK test
func TestStark(t *testing.T) {

	INPUT := big.NewInt(3)
	LOGSTEPS := big.NewInt(13)

	// 64 constants
	constants := make([]*big.Int, 64)
	for i := int64(0); i < 64; i++ {
		b := big.NewInt(i)
		b.Exp(b, big.NewInt(7), nil)
		constants[i] = b.Xor(b, big.NewInt(42))
	}
	f, _ := NewPrimeField(nil)
	two_logsteps := new(big.Int).Exp(big.NewInt(2), LOGSTEPS, nil)

	// Generate STARK proof
	start := time.Now()
	proof, err := mk_mimc_proof(f, INPUT, two_logsteps, constants)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	fmt.Printf("STARK computed in %s\n\n", time.Since(start))

	// Encoding
	fmt.Printf("STARK Proof size: ")
	encoded, _ := rlp.EncodeToBytes(proof)
	TOT := len(encoded)
	L2 := 0
	for layer, comp := range proof.Child {
		encoded_comp, _ := rlp.EncodeToBytes(comp)
		L2 = L2 + len(encoded_comp)
		fmt.Printf(" layer %d: %d bytes | ", layer, len(encoded_comp))
	}
	fmt.Printf("\nApprox proof length: %d bytes (branches), %d bytes (FRI proof), %d bytes (total)\n", TOT-L2, L2, TOT)

	// Decoding
	var prf Proof
	err = rlp.Decode(bytes.NewReader(encoded), &prf)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	// Verify STARK proof
	start = time.Now()
	output := f.mimc(INPUT, two_logsteps, constants)
	verify_mimc_proof(f, INPUT, two_logsteps, constants, output, &prf)
	fmt.Printf("STARK verified in %s\n", time.Since(start))
}
