package main

import "fmt"

func main() {
	var N int
	fmt.Scan(&N)
	var A [N]int
	for i := 0; i < N; i++ {
		fmt.Scan(&A[i])
	}
	
}
