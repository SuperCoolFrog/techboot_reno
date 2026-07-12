package main

type GameState uint64

const (
	Scene1_Init GameState = iota
	Scene1_Start
	Scene1_Animating
	Scene1_Waiting
	Scene1_ExitAnimation
	Scene1_Exiting

	Scene2_Init
	Scene2_CutScene
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
			Scene1_HandleAnimationComplete(game)
		}
	case Scene1_Waiting:
		Scene1_HandleButtonList(game)
	case Scene1_ExitAnimation:
		game.Animations.PlayAnimatedGridExit(game.GridSystem, 5.0, false)
		game.State = Scene1_Exiting
	case Scene1_Exiting:
		if game.Animations.IsPlaying[AnimationIntroGrid] {
			game.Animations.UpdateAnimatedGridExit(game.GridSystem)
		} else {
			game.State = Scene2_Init
		}
	}
}
