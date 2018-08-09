package stark

import (
	"fmt"
	"math/big"
	"sync"
	"time"
)

// Creates an object that includes convenience operations for numbers and polynomials in some prime field
type PrimeField struct {
	modulus *big.Int
}

func NewPrimeField(_modulus *big.Int) (self *PrimeField, err error) {
	var s PrimeField
	s.SetModulus(_modulus)
	return &s, nil
}

func (self *PrimeField) SetModulus(_modulus *big.Int) {
	if _modulus == nil {
		modulus := new(big.Int)
		two_to_256 := new(big.Int)
		two_to_32_times_351 := new(big.Int)
		c := modulus.Mul(two_to_32_times_351.Exp(big.NewInt(2), big.NewInt(32), nil), big.NewInt(351))
		modulus.Sub(two_to_256.Exp(big.NewInt(2), big.NewInt(256), nil), c)
		self.modulus = modulus.Add(modulus, big.NewInt(1))
	} else {
		self.modulus = _modulus
	}
}

func (self *PrimeField) add(x, y *big.Int) (o *big.Int) {
	o = new(big.Int).Add(x, y)
	return new(big.Int).Mod(o, self.modulus)
}

func (self *PrimeField) sub(x, y *big.Int) (o *big.Int) {
	o = new(big.Int).Sub(x, y)
	return new(big.Int).Mod(o, self.modulus)
}

func (self *PrimeField) mul(x, y *big.Int) (o *big.Int) {
	o = new(big.Int).Mul(x, y)
	return new(big.Int).Mod(o, self.modulus)
}

func (self *PrimeField) div(x, y *big.Int) *big.Int {
	return self.mul(x, self.inv(y))
}

func (self *PrimeField) pow(x, y *big.Int) (o *big.Int) {
	o = new(big.Int).Exp(x, y, self.modulus)
	return o
}

// Modular inverse using the extended Euclidean algorithm
func (self *PrimeField) inv(a *big.Int) *big.Int {
	zero := big.NewInt(0)
	one := big.NewInt(1)
	if a.Cmp(zero) == 0 {
		return zero
	}
	lm := big.NewInt(1)
	hm := big.NewInt(0)

	low := a.Mod(a, self.modulus)
	high := self.modulus
	for low.Cmp(one) > 0 {
		r := new(big.Int).Div(high, low)

		// nm = hm-lm*r
		t := new(big.Int)
		nm := new(big.Int).Sub(hm, t.Mul(lm, r))

		// nw = high - low*r
		t.Mul(low, r)
		nw := new(big.Int).Sub(high, t)

		lm, low, hm, high = nm, nw, lm, low
	}
	return hm.Mod(lm, self.modulus)
}

func (self *PrimeField) multi_inv(values []*big.Int) []*big.Int {
	partials := make([]*big.Int, 1+len(values))
	partials[0] = big.NewInt(1)
	outputs := make([]*big.Int, len(values))
	z := big.NewInt(0)
	for i := 0; i < len(values); i++ {
		if values[i].Cmp(z) == 0 {
			partials[i+1] = partials[i]
		} else {
			partials[i+1] = self.mul(partials[i], values[i])
		}
	}
	inv := self.inv(partials[len(partials)-1])
	for i := len(values) - 1; i >= 0; i = i - 1 {
		if values[i].Cmp(z) == 0 {
			outputs[i] = big.NewInt(0)
		} else {
			outputs[i] = self.mul(partials[i], inv)
			inv = self.mul(inv, values[i])
		}
	}
	return outputs
}

// Evaluate a polynomial at a point
func (self *PrimeField) eval_poly_at(poly []*big.Int, x *big.Int) *big.Int {
	o := big.NewInt(0)
	p := big.NewInt(1)
	t1 := new(big.Int)
	t0 := new(big.Int)
	for _, coeff := range poly {
		t1.Mul(p, coeff)
		o.Add(o, t1)
		t0.Mul(p, x)
		p.Mod(t0, self.modulus)
	}
	return t0.Mod(o, self.modulus)
}

// Arithmetic for polynomials
func (self *PrimeField) add_polys(a, b []*big.Int) []*big.Int {
	l := len(a)
	if len(b) > l {
		l = len(b)
	}
	o := make([]*big.Int, 0)
	for i := 0; i < l; i++ {
		av := big.NewInt(0)
		bv := big.NewInt(0)
		if i < len(a) {
			av = a[i]
		}
		if i < len(b) {
			bv = b[i]
		}
		t := new(big.Int)
		t.Add(av, bv)
		o = append(o, t.Mod(t, self.modulus))
	}
	return o
}

func (self *PrimeField) sub_polys(a, b []*big.Int) []*big.Int {
	l := len(a)
	if len(b) > l {
		l = len(b)
	}
	o := make([]*big.Int, 0)
	for i := 0; i < l; i++ {
		av := big.NewInt(0)
		bv := big.NewInt(0)
		if i < len(a) {
			av = a[i]
		}
		if i < len(b) {
			bv = b[i]
		}
		t := new(big.Int)
		t.Sub(av, bv)
		o = append(o, t.Mod(t, self.modulus))
	}
	return o
}

func (self *PrimeField) mul_by_const(a []*big.Int, c *big.Int) []*big.Int {
	o := make([]*big.Int, 0)
	for i := 0; i < len(a); i++ {
		t := new(big.Int)
		t.Mul(a[i], c)
		o = append(o, t.Mod(t, self.modulus))
	}
	return o
}

func (self *PrimeField) mul_polys(a []*big.Int, b []*big.Int) []*big.Int {
	o := make([]*big.Int, len(a)+len(b)-1)
	for i := 0; i < len(a)+len(b)-1; i++ {
		o[i] = big.NewInt(0)
	}
	for i, aval := range a {
		for j, bval := range b {
			t := new(big.Int)
			t.Mul(aval, bval)
			o[i+j].Add(o[i+j], t)
		}
	}
	for i, _ := range o {
		o[i].Mod(o[i], self.modulus)
	}
	return o
}

func (self *PrimeField) div_polys(a []*big.Int, b []*big.Int) (o []*big.Int, err error) {
	if len(a) < len(b) {
		return o, fmt.Errorf("incorrect length of first input")
	}
	atmp := make([]*big.Int, len(a))
	for i, x := range a {
		atmp[i] = new(big.Int).Set(x)
	}
	o = make([]*big.Int, 0)
	apos := len(a) - 1
	bpos := len(b) - 1

	for diff := apos - bpos; diff >= 0; diff -= 1 {
		quot := self.div(atmp[apos], b[bpos])
		o = append([]*big.Int{quot}, o...)
		for i := bpos; i >= 0; i = i - 1 {
			t := new(big.Int).Mul(b[i], quot)
			atmp[diff+i].Sub(atmp[diff+i], t)
		}
		apos -= 1
	}
	for i, x := range o {
		o[i].Mod(x, self.modulus)
	}
	return o, nil
}

// Build a polynomial that returns 0 at all specified xs
func (self *PrimeField) zpoly(xs []*big.Int) []*big.Int {
	root := make([]*big.Int, 0)
	root = append(root, big.NewInt(1))
	t := new(big.Int)
	for _, x := range xs {
		root = append([]*big.Int{big.NewInt(0)}, root...)
		for j := 0; j < len(root)-1; j++ {
			root[j].Sub(root[j], t.Mul(root[j+1], x))
		}
	}
	for i, x := range root {
		root[i].Mod(x, self.modulus)
	}
	return root
}

/*
 Given p+1 y values and x values with no errors, recovers the original p+1 degree polynomial.
 Lagrange interpolation works roughly in the following way.
 1. Suppose you have a set of points, eg. x = [1, 2, 3], y = [2, 5, 10]
 2. For each x, generate a polynomial which equals its corresponding
    y coordinate at that point and 0 at all other points provided.
 3. Add these polynomials together.
*/
func (self *PrimeField) lagrange_interp(xs, ys []*big.Int) []*big.Int {
	// Generate master numerator polynomial, eg. (x - x1) * (x - x2) * ... * (x - xn)
	root := self.zpoly(xs)
	// assert len(root) == len(ys) + 1
	// print(root)
	// Generate per-value numerator polynomials, eg. for x=x2,
	// (x - x1) * (x - x3) * ... * (x - xn), by dividing the master
	// polynomial back by each x coordinate

	// nums = [self.div_polys(root, [-x, 1]) for x in xs]
	nums := make([][]*big.Int, len(xs))
	one := big.NewInt(1)
	negx := new(big.Int)
	for i, x := range xs {
		negx.Neg(x)
		nums[i], _ = self.div_polys(root, []*big.Int{negx, one})
	}

	// Generate denominators by evaluating numerator polys at each x
	// denoms = [self.eval_poly_at(nums[i], xs[i]) for i in range(len(xs))]
	denoms := make([]*big.Int, len(xs))
	for i := 0; i < len(xs); i++ {
		denoms[i] = self.eval_poly_at(nums[i], xs[i])
	}
	invdenoms := self.multi_inv(denoms)

	// Generate output polynomial, which is the sum of the per-value numerator polynomials rescaled to have the right y values
	b := make([]*big.Int, len(ys))
	for i := 0; i < len(ys); i++ {
		b[i] = big.NewInt(0)
	}
	for i := 0; i < len(xs); i++ {
		yslice := self.mul(ys[i], invdenoms[i])
		for j := 0; j < len(ys); j++ {
			if i < len(ys) && j < len(nums[i]) {
				t := new(big.Int).Mul(nums[i][j], yslice)
				b[j].Add(b[j], t)
			}
		}
	}
	for i := 0; i < len(b); i++ {
		b[i].Mod(b[i], self.modulus)
	}
	return b
}

func create_poly2(a, b *big.Int) (o []*big.Int) {
	return []*big.Int{a, b}
}

func create_poly4(a, b, c, d *big.Int) (o []*big.Int) {
	return []*big.Int{a, b, c, d}
}

// Optimized version of the above restricted to deg-2 polynomials
func (self *PrimeField) lagrange_interp_2(xs, ys []*big.Int) []*big.Int {
	m := self.modulus

	t0 := new(big.Int)
	t1 := new(big.Int)
	t2 := new(big.Int)
	one := big.NewInt(1)
	// eq0 = [-xs[1] % m, 1]
	t1.Neg(xs[1])
	t2.Mod(t1, m)
	eq0 := create_poly2(t2, one)
	e0 := self.eval_poly_at(eq0, xs[0])

	// eq1 = [-xs[0] % m, 1]
	t0.Neg(xs[0])
	t1.Mod(t0, m)
	eq1 := create_poly2(t1, one)
	e1 := self.eval_poly_at(eq1, xs[1])

	// invall = self.inv(e0 * e1)
	invall := self.inv(new(big.Int).Mul(e0, e1))

	// inv_y0 = ys[0] * invall * e1
	inv_y0 := new(big.Int)
	t0.Mul(ys[0], invall)
	inv_y0.Mul(t0, e1)

	// inv_y1 = ys[1] * invall * e0
	inv_y1 := new(big.Int)
	t0.Mul(ys[1], invall)
	inv_y1.Mul(t0, e0)

	// [(eq0[i] * inv_y0 + eq1[i] * inv_y1) % m for i in range(2)]
	o := make([]*big.Int, 2)
	for i := 0; i < 2; i++ {
		t0.Mul(eq0[i], inv_y0)
		t1.Mul(eq1[i], inv_y1)
		t2.Add(t0, t1)
		o[i] = new(big.Int).Mod(t2, m)
	}
	return o
}

// Optimized poly evaluation for degree 4
func (self *PrimeField) eval_quartic(p []*big.Int, x *big.Int) *big.Int {
	xsq := new(big.Int).Mul(x, x)
	xsq.Mod(xsq, self.modulus)
	xcb := new(big.Int).Mul(xsq, x)

	o3 := new(big.Int).Mul(p[3], xcb)
	o2 := new(big.Int).Mul(p[2], xsq)
	o1 := new(big.Int).Mul(p[1], x)
	o := new(big.Int).Add(p[0], o1)
	o1.Add(o3, o2)
	xcb.Add(o, o1)
	return o.Mod(xcb, self.modulus)
}

// Optimized version of the above restricted to deg-4 polynomials
func (self *PrimeField) lagrange_interp_4(xs, ys []*big.Int) []*big.Int {
	// x01, x02, x03, x12, x13, x23 := xs[0] * xs[1], xs[0] * xs[2], xs[0] * xs[3], xs[1] * xs[2], xs[1] * xs[3], xs[2] * xs[3]
	x01 := new(big.Int).Mul(xs[0], xs[1])
	x02 := new(big.Int).Mul(xs[0], xs[2])
	x03 := new(big.Int).Mul(xs[0], xs[3])
	x12 := new(big.Int).Mul(xs[1], xs[2])
	x13 := new(big.Int).Mul(xs[1], xs[3])
	x23 := new(big.Int).Mul(xs[2], xs[3])
	m := self.modulus
	t3 := new(big.Int)
	t2 := new(big.Int)
	t0 := new(big.Int)
	s3 := new(big.Int)
	s2 := new(big.Int)
	negxs0 := new(big.Int).Neg(xs[0])
	negxs2 := new(big.Int).Neg(xs[2])
	negxs3 := new(big.Int).Neg(xs[3])
	one := big.NewInt(1)

	//eq0 := create_poly4(-x12*xs[3]%m, (x12 + x13 + x23), -xs[1]-xs[2]-xs[3], 1)
	s3.Mul(x12, negxs3)
	t3.Mod(s3, m)
	s3.Add(x12, x13)
	t2.Add(s3, x23)
	t1 := s3.Neg(xs[1])
	s2.Sub(t3, xs[2])
	t1.Sub(s2, xs[3])
	eq0 := create_poly4(t3, t2, t1, one)

	//eq1 := create_poly4(-x02*xs[3]%m, (x02 + x03 + x23), -xs[0]-xs[2]-xs[3], 1)
	s3.Mul(x02, negxs3)
	t3.Mod(s3, m)
	t2.Add(x02, x03)
	t2.Add(t2, x23)
	s3.Sub(negxs0, xs[2])
	t1.Sub(s3, xs[3])
	eq1 := create_poly4(t3, t2, t1, one)

	//eq2 := create_poly4(-x01*xs[3]%m, (x01 + x03 + x13), -xs[0]-xs[1]-xs[3], 1)
	t3 = new(big.Int).Mul(x01, negxs3)
	t3.Mod(t3, m)
	t2 = new(big.Int).Set(x01)
	t2.Add(t2, x03)
	t2.Add(t2, x13)
	t1.Sub(negxs0, xs[1])
	t1.Sub(t1, xs[3])
	eq2 := create_poly4(t3, t2, t1, one)

	//eq3 := create_poly4(-x01*xs[2]%m, (x01 + x02 + x12), -xs[0]-xs[1]-xs[2], 1)
	t3.Mul(x01, negxs2)
	t3.Mod(t3, m)
	s3.Add(x01, x02)
	t2.Add(s3, x12)
	s3.Sub(negxs0, xs[1])
	t1.Sub(t3, xs[2])
	eq3 := create_poly4(t3, t2, t1, one)

	e0 := self.eval_poly_at(eq0, xs[0])
	e1 := self.eval_poly_at(eq1, xs[1])
	e2 := self.eval_poly_at(eq2, xs[2])
	e3 := self.eval_poly_at(eq3, xs[3])
	e01 := new(big.Int).Mul(e0, e1)
	e23 := new(big.Int).Mul(e2, e3)
	invall := self.inv(new(big.Int).Mul(e01, e23))

	// inv_y0 := ys[0] * invall * e1 * e23 % m
	inv_y0 := new(big.Int).Mul(ys[0], invall)
	s3.Mul(inv_y0, e1)
	s2.Mul(s3, e23)
	inv_y0.Mod(s2, m)

	// inv_y1 := ys[1] * invall * e0 * e23 % m
	inv_y1 := new(big.Int).Mul(ys[1], invall)
	s3.Mul(inv_y1, e0)
	s2.Mul(s3, e23)
	inv_y1.Mod(s2, m)

	// inv_y2 := ys[2] * invall * e01 * e3 % m
	inv_y2 := new(big.Int).Mul(ys[2], invall)
	s3.Mul(inv_y2, e01)
	s2.Mul(s3, e3)
	inv_y2.Mod(s2, m)

	// inv_y3 := ys[3] * invall * e01 * e2 % m
	inv_y3 := new(big.Int).Mul(ys[3], invall)
	s3.Mul(inv_y3, e01)
	s2.Mul(s3, e2)
	inv_y3.Mod(s2, m)

	o := make([]*big.Int, 4)
	for i := 0; i < 4; i++ {
		// [(eq0[i] * inv_y0 + eq1[i] * inv_y1 + eq2[i] * inv_y2 + eq3[i] * inv_y3) % m for i in range(4)]
		t0.Mul(eq0[i], inv_y0)
		t1.Mul(eq1[i], inv_y1)
		t2.Mul(eq2[i], inv_y2)
		t3.Mul(eq3[i], inv_y3)
		s2.Add(t0, t1)
		s3.Add(t2, t3)
		t3.Add(s2, s3)
		o[i] = new(big.Int).Mod(t3, m)
	}
	return o
}

func (self *PrimeField) _simple_ft(vals []*big.Int, roots_of_unity []*big.Int) []*big.Int {
	L := len(roots_of_unity)
	o := make([]*big.Int, L)
	for i := 0; i < L; i++ {
		o[i] = big.NewInt(0)
	}

	for i := 0; i < L; i++ {
		last := big.NewInt(0)
		for j := 0; j < L; j++ {
			t := new(big.Int).Mul(vals[j], roots_of_unity[(i*j)%L])
			o[i].Add(o[i], t)
		}
		last.Mod(last, self.modulus)
	}
	for i := 0; i < L; i++ {
		o[i].Mod(o[i], self.modulus)
	}
	return o
}

func (self *PrimeField) _fft(vals []*big.Int, roots_of_unity []*big.Int) []*big.Int {
	if len(vals) <= 1 {
		return vals
	}

	roots_of_unity2 := len(roots_of_unity) / 2
	root2 := make([]*big.Int, roots_of_unity2)
	vals_div2 := len(vals) / 2
	for i := 0; i < roots_of_unity2; i++ {
		root2[i] = roots_of_unity[i*2]
	}
	o := make([]*big.Int, len(vals))

	var L []*big.Int
	var R []*big.Int
	if len(vals) > 1024 {
		var wg sync.WaitGroup
		y_times_root := make([]*big.Int, vals_div2)
		wg.Add(1)
		go func() {
			lvals := make([]*big.Int, vals_div2)
			for i := 0; i < vals_div2; i++ {
				lvals[i] = vals[i*2]
			}
			L = self._fft(lvals, root2)
			wg.Done()
		}()
		wg.Add(1)
		go func() {
			rvals := make([]*big.Int, vals_div2)
			for i := 0; i < vals_div2; i++ {
				rvals[i] = vals[i*2+1]
			}
			R = self._fft(rvals, root2)
			for i, rval := range R {
				y_times_root[i] = new(big.Int).Mul(rval, roots_of_unity[i])
			}
			wg.Done()
		}()
		wg.Wait()

		wg.Add(1)
		go func() {
			t := new(big.Int)
			for i, x := range L {
				t.Add(x, y_times_root[i])
				o[i] = new(big.Int).Mod(t, self.modulus)
			}
			wg.Done()
		}()
		wg.Add(1)
		go func() {
			t := new(big.Int)
			for i, x := range L {
				t.Sub(x, y_times_root[i])
				o[i+len(L)] = new(big.Int).Mod(t, self.modulus)
			}
			wg.Done()
		}()
		wg.Wait()
	} else {
		lvals := make([]*big.Int, vals_div2)
		for i := 0; i < vals_div2; i++ {
			lvals[i] = vals[i*2]
		}
		L = self._fft(lvals, root2)

		rvals := make([]*big.Int, vals_div2)
		for i := 0; i < vals_div2; i++ {
			rvals[i] = vals[i*2+1]
		}
		R = self._fft(rvals, root2)

		y_times_root := new(big.Int)
		t1 := new(big.Int)
		t2 := new(big.Int)
		for i, x := range L {
			y_times_root.Mul(R[i], roots_of_unity[i])
			t1.Add(x, y_times_root)
			o[i] = new(big.Int).Mod(t1, self.modulus)
			t2.Sub(x, y_times_root)
			o[i+len(L)] = new(big.Int).Mod(t2, self.modulus)
		}
	}
	return o
}

func (self *PrimeField) fft(vals []*big.Int, root_of_unity *big.Int, inv bool) []*big.Int {
	// Build up roots of unity
	start := time.Now()
	rootz := make([]*big.Int, 2)
	rootz[0] = big.NewInt(1)
	rootz[1] = root_of_unity
	one := big.NewInt(1)
	i := 1
	for rootz[i].Cmp(one) != 0 {
		t := new(big.Int).Mul(rootz[i], root_of_unity)
		rootz = append(rootz, t.Mod(t, self.modulus))
		i = i + 1
	}

	// Fill in vals with zeroes if needed
	if len(rootz) > len(vals)+1 {
		extrazeros := make([]*big.Int, (len(rootz) - len(vals) - 1))
		for i := 0; i < len(extrazeros); i++ {
			extrazeros[i] = big.NewInt(0)
		}
		vals = append(vals, extrazeros...)
	}
	if time.Since(start).Seconds() > MIN_SECONDS_BENCHMARK {
		fmt.Printf("    fft setup [%s]\n", time.Since(start))
	}
	if inv {
		// Inverse FFT
		start = time.Now()
		t := new(big.Int).Sub(self.modulus, big.NewInt(2))
		invlen := new(big.Int).Exp(big.NewInt(int64(len(vals))), t, self.modulus)
		irootz := make([]*big.Int, 0)
		for i := len(rootz) - 1; i > 0; i-- {
			irootz = append(irootz, rootz[i])
		}
		if time.Since(start).Seconds() > MIN_SECONDS_BENCHMARK {
			fmt.Printf("    invfft setup2 [%s]\n", time.Since(start))
		}
		start = time.Now()
		res := self._fft(vals, irootz)
		if time.Since(start).Seconds() > MIN_SECONDS_BENCHMARK {
			fmt.Printf("    inv_fft core [%s]\n", time.Since(start))
		}
		start = time.Now()
		o := make([]*big.Int, len(res))
		for i, x := range res {
			q := new(big.Int).Mul(x, invlen)
			q.Mod(q, self.modulus)
			o[i] = q
		}
		if time.Since(start).Seconds() > MIN_SECONDS_BENCHMARK {
			fmt.Printf("    invfft final [%s]\n", time.Since(start))
		}
		return o
	} else {
		// Regular FFT
		start = time.Now()
		res := self._fft(vals, rootz[0:len(rootz)-1])
		if time.Since(start).Seconds() > MIN_SECONDS_BENCHMARK {
			fmt.Printf("    reg_fft core [%s]\n", time.Since(start))
		}
		return res
	}
}

func (self *PrimeField) mul_polys_fft(a []*big.Int, b []*big.Int, root_of_unity *big.Int) []*big.Int {
	x1 := self.fft(a, root_of_unity, false)
	x2 := self.fft(b, root_of_unity, false)
	c := make([]*big.Int, len(x1))
	for i, v1 := range x1 {
		t := new(big.Int).Mul(v1, x2[i])
		c[i] = t.Mod(t, self.modulus)
	}
	return self.fft(c, root_of_unity, true)
}

// Optimized version of the above restricted to deg-4 polynomials
func (self *PrimeField) multi_interp_4(xsets, ysets [][]*big.Int) [][]*big.Int {
	data := make([][][]*big.Int, len(xsets))
	invtargets := make([]*big.Int, len(xsets)*4)
	m := self.modulus
	var wg sync.WaitGroup
	nj := len(xsets) / NUM_CORES
	if nj < 500 {
		nj = 500
	}
	for i := 0; i < len(xsets); i += nj {
		wg.Add(1)
		go func(i0 int, i1 int) {
			if i1 > len(xsets) {
				i1 = len(xsets)
			}
			nxs0 := new(big.Int)
			nxs1 := new(big.Int)
			nxs2 := new(big.Int)
			nxs3 := new(big.Int)
			x01 := new(big.Int)
			x02 := new(big.Int)
			x03 := new(big.Int)
			x12 := new(big.Int)
			x13 := new(big.Int)
			x23 := new(big.Int)
			sum_xs2_xs3 := new(big.Int)
			sub_nxs0_xs1 := new(big.Int)
			//s3 := new(big.Int)
			//s2 := new(big.Int)
			//ONE := big.NewInt(1)

			for j := i0; j < i1; j++ {
				xs := xsets[j] //[]*big.Int

				x01.Mul(xs[0], xs[1])
				x02.Mul(xs[0], xs[2])
				x03.Mul(xs[0], xs[3])
				x12.Mul(xs[1], xs[2])
				x13.Mul(xs[1], xs[3])
				x23.Mul(xs[2], xs[3])
				nxs0.Neg(xs[0])
				nxs1.Neg(xs[1])
				nxs2.Neg(xs[2])
				nxs3.Neg(xs[3])
				sum_xs2_xs3.Add(xs[2], xs[3])
				sub_nxs0_xs1.Sub(nxs0, xs[1])

				// eq0 := create_poly4(-x12 * xs[3] % m, (x12 + x13 + x23), -xs[1]-xs[2]-xs[3], 1)
				t3 := new(big.Int).Mul(nxs3, x12)
				t2 := new(big.Int).Add(x12, x13)
				t1 := new(big.Int).Sub(nxs1, sum_xs2_xs3)
				eq0 := create_poly4(t3.Mod(t3, m), t2.Add(t2, x23), t1, big.NewInt(1))
				invtargets[j*4+0] = self.eval_quartic(eq0, xs[0])

				// eq1 := create_poly4(-x02 * xs[3] % m, (x02 + x03 + x23), -xs[0]-xs[2]-xs[3], 1)
				t3 = new(big.Int).Mul(x02, nxs3)
				t2 = new(big.Int).Add(x02, x03)
				t1 = new(big.Int).Sub(nxs1, sum_xs2_xs3)
				eq1 := create_poly4(t3.Mod(t3, m), t2.Add(t2, x23), t1.Sub(nxs0, sum_xs2_xs3), big.NewInt(1))
				invtargets[j*4+1] = self.eval_quartic(eq1, xs[1])

				// eq2 := create_poly4(-x01 * xs[3] % m, (x01 + x03 + x13), -xs[0]-xs[1]-xs[3], 1)
				t3 = new(big.Int).Mul(x01, nxs3)
				t2 = new(big.Int).Add(x01, x03)
				t1 = new(big.Int).Sub(sub_nxs0_xs1, xs[3])
				eq2 := create_poly4(t3.Mod(t3, m), t2.Add(t2, x13), t1, big.NewInt(1))
				invtargets[j*4+2] = self.eval_quartic(eq2, xs[2])

				// eq3 := create_poly4(-x01 * xs[2] % m, (x01 + x02 + x12), -xs[0]-xs[1]-xs[2], 1)
				t3 = new(big.Int).Mul(x01, nxs2)
				t3.Mod(t3, m)
				t2 = new(big.Int).Add(x01, x02)
				t2.Add(t2, x12)
				t1 = new(big.Int).Sub(sub_nxs0_xs1, xs[2])
				eq3 := create_poly4(t3, t2, t1, big.NewInt(1))
				invtargets[j*4+3] = self.eval_quartic(eq3, xs[3])

				data[j] = [][]*big.Int{ysets[j], eq0, eq1, eq2, eq3}
			}
			wg.Done()
		}(i, i+nj)
	}
	// wait for goroutines to finish

	wg.Wait()

	invalls := self.multi_inv(invtargets)
	nd := len(data)
	dj := 1500
	o := make([][]*big.Int, nd)
	for q := 0; q < nd; q += dj {
		wg.Add(1)
		go func(i0 int, i1 int) {
			for i := i0; i < i1; i++ {
				if i < nd {
					d := data[i]
					ys := d[0]
					eq0, eq1, eq2, eq3 := d[1], d[2], d[3], d[4]
					invallz := invalls[i*4 : i*4+4]
					inv_y0 := new(big.Int).Mul(ys[0], invallz[0])
					inv_y0.Mod(inv_y0, self.modulus)
					inv_y1 := new(big.Int).Mul(ys[1], invallz[1])
					inv_y1.Mod(inv_y1, self.modulus)
					inv_y2 := new(big.Int).Mul(ys[2], invallz[2])
					inv_y2.Mod(inv_y2, self.modulus)
					inv_y3 := new(big.Int).Mul(ys[3], invallz[3])
					inv_y3.Mod(inv_y3, self.modulus)
					e := make([]*big.Int, 4)
					// [(eq0[i] * inv_y0 + eq1[i] * inv_y1 + eq2[i] * inv_y2 + eq3[i] * inv_y3) % m for i in range(4)])
					for j := 0; j < 4; j++ {
						a := new(big.Int).Add(new(big.Int).Mul(eq0[j], inv_y0), new(big.Int).Mul(eq1[j], inv_y1))
						a.Add(a, new(big.Int).Mul(eq2[j], inv_y2))
						a.Add(a, new(big.Int).Mul(eq3[j], inv_y3))
						e[j] = a.Mod(a, self.modulus)
					}
					o[i] = e
				}
			}
			wg.Done()
		}(q, q+dj)
	}

	wg.Wait()
	// TODO: assert o == [self.lagrange_interp_4(xs, ys) for xs, ys in zip(xsets, ysets)]
	return o
}

// Optimized version of the above restricted to deg-4 polynomials
func (self *PrimeField) multi_interp_4new(xsets, ysets [][]*big.Int) [][]*big.Int {
	data := make([][][]*big.Int, len(xsets))
	invtargets := make([]*big.Int, len(xsets)*4)
	start := time.Now()
	m := self.modulus
	var wg sync.WaitGroup
	nj := len(xsets) / NUM_CORES
	if nj < 256 {
		nj = 256
	}
	for i := 0; i < len(xsets); i += nj {
		wg.Add(1)
		go func(j0 int, j1 int) {
			startb := time.Now()
			if j1 > len(xsets) {
				j1 = len(xsets)
			}
			nxs0 := new(big.Int)
			nxs1 := new(big.Int)
			nxs2 := new(big.Int)
			nxs3 := new(big.Int)
			x01 := new(big.Int)
			x02 := new(big.Int)
			x03 := new(big.Int)
			x12 := new(big.Int)
			x13 := new(big.Int)
			x23 := new(big.Int)
			sum_xs2_xs3 := new(big.Int)
			sub_nxs0_xs1 := new(big.Int)
			s3 := new(big.Int)
			s2 := new(big.Int)
			ONE := big.NewInt(1)
			for j := j0; j < j1; j++ {
				xs := xsets[j]

				x01.Mul(xs[0], xs[1])
				x02.Mul(xs[0], xs[2])
				x03.Mul(xs[0], xs[3])
				x12.Mul(xs[1], xs[2])
				x13.Mul(xs[1], xs[3])
				x23.Mul(xs[2], xs[3])
				nxs0.Neg(xs[0])
				nxs1.Neg(xs[1])
				nxs2.Neg(xs[2])
				nxs3.Neg(xs[3])
				sum_xs2_xs3.Add(xs[2], xs[3])
				sub_nxs0_xs1.Sub(nxs0, xs[1])

				// eq0 := create_poly4(-x12 * xs[3] % m, (x12 + x13 + x23), -xs[1]-xs[2]-xs[3], 1)
				s3.Mul(nxs3, x12)
				s2.Add(x12, x13)
				eq0 := create_poly4(new(big.Int).Mod(s3, m), new(big.Int).Add(s2, x23), new(big.Int).Sub(nxs1, sum_xs2_xs3), ONE)
				invtargets[j*4+0] = self.eval_quartic(eq0, xs[0])

				// eq1 := create_poly4(-x02 * xs[3] % m, (x02 + x03 + x23), -xs[0]-xs[2]-xs[3], 1)
				s3.Mul(x02, nxs3)
				s2.Add(x02, x03)
				eq1 := create_poly4(new(big.Int).Mod(s3, m), new(big.Int).Add(s2, x23), new(big.Int).Sub(nxs0, sum_xs2_xs3), ONE)
				invtargets[j*4+1] = self.eval_quartic(eq1, xs[1])

				// eq2 := create_poly4(-x01 * xs[3] % m, (x01 + x03 + x13), -xs[0]-xs[1]-xs[3], 1)
				s3.Mul(x01, nxs3)
				s2.Add(x01, x03)
				eq2 := create_poly4(new(big.Int).Mod(s3, m), new(big.Int).Add(s2, x13), new(big.Int).Sub(sub_nxs0_xs1, xs[3]), ONE)
				invtargets[j*4+2] = self.eval_quartic(eq2, xs[2])

				// eq3 := create_poly4(-x01 * xs[2] % m, (x01 + x02 + x12), -xs[0]-xs[1]-xs[2], 1)
				s3.Mul(x01, nxs2)
				s2.Add(x01, x02)
				eq3 := create_poly4(new(big.Int).Mod(s3, m), new(big.Int).Add(s2, x12), new(big.Int).Sub(sub_nxs0_xs1, xs[2]), ONE)
				invtargets[j*4+3] = self.eval_quartic(eq3, xs[3])

				data[j] = [][]*big.Int{ysets[j], eq0, eq1, eq2, eq3}
			}
			wg.Done()
			if time.Since(startb).Seconds() > MIN_SECONDS_BENCHMARK {
				//	fmt.Printf("--> %d %d [%d] %s\n", j0, j1, nj, time.Since(startb))
			}
		}(i, i+nj)
	}
	// wait for goroutines to finish
	wg.Wait()
	invalls := self.multi_inv(invtargets)
	start = time.Now()
	nd := len(data)
	dj := len(data) / NUM_CORES
	if dj < 500 {
		dj = 500
	}
	o := make([][]*big.Int, nd)
	for q := 0; q < nd; q += dj {
		wg.Add(1)
		go func(i0 int, i1 int) {
			if i1 > nd {
				i1 = nd
			}
			for i := i0; i < i1; i++ {
				d := data[i]
				ys := d[0]
				eq0, eq1, eq2, eq3 := d[1], d[2], d[3], d[4]
				invallz := invalls[i*4 : i*4+4]
				t0 := new(big.Int)
				t1 := new(big.Int)
				t2 := new(big.Int)
				t3 := new(big.Int)
				t0.Mul(ys[0], invallz[0])
				inv_y0 := new(big.Int).Mod(t0, self.modulus)
				t0.Mul(ys[1], invallz[1])
				inv_y1 := new(big.Int).Mod(t0, self.modulus)
				t0.Mul(ys[2], invallz[2])
				inv_y2 := new(big.Int).Mod(t0, self.modulus)
				t0.Mul(ys[3], invallz[3])
				inv_y3 := new(big.Int).Mod(t0, self.modulus)
				e := make([]*big.Int, 4)
				// [(eq0[i] * inv_y0 + eq1[i] * inv_y1 + eq2[i] * inv_y2 + eq3[i] * inv_y3) % m for i in range(4)])
				for j := 0; j < 4; j++ {
					t0.Mul(eq0[j], inv_y0)
					t1.Mul(eq1[j], inv_y1)
					t2.Add(t0, t1)
					t0.Mul(eq2[j], inv_y2)
					t1.Mul(eq3[j], inv_y3)
					t3.Add(t0, t1)
					t0.Add(t2, t3)
					e[j] = t1.Mod(t0, self.modulus)
				}
				o[i] = e
			}
			wg.Done()
		}(q, q+dj)
	}
	wg.Wait()
	if time.Since(start).Seconds() > MIN_SECONDS_BENCHMARK {
		//		fmt.Printf("    multi_interp_4 phase 3b [%s]\n", time.Since(start))
	}
	// TODO: assert o == [self.lagrange_interp_4(xs, ys) for xs, ys in zip(xsets, ysets)]
	return o
}

// Compute a MIMC permutation for some number of steps
func (self *PrimeField) MiMC(inp *big.Int, steps *big.Int, round_constants []*big.Int) *big.Int {
	start := time.Now()
	t := new(big.Int).Set(inp)
	t1 := new(big.Int)
	t2 := new(big.Int)
	nsteps := int(steps.Int64() - 1)
	THREE := big.NewInt(3)
	for i := 0; i < nsteps; i++ {
		// inp = (inp**3 + round_constants[i % len(round_constants)]) % modulus
		t1.Exp(t, THREE, nil)
		t2.Add(t1, round_constants[i%len(round_constants)])
		t.Mod(t2, self.modulus)
	}
	fmt.Printf("MIMC computed in %s\n\n", time.Since(start))
	return t
}
