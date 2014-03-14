package main

import (
	"fmt"
	"runtime"
	"math/rand"
	"math"
	"time"

	"os"
	"runtime/pprof"
)

func MaxI(l []int) int {
	var cm, ci int
	cm = l[0]
	ci = 0
	same := true
	for _,v := range l {
		if v != cm {
			same = false
			break
		}
	}
	if same {
		return rand.Intn(len(l))
	}
	for i := 1; i < len(l); i++ {
		if l[i] > cm {
			ci = i
			cm = l[i]
		}
	}
	return ci
}

func AverPos(l []int) int {
	count := len(l)
	sum := 0
	for _,v := range l {
		if v >= 0 {
			sum += v
		} else {
			count--
		}
	}
	return sum / count
}


//Play 'n' rounds of the game using the given solver
//Returns #wins, best score, worst score, and average score
func PlayN(n int, slvr Solver, utility UtilityFunc) (int,int,int,int) {
	best := 0
	worst := math.MaxInt32
	sum := 0
	wins := 0
	b := NewBoard(4)
	for i := 0; i < n; i++ {
		b.PlaceRandom()
		b.PlaceRandom()
		w,s := slvr(b, utility)
		if w {
			wins++
		}
		if s > best {
			best = s
		}
		if s < worst {
			worst = s
		}
		sum += s
		b.Reset()
	}
	return wins,best,worst,sum/n
}

//Average an array of numbers
func Aver(l []int) int {
	sum := 0
	for _,v := range l {
		sum += v
	}
	return sum/len(l)
}

//Run a given solver algorithm for a set number of trials
func RunTrials(slvr Solver, utility UtilityFunc) {
	done := make(chan bool)
	num_th := 1
	run_per := 20
	wins := make([]int,num_th)
	bests := make([]int, num_th)
	worsts := make([]int, num_th)
	avs := make([]int, num_th)
	for i := 0; i < num_th; i++ {
		go func(n int) {
			wins[n],bests[n],worsts[n],avs[n] = PlayN(run_per, slvr, utility)
			done <- true
		}(i)
	}
	for i := 0; i < num_th; i++ {
		<-done
	}

	fmt.Printf("Wins: %d\n", Aver(wins))
	fmt.Printf("Best: %d\n", Aver(bests))
	fmt.Printf("Worst: %d\n", Aver(worsts))
	fmt.Printf("Average: %d\n", Aver(avs))
}

func main() {
	runtime.GOMAXPROCS(2)
	rand.Seed(time.Now().UnixNano())
	fi,err := os.Create("prof.out")
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	pprof.StartCPUProfile(fi)
	fmt.Println("Lookahead Solver")
	RunTrials(LookaheadSolver, Utility_Corner)
	pprof.StopCPUProfile()

	fmt.Println("Utility Based Solver")
	RunTrials(BestMoveSolver, Utility_Score)
	fmt.Println("LDRD Solver")
	RunTrials(LDRDSolver, nil)
}

//If you just want to play the game...
func plmain() {
	rand.Seed(time.Now().UnixNano())
	b := NewBoard(4)
	b.PlaceRandom()
	b.PlaceRandom()
	b.PrintBoard()
	var s string
	for !b.CheckWin() {
		fmt.Scanf("%s", &s)
		if len(s) == 0 {
			continue
		}
		switch s[0] {
			case 'w':
				b.Round(0)
			case 'a':
				b.Round(3)
			case 's':
				b.Round(2)
			case 'd':
				b.Round(1)
			default:
				fmt.Println("Invalid Key")
		}
		fmt.Printf("Score: %d\n", b.score)
		b.PrintBoard()
	}
}
