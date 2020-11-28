package connect

import (
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

	// Print the json of the torrentData for now.
	td.printInfo()
}
