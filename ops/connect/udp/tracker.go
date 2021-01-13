package udp

import (
	"cybele/ops/connect"
	"cybele/ops/utils"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math/rand"
	"net"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// UDPTracker is a torrent tracker that speaks udp
type UDPTracker struct {
	dest        string // Destination host
	requestBuf  []byte
	responseBuf []byte
	timeout     time.Duration
	conn        net.Conn // Underlying network connection
	retries     uint8    // retransmission retries threshold
	tid         uint32   // Transaction ID
	cid         uint64   // Connection ID
}

func New(host string) *UDPTracker {
	return &UDPTracker{
		dest:        host,
		requestBuf:  make([]byte, 100),
		responseBuf: make([]byte, 100),
		retries:     uint8(8),
	}
}

func generateTid() uint32 {
	return uint32(rand.Int31())
}

// connect dials up to an udp host.
// returns: A network connection and a 64-bit integer representing
// the connection ID.
func (tr *UDPTracker) connect() (net.Conn, uint64, error) {
	tr.tid = generateTid()
	// Prepare the udp packet
	binary.BigEndian.PutUint64(tr.requestBuf[0:], utils.Pid)
	binary.BigEndian.PutUint32(tr.requestBuf[8:], utils.Connect)
	binary.BigEndian.PutUint32(tr.requestBuf[12:], tr.tid)

	// Timeout Formula: 15 * 2 ^ n (n is the number of retries, starts at 0 and up to 8)
	n := uint8(0)
	tr.timeout = 15 * time.Second
	conn, err := net.DialTimeout("udp", tr.dest, tr.timeout)
	if err != nil {
		return nil, 0, err
	}

	for {
		tr.timeout = time.Duration(15*(2^n)) * time.Second
		n++

		conn.SetWriteDeadline(time.Now().Add(tr.timeout))
		nbytes, err := conn.Write(tr.requestBuf)
		if err != nil {
			return conn, 0, err
		} else if nbytes != len(tr.requestBuf) {
			return conn, 0, errors.New("must send 16 bytes")
		}

		conn.SetReadDeadline(time.Now().Add(tr.timeout))
		nbytes, err = conn.Read(tr.responseBuf)
		if err, ok := err.(net.Error); ok && err.Timeout() {
			if n > tr.retries {
				return conn, 0, err
			}

			// retry the transmission
			continue
		} else if err != nil {
			return conn, 0, err
		} else if nbytes < 16 {
			return conn, 0, errors.New("must read 16 bytes")
		}
		break
	}

	cid := binary.BigEndian.Uint64(tr.responseBuf[8:])
	return conn, cid, nil
}

// Announce announces a torrent to a udp tracker
// returns: an AnnounceResponse struct containing the peers
// and other relevant information (e.g. seeders, leechers, etc.).
func (tr *UDPTracker) Announce(r *connect.AnnounceRequest) (*connect.AnnounceResponse, error) {
	// Get the network connection and connection id
	conn, cid, err := tr.connect()
	if err != nil {
		return nil, err
	}

	tid := generateTid()

	// Prepare the udp packet
	binary.BigEndian.PutUint64(tr.requestBuf[0:], cid)
	binary.BigEndian.PutUint32(tr.requestBuf[8:], utils.Announce)
	binary.BigEndian.PutUint32(tr.requestBuf[12:], tid)

	// Copy info_hash and peerId
	copy(tr.requestBuf[16:], r.InfoHash)
	copy(tr.requestBuf[36:], r.PeerID)

	binary.BigEndian.PutUint64(tr.requestBuf[56:], r.Downloaded)
	binary.BigEndian.PutUint64(tr.requestBuf[64:], r.Left)
	binary.BigEndian.PutUint64(tr.requestBuf[72:], r.Uploaded)
	binary.BigEndian.PutUint32(tr.requestBuf[92:], utils.MaxPeers)

	n := uint8(0)
	for {
		tr.timeout = time.Duration(15*(2^n)) * time.Second
		n++

		conn.SetWriteDeadline(time.Now().Add(tr.timeout))
		_, err := conn.Write(tr.requestBuf)
		if err != nil {
			return nil, err
		}

		conn.SetReadDeadline(time.Now().Add(tr.timeout))
		nbytes, err := conn.Read(tr.responseBuf)

		if err, ok := err.(net.Error); ok && err.Timeout() {
			if n > tr.retries {
				return nil, err
			}

			// retry transmission
			continue
		} else if nbytes < 20 {
			return nil, errors.New("must read at least 20 bytes")
		} else if err != nil {
			return nil, err
		}
		break
	}

	interval := time.Duration(binary.BigEndian.Uint32(tr.responseBuf[8:])) * time.Second

	resp := &connect.AnnounceResponse{
		Complete:   binary.BigEndian.Uint32(tr.responseBuf[16:]),
		Incomplete: binary.BigEndian.Uint32(tr.responseBuf[12:]),
		Interval:   interval,
	}

	peersCount := len(tr.responseBuf[20:])
	peers := make([]connect.PeerObject, 0, peersCount)
	offset := 20

	for offset < peersCount {
		ip := make(net.IP, 4)
		ipInt := binary.BigEndian.Uint32(tr.responseBuf[offset:])
		binary.BigEndian.PutUint32(ip, ipInt)
		port := binary.BigEndian.Uint16(tr.responseBuf[offset+4:])

		offset += 6
		peerObj := connect.PeerObject{
			IP:   ip.String(),
			Port: strconv.Itoa(int(port)),
		}
		peers = append(peers, peerObj)
	}

	resp.Peers = peers
	return resp, nil
}

func MakeRequestObject(td *connect.TorrentData) *connect.AnnounceRequest {
	// convert the info_hash back to binary
	digest, err := hex.DecodeString(td.InfoHash)
	if err != nil {
		utils.HandleError(err.Error())
	}

	peerID := make([]byte, 20)
	rand.Read(peerID)

	size, err := strconv.ParseUint(td.Size, 10, 64)
	if err != nil {
		utils.HandleError(err.Error())
	}

	req := &connect.AnnounceRequest{
		InfoHash:   digest[:],
		PeerID:     peerID[:],
		Uploaded:   uint64(0),
		Downloaded: uint64(0),
		Left:       size,
		Port:       uint16(utils.ConnectionPort),
	}
	return req
}
