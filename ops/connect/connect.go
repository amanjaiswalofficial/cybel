package connect

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"cybele/ops/utils"
)

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
	tr := makeRequestObject(td)

	resp, err := getResponse(tr)

	if err != nil {
		utils.HandleError(utils.ErrorConnectingToTracker)
	}

	tr.response = &resp // Using address, as defer deletes the original value
	trPtr := &tr

	// Decoding the response for further use
	dErr := trPtr.decodeResponse()
	if dErr != nil {
		utils.HandleError(utils.ErrorDecodingResponse)
	}

	// TODO: Update here
	fmt.Println(tr.decodedResp)
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
