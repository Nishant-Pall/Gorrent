package main

import (
	"fmt"
	"os"
)

func main() {

	path := "./fixture/.torrent"
	file, err := os.Open(path)

	if err != nil {
		return
	}

	bto, err := OpenTorrent(file)

	if err != nil {
		return
	}

	fmt.Printf("%v \r\n", bto)

}

func buildTrackerURL() {

}
