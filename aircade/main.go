package main

import (
	"embed"
	_ "embed"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var img *ebiten.Image

//go:embed leaf_ranger/leaf_ranger.png
var data embed.FS

func init() {
	var err error
	img, _, err = ebitenutil.NewImageFromFile("leaf_ranger/leaf_ranger.png")
	if err != nil {
		log.Fatal(err)
	}
}

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(img, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Aircade")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
