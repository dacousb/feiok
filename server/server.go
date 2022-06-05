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
			var p packet.Packet

			p.PushString(motd)
			conn.Write(packet.NewPacket(packet.SEND_MOTD_SERVER, p))

		case packet.SEND_PLAYER_CLIENT:
			name, _ := packet.ReadString(conn)
			looking, _ := packet.ReadByte(conn)
			x, _ := packet.ReadFloat(conn)
			y, _ := packet.ReadFloat(conn)

			s.mutex.Lock()
			s.players[uuid] = &Player{x, y, packet.LookingAt(looking), name, time.Now()}
			s.mutex.Unlock()

		case packet.ASK_PLAYERS_CLIENT:
			var p packet.Packet

			s.mutex.Lock()
			p.PushUint32(uint32(len(s.players) - 1))
			for k, v := range s.players {
				if k == uuid {
					continue
				}
				p.PushString(v.name)
				p.PushByte(byte(v.looking))
				p.PushFloat(v.x)
				p.PushFloat(v.y)
				if time.Since(v.last).Seconds() > timeOut {
					delete(s.players, k)
				}
			}
			s.mutex.Unlock()

			conn.Write(packet.NewPacket(packet.SEND_PLAYERS_SERVER, p))

		case packet.SEND_PLANT_CLIENT:
			x, _ := packet.ReadByte(conn)
			y, _ := packet.ReadByte(conn)

			s.mutex.Lock()
			s.tiles[y][x].stage = packet.WHEAT_1
			s.tiles[y][x].last = time.Now()
			s.mutex.Unlock()

		case packet.ASK_PLANT_CLIENT:
			var p packet.Packet

			s.mutex.Lock()
			p.PushByte(byte(s.width))
			p.PushByte(byte(s.height))
			for y := 0; y < s.height; y++ {
				for x := 0; x < s.width; x++ {
					if time.Since(s.tiles[y][x].last).Seconds() > wheatGrowthInterval &&
						s.tiles[y][x].stage > packet.WHEAT_0 && s.tiles[y][x].stage < packet.WHEAT_4 {

						s.tiles[y][x].last = time.Now()
						s.tiles[y][x].stage += 1
					}
					p.PushByte(byte(s.tiles[y][x].stage))
				}
			}
			s.mutex.Unlock()

			conn.Write(packet.NewPacket(packet.SEND_PLANT_SERVER, p))
		}
	}
}
