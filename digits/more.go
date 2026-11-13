package digits

import (
	"math/bits"
	"strconv"
	"strings"

	"github.com/t1sk3/eulerlib/v4/internal/mathx"
	"github.com/t1sk3/go-utils/types/constraints"
)

// Rotations returns the rotations of the decimal digits of |n|, starting
// with |n| itself and moving the leading digit to the end each time:
// Rotations(197) is [197 971 719] (Project Euler #35). There is one per
// digit, so repeated digits give repeated values. A rotation that starts
// with zeros loses them: Rotations(109) is [109 91 910].
func Rotations[E constraints.Integer](n E) []E {
	u := mathx.Abs64(n)
	c := Count(u)
	pow := uint64(1) // the place value of the leading digit
	for range c - 1 {
		pow *= 10
	}
	out := make([]E, c)
	for i := range out {
		out[i] = E(u)
		u = u%pow*10 + u/pow
	}
	return out
}

// Signature returns how often each decimal digit occurs in |n|:
// Signature(n)[d] is the number of digits d. Two numbers are digit
// permutations of each other exactly when their signatures are equal, and
// a signature is comparable, so it works as a map key for grouping numbers
// by their digits (Project Euler #49, #62).
func Signature[E constraints.Integer](n E) [10]uint8 {
	var sig [10]uint8
	u := mathx.Abs64(n)
	for {
		sig[u%10]++
		u /= 10
		if u == 0 {
			return sig
		}
	}
}

// IsPermutation reports whether |a| and |b| have the same decimal digits
// in some order (Project Euler #52, #70).
func IsPermutation[E constraints.Integer](a, b E) bool {
	return Signature(a) == Signature(b)
}

// PowerSum returns the sum of the k-th powers of the decimal digits of
// |n|: PowerSum(1634, 4) = 1 + 1296 + 81 + 256 = 1634 (Project Euler #30;
// k = 2 is the step of the happy-number sequence). 0^0 counts as 1. It
// panics if k < 0.
func PowerSum[E constraints.Integer](n E, k int) E {
	if k < 0 {
		panic("digits: PowerSum with a negative exponent")
	}
	var pow [10]E
	for d := range pow {
		pow[d] = 1
		for range k {
			pow[d] *= E(d)
		}
	}
	u := mathx.Abs64(n)
	var s E
	for {
		s += pow[u%10]
		u /= 10
		if u == 0 {
			return s
		}
	}
}

// LookAndSay returns the next term of the look-and-say sequence after s:
// each run of equal characters is replaced by its length and the
// character, so "1211" becomes "111221" (one 1, one 2, two 1s).
func LookAndSay(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); {
		j := i
		for j < len(s) && s[j] == s[i] {
			j++
		}
		b.WriteString(strconv.Itoa(j - i))
		b.WriteByte(s[i])
		i = j
	}
	return b.String()
}

// ChampernowneDigit returns the n-th digit (counting from 1) after the
// decimal point of Champernowne's constant 0.123456789101112..., the
// positive integers written one after the other (Project Euler #40). It
// panics if n < 1.
func ChampernowneDigit(n int) int {
	if n < 1 {
		panic("digits: ChampernowneDigit needs n >= 1")
	}
	// Skip the blocks of 9 one-digit numbers, 90 two-digit numbers, ...
	rest := uint64(n) - 1 // 0-based offset into the current block
	width, count, first := uint64(1), uint64(9), uint64(1)
	for {
		hi, size := bits.Mul64(width, count)
		if hi != 0 || rest < size {
			break
		}
		rest -= size
		width++
		count *= 10
		first *= 10
	}
	num := first + rest/width
	for range width - 1 - rest%width {
		num /= 10
	}
	return int(num % 10)
}
