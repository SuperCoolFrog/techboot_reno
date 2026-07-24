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
	Scene2_Dialog_1
	Scene2_Init_Dialog_2
	Scene2_Dialog_2
	Scene2_Init_Dialog_3
	Scene2_Dialog_3
	Scene2_Init_Dialog_4
	Scene2_Dialog_4
	Scene2_Init_Dialog_5
	Scene2_Dialog_5
	Scene2_Init_Dialog_6
	Scene2_Dialog_6
	Scene2_Waiting
	Scene2_CleanUp

	Scene3_Init
	Scene3_InputHandlingLoop
	Scene3_CleanUp

	Scene4_Init
	Scene4_StackAnim
	Scene4_Run

	End // End
)

func (game *Game) UpdateState() {
	gs := game.GridSystem
	anims := game.Animations

	switch game.State {
	case Scene1_Init:
		game.State = Scene1_Start
	case Scene1_Start:
		game.State = Scene1_PlayAnimatedGridIntro(Scene1_Start, Scene1_Animating, game.Animations)
	case Scene1_Animating:
		game.State = Scene1_UpdateAnimatedGridIntro(Scene1_Animating, Scene1_Waiting, game.GridSystem, game.Animations)
	case Scene1_Waiting:
		game.State = Scene1_HandleButtonList(Scene1_Waiting, Scene1_ExitAnimation, game)
	case Scene1_ExitAnimation:
		game.State = Scene1_PlayAnimatedGridExit(Scene1_ExitAnimation, Scene1_Exiting, game.Animations)
	case Scene1_Exiting:
		game.State = Scene1_UpdateAnimatedGridExit(Scene1_Exiting, Scene2_Init, game.GridSystem, game.Animations)
	case Scene2_Init:
		Scene2_HandleInit(game.GridSystem, game.Animations)
		PlayDialogAnimation(2.0, false, game.GridSystem, game.Animations)
		game.State = Scene2_Dialog_1
	case Scene2_Dialog_1:
		game.State = Scene2_HandleDialog1(game.GridSystem, game.Animations)
	case Scene2_Init_Dialog_2:
		Scene2_InitDialog(1, game.GridSystem, game.Animations)
		game.State = Scene2_Dialog_2
	case Scene2_Dialog_2:
		game.State = Scene2_HandleDialog(1, []byte("I don't think she ran..."), Scene2_Dialog_2, Scene2_Init_Dialog_3, gs, anims)
	case Scene2_Init_Dialog_3:
		Scene2_InitDialog(2, game.GridSystem, game.Animations)
		game.State = Scene2_Dialog_3
	case Scene2_Dialog_3:
		game.State = Scene2_HandleDialog(2, []byte("I found an open door"), Scene2_Dialog_3, Scene2_Init_Dialog_4, gs, anims)
	case Scene2_Init_Dialog_4:
		Scene2_InitDialog(3, game.GridSystem, game.Animations)
		game.State = Scene2_Dialog_4
	case Scene2_Dialog_4:
		game.State = Scene2_HandleDialog(3, []byte("CONNECT to RABBIT"), Scene2_Dialog_4, Scene2_Init_Dialog_5, gs, anims)
	case Scene2_Init_Dialog_5:
		Scene2_InitDialog(4, game.GridSystem, game.Animations)
		game.State = Scene2_Dialog_5
	case Scene2_Dialog_5:
		game.State = Scene2_HandleDialog(4, []byte("Good Luck..."), Scene2_Dialog_5, Scene2_Init_Dialog_6, gs, anims)
	case Scene2_Init_Dialog_6:
		gridId := anims.GridId[AnimationDialog]
		s2_AddTRenoMsgBuffer(gs, gridId, 5)
		PlayDialogAnimation(2.0, false, gs, anims)
		game.State = Scene2_Dialog_6
	case Scene2_Dialog_6:
		game.State = Scene2_HandleDialog(5, []byte("thanks"), Scene2_Dialog_6, Scene2_Waiting, gs, anims)
	case Scene2_Waiting:
		game.State = Scene2_WaitForEnter(Scene2_Waiting, Scene2_CleanUp)
	case Scene2_CleanUp:
		game.State = Scene2_CleanUpScene(Scene3_Init, gs, anims)
	case Scene3_Init:
		game.State = Scene3_HandleInit(Scene3_Init, Scene3_InputHandlingLoop, gs, anims)
	case Scene3_InputHandlingLoop:
		game.State = Scene3_Update(game.inputRunes, Scene3_InputHandlingLoop, Scene3_CleanUp, game.prologInput, game.prologOutput, gs, anims)
	case Scene3_CleanUp:
		game.State = Scene3_HandleCleaUp(Scene4_Init, gs, anims)
	case Scene4_Init:
		game.State = Scene4_HandleInit(Scene4_Init, Scene4_StackAnim, gs, anims)
	case Scene4_StackAnim:
		game.State = Scene4_UpdateStackAnimation(Scene4_StackAnim, Scene4_Run, gs, anims)
	case Scene4_Run:
		game.State = Scene4_Update(Scene4_Run, End, game.inputRunes, game.prologInput, game.prologOutput, gs, anims)
	}

}
