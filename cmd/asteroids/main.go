package main

import (
	"fmt"
	"image"
	"os"

	_ "image/png"

	"github.com/askeladdk/pancake/graphics2d"

	"github.com/askeladdk/pancake/graphics"
	gl "github.com/askeladdk/pancake/graphics/opengl"
	"github.com/askeladdk/pancake/input"
	"github.com/askeladdk/pancake/mathx"

	"github.com/askeladdk/pancake"
)

// https://opengameart.org/content/asteroids-game-sprites-atlas
// https://opengameart.org/content/purple-planet

func toggleFlag(flags uint32, flag uint32, state bool) uint32 {
	if state {
		return flags | flag
	} else {
		return flags &^ flag
	}
}

func loadTexture(filename string) (*graphics.Texture, error) {
	if f, err := os.Open(filename); err != nil {
		return nil, err
	} else if img, _, err := image.Decode(f); err != nil {
		return nil, err
	} else {
		return graphics.NewTextureFromImage(img, graphics.FilterNearest), nil
	}
}

func run(app pancake.App) error {
	var sheet *graphics.Texture
	var background graphics2d.Sprite

	if tex, err := loadTexture("asteroids-arcade.png"); err != nil {
		return err
	} else {
		sheet = tex
	}

	if tex, err := loadTexture("background.png"); err != nil {
		return err
	} else {
		background = graphics2d.NewSprite(tex, tex.Bounds())
	}

	ship := graphics2d.NewSprite(sheet, image.Rect(0, 0, 32, 32))
	asteroid := graphics2d.NewSprite(sheet, image.Rect(64, 192, 128, 256))

	resolution := app.Bounds().Size()
	projection := mathx.Ortho2D(
		0,
		float32(resolution.X),
		float32(resolution.Y),
		0,
	)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ClearColor(0, 0, 0, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	shipbrain := inputBrain{}

	simulation := simulation{
		sprites: []graphics2d.Sprite{
			background,
			ship,
			asteroid,
		},
		bounds: mathx.Rectangle{
			mathx.Vec2{},
			mathx.FromPoint(resolution),
		},
		entities: []entity{
			entity{
				sprite: background,
				pos:    mathx.FromPoint(app.Bounds().Size().Div(2)),
			},
			entity{
				sprite:        ship,
				pos:           mathx.FromPoint(app.Bounds().Size().Div(2)),
				rot:           -mathx.Tau / 4,
				minrotv:       1,
				maxv:          2,
				turn:          mathx.Tau / 4,
				thrust:        4,
				dampenr:       0.95,
				dampenv:       0.99,
				brain:         shipbrain.frame,
				collisionMask: 1,
				radius:        14,
			},
			// entity{
			// 	sprite:        asteroid,
			// 	pos:           mathx.FromPoint(app.Bounds().Size().Div(4)),
			// 	turn:          mathx.Tau / 256,
			// 	maxv:          1,
			// 	rotv:          1,
			// 	minrotv:       1,
			// 	dampenr:       1,
			// 	dampenv:       1,
			// 	vel:           mathx.FromHeading(mathx.Tau / 5),
			// 	collisionMask: 1,
			// 	radius:        28,
			// },
		},
		drawer:     graphics2d.NewDrawer(1024, graphics2d.Quad),
		shader:     graphics2d.DefaultShader(),
		projection: projection,
	}

	// var mousepos mathx.Vec2

	var keys uint32

	return app.Events(func(event interface{}) error {
		switch ev := event.(type) {
		case pancake.QuitEvent:
			return pancake.Quit
		case pancake.KeyEvent:
			switch ev.Key {
			case input.KeyEscape:
				return pancake.Quit
			case input.KeyA:
				fallthrough
			case input.KeyLeft:
				keys = toggleFlag(keys, 1, ev.Flags.Down())
			case input.KeyD:
				fallthrough
			case input.KeyRight:
				keys = toggleFlag(keys, 2, ev.Flags.Down())
			case input.KeyW:
				fallthrough
			case input.KeyUp:
				keys = toggleFlag(keys, 4, ev.Flags.Down())
			case input.KeyP:
				if ev.Flags.Pressed() {
					simulation.spawnAsteroid()
				}
			}
		// case pancake.MouseMoveEvent:
		// 	mousepos = mathx.FromPoint(ev.Position)
		case pancake.FrameEvent:
			// entity := simulation.At(shipid)
			// heading := mathx.FromHeading(entity.Rot)
			// target := mousepos.Sub(entity.Pos)
			// cross := target.Cross(heading)
			// if mathx.Abs(cross) < 1e-1 {
			// 	entity.Rot = target.Heading()
			// 	entity.RotV = 0
			// } else if cross < 0 {
			// 	entity.RotV = +mathx.Tau / 64
			// } else if cross > 0 {
			// 	entity.RotV = -mathx.Tau / 64
			// }
			if keys&1 != 0 {
				shipbrain.action(actionTurn, -1*float32(ev.DeltaTime))
			}

			if keys&2 != 0 {
				shipbrain.action(actionTurn, +1*float32(ev.DeltaTime))
			}

			if keys&4 != 0 {
				shipbrain.action(actionForward, 1*float32(ev.DeltaTime))
			}

			simulation.frame()
			app.SetTitle(fmt.Sprintf("Asteroids - %d FPS", app.FrameRate()))
		case pancake.DrawEvent:
			app.Begin()
			gl.ClearColor(0, 0, 1, 0)
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

			simulation.draw()

			app.End()
		}

		return nil
	})
}

func main() {
	opt := pancake.Options{
		WindowSize: image.Point{960, 540},
		Resolution: image.Point{640, 360},
		Title:      "Asteroids",
		FrameRate:  60,
	}

	if err := pancake.Main(opt, run); err != nil {
		fmt.Println(err)
	}
}
