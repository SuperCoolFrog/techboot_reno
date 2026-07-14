package main

type Buffer struct {
	Cols, Rows, Capacity   int
	Head, YCursor, XCursor int
	LineOverflow           bool
	History                [][]byte
}

func NewBuffer(cols, rows, capacity int, lineOverflow bool) *Buffer {
	buffer := &Buffer{
		Cols:         cols,
		Rows:         rows,
		Capacity:     capacity,
		Head:         0,
		YCursor:      0,
		XCursor:      0,
		LineOverflow: lineOverflow,
		History:      make([][]byte, capacity),
	}

	buffer.NewLine()

	return buffer
}

func (buffer *Buffer) AppendToBuffer(char byte) {
	if buffer.XCursor > buffer.Cols {
		if buffer.LineOverflow {
			buffer.NewLine()
		} else {
			return
		}
	}

	buffer.History[buffer.YCursor-1][buffer.XCursor] = char
	buffer.XCursor++
}

func (buffer *Buffer) AppendPrePostfix(prefix, postfix []byte) {
	preCount := len(prefix)
	postCount := len(postfix)
	if preCount+postCount+buffer.XCursor > buffer.Cols {
		if buffer.LineOverflow {
			buffer.NewLine()
		} else {
			return
		}
	}
	for i := 0; i < len(prefix); i++ {
		buffer.History[buffer.YCursor-1][i] = prefix[i]
		buffer.XCursor++
	}
	for i := 0; i < len(postfix); i++ {
		buffer.History[buffer.YCursor-1][buffer.XCursor+i] = postfix[i]
	}
}

func (buffer *Buffer) AppendWithPrePostfix(char byte, prefix, postfix []byte) {
	preCount := len(prefix)
	postCount := len(postfix)
	if preCount+postCount+1+buffer.XCursor > buffer.Cols {
		if buffer.LineOverflow {
			buffer.NewLine()
		} else {
			return
		}
	}
	for i := 0; i < len(prefix); i++ {
		buffer.History[buffer.YCursor-1][i] = prefix[i]
	}
	buffer.AppendToBuffer(char)
	for i := 0; i < len(postfix); i++ {
		buffer.History[buffer.YCursor-1][buffer.XCursor+i] = postfix[i]
	}
}

func (buffer *Buffer) TrimEnd() {
	if buffer.XCursor > 0 {
		for i := 0; i < buffer.Cols-buffer.XCursor; i++ {
			buffer.History[buffer.YCursor-1][buffer.XCursor+i] = ' '
		}
		buffer.XCursor--
	}
}

func (buffer *Buffer) TrimEndWithPrePostfix(prefix, postfix []byte) {
	if buffer.XCursor > len(prefix) {
		for i := 0; i < buffer.Cols-buffer.XCursor; i++ {
			buffer.History[buffer.YCursor-1][buffer.XCursor+i] = ' '
		}
		buffer.XCursor--
		for i := 0; i < len(postfix); i++ {
			buffer.History[buffer.YCursor-1][buffer.XCursor+i] = postfix[i]
		}
	}
}

func (buffer *Buffer) NextBuffer() {
	buffer.History[buffer.YCursor] = make([]byte, buffer.Cols)
	buffer.XCursor = 0
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
