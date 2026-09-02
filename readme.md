# eulerlib

A generic Go library of number-theory utilities for solving [Project Euler](https://projecteuler.net/) problems. All functions are generic over Go's built-in numeric types via type constraints.

```
go get github.com/t1sk3/eulerlib/v2
```

Requires **Go 1.23+**.

eulerlib is split into small, focused packages instead of one flat package — import only what you need:

```go
import (
	"github.com/t1sk3/eulerlib/v2/prime_numbers"
	"github.com/t1sk3/eulerlib/v2/num_theory"
	"github.com/t1sk3/eulerlib/v2/vector"
)

prime_numbers.IsPrime(97)        // true
num_theory.Factorize(int64(360)) // map[2:3 3:2 5:1]
vector.NewVector([]int{1, 2, 3}) // *Vector[int]
```

| Package | Import path | Covers |
|---|---|---|
| `prime_numbers` | `.../v2/prime_numbers` | Primality, prime listing/counting/generation, prime iterator |
| `num_theory` | `.../v2/num_theory` | Divisors, GCD/LCM, factorization, powers, modular arithmetic, combinatorics |
| `fibonacci` | `.../v2/fibonacci` | Fibonacci sequences and single terms (incl. `big.Int`), Collatz-adjacent helpers |
| `figurate` | `.../v2/figurate` | Triangular, pentagonal, and hexagonal numbers |
| `sequences` | `.../v2/sequences` | Collatz sequence and length |
| `pythagoras` | `.../v2/pythagoras` | Pythagorean triplet check |
| `vector` | `.../v2/vector` | N-dimensional, 2D, and 3D vector arithmetic; cartesian/hexagonal coordinate operations |
| `etc` | `.../v2/etc` | Generic slice/string/file utilities (min/max, filter/map/reduce, dedup, ranges, base conversion, file I/O) |
| `utils` | `.../v2/utils` | Generic type constraints (`Integer`, `Float`, `Number`, ...) and reflection-based type checks |
| `progress` | `.../v2/progress` | Live progress bar (and open-ended spinner) for `range` loops over a slice or a count |
| `stopwatch` | `.../v2/stopwatch` | Timing helper for benchmarking brute-force runs |
| `memo` | `.../v2/memo` | Generic memoization wrappers for recursive/DP-style brute forces |

Each package also has full [godoc](https://pkg.go.dev/github.com/t1sk3/eulerlib/v2) comments on every exported symbol.

---

## Reference

### Primes — `prime_numbers`

| Function | Description |
|---|---|
| `IsPrime(n)` | Returns true if n is prime |
| `NextPrime(n)` | Returns the next prime after n |
| `ListPrimes(n)` | Returns all primes up to n |
| `ListPrimality(n)` | Sieve of Eratosthenes — returns `[]bool` indexed by number |
| `ListSmallestPrimeFactors(n)` | Sieve — returns `[]E` of each number's smallest prime factor, for O(log n) factorization via `num_theory.FactorizeSPF` |
| `FirstNPrimes(n)` | Returns the first n primes |
| `PrimeCount(s, e)` | Counts primes in the range [s, e] |
| `SumPrimes(s, e)` | Sums all primes in the range [s, e] |
| `PrimeGenerator(limit)` | Channel-based prime generator up to limit |
| `NewPrimeNumberIterator(start?)` | Iterator that steps through primes one at a time |

### Divisors, Factorization & GCD — `num_theory`

| Function | Description |
|---|---|
| `Divisors(n)` | Returns all divisors of n |
| `CountDivisors(n)` | Returns the number of divisors of n |
| `SumDivisors(n)` | Sum of proper divisors of n (excludes n itself) |
| `IsAbundant(n)` | True if `SumDivisors(n) > n` |
| `IsDeficient(n)` | True if `SumDivisors(n) < n` |
| `IsPerfect(n)` | True if `SumDivisors(n) == n` (e.g. 6, 28, 496) |
| `IsAmicable(n)` | True if n and `SumDivisors(n)` form an amicable pair |
| `Factorize(n)` | Prime factorization as `map[E]E{prime: exponent}` |
| `FactorizeSPF(n, spf)` | Prime factorization using a smallest-prime-factor table from `prime_numbers.ListSmallestPrimeFactors` — O(log n) instead of `Factorize`'s O(n) |
| `FactorizeBigInt(n)` | Prime factorization of `*big.Int` as `[][2]*big.Int` |
| `PrimeFactors(n)` | Returns the prime factors of n (with repeats) |
| `PrimeFactorsBigInt(n)` | Prime factors of a `*big.Int` as `[][]int64{prime, exp}` |
| `Gcd(a, b, ...)` | Greatest common divisor |
| `Lcm(a, b, ...)` | Least common multiple |
| `IsCoprime(a, b)` | True if `Gcd(a, b) == 1` |
| `ExtGcd(a, b)` | Extended Euclidean algorithm — returns `(g, x, y)` with `a*x + b*y = g` |
| `ModInverse(a, m)` | Modular multiplicative inverse of a mod m, and whether it exists |
| `CRT(remainders, moduli)` | Chinese Remainder Theorem — solves `x ≡ remainders[i] (mod moduli[i])` for all i (moduli needn't be coprime) |

### Powers & Modular Arithmetic — `num_theory`

| Function | Description |
|---|---|
| `Pow(b, n)` | b^n as an integer |
| `PowBigInt(b, n)` | b^n as `*big.Int` (non-mutating, binary exponentiation) |
| `PowBigFloat(b, n)` | b^n as `*big.Float` (non-mutating, binary exponentiation) |
| `PowMod(x, y, p)` | x^y mod p (requires p < 2³¹) |
| `IsPowerOfTwo(n)` | True if n is a power of two |
| `Binomial(n, k)` | n! / (k! (n-k)!) as `*big.Int` |

### Combinatorics — `num_theory`

| Function | Description |
|---|---|
| `Factorial(n)` | n! as integer |
| `FactorialBigInt(n)` | n! as `*big.Int` |
| `FactorialDigitSum(n)` | Sum of factorials of each digit of n |
| `Permutations(arr)` | All permutations of a slice |
| `PermutationCount(arr)` | Number of distinct permutations |
| `Combinations(set, n)` | All size-n combinations (pass 0 for all sizes) |
| `Partition(n)` | Partition function p(n) — ways to write n as a sum of positive integers, order ignored — as `*big.Int` |

### Number Theory — `num_theory`

| Function | Description |
|---|---|
| `Totient(n)` | Euler's totient φ(n) *(in `etc`)* |
| `ListTotients(n)` | φ(0)…φ(n) via O(n log log n) sieve *(in `etc`)* |
| `ListMobius(n)` | μ(0)…μ(n), the Möbius function, via O(n) linear sieve *(in `etc`)* |
| `DigitSum(n)` | Sum of decimal digits of n |
| `DigitSumString(s)` | Sum of decimal digits of a numeric string (arbitrary length) |
| `DigitCount(n)` | Number of decimal digits in n |
| `IsSquare(n)` | True if n is a perfect square |
| `FloatIsInteger(n)` | True if float has no fractional part |
| `IsPandigital(n)` | True if n contains digits 1–9 exactly once *(in `etc`)* |
| `IsPandigitalInBase(n, b)` | Pandigital check in an arbitrary base *(in `etc`)* |
| `ToRadians(n)` | Converts degrees to radians |

### Figurate Numbers — `figurate`

| Function | Description |
|---|---|
| `NthTriangular(n)` | nth triangular number: n(n+1)/2 |
| `IsTriangular(n)` | True if n is triangular |
| `NthPentagonal(n)` | nth pentagonal number: n(3n−1)/2 |
| `IsPentagonal(n)` | True if n is pentagonal |
| `NthHexagonal(n)` | nth hexagonal number: n(2n−1) |
| `IsHexagonal(n)` | True if n is hexagonal |

### Sequences — `fibonacci` / `sequences`

| Function | Description |
|---|---|
| `Fibonacci(limit)` | Slice of Fibonacci numbers (channel-backed) |
| `GenFibo(limit)` | Channel-based Fibonacci generator |
| `FibonacciSingle(n)` | nth Fibonacci number (F(0)=0, F(1)=1) |
| `FibonacciBig(limit)` | Fibonacci slice as `[]big.Int` |
| `GenFiboBig(limit)` | Channel-based `big.Int` Fibonacci generator |
| `FibonacciSingleBig(n)` | nth Fibonacci number as `*big.Int` (exact, iterative) |
| `Collatz(n)` | Full Collatz sequence from n to 1 *(in `sequences`)* |
| `CollatzLength(n)` | Number of steps in the Collatz sequence *(in `sequences`)* |

### Geometry — `pythagoras`

| Function | Description |
|---|---|
| `IsTriplet(a, b, c)` | True if (a, b, c) is a Pythagorean triplet |

### Vectors — `vector`

`Vector[E]` is an N-dimensional vector backed by a slice of elements. `Vector2[E]` and `Vector3[E]` embed `Vector[E]` to add named `X`/`Y`/`Z` accessors and dimension-specific operations (`Cross`, `Perp`); both require `SignedNumber`.

| Function | Description |
|---|---|
| `NewVector(elements)` | Creates a vector from a slice of elements (copies the slice) |
| `v.Len()` | Number of elements in the vector |
| `v.At(i)` | Element at index i |
| `v.Set(i, e)` | Sets the element at index i |
| `v.Elements()` | Returns a copy of the underlying elements as a slice |
| `v.Swap(i, j)` | Swaps the elements at indices i and j |
| `v.Add(o)` | Element-wise addition, returns a new vector (`o` is a `Vector` value; pass `*other` if you have `*Vector`) |
| `v.Sub(o)` | Element-wise subtraction, returns a new vector (`o` is a `Vector` value; pass `*other` if you have `*Vector`) |
| `v.Mul(o)` | Element-wise multiplication, returns a new vector (`o` is a `Vector` value; pass `*other` if you have `*Vector`) |
| `v.Div(o)` | Element-wise division, returns a new vector (`o` is a `Vector` value; pass `*other` if you have `*Vector`) |
| `v.Dot(o)` | Dot product of two vectors (`o` is a `Vector` value; pass `*other` if you have `*Vector`) |

#### Vector2 / Vector3

| Function | Description |
|---|---|
| `NewVector2(x, y)` | Creates a 2D vector |
| `NewVector3(x, y, z)` | Creates a 3D vector |
| `AsVector2(v)` | Adopts a length-2 `Vector` as a `Vector2` (shares storage, panics if length ≠ 2) |
| `AsVector3(v)` | Adopts a length-3 `Vector` as a `Vector3` (shares storage, panics if length ≠ 3) |
| `v.X()`, `v.Y()` | Component accessors (`Vector2` and `Vector3`) |
| `v.Z()` | Component accessor (`Vector3` only) |
| `v.SetX(x)`, `v.SetY(y)` | Component setters (`Vector2` and `Vector3`) |
| `v3.SetZ(z)` | Component setter (`Vector3` only) |
| `v2.Cross(o)` | Perp-dot product (z-component of the 3D cross product) — `Vector2` only |
| `v2.Perp()` | Returns the perpendicular vector `(-y, x)` — `Vector2` only |
| `v3.Cross(o)` | 3D cross product, returns a new `Vector3` |
| `v.Add(o)`, `v.Sub(o)`, `v.Mul(o)`, `v.Div(o)` | Element-wise arithmetic (`o` is a value; pass `*other` if you have a pointer), returns the same dimension type (`Vector2` and `Vector3`) |

#### Coordinate operations

`HexCoordinate[E]` is a point in axial hexagonal coordinates (q, r); it embeds `Vector2[E]` and adds `Q()`/`R()`/`S()` accessors.

| Function | Description |
|---|---|
| `AreColinear(p1, p2, p3)` | True if `p1`, `p2`, and `p3` all lie on a common line |
| `NewHexCoordinate(q, r)` | Creates an axial hex coordinate |
| `h.Q()`, `h.R()` | Axial coordinate accessors (aliases for the inherited `X()`/`Y()`) |
| `h.S()` | Third cube coordinate, `s = -q-r` |
| `HexAreColinear(p1, p2, p3)` | True if `p1`, `p2`, and `p3` all lie on a common line on the hex grid |
| `Quadrant(p)` | Standard 1–4 quadrant number of `p`, sweeping counterclockwise from the positive x-axis; each quadrant includes its lower (starting) boundary and excludes its upper one, so axis points belong to the quadrant they start (origin is quadrant IV by convention) |
| `SameQuadrant(p1, p2)` | True if `p1` and `p2` lie in the same of the four quadrants returned by `Quadrant` |
| `HexSameQuadrant(p1, p2)` | True if `p1` and `p2` lie in the same quadrant of the axial (q, r) coordinate system |
| `CartesianToHex(p)` | Converts a cartesian point to fractional axial hex coordinates (unit-size, pointy-top hexagons; requires `Float`) |
| `HexToCartesian(h)` | Converts axial hex coordinates to a cartesian point (inverse of `CartesianToHex`; requires `Float`) |
| `RoundHex(h)` | Rounds fractional axial hex coordinates to the nearest hex cell (requires `Float`) |

### Slice Utilities — `etc`

| Function | Description |
|---|---|
| `Sum(slice)` | Sum of all elements |
| `Product(slice)` | Product of all elements |
| `Min(a, b, ...)` | Minimum value |
| `Max(a, b, ...)` | Maximum value |
| `Abs(n)` | Absolute value of a signed integer |
| `SumOfSquares(n)` | 1² + 2² + … + n² |
| `SquareOfSum(n)` | (1 + 2 + … + n)² |
| `Range(start, stop)` | Slice from start to stop (exclusive), stepping by 1 |
| `RangeStep(start, stop, step)` | Range with custom step |
| `Filter(slice, fn)` | Keep elements matching predicate |
| `Map(slice, fn)` | Transform every element |
| `Reduce(slice, fn)` | Fold left, starting from slice[0] *(in `num_theory`)* |
| `ReduceWithInit(init, slice, fn)` | Fold left with an explicit initial value *(in `num_theory`)* |
| `Sort(slice, fn)` | Sort in-place using a less function |
| `Unique(slice)` | Remove duplicates (requires `comparable`) |
| `UniqueCount(slice)` | Map of element → count |
| `RemoveDuplicates(slice)` | Remove duplicates (works on `Comparable` constraint) |
| `RemoveDuplicatesFunc(slice, fn)` | Remove duplicates with custom equality function |
| `SliceContains(slice, v)` | True if slice contains v |
| `SliceContainsAny(slice, v, fn)` | True if any element matches v via fn |
| `CountOccurrenceInSlice(slice, v)` | Count occurrences of v |
| `RemoveFromSlice(slice, i)` | Remove element at index i |
| `GenerateSlice(n, value)` | Create a slice of length n filled with value |
| `JoinSlice(slice)` | Concatenate integer slice into a string |

### String Utilities — `etc`

| Function | Description |
|---|---|
| `IsPalindrome(s)` | True if s reads the same forwards and backwards |
| `ReverseString(s)` | Reverses a string (Unicode-safe) |
| `MakeIntSlice(n)` | Digits of n as a `[]E` slice |
| `DecimalToBase(n, b)` | Converts n to base b (2–62) as a string |

### File I/O — `etc`

| Function | Description |
|---|---|
| `CreateFile(name)` | Create a file |
| `CreateFileWithContent(name, content)` | Create a file and write content to it |
| `ReadFile(name)` | Read a file's content as a string |
| `FileExists(name)` | True if the file exists |

### Progress Bar — `progress`

Requires the `range` keyword (Go 1.23 range-over-func / Go 1.22 range-over-int). Rendering is delegated to [schollz/progressbar](https://github.com/schollz/progressbar).

| Function | Description |
|---|---|
| `ProgressBar(slice, opts...)` | Ranges over a slice, drawing a live progress bar as iteration proceeds |
| `ProgressBarN(n, opts...)` | Ranges over `[0, n)`, drawing a live progress bar as iteration proceeds |
| `Spinner(opts...)` | Ranges over `[0, 1, 2, ...)` with no upper bound, drawing a live spinner (no % or ETA) — for open-ended brute-force search where the total is unknown; break the loop when done |
| `WithWidth(n)` | Option: number of characters in the bar itself (default 10) |
| `WithWriter(w)` | Option: where the bar is rendered (default `os.Stderr`) |
| `WithLabel(s)` | Option: description printed before the bar |

`Option` is an alias for `progressbar.Option`, so any option from the underlying library (colors, spinner style, byte mode, ...) can be passed straight through too.

### Timing — `stopwatch`

| Function | Description |
|---|---|
| `Start(label, opts...)` | Starts a timer, returns a `Stop` func — call (or `defer`) it to print `"<label>: <elapsed>"` and get the elapsed `time.Duration` |
| `Time(label, fn, opts...)` | Runs fn, printing and returning the elapsed time the same way `Start`'s returned func does |
| `WithWriter(w)` | Option: where the elapsed-time line is printed (default `os.Stderr`) |

```go
defer stopwatch.Start("part 2")()
```

### Memoization — `memo`

| Function | Description |
|---|---|
| `Memoize(fn)` | Wraps a single-argument function so repeat calls with the same argument are computed once |
| `Memoize2(fn)` | `Memoize` for a two-argument function — handy for `(position, state)`-style recurrences |

Not safe for concurrent use. For a recursive function, declare the variable first so the memoized wrapper can call itself:

```go
var fib func(n int) int64
fib = memo.Memoize(func(n int) int64 {
	if n < 2 {
		return int64(n)
	}
	return fib(n-1) + fib(n-2)
})
```

### Type Utilities — `utils`

| Function | Description |
|---|---|
| `SameType[A, B]()` | True if A and B are the same type or type family |
| `IsSignedInteger(t)` | Reflection-based type check |
| `IsUnsignedInteger(t)` | Reflection-based type check |
| `IsInteger(t)` | Reflection-based type check |
| `IsFloat(t)` | Reflection-based type check |
| `IsRealNumber(t)` | Reflection-based type check |
| `IsComplex(t)` | Reflection-based type check |
| `IsNumber(t)` | Reflection-based type check |
| `IsComparable(t)` | Reflection-based type check |

---

## Type Constraints — `utils`

```go
SignedInteger   // ~int | ~int8 | ~int16 | ~int32 | ~int64
UnsignedInteger // ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
Integer         // SignedInteger | UnsignedInteger
Float           // ~float32 | ~float64
RealNumber      // Integer | Float
Complex         // ~complex64 | ~complex128
Number          // RealNumber | Complex
Comparable      // Number | ~string | ~bool
SignedNumber    // SignedInteger | Float | Complex
```

---

## Examples

```go
import (
	"fmt"

	"github.com/t1sk3/eulerlib/v2/etc"
	"github.com/t1sk3/eulerlib/v2/fibonacci"
	"github.com/t1sk3/eulerlib/v2/figurate"
	"github.com/t1sk3/eulerlib/v2/num_theory"
	"github.com/t1sk3/eulerlib/v2/prime_numbers"
	"github.com/t1sk3/eulerlib/v2/progress"
	"github.com/t1sk3/eulerlib/v2/sequences"
	"github.com/t1sk3/eulerlib/v2/vector"
)

// Check if 28 is a perfect number
num_theory.IsPerfect(28) // true — SumDivisors(28) = 1+2+4+7+14 = 28

// Find the 10,001st prime
iter := prime_numbers.NewPrimeNumberIterator[int]()
for i := 0; i < 10001; i++ {
    iter.Next()
}
fmt.Println(iter.Current()) // 104743

// Sum of all primes below 2,000,000 (Euler #10)
fmt.Println(prime_numbers.SumPrimes(0, 2_000_000)) // 142913828922

// Factorize a number
num_theory.Factorize(int64(360)) // map[2:3 3:2 5:1]  (2³ × 3² × 5)

// Collatz sequence length for 27
sequences.CollatzLength(27) // 112

// 100th Fibonacci number (exact)
fibonacci.FibonacciSingleBig(100).String() // "354224848179261915075"

// Euler's totient sieve up to 100
phi := etc.ListTotients(100) // phi[i] = φ(i)

// Check figurate numbers
figurate.IsTriangular(55)  // true  (T(10) = 55)
figurate.IsPentagonal(51)  // true  (P(6)  = 51)
figurate.IsHexagonal(45)   // true  (H(5)  = 45)

// N-dimensional vector arithmetic
a := vector.NewVector([]int{1, 2, 3})
b := vector.NewVector([]int{4, 5, 6})
a.Dot(*b) // 32  (1*4 + 2*5 + 3*6)

// 2D cross (perp-dot) product
p := vector.NewVector2(1, 2)
q := vector.NewVector2(3, 4)
p.Cross(*q) // -2  (1*4 - 2*3)

// 3D cross product
i := vector.NewVector3(1, 0, 0)
j := vector.NewVector3(0, 1, 0)
i.Cross(*j) // Vector3{0, 0, 1}

// Live progress bar while ranging over a slice
primes := prime_numbers.ListPrimes(1000)
for i, p := range progress.ProgressBar(primes) {
    _ = i
    _ = p // do something with each prime
}
// 100% |████████████████████████████████████████| (168/168, 9000 it/s)

// ...or over a plain count
for i := range progress.ProgressBarN(1_000_000, progress.WithLabel("euler #10")) {
    _ = i
}
```

---

## License

[MIT](LICENSE)
