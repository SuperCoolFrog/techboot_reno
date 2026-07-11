package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
	"unsafe"
)

type RenderFlag uint8

const (
	RenderFlagNone RenderFlag = iota
	RenderFlagEmpty
	RenderFlagKeyCode
	RenderFlagReserved
	RenderFlagCellSquare
)

type Grid struct {
	Cols        int
	Rows        int
	CellSize    int
	Padding     int
	Capacity    int // Current Capacity for Grid
	MaxCapacity int // Max Capacity - used for reusing grids

	/* SoA - lite*/
	MasterBuffer []byte // stores all of the following in one continuous chunk of mem

	Flags []RenderFlag
	Chars []byte

	IsBuffer     []bool
	BufferCols   []int
	BufferRows   []int
	BufferCursor []int
}

func NewGrid(cols, rows, cellSize, padding int) *Grid {
	g := &Grid{}
	g.Cols = cols
	g.Rows = rows
	g.CellSize = cellSize
	g.Padding = padding
	g.Capacity = cols * rows
	g.MaxCapacity = g.Capacity

	// 1. Calculate the required byte sizes for each array block
	sizeFlags := g.MaxCapacity * int(unsafe.Sizeof(RenderFlag(0)))  // MaxCap * 1
	sizeChars := g.MaxCapacity * int(unsafe.Sizeof(byte(0)))        // MaxCap * 1
	sizeIsBuffer := g.MaxCapacity * int(unsafe.Sizeof(bool(false))) // MaxCap * 1
	sizeBufferCols := g.MaxCapacity * int(unsafe.Sizeof(int(0)))    // MaxCap * 8
	sizeBufferRows := g.MaxCapacity * int(unsafe.Sizeof(int(0)))    // MaxCap * 8
	sizeBufferCursor := g.MaxCapacity * int(unsafe.Sizeof(int(0)))  // MaxCap * 8

	// 2. Sum the absolute total size for the master buffer allocation
	totalByteSize := sizeFlags + sizeChars + sizeIsBuffer + sizeBufferCols + sizeBufferRows + sizeBufferCursor

	// 3. Trigger a single contiguous heap allocation
	g.MasterBuffer = make([]byte, totalByteSize)

	// 4. Set up the sliding pointer to partition the master buffer
	ptr := unsafe.Pointer(&g.MasterBuffer[0])

	// 5. Build individual slice headers targeting the master buffer memory
	g.Flags = unsafe.Slice((*RenderFlag)(ptr), g.MaxCapacity)
	ptr = unsafe.Add(ptr, sizeFlags)

	g.Chars = unsafe.Slice((*byte)(ptr), g.MaxCapacity)
	ptr = unsafe.Add(ptr, sizeChars)

	g.IsBuffer = unsafe.Slice((*bool)(ptr), g.MaxCapacity)
	ptr = unsafe.Add(ptr, sizeIsBuffer)

	g.BufferCols = unsafe.Slice((*int)(ptr), g.MaxCapacity)
	ptr = unsafe.Add(ptr, sizeBufferCols)

	g.BufferRows = unsafe.Slice((*int)(ptr), g.MaxCapacity)
	ptr = unsafe.Add(ptr, sizeBufferRows)

	g.BufferCursor = unsafe.Slice((*int)(ptr), g.MaxCapacity)

	return g
}

func (g *Grid) ResetAndResize(cols, rows, cellSize, padding int) {
	nuCap := cols * rows

	if nuCap > g.MaxCapacity {
		if g.MaxCapacity == 0 {
			g.MaxCapacity = 16
		}
		for nuCap > g.MaxCapacity {
			g.MaxCapacity *= 2
		}

		// 1. Calculate the required byte sizes for each array block
		sizeFlags := g.MaxCapacity * int(unsafe.Sizeof(RenderFlag(0)))  // MaxCap * 1
		sizeChars := g.MaxCapacity * int(unsafe.Sizeof(byte(0)))        // MaxCap * 1
		sizeIsBuffer := g.MaxCapacity * int(unsafe.Sizeof(bool(false))) // MaxCap * 1
		sizeBufferCols := g.MaxCapacity * int(unsafe.Sizeof(int(0)))    // MaxCap * 8
		sizeBufferRows := g.MaxCapacity * int(unsafe.Sizeof(int(0)))    // MaxCap * 8
		sizeBufferCursor := g.MaxCapacity * int(unsafe.Sizeof(int(0)))  // MaxCap * 8

		// 2. Sum the absolute total size for the master buffer allocation
		totalByteSize := sizeFlags + sizeChars + sizeIsBuffer + sizeBufferCols + sizeBufferRows + sizeBufferCursor

		// 3. Trigger a single massive heap allocation
		newMasterBuffer := make([]byte, totalByteSize)

		// 4. Back up old data before re-segmenting (if previous data exists)
		oldFlags := g.Flags
		oldChars := g.Chars
		oldIsBuffer := g.IsBuffer
		oldBufferCols := g.BufferCols
		oldBufferRows := g.BufferRows
		oldBufferCursor := g.BufferCursor

		// 5. Partition the master buffer into slices using slice headers via unsafe pointer conversion
		// We establish safe start pointers using byte index offsets
		ptr := unsafe.Pointer(&newMasterBuffer[0])

		g.Flags = unsafe.Slice((*RenderFlag)(ptr), g.MaxCapacity)
		ptr = unsafe.Add(ptr, sizeFlags)

		g.Chars = unsafe.Slice((*byte)(ptr), g.MaxCapacity)
		ptr = unsafe.Add(ptr, sizeChars)

		g.IsBuffer = unsafe.Slice((*bool)(ptr), g.MaxCapacity)
		ptr = unsafe.Add(ptr, sizeIsBuffer)

		g.BufferCols = unsafe.Slice((*int)(ptr), g.MaxCapacity)
		ptr = unsafe.Add(ptr, sizeBufferCols)

		g.BufferRows = unsafe.Slice((*int)(ptr), g.MaxCapacity)
		ptr = unsafe.Add(ptr, sizeBufferRows)

		g.BufferCursor = unsafe.Slice((*int)(ptr), g.MaxCapacity)

		// 6. Copy over existing data into our new sub-slices
		if len(oldFlags) > 0 {
			copy(g.Flags, oldFlags)
			copy(g.Chars, oldChars)
			copy(g.IsBuffer, oldIsBuffer)
			copy(g.BufferCols, oldBufferCols)
			copy(g.BufferRows, oldBufferRows)
			copy(g.BufferCursor, oldBufferCursor)
		}

		// Keep a reference to the root master slice so Go doesn't GC it
		g.MasterBuffer = newMasterBuffer
	}

	// Clear/Zero out master buffer
	clear(g.MasterBuffer)

	// 7. Strictly limit the active viewable lengths to the current runtime window (nuCap)
	g.Capacity = nuCap
	g.Flags = g.Flags[:nuCap]
	g.Chars = g.Chars[:nuCap]
	g.IsBuffer = g.IsBuffer[:nuCap]
	g.BufferCols = g.BufferCols[:nuCap]
	g.BufferRows = g.BufferRows[:nuCap]
	g.BufferCursor = g.BufferCursor[:nuCap]
}

func (g *Grid) IdxFromXY(x, y int) int {
	return y*g.Cols + x
}

func (g *Grid) XYFromIdx(idx int) (x int, y int) {
	x = idx / g.Cols
	y = idx % g.Cols

	return x, y
}

func (g *Grid) Set(x int, y int, flag RenderFlag, char byte) {
	idx := g.IdxFromXY(x, y)
	g.Flags[idx] = flag
	g.Chars[idx] = char
}

func (g *Grid) SetAllCells(flag RenderFlag, char byte) {
	for i := range g.Capacity {
		g.Flags[i] = flag
		g.Chars[i] = char
	}
}

func (g *Grid) Get(x int, y int) (flag RenderFlag, char byte) {
	idx := g.IdxFromXY(x, y)
	flag = g.Flags[idx]
	char = g.Chars[idx]

	return flag, char
}

func (g *Grid) NewBuffer(X, Y, Cols, Rows int) int {
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

func (g *Grid) BufferAppend(bufferId int, char byte) (success bool) {

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

func (g *Grid) RenderDebug(screen *ebiten.Image) error {
	size := float32(g.CellSize)
	strokeW := float32(0.5)
	clr := color.RGBA{R: 0, G: 255, B: 0, A: 255}
	padding := float32(g.Padding)
	face := &text.GoTextFace{
		Source: fontSrc,
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

func (g *Grid) Render(screen *ebiten.Image) error {
	size := float32(g.CellSize)
	strokeW := float32(0.5)
	clr := color.RGBA{R: 0, G: 255, B: 0, A: 255}
	padding := float32(g.Padding)

	for i := range g.Capacity {
		x := float32(i%g.Cols)*size + padding
		y := float32(math.Trunc(float64(i/g.Cols)))*size + padding

		if g.Flags[i] == RenderFlagCellSquare {
			vector.StrokeRect(screen, x, y, size, size, strokeW, clr, true)
		}
	}

	return nil
}
