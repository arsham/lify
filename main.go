// Package main starts the game.
package main

import (
	"embed"
	"log/slog"

	"github.com/faiface/pixel/pixelgl"
	"github.com/pkg/profile"

	"github.com/arsham/neuragene/game"
	"github.com/arsham/neuragene/internal/config"
)

//go:embed bin
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

	pixelgl.Run(func() {
		g, err := game.NewEngine(env, &assets)
		if err != nil {
			slog.Error("Error running simulation: %w", err)
			return
		}
		g.Run()
	})
}
