package ecu

import (
	"bufio"

	"github.com/andrewdjackson/memsulator/scenarios"

	"github.com/andrewdjackson/memsulator/utils"
	"go.bug.st/serial.v1"
)

// MemsCommandResponse communication pair
type MemsCommandResponse struct {
	Command  []byte
	Response []byte
}

// MemsConnection communtication structure for MEMS
type MemsConnection struct {
	// SerialPort the serial connection
	SerialPort      serial.Port
	portReader      *bufio.Reader
	ECUID           []byte
	command         []byte
	response        []byte
	SendToFCR       chan MemsCommandResponse
	ReceivedFromFCR chan MemsCommandResponse
	Connected       bool
	Initialised     bool
	Exit            bool
	responder       *Responder
}

// NewMemsConnection creates a new mems structure
func NewMemsConnection() *MemsConnection {
	m := &MemsConnection{}
	m.Connected = false
	m.Initialised = false
	m.SendToFCR = make(chan MemsCommandResponse)
	m.ReceivedFromFCR = make(chan MemsCommandResponse)

	return m
}

// Open communiction via serial port
func (mems *MemsConnection) Open(port string) {
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
	} else {
		utils.LogI.Println("Listening on ", port)

		mems.SerialPort = s
		mems.Connected = true
	}
}

// LoadScenario the emulation scenario
func (mems *MemsConnection) LoadScenario(scenario *scenarios.Scenario) {
	mems.responder = NewResponder()
	mems.responder.LoadScenario(scenario)
	utils.LogI.Printf("loaded scenario")
}

// ListenToFCRLoop listens for commands from the FCR
func (mems *MemsConnection) ListenToFCRLoop() {
	var cr MemsCommandResponse

	for {
		// read from the serial port
		cmd := mems.readSerial()

		// if we get a command the send a response
		if len(cmd) > 0 {
			// find the command response
			response := mems.responder.GetECUResponse(cmd)
			//response := mems.Response(cmd)

			if len(response) > 0 {
				// send the response to the FCR
				mems.writeSerial(response)

				cr.Command = cmd
				cr.Response = response

				// send the command / response over the channel
				mems.ReceivedFromFCR <- cr
			} else {
				utils.LogI.Printf("%s unexpected generated response for %x", utils.ECUCommandTrace, cmd)
			}
		}
	}
}

// readSerial read command sent from FCR
// read 1 byte at a time until we have all the expected bytes
func (mems *MemsConnection) readSerial() []byte {
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
func (mems *MemsConnection) writeSerial(data []byte) {
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
