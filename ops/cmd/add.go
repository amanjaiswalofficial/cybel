package cmd

import (
	"cybele/ops/connect"
	"cybele/ops/utils"
	"github.com/spf13/cobra"
	"os"
)

// Register the add command as a subcommand
func init() {
	rootCmd.AddCommand(addCmd)
}

// addCmd add a torrent file (possibly many) to the download queue
var addCmd = &cobra.Command{
	Use:   "add <path/to/file>.torrent",
	Short: "Download files/content from <path/to/file>.torrent",
	Run:   RunAddCmd,
	Args:  cobra.ExactArgs(1),
}

func RunAddCmd(cmd *cobra.Command, args []string) {
	// Show an error is the file doesn't exist
	if _, err := os.Stat(args[0]); os.IsNotExist(err) {
		utils.HandleError(err.Error())
	}
	// Make sure the provided file is actually a .torrent file
	for i := len(args[0]) - 1; i >= 0; i-- {
		if args[0][i] == '.' {
			if args[0][i:] != ".torrent" {
				utils.LogMessage("Not a torrent file")
				os.Exit(1)
			}
		}
	}

	// Process the file and add it to the download queue
	err := connect.WriteJSON(args[0])
	if err != nil {
		utils.HandleError(err.Error())
	}
	utils.LogMessage("Succesfully add", args[0], "to the download queue")
}
