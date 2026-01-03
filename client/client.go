package client

import (
	"fmt"
	"gorrent/handshake"
	"gorrent/message"
	"gorrent/peer"
	"net"
	"time"
)

type Client struct {
	Conn     net.Conn
	Choked   bool
	Bitfield message.BitField
	peer     peer.Peer
	infoHash [20]byte
	peerID   [20]byte
}

func CompleteHandShake(conn net.Conn, peerID, infoHash [20]byte) (*handshake.HandShake, error) {
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetDeadline(time.Time{})

	hs := handshake.New(infoHash, peerID)

	conn.Write(hs.Serialize())

	respHandShake, err := hs.Read(conn)
	if err != nil {
		return nil, err
	}

	ok := respHandShake.ValidateResponse(infoHash)
	if !ok {
		return nil, fmt.Errorf("%v \r\n", "Response handshake validation failed")
	}

	return hs, nil
}

func receiveBitField(conn net.Conn) (message.BitField, error) {
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetDeadline(time.Time{})

	msg, err := message.Read(conn)

	if err != nil {
		return nil, err
	}

	if msg == nil {
		return nil, fmt.Errorf("Expected bitfield but got: %v", msg)
	}

	if msg.ID != message.MsgBitfield {
		return nil, fmt.Errorf("Expected bitfield but got ID: %d", msg.ID)
	}

	return msg.Payload, nil
}

func New(peer peer.Peer, peerID, infoHash [20]byte) (*Client, error) {

	conn, err := net.DialTimeout("tcp", peer.String(), 3*time.Second)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	_, err = CompleteHandShake(conn, peerID, infoHash)

	if err != nil {
		conn.Close()
		return nil, err
	}

	bf, err := receiveBitField(conn)

	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Client{
		Conn:     conn,
		Choked:   true,
		Bitfield: bf,
		peer:     peer,
		infoHash: infoHash,
		peerID:   peerID,
	}, nil
}

func (c *Client) SendInterested() error {
	msg := message.Message{ID: message.MsgInterested}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}
