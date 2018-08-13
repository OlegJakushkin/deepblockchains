package stark

import (
	"bytes"
	"fmt"
	"math/big"
	"sync"
	"time"
)

// Verifies a STARK, using verify_low_degree_proof for each component in the STARK proof
func VerifyProof(f *PrimeField, inp *big.Int, steps *big.Int, round_constants []*big.Int, proof *Proof) error {
	// Verifies the low-degree proofs
	var wg sync.WaitGroup
	var output *big.Int
	wg.Add(1)
	go func() {
		output = f.MiMC(inp, steps, round_constants)
		wg.Done()
	}()

	m_root := proof.Root
	l_root := proof.LRoot
	branches := proof.Branches

	ext_factor := big.NewInt(int64(extension_factor))
	precision := new(big.Int).Mul(steps, ext_factor)

	// Get (steps)th root of unity
	G2 := f.pow(SEVEN, f.div(new(big.Int).Sub(f.modulus, ONE), precision))
	skips := f.div(precision, steps)
	errV := verify_low_degree_proof(f, l_root, G2, proof.Child, int(steps.Int64()*2), ext_factor)
	if errV != nil {
		return errV
	}

	// Gets the polynomial representing the round constants
	skips2 := f.div(steps, big.NewInt(int64(len(round_constants))))

	extf := new(big.Int).Mul(ext_factor, skips2)
	constants_mini_polynomial := f.fft(round_constants, f.pow(G2, extf), true)

	start := time.Now()
	// Performs the spot checks
	k1 := blake(append(m_root, byte(1)))
	k2 := blake(append(m_root, byte(2)))
	k3 := blake(append(m_root, byte(3)))
	k4 := blake(append(m_root, byte(4)))
	positions, err := get_pseudorandom_indices(f, l_root, precision, spot_check_security_factor, int64(extension_factor))
	if err != nil {
		return err
	}
	t := new(big.Int).Sub(steps, ONE)
	t.Mul(t, skips)
	last_step_position := f.pow(G2, t)

	nev := len(positions)
	njmp := min_iterations(nev, -1)

	// Wait for output from f.MiMC
	wg.Wait()
	var finalErr error
	finalErr = nil

	for j := 0; j < nev; j += njmp {
		wg.Add(1)
		go func(i0 int, i1 int) {
			if i1 > nev {
				i1 = nev
			}
			//	start0 := time.Now()
			s0 := new(big.Int)
			s1 := new(big.Int)
			s2 := new(big.Int)
			s3 := new(big.Int)
			p_of_x := new(big.Int)
			p_of_g1x := new(big.Int)
			d_of_x := new(big.Int)
			b_of_x := new(big.Int)

			for i := i0; i < i1; i++ {
				pos := positions[i]
				x := f.pow(G2, pos)
				x_to_the_steps := f.pow(x, steps)
				mbranch1, err := verify_branch(m_root, uint(pos.Uint64()), branches[i*3])
				if err != nil {
					finalErr = err
				}
				s0.Add(pos, skips)
				s1.Mod(s0, precision)
				mbranch2, err := verify_branch(m_root, uint(s1.Uint64()), branches[i*3+1])
				if err != nil {
					finalErr = err
				}
				l_of_x, err := verify_branch(l_root, uint(pos.Uint64()), branches[i*3+2])
				if err != nil {
					finalErr = err
				}

				p_of_x.SetBytes(mbranch1[:32])
				p_of_g1x.SetBytes(mbranch2[:32])
				d_of_x.SetBytes(mbranch1[32:64])
				b_of_x.SetBytes(mbranch1[64:])

				//zvalue := f.div(f.sub(f.pow(x, steps), one), f.sub(x, last_step_position))
				//	k_of_x := f.eval_poly_at(constants_mini_polynomial, f.pow(x, skips2))

				// Check transition constraints C(P(x)) = Z(x) * D(x)
				// (p_of_g1x - p_of_x ** 3 - k_of_x - zvalue * d_of_x) % modulus == 0
				s3.Exp(p_of_x, THREE, nil)
				s0.Sub(p_of_g1x, s3)
				s1.Sub(s0, f.eval_poly_at(constants_mini_polynomial, f.pow(x, skips2)))
				s0.Mul(f.div(f.sub(f.pow(x, steps), ONE), f.sub(x, last_step_position)), d_of_x)
				s2.Sub(s1, s0)
				if s3.Mod(s2, f.modulus).Cmp(ZERO) != 0 {
					finalErr = fmt.Errorf("transition constraint violation")
				}

				// Check boundary constraints B(x) * Q(x) + I(x) = P(x)
				interpolant := f.lagrange_interp_2(create_poly2(ONE, last_step_position), create_poly2(inp, output))
				zeropoly2 := f.mul_polys(create_poly2(NEGONE, ONE), create_poly2(new(big.Int).Mul(NEGONE, last_step_position), ONE))
				// Check if (p_of_x - b_of_x * f.eval_poly_at(zeropoly2, x) - f.eval_poly_at(interpolant, x)) % modulus == 0
				s3.Mul(b_of_x, f.eval_poly_at(zeropoly2, x))
				s0.Sub(p_of_x, s3)
				s1.Sub(s0, f.eval_poly_at(interpolant, x))
				if s2.Mod(s1, f.modulus).Cmp(ZERO) != 0 {
					finalErr = fmt.Errorf("boundary constraint violation")
				}
				// Check correctness of the linear combination
				// Check if (l_of_x - d_of_x - k1 * p_of_x - k2 * p_of_x * x_to_the_steps - k3 * b_of_x - k4 * b_of_x * x_to_the_steps) % modulus == 0
				s0.Sub(BytesToBig(l_of_x), d_of_x)
				s3.Mul(k1, p_of_x)
				s1.Sub(s0, s3)
				s2.Mul(p_of_x, x_to_the_steps)
				s3.Mul(k2, s2)
				s2.Sub(s1, s3)
				s3.Mul(k3, b_of_x)
				s0.Sub(s2, s3)
				s2.Mul(b_of_x, x_to_the_steps)
				s3.Mul(k4, s2)
				s1.Sub(s0, s3)
				s2.Mod(s1, f.modulus)
				if s2.Cmp(ZERO) != 0 {
					finalErr = fmt.Errorf("linear combination violation")
				}
			}
			wg.Done()
			// fmt.Printf("   consistency checks %d %d -- [%s]\n", i0, i1, time.Since(start0))
		}(j, j+njmp)
	}

	wg.Wait()
	fmt.Printf("Verified %d consistency checks [%s]\n", spot_check_security_factor, time.Since(start))

	return finalErr
}

// Verify an FRI proof
func verify_low_degree_proof(f *PrimeField, merkle_root []byte, root_of_unity *big.Int, components []*FriComponent, maxdeg_plus_1 int, exclude_multiples_of *big.Int) error {
	// Calculate which root of unity we're working with
	start := time.Now()
	testval := new(big.Int).Set(root_of_unity)
	roudeg := new(big.Int).Set(ONE)
	for testval.Cmp(ONE) != 0 {
		roudeg.Mul(roudeg, TWO)
		// testval = (testval * testval) % modulus
		testval.Mul(testval, testval)
		testval.Mod(testval, f.modulus)
	}

	// Powers of the given root of unity 1, p, p**2, p**3 such that p**4 = 1
	quartic_roots_of_unity := make([]*big.Int, 4)
	quartic_roots_of_unity[0] = new(big.Int).Set(ONE)
	quartic_roots_of_unity[1] = f.pow(root_of_unity, new(big.Int).Div(roudeg, FOUR))
	quartic_roots_of_unity[2] = f.pow(root_of_unity, new(big.Int).Div(roudeg, TWO))
	t := new(big.Int).Mul(roudeg, THREE)
	t.Div(t, FOUR)
	quartic_roots_of_unity[3] = f.pow(root_of_unity, t)
	if time.Since(start).Seconds() > MIN_SECONDS_BENCHMARK {

		fmt.Printf("  verify_low_degree_proof setup [%s]\n", time.Since(start))
	}
	// Verify the recursive components of the proof
	errc := make(chan error, len(components))
	for lvl, component := range components[0 : len(components)-1] {
		go func(level int, comp *FriComponent, merkle_root []byte, root_of_unity *big.Int, maxdeg_plus_1 int, roudeg *big.Int) {
			start = time.Now()

			root2 := comp.Root
			branches := comp.Branches

			// Calculate the pseudo-random x coordinate
			special_x := new(big.Int).SetBytes(merkle_root)
			special_x.Mod(special_x, f.modulus)

			// Calculate the pseudo-randomly sampled y indices
			modulus := new(big.Int).Div(roudeg, FOUR)
			ys, err := get_pseudorandom_indices(f, root2, modulus, 40, exclude_multiples_of.Int64())
			if err != nil {
				fmt.Printf("Failure get_pseudorandom_indices \n")
				errc <- err
				return
			}
			// For each y coordinate, get the x coordinates on the row, the values on the row, and the value at that y from the column
			xcoords := make([][]*big.Int, 0)
			rows := make([][]*big.Int, 0)
			columnvals := make([]*big.Int, 0)
			t := new(big.Int)
			for i, y := range ys {
				// The x coordinates from the polynomial
				x1 := f.pow(root_of_unity, y)
				a := make([]*big.Int, 4)
				for j := 0; j < 4; j++ {
					t.Mul(quartic_roots_of_unity[j], x1)
					a[j] = new(big.Int).Mod(t, f.modulus)
				}
				xcoords = append(xcoords, a)

				// The values from the original polynomial
				row := make([]*big.Int, 4)
				for j := int64(0); j < 4; j++ {
					idx := y.Int64() + int64(roudeg.Int64()/4*j)
					r, err := verify_branch_int(merkle_root, uint(idx), branches[i][j+1])
					if err != nil {
						errc <- err
						fmt.Printf("Failure VERIFY_BRANCH_INT 1\n")
						return
					}
					row[j] = r
				}

				rows = append(rows, row)

				c, err := verify_branch_int(root2, uint(y.Int64()), branches[i][0])
				if err != nil {
					errc <- err
					fmt.Printf("Failure VERIFY_BRANCH_INT 2\n")
					return
				}
				columnvals = append(columnvals, c)
			}

			// Verify for each selected y coordinate that the four points from the
			// polynomial and the one point from the column that are on that y
			// coordinate are on the same deg < 4 polynomial
			polys := f.multi_interp_4(xcoords, rows)
			for j, p := range polys {
				c := columnvals[j]
				q := f.eval_quartic(p, special_x)
				if q.Cmp(c) != 0 {
					fmt.Printf("Failure QUARTIC MISMATCH\n")

					errc <- fmt.Errorf("quartic mismatch")
				}
			}
			fmt.Printf("Verifying degree (%d) <= %d [%s]\n", level, maxdeg_plus_1, time.Since(start))

			errc <- nil
		}(lvl, component, merkle_root, root_of_unity, maxdeg_plus_1, roudeg)
		// Update constants to check the next proof
		start = time.Now()
		merkle_root = component.Root
		root_of_unity = f.pow(root_of_unity, FOUR)
		maxdeg_plus_1 = maxdeg_plus_1 / 4
		roudeg = f.div(roudeg, FOUR)
	}

	// Verify the direct components of the proof
	comp := components[len(components)-1]
	go func(comp *FriComponent) {
		start := time.Now()
		if maxdeg_plus_1 > 16 {
			errc <- fmt.Errorf("max_degreeplus_1 too high")
			return
		}
		// Check the Merkle root matches up
		mtree := merkelize(comp.Values)
		if bytes.Compare(mtree[1], merkle_root) != 0 {
			fmt.Printf("incorrect merkle root Failure\n")
			errc <- fmt.Errorf("incorrect merkle root")
		}

		// Check the degree of the data
		powers := get_power_cycle(root_of_unity, f.modulus)
		pts := make([]int64, 0)
		if exclude_multiples_of.Cmp(ZERO) > 0 {
			t := new(big.Int)
			for i := int64(0); i < int64(len(comp.Values)); i++ {
				if t.Mod(big.NewInt(i), exclude_multiples_of).Cmp(ZERO) > 0 {
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
			ys = append(ys, BytesToBig(comp.Values[x]))
		}

		poly := f.lagrange_interp(xs, ys)
		for _, x := range pts[maxdeg_plus_1:] {
			q := f.eval_poly_at(poly, powers[x])
			y := BytesToBig(comp.Values[x])
			if q.Cmp(y) != 0 {
				fmt.Printf("Lagrange mismatch Failure\n")
				errc <- fmt.Errorf("Lagrange mismatch")
			}
		}
		fmt.Printf("Verifying degree <= %d [%s]\n", maxdeg_plus_1, time.Since(start))
		errc <- nil
	}(comp)

	for i := 0; i < len(components); i++ {
		err := <-errc
		if err != nil {
			fmt.Printf("FRI proof NOT verified\n")
			return err
		}
	}
	fmt.Printf("FRI proof verified [%s]\n", time.Since(start))
	return nil
}
