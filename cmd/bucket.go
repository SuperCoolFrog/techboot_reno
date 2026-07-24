package main

// Animations
const (
	AnimationStartScene AnimationId = iota
	AnimationDialog
	AnimationScanner
	AnimationMemoryStack

	AnimationCount
)

// GridId
var (
	Grid27x21x48x12 GridID
	Grid40x30x32x16 GridID

	/* Grid27x21x48x12 */
	GridStartScene  GridID
	GridDialogScene GridID
)

func InitBucketItems(gs *GridSystem, anims *AnimationSystem) {
	bucketInitGrids(gs)
	bucketInitIntroAnimation(gs, anims)
	bucketInitDialogAnimation(gs, anims)
}

func bucketInitGrids(gs *GridSystem) {
	Grid27x21x48x12 = gs.AllocateGrid(27, 21, 48, 12, 12)
	gs.SetAllCells(Grid27x21x48x12, CellTypeNone, 0)

	Grid40x30x32x16 = gs.AllocateGrid(40, 30, 32, 16, 16)
	gs.SetAllCells(Grid40x30x32x16, CellTypeNone, 0)
}

func bucketInitIntroAnimation(gs *GridSystem, anims *AnimationSystem) {
	GridStartScene = Grid27x21x48x12

	if anims.IsPlaying[AnimationStartScene] {
		return
	}

	anims.IsPlaying[AnimationStartScene] = false
	anims.Loop[AnimationStartScene] = false
	anims.Timers[AnimationStartScene] = 0.0
	anims.Durations[AnimationStartScene] = 1.0
	// anims.Delay[AnimationGridIntro] = 5.0 // Tried to fix vsync at the beginning but just live with it

	anims.HasGrid[AnimationStartScene] = true
	anims.GridId[AnimationStartScene] = Grid27x21x48x12
}

func bucketInitDialogAnimation(gs *GridSystem, anims *AnimationSystem) {
	GridDialogScene = Grid40x30x32x16

	anims.IsPlaying[AnimationDialog] = false
	anims.Loop[AnimationDialog] = false
	anims.Timers[AnimationDialog] = 0.0
	anims.Durations[AnimationDialog] = 0.0

	//gs.SetAllCells(GridDialogScene, CellTypeEmpty, 0)
	//gs.EnableGrid(GridDialogScene)

	anims.HasGrid[AnimationDialog] = true
	anims.GridId[AnimationDialog] = GridDialogScene
}
