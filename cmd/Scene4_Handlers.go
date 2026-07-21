package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/trealla-prolog/go/trealla"
	// "techboot_reno/cmd/assets"
)

func Scene4_HandleInit(current, next GameState, gs *GridSystem, anims *AnimationSystem) GameState {
	LogBuffer.AppendAll([]byte("Connecting..."))
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
