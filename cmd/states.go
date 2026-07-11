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
		game.Animations.PlayAnimatedGridIntro(1.0, false)
	case Scene1_Animating:
		if !game.Animations.IsPlaying[AnimationGridIntro] {
			game.State = Scene1_Waiting
		}
	}
}
