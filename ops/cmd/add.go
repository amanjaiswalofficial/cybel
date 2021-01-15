package cmd

import (
	"cybele/ops/connect"
	"cybele/ops/connect/udp"
	"cybele/ops/utils"
	
	"net/url"
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
	fExtension := filepath.Ext(args[0])
	if fExtension != ".torrent" {
		utils.LogMessage("Not a torrent file")
		os.Exit(1)
	}

	// Process the file and add it to the download queue
	td, err := connect.WriteJSON(args[0])
	if err != nil {
		utils.HandleError(err.Error())
	}

	var peers []connect.PeerObject
	peers, err = getPeersFromTracker(td.Announce, td)

	// In case the main announce url doesn't work, retry connecting with the urls
	// in announce-list (if there's any)
	if err != nil {
		for i := 0; i < len(td.AnnounceList); i++ {
			peers, err = getPeersFromTracker(td.AnnounceList[i], td)
			if err != nil {
				continue
			}
		}
	}

	if err != nil {
		// None of the urls works,
		utils.HandleError("Torrent urls broken")
	}

	// to change
	for _, peer := range peers {
		utils.LogMessage("IP:", peer.IP)
		utils.LogMessage("Port:", peer.Port)
	}
}

func getPeersFromTracker(rawUrl string, td *connect.TorrentData) ([]connect.PeerObject, error) {
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}

	var peers []connect.PeerObject
	if parsedUrl.Scheme == "udp" {
		peers, err = getPeersFromUDP(parsedUrl.Host, td)
		if err != nil {
			return nil, err
		}

		return peers, nil
	}

	peers, err = getPeersFromHTTP(td)
	if err != nil {
		return nil, err
	}
	return peers, nil
}

func getPeersFromUDP(host string, td *connect.TorrentData) ([]connect.PeerObject, error) {
	tr := udp.New(host)
	req := udp.MakeRequestObject(td)
	resp, err := tr.Announce(req)
	if err != nil {
		return nil, err
	}

	return resp.Peers, nil
}

func getPeersFromHTTP(td *connect.TorrentData) ([]connect.PeerObject, error) {
	return nil, nil
}
