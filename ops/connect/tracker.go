package connect

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"cybele/ops/bencode"
	"cybele/ops/utils"
)

// Tracker represents a udp or http tracker
// TODO: Implement http tracker with this interface
type Tracker interface {
	// Announce connects to (udp or http) tracker
	Announce(r *AnnounceRequest) (*AnnounceResponse, error)
}

type AnnounceRequest struct {
	InfoHash   []byte
	PeerID     []byte
	Port       uint16
	Uploaded   uint64
	Downloaded uint64
	Left       uint64
	Compact    uint8
	TrackerID  string // This is used for future announcements
}

type AnnounceResponse struct {
	Complete    uint32
	Incomplete  uint32
	Interval    time.Duration
	MinInterval time.Duration
	Peers       []PeerObject
	TrackerID   string // Returned by the tracker for future announcements
}

// PeerObject contains structs having IP, Port, PeerID from
type PeerObject struct {
	IP     string
	Port   string
	PeerID string
}

// Object to store tracker related information
type trackerRequest struct {
	url         *url.URL
	response    *string
	DecodedResp struct {
		Complete    int64
		Incomplete  int64
		Interval    int64
		MinInterval int64
		Peers       []PeerObject
	}
}

// Add required params to the url to request to tracker
func (tr trackerRequest) addParamsToTrackerRequest(td TorrentData) {

	infoHash := utils.MakeInfoHash(td.InfoHash)
	params := url.Values{
		"peer_id":    []string{utils.MakePeerID()}, // to change
		"port":       []string{strconv.Itoa(int(utils.ConnectionPort))},
		"uploaded":   []string{"0"},                             // by default, for first request
		"downloaded": []string{"0"},                             // by default, for first request
		"left":       []string{fmt.Sprintf("%d", td.TotalSize)}, // to confirm
	}

	tr.url.RawQuery = params.Encode()
	// Adding info_hash separately to avoid url-encoding and keeping hex-encode
	tr.url.RawQuery = tr.url.RawQuery + "&info_hash=" + infoHash
}

// decodeResponse() is used to decode values from response received from tracker
// it uses bencoding to convert values into human readable format
// returns: trackerRequest struct with updated values for decodedResp
// 			or error if it exists
func (tr *trackerRequest) decodeResponse() (err error) {
	trackerResponse := strings.NewReader(*tr.response)
	decodedResponse, dErr := bencode.Decode(trackerResponse)
	if dErr != nil {
		return dErr
	}

	// TODO: Add error handling for failure (requires discussion)
	// for each type of key, value pair in response
	for key, val := range decodedResponse {
		switch val.(type) {

		case int64:
			/*
				For each of the key with int64 type
				Find the respective struct variable (using FieldByName)
				Then by using SetInt(), set value for the same in tr.decodedResp
			*/
			rresp := reflect.ValueOf(&tr.DecodedResp)
			resStruct := rresp.Elem()
			formattedKey := utils.FormatKey(key)
			targetField := resStruct.FieldByName(strings.Title(formattedKey))
			targetField.SetInt(val.(int64))

		case []interface{}:
			fetchedValues := val.([]interface{})
			var p PeerObject
			for _, values := range fetchedValues {
				assertedFetchedVal := values.(map[string]interface{})
				for key, fetchedVal := range assertedFetchedVal {

					/*
						For each of the key in the peers response from tracker
						i.e. Ip, Port and PeerID
						Dynamically, find FieldByName
						And update the same in tr.decodedResp.peers by appending
						each peerObject
					*/
					rresp := reflect.ValueOf(&p)
					resStruct := rresp.Elem()
					formattedKey := utils.FormatKey(key)
					targetField :=
						resStruct.FieldByName(strings.Title(formattedKey))
					targetField.SetString(fetchedVal.(string))

					tr.DecodedResp.Peers =
						append(tr.DecodedResp.Peers, p)
				}
			}
		default:
			return errors.New(utils.UnknownDecodeKeysEncountered)
		}
	}
	return nil
}
