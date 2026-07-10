package main

import (
	"bytes"
	_ "embed"
	"log"

	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
)

const (
	screenWidth  = 640 * 2
	screenHeight = 480 * 2
	fontSize     = 28.0 * 2
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
}

func (g *Game) Update() error {
	// Update logic if needed
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Clear the screen with a black background
	// screen.Fill(colornames.Black)
	screen.Fill(color.RGBA{R: 0, G: 0, B: 50, A: 255})

	face := &text.GoTextFace{
		Source: arialFontSource,
		Size:   fontSize, // Use consistent font size
	}

	opt := &text.DrawOptions{}
	// opt.GeoM.Translate(0, 0) // Set Position

	// Set text color to green
	opt.ColorScale.ScaleWithColor(color.RGBA{R: 0, G: 255, B: 0, A: 255})

	// text.Draw(screen, "Hello, Ebitengine!", face, opt)

	/* Grid Test */

	padding := float64(fontSize / 4)

	strNums := ""

	// x = 38 cells ;; 37 with padding
	opt.GeoM.Translate(padding, padding) // Set Position
	for i := range 38 {
		strNums = fmt.Sprintf("%s%d", strNums, i%10)
	}
	text.Draw(screen, strNums, face, opt)

	for i := range 38 {
		i1 := i + 1
		opt2 := &text.DrawOptions{}
		opt2.GeoM.Translate(padding, float64(i1*fontSize)+float64(fontSize)/4.0) // Set Position

		// Set text color to green
		opt2.ColorScale.ScaleWithColor(color.RGBA{R: 0, G: 255, B: 0, A: 255})
		text.Draw(screen, fmt.Sprintf("%d", i1%10), face, opt2)
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	windowSizeW, windowSizeH := screenWidth, screenHeight
	ebiten.SetWindowSize(windowSizeW, windowSizeH)
	ebiten.SetWindowTitle("Ebiten Text Example")
	ebiten.SetWindowResizable(true)

	game := &Game{}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
