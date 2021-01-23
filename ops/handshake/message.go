package handshake

const (
	keepAlive = "keep-alive"
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
	message[keepAlive] = &attrib{prefix: []byte{0,0,0,0}, messageID: 1}

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