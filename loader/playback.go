package loader

import (
	log "github.com/sirupsen/logrus"
)

type Playback struct {
	scenario        *Scenario
	position        int
	dataframe80Sent bool
	dataframe7dSent bool
}

func NewPlayback(scenario *Scenario) *Playback {
	return &Playback{
		scenario: scenario,
		position: 0,
	}
}

func (playback *Playback) Start() {
	playback.position = 0
	playback.dataframe80Sent = false
	playback.dataframe7dSent = false
}

func (playback *Playback) PrevDataframe() *Dataframes {
	data := playback.scenario.dataframes[playback.position]
	playback.position -= 1

	if playback.position < 0 {
		playback.Start()
	}

	return data
}

func (playback *Playback) NextDataframe(command []byte) []byte {
	var dataframe []byte

	playback.incrementPosition()

	scenarioLength := playback.scenario.Count - 1

	if playback.position > scenarioLength {
		log.Infof("reached end of scenario (%d/%d), restarting from beginning", playback.position+1, playback.scenario.Count)
		playback.Start()
	}

	if playback.isValidDataframe() {
		if command[0] == 0x80 {
			playback.dataframe80Sent = true
			dataframe = playback.scenario.dataframes[playback.position].Dataframe80
		}

		if command[0] == 0x7d {
			playback.dataframe7dSent = true
			dataframe = playback.scenario.dataframes[playback.position].Dataframe7d
		}

		log.Infof("playback %X (%d/%d)", dataframe, playback.position+1, playback.scenario.Count)
	} else {
		// skip dataframe
		log.Warnf("playback dataframe invalid at %d", playback.position)

		playback.dataframe80Sent = true
		playback.dataframe7dSent = true

		dataframe = playback.NextDataframe(command)
	}

	return dataframe
}

func (playback *Playback) incrementPosition() {
	if playback.dataframe80Sent && playback.dataframe7dSent {
		playback.position += 1
		playback.dataframe80Sent = false
		playback.dataframe7dSent = false
	}
}

func (playback *Playback) isValidDataframe() bool {
	return len(playback.scenario.dataframes[playback.position].Dataframe80) == 29 && len(playback.scenario.dataframes[playback.position].Dataframe7d) == 33
}
