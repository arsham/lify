// Package asset contains the assets for rendering.
package asset

import (
	"fmt"
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

	return a, nil
}

// Sprites returns the map of sprites.
func (a *Manager) Sprites() map[Name]*ebiten.Image {
	return a.sprites
}
