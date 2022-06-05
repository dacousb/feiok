package game

import (
	"bytes"
	"embed"
	"image/png"

	. "github.com/dacousb/feiok/try"
	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/*
var assets embed.FS

var tile *ebiten.Image
var player_b *ebiten.Image
var player_l *ebiten.Image
var player_r *ebiten.Image
var wheat_1 *ebiten.Image
var wheat_2 *ebiten.Image
var wheat_3 *ebiten.Image
var wheat_4 *ebiten.Image

const spriteSize = 32

func readImage(name string) (*ebiten.Image, error) {
	f, err := assets.ReadFile("assets/" + name)
	if err != nil {
		return nil, err
	}
	img, err := png.Decode(bytes.NewReader(f))
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(img), nil
}

func init() {
	tile = Try(readImage("tile.png"))
	player_b = Try(readImage("player_b.png"))
	player_l = Try(readImage("player_l.png"))
	player_r = Try(readImage("player_r.png"))
	wheat_1 = Try(readImage("wheat_1.png"))
	wheat_2 = Try(readImage("wheat_2.png"))
	wheat_3 = Try(readImage("wheat_3.png"))
	wheat_4 = Try(readImage("wheat_4.png"))
}
