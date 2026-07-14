package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"techboot_reno/cmd/assets"
)

const (
	S3GridXCount          int = 40
	S3GridYCount          int = 30
	S3GridCellSize        int = 32
	S3BufferCapacity      int = 23                   // Max 23
	S3BufferInputCapacity int = S3BufferCapacity - 2 // Max 23
	S3RenderHistoryMax    int = S3GridYCount - 2
	S3HistoryY            int = 1
	S3HistoryX            int = 1
)

var (
	GridIdScene3 GridID

	S3HistoryCapacity     = 2000 // @TODO Ring Buffer
	S3History             [][]byte
	S3HistoryHead         int
	S3HistoryCursor       int
	S3CurrentBufferCursor int
)

func Scene3_HandleInit(current, next GameState, gs *GridSystem, anims *AnimationSystem) GameState {
	gridId := gs.AllocateGrid(S3GridXCount, S3GridYCount, S3GridCellSize, 0)
	gs.SetAllCells(gridId, CellTypeEmpty, 0)
	gs.EnableGrid(gridId)

	// Border
	// .Corners
	gs.SetCellSprite(gridId, 0, 0, assets.SpriteIDCornerTopLeft)
	gs.SetCellSprite(gridId, 0, S3GridYCount-1, assets.SpriteIDCornerBottomLeft)
	gs.SetCellSprite(gridId, S3GridXCount-1, 0, assets.SpriteIDCornerTopRight)
	gs.SetCellSprite(gridId, S3GridXCount-1, S3GridYCount-1, assets.SpriteIDCornerBottomRight)

	// .Walls
	// ..Left
	for i := 1; i < S3GridYCount-1; i++ {
		gs.SetCellSprite(gridId, 0, i, assets.SpriteIDVerticalBar)
	}
	// ..Right
	for i := 1; i < S3GridYCount-1; i++ {
		gs.SetCellSprite(gridId, S3GridXCount-1, i, assets.SpriteIDVerticalBar)
	}
	// ..Top
	for i := 1; i < S3GridXCount-1; i++ {
		gs.SetCellSprite(gridId, i, 0, assets.SpriteIDHorizontalBar)
	}
	// ..Bottom
	for i := 1; i < S3GridXCount-1; i++ {
		gs.SetCellSprite(gridId, i, S3GridYCount-1, assets.SpriteIDHorizontalBar)
	}

	// Dividers
	// .Vertical
	gs.SetCellSprite(gridId, 24, 0, assets.SpriteIDDownConnectBar)
	for i := 1; i < S3GridYCount-1; i++ {
		gs.SetCellSprite(gridId, 24, i, assets.SpriteIDVerticalBar)
	}
	gs.SetCellSprite(gridId, 24, S3GridYCount-1, assets.SpriteIDUpConnectBar)
	// .Horizontal
	gs.SetCellSprite(gridId, 24, S3GridYCount-11, assets.SpriteIDRightConnectBar)
	for i := 1; i < S3GridXCount-25; i++ {
		gs.SetCellSprite(gridId, 24+i, S3GridYCount-11, assets.SpriteIDHorizontalBar)
	}
	gs.SetCellSprite(gridId, S3GridXCount-1, S3GridYCount-11, assets.SpriteIDLeftConnectBar)

	S3History = make([][]byte, S3HistoryCapacity)
	S3History[0] = make([]byte, S3BufferCapacity)
	S3HistoryHead = 0
	S3HistoryCursor = 0

	s3_NextBuffer()

	GridIdScene3 = gridId

	return next
}

func Scene3_InputHandler(runes []rune, current, next GameState, gs *GridSystem) GameState {
	for i := 0; i < len(runes); i++ {
		s3_AppendToBuffer(byte(runes[i]))
	}

	s3_UpdateGrid(gs)

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		s3_SubmitLine()
	}

	return current
}

func s3_UpdateGrid(gs *GridSystem) {
	for h := 0; h < S3RenderHistoryMax; h++ {
		//for historyIdx := S3HistoryHead; historyIdx < S3HistoryCursor; historyIdx++ {
		historyIdx := S3HistoryHead + h
		bytes := S3History[historyIdx]
		for i := 0; i < len(bytes); i++ {
			gs.Set(GridIdScene3, S3HistoryX+i, S3HistoryY+h, CellTypeChar, bytes[i])
		}
	}
}

func s3_AppendToBuffer(char byte) {
	if S3CurrentBufferCursor > S3BufferInputCapacity {
		S3History[S3HistoryCursor-1][S3CurrentBufferCursor] = '|'
		return
	}

	S3History[S3HistoryCursor-1][S3CurrentBufferCursor] = char
	S3CurrentBufferCursor++
	S3History[S3HistoryCursor-1][S3CurrentBufferCursor] = '|'
}

func s3_NextBuffer() {
	S3History[S3HistoryCursor] = make([]byte, S3BufferCapacity)
	S3History[S3HistoryCursor][0] = ':'
	S3History[S3HistoryCursor][1] = '|'
	S3CurrentBufferCursor = 1
	S3HistoryCursor++
}

func s3_SubmitLine() {
	S3History[S3HistoryCursor-1][0] = ' '
	S3History[S3HistoryCursor-1][S3CurrentBufferCursor] = ' '
	s3_NextBuffer()

	for S3HistoryCursor-S3HistoryHead > S3RenderHistoryMax {
		S3HistoryHead++
	}

	//@TODO eventually implement ring buffer
}
