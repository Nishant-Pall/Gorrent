package main

import (
	"crypto/sha1"
	"fmt"
	"io"

	"github.com/Nishant-Pall/bengoder"
)

type BencodeInfo struct {
	Length      int
	Name        string
	PieceLength int
	Pieces      string
}

func (bi *BencodeInfo) hash() [20]byte {
	rawInfo := bengoder.Encode(bi)

	h := sha1.Sum(rawInfo)

	return h
}

type BencodeTorrent struct {
	Announce     string
	Comment      string
	CreatedBy    string
	CreationDate int
	Info         BencodeInfo
	infoHash     [20]byte
}

func (bto *BencodeTorrent) OpenTorrent(r io.Reader) error {

	decodedTorrent := bengoder.UnMarshallFile(r)

	bencodeInfo := &BencodeInfo{}

	info, ok := decodedTorrent["info"].(map[string]any)

	if !ok {
		return fmt.Errorf("Error parsing info dict")
	}

	bto.Announce = decodedTorrent["announce"].(string)
	bto.Comment = decodedTorrent["comment"].(string)
	bto.CreatedBy = decodedTorrent["created by"].(string)
	bto.CreationDate = decodedTorrent["creation date"].(int)

	bencodeInfo.Length = info["length"].(int)
	bencodeInfo.Name = info["name"].(string)
	bencodeInfo.PieceLength = info["piece length"].(int)
	bencodeInfo.Pieces = info["pieces"].(string)

	bto.Info = *bencodeInfo

	return nil
}

func (bto *BencodeTorrent) toTorrentFile() (TorrentFile, error) {

	infoHash := bto.Info.hash()

	t := TorrentFile{
		Announce:    bto.Announce,
		InfoHash:    infoHash,
		PieceLength: bto.Info.PieceLength,
		Length:      bto.Info.Length,
		Name:        bto.Info.Name,
	}

	return t, nil
}
