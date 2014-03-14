package main

import (
	"math/rand"
	"fmt"
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
	col []int
	rev bool
	size int
	alt int
}

type Row struct {
	I int
	grid [][]int
	rev bool
	size int
	alt int
}

type Board struct {
	tiles [][]int
	size int
	score int
	r *Row
	c *Column
	//rng *rand.Rand
}

func NewBoard(size int) *Board {
	b := new(Board)
	b.size = size
	b.tiles = make([][]int, size)
	b.r = new(Row)
	b.c = new(Column)
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
	/*
	if b.rng.Intn(100) >= 90 {
		num = 4
	}
	for {
		x := b.rng.Intn(b.size)
		y := b.rng.Intn(b.size)
		if b.tiles[x][y] == 0 {
			b.tiles[x][y] = num
			return
		}
	}
	*/
}

func (b *Board) Reset() {
	b.score = 0
	for _,r := range b.tiles {
		for i,_ := range r {
			r[i] = 0
		}
	}
}

//An 'Array' used to generically talk about rows or columns in
//either direction
type Iter interface {
	At(int) int
	Set(int, int)
	Len() int
	Comp(int,int) bool
}

func Shift(it Iter) (bool, int) {
	change := false
	score := 0
	for i := 0; i < it.Len() - 1; i++ {
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
					if it.Comp(i,j) {
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
	b.r.I = i
	b.r.grid = b.tiles
	b.r.rev = rev
	b.r.size = b.size
	if rev {
		b.r.alt = b.size - 1
	}
	return b.r
}

func (r *Row) At(i int) int {
	if r.rev {
		return r.grid[r.alt - i][r.I]
	} else {
		return r.grid[i][r.I]
	}
}

func (r *Row) Set(i, v int) {
	if r.rev {
		r.grid[r.alt - i][r.I] = v
	} else {
		r.grid[i][r.I] = v
	}
}

func (r *Row) Len() int {
	return r.size
}

func (r *Row) Comp(i,j int) bool {
	if r.rev {
		return r.grid[r.alt - i][r.I] == r.grid[r.alt - j][r.I]
	} else {
		return r.grid[i][r.I] == r.grid[j][r.I]
	}
}

func (b *Board) GetColumn(i int, rev bool) *Column {
	b.c.col = b.tiles[i]
	b.c.rev = rev
	b.c.size = b.size
	if rev {
		b.c.alt = b.size - 1
	}
	return b.c
}

func (c *Column) At(i int) int {
	if c.rev {
		return c.col[c.alt - i]
	} else {
		return c.col[i]
	}
}

func (c *Column) Set(i, v int) {
	if c.rev {
		c.col[c.alt - i] = v
	} else {
		c.col[i] = v
	}
}

func (c *Column) Len() int {
	return c.size
}

func (c *Column) Comp(i,j int) bool {
	if c.rev {
		return c.col[c.alt - i] == c.col[c.alt - j]
	} else {
		return c.col[i] == c.col[j]
	}
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
	//nb := NewBoard(b.size)
	nb := new(Board)
	nb.size = b.size
	nb.c = new(Column)
	nb.r = new(Row)
	nb.tiles = make([][]int, b.size)
	for i,_ := range nb.tiles {
		nb.tiles[i] = make([]int, b.size)
	}
	nb.score = b.score
	//nb.rng = b.rng
	for i,rs := range b.tiles {
		copy(nb.tiles[i], rs)
	}
	return nb
}

//Did the player win?
func (b *Board) CheckWin() bool {
	if b.score < 16000 { //Approx value, too lazy to find a better one
		return false
	}
	for _,r := range b.tiles {
		for _,v := range r {
			if v == 2048 {
				return true
			}
		}
	}
	return false
}

func (b *Board) Shift(i int) bool {
	switch i {
	case UP:
		return b.Up()
	case RIGHT:
		return b.Right()
	case DOWN:
		return b.Down()
	case LEFT:
		return b.Left()
	default:
		fmt.Println("Invalid Key")
		return false
	}
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

func (b *Board) WeightedSum() int {
	sum := 0
	for _,r := range b.tiles {
		for _,v := range r {
			sum += (v * v)
		}
	}
	return sum
}

func (b *Board) SetTo(ob *Board) {
	b.score = ob.score
	b.size = ob.size
	for i,rs := range ob.tiles {
		copy(b.tiles[i], rs)
	}
}

