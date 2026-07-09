package main

import (
	"log"
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/colornames"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

//go:embed fonts/arial.ttf
var arialFont []byte

type Game struct{}

func (g *Game) Update() error {
	// Update logic if needed
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Clear the screen with a black background
	screen.Fill(colornames.Black)

	// Set text color to green
	text.Draw(screen, "Hello, Ebitengine!", ebiten.ArialFont, 50, 50, colornames.Green)
	text.Draw(screen, "This is a test.", ebiten.ArialFont, 50, 100, colornames.Green)
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
