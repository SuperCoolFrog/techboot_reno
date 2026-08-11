package main

type flag uint64

type FlagSystem struct {
	TotalIds   int
	DenseFlags []uint64 //DenseFlags means each index can hold multiple id->flag values
}

func NewFlagSystem(totalIds int) *FlagSystem {
	// 1. Calculate exactly how many uint64 elements we need at compile time.
	// (1000 + 63) / 64 = 16 elements of uint64.
	bitsetStorageSize := (totalIds + 63) / 64

	// 2. Create the fixed-size array type.
	bfs := &FlagSystem{
		TotalIds:   totalIds,
		DenseFlags: make([]uint64, bitsetStorageSize),
	}

	return bfs
}

func (b *FlagSystem) Set(id uint16, value bool) {
	if int(id) >= b.TotalIds {
		return
	}

	idx := id / 64
	bit := id % 64

	if value {
		b.DenseFlags[idx] |= (1 << bit)
	} else {
		b.DenseFlags[idx] &^= (1 << bit) // Using Go's clean bit-clear operator
	}
}

func (b *FlagSystem) Has(id uint16) bool {
	if int(id) >= b.TotalIds {
		return false
	}
	return (b.DenseFlags[id/64] & (1 << (id % 64))) != 0
}
