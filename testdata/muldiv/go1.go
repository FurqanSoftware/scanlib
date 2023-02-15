package main

import "fmt"

func main() {
	var A, B, C, D int
	fmt.Scan(&A, &B, &C, &D)
	var Af, Bf, Cf, Df float64
	fmt.Scan(&Af, &Bf, &Cf, &Df)
	var T string
	if A*B==0 {
		fmt.Scan(&T)
	}
	
}
