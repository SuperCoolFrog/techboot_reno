package assets

type SpriteID uint32

const (
	SpriteIDVerticalBar SpriteID = iota
	SpriteIDHorizontalBar
	SpriteIDRightConnectBar
	SpriteIDLeftConnectBar
	SpriteIDCornerTopLeft
	SpriteIDCornerTopRight
	SpriteIDCornerBottomLeft
	SpriteIDCornerBottomRight
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
