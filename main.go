package main

import (
	"fmt"
	"gorrent/torrentFile"
	"os"
)

func main() {

	path := "./fixture/.torrent"
	file, err := os.Open(path)

	if err != nil {
		return
	}

	tf, err := torrentFile.Open(file)
	if err != nil {
		fmt.Printf("%v \r\n", err)
		return
	}

	err = tf.DownloadToFile("output.torrent")
	if err != nil {
		fmt.Printf("%v \r\n", err)
		return
	}
}
