package bencode

import (
	"path/filepath"
)

type MetaInfo struct {
	Announce     string
	AnnounceList []string
	Info         *InfoDictionary
}

type InfoDictionary struct {
	PieceLength int64
	Pieces      string
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

	meta := &MetaInfo{
		Announce:     announce,
		AnnounceList: unpackedList,
	}

	// unpack info dictionary
	inf := data["info"].(map[string]interface{})
	pieceLength := inf["piece length"].(int64)
	pieces := inf["pieces"].(string)
	name := inf["name"].(string)

	infoDict := &InfoDictionary{
		PieceLength: pieceLength,
		Pieces:      pieces,
		Name:        name,
	}

	// unpack files list
	var fPieces []*FilePiece

	// Multi-file torrent
	if _, ok := inf["files"]; ok {
		files := inf["files"].([]interface{})
		fPieces = make([]*FilePiece, len(files))
		for i := 0; i < len(files); i++ {
			mp := files[i].(map[string]interface{})
			pLength := mp["length"].(int64)
			pathSlice := mp["path"].([]interface{})

			pPath := ""
			for i := 0; i < len(pathSlice); i++ {
				f := pathSlice[i].(string)
				pPath = filepath.Join(pPath, f)
			}

			fPieces[i] = &FilePiece{Length: pLength, Path: pPath}
		}

	} else {
		// Single-file torrent
		fname := inf["name"].(string)
		length := inf["length"].(int64)
		fp := &FilePiece{Path: fname, Length: length}
		fPieces = append(fPieces, fp)
	}

	infoDict.Files = fPieces
	meta.Info = infoDict
	return meta
}
