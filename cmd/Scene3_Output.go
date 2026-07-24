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
	GridIdOutput           GridID
	GridIdAnimationScanner GridID
)

func InitOutputGrid(gs *GridSystem, anims *AnimationSystem) {
	GridIdOutput = gs.AllocateGrid(OutputGridCols, OutputGridRows, OutputGridSize, OutputGridX, OutputGridY)
	gs.SetAllCells(GridIdOutput, CellTypeSquare, 0)
	gs.EnableGrid(GridIdOutput)

	GridIdAnimationScanner = gs.AllocateGrid(AnimationScannerCols, AnimationScannerRows, AnimationScannerSize, OutputGridX, OutputGridY)
	gs.SetAllCells(GridIdAnimationScanner, CellTypeEmpty, 0)
	gs.EnableGrid(GridIdAnimationScanner)

	// for i := 0; i < AnimationScannerCols; i++ {
	// 	gs.SetCellSprite(AnimationScannerGrid, i, 0, assets.SpriteIDHorizontalBar)
	// }

	anims.IsPlaying[GridIdAnimationScanner] = true
	anims.Loop[GridIdAnimationScanner] = true
	anims.Durations[GridIdAnimationScanner] = 10.0
	anims.Timers[GridIdAnimationScanner] = 0.0
	anims.Delay[GridIdAnimationScanner] = 0
}

func UpdateAnimationGrid(gs *GridSystem, anims *AnimationSystem) {
	timer := anims.Timers[GridIdAnimationScanner]
	duration := anims.Durations[GridIdAnimationScanner]
	delay := anims.Delay[GridIdAnimationScanner]

	trueTime := float32(math.Max(float64(timer-delay), 0))
	completedTime := trueTime / duration

	gs.SetAllCells(GridIdAnimationScanner, CellTypeEmpty, 0)

	if completedTime <= 0.5 {
		y := int(float32(AnimationScannerRows-1) * (completedTime / .5))
		for i := 0; i < AnimationScannerCols; i++ {
			gs.SetCellSprite(GridIdAnimationScanner, i, y, assets.SpriteIDHorizontalBar)
		}
	} else {
		y := int(float32(AnimationScannerRows-1) * ((completedTime - .5) / .5))
		for i := 0; i < AnimationScannerCols; i++ {
			gs.SetCellSprite(GridIdAnimationScanner, i, AnimationScannerRows-1-y, assets.SpriteIDHorizontalBar)
		}
	}
}
