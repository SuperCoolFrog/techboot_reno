package main

import (
	"techboot_reno/cmd/assets"
)

func Scene3_HandleInit(current, next GameState, gs *GridSystem, anims *AnimationSystem) GameState {
	gridId := gs.AllocateGrid(40, 30, 32, 16)
	gs.SetAllCells(gridId, CellTypeEmpty, 0)
	gs.EnableGrid(gridId)

	gs.SetCellSprite(gridId, 0, 0, assets.SpriteIDCornerTopLeft)

	return next
}
