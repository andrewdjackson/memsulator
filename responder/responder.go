package responder

import (
	"github.com/andrewdjackson/memsulator/loader"
	"github.com/andrewdjackson/memsulator/utils"
	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
	"time"
)

type Emulator struct {
	connected   bool
	initialised bool
	serialPort  serial.Port
	playback    *loader.Playback
}

func NewEmulator(playback *loader.Playback) *Emulator {
	emulator := &Emulator{
		playback:  playback,
		connected: false,
	}

	return emulator
}

func (emulator *Emulator) Connect(port string) (bool, error) {
	var err error

	emulator.connected = true

	if err = emulator.connectToSerialPort(port); err != nil {
		emulator.connected = false
	}

	return emulator.connected, err
}

// Listen for commands sent to the ECU
// and send back the ECU response
// this is an infinte loop, so needs to be executed as a go process
func (emulator *Emulator) Listen() {
	for {
		command := emulator.readSerial()

		if len(command) > 0 {
			log.Infof("command  %X", command)

			response := emulator.sendResponse(command)
			log.Infof("response %X", response)
		}
	}
}

func (emulator *Emulator) Disconnect() (bool, error) {
	var err error

	if err = emulator.serialPort.Close(); err == nil {
		emulator.connected = false
	}

	return emulator.connected, err
}

func (emulator *Emulator) sendResponse(command []byte) []byte {
	var response []byte

	if command[0] == 0x80 || command[0] == 0x7D {
		response = emulator.playback.NextDataframe(command)
	} else {
		// get appropriate response
		response = getResponse(command)
		log.Infof("found response %X", response)
	}

	return response
}

func (emulator *Emulator) connectToSerialPort(port string) error {
	var err error

	log.Infof("attempting to open serial serialPort %s", port)

	mode := serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	// connect to the ecu
	if emulator.serialPort, err = serial.Open(port, &mode); err != nil {
		log.Errorf("error opening serial port (%s)", err)
	}

	return err
}

func (emulator *Emulator) readSerial() []byte {
	var n int
	var e error

	// serial read buffer
	b := make([]byte, 1)

	//  data frame buffer
	data := make([]byte, 0)

	if emulator.serialPort != nil {
		// wait for a response from MEMS
		n, e = emulator.serialPort.Read(b)

		if e != nil {
			log.Infof("%s error %s", utils.ECUCommandTrace, e)
		} else {
			// append the read bytes to the data frame
			data = append(data, b[:n]...)
		}
	}

	if n > 0 {
		log.Infof("serial read %x (%d)", data, n)
	}

	return data
}

func (emulator *Emulator) writeSerial(data []byte) {
	if emulator.serialPort != nil {
		n, e := emulator.serialPort.Write(data)

		if e != nil {
			log.Infof("send Error %s", e)
		}

		if n > 0 {
			log.Infof("%x (%d)", data, n)
		}
	}
}

func serialWait() {
	time.Sleep(time.Duration(75) * time.Millisecond)
}
