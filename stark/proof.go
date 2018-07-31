package stark

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

const (
	spot_check_security_factor = 80
	extension_factor           = 8
)

type Proof struct {
	Root     []byte
	LRoot    []byte
	Branches [][][]byte
	Child    []*FriComponent
}

// []byte hash, [][]byte branch, [][][]byte set of branches
type FriComponent struct {
	Values   [][]byte
	Root     []byte
	Branches [][][][]byte
}

// Extract pseudorandom indices from entropy
func get_pseudorandom_indices(f *PrimeField, seed []byte, modulus *big.Int, count int, exclude_multiples_of int64) (o []*big.Int, err error) {
	max_modulus := new(big.Int).Exp(big.NewInt(2), big.NewInt(24), nil)
	if modulus.Cmp(max_modulus) > 0 {
		return o, fmt.Errorf("modulus must be less than 2**24")
	}
	data := seed
	for len(data) < count*4 {
		b := data[len(data)-32 : len(data)]
		data = append(data, blakeb(b)...)
	}

	o = make([]*big.Int, 0)
	if exclude_multiples_of == 0 {
		for i := 0; i < count*4; i += 4 {
			t := common.BytesToHash(data[i : i+4]).Big()
			o = append(o, t.Mod(t, f.modulus))
		}
	} else {
		real_modulus := new(big.Int).Mul(modulus, big.NewInt(int64(exclude_multiples_of-1)))
		real_modulus.Div(real_modulus, big.NewInt(exclude_multiples_of))

		for i := 0; i < count*4; i += 4 {
			b := make([]byte, 32)
			copy(b[28:32], data[i:i+4])
			t := common.BytesToHash(b).Big()
			t.Mod(t, real_modulus)
			o = append(o, t)
		}

		for i := 0; i < len(o); i++ {
			s := new(big.Int).Add(o[i], big.NewInt(1))
			t := new(big.Int).Div(o[i], new(big.Int).Sub(big.NewInt(exclude_multiples_of), big.NewInt(1)))
			s.Add(s, t)
			o[i] = s
		}
	}
	return o, nil
}

// Generate a STARK for a MIMC calculation
func mk_mimc_proof(f *PrimeField, inp *big.Int, steps *big.Int, round_constants []*big.Int) (resultProof *Proof, err error) {
	// Some constraints to make our job easier
	two_power_32 := new(big.Int).Exp(big.NewInt(2), big.NewInt(32), nil)
	ef := big.NewInt(int64(extension_factor))
	max_steps := new(big.Int).Div(two_power_32, ef)
	if steps.Cmp(max_steps) > 0 {
		return resultProof, fmt.Errorf("too many steps")
	}
	if !is_a_power_of_2(steps) {
		return resultProof, fmt.Errorf("steps must be power of 2")
	}
	rc := big.NewInt(int64(len(round_constants)))
	if !is_a_power_of_2(rc) {
		return resultProof, fmt.Errorf("len(round_constants) must be power of 2")
	}
	if steps.Cmp(rc) < 0 {
		return resultProof, fmt.Errorf("len(round_constants) too high")
	}

	precision := new(big.Int).Mul(steps, ef)

	// Root of unity such that x^precision=1
	t := new(big.Int).Sub(f.modulus, big.NewInt(1))
	G2 := f.pow(big.NewInt(7), t.Div(t, precision))

	// Root of unity such that x^steps=1
	skips := new(big.Int).Div(precision, steps)
	G1 := f.pow(G2, skips)

	// Powers of the higher-order root of unity
	xs := get_power_cycle(G2, f.modulus)
	last_step_position := xs[(int(steps.Int64())-1)*extension_factor]

	// Generate the computational trace
	computational_trace := make([]*big.Int, 1)
	computational_trace[0] = inp
	for i := 0; i < int(steps.Int64()-1); i++ {
		t := new(big.Int).Exp(computational_trace[len(computational_trace)-1], big.NewInt(3), nil)
		t.Add(t, round_constants[i%len(round_constants)])
		computational_trace = append(computational_trace, t.Mod(t, f.modulus))
	}
	output := computational_trace[len(computational_trace)-1]
	fmt.Printf("Computational trace output: %v\n", output)

	// Interpolate the computational trace into a polynomial P, with each step along a successive power of G1
	computational_trace_polynomial := f.fft(computational_trace, G1, true)
	p_evaluations := f.fft(computational_trace_polynomial, G2, false)
	fmt.Printf("Converted computational steps into a polynomial and low-degree extended it %v\n", len(computational_trace_polynomial))

	skips2 := new(big.Int).Div(steps, rc)
	constants_mini_polynomial := f.fft(round_constants, f.pow(G1, skips2), true)
	// constants_polynomial = [0 if i % skips2 else constants_mini_polynomial[i//skips2] for i in range(steps)]
	constants_polynomial := make([]*big.Int, steps.Uint64())
	for i := int64(0); i < steps.Int64(); i++ {
		t := new(big.Int).Mod(big.NewInt(i), skips2)
		if t.Cmp(big.NewInt(0)) == 0 {
			constants_polynomial[i] = big.NewInt(0)
		} else {
			idx := new(big.Int).Div(big.NewInt(i), skips2)
			constants_polynomial[i] = constants_mini_polynomial[idx.Uint64()]
		}
	}
	constants_mini_extension := f.fft(constants_mini_polynomial, f.pow(G2, skips2), false)
	fmt.Printf("Converted round constants into a polynomial and low-degree extended it 8192 [%d]\n", len(constants_mini_polynomial))

	// Create the composed polynomial such that C(P(x), P(g1*x), K(x)) = P(g1*x) - P(x)**3 - K(x)
	c_of_p_evaluations := make([]*big.Int, precision.Uint64())
	for i := int(0); i < int(precision.Int64()); i++ {
		idx := int64(i+extension_factor) % precision.Int64()
		s0 := new(big.Int).Set(p_evaluations[idx])
		t := new(big.Int).Mul(p_evaluations[i], p_evaluations[i])
		t.Mul(t, p_evaluations[i])
		t.Mod(t, f.modulus)
		s0.Sub(s0, t)
		s0.Sub(s0, constants_mini_extension[i%len(constants_mini_extension)])
		c_of_p_evaluations[i] = s0.Mod(s0, f.modulus)
	}
	fmt.Printf("Computed C(P, K) polynomial c_of_p_evaluations[65534]: %s\n", c_of_p_evaluations[65534])

	// Compute D(x) = C(P(x), P(g1*x), K(x)) / Z(x) ;;;; Z(x) = (x^steps - 1) / (x - x_atlast_step)
	z_num_evaluations := make([]*big.Int, precision.Uint64())
	for i := int(0); i < int(precision.Int64()); i++ {
		t := new(big.Int).Set(xs[(int64(i)*steps.Int64())%precision.Int64()])
		t.Sub(t, big.NewInt(1))
		z_num_evaluations[i] = t
	}
	z_num_inv := f.multi_inv(z_num_evaluations)
	z_den_evaluations := make([]*big.Int, precision.Uint64())
	for i := int(0); i < int(precision.Int64()); i++ {
		z_den_evaluations[i] = new(big.Int).Sub(xs[i], last_step_position)
	}
	d_evaluations := make([]*big.Int, precision.Uint64())
	for i := int(0); i < int(precision.Int64()); i++ {
		cp := c_of_p_evaluations[i]
		zd := z_den_evaluations[i]
		zni := z_num_inv[i]
		t := new(big.Int).Mul(cp, zd)
		t.Mul(t, zni)
		t.Mod(t, f.modulus)
		d_evaluations[i] = t // [cp * zd * zni % modulus for cp, zd, zni in zip(c_of_p_evaluations, z_den_evaluations, z_num_inv)]
	}
	fmt.Printf("Computed D polynomial\n")

	// Compute interpolant of ((1, input), (x_atlast_step, output))
	xsa := make([]*big.Int, 2)
	xsa[0] = big.NewInt(1)
	xsa[1] = new(big.Int).Set(last_step_position)
	ysa := make([]*big.Int, 2)
	ysa[0] = new(big.Int).Set(inp)
	ysa[1] = new(big.Int).Set(output)
	interpolant := f.lagrange_interp_2(xsa, ysa)
	i_evaluations := make([]*big.Int, len(xs))
	for i, x := range xs {
		i_evaluations[i] = f.eval_poly_at(interpolant, x)
	}

	xsa[0] = big.NewInt(-1)
	xsa[1] = big.NewInt(+1)
	ysa[0] = new(big.Int).Mul(last_step_position, big.NewInt(-1))
	ysa[1] = big.NewInt(+1)

	zeropoly2 := f.mul_polys(xsa, ysa)
	inpa := make([]*big.Int, len(xs))
	for i, x := range xs {
		inpa[i] = f.eval_poly_at(zeropoly2, x)
	}
	inv_z2_evaluations := f.multi_inv(inpa)

	b_evaluations := make([]*big.Int, len(inv_z2_evaluations))
	for j := 0; j < len(inv_z2_evaluations); j++ {
		b := new(big.Int).Sub(p_evaluations[j], i_evaluations[j])
		b.Mul(b, inv_z2_evaluations[j])
		b_evaluations[j] = b.Mod(b, f.modulus)
	}

	// Compute their Merkle root
	minp := make([][]byte, 0)
	for j := 0; j < len(p_evaluations); j++ {
		b := append(append(common.BigToHash(p_evaluations[j]).Bytes(), common.BigToHash(d_evaluations[j]).Bytes()...),
			common.BigToHash(b_evaluations[j]).Bytes()...)
		minp = append(minp, b)
	}
	mtree := merkelize(minp)
	fmt.Printf("Computed hash root %x\n", mtree[1])

	// Based on the hashes of P, D and B, we select a random linear combination
	// of P * x^steps, P, B * x^steps, B and D, and prove the low-degreeness of that,
	// instead of proving the low-degreeness of P, B and D separately
	k1 := blake(append(mtree[1][0:32], byte(0x1)))
	k2 := blake(append(mtree[1][0:32], byte(0x2)))
	k3 := blake(append(mtree[1][0:32], byte(0x3)))
	k4 := blake(append(mtree[1][0:32], byte(0x4)))

	// Compute the linear combination. We dont even both calculating it in coefficient form; we just compute the evaluations
	G2_to_the_steps := f.pow(G2, steps)

	powers := make([]*big.Int, 0)
	powers = append(powers, big.NewInt(1))
	for i := int64(1); i < precision.Int64(); i++ {
		t := new(big.Int).Mul(powers[len(powers)-1], G2_to_the_steps)
		t.Mod(t, f.modulus)
		powers = append(powers, t)
	}
	l_evaluations := make([][]byte, precision.Int64())
	for i := int64(0); i < precision.Int64(); i++ {
		t0 := new(big.Int).Set(d_evaluations[i])
		t1 := new(big.Int).Mul(p_evaluations[i], k1)
		t2 := new(big.Int).Mul(p_evaluations[i], k2)
		t2.Mul(t2, powers[i])
		t3 := new(big.Int).Mul(b_evaluations[i], k3)
		t4 := new(big.Int).Mul(b_evaluations[i], k4)
		t4.Mul(t4, powers[i])
		t0.Add(t0, t1)
		t0.Add(t0, t2)
		t0.Add(t0, t3)
		t0.Add(t0, t4)
		t0.Mod(t0, f.modulus)
		l_evaluations[i] = common.BigToHash(t0).Bytes()
	}
	l_mtree := merkelize(l_evaluations)
	fmt.Printf("Computed random linear combination %x\n", l_mtree[1])

	// Do some spot checks of the Merkle tree at pseudo-random coordinates, excluding multiples of `extension_factor`
	samples := (spot_check_security_factor)
	positions, err := get_pseudorandom_indices(f, l_mtree[1], precision, samples, int64(extension_factor))
	if err != nil {
		return resultProof, err
	}
	branches := make([][][]byte, 0)
	for _, pos := range positions {
		b1 := mk_branch(mtree, pos.Int64())
		branches = append(branches, b1)
		t := new(big.Int).Add(pos, skips)
		t.Mod(t, precision)
		b2 := mk_branch(mtree, t.Int64())
		branches = append(branches, b2)
		b3 := mk_branch(l_mtree, pos.Int64())
		branches = append(branches, b3)
	}
	fmt.Printf("Computed spot checks [percentage] samples %d\n", len(branches))

	// Return the Merkle roots of P and D, the spot check Merkle proofs and low-degree proofs of P and D
	var o Proof
	o.Root = mtree[1]     // Merkle Root
	o.LRoot = l_mtree[1]  // Merkle Root
	o.Branches = branches // check
	o.Child = make([]*FriComponent, 0)
	err = prove_low_degree(&o, f, l_evaluations, G2, steps.Int64()*2, ef)
	return &o, nil
}

/*
 Generate an FRI proof that the polynomial that has the specified values at successive powers of the specified root of unity has a degree lower than maxdeg_plus_1
 We use maxdeg+1 instead of maxdeg because it's more mathematically convenient in this case.
*/
func prove_low_degree(p *Proof, f *PrimeField, values [][]byte, root_of_unity *big.Int, maxdeg_plus_1 int64, exclude_multiples_of *big.Int) (err error) {
	fmt.Printf("Proving %d values are degree <= %d\n", len(values), maxdeg_plus_1)

	// If the degree we are checking for is less than or equal to 32, use the polynomial directly as a proof
	if maxdeg_plus_1 <= 16 {
		// return [[x.to_bytes(32, 'big') for x in values]]
		var o FriComponent
		fmt.Printf("Produced FRI proof\n")
		o.Values = values
		p.Child = append(p.Child, &o)
		return nil
	}

	// Calculate the set of x coordinates
	xs := get_power_cycle(root_of_unity, f.modulus)
	if len(values) != len(xs) {
		return fmt.Errorf("Incorrect input length of xs")
	}

	// Put the values into a Merkle tree. This is the root that the proof will be checked against
	m := merkelize(values)

	// Select a pseudo-random x coordinate
	m1 := new(big.Int).SetBytes(m[1])
	m1.Mod(m1, f.modulus)
	special_x := m1

	// Calculate the "column" at that x coordinate
	// (see https://vitalik.ca/general/2017/11/22/starks_part_2.html)
	// We calculate the column by Lagrange-interpolating each row, and not
	// directly from the polynomial, as this is more efficient
	quarter_len := len(xs) / 4
	xsets := make([][]*big.Int, (quarter_len))
	ysets := make([][]*big.Int, (quarter_len))
	for i := 0; i < quarter_len; i++ {
		xsets[i] = make([]*big.Int, 4)
		ysets[i] = make([]*big.Int, 4)
		for j := 0; j < 4; j++ {
			xsets[i][j] = xs[i+quarter_len*j]
			ysets[i][j] = common.BytesToHash(values[i+quarter_len*j]).Big()
		}
	}
	x_polys := f.multi_interp_4(xsets, ysets)

	column := make([][]byte, len(x_polys))
	for i, p := range x_polys {
		c := f.eval_quartic(p, special_x)
		column[i] = common.BytesToHash(c.Bytes()).Bytes()
	}
	m2 := merkelize(column)

	// Pseudo-randomly select y indices to sample
	ys, err := get_pseudorandom_indices(f, m2[1], big.NewInt(int64(len(column))), int(40), int64(exclude_multiples_of.Int64()))
	if err != nil {
		return err
	}

	// Compute the Merkle branches for the values in the polynomial and the column
	branches := make([][][][]byte, 0)
	for _, y := range ys {
		b := make([][][]byte, 0)
		b = append(b, mk_branch(m2, y.Int64()))
		for j := 0; j < 4; j++ {
			idx := new(big.Int).Add(y, big.NewInt(int64(j*len(xs)/4)))
			b = append(b, mk_branch(m, idx.Int64()))
		}
		branches = append(branches, b)
	}

	// This component of the proof
	var o FriComponent
	o.Root = m2[1]
	o.Branches = branches
	p.Child = append(p.Child, &o)
	err = prove_low_degree(p, f, column, f.pow(root_of_unity, big.NewInt(4)), maxdeg_plus_1/4, exclude_multiples_of)
	if err != nil {
		return err
	}
	return nil
}
