package torrent

import (
    "math"
    "cybele/ops/connect"
 )

type Piece struct {
	Index uint32
	Hash  []byte
	Size  uint64
}

func ParsePieces(tr connect.TorrentData) []Piece {
	npieces := uint32(math.Ceil(float64(tr.TotalSize) / float64(tr.PieceSize)))
	pieces := make([]Piece, 0, npieces)

	piecesHash := []byte(tr.PiecesHash)
	for i := uint32(0); i < npieces; i++ {
		begin := uint64(i*tr.PieceSize)
        end := uint64(begin + uint64(tr.PieceSize))

		if end > tr.TotalSize {
			end = tr.TotalSize
		}

		piece := Piece{
			Index: uint32(i),
			Hash:  piecesHash[i*20 : (i+1)*20],
			Size:  end - begin,
		}

		pieces = append(pieces, piece)
	}

	return pieces
}
