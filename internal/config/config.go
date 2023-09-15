// Package config contains the necessary logic to setup the application, like
// logs and environment variables.
package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/kkyr/fig"
)

var (
	version    = "development"
	currentSha = "N/A"
)

// Env contains the environment variables required to run the application.
type Env struct {
	LogLevel int `fig:"log_level" default:"8"`
	UI       struct {
		DPI    int `default:"96"`
		Width  int `default:"1920"`
		Height int `default:"1080"`
	}
}

// Config processes the environment variables and returns the Env object.
func Config() (*Env, error) {
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "version" {
		slog.Info("Starting lify", "version", version, "current_sha", currentSha)
		os.Exit(0)
	}
	var e Env
	err := fig.Load(&e,
		fig.File("config.yaml"),
		fig.UseEnv(""))
	if err != nil {
		return nil, fmt.Errorf("failed to process environment variables: %w", err)
	}

	opts := &slog.HandlerOptions{
		Level: slog.Level(e.LogLevel),
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	SetLogger(slog.New(handler))
	return &e, nil
}
