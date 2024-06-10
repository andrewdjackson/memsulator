package loader

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlayback(t *testing.T) {
	scenario := NewScenario()
	if err := scenario.Load(TEST_FILE); err != nil {
		t.Errorf("%s", err)
	} else {
		assert.Len(t, scenario.dataframes, DATAFRAME_COUNT)

		playback := NewPlayback(scenario)
		assert.Equal(t, playback.position, 0)

		var response []byte

		for i := 0; i < DATAFRAME_COUNT; i += 1 {
			response = playback.NextDataframe([]byte{0x80})
			fmt.Printf("%X\n", response)

			assert.Equal(t, playback.position, i)
			assert.Equal(t, uint8(response[0]), uint8(0x80))
			assert.Equal(t, uint8(response[len(response)-1]), uint8(i))

			response = playback.NextDataframe([]byte{0x7D})
			fmt.Printf("%X\n", response)

			assert.Equal(t, playback.position, i)
			assert.Equal(t, uint8(response[0]), uint8(0x7d))
		}

		// test loop back to start
		response = playback.NextDataframe([]byte{0x80})
		fmt.Printf("%X\n", response)

		assert.Equal(t, playback.position, 0)
		assert.Equal(t, uint8(response[0]), uint8(0x80))
		assert.Equal(t, uint8(response[len(response)-1]), uint8(0))
	}
}
