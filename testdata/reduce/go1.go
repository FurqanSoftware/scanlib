package main

import "fmt"

func main() {
	var N int
	fmt.Scan(&N)
	var s int
	for i := 0; i < N; i++ {
		var A int
		fmt.Scan(&A)
		s = s+A
	}
	
}
