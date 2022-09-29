package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"
)

// Board represents a two-dimensional field of cells.
type Board struct {
	s    [][]bool
	w, h int
}

// NewBoard returns an empty field of the specified width and height.
func NewBoard(w, h int) *Board {
	s := make([][]bool, h)
	for i := range s {
		s[i] = make([]bool, w)
	}
	return &Board{s: s, w: w, h: h}
}

// Set sets the state of the specified cell to the given value.
func (f *Board) Set(x, y int, b bool) {
	f.s[y][x] = b
}

// Active reports whether the specified cell is active.
// If the x or y coordinates are outside the field boundaries they are wrapped
// toroidally. For instance, an x value of -1 is treated as width-1.
func (f *Board) Active(x, y int) bool {
	x += f.w
	x %= f.w
	y += f.h
	y %= f.h
	return f.s[y][x]
}

// Next returns the state of the specified cell at the next time step.
func (f *Board) Next(x, y int) bool {
	// Count the adjacent cells that are active.
	active := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if (j != 0 || i != 0) && f.Active(x+i, y+j) {
				active++
			}
		}
	}
	// Return next state according to the game rules:
	//   exactly 3 neighbors: on,
	//   exactly 2 neighbors: maintain current state,
	//   otherwise: off.
	return active == 3 || active == 2 && f.Active(x, y)
}

// State stores the state of a round of Conway's Game of State.
type State struct {
	a, b *Board
	w, h int
}

// NewState returns a new State game state with a random initial state.
func NewState(w, h int) *State {
	a := NewBoard(w, h)
	for i := 0; i < (w * h / 4); i++ {
		a.Set(rand.Intn(w), rand.Intn(h), true)
	}
	return &State{
		a: a, b: NewBoard(w, h),
		w: w, h: h,
	}
}

// Step advances the game by one instant, recomputing and updating all cells.
func (l *State) Step() {
	// Update the state of the next field (b) from the current field (a).
	for y := 0; y < l.h; y++ {
		for x := 0; x < l.w; x++ {
			l.b.Set(x, y, l.a.Next(x, y))
		}
	}
	// Swap fields a and b.
	l.a, l.b = l.b, l.a
}

// String returns the game board as a string.
func (l *State) String() string {
	var buf bytes.Buffer
	for y := 0; y < l.h; y++ {
		for x := 0; x < l.w; x++ {
			b := byte(' ')
			if l.a.Active(x, y) {
				b = '*'
			}
			buf.WriteByte(b)
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}

func main() {
	l := NewState(40, 15)
	for i := 0; i < 300; i++ {
		l.Step()
		fmt.Print("\x0c", l) // Clear screen and print field.
		time.Sleep(time.Second / 30)
	}
}
