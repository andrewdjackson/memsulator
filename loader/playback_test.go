package loader

import (
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

		for i := 0; i < DATAFRAME_COUNT-1; i += 1 {
			response = playback.NextDataframe([]byte{0x80})

			assert.Equal(t, playback.position, i)
			assert.Equal(t, uint8(response[0]), uint8(0x80))
			assert.Equal(t, uint8(response[len(response)-1]), uint8(i))

			response = playback.NextDataframe([]byte{0x7D})

			assert.Equal(t, playback.position, i)
			assert.Equal(t, uint8(response[0]), uint8(0x7d))
		}
	}
}

func TestPlaybackLoop(t *testing.T) {
	scenario := NewScenario()
	if err := scenario.Load(TEST_FILE); err != nil {
		t.Errorf("%s", err)
	} else {
		assert.Len(t, scenario.dataframes, DATAFRAME_COUNT)

		playback := NewPlayback(scenario)
		var response []byte
		var expectedPosition int
		loopTwiceCount := playback.scenario.Count * 2

		for i := 0; i < loopTwiceCount; i += 1 {
			if i >= playback.scenario.Count {
				expectedPosition = i - playback.scenario.Count
			} else {
				expectedPosition = i
			}

			response = playback.NextDataframe([]byte{0x80})
			assert.Equal(t, uint8(response[0]), uint8(0x80))
			response = playback.NextDataframe([]byte{0x7D})
			assert.Equal(t, uint8(response[0]), uint8(0x7d))
			assert.Equal(t, playback.position, expectedPosition)
		}
	}
}
