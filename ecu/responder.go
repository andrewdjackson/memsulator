package ecu

import (
	"encoding/hex"
	"strings"

	"github.com/andrewdjackson/memsulator/scenarios"
	"github.com/andrewdjackson/memsulator/utils"
)

// PlaybookResponse type
type PlaybookResponse struct {
	dataframe7d []byte
	dataframe80 []byte
}

// Playbook struct
type Playbook struct {
	responses         []PlaybookResponse
	position          int
	count             int
	servedDataframe7d bool
	servedDataframe80 bool
}

// Responder struct
type Responder struct {
	playbook    Playbook
	responseMap map[string][]byte
}

// NewResponder creates an instance of a responder
func NewResponder() *Responder {
	responder := &Responder{}
	responder.responseMap = make(map[string][]byte)

	responder.buildResponseMap()

	return responder
}

// LoadScenario loads a scenario for playing from the ECU
func (responder *Responder) LoadScenario(scenario *scenarios.Scenario) {
	// reset the position of the playbook
	responder.playbook.position = 0
	responder.playbook.count = scenario.Count
	responder.playbook.servedDataframe7d = false
	responder.playbook.servedDataframe80 = false

	// iterate the scenario extracting the raw dataframes into a sequential playbook
	for i := 0; i < scenario.Count; i++ {
		pr := PlaybookResponse{}
		pr.dataframe7d = responder.convertHexStringToByteArray(scenario.Rawdata[i].Dataframe7d)
		pr.dataframe80 = responder.convertHexStringToByteArray(scenario.Rawdata[i].Dataframe80)

		responder.playbook.responses = append(responder.playbook.responses, pr)
	}
}

// GetECUResponse returns an emulated response byte string
func (responder *Responder) GetECUResponse(cmd []byte) []byte {
	var data []byte

	// convert the command code to a string
	command := hex.EncodeToString(cmd)
	command = strings.ToUpper(command)

	// if the command is a dataframe request and we have a response file
	// then use the response file
	if responder.isDataframeRequest(command) {

		position := responder.playbook.position

		if command == "7D" {
			data = responder.playbook.responses[position].dataframe7d
			// truncate to the right size
			data = data[:33]
			responder.playbook.servedDataframe7d = true
		}

		if command == "80" {
			data = responder.playbook.responses[position].dataframe80
			// truncate to the right size
			data = data[:29]
			responder.playbook.servedDataframe80 = true
		}

		// served both dataframes from this position, index on to the next position
		if responder.playbook.servedDataframe7d && responder.playbook.servedDataframe80 {
			responder.playbook.servedDataframe7d = false
			responder.playbook.servedDataframe80 = false

			responder.playbook.position = responder.playbook.position + 1
			utils.LogI.Printf("both dataframes served, indexing to %d of %d", responder.playbook.position, responder.playbook.count)

			// if we've reached the end then loop back to the start
			if responder.playbook.position > responder.playbook.count {
				responder.playbook.position = 0
				utils.LogW.Printf("reached end of scenario, restarting from beginning")
			}
		}
	} else {
		// generate the relevant response
		data = responder.generateECUResponse(command)
	}

	return data
}

// determines where the command code is a dataframe request
func (responder *Responder) isDataframeRequest(command string) bool {
	return (command == "80" || command == "7D")
}

// converts the hex string to a byte array
func (responder *Responder) convertHexStringToByteArray(response string) []byte {
	// convert to byte array
	data, _ := hex.DecodeString(response)

	return data
}

// if we're responding to a command that isn't a dataframe request
// then generate the correct response
func (responder *Responder) generateECUResponse(command string) []byte {
	command = strings.ToUpper(command)

	r := responder.responseMap[command]

	if r == nil {
		r = responder.responseMap["00"]
		copy(r[0:], command)
	}

	utils.LogI.Printf("generating response %x for %s", r, command)
	return r
}

// build the response map for generated responses
func (responder *Responder) buildResponseMap() {
	// Response formats for commands that do not respond with the format [COMMAND][VALUE]
	// Generally these are either part of the initialisation sequence or are ECU data frames
	responder.responseMap["0A"] = []byte{0x0A}
	responder.responseMap["CA"] = []byte{0xCA}
	responder.responseMap["75"] = []byte{0x75}

	// Format for DataFrames starts with [Command Echo][Data Size][Data Bytes (28 for 0x80 and 32 for 0x7D)]
	responder.responseMap["80"] = []byte{0x80, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B}
	responder.responseMap["7D"] = []byte{0x7d, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F}
	responder.responseMap["D0"] = []byte{0xD0, 0x99, 0x00, 0x03, 0x03}

	// heatbeat
	responder.responseMap["F4"] = []byte{0xf4, 0x00}

	// adjustments
	responder.responseMap["7A"] = []byte{0x7a, 0x8a}
	responder.responseMap["7B"] = []byte{0x7b, 0x1e}
	responder.responseMap["7C"] = []byte{0x7c, 0x8a}
	responder.responseMap["79"] = []byte{0x79, 0x8a}
	responder.responseMap["7B"] = []byte{0x7b, 0x8a}
	responder.responseMap["8A"] = []byte{0x8a, 0x23}
	responder.responseMap["89"] = []byte{0x89, 0x23}
	responder.responseMap["92"] = []byte{0x92, 0x80}
	responder.responseMap["91"] = []byte{0x91, 0x80}
	responder.responseMap["94"] = []byte{0x94, 0x80}
	responder.responseMap["93"] = []byte{0x93, 0x80}

	//resets
	responder.responseMap["FA"] = []byte{0xfa, 0x00}
	responder.responseMap["0F"] = []byte{0x0f, 0x00}
	responder.responseMap["CC"] = []byte{0xcc, 0x00}

	// generic response, expect command and single byte response
	responder.responseMap["00"] = []byte{0x00, 0x00}
}
