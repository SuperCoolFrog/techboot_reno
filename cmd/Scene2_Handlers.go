package main

func Scene2_HandleInit(gridSystem *GridSystem, anims *AnimationSystem) {
	animationGridId := anims.GridId[AnimationGrid]
	gridSystem.DisableGrid(animationGridId)
	gridSystem.Set(animationGridId, 0, 1, CellTypeEmpty, ' ')
	gridSystem.Set(animationGridId, 1, 1, CellTypeEmpty, ' ')

	InitDialogAnimation(gridSystem, anims)

	gridId := anims.GridId[AnimationDialog]

	gridSystem.Set(gridId, 0, 0, CellTypeChar, '`')
	gridSystem.Set(gridId, 1, 0, CellTypeChar, 'R')
	gridSystem.Set(gridId, 2, 0, CellTypeChar, 'I')
	gridSystem.Set(gridId, 3, 0, CellTypeChar, 'P')
	gridSystem.Set(gridId, 4, 0, CellTypeChar, 'L')
	gridSystem.Set(gridId, 5, 0, CellTypeChar, '3')
	gridSystem.Set(gridId, 6, 0, CellTypeChar, 'Y')
	gridSystem.Set(gridId, 7, 0, CellTypeChar, '`')
	gridSystem.Set(gridId, 8, 0, CellTypeChar, ':')
	gridSystem.Set(gridId, 9, 0, CellTypeChar, '|')

	//Setup Buffer 1
	gridSystem.NewBuffer(gridId, 9, 0, 22, 1)

	// Init Animations

	// 'Ripl3y':"I think you're right"|
	// 'Ripl3y':"I think they have her"|
	// 'Ripl3y':"I found the only open connection"|
	// 'Ripl3y':"I wired it to RABIT; Connect to Rabbit and you're in"|
	// 'Ripl3y':"Good Luck"|

}

func Scene2_HandleDialog1(gs *GridSystem, anims *AnimationSystem) GameState {
	gridId := anims.GridId[AnimationDialog]

	b1RowY := 0
	buffer1Idx := gs.IdxFromXY(gridId, 10, b1RowY)
	message := []byte("I think your're right")

	UpdateDialogAnimation(gridId, buffer1Idx, message, gs, anims)

	b1Cursor := gs.BufferCursor[buffer1Idx]
	gs.Set(gridId, b1Cursor, b1RowY, CellTypeChar, '|')

	if !anims.IsPlaying[AnimationDialog] {
		return Scene2_Dialog_2
	}

	return Scene2_Dialog_1
}
