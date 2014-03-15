package main

import (
	"fmt"
	"runtime"
	"math/rand"
	"math"
	"time"

	"os"
	"runtime/pprof"

	"github.com/nsf/termbox-go"
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
	num_th := 2
	run_per := 5
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

func aimain() {
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

func draw_all() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.Flush()
}

func print_wide(x, y int, s string) {
	c := termbox.ColorDefault
	for _, r := range s {
		termbox.SetCell(x, y, r, termbox.ColorDefault, c)
		x++
	}
}

func TboxPrintBoard(b *Board) {
	s := fmt.Sprintf("Score: %d", b.score)
	print_wide(0,0,s)
	for i, r := range b.tiles {
		for j, v := range r {
			s = fmt.Sprintf("|%3d", v)
			print_wide(j*4,i+1,s)
		}
		print_wide(16,i,"|")
	}
}
//If you just want to play the game...
func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	rand.Seed(time.Now().UnixNano())
	b := NewBoard(4)
	b.PlaceRandom()
	b.PlaceRandom()
	draw_all()
	TboxPrintBoard(b)
	for !b.CheckWin() {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowUp:
				b.Round(0)
				TboxPrintBoard(b)
			case termbox.KeyArrowRight:
				b.Round(1)
				TboxPrintBoard(b)
			case termbox.KeyArrowDown:
				b.Round(2)
				TboxPrintBoard(b)
			case termbox.KeyArrowLeft:
				b.Round(3)
				TboxPrintBoard(b)
			case termbox.KeyEsc:
				return
			}
		case termbox.EventResize:
			draw_all()
		}
		if b.CheckLoss() {
			print_wide(0,6,"You Lose!")
			return
		}
		termbox.Flush()
	}
	print_wide(0,6,"You win!!")
}
