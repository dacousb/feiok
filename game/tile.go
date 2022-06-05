package game

import (
	"github.com/dacousb/feiok/packet"
	"github.com/hajimehoshi/ebiten/v2"
)

type Tile struct {
	stage packet.WheatStage
}

func (g *Game) loadTiles() {
	g.width, g.height = 16, 16
	g.tiles = make([][]*Tile, g.height)

	for y := 0; y < g.height; y++ {
		g.tiles[y] = make([]*Tile, g.width)

		for x := 0; x < g.width; x++ {
			g.tiles[y][x] = &Tile{stage: packet.WHEAT_0}
		}
	}
}

func (g *Game) drawTiles(screen *ebiten.Image) {
	g.data_mutex.RLock()

	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(g.getIsoCoords(float64(x), float64(y), 0))
			screen.DrawImage(tile, op)

			op = &ebiten.DrawImageOptions{}
			op.GeoM.Translate(g.getIsoCoords(float64(x), float64(y), 1))

			switch g.tiles[y][x].stage {
			case packet.WHEAT_1:
				screen.DrawImage(wheat_1, op)
			case packet.WHEAT_2:
				screen.DrawImage(wheat_2, op)
			case packet.WHEAT_3:
				screen.DrawImage(wheat_3, op)
			case packet.WHEAT_4:
				screen.DrawImage(wheat_4, op)
			}
		}
	}

	g.data_mutex.RUnlock()
}
