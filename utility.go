package main

type UtilityFunc func(*Board) int

func Utility_OpenHeavy(b *Board) int {
	//return (b.OpenCount() * 10) + b.Sum()
	oc := b.OpenCount()
	return (oc * 2) + (b.score * 2)
}

func Utility_Score(b *Board) int {
	return b.score
}

