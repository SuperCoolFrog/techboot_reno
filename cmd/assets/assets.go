package assets

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"os"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	loadMu  sync.Mutex
	loadErr error
	loaded  bool
)

const defaultFilter = ebiten.FilterLinear

// const gameAssetRoot = "./cmd/assets"
const spriteFilePath = "./assets/sprites.png"

var Images = map[SpriteID]*ebiten.Image{}

/*
Loads Images for the requested game.
*/
func Load() error {
	loadMu.Lock()
	defer loadMu.Unlock()

	images, err := loadGameImages()
	if err != nil {
		loadErr = err
		return loadErr
	}

	Images = images
	loaded = true
	loadErr = nil

	return nil
}

func loadGameImages() (map[SpriteID]*ebiten.Image, error) {
	images := make(map[SpriteID]*ebiten.Image)

	spriteSheet, err := decodeImageFile(spriteFilePath)
	if err != nil {
		return nil, fmt.Errorf("load %s: %v", spriteFilePath, err)
	}

	for i := SpriteID(0); i < SpirteIDCount; i++ {
		x1 := int(i) * 48
		images[i] = clippedSubImage(spriteSheet, image.Rect(x1, 0, x1+48, 48))
	}

	return images, nil
}

func decodeImageFile(path string) (*ebiten.Image, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return decodeImageBytes(data)
}

func decodeImageBytes(data []byte) (*ebiten.Image, error) {
	decoded, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return ebiten.NewImageFromImage(decoded), nil
}

func clippedSubImage(source *ebiten.Image, src image.Rectangle) *ebiten.Image {
	if source == nil || src.Empty() {
		return nil
	}
	src = src.Intersect(source.Bounds())
	if src.Empty() {
		return nil
	}
	subImage, ok := source.SubImage(src).(*ebiten.Image)
	if !ok || subImage == nil {
		fmt.Printf("Failed to subimage: %v", src.Bounds())
		return nil
	}

	return subImage
}
