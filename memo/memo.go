// Package memo provides generic memoization wrappers, handy for the
// recursive counting/DP-style brute forces that show up throughout Project
// Euler (and slow to a crawl without a cache).
package memo

import (
	"github.com/t1sk3/go-utils/concurrent"
	"github.com/t1sk3/go-utils/types/pair"
)

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
	cache := make(map[pair.Pair[K1, K2]]V)
	return func(a K1, b K2) V {
		k := pair.Of(a, b)
		if v, ok := cache[k]; ok {
			return v
		}
		v := fn(a, b)
		cache[k] = v
		return v
	}
}

// MemoizeConcurrent is Memoize, but safe for use from several goroutines at
// once (for example inside a stream.Parallel stage). fn runs without any
// lock held, so a recursive fn is fine. Two goroutines asking for the same
// missing key at the same time may both compute it; the first result
// stored wins, and both callers get it.
func MemoizeConcurrent[K comparable, V any](fn func(K) V) func(K) V {
	// Not concurrent.Map.ComputeIfAbsent: that runs fn under the map's
	// write lock, so a recursive fn would deadlock.
	cache := new(concurrent.Map[K, V])
	return func(k K) V {
		if v, ok := cache.Get(k).Value(); ok {
			return v
		}
		v := fn(k)
		actual, _ := cache.PutIfAbsent(k, v)
		return actual
	}
}
