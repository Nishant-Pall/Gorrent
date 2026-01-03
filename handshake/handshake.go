package handshake

import (
	"bytes"
	"fmt"
	"io"
)

type HandShake struct {
	Pstr     string
	InfoHash [20]byte
	PeerID   [20]byte
}

const PSTR = "BitTorrent protocol"

func New(infoHash, peerId [20]byte) *HandShake {
	return &HandShake{
		Pstr:     PSTR,
		InfoHash: infoHash,
		PeerID:   peerId,
	}
}

func (h *HandShake) Serialize() []byte {
	buf := make([]byte, len(h.Pstr)+49)

	buf[0] = byte(len(h.Pstr))

	curr := 1
	curr += copy(buf[curr:], h.Pstr)
	curr += copy(buf[curr:], make([]byte, 8))
	curr += copy(buf[curr:], h.InfoHash[:])
	curr += copy(buf[curr:], h.PeerID[:])

	return buf
}

func (h *HandShake) Read(r io.Reader) (*HandShake, error) {

	lengthBuf := make([]byte, 1)

	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return nil, err
	}

	pstrlen := int(lengthBuf[0])
	if pstrlen == 0 {
		return nil, fmt.Errorf("%v \r\n", "pstrlen cannot be 0")
	}

	handshakeBuf := make([]byte, 48+pstrlen)

	_, err = io.ReadFull(r, handshakeBuf)
	if err != nil {
		return nil, err
	}

	var peerID, infoHash [20]byte
	copy(infoHash[:], handshakeBuf[pstrlen+8:pstrlen+8+20])
	copy(peerID[:], handshakeBuf[pstrlen+20+8:])

	return &HandShake{
		Pstr:     string(handshakeBuf[:pstrlen]),
		InfoHash: infoHash,
		PeerID:   peerID,
	}, nil

}

func (h *HandShake) ValidateResponse(infoHash [20]byte) bool {
	return bytes.Equal(h.InfoHash[:], infoHash[:]) && h.Pstr == PSTR
}
