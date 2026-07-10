package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
)

type RenderFlag uint64

const (
	RenderFlagNone RenderFlag = iota
	RenderFlagEmpty
	RenderFlagKeyCode
	RenderFlag0
)

type Grid struct {
	Cols int
	Rows int
	CellSize int

	/* SoA - lite*/

	Flags []RenderFlag
	Chars []byte
}

func NewGrid(cols, rows, cellSize int) Grid {
	capacity := cols*rows

	grid := Grid{
		Cols: cols,
		Rows: rows,
		CellSize: cellSize,
		Flags: make([]RenderFlag, capacity),
		Chars: make([]byte, capacity),
	} 



	return grid
}

func (g Grid) Set(x int, y int, flag RenderFlag, char byte) {
	idx := y*g.Cols+x
	g.Flags[idx] = flag
	g.Chars[idx] = char
}

func (g Grid) Get(x int, y int) (flag RenderFlag, char byte) {
	idx := y*g.Cols+x
	flag = g.Flags[idx]
	char = g.Chars[idx]

	return flag, char
}

func (g Grid) SetAllCells(flag RenderFlag, char byte) {
	for i := range g.Flags {
		g.Flags[i] = flag
		g.Chars[i] = char
	}
}

func (g Grid) RenderDebug(screen *ebiten.Image) error {
	size := float32(g.CellSize)
	strokeW := float32(0.5)
	clr := color.RGBA{R: 0, G: 255, B: 0, A: 255}

	for i := range g.Flags {
		x := float32(i % g.Cols) * size
		y := float32(math.Trunc(float64(i / g.Cols))) * size

		vector.StrokeRect(screen, x, y, size, size, strokeW, clr, true)
	}

	return nil
}


