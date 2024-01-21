package world

import (
	ent "github.com/saalico/FarOut/entity"
)

type World struct {
	name              string
	camera            *Camera
	matrix            *Matrix
	entities          []*ent
	activeEntities    []*ent
	offScreenEntities []*ent
}

func initWorld(name string) World {
	var w World
	w.name = name
	w.camera = &Camera{}
	w.matrix = &Matrix{}
	w.entities = make([]*ent, 0)
	w.activeEntities = make([]*ent, 0)
	w.offScreenEntities = make([]*ent, 0)
	return w
}
