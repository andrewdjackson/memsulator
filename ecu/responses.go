package ecu

import (
	"bufio"
	"encoding/hex"
	"github.com/andrewdjackson/memsulator/utils"
	"log"
	"os"
	"strings"
)

// UseResponseFile determine whether to return response sequentially from a response file
var UseResponseFile bool
var responseFile []string
var index = 0
var responseMap = make(map[string][]byte)

// package init function
func init() {
	// Response formats for commands that do not respond with the format [COMMAND][VALUE]
	// Generally these are either part of the initialisation sequence or are ECU data frames
	responseMap["0a"] = []byte{0x0A}
	responseMap["ca"] = []byte{0xCA}
	responseMap["75"] = []byte{0x75}

	// Format for DataFrames starts with [Command Echo][Data Size][Data Bytes (28 for 0x80 and 32 for 0x7D)]
	responseMap["80"] = []byte{0x80, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B}
	responseMap["7d"] = []byte{0x7d, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F}
	responseMap["d0"] = []byte{0xD0, 0x99, 0x00, 0x03, 0x03}

	// heatbeat
	responseMap["f4"] = []byte{0xf4, 0x00}

	// adjustments
	responseMap["7a"] = []byte{0x7a, 0x8a}
	responseMap["7c"] = []byte{0x7c, 0x8a}
	responseMap["79"] = []byte{0x79, 0x8a}
	responseMap["7b"] = []byte{0x7b, 0x8a}
	responseMap["8a"] = []byte{0x8a, 0x23}
	responseMap["89"] = []byte{0x89, 0x23}
	responseMap["92"] = []byte{0x92, 0x80}
	responseMap["91"] = []byte{0x91, 0x80}
	responseMap["94"] = []byte{0x94, 0x80}
	responseMap["93"] = []byte{0x93, 0x80}

	//resets
	responseMap["fa"] = []byte{0xfa, 0x00}
	responseMap["0f"] = []byte{0x0f, 0x00}
	responseMap["cc"] = []byte{0xcc, 0x00}

	// generic response, expect command and single byte response
	responseMap["00"] = []byte{0x00, 0x00}
}

// Response returns an emulated response byte string
func (mems *MemsConnection) Response(command []byte) []byte {
	c := hex.EncodeToString(command)

	// if the command is a dataframe request and we have a response file
	// then use the response file
	if c == "80" || c == "7d" {
		utils.LogI.Printf("reading from log file..")
		return mems.getLogResponse(c)
	}

	// otherwise generate the response
	return mems.generateResponse(c)
}

func (mems *MemsConnection) generateResponse(command string) []byte {
	r := responseMap[command]

	if r == nil {
		r = responseMap["00"]
		copy(r[0:], command)
	}

	utils.LogI.Printf("generating response %x", r)
	return r
}

func (mems *MemsConnection) getLogResponse(command string) []byte {
	var data []byte
	var response string

	command = strings.ToUpper(command)

	// load the data if not already loaded
	if len(responseFile) == 0 {
		responseFile, _ = mems.readResponseFile("response.data")
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
			response = strings.ReplaceAll(response, "\n", "")
			response = strings.ReplaceAll(response, "\r", "")

			log.Printf("> %s", response)

			// convert to byte array
			data, err := hex.DecodeString(response)

			// truncate to the right size
			if command == "80" {
				data = data[:29]
			} else {
				data = data[:33]
			}

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

func (mems *MemsConnection) readResponseFile(path string) ([]string, error) {
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
