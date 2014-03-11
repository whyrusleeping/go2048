package main

import (
	"fmt"
	"runtime"
	"math/rand"
	"math"
	"time"
)

const (
	UP = iota
	RIGHT
	DOWN
	LEFT
)

type Column struct {
	I int
	grid [][]int
	rev bool
}

type Row struct {
	I int
	grid [][]int
	rev bool
}

type Board struct {
	tiles [][]int
	size int
	score int
}

func NewBoard(size int) *Board {
	b := new(Board)
	b.size = size
	b.tiles = make([][]int, size)
	for i,_ := range b.tiles {
		b.tiles[i] = make([]int, size)
	}
	return b
}

//Place the next new piece
func (b *Board) PlaceRandom() {
	num := 2
	if rand.Intn(100) >= 90 {
		num = 4
	}
	for {
		x := rand.Intn(b.size)
		y := rand.Intn(b.size)
		if b.tiles[x][y] == 0 {
			b.tiles[x][y] = num
			return
		}
	}
}

//An 'Array' used to generically talk about rows or columns in
//either direction
type Iter interface {
	At(int) int
	Set(int, int)
	Len() int
}

func Shift(it Iter) (bool, int) {
	change := false
	score := 0
	for i := 0; i < it.Len(); i++ {
		if it.At(i) == 0 {
			for j := i + 1; j < it.Len(); j++ {
				if it.At(j) != 0 {
					it.Set(i,it.At(j))
					it.Set(j,0)
					change = true
					break
				}
			}
		}
		if it.At(i) != 0 {
			for j := i + 1; j < it.Len(); j++ {
				if it.At(j) != 0 {
					if it.At(j) == it.At(i) {
						it.Set(i, it.At(i) * 2)
						score += it.At(i)
						it.Set(j, 0)
						change = true
					}
					break
				}
			}
		}
	}
	return change,score
}

func (b *Board) GetRow(i int, rev bool) *Row {
	r := new(Row)
	r.I = i
	r.grid = b.tiles
	r.rev = rev
	return r
}

func (r *Row) At(i int) int {
	if r.rev {
		return r.grid[len(r.grid[r.I]) - (1 + i)][r.I]
	} else {
		return r.grid[i][r.I]
	}
}

func (r *Row) Set(i, v int) {
	if r.rev {
		r.grid[len(r.grid[r.I]) - (1 + i)][r.I] = v
	} else {
		r.grid[i][r.I] = v
	}
}

func (r *Row) Len() int {
	return len(r.grid[r.I])
}

func (b *Board) GetColumn(i int, rev bool) *Column {
	c := new(Column)
	c.I = i
	c.grid = b.tiles
	c.rev = rev
	return c
}

func (c *Column) At(i int) int {
	if c.rev {
		return c.grid[c.I][len(c.grid[c.I]) - (1 + i)]
	} else {
		return c.grid[c.I][i]
	}
}

func (c *Column) Set(i, v int) {
	if c.rev {
		c.grid[c.I][len(c.grid[c.I]) - (1 + i)] = v
	} else {
		c.grid[c.I][i] = v
	}
}

func (c *Column) Len() int {
	return len(c.grid[c.I])
}

func (b *Board) PrintBoard() {
	for _,v := range b.tiles {
		for _,t := range v {
			fmt.Print("|")
			if t == 0 {
				fmt.Print("___")
			} else {
				fmt.Printf("%3d", t)
			}
			fmt.Print("|")
		}
		fmt.Println()
	}
	fmt.Println()
}

func (b *Board) Left() bool {
	change := false
	for i := 0; i < b.size; i++ {
		v,s := Shift(b.GetColumn(i, false))
		b.score += s
		change = v || change
	}
	return change
}

func (b *Board) Right() bool {
	change := false
	for i := 0; i < b.size; i++ {
		v,s := Shift(b.GetColumn(i, true))
		b.score += s
		change = v || change
	}
	return change
}

func (b *Board) Up() bool {
	change := false
	for i := 0; i < b.size; i++ {
		v,s := Shift(b.GetRow(i, false))
		b.score += s
		change = v || change
	}
	return change
}

func (b *Board) Down() bool {
	change := false
	for i := 0; i < b.size; i++ {
		v,s := Shift(b.GetRow(i, true))
		b.score += s
		change = v || change
	}
	return change
}

//Produce a new copy of the board
func (b *Board) Copy() *Board {
	nb := NewBoard(b.size)
	nb.score = b.score
	for i,rs := range b.tiles {
		for j,v := range rs {
			nb.tiles[i][j] = v
		}
	}
	return nb
}

//Did the player win?
func (b *Board) CheckWin() bool {
	for _,r := range b.tiles {
		for _,v := range r {
			if v == 2048 {
				return true
			}
		}
	}
	return false
}

func (b *Board) Round(i int) bool {
	s := false
	switch i {
	case UP:
		if b.Up() {
			b.PlaceRandom()
			s = true
		}
	case RIGHT:
		if b.Right() {
			b.PlaceRandom()
			s = true
		}
	case DOWN:
		if b.Down() {
			b.PlaceRandom()
			s = true
		}
	case LEFT:
		if b.Left() {
			b.PlaceRandom()
			s = true
		}
	default:
		fmt.Println("Invalid Key")
	}
	return s
}

func (b *Board) CheckLoss() bool {
	for _,r := range b.tiles {
		for _,v := range r {
			if v == 0 {
				return false
			}
		}
	}
	l := b.Copy()
	if l.Left() {
		return false
	}
	l = b.Copy()
	if l.Right() {
		return false
	}
	l = b.Copy()
	if l.Up() {
		return false
	}
	l = b.Copy()
	if l.Down() {
		return false
	}
	return true
}

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

func (b *Board) OpenCount() int {
	out := 0
	for _,r := range b.tiles {
		for _,v := range r {
			if v == 0 {
				out++
			}
		}
	}
	return out
}

func (b *Board) Sum() int {
	out := 0
	for _,r := range b.tiles {
		for _,v := range r {
			out += v
		}
	}
	return out
}

func (b *Board) Utility() int {
	//return (b.OpenCount() * 10) + b.Sum()
	return (b.OpenCount() * 4) + (b.Sum() * 6) + b.score
}

type Solver func() (bool, int)

//Left Down Right Down Left Down Right Down... (surprisingly good)
func LDRDSolver() (bool, int) {
	b := NewBoard(4)
	b.PlaceRandom()
	b.PlaceRandom()
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
func BestMoveSolver() (bool, int) {
	b := NewBoard(4)
	b.PlaceRandom()
	b.PlaceRandom()
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
				opts[i] = nb.Utility()
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

//Play 'n' rounds of the game using the given solver
//Returns #wins, best score, worst score, and average score
func PlayN(n int, slvr Solver) (int,int,int,int) {
	best := 0
	worst := math.MaxInt32
	sum := 0
	wins := 0
	for i := 0; i < n; i++ {
		w,s := slvr()
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
func RunTrials(slvr Solver) {
	done := make(chan bool)
	wins := make([]int,4)
	bests := make([]int, 4)
	worsts := make([]int, 4)
	avs := make([]int, 4)
	for i := 0; i < 4; i++ {
		go func(n int) {
			wins[n],bests[n],worsts[n],avs[n] = PlayN(300, slvr)
			done <- true
		}(i)
	}
	for i := 0; i < 4; i++ {
		<-done
	}

	fmt.Printf("Wins: %d\n", Aver(wins))
	fmt.Printf("Best: %d\n", Aver(bests))
	fmt.Printf("Worst: %d\n", Aver(worsts))
	fmt.Printf("Average: %d\n", Aver(avs))
}

func main() {
	runtime.GOMAXPROCS(5)
	rand.Seed(time.Now().UnixNano())
	RunTrials(LDRDSolver)
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
