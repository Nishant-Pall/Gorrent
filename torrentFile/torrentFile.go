package torrentFile

import (
	"crypto/rand"
	"fmt"
	"gorrent/p2p"
	"gorrent/peer"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/Nishant-Pall/bengoder"
)

const Port uint16 = 6881

type bencodeTrackerResponse struct {
	peers    string
	interval int
}

type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
	PeerID      [20]byte
}

func (t *TorrentFile) generatePeerId() ([20]byte, error) {
	var peerId [20]byte
	_, err := rand.Read(peerId[:])

	if err != nil {
		return [20]byte{}, err
	}

	t.PeerID = peerId

	return peerId, nil
}

func (t *TorrentFile) BuildTrackerURL() (string, error) {
	base, err := url.Parse(t.Announce)

	if err != nil {
		return "", err
	}

	peerId, err := t.generatePeerId()

	if err != nil {
		return "", err
	}

	params := url.Values{
		"info_hash":  []string{string(t.InfoHash[:])},
		"peer_id":    []string{string(peerId[:])},
		"port":       []string{strconv.Itoa(int(Port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(t.Length)},
		"event":      []string{"started"},
	}

	base.RawQuery = params.Encode()

	return base.String(), nil
}

func (t *TorrentFile) RequestPeers() ([]peer.Peer, error) {
	client := &http.Client{Timeout: 15 * time.Second}

	url, err := t.BuildTrackerURL()
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	rawResponse, err := bengoder.UnMarshall(resp.Body)
	if err != nil {
		return nil, err
	}

	if _, exists := rawResponse["failure reason"]; exists {
		return nil, fmt.Errorf("Failure: %v \r\n", rawResponse["failure reason"])
	}

	peers, err := peer.Unmarshal([]byte(rawResponse["peers"].(string)))
	if err != nil {
		return nil, err
	}

	return peers, nil
}

func Open(file io.Reader) (*TorrentFile, error) {

	bto := &BencodeTorrent{}

	err := bto.OpenTorrent(file)
	if err != nil {
		return nil, err
	}

	t, err := bto.ToTorrentFile()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (tf *TorrentFile) DownloadToFile(filePath string) error {

	peers, err := tf.RequestPeers()

	if err != nil {
		return err
	}

	torrent := p2p.Torrent{
		Peers:       peers,
		PeerID:      tf.PeerID,
		InfoHash:    tf.InfoHash,
		PieceHashes: tf.PieceHashes,
		PieceLength: tf.PieceLength,
		Length:      tf.Length,
		Name:        tf.Name,
	}

	buf, err := torrent.Download()

	if err != nil {
		return err
	}

	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = outFile.Write(buf)
	if err != nil {
		return err
	}

	return nil
}
