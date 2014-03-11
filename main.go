package main

import (
	"fmt"
	"math/rand"
	"time"
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

func (b *Board) CheckLoss() bool {
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

func main() {
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
				if b.Up() {
					b.PlaceRandom()
				}
			case 'a':
				if b.Left() {
					b.PlaceRandom()
				}
			case 's':
				if b.Down() {
					b.PlaceRandom()
				}
			case 'd':
				if b.Right() {
					b.PlaceRandom()
				}
			default:
				fmt.Println("Invalid Key")
		}
		fmt.Printf("Score: %d\n", b.score)
		b.PrintBoard()
	}
}
