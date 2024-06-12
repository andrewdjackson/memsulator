package responder

import (
	"github.com/andrewdjackson/memsulator/loader"
	"github.com/stretchr/testify/assert"
	"testing"
)

const TEST_FILE = "../scenarios/test.csv"

func TestEmulator_sendResponse(t *testing.T) {
	scenario := loader.NewScenario()
	if err := scenario.Load(TEST_FILE); err != nil {
		t.Errorf("%s", err)
	} else {
		playback := loader.NewPlayback(scenario)
		emulator := NewEmulator(playback)

		response := emulator.sendResponse([]byte{0x80})
		assert.Equal(t, uint8(response[0]), uint8(0x80))
		assert.Equal(t, uint8(response[len(response)-1]), uint8(0))

		response = emulator.sendResponse([]byte{0x7d})
		assert.Equal(t, uint8(response[0]), uint8(0x7d))
		assert.Equal(t, uint8(response[len(response)-1]), uint8(0))

		// initialisation byte
		response = emulator.sendResponse([]byte{0xCA})
		assert.Equal(t, uint8(response[0]), uint8(0xCA))
		assert.Len(t, response, 1)
	}
}

func TestEmulator_Connect(t *testing.T) {
	var err error
	var connected bool

	vport := NewVirtualSerialPort()
	err = vport.CreateVirtualPorts()
	assert.Equal(t, err, nil)

	scenario := loader.NewScenario()
	playback := loader.NewPlayback(scenario)
	emulator := NewEmulator(playback)

	connected, err = emulator.Connect(vport.fcrPort)
	assert.Equal(t, err, nil)
	assert.True(t, connected)

	connected, err = emulator.Disconnect()
	assert.Equal(t, err, nil)
	assert.False(t, connected)
}

func TestEmulator_ReceiveAndSend(t *testing.T) {
	var err error
	var connected bool

	vport := NewVirtualSerialPort()
	err = vport.CreateVirtualPorts()
	assert.Equal(t, err, nil)

	scenario := loader.NewScenario()
	playback := loader.NewPlayback(scenario)
	emulator := NewEmulator(playback)

	connected, err = emulator.Connect(vport.fcrPort)
	assert.Equal(t, err, nil)
	assert.True(t, connected)

	connected, err = emulator.Disconnect()
	assert.Equal(t, err, nil)
	assert.False(t, connected)
}
