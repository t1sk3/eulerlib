# eulerlib

A generic Go library of number-theory utilities for solving [Project Euler](https://projecteuler.net/) problems. All functions are generic over Go's built-in numeric types via type constraints.

```
go get github.com/t1sk3/eulerlib
```

Requires **Go 1.23+**.

---

## Reference

### Primes

| Function | Description |
|---|---|
| `IsPrime(n)` | Returns true if n is prime |
| `NextPrime(n)` | Returns the next prime after n |
| `ListPrimes(n)` | Returns all primes up to n |
| `ListPrimality(n)` | Sieve of Eratosthenes — returns `[]bool` indexed by number |
| `FirstNPrimes(n)` | Returns the first n primes |
| `PrimeCount(s, e)` | Counts primes in the range [s, e] |
| `SumPrimes(s, e)` | Sums all primes in the range [s, e] |
| `PrimeGenerator(limit)` | Channel-based prime generator up to limit |
| `NewPrimeNumberIterator(start?)` | Iterator that steps through primes one at a time |
| `PrimeFactors(n)` | Returns the prime factors of n (with repeats) |
| `PrimeFactorsBigInt(n)` | Prime factors of a `*big.Int` as `[][]int64{prime, exp}` |

### Divisors

| Function | Description |
|---|---|
| `Divisors(n)` | Returns all divisors of n |
| `CountDivisors(n)` | Returns the number of divisors of n |
| `SumDivisors(n)` | Sum of proper divisors of n (excludes n itself) |
| `IsAbundant(n)` | True if `SumDivisors(n) > n` |
| `IsDeficient(n)` | True if `SumDivisors(n) < n` |
| `IsPerfect(n)` | True if `SumDivisors(n) == n` (e.g. 6, 28, 496) |
| `IsAmicable(n)` | True if n and `SumDivisors(n)` form an amicable pair |

### Factorization & GCD

| Function | Description |
|---|---|
| `Factorize(n)` | Prime factorization as `map[E]E{prime: exponent}` |
| `FactorizeBigInt(n)` | Prime factorization of `*big.Int` as `[][2]*big.Int` |
| `Gcd(a, b, ...)` | Greatest common divisor |
| `Lcm(a, b, ...)` | Least common multiple |

### Powers & Modular Arithmetic

| Function | Description |
|---|---|
| `Pow(b, n)` | b^n as an integer |
| `PowBigInt(b, n)` | b^n as `*big.Int` (non-mutating, binary exponentiation) |
| `PowBigFloat(b, n)` | b^n as `*big.Float` (non-mutating, binary exponentiation) |
| `PowMod(x, y, p)` | x^y mod p (requires p < 2³¹) |
| `IsPowerOfTwo(n)` | True if n is a power of two |
| `Binomial(n, k)` | n! / (k! (n-k)!) as `*big.Int` |

### Combinatorics

| Function | Description |
|---|---|
| `Factorial(n)` | n! as integer |
| `FactorialBigInt(n)` | n! as `*big.Int` |
| `FactorialDigitSum(n)` | Sum of factorials of each digit of n |
| `Permutations(arr)` | All permutations of a slice |
| `PermutationCount(arr)` | Number of distinct permutations |
| `Combinations(set, n)` | All size-n combinations (pass 0 for all sizes) |

### Number Theory

| Function | Description |
|---|---|
| `Totient(n)` | Euler's totient φ(n) |
| `ListTotients(n)` | φ(0)…φ(n) via O(n log log n) sieve |
| `DigitSum(n)` | Sum of decimal digits of n |
| `DigitSumString(s)` | Sum of decimal digits of a numeric string (arbitrary length) |
| `IsSquare(n)` | True if n is a perfect square |
| `FloatIsInteger(n)` | True if float has no fractional part |
| `DecimalToBase(n, b)` | Converts n to base b (2–62) as a string |
| `IsPandigital(n)` | True if n contains digits 1–9 exactly once |
| `IsPandigitalInBase(n, b)` | Pandigital check in an arbitrary base |

### Figurate Numbers

| Function | Description |
|---|---|
| `NthTriangular(n)` | nth triangular number: n(n+1)/2 |
| `IsTriangular(n)` | True if n is triangular |
| `NthPentagonal(n)` | nth pentagonal number: n(3n−1)/2 |
| `IsPentagonal(n)` | True if n is pentagonal |
| `NthHexagonal(n)` | nth hexagonal number: n(2n−1) |
| `IsHexagonal(n)` | True if n is hexagonal |

### Sequences

| Function | Description |
|---|---|
| `Fibonacci(limit)` | Slice of Fibonacci numbers (channel-backed) |
| `GenFibo(limit)` | Channel-based Fibonacci generator |
| `FibonacciSingle(n)` | nth Fibonacci number (F(0)=0, F(1)=1) |
| `FibonacciBig(limit)` | Fibonacci slice as `[]big.Int` |
| `GenFiboBig(limit)` | Channel-based `big.Int` Fibonacci generator |
| `FibonacciSingleBig(n)` | nth Fibonacci number as `*big.Int` (exact, iterative) |
| `Collatz(n)` | Full Collatz sequence from n to 1 |
| `CollatzLength(n)` | Number of steps in the Collatz sequence |

### Geometry

| Function | Description |
|---|---|
| `IsTriplet(a, b, c)` | True if (a, b, c) is a Pythagorean triplet |
| `ToRadians(n)` | Converts degrees to radians |

### Vectors

`Vector[E]` is an N-dimensional vector backed by a slice of elements. `Vector2[E]` and `Vector3[E]` embed `Vector[E]` to add named `X`/`Y`/`Z` accessors and dimension-specific operations (`Cross`, `Perp`); both require `SignedNumber`.

| Function | Description |
|---|---|
| `NewVector(elements)` | Creates a vector from a slice of elements (copies the slice) |
| `v.Len()` | Number of elements in the vector |
| `v.At(i)` | Element at index i |
| `v.Set(i, e)` | Sets the element at index i |
| `v.Elements()` | Returns a copy of the underlying elements as a slice |
| `v.Swap(i, j)` | Swaps the elements at indices i and j |
| `v.Add(o)` | Element-wise addition, returns a new vector |
| `v.Sub(o)` | Element-wise subtraction, returns a new vector |
| `v.Mul(o)` | Element-wise multiplication, returns a new vector |
| `v.Div(o)` | Element-wise division, returns a new vector |
| `v.Dot(o)` | Dot product of two vectors |

#### Vector2 / Vector3

| Function | Description |
|---|---|
| `NewVector2(x, y)` | Creates a 2D vector |
| `NewVector3(x, y, z)` | Creates a 3D vector |
| `AsVector2(v)` | Adopts a length-2 `Vector` as a `Vector2` (shares storage, panics if length ≠ 2) |
| `AsVector3(v)` | Adopts a length-3 `Vector` as a `Vector3` (shares storage, panics if length ≠ 3) |
| `v.X()`, `v.Y()` | Component accessors (`Vector2` and `Vector3`) |
| `v.Z()` | Component accessor (`Vector3` only) |
| `v.SetX(x)`, `v.SetY(y)` | Component setters (`Vector2` only) |
| `v2.Cross(o)` | Perp-dot product (z-component of the 3D cross product) — `Vector2` only |
| `v2.Perp()` | Returns the perpendicular vector `(-y, x)` — `Vector2` only |
| `v3.Cross(o)` | 3D cross product, returns a new `Vector3` |
| `v.Add(o)`, `v.Sub(o)` | Element-wise addition/subtraction, returns the same dimension type |

### Slice Utilities

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
| `Reduce(slice, fn)` | Fold left, starting from slice[0] |
| `ReduceWithInit(init, slice, fn)` | Fold left with an explicit initial value |
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

### String Utilities

| Function | Description |
|---|---|
| `IsPalindrome(s)` | True if s reads the same forwards and backwards |
| `ReverseString(s)` | Reverses a string (Unicode-safe) |
| `MakeIntSlice(n)` | Digits of n as a `[]E` slice |

### File I/O

| Function | Description |
|---|---|
| `CreateFile(name)` | Create a file |
| `CreateFileWithContent(name, content)` | Create a file and write content to it |
| `ReadFile(name)` | Read a file's content as a string |
| `FileExists(name)` | True if the file exists |

### Type Utilities

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

## Type Constraints

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
// Check if 28 is a perfect number
eulerlib.IsPerfect(28) // true — SumDivisors(28) = 1+2+4+7+14 = 28

// Find the 10,001st prime
iter := eulerlib.NewPrimeNumberIterator[int]()
for i := 0; i < 10001; i++ {
    iter.Next()
}
fmt.Println(iter.Current()) // 104743

// Sum of all primes below 2,000,000 (Euler #10)
fmt.Println(eulerlib.SumPrimes(0, 2_000_000)) // 142913828922

// Factorize a number
eulerlib.Factorize(int64(360)) // map[2:3 3:2 5:1]  (2³ × 3² × 5)

// Collatz sequence length for 27
eulerlib.CollatzLength(27) // 112

// 100th Fibonacci number (exact)
eulerlib.FibonacciSingleBig(100).String() // "354224848179261915075"

// Euler's totient sieve up to 100
phi := eulerlib.ListTotients(100) // phi[i] = φ(i)

// Check figurate numbers
eulerlib.IsTriangular(55)  // true  (T(10) = 55)
eulerlib.IsPentagonal(51)  // true  (P(6)  = 51)
eulerlib.IsHexagonal(45)   // true  (H(5)  = 45)

// N-dimensional vector arithmetic
a := eulerlib.NewVector([]int{1, 2, 3})
b := eulerlib.NewVector([]int{4, 5, 6})
a.Dot(*b) // 32  (1*4 + 2*5 + 3*6)

// 2D cross (perp-dot) product
p := eulerlib.NewVector2(1, 2)
q := eulerlib.NewVector2(3, 4)
p.Cross(*q) // -2  (1*4 - 2*3)

// 3D cross product
i := eulerlib.NewVector3(1, 0, 0)
j := eulerlib.NewVector3(0, 1, 0)
i.Cross(*j) // Vector3{0, 0, 1}
```
