package cmd

import (
	"cybele/ops/connect"
	"cybele/ops/connect/udp"
	"cybele/ops/handshake"
	"cybele/ops/utils"
	"net/url"
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
	bs, err := utils.ReadFileFromPath(jsonPath)
	torrent := connect.ReadJSONFromByteSlice(bs)

	peers, err := GetPeers(&torrent)
	if err != nil {
		utils.HandleError(err.Error())
	}

	peerID := utils.MakePeerID()

	var hs handshake.Handshake

	hsStr := hs.GetString([]byte(torrent.InfoHash), []byte(peerID))

	handshake.DoHandshake(hsStr, []byte(torrent.InfoHash), peers)
}

func GetPeers(torrent *connect.TorrentData) ([]connect.PeerObject, error) {
	urlp, err := url.Parse(torrent.Announce)
	if err != nil {
		return nil, err
	}

	var peers []connect.PeerObject

	if urlp.Scheme == "udp" {
		tracker := udp.New(urlp.Host)
		req := udp.MakeRequestObject(torrent)
		resp, err := tracker.Announce(req)
		if err != nil {
			return nil, err
		}
		peers = resp.Peers
	} else {
		resp := connect.ConnectToTracker(*torrent)
		peers = resp.DecodedResp.Peers
	}

	return peers, nil
}
