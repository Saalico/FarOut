package entity

import (
	"github.com/hajimehoshi/ebiten/v2"
	"sprite"
)

type coordinates struct {
	x, y float64
}

type dimensions struct {
	w, h float64
}

type vector struct {
	dx, dy float64
}

var (
	min     = dimensions{1, 1}
	vvsmall = dimensions{2, 2}
	vsmall  = dimensions{4, 4}
	small   = dimensions{8, 8}
	medium  = dimensions{16, 16}
	large   = dimensions{32, 32}
	vlarge  = dimensions{64, 64}
	vvlarge = dimensions{128, 128}
	max     = dimensions{256, 256}
)

type Direction int

const (
	Up    Direction = 0
	Right           = 1
	Down            = 2
	Left            = 3
)

var Directions = [4]string{"Up", "Right", "Down", "Left"}

type vitals struct {
	health, durability, stamina, hunger int
}

type action struct {
	frames   sprite.Animation
	canStop  bool
	velocity vector
}

type status struct {
	vitals
	coordinates
	dimensions
	facing       Direction
	primedAction *action
	special      map[string]bool
}
type This struct {
	name string
	status
	actions map[string]*action
	sprite  sprite.ImageData
}

func (e *This) init(name string, health, stamina int, size dimensions, atlas sprite.Atlas) *This {
	e.name = name
	e.status.vitals.health = health
	e.status.vitals.stamina = stamina
	e.status.vitals.hunger = 100
	e.status.vitals.durability = 1
	e.sprite.Atlas = atlas
	e.sprite.Image = atlas[0][0]
	e.actions = make(map[string]*action)
	e.status.primedAction = nil
	e.status.dimensions.w, e.status.dimensions.h = size.w, size.h
	e.status.coordinates = coordinates{0, 0}

	//TODO Change vector allow function that refers to global game speed.
	var vectors = []vector{{0, -2.5}, {2.5, 0}, {0, 2.5}, {-2.5, 0}}
	for i := range Directions {
		e.initAction("idle"+Directions[i], e.sprite.Atlas[i], 0, 2)
		e.actions["idle"+Directions[i]].velocity = vector{0, 0}

		e.initAction("walk"+Directions[i], e.sprite.Atlas[i], 2, 4)
		e.actions["walk"+Directions[i]].velocity = vectors[i]

	}
	e.status.primedAction = e.actions["idleDown"]
	e.sprite.Op = ebiten.DrawImageOptions{}
	e.sprite.Op.GeoM.Scale(2.5, 2.5)
	return e

}

func getFacing(v vector) int {
	var d int
	if v.dx > 0 && v.dy == 0 {
		d = int(Right)
	} else if v.dx < 0 && v.dy == 0 {
		d = int(Left)
	} else if v.dy < 0 && v.dx == 0 {
		d = int(Up)
	} else if v.dy > 0 && v.dx == 0 {
		d = int(Down)
	} else {
		return -1
	}
	return d
}

func (c *This) initAction(name string, frames sprite.Animation, actionStart, actionEnd int) {
	frames = frames[actionStart:actionEnd]
	c.actions[name] = &(action{
		frames,
		true,
		vector{},
	})
}

func (c *This) Prime(a *action) error {
	c.status.primedAction = a
	return nil
}

func (e *This) drawAnimationFrame(screen *ebiten.Image, a *action, frame int) {
	e.sprite.Op.GeoM.Translate(a.velocity.dx, a.velocity.dy)
	screen.DrawImage(a.frames[frame], &e.sprite.Op)
	e.sprite.Image = a.frames[frame]
}

func (e *This) Do(screen *ebiten.Image, a *action, frame int) {
	e.drawAnimationFrame(screen, a, frame)
	e.setFacing()
}

func (e *This) Status() status {
	return e.status
}

func (c *This) setFacing() {
	d := getFacing(c.primedAction.velocity)
	if d == -1 {
		return
	}
	c.status.facing = Direction(d)
}
