package handshake

import "encoding/binary"

type messageID int8

const (
	keepAlive    messageID = 0
	choke        messageID = 0
	unchoke      messageID = 1
	interested   messageID = 2
	notIntersted messageID = 3
	have         messageID = 4
	bitField     messageID = 5
	request      messageID = 6
	piece        messageID = 7
	cancel       messageID = 8
)

type Message struct {
	MsgID   messageID
	Payload []byte
}

// convert the message into a byteslice to be sent to peers
// based on the key provided
func (msg Message) serializeMessage(prefixLen messageID) []byte {

	data := make([]byte, prefixLen+4)

	binary.BigEndian.PutUint32(data[0:4], uint32(prefixLen))
	data[4] = byte(prefixLen)

	return data
}
