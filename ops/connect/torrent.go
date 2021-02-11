package connect

import (
	"cybele/ops/bencode"
	"cybele/ops/utils"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// TorrentData : A struct for getting the json data
type TorrentData struct {
	Name         string   `json:"name"`
	Filename     string   `json:"filename"`
	InfoHash     string   `json:"info_hash"`
	PieceSize    uint32   `json:"piecesize"`
	TotalSize    uint64   `json:"totalsize"`
	Announce     string   `json:"announce"`
	AnnounceList []string `json:"announce_list"`
	Files        []File   `json:"files"`
	PiecesHash   []byte   `json:"pieces"`
}

// File struct for getting torrent metainfo files field
type File struct {
	Path   string `json:"path"`
	Length uint64 `json:"length"`
}

// PrintInfo properties from the torrent file
func (td TorrentData) PrintInfo() {
	fmt.Printf("name: %+v\n", td.Name)
	fmt.Printf("filename: %+v\n", td.Filename)
	fmt.Printf("infoHash: %+v\n", td.InfoHash)
	fmt.Printf("piece size (in bytes): %d\n", td.PieceSize)
	fmt.Printf("total size (in bytes): %d\n", td.TotalSize)
	fmt.Printf("announceURL: %+v\n", td.Announce)
	fmt.Printf("files:\n")
	for _, f := range td.Files {
		fmt.Printf("Path: %s\n", f.Path)
		fmt.Printf("Length: %d\n\n", f.Length)
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
func WriteJSON(path string) (*TorrentData, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dec, err := bencode.Decode(f)
	if err != nil {
		return nil, err
	}

	meta := bencode.Unpack(dec)
	files := make([]File, 0, len(meta.Info.Files))
	totalSize := uint64(0)

	for i := 0; i < len(meta.Info.Files); i++ {
		f := File{Path: meta.Info.Files[i].Path, Length: uint64(meta.Info.Files[i].Length)}
		totalSize += f.Length
		files = append(files, f)
	}

	hash := utils.ComputeInfoHash(path)

	fname := filepath.Base(path)
	ext := filepath.Ext(fname)

	td := TorrentData{
		Name:         meta.Info.Name,
		Filename:     fname,
		InfoHash:     hash,
		PieceSize:    uint32(meta.Info.PieceLength),
		TotalSize:    totalSize,
		Announce:     meta.Announce,
		AnnounceList: meta.AnnounceList,
		Files:        files,
		PiecesHash:   []byte(meta.Info.Pieces),
	}

	rawBytes, err := json.MarshalIndent(td, "", "\t")
	if err != nil {
		return nil, errors.New(utils.ErrorMarshaling)
	}

	err = utils.AddToCache(fname[:len(fname)-len(ext)], rawBytes)
	if err != nil {
		return nil, err
	}

	return &td, nil
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
