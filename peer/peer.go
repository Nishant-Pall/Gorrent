package peer

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
)

type Peer struct {
	IP   net.IP
	Port uint16
}

func (p *Peer) String() string {
	return net.JoinHostPort(string(p.IP.String()), strconv.Itoa(int(p.Port)))
}

func Unmarshal(peersBlob []byte) ([]Peer, error) {
	const peerSize = 6
	numPeers := len(peersBlob) / peerSize

	if len(peersBlob)%numPeers != 0 {
		return nil, fmt.Errorf("Received malformed peers")
	}

	peers := make([]Peer, numPeers)

	for i := range numPeers {
		offset := i * peerSize
		peers[i].IP = net.IP(peersBlob[offset : offset+4])
		peers[i].Port = binary.BigEndian.Uint16(peersBlob[offset+4 : offset+6])
	}
	return peers, nil
}
