package utils

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
)
