package main

import (
	"bytes"
	_ "embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/colornames"
	"image/color"
)

const (
	screenWidth  = 640
	screenHeight = 480
	fontSize     = 28.0 // Consistent font size
)

//go:embed fonts/arial.ttf
var arialFontBytes []byte
var arialFontSource *text.GoTextFaceSource

func init() {
	src, err := text.NewGoTextFaceSource(bytes.NewReader(arialFontBytes))
	if err != nil {
		log.Fatalf("Error creating font face: %v", err)
	}
	arialFontSource = src
}

type Game struct{}

func (g *Game) Update() error {
	// Update logic if needed
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Clear the screen with a black background
	screen.Fill(colornames.Black)

	face := &text.GoTextFace{
		Source: arialFontSource,
		Size:   fontSize, // Use consistent font size
	}

	op := &text.DrawOptions{}
	op.GeoM.Translate(100, 150) // Set Position
	op.ColorScale.ScaleWithColor(color.RGBA{R: 255, G: 99, B: 71, A: 255})

	// Set text color to green
	text.Draw(screen, "Hello, Ebitengine!", face, op)

	op2 := &text.DrawOptions{}
	op2.GeoM.Translate(100, 250) // Set Position
	op2.ColorScale.ScaleWithColor(color.RGBA{R: 255, G: 99, B: 71, A: 255})
	text.Draw(screen, "This is a test.", face, op2)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Ebiten Text Example")

	game := &Game{}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
