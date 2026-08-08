package main

import (
	"fmt"
	"unsafe"
)

type PathId uint32

type RouteType uint32

const (
	RouteTypeVOff RouteType = iota
	RouteTypeVOn
	RouteTypeHOff
	RouteTypeHOn
)

/*
Only do linear paths for now and then combine linear paths for turns

gcd(1280, 960) = 320

320+1 = 321 = all lattice points between (0,0) and (1280, 960)

This means that is the maximum number of points on a path that could be considered per puzzle.
This assuming you are working on a grid 1px CellSize
Since this will most likely never be the case, I will assume 50points per puzzle and I should be safe within output grid

Keep in mind though, if I ever do something with full screen pathing, I may need to increase this.

*/

// Path System HardCoded 50points per Path
type PathSystem struct {
	PointsPerPath int
	TotalPoints   int
	TotalPaths    int
	MasterChunk   []byte

	StartX        []int
	StartY        []int
	EndX          []int
	EndY          []int
	PointsCount   []int
	PointsOffsets []int

	PointsX []int
	PointsY []int

	NextPathId PathId
}

func NewPathSystem(totalPaths int) *PathSystem {
	PointsPerPath := 50
	TotalPoints := (PointsPerPath * totalPaths)

	ps := &PathSystem{
		TotalPaths:    totalPaths,
		PointsPerPath: PointsPerPath,
		TotalPoints:   TotalPoints,
		// StartX:     make([]int, totalPaths),
		// StartY:     make([]int, totalPaths),
		// EndX:       make([]int, totalPaths),
		// EndY:       make([]int, totalPaths),

		NextPathId: 0,
	}

	sizeStartX := totalPaths * int(unsafe.Sizeof(int(0)))
	sizeStartY := totalPaths * int(unsafe.Sizeof(int(0)))
	sizeEndX := totalPaths * int(unsafe.Sizeof(int(0)))
	sizeEndY := totalPaths * int(unsafe.Sizeof(int(0)))
	sizePointsCount := totalPaths * int(unsafe.Sizeof(int(0)))
	sizePointsOffsets := totalPaths * int(unsafe.Sizeof(int(0)))
	sizePointsX := TotalPoints * int(unsafe.Sizeof(int(0)))
	sizePointsY := TotalPoints * int(unsafe.Sizeof(int(0)))

	totalSize := sizeStartX + sizeStartY + sizeEndX + sizeEndY +
		sizePointsCount + sizePointsOffsets + sizePointsX + sizePointsY

	ps.MasterChunk = make([]byte, totalSize)
	ptr := unsafe.Pointer(&ps.MasterChunk[0])

	ps.StartX = unsafe.Slice((*int)(ptr), totalPaths)
	ptr = unsafe.Add(ptr, sizeStartX)

	ps.StartY = unsafe.Slice((*int)(ptr), totalPaths)
	ptr = unsafe.Add(ptr, sizeStartY)

	ps.EndX = unsafe.Slice((*int)(ptr), totalPaths)
	ptr = unsafe.Add(ptr, sizeEndX)

	ps.EndY = unsafe.Slice((*int)(ptr), totalPaths)
	ptr = unsafe.Add(ptr, sizeEndY)

	ps.PointsCount = unsafe.Slice((*int)(ptr), totalPaths)
	ptr = unsafe.Add(ptr, sizePointsCount)

	ps.PointsOffsets = unsafe.Slice((*int)(ptr), totalPaths)
	ptr = unsafe.Add(ptr, sizePointsOffsets)

	ps.PointsX = unsafe.Slice((*int)(ptr), TotalPoints)
	ptr = unsafe.Add(ptr, sizePointsX)

	ps.PointsY = unsafe.Slice((*int)(ptr), TotalPoints)
	// ptr = unsafe.Add(ptr, sizePointsY)

	return ps
}

func (ps *PathSystem) NewPath(startX, startY, endX, endY int) (PathId, error) {
	if int(ps.NextPathId) >= ps.TotalPaths {
		return 0, fmt.Errorf("Unable to allocate any more paths; Max value reached")
	}

	id := ps.NextPathId
	ps.NextPathId++

	offset := 0
	if id > 0 {
		offset = ps.PointsOffsets[id-1] + ps.PointsPerPath
	}
	ps.PointsOffsets[id] = offset

	ps.StartX[id] = startX
	ps.StartY[id] = startY
	ps.EndX[id] = endX
	ps.EndY[id] = endY

	itrX := 1
	if startX > endX {
		itrX = -1
	}

	itrY := 1
	if startY > endY {
		itrY = -1
	}

	x := startX
	y := startY
	count := 0

	for (x != endX || y != endY) && count < ps.PointsPerPath {
		ps.PointsX[offset+count] = x
		ps.PointsY[offset+count] = y

		count++

		if x != endX {
			x += itrX
		}
		if y != endY {
			y += itrY
		}
	}

	// Add last endpoint
	ps.PointsX[offset+count] = endX
	ps.PointsY[offset+count] = endY

	count++

	ps.PointsCount[id] = count

	return id, nil
}

func (ps *PathSystem) GetXYPoints(pathId PathId) (XPoints []int, YPoints []int) {
	offset := ps.PointsOffsets[pathId]
	count := ps.PointsCount[pathId]

	return ps.PointsX[offset : offset+count], ps.PointsY[offset : offset+count]
}
