package main

// Animations
const (
	AnimationStartScene AnimationId = iota
	AnimationDialog
	AnimationScanner
	AnimationMemoryStack

	AnimationCount
)

type Bucket struct {
	Grid27x21x48x12 GridID
	Grid42x30x30x16 GridID
	/* Grid27x21x48x12 */
	GridStartScene GridID
	/* Grid40x30x32x16 */
	GridDialogScene GridID

	Buffer42x30x30xfalse BufferID
	/* Buffer for Scene2 Dialog */
	BufferDialogScene BufferID

	/* This is GameState as Idx, some repeating states need iterator i.e. animations/cutscenes */
	SceneStateItr []int
}

func InitBucketItems(gs *GridSystem, anims *AnimationSystem, bs *BufferSystem) Bucket {
	bucket := &Bucket{
		SceneStateItr: make([]int, GameStateCount),
	}
	bucketInitGrids(gs, bucket)
	bucketInitBuffers(bs, bucket)
	bucketInitIntroAnimation(gs, anims, bucket)
	bucketInitDialogAnimation(gs, anims, bucket)

	return *bucket
}

func bucketInitGrids(gs *GridSystem, bucket *Bucket) {
	bucket.Grid27x21x48x12 = gs.AllocateGrid(27, 21, 48, 12, 12)
	gs.SetAllCells(bucket.Grid27x21x48x12, CellTypeNone, 0)

	bucket.Grid42x30x30x16 = gs.AllocateGrid(42, 30, 30, 16, 16)
	gs.SetAllCells(bucket.Grid42x30x30x16, CellTypeNone, 0)
}

func bucketInitBuffers(bs *BufferSystem, bucket *Bucket) {
	bucket.Buffer42x30x30xfalse = bs.AllocateBuffer(42, 30, 30, false)
}

func bucketInitIntroAnimation(gs *GridSystem, anims *AnimationSystem, bucket *Bucket) {
	bucket.GridStartScene = bucket.Grid27x21x48x12

	anims.IsPlaying[AnimationStartScene] = false
	anims.Loop[AnimationStartScene] = false
	anims.Timers[AnimationStartScene] = 0.0
	anims.Durations[AnimationStartScene] = 1.0
	// anims.Delay[AnimationGridIntro] = 5.0 // Tried to fix vsync at the beginning but just live with it

	anims.HasGrid[AnimationStartScene] = true
	anims.GridId[AnimationStartScene] = bucket.Grid27x21x48x12
}

func bucketInitDialogAnimation(gs *GridSystem, anims *AnimationSystem, bucket *Bucket) {
	bucket.GridDialogScene = bucket.Grid42x30x30x16
	bucket.BufferDialogScene = bucket.Buffer42x30x30xfalse

	anims.IsPlaying[AnimationDialog] = false
	anims.Loop[AnimationDialog] = false
	anims.Timers[AnimationDialog] = 0.0
	anims.Durations[AnimationDialog] = 0.0

	//gs.SetAllCells(GridDialogScene, CellTypeEmpty, 0)
	//gs.EnableGrid(GridDialogScene)

	anims.HasGrid[AnimationDialog] = true
	anims.GridId[AnimationDialog] = bucket.GridDialogScene
}
