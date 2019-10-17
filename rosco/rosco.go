package rosco

import (
	"bufio"
	"encoding/hex"
	"os"
	"strings"
)

// UseResponseFile determine whether to return response sequentially from a response file
var UseResponseFile bool

var responseFile []string
var rosco = make(map[string][]byte)
var index = 0

func init() {
	UseResponseFile = false

	//rosco["0a"] = []byte{0x0A}
	rosco["80"] = []byte{0x80, 0x1C, 0x03, 0x5B, 0x8B, 0xFF, 0x56, 0xFF, 0x22, 0x8B, 0x1D, 0x00, 0x10, 0x01, 0x00, 0x00, 0x00, 0x24, 0x90, 0x2E, 0x00, 0x03, 0x00, 0x48, 0x06, 0x61, 0x10, 0x00, 0x00}
	rosco["7d"] = []byte{0x7d, 0x20, 0x10, 0x0D, 0xFF, 0x92, 0x00, 0x69, 0xFF, 0xFF, 0x00, 0x00, 0x96, 0x64, 0x00, 0xFF, 0x34, 0xFF, 0xFF, 0x30, 0x80, 0x7F, 0xFE, 0xFF, 0x19, 0x00, 0x1E, 0x80, 0x26, 0x40, 0x34, 0xC0, 0x1A}
	rosco["d0"] = []byte{0xD0, 0x99, 0x00, 0x03, 0x03}
	rosco["ca"] = []byte{0xCA}
	rosco["75"] = []byte{0x75}
	rosco["f0"] = []byte{0xF0, 0x00}
	rosco["f4"] = []byte{0xF4, 0x00}
}

// Response returns an emulated response byte string
func Response(command []byte) []byte {
	c := hex.EncodeToString(command)

	// if the command is a dataframe request and we have a response file
	// then use the response file
	if c == "80" || c == "7d" {
		//if UseResponseFile {
		return getLogResponse(c)
		//}
	}

	// otherwise generate the response
	return generateResponse(c)
}

func generateResponse(command string) []byte {
	return rosco[command]
}

func getLogResponse(command string) []byte {
	var data []byte
	var response string

	command = strings.ToUpper(command)

	// load the data if not already loaded
	if len(responseFile) == 0 {
		responseFile, _ = readResponseFile("response.log")
	}

	for i := index; i < len(responseFile); i++ {
		response = responseFile[i]

		// command matches
		if strings.HasPrefix(response, command) {
			// remove the : character
			response = strings.ReplaceAll(response, ":", "")
			// remove all the spaces
			response = strings.ReplaceAll(response, " ", "")
			// remove all the LF
			response = strings.TrimSuffix(response, "\n\r")

			// convert to byte array
			data, err := hex.DecodeString(response)

			if err != nil {
				panic(err)
			}

			// increment index
			index = i

			if index > len(responseFile) {
				// loop back to start
				index = 0
			}

			return data
		}
	}

	return data
}

func readResponseFile(path string) ([]string, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}
