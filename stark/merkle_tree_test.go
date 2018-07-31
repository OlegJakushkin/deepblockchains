package stark

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
)

func TestMerkleTree(t *testing.T) {
	nitems := int64(128)
	o := make([][]byte, nitems)
	for x := int64(0); x < nitems; x++ {
		o[x] = big.NewInt(x).Bytes()
	}
	mtree := merkelize(o)
	b := mk_branch(mtree, 59)
	res, err := verify_branch(mtree[1], 59, b)
	if err != nil {
		t.Fatalf("Merkle tree failure %v", err)
	}
	fmt.Printf("Merkle tree works %v", res)
}

func TestBlake(t *testing.T) {
	res := blakeb([]byte("Hello"))
	correct_res, _ := hex.DecodeString("f73a5fbf881f89b814871f46e26ad3fa37cb2921c5e8561618639015b3ccbb71")

	if bytes.Compare(res, correct_res) != 0 {
		t.Fatalf("Blake failure: got %x, expected %x\n", res, correct_res)
	}
}
