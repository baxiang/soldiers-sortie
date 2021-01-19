package main

import "fmt"

func main() {
	n := 15
	m := 12 * n
	var sum float64 = 0
	var percent float64 = 4
	var v float64 = 20000
	for i := 0; i < m; i++ {
		sum *= (percent/12 + 100) / 100
		sum += v
	}
	fmt.Printf("%.2f, %.2f, %d\n", sum, sum/(v*float64(m)), 300000*n)
}
