package main

import (
	"flag"
	"fmt"

	"github.com/mohae/joefriday/cpu/cpuinfo"
)

func fibo(n int) int {
	if n < 2 {
		return n
	}
	return fibo(n-2) + fibo(n-1)
}

func main() {
	t := flag.Bool("cpu", false, "")
	n := flag.Int("n", 1, "")

	if *t {
		info, _ := cpuinfo.Get()
		fmt.Printf("%+#v", info)
	}

	flag.Parse()

	fmt.Println(fibo(*n))
}
