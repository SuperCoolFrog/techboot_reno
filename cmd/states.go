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
	Scene2_Dialog
	Scene2_Waiting
	Scene2_CleanUp

	Scene3_Init
	Scene3_UpdateLoop
	Scene3_CleanUp

	Scene4_Init
	Scene4_StackAnim
	Scene4_Run

	End // End

	GameStateCount
)

func (game *Game) UpdateState() {
	// gs := game.GridSystem
	// anims := game.Animations

	switch game.State {
	case Scene1_Init:
		game.State = Scene1_Start
	case Scene1_Start:
		game.State = Scene1_PlayAnimatedGridIntro(Scene1_Start, Scene1_Animating, game)
	case Scene1_Animating:
		game.State = Scene1_UpdateAnimatedGridIntro(Scene1_Animating, Scene1_Waiting, game)
	case Scene1_Waiting:
		game.State = Scene1_HandleButtonList(Scene1_Waiting, Scene1_ExitAnimation, game)
	case Scene1_ExitAnimation:
		game.State = Scene1_PlayAnimatedGridExit(Scene1_ExitAnimation, Scene1_Exiting, game.Animations)
	case Scene1_Exiting:
		game.State = Scene1_UpdateAnimatedGridExit(Scene1_Exiting, Scene2_Init, game)
	case Scene2_Init:
		game.State = Scene2_HandleInit(Scene2_Init, Scene2_Dialog, game)
	case Scene2_Dialog:
		game.State = Scene2_HandleAllDialog(Scene2_Dialog, Scene2_Waiting, game)
	case Scene2_Waiting:
		game.State = Scene2_WaitForEnter(Scene2_Waiting, Scene2_CleanUp, game)
	case Scene2_CleanUp:
		game.State = Scene2_CleanUpScene(Scene3_Init, game)
	case Scene3_Init:
		game.State = Scene3_HandleInit(Scene3_Init, Scene3_UpdateLoop, game)
	case Scene3_UpdateLoop:
		game.State = Scene3_Update(Scene3_UpdateLoop, Scene3_CleanUp, game)
		/**
		case Scene3_CleanUp:
			game.State = Scene3_HandleCleaUp(Scene4_Init, gs, anims)
		case Scene4_Init:
			game.State = Scene4_HandleInit(Scene4_Init, Scene4_StackAnim, gs, anims)
		case Scene4_StackAnim:
			game.State = Scene4_UpdateStackAnimation(Scene4_StackAnim, Scene4_Run, gs, anims)
		case Scene4_Run:
			game.State = Scene4_Update(Scene4_Run, End, game.inputRunes, game.prologInput, game.prologOutput, gs, anims)
			**/
	}

}
