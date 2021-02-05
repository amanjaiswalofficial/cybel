package handshake

import (
	"cybele/ops/utils"
	"fmt"
	"net"
)

// KeepTalkingToPeer is a method to keep the messaging going to peers on the
// tcp once created connection with them
func KeepTalkingToPeer(conn net.Conn) {

	//clean the stream by reading the received bitfield message
	buffer := utils.ReadData(conn, 4)
	remainderStreamLength := int(buffer[3] + 5) // the later 5 bytes are another empty bitfield msg
	utils.ReadData(conn, remainderStreamLength)

	var msg Message
	
	// sending unchoked and interested messages
	data := msg.serializeMessage(unchoke, nil)
	utils.WriteData(conn, data)
	
	data = msg.serializeMessage(interested, nil)
	utils.WriteData(conn, data)

	dataNew := utils.ReadData(conn, 5)
	fmt.Println(dataNew)
	conn.Close()
}
