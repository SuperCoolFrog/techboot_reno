package main

import (
	"math"
	"techboot_reno/cmd/assets"
)

const (
	OutputGridCols int = S3GridXCount - DividerX - 2
	OutputGridRows int = DividerY - 1
	OutputGridSize int = 20
	OutputGridX    int = OutputGridSize * (DividerX + 1)
	OutputGridY    int = OutputGridSize

	AnimationScannerCols int = OutputGridCols * 4
	AnimationScannerRows int = OutputGridRows * 4
	AnimationScannerSize int = OutputGridSize / 4
)

var (
	OutputGridId         GridID
	AnimationScannerGrid GridID
)

func InitOutputGrid(gs *GridSystem, anims *AnimationSystem) {
	OutputGridId = gs.AllocateGrid(OutputGridCols, OutputGridRows, OutputGridSize, OutputGridX, OutputGridY)
	gs.SetAllCells(OutputGridId, CellTypeSquare, 0)
	gs.EnableGrid(OutputGridId)

	AnimationScannerGrid = gs.AllocateGrid(AnimationScannerCols, AnimationScannerRows, AnimationScannerSize, OutputGridX, OutputGridY)
	gs.SetAllCells(AnimationScannerGrid, CellTypeEmpty, 0)
	gs.EnableGrid(AnimationScannerGrid)

	// for i := 0; i < AnimationScannerCols; i++ {
	// 	gs.SetCellSprite(AnimationScannerGrid, i, 0, assets.SpriteIDHorizontalBar)
	// }

	anims.IsPlaying[AnimationScannerGrid] = true
	anims.Loop[AnimationScannerGrid] = true
	anims.Durations[AnimationScannerGrid] = 10.0
	anims.Timers[AnimationScannerGrid] = 0.0
	anims.Delay[AnimationScannerGrid] = 0
}

func UpdateAnimationGrid(gs *GridSystem, anims *AnimationSystem) {
	timer := anims.Timers[AnimationScannerGrid]
	duration := anims.Durations[AnimationScannerGrid]
	delay := anims.Delay[AnimationScannerGrid]

	trueTime := float32(math.Max(float64(timer-delay), 0))
	completedTime := trueTime / duration

	gs.SetAllCells(AnimationScannerGrid, CellTypeEmpty, 0)

	if completedTime <= 0.5 {
		y := int(float32(AnimationScannerRows-1) * (completedTime / .5))
		for i := 0; i < AnimationScannerCols; i++ {
			gs.SetCellSprite(AnimationScannerGrid, i, y, assets.SpriteIDHorizontalBar)
		}
	} else {
		y := int(float32(AnimationScannerRows-1) * ((completedTime - .5) / .5))
		for i := 0; i < AnimationScannerCols; i++ {
			gs.SetCellSprite(AnimationScannerGrid, i, AnimationScannerRows-1-y, assets.SpriteIDHorizontalBar)
		}
	}
}
