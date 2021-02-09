package cmd

import (
	"cybele/ops/connect"
	"cybele/ops/utils"

	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// Register the add command as a subcommand
func init() {
	rootCmd.AddCommand(addCmd)
}

// addCmd add a torrent file (possibly many) to the download queue
var addCmd = &cobra.Command{
	Use:   "add <path/to/file>.torrent",
	Short: "Add torrents to the download queue",
	Run:   RunAddCmd,
	Args:  cobra.ExactArgs(1),
}

func RunAddCmd(cmd *cobra.Command, args []string) {
	// Show an error is the file doesn't exist
	if _, err := os.Stat(args[0]); os.IsNotExist(err) {
		utils.HandleError(err.Error())
	}

	// Make sure the provided file is actually a .torrent file
	fExtension := filepath.Ext(args[0])
	if fExtension != ".torrent" {
		utils.HandleError("Not a torrent file")
	}

	// Process the file and add it to the download queue
	torrent, err := connect.WriteJSON(args[0])
	if err != nil {
		utils.HandleError(err.Error())
	}

	utils.LogMessage(torrent.Filename, "added to download queue")
}
