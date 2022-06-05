package packet

import (
	"encoding/binary"
	"io"
	"math"
	"net"
)

const (
	ASK_MOTD_CLIENT     byte = iota // the client asks the server for its MOTD
	SEND_MOTD_SERVER                // the server sends its MOTD
	SEND_PLAYER_CLIENT              // the client sends info about its player
	ASK_PLAYERS_CLIENT              // the client asks for the player list
	SEND_PLAYERS_SERVER             // the server sends the player list
	SEND_PLANT_CLIENT               // the client sends a new planted tile
	ASK_PLANT_CLIENT                // the client asks for the tiles
	SEND_PLANT_SERVER               // the server sends the tiles
)

type LookingAt byte

const (
	LOOKING_B LookingAt = iota
	LOOKING_L
	LOOKING_R
)

type WheatStage byte

const (
	WHEAT_0 WheatStage = iota
	WHEAT_1
	WHEAT_2
	WHEAT_3
	WHEAT_4
)

type Packet []byte

func (p *Packet) PushFloat(n float64) {
	f := make([]byte, 8)
	binary.LittleEndian.PutUint64(f, math.Float64bits(n))
	*p = append(*p, f...)
}

func (p *Packet) PushUint32(n uint32) {
	i := make([]byte, 4)
	binary.LittleEndian.PutUint32(i, uint32(n))
	*p = append(*p, i...)
}

func (p *Packet) PushByte(n byte) {
	*p = append(*p, n)
}

func (p *Packet) PushString(s string) {
	p.PushUint32(uint32(len(s)))
	*p = append(*p, []byte(s)...)
}

func ReadPrefix(conn net.Conn) (byte, error) {
	buff := make([]byte, 1)
	_, err := io.ReadFull(conn, buff)
	if err != nil {
		return 0, err
	}

	return buff[0], nil
}

func ReadFloat(conn net.Conn) (float64, error) {
	buff := make([]byte, 8)
	_, err := io.ReadFull(conn, buff)
	if err != nil {
		return 0, err
	}

	return math.Float64frombits(binary.LittleEndian.Uint64(buff)), nil
}

func ReadUint32(conn net.Conn) (uint32, error) {
	buff := make([]byte, 4)
	_, err := io.ReadFull(conn, buff)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint32(buff), nil
}

func ReadByte(conn net.Conn) (byte, error) {
	buff := make([]byte, 1)
	_, err := io.ReadFull(conn, buff)
	if err != nil {
		return 0, err
	}

	return buff[0], nil
}

func ReadString(conn net.Conn) (string, error) {
	l, err := ReadUint32(conn)
	if err != nil {
		return "", err
	}

	buff := make([]byte, l)
	_, err = io.ReadFull(conn, buff)
	if err != nil {
		return "", err
	}

	return string(buff), nil
}

func NewPacket(prefix byte, bytes []byte) []byte {
	return append([]byte{prefix}, bytes...)
}
