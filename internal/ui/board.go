// Package ui contains logic for interacting and rendering things on the
// screen.
package ui

import (
	"fmt"
	_ "image/png" // Required for loading png files.
	"path/filepath"

	"github.com/oakmound/oak/v4/render"

	"github.com/arsham/lify/internal/config"
)

// Board contains the logic for interacting and rendering everything on the
// board in the game. This is the only object that is allowed to Resolve()
// resources. No other objects is allowed to make that decision, but they can
// cascade the Resolve() if they manage another object themselves.
type Board struct {
	assets map[Asset]*render.Sprite
	fonts  map[Asset]*render.Font
	x, y   int
}

// NewBoard returns a new instance of the board. It loads all the resources
// into memory and returns an error if any of the resources can't be loaded.
func NewBoard(env *config.Env) *Board {
	return &Board{
		x:      env.UI.Width,
		y:      env.UI.Height,
		assets: make(map[Asset]*render.Sprite, 100),
		fonts:  make(map[Asset]*render.Font, 5),
	}
}

// Load loads the assets.
func (b *Board) Load() error {
	herb1, err := render.LoadSprite(filepath.Join("assets", "images", "herb", "herb1.png"))
	if err != nil {
		return fmt.Errorf("failed loading %s: %w", "herb1.png", err)
	}
	b.assets[AssetHerb1] = herb1

	font := render.DefaultFont()
	font, err = font.RegenerateWith(func(g render.FontGenerator) render.FontGenerator {
		g.DPI = 120
		g.Size = 20
		return g
	})
	if err != nil {
		font = render.DefaultFont()
	}
	b.fonts[AssetFontInfo] = font
	return nil
}

// Asset returns an asset from cache. It returns an error if the asset is not
// loaded.
func (b *Board) Asset(a Asset) (*render.Sprite, error) {
	v, ok := b.assets[a]
	if !ok {
		return nil, fmt.Errorf("asset %s not found", a)
	}
	return v, nil
}

// Font returns a font from cache. It returns a default font if the font is not
// loaded.
func (b *Board) Font(a Asset) *render.Font {
	v, ok := b.fonts[a]
	if !ok {
		return render.DefaultFont()
	}
	return v
}
