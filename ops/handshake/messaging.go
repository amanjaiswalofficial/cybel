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

	data := msg.serializeMessage(keepAlive)

	_, err := conn.Write(data)
	if err != nil {
		fmt.Print(err)
	}

	buffer := make([]byte, 100)
	_, errorVal := io.ReadFull(conn, buffer)
	if errorVal != nil {
		fmt.Print(errorVal)
	}

	fmt.Println(buffer)
	fmt.Println(string(buffer))
	fmt.Println()
}
