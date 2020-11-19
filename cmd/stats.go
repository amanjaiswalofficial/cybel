package cmd

import (
	"github.com/spf13/cobra"
)

// Register stats command as a subcommand
func init() {
	statsCmd.Flags().StringSliceP("files", "f", []string{}, "collection of torrent files")
	rootCmd.AddCommand(statsCmd)
}

// Stats subcommand show torrent files statistics.
var statsCmd = &cobra.Command{
	Use:   "stats [-f FILE.torrent... ]",
	Short: "Show torrent files statistics",
	Long: `Show torrent files statistics
           if [-f | --files] is set, show statistics for the files provided`,
	Run: func(cmd *cobra.Command, args []string) {},
}
