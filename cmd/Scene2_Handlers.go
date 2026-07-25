package main

import (
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

	// Set iterator for dialog animations
	game.Bucket.SceneStateItr[next] = 0

	game.GridSystem.EnableGrid(game.Bucket.GridDialogScene)

	// s2_AddRipMsgBuffer(game.GridSystem, game.Bucket.GridDialogScene, 0)

	// PlayDialogAnimation(game.GridSystem, game.Animations)

	// bucket.GridDialogScene = bucket.Grid40x30x32x16
	// bucket.BufferDialogScene = bucket.Buffer40x30x30xfalse

	txt := DialogText[0]

	for i := 0; i < len(txt); i++ {
		game.Buffers.AppendWithDecor(game.Bucket.BufferDialogScene, txt[i], DecorMap[0])
	}

	game.Buffers.DrawToGrid(game.Bucket.BufferDialogScene, game.Bucket.GridDialogScene, 0, 0, game.GridSystem)

	return next
}

func Scene2_HandleAllDialog(current, next GameState, game *Game) GameState {

	// return current //
	return next
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
