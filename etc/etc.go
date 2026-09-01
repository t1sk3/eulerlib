// Package etc is a grab-bag of generic slice, string, and file utilities
// used throughout eulerlib (and handy on their own): min/max, filter/map/
// reduce-style helpers, deduplication, ranges, base conversion, and small
// file I/O wrappers.
package etc

import (
	"os"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/t1sk3/eulerlib/v2/utils"
)

// Checks whether the given integer is pandigital in base 10
func IsPandigital[E utils.Integer](n E) bool {
	digits := make(map[E]bool)
	for n > 0 {
		digits[n%10] = true
		n /= 10
	}
	for i := E(1); i <= 9; i++ {
		if !digits[i] {
			return false
		}
	}
	return true
}

// Checks whether the given integer is pandigital in the given base
func IsPandigitalInBaseString[E utils.Integer](s string, b E) bool {
	digits := make(map[E]bool)
	for _, e := range s {
		d, _ := strconv.Atoi(string(e))
		if d >= int(b) {
			return false
		}
		digits[E(d)] = true
	}
	for i := E(0); i < b; i++ {
		if !digits[i] {
			return false
		}
	}
	return true
}

// Checks whether the given integer is pandigital in the given base
// It converts the integer to the given base and then checks whether it is pandigital
func IsPandigitalInBase[E utils.Integer](n E, b E) bool {
	if b > 36 {
		return false
	}
	return IsPandigitalInBaseString(DecimalToBase(n, b), b)
}

// creates a slice containing all digits of the given integer as individual integers
func MakeIntSlice[E utils.Integer](n E) []E {
	res := []E{}
	n_string := DecimalToBase(n, 10)
	var temp int
	for _, e := range n_string {
		temp, _ = strconv.Atoi(string(e))
		res = append(res, E(temp))
	}
	return res
}

// Checks whether the given stringis a palindrome
func IsPalindrome(s string) bool {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes) == s
}

// returns the minimum of all given numbers
func Min[E utils.RealNumber](numbers ...E) E {
	if len(numbers) == 0 {
		panic("at least one number is required in Min[E]")
	}
	m := numbers[0]
	for _, n := range numbers {
		if n < m {
			m = n
		}
	}
	return m
}

// returns the maximum of all given numbers
func Max[E utils.RealNumber](numbers ...E) E {
	if len(numbers) == 0 {
		panic("at least one number is required in Max[E]")
	}
	m := numbers[0]
	for _, n := range numbers {
		if n > m {
			m = n
		}
	}
	return m
}

// joins a slice of integers into a single string
func JoinSlice[E utils.Integer](s []E) string {
	res := ""
	for _, e := range s {
		res += strconv.Itoa(int(e))
	}
	return res
}

// joins a slice of strings into a single string
func JoinSliceString(s []string) string {
	return strings.Join(s, "")
}

// removes any duplicates fm a slice
func RemoveDuplicates[E utils.Comparable](s []E) []E {
	res := []E{}
	keys := make(map[E]bool)

	for _, entry := range s {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			res = append(res, entry)
		}
	}
	return res
}

// RemoveDuplicatesFunc removes duplicates from s using equality function f.
func RemoveDuplicatesFunc[E any](s []E, f func(E, E) bool) []E {
	list := []E{}
	for _, item := range s {
		if !SliceContainsAny(list, item, f) {
			list = append(list, item)
		}
	}
	return list
}

// RemoveDuplicates2 is deprecated. Use RemoveDuplicatesFunc instead.
//
// Deprecated: Use RemoveDuplicatesFunc.
func RemoveDuplicates2[E any](s []E, f func(E, E) bool) []E {
	return RemoveDuplicatesFunc(s, f)
}

// Removes duplicates from a slice of strings
func RemoveDuplicateSlices(s [][]string) (res [][]string) {
	var tmp bool
	for _, e := range s {
		tmp = true
		for _, e2 := range res {
			if reflect.DeepEqual(e, e2) {
				tmp = false
				break
			}
		}
		if tmp {
			res = append(res, e)
		}
	}
	return
}

// returns the sum of the integers in the given slice
func Sum[E utils.Number](lst []E) E {
	res := E(0)
	for _, element := range lst {
		res += element
	}
	return res
}

// Product returns the product of all elements in lst. Returns 1 for an empty slice.
func Product[E utils.Number](lst []E) E {
	res := E(1)
	for _, element := range lst {
		res *= element
	}
	return res
}

// DecimalToBase converts n to base b (2..62) and returns the result as a string.
// Returns "0" for n==0, prepends "-" for negative n, returns "" for invalid base.
func DecimalToBase[E utils.Integer, F utils.Integer](n E, b F) (res string) {
	alphabet := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	n64 := int64(n)
	b64 := int64(b)

	if b64 < 2 || b64 > 62 {
		return
	}
	if n64 == 0 {
		return "0"
	}
	negative := n64 < 0
	if negative {
		n64 = -n64
	}
	for n64 != 0 {
		res = string(alphabet[n64%b64]) + res
		n64 /= b64
	}
	if negative {
		res = "-" + res
	}
	return
}

// ReverseString reverses a string
func ReverseString(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// Counts the occurence of a given in a slice
func CountOccurrenceInSlice[E utils.RealNumber](s []E, p E) (res E) {
	for _, e := range s {
		if e == p {
			res++
		}
	}
	return
}

// Removes element from slice at index
func RemoveFromSlice[E utils.RealNumber](slice []E, s int) []E {
	return append(slice[:s], slice[s+1:]...)
}

// Totient returns Euler's totient φ(n): the count of integers in [1,n] coprime to n.
func Totient[E utils.Integer](n E) E {
	res := n
	for i := E(2); i*i <= n; i++ {
		if n%i == 0 {
			for n%i == 0 {
				n /= i
			}
			res -= res / i
		}
	}
	if n > 1 {
		res -= res / n
	}
	return res
}

// ListTotients returns a slice of length n+1 where res[i] = φ(i) for i in [0,n].
// Uses a sieve for O(n log log n) performance.
func ListTotients[E utils.Integer](n E) []E {
	res := make([]E, n+1)
	for i := E(0); i <= n; i++ {
		res[i] = i
	}
	for i := E(2); i <= n; i++ {
		if res[i] == i { // i is prime
			for j := i; j <= n; j += i {
				res[j] -= res[j] / i
			}
		}
	}
	return res
}

// ListMobius returns a slice of length n+1 where res[i] = μ(i), the Möbius
// function, for i in [0,n]. μ(0) is defined as 0; for n >= 1, μ(1) = 1,
// μ(i) = 0 if i has a squared prime factor, otherwise μ(i) = (-1)^k where k
// is the number of distinct prime factors of i. Uses a linear sieve for
// O(n) performance. Requires a signed integer type since μ can be -1.
func ListMobius[E utils.SignedInteger](n E) []E {
	if n < 0 {
		panic("n must be positive")
	}
	mu := make([]E, n+1)
	isComposite := make([]bool, n+1)
	var primes []E
	if n >= 1 {
		mu[1] = 1
	}
	for i := E(2); i <= n; i++ {
		if !isComposite[i] {
			primes = append(primes, i)
			mu[i] = -1
		}
		for _, p := range primes {
			if i*p > n {
				break
			}
			isComposite[i*p] = true
			if i%p == 0 {
				mu[i*p] = 0
				break
			}
			mu[i*p] = -mu[i]
		}
	}
	return mu
}

// Range returns a slice from start (inclusive) to stop (exclusive), stepping by 1.
// If start > stop, it counts down.
func Range[E utils.RealNumber](start, stop E) (res []E) {
	if start < stop {
		for i := start; i < stop; i++ {
			res = append(res, i)
		}
	} else if start > stop {
		for i := start; i > stop; i-- {
			res = append(res, i)
		}

	}
	return
}

// Returns a slice of integers from start to stop with the given step
func RangeStep[E utils.RealNumber](start, stop, step E) (res []E) {
	if step == 0 {
		return
	}
	if start < stop {
		for i := start; i < stop; i += step {
			res = append(res, i)
		}
	} else if start > stop {
		for i := start; i > stop; i -= step {
			res = append(res, i)
		}
	}
	return
}

// Checks whether the given slice contains the given element
func SliceContains[E utils.Comparable](s []E, e E) bool {
	return slices.Contains(s, e)
}

// SliceContainsAny returns true if any element of s equals e according to f.
func SliceContainsAny[E any](s []E, e E, f func(E, E) bool) bool {
	for _, v := range s {
		if f(v, e) {
			return true
		}
	}
	return false
}

// generates a slice of length n with all elements set to the given value
func GenerateSlice[E utils.Integer, F utils.Comparable](n E, value F) (res []F) {
	for i := E(0); i < n; i++ {
		res = append(res, value)
	}
	return
}

// filters a slice based on the given function
func Filter[E any](s []E, f func(E) bool) (res []E) {
	for _, e := range s {
		if f(e) {
			res = append(res, e)
		}
	}
	return
}

// maps a slice based on the given function
func Map[E any, F any](s []E, f func(E) F) (res []F) {
	for _, e := range s {
		res = append(res, f(e))
	}
	return
}

// sorts a slice based on the given function
func Sort[E any](s []E, f func(E, E) bool) []E {
	sort.Slice(s, func(i, j int) bool {
		return f(s[i], s[j])
	})
	return s
}

// returns a slice of unique elements from the given slice
func Unique[E comparable](s []E) []E {
	keys := make(map[E]bool)
	list := []E{}
	for _, entry := range s {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// returns a slice of unique elements and their counts from the given slice
func UniqueCount[E comparable](s []E) map[E]int {
	keys := make(map[E]int)
	for _, entry := range s {
		keys[entry]++
	}
	return keys
}

// creates a file with the given name
func CreateFile(name string) (*os.File, error) {
	return os.Create(name)
}

// creates a file with the given name and writes the given content to it
func CreateFileWithContent(name string, content string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	return err
}

// reads the content of the file with the given name and returns it as a string
func ReadFile(name string) (string, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// checks whether a file with the given name exists
func FileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

// Abs returns the absolute value of n.
func Abs[E utils.SignedInteger](n E) E {
	if n < 0 {
		return -n
	}
	return n
}

// SumOfSquares returns 1² + 2² + ... + n².
func SumOfSquares[E utils.Integer](n E) E {
	return n * (n + 1) * (2*n + 1) / 6
}

// SquareOfSum returns (1 + 2 + ... + n)².
func SquareOfSum[E utils.Integer](n E) E {
	s := n * (n + 1) / 2
	return s * s
}
