package responder

import (
	"github.com/andrewdjackson/memsulator/loader"
	"github.com/andrewdjackson/memsulator/utils"
	log "github.com/sirupsen/logrus"
	"go.bug.st/serial.v1"
	"time"
)

type Emulator struct {
	connected   bool
	initialised bool
	serialPort  serial.Port
	playback    *loader.Playback
}

func NewEmulator(playback *loader.Playback) *Emulator {
	return &Emulator{
		playback:  playback,
		connected: false,
	}
}

func (emulator *Emulator) Connect(port string) (bool, error) {
	var err error

	emulator.connected = true

	if err = emulator.connectToSerialPort(port); err != nil {
		emulator.connected = false
	}

	return emulator.connected, err
}

func (emulator *Emulator) ReceiveAndSend() ([]byte, error) {
	var err error

	command := emulator.waitForCommand()
	log.Infof("%X", command)

	if command[0] == 0xCA {
		// initialisation required
		emulator.initialise()
	}

	response := emulator.sendResponse(command)

	return response, err
}

func (emulator *Emulator) Disconnect() (bool, error) {
	var err error

	if err = emulator.serialPort.Close(); err == nil {
		emulator.connected = false
	}

	return emulator.connected, err
}

func (emulator *Emulator) initialise() bool {
	var command []byte

	// CA <- CA
	command = []byte{0xCA}
	emulator.writeSerial(command)
	// 75 <- 75
	command = emulator.readSerialByte()
	emulator.writeSerial(command)
	// F4 <- F4
	command = emulator.readSerialByte()
	emulator.writeSerial(command)
	// D0 <- D099000203
	command = emulator.readSerialByte()
	emulator.writeSerial([]byte{0xD0, 0x99, 0x00, 0x02, 0x03})

	return true
}

func (emulator *Emulator) sendResponse(command []byte) []byte {
	var response []byte

	if command[0] == 0x80 || command[0] == 0x7D {
		response = emulator.playback.NextDataframe(command)
	} else {
		// get appropriate response
		response = []byte{command[0], 0x00}
	}

	return response
}

func (emulator *Emulator) waitForCommand() []byte {
	var command []byte

	for {
		command = emulator.readSerialByte()
		break
	}

	return command
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

func (emulator *Emulator) readSerialByte() []byte {
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

func (emulator *Emulator) readSerial() []byte {
	var n int
	var err error

	// serial read buffer
	b := make([]byte, 1)

	//  data frame buffer
	data := make([]byte, 0)

	if emulator.serialPort != nil {
		n, err = emulator.serialPort.Read(b)

		if err != nil {
			log.Infof("error %s", err)
		} else {
			// append the read bytes to the data frame
			data = append(data, b[:n]...)
		}
	}

	if n > 0 {
		log.Infof("%x (%d)", data, n)
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
