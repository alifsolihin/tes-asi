package main

import (
	"fmt"
	"math"
)

func factorial(n int) int {
	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}
	return result
}

func powerOfTwo(n int) int {
	return int(math.Pow(2, float64(n)))
}

func f(n int) int {
	num := factorial(n)
	denom := powerOfTwo(n)
	return int(math.Ceil(float64(num) / float64(denom)))
}

func main() {
	for i := 0; i <= 10; i++ {
		fmt.Printf("f(%d) = %d\n", i, f(i))
	}
}
