package main

import (
	"image/color"
	"math/rand"

	"github.com/askeladdk/pancake/graphics"
	"github.com/askeladdk/pancake/graphics2d"
	"github.com/askeladdk/pancake/mathx"
)

type actionCode int

const (
	actionForward actionCode = iota
	actionTurn
)

type action struct {
	code  actionCode
	value float32
}

type brainFunc func(*entity)

type inputBrain struct {
	actions []action
}

func (b *inputBrain) action(code actionCode, value float32) {
	b.actions = append(b.actions, action{code, value})
}

func (b *inputBrain) frame(e *entity) {
	for _, a := range b.actions {
		switch a.code {
		case actionForward:
			acc := mathx.FromHeading(e.rot).Mul(a.value * e.thrust)
			vel := e.vel.Add(acc)
			if vel.Len() > e.maxv {
				vel = vel.Unit().Mul(e.maxv)
			}
			e.vel = vel
		case actionTurn:
			e.rotv = e.turn * a.value
		}
	}
	b.actions = b.actions[:0]
}

type entity struct {
	sprite        graphics2d.Sprite // sprite
	pos           mathx.Vec2        // position
	vel           mathx.Vec2        // velocity
	rot           float32           // rotation
	rotv          float32           // rotational velocity
	dampenv       float32           // velocity dampening
	dampenr       float32           // rotational dampening
	maxv          float32           // maximum velocity
	minrotv       float32           // minimum rotational velocity
	turn          float32           // turn rate
	thrust        float32           // thrust speed
	brain         brainFunc         // intelligence
	collisionMask uint32
	radius        float32
}

type simulation struct {
	sprites    []graphics2d.Sprite
	bounds     mathx.Rectangle
	entities   []entity
	drawer     *graphics2d.Drawer
	shader     *graphics.ShaderProgram
	projection mathx.Mat4
}

type entities struct {
	s *simulation
	i int
	j int
}

func (es *entities) Len() int {
	return es.j - es.i
}

func (es *entities) Slice(i, j int) graphics2d.InstanceSlice {
	return &entities{
		s: es.s,
		i: es.i + i,
		j: es.i + j,
	}
}

func (es *entities) Color(i int) color.NRGBA {
	return color.NRGBA{0xff, 0xff, 0xff, 0xff}
}

func (es *entities) at(i int) entity {
	return es.s.entities[es.i+i]
}

func (es *entities) Texture(i int) *graphics.Texture {
	return es.at(i).sprite.Texture
}

func (es entities) TextureRegion(i int) mathx.Aff3 {
	return es.at(i).sprite.Region
}

func (es entities) ModelView(i int) mathx.Aff3 {
	e := es.at(i)
	return mathx.
		ScaleAff3(e.sprite.Size).
		Rotated(e.rot).
		Translated(e.pos)
}

func (s *simulation) collisionTest() {
	for i := 0; i < len(s.entities); i++ {
		a := s.at(i)
		if a.collisionMask == 0 {
			continue
		}
		for j := i + 1; j < len(s.entities); j++ {
			b := s.at(j)
			if a.collisionMask&b.collisionMask != 0 {
				c0 := mathx.Circle{a.pos, a.radius}
				c1 := mathx.Circle{b.pos, b.radius}
				if c0.IntersectsCircle(c1) {
					v := a.pos.Sub(b.pos).Unit()
					a.vel = v.Mul(a.maxv)
					b.vel = v.Mul(b.maxv).Neg()
					a.rotv += mathx.Tau / 64
					b.rotv -= mathx.Tau / 64
				}
			}
		}
	}
}

func (s *simulation) frame() {
	for i, e := range s.entities {
		if e.brain != nil {
			e.brain(&e)
			s.entities[i] = e
		}
	}

	s.collisionTest()

	for i, e := range s.entities {
		b := s.bounds.Expand(e.sprite.Size.Mul(0.5))
		e.pos = e.pos.Add(e.vel).Wrap(b)
		e.vel = e.vel.Mul(e.dampenv)
		e.rotv = mathx.Clamp(e.rotv*e.dampenr, -e.minrotv, e.minrotv)
		e.rot = mathx.Mod(e.rot+e.rotv*e.turn, mathx.Tau)
		s.entities[i] = e
	}
}

func (s *simulation) draw() {
	s.shader.Begin()
	s.shader.SetUniform("u_Projection", s.projection)
	s.drawer.Draw(&entities{s, 0, len(s.entities)})
	s.shader.End()
}

func (s *simulation) at(i int) *entity {
	return &s.entities[i]
}

func (s *simulation) spawnAsteroid() {
	pos := mathx.Vec2{
		rand.Float32(),
		rand.Float32(),
	}.MulVec2(s.bounds.Max)

	s.entities = append(s.entities, entity{
		sprite:        s.sprites[2],
		pos:           pos,
		turn:          mathx.Tau / 64 * (2*rand.Float32() - 1),
		maxv:          1,
		rotv:          1,
		minrotv:       rand.Float32(),
		dampenr:       .99,
		dampenv:       1,
		vel:           mathx.FromHeading(mathx.Tau * rand.Float32()),
		collisionMask: 1,
		radius:        28,
	})
}
