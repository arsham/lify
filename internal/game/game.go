// Package game contains the necessary logic for interacting with the game
// through UX.
package game

import (
	"github.com/arsham/lify/internal/config"
	"github.com/hajimehoshi/ebiten/v2"
)

// Game controls how the game plays out.
type Game struct {
	env *config.Env
}

// Start initialises the game, starts it and draws the UI.
func Start(env *config.Env) error {
	g := &Game{
		env: env,
	}
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetWindowSize(1024, 800)
	return ebiten.RunGame(g)
}

// Update updates a game by one tick. The given argument represents a screen
// image.
//
// You can assume that Update is always called TPS-times per second (60 by
// default), and you can assume that the time delta between two Updates is
// always 1 / TPS [s] (1/60[s] by default).
//
// An actual TPS is available by ActualTPS(), and the result might slightly
// differ from your expected TPS, but still, your game logic should stick to
// the fixed time delta and should not rely on ActualTPS() value. This API is
// for just measurement and/or debugging. In the long run, the number of Update
// calls should be adjusted based on the set TPS on average.
//
// In the first frame, it is ensured that Update is called at least once before
// Draw. You can use Update to initialise the game state.
//
// If the error returned is nil, game execution proceeds normally. If the error
// returned is Termination, game execution halts, but does not return an error
// from RunGame. If the error returned is any other non-nil value, game
// execution halts and the error is returned from RunGame.
func (r *Game) Update() error {
	return nil
}

// Draw draws the game screen by one frame.
//
// The give argument represents a screen image. The updated content is adopted
// as the game screen.
//
// The frequency of Draw calls depends on the user's environment, especially
// the monitors refresh rate. For portability, you should not put your game
// logic in Draw in general.
func (r *Game) Draw(*ebiten.Image) {
}

// Layout accepts a native outside size in device-independent pixels and
// returns the game's logical screen size.
//
// Even though the outside size and the screen size differ, the rendering scale
// is automatically adjusted to fit with the outside.
//
// Layout is called almost every frame.
//
// You can return a fixed screen size if you don't care, or you can also return
// a calculated screen size adjusted with the given outside size.
func (r *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
