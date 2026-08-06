package main

import (
	"errors"
	"unsafe"
)

type JunctionSystem struct {
	TotalParents         int
	TotalChildren        int
	MaxChildrenPerParent int
	MasterChunk          []byte

	Parents            []uint32
	ParentIsSet        []bool
	ChildrenCount      []int
	ChildrenCursor     []int
	Children           []uint32
	NextChildrenOffset int
	ChildrenOffsets    []int
}

func NewJunctionSystem(numberOfParents, maxChildrenPerParent int) *JunctionSystem {
	jx := &JunctionSystem{
		TotalParents:         numberOfParents,
		TotalChildren:        numberOfParents * maxChildrenPerParent,
		MaxChildrenPerParent: maxChildrenPerParent,
		ParentIsSet:          make([]bool, numberOfParents),
		ChildrenCount:        make([]int, numberOfParents),
		ChildrenCursor:       make([]int, numberOfParents),
		ChildrenOffsets:      make([]int, numberOfParents),
		NextChildrenOffset:   0,
	}

	sizeParents := numberOfParents * int(unsafe.Sizeof(uint32(0)))
	sizeChildren := jx.TotalChildren * int(unsafe.Sizeof(uint32(0)))

	totalSize := sizeParents + sizeChildren

	jx.MasterChunk = make([]byte, totalSize)
	ptr := unsafe.Pointer(&jx.MasterChunk[0])

	jx.Parents = unsafe.Slice((*uint32)(ptr), numberOfParents)
	ptr = unsafe.Add(ptr, numberOfParents)

	jx.Children = unsafe.Slice((*uint32)(ptr), jx.TotalChildren)

	return jx
}

func (jx *JunctionSystem) AddParent(parentId uint32, childrenCount int) error {
	if jx.ParentIsSet[parentId] {
		return errors.New("Parent has previously been set.")
	}

	if childrenCount > jx.MaxChildrenPerParent {
		return errors.New("ChildrenCount exceeds max allowed per parent")
	}

	if jx.NextChildrenOffset+childrenCount >= jx.TotalChildren {
		return errors.New("Not enough room to allocate children")
	}

	jx.ParentIsSet[parentId] = true
	jx.ChildrenCount[parentId] = childrenCount
	jx.ChildrenCursor[parentId] = 0
	jx.ChildrenOffsets[parentId] = jx.NextChildrenOffset

	jx.NextChildrenOffset = jx.NextChildrenOffset + childrenCount

	return nil
}

func (jx *JunctionSystem) AddChild(parentId uint32, childId uint32) error {
	if !jx.ParentIsSet[parentId] {
		return errors.New("Parent has not yet been set using AddParent/2.")
	}

	if jx.ChildrenCursor[parentId] >= jx.ChildrenCount[parentId] {
		return errors.New("No room left to assign a child")
	}

	cursor := jx.ChildrenCursor[parentId]
	offset := jx.ChildrenOffsets[parentId]

	jx.Children[offset+cursor] = childId

	jx.ChildrenCursor[parentId]++

	return nil
}

func (jx *JunctionSystem) GetChildren(parentId uint32) ([]uint32, error) {
	if !jx.ParentIsSet[parentId] {
		return nil, errors.New("Matching parentId has not been added to junction system")
	}

	cursor := jx.ChildrenCursor[parentId]
	offset := jx.ChildrenOffsets[parentId]

	children := jx.Children[offset : offset+cursor]

	return children, nil
}
