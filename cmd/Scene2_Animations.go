package main

func InitDialogAnimation(gs *GridSystem, anims *AnimationSystem) {
	anims.IsPlaying[AnimationDialog] = false
	anims.Loop[AnimationDialog] = false
	anims.Timers[AnimationDialog] = 0.0
	anims.Durations[AnimationDialog] = 0.0

	gridId := gs.AllocateGrid(40, 30, 32, 12)
	gs.SetAllCells(gridId, CellTypeSquare, 0)
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

	anims.IsPlaying[AnimationDialog] = false
	anims.Loop[AnimationDialog] = loop
	anims.Timers[AnimationDialog] = 0.0
	anims.Durations[AnimationDialog] = duration
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

	currentCount := bCursor - bufferIdx
	total := len(message)

	targetCount := int(float32(total) * completedTime)

	for currentCount < targetCount {
		gs.BufferAppend(gridId, bufferIdx, message[bCursor])
		currentCount++
		bCursor = gs.BufferCursor[bufferIdx]
	}
}
