package main

import (
	"gorrent/torrentFile"
	"io"
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
	tf, err := torrentFile.Open(file)

	if err != nil {
		return
	}

	err = tf.DownloadToFile("output.torrent")
	if err != nil {
		return
	}
}
