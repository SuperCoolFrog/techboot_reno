package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"math"
	"techboot_reno/cmd/assets"
)

func Scene4_HandleInit(current, next GameState, game *Game) GameState {

	game.gs.EnableGrid(game.b.GridOutput)

	game.bs.AppendAll(game.b.BufferLogs, []byte("Connecting..."))
	game.bs.NewLine(game.b.BufferLogs)
	game.bs.DrawToGrid(game.b.BufferLogs, game.b.GridMainUI, game.b.LogBufferColIdx, game.b.LogBufferRowIdx, game.gs)

	game.ans.IsPlaying[AnimationMemoryStack] = true
	game.ans.Loop[AnimationMemoryStack] = false
	game.ans.Durations[AnimationMemoryStack] = 5.0
	game.ans.Timers[AnimationMemoryStack] = 0.0
	game.ans.Delay[AnimationMemoryStack] = 0

	return next
}

func Scene4_UpdateStackAnimation(current, next GameState, game *Game) GameState {
	timer := game.ans.Timers[AnimationMemoryStack]
	duration := game.ans.Durations[AnimationMemoryStack]
	delay := game.ans.Delay[AnimationMemoryStack]

	trueTime := float32(math.Max(float64(timer-delay), 0))
	completedTime := trueTime / duration

	// gs.SetAllCells(AnimationScannerGrid, CellTypeEmpty, 0)
	OutputGridCols := game.gs.Cols[game.b.GridOutput]
	OutputGridRows := game.gs.Rows[game.b.GridOutput]

	x := OutputGridCols / 2
	y := int(float32(OutputGridRows-1) * completedTime)
	for i := 0; i < OutputGridRows; i++ {
		iy := OutputGridRows - 1 - y
		game.gs.SetCellSprite(game.b.GridOutput, x-1, iy, assets.SpriteIDSquare)
		game.gs.SetCellSprite(game.b.GridOutput, x, iy, assets.SpriteIDSquare)
		game.gs.SetCellSprite(game.b.GridOutput, x+1, iy, assets.SpriteIDSquare)
	}

	if game.ans.IsPlaying[AnimationMemoryStack] {
		game.bs.DrawToGrid(game.b.BufferLogs, game.b.GridMainUI, game.b.LogBufferColIdx, game.b.LogBufferRowIdx, game.gs)
		return current
	}

	// LogBuffer.AppendAll([]byte("Connection Successful"))
	game.bs.AppendAll(game.b.BufferLogs, []byte("Connection Successful"))
	game.bs.NewLine(game.b.BufferLogs)

	game.gs.SetAllCells(game.b.GridOutput, CellTypeEmpty, 0)

	return next
}

func Scene4_Update(current, next GameState, game *Game) GameState {
	for i := 0; i < len(game.inputRunes); i++ {
		err := game.bs.AppendWithDecor(game.b.BufferCommands, byte(game.inputRunes[i]), CmdBufferDecor)
		if err != nil {
			fmt.Printf("Error Appending Last Rune: %v\n", err)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		game.bs.TrimDecor(game.b.BufferCommands, CmdBufferDecor)
		game.bs.NewLine(game.b.BufferCommands)
		game.bs.AppendDecorators(game.b.BufferCommands, CmdBufferDecor)

		if lastLine, lineError := game.bs.GetLastBufferLine(game.b.BufferCommands); lineError == nil {
			ParseInput(lastLine, game.prologInput)
		} else {
			fmt.Printf("Error Getting lastLine: %v", lineError)
		}
	}
	if utilDebouncedKeyPressed(ebiten.KeyBackspace) {
		game.bs.DecrementCursorWithDecor(game.b.BufferCommands, CmdBufferDecor)
	}

	err := game.bs.DrawToGridWithDecor(game.b.BufferCommands, game.b.GridMainUI, 1, 1, CmdBufferDecor, game.gs)
	if err != nil {
		fmt.Printf("Error Drawing to BufferCommands: %v", err)
	}

	// LogBuff
	game.bs.DrawToGrid(game.b.BufferLogs, game.b.GridMainUI, 37, 38, game.gs)

	state := current

loop:
	for {
		select {
		case cmd := <-game.prologOutput:

			switch cmd.ResultType {
			case CommandList:
				game.gs.SetAllCells(game.b.GridOutput, CellTypeEmpty, 0)

				game.gs.SetRowCells(game.b.GridOutput, 0, CellTypeChar, cmd.Command)
				game.gs.SetRowSprites(game.b.GridOutput, 1, assets.SpriteIDHorizontalBar)

				for i := 0; i < len(cmd.Items); i++ {
					b := game.b.CommandBytes(cmd.Items[i], game.bs)
					game.gs.SetCellSprite(game.b.GridOutput, 0, i+2, assets.SpriteIDRightConnectBar)
					game.gs.SetRowCellsAtCol(game.b.GridOutput, 1, i+2, CellTypeChar, b)
				}
			case CommandPuzzle:
				state = next
				game.bs.AppendAll(game.b.BufferLogs, []byte("Connection Failed"))
				game.bs.NewLine(game.b.BufferLogs)
			case CommandConnectFalse:
				game.bs.AppendAll(game.b.BufferLogs, []byte("Connection Failed"))
				game.bs.NewLine(game.b.BufferLogs)
			case CommandInvalid:
				game.bs.AppendAll(game.b.BufferLogs, []byte("Invalid Command"))
				game.bs.NewLine(game.b.BufferLogs)
			default:
				game.bs.AppendAll(game.b.BufferLogs, []byte("Command did nothing"))
				game.bs.NewLine(game.b.BufferLogs)
			}
		default:
			break loop // nothing left in the queue for this frame
		}
	}

	return state
}

func Scene4_SetupPuzzle(current, next GameState, game *Game) GameState {
	game.gs.SetAllCells(game.b.GridOutput, CellTypeEmpty, 0)

	var puzzleId PuzzleId

	if !game.pz.HasPuzzleAssignment(next) {
		puzzleId, pzError := game.pz.GetUnassignedIntroPuzzle()

		if pzError != nil {
			panic(pzError)
		}

		game.pz.AssignPuzzle(puzzleId, next)
	}

	s4DrawPuzzle(puzzleId, game)

	s4StartPathAnimation(puzzleId, game)

	return next
}

func s4DrawPuzzle(puzzleId PuzzleId, game *Game) {
	gates := game.pz.GetPuzzleGates(puzzleId)

	for i := 0; i < len(gates); i++ {
		game.pz.DrawGate(puzzleId, i, game.b.GridOutput, game.gs)
	}

	paths, errC := game.jxpp.GetChildren(uint32(puzzleId))
	if errC != nil {
		panic(errC)
	}

	for i := 0; i < len(paths); i++ {
		pathId := paths[i]
		pX := game.ps.StartX[pathId]
		pY := game.ps.StartY[pathId]

		game.gs.SetCellSprite(game.b.GridOutput, pX, pY, assets.SpriteIDCarrotUp)
	}

	game.gs.SetCellSprite(game.b.GridOutput, game.gs.Cols[game.b.GridOutput]/2, 0, assets.SpriteIDCarrotUp)
}

func s4StartPathAnimation(puzzleId PuzzleId, game *Game) {
	game.ans.IsPlaying[AnimationPath] = true
	game.ans.Timers[AnimationPath] = 0.0

}

func s4ResetFlags(puzzleId PuzzleId, game *Game) {
	game.pac.Set(uint16(puzzleId), false)

	gates := game.pz.GetAttemptGates(puzzleId)
	for i := 0; i < len(gates); i++ {
		gateId := game.pz.GetGateId(puzzleId, i)
		gateId16 := uint16(gateId)
		game.gac.Set(gateId16, false)
		game.gap.Set(gateId16, false)
	}
}

func s4AnimateGatePaths(puzzleId PuzzleId, game *Game) (isPlaying bool) {
	gates := game.pz.GetAttemptGates(puzzleId)

	for i := 0; i < len(gates); i++ {
		gateId := game.pz.GetGateId(puzzleId, i)
		gateId16 := uint16(gateId)
		gate := gates[i]

		if game.gac.Has(gateId16) {
			continue
		}

		paths := game.jxgp.GetChildren(uint16(gateId), uint16(gate))

		if len(paths) == 0 {
			game.gac.Set(gateId16, true)
			continue
		}

		if !game.gap.Has(gateId16) {
			s4StartPathAnimation(puzzleId, game)
			game.gap.Set(gateId16, true)
		}

		if !game.ans.IsPlaying[AnimationPath] {
			game.gac.Set(gateId16, true)
		}

		timer := game.ans.Timers[AnimationPath]
		duration := game.ans.Durations[AnimationPath]
		delay := game.ans.Delay[AnimationPath]

		trueTime := float32(math.Max(float64(timer-delay), 0))
		completedTime := trueTime / duration

		for i := 0; i < len(paths); i++ {
			pathId := paths[i]
			x1 := game.ps.StartX[pathId]
			y1 := game.ps.StartY[pathId]

			xPts, yPts := game.ps.GetXYPoints(PathId(pathId))
			count := game.ps.PointsCount[pathId]

			end := int(completedTime * float32(count+1)) // need to add +1 because for uses <, so need to include top idx

			for i := 0; i < end; i++ {
				x := xPts[i]
				y := yPts[i]

				game.gs.SetCellSprite(game.b.GridOutput, x, y, assets.SpriteIDSquare)
			}

			game.gs.SetCellSprite(game.b.GridOutput, x1, y1, assets.SpriteIDCarrotUp)
		}

		// Gate is still playing path animations
		return true
	}

	return false
}

func s4AnimatePath(puzzleId PuzzleId, state GameState, game *Game) (isPlaying bool) {
	if !game.ans.IsPlaying[AnimationPath] {
		game.pac.Set(uint16(puzzleId), true)
	}

	if game.pac.Has(uint16(puzzleId)) {
		return s4AnimateGatePaths(puzzleId, game)
	}

	timer := game.ans.Timers[AnimationPath]
	duration := game.ans.Durations[AnimationPath]
	delay := game.ans.Delay[AnimationPath]

	trueTime := float32(math.Max(float64(timer-delay), 0))
	completedTime := trueTime / duration

	paths, errC := game.jxpp.GetChildren(uint32(puzzleId))
	if errC != nil {
		panic(errC)
	}

	for i := 0; i < len(paths); i++ {
		pathId := paths[i]
		x1 := game.ps.StartX[pathId]
		y1 := game.ps.StartY[pathId]

		xPts, yPts := game.ps.GetXYPoints(PathId(pathId))
		count := game.ps.PointsCount[pathId]

		end := int(completedTime * float32(count+1)) // need to add +1 because for uses <, so need to include top idx

		for i := 0; i < end; i++ {
			x := xPts[i]
			y := yPts[i]

			game.gs.SetCellSprite(game.b.GridOutput, x, y, assets.SpriteIDSquare)
		}

		game.gs.SetCellSprite(game.b.GridOutput, x1, y1, assets.SpriteIDCarrotUp)
	}

	return game.ans.IsPlaying[AnimationPath]
}

func s4ClearPathSprites(puzzleId PuzzleId, game *Game) {
	paths, errC := game.jxpp.GetChildren(uint32(puzzleId))

	if errC != nil {
		panic(errC)
	}

	for i := 0; i < len(paths); i++ {
		pathId := paths[i]
		xPts, yPts := game.ps.GetXYPoints(PathId(pathId))

		for i := 0; i < len(xPts); i++ {
			x := xPts[i]
			y := yPts[i]

			game.gs.Set(game.b.GridOutput, x, y, CellTypeEmpty, ' ')
		}
	}

	s4DrawPuzzle(puzzleId, game)
}

func Scene4_Puzzling(current, next GameState, game *Game) GameState {
	puzzleId, assignmentError := game.pz.GetPuzzleAssignment(current)
	if assignmentError != nil {
		panic(fmt.Errorf("Error getting puzzle assigment: %v\n", assignmentError))
	}

	if s4AnimatePath(puzzleId, current, game) {
		return current
	} else if game.pcn.Has(uint16(puzzleId)) { // lock next until connect is called
		game.bs.AppendAll(game.b.BufferLogs, []byte("Connection Successful"))
		game.bs.NewLine(game.b.BufferLogs)
		game.bs.DrawToGrid(game.b.BufferLogs, game.b.GridMainUI, 37, 38, game.gs)
		return next
	} else {
		s4ClearPathSprites(puzzleId, game)
	}

	for i := 0; i < len(game.inputRunes); i++ {
		err := game.bs.AppendWithDecor(game.b.BufferCommands, byte(game.inputRunes[i]), CmdBufferDecor)
		if err != nil {
			fmt.Printf("Error Appending Last Rune: %v\n", err)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		game.bs.TrimDecor(game.b.BufferCommands, CmdBufferDecor)
		game.bs.NewLine(game.b.BufferCommands)
		game.bs.AppendDecorators(game.b.BufferCommands, CmdBufferDecor)

		if lastLine, lineError := game.bs.GetLastBufferLine(game.b.BufferCommands); lineError == nil {
			ParseInput(lastLine, game.prologInput)
		} else {
			fmt.Printf("Error Getting lastLine: %v", lineError)
		}
	}
	if utilDebouncedKeyPressed(ebiten.KeyBackspace) {
		game.bs.DecrementCursorWithDecor(game.b.BufferCommands, CmdBufferDecor)
	}

	err := game.bs.DrawToGridWithDecor(game.b.BufferCommands, game.b.GridMainUI, 1, 1, CmdBufferDecor, game.gs)
	if err != nil {
		fmt.Printf("Error Drawing to BufferCommands: %v", err)
	}

	// LogBuff
	game.bs.DrawToGrid(game.b.BufferLogs, game.b.GridMainUI, 37, 38, game.gs)

	state := current

loop:
	for {
		select {
		case cmd := <-game.prologOutput:
			switch cmd.ResultType {
			case CommandList:
				game.gs.SetAllCells(game.b.GridOutput, CellTypeEmpty, 0)

				game.gs.SetRowCells(game.b.GridOutput, 0, CellTypeChar, cmd.Command)
				game.gs.SetRowSprites(game.b.GridOutput, 1, assets.SpriteIDHorizontalBar)

				for i := 0; i < len(cmd.Items); i++ {
					b := game.b.CommandBytes(cmd.Items[i], game.bs)
					game.gs.SetCellSprite(game.b.GridOutput, 0, i+2, assets.SpriteIDRightConnectBar)
					game.gs.SetRowCellsAtCol(game.b.GridOutput, 1, i+2, CellTypeChar, b)
				}
			case CommandPuzzle:
				game.gs.SetAllCells(game.b.GridOutput, CellTypeEmpty, 0)
				s4DrawPuzzle(puzzleId, game)
				s4ResetFlags(puzzleId, game)
				s4StartPathAnimation(puzzleId, game)

				if game.pz.IsPuzzleSolved(puzzleId) {
					game.pcn.Set(uint16(puzzleId), true)
				} else {
					game.bs.AppendAll(game.b.BufferLogs, []byte("Connection Failed"))
					game.bs.NewLine(game.b.BufferLogs)
				}

			case CommandSet:
				puzzleId, assignmentError := game.pz.GetPuzzleAssignment(current)

				if assignmentError != nil {
					fmt.Printf("Error getting puzzle assigment: %v\n", assignmentError)
					break loop
				}

				for i := 0; i < len(cmd.Items); i++ {
					var gateType GateType

					switch CommandId(cmd.Items[i]) {

					case CommandGateTypeJoin:
						gateType = GateJoin
					case CommandGateTypePair:
						gateType = GatePair
					case CommandGateTypeSplit:
						gateType = GateSplit
					case CommandGateTypePass:
						gateType = GatePass
					default:
						gateType = GateUnknown
					}

					gateIdx := cmd.ValuesInt[i]

					err := game.pz.SetAttemptGate(puzzleId, gateIdx, gateType)
					if err != nil {
						fmt.Printf("Error setting gatetype: %d ; %d ; %d ;\n%v\n", puzzleId, gateIdx, gateType, err)
					}

					game.pz.DrawGate(puzzleId, i, game.b.GridOutput, game.gs)
				}
			case CommandInvalid:
				// fmt.Printf("Invalid!\n")
				game.bs.AppendAll(game.b.BufferLogs, []byte("Invalid Command"))
				game.bs.NewLine(game.b.BufferLogs)
			default:
				game.bs.AppendAll(game.b.BufferLogs, []byte("Command did nothing"))
				game.bs.NewLine(game.b.BufferLogs)
			}
		default:
			break loop // nothing left in the queue for this frame
		}
	}

	return state
}
