package game

import (
	"log"
	"math"
	"net"
	"time"

	"github.com/dacousb/feiok/packet"
	. "github.com/dacousb/feiok/try"
)

const serverInterval = 15

func (g *Game) setHost(host string) {
	conn := Try(net.Dial("tcp", host))
	g.host = host
	g.conn = conn

	g.sendPacket(packet.NewPacket(packet.ASK_MOTD_CLIENT, packet.Packet{}))
}

func (g *Game) askPlayers() {
	for {
		time.Sleep(serverInterval * time.Millisecond)
		g.sendPacket(packet.NewPacket(packet.ASK_PLAYERS_CLIENT, packet.Packet{}))
	}
}

func (g *Game) askPlant() {
	for {
		time.Sleep(serverInterval * time.Millisecond)
		g.sendPacket(packet.NewPacket(packet.ASK_PLANT_CLIENT, packet.Packet{}))
	}
}

func (g *Game) askStats() {
	for {
		time.Sleep(serverInterval * time.Millisecond)
		g.sendPacket(packet.NewPacket(packet.ASK_STATS_CLIENT, packet.Packet{}))
	}
}

func (g *Game) sendPlayer() {
	var p packet.Packet
	p.PushString(g.main.name)
	p.PushByte(byte(g.main.looking))
	p.PushFloat(g.main.x)
	p.PushFloat(g.main.y)

	g.sendPacket(packet.NewPacket(packet.SEND_PLAYER_CLIENT, p))
}

func (g *Game) sendPlant() {
	var p packet.Packet
	p.PushByte(byte(math.Round(g.main.x)))
	p.PushByte(byte(math.Round(g.main.y)))

	g.sendPacket(packet.NewPacket(packet.SEND_PLANT_CLIENT, p))
}

func (g *Game) sendHarvest() {
	var p packet.Packet
	p.PushByte(byte(math.Round(g.main.x)))
	p.PushByte(byte(math.Round(g.main.y)))

	g.sendPacket(packet.NewPacket(packet.SEND_HARVEST_CLIENT, p))
}

func (g *Game) sendPacket(b []byte) {
	g.conn_mutex.Lock()
	g.conn.Write(b)
	g.conn_mutex.Unlock()
}

func (g *Game) responsePool() {
	for {
		prefix := Try(packet.ReadByte(g.conn))

		switch prefix {
		case packet.SEND_MOTD_SERVER:
			g.motd = Try(packet.ReadString(g.conn))

		case packet.SEND_PLAYERS_SERVER:
			length := Try(packet.ReadUint32(g.conn))

			g.data_mutex.Lock()
			g.players = []*Player{}
			for i := 0; i < int(length); i++ {
				name := Try(packet.ReadString(g.conn))
				looking := Try(packet.ReadByte(g.conn))
				x := Try(packet.ReadFloat(g.conn))
				y := Try(packet.ReadFloat(g.conn))

				g.players = append(g.players, &Player{x, y, packet.LookingAt(looking), name})
			}
			g.data_mutex.Unlock()

		case packet.SEND_PLANT_SERVER:
			x := int(Try(packet.ReadByte(g.conn)))
			y := int(Try(packet.ReadByte(g.conn)))

			g.data_mutex.Lock()
			for yi := 0; yi < y; yi++ {
				for xi := 0; xi < x; xi++ {
					g.tiles[yi][xi].stage = packet.WheatStage(Try(packet.ReadByte(g.conn)))
				}
			}
			g.data_mutex.Unlock()

		case packet.SEND_STATS_SERVER:
			g.wheat = Try(packet.ReadUint32(g.conn))

		default:
			log.Fatalf("got an unknown packet prefix (%d)", prefix)
		}
	}
}
