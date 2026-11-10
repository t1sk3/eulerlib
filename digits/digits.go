// Package digits works with the digits of integers: digit slices and sums,
// reversal, palindromes, pandigital checks, concatenation and base
// conversion.
//
// Functions that take a generic integer work on |n| for negative n (and are
// exact for the minimum value of a signed type) unless their documentation
// says otherwise. Functions that build a number (FromDigits, Reverse,
// Concat) silently overflow E when the result does not fit, like ordinary
// Go arithmetic.
package digits

import (
	"strings"

	"github.com/t1sk3/eulerlib/v4/internal/mathx"
	"github.com/t1sk3/go-utils/stream"
	"github.com/t1sk3/go-utils/types/constraints"
)

// alphabet maps digit values to characters for ToBase and
// IsPandigitalString, for bases up to 62. It is the same alphabet as
// eulerlib v3's etc.DecimalToBase.
const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func checkBase(base int, max int) {
	if base < 2 || base > max {
		panic("digits: base out of range")
	}
}

// Of returns the decimal digits of |n|, most significant first. Of(0) is
// [0].
func Of[E constraints.Integer](n E) []E {
	return OfBase(n, 10)
}

// OfBase returns the digits of |n| in the given base, most significant
// first. OfBase(0, b) is [0]. It panics if base < 2, or if base-1 does not
// fit in E.
func OfBase[E constraints.Integer](n E, base int) []E {
	if base < 2 || uint64(base-1) > uint64(mathx.MaxOf[E]()) {
		panic("digits: base out of range")
	}
	u, b := mathx.Abs64(n), uint64(base)
	if u == 0 {
		return []E{0}
	}
	var buf [64]E
	i := len(buf)
	for u > 0 {
		i--
		buf[i] = E(u % b)
		u /= b
	}
	out := make([]E, len(buf)-i)
	copy(out, buf[i:])
	return out
}

// FromDigits returns the number whose decimal digits are ds, most
// significant first: the inverse of Of. FromDigits(nil) is 0.
func FromDigits[E constraints.Integer](ds []E) E {
	var n E
	for _, d := range ds {
		n = n*10 + d
	}
	return n
}

// Sum returns the sum of the decimal digits of |n|.
func Sum[E constraints.Integer](n E) E {
	u := mathx.Abs64(n)
	var s E
	for u > 0 {
		s += E(u % 10)
		u /= 10
	}
	return s
}

// SumString returns the sum of the decimal digits in s, ignoring every
// other character. It works for digit strings of any length, such as the
// output of (*big.Int).String.
func SumString(s string) int64 {
	var res int64
	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			res += int64(ch - '0')
		}
	}
	return res
}

// Count returns the number of decimal digits of |n|. Count(0) is 1.
func Count[E constraints.Integer](n E) int {
	u := mathx.Abs64(n)
	c := 1
	for u >= 10 {
		c++
		u /= 10
	}
	return c
}

// Reverse returns n with its decimal digits reversed, keeping its sign:
// Reverse(1230) is 321 and Reverse(-12) is -21.
func Reverse[E constraints.Integer](n E) E {
	u := mathx.Abs64(n)
	var r E
	for u > 0 {
		r = r*10 + E(u%10)
		u /= 10
	}
	if n < 0 {
		return -r
	}
	return r
}

// IsPalindrome reports whether the decimal digits of |n| read the same
// forwards and backwards.
func IsPalindrome[E constraints.Integer](n E) bool {
	return IsPalindromeInBase(n, 10)
}

// IsPalindromeInBase reports whether the digits of |n| in the given base
// read the same forwards and backwards (Project Euler #36 asks for numbers
// that are palindromes in base 10 and base 2). It panics if base < 2.
func IsPalindromeInBase[E constraints.Integer](n E, base int) bool {
	checkBase(base, int(^uint(0)>>1))
	u, b := mathx.Abs64(n), uint64(base)
	// Reverse in uint64, stopping early once the reversed half reaches the
	// remaining half, so nothing can overflow.
	if u != 0 && u%b == 0 {
		return false
	}
	var r uint64
	for u > r {
		r = r*b + u%b
		u /= b
	}
	return u == r || u == r/b
}

// IsPalindromeString reports whether s reads the same forwards and
// backwards, rune by rune.
func IsPalindromeString(s string) bool {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		if r[i] != r[j] {
			return false
		}
	}
	return true
}

// IsPandigital reports whether the decimal digits of |n| are 1 through 9,
// each exactly once (Project Euler's "1 to 9 pandigital").
func IsPandigital[E constraints.Integer](n E) bool {
	return IsPandigitalN(n, 9)
}

// IsPandigitalN reports whether the decimal digits of |n| are 1 through k,
// each exactly once (Project Euler #41's "n-digit pandigital"). It panics
// unless 1 <= k <= 9.
func IsPandigitalN[E constraints.Integer](n E, k int) bool {
	if k < 1 || k > 9 {
		panic("digits: IsPandigitalN needs 1 <= k <= 9")
	}
	u := mathx.Abs64(n)
	var seen uint16
	count := 0
	for u > 0 {
		d := u % 10
		if d == 0 || d > uint64(k) || seen&(1<<d) != 0 {
			return false
		}
		seen |= 1 << d
		count++
		u /= 10
	}
	return count == k
}

// IsPandigitalInBase reports whether every digit 0 through base-1 occurs at
// least once among the digits of |n| in that base (the usual mathematical
// definition: 1023456789 is pandigital in base 10, and so is 10234567899).
// It panics unless 2 <= base <= 64.
func IsPandigitalInBase[E constraints.Integer](n E, base int) bool {
	checkBase(base, 64)
	u, b := mathx.Abs64(n), uint64(base)
	var seen uint64
	for {
		seen |= 1 << (u % b)
		u /= b
		if u == 0 {
			break
		}
	}
	return seen == allDigits(base)
}

// IsPandigitalString reports whether s, read as digits in the given base
// (using ToBase's alphabet), contains every digit 0 through base-1 at least
// once and nothing else. It returns false if s contains a character that is
// not a digit of that base. It panics unless 2 <= base <= 62.
func IsPandigitalString(s string, base int) bool {
	checkBase(base, len(alphabet))
	var seen uint64
	for _, ch := range s {
		d := strings.IndexRune(alphabet[:base], ch)
		if d < 0 {
			return false
		}
		seen |= 1 << d
	}
	return seen == allDigits(base)
}

func allDigits(base int) uint64 {
	if base == 64 {
		return ^uint64(0)
	}
	return 1<<base - 1
}

// Concat returns the number formed by writing the decimal digits of each of
// ns one after the other: Concat(1, 23, 4) is 1234 (Project Euler #38's
// "concatenated product"). A 0 contributes a single 0 digit. It panics if
// any of ns is negative.
func Concat[E constraints.Integer](ns ...E) E {
	var r E
	for _, n := range ns {
		if n < 0 {
			panic("digits: Concat of a negative number")
		}
		for range Count(n) {
			r *= 10
		}
		r += n
	}
	return r
}

// ConcatString is Concat, returning the digits as a string instead of a
// number, so it can't overflow: ConcatString(1, 23, 4) is "1234". It panics
// if any of ns is negative.
func ConcatString[E constraints.Integer](ns ...E) string {
	var b strings.Builder
	for _, n := range ns {
		if n < 0 {
			panic("digits: ConcatString of a negative number")
		}
		b.WriteString(ToBase(n, 10))
	}
	return b.String()
}

// ToBase returns n written in the given base, using the digits 0-9, then
// A-Z, then a-z, with a leading "-" for negative n. ToBase(0, b) is "0".
// For bases up to 36 it matches strconv.FormatInt except that letters are
// upper case. It panics unless 2 <= base <= 62.
func ToBase[E constraints.Integer](n E, base int) string {
	checkBase(base, len(alphabet))
	u, b := mathx.Abs64(n), uint64(base)
	var buf [65]byte
	i := len(buf)
	for {
		i--
		buf[i] = alphabet[u%b]
		u /= b
		if u == 0 {
			break
		}
	}
	if n < 0 {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

// SumOf returns a collector that adds up the digit sums of a stream's
// elements:
//
//	stream.RangeClosed(1, 10).Collect(digits.SumOf[int]()) // 46
func SumOf[E constraints.Integer]() stream.Collector[E, E] {
	return stream.NewCollector(
		func() E { return 0 },
		func(acc, v E) E { return acc + Sum(v) },
		stream.Identity[E],
	)
}
