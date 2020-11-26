package bencode

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

type MetaInfo struct {
	Announce     string
	AnnounceList []string
	CreationDate time.Time
	Comment      string
	CreatedBy    string
	Encoding     string
	Info         *InfoDictionary
}

type InfoDictionary struct {
	PieceLength int64
	Pieces      string
	Private     int64
	Name        string
	Files       []*FilePiece
}

type FilePiece struct {
	Length int64
	Path   string
}

// Unpack map's data into a metainfo structure to actually make use of the data.
func Unpack(data map[string]interface{}) *MetaInfo {
	announce := data["announce"].(string)
	announceList := data["announce-list"].([]interface{})
	unpackedList := make([]string, 0, len(announceList))
	for i := 0; i < len(announceList); i++ {
		announceURL := announceList[i].([]interface{})
		unpackedList = append(unpackedList, announceURL[0].(string))
	}

	// Convert unix timestamp to datetime
	sec := data["creation date"].(int64)
	tm := time.Unix(sec, 0)
	encoding := data["encoding"].(string)
	createdBy := data["created by"].(string)

	meta := &MetaInfo{
		Announce:     announce,
		AnnounceList: unpackedList,
		CreatedBy:    createdBy,
		CreationDate: tm,
		Encoding:     encoding,
	}

	// unpack info dictionary
	inf := data["info"].(map[string]interface{})
	pieceLength := inf["piece length"].(int64)
	pieces := inf["pieces"].(string)
	private := inf["private"].(int64)
	name := inf["name"].(string)

	infoDict := &InfoDictionary{
		PieceLength: pieceLength,
		Pieces:      pieces,
		Private:     private,
		Name:        name,
	}

	// unpack files list
	files := inf["files"].([]interface{})
	fPieces := make([]*FilePiece, 0, len(files))
	for i := 0; i < len(files); i++ {
		mp := files[i].(map[string]interface{})
		pLength := mp["length"].(int64)
		pathSlice := mp["path"].([]interface{})

		pPath := ""
		for i := 0; i < len(pathSlice); i++ {
			f := pathSlice[i].(string)
			pPath = filepath.Join(pPath, f)
		}

		fPiece := &FilePiece{
			Length: pLength,
			Path:   pPath,
		}
		fPieces = append(fPieces, fPiece)
	}

	infoDict.Files = fPieces
	meta.Info = infoDict
	return meta
}

// Cache the torrent file on disk in json format, avoiding the overhead of re-decoding the file everytime
// we add the torrent file to the download queue (This situation only happens if we want to re-download the torrent).
func CacheFile(metadata *MetaInfo) error {
	data, err := json.MarshalIndent(metadata, "", " ")

	f, err := os.Create(metadata.Info.Name + ".json")
	if err != nil {
		return nil
	}

	n, err := f.Write(data)
	if err != nil {
		return err
	} else if n != len(data) {
		return errors.New("Write failed")
	}

	return nil
}
