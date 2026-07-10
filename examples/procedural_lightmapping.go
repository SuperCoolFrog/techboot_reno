package main

import (
	"image"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 640
	screenHeight = 480
	// **baseGlowPadding**: The default pixel radius used when baking the asset at startup
	baseGlowPadding = 50
)

type Game struct {
	glowTexture *ebiten.Image
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{12, 12, 18, 255})

	// =========================================================
	// **DYNAMIC CONTROLS**: Tweak these values to change the bloom
	// =========================================================
	targetRadius := 45.0  // **Radius**: Adjusts size in pixels (Default asset was 50.0)
	targetIntensity := .5 // **Intensity**: 1.0 is normal, 2.0 is double brightness, 0.5 is faint

	// Square configurations
	sqX, sqY := 80.0, 140.0
	sqW, sqH := 200.0, 200.0

	// 1. **Calculate Matrix Scale**: Find the ratio between current target radius and base baked radius
	scaleFactor := targetRadius / baseGlowPadding

	glowOpts := &ebiten.DrawImageOptions{}
	glowOpts.Blend = ebiten.BlendLighter

	// 2. **Apply Sizing Matrix**: Scale the glow card texture dynamically on the GPU
	glowOpts.GeoM.Scale(scaleFactor, scaleFactor)

	// 3. **Center Alignment Offset**: Offset calculation to keep the scaled glow perfectly centered over the square
	// Center calculation formula: SquarePosition - (GlowPadding * Scale)
	offsetX := sqX - (baseGlowPadding * scaleFactor)
	offsetY := sqY - (baseGlowPadding * scaleFactor)

	// If the square dimensions change, account for the scaling layout differences
	offsetX += (sqW - (sqW * scaleFactor)) / 2.0
	offsetY += (sqH - (sqH * scaleFactor)) / 2.0

	glowOpts.GeoM.Translate(offsetX, offsetY)

	// 4. **Apply Intensity Matrix**: Scale the color channels to blow out the brightness value explicitly
	glowOpts.ColorScale.Scale(
		float32(targetIntensity), // Red
		float32(targetIntensity), // Green
		float32(targetIntensity), // Blue
		float32(targetIntensity), // Alpha channel transparency multiplier
	)

	// 5. **Render pass**: Draw the customized light layer, then drop the crisp asset on top
	screen.DrawImage(g.glowTexture, glowOpts)
	vectorDrawSquare(screen, int(sqX), int(sqY), int(sqW), int(sqH), color.RGBA{0, 255, 255, 255})

	// Flat reference block on the right
	vectorDrawSquare(screen, 360, 140, 200, 200, color.RGBA{0, 255, 255, 255})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func generateGlowCard(baseW, baseH, glowRadius int, r, g, b uint8) *ebiten.Image {
	cardSizeW := baseW + (glowRadius * 2)
	cardSizeH := baseH + (glowRadius * 2)
	img := image.NewRGBA(image.Rect(0, 0, cardSizeW, cardSizeH))

	for x := 0; x < cardSizeW; x++ {
		for y := 0; y < cardSizeH; y++ {
			var dx, dy float64
			if x < glowRadius {
				dx = float64(glowRadius - x)
			} else if x >= glowRadius+baseW {
				dx = float64(x - (glowRadius + baseW - 1))
			} else {
				dx = 0
			}

			if y < glowRadius {
				dy = float64(glowRadius - y)
			} else if y >= glowRadius+baseH {
				dy = float64(y - (glowRadius + baseH - 1))
			} else {
				dy = 0
			}

			dist := math.Sqrt(dx*dx + dy*dy)
			if dist <= float64(glowRadius) {
				factor := math.Pow(1.0-(dist/float64(glowRadius)), 2.5)
				alpha := uint8(factor * 255.0)
				img.Set(x, y, color.RGBA{
					R: uint8(float64(r) * factor),
					G: uint8(float64(g) * factor),
					B: uint8(float64(b) * factor),
					A: alpha,
				})
			}
		}
	}
	return ebiten.NewImageFromImage(img)
}

func main() {
	// Pre-generate the base 1:1 light reference map (200x200 canvas core + 50px asset padding)
	glowTex := generateGlowCard(200, 200, baseGlowPadding, 0, 255, 255)

	g := &Game{
		glowTexture: glowTex,
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Dynamic Sizing & Intensity controls")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func vectorDrawSquare(img *ebiten.Image, x, y, w, h int, clr color.Color) {
	rect := image.Rect(x, y, x+w, y+h)
	sub := img.SubImage(rect).(*ebiten.Image)
	sub.Fill(clr)
}
