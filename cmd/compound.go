package main

import (
	"errors"
	"unsafe"
)

/**

i.e. maxId 3, maxChildrent 2

[
	[ [], [] ],[ [], [] ],[ [], [] ],
]
**/

type CompoundSystem struct {
	MaxId       int
	MaxChildren int
	MasterChunk []byte

	Children      []uint16
	ChildrenCount []int
}

func NewCompoundSystem(maxId, maxChildren int) *CompoundSystem {
	cs := &CompoundSystem{
		MaxId:       maxId,
		MaxChildren: maxChildren,
	}

	totalIdCombos := maxId * maxId
	totalChildren := totalIdCombos * maxChildren

	sizeChildrenCount := totalIdCombos * int(unsafe.Sizeof(int(0))) // 8 bytes
	sizeChildren := totalChildren * int(unsafe.Sizeof(uint16(0)))   // 4 bytes

	totalSize := sizeChildrenCount + sizeChildren

	cs.MasterChunk = make([]byte, totalSize)
	ptr := unsafe.Pointer(&cs.MasterChunk[0])

	cs.ChildrenCount = unsafe.Slice((*int)(ptr), totalIdCombos)
	ptr = unsafe.Add(ptr, sizeChildrenCount)

	cs.Children = unsafe.Slice((*uint16)(ptr), totalChildren)
	// ptr = unsafe.Add(ptr, sizeChildren)

	return cs
}

func (cs *CompoundSystem) AddChild(idA, idB uint16, child uint16) error {
	idx := (int(idA) * cs.MaxId) + int(idB)
	count := cs.ChildrenCount[idx]

	if count >= cs.MaxChildren {
		return errors.New("Max children already exist")
	}

	cs.Children[idx+count] = child

	cs.ChildrenCount[idx]++

	return nil
}

func (cs *CompoundSystem) GetChildren(idA, idB uint16) []uint16 {
	idx := (int(idA) * cs.MaxId) + int(idB)
	count := cs.ChildrenCount[idx]

	return cs.Children[idx : idx+count]
}
