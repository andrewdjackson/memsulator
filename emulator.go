package main

import (
	"flag"
	"github.com/andrewdjackson/memsulator/loader"
	"github.com/andrewdjackson/memsulator/responder"
	log "github.com/sirupsen/logrus"
	"strings"
)

func main() {
	file := flag.String("scenario", "scenarios/default.csv", "scenario file to run")
	port := flag.String("port", "/dev/serial0", "serial communication port")
	flag.Parse()

	scenario := loader.NewScenario()
	if err := scenario.Load(*file); err == nil {

		if isVirtualSerialPort(*port) {
			if err = createVirtualSerialPort(); err != nil {
				log.Fatalf("unable to connect (%s)", err)
			}
		}

		playback := loader.NewPlayback(scenario)
		playback.Start()

		emulator := responder.NewEmulator(playback)
		if connected, err := emulator.Connect(*port); err == nil {
			if connected {
				go emulator.Listen()

				select {}
			} else {
				log.Fatalf("unable to connect (%s)", err)
			}
		}
	}
}

func isVirtualSerialPort(port string) bool {
	return strings.HasSuffix(port, "ttycodereader")
}

func createVirtualSerialPort() error {
	virtualPort := responder.NewVirtualSerialPort()
	return virtualPort.CreateVirtualPorts()
}
