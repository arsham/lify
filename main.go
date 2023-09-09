// Package main starts the game.
package main

import (
	"embed"
	"log/slog"
	"os"

	"github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/dlog"

	"github.com/arsham/lify/internal/config"
	"github.com/arsham/lify/internal/game"
)

//go:embed assets
var assets embed.FS

func main() {
	env, err := config.Config()
	if err != nil {
		slog.Error("Failed getting configuration: %w", err)
	}

	oak.SetFS(assets)
	if err := game.Start(env); err != nil {
		dlog.Error("starting lify:", err)
		os.Exit(1)
	}
}
