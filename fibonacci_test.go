package eulerlib

import (
	"math/big"
	"testing"
)

func collectInt64(ch <-chan int64) []int64 {
	res := []int64{}
	for v := range ch {
		res = append(res, v)
	}
	return res
}

func collectBigInt(ch <-chan big.Int) []big.Int {
	res := []big.Int{}
	for v := range ch {
		res = append(res, v)
	}
	return res
}

func TestGenFibo(t *testing.T) {
	got := collectInt64(GenFibo(3))
	want := []int64{1, 2, 3, 5, 8, 13, 21}
	if len(got) != len(want) {
		t.Fatalf("GenFibo(3) length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("GenFibo(3)[%d] = %d, want %d", i, got[i], want[i])
		}
	}
}

func TestFibonacci(t *testing.T) {
	got := Fibonacci(0)
	want := []int64{1, 2, 3, 5}
	if len(got) != len(want) {
		t.Fatalf("Fibonacci(0) length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Fibonacci(0)[%d] = %d, want %d", i, got[i], want[i])
		}
	}
}

func TestFibonacciSingle(t *testing.T) {
	testCases := []struct {
		n    int64
		want int64
	}{
		{0, 0},
		{1, 1},
		{2, 1},
		{10, 55},
	}

	for _, tc := range testCases {
		got := FibonacciSingle(tc.n)
		if got != tc.want {
			t.Errorf("FibonacciSingle(%d) = %d, want %d", tc.n, got, tc.want)
		}
	}
}

func TestGenFiboBig(t *testing.T) {
	got := collectBigInt(GenFiboBig(0))
	want := []int64{0, 1, 1, 2}
	if len(got) != len(want) {
		t.Fatalf("GenFiboBig(0) length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i].Int64() != want[i] {
			t.Fatalf("GenFiboBig(0)[%d] = %d, want %d", i, got[i].Int64(), want[i])
		}
	}
}

func TestFibonacciBig(t *testing.T) {
	got := FibonacciBig(0)
	want := []int64{0, 1, 1, 2, 3, 5}
	if len(got) != len(want) {
		t.Fatalf("FibonacciBig(0) length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i].Int64() != want[i] {
			t.Fatalf("FibonacciBig(0)[%d] = %d, want %d", i, got[i].Int64(), want[i])
		}
	}
}

func TestFibonacciSingleBig(t *testing.T) {
	got := FibonacciSingleBig(10)
	if got.Int64() != 55 {
		t.Fatalf("FibonacciSingleBig(10) = %d, want 55", got.Int64())
	}
}

func TestPowBigFloat(t *testing.T) {
	base := big.NewFloat(3)
	got := powBigFloat(base, 4)
	f, _ := got.Float64()
	if f != 81 {
		t.Fatalf("powBigFloat(3, 4) = %v, want 81", f)
	}
}

func TestZero(t *testing.T) {
	z := zero()
	if z.Prec() != 256 {
		t.Fatalf("zero().Prec() = %d, want 256", z.Prec())
	}
	f, _ := z.Float64()
	if f != 0 {
		t.Fatalf("zero() = %v, want 0", f)
	}
}
