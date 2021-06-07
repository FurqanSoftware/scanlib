package main

import "fmt"

func main() {
	var N int
	fmt.Scan(&N)
	var S [N]string
	for i := 0; i < N; i++ {
		fmt.Scan(&S[i])
	}
	
}
