package cmd

import (
	"cybele/ops/connect"
	"cybele/ops/handshake"
	"cybele/ops/utils"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

// startCmd starts a torrent download.
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a single or multiple torrents downloads",
	Run:   RunStartCmd,
	Args:  cobra.RangeArgs(1, 20),
}

// Example: cybele start "Wonder Woman 1984 (2020) [1080p] [WEBRip] [5.1] [YTS.MX]"
// RunStartCmd prints the list of files already in the queue
func RunStartCmd(cmd *cobra.Command, args []string) {

	fileName := strings.Join(args, " ")
	fileName = fileName + ".json"
	jsonPath := filepath.Join(utils.CybeleCachePath, fileName)
	trackerObject, torrentData := connect.FetchDetailsFromTorrent(jsonPath)

	peerID := utils.MakePeerID()

	var hs handshake.Handshake

	hsStr := hs.GetString([]byte(torrentData.InfoHash), []byte(peerID))

	handshake.DoHandshake(hsStr, []byte(torrentData.InfoHash), trackerObject.DecodedResp.Peers)

}
