package main

func (anims *AnimationSystem) PlayAnimatedGridIntro(gridSystem *GridSystem, duration float32, loop bool) {
	if anims.IsPlaying[AnimationIntroGrid] {
		return
	}

	anims.IsPlaying[AnimationIntroGrid] = true
	anims.Loop[AnimationIntroGrid] = loop
	anims.Timers[AnimationIntroGrid] = 0.0
	anims.Durations[AnimationIntroGrid] = duration
	// anims.Delay[AnimationGridIntro] = 5.0 // Tried to fix vsync at the beginning but just live with it

	gridId := gridSystem.AllocateGrid(27, 21, 48, 12)

	gridSystem.SetAllCells(gridId, CellTypeNone, 0)

	anims.HasGrid[AnimationIntroGrid] = true
	anims.GridId[AnimationIntroGrid] = gridId
}

func (anims *AnimationSystem) UpdateAnimatedGridIntro(gridSystem *GridSystem) {
	timer := anims.Timers[AnimationIntroGrid]
	duration := float64(anims.Durations[AnimationIntroGrid])
	delay := anims.Delay[AnimationIntroGrid]
	gridId := anims.GridId[AnimationIntroGrid]
	cols := gridSystem.Cols[gridId]
	rows := gridSystem.Rows[gridId]

	// steps := 2.0
	//stepDuration := duration / steps
	trueTime := float64(timer - delay)

	// Step 1
	step1Duration := duration * 0.5

	completed1 := trueTime / step1Duration

	if completed1 < 1.0 {
		maxCol := int(float64(cols) * completed1)
		maxRow := int(float64(rows) * completed1)

		for x := 0; x < maxCol; x++ {
			// grid.Set(x, 0, RenderFlagCellSquare, 0)
			for y := 0; y < maxRow; y++ {
				gridSystem.Set(gridId, x, y, CellTypeSquare, 0)
			}
		}

		return
	}

	// Step 2
	step2Duration := step1Duration + duration*0.5

	completed2 := trueTime / step2Duration

	//	Start?

	if completed2 < 1.0 {
		halfCols := float64(cols) / 2.0
		halfRows := float64(rows) / 2.0

		targetCutCols := float64(cols) / 3.5
		targetCutRows := float64(rows) / 2.5

		cutCols := targetCutCols * completed2
		cutRows := targetCutRows * completed2

		minCol := int(halfCols) - int(cutCols)
		minRow := int(halfRows) - int(cutRows)
		maxCol := int(halfCols) + int(cutCols)
		maxRow := int(halfRows) + int(cutRows)

		for x := minCol; x < maxCol; x++ {
			for y := minRow; y < maxRow; y++ {
				gridSystem.Set(gridId, x, y, CellTypeEmpty, 0)
			}
		}
	}
}

func (anims *AnimationSystem) PlayAnimatedGridExit(gridSystem *GridSystem, duration float32, loop bool) {
	if anims.IsPlaying[AnimationIntroGrid] {
		return
	}

	anims.IsPlaying[AnimationIntroGrid] = true
	anims.Loop[AnimationIntroGrid] = loop
	anims.Timers[AnimationIntroGrid] = 0.0
	anims.Durations[AnimationIntroGrid] = duration

	// Assume Grid exists from intro
}

func (anims *AnimationSystem) UpdateAnimatedGridExit(gridSystem *GridSystem) {
	timer := anims.Timers[AnimationIntroGrid]
	duration := float64(anims.Durations[AnimationIntroGrid])
	delay := anims.Delay[AnimationIntroGrid]
	gridId := anims.GridId[AnimationIntroGrid]

	trueTime := float64(timer - delay)

	// Step 1
	step1Duration := duration * 0.32

	completed1 := trueTime / step1Duration

	if completed1 < 1.0 {

		if completed1 >= .95 {
			gridSystem.Set(gridId, 9, 3, CellTypeEmpty, ' ')
		}
		if completed1 >= .9 {
			gridSystem.Set(gridId, 10, 3, CellTypeEmpty, ' ')
		}
		if completed1 >= .85 {
			gridSystem.Set(gridId, 11, 3, CellTypeEmpty, ' ')
		}
		if completed1 >= .8 {
			gridSystem.Set(gridId, 12, 3, CellTypeEmpty, ' ')
		}
		if completed1 >= .75 {
			gridSystem.Set(gridId, 13, 3, CellTypeEmpty, ' ')
		}
		if completed1 >= .7 {
			gridSystem.Set(gridId, 14, 3, CellTypeEmpty, ' ')
		}
		if completed1 >= .65 {
			gridSystem.Set(gridId, 15, 3, CellTypeEmpty, ' ')
		}
		if completed1 >= .6 {
			gridSystem.Set(gridId, 16, 3, CellTypeEmpty, ' ')
		}

		if completed1 >= .55 {
			gridSystem.Set(gridId, 11, 4, CellTypeEmpty, ' ')
		}
		if completed1 >= .5 {
			gridSystem.Set(gridId, 12, 4, CellTypeEmpty, ' ')
		}
		if completed1 >= .45 {
			gridSystem.Set(gridId, 13, 4, CellTypeEmpty, ' ')
		}

		if completed1 >= .4 {
			gridSystem.Set(gridId, 14, 4, CellTypeEmpty, ' ')
		}

		if completed1 >= .35 {
			gridSystem.Set(gridId, 8, 9, CellTypeEmpty, ' ')
		}
		if completed1 >= .3 {
			gridSystem.Set(gridId, 9, 9, CellTypeEmpty, ' ')
		}
		if completed1 >= .25 {
			gridSystem.Set(gridId, 10, 9, CellTypeEmpty, ' ')
		}
		if completed1 >= .2 {
			gridSystem.Set(gridId, 11, 9, CellTypeEmpty, ' ')
		}
		if completed1 >= .15 {
			gridSystem.Set(gridId, 12, 9, CellTypeEmpty, ' ')
		}
		if completed1 >= .1 {
			gridSystem.Set(gridId, 13, 9, CellTypeEmpty, ' ')
		}
		if completed1 >= .05 {
			gridSystem.Set(gridId, 14, 9, CellTypeEmpty, ' ')
		}

		return
	}

	// Step 2
	cols := gridSystem.Cols[gridId]
	rows := gridSystem.Rows[gridId]

	step2Duration := step1Duration + duration*0.32
	completed2 := trueTime / step2Duration

	//	Start

	if completed2 < 1.0 {
		halfCols := float64(cols) / 2.0
		halfRows := float64(rows) / 2.0

		targetCutCols := float64(cols) / 2
		targetCutRows := float64(rows) / 2

		cutCols := targetCutCols * completed2
		cutRows := targetCutRows * completed2

		minCol := int(halfCols) - int(cutCols)
		minRow := int(halfRows) - int(cutRows)
		maxCol := int(halfCols) + int(cutCols)
		maxRow := int(halfRows) + int(cutRows)

		for x := minCol; x < maxCol; x++ {
			for y := minRow; y < maxRow; y++ {
				gridSystem.Set(gridId, x, y, CellTypeEmpty, 0)
			}
		}

		gridSystem.Set(gridId, 7, 8, CellTypeChar, ':')
		gridSystem.Set(gridId, 8, 8, CellTypeChar, '[')
		gridSystem.Set(gridId, 9, 8, CellTypeChar, 'S')
		gridSystem.Set(gridId, 10, 8, CellTypeChar, 'T')
		gridSystem.Set(gridId, 11, 8, CellTypeChar, 'A')
		gridSystem.Set(gridId, 12, 8, CellTypeChar, 'R')
		gridSystem.Set(gridId, 13, 8, CellTypeChar, 'T')
		gridSystem.Set(gridId, 14, 8, CellTypeChar, ']')
		gridSystem.Set(gridId, 15, 8, CellTypeChar, '|')
		return
	}

	// Step 3
	step3Duration := step2Duration + duration*0.36
	completed3 := trueTime / step3Duration

	if completed3 <= 1.0 {
		if ct, _ := gridSystem.Get(gridId, 9, 8); completed3 >= .9 && ct != CellTypeEmpty {
			gridSystem.Set(gridId, 8, 8, CellTypeChar, '|')
			gridSystem.Set(gridId, 9, 8, CellTypeEmpty, ' ')
		}
		if ct, _ := gridSystem.Get(gridId, 10, 8); completed3 >= .8 && ct != CellTypeEmpty {
			gridSystem.Set(gridId, 9, 8, CellTypeChar, '|')
			gridSystem.Set(gridId, 10, 8, CellTypeEmpty, ' ')
		}
		if ct, _ := gridSystem.Get(gridId, 11, 8); completed3 >= .7 && ct != CellTypeEmpty {
			gridSystem.Set(gridId, 10, 8, CellTypeChar, '|')
			gridSystem.Set(gridId, 11, 8, CellTypeEmpty, ' ')
		}
		if ct, _ := gridSystem.Get(gridId, 12, 8); completed3 >= .6 && ct != CellTypeEmpty {
			gridSystem.Set(gridId, 11, 8, CellTypeChar, '|')
			gridSystem.Set(gridId, 12, 8, CellTypeEmpty, ' ')
		}
		if ct, _ := gridSystem.Get(gridId, 13, 8); completed3 >= .5 && ct != CellTypeEmpty {
			gridSystem.Set(gridId, 12, 8, CellTypeChar, '|')
			gridSystem.Set(gridId, 13, 8, CellTypeEmpty, ' ')
		}
		if ct, _ := gridSystem.Get(gridId, 14, 8); completed3 >= .4 && ct != CellTypeEmpty {
			gridSystem.Set(gridId, 13, 8, CellTypeChar, '|')
			gridSystem.Set(gridId, 14, 8, CellTypeEmpty, ' ')
		}
		if ct, _ := gridSystem.Get(gridId, 15, 8); completed3 >= .2 && ct != CellTypeEmpty {
			gridSystem.Set(gridId, 14, 8, CellTypeChar, '|')
			gridSystem.Set(gridId, 15, 8, CellTypeEmpty, ' ')
		}

	}

}
