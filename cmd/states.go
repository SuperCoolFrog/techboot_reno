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

	End // End
)

func (game *Game) UpdateState() {
	gs := game.GridSystem
	anims := game.Animations

	switch game.State {
	case Scene1_Init:
		game.State = Scene1_Start
	case Scene1_Start:
		game.State = Scene1_Animating
		game.Animations.PlayAnimatedGridIntro(game.GridSystem, 1.0, false)
	case Scene1_Animating:
		if game.Animations.IsPlaying[AnimationGrid] {
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
		if game.Animations.IsPlaying[AnimationGrid] {
			game.Animations.UpdateAnimatedGridExit(game.GridSystem)
		} else {
			game.State = Scene2_Init
		}
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
		game.State = Scene2_HandleDialog(1, []byte("I don't this she ran..."), Scene2_Dialog_2, Scene2_Init_Dialog_3, gs, anims)
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
		game.State = Scene2_HandleDialog(5, []byte("thanks"), Scene2_Dialog_6, End, gs, anims)
	}
}
