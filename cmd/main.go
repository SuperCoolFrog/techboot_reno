package main

import (
	"bytes"
	_ "embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
)

const (
	// Grid ends up being 26,19
	screenWidth   = 640 * 2
	screenHeight  = 480 * 2
	fontSize      = 24.0 * 2.0
	screenPadding = fontSize / 4
)

//go:embed fonts/Courierprime_1OVL.ttf
var fontBytes []byte
var fontSrc *text.GoTextFaceSource

func init() {
	src, err := text.NewGoTextFaceSource(bytes.NewReader(fontBytes))
	if err != nil {
		log.Fatalf("Error creating font face: %v", err)
	}
	fontSrc = src
}

type Game struct {
	Grid         *Grid
	inputRunes   []rune
	testBufferId int
}

func (g *Game) Update() error {

	// Update logic if needed
	g.inputRunes = ebiten.AppendInputChars(g.inputRunes[:0])

	if len(g.inputRunes) > 0 {
		// g.Grid.Set(0, 0, RenderFlagKeyCode, byte(g.inputRunes[0]))
		g.Grid.BufferAppend(g.testBufferId, byte(g.inputRunes[0]))
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Clear the screen with a black background
	screen.Fill(color.RGBA{R: 0, G: 0, B: 0, A: 255})
	// screen.Fill(color.RGBA{R: 0, G: 0, B: 50, A: 255})

	err := g.Grid.RenderDebug(screen)
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	windowSizeW, windowSizeH := screenWidth, screenHeight
	ebiten.SetWindowSize(windowSizeW, windowSizeH)
	ebiten.SetWindowTitle("Techboot Reno - Cyber Crawler")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(60) // Locks Update cycles to 60Hz natively

	game := &Game{
		Grid: NewGrid(26, 20, fontSize, screenPadding),
	}
	game.testBufferId = game.Grid.NewBuffer(1, 1, 5, 2)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
