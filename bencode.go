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

func (bi *BencodeInfo) splitPieceHashes() ([][20]byte, error) {
	hashLen := 20
	buf := []byte(bi.Pieces)

	if len(buf)%hashLen != 0 {
		return nil, fmt.Errorf("Received malformed pieces of length %d", len(buf))
	}

	numHashes := len(buf) / hashLen
	hashes := make([][20]byte, numHashes)

	for i := range numHashes {
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])
	}

	return hashes, nil

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

	pieceHashes, err := bto.Info.splitPieceHashes()

	if err != nil {
		return TorrentFile{}, err
	}

	t := TorrentFile{
		Announce:    bto.Announce,
		InfoHash:    infoHash,
		PieceLength: bto.Info.PieceLength,
		PieceHashes: pieceHashes,
		Length:      bto.Info.Length,
		Name:        bto.Info.Name,
	}

	return t, nil
}
