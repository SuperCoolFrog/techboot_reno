package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"techboot_reno/cmd/assets"
)

const (
	S3GridXCount   int = 64
	S3GridYCount   int = 48
	S3GridCellSize int = 20

	CommandBufferCapacity int = 2000 // @TODO Ring Buffer
	CommandBufferCols     int = 35
	CommandBufferRows     int = S3GridYCount - 2
	CommandBufferX        int = 1
	CommandBufferY        int = 1
)

var (
	GridIdScene3   GridID
	CommandBuffer  *Buffer
	CmdBufferDecor = BufferDecorator{
		Prefix:  []byte{':'},
		Postfix: []byte{'|'},
	}
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

	GridIdScene3 = gridId

	CommandBuffer = NewBuffer(CommandBufferCols, CommandBufferRows, CommandBufferCapacity, false)
	CommandBuffer.AppendDecorators(CmdBufferDecor)

	return next
}

func Scene3_InputHandler(runes []rune, current, next GameState, gs *GridSystem) GameState {
	for i := 0; i < len(runes); i++ {
		CommandBuffer.AppendWithDecor(byte(runes[i]), CmdBufferDecor)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		CommandBuffer.TrimDecor(CmdBufferDecor)
		CommandBuffer.NewLine()
		CommandBuffer.AppendDecorators(CmdBufferDecor)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		CommandBuffer.DecrementCursorWithDecor(CmdBufferDecor)
	}

	CommandBuffer.DrawToGrid(GridIdScene3, CommandBufferX, CommandBufferY, gs)

	return current
}
