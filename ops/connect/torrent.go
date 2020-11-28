package connect

import (
	"fmt"
	"encoding/json"
)

// TorrentData : A struct for getting the json data
type TorrentData struct {
	Name string `json:"name"`
	Filename string `json:"filename"`
	Comment string `json:"comment"`
	Date string `json:"data"`
	CreatedBy string `json:"created_by"`
	InfoHash string `json:"info_hash"`
	Size string `json:"size"`
	Announce string `json:"announce"`
	AnnounceList []string `json:"announce_list"`
	Files []string `json:"files"`

}

// Print properties from the torrent file
func (td TorrentData) printInfo() {
	fmt.Printf("name: %+v\n", td.Name)
	fmt.Printf("filename: %+v\n", td.Filename)
	fmt.Printf("infoHash: %+v\n", td.InfoHash)
	fmt.Printf("size(in bytes): %+v\n", td.Size)
	fmt.Printf("announceURL: %+v\n", td.Announce)
	fmt.Printf("files:\n")
	for fileName := range td.Files {
		fmt.Println(fileName)
	}
}

// ReadJSONFromByteSlice accepts a byteslice
// Tries to convert it into a json format
// returns: json of type TorrentData
func ReadJSONFromByteSlice(data []byte) (TorrentData) {
	var td TorrentData
	json.Unmarshal(data, &td)
	return td
}

// IsEmpty checks if a TorrentData struct is empty
// Checks for null values for different variables
// returns: true if empty, otherwise false
func (td TorrentData) IsEmpty() bool {
	if td.Filename == "" {
		return true
	} else if len(td.Files) == 0 {
		return true
	}
	return false
}
