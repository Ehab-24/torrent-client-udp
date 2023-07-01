package torrent

import (
	"os"

	"github.com/jackpal/bencode-go"
)

type file struct {
	Length int `bencode:"length"`
	Path   int `bencode:"path"`
}

type Info struct {
	Name        string `bencode:"name"`
	Length      int    `bencode:"length"`
	PieceLength int    `bencode:"piece length"`
	Pieces      string `bencode:"pieces"`
	Files       []file `bencode:"files"`
}

type Torrent struct {
	Announce string `bencode:"announce"`
	Info     Info   `bencode:"Info"`
}

func Open(path string) (*Torrent, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	// Parse the file from bencode
	var torrent Torrent
	if err := bencode.Unmarshal(file, &torrent); err != nil {
		return nil, err
	}

	return &torrent, nil
}
