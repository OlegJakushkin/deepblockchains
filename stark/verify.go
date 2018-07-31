package stark

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Verifies a STARK, using verify_low_degree_proof for each component in the STARK proof
func verify_mimc_proof(f *PrimeField, inp *big.Int, steps *big.Int, round_constants []*big.Int, output *big.Int, proof *Proof) error {
	m_root := proof.Root
	l_root := proof.LRoot
	branches := proof.Branches

	ext_factor := big.NewInt(int64(extension_factor))
	precision := new(big.Int).Mul(steps, ext_factor)

	// Get (steps)th root of unity
	G2 := f.pow(big.NewInt(7), f.div(new(big.Int).Sub(f.modulus, big.NewInt(1)), precision))
	skips := f.div(precision, steps)

	// Gets the polynomial representing the round constants
	skips2 := f.div(steps, big.NewInt(int64(len(round_constants))))

	extf := new(big.Int).Mul(ext_factor, skips2)
	constants_mini_polynomial := f.fft(round_constants, f.pow(G2, extf), true)

	// Verifies the low-degree proofs
	err := verify_low_degree_proof(f, l_root, G2, proof.Child, int(steps.Int64()*2), ext_factor)
	if err != nil {
		return err
	}
	// Performs the spot checks
	k1 := (blake(append(m_root, byte(1))))
	k2 := (blake(append(m_root, byte(2))))
	k3 := (blake(append(m_root, byte(3))))
	k4 := (blake(append(m_root, byte(4))))

	samples := spot_check_security_factor
	positions, err := get_pseudorandom_indices(f, l_root, precision, samples, int64(extension_factor))
	if err != nil {
		return err
	}
	t := new(big.Int).Sub(steps, big.NewInt(1))
	t.Mul(t, skips)
	last_step_position := f.pow(G2, t)
	for i, pos := range positions {
		x := f.pow(G2, pos)
		x_to_the_steps := f.pow(x, steps)
		mbranch1, err := verify_branch(m_root, uint(pos.Uint64()), branches[i*3])
		if err != nil {
			return err
		}
		t := new(big.Int).Add(pos, skips)
		t.Mod(t, precision)
		mbranch2, err := verify_branch(m_root, uint(t.Uint64()), branches[i*3+1])
		l_of_x, err := verify_branch(l_root, uint(pos.Uint64()), branches[i*3+2])

		p_of_x := common.BytesToHash(mbranch1[:32]).Big()
		p_of_g1x := common.BytesToHash(mbranch2[:32]).Big()
		d_of_x := common.BytesToHash(mbranch1[32:64]).Big()
		b_of_x := common.BytesToHash(mbranch1[64:]).Big()

		zvalue := f.div(f.sub(f.pow(x, steps), big.NewInt(1)), f.sub(x, last_step_position))
		k_of_x := f.eval_poly_at(constants_mini_polynomial, f.pow(x, skips2))

		// Check transition constraints C(P(x)) = Z(x) * D(x)
		// (p_of_g1x - p_of_x ** 3 - k_of_x - zvalue * d_of_x) % modulus == 0
		s0 := new(big.Int).Sub(p_of_g1x, new(big.Int).Exp(p_of_x, big.NewInt(3), nil))
		s0.Sub(s0, k_of_x)
		s0.Sub(s0, new(big.Int).Mul(zvalue, d_of_x))
		s0.Mod(s0, f.modulus)
		if s0.Cmp(big.NewInt(0)) != 0 {
			return fmt.Errorf("transition constraint violation")
		}

		// Check boundary constraints B(x) * Q(x) + I(x) = P(x)
		a := create_poly2(big.NewInt(1), last_step_position)
		b := create_poly2(inp, output)
		interpolant := f.lagrange_interp_2(a, b)
		c := create_poly2(big.NewInt(-1), big.NewInt(1))
		d := create_poly2(new(big.Int).Mul(big.NewInt(-1), last_step_position), big.NewInt(1))
		zeropoly2 := f.mul_polys(c, d)
		// Check if (p_of_x - b_of_x * f.eval_poly_at(zeropoly2, x) - f.eval_poly_at(interpolant, x)) % modulus == 0
		s0.Sub(p_of_x, new(big.Int).Mul(b_of_x, f.eval_poly_at(zeropoly2, x)))
		s0.Sub(s0, f.eval_poly_at(interpolant, x))
		s0.Mod(s0, f.modulus)
		if s0.Cmp(big.NewInt(0)) != 0 {
			return fmt.Errorf("boundary constraint violation")
		}
		// Check correctness of the linear combination
		// Check if (l_of_x - d_of_x - k1 * p_of_x - k2 * p_of_x * x_to_the_steps - k3 * b_of_x - k4 * b_of_x * x_to_the_steps) % modulus == 0
		s0.Sub(common.BytesToHash(l_of_x).Big(), d_of_x)
		s0.Sub(s0, new(big.Int).Mul(k1, p_of_x))
		s0.Sub(s0, new(big.Int).Mul(k2, new(big.Int).Mul(p_of_x, x_to_the_steps)))
		s0.Sub(s0, new(big.Int).Mul(k3, b_of_x))
		s0.Sub(s0, new(big.Int).Mul(k4, new(big.Int).Mul(b_of_x, x_to_the_steps)))
		s0.Mod(s0, f.modulus)
		if s0.Cmp(big.NewInt(0)) != 0 {
			return fmt.Errorf("linear combination violation")
		}
	}
	fmt.Printf("Verified %d consistency checks\n", spot_check_security_factor)
	return nil
}

// Verify an FRI proof
func verify_low_degree_proof(f *PrimeField, merkle_root []byte, root_of_unity *big.Int, components []*FriComponent, maxdeg_plus_1 int, exclude_multiples_of *big.Int) error {
	// Calculate which root of unity we're working with
	testval := new(big.Int).Set(root_of_unity)
	roudeg := big.NewInt(1)
	for testval.Cmp(big.NewInt(1)) != 0 {
		roudeg.Mul(roudeg, big.NewInt(2))
		// testval = (testval * testval) % modulus
		testval.Mul(testval, testval)
		testval.Mod(testval, f.modulus)
	}

	// Powers of the given root of unity 1, p, p**2, p**3 such that p**4 = 1
	quartic_roots_of_unity := make([]*big.Int, 4)
	quartic_roots_of_unity[0] = big.NewInt(1)
	quartic_roots_of_unity[1] = f.pow(root_of_unity, new(big.Int).Div(roudeg, big.NewInt(4)))
	quartic_roots_of_unity[2] = f.pow(root_of_unity, new(big.Int).Div(roudeg, big.NewInt(2)))
	t := new(big.Int).Mul(roudeg, big.NewInt(3))
	t.Div(t, big.NewInt(4))
	quartic_roots_of_unity[3] = f.pow(root_of_unity, t)

	// Verify the recursive components of the proof
	for level, comp := range components[0 : len(components)-1] {
		fmt.Printf("Verifying degree (%d) <= %d\n", level, maxdeg_plus_1)
		root2 := comp.Root
		branches := comp.Branches

		// Calculate the pseudo-random x coordinate
		special_x := new(big.Int).SetBytes(merkle_root)
		special_x.Mod(special_x, f.modulus)

		// Calculate the pseudo-randomly sampled y indices
		modulus := new(big.Int).Div(roudeg, big.NewInt(4))
		ys, err := get_pseudorandom_indices(f, root2, modulus, 40, exclude_multiples_of.Int64())
		if err != nil {
			return err
		}
		// For each y coordinate, get the x coordinates on the row, the values on the row, and the value at that y from the column
		xcoords := make([][]*big.Int, 0)
		rows := make([][]*big.Int, 0)
		columnvals := make([]*big.Int, 0)
		for i, y := range ys {
			// The x coordinates from the polynomial
			x1 := f.pow(root_of_unity, y)
			a := make([]*big.Int, 4)
			for j := 0; j < 4; j++ {
				t := new(big.Int).Mul(quartic_roots_of_unity[j], x1)
				t.Mod(t, f.modulus)
				a[j] = t
			}
			xcoords = append(xcoords, a)

			// The values from the original polynomial
			row := make([]*big.Int, 4)
			for j := int64(0); j < 4; j++ {
				idx := y.Int64() + int64(roudeg.Int64()/4*j)
				r, err := verify_branch_int(merkle_root, uint(idx), branches[i][j+1])
				if err != nil {
					return err
				}
				row[j] = r
			}

			rows = append(rows, row)

			c, err := verify_branch_int(root2, uint(y.Int64()), branches[i][0])
			if err != nil {
				return err
			}
			columnvals = append(columnvals, c)
		}
		// Verify for each selected y coordinate that the four points from the
		// polynomial and the one point from the column that are on that y
		// coordinate are on the same deg < 4 polynomial
		polys := f.multi_interp_4(xcoords, rows)
		for j, p := range polys { // zip(polys, columnvals) {
			c := columnvals[j]
			q := f.eval_quartic(p, special_x)
			if q.Cmp(c) != 0 {
				return fmt.Errorf("quartic mismatch")
			}
		}
		// Update constants to check the next proof
		merkle_root = root2
		root_of_unity = f.pow(root_of_unity, big.NewInt(4))
		maxdeg_plus_1 = maxdeg_plus_1 / 4
		roudeg = f.div(roudeg, big.NewInt(4))
	}
	// Verify the direct components of the proof
	comp := components[len(components)-1]
	if maxdeg_plus_1 > 16 {
		return fmt.Errorf("max_degreeplus_1 too high")
	}

	fmt.Printf("Verifying degree <= %d\n", maxdeg_plus_1)
	// Check the Merkle root matches up
	mtree := merkelize(comp.Values)
	if bytes.Compare(mtree[1], merkle_root) != 0 {
		return fmt.Errorf("incorrect merkle root")
	}

	// Check the degree of the data
	powers := get_power_cycle(root_of_unity, f.modulus)
	pts := make([]int64, 0)
	if exclude_multiples_of.Cmp(big.NewInt(0)) > 0 {
		for i := int64(0); i < int64(len(comp.Values)); i++ {
			x := big.NewInt(i)
			x.Mod(x, exclude_multiples_of)
			if x.Cmp(big.NewInt(0)) > 0 {
				pts = append(pts, i)
			}
		}
	} else {
		for i := int64(0); i < int64(len(comp.Values)); i++ {
			pts = append(pts, i)
		}
	}

	// check points
	xs := make([]*big.Int, 0)
	ys := make([]*big.Int, 0)
	for _, x := range pts[:maxdeg_plus_1] {
		xs = append(xs, powers[x])
		ys = append(ys, common.BytesToHash(comp.Values[x]).Big())
	}

	poly := f.lagrange_interp(xs, ys)
	for _, x := range pts[maxdeg_plus_1:] {
		q := f.eval_poly_at(poly, powers[x])
		y := common.BytesToHash(comp.Values[x]).Big()
		if q.Cmp(y) != 0 {
			return fmt.Errorf("Lagrange mismatch")
		}
	}
	fmt.Printf("FRI proof verified (%d pts)\n", len(pts))

	return nil
}
