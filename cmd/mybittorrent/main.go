package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func decodeBencode(bencodedString string) (interface{}, int, error) {
	if unicode.IsDigit(rune(bencodedString[0])) {
		// string case
		var firstColonIndex int

		for i := 0; i < len(bencodedString); i++ {
			if bencodedString[i] == ':' {
				firstColonIndex = i
				break
			}
		}

		lengthStr := bencodedString[:firstColonIndex]

		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return "", 0, err
		}

		untilIndex := firstColonIndex + 1 + length
		return bencodedString[firstColonIndex+1 : untilIndex], untilIndex, nil
	} else if strings.HasPrefix(bencodedString, "i") {
		// integers case
		var endIndex int

		for i := 0; i < len(bencodedString); i++ {
			if bencodedString[i] == 'e' {
				endIndex = i
				break
			}
		}

		num, err := strconv.Atoi(bencodedString[1:endIndex])
		if err != nil {
			return "", 0, err
		}

		return num, endIndex + 1, nil
	} else if strings.HasPrefix(bencodedString, "l") {
		// list case
		in := strings.TrimPrefix(bencodedString, "l")

		var (
			ret        = []interface{}{}
			untilIndex int
		)
		for {
			if in[0] == 'e' {
				break
			}

			decoded, nextIndex, err := decodeBencode(in)
			if err != nil {
				return "", 0, err
			}
			ret = append(ret, decoded)

			in = in[nextIndex:]
			untilIndex += nextIndex
		}

		return ret, untilIndex + 1, nil
	} else if strings.HasPrefix(bencodedString, "d") {
		// dictionary case
		in := strings.TrimPrefix(bencodedString, "d")

		var (
			ret        = map[string]interface{}{}
			key        string
			untilIndex int
		)
		for {
			if in[0] == 'e' {
				break
			}

			decoded, nextIndex, err := decodeBencode(in)
			if err != nil {
				return "", 0, err
			}
			if key == "" {
				key = decoded.(string)
			} else {
				ret[key] = decoded
				key = ""
			}

			in = in[nextIndex:]
			untilIndex += nextIndex
		}

		return ret, untilIndex + 1, nil
	} else {
		return "", 0, fmt.Errorf("unexpected format")
	}
}
func main() {
	command := os.Args[1]

	switch command {
	case "decode":
		bencodedValue := os.Args[2]

		decoded, _, err := decodeBencode(bencodedValue)
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
