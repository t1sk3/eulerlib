package eulerlib

import (
	"reflect"
	"slices"
	"sort"
	"strconv"
)

// Checks whether the given integer is pandigital in base 10
func IsPandigital[E Integer](n E) bool {
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
func IsPandigitalInBaseString[E Integer](s string, b E) bool {
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
func IsPandigitalInBase[E Integer](n E, b E) bool {
	if b > 36 {
		return false
	}
	return IsPandigitalInBaseString(DecimalToBase(n, b), b)
}

// creates a slice containing all digits of the given integer as individual integers
func MakeIntSlice[E Integer](n E) []E {
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
func Min[E Number](numbers ...E) E {
	if len(numbers) == 0 {
		return E(0)
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
func Max[E Number](numbers ...E) E {
	if len(numbers) == 0 {
		return E(0)
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
func JoinSlice[E Integer](s []E) string {
	res := ""
	for _, e := range s {
		res += strconv.Itoa(int(e))
	}
	return res
}

// joins a slice of strings into a single string
func JoinSliceString(s []string) string {
	res := ""
	for _, e := range s {
		res += e
	}
	return res
}

// removes any duplicates fm a slice
func RemoveDuplicates[E AnyComparable](s []E) []E {
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

func RemoveDuplicates2[E any](s []E) []E {
	list := []E{}
	for _, item := range s {
		if !SliceContainsAny(list, item) {
			list = append(list, item)
		}
	}

	return list
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

// returns the sum of theintegers in the given slice
func Sum[E Number](lst []E) E {
	res := E(0)
	for _, element := range lst {
		res += element
	}
	return res
}

// Converts the given integer to the give base and returns it as a string, max base is 36
func DecimalToBase[E Integer, F Integer](n E, b F) (res string) {
	alphabet := "0123456789abcdefghijklmnopqrstuvwxyz"
	n64 := int64(n)
	b64 := int64(b)

	if b64 < 2 || b64 > 36 {
		return
	}

	for n64 != 0 {
		res = string(alphabet[n64%b64]) + res
		n64 /= b64
	}

	return
}

// Returns a substring of the given string from the beginning up to a given index
func Substring(s string, n int) string {
	return s[:n]
}

// Reverses a string
func ReverseString(s string) (res string) {
	for i := len(s) - 1; i >= 0; i-- {
		res += string(s[i])
	}
	return
}

// Counts the occurence of a given in a slice
func CountOccurrenceInSlice[E Number](s []E, p E) (res E) {
	for _, e := range s {
		if e == p {
			res++
		}
	}
	return
}

// Removes element from slice at index
func RemoveFromSlice[E Number](slice []E, s int) []E {
	return append(slice[:s], slice[s+1:]...)
}

// Returns the number of digits in the given integer
func Totient[E Integer](n E) E {
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

// Lists totients
func ListTotients[E Integer](n E) []E {
	res := make([]E, n+1)
	for i := E(1); i <= n; i++ {
		res[i] = Totient(i)
	}
	return res
}

// Returns the number of digits in the given integer
func Range[E Number](start, stop E) (res []E) {
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
func RangeStep[E Number](start, stop, step E) (res []E) {
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
func SliceContains[E AnyComparable](s []E, e E) bool {
	return slices.Contains(s, e)
}

func SliceContainsAny[E any](s []E, e E) bool {
	for _, v := range s {
		if reflect.DeepEqual(v, e) {
			return true
		}
	}
	return false
}

// generates a slice of length n with all elements set to the given value
func GenerateSlice[E Integer, F AnyComparable](n E, value F) (res []F) {
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
