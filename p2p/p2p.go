package p2p

import (
	"gorrent/client"
	"gorrent/message"
	"gorrent/peer"
	"log"
)

type Torrent struct {
	Peers       []peer.Peer
	PeerID      [20]byte
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

func (t *Torrent) Download() error {
	client, err := client.New(t.Peers[0], t.PeerID, t.InfoHash)

	if err != nil {
		return err
	}
	defer client.Conn.Close()
	log.Printf("Completed handshake with %s\n", t.Peers[0].IP)

	client.SendUnchoke()
	message.Read(client.Conn)

	client.SendInterested()
	message.Read(client.Conn)

	return nil
}
