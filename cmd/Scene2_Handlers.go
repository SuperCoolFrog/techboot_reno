package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
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

func Scene2_WaitForEnter(current, next GameState, game *Game) GameState {

	game.GridSystem.Set(game.Bucket.GridDialogScene, 16, 15, CellTypeChar, '[')
	game.GridSystem.Set(game.Bucket.GridDialogScene, 17, 15, CellTypeChar, 'E')
	game.GridSystem.Set(game.Bucket.GridDialogScene, 18, 15, CellTypeChar, 'N')
	game.GridSystem.Set(game.Bucket.GridDialogScene, 19, 15, CellTypeChar, 'T')
	game.GridSystem.Set(game.Bucket.GridDialogScene, 20, 15, CellTypeChar, 'E')
	game.GridSystem.Set(game.Bucket.GridDialogScene, 21, 15, CellTypeChar, 'R')
	game.GridSystem.Set(game.Bucket.GridDialogScene, 22, 15, CellTypeChar, ']')

	if game.MouseMoved {
		mx, my := ebiten.CursorPosition()
		mouseRect := image.Rect(mx, my, mx+1, my+1)

		if mouseRect.Overlaps(game.GridSystem.GridRectangle(game.Bucket.GridDialogScene, 16, 15, 7, 1)) {
			game.GridSystem.Set(game.Bucket.GridDialogScene, 15, 15, CellTypeChar, ':')
			game.GridSystem.Set(game.Bucket.GridDialogScene, 23, 15, CellTypeChar, '|')
		} else {
			game.GridSystem.Set(game.Bucket.GridDialogScene, 15, 15, CellTypeEmpty, ' ')
			game.GridSystem.Set(game.Bucket.GridDialogScene, 23, 15, CellTypeEmpty, ' ')
		}
	}

	acceptSelected := inpututil.IsKeyJustPressed(ebiten.KeyEnter) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)

	if acceptSelected {
		return next
	}

	return current
}

func Scene2_CleanUpScene(next GameState, game *Game) GameState {
	game.GridSystem.DisableGrid(game.Bucket.GridDialogScene)
	return next
}
