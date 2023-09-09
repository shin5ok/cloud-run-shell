package main

import (
	"flag"
	"fmt"
)

func fibo(n int) int {
	if n < 2 {
		return n
	}
	return fibo(n-2) + fibo(n-1)
}

func main() {
	n := flag.Int("n", 1, "")
	flag.Parse()

	fmt.Println(fibo(*n))
}
