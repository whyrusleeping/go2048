package main

import "fmt"

type Solver func(*Board, UtilityFunc) (bool, int)

func LookaheadSolver(b *Board, utility UtilityFunc) (bool, int) {
	opts := make([]int, 4)
	turn := 0
	for !b.CheckWin() {
		turn++
		for i := 0; i < 4; i++ {
			for j := 0; j < 16; j++ {
				nb := b.Copy()
				c := nb.Round(i)
				if !c || nb.CheckLoss() {
				} else {
					/*
					//Initial attempt, score based heuristic
					subo := make([]int, 4)
					for j := 0; j < 4; j++ {
						for try := 0; try < 16; try++ {
							onb := b.Copy()
							w := onb.Round(j)
							if !w {
								subo[j] = -1
							} else {
								subo[j] += utility(onb)
							}
						}
					}
					opts[i] = Aver(subo)
					*/
					opts[i] += utility(nb)
				}
			}
		}
		var act bool
		var i int
		fmt.Println(opts)
		for i = 0; i < 4; i++ {
			choice := MaxI(opts)
			act = b.Round(choice)
			if !act {
				opts[choice] = 0
			} else {
				break
			}
		}
		if i == 4 {
			fmt.Println("Something is very wrong...")
			b.PrintBoard()
			for {}
		}
		if !act {
			fmt.Println("Chose bad move...")
			fmt.Println(opts)
		}

		if b.CheckLoss() {
			return false, b.score
		}
		if turn % 1000 == 0 {
			fmt.Println(turn)
			b.PrintBoard()
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
