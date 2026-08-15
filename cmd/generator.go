package main

import (
	"context"
	"fmt"

	"github.com/trealla-prolog/go/trealla"
)

/**
import (
	"context"
	"fmt"

	"github.com/trealla-prolog/go/trealla"
	"strings"
)

type CommandResponse struct {
	ResultType CommandId
	Items      []CommandId
	ValuesInt  []int
	Command    []byte
}
**/

type PuzzleConfig struct {
	trealla.Functor

	ConfigId     int
	GameState    int
	PuzzleTypeId int
	GateCount    int
	MarkerCount  int

	PuzzlePaths  [][]int
	ValidGates   [][]int
	PuzzleGates  [][]int
	AttemptGates [][]int
	GatePaths    [][]int
}

type PuzzleConfigResult struct {
	Configs []PuzzleConfig `prolog:"Configs"`
}

/*
Placeholder for now until puzzle generator is created

** Be sure that bucket is initialized
*/
func GeneratePuzzles(g *Game) error {
	// Initialize a clean, zero-dependency Trealla WebAssembly instance
	pl, err := trealla.New()
	if err != nil {
		return fmt.Errorf("Failed to boot Trealla: %v\n", err)
	}

	ctx := context.Background()

	if err := pl.ConsultText(ctx, "user", g.puzzlespl); err != nil {
		panic(err)
	}

	fmt.Printf("✅ Puzzles loaded and ready for query!\n")

	queryStr := "all_configs(Configs)."
	q, err11 := pl.QueryOnce(ctx, queryStr)
	if err11 != nil {
		panic(err11)
	}

	var r PuzzleConfigResult

	if errv := q.Solution.Scan(&r); errv != nil {
		panic(errv)
	}

	fmt.Printf("Loaded Puzzle Data: %v\n", r)

	errIntros := genIntroPuzzles(g)
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
