package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	// "github.com/trealla-prolog/go/trealla"
	"math"
	"techboot_reno/cmd/assets"
)

var (
	CmdBufferDecor = BufferDecorator{
		Prefix:  []byte{':'},
		Postfix: []byte{'|'},
	}
)

func Scene3_HandleInit(current, next GameState, game *Game) GameState {
	game.gs.EnableGrid(game.b.GridMainUI)

	game.gs.SetAllCells(game.b.GridOutput, CellTypeSquare, 0)
	game.gs.EnableGrid(game.b.GridOutput)

	game.gs.SetAllCells(game.b.GridOutputPrecision, CellTypeEmpty, 0)
	game.gs.EnableGrid(game.b.GridOutputPrecision)

	// Border

	game.bs.AppendDecorators(game.b.BufferCommands, CmdBufferDecor)

	game.bs.AppendAll(game.b.BufferLogs, []byte("Type: connect rabbit="))
	game.bs.NewLine(game.b.BufferLogs)

	game.ans.IsPlaying[AnimationScanner] = true
	game.ans.Loop[AnimationScanner] = true

	return next
}

func Scene3_Update(current, next GameState, game *Game) GameState {

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

	scene3UpdateAnimationGrid(game)

	state := current

loop:
	for {
		select {
		case cmd := <-game.prologOutput:
			fmt.Printf("Commands: %v\n", cmd)

			switch cmd.ResultType {
			case CommandConnectTrue:
				fmt.Printf("Connection Made!\n")
				state = next
			case CommandConnectFalse:
				// fmt.Printf("Connection Failed!\n")
				game.bs.AppendAll(game.b.BufferLogs, []byte("Connection Failed"))
				game.bs.NewLine(game.b.BufferLogs)
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

func Scene3_HandleCleaUp(next GameState, game *Game) GameState {
	game.ans.IsPlaying[AnimationScanner] = false
	game.ans.Loop[AnimationScanner] = false

	game.gs.SetAllCells(game.b.GridOutput, CellTypeEmpty, 0)
	game.gs.DisableGrid(game.b.GridOutput)
	game.gs.DisableGrid(game.b.GridOutputPrecision)

	return next
}

func scene3UpdateAnimationGrid(game *Game) {
	timer := game.ans.Timers[AnimationScanner]
	duration := game.ans.Durations[AnimationScanner]
	delay := game.ans.Delay[AnimationScanner]

	trueTime := float32(math.Max(float64(timer-delay), 0))
	completedTime := trueTime / duration

	game.gs.SetAllCells(game.b.GridOutputPrecision, CellTypeEmpty, 0)

	scannerGridRows := game.gs.Rows[game.b.GridOutputPrecision]
	scannerGridCols := game.gs.Cols[game.b.GridOutputPrecision]

	if completedTime <= 0.5 {
		y := int(float32(scannerGridRows-1) * (completedTime / .5))
		for i := 0; i < scannerGridCols; i++ {
			game.gs.SetCellSprite(game.b.GridOutputPrecision, i, y, assets.SpriteIDHorizontalBar)
		}
	} else {
		y := int(float32(scannerGridRows-1) * ((completedTime - .5) / .5))
		for i := 0; i < scannerGridCols; i++ {
			game.gs.SetCellSprite(game.b.GridOutputPrecision, i, scannerGridRows-1-y, assets.SpriteIDHorizontalBar)
		}
	}
}
