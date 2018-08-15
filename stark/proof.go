package stark

import (
	"fmt"
	"math/big"
	"sync"
	"time"
)

const (
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
	max_modulus := new(big.Int).Exp(TWO, TWENTYFOUR, nil)
	if modulus.Cmp(max_modulus) > 0 {
		return o, fmt.Errorf("modulus must be less than 2**24")
	}
	data := seed
	for len(data) < count*4 {
		data = append(data, blakeb(data[len(data)-32:len(data)])...)
	}

	var wg sync.WaitGroup
	o = make([]*big.Int, count)
	njmp := min_iterations(count, 3)
	if exclude_multiples_of == 0 {
		wg.Add(1)
		for j := 0; j < count; j += njmp {
			go func(i0 int, i1 int) {
				if i1 > count {
					i1 = count
				}
				t := new(big.Int)
				for i := i0; i < i1; i++ {
					o[i] = t.Mod(BytesToBig(data[i*4:i*4+4]), f.modulus)
				}
				wg.Done()
			}(j, j+njmp)
		}
	} else {
		wg.Add(1)
		exclude_multiples_of_big := big.NewInt(exclude_multiples_of)
		real_modulus := new(big.Int).Mul(modulus, big.NewInt(int64(exclude_multiples_of-1)))
		real_modulus.Div(real_modulus, exclude_multiples_of_big)
		exclude_multiples_of_minus_1 := new(big.Int).Sub(exclude_multiples_of_big, ONE)
		for j := 0; j < count; j += njmp {
			go func(i0 int, i1 int) {
				if i1 > count {
					i1 = count
				}
				t0 := new(big.Int)
				t1 := new(big.Int)
				t := new(big.Int)
				for i := i0; i < i1; i++ {
					t1.Mod(BytesToBig(data[i*4:i*4+4]), real_modulus)
					o[i] = new(big.Int).Add(t0.Add(t1, ONE), t.Div(t1, exclude_multiples_of_minus_1))
				}
				wg.Done()
			}(j, j+njmp)
		}
		wg.Wait()
	}
	return o, nil
}

func interpolant_inpa(ch_i_evaluations chan IndexedInt, ch_inpa chan IndexedInt, _xs []*big.Int, interpolant []*big.Int, zeropoly2 []*big.Int, i0 int) {
	f, _ := NewPrimeField(nil)
	for _i, xsi := range _xs {
		i := i0 + _i
		ch_i_evaluations <- IndexedInt{idx: i, data: f.eval_poly_at(interpolant, xsi)}
		ch_inpa <- IndexedInt{idx: i, data: f.eval_poly_at(zeropoly2, xsi)}
	}
}

func constrained_poly(ch_c_of_p_evaluations chan IndexedInt, ch_z_num_evaluations chan IndexedInt, ch_z_den_evaluations chan IndexedInt,
	p_evaluations []*big.Int, xs []*big.Int, constants_mini_extension []*big.Int, last_step_position *big.Int, steps int64, precision int64, i0 int, i1 int) {
	s2 := new(big.Int)
	s1 := new(big.Int)
	t := new(big.Int)
	f, _ := NewPrimeField(nil)
	for i := i0; i < i1; i++ {
		idx := int64(i+extension_factor) % precision
		s1.Mod(s2.Mul(t.Mul(p_evaluations[i], p_evaluations[i]), p_evaluations[i]), f.modulus)
		ch_c_of_p_evaluations <- IndexedInt{idx: i, data: new(big.Int).Mod(s1.Sub(s2.Sub(p_evaluations[idx], s1), constants_mini_extension[i%len(constants_mini_extension)]), f.modulus)}
		ch_z_num_evaluations <- IndexedInt{idx: i, data: new(big.Int).Sub(xs[(int64(i)*steps)%precision], ONE)}
		ch_z_den_evaluations <- IndexedInt{idx: i, data: new(big.Int).Sub(xs[i], last_step_position)}
	}
}

// used in (6b)
func compute_d_evaluations(ch_d_evaluations chan IndexedInt, c_of_p_evaluations []*big.Int, z_den_evaluations []*big.Int, z_num_inv []*big.Int, i0 int) {
	f, _ := NewPrimeField(nil)
	t0 := new(big.Int)
	t1 := new(big.Int)
	for _i, c_of_p_evaluation := range c_of_p_evaluations {
		// [cp * zd * zni % modulus for cp, zd, zni in zip(c_of_p_evaluations, z_den_evaluations, z_num_inv)]
		ch_d_evaluations <- IndexedInt{idx: i0 + _i, data: new(big.Int).Mod(t1.Mul(t0.Mul(c_of_p_evaluation, z_den_evaluations[_i]), z_num_inv[_i]), f.modulus)}
	}
}

// used in (7)
func compute_b_evaluations(ch_b_evaluations chan IndexedInt, ch_merkle_inp chan IndexedBytes, p_evaluations []*big.Int, d_evaluations []*big.Int, i_evaluations []*big.Int, inv_z2_evaluations []*big.Int, i0 int) {
	b := new(big.Int)
	t := new(big.Int)
	f, _ := NewPrimeField(nil)
	for _i, p_evaluation := range p_evaluations {
		i := i0 + _i
		b_evaluation := new(big.Int).Mod(t.Mul(b.Sub(p_evaluation, i_evaluations[_i]), inv_z2_evaluations[_i]), f.modulus)
		ch_b_evaluations <- IndexedInt{idx: i, data: b_evaluation}
		ch_merkle_inp <- IndexedBytes{idx: i, data: append(append(BigToBytes(p_evaluation), BigToBytes(d_evaluations[_i])...), BigToBytes(b_evaluation)...)}
	}
}

// used in (9)
func compute_random_linear_combination(ch_l_evaluations chan IndexedBytes, b_evaluations []*big.Int, p_evaluations []*big.Int, d_evaluations []*big.Int, powers []*big.Int, k1 *big.Int, k2 *big.Int, k3 *big.Int, k4 *big.Int, i0 int) {
	f, _ := NewPrimeField(nil)
	t0 := new(big.Int)
	t1 := new(big.Int)
	t2 := new(big.Int)
	t3 := new(big.Int)
	for _i, b_evaluation := range b_evaluations {
		i := i0 + _i
		/* l_evaluations = [d_evaluations[i] +
		p_evaluations[i] * (k1 + k2 * powers[i]) +
		b_evaluations[i] * (k3 + powers[i] * k4)] % modulus
		for i in range(precision)]*/
		// t2 = p_evaluations[i] * (k1 + k2 * powers[i] )
		t2.Mul(t1.Add(t0.Mul(powers[i], k2), k1), p_evaluations[_i])
		// t3 = b_evaluations[i] * (k3 + powers[i] * k4)
		t3.Mul(b_evaluation, t0.Add(t1.Mul(powers[i], k4), k3))
		// t0 = d_evaluations[i] + t2 above + t3 above
		ch_l_evaluations <- IndexedBytes{idx: i, data: BigToBytes(t1.Mod(t0.Add(d_evaluations[_i], t1.Add(t2, t3)), f.modulus))}
	}
}

// used in (12)
func compute_branches(ch_branches chan IndexedBranchSet, mtree [][]byte, l_mtree [][]byte, positions []*big.Int, skips *big.Int, precision *big.Int, i0 int) {
	t0 := new(big.Int)
	t1 := new(big.Int)
	for _i, pos := range positions {
		i := i0 + _i
		branchset := make([][][]byte, 3)
		branchset[0] = mk_branch(mtree, pos.Int64())
		branchset[1] = mk_branch(mtree, t1.Mod(t0.Add(pos, skips), precision).Int64())
		branchset[2] = mk_branch(l_mtree, pos.Int64())
		ch_branches <- IndexedBranchSet{idx: i, data: branchset}
	}
}

func compute_column(ch_column chan IndexedBytes, x_polys [][]*big.Int, special_x *big.Int, i0 int) {
	f, _ := NewPrimeField(nil)
	for _i, x_poly := range x_polys {
		ch_column <- IndexedBytes{idx: i0 + _i, data: BigToBytes(f.eval_quartic(x_poly, special_x))}
	}
}

// Generate a STARK for a MIMC calculation
func NewProof(f *PrimeField, inp *big.Int, steps *big.Int, round_constants []*big.Int) (resultProof *Proof, err error) {
	// (0) setup variables, check constraints
	start0 := time.Now()
	two_power_32 := new(big.Int).Exp(TWO, THIRTYTWO, nil)
	extension_factor_big := big.NewInt(int64(extension_factor))
	precision := new(big.Int).Mul(steps, extension_factor_big)
	nprecision := int(precision.Int64())
	max_steps := new(big.Int).Div(two_power_32, extension_factor_big)
	if steps.Cmp(max_steps) > 0 {
		return resultProof, fmt.Errorf("too many steps")
	}
	if !is_a_power_of_2(steps) {
		return resultProof, fmt.Errorf("steps must be power of 2")
	}
	round_constants_big := big.NewInt(int64(len(round_constants)))
	if !is_a_power_of_2(round_constants_big) {
		return resultProof, fmt.Errorf("len(round_constants) must be power of 2")
	}
	if steps.Cmp(round_constants_big) < 0 {
		return resultProof, fmt.Errorf("len(round_constants) too high")
	}

	// Root of unity such that x^precision=1
	modulus_minus_one := new(big.Int).Sub(f.modulus, ONE)
	G2 := f.pow(SEVEN, new(big.Int).Div(modulus_minus_one, precision))
	G2_to_the_steps := f.pow(G2, steps)
	powers := make([]*big.Int, nprecision+1)

	// Root of unity such that x^steps=1
	skips := new(big.Int).Div(precision, steps)
	G1 := f.pow(G2, skips)
	fmt.Printf("(0) Setup: %s\n", time.Since(start0))
	var vg sync.WaitGroup
	var xs []*big.Int
	var last_step_position *big.Int

	// (1) Powers of the higher-order root of unity
	// inputs: G2, modulus
	// outputs: xs, powers
	vg.Add(1)
	go func() {
		start := time.Now()
		xs = get_power_cycle(G2, f.modulus)
		last_step_position = xs[(int(steps.Int64())-1)*extension_factor]
		vg.Done()

		powers[0] = new(big.Int).Set(ONE)
		t0 := new(big.Int)
		for i := 0; i < nprecision; i++ {
			powers[i+1] = new(big.Int).Mod(t0.Mul(powers[i], G2_to_the_steps), f.modulus)
		}
		fmt.Printf("(1) Powers of the higher-order root of unity [%s => %s]\n", time.Since(start), time.Since(start0))
	}()

	// (2) Generate the computational trace, Interpolate the computational trace into a polynomial P, with each step along a successive power of G1
	// inputs: round_constants
	// outputs: computational_trace, computational_trace_polynomial, p_evaluations, output
	computational_trace := make([]*big.Int, steps.Int64())
	var computational_trace_polynomial []*big.Int

	skips2 := new(big.Int).Div(steps, round_constants_big)
	var output *big.Int
	var p_evaluations []*big.Int
	vg.Add(1)
	go func() {
		start := time.Now()
		computational_trace[0] = inp
		t0 := new(big.Int)
		t1 := new(big.Int)
		for i := 1; i < int(steps.Int64()); i++ {
			computational_trace[i] = new(big.Int).Mod(t0.Add(t1.Exp(computational_trace[i-1], THREE, nil), round_constants[(i-1)%len(round_constants)]), f.modulus)
		}
		output = computational_trace[len(computational_trace)-1]
		fmt.Printf("(2a) Computational trace output (%d steps) [%s => %s]\n", steps, time.Since(start), time.Since(start0))

		start = time.Now()
		computational_trace_polynomial = f.fft(computational_trace, G1, true)
		fmt.Printf("(2b) Converted computational steps into a polynomial [%s => %s]\n", time.Since(start), time.Since(start0))

		start = time.Now()
		p_evaluations = f.fft(computational_trace_polynomial, G2, false)
		fmt.Printf("(2c) Extended it into p_evaluations [%s => %s]\n", time.Since(start), time.Since(start0))
		vg.Done()
	}()

	// (3) Convert round constants into a polynomial
	// inputs: round_constants
	//  outputs: constants_mini_polynomial, constants_mini_extension
	var constants_mini_polynomial []*big.Int
	var constants_mini_extension []*big.Int
	vg.Add(1)
	go func() {
		start := time.Now()
		constants_mini_polynomial = f.fft(round_constants, f.pow(G1, skips2), true)
		constants_mini_extension = f.fft(constants_mini_polynomial, f.pow(G2, skips2), false)
		fmt.Printf("(3) Converted round constants into a polynomial and low-degree extended it [%s => %s]\n", time.Since(start), time.Since(start0))
		vg.Done()
	}()

	// waiting for (1)+(2)+(3)
	vg.Wait()

	// (4) Compute interpolant of ((1, input), (x_atlast_step, output))
	// input: xs from (1), output from (2)
	// output: i_evaluations, inv_z2_evaluations
	var wg sync.WaitGroup
	start4 := time.Now()
	zeropoly2 := f.mul_polys([]*big.Int{NEGONE, ONE}, []*big.Int{new(big.Int).Neg(last_step_position), ONE})
	interpolant := f.lagrange_interp_2([]*big.Int{ONE, last_step_position}, []*big.Int{inp, output})
	nxs := len(xs)
	njmp := min_iterations(nxs, 3)
	i_evaluations := make([]*big.Int, nxs)
	inpa := make([]*big.Int, nxs)
	ch_i_evaluations := make(chan IndexedInt, nxs)
	ch_inpa := make(chan IndexedInt, nxs)
	for i := 0; i < len(xs); i += njmp {
		i1 := i + njmp
		if i1 > nxs {
			i1 = nxs
		}
		go interpolant_inpa(ch_i_evaluations, ch_inpa, xs[i:i1], interpolant, zeropoly2, i)
	}
	var inv_z2_evaluations []*big.Int
	wg.Add(1)
	go func() {
		for i := 0; i < nxs; i++ {
			i_eval := <-ch_i_evaluations
			i_evaluations[i_eval.idx] = i_eval.data
			a := <-ch_inpa
			inpa[a.idx] = a.data
		}
		fmt.Printf("(4a) Computed i_evaluations [%s]\n", time.Since(start4))

		start4 = time.Now()
		inv_z2_evaluations = f.multi_inv(inpa)
		fmt.Printf("(4b) Computed inv_z2_evaluations [%s => %s]\n", time.Since(start4), time.Since(start0))
		wg.Done()
	}()

	// (5) Create  composed polynomial such that C(P(x), P(g1*x), K(x)) = P(g1*x) - P(x)**3 - K(x)
	//  inputs: p_evaluations from (2), constants_mini_extension from (3)
	//  outputs: c_of_p_evaluations, z_num_evaluations, z_den_evaluations
	dprecision := min_iterations(nprecision, 3)
	c_of_p_evaluations := make([]*big.Int, nprecision)
	z_num_evaluations := make([]*big.Int, nprecision)
	z_den_evaluations := make([]*big.Int, nprecision)

	ch_c_of_p_evaluations := make(chan IndexedInt, nprecision)
	ch_z_num_evaluations := make(chan IndexedInt, nprecision)
	ch_z_den_evaluations := make(chan IndexedInt, nprecision)
	wg.Add(1)
	go func() {
		start := time.Now()
		for i := int(0); i < nprecision; i += dprecision {
			i1 := i + dprecision
			if i1 > nprecision {
				i1 = nprecision
			}
			go constrained_poly(ch_c_of_p_evaluations, ch_z_num_evaluations, ch_z_den_evaluations, p_evaluations, xs, constants_mini_extension, last_step_position, steps.Int64(), precision.Int64(), i, i1)
		}
		for i := 0; i < nprecision; i++ {
			d := <-ch_c_of_p_evaluations
			c_of_p_evaluations[d.idx] = d.data
			d = <-ch_z_num_evaluations
			z_num_evaluations[d.idx] = d.data
			d = <-ch_z_den_evaluations
			z_den_evaluations[d.idx] = d.data
		}
		fmt.Printf("(5) Computed C(P, K) polynomial [%s => %s]\n", time.Since(start), time.Since(start0))
		wg.Done()
	}()

	// waiting for (4)+(5)
	wg.Wait()

	// (6) Compute D(x) = C(P(x), P(g1*x), K(x)) / Z(x) ;;;; Z(x) = (x^steps - 1) / (x - x_atlast_step)
	// inputs: z_num_evaluations from (5)
	// outputs: z_num_inv
	start := time.Now()
	z_num_inv := f.multi_inv(z_num_evaluations)
	fmt.Printf("(6a) Computed z_num_inv [%s => %s]\n", time.Since(start), time.Since(start0))

	// (6b) Compute d_evaluations
	// inputs: c_of_p_evaluations, z_den_evaluations, z_num_inv; output: d_evaluations
	d_evaluations := make([]*big.Int, nprecision)
	ch_d_evaluations := make(chan IndexedInt, nprecision)
	start = time.Now()
	for i := int(0); i < nprecision; i += dprecision {
		i1 := i + dprecision
		if i1 > nprecision {
			i1 = nprecision
		}
		go compute_d_evaluations(ch_d_evaluations, c_of_p_evaluations[i:i1], z_den_evaluations[i:i1], z_num_inv[i:i1], i)
	}
	for i := 0; i < nprecision; i++ {
		d := <-ch_d_evaluations
		d_evaluations[d.idx] = d.data
	}
	fmt.Printf("(6b) Computed d_evaluations [%s => %s]\n", time.Since(start), time.Since(start0))

	// (7) Compute b_evaluations and Merkle Root
	// input: inv_z2_evaluations from (4b); p_evaluations from (2), d_evaluations from (6b)
	// output: merkle_inp, b_evaluations
	start = time.Now()
	nev := len(inv_z2_evaluations)
	merkle_inp := make([][]byte, len(p_evaluations))
	njmp = min_iterations(nev, 2)
	b_evaluations := make([]*big.Int, nev)
	ch_merkle_inp := make(chan IndexedBytes, nev)
	ch_b_evaluations := make(chan IndexedInt, nev)
	for i := 0; i < nev; i += njmp {
		i1 := i + njmp
		if i1 > nev {
			i1 = nev
		}
		go compute_b_evaluations(ch_b_evaluations, ch_merkle_inp, p_evaluations[i:i1], d_evaluations[i:i1], i_evaluations[i:i1], inv_z2_evaluations[i:i1], i)
	}

	for i := 0; i < nev; i++ {
		d := <-ch_b_evaluations
		b_evaluations[d.idx] = d.data
		d0 := <-ch_merkle_inp
		merkle_inp[d0.idx] = d0.data
	}
	fmt.Printf("(7) Computed b_evaluations, MerkleTree input [%s => %s] \n", time.Since(start), time.Since(start0))

	// (8) Compute merkle tree
	// input: merkle_inp from (7)
	// output: mtree[1], k1, k2, k3, k4
	start = time.Now()
	mtree := merkelize(merkle_inp)
	k1 := blake(append(mtree[1][0:32], byte(0x1)))
	k2 := blake(append(mtree[1][0:32], byte(0x2)))
	k3 := blake(append(mtree[1][0:32], byte(0x3)))
	k4 := blake(append(mtree[1][0:32], byte(0x4)))
	fmt.Printf("(8) Compute mtree[1] %x [%s => %s]\n", mtree[1], time.Since(start), time.Since(start0))

	// (9) Based on the hashes of P, D and B, we select a random linear combination of P * x^steps, P, B * x^steps, B and D, and prove the low-degreeness of that, instead of proving the low-degreeness of P, B and D separately
	// input: d_evaluations from (6), p_evaluations from (2), b_evaluations from (8) + k1, k2, k3, k4 from (9)
	// output: l_evaluations
	// Compute the linear combination. We dont even both calculating it in coefficient form; we just compute the evaluations
	start = time.Now()
	njmp = min_iterations(nprecision, 4)
	l_evaluations := make([][]byte, precision.Int64())
	ch_l_evaluations := make(chan IndexedBytes, nprecision)
	for i := 0; i < nprecision; i += njmp {
		i1 := i + njmp
		if i1 > nprecision {
			i1 = nprecision
		}
		go compute_random_linear_combination(ch_l_evaluations, b_evaluations[i:i1], p_evaluations[i:i1], d_evaluations[i:i1], powers, k1, k2, k3, k4, i)
	}

	for i := 0; i < nprecision; i++ {
		d := <-ch_l_evaluations
		l_evaluations[d.idx] = d.data
	}
	fmt.Printf("(9) Computed random linear combination [%s => %s]\n", time.Since(start), time.Since(start0))

	// (10) Do some spot checks of the Merkle tree at pseudo-random coordinates, excluding multiples of `extension_factor`
	//      Setup proof structure with the Merkle roots of P and D
	// input: l_evaluations from (10)
	// output: l_mtree/p.LRoot, positions
	var p Proof
	p.Root = mtree[1] // Merkle Root
	p.Child = make([]*FriComponent, 0)
	var l_mtree [][]byte
	var positions []*big.Int
	var err2 error
	wg.Add(1)
	go func() {
		l_mtree = merkelize(l_evaluations)
		positions, err2 = get_pseudorandom_indices(f, l_mtree[1], precision, spot_check_security_factor, int64(extension_factor))
		p.LRoot = l_mtree[1] // Merkle Root
		wg.Done()
		fmt.Printf("(10) Merkelized l_evaluations, Setup Spot check positions [%s => %s]\n", time.Since(start), time.Since(start0))
	}()

	// (11) Recursive prove_low_degree proofs!
	// input: l_evaluations
	// output: p.Child[:]
	wg.Add(1)
	go func() {
		start := time.Now()
		err = prove_low_degree(&p, f, l_evaluations, G2, steps.Int64()*2, extension_factor_big)
		wg.Done()
		fmt.Printf("(11) Finished prove_low_degree [%s => %s]\n", time.Since(start), time.Since(start0))
	}()

	// waiting for (10)+(11)
	wg.Wait()

	// (12) Setup proof structure wit the spot check Merkle proofs and low-degree proofs of P and D
	// input: l_mtree+positions from (11), mtree from (10)
	// output: p.Branches
	start = time.Now()
	nev = len(positions)
	njmp = min_iterations(nev, 3)
	branches := make([][][]byte, len(positions)*3)
	ch_branches := make(chan IndexedBranchSet, nev)
	for i := 0; i < nev; i += njmp {
		i1 := i + njmp
		if i1 > nev {
			i1 = nev
		}
		go compute_branches(ch_branches, mtree, l_mtree, positions[i:i1], skips, precision, i)
	}
	for i := 0; i < nev; i++ {
		d := <-ch_branches
		branches[d.idx*3+0] = d.data[0]
		branches[d.idx*3+1] = d.data[1]
		branches[d.idx*3+2] = d.data[2]
	}
	p.Branches = branches
	fmt.Printf("(12) Finalized branches [%s => %s]\n", time.Since(start), time.Since(start0))
	return &p, err
}

/*
 Generate an FRI proof that the polynomial that has the specified values at successive powers of the specified root of unity has a degree lower than maxdeg_plus_1
 We use maxdeg+1 instead of maxdeg because it's more mathematically convenient in this case.
*/
func prove_low_degree(p *Proof, f *PrimeField, values [][]byte, root_of_unity *big.Int, maxdeg_plus_1 int64, exclude_multiples_of *big.Int) (err error) {
	var wg sync.WaitGroup

	// If the degree we are checking for is less than or equal to 32, use the polynomial directly as a proof
	if maxdeg_plus_1 <= 16 {
		// return [[x.to_bytes(32, 'big') for x in values]]
		var o FriComponent
		o.Values = values
		p.Child = append(p.Child, &o)
		return nil
	}

	// (1) Calculate the set of x coordinates
	// input: root_of_unity, values
	// output: xs, x_polys
	var xs []*big.Int
	var xsets [][]*big.Int
	var ysets [][]*big.Int
	var x_polys [][]*big.Int
	wg.Add(1)
	go func() {
		start := time.Now()
		xs = get_power_cycle(root_of_unity, f.modulus)
		fmt.Printf("  (%d-1a) Calculate the set of x coordinates [%s]\n", maxdeg_plus_1, time.Since(start))

		start = time.Now()
		quarter_len := len(xs) / 4
		xsets = make([][]*big.Int, quarter_len)
		ysets = make([][]*big.Int, quarter_len)
		for i := 0; i < quarter_len; i++ {
			xsets[i] = make([]*big.Int, 4)
			ysets[i] = make([]*big.Int, 4)
			for j := 0; j < 4; j++ {
				xsets[i][j] = xs[i+quarter_len*j]
				ysets[i][j] = BytesToBig(values[i+quarter_len*j])
			}
		}
		fmt.Printf("  (%d-1b) Setup xsets, ysets [%s]\n", maxdeg_plus_1, time.Since(start))

		start = time.Now()
		x_polys = f.multi_interp_4(xsets, ysets)
		fmt.Printf("  (%d-1c) Computed x_polys [%s]\n", maxdeg_plus_1, time.Since(start))
		wg.Done()
	}()

	// (2) Merkelize values - This is the root that the proof will be checked against
	// input: values from function input
	// output: special_x
	var special_x *big.Int
	var m [][]byte
	wg.Add(1)
	go func() {
		start := time.Now()
		m = merkelize(values)
		special_x = new(big.Int).Mod(new(big.Int).SetBytes(m[1]), f.modulus)
		wg.Done()
		fmt.Printf("  (%d-2) Merkelize values [%s]\n", maxdeg_plus_1, time.Since(start))
	}()

	// wait for (1)+(2)
	wg.Wait()
	if len(values) != len(xs) {
		return fmt.Errorf("Incorrect input length of xs")
	}

	// (3) Calculate the "column" at that x coordinate (see https://vitalik.ca/general/2017/11/22/starks_part_2.html)
	// We calculate the column by Lagrange-interpolating each row, and not directly from the polynomial, as this is more efficient
	// input: x_polys from (1), special_x from (2)
	// output: column
	start := time.Now()
	nev := len(x_polys)
	column := make([][]byte, nev)
	ch_column := make(chan IndexedBytes, nev)
	nj := min_iterations(nev, 4)
	for i := 0; i < nev; i += nj {
		i1 := i + nj
		if i1 > nev {
			i1 = nev
		}
		go compute_column(ch_column, x_polys[i:i1], special_x, i)
	}
	for i := 0; i < nev; i++ {
		d := <-ch_column
		column[d.idx] = d.data
	}
	fmt.Printf("  (%d-3) Computed column [%s]\n", maxdeg_plus_1, time.Since(start))

	// (4) Build NEXT component of the proof
	// input: column from (3)
	// output: p.Child.{Branches, Root}
	var component FriComponent
	p.Child = append(p.Child, &component)

	wg.Add(1)
	go func() {
		start := time.Now()
		err = prove_low_degree(p, f, column, f.pow(root_of_unity, FOUR), maxdeg_plus_1/4, exclude_multiples_of)
		wg.Done()
		fmt.Printf("  (%d-4) Computed prove_low_degree [%s]\n", maxdeg_plus_1, time.Since(start))
	}()
	if err != nil {
		return err
	}

	// (5) Merkelize column and Pseudo-randomly select y indices to sample
	// input: column from (3)
	// output: component.{Branches, Root}
	var err2 error
	var m2 [][]byte
	var ys []*big.Int
	branches := make([][][][]byte, 0)
	wg.Add(1)
	go func() {
		start := time.Now()
		m2 = merkelize(column)
		component.Root = m2[1]

		ys, err2 = get_pseudorandom_indices(f, m2[1], big.NewInt(int64(nev)), 40, int64(exclude_multiples_of.Int64()))
		for _, y := range ys {
			b := make([][][]byte, 5)
			b[0] = mk_branch(m2, y.Int64())
			for j := 0; j < 4; j++ {
				b[j+1] = mk_branch(m, y.Int64()+int64(j*len(xs)/4))
			}
			branches = append(branches, b)
		}
		wg.Done()
		component.Branches = branches
		fmt.Printf("  (%d-5) Computed branches [%s]\n", maxdeg_plus_1, time.Since(start))
	}()
	if err2 != nil {
		return err2
	}

	// waiting for (4)+(5)
	wg.Wait()
	return nil
}
