// Package memo provides generic memoization wrappers, handy for the
// recursive counting/DP-style brute forces that show up throughout Project
// Euler (and slow to a crawl without a cache).
package memo

// Memoize wraps fn so repeated calls with the same argument are computed
// only once. It is not safe for concurrent use from multiple goroutines.
//
// For a recursive function, declare the function variable first so the
// memoized wrapper can call itself:
//
//	var fib func(n int) int64
//	fib = memo.Memoize(func(n int) int64 {
//	    if n < 2 {
//	        return int64(n)
//	    }
//	    return fib(n-1) + fib(n-2)
//	})
func Memoize[K comparable, V any](fn func(K) V) func(K) V {
	cache := make(map[K]V)
	return func(k K) V {
		if v, ok := cache[k]; ok {
			return v
		}
		v := fn(k)
		cache[k] = v
		return v
	}
}

// key2 combines two comparable values into a single comparable map key.
type key2[K1, K2 comparable] struct {
	a K1
	b K2
}

// Memoize2 is Memoize for a two-argument function, useful for the
// (position, state) style recurrences common in digit-DP and counting
// problems:
//
//	var ways func(pos, remaining int) int64
//	ways = memo.Memoize2(func(pos, remaining int) int64 {
//	    ...
//	    return ways(pos+1, remaining-k)
//	})
func Memoize2[K1, K2 comparable, V any](fn func(K1, K2) V) func(K1, K2) V {
	cache := make(map[key2[K1, K2]]V)
	return func(a K1, b K2) V {
		k := key2[K1, K2]{a, b}
		if v, ok := cache[k]; ok {
			return v
		}
		v := fn(a, b)
		cache[k] = v
		return v
	}
}
