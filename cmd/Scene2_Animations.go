package main

func InitDialogAnimation(gs *GridSystem, anims *AnimationSystem) {
	anims.IsPlaying[AnimationDialog] = false
	anims.Loop[AnimationDialog] = false
	anims.Timers[AnimationDialog] = 0.0
	anims.Durations[AnimationDialog] = 0.0

	gridId := gs.AllocateGrid(40, 30, 32, 16)
	gs.SetAllCells(gridId, CellTypeEmpty, 0)
	gs.EnableGrid(gridId)

	anims.HasGrid[AnimationDialog] = true
	anims.GridId[AnimationDialog] = gridId
}

func PlayDialogAnimation(
	duration float32,
	loop bool,
	gs *GridSystem,
	anims *AnimationSystem,
) {
	if anims.IsPlaying[AnimationDialog] {
		return
	}

	anims.IsPlaying[AnimationDialog] = true
	anims.Loop[AnimationDialog] = loop
	anims.Timers[AnimationDialog] = 0.0
	anims.Durations[AnimationDialog] = duration
	anims.Delay[AnimationDialog] = 2
}

func UpdateDialogAnimation(
	gridId GridID, bufferIdx int, message []byte, gs *GridSystem, anims *AnimationSystem,
) {
	timer := anims.Timers[AnimationDialog]
	duration := anims.Durations[AnimationDialog]
	delay := anims.Delay[AnimationDialog]

	trueTime := timer - delay
	completedTime := trueTime / duration

	bCursor := gs.BufferCursor[bufferIdx]

	total := len(message)

	targetCount := int(float32(total) * completedTime)

	for bCursor <= targetCount {
		gs.BufferAppend(gridId, bufferIdx, message[bCursor])
		bCursor = gs.BufferCursor[bufferIdx]
	}
}
