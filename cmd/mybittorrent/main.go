package main

import (
	// Uncomment this line to pass the first stage
	// "encoding/json"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"unicode"
	// bencode "github.com/jackpal/bencode-go" // Available if you need it!
)

// Example:
// - 5:hello -> hello
// - 10:hello12345 -> hello12345
// func decodeBencode(bencodedString string) (interface{}, error) {

// 	if strings.HasPrefix(bencodedString, "i") && strings.HasSuffix(bencodedString, "e") {
// 		return strconv.Atoi(bencodedString[1 : len(bencodedString)-1])
// 	}
// 	if strings.HasPrefix(bencodedString, "l") && strings.HasSuffix(bencodedString, "e") {
// 		var decodedList []interface{}
// 		for _, element := range strings.Split(bencodedString[1:len(bencodedString)-1], "e") {
// 			decodedElement, err := decodeBencode(element + "e")
// 			if err != nil {
// 				return nil, err
// 			}
// 			decodedList = append(decodedList, decodedElement)
// 		}
// 		return decodedList, nil
// 	}
// 	if unicode.IsDigit(rune(bencodedString[0])) {
// 		var firstColonIndex int

// 		for i := 0; i < len(bencodedString); i++ {
// 			if bencodedString[i] == ':' {
// 				firstColonIndex = i
// 				break
// 			}
// 		}

// 		lengthStr := bencodedString[:firstColonIndex]

// 		length, err := strconv.Atoi(lengthStr)
// 		if err != nil {
// 			return "", err
// 		}

// 		return bencodedString[firstColonIndex+1 : firstColonIndex+1+length], nil
// 	} else {
// 		return "", fmt.Errorf("Only strings are supported at the moment")
// 	}

// }

func DecodeString(bencodedString string) (string, error) {
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
		return "", err
	}
	return bencodedString[firstColonIndex+1 : firstColonIndex+1+length], nil
}
func DecodeInteger(bencodedString string) (int, error) {
	var endAt int
	for i := 0; i < len(bencodedString); i++ {
		if bencodedString[i] == 'e' {
			endAt = i
			break
		}
	}
	stringToDecode := bencodedString[1:endAt]
	return strconv.Atoi(stringToDecode)
}
func DecodeList(bencodedString string) ([]interface{}, error) {
	values := []interface{}{}
	if len(bencodedString) < 2 {
		return values, fmt.Errorf("bencoded lists are at least two characters long")
	}
	if bencodedString[0] != 'l' {
		return values, fmt.Errorf("bencoded lists start with l")
	}
	if len(bencodedString) == 2 {
		return values, nil
	}
	position := 0
	bencodedList := bencodedString[1 : len(bencodedString)-1]
	for position < len(bencodedList)-1 {
		listToDecode := bencodedList[position:]
		element, err := decodeBencode(listToDecode)
		if err != nil {
			return []interface{}{}, err
		}
		values = append(values, element)
		position += len(fmt.Sprint(element)) + 2
		if position >= len(bencodedList) || bencodedList[position] == 'e' {
			break
		}
	}
	return values, nil
}
func decodeBencode(bencodedString string) (interface{}, error) {
	if len(bencodedString) == 0 {
		return "", fmt.Errorf("empty string provided")
	}
	firstChar := rune(bencodedString[0])
	switch {
	case unicode.IsDigit(firstChar):
		return DecodeString(bencodedString)
	case firstChar == 'i':
		return DecodeInteger(bencodedString)
	case firstChar == 'l':
		return DecodeList(bencodedString)
	default:
		return nil, fmt.Errorf("unexpected rune")
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	// fmt.Println("Logs from your program will appear here!")

	command := os.Args[1]

	if command == "decode" {
		// Uncomment this block to pass the first stage

		bencodedValue := os.Args[2]

		decoded, err := decodeBencode(bencodedValue)
		if err != nil {
			fmt.Println(err)
			return
		}
		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
