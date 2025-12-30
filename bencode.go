package main

import (
	"fmt"
	"io"

	"github.com/Nishant-Pall/bengoder"
)

type BencodeTorrent struct {
	Announce     string
	Comment      string
	CreatedBy    string
	CreationDate int
	Info         BencodeInfo
}

type BencodeInfo struct {
	Length      int
	Name        string
	PieceLength int
	Pieces      string
}

func OpenTorrent(r io.Reader) (*BencodeTorrent, error) {

	reader := bengoder.NewResp(r)

	val, _ := reader.Decode()
	m, ok := val.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("Error parsing to map")
	}

	bencodeInfo := &BencodeInfo{}

	info := m["info"].(map[string]any)

	bencodeInfo.Length = info["length"].(int)
	bencodeInfo.Name = info["name"].(string)
	bencodeInfo.PieceLength = info["piece length"].(int)
	bencodeInfo.Pieces = info["pieces"].(string)

	bencodeTorrent := &BencodeTorrent{}

	bencodeTorrent.Announce = m["announce"].(string)
	bencodeTorrent.Comment = m["comment"].(string)
	bencodeTorrent.CreatedBy = m["created by"].(string)
	bencodeTorrent.CreationDate = m["creation date"].(int)
	bencodeTorrent.Info = *bencodeInfo

	return bencodeTorrent, nil
}
