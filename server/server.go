package server

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/dacousb/feiok/packet"
	. "github.com/dacousb/feiok/try"
	"github.com/google/uuid"
)

const (
	motd    = "feiok? prob not"
	timeOut = 10
)

const (
	wheatGrowthInterval = 10
)

type Tile struct {
	stage packet.WheatStage
	last  time.Time
}

type Player struct {
	x, y    float64
	looking packet.LookingAt
	name    string

	last time.Time
}

type Server struct {
	players       map[string]*Player
	tiles         [][]*Tile
	width, height int
	wheat         uint32

	mutex sync.RWMutex
}

func New() *Server {
	s := &Server{
		players: make(map[string]*Player),
		width:   16,
		height:  16,
		mutex:   sync.RWMutex{},
	}

	s.tiles = make([][]*Tile, s.height)
	for y := 0; y < s.height; y++ {
		s.tiles[y] = make([]*Tile, s.width)
		for x := 0; x < s.width; x++ {
			s.tiles[y][x] = &Tile{stage: packet.WHEAT_0}
		}
	}

	return s
}

func (s *Server) Run() {
	ln := Try(net.Listen("tcp", ":2022"))
	log.Println("listening on port :2022")
	for {
		conn := Try(ln.Accept())
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	uuid := uuid.New().String()
	log.Printf("new connection with %s\n", uuid)

	for {
		prefix, err := packet.ReadPrefix(conn)
		if err != nil {
			log.Printf("closed connection with %s\n", uuid)
			conn.Close()
			return
		}

		switch prefix {
		case packet.ASK_MOTD_CLIENT:
			s.askMotdClient(conn)
		case packet.SEND_PLAYER_CLIENT:
			s.sendPlayerClient(conn, uuid)
		case packet.ASK_PLAYERS_CLIENT:
			s.askPlayersClient(conn, uuid)
		case packet.SEND_PLANT_CLIENT:
			s.sendPlantClient(conn)
		case packet.ASK_PLANT_CLIENT:
			s.askPlantClient(conn)
		case packet.SEND_HARVEST_CLIENT:
			s.sendHarvestClient(conn)
		case packet.ASK_STATS_CLIENT:
			s.askStatsClient(conn)
		}
	}
}
