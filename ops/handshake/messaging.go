package handshake

import (
	"fmt"
	"io"
	"net"
)

// KeepTalkingToPeer is a method to keep the messaging going to peers on the 
// tcp once created connection with them
func KeepTalkingToPeer(conn net.Conn) {

	msg := makeMessage()

	data := msg.serializeMessage(unchoke)

	_, err := conn.Write(data)
	if err != nil {
		fmt.Print(err)
	}

	data = msg.serializeMessage(interested)

	_, errors := conn.Write(data)
	if errors != nil {
		fmt.Println(err)
	}

	buffer := make([]byte, 100)
	_, errorVal := io.ReadFull(conn, buffer)
	if errorVal != nil {
		fmt.Print(errorVal)
	}

	fmt.Println(string(buffer))
}
