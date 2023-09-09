// Package game contains the necessary logic for interacting with the game
// through UX.
package game

import (
	"context"
	"fmt"
	"image/color"
	"time"

	"github.com/disintegration/gift"
	"github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/dlog"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/key"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/render/mod"
	"github.com/oakmound/oak/v4/scene"

	"github.com/arsham/lify/internal/config"
	"github.com/arsham/lify/internal/ui"
)

// Game controls how the game plays out.
type Game struct {
	env *config.Env
	ui  *ui.Board
}

// Start initialises the game, starts it and draws the UI.
func Start(env *config.Env) error {
	b := ui.NewBoard(env)
	g := &Game{
		env: env,
		ui:  b,
	}

	render.SetDrawStack(
		render.NewCompositeR(),
		render.NewDynamicHeap(),
		render.NewStaticHeap(),
	)

	win := oak.NewWindow()
	err := win.AddScene(sceneLify, g.startLifyScene(win))
	if err != nil {
		return fmt.Errorf("adding scene: %w", err)
	}

	err = win.AddScene(sceneLoading, g.loadingScene(win, b))
	if err != nil {
		return fmt.Errorf("adding scene: %w", err)
	}

	return win.Init(sceneLoading, func(c oak.Config) (oak.Config, error) {
		c.FrameRate = 60
		c.DrawFrameRate = 60
		c.Screen = oak.Screen{
			Width:  env.UI.Width,
			Height: env.UI.Height,
			Scale:  1,
		}
		c.Debug = oak.Debug{
			Level: "Info",
		}
		c.Title = "Lify Simulator"
		c.TrackInputChanges = true
		c.LoadBuiltinCommands = true
		c.TopMost = true
		c.BatchLoad = false

		return c, nil
	})
}

func (g *Game) startLifyScene(win *oak.Window) scene.Scene {
	return scene.Scene{
		Start: func(ctx *scene.Context) {
			event.GlobalBind(ctx, key.Down(key.Q), func(key.Event) event.Response {
				ctx.Window.Quit()
				return 0
			})
			win.ParentContext = context.WithValue(context.Background(), preLoadTimeStr, time.Now())
			screen := render.NewColorBoxM(win.Bounds().X(), win.Bounds().Y(), color.RGBA{0, 0, 0, 0})
			mid := win.Bounds().DivConst(2)

			herb, err := g.ui.Asset(ui.AssetHerb1)
			if err != nil {
				dlog.Error("Failed getting asset", err)
				ctx.Window.Quit()
				return
			}
			identM := herb.Modify(mod.ResizeToFit(64, 64, gift.CubicResampling))
			identM.Draw(screen, float64(mid.X()), float64(mid.Y()))
			_, err = render.Draw(screen)
			if err != nil {
				dlog.Error("Failed rendering text:", err)
				ctx.Window.Quit()
				return
			}
			win.SetLoadingRenderable(screen)
		},
		End: func() (string, *scene.Result) {
			return sceneLify, nil
		},
	}
}

func (g *Game) loadingScene(win *oak.Window, b *ui.Board) scene.Scene {
	return scene.Scene{
		Start: func(ctx *scene.Context) {
			err := win.SetFullScreen(true)
			if err != nil {
				dlog.Error("Failed setting full screen failed:", err)
			}
			titleText := render.NewText("Loading assets...", 0, 0)
			titleText.SetFont(g.ui.Font(ui.AssetFontInfo))
			putCentre(ctx, titleText, axixXY)
			_, err = render.Draw(titleText)
			if err != nil {
				dlog.Error("Failed rendering text:", err)
				ctx.Window.Quit()
				return
			}

			event.GlobalBind(ctx, key.Down(key.Q), func(key.Event) event.Response {
				ctx.Window.Quit()
				return 0
			})

			go func() {
				err := b.Load()
				if err != nil {
					dlog.Error("Failed loading assets:", err)
					ctx.Window.Quit()
					return
				}

				titleText.SetString("Assets have been loaded")
				titleText.SetFont(g.ui.Font(ui.AssetFontInfo))
				putCentre(ctx, titleText, axixXY)
				bounds := ctx.Window.Bounds()
				instructions := render.NewText("Press Enter to start, or press Q to quit", 0, float64(bounds.Y()*3/4))
				instructions.SetFont(g.ui.Font(ui.AssetFontInfo))
				putCentre(ctx, instructions, axixX)

				_, err = render.Draw(instructions)
				if err != nil {
					dlog.Error("Failed rendering text:", err)
					ctx.Window.Quit()
					return
				}

				event.GlobalBind(ctx, key.AnyDown, func(key.Event) event.Response {
					ctx.Window.NextScene()
					return 0
				})
			}()
		},
		End: func() (string, *scene.Result) {
			return sceneLify, nil
		},
	}
}
