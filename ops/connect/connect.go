package connect

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"cybele/ops/utils"
	"cybele/ops/bencode"
)

// Object to store tracker related information
type trackerRequest struct {
	url      *url.URL
	respBody *string
}

// FetchDetailsFromTorrent is used to read json from a file
// Then fetch the list of peers by making an HTTP request
// To the tracker
// returns: none
func FetchDetailsFromTorrent(path string) {
	bs, err := utils.ReadFileFromPath(path)
	if err != nil {
		utils.HandleError(err.Error())
	}

	td := ReadJSONFromByteSlice(bs)
	if td.IsEmpty() {
		utils.HandleError(utils.ErrorReadingJSON)
	}

	connectToTracker(td)
}

// Connect to tracker and retrieve list of peers
func connectToTracker(td TorrentData) {

	tReqObj := makeRequestObject(td)
	resp, err := getResponse(tReqObj)

	if err != nil {
		utils.HandleError(utils.ErrorConnectingToTracker)
	}

	tReqObj.respBody = &resp // Using address, as defer deletes the original value
	trackerResponse := strings.NewReader(*tReqObj.respBody)
	decodedResponse, errs := bencode.Decode(trackerResponse)
	if errs != nil {
		utils.LogMessage(errs.Error())
	}

	// TO UPDATE
	fmt.Print(decodedResponse)


}

// Make request object of type trackerRequest from TorrentData
func makeRequestObject(td TorrentData) trackerRequest {
	var tr trackerRequest
	u, err := url.Parse(td.Announce)
	if err != nil {
		utils.HandleError(utils.ErrorParsingAnnounceURL)
	}
	tr.url = u
	tr.addParamsToTrackerRequest(td)
	return tr

}

// Add required params to the url to request to tracker
func (tr trackerRequest) addParamsToTrackerRequest(td TorrentData) {

	infoHash := utils.MakeInfoHash(td.InfoHash)
	params := url.Values{
		"peer_id":    []string{string("-AA1111-123456789012")}, // to change
		"port":       []string{strconv.Itoa(int(utils.ConnectionPort))},
		"uploaded":   []string{"0"},     // by default, for first request
		"downloaded": []string{"0"},     // by default, for first request
		"left":       []string{td.Size}, // to confirm
		//"compact":    []string{"1"},     // by default, for the first request
	}

	tr.url.RawQuery = params.Encode()
	// Adding info_hash separately to avoid url-encoding and keeping hex-encode
	tr.url.RawQuery = tr.url.RawQuery + "&info_hash=" + infoHash
}

// Get response for the GET request to tracker URL with required params
func getResponse(tr trackerRequest) (string, error) {

	utils.LogMessage("Connecting:", tr.url.String(), "\n")
	resp, err := http.Get(tr.url.String()) // Make GET Request to tracker for response
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body) // Decode to get response body
	if err != nil {
		return "", err
	}

	defer resp.Body.Close() // Close the what?

	return string(body), nil
}
