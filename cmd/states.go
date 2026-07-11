package main

type GameState uint64

const (
	Scene1_Init GameState = iota
	Scene1_Start
	Scene1_Animating
	Scene1_Waiting
)

func (game *Game) UpdateState() {
	switch game.State {
	case Scene1_Init:
		game.State = Scene1_Start
	case Scene1_Start:
		game.State = Scene1_Animating
		game.Animations.PlayAnimatedGridIntro(game.GridSystem, 1.0, false)
	case Scene1_Animating:
		if game.Animations.IsPlaying[AnimationIntroGrid] {
			game.Animations.UpdateAnimatedGridIntro(game.GridSystem)
		} else {
			game.State = Scene1_Waiting
			game.Animations.Offset = 1 // Intro animation is 0
			game.GridSystem.EnableGrid(game.Animations.GridId[AnimationIntroGrid])
		}
	case Scene1_Waiting:
	}
}
