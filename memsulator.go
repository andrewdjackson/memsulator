package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/andrewdjackson/memsulator/ecu"
	"github.com/andrewdjackson/memsulator/utils"
	"github.com/mitchellh/go-homedir"
	"github.com/tarm/serial"
)

var home, _ = homedir.Dir()
var ecuPort = home + "/ttyecu"
var fcrPort = home + "/ttycodereader"
var virtualPortChan = make(chan bool)

type mems struct {
	s *serial.Port
}

// fileExists reports whether the named file or directory exists.
func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func createVirtualSerialPorts() {
	var cmd *exec.Cmd
	utils.LogI.Printf("creating virtual ports")

	// socat -d -d pty,link=ttycodereader,raw,echo=0 pty,link=ttyecu,raw,echo=0"
	binary, lookErr := exec.LookPath("socat")
	if lookErr != nil {
		panic(lookErr)
	}

	args := []string{"-d", "-d", "pty,link=" + fcrPort + ",raw,echo=0", "pty,link=" + ecuPort + ",raw,echo=0"}
	env := os.Environ()
	cmd = exec.Command(binary)
	cmd.Args = args
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		panic(err)
	}

	fmt.Println("Created virtual serial ports")
}

// CreateVirtualPorts for the FCR to connect to
func CreateVirtualPorts() {
	createVirtualSerialPorts()

	for {
		if fileExists(ecuPort) {
			utils.LogI.Println("virtual serial ports ready")
			virtualPortChan <- true
			break
		}
	}
}

func main() {
	// wait for the virtual serial ports to be created
	go CreateVirtualPorts()

	// this blocks until the port is ready
	ready := <-virtualPortChan

	if ready {
		mems := ecu.NewMemsConnection()
		mems.Open(fcrPort)

		// listen for commands from the FCR
		go mems.ListenToFCRLoop()

		for {
			// wait for the response
			cr := <-mems.ReceivedFromFCR
			// log unexpected responses
			if len(cr.Response) == 0 {
				utils.LogI.Printf("unexepected response for command %x (%x)", cr.Command, cr.Response)
			}
		}
	}
}
