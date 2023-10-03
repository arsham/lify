// Package asset contains the assets for rendering.
package asset

import (
	"fmt"
	"image"
	_ "image/png" // This is needed for decoding png files.
	"io/fs"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Name is the name of an asset. Each asset has an accompanied batch.
//
//go:generate stringer -type=Name -output=asset_string.go
type Name int

// These are asset names.
const (
	Ant Name = iota + 1
	FruitApple
)

// Manager holds the assets for rendering.
type Manager struct {
	// sprites contains the sprites for rendering.
	sprites map[Name]*ebiten.Image
	fs      fs.FS
}

// New creates an AssetManager and loads all the assets into it.
func New(filesystem fs.FS) (*Manager, error) {
	a := &Manager{
		sprites: make(map[Name]*ebiten.Image, 10),
		fs:      filesystem,
	}

	antPic, _, err := ebitenutil.NewImageFromFileSystem(a.fs, filepath.Join("assets", "images", "ant.png"))
	if err != nil {
		return nil, fmt.Errorf("loading asset: %w", err)
	}
	a.sprites[Ant] = antPic

	fruitSheet, _, err := ebitenutil.NewImageFromFileSystem(a.fs, filepath.Join("assets", "images", "food", "fruits.png"))
	if err != nil {
		return nil, fmt.Errorf("loading asset: %w", err)
	}

	fruits := splitSpriteSheetImages(fruitSheet, 1, 16, 16)

	a.sprites[FruitApple] = fruits[0]

	return a, nil
}

// LoadSpritesheet returns n sub images from the given input image.
func splitSpriteSheetImages(spritesheet *ebiten.Image, n, width, height int) []*ebiten.Image {
	sprites := []*ebiten.Image{}

	for i := 0; i < n; i++ {
		dimensions := image.Rect(i*width, 0, (i+1)*width, height)
		sprite := ebiten.NewImageFromImage(spritesheet.SubImage(dimensions))
		sprites = append(sprites, sprite)
	}

	return sprites
}

// Sprites returns the map of sprites.
func (a *Manager) Sprites() map[Name]*ebiten.Image {
	return a.sprites
}
