package ecu

import (
	"bufio"
	"encoding/hex"
	"github.com/andrewdjackson/memsulator/scenarios"
	"github.com/andrewdjackson/memsulator/utils"
	"go.bug.st/serial"
	"strings"
)

// MemsConnection communication structure for MEMS
type ECU16 struct {
	// SerialPort the serial connection
	SerialPort      serial.Port
	portReader      *bufio.Reader
	ECUID           []byte
	command         []byte
	response        []byte
	SendToFCR       chan CommandResponse
	ReceivedFromFCR chan CommandResponse
	Connected       bool
	Initialised     bool
	Exit            bool
	responder       *Responder
	MemsVersion     string
}

// NewMemsConnection creates a new mems structure
func NewECU16() *ECU16 {
	m := &ECU16{}
	m.Connected = false
	m.Initialised = false
	m.SendToFCR = make(chan CommandResponse)
	m.ReceivedFromFCR = make(chan CommandResponse)

	return m
}

// Open communication via serial port
func (mems *ECU16) Open(port string) {
	// connect to the ecu
	mode := serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	utils.LogI.Println("Opening ", port)

	s, err := serial.Open(port, &mode)
	//s, err := serial.OpenPort(c)
	if err != nil {
		utils.LogI.Printf("%s", err)
		mems.Connected = false
	} else {
		utils.LogI.Println("Listening on ", port)

		mems.SerialPort = s
		mems.Connected = true
	}
}

// LoadScenario the emulation scenario
func (mems *ECU16) LoadScenario(scenario *scenarios.Scenario) {
	mems.responder = NewResponder()
	mems.responder.LoadScenario(scenario)
	utils.LogI.Printf("loaded scenario")
}

// Listen listens for commands from the FCR
func (mems *ECU16) Listen() {
	var cr CommandResponse

	for {
		// read from the serial port
		cmd := mems.readSerial()

		// if we get a command the send a response
		if len(cmd) > 0 {
			// find the command response
			response := mems.responder.GetECUResponse(cmd)

			if len(response) > 0 {
				cr.Command = cmd
				cr.Response = response

				mems.sendResponse(cr)

				// send the command / response over the channel
				mems.ReceivedFromFCR <- cr
			} else {
				utils.LogI.Printf("%s unexpected generated response for %x", utils.ECUCommandTrace, cmd)
			}
		}
	}
}

func (mems *ECU16) sendResponse(cr CommandResponse) {
	// ignore 7D requests if the MEMS is Version 1.3
	if mems.MemsVersion == "1.3" {
		cmd := hex.EncodeToString(cr.Command)
		if strings.ToUpper(cmd) == "7D" {
			utils.LogI.Printf("0x7d command ignored by MEMS 1.3")
			return
		}
	}

	// send the response to the FCR
	mems.writeSerial(cr.Response)
}

// readSerial read command sent from FCR
// read 1 byte at a time until we have all the expected bytes
func (mems *ECU16) readSerial() []byte {
	var n int
	var e error

	//size := mems.getResponseSize(mems.command)
	size := 1

	// serial read buffer
	b := make([]byte, 1)

	//  data frame buffer
	data := make([]byte, 0)

	if mems.SerialPort != nil {

		// read all the expected bytes before returning the data
		for count := 0; count < size; {
			// wait for a response from MEMS
			n, e = mems.SerialPort.Read(b)

			if e != nil {
				utils.LogI.Printf("%s error %s", utils.ECUCommandTrace, e)
			} else {
				// append the read bytes to the data frame
				data = append(data, b[:n]...)
			}

			// increment by the number of bytes read
			count = count + n
			if count > size {
				utils.LogI.Printf("%s data frame size mismatch (received %d, expected %d)", utils.ECUCommandTrace, count, size)
			}
		}
	}

	if n > 0 {
		utils.LogI.Printf("%s %x (%d)", utils.ECUCommandTrace, data, n)
		mems.command = data
	}

	return data
}

// writeSerial write response from ECU to FCR
func (mems *ECU16) writeSerial(data []byte) {
	if mems.SerialPort != nil {
		// save the sent response
		mems.response = data

		// write the response to the code reader
		n, e := mems.SerialPort.Write(data)

		if e != nil {
			utils.LogI.Printf("%s send Error %s", utils.ECUResponseTrace, e)
		}

		if n > 0 {
			utils.LogI.Printf("%s %x (%d)", utils.ECUResponseTrace, data, n)
		}
	}
}
