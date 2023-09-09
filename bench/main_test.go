package main

import "testing"

func BenchmarkFibo(b *testing.B) {
	fibo(b.N)
}
