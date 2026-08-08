package main

import (
	"fmt"
	"techboot_reno/cmd/assets"
	"unsafe"
)

type PuzzleId uint32

type GateType uint32

const (
	GateEmpty GateType = iota
	GateUnknown
	GateJoin
	GatePair
	GateSplit
	GatePass
)

type PuzzleSystem struct {
	TotalPuzzles int
	TotalGates   int
	MasterChunk  []byte

	IntroPuzzlesCount int
	EasyPuzzlesCount  int
	MedPuzzlesCount   int
	HardPuzzlesCount  int

	IntroPuzzles []PuzzleId
	EasyPuzzles  []PuzzleId
	MedPuzzles   []PuzzleId
	HardPuzzles  []PuzzleId

	// Per-Puzzle
	PuzzleGateCounts []int // Tracks exactly how many gates a specific PuzzleId has
	PuzzleIsAssigned []bool
	PuzzleAssignment []GameState

	// Stored in Chunk
	ValidGates     []GateType // Valid Gates that solve the puzzle
	PuzzleGates    []GateType // This is the gates presented to player
	AttemptedGates []GateType // This is the player's last attempt
	GatesOffsets   []int
	GateX          []int
	GateY          []int

	NextPuzzleId PuzzleId
}

func NewPuzzleSystem(totalGates, introPuzzles, easyPuzzles, medPuzzles, hardPuzzles int) *PuzzleSystem {
	totalPuzzles := introPuzzles + easyPuzzles + medPuzzles + hardPuzzles

	ps := &PuzzleSystem{
		TotalPuzzles:      totalPuzzles,
		TotalGates:        totalGates,
		IntroPuzzlesCount: introPuzzles,
		EasyPuzzlesCount:  easyPuzzles,
		MedPuzzlesCount:   medPuzzles,
		HardPuzzlesCount:  hardPuzzles,
		IntroPuzzles:      make([]PuzzleId, introPuzzles),
		EasyPuzzles:       make([]PuzzleId, easyPuzzles),
		MedPuzzles:        make([]PuzzleId, medPuzzles),
		HardPuzzles:       make([]PuzzleId, hardPuzzles),

		NextPuzzleId: 0,

		PuzzleGateCounts: make([]int, totalPuzzles),
		PuzzleIsAssigned: make([]bool, totalPuzzles),
		PuzzleAssignment: make([]GameState, totalPuzzles),

		// Standard Go slice allocations for non-contiguous offsets
		GatesOffsets: make([]int, totalGates),
	}

	// 1. Compute Byte Sizes Grouped by Memory Alignment Requirements

	// 8-Byte Types (ints)
	sizeGateX := totalGates * int(unsafe.Sizeof(int(0)))
	sizeGateY := totalGates * int(unsafe.Sizeof(int(0)))

	// 4-Byte Types (uint32)
	sizeGates := totalGates * int(unsafe.Sizeof(GateType(0)))

	// 2. Allocate the Master Chunk
	totalSize := sizeGateX + sizeGateY + (sizeGates * 3) // Valid,Puzzle,Attempted,

	ps.MasterChunk = make([]byte, totalSize)
	ptr := unsafe.Pointer(&ps.MasterChunk[0])

	// 3. Slice 8-Byte Fields First (Strictest Alignment)
	ps.GateX = unsafe.Slice((*int)(ptr), totalGates)
	ptr = unsafe.Add(ptr, sizeGateX)

	ps.GateY = unsafe.Slice((*int)(ptr), totalGates)
	ptr = unsafe.Add(ptr, sizeGateY)

	// 4. Slice 4-Byte Fields (Guaranteed 4-byte aligned because previous sizes are multiples of 8)
	ps.ValidGates = unsafe.Slice((*GateType)(ptr), totalGates)
	ptr = unsafe.Add(ptr, sizeGates)

	ps.PuzzleGates = unsafe.Slice((*GateType)(ptr), totalGates)
	ptr = unsafe.Add(ptr, sizeGates)

	ps.AttemptedGates = unsafe.Slice((*GateType)(ptr), totalGates)
	// ptr = unsafe.Add(ptr, sizeGates)

	return ps
}

func (ps *PuzzleSystem) AllocatePuzzle(numberOfGates int) (PuzzleId, error) {
	if int(ps.NextPuzzleId) >= ps.TotalPuzzles {
		return 0, fmt.Errorf("Max Puzzles Reach, Unable to create new puzzle")
	}

	id := ps.NextPuzzleId

	ps.NextPuzzleId++

	var startOffset int = 0

	if id > 0 {
		prevID := id - 1
		startOffset = ps.GatesOffsets[prevID] + ps.PuzzleGateCounts[prevID]
	}

	if startOffset+numberOfGates > ps.TotalGates {
		return 0, fmt.Errorf("Number of gates(%d) would exceed allotted TotalGates (%d)", numberOfGates, ps.TotalGates)
	}

	ps.PuzzleGateCounts[id] = numberOfGates
	ps.GatesOffsets[id] = startOffset

	return id, nil
}

func (ps *PuzzleSystem) SetValidGate(puzzleId PuzzleId, gateIdx, gateX, gateY int, gate GateType) error {
	offset := ps.GatesOffsets[puzzleId]

	if gateIdx >= ps.PuzzleGateCounts[puzzleId] {
		return fmt.Errorf("GateIdx out of bounds: Given %d ; Count %d", gateIdx, ps.PuzzleGateCounts[puzzleId])
	}

	idx := offset + gateIdx
	ps.ValidGates[idx] = gate
	ps.GateX[idx] = gateX
	ps.GateY[idx] = gateY

	return nil
}

func (ps *PuzzleSystem) SetPuzzleGate(puzzleId PuzzleId, gateIdx, gateX, gateY int, gate GateType) error {
	offset := ps.GatesOffsets[puzzleId]

	if gateIdx >= ps.PuzzleGateCounts[puzzleId] {
		return fmt.Errorf("GateIdx out of bounds: Given %d ; Count %d", gateIdx, ps.PuzzleGateCounts[puzzleId])
	}

	idx := offset + gateIdx
	ps.PuzzleGates[idx] = gate
	ps.GateX[idx] = gateX
	ps.GateY[idx] = gateY

	return nil
}

func (ps *PuzzleSystem) SetGateType(puzzleId PuzzleId, gateIdx int, gate GateType) error {
	offset := ps.GatesOffsets[puzzleId]

	if gateIdx >= ps.PuzzleGateCounts[puzzleId] {
		return fmt.Errorf("GateIdx out of bounds: Given %d ; Count %d", gateIdx, ps.PuzzleGateCounts[puzzleId])
	}

	idx := offset + gateIdx
	ps.PuzzleGates[idx] = gate

	return nil
}

func (ps *PuzzleSystem) SetAttemptedGate(puzzleId PuzzleId, gateIdx, gateX, gateY int, gate GateType) error {
	offset := ps.GatesOffsets[puzzleId]

	if gateIdx >= ps.PuzzleGateCounts[puzzleId] {
		return fmt.Errorf("GateIdx out of bounds: Given %d ; Count %d", gateIdx, ps.PuzzleGateCounts[puzzleId])
	}

	idx := offset + gateIdx
	ps.AttemptedGates[idx] = gate
	ps.GateX[idx] = gateX
	ps.GateY[idx] = gateY

	return nil
}

func (ps *PuzzleSystem) GetValidGates(puzzleId PuzzleId) []GateType {
	offset := ps.GatesOffsets[puzzleId]
	count := ps.PuzzleGateCounts[puzzleId]

	start := offset
	end := start + count

	return ps.ValidGates[start:end]
}

func (ps *PuzzleSystem) GetPuzzleGates(puzzleId PuzzleId) []GateType {
	offset := ps.GatesOffsets[puzzleId]
	count := ps.PuzzleGateCounts[puzzleId]

	start := offset
	end := start + count

	return ps.PuzzleGates[start:end]
}

func (ps *PuzzleSystem) GetAttemptGates(puzzleId PuzzleId) []GateType {
	offset := ps.GatesOffsets[puzzleId]
	count := ps.PuzzleGateCounts[puzzleId]

	start := offset
	end := start + count

	return ps.AttemptedGates[start:end]
}

func (ps *PuzzleSystem) GetGatesX(puzzleId PuzzleId) []int {
	offset := ps.GatesOffsets[puzzleId]
	count := ps.PuzzleGateCounts[puzzleId]

	start := offset
	end := start + count

	return ps.GateX[start:end]
}

func (ps *PuzzleSystem) GetGatesY(puzzleId PuzzleId) []int {
	offset := ps.GatesOffsets[puzzleId]
	count := ps.PuzzleGateCounts[puzzleId]

	start := offset
	end := start + count

	return ps.GateY[start:end]
}

func (ps *PuzzleSystem) IsPuzzleSolved(puzzleId PuzzleId) bool {
	attempt := ps.GetAttemptGates(puzzleId)
	valid := ps.GetValidGates(puzzleId)

	for i := 0; i < len(attempt); i++ {
		if attempt[i] != valid[i] {
			return false
		}
	}

	return true
}

func (ps *PuzzleSystem) GetUnassignedIntroPuzzle() (PuzzleId, error) {
	for i := 0; i < len(ps.IntroPuzzles); i++ {
		id := ps.IntroPuzzles[i]
		if !ps.PuzzleIsAssigned[id] {
			return id, nil
		}
	}

	return 0, fmt.Errorf("No available intro puzzle")
}

// @TODO move this into junction system
func (ps *PuzzleSystem) AssignPuzzle(puzzleId PuzzleId, state GameState) {
	ps.PuzzleIsAssigned[puzzleId] = true
	ps.PuzzleAssignment[puzzleId] = state
}

func (ps *PuzzleSystem) GetPuzzleAssignment(state GameState) (PuzzleId, error) {
	for i := PuzzleId(0); i < PuzzleId(ps.NextPuzzleId); i++ {
		if ps.PuzzleIsAssigned[i] && ps.PuzzleAssignment[i] == state {
			return i, nil
		}
	}

	return 0, fmt.Errorf("No assignment for state")
}

func (ps *PuzzleSystem) HasPuzzleAssignment(state GameState) bool {
	for i := PuzzleId(0); i < PuzzleId(ps.NextPuzzleId); i++ {
		if ps.PuzzleIsAssigned[i] && ps.PuzzleAssignment[i] == state {
			return true
		}
	}

	return false
}

func (ps *PuzzleSystem) DrawGate(puzzleId PuzzleId, gateIdx int, gridId GridID, gs *GridSystem) error {

	gates := ps.GetPuzzleGates(puzzleId)
	gatesX := ps.GetGatesX(puzzleId)
	gatesY := ps.GetGatesY(puzzleId)

	for i := 0; i < len(gates); i++ {
		gateType := gates[i]
		x := gatesX[i]
		y := gatesY[i]

		gs.SetCellSprite(gridId, x, y-1, assets.SpriteIDHorizontalBar)
		gs.SetCellSprite(gridId, x, y+1, assets.SpriteIDHorizontalBar)
		gs.SetCellSprite(gridId, x-1, y-1, assets.SpriteIDHorizontalBar)
		gs.SetCellSprite(gridId, x+1, y+1, assets.SpriteIDHorizontalBar)
		gs.SetCellSprite(gridId, x+1, y-1, assets.SpriteIDHorizontalBar)
		gs.SetCellSprite(gridId, x-1, y+1, assets.SpriteIDHorizontalBar)

		gs.SetCellSprite(gridId, x-2, y, assets.SpriteIDVerticalBar)
		gs.SetCellSprite(gridId, x+2, y, assets.SpriteIDVerticalBar)

		if i < 9 {
			gs.Set(gridId, x-3, y, CellTypeChar, byte('0'+(i+1)))
		}

		hasSymbol := false
		var cornerSprite assets.SpriteID

		switch gateType {
		case GateUnknown:
			gs.Set(gridId, x, y, CellTypeChar, '?')
		case GateJoin:
			gs.Set(gridId, x, y, CellTypeChar, 'J')
			hasSymbol = true
			cornerSprite = assets.SpriteIDSquare
		case GatePair:
			gs.Set(gridId, x, y, CellTypeChar, 'P')
			hasSymbol = true
			cornerSprite = assets.SpriteIDDiamond
		case GateSplit:
			gs.Set(gridId, x, y, CellTypeChar, 'S')
			hasSymbol = true
			cornerSprite = assets.SpriteIDCarrotUp
		case GatePass:
			gs.Set(gridId, x, y, CellTypeChar, 'O')
			hasSymbol = true
			cornerSprite = assets.SpriteIDCircle
		}

		if hasSymbol {
			gs.SetCellSprite(gridId, x-2, y-1, cornerSprite)
			gs.SetCellSprite(gridId, x-2, y+1, cornerSprite)
			gs.SetCellSprite(gridId, x+2, y-1, cornerSprite)
			gs.SetCellSprite(gridId, x+2, y+1, cornerSprite)
		}
	}

	return nil
}
