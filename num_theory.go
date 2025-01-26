package eulerlib

import (
	"math"
	"math/big"
	"math/bits"
	"strconv"
)

// returns the number of divisors the given integer has
func CountDivisors[E Integer](n E) E {
	count := E(0)
	end := E(math.Sqrt(float64(n)))
	for i := E(1); i < end; i++ {
		if n%i == 0 {
			count += 2
		}
	}
	if end*end == n {
		count++
	}
	return count
}

// returns all divisors of the given integer
func Divisors[E Integer](n E) []E {
	end := E(math.Sqrt(float64(n)))
	divisors := []E{1}

	for i := E(2); i <= end; i++ {
		if n%i == 0 {
			divisors = append(divisors, E(i))
			divisors = append(divisors, E(n/i))
		}
	}
	divisors = append(divisors, E(n))
	return divisors
}

// returns a slice with all permutations of the given slice
func Permutations[E AnyComparable](arr []E) [][]E {
	var helper func([]E, int)
	res := [][]E{}

	helper = func(arr []E, n int) {
		if n == 1 {
			tmp := make([]E, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
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

// Returns all combinations of the elements in the given slice
func Combinations[E AnyComparable, F Integer](set []E, n F) (subsets [][]E) { // https://github.com/mxschmitt/golang-combinations/blob/master/combinations.go
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

		for object := uint(0); object < length; object++ {
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
func Factorial[E Integer](n E) E {
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
func FactorialDigitSum[E Integer](n E) E {
	s := strconv.Itoa(int(n))
	res := E(0)

	for j := 0; j < len(s); j++ {
		i, _ := strconv.Atoi(string(s[j]))
		res += Factorial(E(i))
	}
	return res
}

func DigitSum[E Integer](n E) (res E) {

	for n > 0 {
		res += n % 10
		n /= 10
	}

	return
}

func DigitSumString(s string) int64 {
	tmp, _ := strconv.Atoi(s)
	return DigitSum(int64(tmp))
}

// Calculates greatest common divisor for the given integers
func Gcd[E Integer](args ...E) E {
outer:
	for i := args[0]; ; i-- {
		for _, v := range args {
			if v%i != 0 {
				continue outer
			}
		}
		return i
	}
}

func Factorize[E Integer](n E) map[E]E {
	factors := make(map[E]E)
	for i := E(2); i <= n; i++ {
		for n%i == 0 {
			factors[i]++
			n /= i
		}
	}
	return factors
}

func FactorizeBigInt(n *big.Int) map[int64]int64 {
	factors := make(map[int64]int64)
	i := big.NewInt(2)
	for i.Cmp(n) == -1 || i.Cmp(n) == 0 {
		for n.Mod(n, i).Cmp(big.NewInt(0)) == 0 {
			factors[i.Int64()]++
			n.Div(n, i)
		}
		i.Add(i, big.NewInt(1))
	}
	return factors
}

// Calculates all primefactors of the given number
func PrimeFactors[E Integer](n E) []E {
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

// Return x^y % m
func PowMod[E Integer](x, y, p E) int64 {
	res := int64(1)
	for y > 0 {
		if y&1 == 1 {
			res = (res * int64(x)) % int64(p)
		}
		y >>= 1
		x = (x * x) % p
	}
	return res
}

// Returns the given angle (in degrees) in radians.
func ToRadians(n float64) float64 {
	return n * math.Pi / 180
}

// Returns n!/(k!(n-k)!)
func Binomial[E Integer](n E, k E) *big.Int {
	prod := big.NewInt(1)
	for i := E(0); i < k; i++ {
		prod.Mul(prod, big.NewInt(int64(n-i)))
	}
	return prod.Div(prod, FactorialBigInt(int64(k)))
}

// Calculates all primefactors of the given Big Integer
func PrimeFactorsBigInt(n *big.Int) (primefs [][]int64) {
	var i2 *big.Int
	tmp, _ := new(big.Int).SetString(n.Text(10), 10)
	two := big.NewInt(2)
	tmp2 := []int64{0, 0}
	for tmp.Mod(n, two).Text(10) == "0" {
		tmp2[0] = 2
		tmp2[1]++
		n.Div(n, two)
		tmp, _ = new(big.Int).SetString(n.Text(10), 10)
	}

	primefs = append(primefs, tmp2)

	i := big.NewInt(3)
	itmp, _ := new(big.Int).SetString(i.Text(10), 10)
	for itmp.Mul(itmp, itmp).Cmp(n) == 0 || itmp.Mul(itmp, itmp).Cmp(n) == -1 {
		tmp2 = []int64{0, 0}
		i2, _ = new(big.Int).SetString(i.Text(10), 10)
		for tmp.Mod(n, i).Text(10) == "0" {
			tmp2[0] = i.Int64()
			tmp2[1]++
			n.Div(n, i2)
			tmp, _ = new(big.Int).SetString(n.Text(10), 10)
		}
		primefs = append(primefs, tmp2)
		i.Add(i, two)
		itmp = i
	}

	if n.Cmp(two) == 1 {
		primefs = append(primefs, []int64{n.Int64(), 1})
	}

	return primefs
}

// Calculates b^n using integers
func Pow[E Integer](b E, n E) E {
	if n == 1 {
		return b
	}
	res := E(1)
	for i := E(1); i <= n; i++ {
		res *= b
	}
	return res
}

// Calculates b^n as a Big Integer
func PowBigInt(b *big.Int, n int64) *big.Int {
	if n == 1 {
		return b
	}
	tmp := big.NewInt(b.Int64())
	for i := int64(1); i <= n; i++ {
		b.Mul(b, tmp)
	}
	return b
}

// Calculates b^n as a Big Float
func PowBigFloat(b *big.Float, n int64) *big.Float {
	tmp := new(big.Float).Copy(b)
	for i := int64(1); i < n; i++ {
		b.Mul(b, tmp)
	}
	return b
}

// Checks whether the given number is a power of 2
func IsPowerOfTwo[E Integer](n E) bool {
	return n&(n-1) == 0
}

// Checks whether or not the given number is a perfect square
func IsSquare[E Integer](n E) bool {
	s := E(math.Sqrt(float64(n)))
	return s*s == n
}

func IsInteger[E Float](n E) bool {
	return E(math.Floor(float64(n))) == n
}

func Lcd[E Integer](nums ...E) E {
	res := E(1)
	for _, v := range nums {
		res = res * v / Gcd(res, v)
	}
	return res
}

func Lcm[E Integer](nums ...E) E {
	res := E(1)
	for _, v := range nums {
		res = res * v / Gcd(res, v)
	}
	return res
}

// this function reduces a slice of integers using the given function
func Reduce[E Integer](f func(E, E) E, nums []E) E {
	res := nums[0]
	for _, v := range nums[1:] {
		res = f(res, v)
	}
	return res
}

// MaxInSlice: this functions returns the max value in the given slice (redirects to Max function)
// Deprecated: Use Max instead
func MaxInSlice[E Integer](nums []E) E {
	return Max(nums...)
}
