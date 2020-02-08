package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/askeladdk/pancake/graphics"
	gl "github.com/askeladdk/pancake/graphics/opengl"
	"github.com/askeladdk/pancake/input"
	"github.com/askeladdk/pancake/mathx"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/askeladdk/pancake"
)

// https://opengameart.org/content/asteroids-game-sprites-atlas

type Entity struct {
	Image *Image
	Pos   mgl32.Vec2
	Vel   mgl32.Vec2
	Acc   mgl32.Vec2
	// Rot   float32
	// RVel  float32
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
	for _, e := range entities {
		e.Pos = e.Pos.Add(e.Vel)
		e.Vel = e.Vel.Add(e.Acc)
	}

	return entities
}

func run(app pancake.App) error {
	var sprites *Image
	var spriteDrawer *SpriteDrawer

	if img, err := LoadPNG("asteroids-arcade.png"); err != nil {
		return err
	} else {
		sprites = img
	}

	ship := sprites.SubImage(image.Rect(0, 0, 32, 32))

	position := mgl32.Vec2{160, 90}

	if sd, err := NewSpriteDrawer(app.Bounds().Size()); err != nil {
		return err
	} else {
		spriteDrawer = sd
	}

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ClearColor(0, 0, 0, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	var velocity mgl32.Vec2

	return app.Events(func(event interface{}) error {
		switch ev := event.(type) {
		case pancake.QuitEvent:
			return pancake.Quit
		case pancake.KeyEvent:
			switch ev.Key {
			case input.KeyEscape:
				return pancake.Quit
			case input.KeyLeft:
				velocity[0] = 1
			case input.KeyRight:
				velocity[0] = 0
			case input.KeyUp:
				velocity[1] = 1
			case input.KeyDown:
				velocity[1] = 0
			}
		case pancake.FrameEvent:
			app.SetTitle(fmt.Sprintf("Asteroids - %d FPS", app.FrameRate()))
		case pancake.DrawEvent:
			app.Begin()
			gl.ClearColor(0, 0, 1, 0)
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

			spriteDrawer.Begin()

			spriteDrawer.DrawImage(ship, position)

			spriteDrawer.End()

			app.End()
		}

		return nil
	})
}

func main() {
	opt := pancake.Options{
		WindowSize: image.Point{960, 540},
		Resolution: image.Point{320, 180},
		Title:      "Asteroids",
		FrameRate:  60,
	}

	if err := pancake.Main(opt, run); err != nil {
		fmt.Println(err)
	}
}
