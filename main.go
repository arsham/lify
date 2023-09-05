// Package main starts the game.
package main

import (
	"log/slog"
	"os"

	"github.com/arsham/lify/internal/config"
	"github.com/arsham/lify/internal/game"
)

func main() {
	env, err := config.Config()
	if err != nil {
		slog.Error("Failed getting configuration: %w", err)
	}
	if err := game.Start(env); err != nil {
		config.Logger().Error("starting lify: %w", err)
		os.Exit(1)
	}
}
