package main

import (
	"fmt"
	"io"
)

func main() {
	var cmd string
	for {
		if _, err := fmt.Scan(&cmd); err == io.EOF {
			break
		}
		var param int
		if cmd=="PUSH"||cmd=="REPEAT" {
			fmt.Scan(&param)
		}
		if cmd=="REPEAT" {
		}
	}
	
}
