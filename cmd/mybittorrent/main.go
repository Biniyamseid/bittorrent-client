package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jackpal/bencode-go"
)

type FileInfo struct {
	Announce string
	Info     struct {
		Length int64
	}
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
		Info: struct {
			Length int64
		}{
			Length: torrInfoMap["info"].(map[string]interface{})["length"].(int64),
		},
	}

	return fileInfo, nil
}

//func ParseFile(torrentFile string) (*FileInfo, error) {
//	data, err := os.ReadFile(torrentFile)
//	if err != nil {
//		return nil, err
//	}
//
//	reader := strings.NewReader(string(data))
//	torrInfo, err := bencode.Decode(reader)
//	if err != nil {
//		return nil, err
//	}
//
//	torrInfoMap := torrInfo.(map[string]interface{})
//	fileInfo := &FileInfo{
//		Announce: torrInfoMap["announce"].(string),
//		Info: struct {
//			Length int
//		}{
//			Length: torrInfoMap["info"].(map[string]interface{})["length"].(int),
//		},
//	}
//
//	return fileInfo, nil
//}

func Decode(bencodedValue string) (interface{}, error) {
	decoded, err := decodeBencode(bencodedValue)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

func decodeBencode(bencodedString string) (interface{}, error) {
	reader := strings.NewReader(bencodedString)
	return bencode.Decode(reader)
}
func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	command := os.Args[1]

	switch command {
	case "decode":
		bencodedValue := os.Args[2]

		decoded, err := decodeBencode(bencodedValue)
		if err != nil {

			fmt.Println(err)
			return
		}

		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))
		// case "info":
		// 	filename := os.Args[2]
		// 	data, err := os.ReadFile(filename)
		// 	panicIf(err)
		// 	torrContent := string(data)
		// 	torrInfo, err := bencode.Decode(strings.NewReader(torrContent))
		// 	panicIf(err)
		// 	url := torrInfo.(map[string]interface{})["announce"].(string)
		// 	info := torrInfo.(map[string]interface{})["info"].(map[string]interface{})
		// 	length := info["length"].(int)

		// 	fmt.Printf("Tracker URL: %s\nLength: %d\n", url, length)

	case "info":
		torrentFile := os.Args[2]
		if len(torrentFile) == 0 {
			log.Fatalf("No argument provided for 'info'")
		}
		fileInfo, err := ParseFile(torrentFile)
		if err != nil {
			log.Fatalf("Failed to parse %q: %v", torrentFile, err)
		}
		fmt.Printf("Tracker URL: %s\n", fileInfo.Announce)
		fmt.Printf("Length: %d", fileInfo.Info.Length)

	default:
		fmt.Println("Unknown command: " + command)
		os.Exit(1)

	}

}
