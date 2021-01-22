package cmd

import (
	"cybele/ops/connect"
	"cybele/ops/utils"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

// listCmd lists the currrent queue of torrents added
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "List items in current queue of torrents to download",
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

	var hs connect.Handshake

	hsStr := hs.GetString([]byte(torrentData.InfoHash), []byte(peerID))

	connect.DoHandshake(hsStr, []byte(torrentData.InfoHash), trackerObject.DecodedResp.Peers)

	// for _, peerObj := range trackerObject.DecodedResp.Peers {
	// 	printString := fmt.Sprintf("will connect to %v:%v", peerObj.IP, peerObj.Port)

	// }
}
