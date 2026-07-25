package main

import (
	"fmt"
	"unsafe"
)

type BufferID int

type BufferDecorator struct {
	Prefix  []byte
	Postfix []byte
}

type Buffer struct {
	Cols, Rows, Capacity int
	Head, YCursor        int
	LineOverflow         bool
	History              [][]byte
	XCursors             []int
	internalCapacity     int
}

/*
Developed this way so that we don't have to pass around pointers.  We can just pass around IDs.

IMPORTANT: Cursors are were the Next value goes i.e. current row would be YCursor-1
*/
type BufferSystem struct {
	MaxTotalBytes int
	MasterChunk   []byte // All histories for every buffer

	Histories []byte
	XCursors  []int

	HistoryOffsets  []int
	XCursorsOffsets []int

	TotalRows     []int //Capacity i.e. 2000
	ActiveRows    []int //Displayed/output
	Cols          []int
	ActiveRowHead []int
	YCursor       []int
	LineOverflow  []bool
	NextBufferId  BufferID
}

func NewBufferSystem(maxTotalBytes int, maxBuffersCount int) *BufferSystem {

	/**

	2000 lines, 2000 rows, 50 cols

	100 lines, 100 rows, 10 cols

	maxTotalBytes 100X10 + 2000*50
	100000 + 1000
	101000
	**/

	bs := &BufferSystem{
		MaxTotalBytes:   maxTotalBytes,
		TotalRows:       make([]int, maxBuffersCount),
		ActiveRows:      make([]int, maxBuffersCount),
		Cols:            make([]int, maxBuffersCount),
		ActiveRowHead:   make([]int, maxBuffersCount),
		YCursor:         make([]int, maxBuffersCount),
		LineOverflow:    make([]bool, maxBuffersCount),
		XCursorsOffsets: make([]int, maxBuffersCount),
		HistoryOffsets:  make([]int, maxBuffersCount),
		NextBufferId:    0,
	}

	sizeXCursors := maxTotalBytes * int(unsafe.Sizeof(int(0)))

	bs.MasterChunk = make([]byte, sizeXCursors+maxTotalBytes)
	ptr := unsafe.Pointer(&bs.MasterChunk[0])

	bs.XCursors = unsafe.Slice((*int)(ptr), maxTotalBytes)
	ptr = unsafe.Add(ptr, sizeXCursors)

	bs.Histories = unsafe.Slice((*byte)(ptr), maxTotalBytes)

	return bs
}

func (bs *BufferSystem) AllocateBuffer(cols, activeRows, totalRows int, lineOverflow bool) BufferID {
	id := bs.NextBufferId
	bs.NextBufferId++

	var historyOffset int = 0
	var xCursorsOffset int = 0

	if id > 0 {
		prevId := id - 1
		historyOffset = bs.HistoryOffsets[prevId] + bs.Cols[prevId]*bs.TotalRows[prevId]
		xCursorsOffset = bs.XCursorsOffsets[prevId] + bs.TotalRows[prevId]
	}

	bs.HistoryOffsets[id] = historyOffset
	bs.XCursorsOffsets[id] = xCursorsOffset
	bs.Cols[id] = cols
	bs.ActiveRows[id] = activeRows
	bs.TotalRows[id] = totalRows
	bs.LineOverflow[id] = lineOverflow
	bs.YCursor[id] = 1
	bs.XCursors[xCursorsOffset] = 1
	bs.ActiveRowHead[id] = 0

	return id
}

func (bs *BufferSystem) GetCurrentHistoryIdx(id BufferID) int {
	historyOffset := bs.HistoryOffsets[id]

	yCursor := bs.YCursor[id]

	if yCursor > 0 {
		return historyOffset + (yCursor-1)*bs.Cols[id]
	}

	return historyOffset
}

func (bs *BufferSystem) GetXCursorIdx(id BufferID) int {
	yCursor := bs.YCursor[id]
	xCursorOffset := bs.XCursorsOffsets[id]

	if yCursor > 0 {
		return xCursorOffset + yCursor - 1
	}

	return xCursorOffset
}

func (bs *BufferSystem) GetXCursor(id BufferID) int {
	return bs.XCursors[bs.GetXCursorIdx(id)]
}

func (bs *BufferSystem) IncrementXCursor(id BufferID) int {
	idx := bs.GetXCursorIdx(id)
	bs.XCursors[idx]++
	return bs.XCursors[idx]
}

func (bs *BufferSystem) DecrementXCursor(id BufferID) int {
	idx := bs.GetXCursorIdx(id)
	bs.XCursors[idx]--
	return bs.XCursors[idx]
}

func (bs *BufferSystem) NextBuffer(id BufferID) {
	bs.YCursor[id]++
	bs.IncrementXCursor(id) // Always start cursor at 1
}

func (bs *BufferSystem) NewLine(id BufferID) {
	bs.NextBuffer(id)

	if bs.TotalRows[id] > bs.ActiveRows[id] {
		for bs.YCursor[id]-bs.ActiveRowHead[id] > bs.ActiveRows[id] {
			bs.ActiveRowHead[id]++
		}
	}

	//@TODO eventually implement ring buffer
}

func (bs *BufferSystem) Append(id BufferID, char byte) {
	if bs.GetXCursor(id) >= bs.Cols[id] {
		if bs.LineOverflow[id] {
			bs.NewLine(id)
		} else {
			return
		}
	}

	historyIdx := bs.GetCurrentHistoryIdx(id) + bs.GetXCursor(id)

	bs.Histories[historyIdx] = char
	bs.IncrementXCursor(id)
}

func (bs *BufferSystem) AppendAll(id BufferID, chars []byte) {
	for i := 0; i < len(chars); i++ {
		bs.Append(id, chars[i])
	}
}

func (bs *BufferSystem) AppendDecorators(id BufferID, decor BufferDecorator) {
	preCount := len(decor.Prefix)
	postCount := len(decor.Postfix)

	if preCount+postCount+bs.GetXCursor(id) > bs.Cols[id] {
		if bs.LineOverflow[id] {
			bs.NewLine(id)
		} else {
			return
		}
	}

	historyIdxRow0 := bs.GetCurrentHistoryIdx(id)

	for i := 0; i < len(decor.Prefix); i++ {
		bs.Histories[historyIdxRow0+i] = decor.Prefix[i]

	}

	for bs.GetXCursor(id) <= len(decor.Prefix) {
		bs.IncrementXCursor(id)
	}

	// Postfix is added but doesn't effect XCursor.  Meant to be overwritten
	for i := 0; i < len(decor.Postfix); i++ {
		postFixIdx := historyIdxRow0 + bs.GetXCursor(id) + i
		if postFixIdx < bs.Cols[id] {
			bs.Histories[postFixIdx] = decor.Postfix[i]
		} else {
			break
		}
	}
}

func (bs *BufferSystem) AppendWithDecor(id BufferID, char byte, decor BufferDecorator) {
	preCount := len(decor.Prefix)
	postCount := len(decor.Postfix)

	if preCount+postCount+1+bs.GetXCursor(id) >= bs.Cols[id] {
		if bs.LineOverflow[id] {
			bs.NewLine(id)
		} else {
			return
		}
	}

	historyIdxRow0 := bs.GetCurrentHistoryIdx(id)

	for i := 0; i < len(decor.Prefix); i++ {
		bs.Histories[historyIdxRow0+i] = decor.Prefix[i]
	}

	bs.Append(id, char)

	fmt.Printf("Post Append XCursor %d\n", bs.GetXCursor(id))

	// Postfix is added but doesn't effect XCursor.  Meant to be overwritten
	for i := 0; i < len(decor.Postfix); i++ {
		postFixIdx := historyIdxRow0 + bs.GetXCursor(id) + i
		if postFixIdx < bs.Cols[id] {
			bs.Histories[postFixIdx] = decor.Postfix[i]
		} else {
			break
		}
	}
}

func (bs *BufferSystem) DecrementCursor(id BufferID) {
	if bs.GetXCursor(id) > 0 {
		historyIdx := bs.GetCurrentHistoryIdx(id)
		for i := 0; i < bs.Cols[id]-bs.GetXCursor(id); i++ {
			bs.Histories[historyIdx+bs.GetXCursor(id)+i] = ' '
		}
		bs.DecrementXCursor(id)
	}
}

func (bs *BufferSystem) TrimDecor(id BufferID, decor BufferDecorator) {
	historyIdx := bs.GetCurrentHistoryIdx(id)
	for i := 0; i < len(decor.Prefix); i++ {
		if bs.Histories[historyIdx+i] == decor.Prefix[i] {
			bs.Histories[historyIdx+i] = ' '
		}
	}

	for i := 0; i < len(decor.Postfix); i++ {
		bs.Histories[historyIdx+bs.GetXCursor(id)+i] = ' '
	}
}

func (bs *BufferSystem) DecrementCursorWithDecor(id BufferID, decor BufferDecorator) {
	historyIdx := bs.GetCurrentHistoryIdx(id)

	if bs.GetXCursor(id) > len(decor.Prefix) {
		for i := 0; i < bs.Cols[id]-bs.GetXCursor(id); i++ {
			bs.Histories[historyIdx+bs.GetXCursor(id)+i] = ' '
		}

		bs.DecrementXCursor(id)

		for i := 0; i < len(decor.Postfix); i++ {
			bs.Histories[historyIdx+bs.GetXCursor(id)+i] = decor.Postfix[i]
		}
	}
}

func (bs *BufferSystem) GetBufferRow(id BufferID, rowIdx int) ([]byte, bool) {
	if rowIdx < bs.TotalRows[id] {
		historyOffset := bs.HistoryOffsets[id]
		xCursorOffset := bs.XCursorsOffsets[id]

		historyIdx := historyOffset + rowIdx*bs.Cols[id]
		xCursor := bs.XCursors[xCursorOffset+rowIdx]

		fmt.Printf("GetBRVals rowIdx: %d ;; HO: %d ;; XCO: %d ;; HIdx: %d ;; XCur: %d \n", rowIdx, historyOffset, xCursorOffset, historyIdx, xCursor)

		return bs.Histories[historyIdx:xCursor], true
	}

	return []byte{}, false
}

func (bs *BufferSystem) GetLastBufferLine(id BufferID) ([]byte, bool) {
	if bs.YCursor[id] > 1 {
		lastY := bs.YCursor[id] - 2

		return bs.GetBufferRow(id, lastY)
	}

	return []byte{}, false
}

func (bs *BufferSystem) DrawToGrid(bufferId BufferID, gridId GridID, x, y int, gs *GridSystem) {
	historyOffset := bs.HistoryOffsets[bufferId]
	rowStart := historyOffset + bs.ActiveRowHead[bufferId]

	fmt.Printf("Drw:: rowStart %d ;; histOff %d\n", rowStart, historyOffset)

	for r := 0; r < bs.YCursor[bufferId]; r++ {
		rowIdx := rowStart + r

		bufferBytes, valid := bs.GetBufferRow(bufferId, rowIdx)

		if !valid {
			continue
		}

		for i := 0; i < len(bufferBytes); i++ {
			gs.Set(gridId, x+i, y+r, CellTypeChar, bufferBytes[i])
		}
	}
}
