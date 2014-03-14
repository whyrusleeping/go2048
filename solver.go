package main

import "fmt"

type Solver func(*Board, UtilityFunc) (bool, int)

func UtilityForMove(b *Board, utility UtilityFunc, direction, depth int) int {
	rval := 0
	count := 6
	nb := b.Copy()
	c := nb.Shift(direction)
	if !c {
		return 0
	}
	modnb := nb.Copy()
	for j := 0; j < count; j++ {
		nb.PlaceRandom()
		if nb.CheckLoss() {
			rval += 1
		} else {
			snb := nb.Copy()
			for k := 0; k < 4; k++ {
				if depth > 0 {
					rval += UtilityForMove(snb, utility, k, depth - 1)
				} else {
					mov := snb.Shift(k)
					if mov {
						rval += utility(snb)
						snb.SetTo(nb)
					}
				}
			}
		}
		nb.SetTo(modnb)
	}
	return rval / count
}

func LookaheadSolver(b *Board, utility UtilityFunc) (bool, int) {
	for !b.CheckWin() {
		opts := make([]int, 4)
		for i := 0; i < 4; i++ {
			opts[i] = UtilityForMove(b, utility, i, 1)
		}
		act := b.Round(MaxI(opts))
		if !act {
			fmt.Println("Chose bad move...")
			fmt.Println(opts)
			fmt.Println(MaxI(opts))
			b.PrintBoard()
		}

		if b.CheckLoss() {
			return false, b.score
		}
	}
	return true, b.score
}

//Left Down Right Down Left Down Right Down... (surprisingly good)
func LDRDSolver(b *Board, util UtilityFunc) (bool, int) {
	for !b.CheckWin() {
		a := b.Round(LEFT)
		//b.PrintBoard()
		h := b.Round(DOWN)
		//b.PrintBoard()
		c := b.Round(RIGHT)
		//b.PrintBoard()
		d := b.Round(DOWN)
		//b.PrintBoard()
		if !(a || h || c || d) {
			b.Round(UP)
			//b.PrintBoard()
		}
		if b.CheckLoss() {
			//fmt.Println("Computer Lost!")
			//fmt.Printf("Final Score: %d\n", b.score)
			return false,b.score
		}
	}
	//fmt.Println("Computer won!!")
	return true,b.score
}

//Looks for the best next move based on a heuristic
func BestMoveSolver(b *Board, util UtilityFunc) (bool, int) {
	opts := make([]int, 4)
	for !b.CheckWin() {
		for i := 0; i < 4; i++ {
			nb := b.Copy()
			c := nb.Round(i)
			if !c || nb.CheckLoss() {
				opts[i] = -1
			} else {
				//Initial attempt, score based heuristic
				//opts[i] = nb.score
				opts[i] = util(nb)
			}
		}
		b.Round(MaxI(opts))

		/*
		fmt.Println()
		b.PrintBoard()
		*/

		if b.CheckLoss() {
			//fmt.Printf("Computer Lost! Score: %d\n", b.score)
			return false, b.score
		}
	}
	//fmt.Println("Computer Won!")
	return true, b.score
}
