// Package roman converts between integers and Roman numerals (Project
// Euler #89).
package roman

import (
	"fmt"
	"strings"
)

// numerals are the symbols and subtractive pairs of the minimal form, from
// the largest value down.
var numerals = []struct {
	value  int
	symbol string
}{
	{1000, "M"}, {900, "CM"}, {500, "D"}, {400, "CD"},
	{100, "C"}, {90, "XC"}, {50, "L"}, {40, "XL"},
	{10, "X"}, {9, "IX"}, {5, "V"}, {4, "IV"}, {1, "I"},
}

// Format returns n as a Roman numeral in its minimal form: the fewest
// symbols, using the subtractive pairs IV, IX, XL, XC, CD and CM.
// Format(1994) is "MCMXCIV". There is no symbol above M, so every
// thousand is one M: Format(5000) is "MMMMM". It panics if n < 1 (the
// Romans had no zero).
func Format(n int) string {
	if n < 1 {
		panic("roman: Format needs n >= 1")
	}
	var b strings.Builder
	for _, num := range numerals {
		for n >= num.value {
			b.WriteString(num.symbol)
			n -= num.value
		}
	}
	return b.String()
}

func value(c byte) int {
	switch c {
	case 'I':
		return 1
	case 'V':
		return 5
	case 'X':
		return 10
	case 'L':
		return 50
	case 'C':
		return 100
	case 'D':
		return 500
	case 'M':
		return 1000
	}
	return 0
}

// Parse returns the value of the Roman numeral s. It reads the minimal
// form and also the longer, valid-but-wasteful spellings of Project Euler
// #89, such as "IIII" for 4 or "VIIII" for 9: symbols are added up from
// left to right, except that a symbol standing directly before a larger
// one is subtracted from it. Upper case only.
//
// It returns an error if s is empty, contains another character, or is
// out of order: only the six standard subtractive pairs (IV, IX, XL, XC,
// CD, CM) are accepted, the values must not increase otherwise, and what
// follows a pair must be smaller than the symbol it subtracted (so "IXI"
// and "XLX" are errors). It does not limit how often a symbol repeats.
func Parse(s string) (int, error) {
	if s == "" {
		return 0, fmt.Errorf("roman: empty numeral")
	}
	total := 0
	limit := 1000 // the largest value the next symbol or pair may have
	for i := 0; i < len(s); i++ {
		v := value(s[i])
		if v == 0 {
			return 0, fmt.Errorf("roman: invalid character %q in %q", s[i], s)
		}
		next := v // the limit for whatever follows
		if i+1 < len(s) {
			if big := value(s[i+1]); big > v {
				// Only I, X and C subtract, and only from the next two
				// symbols up.
				if (v != 1 && v != 10 && v != 100) || (big != 5*v && big != 10*v) {
					return 0, fmt.Errorf("roman: invalid subtractive pair %q in %q", s[i:i+2], s)
				}
				next = v - 1
				v = big - v
				i++
			}
		}
		if v > limit {
			return 0, fmt.Errorf("roman: symbols out of order at position %d in %q", i, s)
		}
		total += v
		limit = next
	}
	return total, nil
}
