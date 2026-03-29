package main

import "fmt"

func main() {
	var N int
	fmt.Scan(&N)
	var A [N]int
	var s int
	for i := 0; i < N; i++ {
		fmt.Scan(&A[i])
		s = s+A[i]
	}
	
}
