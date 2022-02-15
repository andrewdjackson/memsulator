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
	responder.responseMap["D1"] = []byte{0xD1, 0x41, 0x42, 0x4E, 0x4D, 0x50, 0x30, 0x30, 0x33, 0x99, 0x00, 0x03, 0x03}

	// heatbeat
	responder.responseMap["F4"] = []byte{0xf4, 0x00}

	// adjustments
	responder.responseMap["79"] = []byte{0x79, 0x8b} // increment STFT (default is 138)
	responder.responseMap["7A"] = []byte{0x7a, 0x89} // decrement STFT (default is 138)
	responder.responseMap["7B"] = []byte{0x7b, 0x1f} // increment LTFT (default is 30)
	responder.responseMap["7C"] = []byte{0x7c, 0x1d} // decrement LTFT (default is 30)
	responder.responseMap["89"] = []byte{0x89, 0x24} // increment Idle Decay (default is 35)
	responder.responseMap["8A"] = []byte{0x8a, 0x22} // decrement Idle Decay (default is 35)
	responder.responseMap["91"] = []byte{0x91, 0x81} // increment Idle Speed  (default is 128)
	responder.responseMap["92"] = []byte{0x92, 0x7f} // decrement Idle Speed (default is 128)
	responder.responseMap["93"] = []byte{0x93, 0x81} // increment Ignition Advance Offset (default is 128)
	responder.responseMap["94"] = []byte{0x94, 0x7f} // decrement Ignition Advance Offset (default is 128)
	responder.responseMap["FD"] = []byte{0xfd, 0x81} // increment IAC (default is 128)
	responder.responseMap["FE"] = []byte{0xfe, 0x7f} // decrement IAC (default is 128)

	//resets
	responder.responseMap["0F"] = []byte{0x0f, 0x00} // clear all adjustments
	responder.responseMap["CC"] = []byte{0xcc, 0x00} // clear faults
	responder.responseMap["FA"] = []byte{0xfa, 0x00} // clear all computed and learnt settings
	responder.responseMap["FB"] = []byte{0xfb, 0x80} // Idle Air Control position

	// actuators
	responder.responseMap["11"] = []byte{0x11, 0x00} // fuel pump on
	responder.responseMap["01"] = []byte{0x01, 0x00} // fuel pump off
	responder.responseMap["12"] = []byte{0x12, 0x00} // ptc relay on
	responder.responseMap["02"] = []byte{0x02, 0x00} // ptc relay off
	responder.responseMap["13"] = []byte{0x13, 0x00} // ac relay on
	responder.responseMap["03"] = []byte{0x03, 0x00} // ac relay off
	responder.responseMap["18"] = []byte{0x18, 0x00} // purge valve on
	responder.responseMap["08"] = []byte{0x08, 0x00} // purge vavle off
	responder.responseMap["19"] = []byte{0x19, 0x00} // O2 heater on
	responder.responseMap["09"] = []byte{0x09, 0x00} // O2 heater off
	responder.responseMap["1B"] = []byte{0x1b, 0x00} // boost valve on
	responder.responseMap["0B"] = []byte{0x0b, 0x00} // boost valve off
	responder.responseMap["1D"] = []byte{0x1d}       // fan 1 on
	responder.responseMap["0D"] = []byte{0x0d, 0x00} // fan 1 off
	responder.responseMap["1E"] = []byte{0x1e}       // fan 2 on
	responder.responseMap["0E"] = []byte{0x0e, 0x00} // fan 2 off
	responder.responseMap["EF"] = []byte{0xef, 0x03} // test mpi injectors
	responder.responseMap["F7"] = []byte{0xf7, 0x03} // test injectors
	responder.responseMap["F8"] = []byte{0xf8, 0x02} // fire coil

	// unknown command responses
	responder.responseMap["65"] = []byte{0x65, 0x00}
	responder.responseMap["6D"] = []byte{0x6d, 0x00}
	responder.responseMap["7E"] = []byte{0x7e, 0x08}
	responder.responseMap["7F"] = []byte{0x7f, 0x05}
	responder.responseMap["82"] = []byte{0x82, 0x09, 0x9E, 0x1D, 0x00, 0x00, 0x60, 0x05, 0xFF, 0xFF}
	responder.responseMap["CB"] = []byte{0xcb, 0x00}
	responder.responseMap["CD"] = []byte{0xcd, 0x01}
	responder.responseMap["D2"] = []byte{0xd2, 0x02, 0x01, 0x00, 0x01}
	responder.responseMap["D3"] = []byte{0xd3, 0x02, 0x01, 0x00, 0x02}
	responder.responseMap["E7"] = []byte{0xe7, 0x02}
	responder.responseMap["E8"] = []byte{0xe8, 0x05, 0x26, 0x01, 0x00, 0x01}
	responder.responseMap["ED"] = []byte{0xed, 0x00}
	responder.responseMap["EE"] = []byte{0xee, 0x00}
	responder.responseMap["F0"] = []byte{0xf0, 0x05}
	responder.responseMap["F3"] = []byte{0xf3, 0x00}
	responder.responseMap["F5"] = []byte{0xf5, 0x00}
	responder.responseMap["F6"] = []byte{0xf6, 0x00}
	responder.responseMap["FC"] = []byte{0xfc, 0x00}

	// generic response, expect command and single byte response
	responder.responseMap["00"] = []byte{0x00, 0x00}
}
