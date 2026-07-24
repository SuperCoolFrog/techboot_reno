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

	/* Grid27x21x48x12 */
	GridStartScene GridID
)

func InitBucketItems(gs *GridSystem, anims *AnimationSystem) {

	bucketInitIntroAnimation(gs, anims)

}

func bucketInitIntroAnimation(gs *GridSystem, anims *AnimationSystem) {
	Grid27x21x48x12 = gs.AllocateGrid(27, 21, 48, 12, 12)
	GridStartScene = Grid27x21x48x12

	if anims.IsPlaying[AnimationStartScene] {
		return
	}

	anims.IsPlaying[AnimationStartScene] = false
	anims.Loop[AnimationStartScene] = false
	anims.Timers[AnimationStartScene] = 0.0
	anims.Durations[AnimationStartScene] = 1.0
	// anims.Delay[AnimationGridIntro] = 5.0 // Tried to fix vsync at the beginning but just live with it

	gs.SetAllCells(Grid27x21x48x12, CellTypeNone, 0)

	anims.HasGrid[AnimationStartScene] = true
	anims.GridId[AnimationStartScene] = Grid27x21x48x12
}
