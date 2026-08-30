package fibonacci

import (
	"math/big"

	"github.com/t1sk3/eulerlib/v2/utils"
)

// Generates the Fibonaccisequence
func GenFibo(limit int64) <-chan int64 {
	chnl := make(chan int64)
	fibos := []int64{1, 2, 3}
	go func() {
		for i := 0; int64(i) <= limit+3; i++ {
			fibos = append(fibos, fibos[len(fibos)-1]+fibos[len(fibos)-2])
			chnl <- fibos[i]
		}

		close(chnl)
	}()
	return chnl
}

// Returns a slice with the first x Fibonacci numbers
func Fibonacci(limit int64) (res []int64) {
	for x := range GenFibo(limit) {
		res = append(res, x)
	}
	return
}

// FibonacciSingle returns the nth fibonacci number.
// The result is the same as the last element of Fibonacci(n)
// but calculated while using less memory
func FibonacciSingle[E utils.Integer](n E) int64 {
	if n < 2 {
		return int64(n)
	}
	it := E(2)
	fibs := []E{1, 1}
	for it < n {
		tmp := fibs[0]
		fibs = fibs[1:]
		fibs = append(fibs, fibs[0]+tmp)
		it++
	}
	return int64(fibs[1])
}

// Generates the Fibonaccisequence in Big Integer
func GenFiboBig(limit int64) <-chan big.Int {
	chnl := make(chan big.Int)
	fibos := []big.Int{*big.NewInt(0), *big.NewInt(1)}
	go func() {
		for i := int64(0); i <= limit+3; i++ {
			fibos = append(fibos, *new(big.Int).Add(&fibos[len(fibos)-1], &fibos[len(fibos)-2]))
			chnl <- fibos[i]
		}

		close(chnl)
	}()
	return chnl
}

// Returns a slice with the first x Fibonacci numbers in Big Integer
func FibonacciBig(limit int64) (res []big.Int) {
	res = []big.Int{*big.NewInt(0), *big.NewInt(1)}
	for i := int64(0); i <= limit+3; i++ {
		res = append(res, *new(big.Int).Add(&res[len(res)-1], &res[len(res)-2]))
	}
	return
}

// FibonacciSingleBig returns the nth Fibonacci number as a *big.Int.
// F(0)=0, F(1)=1, F(2)=1, ... Uses an iterative approach for exact precision.
func FibonacciSingleBig(n int64) *big.Int {
	if n < 0 {
		panic("FibonacciSingleBig: n must be non-negative")
	}
	if n == 0 {
		return big.NewInt(0)
	}
	a := big.NewInt(0)
	b := big.NewInt(1)
	for i := int64(1); i < n; i++ {
		a.Add(a, b)
		a, b = b, a
	}
	return new(big.Int).Set(b)
}
