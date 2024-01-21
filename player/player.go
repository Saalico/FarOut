package player

import (
	"camera"
	"ent"
	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	playerID       uint8
	selectedEntity *ent.Entity
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
