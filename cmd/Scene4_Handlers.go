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

	game.gs.EnableGrid(game.b.GridOutput)

	game.bs.AppendAll(game.b.BufferLogs, []byte("Connecting..."))
	game.bs.NewLine(game.b.BufferLogs)
	game.bs.DrawToGrid(game.b.BufferLogs, game.b.GridMainUI, game.b.LogBufferColIdx, game.b.LogBufferRowIdx, game.gs)

	game.ans.IsPlaying[AnimationMemoryStack] = true
	game.ans.Loop[AnimationMemoryStack] = false
	game.ans.Durations[AnimationMemoryStack] = 5.0
	game.ans.Timers[AnimationMemoryStack] = 0.0
	game.ans.Delay[AnimationMemoryStack] = 0

	return next
}

func Scene4_UpdateStackAnimation(current, next GameState, game *Game) GameState {
	timer := game.ans.Timers[AnimationMemoryStack]
	duration := game.ans.Durations[AnimationMemoryStack]
	delay := game.ans.Delay[AnimationMemoryStack]

	trueTime := float32(math.Max(float64(timer-delay), 0))
	completedTime := trueTime / duration

	// gs.SetAllCells(AnimationScannerGrid, CellTypeEmpty, 0)
	OutputGridCols := game.gs.Cols[game.b.GridOutput]
	OutputGridRows := game.gs.Rows[game.b.GridOutput]

	x := OutputGridCols / 2
	y := int(float32(OutputGridRows-1) * completedTime)
	for i := 0; i < OutputGridRows; i++ {
		iy := OutputGridRows - 1 - y
		game.gs.SetCellSprite(game.b.GridOutput, x-1, iy, assets.SpriteIDSquare)
		game.gs.SetCellSprite(game.b.GridOutput, x, iy, assets.SpriteIDSquare)
		game.gs.SetCellSprite(game.b.GridOutput, x+1, iy, assets.SpriteIDSquare)
	}

	if game.ans.IsPlaying[AnimationMemoryStack] {
		game.bs.DrawToGrid(game.b.BufferLogs, game.b.GridMainUI, game.b.LogBufferColIdx, game.b.LogBufferRowIdx, game.gs)
		return current
	}

	// LogBuffer.AppendAll([]byte("Connection Successful"))
	game.bs.AppendAll(game.b.BufferLogs, []byte("Connection Successful"))
	game.bs.NewLine(game.b.BufferLogs)

	game.gs.SetAllCells(game.b.GridOutput, CellTypeEmpty, 0)

	return next
}

// func Scene4_Update(current, next GameState, runes []rune, input chan []byte, commands chan trealla.Atom, gs *GridSystem, anims *AnimationSystem) GameState {
func Scene4_Update(current, next GameState, game *Game) GameState {
	for i := 0; i < len(game.inputRunes); i++ {
		err := game.bs.AppendWithDecor(game.b.BufferCommands, byte(game.inputRunes[i]), CmdBufferDecor)
		if err != nil {
			fmt.Printf("Error Appending Last Rune: %v\n", err)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		game.bs.TrimDecor(game.b.BufferCommands, CmdBufferDecor)
		game.bs.NewLine(game.b.BufferCommands)
		game.bs.AppendDecorators(game.b.BufferCommands, CmdBufferDecor)

		if lastLine, lineError := game.bs.GetLastBufferLine(game.b.BufferCommands); lineError == nil {
			ParseInput(lastLine, game.prologInput)
		} else {
			fmt.Printf("Error Getting lastLine: %v", lineError)
		}
	}
	if utilDebouncedKeyPressed(ebiten.KeyBackspace) {
		game.bs.DecrementCursorWithDecor(game.b.BufferCommands, CmdBufferDecor)
	}

	err := game.bs.DrawToGridWithDecor(game.b.BufferCommands, game.b.GridMainUI, 1, 1, CmdBufferDecor, game.gs)
	if err != nil {
		fmt.Printf("Error Drawing to BufferCommands: %v", err)
	}

	// LogBuff
	game.bs.DrawToGrid(game.b.BufferLogs, game.b.GridMainUI, 37, 38, game.gs)

	state := current

loop:
	for {
		select {
		case cmd := <-game.prologOutput:
			fmt.Printf("Commands: %v\n", cmd)

			switch cmd.ResultType {
			case CommandList:
				for i := 0; i < len(cmd.Items); i++ {
					b := game.b.CommandBytes(cmd.Items[i], game.bs)
					fmt.Printf("%d: %d: %s\n", i, len(b), b)
				}
			case CommandInvalid:
				// fmt.Printf("Invalid!\n")
				game.bs.AppendAll(game.b.BufferLogs, []byte("Invalid Command"))
				game.bs.NewLine(game.b.BufferLogs)
			}
		default:
			break loop // nothing left in the queue for this frame
		}
	}

	return state
}
