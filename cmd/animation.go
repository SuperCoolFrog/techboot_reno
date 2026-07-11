package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
)

type AnimationId uint32

const (
	AnimationIntroGrid AnimationId = iota
)

const (
	deltaTime = 1.0 / 60.0 // constant because locked to 60hz in ebitengine
)

type AnimationSystem struct {
	Id        []AnimationId
	IsPlaying []bool
	Loop      []bool
	Delay     []float32
	Timers    []float32
	Durations []float32
	HasGrid   []bool
	GridId    []GridID
	Offset    int
}

func NewAnimationSystem() *AnimationSystem {
	/*
		0 = Intro --One Time Animation
		1-10 = typing

		-----
		11
	*/
	const AnimationCount = 11

	a := &AnimationSystem{
		Id:        make([]AnimationId, AnimationCount),
		IsPlaying: make([]bool, AnimationCount),
		Loop:      make([]bool, AnimationCount),
		Delay:     make([]float32, AnimationCount),
		Timers:    make([]float32, AnimationCount),
		Durations: make([]float32, AnimationCount),
		HasGrid:   make([]bool, AnimationCount),
		GridId:    make([]GridID, AnimationCount),
		Offset:    0,
	}

	return a
}

func (anims *AnimationSystem) Update() {
	for animation := anims.Offset; animation < len(anims.Id); animation++ {

		if !anims.IsPlaying[animation] {
			continue
		}
		timer := anims.Timers[animation]
		duration := anims.Durations[animation]
		delay := anims.Delay[animation]

		nuTimer := timer + deltaTime
		anims.Timers[animation] = nuTimer

		if nuTimer < delay {
			continue
		}

		if (nuTimer - delay) >= duration {
			if anims.Loop[animation] {
				anims.Timers[animation] = nuTimer - duration
			} else {
				anims.Timers[animation] = 0.0
				anims.IsPlaying[animation] = false
			}
		}
	}
}

func (anims *AnimationSystem) Render(screen *ebiten.Image, gridSystem *GridSystem) {
	for animation := anims.Offset; animation < len(anims.Id); animation++ {
		if !anims.IsPlaying[animation] {
			continue
		}

		if anims.HasGrid[animation] {
			fmt.Printf("Rendering")
			gridId := anims.GridId[animation]
			gridSystem.Render(screen, gridId)
		}
	}
}

/*
Definitions
*/

func (anims *AnimationSystem) PlayAnimatedGridIntro(gridSystem *GridSystem, duration float32, loop bool) {
	if anims.IsPlaying[AnimationIntroGrid] {
		return
	}

	anims.IsPlaying[AnimationIntroGrid] = true
	anims.Loop[AnimationIntroGrid] = loop
	anims.Timers[AnimationIntroGrid] = 0.0
	anims.Durations[AnimationIntroGrid] = duration
	// anims.Delay[AnimationGridIntro] = 5.0 // Tried to fix vsync at the beginning but just live with it

	gridId := gridSystem.AllocateGrid(27, 21, 48, 0)

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

// func (anims *AnimationSystem) PlayTypingAnimation(grid *Grid, duration float32, string []byte) {
// 	if anims.IsPlaying[AnimationTyping] {
// 		return
// 	}
//
// 	anims.IsPlaying[AnimationTyping] = true
// 	anims.Loop[AnimationTyping] = false
// 	anims.Timers[AnimationTyping] = 0.0
// 	anims.Durations[AnimationTyping] = duration
//
// 	anims.GridAnimationIntro.ResetAndResize(27, 21, 48, 0)
// 	anims.GridAnimationIntro.SetAllCells(CellTypeNone, 0)
// }
