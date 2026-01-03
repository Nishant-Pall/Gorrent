package message

import (
	"encoding/binary"
	"io"
)

type BitField []byte

func (bf BitField) HasPiece(index int) bool {
	byteIndex := index / 8
	offset := index % 8

	return bf[byteIndex]>>(7-offset)&1 != 0
}

type messageID uint8

const (
	MsgChoke         messageID = 0
	MsgUnchoke       messageID = 1
	MsgInterested    messageID = 2
	MsgNotInterested messageID = 3
	MsgHave          messageID = 4
	MsgBitfield      messageID = 5
	MsgRequest       messageID = 6
	MsgPiece         messageID = 7
	MsgCancel        messageID = 8
)

type Message struct {
	ID      messageID
	Payload BitField
}

func FormatRequest(index, offset, length int) *Message {
	payload := make([]byte, 12)

	binary.BigEndian.PutUint32(payload[0:4], uint32(index))
	binary.BigEndian.PutUint32(payload[4:8], uint32(offset))
	binary.BigEndian.PutUint32(payload[8:12], uint32(length))
	return &Message{ID: MsgRequest, Payload: payload}
}

func (m *Message) Serialize() []byte {
	if m == nil {
		return make([]byte, 4)
	}

	length := uint32(len(m.Payload) + 1)

	buf := make([]byte, length+4)
	binary.BigEndian.PutUint32(buf[0:4], length)
	buf[4] = byte(m.ID)
	copy(buf[5:], m.Payload)

	return buf
}

func Read(r io.Reader) (*Message, error) {
	lengthBuff := make([]byte, 4)

	_, err := io.ReadFull(r, lengthBuff)
	if err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(lengthBuff)

	messageBuf := make([]byte, length)
	_, err = io.ReadFull(r, messageBuf)
	if err != nil {
		return nil, err
	}

	return &Message{
		ID:      messageID(messageBuf[0]),
		Payload: messageBuf[1:],
	}, nil
}

func ParsePiece(index int, buf []byte, msg *Message) int {
	parsedIndex := int(binary.BigEndian.Uint32(msg.Payload[0:4]))

	if parsedIndex != index {
		return 0
	}

	begin := int(binary.BigEndian.Uint32(msg.Payload[4:8]))
	data := msg.Payload[8:]
	copy(buf[begin:], data)

	return len(data)
}
