package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
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
	Padding int

	/* SoA - lite*/

	Flags []RenderFlag
	Chars []byte
}

func NewGrid(cols, rows, cellSize, padding int) Grid {
	capacity := cols*rows

	grid := Grid{
		Cols: cols,
		Rows: rows,
		CellSize: cellSize,
		Padding: padding,
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
	padding := float32(g.Padding)
	face := &text.GoTextFace{
		Source: arialFontSource,
		Size:   fontSize, // Use consistent font size
	}

	for i := range g.Flags {
		x := float32(i % g.Cols) * size + padding
		y := float32(math.Trunc(float64(i / g.Cols))) * size + padding

		vector.StrokeRect(screen, x, y, size, size, strokeW, clr, true)

		if g.Flags[i] == RenderFlagKeyCode {
			opt := &text.DrawOptions{}

			charStr := string(g.Chars[i])

			// Set text color to green
			opt.ColorScale.ScaleWithColor(color.RGBA{R: 0, G: 255, B: 0, A: 255})

			w, h:= text.Measure(charStr, face, 0.0) 

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


