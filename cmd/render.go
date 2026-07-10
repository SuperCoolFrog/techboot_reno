package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
)

type RenderFlag uint64

const (
	RenderFlagNone RenderFlag = iota
	RenderFlagEmpty
	RenderFlagKeyCode
	RenderFlagReserved
)

type Grid struct {
	Cols     int
	Rows     int
	CellSize int
	Padding  int

	/* SoA - lite*/

	Flags []RenderFlag
	Chars []byte

	IsBuffer     []bool
	BufferCols   []int
	BufferRows   []int
	BufferCursor []int
}

func NewGrid(cols, rows, cellSize, padding int) Grid {
	capacity := cols * rows

	grid := Grid{
		Cols:         cols,
		Rows:         rows,
		CellSize:     cellSize,
		Padding:      padding,
		Flags:        make([]RenderFlag, capacity),
		Chars:        make([]byte, capacity),
		IsBuffer:     make([]bool, capacity),
		BufferCols:   make([]int, capacity),
		BufferRows:   make([]int, capacity),
		BufferCursor: make([]int, capacity),
	}

	return grid
}

func (g Grid) IdxFromXY(x, y int) int {
	return y*g.Cols + x
}

func (g Grid) XYFromIdx(idx int) (x int, y int) {
	x = idx / g.Cols
	y = idx % g.Cols

	return x, y
}

func (g Grid) Set(x int, y int, flag RenderFlag, char byte) {
	idx := g.IdxFromXY(x, y)
	g.Flags[idx] = flag
	g.Chars[idx] = char
}

func (g Grid) Get(x int, y int) (flag RenderFlag, char byte) {
	idx := g.IdxFromXY(x, y)
	flag = g.Flags[idx]
	char = g.Chars[idx]

	return flag, char
}

func (g Grid) NewBuffer(X, Y, Cols, Rows int) int {
	for bufferX := range Cols {
		for bufferY := range Rows {
			x := X + bufferX
			y := Y + bufferY

			g.Set(x, y, RenderFlagReserved, '_')
		}
	}

	idx := g.IdxFromXY(X, Y)
	g.IsBuffer[idx] = true
	g.BufferCols[idx] = Cols
	g.BufferRows[idx] = Rows
	g.BufferCursor[idx] = 0

	return idx
}

func (g Grid) BufferAppend(bufferId int, char byte) (success bool) {

	if !g.IsBuffer[bufferId] {
		return false
	}

	cursor := g.BufferCursor[bufferId]
	cols := g.BufferCols[bufferId]
	rows := g.BufferRows[bufferId]
	size := cols * rows

	if cursor >= size {
		return false
	}

	cursorX := cursor % cols
	cursorY := cursor / cols

	x, y := g.XYFromIdx(bufferId)

	idx := g.IdxFromXY(x+cursorX, y+cursorY)

	g.Chars[idx] = char
	g.BufferCursor[bufferId] = cursor + 1

	return true
}

func (g Grid) RenderDebug(screen *ebiten.Image) error {
	size := float32(g.CellSize)
	strokeW := float32(0.5)
	clr := color.RGBA{R: 0, G: 255, B: 0, A: 255}
	padding := float32(g.Padding)
	face := &text.GoTextFace{
		Source: arialFontSource,
		Size:   fontSize, // Use consistent font size
	}

	for i := range g.Flags {
		x := float32(i%g.Cols)*size + padding
		y := float32(math.Trunc(float64(i/g.Cols)))*size + padding

		vector.StrokeRect(screen, x, y, size, size, strokeW, clr, true)

		// @TODO Remove from RenderDebug
		if g.Chars[i] != 0 {
			opt := &text.DrawOptions{}

			charStr := string(g.Chars[i])

			// Set text color to green
			opt.ColorScale.ScaleWithColor(color.RGBA{R: 0, G: 255, B: 0, A: 255})

			w, h := text.Measure(charStr, face, 0.0)

			// centerX := x + size/2
			// centerY := y + size/2

			charX := x + (size-float32(w))/2
			charY := y + (size-float32(h))/2

			opt.GeoM.Translate(float64(charX), float64(charY)) // Set Position

			text.Draw(screen, charStr, face, opt)
		}
	}

	return nil
}
