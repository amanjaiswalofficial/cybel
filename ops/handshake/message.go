package handshake

import "encoding/binary"

type messageID int8
type payload []byte

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
	Payload payload
}

// convert the message into a byteslice to be sent to peers
// based on the key provided
func (msg Message) serializeMessage(prefixLen messageID, pLoad payload) []byte {

	msg.MsgID = prefixLen
	
	if pLoad != nil {
		pLoad = make([]byte, 0)
	}

	data := make([]byte, len(pLoad)+1+4)

	binary.BigEndian.PutUint32(data[0:4], uint32(len(pLoad)+1))
	data[4] = byte(msg.MsgID)

	return data
}
