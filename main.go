package main

import (
	"flag"
	"fmt"
	"gorrent/torrentFile"
	"os"
)

func main() {

	pathPtr := flag.String("inpath", "", "Path to torrent file")
	outputPathPtr := flag.String("outpath", "", "Path/Name to output file")
	flag.Parse()

	if *pathPtr == "" {
		fmt.Printf("Path (-inpath) to torrent file cannot be empty")
		return
	}

	if *outputPathPtr == "" {
		fmt.Printf("Path/name (-outpath) of destination file cannot be empty")
		return
	}

	file, err := os.Open(*pathPtr)

	if err != nil {
		fmt.Println(err)
		return
	}

	tf, err := torrentFile.Open(file)
	if err != nil {
		fmt.Printf("%v \r\n", err)
		return
	}

	err = tf.DownloadToFile(*outputPathPtr)
	if err != nil {
		fmt.Printf("%v \r\n", err)
		return
	}
}
