package num_theory

import (
	"math"
	"math/big"
	"math/bits"
	"strconv"

	"github.com/t1sk3/eulerlib/etc"
	"github.com/t1sk3/eulerlib/utils"
)

// CountDivisors returns the number of divisors the given integer has.
func CountDivisors[E utils.Integer](n E) E {
	count := E(0)
	end := E(math.Sqrt(float64(n)))
	for i := E(1); i <= end; i++ {
		if n%i == 0 {
			if i == n/i {
				count++
			} else {
				count += 2
			}
		}
	}
	return count
}

// returns all divisors of the given integer
func Divisors[E utils.Integer](n E) []E {
	end := E(math.Sqrt(float64(n)))
	divisors := []E{1}

	for i := E(2); i <= end; i++ {
		if n%i == 0 {
			divisors = append(divisors, E(i))
			divisors = append(divisors, E(n/i))
		}
	}
	divisors = append(divisors, E(n))
	return etc.RemoveDuplicates(divisors)
}

// returns a slice with all permutations of the given slice
func Permutations[E utils.Comparable](arr []E) [][]E {
	var helper func([]E, int)
	res := [][]E{}

	helper = func(arr []E, n int) {
		if n == 1 {
			tmp := make([]E, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := range n {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
}

// returns the amount of permutations of the given slice
func PermutationCount[E comparable](n []E) int {
	elements := etc.UniqueCount(n)
	res := Factorial(len(n))
	for _, v := range elements {
		res /= Factorial(v)
	}
	return res
}

// Returns all combinations of the elements in the given slice
func Combinations[E any, F utils.Integer](set []E, n F) (subsets [][]E) { // https://github.com/mxschmitt/golang-combinations/blob/master/combinations.go
	length := uint(len(set))

	if n > F(len(set)) {
		n = F(len(set))
	}

	// Go through all possible combinations of objects
	// from 1 (only first object in subset) to 2^length (all objects in subset)
	for subsetBits := 1; subsetBits < (1 << length); subsetBits++ {
		if n > 0 && F(bits.OnesCount(uint(subsetBits))) != n {
			continue
		}

		var subset []E

		for object := range length {
			// checks if object is contained in subset
			// by checking if bit 'object' is set in subsetBits
			if (subsetBits>>object)&1 == 1 {
				// add object to subset
				subset = append(subset, set[object])
			}
		}
		// add subset to subsets
		subsets = append(subsets, subset)
	}
	return subsets
}

// Calculates the factorial of the given integer
func Factorial[E utils.Integer](n E) E {
	if n == 0 {
		return 1
	}
	res := E(1)
	for i := E(2); i < n+1; i++ {
		res *= i
	}
	return res
}

// Return the factorial as a Big Integer
func FactorialBigInt(n int64) *big.Int {
	res := big.NewInt(1)
	for i := n; i > 1; i-- {
		res.Mul(res, big.NewInt(i))
	}
	return res
}

// Calculates factorial digital sum
func FactorialDigitSum[E utils.Integer](n E) E {
	s := strconv.Itoa(int(n))
	res := E(0)

	for j := 0; j < len(s); j++ {
		i, _ := strconv.Atoi(string(s[j]))
		res += Factorial(E(i))
	}
	return res
}

func DigitSum[E utils.Integer](n E) (res E) {

	for n > 0 {
		res += n % 10
		n /= 10
	}

	return
}

// DigitSumString returns the sum of the decimal digits of the string s.
// Works correctly for arbitrarily large digit strings.
func DigitSumString(s string) int64 {
	var res int64
	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			res += int64(ch - '0')
		}
	}
	return res
}

// Gcd calculates the greatest common divisor for the given integers. Panics if no args given.
func Gcd[E utils.Integer](args ...E) E {
	if len(args) == 0 {
		panic("Gcd: at least one argument required")
	}
	res := args[0]
	for _, v := range args[1:] {
		res = gcd(res, v)
	}
	return res
}

func gcd[E utils.Integer](a, b E) E {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// Factorize returns the prime factorization of n as a map of prime → exponent.
// Returns an empty map for n <= 1.
func Factorize[E utils.Integer](n E) map[E]E {
	factors := make(map[E]E)
	for i := E(2); i <= n; i++ {
		for n%i == 0 {
			factors[i]++
			n /= i
		}
	}
	return factors
}

// FactorizeBigInt returns the prime factorization of n as a slice of [factor, exponent] pairs.
func FactorizeBigInt(n *big.Int) [][2]*big.Int {
	var factors [][2]*big.Int
	n = new(big.Int).Set(n)
	two := big.NewInt(2)
	i := new(big.Int).Set(two)
	for {
		sq := new(big.Int).Mul(i, i)
		if sq.Cmp(n) > 0 {
			break
		}
		if new(big.Int).Mod(n, i).Sign() == 0 {
			exp := big.NewInt(0)
			for new(big.Int).Mod(n, i).Sign() == 0 {
				n.Div(n, i)
				exp.Add(exp, big.NewInt(1))
			}
			fi := new(big.Int).Set(i)
			factors = append(factors, [2]*big.Int{fi, exp})
		}
		i.Add(i, big.NewInt(1))
	}
	if n.Cmp(big.NewInt(1)) > 0 {
		factors = append(factors, [2]*big.Int{new(big.Int).Set(n), big.NewInt(1)})
	}
	return factors
}

// Calculates all primefactors of the given number
func PrimeFactors[E utils.Integer](n E) []E {
	primefs := []E{}
	for n%2 == 0 {
		primefs = append(primefs, 2)
		n = n / 2
	}

	for i := E(3); i*i <= n; i = i + 2 {
		for n%i == 0 {
			primefs = append(primefs, i)
			n = n / i
		}
	}

	if n > 2 {
		primefs = append(primefs, n)
	}

	return primefs
}

// PowMod returns x^y % p. Requires p < 2^31 to avoid int64 overflow during squaring.
func PowMod[E utils.Integer](x, y, p E) int64 {
	res := int64(1)
	xm := int64(x) % int64(p)
	for y > 0 {
		if y&1 == 1 {
			res = (res * xm) % int64(p)
		}
		y >>= 1
		xm = (xm * xm) % int64(p)
	}
	return res
}

// Returns the given angle (in degrees) in radians.
func ToRadians(n float64) float64 {
	return n * math.Pi / 180
}

// Returns n!/(k!(n-k)!)
func Binomial[E utils.Integer](n E, k E) *big.Int {
	prod := big.NewInt(1)
	for i := E(0); i < k; i++ {
		prod.Mul(prod, big.NewInt(int64(n-i)))
	}
	return prod.Div(prod, FactorialBigInt(int64(k)))
}

// PrimeFactorsBigInt returns the prime factorization of n as [][]int64,
// where each element is [prime, exponent]. Does not modify n.
func PrimeFactorsBigInt(n *big.Int) (primefs [][]int64) {
	n = new(big.Int).Set(n)
	two := big.NewInt(2)
	exp := int64(0)
	for new(big.Int).Mod(n, two).Sign() == 0 {
		exp++
		n.Div(n, two)
	}
	if exp > 0 {
		primefs = append(primefs, []int64{2, exp})
	}
	i := big.NewInt(3)
	for {
		sq := new(big.Int).Mul(i, i)
		if sq.Cmp(n) > 0 {
			break
		}
		exp = 0
		for new(big.Int).Mod(n, i).Sign() == 0 {
			exp++
			n.Div(n, i)
		}
		if exp > 0 {
			primefs = append(primefs, []int64{i.Int64(), exp})
		}
		i.Add(i, two)
	}
	if n.Cmp(big.NewInt(1)) > 0 {
		primefs = append(primefs, []int64{n.Int64(), 1})
	}
	return primefs
}

// Pow calculates b^n using integers. Panics for negative n.
func Pow[E utils.Integer](b E, n E) E {
	if n < 0 {
		panic("Pow: negative exponent")
	}
	res := E(1)
	for i := E(1); i <= n; i++ {
		res *= b
	}
	return res
}

// PowBigInt returns b^n as a new *big.Int using binary exponentiation. Does not modify b.
func PowBigInt(b *big.Int, n int64) *big.Int {
	if n == 0 {
		return big.NewInt(1)
	}
	result := big.NewInt(1)
	base := new(big.Int).Set(b)
	for n > 0 {
		if n&1 == 1 {
			result.Mul(result, base)
		}
		base.Mul(base, base)
		n >>= 1
	}
	return result
}

// PowBigFloat returns b^n as a new *big.Float using binary exponentiation. Does not modify b.
func PowBigFloat(b *big.Float, n int64) *big.Float {
	if n == 0 {
		return new(big.Float).SetInt64(1)
	}
	result := new(big.Float).SetInt64(1)
	base := new(big.Float).Copy(b)
	for n > 0 {
		if n&1 == 1 {
			result.Mul(result, base)
		}
		base.Mul(base, base)
		n >>= 1
	}
	return result
}

// Checks whether the given number is a power of 2
func IsPowerOfTwo[E utils.Integer](n E) bool {
	return n > 0 && n&(n-1) == 0
}

// Checks whether or not the given number is a perfect square
func IsSquare[E utils.Integer](n E) bool {
	if n < 0 {
		return false
	}

	// Quick elimination using hexadecimal bitmask
	// Perfect squares only end in 0, 1, 4, or 9 in base 16
	if (0x0202021202030213>>(n&0x3F))&1 == 0 {
		return false
	}

	// Fast binary search for the square root
	low, high := E(0), n
	for low <= high {
		mid := (low + high) / 2
		sq := mid * mid
		if sq == n {
			return true
		} else if sq < n {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return false
}

func FloatIsInteger[E utils.Float](n E) bool {
	return E(math.Floor(float64(n))) == n
}

func Lcm[E utils.Integer](nums ...E) E {
	res := E(1)
	for _, v := range nums {
		res = res * v / Gcd(res, v)
	}
	return res
}

// Lcd is deprecated. Use Lcm instead.
//
// Deprecated: Use Lcm.
func Lcd[E utils.Integer](nums ...E) E {
	return Lcm(nums...)
}

// Reduce applies f cumulatively to the elements of nums, starting with nums[0].
// Panics if nums is empty.
func Reduce[E any](nums []E, f func(E, E) E) E {
	if len(nums) == 0 {
		panic("Reduce: nums must not be empty")
	}
	res := nums[0]
	for _, v := range nums[1:] {
		res = f(res, v)
	}
	return res
}

// ReduceWithInit applies f cumulatively to the elements of nums, starting with initial.
func ReduceWithInit[E any](initial E, nums []E, f func(E, E) E) E {
	res := initial
	for _, v := range nums {
		res = f(res, v)
	}
	return res
}

// MaxInSlice returns the maximum value in the given slice.
//
// Deprecated: Use Max instead.
func MaxInSlice[E utils.Integer](nums []E) E {
	return etc.Max(nums...)
}

// SumDivisors returns the sum of proper divisors of n (all divisors excluding n itself).
func SumDivisors[E utils.Integer](n E) E {
	if n <= 1 {
		return 0
	}
	sum := E(1)
	end := E(math.Sqrt(float64(n)))
	for i := E(2); i <= end; i++ {
		if n%i == 0 {
			sum += i
			if i != n/i {
				sum += n / i
			}
		}
	}
	return sum
}

// IsAbundant returns true if the sum of proper divisors of n exceeds n.
func IsAbundant[E utils.Integer](n E) bool {
	return SumDivisors(n) > n
}

// IsDeficient returns true if the sum of proper divisors of n is less than n.
func IsDeficient[E utils.Integer](n E) bool {
	return SumDivisors(n) < n
}

// IsPerfect returns true if the sum of proper divisors of n equals n.
func IsPerfect[E utils.Integer](n E) bool {
	return SumDivisors(n) == n
}

// IsAmicable returns true if SumDivisors(SumDivisors(a)) == a and SumDivisors(a) != a.
func IsAmicable[E utils.Integer](a E) bool {
	b := SumDivisors(a)
	return b != a && SumDivisors(b) == a
}
