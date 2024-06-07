package loader

import (
	"testing"
)

func TestScenario(t *testing.T) {
	scenario := NewScenario()
	scenario.Load("../scenarios/thyberg.csv")

	t.Errorf("it went wrong")
}
