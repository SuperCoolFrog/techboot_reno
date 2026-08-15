package main

/*
Placeholder for now until puzzle generator is created

** Be sure that bucket is initialized
*/
func GeneratePuzzles(game *Game) error {
	errIntros := genIntroPuzzles(game)
	if errIntros != nil {
		return errIntros
	}

	return nil
}

/*
	type PuzzleType uint8
	0,1,2,3


	Puzzles:
		ConfigId, State, PuzzleType, GateCount, MarkerCount

	Gates:
		ConfigId, X, Y




	TotalPuzzles
*/

func genIntroPuzzles(game *Game) error {
	// Scene -> Puzzle Count
	game.jxsp.AddParent(uint32(Scene4_Puzzle), 2)

	// #region Intro Puzzle 1

	id, errorIntro := game.pz.AllocatePuzzle(1, 1)
	if errorIntro != nil {
		return errorIntro
	}

	game.pz.IntroPuzzles[0] = id
	game.pac.Set(uint16(id), false)
	game.pcn.Set(uint16(id), false)

	game.jxsp.AddChild(uint32(Scene4_Puzzle), uint32(id))

	cols := game.gs.Cols[game.b.GridOutput]
	rows := game.gs.Rows[game.b.GridOutput]
	g1X := cols / 2
	g1Y := rows / 2

	gateIdx0 := 0
	gateId0 := game.pz.GetGateId(id, gateIdx0)

	game.pz.SetValidGate(id, gateIdx0, g1X, g1Y, GatePass)
	game.pz.SetPuzzleGate(id, gateIdx0, g1X, g1Y, GateUnknown)
	game.pz.SetAttemptedGate(id, gateIdx0, g1X, g1Y, GateUnknown)
	game.pz.SetMarker(id, 0, cols/2, 0, PuzzleMarkerYes)

	game.pz.PuzzleGateCounts[id] = 1

	game.jxpp.AddParent(uint32(id), 1)

	pathIn, err := game.ps.NewPath(cols/2, rows-1, cols/2, g1Y+2)
	if err != nil {
		return err
	}
	game.jxpp.AddChild(uint32(id), uint32(pathIn))

	// Gate Paths

	// @TODO this path gets added to gate1+GateTypePass for jxgp
	pathPass, errpj := game.ps.NewPath(cols/2, g1Y-2, cols/2, 0)
	if errpj != nil {
		return errpj
	}
	game.jxgp.AddChild(uint16(gateId0), uint16(GatePass), uint16(pathPass))
	game.gap.Set(uint16(gateId0), false)
	game.gac.Set(uint16(gateId0), false)

	// #endregion Intro Puzzle 1

	return nil
}
