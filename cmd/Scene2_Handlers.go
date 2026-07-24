package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func s2_AddRipMsgBuffer(gs *GridSystem, gridId GridID, y int) {
	gs.Set(gridId, 0, y, CellTypeChar, '`')
	gs.Set(gridId, 1, y, CellTypeChar, 'R')
	gs.Set(gridId, 2, y, CellTypeChar, 'I')
	gs.Set(gridId, 3, y, CellTypeChar, 'P')
	gs.Set(gridId, 4, y, CellTypeChar, 'L')
	gs.Set(gridId, 5, y, CellTypeChar, '3')
	gs.Set(gridId, 6, y, CellTypeChar, 'Y')
	gs.Set(gridId, 7, y, CellTypeChar, '`')
	gs.Set(gridId, 8, y, CellTypeChar, ':')

	//Setup Buffer
	gs.NewBuffer(gridId, 9, y, gs.Cols[gridId]-10, 1)
}

func s2_AddTRenoMsgBuffer(gs *GridSystem, gridId GridID, y int) {
	gs.Set(gridId, 0, y, CellTypeChar, '`')
	gs.Set(gridId, 1, y, CellTypeChar, 'T')
	gs.Set(gridId, 2, y, CellTypeChar, '.')
	gs.Set(gridId, 3, y, CellTypeChar, 'R')
	gs.Set(gridId, 4, y, CellTypeChar, 'E')
	gs.Set(gridId, 5, y, CellTypeChar, 'N')
	gs.Set(gridId, 6, y, CellTypeChar, '0')
	gs.Set(gridId, 7, y, CellTypeChar, '`')
	gs.Set(gridId, 8, y, CellTypeChar, ':')

	//Setup Buffer
	gs.NewBuffer(gridId, 9, y, gs.Cols[gridId]-10, 1)
}

func Scene2_HandleInit(gridSystem *GridSystem, anims *AnimationSystem) {
	animationGridId := anims.GridId[AnimationStartScene]
	gridSystem.DisableGrid(animationGridId)
	gridSystem.Set(animationGridId, 0, 1, CellTypeEmpty, ' ')
	gridSystem.Set(animationGridId, 1, 1, CellTypeEmpty, ' ')

	InitDialogAnimation(gridSystem, anims)

	gridId := anims.GridId[AnimationDialog]

	s2_AddRipMsgBuffer(gridSystem, gridId, 0)

	// Init Animations

	// 'Ripl3y':"I think you're right"|
	// 'Ripl3y':"I don't think she ran"|
	// 'Ripl3y':"I found an open door"|
	// 'Ripl3y':"CONNECT to RABBIT and you're in"|
	// 'Ripl3y':"Good Luck"|

}

func Scene2_HandleDialog1(gs *GridSystem, anims *AnimationSystem) GameState {
	gridId := anims.GridId[AnimationDialog]

	buffer1Idx := gs.IdxFromXY(gridId, 9, 0)
	message := []byte("I think you're right") // maybe to []byte{'I', ...} instead.

	UpdateDialogAnimation(gridId, buffer1Idx, message, gs, anims)

	b1Cursor := gs.BufferCursor[buffer1Idx]

	if anims.IsPlaying[AnimationDialog] {
		gs.Set(gridId, 9+b1Cursor, 0, CellTypeReserved, '|')
	} else {
		gs.Set(gridId, 9+b1Cursor, 0, CellTypeReserved, ' ')
		return Scene2_Init_Dialog_2
	}

	return Scene2_Dialog_1
}

func Scene2_InitDialog(bufferY int, gs *GridSystem, anims *AnimationSystem) {
	gridId := anims.GridId[AnimationDialog]
	s2_AddRipMsgBuffer(gs, gridId, bufferY)
	PlayDialogAnimation(2.0, false, gs, anims)
}

func Scene2_HandleDialog(bufferY int, message []byte, currentState, nextState GameState, gs *GridSystem, anims *AnimationSystem) GameState {
	gridId := anims.GridId[AnimationDialog]

	buffer1Idx := gs.IdxFromXY(gridId, 9, bufferY)
	//message := []byte("I don't this she ran...")

	UpdateDialogAnimation(gridId, buffer1Idx, message, gs, anims)

	b1Cursor := gs.BufferCursor[buffer1Idx]

	if anims.IsPlaying[AnimationDialog] {
		gs.Set(gridId, 9+b1Cursor, bufferY, CellTypeReserved, '|')
	} else {
		gs.Set(gridId, 9+b1Cursor, bufferY, CellTypeReserved, ' ')
		return nextState
	}

	return currentState
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
