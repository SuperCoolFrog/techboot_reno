package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Animation uint32

const (
	AnimationGridIntro Animation = iota

	AnimationCount // Always last for looping
)

const (
	deltaTime = 1.0 / 60.0 // constant because locked to 60hz in ebitengine
)

type Animations struct {
	IsPlaying          []bool
	Loop               []bool
	Delay              []float32
	Timers             []float32
	Durations          []float32
	GridAnimationIntro *Grid
}

func NewAnimations() *Animations {
	a := &Animations{
		IsPlaying:          make([]bool, AnimationCount),
		Loop:               make([]bool, AnimationCount),
		Delay:              make([]float32, AnimationCount),
		Timers:             make([]float32, AnimationCount),
		Durations:          make([]float32, AnimationCount),
		GridAnimationIntro: NewGrid(27, 21, 48, 0),
	}

	return a
}

func (anims *Animations) Update() {
	for animation := range AnimationCount {

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

		// Inidividual Updates
		switch animation {
		case AnimationGridIntro:
			anims.updateAnimatedGridIntro()
		}
	}
}

func (anims *Animations) Render(screen *ebiten.Image) {
	for animation := range AnimationCount {
		if !anims.IsPlaying[animation] {
			continue
		}

		// Inidividual Renders
		switch animation {
		case AnimationGridIntro:
			anims.GridAnimationIntro.Render(screen)
		}
	}
}

func (anims *Animations) PlayAnimatedGridIntro(duration float32, loop bool) {
	if anims.IsPlaying[AnimationGridIntro] {
		return
	}

	anims.IsPlaying[AnimationGridIntro] = true
	anims.Loop[AnimationGridIntro] = loop
	anims.Timers[AnimationGridIntro] = 0.0
	anims.Durations[AnimationGridIntro] = duration
	// anims.Delay[AnimationGridIntro] = 5.0 // Tried to fix vsync at the beginning but just live with it

	anims.GridAnimationIntro.ResetAndResize(27, 21, 48, 0)
	anims.GridAnimationIntro.SetAllCells(RenderFlagNone, 0)
}

func (anims *Animations) updateAnimatedGridIntro() {
	timer := anims.Timers[AnimationGridIntro]
	duration := anims.Durations[AnimationGridIntro]
	delay := anims.Delay[AnimationGridIntro]

	grid := anims.GridAnimationIntro
	completed := float64(timer-delay) / float64(duration)

	maxCol := int(float64(grid.Cols) * completed)
	maxRow := int(float64(grid.Rows) * completed)

	for x := 0; x < maxCol; x++ {
		// grid.Set(x, 0, RenderFlagCellSquare, 0)
		for y := 0; y < maxRow; y++ {
			grid.Set(x, y, RenderFlagCellSquare, 0)
		}
	}

}
