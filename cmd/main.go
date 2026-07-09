package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"log"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

type Game struct{}

func (g *Game) Update() error {
	//Update logic if needed
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Clear the screen with a black background
	screen.Fill(color.RGBA(0, 0, 0, 255))

	// Set text color to green
	text.Draw(screen, "Hello, Ebitengine!", basicfont.Face7x13Dot, 50, 50, color.RGBA{0, 255, 0, 255})
	text.Draw(screen, "This is a test.", basicfont.Face7x13Dot, 50, 100, color.RGBA{0, 255, 0, 255})
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
