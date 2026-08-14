package main

import (
	"fmt"
	"techboot_reno/cmd/assets"
	"unsafe"
)

type PuzzleId uint16

type GateId uint16 // offset+idx
type GateType uint16

const (
	GateEmpty GateType = iota
	GateUnknown
	GateJoin
	GateSplit
	GatePass

	GateTypeCount
)

type PuzzleMarker byte

const (
	PuzzleMarkerNone PuzzleMarker = 0
	PuzzleMarkerYes  PuzzleMarker = 1 << 0
	PuzzleMarkerNo   PuzzleMarker = 2 << 0
)

type PuzzleSystem struct {
	TotalPuzzles int
	TotalGates   int
	TotalMarkers int
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

	MarkersCounts []int

	// Stored in Chunk
	ValidGates     []GateType // Valid Gates that solve the puzzle
	PuzzleGates    []GateType // This is the gates presented to player
	AttemptedGates []GateType // This is the player's last attempt
	GatesOffsets   []int
	GateX          []int
	GateY          []int

	Markers        []PuzzleMarker
	MarkerX        []int
	MarkerY        []int
	MarkersOffsets []int

	NextPuzzleId PuzzleId
}

func NewPuzzleSystem(totalGates, maxMarkersPerPuzzle, introPuzzles, easyPuzzles, medPuzzles, hardPuzzles int) *PuzzleSystem {
	totalPuzzles := introPuzzles + easyPuzzles + medPuzzles + hardPuzzles
	totalMarkers := totalPuzzles + maxMarkersPerPuzzle

	ps := &PuzzleSystem{
		TotalPuzzles:      totalPuzzles,
		TotalGates:        totalGates,
		TotalMarkers:      totalMarkers,
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
		MarkersCounts:    make([]int, totalPuzzles),

		// Standard Go slice allocations for non-contiguous offsets
		GatesOffsets:   make([]int, totalPuzzles),
		MarkersOffsets: make([]int, totalPuzzles),
	}

	// 1. Compute Byte Sizes Grouped by Memory Alignment Requirements

	// 8-Byte Types (ints)
	sizeGateX := totalGates * int(unsafe.Sizeof(int(0)))
	sizeGateY := totalGates * int(unsafe.Sizeof(int(0)))

	sizeMarkerX := totalMarkers * int(unsafe.Sizeof(int(0)))
	sizeMarkerY := totalMarkers * int(unsafe.Sizeof(int(0)))

	// 4-Byte Types (uint32)
	sizeGates := totalGates * int(unsafe.Sizeof(GateType(0)))

	sizeMarkers := totalMarkers * int(unsafe.Sizeof(PuzzleMarker(0)))

	// 2. Allocate the Master Chunk
	totalSize := sizeMarkers + sizeMarkerX + sizeMarkerY +
		sizeGateX + sizeGateY + (sizeGates * 3) // Valid,Puzzle,Attempted,

	ps.MasterChunk = make([]byte, totalSize)
	ptr := unsafe.Pointer(&ps.MasterChunk[0])

	// 3. Slice 8-Byte Fields First (Strictest Alignment)
	ps.GateX = unsafe.Slice((*int)(ptr), totalGates)
	ptr = unsafe.Add(ptr, sizeGateX)

	ps.GateY = unsafe.Slice((*int)(ptr), totalGates)
	ptr = unsafe.Add(ptr, sizeGateY)

	ps.MarkerX = unsafe.Slice((*int)(ptr), totalMarkers)
	ptr = unsafe.Add(ptr, sizeMarkerX)

	ps.MarkerY = unsafe.Slice((*int)(ptr), totalMarkers)
	ptr = unsafe.Add(ptr, sizeMarkerY)

	// 4. Slice 4-Byte Fields (Guaranteed 4-byte aligned because previous sizes are multiples of 8)
	ps.ValidGates = unsafe.Slice((*GateType)(ptr), totalGates)
	ptr = unsafe.Add(ptr, sizeGates)

	ps.PuzzleGates = unsafe.Slice((*GateType)(ptr), totalGates)
	ptr = unsafe.Add(ptr, sizeGates)

	ps.AttemptedGates = unsafe.Slice((*GateType)(ptr), totalGates)
	ptr = unsafe.Add(ptr, sizeGates)

	// 5. 1 Byte Fields
	ps.Markers = unsafe.Slice((*PuzzleMarker)(ptr), totalMarkers)
	// ptr = unsafe.Add(ptr, sizeMarkers)

	return ps
}

func (ps *PuzzleSystem) AllocatePuzzle(numberOfGates, numberOfMarkers int) (PuzzleId, error) {
	if int(ps.NextPuzzleId) >= ps.TotalPuzzles {
		return 0, fmt.Errorf("Max Puzzles Reach, Unable to create new puzzle")
	}

	id := ps.NextPuzzleId

	ps.NextPuzzleId++

	var gateOffset int = 0
	var markerOffset int = 0

	if id > 0 {
		prevID := id - 1
		gateOffset = ps.GatesOffsets[prevID] + ps.PuzzleGateCounts[prevID]
		markerOffset = ps.MarkersOffsets[prevID] + ps.MarkersCounts[prevID]
	}

	if gateOffset+numberOfGates > ps.TotalGates {
		return 0, fmt.Errorf("Number of gates(%d) would exceed allotted TotalGates (%d)", numberOfGates, ps.TotalGates)
	}

	if markerOffset+numberOfMarkers > ps.TotalMarkers {
		return 0, fmt.Errorf("Number of Markers(%d) would exceed allotted TotalMarkers (%d)", numberOfMarkers, ps.TotalMarkers)
	}

	ps.PuzzleGateCounts[id] = numberOfGates
	ps.GatesOffsets[id] = gateOffset
	ps.MarkersCounts[id] = numberOfGates
	ps.MarkersOffsets[id] = markerOffset

	return id, nil
}

func (ps *PuzzleSystem) GetGateId(puzzleId PuzzleId, gateIdx int) GateId {
	offset := ps.GatesOffsets[puzzleId]
	return GateId(offset + gateIdx)
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

func (ps *PuzzleSystem) SetAttemptGate(puzzleId PuzzleId, gateIdx int, gate GateType) error {
	offset := ps.GatesOffsets[puzzleId]

	if gateIdx >= ps.PuzzleGateCounts[puzzleId] {
		return fmt.Errorf("GateIdx out of bounds: Given %d ; Count %d", gateIdx, ps.PuzzleGateCounts[puzzleId])
	}

	idx := offset + gateIdx
	ps.AttemptedGates[idx] = gate

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

func (ps *PuzzleSystem) SetMarker(puzzleId PuzzleId, markerIdx, x, y int, marker PuzzleMarker) error {
	offset := ps.MarkersOffsets[puzzleId]

	if markerIdx >= ps.MarkersCounts[puzzleId] {
		return fmt.Errorf("MarkerIdx out of bounds: Given %d ; Count %d", markerIdx, ps.MarkersCounts[puzzleId])
	}

	idx := offset + markerIdx
	ps.Markers[idx] = marker
	ps.MarkerX[idx] = x
	ps.MarkerY[idx] = y

	return nil
}

func (ps *PuzzleSystem) GetMarkers(puzzleId PuzzleId) []PuzzleMarker {
	offset := ps.MarkersOffsets[puzzleId]
	count := ps.MarkersCounts[puzzleId]

	start := offset
	end := start + count

	return ps.Markers[start:end]
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

	gates := ps.GetAttemptGates(puzzleId)
	gatesX := ps.GetGatesX(puzzleId)
	gatesY := ps.GetGatesY(puzzleId)

	for i := 0; i < len(gates); i++ {
		gateType := gates[i]
		x := gatesX[i]
		y := gatesY[i]

		gs.SetCellSprite(gridId, x, y-1, assets.SpriteIDHorizontalBar)
		gs.SetCellSprite(gridId, x, y+1, assets.SpriteIDHorizontalBar)

		// gs.SetCellSprite(gridId, x-1, y-1, assets.SpriteIDHorizontalBar)
		// gs.SetCellSprite(gridId, x+1, y+1, assets.SpriteIDHorizontalBar)
		// gs.SetCellSprite(gridId, x+1, y-1, assets.SpriteIDHorizontalBar)
		// gs.SetCellSprite(gridId, x-1, y+1, assets.SpriteIDHorizontalBar)

		gs.SetCellSprite(gridId, x-1, y, assets.SpriteIDVerticalBar)
		gs.SetCellSprite(gridId, x+1, y, assets.SpriteIDVerticalBar)

		if i < 9 {
			gs.Set(gridId, x-2, y, CellTypeChar, byte('0'+(i+1)))
		}

		hasSymbol := false
		var cornerSprite, gateSprite assets.SpriteID

		switch gateType {
		case GateJoin:
			hasSymbol = true
			cornerSprite = assets.SpriteIDSquare
			gateSprite = assets.SpriteIDGateJoin
		case GateSplit:
			hasSymbol = true
			cornerSprite = assets.SpriteIDDiamond
			gateSprite = assets.SpriteIDGateSplit
		case GatePass:
			hasSymbol = true
			cornerSprite = assets.SpriteIDCircle
			gateSprite = assets.SpriteIDGatePass
		default:
			gs.Set(gridId, x, y, CellTypeChar, '?')
		}

		if hasSymbol {
			gs.SetCellSprite(gridId, x, y, gateSprite)
			gs.SetCellSprite(gridId, x-1, y-1, cornerSprite)
			gs.SetCellSprite(gridId, x-1, y+1, cornerSprite)
			gs.SetCellSprite(gridId, x+1, y-1, cornerSprite)
			gs.SetCellSprite(gridId, x+1, y+1, cornerSprite)
		}
	}

	return nil
}
