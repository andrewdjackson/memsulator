package loader

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const TEST_FILE = "../scenarios/test.csv"
const DATAFRAME_COUNT = 8

func TestLoadScenario(t *testing.T) {
	scenario := NewScenario()
	if err := scenario.Load(TEST_FILE); err != nil {
		t.Errorf("%s", err)
	} else {
		assert.Len(t, scenario.dataframes, DATAFRAME_COUNT)

		for i := 0; i < DATAFRAME_COUNT; i += 1 {
			df := scenario.dataframes[i]
			assert.Equal(t, uint8(df.Dataframe80[0]), uint8(0x80))
			assert.Equal(t, uint8(df.Dataframe7d[0]), uint8(0x7d))

			// validate last byte of dataframe is sequential for tests
			assert.Equal(t, uint8(df.Dataframe80[len(df.Dataframe80)-1]), uint8(i))
			assert.Equal(t, uint8(df.Dataframe7d[len(df.Dataframe7d)-1]), uint8(i))
		}
	}
}
