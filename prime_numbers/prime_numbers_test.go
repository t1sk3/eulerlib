package prime_numbers

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

func TestListPrimes(t *testing.T) {
	testCases := []struct {
		n    int64
		want []int64
	}{
		{0, []int64{}},
		{1, []int64{}},
		{2, []int64{2}},
		{3, []int64{2, 3}},
		{4, []int64{2, 3}},
		{5, []int64{2, 3, 5}},
	}
	for _, tc := range testCases {
		got := ListPrimes(tc.n)
		if len(got) != len(tc.want) {
			t.Errorf("ListPrimes(%d) = %v, want %v", tc.n, got, tc.want)
			continue
		}
		for i := range got {
			if got[i] != tc.want[i] {
				t.Errorf("ListPrimes(%d)[%d] = %d, want %d", tc.n, i, got[i], tc.want[i])
			}
		}
	}
}

// check prime generotor channel
func TestPrimeGenerator(t *testing.T) {
	limit := int64(1000)
	generator := PrimeGenerator(limit)

	count := 0
	for prime := range generator {
		if prime > limit {
			t.Errorf("PrimeGenerator exceeded limit: %d > %d", prime, limit)
		}
		count++
	}

	if count == 0 {
		t.Errorf("PrimeGenerator did not yield any primes")
	}
	// Check if the last prime is less than or equal to the limit
	if prime := <-generator; prime > limit {
		t.Errorf("Last prime from PrimeGenerator exceeded limit: %d > %d", prime, limit)
	}
	// Check if the generator is closed properly
	if _, ok := <-generator; ok {
		t.Errorf("PrimeGenerator channel is not closed properly")
	}
	// Check if the generator produces unique primes
	seen := make(map[int64]bool)
	for prime := range generator {
		if seen[prime] {
			t.Errorf("PrimeGenerator produced duplicate prime: %d", prime)
		}
		seen[prime] = true
	}
	lastPrime := int64(999999999999)
	// Check if the generator produces primes in ascending order
	for prime := range generator {
		if len(seen) > 0 && prime <= lastPrime {
			t.Errorf("PrimeGenerator produced primes out of order: %d <= %d", prime, lastPrime)
		}
		lastPrime = prime
	}
	// Check if the generator produces primes up to the limit
	for prime := range generator {
		if prime > limit {
			t.Errorf("PrimeGenerator produced prime out of range: %d > %d", prime, limit)
		}
	}
	// Check if the generator produces primes in ascending order
	for prime := range generator {
		if len(seen) > 0 && prime <= lastPrime {
			t.Errorf("PrimeGenerator produced primes out of order: %d <= %d", prime, lastPrime)
		}
		lastPrime = prime
	}
	// Check if the generator produces primes in ascending order
	for prime := range generator {
		if len(seen) > 0 && prime <= lastPrime {
			t.Errorf("PrimeGenerator produced primes out of order: %d <= %d", prime, lastPrime)
		}
		lastPrime = prime
	}
	// Check if the generator produces primes in ascending order
	for prime := range generator {
		if len(seen) > 0 && prime <= lastPrime {
			t.Errorf("PrimeGenerator produced primes out of order: %d <= %d", prime, lastPrime)
		}
		lastPrime = prime
	}
}

func TestSumPrimes(t *testing.T) {
	limit := 2000000
	got := SumPrimes(0, limit)
	want := 0
	for i := 2; i < limit; i++ {
		if IsPrime(i) {
			want += i
		}
	}
	if got != want {
		t.Errorf("SumPrimes(%d) == %d, want %d", limit, got, want)
	}

	got = SumPrimes(limit/2, limit)
	want = 0
	for i := limit / 2; i < limit; i++ {
		if IsPrime(i) {
			want += i
		}
	}
	if got != want {
		t.Errorf("SumPrimes(%d, %d) == %d, want %d", limit/2, limit, got, want)
	}
	got = SumPrimes(0, 0)
	want = 0
	if got != want {
		t.Errorf("SumPrimes(0, 0) == %d, want %d", got, want)
	}
	got = SumPrimes(0, 1)
	want = 0
	if got != want {
		t.Errorf("SumPrimes(0, 1) == %d, want %d", got, want)
	}
}

func TestListSmallestPrimeFactors(t *testing.T) {
	const n = 1000
	spf := ListSmallestPrimeFactors(n)

	if len(spf) != n+1 {
		t.Fatalf("len(spf) == %d, want %d", len(spf), n+1)
	}
	if spf[0] != 0 || spf[1] != 0 {
		t.Errorf("spf[0], spf[1] == %d, %d, want 0, 0", spf[0], spf[1])
	}
	for i := 2; i <= n; i++ {
		p := spf[i]
		if !IsPrime(p) {
			t.Fatalf("spf[%d] = %d is not prime", i, p)
		}
		if i%p != 0 {
			t.Fatalf("spf[%d] = %d does not divide %d", i, p, i)
		}
		for q := 2; q < p; q++ {
			if i%q == 0 {
				t.Fatalf("spf[%d] = %d, but %d is a smaller factor", i, p, q)
			}
		}
	}
}
