package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
	"math"
	"techboot_reno/cmd/assets"
	"unsafe"
)

type GridID uint32
type GridCellType byte

const (
	CellTypeNone GridCellType = iota
	CellTypeEmpty
	CellTypeChar
	CellTypeReserved
	CellTypeSquare
	CellTypeSprite
)

type GridSystem struct {
	MaxTotalCells int
	MasterBuffer  []byte

	// Flat, parallel master slices targeting the MasterBuffer
	CellTypes    []GridCellType
	Chars        []byte
	SpriteIDs    []assets.SpriteID
	IsBuffer     []bool
	BufferCols   []int
	BufferRows   []int
	BufferCursor []int

	// Metadata tracking tables (indexed directly by GridId)
	Offsets   []int
	Counts    []int
	Rows      []int
	Cols      []int
	CellSizes []int
	OffsetX   []int
	OffsetY   []int
	IsActive  []bool

	NextGridID GridID
}

func NewGridSystem(maxTotalCells int, maxGrid int) *GridSystem {
	gs := &GridSystem{
		MaxTotalCells: maxTotalCells,
		Offsets:       make([]int, maxGrid),
		Counts:        make([]int, maxGrid),
		Rows:          make([]int, maxGrid),
		Cols:          make([]int, maxGrid),
		CellSizes:     make([]int, maxGrid),
		OffsetX:       make([]int, maxGrid),
		OffsetY:       make([]int, maxGrid),
		IsActive:      make([]bool, maxGrid),
		NextGridID:    0,
	}

	// 1. Calculate byte sizes based on the absolute maximum total cell capacity
	sizeCellTypes := maxTotalCells * int(unsafe.Sizeof(GridCellType(0)))
	sizeChars := maxTotalCells * int(unsafe.Sizeof(byte(0)))
	sizeSpriteID := maxTotalCells * int(unsafe.Sizeof(assets.SpriteID(0)))
	sizeIsBuffer := maxTotalCells * int(unsafe.Sizeof(bool(false)))
	sizeCols := maxTotalCells * int(unsafe.Sizeof(int(0)))
	sizeRows := maxTotalCells * int(unsafe.Sizeof(int(0)))
	sizeCursor := maxTotalCells * int(unsafe.Sizeof(int(0)))

	totalByteSize := sizeCellTypes + sizeChars + sizeIsBuffer + sizeCols + sizeRows + sizeCursor

	// 2. Allocate the single massive global block
	gs.MasterBuffer = make([]byte, totalByteSize)
	ptr := unsafe.Pointer(&gs.MasterBuffer[0])

	// 3. Slice up the global memory into parallel arrays
	gs.CellTypes = unsafe.Slice((*GridCellType)(ptr), maxTotalCells)
	ptr = unsafe.Add(ptr, sizeCellTypes)

	gs.Chars = unsafe.Slice((*byte)(ptr), maxTotalCells)
	ptr = unsafe.Add(ptr, sizeChars)

	gs.SpriteIDs = unsafe.Slice((*assets.SpriteID)(ptr), maxTotalCells)
	ptr = unsafe.Add(ptr, sizeSpriteID)

	gs.IsBuffer = unsafe.Slice((*bool)(ptr), maxTotalCells)
	ptr = unsafe.Add(ptr, sizeIsBuffer)

	gs.BufferCols = unsafe.Slice((*int)(ptr), maxTotalCells)
	ptr = unsafe.Add(ptr, sizeCols)

	gs.BufferRows = unsafe.Slice((*int)(ptr), maxTotalCells)
	ptr = unsafe.Add(ptr, sizeRows)

	gs.BufferCursor = unsafe.Slice((*int)(ptr), maxTotalCells)

	return gs
}

func (gs *GridSystem) AllocateGrid(cols, rows, cellSize, offsetX, offsetY int) GridID {
	id := gs.NextGridID
	gs.NextGridID++

	cellCount := cols * rows

	var startOffset int = 0
	if id > 0 {
		prevID := id - 1
		startOffset = gs.Offsets[prevID] + gs.Counts[prevID]
	}

	if startOffset+cellCount > gs.MaxTotalCells {
		panic("GridSystem out of total cell capacity")
	}

	gs.Offsets[id] = startOffset
	gs.Counts[id] = cellCount
	gs.Cols[id] = cols
	gs.Rows[id] = rows
	gs.CellSizes[id] = cellSize
	gs.OffsetX[id] = offsetX
	gs.OffsetY[id] = offsetY
	gs.IsActive[id] = false

	return id
}

func (gs *GridSystem) EnableGrid(gridId GridID) {
	gs.IsActive[gridId] = true
}

func (gs *GridSystem) DisableGrid(gridId GridID) {
	gs.IsActive[gridId] = false
}

func (gs *GridSystem) GridRectangle(gridId GridID, x1, y1, colCount, rowCount int) image.Rectangle {
	size := gs.CellSizes[gridId]
	offsetX := gs.OffsetX[gridId]
	offsetY := gs.OffsetY[gridId]

	posX := x1*size + offsetX
	posY := y1*size + offsetY

	w := posX + colCount*size
	h := posY + rowCount*size

	return image.Rect(posX, posY, w, h)
}

func (gs *GridSystem) IdxFromXY(gridId GridID, x, y int) int {
	cols := gs.Cols[gridId]
	offset := gs.Offsets[gridId]
	return offset + (y*cols + x)
}

func (gs *GridSystem) XYFromBufferIdx(gridId GridID, globalIdx int) (x int, y int) {
	offset := gs.Offsets[gridId]
	idx := globalIdx - offset // this assumes the bufferIdx is a global/master buffer idx not relative to

	cols := gs.Cols[gridId]

	x = idx / cols
	y = idx % cols

	return x, y
}

func (gs *GridSystem) Set(gridId GridID, x int, y int, flag GridCellType, char byte) {
	// Guard rails to protect neighboring grid data
	if x < 0 || x >= gs.Cols[gridId] || y < 0 || y >= gs.Rows[gridId] {
		panic("Grid cell coordinates out of bounds! " + string(x) + "," + string(y))
	}

	idx := gs.IdxFromXY(gridId, x, y)

	gs.CellTypes[idx] = flag
	gs.Chars[idx] = char
}

func (gs *GridSystem) SetCellSprite(gridId GridID, x int, y int, spriteId assets.SpriteID) {
	// Guard rails to protect neighboring grid data
	if x < 0 || x >= gs.Cols[gridId] || y < 0 || y >= gs.Rows[gridId] {
		panic("Grid cell coordinates out of bounds!")
	}

	idx := gs.IdxFromXY(gridId, x, y)

	gs.CellTypes[idx] = CellTypeSprite
	gs.SpriteIDs[idx] = spriteId
}

func (gs *GridSystem) SetAllCells(gridId GridID, cellType GridCellType, char byte) {
	offset := gs.Offsets[gridId]
	count := gs.Counts[gridId]

	gridChars := gs.Chars[offset : offset+count]
	gridCellTypes := gs.CellTypes[offset : offset+count]

	for i := 0; i < len(gridChars); i++ {
		gridCellTypes[i] = cellType
		gridChars[i] = char
	}
}

func (gs *GridSystem) Get(gridId GridID, x int, y int) (cellType GridCellType, char byte) {
	idx := gs.IdxFromXY(gridId, x, y)
	cellType = gs.CellTypes[idx]
	char = gs.Chars[idx]

	return cellType, char
}

func (gs *GridSystem) GetCellSprite(gridId GridID, x int, y int) (image *ebiten.Image) {
	idx := gs.IdxFromXY(gridId, x, y)
	spriteId := gs.SpriteIDs[idx]

	img, found := assets.Images[spriteId]
	if !found {
		panic("Could not find sprite for spriteId: " + string(spriteId))
	}

	return img
}

func (gs *GridSystem) NewBuffer(gridId GridID, bufferX, bufferY, bufferCols, bufferRows int) (globalIdx int) {
	gridCols := gs.Cols[gridId]
	gridRows := gs.Rows[gridId]

	// 1. Precise 2D Matrix Guard Rails
	// Check if the starting point is valid
	if bufferX < 0 || bufferY < 0 {
		panic("Buffer starting coordinates cannot be negative!")
	}
	// Check if the sub-buffer physically bleeds out of the grid's right or bottom edges
	if bufferX+bufferCols > gridCols || bufferY+bufferRows > gridRows {
		panic("Sub-buffer geometry overflows the grid boundaries!")
	}

	// 2. Corrected Target Loops (Iterate from start coordinate to start + dimension size)
	endX := bufferX + bufferCols
	endY := bufferY + bufferRows

	for y := bufferY; y < endY; y++ {
		for x := bufferX; x < endX; x++ {
			gs.Set(gridId, x, y, CellTypeReserved, ' ')
		}
	}

	// 3. Store metadata at the buffer's root cell location
	bufferIdx := gs.IdxFromXY(gridId, bufferX, bufferY)

	gs.IsBuffer[bufferIdx] = true
	gs.BufferCols[bufferIdx] = bufferCols
	gs.BufferRows[bufferIdx] = bufferRows
	gs.BufferCursor[bufferIdx] = 0

	return bufferIdx
}

func (gs *GridSystem) BufferAppend(gridId GridID, bufferId int, char byte) (success bool) {
	if !gs.IsBuffer[bufferId] {
		return false
	}

	cursor := gs.BufferCursor[bufferId]
	cols := gs.BufferCols[bufferId]
	rows := gs.BufferRows[bufferId]
	size := cols * rows

	if cursor >= size {
		return false
	}

	cursorX := cursor % cols
	cursorY := cursor / cols

	// 2. Compute the exact linear step within the parent grid matrix
	// Moving down a row in the sub-buffer means stepping forward by a full parent grid width
	gridCols := gs.Cols[gridId]
	linearStep := (cursorY * gridCols) + cursorX

	targetIdx := bufferId + linearStep

	gs.Chars[targetIdx] = char
	gs.BufferCursor[bufferId] = cursor + 1

	return true
}

func (gs *GridSystem) RenderDebug(screen *ebiten.Image, gridID GridID) error {
	offset := gs.Offsets[gridID]
	count := gs.Counts[gridID]
	size := float32(gs.CellSizes[gridID])
	offsetX := float32(gs.OffsetX[gridID])
	offsetY := float32(gs.OffsetY[gridID])
	cols := gs.Cols[gridID]

	chars := gs.Chars[offset : offset+count]
	cellTypes := gs.CellTypes[offset : offset+count]

	strokeW := float32(0.5)
	clr := color.RGBA{R: 0, G: 255, B: 0, A: 255}
	face := &text.GoTextFace{
		Source: fontSrc,
		Size:   fontSize, // Use consistent font size
	}

	for i := 0; i < len(cellTypes); i++ {
		x := float32(i%cols)*size + offsetX
		y := float32(math.Trunc(float64(i/cols)))*size + offsetY

		vector.StrokeRect(screen, x, y, size, size, strokeW, clr, true)

		// @TODO Remove from RenderDebug
		if chars[i] != 0 {
			opt := &text.DrawOptions{}

			charStr := string(chars[i])

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

func (gs *GridSystem) Render(screen *ebiten.Image) error {
	for id := GridID(0); id < gs.NextGridID; id++ {
		if gs.IsActive[id] {
			gs.RenderGrid(screen, id)
		}
	}

	return nil
}

func (gs *GridSystem) RenderGrid(screen *ebiten.Image, gridID GridID) error {
	offset := gs.Offsets[gridID]
	count := gs.Counts[gridID]
	size := float32(gs.CellSizes[gridID])
	offsetX := float32(gs.OffsetX[gridID])
	offsetY := float32(gs.OffsetY[gridID])
	cols := gs.Cols[gridID]
	capacity := gs.Counts[gridID]

	chars := gs.Chars[offset : offset+count]
	cellTypes := gs.CellTypes[offset : offset+count]
	spriteIds := gs.SpriteIDs[offset : offset+count]

	strokeW := float32(0.5)
	clr := color.RGBA{R: 0, G: 255, B: 0, A: 255}

	face := &text.GoTextFace{
		Source: fontSrc,
		Size:   float64(size),
	}

	for i := range capacity {
		x := float32(i%cols)*size + offsetX
		y := float32(math.Trunc(float64(i/cols)))*size + offsetY

		switch cellTypes[i] {
		case CellTypeSprite:
			RenderSpriteCell(screen, x, y, size, spriteIds[i])
		case CellTypeSquare:
			vector.StrokeRect(screen, x, y, size, size, strokeW, clr, true)
		case CellTypeReserved:
			opt := &text.DrawOptions{}
			charStr := string(chars[i])
			opt.ColorScale.ScaleWithColor(color.RGBA{R: 0, G: 255, B: 0, A: 255})
			w, h := text.Measure(charStr, face, 0.0)
			charX := x + (size-float32(w))/2
			charY := y + (size-float32(h))/2
			opt.GeoM.Translate(float64(charX), float64(charY)) // Set Position
			text.Draw(screen, charStr, face, opt)
		case CellTypeChar:
			opt := &text.DrawOptions{}
			charStr := string(chars[i])
			opt.ColorScale.ScaleWithColor(color.RGBA{R: 0, G: 255, B: 0, A: 255})
			w, h := text.Measure(charStr, face, 0.0)
			charX := x + (size-float32(w))/2
			charY := y + (size-float32(h))/2
			opt.GeoM.Translate(float64(charX), float64(charY)) // Set Position
			text.Draw(screen, charStr, face, opt)
		}
	}

	return nil
}

func RenderSpriteCell(screen *ebiten.Image, x, y, size float32, spriteId assets.SpriteID) {
	sprite, found := assets.Images[spriteId]
	if !found {
		panic("Cell image not found for spriteId: " + string(spriteId))
	}

	x64 := float64(x)
	y64 := float64(y)
	size64 := float64(size)

	spriteRect := sprite.Bounds()
	spriteW := float64(spriteRect.Dx())
	spriteH := float64(spriteRect.Dy())

	op := &ebiten.DrawImageOptions{}

	scaleX := size64 / spriteW
	scaleY := size64 / spriteH
	op.GeoM.Scale(scaleX, scaleY)

	scaledW := (spriteW * scaleX)
	scaledH := (spriteW * scaleY)

	spriteX := (x64 + size64/2) - scaledW/2
	spriteY := (y64 + size64/2) - scaledH/2
	op.GeoM.Translate(spriteX, spriteY)

	op.ColorScale.ScaleWithColor(fontColor)

	screen.DrawImage(sprite, op)
}
