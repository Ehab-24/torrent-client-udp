package main

import (
	"log"
	"net/url"
	"os"
	"runtime"
	"strconv"

	"github.com/Ehab-24/torrent-udp/node"
	"github.com/Ehab-24/torrent-udp/torrent"
)

func main() {

	// Path to .torrent file must be provided in the os args
	if len(os.Args) != 2 {
		log.Fatal("missing required argument \"path\"")
	}
	path := os.Args[1]

	// Read and parse the torrent file
	torrent, err := torrent.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	announceURL, err := url.Parse(torrent.Announce)
	if err != nil {
		log.Fatal(err)
	}
	port, err := strconv.Atoi(announceURL.Port())
	if err != nil {
		log.Fatal(err)
	}

	// Create a new Node instance for this client
	node, err := node.New(announceURL.Hostname(), uint16(port))
	if err != nil {
		log.Fatal(err)
	}
	defer node.Close()

	// Send connection message to the tracker
	if err := node.SendConnect(); err != nil {
		log.Fatal(err)
	}

	for runtime.NumGoroutine() > 1 {
	}
}
