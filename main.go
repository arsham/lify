// Package main starts the game.
package main

import (
	"embed"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pkg/profile"

	"github.com/arsham/neuragene/internal/config"
	"github.com/arsham/neuragene/internal/game"
)

//go:embed assets
var assets embed.FS

func main() {
	env, err := config.Config()
	if err != nil {
		slog.Error("Failed getting configuration: %w", err)
	}
	defer profile.Start(
		profile.CPUProfile,
		profile.ProfilePath("./tmp/profiles"),
		profile.NoShutdownHook,
	).Stop()

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	g, err := game.NewEngine(env, &assets)
	if err != nil {
		slog.Error("Error running simulation: %w", err)
		return
	}
	if err := ebiten.RunGame(g); err != nil {
		slog.Error(err.Error())
	}
}
