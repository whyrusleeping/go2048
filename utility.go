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

func Utility_Corner(b *Board) int {
	val := 0
	max := 0
	for i,r := range b.tiles {
		for j,v := range r {
			if v != 0 {
				if v > max {
					max = v
				}
				posw := (j+1) + ((i + 1) * 3)
				val += (posw * v)
			}
		}
	}
	if b.tiles[b.size-1][b.size-1] != max {
		val -= max * 3
	}
	return val
}
