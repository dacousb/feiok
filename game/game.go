package game

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/dacousb/feiok/packet"
)

const (
	screenWidth  = 960
	screenHeight = 540
)

type Game struct {
	tiles         [][]*Tile
	main          *Player
	players       []*Player
	width, height int
	wheat         uint32

	host       string
	motd       string
	conn       net.Conn
	conn_mutex sync.RWMutex
	data_mutex sync.RWMutex
}

func New() *Game {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("feiok")

	game := &Game{}
	game.loadTiles()
	game.loadMain()

	return game
}

func (g *Game) Run(host string) {
	g.setHost(host)
	g.sendPlayer()

	go g.askPlayers()
	go g.askPlant()
	go g.askStats()

	go g.responsePool()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.main.y -= 0.1
		g.main.looking = packet.LOOKING_B
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.main.x -= 0.1
		g.main.looking = packet.LOOKING_B
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.main.y += 0.1
		g.main.looking = packet.LOOKING_L
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.main.x += 0.1
		g.main.looking = packet.LOOKING_R
	}
	g.fixPosition()
	g.sendPlayer()

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.sendPlant()
	}

	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		g.sendHarvest()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawTiles(screen)
	g.drawPlayers(screen)
	ebitenutil.DebugPrint(screen,
		fmt.Sprintf("fps: %0.f\nhost: %s\nmotd: %s\nx: %0.2f y: %0.2f\nwheat: %d",
			ebiten.CurrentFPS(), g.host, g.motd, g.main.x, g.main.y, g.wheat))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
