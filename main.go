// Package main starts the game.
package main

import (
	"fmt"
	"log/slog"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/pkg/profile"

	"github.com/arsham/neuragene/entity"
	"github.com/arsham/neuragene/game"
	"github.com/arsham/neuragene/internal/config"
	"github.com/arsham/neuragene/scene"
	"github.com/arsham/neuragene/system"
)

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
		err := run(env)
		if err != nil {
			slog.Error("Error running simulation: %w", err)
		}
	})
}

func run(env *config.Env) error {
	cfg := pixelgl.WindowConfig{
		Title:     "Neuragene",
		Bounds:    pixel.R(0, 0, float64(env.UI.Width), float64(env.UI.Height)),
		VSync:     true,
		Resizable: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		return fmt.Errorf("creating new window: %w", err)
	}

	em := entity.NewManager(10)
	sm := system.NewManager(4)
	g := &game.Engine{
		Window:       win,
		Entities:     em,
		Systems:      sm,
		CurrentScene: scene.PlayScene,
	}
	g.Scenes = map[scene.Type]game.Scene{
		scene.PlayScene: scene.NewPlay(g),
	}
	err = g.Setup()
	if err != nil {
		return fmt.Errorf("setting up the engine: %w", err)
	}
	g.Run()
	return nil
}
