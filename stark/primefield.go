package stark

import (
	"fmt"
	"math/big"
	"sync"
	"time"
)

var (
	ZERO        = big.NewInt(0)
	NEGONE      = big.NewInt(-1)
	ONE         = big.NewInt(1)
	TWO         = big.NewInt(2)
	THREE       = big.NewInt(3)
	FOUR        = big.NewInt(4)
	FIVE        = big.NewInt(5)
	SIX         = big.NewInt(6)
	SEVEN       = big.NewInt(7)
	TWENTYFOUR  = big.NewInt(24)
	THIRTYTWO   = big.NewInt(32)
	TWOFIFTYSIX = big.NewInt(256)
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
		c := modulus.Mul(two_to_32_times_351.Exp(TWO, THIRTYTWO, nil), big.NewInt(351))
		modulus.Sub(two_to_256.Exp(TWO, TWOFIFTYSIX, nil), c)
		self.modulus = modulus.Add(modulus, ONE)
	} else {
		self.modulus = _modulus
	}
}

func (self *PrimeField) add(x, y *big.Int) (o *big.Int) {
	return new(big.Int).Mod(new(big.Int).Add(x, y), self.modulus)
}

func (self *PrimeField) sub(x, y *big.Int) (o *big.Int) {
	return new(big.Int).Mod(new(big.Int).Sub(x, y), self.modulus)
}

func (self *PrimeField) mul(x, y *big.Int) (o *big.Int) {
	return new(big.Int).Mod(new(big.Int).Mul(x, y), self.modulus)
}

func (self *PrimeField) div(x, y *big.Int) *big.Int {
	return self.mul(x, self.inv(y))
}

func (self *PrimeField) pow(x, y *big.Int) (o *big.Int) {
	return new(big.Int).Exp(x, y, self.modulus)
}

// Modular inverse using the extended Euclidean algorithm
func (self *PrimeField) inv(a *big.Int) *big.Int {
	if a.Cmp(ZERO) == 0 {
		return ZERO
	}
	lm := new(big.Int).Set(ONE)
	hm := new(big.Int)

	low := a.Mod(a, self.modulus)
	high := self.modulus
	t := new(big.Int)
	r := new(big.Int)
	for low.Cmp(ONE) > 0 {
		r.Div(high, low)
		// nm = hm-lm*r
		nm := new(big.Int).Sub(hm, t.Mul(lm, r))
		// nw = high - low*r
		nw := new(big.Int).Sub(high, t.Mul(low, r))
		lm, low, hm, high = nm, nw, lm, low
	}
	return hm.Mod(lm, self.modulus)
}

func (self *PrimeField) multi_inv(values []*big.Int) []*big.Int {
	partials := make([]*big.Int, 1+len(values))
	partials[0] = new(big.Int).Set(ONE)
	outputs := make([]*big.Int, len(values))

	for i := 0; i < len(values); i++ {
		if values[i].Cmp(ZERO) == 0 {
			partials[i+1] = partials[i]
		} else {
			partials[i+1] = self.mul(partials[i], values[i])
		}
	}
	inv := self.inv(partials[len(partials)-1])
	for i := len(values) - 1; i >= 0; i = i - 1 {
		if values[i].Cmp(ZERO) == 0 {
			outputs[i] = new(big.Int)
		} else {
			outputs[i] = self.mul(partials[i], inv)
			inv = self.mul(inv, values[i])
		}
	}
	return outputs
}

// Evaluate a polynomial at a point
func (self *PrimeField) eval_poly_at(poly []*big.Int, x *big.Int) *big.Int {
	o := new(big.Int)
	p := new(big.Int).Set(ONE)
	t1 := new(big.Int)
	t0 := new(big.Int)
	for _, coeff := range poly {
		o.Add(o, t1.Mul(p, coeff))
		p.Mod(t0.Mul(p, x), self.modulus)
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
	t := new(big.Int)
	for i := 0; i < l; i++ {
		av := new(big.Int)
		bv := new(big.Int)
		if i < len(a) {
			av = a[i]
		}
		if i < len(b) {
			bv = b[i]
		}
		o = append(o, new(big.Int).Mod(t.Add(av, bv), self.modulus))
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
		av := new(big.Int)
		bv := new(big.Int)
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
	t := new(big.Int)
	for i := 0; i < len(a); i++ {
		o = append(o, new(big.Int).Mod(t.Mul(a[i], c), self.modulus))
	}
	return o
}

func (self *PrimeField) mul_polys(a []*big.Int, b []*big.Int) []*big.Int {
	o := make([]*big.Int, len(a)+len(b)-1)
	for i := 0; i < len(a)+len(b)-1; i++ {
		o[i] = new(big.Int)
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
	root = append(root, ONE)
	t := new(big.Int)
	for _, x := range xs {
		root = append([]*big.Int{new(big.Int)}, root...)
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
	negx := new(big.Int)
	for i, x := range xs {
		nums[i], _ = self.div_polys(root, []*big.Int{negx.Neg(x), ONE})
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
		b[i] = new(big.Int)
	}
	t := new(big.Int)
	for i := 0; i < len(xs); i++ {
		yslice := self.mul(ys[i], invdenoms[i])
		for j := 0; j < len(ys); j++ {
			if i < len(ys) && j < len(nums[i]) {
				t.Mul(nums[i][j], yslice)
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

	// eq0 = [-xs[1] % m, 1]
	eq0 := create_poly2(t2.Mod(t1.Neg(xs[1]), m), ONE)
	e0 := self.eval_poly_at(eq0, xs[0])

	// eq1 = [-xs[0] % m, 1]
	eq1 := create_poly2(t1.Mod(t0.Neg(xs[0]), m), ONE)
	e1 := self.eval_poly_at(eq1, xs[1])

	// invall = self.inv(e0 * e1)
	invall := self.inv(t0.Mul(e0, e1))

	// inv_y0 = ys[0] * invall * e1
	inv_y0 := new(big.Int).Mul(t0.Mul(ys[0], invall), e1)

	// inv_y1 = ys[1] * invall * e0
	inv_y1 := new(big.Int).Mul(t0.Mul(ys[1], invall), e0)

	// [(eq0[i] * inv_y0 + eq1[i] * inv_y1) % m for i in range(2)]
	o := make([]*big.Int, 2)
	for i := 0; i < 2; i++ {
		o[i] = new(big.Int).Mod(t2.Add(t0.Mul(eq0[i], inv_y0), t1.Mul(eq1[i], inv_y1)), m)
	}
	return o
}

// Optimized poly evaluation for degree 4
func (self *PrimeField) eval_quartic(p []*big.Int, x *big.Int) *big.Int {
	xsq := new(big.Int).Mul(x, x)
	xcb := new(big.Int).Mul(xsq, x)
	t3 := new(big.Int)
	xsq.Add(t3.Mul(p[3], xcb), new(big.Int).Mul(p[2], xsq))
	xcb.Add(new(big.Int).Add(p[0], t3.Mul(p[1], x)), xsq)
	return xcb.Mod(xcb, self.modulus)
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
	t1 := new(big.Int)
	t0 := new(big.Int)
	s3 := new(big.Int)
	s2 := new(big.Int)
	negxs0 := new(big.Int).Neg(xs[0])
	negxs1 := new(big.Int).Neg(xs[1])
	negxs2 := new(big.Int).Neg(xs[2])
	negxs3 := new(big.Int).Neg(xs[3])

	//eq0 := create_poly4(-x12*xs[3]%m, (x12 + x13 + x23), -xs[1]-xs[2]-xs[3], 1)
	eq0 := create_poly4(new(big.Int).Mod(t3.Mul(x12, negxs3), m), new(big.Int).Add(t2.Add(x12, x13), x23), new(big.Int).Sub(negxs1, t0.Add(xs[3], xs[2])), ONE)

	//eq1 := create_poly4(-x02*xs[3]%m, (x02 + x03 + x23), -xs[0]-xs[2]-xs[3], 1)
	eq1 := create_poly4(new(big.Int).Mod(t3.Mul(x02, negxs3), m), new(big.Int).Add(t2.Add(x02, x03), x23), new(big.Int).Sub(t1.Sub(negxs0, xs[2]), xs[3]), ONE)

	//eq2 := create_poly4(-x01*xs[3]%m, (x01 + x03 + x13), -xs[0]-xs[1]-xs[3], 1)
	eq2 := create_poly4(new(big.Int).Mod(t3.Mul(x01, negxs3), m), new(big.Int).Add(t2.Add(x01, x03), x13), new(big.Int).Sub(negxs0, t1.Add(xs[1], xs[3])), ONE)

	//eq3 := create_poly4(-x01*xs[2]%m, (x01 + x02 + x12), -xs[0]-xs[1]-xs[2], 1)
	eq3 := create_poly4(new(big.Int).Mod(t3.Mul(x01, negxs2), m), new(big.Int).Add(t2.Add(x01, x02), x12), new(big.Int).Sub(negxs0, t1.Add(xs[1], xs[2])), ONE)

	e0 := self.eval_poly_at(eq0, xs[0])
	e1 := self.eval_poly_at(eq1, xs[1])
	e2 := self.eval_poly_at(eq2, xs[2])
	e3 := self.eval_poly_at(eq3, xs[3])
	e01 := new(big.Int).Mul(e0, e1)
	e23 := new(big.Int).Mul(e2, e3)
	invall := self.inv(t0.Mul(e01, e23))

	// inv_y0 := ys[0] * invall * e1 * e23 % m
	inv_y0 := new(big.Int).Mod(t2.Mul(t0.Mul(ys[0], invall), t1.Mul(e1, e23)), m)

	// inv_y1 := ys[1] * invall * e0 * e23 % m
	inv_y1 := new(big.Int).Mod(t2.Mul(t0.Mul(ys[1], invall), t1.Mul(e0, e23)), m)

	// inv_y2 := ys[2] * invall * e01 * e3 % m
	inv_y2 := new(big.Int).Mod(t2.Mul(t0.Mul(ys[2], invall), t1.Mul(e01, e3)), m)

	// inv_y3 := ys[3] * invall * e01 * e2 % m
	inv_y3 := new(big.Int).Mod(t2.Mul(t0.Mul(ys[3], invall), t1.Mul(e01, e2)), m)

	o := make([]*big.Int, 4)
	for i := 0; i < 4; i++ {
		// [(eq0[i] * inv_y0 + eq1[i] * inv_y1 + eq2[i] * inv_y2 + eq3[i] * inv_y3) % m for i in range(4)]
		t3.Add(s2.Add(t0.Mul(eq0[i], inv_y0), t1.Mul(eq1[i], inv_y1)), s3.Add(t2.Mul(eq2[i], inv_y2), t3.Mul(eq3[i], inv_y3)))
		o[i] = new(big.Int).Mod(t3, m)
	}
	return o
}

func (self *PrimeField) _simple_ft(vals []*big.Int, roots_of_unity []*big.Int) []*big.Int {
	L := len(roots_of_unity)
	o := make([]*big.Int, L)
	for i := 0; i < L; i++ {
		o[i] = new(big.Int)
	}

	for i := 0; i < L; i++ {
		last := new(big.Int)
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
	if len(vals) >= 1024 {
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
			o[i] = new(big.Int).Mod(t1.Add(x, y_times_root), self.modulus)
			o[i+len(L)] = new(big.Int).Mod(t2.Sub(x, y_times_root), self.modulus)
		}
	}
	return o
}

func (self *PrimeField) fft(vals []*big.Int, root_of_unity *big.Int, inv bool) []*big.Int {
	// Build up roots of unity
	start := time.Now()
	rootz := make([]*big.Int, 2)
	rootz[0] = new(big.Int).Set(ONE)
	rootz[1] = root_of_unity

	i := 1
	for rootz[i].Cmp(ONE) != 0 {
		t := new(big.Int).Mul(rootz[i], root_of_unity)
		rootz = append(rootz, t.Mod(t, self.modulus))
		i = i + 1
	}

	// Fill in vals with zeroes if needed
	if len(rootz) > len(vals)+1 {
		extrazeros := make([]*big.Int, (len(rootz) - len(vals) - 1))
		for i := 0; i < len(extrazeros); i++ {
			extrazeros[i] = new(big.Int)
		}
		vals = append(vals, extrazeros...)
	}
	if time.Since(start).Seconds() > MIN_SECONDS_BENCHMARK {
		//	fmt.Printf("    fft setup [%s]\n", time.Since(start))
	}
	if inv {
		// Inverse FFT
		start = time.Now()
		t := new(big.Int).Sub(self.modulus, TWO)
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
			//	fmt.Printf("    inv_fft core [%s]\n", time.Since(start))
		}
		start = time.Now()
		o := make([]*big.Int, len(res))
		q := new(big.Int)
		for i, x := range res {
			q.Mul(x, invlen)
			o[i] = new(big.Int).Mod(q, self.modulus)
		}
		if time.Since(start).Seconds() > MIN_SECONDS_BENCHMARK {
			//	fmt.Printf("    invfft final [%s]\n", time.Since(start))
		}
		return o
	} else {
		// Regular FFT
		start = time.Now()
		res := self._fft(vals, rootz[0:len(rootz)-1])
		if time.Since(start).Seconds() > MIN_SECONDS_BENCHMARK {
			//	fmt.Printf("    reg_fft core [%s]\n", time.Since(start))
		}
		return res
	}
}

func (self *PrimeField) mul_polys_fft(a []*big.Int, b []*big.Int, root_of_unity *big.Int) []*big.Int {
	x1 := self.fft(a, root_of_unity, false)
	x2 := self.fft(b, root_of_unity, false)
	c := make([]*big.Int, len(x1))
	t := new(big.Int)
	for i, v1 := range x1 {
		t.Mul(v1, x2[i])
		c[i] = new(big.Int).Mod(t, self.modulus)
	}
	return self.fft(c, root_of_unity, true)
}

// Optimized version of the above restricted to deg-4 polynomials
func _multi_interp_4(ch chan IndexedIntArrays, chinvtarget chan IndexedInt, _xsets [][]*big.Int, _ysets [][]*big.Int, m *big.Int, i0 int) {
	f, _ := NewPrimeField(nil)
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
	t3 := new(big.Int)
	t2 := new(big.Int)
	sum_xs2_xs3 := new(big.Int)
	sub_nxs0_xs1 := new(big.Int)
	for _j, xs := range _xsets {
		j := _j + i0
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

		// eq[0] := create_poly4(-x12 * xs[3] % m, (x12 + x13 + x23), -xs[1]-xs[2]-xs[3], 1)
		// eq[1] := create_poly4(-x02 * xs[3] % m, (x02 + x03 + x23), -xs[0]-xs[2]-xs[3], 1)
		// eq[2] := create_poly4(-x01 * xs[3] % m, (x01 + x03 + x13), -xs[0]-xs[1]-xs[3], 1)
		// eq[3] := create_poly4(-x01 * xs[2] % m, (x01 + x02 + x12), -xs[0]-xs[1]-xs[2], 1)
		eq := [][]*big.Int{
			create_poly4(new(big.Int).Mod(t3.Mul(nxs3, x12), m), new(big.Int).Add(t2.Add(x12, x13), x23), new(big.Int).Sub(nxs1, sum_xs2_xs3), ONE),
			create_poly4(new(big.Int).Mod(t3.Mul(x02, nxs3), m), new(big.Int).Add(t2.Add(x02, x03), x23), new(big.Int).Sub(nxs0, sum_xs2_xs3), ONE),
			create_poly4(new(big.Int).Mod(t3.Mul(x01, nxs3), m), new(big.Int).Add(t2.Add(x01, x03), x13), new(big.Int).Sub(sub_nxs0_xs1, xs[3]), ONE),
			create_poly4(new(big.Int).Mod(t3.Mul(x01, nxs2), m), new(big.Int).Add(t2.Add(x01, x02), x12), new(big.Int).Sub(sub_nxs0_xs1, xs[2]), ONE),
		}

		ch <- IndexedIntArrays{idx: j, data: [][]*big.Int{_ysets[_j], eq[0], eq[1], eq[2], eq[3]}}
		for k := 0; k < 4; k++ {
			chinvtarget <- IndexedInt{idx: j*4 + k, data: f.eval_quartic(eq[k], xs[k])}
		}
	}
}

func _multi_interp_4_part2(ch chan IndexedIntArray, data [][][]*big.Int, invalls []*big.Int, i0 int) {
	f, _ := NewPrimeField(nil)
	t0 := new(big.Int)
	t1 := new(big.Int)
	t2 := new(big.Int)
	t3 := new(big.Int)
	inv_y0 := new(big.Int)
	inv_y1 := new(big.Int)
	inv_y2 := new(big.Int)
	inv_y3 := new(big.Int)
	for _i, d := range data {
		i := i0 + _i
		ys, eq0, eq1, eq2, eq3 := d[0], d[1], d[2], d[3], d[4]
		invallz := invalls[i*4 : i*4+4]
		inv_y0.Mod(t0.Mul(ys[0], invallz[0]), f.modulus)
		inv_y1.Mod(t1.Mul(ys[1], invallz[1]), f.modulus)
		inv_y2.Mod(t2.Mul(ys[2], invallz[2]), f.modulus)
		inv_y3.Mod(t3.Mul(ys[3], invallz[3]), f.modulus)
		e := make([]*big.Int, 4)
		// [(eq0[i] * inv_y0 + eq1[i] * inv_y1 + eq2[i] * inv_y2 + eq3[i] * inv_y3) % m for i in range(4)])
		for j := 0; j < 4; j++ {
			e[j] = new(big.Int).Mod(t1.Add(t0.Add(t1.Mul(eq0[j], inv_y0), t2.Mul(eq1[j], inv_y1)), t3.Add(t1.Mul(eq2[j], inv_y2), t2.Mul(eq3[j], inv_y3))), f.modulus)
		}
		ch <- IndexedIntArray{idx: i, data: e}
	}
}

func (self *PrimeField) multi_interp_4(xsets, ysets [][]*big.Int) [][]*big.Int {
	data := make([][][]*big.Int, len(xsets))
	invtargets := make([]*big.Int, len(xsets)*4)
	nev := len(xsets)
	nj := min_iterations(nev, 3)
	ch := make(chan IndexedIntArrays, nev)
	ch2 := make(chan IndexedIntArray, nev)
	chinvtarget := make(chan IndexedInt, nev*4)
	for i := 0; i < nev; i += nj {
		i1 := i + nj
		if i1 > nev {
			i1 = nev
		}
		go _multi_interp_4(ch, chinvtarget, xsets[i:i1], ysets[i:i1], self.modulus, i)
	}

	for i := 0; i < nev; i++ {
		d := <-ch
		data[d.idx] = d.data
		for k := 0; k < 4; k++ {
			e := <-chinvtarget
			invtargets[e.idx] = e.data
		}
	}

	invalls := self.multi_inv(invtargets)

	o := make([][]*big.Int, nev)
	for i := 0; i < nev; i += nj {
		i1 := i + nj
		if i1 > nev {
			i1 = nev
		}
		go _multi_interp_4_part2(ch2, data[i:i1], invalls, i)
	}
	for i := 0; i < nev; i++ {
		d := <-ch2
		o[d.idx] = d.data
	}
	// TODO: assert o == [self.lagrange_interp_4(xs, ys) for xs, ys in zip(xsets, ysets)]
	return o
}

// Compute a MIMC permutation for some number of steps
func (self *PrimeField) MiMC(input *big.Int, steps *big.Int, round_constants []*big.Int) *big.Int {
	inp := new(big.Int).Set(input)
	t1 := new(big.Int)
	t2 := new(big.Int)
	nsteps := int(steps.Int64() - 1)
	for i := 0; i < nsteps; i++ {
		// inp = (inp**3 + round_constants[i % len(round_constants)]) % modulus
		t1.Exp(inp, THREE, nil)
		t2.Add(t1, round_constants[i%len(round_constants)])
		inp.Mod(t2, self.modulus)
	}
	return inp
}
