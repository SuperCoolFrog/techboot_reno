package main

// import (
// 	"unsafe"
// )

type PathId uint32

type RouteType uint32

const (
	RouteTypeVOff RouteType = iota
	RouteTypeVOn
	RouteTypeHOff
	RouteTypeHOn
)

type PathSystem struct {
	TotalPaths  int
	MasterChunk []byte

	RoutingSegment     []int
	RoutingSegsOffsets []int
	RoutingSegX        []int
	RoutingSegY        []int
	RoutingSegType     []RouteType

	Paths []int

	NextPathId PathId
}

func NewPathSystem(totalPaths, maxSegmentsPerPath int) *PathSystem {
	ps := &PathSystem{}

	// total := totalPaths * maxSegmentsPerPath

	// sizeRoutings := total * int(unsafe.Sizeof(int(0)))
	// sizeRoutingX := total * int(unsafe.Sizeof(int(0)))
	// sizeRoutingY := total * int(unsafe.Sizeof(int(0)))

	return ps
}
