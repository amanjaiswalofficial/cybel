package utils

import (
	"os"
	"strings"
)

const (
	// ErrorReadingJSON is displayed on failure
	ErrorReadingJSON = "Error Reading JSON from path"
	// ErrorMarshaling is used when unable to marshal torrent data to json
	ErrorMarshaling = "Error Marshaling data to JSON"
	// ConnectionPort is used as url params for tracker request
	ConnectionPort = 6889
	// ErrorParsingAnnounceURL is used when unable to parse url
	ErrorParsingAnnounceURL = "Error Parsing the announce URL for torrent"
	// ErrorConnectingToTracker is used when error connecting to tracker
	ErrorConnectingToTracker = "Error Connecting to Tracker"
	// ErrorDecodingResponse is used when response from tracker couldn't be decoded
	ErrorDecodingResponse = "Error encountered while decoding response from tracker"
	// UnknownDecodeKeysEncountered is used when keys from tracker response couldn't be handled
	UnknownDecodeKeysEncountered = "Error encountered while decoding keys from tracker response"
	// UDPTimeout is used when the number of retransmission surpass the threshold
	UDPTimeout = "Transmission timed out"
	// Protocol ID (magic constant used by the udp tracker)
	Pid = uint64(0x41727101980)
	// UDP Request Actions
	Connect  = uint32(0)
	Announce = uint32(1)
	// Max number of peers from response
	MaxPeers = uint32(10)
)

// CybeleCachePath is where all json files will reside
// as well as the files that are in the download queue.
var CybeleCachePath = strings.Join([]string{os.Getenv("HOME"), ".cache", "cybele"}, "/")
// QueueFileName is used to keep track of files in queue
var QueueFileName = "queue"
