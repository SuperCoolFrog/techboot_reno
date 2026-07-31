package main

// Output Grid sizes Grid26x36x20

/**


**/

type PuzzleId uint32

type GateType uint32

const (
	GateAnd GateType = iota
	GateOr
)

type PuzzleSystem struct {
	TotalPuzzles int
	MasterChunk  []byte

	Gates        []GateType
	GatesOffsets []int

	GateCellType    []GridCellType
	GateX           []int
	GateY           []int
	GateCellOffsets []int

	Currents        []int
	CurrentsOffsets []int
	CurrentFromCol  []int
	CurrentFromRow  []int
	CurrentToCol    []int
	CurrentToRow    []int

	Cols int // Match Output Cols
	Rows int // Match Output Rows

	NextPuzzleId PuzzleId
}

func NewPuzzleSystem(easyPuzzlesCount, medPuzzlesCount, hardPuzzlesCount int) *PuzzleSystem {
	ps := &PuzzleSystem{
		TotalPuzzles: easyPuzzlesCount + medPuzzlesCount + hardPuzzlesCount,
		Cols:         26,
		Rows:         36,
	}

	return ps
}

func (ps *PuzzleSystem) GetRandomPuzzle(puzzleCommandId CommandId) {}
