package connect

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"cybele/ops/utils"
)

// Object to store tracker related information
type trackerRequest struct {
	url *url.URL 
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
	utils.LogMessage("Response:",*tReqObj.respBody)
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
	params := url.Values{
		"info_hash":  []string{string(td.InfoHash)},
		"peer_id":    []string{string("-AA1111-123456789012")}, // to change
		"port":       []string{strconv.Itoa(int(utils.ConnectionPort))},
		"uploaded":   []string{"0"},     // by default, for first request
		"downloaded": []string{"0"},     // by default, for first request
		"left":       []string{td.Size}, // to confirm
		"compact":    []string{"1"},     // by default, for the first request
	}
	tr.url.RawQuery = params.Encode()
}

// Get response for the GET request to tracker URL with required params
func getResponse(tr trackerRequest) (string, error){

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
