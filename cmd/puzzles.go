package main

import (
	"fmt"
	"unsafe"
)

type PuzzleId uint32
type GateType uint32

const (
	GateEmpty GateType = iota
	GateUnknown
	GateAnd
	GateOr
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

	// --- Global Per-Puzzle Properties ---
	PuzzleGateCounts []int // Tracks exactly how many gates a specific PuzzleId has

	ValidGates     []GateType // Valid Gates that solve the puzzle
	PuzzleGates    []GateType // This is the gates presented to player
	AttemptedGates []GateType // This is the player's last attempt
	GatesOffsets   []int

	GateX []int
	GateY []int

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

func (ps *PuzzleSystem) InitializeGamePuzzles() error {
	/* Placeholder for now until puzzle generator is created */

	id, errorIntro := ps.AllocatePuzzle(1)

	if errorIntro != nil {
		return errorIntro
	}

	ps.IntroPuzzles[0] = id

	ps.SetValidGate(id, 0, 10, 10, GateOr)
	ps.SetPuzzleGate(id, 0, 10, 10, GateUnknown)
	ps.SetAttemptedGate(id, 0, 10, 10, GateUnknown)

	return nil
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
	ps.GateX[idx] = gateY
	ps.GateY[idx] = gateX

	return nil
}

func (ps *PuzzleSystem) SetPuzzleGate(puzzleId PuzzleId, gateIdx, gateX, gateY int, gate GateType) error {
	offset := ps.GatesOffsets[puzzleId]

	if gateIdx >= ps.PuzzleGateCounts[puzzleId] {
		return fmt.Errorf("GateIdx out of bounds: Given %d ; Count %d", gateIdx, ps.PuzzleGateCounts[puzzleId])
	}

	idx := offset + gateIdx
	ps.PuzzleGates[idx] = gate
	ps.GateX[idx] = gateY
	ps.GateY[idx] = gateX

	return nil
}

func (ps *PuzzleSystem) SetAttemptedGate(puzzleId PuzzleId, gateIdx, gateX, gateY int, gate GateType) error {
	offset := ps.GatesOffsets[puzzleId]

	if gateIdx >= ps.PuzzleGateCounts[puzzleId] {
		return fmt.Errorf("GateIdx out of bounds: Given %d ; Count %d", gateIdx, ps.PuzzleGateCounts[puzzleId])
	}

	idx := offset + gateIdx
	ps.AttemptedGates[idx] = gate
	ps.GateX[idx] = gateY
	ps.GateY[idx] = gateX

	return nil
}
