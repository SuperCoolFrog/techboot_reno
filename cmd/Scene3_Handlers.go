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
	game.GridSystem.EnableGrid(game.Bucket.GridMainUI)

	game.GridSystem.SetAllCells(game.Bucket.GridOutput, CellTypeSquare, 0)
	game.GridSystem.EnableGrid(game.Bucket.GridOutput)

	game.GridSystem.SetAllCells(game.Bucket.GridOutputPrecision, CellTypeEmpty, 0)
	game.GridSystem.EnableGrid(game.Bucket.GridOutputPrecision)

	// Border

	game.Buffers.AppendDecorators(game.Bucket.BufferCommands, CmdBufferDecor)

	game.Buffers.AppendAll(game.Bucket.BufferLogs, []byte("Type: connect rabbit="))
	game.Buffers.NewLine(game.Bucket.BufferLogs)

	game.Animations.IsPlaying[AnimationScanner] = true
	game.Animations.Loop[AnimationScanner] = true

	return next
}

func Scene3_Update(current, next GameState, game *Game) GameState {

	for i := 0; i < len(game.inputRunes); i++ {
		err := game.Buffers.AppendWithDecor(game.Bucket.BufferCommands, byte(game.inputRunes[i]), CmdBufferDecor)
		if err != nil {
			fmt.Printf("Error Appending Last Rune: %v\n", err)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		game.Buffers.TrimDecor(game.Bucket.BufferCommands, CmdBufferDecor)
		game.Buffers.NewLine(game.Bucket.BufferCommands)
		game.Buffers.AppendDecorators(game.Bucket.BufferCommands, CmdBufferDecor)

		if lastLine, lineError := game.Buffers.GetLastBufferLine(game.Bucket.BufferCommands); lineError == nil {
			ParseInput(lastLine, game.prologInput)
		} else {
			fmt.Printf("Error Getting lastLine: %v", lineError)
		}
	}
	if utilDebouncedKeyPressed(ebiten.KeyBackspace) {
		game.Buffers.DecrementCursorWithDecor(game.Bucket.BufferCommands, CmdBufferDecor)
	}

	err := game.Buffers.DrawToGridWithDecor(game.Bucket.BufferCommands, game.Bucket.GridMainUI, 1, 1, CmdBufferDecor, game.GridSystem)
	if err != nil {
		fmt.Printf("Error Drawing to BufferCommands: %v", err)
	}

	// LogBuff
	game.Buffers.DrawToGrid(game.Bucket.BufferLogs, game.Bucket.GridMainUI, 37, 38, game.GridSystem)

	scene3UpdateAnimationGrid(game)

	state := current

loop:
	for {
		select {
		case cmd := <-game.prologOutput:
			fmt.Printf("Commands: %v\n", cmd)

			switch cmd {
			case AtomConnectTrue:
				fmt.Printf("Connection Made!\n")
				state = next
			case AtomConnectFalse:
				// fmt.Printf("Connection Failed!\n")
				game.Buffers.AppendAll(game.Bucket.BufferLogs, []byte("Connection Failed"))
				game.Buffers.NewLine(game.Bucket.BufferLogs)
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

	game.Animations.IsPlaying[AnimationScanner] = false
	game.Animations.Loop[AnimationScanner] = false

	game.GridSystem.SetAllCells(game.Bucket.GridOutput, CellTypeEmpty, 0)
	game.GridSystem.DisableGrid(game.Bucket.GridOutput)

	return next
}

func scene3UpdateAnimationGrid(game *Game) {
	timer := game.Animations.Timers[AnimationScanner]
	duration := game.Animations.Durations[AnimationScanner]
	delay := game.Animations.Delay[AnimationScanner]

	trueTime := float32(math.Max(float64(timer-delay), 0))
	completedTime := trueTime / duration

	game.GridSystem.SetAllCells(game.Bucket.GridOutputPrecision, CellTypeEmpty, 0)

	scannerGridRows := game.GridSystem.Rows[game.Bucket.GridOutputPrecision]
	scannerGridCols := game.GridSystem.Cols[game.Bucket.GridOutputPrecision]

	if completedTime <= 0.5 {
		y := int(float32(scannerGridRows-1) * (completedTime / .5))
		for i := 0; i < scannerGridCols; i++ {
			game.GridSystem.SetCellSprite(game.Bucket.GridOutputPrecision, i, y, assets.SpriteIDHorizontalBar)
		}
	} else {
		y := int(float32(scannerGridRows-1) * ((completedTime - .5) / .5))
		for i := 0; i < scannerGridCols; i++ {
			game.GridSystem.SetCellSprite(game.Bucket.GridOutputPrecision, i, scannerGridRows-1-y, assets.SpriteIDHorizontalBar)
		}
	}
}
