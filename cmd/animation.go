package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type AnimationId uint32

const (
	AnimationGrid AnimationId = iota
	AnimationDialog
	AnimationScanner
	AnimationMemoryStack
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
			gridId := anims.GridId[animation]
			gridSystem.RenderGrid(screen, gridId)
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
