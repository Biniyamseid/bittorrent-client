package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/jackpal/bencode-go"
)

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

	default:
		fmt.Println("Unknown command: " + command)
		os.Exit(1)

	}

}
