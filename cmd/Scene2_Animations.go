package main

func PlayDialogAnimation(
	gs *GridSystem,
	anims *AnimationSystem,
) {

	if anims.IsPlaying[AnimationDialog] {
		return
	}

	anims.IsPlaying[AnimationDialog] = true
	anims.Loop[AnimationDialog] = false
	anims.Timers[AnimationDialog] = 0.0
	anims.Durations[AnimationDialog] = 2.0
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
