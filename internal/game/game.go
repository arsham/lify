// Package game contains the necessary logic for interacting with the game
// through UX.
package game

import (
	"fmt"

	"github.com/arsham/lify/internal/config"
	"github.com/arsham/lify/internal/ui"
)

// Start initialises the game, starts it and draws the UI.
func Start(env *config.Env) error {
	b := ui.NewBoard(env)
	s, err := ui.NewScene(env, b)
	if err != nil {
		return fmt.Errorf("creating a new scene: %w", err)
	}

	return s.Start()
}
