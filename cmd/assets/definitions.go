package assets

type SpriteID uint32

const (
	SpriteIDVerticalBar SpriteID = iota
	SpriteIDHorizontalBar
	SpriteIDRightConnectBar
	SpriteIDLeftConnectBar
	SpriteIDCornerTopLeft
	SpriteIDCornerTopRight
	SpriteIDBottomLeft
	SpriteIDBottomRight
	SpriteIDCircle
	SpriteIDSquare
	SpriteIDCarrotSE
	SpriteIDCarrotNW
	SpriteIDCarrotNE
	SpriteIDCarrotSW
	SpriteIDDiamond
	SpriteIDCarrotRight
	SpriteIDCarrotLeft
	SpriteIDCarrotDown
	SpriteIDCarrotUp

	SpirteIDCount // Always Last
)
