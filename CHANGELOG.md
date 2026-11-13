# Changelog

All notable changes to eulerlib are documented here. This project follows [Semantic Versioning](https://semver.org/); because Go modules encode breaking changes in the import path, a major bump also means the module path changes (e.g. `.../v2` → `.../v3`).

## Unreleased

Fixes from a review of v4.0.0. No API changes.

- `plot.Scatter2D`/`Scatter3D` embedded titles, axis captions, colors and point labels in the page's `<script>` without escaping, so a string containing `</script>` could inject markup. They are now escaped. The package documentation no longer calls the output self-contained: it loads ECharts from go-echarts' asset host.
- `plot.SaveCSV2D`/`SaveCSV3D` ignored an error from closing the file.
- `pythagoras.PrimitiveTriplets` and `TripletsWithPerimeterAtMost` overflowed their bound checks for narrow element types and yielded invalid triplets (for `int32`, from a bound of about 2^30 and 2^29 respectively).
- `numtheory.SumOfSquares` overflowed an intermediate product, so it was wrong for results that fit (for `int64`, from n ≈ 1.66 million). It now returns 0 for n <= 0.
- `numtheory.IsAbundant`, `IsDeficient` and `IsPerfect` were wrong when the divisor sum did not fit in the element type, and `IsAbundant` was true for negative numbers.
- `numtheory.CRT` returned a wrapped modulus when lcm(moduli) did not fit in the element type; it now returns an empty Optional.
- `combinatorics.Binomial` returned 0 for a `uint64` n above `MaxInt64`.
- `numtheory.Gcd` documents its one negative result (a gcd of 2^(bits-1) in a signed type).

## v4.0.0

The migration table, with code, is in [readme.md](readme.md#upgrading-from-v3).

### Breaking

- The module path is now `github.com/t1sk3/eulerlib/v4`.
- Packages renamed to Go style: `prime_numbers` → `primes`, `num_theory` → `numtheory` (combinatorics moved to the new `combinatorics` package, digit helpers to `digits`).
- Removed `etc` (its Euler-specific parts moved to `digits` and `numtheory`; the rest is covered by go-utils `stream`, `slices` and builtins), `utils` (obsolete reflection checks) and `vector` (use go-utils `types/vector`; coordinate helpers moved to `geometry`). Complex-valued vectors are no longer supported.
- Lazy and infinite sequences are `stream.Stream[E]`; key/value results are `iter.Seq2`. Removed every channel-returning function (`PrimeGenerator`, `GenFibo`, `GenFiboBig`) and `PrimeNumberIterator`.
- Factorizations use `int` exponents: `Factorize` returns `map[E]int`; `FactorizeBigInt` and `PrimeFactorsBigInt` became one `FactorizeBig` returning `[]pair.Pair[*big.Int, int]`.
- `PowMod` returns `E` instead of `int64`; `Pow` takes an `int` exponent; `CountDivisors` and `digits.Count` return `int`.
- `Permutations`/`Combinations` are lazy streams that never modify their input; `Combinations(set, 0)` no longer means "all sizes" (use `Subsets`).
- Fibonacci numbers start at F(0) = 0 everywhere (`GenFibo` started at 1, 2, 3).
- `progress.ProgressBar` → `progress.Slice`, `progress.ProgressBarN` → `progress.N`.
- `vector.HexCoordinate` (pointer, embedding `Vector2`) → `geometry.Hex[E]{Q, R}`, a comparable value type.
- Behaviour changes: `digits.IsPandigital` requires digits 1–9 **exactly once** (v3's `etc.IsPandigital` only checked presence); `figurate.IsPentagonal(0)` and `IsHexagonal(0)` are now true, like `IsTriangular(0)`; `AmicablePairs(limit)` requires both members below `limit`; `Gcd()` returns 0 and `Lcm()` returns 1 instead of panicking; `digits.ToBase` panics on a base outside 2..62 instead of returning `""`; `sequences.Collatz` panics instead of silently overflowing E; `FloatIsInteger(±Inf)` is false.

### Added

- `primes`: `PrevPrime`, `Sieve` (bitset, with `Has`, `IsPrime`, `Count`, `All`, `Stream`), `All`/`From`/`Between` streams backed by a segmented sieve, `UpTo`, `First`, `Count`, `Sum`.
- `numtheory`: `Factors` (`iter.Seq2`), `FactorsSPF`, `DistinctPrimeFactors`, `SigmaSieve`, `MulMod`, `Isqrt`, predicate builders `DivisibleBy`, `DivisibleByAny`, `CoprimeTo`, collectors `Product`, `GcdOf`, `LcmOf`.
- `combinatorics`: `LexPermutations`, `NthPermutation`, `Subsets`, `CartesianProduct`, `PartitionTable`.
- `digits`: `OfBase`, `FromDigits`, `Reverse`, `IsPalindromeInBase`, `IsPandigitalN`, `Concat`, `ConcatString`, `SumOf` collector.
- `fibonacci.IndexOfFirstWithDigits`; `figurate.NthSquare`, `NthPolygonal` and streams for every kind; `sequences.CollatzLengths`; `pythagoras.Triplet`, `PrimitiveTriplets`, `TripletsWithPerimeterAtMost`; `geometry.HexRing`, `Hex.Distance`, `Hex.Neighbors`; `memo.MemoizeConcurrent`; `progress.Seq`, `progress.Stream`; `plot.PointFromVec2/Vec3/Pair`, `Points2D`, `Points3D`.
- Runnable `Example` tests in every package, and `example_euler_test.go` solving Euler #1, #5, #7, #10, #14, #21, #24 and #45, each checked against a brute-force solution.

### Changed

- `IsPrime` is a deterministic Miller–Rabin test, exact for all 64-bit values.
- `Factors`, `Factorize`, `Divisors`, `CountDivisors`, `SumDivisors` and `Totient` factorize with trial division plus Pollard's rho, so they are fast for every 64-bit input.
- `Divisors` is sorted ascending.
- `NthBig` (was `FibonacciSingleBig`) uses fast doubling.
- `LongestCollatz` uses a memoised length table.
- `CRT` combines congruences with overflow-free modular multiplication.

### Fixed

- `IsPrime`, `IsSquare` and `figurate.Is*` used `math.Sqrt` on float64 (exact only up to 2^53), and `figurate.Is*` overflowed computing `8n+1`/`24n+1`; all are now exact for every value of their type.
- `PowMod` required `p < 2^31`; it is now exact for any modulus.
- `Combinations` used a bitmask and broke for 64 or more elements.
- `Permutations` modified the caller's slice.
- `PrimeCount` was O(n√n) and `Primes`/`NextPrime` used trial division per candidate; the segmented sieve makes the first million primes ~260× faster (5.1 s → 19.5 ms).
- The channel generators leaked a goroutine when the caller stopped early.
- `etc.DecimalToBase` was wrong for `MinInt64` and for `uint64` values above `MaxInt64`; `digits.ToBase` is exact.
- `etc.IsPandigitalInBaseString` could not read letter digits (bases above 10).
- `CRT` could overflow for moduli whose product is close to `MaxInt64`.
- `num_theory.AmicablePairsBelow` could return pairs whose larger member was not below the limit.

## v3.1.0

No exported signatures change in this release.

### Changed

- The minimum Go version is now **1.27** (`go 1.27.0` in `go.mod`). Go 1.24 no longer receives security fixes, and CI already only ran on 1.27, so the old minimum was never tested.
- go-utils is upgraded from v0.1.0 to **v0.6.1**. The separate `github.com/t1sk3/go-utils/types` module is no longer required: since go-utils v0.5.1 those packages live in the root module, and requiring both caused "ambiguous import" errors for projects that also depend on a recent go-utils. If your own `go.mod` still lists it, run `go get github.com/t1sk3/go-utils/types@none`.
- go-utils v0.6.1 removed its package-level `stream.Map`/`stream.Fold`/... functions in favour of generic methods, so eulerlib's internal calls now use the method form (`s.Map(f)`, `s.Fold(init, f)`). No behaviour change.
- Dependency refresh: `go-echarts` v2.7.3, `schollz/progressbar` v3.19.1, `golang.org/x/sys` v0.50.0, `golang.org/x/term` v0.48.0.
- `etc.Sort` uses `slices.SortFunc` instead of `sort.Slice` (still not stable).
- Mechanical modernisations from `go fix` (`any`, `reflect.TypeFor`, `range n`, `strings.Builder`).
- CI: `actions/setup-go@v5`, Go version taken from `go.mod` plus a `stable` job; checks `gofmt`, `go mod tidy -diff`, `go vet`, `go fix -diff`, and runs tests with `-race`.

### Fixed

- `prime_numbers.Primes()` was not reusable: the stream shared one iterator, so consuming it a second time continued where the first left off (`[2 3 5]`, then `[7 11 13]`). Each consumption now starts again at 2.
- `num_theory.Factorize` looped up to `n` instead of `sqrt(n)`, so factorizing a large prime took O(n) steps (seconds for 10^9+7). It now stops at `sqrt(n)` (written as `i <= n/i` so it can't overflow) and is cross-checked against the old implementation for every n in [0, 5000].
- `progress.ProgressBar`, `ProgressBarN` and `Spinner` created (and drew) their bar when the function was called rather than when iteration started, and ranging over the returned sequence twice reused an already finished bar. The bar is now created inside the sequence, so nothing is drawn before the loop starts and every loop gets a fresh bar.

## v3.0.0

### Breaking

- The module path changed from `github.com/t1sk3/eulerlib/v2` to `github.com/t1sk3/eulerlib/v3`. Update every import accordingly.
- `utils.Integer`, `utils.RealNumber`, `utils.Comparable`, `utils.SignedInteger`, `utils.UnsignedInteger`, `utils.Float`, `utils.Complex`, `utils.Number`, and `utils.SignedNumber` were removed. Every exported function that was constrained by one of them now uses the equivalent type from [`github.com/t1sk3/go-utils/types/constraints`](https://github.com/t1sk3/go-utils) — see the migration table in [readme.md](readme.md#type-constraints-and-upgrading-from-v2). Calling eulerlib's functions with concrete types (`int`, `int64`, `float64`, ...) is unaffected; only code that referenced eulerlib's constraint types directly needs updating.
- `num_theory.ModInverse(a, m)` now returns `optional.Optional[E]` instead of `(E, bool)`.
- `num_theory.CRT(remainders, moduli)` now returns `optional.Optional[pair.Pair[E, E]]` (the pair is `(x, M)`) instead of `(E, E, bool)`.
- `utils`'s reflection-based type checks (`SameType`, `IsInteger`, ...) are unchanged and remain in the `utils` package.

### Added

- `prime_numbers.Primes[E]()` — an infinite, lazily-computed `stream.Stream[E]` of primes, composable with `.Limit()`, `.TakeWhile()`, `.Filter()`, etc. without an upper bound up front.
- `etc.FirstMatching(start, step, predicate)` — generic "find the first n such that..." search, returning `optional.Optional[E]`.
- `sequences.LongestCollatz(limit)` — the starting number under `limit` with the longest Collatz chain (Project Euler #14), as `pair.Pair[E, int]`.
- `num_theory.AmicablePairsBelow(limit)` — every amicable pair under `limit` (Project Euler #21), as `[]pair.Pair[E, E]`.
- `etc.MostFrequent(slice)` — the most common element and its count, as `optional.Optional[pair.Pair[E, int]]`.

### Changed

- eulerlib now depends on [`github.com/t1sk3/go-utils`](https://github.com/t1sk3/go-utils) for generic constraints, `Optional`, `Pair`, `Set`, and `Stream`.
- `go.mod`'s `go` directive moved from `1.23` to `1.24.8` (required by go-utils).
- `memo.Memoize2` uses `pair.Pair` as its cache key instead of a private struct.
- `etc.Unique`/`etc.RemoveDuplicates` use `set.Set` for membership tracking; `etc.Sum` delegates to `stream.Fold`.
- `prime_numbers.FirstNPrimes` is implemented with `stream.Generate(...).Limit(n).ToSlice()`.

## v2.1.1 and earlier

Released before this changelog was introduced — see the [git history](https://github.com/t1sk3/eulerlib/commits/main) and [tags](https://github.com/t1sk3/eulerlib/tags).
