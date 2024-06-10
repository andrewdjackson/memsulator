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

	if playback.dataframe80Sent && playback.dataframe7dSent {
		playback.position += 1
		playback.dataframe80Sent = false
		playback.dataframe7dSent = false
	}

	if playback.position >= playback.scenario.Count {
		log.Infof("reached end of scenario, restarting from beginning")
		playback.Start()
	}

	if command[0] == 0x80 {
		playback.dataframe80Sent = true
		dataframe = playback.scenario.dataframes[playback.position].Dataframe80
	}

	if command[0] == 0x7d {
		playback.dataframe7dSent = true
		dataframe = playback.scenario.dataframes[playback.position].Dataframe7d
	}

	return dataframe
}
