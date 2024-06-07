package loader

import (
	"testing"
)

func TestLoadScenario(t *testing.T) {
	scenario := NewScenario()
	if err := scenario.Load("../scenarios/thyberg.csv"); err != nil {
		t.Errorf("%s", err)
	} else {
		if len(scenario.dataframes) == 0 {
			t.Errorf("no dataframes in scenario")
		}
	}
}
