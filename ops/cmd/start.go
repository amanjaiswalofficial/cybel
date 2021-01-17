package cmd

import (
	"cybele/ops/utils"
	"cybele/ops/connect"
	"path/filepath"

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
	Args:  cobra.ExactArgs(1),
}

// RunStartCmd prints the list of files already in the queue
func RunStartCmd(cmd *cobra.Command, args []string) {

	fileName := "Wonder Woman 1984 (2020) [1080p] [WEBRip] [5.1] [YTS.MX]"
	fileName = fileName + ".json"
	jsonPath := filepath.Join(utils.CybeleCachePath, fileName)
	utils.LogMessage("Using Wonder Woman 1984 (2020) [1080p] [WEBRip] [5.1] [YTS.MX].json")
	connect.FetchDetailsFromTorrent(jsonPath)

}
