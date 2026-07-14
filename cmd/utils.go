package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func utilDebouncedKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 30 // Frames before continuous backspacing starts (e.g., 0.5s)
		interval = 3  // Frames between each backspace after delay
	)

	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true // Initial press triggers instantly
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true // Subsequent repeating triggers
	}
	return false
}
