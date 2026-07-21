package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/trealla-prolog/go/trealla"
	"math"
	"techboot_reno/cmd/assets"
)

func Scene4_HandleInit(current, next GameState, gs *GridSystem, anims *AnimationSystem) GameState {
	LogBuffer.AppendAll([]byte("Connecting..."))
	LogBuffer.NewLine()
	LogBuffer.DrawToGrid(GridIdScene3, LogBufferX, LogBufferY, gs)

	anims.IsPlaying[AnimationMemoryStack] = true
	anims.Loop[AnimationMemoryStack] = false
	anims.Durations[AnimationMemoryStack] = 5.0
	anims.Timers[AnimationMemoryStack] = 0.0
	anims.Delay[AnimationMemoryStack] = 0

	return next
}

func Scene4_UpdateStackAnimation(current, next GameState, gs *GridSystem, anims *AnimationSystem) GameState {
	timer := anims.Timers[AnimationMemoryStack]
	duration := anims.Durations[AnimationMemoryStack]
	delay := anims.Delay[AnimationMemoryStack]

	trueTime := float32(math.Max(float64(timer-delay), 0))
	completedTime := trueTime / duration

	// gs.SetAllCells(AnimationScannerGrid, CellTypeEmpty, 0)

	x := OutputGridCols / 2
	y := int(float32(OutputGridRows-1) * completedTime)
	for i := 0; i < OutputGridRows; i++ {
		iy := OutputGridRows - 1 - y
		gs.SetCellSprite(OutputGridId, x-1, iy, assets.SpriteIDSquare)
		gs.SetCellSprite(OutputGridId, x, iy, assets.SpriteIDSquare)
		gs.SetCellSprite(OutputGridId, x+1, iy, assets.SpriteIDSquare)
	}

	if anims.IsPlaying[AnimationMemoryStack] {
		LogBuffer.DrawToGrid(GridIdScene3, LogBufferX, LogBufferY, gs)
		return current
	}

	LogBuffer.AppendAll([]byte("Connection Successful"))
	LogBuffer.NewLine()

	return next
}

func Scene4_Update(current, next GameState, runes []rune, input chan []byte, commands chan trealla.Atom, gs *GridSystem, anims *AnimationSystem) GameState {
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

			// switch cmd {
			// case AtomConnectTrue:
			// 	fmt.Printf("Connection Made!\n")
			// 	state = next
			// case AtomConnectFalse:
			// 	// fmt.Printf("Connection Failed!\n")
			// 	LogBuffer.AppendAll([]byte("Connection Failed"))
			// 	LogBuffer.NewLine()
			// case AtomInvalid:
			// 	// fmt.Printf("Invalid!\n")
			// 	LogBuffer.AppendAll([]byte("Invalid Command"))
			// 	LogBuffer.NewLine()
			// }
		default:
			break loop // nothing left in the queue for this frame
		}
	}

	return state
}
