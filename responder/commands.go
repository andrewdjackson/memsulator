package responder

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

var responseMap = make(map[string][]byte)

// package init function
func init() {
	// Response formats for commands that do not respond with the format [COMMAND][VALUE]
	// Generally these are either part of the initialisation sequence or are ECU data frames
	responseMap["0A"] = []byte{0x0A}
	responseMap["CA"] = []byte{0xCA}
	responseMap["75"] = []byte{0x75}

	// Format for DataFrames starts with [Command Echo][Data Size][Data Bytes (28 for 0x80 and 32 for 0x7D)]
	responseMap["80"] = []byte{0x80, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B}
	responseMap["7D"] = []byte{0x7d, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F}
	responseMap["D0"] = []byte{0xD0, 0x99, 0x00, 0x02, 0x03}
	responseMap["D1"] = []byte{0xD1, 0x41, 0x42, 0x4E, 0x4D, 0x50, 0x30, 0x30, 0x32, 0x99, 0x00, 0x02, 0x03, 0x41, 0x42}

	// heartbeat
	responseMap["F4"] = []byte{0xf4, 0x00}

	// adjustments
	responseMap["79"] = []byte{0x79, 0x8b} // increment STFT (default is 138)
	responseMap["7A"] = []byte{0x7a, 0x89} // decrement STFT (default is 138)
	responseMap["7B"] = []byte{0x7b, 0x1f} // increment LTFT (default is 30)
	responseMap["7C"] = []byte{0x7c, 0x1d} // decrement LTFT (default is 30)
	responseMap["89"] = []byte{0x89, 0x24} // increment Idle Decay (default is 35)
	responseMap["8A"] = []byte{0x8a, 0x22} // decrement Idle Decay (default is 35)
	responseMap["91"] = []byte{0x91, 0x81} // increment Idle Speed  (default is 128)
	responseMap["92"] = []byte{0x92, 0x7f} // decrement Idle Speed (default is 128)
	responseMap["93"] = []byte{0x93, 0x81} // increment Ignition Advance Offset (default is 128)
	responseMap["94"] = []byte{0x94, 0x7f} // decrement Ignition Advance Offset (default is 128)
	responseMap["FD"] = []byte{0xfd, 0x81} // increment IAC (default is 128)
	responseMap["FE"] = []byte{0xfe, 0x7f} // decrement IAC (default is 128)

	//resets
	responseMap["0F"] = []byte{0x0f, 0x00} // clear all adjustments
	responseMap["CC"] = []byte{0xcc, 0x00} // clear faults
	responseMap["FA"] = []byte{0xfa, 0x00} // clear all computed and learnt settings
	responseMap["FB"] = []byte{0xfb, 0x80} // Idle Air Control Position

	// actuators
	responseMap["10"] = []byte{0x10, 0x00} // temperature gauge on
	responseMap["00"] = []byte{0x00, 0x00} // temperature gauge off
	responseMap["11"] = []byte{0x11, 0x00} // fuel pump on
	responseMap["01"] = []byte{0x01, 0x00} // fuel pump off
	responseMap["12"] = []byte{0x12, 0x00} // ptc relay on
	responseMap["02"] = []byte{0x02, 0x00} // ptc relay off
	responseMap["13"] = []byte{0x13, 0x00} // ac relay on
	responseMap["03"] = []byte{0x03, 0x00} // ac relay off
	responseMap["18"] = []byte{0x18, 0x00} // purge valve on
	responseMap["08"] = []byte{0x08, 0x00} // purge vavle off
	responseMap["19"] = []byte{0x19, 0x00} // O2 heater on
	responseMap["09"] = []byte{0x09, 0x00} // O2 heater off
	responseMap["1B"] = []byte{0x1b, 0x00} // boost valve on
	responseMap["0B"] = []byte{0x0b, 0x00} // boost valve off
	responseMap["1D"] = []byte{0x1d, 0x00} // fan 1 on
	responseMap["0D"] = []byte{0x0d, 0x00} // fan 1 off
	responseMap["1E"] = []byte{0x1e, 0x00} // fan 2 on
	responseMap["0E"] = []byte{0x0e, 0x00} // fan 2 off
	responseMap["EF"] = []byte{0xef, 0x03} // test mpi injectors
	responseMap["F7"] = []byte{0xf7, 0x03} // test injectors
	responseMap["F8"] = []byte{0xf8, 0x02} // fire coil

	// unknown command Responses
	responseMap["65"] = []byte{0x65, 0x00}
	responseMap["6D"] = []byte{0x6d, 0x00}
	responseMap["7E"] = []byte{0x7e, 0x08}
	responseMap["7F"] = []byte{0x7f, 0x05}
	responseMap["82"] = []byte{0x82, 0x09, 0x9E, 0x1D, 0x00, 0x00, 0x60, 0x05, 0xFF, 0xFF}
	responseMap["CB"] = []byte{0xcb, 0x00}
	responseMap["CD"] = []byte{0xcd, 0x01}
	responseMap["D2"] = []byte{0xd2, 0x02, 0x01, 0x00, 0x01}
	responseMap["D3"] = []byte{0xd3, 0x02, 0x01, 0x00, 0x02}
	responseMap["E7"] = []byte{0xe7, 0x02}
	responseMap["E8"] = []byte{0xe8, 0x05, 0x26, 0x01, 0x00, 0x01}
	responseMap["ED"] = []byte{0xed, 0x00}
	responseMap["EE"] = []byte{0xee, 0x00}
	responseMap["F0"] = []byte{0xf0, 0x05}
	responseMap["F3"] = []byte{0xf3, 0x00}
	responseMap["F5"] = []byte{0xf5, 0x00}
	responseMap["F6"] = []byte{0xf6, 0x00}
	responseMap["FC"] = []byte{0xfc, 0x00}
}

func getResponse(command []byte) []byte {
	c := fmt.Sprintf("%X", command)
	if r := responseMap[c]; r == nil {
		// no mapped response, echo the command
		log.Warnf("unable to find command %s", c)
		return []byte{command[0], 0x00}
	} else {
		return r
	}
}
