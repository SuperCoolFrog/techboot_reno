package main

func (anims *AnimationSystem) PlayAnimatedGridIntro(gridSystem *GridSystem, duration float32, loop bool) {
	if anims.IsPlaying[AnimationGrid] {
		return
	}

	anims.IsPlaying[AnimationGrid] = true
	anims.Loop[AnimationGrid] = loop
	anims.Timers[AnimationGrid] = 0.0
	anims.Durations[AnimationGrid] = duration
	// anims.Delay[AnimationGridIntro] = 5.0 // Tried to fix vsync at the beginning but just live with it

	gridId := gridSystem.AllocateGrid(27, 21, 48, 12, 12)

	gridSystem.SetAllCells(gridId, CellTypeNone, 0)

	anims.HasGrid[AnimationGrid] = true
	anims.GridId[AnimationGrid] = gridId
}

func (anims *AnimationSystem) UpdateAnimatedGridIntro(gridSystem *GridSystem) {
	timer := anims.Timers[AnimationGrid]
	duration := float64(anims.Durations[AnimationGrid])
	delay := anims.Delay[AnimationGrid]
	gridId := anims.GridId[AnimationGrid]
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
	if anims.IsPlaying[AnimationGrid] {
		return
	}

	anims.IsPlaying[AnimationGrid] = true
	anims.Loop[AnimationGrid] = loop
	anims.Timers[AnimationGrid] = 0.0
	anims.Durations[AnimationGrid] = duration

	// Assume Grid exists from intro
}

func (anims *AnimationSystem) UpdateAnimatedGridExit(gridSystem *GridSystem) {
	timer := anims.Timers[AnimationGrid]
	duration := float64(anims.Durations[AnimationGrid])
	delay := anims.Delay[AnimationGrid]
	gridId := anims.GridId[AnimationGrid]
	cols := gridSystem.Cols[gridId]
	rows := gridSystem.Rows[gridId]

	trueTime := float64(timer - delay)

	// Step 1
	step1Duration := duration * 0.25

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
	step2Duration := step1Duration + duration*0.25
	completed2 := (trueTime - step1Duration) / (duration * 0.25)

	//	Start

	if completed2 < 1.0 {
		halfCols := float64(cols) / 2.0
		halfRows := float64(rows) / 2.0
		qtrCols := halfCols / 2.0
		qtrRows := halfRows / 2.0

		cutCols := qtrCols * completed2
		cutRows := qtrRows * completed2

		minCol := int(halfCols-qtrCols) - int(cutCols)
		minRow := int(halfRows-qtrRows) - int(cutRows)
		maxCol := int(halfCols+qtrCols) + int(cutCols)
		maxRow := int(halfRows+qtrRows) + int(cutRows)

		for x := minCol; x < maxCol; x++ {
			for y := minRow; y < maxRow; y++ {
				if x >= 0 && x < cols && y >= 0 && y < rows {
					gridSystem.Set(gridId, x, y, CellTypeEmpty, 0)
				}
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
	step3Duration := step2Duration + duration*0.25
	completed3 := (trueTime - step2Duration) / (duration * 0.25)

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
		if ct, _ := gridSystem.Get(gridId, 14, 8); completed3 >= .3 && ct != CellTypeEmpty {
			gridSystem.Set(gridId, 13, 8, CellTypeChar, '|')
			gridSystem.Set(gridId, 14, 8, CellTypeEmpty, ' ')
		}
		if ct, _ := gridSystem.Get(gridId, 15, 8); completed3 >= .1 && ct != CellTypeEmpty {
			gridSystem.Set(gridId, 14, 8, CellTypeChar, '|')
			gridSystem.Set(gridId, 15, 8, CellTypeEmpty, ' ')
		}

		return
	}

	// Step 4
	gridSystem.Set(gridId, 8, 8, CellTypeEmpty, ' ')
	step4Duration := step3Duration + duration*.125
	completed4 := (trueTime - step3Duration) / (duration * 0.125)

	if completed4 < 1.0 {

		for x := 1; x <= int(float64(8)*completed4); x++ {
			if x <= 6 {
				gridSystem.Set(gridId, 6-x, 8, CellTypeChar, ':')
				gridSystem.Set(gridId, 7-x, 8, CellTypeChar, '|')
				gridSystem.Set(gridId, 8-x, 8, CellTypeEmpty, ' ')
			} else {
				gridSystem.Set(gridId, 0, 8, CellTypeChar, ':')
				gridSystem.Set(gridId, 1, 8, CellTypeChar, '|')
				gridSystem.Set(gridId, 2, 8, CellTypeChar, ' ')
			}

		}
		return
	}

	// Step 5
	completed5 := (trueTime - step4Duration) / (duration * 0.125)

	if completed5 < 1.0 {
		for y := 1; y <= int(float64(8)*completed5); y++ {
			if y < 8 {
				gridSystem.Set(gridId, 0, 8-y, CellTypeChar, ':')
				gridSystem.Set(gridId, 1, 8-y, CellTypeChar, '|')
				gridSystem.Set(gridId, 0, 8-y+1, CellTypeEmpty, ' ')
				gridSystem.Set(gridId, 1, 8-y+1, CellTypeEmpty, ' ')
			} else {
				gridSystem.Set(gridId, 0, 0, CellTypeChar, ':')
				gridSystem.Set(gridId, 1, 0, CellTypeChar, '|')
				gridSystem.Set(gridId, 0, 1, CellTypeEmpty, ' ')
				gridSystem.Set(gridId, 1, 1, CellTypeEmpty, ' ')
			}
		}
	}
}
