package main

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
}

func NewBuffer(cols, rows, capacity int, lineOverflow bool) *Buffer {
	buffer := &Buffer{
		Cols:         cols,
		Rows:         rows,
		Capacity:     capacity,
		Head:         0,
		YCursor:      0,
		LineOverflow: lineOverflow,
		History:      make([][]byte, capacity),
		XCursors:     make([]int, capacity),
	}

	buffer.NewLine()

	return buffer
}

func (buffer *Buffer) GetXCursor() int {
	roXCursor := 0

	if buffer.YCursor > 0 {
		roXCursor = buffer.XCursors[buffer.YCursor-1]
	}
	return roXCursor
}
func (buffer *Buffer) IncrementXCursor() int {
	if buffer.YCursor > 0 {
		buffer.XCursors[buffer.YCursor-1]++
		return buffer.XCursors[buffer.YCursor-1]
	}

	buffer.XCursors[buffer.YCursor]++
	return buffer.XCursors[buffer.YCursor]
}

func (buffer *Buffer) DecrementXCursor() int {
	if buffer.YCursor > 0 {
		buffer.XCursors[buffer.YCursor-1]--
		return buffer.XCursors[buffer.YCursor-1]
	}

	buffer.XCursors[buffer.YCursor]--
	return buffer.XCursors[buffer.YCursor]
}

func (buffer *Buffer) Append(char byte) {
	if buffer.GetXCursor() >= buffer.Cols {
		if buffer.LineOverflow {
			buffer.NewLine()
		} else {
			return
		}
	}

	buffer.History[buffer.YCursor-1][buffer.GetXCursor()] = char
	buffer.IncrementXCursor()
}

func (buffer *Buffer) AppendAll(chars []byte) {
	for i := 0; i < len(chars); i++ {
		buffer.Append(chars[i])
	}
}

func (buffer *Buffer) AppendDecorators(decor BufferDecorator) {
	preCount := len(decor.Prefix)
	postCount := len(decor.Postfix)
	if preCount+postCount+buffer.GetXCursor() > buffer.Cols {
		if buffer.LineOverflow {
			buffer.NewLine()
		} else {
			return
		}
	}
	for i := 0; i < len(decor.Prefix); i++ {
		buffer.History[buffer.YCursor-1][i] = decor.Prefix[i]
		buffer.IncrementXCursor()
	}

	for i := 0; i < len(decor.Postfix); i++ {
		buffer.History[buffer.YCursor-1][buffer.GetXCursor()+i] = decor.Postfix[i]
	}
}

func (buffer *Buffer) AppendWithDecor(char byte, decor BufferDecorator) {
	preCount := len(decor.Prefix)
	postCount := len(decor.Postfix)
	if preCount+postCount+1+buffer.GetXCursor() >= buffer.Cols {
		if buffer.LineOverflow {
			buffer.NewLine()
		} else {
			return
		}
	}
	for i := 0; i < len(decor.Prefix); i++ {
		buffer.History[buffer.YCursor-1][i] = decor.Prefix[i]
	}
	buffer.Append(char)
	for i := 0; i < len(decor.Postfix); i++ {
		buffer.History[buffer.YCursor-1][buffer.GetXCursor()+i] = decor.Postfix[i]
	}
}

func (buffer *Buffer) DecrementCursor() {
	if buffer.GetXCursor() > 0 {
		for i := 0; i < buffer.Cols-buffer.GetXCursor(); i++ {
			buffer.History[buffer.YCursor-1][buffer.GetXCursor()+i] = ' '
		}
		buffer.DecrementXCursor()
	}
}

func (buffer *Buffer) TrimDecor(decor BufferDecorator) {
	for i := 0; i < len(decor.Prefix); i++ {
		if buffer.History[buffer.YCursor-1][i] == decor.Prefix[i] {
			buffer.History[buffer.YCursor-1][i] = ' '
		}
	}
	for i := 0; i < len(decor.Postfix); i++ {
		buffer.History[buffer.YCursor-1][buffer.GetXCursor()+i] = ' '
	}
}

func (buffer *Buffer) DecrementCursorWithDecor(decor BufferDecorator) {
	if buffer.GetXCursor() > len(decor.Prefix) {
		for i := 0; i < buffer.Cols-buffer.GetXCursor(); i++ {
			buffer.History[buffer.YCursor-1][buffer.GetXCursor()+i] = ' '
		}
		buffer.DecrementXCursor()
		for i := 0; i < len(decor.Postfix); i++ {
			buffer.History[buffer.YCursor-1][buffer.GetXCursor()+i] = decor.Postfix[i]
		}
	}
}

func (buffer *Buffer) NextBuffer() {
	buffer.History[buffer.YCursor] = make([]byte, buffer.Cols)
	buffer.XCursors[buffer.YCursor] = 0
	buffer.YCursor++
}

func (buffer *Buffer) NewLine() {
	buffer.NextBuffer()

	if buffer.Capacity > buffer.Rows {
		for buffer.YCursor-buffer.Head > buffer.Rows {
			buffer.Head++
		}
	}

	//@TODO eventually implement ring buffer
}

func (buffer *Buffer) GetLastBufferLine() ([]byte, bool) {
	if buffer.YCursor > 1 {
		lastY := buffer.YCursor - 2
		b := buffer.History[lastY]
		xCursor := buffer.XCursors[lastY]

		return b[0:xCursor], true
	}

	return []byte{}, false
}

func (buffer *Buffer) DrawToGrid(gridId GridID, x, y int, gs *GridSystem) {
	for r := 0; r < buffer.Rows; r++ {
		//for historyIdx := S3HistoryHead; historyIdx < S3HistoryCursor; historyIdx++ {
		historyIdx := buffer.Head + r
		bytes := buffer.History[historyIdx]

		for i := 0; i < len(bytes); i++ {
			gs.Set(gridId, x+i, y+r, CellTypeChar, bytes[i])
		}
	}

	// if S3HistoryHead > 0 {
	// 	gs.SetCellSprite(GridIdScene3, 1, 1, assets.SpriteIDCarrotUp)
	// }
}
