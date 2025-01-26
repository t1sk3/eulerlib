package eulerlib

import (
	"testing"

	"github.com/fxtlabs/primes"
)

func TestPrimeNumbers(t *testing.T) {
	for i := 1; i < 1000000; i++ {
		got := IsPrime(i)
		want := primes.IsPrime(i)
		if got != want {
			t.Errorf("IsPrime(%d) == %t, want %t", i, got, want)
		}
	}
}

func TestNextPrime(t *testing.T) {
	testNums := []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 20}
	want := []int64{2, 2, 3, 5, 5, 7, 7, 11, 11, 23}
	for i, num := range testNums {
		got := NextPrime(num)
		if got != want[i] {
			t.Errorf("NextPrime(%d) == %d, want %d", num, got, want[i])
		}
	}
}

func TestPrimeNumberIterator(t *testing.T) {
	iter := NewPrimeNumberIterator(5000)
	for iter.Current() < 1000000 {
		want := NextPrime(iter.Current())
		got := iter.Next()
		if got != want {
			t.Errorf("NextPrime(%d) == %d, want %d", iter.Current(), got, want)
		}
	}
}
