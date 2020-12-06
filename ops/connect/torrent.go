package connect

import (
	"cybele/ops/bencode"
	"cybele/ops/utils"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

// TorrentData : A struct for getting the json data
type TorrentData struct {
	Name         string   `json:"name"`
	Filename     string   `json:"filename"`
	Comment      string   `json:"comment"`
	Date         string   `json:"date"`
	CreatedBy    string   `json:"created_by"`
	InfoHash     string   `json:"info_hash"`
	Size         string   `json:"size"`
	Announce     string   `json:"announce"`
	AnnounceList []string `json:"announce_list"`
	Files        []string `json:"files"`
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
func ReadJSONFromByteSlice(data []byte) TorrentData {
	var td TorrentData
	json.Unmarshal(data, &td)
	return td
}

// WriteJSON takes a torrent file path,
// decodes the bencoded data and encode
// it as json, then writes that json
// out to a file. returns: an error if any.
func WriteJSON(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	dec, err := bencode.Decode(f)
	if err != nil {
		return err
	}

	meta := bencode.Unpack(dec)
	files := make([]string, 0, len(meta.Info.Files))

	for i := 0; i < len(meta.Info.Files); i++ {
		files = append(files, meta.Info.Files[i].Path)
	}

	hash := utils.ComputeInfoHash(path)
	size := fmt.Sprintf("%d", meta.Info.PieceLength)

	td := TorrentData{
		Name:         meta.Info.Name,
		Filename:     strings.Join([]string{meta.Info.Name, ".json"}, ""),
		Date:         meta.CreationDate.String(),
		Comment:      meta.Comment,
		CreatedBy:    meta.CreatedBy,
		InfoHash:     hash,
		Size:         size,
		Announce:     meta.Announce,
		AnnounceList: meta.AnnounceList,
		Files:        files,
	}

	rawBytes, err := json.MarshalIndent(td, "", "\t")
	if err != nil {
		return errors.New(utils.ErrorMarshaling)
	}

	err = utils.AddToCache(td.Filename, rawBytes)
	if err != nil {
		return err
	}

	return nil
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
