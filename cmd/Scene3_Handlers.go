package main

import (
	"techboot_reno/cmd/assets"
)

func Scene3_HandleInit(current, next GameState, gs *GridSystem, anims *AnimationSystem) GameState {
	gridId := gs.AllocateGrid(40, 30, 32, 0)
	gs.SetAllCells(gridId, CellTypeEmpty, 0)
	gs.EnableGrid(gridId)

	xCount := 40
	yCount := 30

	// Border
	// .Corners
	gs.SetCellSprite(gridId, 0, 0, assets.SpriteIDCornerTopLeft)
	gs.SetCellSprite(gridId, 0, yCount-1, assets.SpriteIDCornerBottomLeft)
	gs.SetCellSprite(gridId, xCount-1, 0, assets.SpriteIDCornerTopRight)
	gs.SetCellSprite(gridId, xCount-1, yCount-1, assets.SpriteIDCornerBottomRight)

	// .Walls
	// ..Left
	for i := 1; i < yCount-1; i++ {
		gs.SetCellSprite(gridId, 0, i, assets.SpriteIDVerticalBar)
	}
	// ..Right
	for i := 1; i < yCount-1; i++ {
		gs.SetCellSprite(gridId, xCount-1, i, assets.SpriteIDVerticalBar)
	}
	// ..Top
	for i := 1; i < xCount-1; i++ {
		gs.SetCellSprite(gridId, i, 0, assets.SpriteIDHorizontalBar)
	}
	// ..Bottom
	for i := 1; i < xCount-1; i++ {
		gs.SetCellSprite(gridId, i, yCount-1, assets.SpriteIDHorizontalBar)
	}

	// Dividers
	// .Vertical
	gs.SetCellSprite(gridId, 24, 0, assets.SpriteIDDownConnectBar)
	for i := 1; i < yCount-1; i++ {
		gs.SetCellSprite(gridId, 24, i, assets.SpriteIDVerticalBar)
	}
	gs.SetCellSprite(gridId, 24, yCount-1, assets.SpriteIDUpConnectBar)
	// .Horizontal
	gs.SetCellSprite(gridId, 24, yCount-11, assets.SpriteIDRightConnectBar)
	for i := 1; i < xCount-25; i++ {
		gs.SetCellSprite(gridId, 24+i, yCount-11, assets.SpriteIDHorizontalBar)
	}
	gs.SetCellSprite(gridId, xCount-1, yCount-11, assets.SpriteIDLeftConnectBar)

	return next
}
