package stark

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func merkelize(L [][]byte) [][]byte {
	LH := make([][]byte, len(L))
	for i, v := range L {
		LH[i] = v
	}
	nodes := make([][]byte, len(L))
	nodes = append(nodes, LH...)
	for i := len(L) - 1; i >= 0; i-- {
		nodes[i] = blakeb(append(nodes[i*2], nodes[i*2+1]...))
		// fmt.Printf("MERKLIZE node %d %x <-- %d:%x %d:%x\n", i, nodes[i], i*2, nodes[i*2], i*2+1, nodes[i*2+1])
	}
	return nodes
}

func mk_branch(tree [][]byte, index int64) [][]byte {
	index += int64(len(tree)) / 2
	o := make([][]byte, 1)
	o[0] = tree[index]
	for index > 1 {
		o = append(o, tree[index^1])
		index = index / 2
	}
	return o
}

func verify_branch_int(root []byte, index uint, proof [][]byte) (res *big.Int, err error) {
	res_byte, err := verify_branch(root, index, proof)
	if err != nil {
		return res, err
	}
	res = common.BytesToHash(res_byte).Big()
	return res, nil
}

func verify_branch(root []byte, index uint, proof [][]byte) (res []byte, err error) {
	q := 1 << uint(len(proof)) // 2**len(proof)
	index += uint(q)
	v := proof[0]
	for _, p := range proof[1:] {
		if index%2 > 0 {
			v = blakeb(append(p, v...))
		} else {
			v = blakeb(append(v, p...))
		}
		index = index / 2
	}
	if bytes.Compare(v, root) != 0 {
		return res, fmt.Errorf("Mismatch root, got:[%x] expected:[%x]", v, root)
	}
	return proof[0], nil
}
