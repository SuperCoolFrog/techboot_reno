package main

import (
	"techboot_reno/cmd/assets"
)

func Scene3_HandleInit(current, next GameState, gs *GridSystem, anims *AnimationSystem) GameState {
	gridId := gs.AllocateGrid(40, 30, 32, 0)
	gs.SetAllCells(gridId, CellTypeEmpty, 0)
	gs.EnableGrid(gridId)

	// for i := 0; i < 40; i++ {
	// 	gs.Set(gridId, i, 0, CellTypeChar, byte(i%10)+'0')
	// }

	// Output Window
	// .Top
	gs.SetCellSprite(gridId, 24, 1, assets.SpriteIDCornerTopLeft)
	for i := 1; i < 15; i++ {
		gs.SetCellSprite(gridId, 24+i, 1, assets.SpriteIDHorizontalBar)
	}
	gs.SetCellSprite(gridId, 39, 1, assets.SpriteIDCornerTopRight)
	// .Left
	for i := 1; i < 15; i++ {
		gs.SetCellSprite(gridId, 24, 1+i, assets.SpriteIDVerticalBar)
	}
	// .Right
	for i := 1; i < 15; i++ {
		gs.SetCellSprite(gridId, 39, 1+i, assets.SpriteIDVerticalBar)
	}
	// .Bottom
	gs.SetCellSprite(gridId, 24, 16, assets.SpriteIDCornerBottomLeft)
	for i := 1; i < 15; i++ {
		gs.SetCellSprite(gridId, 24+i, 16, assets.SpriteIDHorizontalBar)
	}
	gs.SetCellSprite(gridId, 39, 16, assets.SpriteIDCornerBottomRight)

	return next
}
