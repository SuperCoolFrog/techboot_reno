package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
)

func Scene1_HandleAnimationComplete(game *Game) {
	gridId := game.Animations.GridId[AnimationIntroGrid]
	game.State = Scene1_Waiting
	game.Animations.Offset = 1 // Intro animation is 0
	game.GridSystem.EnableGrid(gridId)
	game.GridSystem.Set(gridId, 9, 3, CellTypeChar, 'T')
	game.GridSystem.Set(gridId, 10, 3, CellTypeChar, 'E')
	game.GridSystem.Set(gridId, 11, 3, CellTypeChar, 'C')
	game.GridSystem.Set(gridId, 12, 3, CellTypeChar, 'H')
	game.GridSystem.Set(gridId, 13, 3, CellTypeChar, 'B')
	game.GridSystem.Set(gridId, 14, 3, CellTypeChar, '0')
	game.GridSystem.Set(gridId, 15, 3, CellTypeChar, '0')
	game.GridSystem.Set(gridId, 16, 3, CellTypeChar, 'T')

	game.GridSystem.Set(gridId, 11, 4, CellTypeChar, 'R')
	game.GridSystem.Set(gridId, 12, 4, CellTypeChar, 'E')
	game.GridSystem.Set(gridId, 13, 4, CellTypeChar, 'N')
	game.GridSystem.Set(gridId, 14, 4, CellTypeChar, '0')
}

func Scene1_HandleButtonList(game *Game) {
	gridId := game.Animations.GridId[AnimationIntroGrid]

	startFocusCellType, _ := game.GridSystem.Get(gridId, 7, 8)
	exitFocusCellType, _ := game.GridSystem.Get(gridId, 7, 9)

	if startFocusCellType == CellTypeEmpty && exitFocusCellType == CellTypeEmpty {
		// Neither Selected so select start
		// start
		game.GridSystem.Set(gridId, 7, 8, CellTypeChar, ':')
		game.GridSystem.Set(gridId, 15, 8, CellTypeChar, '|')

		startFocusCellType = CellTypeChar

		// exit
		game.GridSystem.Set(gridId, 7, 9, CellTypeEmpty, ' ')
		game.GridSystem.Set(gridId, 15, 9, CellTypeEmpty, ' ')
	}

	downKeyPressed := inpututil.IsKeyJustPressed(ebiten.KeyDown) ||
		inpututil.IsKeyJustPressed(ebiten.KeyJ) ||
		inpututil.IsKeyJustPressed(ebiten.KeyS)

	upKeyPressed := inpututil.IsKeyJustPressed(ebiten.KeyUp) ||
		inpututil.IsKeyJustPressed(ebiten.KeyK) ||
		inpututil.IsKeyJustPressed(ebiten.KeyW)

	if downKeyPressed && startFocusCellType != CellTypeEmpty {
		// start
		game.GridSystem.Set(gridId, 7, 8, CellTypeEmpty, ' ')
		game.GridSystem.Set(gridId, 15, 8, CellTypeEmpty, ' ')

		startFocusCellType = CellTypeEmpty

		// exit
		game.GridSystem.Set(gridId, 7, 9, CellTypeChar, ':')
		game.GridSystem.Set(gridId, 15, 9, CellTypeChar, '|')
		exitFocusCellType = CellTypeChar
	}

	if upKeyPressed && exitFocusCellType != CellTypeEmpty {
		// start
		game.GridSystem.Set(gridId, 7, 8, CellTypeChar, ':')
		game.GridSystem.Set(gridId, 15, 8, CellTypeChar, '|')
		startFocusCellType = CellTypeChar

		// exit
		game.GridSystem.Set(gridId, 7, 9, CellTypeEmpty, ' ')
		game.GridSystem.Set(gridId, 15, 9, CellTypeEmpty, ' ')
		exitFocusCellType = CellTypeEmpty
	}

	// keyPressed := upKeyPressed || downKeyPressed

	// Mouse
	mx, my := ebiten.CursorPosition()
	mouseRect := image.Rect(mx, my, mx+1, my+1)

	// Start
	if game.MouseMoved && mouseRect.Overlaps(game.GridSystem.GridRectangle(gridId, 8, 8, 7, 1)) {
		// start
		game.GridSystem.Set(gridId, 7, 8, CellTypeChar, ':')
		game.GridSystem.Set(gridId, 15, 8, CellTypeChar, '|')
		startFocusCellType = CellTypeChar

		// exit
		game.GridSystem.Set(gridId, 7, 9, CellTypeEmpty, ' ')
		game.GridSystem.Set(gridId, 15, 9, CellTypeEmpty, ' ')
		exitFocusCellType = CellTypeEmpty
	}

	// Exit
	if game.MouseMoved && mouseRect.Overlaps(game.GridSystem.GridRectangle(gridId, 8, 9, 7, 1)) {
		// start
		game.GridSystem.Set(gridId, 7, 8, CellTypeEmpty, ' ')
		game.GridSystem.Set(gridId, 15, 8, CellTypeEmpty, ' ')
		startFocusCellType = CellTypeEmpty

		// exit
		game.GridSystem.Set(gridId, 7, 9, CellTypeChar, ':')
		game.GridSystem.Set(gridId, 15, 9, CellTypeChar, '|')
		exitFocusCellType = CellTypeChar
	}

	game.GridSystem.Set(gridId, 8, 8, CellTypeChar, '[')
	game.GridSystem.Set(gridId, 9, 8, CellTypeChar, 'S')
	game.GridSystem.Set(gridId, 10, 8, CellTypeChar, 'T')
	game.GridSystem.Set(gridId, 11, 8, CellTypeChar, 'A')
	game.GridSystem.Set(gridId, 12, 8, CellTypeChar, 'R')
	game.GridSystem.Set(gridId, 13, 8, CellTypeChar, 'T')
	game.GridSystem.Set(gridId, 14, 8, CellTypeChar, ']')

	game.GridSystem.Set(gridId, 8, 9, CellTypeChar, '[')
	game.GridSystem.Set(gridId, 9, 9, CellTypeChar, 'E')
	game.GridSystem.Set(gridId, 10, 9, CellTypeChar, 'X')
	game.GridSystem.Set(gridId, 11, 9, CellTypeChar, 'I')
	game.GridSystem.Set(gridId, 12, 9, CellTypeChar, 'T')
	game.GridSystem.Set(gridId, 13, 9, CellTypeChar, ' ')
	game.GridSystem.Set(gridId, 14, 9, CellTypeChar, ']')
}
