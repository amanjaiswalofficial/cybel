package bencode

import (
	"bufio"
	"errors"
	"io"
	"strconv"
)

// Package bencode maps a .torrent file into Go's map

// Decode takes an io.Reader (bencoded file) and return a Go's map. Otherwise, if any error is encountered,
// a nil map and a error is returned.
func Decode(r io.Reader) (map[string]interface{}, error) {
	buf := bufio.NewReader(r)

	if firstByte, err := buf.ReadByte(); err != nil {
		return nil, err
	} else if firstByte != 'd' {
		return nil, errors.New("bencode data must begin with a dictionary")
	}

	mp, err := decodeDict(buf)
	if err != nil {
		return nil, errors.New("Decode failed")
	}
	return mp, nil
}

// decodeDict decodes a bencoded dictionary.
func decodeDict(buf *bufio.Reader) (map[string]interface{}, error) {
	dict := make(map[string]interface{})
	for {
		key, err := decodeString(buf)
		if err != nil {
			return nil, err
		}

		chr, err := buf.ReadByte()
		if err != nil {
			return nil, err
		}

		val, err := decodeType(buf, chr)
		if err != nil {
			return nil, err
		}

		dict[key] = val

		next, err := buf.ReadByte()
		if err != nil {
			return nil, err
		}
		// if the byte e is reached, it means the dictionary ends here.
		if next == 'e' {
			break
		} else if err := buf.UnreadByte(); err != nil {
			return nil, err
		}

	}

	return dict, nil
}

// decodeString, decodes a bencoded string.
// which have the following scheme: <string length encoded in base 10 ASCII>:<string data>
func decodeString(buf *bufio.Reader) (string, error) {
	length, err := readIntUntil(buf, ':')

	var stringLen int64
	var ok bool
	if stringLen, ok = length.(int64); !ok {
		return "", errors.New("len overflow")
	}

	if stringLen < 0 {
		return "", errors.New("key length cannot be negative")
	}

	buffer := make([]byte, stringLen)
	_, err = io.ReadFull(buf, buffer)
	return string(buffer), err
}

// Reads a byte string representing an integer number until it reachs delim, then it converts it to an int64.
// this function is a helper for other funcs such as decodeString(), decodeInt() and decodeType()
func readIntUntil(buf *bufio.Reader, delim byte) (interface{}, error) {
	slice, err := buf.ReadSlice(delim)
	if err != nil {
		return nil, err
	}

	data := string(slice[:len(slice)-1])
	if num, err := strconv.ParseInt(data, 10, 64); err == nil {
		return num, nil
	}
	return nil, err
}

func decodeType(buf *bufio.Reader, t byte) (interface{}, error) {
	var value interface{}
	var err error
	switch t {
	case 'i':
		value, err = decodeInt(buf)
	case 'l':
		value, err = decodeList(buf)
	case 'd':
		value, err = decodeDict(buf)
	default:
		if err = buf.UnreadByte(); err != nil {
			return nil, err
		}

		value, err = decodeString(buf)
	}

	return value, err
}

// decodeInt decodes a bencoded interger,
// which has the following scheme: i<integer encoded in base 10 ASCII>e
func decodeInt(buf *bufio.Reader) (interface{}, error) {
	return readIntUntil(buf, 'e')
}

func decodeList(buf *bufio.Reader) ([]interface{}, error) {
	var list []interface{}
	for {
		b, err := buf.ReadByte()
		if err != nil {
			return nil, err
		}
		// if the byte e is reached, it means the list ends here
		if b == 'e' {
			break
		}

		value, err := decodeType(buf, b)
		if err != nil {
			return nil, err
		}

		list = append(list, value)
	}

	return list, nil
}
