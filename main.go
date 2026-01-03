package main

import (
	"fmt"
	"gorrent/client"
	"gorrent/torrentFile"
	"io"
	"log"
	"os"
)

func main() {

	path := "./fixture/.torrent"
	file, err := os.Open(path)

	if err != nil {
		return
	}

	downloadFiles(file)
}

func downloadFiles(file io.Reader) {
	bto := &torrentFile.BencodeTorrent{}
	bto.OpenTorrent(file)

	t, _ := bto.ToTorrentFile()
	url, _ := t.BuildTrackerURL()

	peers, err := t.RequestPeers(url)

	if err != nil {
		fmt.Printf("%v \r\n", err)
		return
	}

	client, err := client.New(peers[0], t.PeerID, t.InfoHash)

	if err != nil {
		fmt.Printf("%v \r\n", err)
		return
	}
	defer client.Conn.Close()
	log.Printf("Completed handshake with %s\n", peers[0].IP)

	client.SendInterested()
}
