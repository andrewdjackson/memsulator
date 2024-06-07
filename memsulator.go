package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/andrewdjackson/memsulator/ecu"
	"github.com/andrewdjackson/memsulator/scenarios"
	"github.com/andrewdjackson/memsulator/utils"
	"github.com/mitchellh/go-homedir"
)

// Memsulator instance struct
type Memsulator struct {
	scenario        *scenarios.Scenario
	homefolder      string
	ecuPort         string
	fcrPort         string
	memsVersion     string
	virtualPortChan chan bool
}

// NewMemsulator creates a new instance
func NewMemsulator() *Memsulator {
	memsulator := &Memsulator{}

	memsulator.homefolder, _ = homedir.Dir()
	memsulator.ecuPort = filepath.ToSlash(memsulator.homefolder + "/ttyecu")
	memsulator.fcrPort = filepath.ToSlash(memsulator.homefolder + "/ttycodereader")
	memsulator.virtualPortChan = make(chan bool)

	return memsulator
}

// fileExists reports whether the named file or directory exists.
func (memsulator *Memsulator) fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func (memsulator *Memsulator) createVirtualSerialPorts() {
	var cmd *exec.Cmd
	utils.LogI.Printf("creating virtual ports")

	// socat -d -d pty,link=ttycodereader,raw,echo=0 pty,link=ttyecu,raw,echo=0"
	binary, lookErr := exec.LookPath("socat")
	if lookErr != nil {
		utils.LogE.Fatalf("unable to find socat command, brew install socat? (%s)", lookErr)
	}

	args := []string{"-d", "-d", "pty,link='" + memsulator.fcrPort + "',cfmakeraw,ignbrk=1,igncr=1,ignpar=1", "pty,link='" + memsulator.ecuPort + "',cfmakeraw,ignbrk=1,igncr=1,ignpar=1"}
	env := os.Environ()
	cmd = exec.Command(binary)
	cmd.Args = args
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()

	if err != nil {
		utils.LogE.Fatalf("cmd.Run() failed with %s", err)
	}

	utils.LogI.Printf("created virtual serial ports (%s)", cmd.String())
}

// CreateVirtualPorts for the FCR to connect to
func (memsulator *Memsulator) CreateVirtualPorts() {
	memsulator.createVirtualSerialPorts()

	for {
		if memsulator.fileExists(memsulator.ecuPort) {
			utils.LogI.Println("virtual serial ports ready")
			memsulator.virtualPortChan <- true
			break
		}
	}
}

func (memsulator *Memsulator) startECU() {
	var ready bool

	if strings.HasSuffix(memsulator.fcrPort, "ttycodereader") {
		// wait for the virtual serial ports to be created
		go memsulator.CreateVirtualPorts()

		// this blocks until the port is ready
		ready = <-memsulator.virtualPortChan
	} else {
		utils.LogI.Printf("opening port %s", memsulator.fcrPort)
		ready = true
	}

	if ready {
		mems := ecu.NewMemsConnection()
		mems.MemsVersion = memsulator.memsVersion
		mems.LoadScenario(memsulator.scenario)
		mems.Open(memsulator.fcrPort)

		if mems.Connected {
			// listen for commands from the FCR
			go mems.Listen()

			for {
				// wait for the response
				cr := <-mems.ReceivedFromFCR
				// log unexpected responses
				if len(cr.Response) == 0 {
					utils.LogI.Printf("unexepected response for command %x (%x)", cr.Command, cr.Response)
				}
			}
		} else {
			utils.LogE.Fatal("unable to connect, terminating")
		}
	}
}

/*
func main() {
	scenefile := flag.String("scenario", "scenarios/default.csv", "scenario file to run")
	port := flag.String("port", "", "serial communication port")
	memsVersion := flag.String("mems", "1.6", "mems version (1.3, 1.6 or 1.9)")
	flag.Parse()

	scenario := scenarios.NewScenario()

	// run the specified scenario
	utils.LogI.Printf("running scenario..")

	scenario.Load(*scenefile)
	memsulator := NewMemsulator()
	if *port != "" {
		memsulator.fcrPort = *port
	}
	if *memsVersion != "" {
		memsulator.memsVersion = *memsVersion

	}
	memsulator.scenario = scenario
	memsulator.startECU()
}
*/
