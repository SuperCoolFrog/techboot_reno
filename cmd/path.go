package main

import (
	"fmt"
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
*/
type PathSystem struct {
	TotalPaths int

	StartX []int
	StartY []int
	EndX   []int
	EndY   []int

	NextPathId PathId
}

func NewPathSystem(totalPaths int) *PathSystem {
	ps := &PathSystem{
		TotalPaths: totalPaths,
		StartX:     make([]int, totalPaths),
		StartY:     make([]int, totalPaths),
		EndX:       make([]int, totalPaths),
		EndY:       make([]int, totalPaths),

		NextPathId: 0,
	}

	return ps
}

func (ps *PathSystem) NewPath(startX, startY, endX, endY int) (PathId, error) {
	if int(ps.NextPathId) >= ps.TotalPaths {
		return 0, fmt.Errorf("Unable to allocate any more paths; Max value reached")
	}

	id := ps.NextPathId
	ps.NextPathId++

	ps.StartX[id] = startX
	ps.StartY[id] = startY
	ps.EndX[id] = endX
	ps.EndY[id] = endY

	return id, nil
}
