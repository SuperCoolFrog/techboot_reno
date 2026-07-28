package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"math"
	"techboot_reno/cmd/assets"
)

var (
	OutputBuffer *Buffer
)

func Scene4_HandleInit(current, next GameState, game *Game) GameState {

	game.GridSystem.EnableGrid(game.Bucket.GridOutput)

	game.Buffers.AppendAll(game.Bucket.BufferLogs, []byte("Connecting..."))
	game.Buffers.NewLine(game.Bucket.BufferLogs)
	game.Buffers.DrawToGrid(game.Bucket.BufferLogs, game.Bucket.GridMainUI, game.Bucket.LogBufferColIdx, game.Bucket.LogBufferRowIdx, game.GridSystem)

	game.Animations.IsPlaying[AnimationMemoryStack] = true
	game.Animations.Loop[AnimationMemoryStack] = false
	game.Animations.Durations[AnimationMemoryStack] = 5.0
	game.Animations.Timers[AnimationMemoryStack] = 0.0
	game.Animations.Delay[AnimationMemoryStack] = 0

	return next
}

func Scene4_UpdateStackAnimation(current, next GameState, game *Game) GameState {
	timer := game.Animations.Timers[AnimationMemoryStack]
	duration := game.Animations.Durations[AnimationMemoryStack]
	delay := game.Animations.Delay[AnimationMemoryStack]

	trueTime := float32(math.Max(float64(timer-delay), 0))
	completedTime := trueTime / duration

	// gs.SetAllCells(AnimationScannerGrid, CellTypeEmpty, 0)
	OutputGridCols := game.GridSystem.Cols[game.Bucket.GridOutput]
	OutputGridRows := game.GridSystem.Rows[game.Bucket.GridOutput]

	x := OutputGridCols / 2
	y := int(float32(OutputGridRows-1) * completedTime)
	for i := 0; i < OutputGridRows; i++ {
		iy := OutputGridRows - 1 - y
		game.GridSystem.SetCellSprite(game.Bucket.GridOutput, x-1, iy, assets.SpriteIDSquare)
		game.GridSystem.SetCellSprite(game.Bucket.GridOutput, x, iy, assets.SpriteIDSquare)
		game.GridSystem.SetCellSprite(game.Bucket.GridOutput, x+1, iy, assets.SpriteIDSquare)
	}

	if game.Animations.IsPlaying[AnimationMemoryStack] {
		game.Buffers.DrawToGrid(game.Bucket.BufferLogs, game.Bucket.GridMainUI, game.Bucket.LogBufferColIdx, game.Bucket.LogBufferRowIdx, game.GridSystem)
		return current
	}

	// LogBuffer.AppendAll([]byte("Connection Successful"))
	game.Buffers.AppendAll(game.Bucket.BufferLogs, []byte("Connection Successful"))
	game.Buffers.NewLine(game.Bucket.BufferLogs)

	game.GridSystem.SetAllCells(game.Bucket.GridOutput, CellTypeEmpty, 0)

	return next
}

// func Scene4_Update(current, next GameState, runes []rune, input chan []byte, commands chan trealla.Atom, gs *GridSystem, anims *AnimationSystem) GameState {
func Scene4_Update(current, next GameState, game *Game) GameState {
	for i := 0; i < len(game.inputRunes); i++ {
		game.Buffers.AppendWithDecor(game.Bucket.BufferCommands, byte(game.inputRunes[i]), CmdBufferDecor)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		//CommandBuffer
		game.Buffers.TrimDecor(game.Bucket.BufferCommands, CmdBufferDecor)
		game.Buffers.NewLine(game.Bucket.BufferCommands)
		game.Buffers.AppendDecorators(game.Bucket.BufferCommands, CmdBufferDecor)

		fmt.Printf("Enter\n")
		if lastLine, lineError := game.Buffers.GetLastBufferLine(game.Bucket.BufferCommands); lineError == nil {
			ParseInput(lastLine, game.prologInput)
		}
	}
	if utilDebouncedKeyPressed(ebiten.KeyBackspace) {
		game.Buffers.DecrementCursorWithDecor(game.Bucket.BufferCommands, CmdBufferDecor)
	}

	game.Buffers.DrawToGrid(game.Bucket.BufferCommands, game.Bucket.GridMainUI, 1, 1, game.GridSystem)

	game.Buffers.DrawToGrid(game.Bucket.BufferLogs, game.Bucket.GridMainUI, game.Bucket.LogBufferColIdx, game.Bucket.LogBufferRowIdx, game.GridSystem)

	state := current

loop:
	for {
		select {
		case cmd := <-game.prologOutput:
			fmt.Printf("Commands: %v\n", cmd)

			switch cmd {
			case AtomList:
				// Display list items
			case AtomInvalid:
				// fmt.Printf("Invalid!\n")
				game.Buffers.AppendAll(game.Bucket.BufferLogs, []byte("Invalid Command"))
				game.Buffers.NewLine(game.Bucket.BufferLogs)
			}
		default:
			break loop // nothing left in the queue for this frame
		}
	}

	return state
}
