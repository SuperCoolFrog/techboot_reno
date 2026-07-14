package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"techboot_reno/cmd/assets"
)

const (
	S3GridXCount          int = 64
	S3GridYCount          int = 48
	S3GridCellSize        int = 20
	S3BufferCapacity      int = 35
	S3BufferInputCapacity int = S3BufferCapacity - 2
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
	dividerX := 36

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
	// ...Commands
	hdrCmdX := dividerX/2 - 4
	gs.Set(gridId, hdrCmdX+1, 0, CellTypeChar, 'C')
	gs.Set(gridId, hdrCmdX+2, 0, CellTypeChar, 'O')
	gs.Set(gridId, hdrCmdX+3, 0, CellTypeChar, 'M')
	gs.Set(gridId, hdrCmdX+4, 0, CellTypeChar, 'M')
	gs.Set(gridId, hdrCmdX+5, 0, CellTypeChar, 'A')
	gs.Set(gridId, hdrCmdX+6, 0, CellTypeChar, 'N')
	gs.Set(gridId, hdrCmdX+7, 0, CellTypeChar, 'D')
	gs.Set(gridId, hdrCmdX+8, 0, CellTypeChar, 'S')

	// ...Output
	rightPanelHeaderX := S3GridXCount - dividerX/2
	gs.Set(gridId, rightPanelHeaderX+1, 0, CellTypeChar, 'O')
	gs.Set(gridId, rightPanelHeaderX+2, 0, CellTypeChar, 'U')
	gs.Set(gridId, rightPanelHeaderX+3, 0, CellTypeChar, 'T')
	gs.Set(gridId, rightPanelHeaderX+4, 0, CellTypeChar, 'P')
	gs.Set(gridId, rightPanelHeaderX+5, 0, CellTypeChar, 'U')
	gs.Set(gridId, rightPanelHeaderX+6, 0, CellTypeChar, 'T')

	// ..Bottom
	for i := 1; i < S3GridXCount-1; i++ {
		gs.SetCellSprite(gridId, i, S3GridYCount-1, assets.SpriteIDHorizontalBar)
	}

	// Dividers
	// .Vertical
	gs.SetCellSprite(gridId, dividerX, 0, assets.SpriteIDDownConnectBar)
	for i := 1; i < S3GridYCount-1; i++ {
		gs.SetCellSprite(gridId, dividerX, i, assets.SpriteIDVerticalBar)
	}
	gs.SetCellSprite(gridId, dividerX, S3GridYCount-1, assets.SpriteIDUpConnectBar)
	// .Horizontal
	verticalY := S3GridYCount - 11
	gs.SetCellSprite(gridId, dividerX, verticalY, assets.SpriteIDRightConnectBar)
	for i := 1; i < S3GridXCount-dividerX; i++ {
		gs.SetCellSprite(gridId, dividerX+i, verticalY, assets.SpriteIDHorizontalBar)
	}
	gs.SetCellSprite(gridId, S3GridXCount-1, verticalY, assets.SpriteIDLeftConnectBar)
	// ...Logs
	gs.Set(gridId, rightPanelHeaderX+1, verticalY, CellTypeChar, 'L')
	gs.Set(gridId, rightPanelHeaderX+2, verticalY, CellTypeChar, 'O')
	gs.Set(gridId, rightPanelHeaderX+3, verticalY, CellTypeChar, 'G')
	gs.Set(gridId, rightPanelHeaderX+4, verticalY, CellTypeChar, 'S')

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
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		if S3CurrentBufferCursor > 1 {
			S3History[S3HistoryCursor-1][S3CurrentBufferCursor] = ' '
			S3CurrentBufferCursor--
			S3History[S3HistoryCursor-1][S3CurrentBufferCursor] = '|'
		}
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

	if S3HistoryHead > 0 {
		gs.SetCellSprite(GridIdScene3, 1, 1, assets.SpriteIDCarrotUp)
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
