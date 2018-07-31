package stark

import (
	"fmt"
	"math/big"
	"strings"
	"testing"
)

func TestFFT(t *testing.T) {
	f, _ := NewPrimeField(big.NewInt(337))
	vals := []*big.Int{big.NewInt(3), big.NewInt(1), big.NewInt(4), big.NewInt(1), big.NewInt(5), big.NewInt(9), big.NewInt(2), big.NewInt(6)}
	root_of_unity := big.NewInt(85)

	// fft.fft([3,1,4,1,5,9,2,6], 337, 85, inv=True) => [46 169 29 149 126 262 140 93]
	expected := "[46 169 29 149 126 262 140 93]"
	res := f.fft(vals, root_of_unity, true)
	str := fmt.Sprintf("%s", res)
	if strings.Compare(str, expected) != 0 {
		t.Fatalf("FFT failure")
	}
	for i := int64(0); i < 8; i++ {
		ti := big.NewInt(i)
		q := f.eval_poly_at(res, f.pow(big.NewInt(85), ti))
		if q.Cmp(vals[i]) != 0 {
			t.Fatalf("FFT inv Failure")
		} else {
			fmt.Printf("r(%v) = %v\n", ti, q)
		}
	}
	fmt.Printf("fft.fft([3,1,4,1,5,9,2,6], 337, 85, inv=True)=%v SUCCESS\n", res)

	// p1 = 4 + 5x, p2 = 1 + 2x + 3*x*x
	p1 := []*big.Int{big.NewInt(4), big.NewInt(5)}
	p2 := []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}
	p_mul := f.mul_polys_fft(p1, p2, root_of_unity)
	p_mul2 := f.mul_polys(p1, p2)
	for i, coeff2 := range p_mul2 {
		if coeff2.Cmp(p_mul[i]) != 0 {
			fmt.Printf("mul_polys_fft failure\n")
		}
	}
	fmt.Printf("p_mul: %v p_mul2: %v SUCCESS\n", p_mul, p_mul2)
}

func TestPrimeField(t *testing.T) {
	f, _ := NewPrimeField(big.NewInt(337))

	// add test
	q := f.add(big.NewInt(300), big.NewInt(40))
	if q.Cmp(big.NewInt(3)) != 0 {
		t.Fatalf("add problem\n")
	}
	fmt.Printf("add(300,40) == 3 PASS\n")

	// sub test
	q = f.sub(big.NewInt(1), big.NewInt(3))
	if q.Cmp(big.NewInt(335)) != 0 {
		t.Fatalf("add problem\n")
	}
	fmt.Printf("sub(1,3) == 335 PASS\n")

	// mul test
	q = f.mul(big.NewInt(5), big.NewInt(100))
	if q.Cmp(big.NewInt(163)) != 0 {
		t.Fatalf("mul problem\n")
	}
	fmt.Printf("mul(5,100) == 163 PASS\n")

	// pow test
	q = f.pow(big.NewInt(2), big.NewInt(9))
	if q.Cmp(big.NewInt(175)) != 0 {
		t.Fatalf("pow problem\n")
	}
	fmt.Printf("pow(2,9) == 175 PASS\n")

	// inv test
	f.SetModulus(big.NewInt(31))
	q = f.inv(big.NewInt(12))
	fmt.Printf("inv(12) == %v PASS\n", q)
	if q.Cmp(big.NewInt(13)) != 0 {
		t.Fatalf("inv problem\n")
	}

	// div test (NOT WORKING)
	q = f.div(big.NewInt(10), big.NewInt(2))
	if q.Cmp(big.NewInt(5)) != 0 {
		t.Fatalf("pow problem\n")
	}
	fmt.Printf("div(10,2) == 2 PASS\n")

	// eval_poly_at test
	p := []*big.Int{big.NewInt(4), big.NewInt(5), big.NewInt(6)}
	q = f.eval_poly_at(p, big.NewInt(2))
	if q.Cmp(big.NewInt(7)) != 0 {
		t.Fatalf("pow problem\n")
	}
	// 6 * 22 + 5 * 2 + 4 = 38, 38 mod 31 = 7.
	fmt.Printf("p: %v p(2) = 7\n", p)

	// zpoly / eval_poly_at test
	xs := []*big.Int{big.NewInt(4), big.NewInt(5), big.NewInt(6)}
	p = f.zpoly(xs)
	fmt.Printf("p: %v\n", p)
	for _, x := range xs {
		q = f.eval_poly_at(p, x)
		if q.Cmp(big.NewInt(0)) != 0 {
			t.Fatalf("zpoly problem\n")
		}
		fmt.Printf("p(%v) = %v\n", x, q)
	}
	fmt.Printf("\n\n")

	//	a(x) = 5x^2 + 4x + 6, b = 2x + 1
	a := []*big.Int{big.NewInt(6), big.NewInt(4), big.NewInt(5)}
	b := []*big.Int{big.NewInt(1), big.NewInt(2)}
	f.SetModulus(big.NewInt(7))
	p_div, _ := f.div_polys(a, b)

	// test add_polys, sub_polys, mul_polys, div_polys, mul_by_const
	// p1 = 4 + 5x, p2 = 1 + 2x + 3*x*x
	p1 := []*big.Int{big.NewInt(4), big.NewInt(5)}
	p2 := []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}
	p_add := f.add_polys(p1, p2)             // 5 + 7x + 3*x*x  5+7*4+3*64
	p_sub := f.sub_polys(p1, p2)             // 3 + 3x - 3*x*x
	p_mul := f.mul_polys(p1, p2)             //
	p_3 := f.mul_by_const(p2, big.NewInt(3)) // 3 + 6x + 9*x*x
	for _, x := range xs {
		q_p1 := f.eval_poly_at(p1, x)
		q_p2 := f.eval_poly_at(p2, x)
		q_add := f.eval_poly_at(p_add, x)
		q_sub := f.eval_poly_at(p_sub, x)
		q_mul := f.eval_poly_at(p_mul, x)
		q_div := f.eval_poly_at(p_div, x)
		q_3 := f.eval_poly_at(p_3, x)

		fmt.Printf("p1=%v (%v) = %v\n", p1, x, q_p1)
		fmt.Printf("p2=%v (%v) = %v\n", p2, x, q_p2)
		fmt.Printf("p1+p2=%v p1+p2(%v) = %v\n", p_add, x, q_add)
		fmt.Printf("p1-p2=%v p1-p2(%v) = %v\n", p_sub, x, q_sub)
		fmt.Printf("p1*p2=%v p1*p2(%v) = %v\n", p_mul, x, q_mul)
		fmt.Printf("p2/p1=%v p2/p1(%v) = %v\n", p_div, x, q_div)
		fmt.Printf("3*p2=%v p1-p2(%v) = %v\n\n", p_3, x, q_3)
	}

	// lagrange_interp
	xs0 := []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}
	ys0 := []*big.Int{big.NewInt(1), big.NewInt(4), big.NewInt(9)}
	f.SetModulus(big.NewInt(337))
	lp := f.lagrange_interp(xs0, ys0)
	lp2 := f.lagrange_interp_2(xs0, ys0)
	xs0 = append(xs0, big.NewInt(4))
	ys0 = append(ys0, big.NewInt(16))
	lp4 := f.lagrange_interp_4(xs0, ys0)
	fmt.Printf("lp=%v\n", lp)
	fmt.Printf("lp2=%v\n", lp2)
	fmt.Printf("lp4=%v\n", lp4)
	for _, x := range xs0 {
		q := f.eval_poly_at(lp, x)
		q2 := f.eval_poly_at(lp2, x)
		q4 := f.eval_poly_at(lp4, x)
		fmt.Printf("lp(%v) =%10s\t", x, q)
		fmt.Printf("lp2(%v)=%10s\t", x, q2)
		fmt.Printf("lp4(%v)=%10s\n", x, q4)
	}
}
