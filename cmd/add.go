package cmd

import (
	"github.com/spf13/cobra"
)

// Register the add command as a subcommand
func init() {
	rootCmd.AddCommand(addCmd)
}

// addCmd add a torrent file (possibly many) to the download queue
var addCmd = &cobra.Command{
	Use:   "add <path/to/file>.torrent",
	Short: "Download files/content from <path/to/file>.torrent",
	Run:   func(cmd *cobra.Command, args []string) {},
	Args:  cobra.ExactArgs(1),
}
