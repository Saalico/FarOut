package sprite

import (
	"bytes"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type spritesheet []byte

type Animation []*ebiten.Image

type Atlas []Animation

type ImageData struct {
	Image *ebiten.Image
	Op    ebiten.DrawImageOptions
	Atlas
}

func buildSpriteAtlas(rows, columns int, sprite spritesheet) Atlas {
	img, _, err := image.Decode(bytes.NewReader(sprite))
	if err != nil {
		log.Fatal(err)
	}
	origEbitenImage := ebiten.NewImageFromImage(img)
	imgs := make([]Animation, rows)
	for i := range imgs {
		imgs[i] = make(Animation, columns)
		for j := range imgs[i] {
			spriteWidth := origEbitenImage.Bounds().Dx() / columns
			spriteHeight := origEbitenImage.Bounds().Dy() / rows
			offsetX := j * spriteWidth
			offsetY := i * spriteHeight
			imgs[i][j] = origEbitenImage.SubImage(
				image.Rect(offsetX, offsetY, offsetX+spriteWidth, offsetY+spriteHeight)).(*ebiten.Image)
		}
	}
	return imgs
}
