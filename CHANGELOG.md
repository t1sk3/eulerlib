# Changelog

All notable changes to eulerlib are documented here. This project follows [Semantic Versioning](https://semver.org/); because Go modules encode breaking changes in the import path, a major bump also means the module path changes (e.g. `.../v2` → `.../v3`).

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
