package handshake

import (
	"bytes"
	"errors"
	"io"
	"net"
	"strings"
	"time"

	"cybele/ops/utils"
	"cybele/ops/connect"
)

type Handshake struct {
	Pstr     string
	InfoHash []byte
	PeerID   []byte
}

// GetString returns the byte-slice for the handshake string
func (Handshake) GetString(infoHash []byte, peerID []byte) []byte {

	hs := &Handshake{
		Pstr:     "BitTorrent protocol",
		InfoHash: infoHash,
		PeerID:   peerID,
	}

	hsStr := hs.Serialize()
	return hsStr

}

// Serialize normalises the Handshake struct to make a byte-slice
func (hs *Handshake) Serialize() []byte {
	buf := make([]byte, len(hs.Pstr)+49)
	buf[0] = byte(len(hs.Pstr))
	offset := 1
	offset += copy(buf[offset:], hs.Pstr)
	offset += copy(buf[offset:], make([]byte, 8)) // 8 reserved bytes
	offset += copy(buf[offset:], hs.InfoHash[:])
	offset += copy(buf[offset:], hs.PeerID[:])

	return buf
}

// DoHandshake is used to connect with several peers for a specific file
// depending on the infoHash (extracted from torrent file)
func DoHandshake(hsStr []byte, infoHash []byte, peers []connect.PeerObject) {

	for i := 0; i < len(peers); i++ {

		address := strings.Join([]string{peers[i].IP, peers[i].Port}, ":")
		conn, err := net.DialTimeout("tcp", address, 3*time.Second)
		if err != nil {
			continue
		}

		_, err = conn.Write(hsStr)
		if err != nil {
			conn.Close()
			continue
		}

		res, err := ReadPostHandshake(conn)
		if err != nil {
			conn.Close()
			continue
		}

		if !bytes.Equal(res.InfoHash[:], infoHash[:20]) {
			conn.Close()
			continue
		}

		KeepTalkingToPeer(conn)
	}

}

// ReadPostHandshake is used to read data from a connection object in
// required format, and return a handshake struct from the response
func ReadPostHandshake(r io.Reader) (*Handshake, error) {
	length := make([]byte, 1)
	_, err := io.ReadFull(r, length)
	if err != nil {
		return nil, err
	}

	pstrlen := int(length[0])

	if pstrlen == 0 {
		return nil, errors.New(utils.ZeroLengthError)
	}

	buf := make([]byte, pstrlen+48)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}

	infoHash, peerID := make([]byte, 20), make([]byte, 20)

	copy(infoHash[:], buf[pstrlen+8:pstrlen+8+20])
	copy(peerID[:], buf[pstrlen+8+20:])

	h := &Handshake{
		Pstr:     string(buf[0:pstrlen]),
		InfoHash: infoHash,
		PeerID:   peerID,
	}

	return h, nil
}
