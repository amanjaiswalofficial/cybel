package cmd

import (
	"cybele/ops/utils"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

// listCmd lists the currrent queue of torrents added
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List items in current queue of torrents to download",
	Run:   RunListCmd,
}

// RunListCmd prints the list of files already in the queue
func RunListCmd(cmd *cobra.Command, args []string) {

	fileNames, err := GetQueueFiles()
	if err != nil {
		fmt.Println(err.Error())
	}
	for i := 0; i < len(fileNames)-1; i++ {
		fmt.Printf("%v.%v\n", i+1, fileNames[i])
	}

}

// GetQueueFiles reads files from queue and returns them
func GetQueueFiles() ([]string, error) {

	// Read queue from path
	filePath := filepath.Join(utils.CybeleCachePath, utils.QueueFileName)
	
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		utils.LogMessage(err.Error())
		return nil, errors.New("Error reading queue")
	}
	fileList := strings.Split(string(data), "\n")

	var fileNames []string
	for _, fileName := range fileList {
		extension := filepath.Ext(fileName)
		name := fileName[0 : len(fileName)-len(extension)]
		fileNames = append(fileNames, name)
	}

	return fileNames, nil

}
