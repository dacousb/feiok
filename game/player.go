package game

import (
	"fmt"
	"math/rand"

	"github.com/dacousb/feiok/packet"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Player struct {
	x, y    float64
	looking packet.LookingAt
	name    string
}

func (g *Game) loadMain() {
	g.main = &Player{
		x: 0, y: 0,
		looking: packet.LOOKING_L,
		name:    fmt.Sprintf("guest%d", rand.Intn(11)),
	}
}

func (g *Game) drawPlayers(screen *ebiten.Image) {
	g.data_mutex.RLock()

	for i := -1; i < len(g.players); i++ {
		var player *Player
		if i == -1 {
			player = g.main
		} else {
			player = g.players[i]
		}

		op := &ebiten.DrawImageOptions{}
		x, y := g.getIsoCoords(player.x, player.y, 1)
		op.GeoM.Translate(x, y)

		screen.DrawImage(getSprite(player.looking), op)
		ebitenutil.DebugPrintAt(screen, player.name, int(x), int(y)-15)
	}

	g.data_mutex.RUnlock()
}

func (g *Game) getIsoCoords(x, y, z float64) (float64, float64) {
	x_iso := (spriteSize * ((x - y) / 2))
	y_iso := (spriteSize * ((x + y) / 4)) - z*spriteSize/2
	return x_iso + (screenWidth-spriteSize)/2,
		y_iso + (screenHeight-float64(spriteSize*g.height)/2)/2
}

func (g *Game) fixPosition() {
	if g.main.x > float64(g.width)-1 {
		g.main.x = float64(g.width) - 1
	} else if g.main.x < 0 {
		g.main.x = 0
	}
	if g.main.y > float64(g.height)-1 {
		g.main.y = float64(g.height) - 1
	} else if g.main.y < 0 {
		g.main.y = 0
	}
}

func getSprite(b packet.LookingAt) *ebiten.Image {
	switch b {
	case packet.LOOKING_B:
		return player_b
	case packet.LOOKING_L:
		return player_l
	case packet.LOOKING_R:
		return player_r
	default:
		return player_l
	}
}
