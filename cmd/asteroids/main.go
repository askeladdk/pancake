package main

import (
	"fmt"
	"image"
	"image/color"
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

type Entity struct {
	Image  graphics2d.Sprite
	Pos    mathx.Vec2
	Vel    mathx.Vec2
	Acc    mathx.Vec2
	Dampen float32
	Rot    float32
	RotV   float32
	Alive  bool
}

type Entities []Entity

func (es Entities) Len() int {
	return len(es)
}

func (es Entities) Slice(i, j int) graphics2d.InstanceSlice {
	return es[i:j]
}

func (es Entities) Color(i int) color.NRGBA {
	return color.NRGBA{0xff, 0xff, 0xff, 0xff}
}

func (es Entities) Texture(i int) *graphics.Texture {
	return es[i].Image.Texture
}

func (es Entities) TextureRegion(i int) mathx.Aff3 {
	return es[i].Image.Region
}

func (es Entities) ModelView(i int) mathx.Aff3 {
	return mathx.
		ScaleAff3(es[i].Image.Size).
		Rotated(es[i].Rot).
		Translated(es[i].Pos)
}

type World struct {
	Bounds image.Rectangle
}

func (entities Entities) Frame() Entities {
	for i, e := range entities {
		e.Vel = e.Vel.Add(e.Acc)
		e.Acc = e.Acc.Mul(e.Dampen)
		e.Vel[0] = mathx.Clamp(e.Vel[0], -2, 2)
		e.Vel[1] = mathx.Clamp(e.Vel[1], -2, 2)
		e.Pos = e.Pos.Add(e.Vel)
		e.Rot = e.Rot + e.RotV
		entities[i] = e
	}

	return entities
}

func LoadTexture(filename string) (*graphics.Texture, error) {
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

	if tex, err := LoadTexture("asteroids-arcade.png"); err != nil {
		return err
	} else {
		sheet = tex
	}

	if tex, err := LoadTexture("background.png"); err != nil {
		return err
	} else {
		background = graphics2d.NewSprite(tex, tex.Bounds())
	}

	ship := graphics2d.NewSprite(sheet, image.Rect(0, 0, 32, 32))

	drawer := graphics2d.NewDrawer(1024, graphics2d.Quad)
	program := graphics2d.DefaultShader()

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

	entities := Entities{
		Entity{
			Image: background,
			Pos:   mathx.FromPoint(app.Bounds().Size().Div(2)),
		},
		Entity{
			Image:  ship,
			Pos:    mathx.FromPoint(app.Bounds().Size().Div(2)),
			Rot:    -mathx.Tau / 4,
			Dampen: 0,
		},
	}

	shipid := 1

	var mousepos mathx.Vec2

	return app.Events(func(event interface{}) error {
		switch ev := event.(type) {
		case pancake.QuitEvent:
			return pancake.Quit
		case pancake.KeyEvent:
			switch ev.Key {
			case input.KeyEscape:
				return pancake.Quit
			// case input.KeyLeft:
			// 	if ev.Flags.Down() {
			// 		entities[shipid].Rot -= mathx.Tau / 32
			// 	}
			// case input.KeyRight:
			// 	if ev.Flags.Down() {
			// 		entities[shipid].Rot += mathx.Tau / 32
			// 	}
			case input.KeyW:
				fallthrough
			case input.KeyUp:
				if ev.Flags.Down() {
					acc := mathx.FromHeading(entities[shipid].Rot)
					entities[shipid].Acc = acc
				}
			}
		case pancake.MouseMoveEvent:
			mousepos = mathx.FromPoint(ev.Position)
		case pancake.FrameEvent:
			heading := mathx.FromHeading(entities[shipid].Rot)
			target := mousepos.Sub(entities[shipid].Pos)
			cross := target.Cross(heading)
			if mathx.Abs(cross) < 1e-3 {
				entities[shipid].RotV = 0
			} else if cross < 0 {
				entities[shipid].RotV = +mathx.Tau / 64
			} else if cross > 0 {
				entities[shipid].RotV = -mathx.Tau / 64
			}
			entities = entities.Frame()
			app.SetTitle(fmt.Sprintf("Asteroids - %d FPS | %v", app.FrameRate(), mousepos))
		case pancake.DrawEvent:
			app.Begin()
			gl.ClearColor(0, 0, 1, 0)
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

			program.Begin()
			program.SetUniform("u_Projection", projection)
			drawer.Draw(entities)
			program.End()

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
