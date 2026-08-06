package main

import (
	"fmt"
	"techboot_reno/cmd/assets"
)

// Animations
const (
	AnimationStartScene AnimationId = iota
	AnimationDialog
	AnimationScanner
	AnimationMemoryStack

	AnimationCount
)

type CommandId uint32

const (
	CommandInvalid CommandId = iota
	CommandConnectTrue
	CommandConnectFalse
	CommandList
	CommandFiles
	CommandNetworks
	CommandRoy1Fn
	CommandRoy2Fn
	CommandRoy3Fn
	CommandBreach
	CommandLobby
	CommandConnect
	CommandListSpecific
	CommandPuzzle
	CommandPuzzleIntro
	CommandPuzzleEasy
	CommandPuzzleMed
	CommandPuzzleHard
	CommandSet
	CommandGateTypeAnd
	CommandGateTypeOr

	CommandsCount
)

type Bucket struct {
	Grid27x21x48x12      GridID
	Grid42x30x30x16      GridID
	Grid64x48x20x0       GridID
	Grid26x36x20x740x20  GridID
	Grid104x144x5x740x20 GridID

	/* Grid27x21x48x12 */
	GridStartScene GridID
	/* Grid40x30x32x16 */
	GridDialogScene GridID
	/* Grid64x48x20x0 */
	GridMainUI GridID
	/* Grid26x36x29x740x20 */
	GridOutput GridID
	/* Grid52x72x5x740x20 */
	GridOutputPrecision GridID

	Buffer42x30x30xfalse   BufferID
	Buffer35x46x2000xfalse BufferID
	Buffer26x9x200xtrue    BufferID //48-48-11-2

	/* Buffer for Scene2 Dialog */
	BufferDialogScene BufferID
	/* Buffer for entering Commands */
	BufferCommands       BufferID
	BufferLogs           BufferID
	BufferCommandOutputs BufferID // Grid26x36x20x740x20

	/* Row mapping for command output list values */
	CommandOutputStartRowIdx []int // CommandId -> rowIdx

	/* Stick it in here for now until I think more about these type of values */
	LogBufferColIdx int
	LogBufferRowIdx int

	/* This is GameState as Idx, some repeating states need iterator i.e. animations/cutscenes */
	SceneStateItr []int
}

/* Pass in Game with Systems initialized*/
func InitBucketItems(game *Game) Bucket {
	bucket := &Bucket{
		SceneStateItr:            make([]int, GameStateCount),
		CommandOutputStartRowIdx: make([]int, int(CommandsCount)),
		LogBufferColIdx:          37,
		LogBufferRowIdx:          38,
	}

	bucketInitGrids(game.gs, bucket)
	bucketInitBuffers(game.bs, bucket)

	bucketInitIntroAnimation(game.gs, game.ans, bucket)
	bucketInitDialogAnimation(game.gs, game.ans, bucket)
	bucketInitScannerAnimations(game.gs, game.ans, bucket)

	bucketInitMainUI(game.gs, game.bs, bucket)

	bucketInitCommandStrings(game.bs, bucket)

	bucketInitializeGamePuzzles(game.pz, game.ps, game.jxpp, game.gs, bucket)

	return *bucket
}

func bucketInitGrids(gs *GridSystem, bucket *Bucket) {
	bucket.Grid27x21x48x12 = gs.AllocateGrid(27, 21, 48, 12, 12)
	gs.SetAllCells(bucket.Grid27x21x48x12, CellTypeNone, 0)

	bucket.Grid42x30x30x16 = gs.AllocateGrid(42, 30, 30, 16, 16)
	gs.SetAllCells(bucket.Grid42x30x30x16, CellTypeNone, 0)

	bucket.Grid64x48x20x0 = gs.AllocateGrid(64, 48, 20, 0, 0)
	gs.SetAllCells(bucket.Grid64x48x20x0, CellTypeNone, 0)

	bucket.Grid26x36x20x740x20 = gs.AllocateGrid(26, 36, 20, 740, 20)
	gs.SetAllCells(bucket.Grid26x36x20x740x20, CellTypeNone, 0)

	bucket.Grid104x144x5x740x20 = gs.AllocateGrid(104, 144, 5, 740, 20)
	gs.SetAllCells(bucket.Grid104x144x5x740x20, CellTypeNone, 0)
}

func bucketInitBuffers(bs *BufferSystem, bucket *Bucket) {
	bucket.Buffer42x30x30xfalse = bs.AllocateBuffer(42, 30, 30, false)
	bucket.Buffer35x46x2000xfalse = bs.AllocateBuffer(35, 46, 2000, false)
	bucket.Buffer26x9x200xtrue = bs.AllocateBuffer(26, 9, 200, true)
	bucket.BufferCommandOutputs = bs.AllocateBuffer(26, int(CommandsCount), int(CommandsCount), false)
}

func bucketInitIntroAnimation(gs *GridSystem, anims *AnimationSystem, bucket *Bucket) {
	bucket.GridStartScene = bucket.Grid27x21x48x12

	anims.IsPlaying[AnimationStartScene] = false
	anims.Loop[AnimationStartScene] = false
	anims.Timers[AnimationStartScene] = 0.0
	anims.Durations[AnimationStartScene] = 1.0
	// anims.Delay[AnimationGridIntro] = 5.0 // Tried to fix vsync at the beginning but just live with it

	anims.HasGrid[AnimationStartScene] = true
}

func bucketInitDialogAnimation(gs *GridSystem, anims *AnimationSystem, bucket *Bucket) {
	bucket.GridDialogScene = bucket.Grid42x30x30x16
	bucket.BufferDialogScene = bucket.Buffer42x30x30xfalse

	anims.IsPlaying[AnimationDialog] = false
	anims.Loop[AnimationDialog] = false
	anims.Timers[AnimationDialog] = 0.0
	anims.Durations[AnimationDialog] = 0.0

	//gs.SetAllCells(GridDialogScene, CellTypeEmpty, 0)
	//gs.EnableGrid(GridDialogScene)

	anims.HasGrid[AnimationDialog] = true
}

func bucketInitScannerAnimations(gs *GridSystem, anims *AnimationSystem, bucket *Bucket) {
	anims.IsPlaying[AnimationScanner] = false
	anims.Loop[AnimationScanner] = true
	anims.Durations[AnimationScanner] = 10.0
	anims.Timers[AnimationScanner] = 0.0
	anims.Delay[AnimationScanner] = 0
}

func bucketInitMainUI(gs *GridSystem, bs *BufferSystem, bucket *Bucket) {
	bucket.GridMainUI = bucket.Grid64x48x20x0
	bucket.GridOutput = bucket.Grid26x36x20x740x20
	bucket.GridOutputPrecision = bucket.Grid104x144x5x740x20
	bucket.BufferCommands = bucket.Buffer35x46x2000xfalse
	bucket.BufferLogs = bucket.Buffer26x9x200xtrue

	gs.SetAllCells(bucket.GridMainUI, CellTypeEmpty, 0)
	// gs.EnableGrid(bucket.GridMainUI)

	rows := gs.Rows[bucket.GridMainUI]
	cols := gs.Cols[bucket.GridMainUI]

	// Border

	// .Corners
	gs.SetCellSprite(bucket.GridMainUI, 0, 0, assets.SpriteIDCornerTopLeft)
	gs.SetCellSprite(bucket.GridMainUI, 0, rows-1, assets.SpriteIDCornerBottomLeft)
	gs.SetCellSprite(bucket.GridMainUI, cols-1, 0, assets.SpriteIDCornerTopRight)
	gs.SetCellSprite(bucket.GridMainUI, cols-1, rows-1, assets.SpriteIDCornerBottomRight)

	// .Walls
	// ..Left
	for i := 1; i < rows-1; i++ {
		gs.SetCellSprite(bucket.GridMainUI, 0, i, assets.SpriteIDVerticalBar)
	}
	// ..Right
	for i := 1; i < rows-1; i++ {
		gs.SetCellSprite(bucket.GridMainUI, cols-1, i, assets.SpriteIDVerticalBar)
	}
	// ..Top
	for i := 1; i < cols-1; i++ {
		gs.SetCellSprite(bucket.GridMainUI, i, 0, assets.SpriteIDHorizontalBar)
	}

	DividerX := 36
	DividerY := rows - 11

	// ...Commands
	hdrCmdX := DividerX/2 - 4
	gs.Set(bucket.GridMainUI, hdrCmdX+1, 0, CellTypeChar, 'C')
	gs.Set(bucket.GridMainUI, hdrCmdX+2, 0, CellTypeChar, 'O')
	gs.Set(bucket.GridMainUI, hdrCmdX+3, 0, CellTypeChar, 'M')
	gs.Set(bucket.GridMainUI, hdrCmdX+4, 0, CellTypeChar, 'M')
	gs.Set(bucket.GridMainUI, hdrCmdX+5, 0, CellTypeChar, 'A')
	gs.Set(bucket.GridMainUI, hdrCmdX+6, 0, CellTypeChar, 'N')
	gs.Set(bucket.GridMainUI, hdrCmdX+7, 0, CellTypeChar, 'D')
	gs.Set(bucket.GridMainUI, hdrCmdX+8, 0, CellTypeChar, 'S')

	// ...Output
	rightPanelHeaderX := cols - DividerX/2
	gs.Set(bucket.GridMainUI, rightPanelHeaderX+1, 0, CellTypeChar, 'O')
	gs.Set(bucket.GridMainUI, rightPanelHeaderX+2, 0, CellTypeChar, 'U')
	gs.Set(bucket.GridMainUI, rightPanelHeaderX+3, 0, CellTypeChar, 'T')
	gs.Set(bucket.GridMainUI, rightPanelHeaderX+4, 0, CellTypeChar, 'P')
	gs.Set(bucket.GridMainUI, rightPanelHeaderX+5, 0, CellTypeChar, 'U')
	gs.Set(bucket.GridMainUI, rightPanelHeaderX+6, 0, CellTypeChar, 'T')

	// ..Bottom
	for i := 1; i < cols-1; i++ {
		gs.SetCellSprite(bucket.GridMainUI, i, rows-1, assets.SpriteIDHorizontalBar)
	}

	// Dividers
	// .Vertical
	gs.SetCellSprite(bucket.GridMainUI, DividerX, 0, assets.SpriteIDDownConnectBar)
	for i := 1; i < rows-1; i++ {
		gs.SetCellSprite(bucket.GridMainUI, DividerX, i, assets.SpriteIDVerticalBar)
	}
	gs.SetCellSprite(bucket.GridMainUI, DividerX, rows-1, assets.SpriteIDUpConnectBar)
	// .Horizontal
	// verticalY := S3GridYCount - 11
	gs.SetCellSprite(bucket.GridMainUI, DividerX, DividerY, assets.SpriteIDRightConnectBar)
	for i := 1; i < cols-DividerX; i++ {
		gs.SetCellSprite(bucket.GridMainUI, DividerX+i, DividerY, assets.SpriteIDHorizontalBar)
	}
	gs.SetCellSprite(bucket.GridMainUI, cols-1, DividerY, assets.SpriteIDLeftConnectBar)
	// ...Logs
	gs.Set(bucket.GridMainUI, rightPanelHeaderX+1, DividerY, CellTypeChar, 'L')
	gs.Set(bucket.GridMainUI, rightPanelHeaderX+2, DividerY, CellTypeChar, 'O')
	gs.Set(bucket.GridMainUI, rightPanelHeaderX+3, DividerY, CellTypeChar, 'G')
	gs.Set(bucket.GridMainUI, rightPanelHeaderX+4, DividerY, CellTypeChar, 'S')

	/* Did not add these in.  I feel that should  be part of init */
	// CommandBuffer.AppendDecorators(CmdBufferDecor)
	// LogBuffer.AppendAll([]byte("Type: connect rabbit="))
}

func (bucket *Bucket) CommandBytes(commandId CommandId, bs *BufferSystem) []byte {
	row := bucket.CommandOutputStartRowIdx[commandId]
	d, err := bs.GetBufferRow(bucket.BufferCommandOutputs, row)

	if err == nil {
		return d
	}

	fmt.Printf("Error Retrieving Command Bytes: %d : %v\n", commandId, err)

	return []byte{}
}

func bucketInitCommandStrings(bs *BufferSystem, bucket *Bucket) {
	bucketAddCommandStrings(CommandFiles, []byte("[Files]:"), bs, bucket)
	bucketAddCommandStrings(CommandNetworks, []byte("[Networks]:"), bs, bucket)
	bucketAddCommandStrings(CommandConnect, []byte(":connect {name}="), bs, bucket)
	bucketAddCommandStrings(CommandList, []byte(":list="), bs, bucket)
	bucketAddCommandStrings(CommandListSpecific, []byte(":list {name}="), bs, bucket)
	bucketAddCommandStrings(CommandRoy1Fn, []byte("Email: Sorry"), bs, bucket)
	bucketAddCommandStrings(CommandRoy2Fn, []byte("Email: RE: Sorry"), bs, bucket)
	bucketAddCommandStrings(CommandRoy3Fn, []byte("Email: Miami"), bs, bucket)
	bucketAddCommandStrings(CommandLobby, []byte("Lobby"), bs, bucket)
	bucketAddCommandStrings(CommandSet, []byte(":set {gateId} {and;or}="), bs, bucket)
}

func bucketAddCommandStrings(commandId CommandId, val []byte, bs *BufferSystem, bucket *Bucket) {
	bucket.CommandOutputStartRowIdx[commandId] = bs.YCursor[bucket.BufferCommandOutputs] - 1

	bs.AppendAll(bucket.BufferCommandOutputs, val)
	bs.NewLine(bucket.BufferCommandOutputs)
}

func bucketInitializeGamePuzzles(pz *PuzzleSystem, ps *PathSystem, jxpp *JunctionSystem, gs *GridSystem, bucket *Bucket) error {
	/* Placeholder for now until puzzle generator is created */

	id, errorIntro := pz.AllocatePuzzle(1)

	if errorIntro != nil {
		return errorIntro
	}

	pz.IntroPuzzles[0] = id

	pz.SetValidGate(id, 0, 13, 16, GateOr)
	pz.SetPuzzleGate(id, 0, 13, 16, GateUnknown)
	pz.SetAttemptedGate(id, 0, 13, 16, GateUnknown)

	pz.PuzzleGateCounts[id] = 1

	cols := gs.Cols[bucket.GridOutputPrecision]
	rows := gs.Rows[bucket.GridOutputPrecision]

	// path1StartX, path1StartY := gs.GridXYToScreenSpace(bucket.GridOutput, cols/2-1, rows+1)
	path1, err := ps.NewPath(cols/2+1, rows-1, 0.0, 0.0)

	if err != nil {
		return err
	}

	jxpp.AddParent(uint32(id), 1)
	jxpp.AddChild(uint32(id), uint32(path1))

	return nil
}
