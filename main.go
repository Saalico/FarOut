package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/mobile/asset"
)

type Matrix struct {
	nodes  map[coordinates]*Entity
	sector quadtree.Quadtree
}

type World struct {
	name              string
	camera            *Camera
	matrix            *Matrix
	entities          []*entity
	activeEntities    []*entity
	offScreenEntities []*entity
}

func initWorld(name string) World {
	var w World
	w.name = name
	w.camera = &Camera{}
	w.matrix = &Matrix{}
	w.entities = make([]*Entity, 0)
	w.activeEntities = make([]*Entity, 0)
	w.offScreenEntities = make([]*Entity, 0)
	return w
}

type Player struct {
	playerID       uint8
	selectedEntity *Entity
	keyMap         map[ebiten.Key]string
	camOp          cameraOptions
}

func (p *Player) getInput() error {
	c := p.selectedEntity
	keys := make([]ebiten.Key, 0)
	keys = inpututil.AppendPressedKeys(keys)
	if len(keys) < 1 {
		c.status.primedAction = c.actions["idle"+Directions[c.status.facing]]
		return nil
	}
	_, ok := p.keyMap[keys[len(keys)-1]]
	if ok && p.selectedEntity.primedAction.canStop {
		c.status.primedAction = c.actions[p.keyMap[keys[len(keys)-1]]]
		return nil
	}
	return nil
}

//func initCameraPanControls() {
//	for i := range Directions {
//		var c Camera
//		dir := Directions[i]
//		switch dir {
//		case "Up":
//		case "Right":
//		case "Down":
//		case "Left":
//
//		}
//	}
//
//}

func initUI() {
}

type Camera struct {
	player         *Player
	selectedEntity *Entity
	offset         vector
	zoom           float32
	actions        []*action
	dimensions
	coordinates
}

type cameraOptions struct {
	tetherRange int
	zoom        float32
	locked      bool
}

func (c *Camera) getOffset() (float64, float64) {
	return c.offset.dx, c.offset.dy
}

func (c *Camera) snapToSelectedEntity() {
	if c.player.selectedEntity != c.selectedEntity {
		c.selectedEntity = c.player.selectedEntity
		c.coordinates = c.player.selectedEntity.coordinates
	}
}

type Game struct {
	world   World
	players map[string]Player
	count   int
	inited  bool
	options map[string]float32
}

func (g *Game) init() {
	defer func() {
		g.inited = true
	}()
	g.options = make(map[string]float32)
	g.options["speed"] = 1.0

	g.world = initWorld("dev")
	//init all player entities
	benj := initEntity("Benj", 100, 100, medium, buildSpriteAtlas(4, 4, asset.BasicCharacterSpritesheet_png))

	g.players = make(map[string]Player)
	var p1 Player
	p1.playerID = 1
	p1.selectedEntity = &benj
	p1.keyMap = defaultKeyMap
	g.players["p1"] = p1

	g.world.entities = append(g.world.entities, &benj)

}

func (g *Game) Update() error {
	if !g.inited {
		g.init()
	}

	// Poll game buffer for player entities that are taking an action
	// In reality this is just going to get the last player input for their selected Entity
	g.world.activeEntities = g.world.activeEntities[:0]

	for _, p := range g.players {
		p.getInput()
		if p.selectedEntity.status.primedAction != nil {
			g.world.activeEntities = append(g.world.activeEntities, p.selectedEntity)
		}
	}

	g.count += 1 * int(g.options["speed"])
	g.count %= 10
	return nil
}

func (g *Game) showAction(c Entity, a action, screen *ebiten.Image) {
	op := &c.op
	for _, img := range a.frames {
		sprite := img
		screen.DrawImage(sprite, op)
		op.GeoM.Translate(c.dimensions.w, 0)
	}
}
func (g *Game) showSpritesheet(entity Entity, screen *ebiten.Image, offset dimensions) {
	op := &entity.op
	op.GeoM.Translate(offset.w, offset.h)
	for _, ani := range entity.atlas {
		for _, img := range ani {
			sprite := img
			screen.DrawImage(sprite, op)
			op.GeoM.Translate(entity.dimensions.w, 0)
		}

		op.GeoM.Translate(-entity.dimensions.w*float64(len(ani)), 0)
		op.GeoM.Translate(0, entity.dimensions.h+2)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	//g.showAction(*c, *c.actions["walkUp"], screen)
	if len(g.world.activeEntities) > 0 {
		for _, c := range g.world.activeEntities {
			i := (g.count / 6) % len(c.status.primedAction.frames)
			c.do(screen, c.status.primedAction, i)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	w, h := ebiten.ScreenSizeInFullscreen()
	ebiten.SetFullscreen(true)
	ebiten.SetWindowSize(w, h)
	return w, h
}

func main() {
	ebiten.SetWindowTitle("Sprites (Ebitengine Demo)")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}

}
