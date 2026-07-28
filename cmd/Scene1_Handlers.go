package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
)

func Scene1_HandleAnimationComplete(game *Game) {
	game.gs.EnableGrid(game.b.GridStartScene)
	game.gs.Set(game.b.GridStartScene, 9, 3, CellTypeChar, 'T')
	game.gs.Set(game.b.GridStartScene, 10, 3, CellTypeChar, 'E')
	game.gs.Set(game.b.GridStartScene, 11, 3, CellTypeChar, 'C')
	game.gs.Set(game.b.GridStartScene, 12, 3, CellTypeChar, 'H')
	game.gs.Set(game.b.GridStartScene, 13, 3, CellTypeChar, 'B')
	game.gs.Set(game.b.GridStartScene, 14, 3, CellTypeChar, '0')
	game.gs.Set(game.b.GridStartScene, 15, 3, CellTypeChar, '0')
	game.gs.Set(game.b.GridStartScene, 16, 3, CellTypeChar, 'T')

	game.gs.Set(game.b.GridStartScene, 11, 4, CellTypeChar, 'R')
	game.gs.Set(game.b.GridStartScene, 12, 4, CellTypeChar, 'E')
	game.gs.Set(game.b.GridStartScene, 13, 4, CellTypeChar, 'N')
	game.gs.Set(game.b.GridStartScene, 14, 4, CellTypeChar, '0')
}

func Scene1_HandleButtonList(current, next GameState, game *Game) GameState {
	gridId := game.b.GridStartScene

	startFocusCellType, _ := game.gs.Get(gridId, 7, 8)
	exitFocusCellType, _ := game.gs.Get(gridId, 7, 9)

	startSelected := startFocusCellType != CellTypeEmpty
	exitSelected := exitFocusCellType != CellTypeEmpty

	if !(startSelected || exitSelected) {
		// Neither Selected so select start
		// start
		game.gs.Set(gridId, 7, 8, CellTypeChar, ':')
		game.gs.Set(gridId, 15, 8, CellTypeChar, '|')

		// exit
		game.gs.Set(gridId, 7, 9, CellTypeEmpty, ' ')
		game.gs.Set(gridId, 15, 9, CellTypeEmpty, ' ')

		startSelected = true
		exitSelected = false
	}

	downKeyPressed := inpututil.IsKeyJustPressed(ebiten.KeyDown) ||
		inpututil.IsKeyJustPressed(ebiten.KeyJ) ||
		inpututil.IsKeyJustPressed(ebiten.KeyS)

	upKeyPressed := inpututil.IsKeyJustPressed(ebiten.KeyUp) ||
		inpututil.IsKeyJustPressed(ebiten.KeyK) ||
		inpututil.IsKeyJustPressed(ebiten.KeyW)

	if downKeyPressed && startSelected {
		// start
		game.gs.Set(gridId, 7, 8, CellTypeEmpty, ' ')
		game.gs.Set(gridId, 15, 8, CellTypeEmpty, ' ')

		// exit
		game.gs.Set(gridId, 7, 9, CellTypeChar, ':')
		game.gs.Set(gridId, 15, 9, CellTypeChar, '|')

		startSelected = false
		exitSelected = true
	}

	if upKeyPressed && exitSelected {
		// start
		game.gs.Set(gridId, 7, 8, CellTypeChar, ':')
		game.gs.Set(gridId, 15, 8, CellTypeChar, '|')

		// exit
		game.gs.Set(gridId, 7, 9, CellTypeEmpty, ' ')
		game.gs.Set(gridId, 15, 9, CellTypeEmpty, ' ')

		startSelected = true
		exitSelected = false
	}

	// Mouse
	mx, my := ebiten.CursorPosition()
	mouseRect := image.Rect(mx, my, mx+1, my+1)

	// Start
	if game.MouseMoved && mouseRect.Overlaps(game.gs.GridRectangle(gridId, 8, 8, 7, 1)) {
		// start
		game.gs.Set(gridId, 7, 8, CellTypeChar, ':')
		game.gs.Set(gridId, 15, 8, CellTypeChar, '|')

		// exit
		game.gs.Set(gridId, 7, 9, CellTypeEmpty, ' ')
		game.gs.Set(gridId, 15, 9, CellTypeEmpty, ' ')

		startSelected = true
		exitSelected = false
	}

	// Exit
	if game.MouseMoved && mouseRect.Overlaps(game.gs.GridRectangle(gridId, 8, 9, 7, 1)) {
		// start
		game.gs.Set(gridId, 7, 8, CellTypeEmpty, ' ')
		game.gs.Set(gridId, 15, 8, CellTypeEmpty, ' ')

		// exit
		game.gs.Set(gridId, 7, 9, CellTypeChar, ':')
		game.gs.Set(gridId, 15, 9, CellTypeChar, '|')

		startSelected = false
		exitSelected = true
	}

	acceptSelected := inpututil.IsKeyJustPressed(ebiten.KeyEnter) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)

	if acceptSelected && startSelected {
		return next
	}

	if acceptSelected && exitSelected {
		game.Exit = true
		return current
	}

	game.gs.Set(gridId, 8, 8, CellTypeChar, '[')
	game.gs.Set(gridId, 9, 8, CellTypeChar, 'S')
	game.gs.Set(gridId, 10, 8, CellTypeChar, 'T')
	game.gs.Set(gridId, 11, 8, CellTypeChar, 'A')
	game.gs.Set(gridId, 12, 8, CellTypeChar, 'R')
	game.gs.Set(gridId, 13, 8, CellTypeChar, 'T')
	game.gs.Set(gridId, 14, 8, CellTypeChar, ']')

	game.gs.Set(gridId, 8, 9, CellTypeChar, '[')
	game.gs.Set(gridId, 9, 9, CellTypeChar, 'E')
	game.gs.Set(gridId, 10, 9, CellTypeChar, 'X')
	game.gs.Set(gridId, 11, 9, CellTypeChar, 'I')
	game.gs.Set(gridId, 12, 9, CellTypeChar, 'T')
	game.gs.Set(gridId, 13, 9, CellTypeChar, ' ')
	game.gs.Set(gridId, 14, 9, CellTypeChar, ']')

	return current
}

func Scene1_HandleExit(game *Game) {}
