package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/trealla-prolog/go/trealla"
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

	DividerX int = 36
	DividerY int = S3GridYCount - 11

	LogBufferCapacity int = 1000
	LogBufferCols     int = S3GridXCount - DividerX - 2
	LogBufferRows     int = S3GridYCount - DividerY - 2
	LogBufferX        int = DividerX + 1
	LogBufferY        int = DividerY + 1
)

var (
	GridIdScene3   GridID
	CommandBuffer  *Buffer
	CmdBufferDecor = BufferDecorator{
		Prefix:  []byte{':'},
		Postfix: []byte{'|'},
	}

	LogBuffer *Buffer
)

func Scene3_HandleInit(current, next GameState, gs *GridSystem, anims *AnimationSystem) GameState {
	gridId := gs.AllocateGrid(S3GridXCount, S3GridYCount, S3GridCellSize, 0, 0)
	gs.SetAllCells(gridId, CellTypeEmpty, 0)
	gs.EnableGrid(gridId)

	InitOutputGrid(gs, anims)

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
	// ...Commands
	hdrCmdX := DividerX/2 - 4
	gs.Set(gridId, hdrCmdX+1, 0, CellTypeChar, 'C')
	gs.Set(gridId, hdrCmdX+2, 0, CellTypeChar, 'O')
	gs.Set(gridId, hdrCmdX+3, 0, CellTypeChar, 'M')
	gs.Set(gridId, hdrCmdX+4, 0, CellTypeChar, 'M')
	gs.Set(gridId, hdrCmdX+5, 0, CellTypeChar, 'A')
	gs.Set(gridId, hdrCmdX+6, 0, CellTypeChar, 'N')
	gs.Set(gridId, hdrCmdX+7, 0, CellTypeChar, 'D')
	gs.Set(gridId, hdrCmdX+8, 0, CellTypeChar, 'S')

	// ...Output
	rightPanelHeaderX := S3GridXCount - DividerX/2
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
	gs.SetCellSprite(gridId, DividerX, 0, assets.SpriteIDDownConnectBar)
	for i := 1; i < S3GridYCount-1; i++ {
		gs.SetCellSprite(gridId, DividerX, i, assets.SpriteIDVerticalBar)
	}
	gs.SetCellSprite(gridId, DividerX, S3GridYCount-1, assets.SpriteIDUpConnectBar)
	// .Horizontal
	// verticalY := S3GridYCount - 11
	gs.SetCellSprite(gridId, DividerX, DividerY, assets.SpriteIDRightConnectBar)
	for i := 1; i < S3GridXCount-DividerX; i++ {
		gs.SetCellSprite(gridId, DividerX+i, DividerY, assets.SpriteIDHorizontalBar)
	}
	gs.SetCellSprite(gridId, S3GridXCount-1, DividerY, assets.SpriteIDLeftConnectBar)
	// ...Logs
	gs.Set(gridId, rightPanelHeaderX+1, DividerY, CellTypeChar, 'L')
	gs.Set(gridId, rightPanelHeaderX+2, DividerY, CellTypeChar, 'O')
	gs.Set(gridId, rightPanelHeaderX+3, DividerY, CellTypeChar, 'G')
	gs.Set(gridId, rightPanelHeaderX+4, DividerY, CellTypeChar, 'S')

	GridIdScene3 = gridId

	CommandBuffer = NewBuffer(CommandBufferCols, CommandBufferRows, CommandBufferCapacity, false)
	CommandBuffer.AppendDecorators(CmdBufferDecor)

	LogBuffer = NewBuffer(LogBufferCols, LogBufferRows, LogBufferCapacity, false)
	LogBuffer.AppendAll([]byte("Type: connect rabbit="))
	LogBuffer.NewLine()

	return next
}

func Scene3_Update(runes []rune, current, next GameState, input chan []byte, commands chan trealla.Atom, gs *GridSystem, anims *AnimationSystem) GameState {
	for i := 0; i < len(runes); i++ {
		CommandBuffer.AppendWithDecor(byte(runes[i]), CmdBufferDecor)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		CommandBuffer.TrimDecor(CmdBufferDecor)
		CommandBuffer.NewLine()
		CommandBuffer.AppendDecorators(CmdBufferDecor)

		fmt.Printf("Enter\n")
		// ParseInput([]byte("CONNECT 42="), input)
		if lastLine, exists := CommandBuffer.GetLastBufferLine(); exists {
			ParseInput(lastLine, input)
		}
	}
	if utilDebouncedKeyPressed(ebiten.KeyBackspace) {
		CommandBuffer.DecrementCursorWithDecor(CmdBufferDecor)
	}

	CommandBuffer.DrawToGrid(GridIdScene3, CommandBufferX, CommandBufferY, gs)

	LogBuffer.DrawToGrid(GridIdScene3, LogBufferX, LogBufferY, gs)

	UpdateAnimationGrid(gs, anims)

	state := current

loop:
	for {
		select {
		case cmd := <-commands:
			fmt.Printf("Commands: %v\n", cmd)

			switch cmd {
			case AtomConnectTrue:
				fmt.Printf("Connection Made!\n")
				state = next
			case AtomConnectFalse:
				// fmt.Printf("Connection Failed!\n")
				LogBuffer.AppendAll([]byte("Connection Failed"))
				LogBuffer.NewLine()
			case AtomInvalid:
				// fmt.Printf("Invalid!\n")
				LogBuffer.AppendAll([]byte("Invalid Command"))
				LogBuffer.NewLine()
			}
		default:
			break loop // nothing left in the queue for this frame
		}
	}

	return state
}

func ParseInput(input []byte, parserInput chan []byte) {
	fmt.Printf("Input: %s ;; %s\n", input, input[len(input)-1])

	if input[len(input)-1] != '=' {
		// Command not entered
		return
	}

	fmt.Printf("Passed\n")

	// Allocate or copy a standalone slice for the background worker
	// to prevent data races on the local inputBuffer
	commandCopy := make([]byte, len(input))
	copy(commandCopy, input)
	// Ship the bytes off the render thread instantly
	select {
	case parserInput <- commandCopy:
	default:
		// Dropped if worker queue is completely choked
		fmt.Printf("Dropped\n")
	}

	// Placeholder for now until using prolog parsing
}

func Scene3_HandleCleaUp(next GameState, gs *GridSystem, anims *AnimationSystem) GameState {

	anims.IsPlaying[AnimationScannerGrid] = false
	anims.Loop[AnimationScannerGrid] = false

	gs.SetAllCells(OutputGridId, CellTypeEmpty, 0)
	gs.SetAllCells(AnimationScannerGrid, CellTypeEmpty, 0)
	gs.DisableGrid(AnimationScannerGrid)

	return next
}
