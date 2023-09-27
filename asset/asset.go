// Package asset contains the assets for rendering.
package asset

import (
	"fmt"
	"image"

	// Image/png is needed for decoding png files.
	_ "image/png"
	"io/fs"
	"path/filepath"

	"github.com/faiface/pixel"
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
	sprites map[Name]*pixel.Sprite
	// batches mirror the sprites.
	batches map[Name]*pixel.Batch
	fs      fs.FS
}

// New creates an AssetManager and loads all the assets into it.
func New(filesystem fs.FS) (*Manager, error) {
	a := &Manager{
		sprites: make(map[Name]*pixel.Sprite, 10),
		batches: make(map[Name]*pixel.Batch, 10),
		fs:      filesystem,
	}

	antPic, err := a.loadPicture(filepath.Join("bin", "images", "ant", "ant.png"))
	if err != nil {
		return nil, fmt.Errorf("loading asset: %w", err)
	}

	antSprite := pixel.NewSprite(antPic, antPic.Bounds())
	a.sprites[Ant] = antSprite
	antBatch := pixel.NewBatch(&pixel.TrianglesData{}, antPic)
	a.batches[Ant] = antBatch

	return a, nil
}

// Batches returns the map of batches.
func (a *Manager) Batches() map[Name]*pixel.Batch {
	return a.batches
}

// Sprites returns the map of sprites.
func (a *Manager) Sprites() map[Name]*pixel.Sprite {
	return a.sprites
}

func (a *Manager) loadPicture(path string) (pixel.Picture, error) {
	file, err := a.fs.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}
