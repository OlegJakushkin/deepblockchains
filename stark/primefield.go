package stark

import (
	"fmt"
	"math/big"
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
	return o.Mod(o, self.modulus)
}

func (self *PrimeField) sub(x, y *big.Int) (o *big.Int) {
	o = new(big.Int).Sub(x, y)
	return o.Mod(o, self.modulus)
}

func (self *PrimeField) mul(x, y *big.Int) (o *big.Int) {
	o = new(big.Int).Mul(x, y)
	return o.Mod(o, self.modulus)
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
		r := new(big.Int)
		r.Div(high, low)

		// nm = hm-lm*r
		t := new(big.Int)
		nm := new(big.Int)
		nm.Sub(hm, t.Mul(lm, r))

		// nw = high - low*r
		t.Mul(low, r)
		nw := new(big.Int)
		nw.Sub(high, t)

		lm, low, hm, high = nm, nw, lm, low
	}
	return lm.Mod(lm, self.modulus)
}

func (self *PrimeField) multi_inv(values []*big.Int) []*big.Int {
	partials := make([]*big.Int, 0)
	partials = append(partials, big.NewInt(1))
	outputs := make([]*big.Int, len(values))

	for i := 0; i < len(values); i++ {
		if values[i].Cmp(big.NewInt(0)) == 0 {
			partials = append(partials, self.mul(partials[len(partials)-1], big.NewInt(1)))
		} else {
			partials = append(partials, self.mul(partials[len(partials)-1], values[i]))
		}
		outputs[i] = big.NewInt(0)
	}
	inv := self.inv(partials[len(partials)-1])
	for i := len(values); i > 0; i = i - 1 {
		if values[i-1].Cmp(big.NewInt(0)) == 0 {
			outputs[i-1] = big.NewInt(0)
		} else {
			outputs[i-1] = self.mul(partials[i-1], inv)
			inv = self.mul(inv, values[i-1])
		}
	}
	return outputs
}

// Evaluate a polynomial at a point
func (self *PrimeField) eval_poly_at(poly []*big.Int, x *big.Int) *big.Int {
	o := big.NewInt(0)
	p := big.NewInt(1)
	for _, coeff := range poly {
		// o += coeff * p
		o.Add(o, new(big.Int).Mul(p, coeff))
		t := new(big.Int).Mul(p, x)
		p = t.Mod(t, self.modulus)
	}
	return o.Mod(o, self.modulus)
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
	for _, x := range xs {
		root = append([]*big.Int{big.NewInt(0)}, root...)
		for j := 0; j < len(root)-1; j++ {
			t := new(big.Int)
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
	for i, x := range xs {
		p := make([]*big.Int, 2)
		p[0] = new(big.Int).Set(x)
		p[0].Mul(p[0], big.NewInt(-1))
		p[1] = big.NewInt(1)
		nums[i], _ = self.div_polys(root, p)
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
	o = append(o, a)
	o = append(o, b)
	return o
}

func create_poly4(a, b, c, d *big.Int) (o []*big.Int) {
	o = append(o, a)
	o = append(o, b)
	o = append(o, c)
	o = append(o, d)
	return o
}

// Optimized version of the above restricted to deg-2 polynomials
func (self *PrimeField) lagrange_interp_2(xs, ys []*big.Int) []*big.Int {
	m := self.modulus

	// eq0 = [-xs[1] % m, 1]
	t1 := new(big.Int).Set(xs[1])
	t1.Mul(t1, big.NewInt(-1))
	t1.Mod(t1, m)
	eq0 := create_poly2(t1, big.NewInt(1))
	e0 := self.eval_poly_at(eq0, xs[0])

	// eq1 = [-xs[0] % m, 1]
	t0 := new(big.Int).Set(xs[0])
	t0.Mul(t0, big.NewInt(-1))
	t0.Mod(t0, m)
	eq1 := create_poly2(t0, big.NewInt(1))
	e1 := self.eval_poly_at(eq1, xs[1])

	// invall = self.inv(e0 * e1)
	t := new(big.Int).Mul(e0, e1)
	invall := self.inv(t)

	// inv_y0 = ys[0] * invall * e1
	inv_y0 := new(big.Int).Set(ys[0])
	inv_y0.Mul(inv_y0, invall)
	inv_y0.Mul(inv_y0, e1)

	// inv_y1 = ys[1] * invall * e0
	inv_y1 := new(big.Int).Set(ys[1])
	inv_y1.Mul(inv_y1, invall)
	inv_y1.Mul(inv_y1, e0)

	// [(eq0[i] * inv_y0 + eq1[i] * inv_y1) % m for i in range(2)]
	o := make([]*big.Int, 2)
	for i := 0; i < 2; i++ {
		t0 := new(big.Int).Mul(eq0[i], inv_y0)
		t1 := new(big.Int).Mul(eq1[i], inv_y1)
		s := new(big.Int).Add(t0, t1)
		s.Mod(s, m)
		o[i] = s
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
	o.Add(o, o2)
	o.Add(o, o3)
	return o.Mod(o, self.modulus)
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

	//eq0 := create_poly4(-x12*xs[3]%m, (x12 + x13 + x23), -xs[1]-xs[2]-xs[3], 1)
	t3 := new(big.Int).Set(x12)
	t3.Mul(t3, xs[3])
	t3.Mul(t3, big.NewInt(-1))
	t3.Mod(t3, m)
	t2 := new(big.Int).Set(x12)
	t2.Add(t2, x13)
	t2.Add(t2, x23)
	t1 := big.NewInt(0)
	t1.Sub(t1, xs[1])
	t1.Sub(t1, xs[2])
	t1.Sub(t1, xs[3])
	eq0 := create_poly4(t3, t2, t1, big.NewInt(1))

	//eq1 := create_poly4(-x02*xs[3]%m, (x02 + x03 + x23), -xs[0]-xs[2]-xs[3], 1)
	t3 = new(big.Int).Set(x02)
	t3.Mul(t3, xs[3])
	t3.Mul(t3, big.NewInt(-1))
	t3.Mod(t3, m)
	t2 = new(big.Int).Set(x02)
	t2.Add(t2, x03)
	t2.Add(t2, x23)
	t1 = big.NewInt(0)
	t1.Sub(t1, xs[0])
	t1.Sub(t1, xs[2])
	t1.Sub(t1, xs[3])
	eq1 := create_poly4(t3, t2, t1, big.NewInt(1))

	//eq2 := create_poly4(-x01*xs[3]%m, (x01 + x03 + x13), -xs[0]-xs[1]-xs[3], 1)
	t3 = new(big.Int).Set(x01)
	t3.Mul(t3, xs[3])
	t3.Mul(t3, big.NewInt(-1))
	t3.Mod(t3, m)
	t2 = new(big.Int).Set(x01)
	t2.Add(t2, x03)
	t2.Add(t2, x13)
	t1 = big.NewInt(0)
	t1.Sub(t1, xs[0])
	t1.Sub(t1, xs[1])
	t1.Sub(t1, xs[3])
	eq2 := create_poly4(t3, t2, t1, big.NewInt(1))

	//eq3 := create_poly4(-x01*xs[2]%m, (x01 + x02 + x12), -xs[0]-xs[1]-xs[2], 1)
	t3 = new(big.Int).Set(x01)
	t3.Mul(t3, xs[2])
	t3.Mul(t3, big.NewInt(-1))
	t3.Mod(t3, m)
	t2 = new(big.Int).Set(x01)
	t2.Add(t2, x02)
	t2.Add(t2, x12)
	t1 = big.NewInt(0)
	t1.Sub(t1, xs[0])
	t1.Sub(t1, xs[1])
	t1.Sub(t1, xs[2])
	eq3 := create_poly4(t3, t2, t1, big.NewInt(1))

	e0 := self.eval_poly_at(eq0, xs[0])
	e1 := self.eval_poly_at(eq1, xs[1])
	e2 := self.eval_poly_at(eq2, xs[2])
	e3 := self.eval_poly_at(eq3, xs[3])
	e01 := new(big.Int).Mul(e0, e1)
	e23 := new(big.Int).Mul(e2, e3)
	invall := self.inv(new(big.Int).Mul(e01, e23))

	// inv_y0 := ys[0] * invall * e1 * e23 % m
	inv_y0 := new(big.Int).Set(ys[0])
	inv_y0.Mul(inv_y0, invall)
	inv_y0.Mul(inv_y0, e1)
	inv_y0.Mul(inv_y0, e23)
	inv_y0.Mod(inv_y0, m)

	// inv_y1 := ys[1] * invall * e0 * e23 % m
	inv_y1 := new(big.Int).Set(ys[1])
	inv_y1.Mul(inv_y1, invall)
	inv_y1.Mul(inv_y1, e0)
	inv_y1.Mul(inv_y1, e23)
	inv_y1.Mod(inv_y1, m)

	// inv_y2 := ys[2] * invall * e01 * e3 % m
	inv_y2 := new(big.Int).Set(ys[2])
	inv_y2.Mul(inv_y2, invall)
	inv_y2.Mul(inv_y2, e01)
	inv_y2.Mul(inv_y2, e3)
	inv_y2.Mod(inv_y2, m)

	// inv_y3 := ys[3] * invall * e01 * e2 % m
	inv_y3 := new(big.Int).Set(ys[3])
	inv_y3.Mul(inv_y3, invall)
	inv_y3.Mul(inv_y3, e01)
	inv_y3.Mul(inv_y3, e2)
	inv_y3.Mod(inv_y3, m)

	o := make([]*big.Int, 4)
	for i := 0; i < 4; i++ {
		// [(eq0[i] * inv_y0 + eq1[i] * inv_y1 + eq2[i] * inv_y2 + eq3[i] * inv_y3) % m for i in range(4)]
		t0 := new(big.Int).Mul(eq0[i], inv_y0)
		t1 := new(big.Int).Mul(eq1[i], inv_y1)
		t2 := new(big.Int).Mul(eq2[i], inv_y2)
		t3 := new(big.Int).Mul(eq3[i], inv_y3)
		s := new(big.Int).Add(t0, t1)
		s.Add(s, t2)
		s.Add(s, t3)
		o[i] = s.Mod(s, m)
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

	lvals := make([]*big.Int, len(vals)/2)
	rvals := make([]*big.Int, len(vals)/2)
	root2 := make([]*big.Int, len(roots_of_unity)/2)
	for i := 0; i < len(vals)/2; i++ {
		lvals[i] = vals[i*2]
		rvals[i] = vals[i*2+1]
	}
	for i := 0; i < len(roots_of_unity)/2; i++ {
		root2[i] = roots_of_unity[i*2]
	}
	L := self._fft(lvals, root2)
	R := self._fft(rvals, root2)
	o := make([]*big.Int, len(vals))
	for i := 0; i < len(vals); i++ {
		o[i] = big.NewInt(0)
	}
	for i, x := range L {
		y := R[i]
		y_times_root := new(big.Int).Mul(y, roots_of_unity[i])
		t := new(big.Int).Add(x, y_times_root)
		o[i] = t.Mod(t, self.modulus)
		t2 := new(big.Int).Sub(x, y_times_root)
		o[i+len(L)] = t2.Mod(t2, self.modulus)
	}
	return o
}

func (self *PrimeField) fft(vals []*big.Int, root_of_unity *big.Int, inv bool) []*big.Int {
	// Build up roots of unity
	rootz := make([]*big.Int, 2)
	rootz[0] = big.NewInt(1)
	rootz[1] = root_of_unity
	for rootz[len(rootz)-1].Cmp(big.NewInt(1)) != 0 {
		t := new(big.Int).Mul(rootz[len(rootz)-1], root_of_unity)
		t.Mod(t, self.modulus)
		rootz = append(rootz, t)
	}

	// Fill in vals with zeroes if needed
	if len(rootz) > len(vals)+1 {
		extrazeros := make([]*big.Int, (len(rootz) - len(vals) - 1))
		for i := 0; i < len(extrazeros); i++ {
			extrazeros[i] = big.NewInt(0)
		}
		vals = append(vals, extrazeros...)
	}

	if inv {
		// Inverse FFT
		t := new(big.Int).Sub(self.modulus, big.NewInt(2))
		invlen := new(big.Int).Exp(big.NewInt(int64(len(vals))), t, self.modulus)
		irootz := make([]*big.Int, 0)
		for i := len(rootz) - 1; i > 0; i-- {
			irootz = append(irootz, rootz[i])
		}
		res := self._fft(vals, irootz)
		o := make([]*big.Int, len(res))
		for i, x := range res {
			q := new(big.Int).Mul(x, invlen)
			q.Mod(q, self.modulus)
			o[i] = q
		}
		return o
	} else {
		// Regular FFT
		return self._fft(vals, rootz[0:len(rootz)-1])
	}
}

func (self *PrimeField) mul_polys_fft(a []*big.Int, b []*big.Int, root_of_unity *big.Int) []*big.Int {
	x1 := self.fft(a, root_of_unity, false)
	x2 := self.fft(b, root_of_unity, false)
	c := make([]*big.Int, len(x1))
	for i, v1 := range x1 {
		v2 := x2[i]
		t := new(big.Int).Mul(v1, v2)
		t.Mod(t, self.modulus)
		c[i] = t
	}
	return self.fft(c, root_of_unity, true)
}

// Optimized version of the above restricted to deg-4 polynomials
func (self *PrimeField) multi_interp_4(xsets, ysets [][]*big.Int) [][]*big.Int {
	data := make([][][]*big.Int, 0)
	invtargets := make([]*big.Int, 0)

	for i, xs := range xsets {
		ys := ysets[i]
		x01 := new(big.Int).Mul(xs[0], xs[1])
		x02 := new(big.Int).Mul(xs[0], xs[2])
		x03 := new(big.Int).Mul(xs[0], xs[3])
		x12 := new(big.Int).Mul(xs[1], xs[2])
		x13 := new(big.Int).Mul(xs[1], xs[3])
		x23 := new(big.Int).Mul(xs[2], xs[3])
		m := self.modulus

		// eq0 := create_poly4(-x12 * xs[3] % m, (x12 + x13 + x23), -xs[1]-xs[2]-xs[3], 1)
		t3 := new(big.Int).Mul(x12, big.NewInt(-1))
		t3.Mul(t3, xs[3])
		t3.Mod(t3, m)
		t2 := new(big.Int).Add(x12, x13)
		t2.Add(t2, x23)
		t1 := big.NewInt(0)
		t1.Sub(t1, xs[1])
		t1.Sub(t1, xs[2])
		t1.Sub(t1, xs[3])
		eq0 := create_poly4(t3, t2, t1, big.NewInt(1))

		// eq1 := create_poly4(-x02 * xs[3] % m, (x02 + x03 + x23), -xs[0]-xs[2]-xs[3], 1)
		t3 = new(big.Int).Mul(x02, big.NewInt(-1))
		t3.Mul(t3, xs[3])
		t3.Mod(t3, m)
		t2 = new(big.Int).Add(x02, x03)
		t2.Add(t2, x23)
		t1 = big.NewInt(0)
		t1.Sub(t1, xs[0])
		t1.Sub(t1, xs[2])
		t1.Sub(t1, xs[3])
		eq1 := create_poly4(t3, t2, t1, big.NewInt(1))

		// eq2 := create_poly4(-x01 * xs[3] % m, (x01 + x03 + x13), -xs[0]-xs[1]-xs[3], 1)
		t3 = new(big.Int).Mul(x01, big.NewInt(-1))
		t3.Mul(t3, xs[3])
		t3.Mod(t3, m)
		t2 = new(big.Int).Add(x01, x03)
		t2.Add(t2, x13)
		t1 = big.NewInt(0)
		t1.Sub(t1, xs[0])
		t1.Sub(t1, xs[1])
		t1.Sub(t1, xs[3])
		eq2 := create_poly4(t3, t2, t1, big.NewInt(1))

		// eq3 := create_poly4(-x01 * xs[2] % m, (x01 + x02 + x12), -xs[0]-xs[1]-xs[2], 1)
		t3 = new(big.Int).Mul(x01, big.NewInt(-1))
		t3.Mul(t3, xs[2])
		t3.Mod(t3, m)
		t2 = new(big.Int).Add(x01, x02)
		t2.Add(t2, x12)
		t1 = big.NewInt(0)
		t1.Sub(t1, xs[0])
		t1.Sub(t1, xs[1])
		t1.Sub(t1, xs[2])
		eq3 := create_poly4(t3, t2, t1, big.NewInt(1))

		e0 := self.eval_quartic(eq0, xs[0])
		e1 := self.eval_quartic(eq1, xs[1])
		e2 := self.eval_quartic(eq2, xs[2])
		e3 := self.eval_quartic(eq3, xs[3])
		d := make([][]*big.Int, 5)
		d[0], d[1], d[2], d[3], d[4] = ys, eq0, eq1, eq2, eq3
		//fmt.Printf("e0 %s e1 %s e2 %s e3 %s\n", e0, e1, e2, e3)
		data = append(data, d)
		invtargets = append(invtargets, e0)
		invtargets = append(invtargets, e1)
		invtargets = append(invtargets, e2)
		invtargets = append(invtargets, e3)

	}
	invalls := self.multi_inv(invtargets)

	o := make([][]*big.Int, 0)
	for i, d := range data {
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
			e0 := new(big.Int).Mul(eq0[j], inv_y0)
			e1 := new(big.Int).Mul(eq1[j], inv_y1)
			e2 := new(big.Int).Mul(eq2[j], inv_y2)
			e3 := new(big.Int).Mul(eq3[j], inv_y3)
			a := new(big.Int).Add(e0, e1)
			a.Add(a, e2)
			a.Add(a, e3)
			a.Mod(a, self.modulus)
			e[j] = a
		}
		o = append(o, e)
	}
	// TODO: assert o == [self.lagrange_interp_4(xs, ys) for xs, ys in zip(xsets, ysets)]
	return o
}

// Compute a MIMC permutation for some number of steps
func (self *PrimeField) mimc(inp *big.Int, steps *big.Int, round_constants []*big.Int) *big.Int {
	start := time.Now()
	t := new(big.Int).Set(inp)
	for i := 0; i < int(steps.Int64()-1); i++ {
		//			inp = (inp**3 + round_constants[i % len(round_constants)]) % modulus
		t.Exp(t, big.NewInt(3), nil)
		t.Add(t, round_constants[i%len(round_constants)])
		t.Mod(t, self.modulus)
	}
	fmt.Printf("MIMC computed in %s\n\n", time.Since(start))
	return t
}
