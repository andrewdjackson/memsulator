package main

import (
	"flag"
	"github.com/andrewdjackson/memsulator/loader"
	"github.com/andrewdjackson/memsulator/responder"
	log "github.com/sirupsen/logrus"
)

func main() {
	file := flag.String("scenario", "scenarios/default.csv", "scenario file to run")
	port := flag.String("port", "/dev/serial0", "serial communication port")
	flag.Parse()

	scenario := loader.NewScenario()
	scenario.Load(*file)

	playback := loader.NewPlayback(scenario)
	playback.Start()

	emulator := responder.NewEmulator(playback)
	if connected, err := emulator.Connect(*port); err == nil {
		if connected {
			for {
				if data, err := emulator.ReceiveAndSend(); err != nil {
					log.Warnf("error %s", err)
				} else {
					log.Infof("data %x", data)
				}
			}
		}
	}

}