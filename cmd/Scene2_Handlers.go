package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const Scene2SpeakerR = 0
const Scene2SpeakerT = 1
const Scene2StateYIdx = 0

var (
	Scene2_DialogIdx     = [6]int{0, 1, 2, 3, 4, 5}
	Scene2_DialogSpeaker = [6]int{
		Scene2SpeakerR,
		Scene2SpeakerR,
		Scene2SpeakerR,
		Scene2SpeakerR,
		Scene2SpeakerR,
		Scene2SpeakerT,
	}

	RipleyDecor = BufferDecorator{
		Prefix:  []byte("`RIPL3Y`:"),
		Postfix: []byte{},
	}
	TRenoDecor = BufferDecorator{
		Prefix:  []byte("`T.REN0`:"),
		Postfix: []byte{},
	}

	DecorMap = []BufferDecorator{
		RipleyDecor,
		TRenoDecor,
	}
)

func Scene2_HandleInit(current, next GameState, game *Game) GameState {
	// Clean Up Previous Grid
	game.GridSystem.DisableGrid(game.Bucket.GridStartScene)
	game.GridSystem.Set(game.Bucket.GridStartScene, 0, 1, CellTypeEmpty, ' ')
	game.GridSystem.Set(game.Bucket.GridStartScene, 1, 1, CellTypeEmpty, ' ')
	// /Clean Up Previous Grid

	game.GridSystem.EnableGrid(game.Bucket.GridDialogScene)

	// Set iterator for dialog animations
	game.Bucket.SceneStateItr[next] = 0

	game.Animations.Delay[AnimationDialog] = 1.0
	game.Animations.Timers[AnimationDialog] = 0.0
	game.Animations.Durations[AnimationDialog] = 3.0
	game.Animations.IsPlaying[AnimationDialog] = true

	return next
}

func Scene2_HandleAllDialog(current, next GameState, game *Game) GameState {
	txt := DialogText[game.Bucket.SceneStateItr[current]]
	speaker := Scene2_DialogSpeaker[game.Bucket.SceneStateItr[current]]
	decor := DecorMap[speaker]

	delay := game.Animations.Delay[AnimationDialog]
	completed := game.Animations.Timers[AnimationDialog]
	duration := game.Animations.Durations[AnimationDialog]

	ratio := (completed - delay) / duration

	if !game.Animations.IsPlaying[AnimationDialog] {
		game.Bucket.SceneStateItr[current]++
		game.Buffers.NewLine(game.Bucket.BufferDialogScene)

		if game.Bucket.SceneStateItr[current] >= len(Scene2_DialogIdx) {
			return next
		}

		game.Animations.Timers[AnimationDialog] = 0.0
		game.Animations.Durations[AnimationDialog] = 3.0
		game.Animations.IsPlaying[AnimationDialog] = true

		return current
	}

	xCursor := game.Buffers.GetXCursor(game.Bucket.BufferDialogScene)

	if xCursor == 0 {
		game.Buffers.AppendDecorators(game.Bucket.BufferDialogScene, decor)
		xCursor = game.Buffers.GetXCursor(game.Bucket.BufferDialogScene)
	}

	prefixLen := len(decor.Prefix)
	appendedByteSize := float32(xCursor - prefixLen - 1)
	totalTxtSize := float32(len(txt))

	var appendError error

	for (appendedByteSize / totalTxtSize) < ratio {
		appendError = game.Buffers.AppendWithDecor(game.Bucket.BufferDialogScene, txt[int(appendedByteSize)], decor)
		if appendError != nil {
			panic(fmt.Sprintf("ErrorAppending: %v", appendError))
			break
		}

		xCursor = game.Buffers.GetXCursor(game.Bucket.BufferDialogScene)
		appendedByteSize = float32(xCursor - prefixLen)
	}

	drawError := game.Buffers.DrawToGrid(game.Bucket.BufferDialogScene, game.Bucket.GridDialogScene, 0, 0, game.GridSystem)
	if drawError != nil {
		panic(drawError)
	}

	return current
}

func Scene2_WaitForEnter(current, next GameState) GameState {
	acceptSelected := inpututil.IsKeyJustPressed(ebiten.KeyEnter) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)

	if acceptSelected {
		return next
	}

	return current
}

func Scene2_CleanUpScene(next GameState, gs *GridSystem, anims *AnimationSystem) GameState {
	dialogGridId := anims.GridId[AnimationDialog]
	gs.DisableGrid(dialogGridId)
	// @TODO this also sets the buffers CellTypes.  Need to consider how to handle that
	// gs.SetAllCells(dialogGridId, CellTypeEmpty, ' ')

	return next
}
