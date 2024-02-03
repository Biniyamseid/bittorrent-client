package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/jackpal/bencode-go"
)

type FileInfo struct {
	Announce string
	Info     map[string]interface{}
}

func ParseFile(torrentFile string) (*FileInfo, error) {
	data, err := os.ReadFile(torrentFile)
	if err != nil {
		return nil, err
	}

	reader := strings.NewReader(string(data))
	torrInfo, err := bencode.Decode(reader)
	if err != nil {
		return nil, err
	}

	torrInfoMap := torrInfo.(map[string]interface{})
	fileInfo := &FileInfo{
		Announce: torrInfoMap["announce"].(string),
		Info:     torrInfoMap["info"].(map[string]interface{}),
	}

	return fileInfo, nil
}

// encode
func EncodeBytes(value interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := encode(&buf, reflect.ValueOf(value))
	return buf.Bytes(), err
}

func encode(w io.Writer, value reflect.Value) error {
	switch value.Kind() {
	case reflect.String:
		return encodeString(w, value.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return encodeInt(w, value.Int())
	case reflect.Slice:
		return encodeList(w, value)
	case reflect.Map:
		return encodeDict(w, value)
	case reflect.Interface:
		// Use a type switch to handle the actual type of the value
		switch v := value.Interface().(type) {
		case string:
			return encodeString(w, v)
		case int, int8, int16, int32, int64:
			return encodeInt(w, reflect.ValueOf(v).Int())
		case []interface{}:
			return encodeList(w, reflect.ValueOf(v))
		case map[string]interface{}:
			return encodeDict(w, reflect.ValueOf(v))
		default:
			return fmt.Errorf("unsupported type: %T", v)
		}
	default:
		return fmt.Errorf("unsupported type: %s", value.Type())
	}
}

func encodeString(w io.Writer, s string) error {
	_, err := fmt.Fprintf(w, "%d:%s", len(s), s)
	return err
}

func encodeInt(w io.Writer, i int64) error {
	_, err := fmt.Fprintf(w, "i%de", i)
	return err
}

func encodeList(w io.Writer, v reflect.Value) error {
	_, err := fmt.Fprint(w, "l")
	if err != nil {
		return err
	}
	for i := 0; i < v.Len(); i++ {
		err := encode(w, v.Index(i))
		if err != nil {
			return err
		}
	}
	_, err = fmt.Fprint(w, "e")
	return err
}

func encodeDict(w io.Writer, v reflect.Value) error {
	_, err := fmt.Fprint(w, "d")
	if err != nil {
		return err
	}
	keys := v.MapKeys()
	// Convert the keys to a slice of strings
	strKeys := make([]string, len(keys))
	for i, key := range keys {
		strKeys[i] = key.String()
	}
	// Sort the slice of strings
	sort.Strings(strKeys)
	// Iterate over the sorted slice to encode the keys and their corresponding values
	for _, strKey := range strKeys {
		key := reflect.ValueOf(strKey)
		err := encode(w, key)
		if err != nil {
			return err
		}
		err = encode(w, v.MapIndex(key))
		if err != nil {
			return err
		}
	}
	_, err = fmt.Fprint(w, "e")
	return err
}
func main() {
	command := os.Args[1]

	switch command {
	case "decode":
		bencodedValue := os.Args[2]

		decoded, err := bencode.Decode(strings.NewReader(bencodedValue))
		if err != nil {

			fmt.Println(err)
			return
		}

		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))

	case "info":
		torrentFile := os.Args[2]
		if len(torrentFile) == 0 {
			log.Fatalf("No argument provided for 'info'")
		}
		fileInfo, err := ParseFile(torrentFile)
		if err != nil {
			log.Fatalf("Failed to parse %q: %v", torrentFile, err)
		}

		// Bencode the info dictionary
		bencodedInfo, err := EncodeBytes(fileInfo.Info)
		if err != nil {
			log.Fatalf("Failed to bencode info: %v", err)
		}

		// Calculate the SHA-1 hash
		hasher := sha1.New()
		hasher.Write(bencodedInfo)
		infoHash := hex.EncodeToString(hasher.Sum(nil))

		fmt.Printf("Tracker URL: %s\n", fileInfo.Announce)
		fmt.Printf("Length: %d\n", fileInfo.Info["length"].(int64))
		fmt.Printf("Info Hash: %s\n", infoHash)

	default:
		fmt.Println("Unknown command: " + command)
		os.Exit(1)

	}

}
