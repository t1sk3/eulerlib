package eulerlib

import (
	"math"
	"math/big"
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

func FibonacciSingle[E Integer](n E) int64 {
	if n < 2 {
		return int64(n)
	}
	nf := float64(n)
	return int64((math.Pow((1+math.Sqrt(5))/2, nf) - math.Pow((1-math.Sqrt(5))/2, nf)) / math.Sqrt(5))
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

func FibonacciSingleBig(n int64) *big.Int {
	if n < 2 {
		return big.NewInt(n)
	}
	temp1 := big.NewFloat(1 + math.Sqrt(5))
	temp2 := big.NewFloat(1 - math.Sqrt(5))
	temp3 := powBigFloat(big.NewFloat(2), uint64(n))
	temp3.Mul(temp3, big.NewFloat(math.Sqrt(5)))

	temp1 = powBigFloat(temp1, uint64(n))
	temp2 = powBigFloat(temp2, uint64(n))
	temp1.Sub(temp1, temp2)
	temp1.Quo(temp1, temp3)

	res, _ := temp1.Int(nil)
	return res
}

func powBigFloat(a *big.Float, e uint64) *big.Float {
	result := zero().Copy(a)
	for i := uint64(0); i < e-1; i++ {
		result.Mul(result, a)
	}
	return result
}

func zero() *big.Float {
	r := big.NewFloat(0.0)
	r.SetPrec(256)
	return r
}
