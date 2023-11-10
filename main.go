package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"src/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 320
	screenHeight = 240
	maxAngle     = 360
)

type Direction int

func castArray[A any](l map[int]int, s []A) []A {
	var b = make([]A, len(s))
	fmt.Println(len(b), len(s))
	for i, v := range s {
		b[l[i]] = v
		fmt.Println(l[i])
	}
	return b
}

const (
	Up    Direction = 0
	Right           = 1
	Down            = 2
	Left            = 3
)

var Directions = [4]string{"Up", "Right", "Down", "Left"}

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

var defaultKeyMap = map[ebiten.Key]string{
	ebiten.KeyW: "walkUp",
	ebiten.KeyA: "walkLeft",
	ebiten.KeyS: "walkDown",
	ebiten.KeyD: "walkRight",
}

type spritesheet []byte

type animation []*ebiten.Image

type atlas []animation

type action struct {
	frames   animation
	canStop  bool
	velocity vector
}

type Game struct {
	Characters       []Character
	players          map[string]Player
	needUpdate       bool
	primedCharacters []*Character
	speed            int
	count            int
	inited           bool
}

func buildSpriteAtlas(rows, columns int, sprite spritesheet) atlas {
	img, _, err := image.Decode(bytes.NewReader(sprite))
	if err != nil {
		log.Fatal(err)
	}
	origEbitenImage := ebiten.NewImageFromImage(img)
	imgs := make([]animation, rows)
	for i := range imgs {
		imgs[i] = make(animation, columns)
		for j := range imgs[i] {
			spriteWidth := origEbitenImage.Bounds().Dx() / columns
			spriteHeight := origEbitenImage.Bounds().Dy() / rows
			offsetX := j * spriteWidth
			offsetY := i * spriteHeight
			imgs[i][j] = origEbitenImage.SubImage(
				image.Rect(offsetX, offsetY, offsetX+spriteWidth, offsetY+spriteHeight)).(*ebiten.Image)

		}
	}
	return castArray(map[int]int{0: 2, 1: 0, 2: 3, 3: 1}, imgs)
}

type Character struct {
	name string
	status
	actions map[string]*action
	imageData
}

type vitals struct {
	health, stamina, hunger int
}

type imageData struct {
	image *ebiten.Image
	op    ebiten.DrawImageOptions
	atlas
}

type status struct {
	vitals
	coordinates
	dimensions
	facing       Direction
	primedAction *action
	special      map[string]bool
}

func initChar(name string, health, stamina int, size dimensions, atlas atlas) Character {
	var c Character
	c.name = name
	c.status.vitals.health = health
	c.status.vitals.stamina = stamina
	c.status.vitals.hunger = 100
	c.imageData.atlas = atlas
	c.imageData.image = atlas[0][0]
	c.actions = make(map[string]*action)
	c.status.primedAction = nil
	c.status.dimensions.w, c.status.dimensions.h = size.w, size.h
	c.status.coordinates = coordinates{0, 0}

	var vectors = []vector{{0, -2.5}, {2.5, 0}, {0, 2.5}, {-2.5, 0}}
	for i := range Directions {
		c.initAction("idle"+Directions[i], c.atlas[i], 0, 2)
		c.actions["idle"+Directions[i]].velocity = vector{0, 0}

		c.initAction("walk"+Directions[i], c.atlas[i], 2, 4)
		c.actions["walk"+Directions[i]].velocity = vectors[i]

	}
	c.status.primedAction = c.actions["idleDown"]
	return c

}

func (c *Character) initAction(name string, frames animation, actionStart, actionEnd int) {
	frames = frames[actionStart:actionEnd]
	c.actions[name] = &(action{
		frames,
		true,
		vector{},
	})
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

type Player struct {
	playerID          uint8
	selectedCharacter *Character
	keyMap            map[ebiten.Key]string
}

func (p *Player) getInput() error {
	c := p.selectedCharacter
	keys := make([]ebiten.Key, 0)
	keys = inpututil.AppendPressedKeys(keys)
	if len(keys) < 1 {
		c.status.primedAction = c.actions["idle"+Directions[c.status.facing]]
		fmt.Println("idle" + Directions[c.status.facing])
		return nil
	}
	_, ok := p.keyMap[keys[len(keys)-1]]
	if ok && p.selectedCharacter.primedAction.canStop {
		c.status.primedAction = c.actions[p.keyMap[keys[len(keys)-1]]]
		return nil
	}
	return nil
}

func (g *Game) init() {
	defer func() {
		g.inited = true
	}()
	g.speed = 1

	//init all player entities
	benj := initChar("Benj", 100, 100, dimensions{16, 16}, buildSpriteAtlas(4, 4, images.BasicCharacterSpritesheet_png))

	g.players = make(map[string]Player)

	var p1 Player
	p1.playerID = 1
	p1.selectedCharacter = &benj
	p1.keyMap = defaultKeyMap
	g.players["p1"] = p1
	g.primedCharacters = make([]*Character, 0)
}

func (c *Character) primeAction(a *action) error {
	c.status.primedAction = a
	return nil
}

func (c *Character) drawAnimationFrame(screen *ebiten.Image, a *action, frame int) {
	c.op.GeoM.Translate(a.velocity.dx, a.velocity.dy)
	screen.DrawImage(a.frames[frame], &c.op)
	c.image = a.frames[frame]
}

func (c *Character) do(screen *ebiten.Image, a *action, frame int) {
	c.drawAnimationFrame(screen, a, frame)
	c.setFacing()
}

func (c *Character) setFacing() {
	d := getFacing(c.primedAction.velocity)
	if d == -1 {
		return
	}
	c.status.facing = Direction(d)
	fmt.Println(c.status.facing)
}

func (g *Game) Update() error {
	if !g.inited {
		g.init()
	}

	g.primedCharacters = g.primedCharacters[:0]
	for _, p := range g.players {
		p.getInput()
		if p.selectedCharacter.status.primedAction != nil {
			g.primedCharacters = append(g.primedCharacters, p.selectedCharacter)
		}
	}

	g.count += 1 * g.speed
	g.count %= 10
	return nil
}
func (g *Game) showAction(c Character, a action, screen *ebiten.Image) {
	op := &c.op
	for _, img := range a.frames {
		sprite := img
		screen.DrawImage(sprite, op)
		op.GeoM.Translate(c.dimensions.w, 0)
	}
}
func (g *Game) showSpritesheet(character Character, screen *ebiten.Image, offset dimensions) {
	op := &character.op
	op.GeoM.Translate(offset.w, offset.h)
	for _, ani := range character.atlas {
		for _, img := range ani {
			sprite := img
			screen.DrawImage(sprite, op)
			op.GeoM.Translate(character.dimensions.w, 0)
		}

		op.GeoM.Translate(-character.dimensions.w*float64(len(ani)), 0)
		op.GeoM.Translate(0, character.dimensions.h+2)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	//g.showAction(*c, *c.actions["walkUp"], screen)
	if len(g.primedCharacters) > 0 {
		for _, c := range g.primedCharacters {
			i := (g.count / 6) % len(c.status.primedAction.frames)
			c.do(screen, c.status.primedAction, i)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Sprites (Ebitengine Demo)")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}

}
