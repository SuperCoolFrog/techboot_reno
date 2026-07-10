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
	// Grid ends up being 44,20
	screenWidth   = 640 * 2
	screenHeight  = 480 * 2
	fontSize      = 24.0 * 2.0
	screenPadding = fontSize / 4
)

//go:embed fonts/Courierprime_1OVL.ttf
var arialFontBytes []byte
var arialFontSource *text.GoTextFaceSource

func init() {
	src, err := text.NewGoTextFaceSource(bytes.NewReader(arialFontBytes))
	if err != nil {
		log.Fatalf("Error creating font face: %v", err)
	}
	arialFontSource = src
}

type Game struct {
	Grid Grid
}

func (g *Game) Update() error {
	// Update logic if needed
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
	// ebiten.SetWindowResizable(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := &Game{
		Grid: NewGrid(44, 20, fontSize),
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
