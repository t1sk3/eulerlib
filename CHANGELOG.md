# Changelog

All notable changes to eulerlib are documented here. This project follows [Semantic Versioning](https://semver.org/); because Go modules encode breaking changes in the import path, a major bump also means the module path changes (e.g. `.../v2` → `.../v3`).

## Unreleased (v3.1.0)

No exported signatures change in this release.

### Changed

- The minimum Go version is now **1.27** (`go 1.27.0` in `go.mod`). Go 1.24 no longer receives security fixes, and CI already only ran on 1.27, so the old minimum was never tested.
- go-utils is upgraded from v0.1.0 to **v0.6.1**. The separate `github.com/t1sk3/go-utils/types` module is no longer required: since go-utils v0.5.1 those packages live in the root module, and requiring both caused "ambiguous import" errors for projects that also depend on a recent go-utils. If your own `go.mod` still lists it, run `go get github.com/t1sk3/go-utils/types@none`.
- go-utils v0.6.1 removed its package-level `stream.Map`/`stream.Fold`/... functions in favour of generic methods, so eulerlib's internal calls now use the method form (`s.Map(f)`, `s.Fold(init, f)`). No behaviour change.
- Dependency refresh: `go-echarts` v2.7.3, `schollz/progressbar` v3.19.1, `golang.org/x/sys` v0.50.0, `golang.org/x/term` v0.48.0.
- `etc.Sort` uses `slices.SortFunc` instead of `sort.Slice` (still not stable).
- Mechanical modernisations from `go fix` (`any`, `reflect.TypeFor`, `range n`, `strings.Builder`).
- CI: `actions/setup-go@v5`, Go version taken from `go.mod` plus a `stable` job; checks `gofmt`, `go mod tidy -diff`, `go vet`, `go fix -diff`, and runs tests with `-race`.

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
