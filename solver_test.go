package main

import "testing"

func BenchmarkLookahead(b *testing.B) {
	PlayN(b.N, LookaheadSolver, Utility_Score)
}
