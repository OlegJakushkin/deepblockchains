package stark

import (
	"fmt"
	"math/big"
	"sync"
	"time"
)

const (
	NUM_CORES                  = 8
	MIN_SECONDS_BENCHMARK      = 0.0025
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

	var wg sync.WaitGroup
	o = make([]*big.Int, count)
	njmp := count / NUM_CORES
	if njmp < 500 {
		njmp = 500
	}
	if exclude_multiples_of == 0 {
		wg.Add(1)
		for j := 0; j < count; j += njmp {
			go func(i0 int, i1 int) {
				for i := 0; i < count*4; i += 4 {
					t := BytesToBig(data[i : i+4])
					o = append(o, t.Mod(t, f.modulus))
				}
				wg.Done()
			}(j, j+njmp)
		}
	} else {
		wg.Add(1)
		real_modulus := new(big.Int).Mul(modulus, big.NewInt(int64(exclude_multiples_of-1)))
		real_modulus.Div(real_modulus, big.NewInt(exclude_multiples_of))

		for j := 0; j < count; j += njmp {
			go func(i0 int, i1 int) {
				for i := i0; i < i1; i++ {
					if i < count {
						t := BytesToBig(data[i*4 : i*4+4])
						o[i] = t.Mod(t, real_modulus)
						s := new(big.Int).Add(o[i], big.NewInt(1))
						t = new(big.Int).Div(o[i], new(big.Int).Sub(big.NewInt(exclude_multiples_of), big.NewInt(1)))
						s.Add(s, t)
						o[i] = s
					}
				}
				wg.Done()
			}(j, j+njmp)
		}
		wg.Wait()
	}
	return o, nil
}

// Generate a STARK for a MIMC calculation
func NewProof(f *PrimeField, inp *big.Int, steps *big.Int, round_constants []*big.Int) (resultProof *Proof, err error) {
	// Some constraints to make our job easier
	start := time.Now()
	two_power_32 := new(big.Int).Exp(big.NewInt(2), big.NewInt(32), nil)
	ef := big.NewInt(int64(extension_factor))
	precision := new(big.Int).Mul(steps, ef)
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

	// Root of unity such that x^precision=1
	t := new(big.Int).Sub(f.modulus, big.NewInt(1))
	G2 := f.pow(big.NewInt(7), t.Div(t, precision))

	// Root of unity such that x^steps=1
	skips := new(big.Int).Div(precision, steps)
	G1 := f.pow(G2, skips)
	if time.Since(start).Seconds() > MIN_SECONDS_BENCHMARK {
		fmt.Printf("Setup: %s\n", time.Since(start))
	}
	var vg sync.WaitGroup
	var xs []*big.Int
	var last_step_position *big.Int

	// Powers of the higher-order root of unity
	// inputs: G2, modulus  outputs: xs
	vg.Add(1)
	go func() {
		xs = get_power_cycle(G2, f.modulus)
		last_step_position = xs[(int(steps.Int64())-1)*extension_factor]
		vg.Done()
	}()

	// Generate the computational trace
	// inputs: round_constants  outputs: computational_trace, output
	start = time.Now()
	computational_trace := make([]*big.Int, steps.Int64())
	computational_trace[0] = inp
	THREE := big.NewInt(3)
	for i := 1; i < int(steps.Int64()); i++ {
		t := new(big.Int).Exp(computational_trace[i-1], THREE, nil)
		t.Add(t, round_constants[(i-1)%len(round_constants)])
		computational_trace[i] = t.Mod(t, f.modulus)
	}
	output := computational_trace[len(computational_trace)-1]
	fmt.Printf("Computational trace output (%d steps) [%s]\n", steps, time.Since(start))

	// Interpolate the computational trace into a polynomial P, with each step along a successive power of G1
	// inputs: computational_trace  outputs: computational_trace_polynomial, p_evaluations
	start = time.Now()
	computational_trace_polynomial := f.fft(computational_trace, G1, true)
	fmt.Printf("Converted computational steps into a polynomial [%s]\n", time.Since(start))
	var p_evaluations []*big.Int
	vg.Add(1)
	go func() {
		start = time.Now()
		p_evaluations = f.fft(computational_trace_polynomial, G2, false)
		vg.Done()
		fmt.Printf("Extended it [%s]\n", time.Since(start))
	}()
	vg.Wait()

	// Convert round constants into a polynomial
	// inputs: round_constants  outputs: constants_mini_polynomial,  constants_polynomial, constants_mini_extension
	start = time.Now()
	skips2 := new(big.Int).Div(steps, rc)
	constants_mini_polynomial := f.fft(round_constants, f.pow(G1, skips2), true)
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
	fmt.Printf("Converted round constants into a polynomial and low-degree extended it [%s]\n", time.Since(start))
	vg.Wait()

	// Create the composed polynomial such that C(P(x), P(g1*x), K(x)) = P(g1*x) - P(x)**3 - K(x)
	//  inputs: p_evaluations, constants_mini_extension outputs: c_of_p_evaluations, z_num_evaluations, z_den_evaluations
	start = time.Now()
	var wg sync.WaitGroup
	nprecision := int(precision.Int64())
	dprecision := nprecision / NUM_CORES
	if dprecision < 500 {
		dprecision = 500
	}
	c_of_p_evaluations := make([]*big.Int, nprecision)
	z_num_evaluations := make([]*big.Int, nprecision)
	z_den_evaluations := make([]*big.Int, nprecision)
	d_evaluations := make([]*big.Int, nprecision)
	for j := int(0); j < nprecision; j += dprecision {
		wg.Add(1)
		go func(i0 int, i1 int) {
			for i := i0; i < i1; i++ {
				if i < nprecision {
					idx := int64(i+extension_factor) % precision.Int64()
					s0 := new(big.Int).Set(p_evaluations[idx])
					t := new(big.Int).Mul(p_evaluations[i], p_evaluations[i])
					t.Mul(t, p_evaluations[i])
					t.Mod(t, f.modulus)
					s0.Sub(s0, t)
					s0.Sub(s0, constants_mini_extension[i%len(constants_mini_extension)])
					c_of_p_evaluations[i] = s0.Mod(s0, f.modulus)
					t = new(big.Int).Set(xs[(int64(i)*steps.Int64())%precision.Int64()])
					z_num_evaluations[i] = t.Sub(t, big.NewInt(1))
					z_den_evaluations[i] = new(big.Int).Sub(xs[i], last_step_position)
				}
			}
			wg.Done()
		}(j, j+dprecision)
	}
	wg.Wait()
	fmt.Printf("Computed C(P, K) polynomial [%s]\n", time.Since(start))

	// Compute D(x) = C(P(x), P(g1*x), K(x)) / Z(x) ;;;; Z(x) = (x^steps - 1) / (x - x_atlast_step)
	// inputs: z_num_evaluations  outputs: z_num_inv
	start = time.Now()
	z_num_inv := f.multi_inv(z_num_evaluations)
	fmt.Printf("Computed D polynomial [%s]\n", time.Since(start))

	start = time.Now()
	for j := int(0); j < nprecision; j += dprecision {
		wg.Add(1)
		go func(i0 int, i1 int) {
			for i := i0; i < i1; i++ {
				if i < nprecision {
					t := new(big.Int).Mul(c_of_p_evaluations[i], z_den_evaluations[i])
					t.Mul(t, z_num_inv[i])
					d_evaluations[i] = t.Mod(t, f.modulus) // [cp * zd * zni % modulus for cp, zd, zni in zip(c_of_p_evaluations, z_den_evaluations, z_num_inv)]
				}
			}
			wg.Done()
		}(j, j+dprecision)
	}
	wg.Wait()
	fmt.Printf("Computed d_evaluations [%s]\n", time.Since(start))

	// Compute interpolant of ((1, input), (x_atlast_step, output))
	start = time.Now()
	xsa := make([]*big.Int, 2)
	xsa[0] = big.NewInt(1)
	xsa[1] = new(big.Int).Set(last_step_position)
	ysa := make([]*big.Int, 2)
	ysa[0] = new(big.Int).Set(inp)
	ysa[1] = new(big.Int).Set(output)
	interpolant := f.lagrange_interp_2(xsa, ysa)
	nxs := len(xs)
	njmp := nxs / NUM_CORES
	if njmp < 500 {
		njmp = 500
	}
	i_evaluations := make([]*big.Int, nxs)
	for i := 0; i < len(xs); i += njmp {
		wg.Add(1)
		go func(i0 int, i1 int) {
			if i1 > nxs {
				i1 = nxs
			}
			for j := i0; j < i1; j++ {
				i_evaluations[j] = f.eval_poly_at(interpolant, xs[j])
			}
			wg.Done()
		}(i, i+njmp)
	}
	wg.Wait()
	fmt.Printf("Computed interpolant [%s]\n", time.Since(start))

	start = time.Now()
	xsa[0] = big.NewInt(-1)
	xsa[1] = big.NewInt(+1)
	ysa[0] = new(big.Int).Neg(last_step_position)
	ysa[1] = big.NewInt(+1)
	zeropoly2 := f.mul_polys(xsa, ysa)

	inpa := make([]*big.Int, len(xs))
	nev := len(xs)
	njmp = nev / NUM_CORES
	for j := 0; j < len(xs); j += njmp {
		wg.Add(1)
		go func(i0 int, i1 int) {
			// startb := time.Now()
			for i := i0; i < i1 && i < nev; i++ {
				inpa[i] = f.eval_poly_at(zeropoly2, xs[i])
			}
			wg.Done()
			//	fmt.Printf("  computed inpa %d-%d [%s] \n", i0, i1, time.Since(startb))
		}(j, j+njmp)
	}
	wg.Wait()

	inv_z2_evaluations := f.multi_inv(inpa)
	nev = len(inv_z2_evaluations)
	fmt.Printf("Computed inv_z2_evaluations [%s]\n", time.Since(start))

	start = time.Now()
	njmp = nev / NUM_CORES
	if njmp < 500 {
		njmp = 500
	}

	b_evaluations := make([]*big.Int, nev)
	for i := 0; i < nev; i += njmp {
		wg.Add(1)
		go func(i0 int, i1 int) {
			if i1 > nev {
				i1 = nev
			}
			for j := i0; j < i1; j++ {
				b := new(big.Int).Sub(p_evaluations[j], i_evaluations[j])
				b.Mul(b, inv_z2_evaluations[j])
				b_evaluations[j] = b.Mod(b, f.modulus)
			}
			wg.Done()
		}(i, i+njmp)
	}
	wg.Wait()
	fmt.Printf("Computed b_evaluations [%s] \n", time.Since(start))

	// Compute their Merkle root
	start = time.Now()
	nev = len(p_evaluations)
	minp := make([][]byte, len(p_evaluations))
	for j := 0; j < nev; j++ {
		minp[j] = append(append(BigToBytes(p_evaluations[j]), BigToBytes(d_evaluations[j])...), BigToBytes(b_evaluations[j])...)
	}
	mtree := merkelize(minp)
	fmt.Printf("Computed hash root [%s] mtree[1] %x\n", time.Since(start), mtree[1])
	// Based on the hashes of P, D and B, we select a random linear combination
	// of P * x^steps, P, B * x^steps, B and D, and prove the low-degreeness of that,
	// instead of proving the low-degreeness of P, B and D separately
	start = time.Now()
	k1 := blake(append(mtree[1][0:32], byte(0x1)))
	k2 := blake(append(mtree[1][0:32], byte(0x2)))
	k3 := blake(append(mtree[1][0:32], byte(0x3)))
	k4 := blake(append(mtree[1][0:32], byte(0x4)))

	// Compute the linear combination. We dont even both calculating it in coefficient form; we just compute the evaluations
	G2_to_the_steps := f.pow(G2, steps)
	powers := make([]*big.Int, nprecision+1)
	powers[0] = big.NewInt(1)
	for i := 0; i < nprecision; i++ {
		t := new(big.Int).Mul(powers[i], G2_to_the_steps)
		powers[i+1] = t.Mod(t, f.modulus)
	}
	njmp = nprecision / NUM_CORES
	if njmp < 500 {
		njmp = 500
	}
	l_evaluations := make([][]byte, precision.Int64())
	for j := int(0); j < nprecision; j += njmp {
		wg.Add(1)
		go func(i0 int, i1 int) {
			for i := i0; i < i1; i++ {
				t0 := new(big.Int).Set(d_evaluations[i])
				t2 := new(big.Int).Mul(p_evaluations[i], k2)
				t4 := new(big.Int).Mul(b_evaluations[i], k4)
				t0.Add(t0, new(big.Int).Mul(p_evaluations[i], k1))
				t0.Add(t0, t2.Mul(t2, powers[i]))
				t0.Add(t0, new(big.Int).Mul(b_evaluations[i], k3))
				t0.Add(t0, t4.Mul(t4, powers[i]))
				t0.Mod(t0, f.modulus)
				l_evaluations[i] = BigToBytes(t0)
			}
			wg.Done()
		}(j, j+njmp)
	}
	wg.Wait()

	// use the l_evaluations:
	var l_mtree [][]byte
	wg.Add(1)
	go func() {
		l_mtree = merkelize(l_evaluations)
		wg.Done()
	}()
	wg.Wait()
	fmt.Printf("Computed random linear combination %x [%s]\n", l_mtree[1], time.Since(start))

	// Do some spot checks of the Merkle tree at pseudo-random coordinates, excluding multiples of `extension_factor`
	start = time.Now()
	samples := spot_check_security_factor
	positions, err := get_pseudorandom_indices(f, l_mtree[1], precision, samples, int64(extension_factor))
	if err != nil {
		return resultProof, err
	}
	nev = len(positions)
	njmp = nev / NUM_CORES
	if njmp < 500 {
		njmp = 500
	}
	branches := make([][][]byte, len(positions)*3)
	for j := 0; j < nev; j += njmp {
		wg.Add(1)
		go func(i0 int, i1 int) {
			for i := i0; i < i1; i++ {
				if i < nev {
					pos := positions[i]
					branches[i*3+0] = mk_branch(mtree, pos.Int64())
					t := new(big.Int).Add(pos, skips)
					t.Mod(t, precision)
					branches[i*3+1] = mk_branch(mtree, t.Int64())
					branches[i*3+2] = mk_branch(l_mtree, pos.Int64())
				}
			}
			wg.Done()
		}(j, j+njmp)
	}
	// Return the Merkle roots of P and D, the spot check Merkle proofs and low-degree proofs of P and D
	var o Proof
	o.Root = mtree[1]    // Merkle Root
	o.LRoot = l_mtree[1] // Merkle Root
	o.Child = make([]*FriComponent, 0)

	err = prove_low_degree(&o, f, l_evaluations, G2, steps.Int64()*2, ef)
	wg.Wait()
	o.Branches = branches // check
	return &o, nil
}

/*
 Generate an FRI proof that the polynomial that has the specified values at successive powers of the specified root of unity has a degree lower than maxdeg_plus_1
 We use maxdeg+1 instead of maxdeg because it's more mathematically convenient in this case.
*/
func prove_low_degree(p *Proof, f *PrimeField, values [][]byte, root_of_unity *big.Int, maxdeg_plus_1 int64, exclude_multiples_of *big.Int) (err error) {
	start := time.Now()

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

	m1 := new(big.Int).SetBytes(m[1])
	special_x := new(big.Int).Mod(m1, f.modulus)

	// Calculate the "column" at that x coordinate
	// (see https://vitalik.ca/general/2017/11/22/starks_part_2.html)
	// We calculate the column by Lagrange-interpolating each row, and not
	// directly from the polynomial, as this is more efficient
	quarter_len := len(xs) / 4
	xsets := make([][]*big.Int, quarter_len)
	ysets := make([][]*big.Int, quarter_len)
	for i := 0; i < quarter_len; i++ {
		xsets[i] = make([]*big.Int, 4)
		ysets[i] = make([]*big.Int, 4)
		for j := 0; j < 4; j++ {
			xsets[i][j] = xs[i+quarter_len*j]
			ysets[i][j] = BytesToBig(values[i+quarter_len*j])
		}
	}

	x_polys := f.multi_interp_4(xsets, ysets)

	column := make([][]byte, len(x_polys))
	nev := len(x_polys)
	nj := nev / NUM_CORES
	if nj < 1024 {
		nj = 1024
	}
	var wg sync.WaitGroup
	for j := 0; j < len(x_polys); j += nj {
		wg.Add(1)
		go func(i0 int, i1 int) {
			if i1 > nev {
				i1 = nev
			}
			for i := i0; i < i1; i++ {
				column[i] = BigToBytes(f.eval_quartic(x_polys[i], special_x))
			}
			wg.Done()
		}(j, j+nj)
	}
	wg.Wait()
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
	fmt.Printf("Proving %d values are degree <= %d [%s] %x\n", len(values), maxdeg_plus_1, time.Since(start), m2[1])

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
