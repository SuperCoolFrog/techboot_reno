package main

type Animation uint32

const (
	AnimationGridIntro Animation = iota

	AnimationCount // Always last for looping
)

const (
	deltaTime = 1.0 / 60.0 // constant because locked to 60hz in ebitengine
)

type Animations struct {
	IsPlaying []bool
	Loop      []bool
	Timers    []float32
	Durations []float32

	introGrid Grid
}

func NewAnimations() Animations {
	a := Animations{
		IsPlaying: make([]bool, AnimationCount),
		Loop:      make([]bool, AnimationCount),
		Timers:    make([]float32, AnimationCount),
		Durations: make([]float32, AnimationCount),
	}

	return a
}

func (anims Animations) Update() {

	for animation := range AnimationCount {

		if !anims.IsPlaying[animation] {
			continue
		}
		timer := anims.Timers[animation]
		duration := anims.Durations[animation]

		nuTimer := timer + deltaTime

		anims.Timers[animation] = nuTimer

		if nuTimer >= duration {
			if anims.Loop[animation] {
				anims.Timers[animation] = nuTimer - duration
			} else {
				anims.Timers[animation] = 0.0
				anims.IsPlaying[animation] = false
			}
		}
	}
}

/* */

func (anims Animations) PlayAnimatedGridIntro(duration float32, loop bool) {
	if anims.IsPlaying[AnimationGridIntro] {
		return
	}

	anims.IsPlaying[AnimationGridIntro] = true
	anims.Loop[AnimationGridIntro] = loop
	anims.Timers[AnimationGridIntro] = 0.0
	anims.Durations[AnimationGridIntro] = duration
}
