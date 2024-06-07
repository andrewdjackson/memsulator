package loader

import (
	"fmt"
	"testing"
)

func TestPlayback(t *testing.T) {
	scenario := NewScenario()
	if err := scenario.Load("../scenarios/thyberg.csv"); err != nil {
		t.Errorf("%s", err)
	} else {
		if len(scenario.dataframes) == 0 {
			t.Errorf("no dataframes in scenario")
		}
	}

	playback := NewPlayback(scenario)
	if playback.position != 0 {
		t.Errorf("not initialised")
	}

	var response []byte

	response = playback.NextDataframe([]byte{0x80})
	fmt.Printf("%X\n", response)
	if playback.position != 0 {
		t.Errorf("not incremented")
	}
	if response[0] != 0x80 {
		t.Errorf("not dataframe")
	}

	response = playback.NextDataframe([]byte{0x7D})
	fmt.Printf("%X\n", response)
	if playback.position != 0 {
		t.Errorf("not incremented")
	}
	if response[0] != 0x7D {
		t.Errorf("not dataframe")
	}

	response = playback.NextDataframe([]byte{0x80})
	fmt.Printf("%X\n", response)
	if playback.position != 1 {
		t.Errorf("not incremented")
	}
	if response[0] != 0x80 {
		t.Errorf("not dataframe")
	}

	response = playback.NextDataframe([]byte{0x7D})
	fmt.Printf("%X\n", response)
	if playback.position != 1 {
		t.Errorf("not incremented")
	}
	if response[0] != 0x7D {
		t.Errorf("not dataframe")
	}
}
