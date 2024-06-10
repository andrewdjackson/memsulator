package responder

import (
	"github.com/andrewdjackson/memsulator/loader"
	"github.com/stretchr/testify/assert"
	"testing"
)

const TEST_FILE = "../scenarios/test.csv"

func TestSendResponse(t *testing.T) {
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
	}
}
