package main

import (
	"flag"
	"math/rand"
	"time"

	"github.com/dacousb/feiok/game"
	"github.com/dacousb/feiok/server"
)

func main() {
	start := flag.Bool("server", false, "start a new server and listen on a port for packets")
	host := flag.String("host", "", "launch the game and connect to the desired host")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	if *start {
		server.New().Run()
		return
	}

	if *host == "" {
		flag.PrintDefaults()
		return
	}

	game.New().Run(*host)
}
