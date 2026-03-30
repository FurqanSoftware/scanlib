package main

import "fmt"

func main() {
	var A, B int
	fmt.Scan(&A, &B)
	var lo int
	var hi int
	lo = min(A, B)
	hi = max(A, B)
	
}
