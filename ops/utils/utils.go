package utils

import (
	"bytes"
	"crypto/sha1"
	"cybele/ops/bencode"
	"encoding/hex"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	enc "github.com/jackpal/bencode-go"
)

// ReadFileFromPath takes input as a path to read file from
// returns: error if exists, otherwise returns file
func ReadFileFromPath(path string) ([]byte, error) {
	readData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return readData, nil
}

// HandleError is used to ensure single place to handle all errors
// After displaying the message, exits the program
func HandleError(err string) {
	LogMessage(err)
	os.Exit(1)
}

// LogMessage is used to log messages to the console
func LogMessage(messages ...string) {
	if os.Getenv("LOGGING") == "true" {
		log.Println(strings.Join(messages, " "))
	}
}

// MakeInfoHash takes input as the hash from the torrentData
// And converts into a format Tracker can understand
// returns: hash in hex format
func MakeInfoHash(basicHash string) string {
	var resultHash string

	var convertedHashArray []string
	// Ex - basicHash = "5149527e0e68e9f9a7f104b7b35dd1ea0f04b4bd"
	totalLen := len(basicHash)
	for i := 0; i < totalLen; i += 2 {
		val, _ := strconv.ParseInt(string(basicHash[i:i+2]), 16, 16)
		var res string
		if val <= 127 {
			res = url.QueryEscape(string(val))
			if string(res[0]) == "%" {
				res = "%" + strings.ToLower(res[1:])
			}
		} else {
			res = "%" + string(basicHash[i:i+2])
		}
		convertedHashArray = append(convertedHashArray, res)
	}
	resultHash = strings.Join(convertedHashArray, "")
	return resultHash
}

// FormatKey is used to properly format strings to expected struture
// Ex: Takes max interval and returns maxInterval
// returns: formatted value for a string
func FormatKey(key string) string {

	keySplit := strings.Split(key, " ")
	for i := 0; i < len(keySplit); i++ {
		if len(keySplit[i]) == 2 {
			keySplit[i] = strings.ToUpper(keySplit[i])
		} else {
			keySplit[i] = strings.Title(keySplit[i])
		}

	}
	formattedKey := strings.Join(keySplit, "")
	return formattedKey
}

// ComputeInfoHash takes a torrent file path
// and computes a SHA1 hash over the info
// dictionary. returns: sha1 hash encoded in hexadecimal
func ComputeInfoHash(path string) string {
	f, err := os.Open(path)
	if err != nil {
		HandleError(err.Error())
	}

	raw, err := bencode.Decode(f)
	if err != nil {
		HandleError(err.Error())
	}

	newBuf := new(bytes.Buffer)

	err = enc.Marshal(newBuf, raw["info"])
	if err != nil {
		HandleError(err.Error())
	}

	binaryHash := sha1.Sum(newBuf.Bytes())
	return hex.EncodeToString(binaryHash[:])
}

// AddTo cache takes input as the torrent filename,
// the torrent file data (in its byte form) and write
// it to the cache directory (and to the download queue
// as well). returns: an error if any.
func AddToCache(filename string, data []byte) error {
	// Check if the cache directory exists, if not, create it
	if _, err := os.Stat(CybeleCachePath); os.IsNotExist(err) {
		// Create the directory (and its parents) with unix permission bits
		err = os.MkdirAll(CybeleCachePath, 0777)
		if err != nil {
			return err
		}
	}

	fname := strings.Join([]string{filename, ".json"}, "")
	fpath := filepath.Join(CybeleCachePath, fname)
	f, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer f.Close()

	f.Write(data)

	// Create the download queue file.
	queuePath := filepath.Join(CybeleCachePath, "queue")
	queueFile, err := os.OpenFile(queuePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)

	if err != nil {
		return err
	}
	defer queueFile.Close()

	queueFile.Write([]byte(filename + "\n"))
	return nil
}

func MakePeerID() string {
	return string("-AA1111-123456789012")
}
