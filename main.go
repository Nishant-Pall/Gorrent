package main

import (
	"os"
)

func main() {

	path := "./fixture/.torrent"
	file, err := os.Open(path)

	if err != nil {
		return
	}

	bto := &BencodeTorrent{}
	bto.OpenTorrent(file)

	bto.toTorrentFile()
	// fmt.Printf("%v \r\n", bto.Announce)

}
