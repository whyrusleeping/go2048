package main

import "testing"

func BenchmarkLookaheadSolve(b *testing.B) {
	PlayN(b.N, LookaheadSolver, Utility_Score)
}
