package server

import (
	"net"
	"time"

	"github.com/dacousb/feiok/packet"
)

func (s *Server) askMotdClient(conn net.Conn) {
	var p packet.Packet

	p.PushString(motd)
	conn.Write(packet.NewPacket(packet.SEND_MOTD_SERVER, p))
}

func (s *Server) sendPlayerClient(conn net.Conn, uuid string) {
	name, _ := packet.ReadString(conn)
	looking, _ := packet.ReadByte(conn)
	x, _ := packet.ReadFloat(conn)
	y, _ := packet.ReadFloat(conn)

	s.mutex.Lock()
	s.players[uuid] = &Player{x, y, packet.LookingAt(looking), name, time.Now()}
	s.mutex.Unlock()
}

func (s *Server) askPlayersClient(conn net.Conn, uuid string) {
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
}

func (s *Server) sendPlantClient(conn net.Conn) {
	x, _ := packet.ReadByte(conn)
	y, _ := packet.ReadByte(conn)

	s.mutex.Lock()
	s.tiles[y][x].stage = packet.WHEAT_1
	s.tiles[y][x].last = time.Now()
	s.mutex.Unlock()
}

func (s *Server) askPlantClient(conn net.Conn) {
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

func (s *Server) sendHarvestClient(conn net.Conn) {
	x, _ := packet.ReadByte(conn)
	y, _ := packet.ReadByte(conn)

	s.mutex.Lock()
	if s.tiles[y][x].stage == packet.WHEAT_4 {
		s.wheat += 1
		s.tiles[y][x].stage = packet.WHEAT_0
	}
	s.mutex.Unlock()
}

func (s *Server) askStatsClient(conn net.Conn) {
	var p packet.Packet

	s.mutex.RLock()
	p.PushUint32(s.wheat)
	s.mutex.RUnlock()

	conn.Write(packet.NewPacket(packet.SEND_STATS_SERVER, p))
}
