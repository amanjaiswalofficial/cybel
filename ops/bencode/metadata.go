package bencode

import (
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

	// Optional fields (default value if any of the fields doesn't exist)

	// Convert unix timestamp to datetime
	var sec int64
	var encoding, createdBy string
	if _, ok := data["creation date"]; ok {
		sec = data["creation date"].(int64)
	}

	tm := time.Unix(sec, 0)
	if _, ok := data["encoding"]; ok {
		encoding = data["encoding"].(string)
	}

	if _, ok := data["created by"]; ok {
		createdBy = data["created by"].(string)
	}

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
	// Optional field (default if doesnt exist)
	var private int64
	if _, ok := inf["private"]; ok {
		private = inf["private"].(int64)
	}
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
