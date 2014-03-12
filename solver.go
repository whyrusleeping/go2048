package main

import "fmt"

type Solver func(*Board, UtilityFunc) (bool, int)

func LookaheadSolver(b *Board, utility UtilityFunc) (bool, int) {
	for !b.CheckWin() {
		opts := make([]int, 4)
		for i := 0; i < 4; i++ {
			for j := 0; j < 16; j++ {
				nb := b.Copy()
				c := nb.Round(i)
				if !c || nb.CheckLoss() {
					break
				} else {
					//Initial attempt, score based heuristic
					for k := 0; k < 4; k++ {
						snb := nb.Copy()
						mov := snb.Round(k)
						if mov {
							opts[i] += utility(snb)
						}
					}
				}
			}
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
