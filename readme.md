# eulerlib

A generic Go library of number-theory utilities for solving [Project Euler](https://projecteuler.net/) problems. Functions are generic over Go's integer and float types, lazy sequences are [go-utils](https://github.com/t1sk3/go-utils) `stream.Stream`s, and everything is exact over the full range of 64-bit integers.

```
go get github.com/t1sk3/eulerlib/v4
```

Requires **Go 1.27+**.

> **v4 is a breaking change from v3**: packages were renamed and reorganised, and `etc`, `utils` and `vector` are gone. See [Upgrading from v3](#upgrading-from-v3) and [CHANGELOG.md](CHANGELOG.md).

```go
import (
	"github.com/t1sk3/eulerlib/v4/numtheory"
	"github.com/t1sk3/eulerlib/v4/primes"
	"github.com/t1sk3/go-utils/stream"
)

primes.IsPrime(uint64(18446744073709551557))         // true (largest 64-bit prime)
numtheory.Factorize(int64(360))                       // map[2:3 3:2 5:1]
stream.Sum(primes.Between(2, 2_000_000))              // 142913828922 (Euler #10)
primes.All[int]().Skip(10_000).FindFirst().Get()      // 104743 (Euler #7)
```

| Package | Covers |
|---|---|
| [`primes`](primes) | Primality (deterministic Miller–Rabin), bitset `Sieve`, segmented-sieve prime streams, smallest-factor table, prime counting π(n) |
| [`numtheory`](numtheory) | Divisors and divisor sums, factorization (Pollard's rho), gcd/lcm, modular arithmetic, powers and roots, totient/Möbius, multiplicative order, modular square roots, continued fractions and Pell's equation, predicate builders, collectors |
| [`combinatorics`](combinatorics) | Factorial, binomial (also mod p), Catalan/Stirling/Bell numbers, partitions, coin-change counts, lazy permutations/combinations/subsets/Cartesian products |
| [`digits`](digits) | Digit slices and sums, reversal, rotations, digit permutations, palindromes, pandigital checks, concatenation, base conversion |
| [`fibonacci`](fibonacci) | Fibonacci streams and n-th terms (incl. `big.Int`, fast doubling) |
| [`figurate`](figurate) | Triangular, square, pentagonal, hexagonal and s-gonal numbers, exact tests, streams |
| [`sequences`](sequences) | Collatz sequence and lengths, Farey sequence, happy numbers, digit chains, linear recurrences mod m |
| [`paths`](paths) | Best path sums through triangles and grids |
| [`roman`](roman) | Roman numerals |
| [`pythagoras`](pythagoras) | Pythagorean `Triplet`, triplet streams via Euclid's formula |
| [`geometry`](geometry) | Colinearity and quadrants on go-utils `vector.Vec2`, axial hex-grid coordinates |
| [`memo`](memo) | Memoization wrappers for one to three arguments, including a goroutine-safe one |
| [`progress`](progress) | Live progress bars over slices, counts, `iter.Seq`s and streams |
| [`stopwatch`](stopwatch) | Timing helper for brute-force runs |
| [`plot`](plot) | Interactive HTML scatter plots (2D/3D) and CSV export |

Every package has runnable examples on [pkg.go.dev](https://pkg.go.dev/github.com/t1sk3/eulerlib/v4), and [`example_euler_test.go`](example_euler_test.go) solves Euler #1, #5, #7, #10, #14, #21, #24 and #45, checking each answer against a brute-force solution.

---

## Working with streams

Infinite or open-ended sequences are `stream.Stream[E]` values. They are lazy, can be consumed any number of times (each consumption starts over), and stop at the largest value of their element type instead of overflowing. Bound them with `Limit` or `TakeWhile`:

```go
primes.All[int]().TakeWhile(stream.LessThan(30)).ToSlice() // [2 3 5 7 11 13 17 19 23 29]
fibonacci.All[uint8]().FindLast().Get()                    // 233, the largest that fits
```

Key/value results are `iter.Seq2`, so you can `range` over them or collect them into a map:

```go
for p, e := range numtheory.Factors(360) { ... } // (2, 3), (3, 2), (5, 1), ascending
maps.Collect(numtheory.Factors(360))             // map[2:3 3:2 5:1]
```

Reductions that come up often are collectors, and predicates are plain `func(E) bool` (so they also work with `slices.IndexFunc` and friends):

```go
stream.RangeClosed(1, 20).Collect(numtheory.LcmOf[int]())                   // 232792560 (Euler #5)
stream.Sum(stream.Range(1, 1000).Filter(numtheory.DivisibleByAny(3, 5)))    // 233168    (Euler #1)
figurate.Hexagonals[int64]().Filter(figurate.IsPentagonal[int64]).Skip(2).FindFirst() // Euler #45
```

Enumerations yield a fresh slice per element and never touch your input, so they are safe to keep or modify, and are only generated as far as you read:

```go
for p := range combinatorics.Permutations(items).All() {
	if check(p) {
		break // nothing past this permutation is computed
	}
}
```

For expensive per-element work, go-utils' parallel stages apply directly: `s.Parallel(runtime.GOMAXPROCS(0)).Map(expensive).Sequential()`. That only pays off when each element costs more than a couple of microseconds.

`stream.Sum`, `stream.Max`, `stream.Distinct` and similar are package functions, because they need a constraint on the element type. Everything else (`Map`, `Filter`, `Fold`, `Collect`, ...) is a method. See the [go-utils stream docs](https://pkg.go.dev/github.com/t1sk3/go-utils/stream).

---

## Reference

Invalid arguments that are programmer errors (a zero modulus, a negative exponent, a base out of range) panic, and each function documents its panics. Arithmetic that builds a number larger than E wraps silently, as ordinary Go arithmetic does, unless the function says otherwise.

### Primes — `primes`

| Function | Description |
|---|---|
| `IsPrime(n)` | Deterministic Miller–Rabin: exact for every value of every integer type |
| `NextPrime(n)` / `PrevPrime(n)` | Smallest prime > n (panics if it doesn't fit in E) / largest prime < n as `Optional` |
| `All[E]()`, `From(n)`, `Between(lo, hi)` | Streams of primes (segmented sieve; Miller–Rabin above 2^48) |
| `UpTo(n)`, `First[E](n)` | Primes ≤ n / the first n primes, as slices |
| `Count(lo, hi)`, `Sum(lo, hi)` | Count / sum of primes in `[lo, hi]` |
| `NewSieve(limit)` | Bitset sieve with `Has`, `IsPrime` (usable as a predicate), `Count`, `All`, `Stream`, `Limit` |
| `SmallestFactorSieve(n)` | `spf[i]` = smallest prime factor of i, for `numtheory.FactorizeSPF` |
| `Pi(n)`, `SumUpTo(n)` | Count / sum (`*big.Int`) of the primes ≤ n without enumerating them: O(n^¾) time, O(√n) memory |
| `Nth[E](n)` | The n-th prime, `Nth(1)` = 2 (Euler #7) |
| `IsPrimeBig(n)` | Primality of a `*big.Int` |
| `IsCircular(n)`, `IsTruncatable(n)`, `IsLeftTruncatable`, `IsRightTruncatable` | Every rotation / truncation is prime (Euler #35, #37) |

### Number theory — `numtheory`

| Function | Description |
|---|---|
| `Divisors(n)`, `CountDivisors(n)`, `SumDivisors(n)` | Ascending divisors / their count / sum of proper divisors |
| `SigmaSieve(n)` | `SumDivisors(i)` for every i ≤ n in O(n log n) |
| `IsAbundant`, `IsDeficient`, `IsPerfect`, `IsAmicable` | Divisor-sum predicates |
| `AmicablePairs(limit)` | Stream of amicable pairs with both members < limit |
| `Factors(n)` | `iter.Seq2[prime, exponent]`, ascending; fast for every 64-bit n |
| `Factorize(n)`, `FactorizeSPF(n, spf)`, `FactorsSPF` | Factorization as `map[E]int` / using a smallest-factor table |
| `PrimeFactors(n)`, `DistinctPrimeFactors(n)` | With and without repeats, ascending |
| `FactorizeBig(n)` | `[]pair.Pair[*big.Int, int]` |
| `Gcd(...)`, `Lcm(...)`, `IsCoprime`, `ExtGcd` | Never negative; `Gcd()` = 0, `Lcm()` = 1 |
| `ModInverse(a, m)`, `CRT(rs, ms)` | Return `optional.Optional` |
| `PowMod(x, y, m)`, `MulMod(a, b, m)` | Exact for any modulus (128-bit intermediate products) |
| `Pow`, `PowBig`, `PowBigFloat`, `IsPowerOfTwo`, `IsSquare`, `Isqrt`, `FloatIsInteger` | Powers and roots; `IsSquare`/`Isqrt` are exact |
| `Totient`, `TotientSieve`, `MobiusSieve`, `DivisorCountSieve` | φ(n), φ(0..n), μ(0..n), d(0..n) |
| `Iroot(n, k)`, `PerfectPower(n)` | Exact integer k-th root / n as base^exp with the largest exp |
| `SumOfTwoSquares(n)` | `(a, b)` with a² + b² = n, as `Optional` |
| `MultiplicativeOrder(a, m)`, `RecurringCycleLength(d)` | Order of a mod m as `Optional` / cycle length of 1/d in decimal (Euler #26) |
| `Jacobi(a, n)`, `ModSqrt(a, p)`, `PrimitiveRoot(p)` | Jacobi/Legendre symbol, Tonelli–Shanks square root as `Optional`, smallest primitive root |
| `SqrtContinuedFraction(n)`, `SqrtCFTerms(n)`, `ECFTerms()` | Continued fraction of √n as `(a0, period)` / as an endless `iter.Seq` / of e (Euler #64, #65) |
| `Convergents(terms)` | `iter.Seq2` of numerator and denominator (`*big.Int`) for any term sequence (Euler #57, #65) |
| `Pell(d)` | Fundamental solution of x² − d·y² = 1 as `Optional` (Euler #66) |
| `Abs`, `SumOfSquares`, `SquareOfSum`, `ToRadians` | Small helpers |
| `DivisibleBy(k)`, `DivisibleByAny(ks...)`, `CoprimeTo(k)` | Predicate builders for `Filter` |
| `Product[E]()`, `GcdOf[E]()`, `LcmOf[E]()` | Collectors for `Collect` |

### Combinatorics — `combinatorics`

| Function | Description |
|---|---|
| `Factorial`, `FactorialBig`, `Binomial` (big), `Partition`, `PartitionTable` | Counting |
| `Catalan(n)`, `Derangements(n)`, `Bell(n)`, `Stirling1(n, k)`, `Stirling2(n, k)` | Counting sequences, as `*big.Int` |
| `WaysToSum(target, parts)` | Coin-change count: ways to write target as a sum of the parts (Euler #31, #76, #77) |
| `BinomialMod(n, k, p)`, `FactorialsMod(n, p)` | C(n, k) mod a prime (Lucas) / tables of i! and its inverse mod p |
| `PermutationCount(items)` | Distinct orderings of a multiset |
| `DigitFactorialSum(n)` | Sum of the factorials of n's digits (Euler #34) |
| `Permutations(items)` | Stream, Heap's algorithm |
| `LexPermutations(items)` | Stream, lexicographic, skips duplicates |
| `NthPermutation(items, k)` | k-th lexicographic permutation without generating the others (Euler #24) |
| `Combinations(items, k)`, `Subsets(items)`, `CartesianProduct(sets...)` | Streams; no size limit |
| `KPermutations(items, k)`, `CombinationsWithReplacement(items, k)` | Streams of ordered k-selections / k-multisets |
| `Partitions(n)` | Stream of the integer partitions of n, parts descending |

### Digits — `digits`

| Function | Description |
|---|---|
| `Of(n)`, `OfBase(n, b)`, `FromDigits(ds)` | Digit slices, most significant first |
| `Sum(n)`, `SumString(s)`, `Count(n)`, `SumOf[E]()` | Digit sum / digit sum of a string / digit count / collector |
| `Reverse(n)` | 1230 → 321 |
| `Rotations(n)` | 197 → [197 971 719] (Euler #35) |
| `Signature(n)`, `IsPermutation(a, b)` | Digit counts as a comparable map key / same digits in another order (Euler #49, #52, #62) |
| `PowerSum(n, k)` | Sum of the k-th powers of the digits (Euler #30) |
| `LookAndSay(s)`, `ChampernowneDigit(n)` | Next look-and-say term / n-th digit of 0.123456789101112… (Euler #40) |
| `IsPalindrome(n)`, `IsPalindromeInBase(n, b)`, `IsPalindromeString(s)` | Palindromes (Euler #4, #36) |
| `IsPandigital(n)`, `IsPandigitalN(n, k)` | Digits 1..9 / 1..k each exactly once |
| `IsPandigitalInBase(n, b)`, `IsPandigitalString(s, b)` | Every digit 0..b-1 at least once |
| `Concat(ns...)`, `ConcatString(ns...)` | `Concat(1, 23, 4)` = 1234 |
| `ToBase(n, b)` | Bases 2..62 (`0-9A-Za-z`) |

### Sequences — `fibonacci`, `figurate`, `sequences`, `pythagoras`

| Function | Description |
|---|---|
| `fibonacci.All[E]()`, `AllBig()`, `First`, `Nth`, `NthBig`, `IndexOfFirstWithDigits` | F(0) = 0, F(1) = 1; `NthBig` uses fast doubling |
| `figurate.NthTriangular`, `NthSquare`, `NthPentagonal`, `NthHexagonal`, `NthPolygonal(s, n)` | n-th figurate numbers |
| `figurate.IsTriangular`, `IsPentagonal`, `IsHexagonal`, `IsPolygonal(s, n)` | Exact for every value of E |
| `figurate.PolygonalIndex(s, n)` | The k with P(s, k) = n, as `Optional` |
| `figurate.Triangulars`, `Squares`, `Pentagonals`, `Hexagonals`, `Polygonals(s)` | Streams starting at n = 1 |
| `sequences.Collatz(n)`, `CollatzLength(n)`, `CollatzLengths(limit)`, `LongestCollatz(limit)` | Collatz (Euler #14) |
| `sequences.Farey(n)` | Stream of the reduced fractions in [0, 1] with denominator ≤ n, ascending (Euler #71–73) |
| `sequences.IsHappy(n)`, `ChainLength(start, next)` | Happy numbers / distinct terms before a chain repeats (Euler #74, #92) |
| `sequences.LinearRecurrenceMod(coeffs, initial, n, m)` | n-th term of a linear recurrence mod m by matrix exponentiation |
| `pythagoras.Triplet{A, B, C}`, `IsTriplet`, `PrimitiveTriplets(maxC)`, `TripletsWithPerimeterAtMost(p)` | Pythagorean triplets (Euler #9, #39, #75) |

### Paths and numerals — `paths`, `roman`

| Function | Description |
|---|---|
| `paths.TriangleMax(tri)` | Largest top-to-bottom path sum in a triangle (Euler #18, #67) |
| `paths.GridMinTwoWays(g)` | Smallest path sum, top-left to bottom-right, moving right and down (Euler #81) |
| `paths.GridMinThreeWays(g)` | Left column to right column, moving up, down and right (Euler #82) |
| `paths.GridMinFourWays(g)` | Top-left to bottom-right, moving in all four directions, Dijkstra (Euler #83) |
| `roman.Format(n)`, `roman.Parse(s)` | Minimal Roman numeral / value of a numeral, minimal or not (Euler #89) |

### Geometry — `geometry`

Built on [go-utils' `vector`](https://pkg.go.dev/github.com/t1sk3/go-utils/types/vector) (`Vec2`, `Vec3`, `VecN`), which provides all vector arithmetic.

| Function | Description |
|---|---|
| `AreColinear(p1, p2, p3)` | Three `vector.Vec2` points on one line |
| `Quadrant(p)`, `SameQuadrant(p1, p2)` | Quadrant 1–4; axes belong to the quadrant they start, the origin is IV |
| `Hex[E]{Q, R}`, `HexOf(q, r)` | Axial hex coordinate: `S`, `Add`, `Sub`, `Scale`, `Distance`, `Neighbors`, `ToVec2`, `String`; comparable, JSON `{"q", "r"}` |
| `HexAreColinear`, `HexSameQuadrant`, `HexRing(center, radius)` | Hex-grid helpers |
| `CartesianToHex`, `HexToCartesian`, `RoundHex` | Conversion for unit pointy-top hexagons |

### Progress bars — `progress`

| Function | Description |
|---|---|
| `Slice(s, opts...)` | `iter.Seq2[int, T]` over a slice |
| `N(n, opts...)` | `iter.Seq[int]` over `[0, n)` |
| `Spinner(opts...)` | Unbounded count with a spinner |
| `Seq(seq, total, opts...)`, `Stream(s, total, opts...)` | Wrap any `iter.Seq` or stream; `total < 0` draws a spinner |
| `WithWidth`, `WithWriter`, `WithLabel` | Options; any `progressbar.Option` also works |

The bar is created when iteration starts, so a sequence can be ranged over more than once.

```go
total := primes.Count(2, 2_000_000)
sum := stream.Sum(progress.Stream(primes.Between(2, 2_000_000), total))
```

### Timing — `stopwatch`

`Start(label, opts...)` returns a stop func that prints `"<label>: <elapsed>"` and returns the duration; `Time(label, fn, opts...)` times `fn`. `WithWriter(w)` redirects the output.

```go
defer stopwatch.Start("part 2")()
```

### Memoization — `memo`

| Function | Description |
|---|---|
| `Memoize(fn)`, `Memoize2(fn)`, `Memoize3(fn)` | Cache a one-, two- or three-argument function (not goroutine-safe) |
| `MemoizeWithClear(fn)` | `Memoize` plus a function that empties the cache |
| `MemoizeConcurrent(fn)` | Goroutine-safe; `fn` runs without a lock held, so recursion is fine |

```go
var fib func(n int) int64
fib = memo.Memoize(func(n int) int64 {
	if n < 2 {
		return int64(n)
	}
	return fib(n-1) + fib(n-2)
})
```

### Plotting — `plot`

| Function | Description |
|---|---|
| `Scatter2D(points, path, opts...)`, `Scatter3D(...)` | Interactive HTML scatter plots (ECharts via [go-echarts](https://github.com/go-echarts/go-echarts)); returns the HTML and writes it to `path` if non-empty |
| `SaveCSV2D`, `SaveCSV3D` | CSV export |
| `PointFromVec2`, `PointFromVec3`, `PointFromPair`, `Points2D(seq)`, `Points3D(seq)` | Build points from go-utils vectors, pairs and streams |

Options: `WithTitle`, `WithXLabel`, `WithYLabel`, `WithZLabel`, `WithSize`, `WithPointSize`, `WithColor`, `WithTheme`.

```go
pts := plot.Points2D(stream.RangeClosed(1, 10).
	Map(func(n int) vector.Vec2[int] { return vector.Of2(n, n*n) }).All())
plot.Scatter2D(pts, "squares.html", plot.WithTitle("y = x^2"))
```

---

## Upgrading from v3

Change imports from `github.com/t1sk3/eulerlib/v3/...` to `github.com/t1sk3/eulerlib/v4/...`, then use the tables below. If your `go.mod` still requires `github.com/t1sk3/go-utils/types`, run `go get github.com/t1sk3/go-utils/types@none`.

| v3 package | v4 |
|---|---|
| `prime_numbers` | `primes` |
| `num_theory` | `numtheory`, `combinatorics`, `digits` |
| `etc` | `digits`, `numtheory`, go-utils `stream`, the standard library (below) |
| `vector` | go-utils `types/vector` + `geometry` |
| `utils` | removed: use type parameters with `constraints.*` |

### `etc`

| v3 | v4 |
|---|---|
| `Min(a, b, …)`, `Max(…)` | builtin `min`, `max`; `slices.Min`/`slices.Max` for slices |
| `Sum(s)` | `stream.Sum(stream.FromSlice(s))`; for complex numbers `stream.FromSlice(s).Fold(0, func(a, v E) E { return a + v })` |
| `Product(s)` | `stream.FromSlice(s).Collect(numtheory.Product[E]())` |
| `Filter(s, f)`, `Map(s, f)` | `stream.FromSlice(s).Filter(f).ToSlice()`, `stream.FromSlice(s).Map(f).ToSlice()` |
| `Sort(s, less)` | `slices.SortFunc(s, cmp)` (the comparator returns an `int`) |
| `Unique(s)`, `RemoveDuplicates(s)` | `stream.Distinct(stream.FromSlice(s)).ToSlice()` (keeps first occurrences) |
| `RemoveDuplicatesFunc`, `RemoveDuplicates2`, `RemoveDuplicateSlices` | `stream.FromSlice(s).DistinctBy(key).ToSlice()` |
| `UniqueCount(s)` | `stream.FromSlice(s).CountBy(stream.Identity[E])` |
| `MostFrequent(s)` | `counts := stream.FromSlice(s).CountBy(stream.Identity[E])`, then `stream.FromMap(counts).Max(func(a, b pair.Pair[E, int]) int { return a.Second - b.Second })` (ties are resolved arbitrarily) |
| `SliceContains(s, v)`, `SliceContainsAny(s, v, f)` | `slices.Contains(s, v)`, `slices.ContainsFunc(s, func(x E) bool { return f(x, v) })` |
| `CountOccurrenceInSlice(s, v)` | `stream.FromSlice(s).Filter(stream.Equal(v)).Count()` |
| `RemoveFromSlice(s, i)` | `slices.Delete(s, i, i+1)` |
| `GenerateSlice(n, v)` | `slices.Repeat([]E{v}, n)` |
| `Range(a, b)`, `RangeStep(a, b, step)` | `stream.Range(a, b).ToSlice()`; `stream.IterateWhile(a, func(x E) bool { return x < b }, func(x E) E { return x + step }).ToSlice()` |
| `MakeSliceOf(n, f)` | `stream.Range(0, n).Map(f).ToSlice()` |
| `FirstMatching(start, step, p)` | `stream.Iterate(start, func(x E) E { return x + step }).Filter(p).FindFirst()` |
| `JoinSliceString(s)`, `JoinSlice(s)` | `strings.Join(s, "")`, `digits.ConcatString(s...)` (or `digits.Concat` for a number) |
| `ReverseString(s)` | `r := []rune(s); slices.Reverse(r); string(r)`; for numbers `digits.Reverse(n)` |
| `IsPalindrome(s)` | `digits.IsPalindromeString(s)`; numbers: `digits.IsPalindrome(n)`, `digits.IsPalindromeInBase(n, b)` |
| `MakeIntSlice(n)` | `digits.Of(n)` |
| `IsPandigital(n)` | `digits.IsPandigital(n)`. **Stricter:** digits 1–9 exactly once (v3 only checked that each was present) |
| `IsPandigitalInBase(n, b)`, `IsPandigitalInBaseString(s, b)` | `digits.IsPandigitalInBase(n, b)`, `digits.IsPandigitalString(s, b)` |
| `DecimalToBase(n, b)` | `digits.ToBase(n, b)`; panics on a base outside 2..62 instead of returning `""` |
| `Totient`, `Phi`, `ListTotients`, `ListMobius` | `numtheory.Totient`, `numtheory.TotientSieve`, `numtheory.MobiusSieve` (`Phi` removed) |
| `Abs`, `SumOfSquares`, `SquareOfSum` | `numtheory.Abs`, `numtheory.SumOfSquares`, `numtheory.SquareOfSum` |
| `CreateFile`, `CreateFileWithContent`, `ReadFile`, `FileExists` | `os.Create`, `os.WriteFile`, `os.ReadFile`, `os.Stat` |

### `num_theory`

| v3 | v4 |
|---|---|
| `Reduce(s, f)`, `ReduceWithInit(init, s, f)` | `stream.FromSlice(s).Reduce(f)`, `stream.FromSlice(s).Fold(init, f)` |
| `MaxInSlice`, `Lcd` | `slices.Max`, `numtheory.Lcm` |
| `Factorize(n) map[E]E` | `numtheory.Factorize(n) map[E]int`, or `numtheory.Factors(n)` (`iter.Seq2`) |
| `FactorizeSPF(n, spf)` | `numtheory.FactorizeSPF(n, spf) map[E]int`, or `numtheory.FactorsSPF` |
| `FactorizeBigInt`, `PrimeFactorsBigInt` | `numtheory.FactorizeBig(n) []pair.Pair[*big.Int, int]` |
| `PowMod(x, y, p) int64` | `numtheory.PowMod(x, y, m) E`, exact for any modulus |
| `PowBigInt` | `numtheory.PowBig` |
| `Pow(b, n E)` | `numtheory.Pow(b E, n int)` |
| `Permutations(arr) [][]E` | `combinatorics.Permutations(arr)` (stream, doesn't modify `arr`) |
| `Combinations(set, n)` | `combinatorics.Combinations(set, k)`; all sizes: `combinatorics.Subsets(set)` |
| `PermutationCount`, `Factorial`, `Binomial`, `Partition` | `combinatorics.*`, same names |
| `FactorialBigInt` | `combinatorics.FactorialBig` |
| `FactorialDigitSum` | `combinatorics.DigitFactorialSum` |
| `DigitSum`, `DigitSumString`, `DigitCount` | `digits.Sum`, `digits.SumString`, `digits.Count` (returns `int`) |
| `AmicablePairsBelow(limit) []Pair` | `numtheory.AmicablePairs(limit)` (stream; both members now < limit) |
| `Divisors`, `CountDivisors`, `SumDivisors`, `Gcd`, `Lcm`, … | `numtheory.*`. `Divisors` is now sorted, `CountDivisors` returns `int`, `Gcd()` returns 0 instead of panicking |

### `prime_numbers`

| v3 | v4 |
|---|---|
| `IsPrime`, `NextPrime` | `primes.IsPrime`, `primes.NextPrime` |
| `ListPrimality(n)` | `primes.NewSieve(n)` (`Has`, `All`, `Stream`, `Count`) |
| `ListPrimes(n)` | `primes.UpTo(n)` |
| `ListSmallestPrimeFactors(n)` | `primes.SmallestFactorSieve(n)` |
| `FirstNPrimes(n)` | `primes.First[E](n)` |
| `Primes()` | `primes.All[E]()` |
| `PrimeCount(s, e)`, `SumPrimes(s, e)` | `primes.Count(s, e)`, `primes.Sum(s, e)` |
| `PrimeGenerator(limit)` | `primes.All[E]().TakeWhile(stream.AtMost(limit)).All()` |
| `PrimeNumberIterator`, `NewPrimeNumberIterator` | `primes.From(n)`, or `iter.Pull(primes.All[E]().All())` if you need a `next()` function |

### `fibonacci`, `figurate`, `sequences`, `pythagoras`

| v3 | v4 |
|---|---|
| `GenFibo`, `GenFiboBig` (channels) | `fibonacci.All[E]()`, `fibonacci.AllBig()` |
| `Fibonacci(limit)`, `FibonacciBig(limit)` | `fibonacci.First[E](n)`; note v4 starts at F(0) = 0 |
| `FibonacciSingle(n) int64`, `FibonacciSingleBig(n)` | `fibonacci.Nth[E](n) E`, `fibonacci.NthBig(n)` |
| `figurate.IsPentagonal(0)`, `IsHexagonal(0)` | now `true`, like `IsTriangular(0)` |
| `sequences.Collatz(n) []E` | `sequences.Collatz(n).ToSlice()` (stream; panics instead of overflowing E) |
| `pythagoras.IsTriplet` | unchanged; new `Triplet` type and triplet streams |

### `vector`

| v3 | v4 |
|---|---|
| `vector.NewVector([]E)` / `*Vector[E]` | `vector.Of(xs...)` → `vector.VecN[E]` (a slice: `v[i]`, `len(v)`, `range v`) |
| `.Len()`, `.At(i)`, `.Set(i, e)`, `.Elements()`, `.Swap(i, j)` | `len(v)`, `v[i]`, `v[i] = e`, `v.Clone()`, `v[i], v[j] = v[j], v[i]` |
| `NewVector2(x, y)`, `NewVector3(x, y, z)` | `vector.Of2(x, y)`, `vector.Of3(x, y, z)` (values, with `X`, `Y`, `Z` fields) |
| `.X()`, `.SetX(x)` | `.X`, `v.X = x` |
| `AsVector2(v)`, `AsVector3(v)` | `v.ToVec2()`, `v.ToVec3()` |
| `AreColinear`, `Quadrant`, `SameQuadrant` | `geometry.*`, taking `vector.Vec2` values |
| `HexCoordinate`, `NewHexCoordinate(q, r)` | `geometry.Hex[E]{Q, R}`, `geometry.HexOf(q, r)` |
| `HexAreColinear`, `HexSameQuadrant`, `CartesianToHex`, `HexToCartesian`, `RoundHex` | `geometry.*`, same names |

go-utils vectors use `constraints.Number`, so complex-valued vectors are no longer supported.

### `progress`

| v3 | v4 |
|---|---|
| `ProgressBar(slice)` | `progress.Slice(slice)` |
| `ProgressBarN(n)` | `progress.N(n)` |
| `Spinner()` | unchanged |

---

## License

[MIT](LICENSE)
