package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func getCPUInfo() []string {
	f, _ := os.Open("/proc/cpuinfo")
	s := bufio.NewScanner(f)
	info := make([]string, 128)
	for s.Scan() {
		info = append(info, s.Text())
	}
	return info
}

func fibo(n int) int {
	if n < 2 {
		return n
	}
	return fibo(n-2) + fibo(n-1)
}

func init() {
	info := getCPUInfo()
	fmt.Print(strings.Join(info, "\n"))
	os.Exit(0)
}

func main() {
	t := flag.Bool("cpu", false, "")
	n := flag.Int("n", 1, "")
	flag.Parse()

	if *t {
		info := getCPUInfo()
		fmt.Print(strings.Join(info, "\n"))
		os.Exit(0)
	}

	fmt.Println(fibo(*n))
}
