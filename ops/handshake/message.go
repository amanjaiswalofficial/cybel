package handshake

const (
	keepAlive = "keep-alive"
	choke = "choke"
	unchoke = "unchoke"
	interested = "interested"
	notIntersted = "not-interested"
	have = "have"
	request = "request"
	cancel = "cancel"
	piece = "piece"
)

type attrib struct {
	prefix []byte
	messageID byte
	payload []byte
}

// Message is a basic message type for sending messages 
// via tcp to peers requesting data
type Message map[string]*attrib

// return a Message map, containing attribs structs as value
// with different attrib for different key constants
func makeMessage() Message {

	message := make(Message)
	message[keepAlive] = &attrib{prefix: []byte{0,0,0,0}, messageID: 0}
	message[choke] = &attrib{prefix: []byte{0,0,0,1}, messageID: 0}
	message[unchoke] = &attrib{prefix: []byte{0,0,0,1}, messageID: 1}
	message[interested] = &attrib{prefix: []byte{0,0,0,1}, messageID: 2}
	message[notIntersted] = &attrib{prefix: []byte{0,0,0,1}, messageID: 3}
	message[have] = &attrib{prefix: []byte{0,0,0,5}, messageID: 5}
	message[request] = &attrib{prefix: []byte{0,0,1,3}, messageID:6}
	message[piece] = &attrib{prefix: []byte{0,0,0,9}, messageID:7}
	message[cancel] = &attrib{prefix: []byte{0,0,1,3}, messageID:8}

	return message

}

// convert the message into a byteslice to be sent to peers
// based on the key provided
func (msg Message) serializeMessage(key string) []byte {
	data := make([]byte, 10)
	copy(data[0:4], msg[key].prefix)
	data[4] = msg[key].messageID
	
	return data
}