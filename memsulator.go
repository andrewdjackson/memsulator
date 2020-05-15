package main

import (
	"os"
	"os/exec"

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
	virtualPortChan chan bool
}

// NewMemsulator creates a new instance
func NewMemsulator() *Memsulator {
	memsulator := &Memsulator{}

	memsulator.homefolder, _ = homedir.Dir()
	memsulator.ecuPort = memsulator.homefolder + "/ttyecu"
	memsulator.fcrPort = memsulator.homefolder + "/ttycodereader"
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

	args := []string{"-d", "-d", "-v", "pty,link=" + memsulator.fcrPort + ",cfmakeraw,ignbrk=1,igncr=1,ignpar=1", "pty,link=" + memsulator.ecuPort + ",cfmakeraw,ignbrk=1,igncr=1,ignpar=1"}
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

	utils.LogI.Println("created virtual serial ports")
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
	// wait for the virtual serial ports to be created
	go memsulator.CreateVirtualPorts()

	// this blocks until the port is ready
	ready := <-memsulator.virtualPortChan

	if ready {
		mems := ecu.NewMemsConnection()
		mems.LoadScenario(memsulator.scenario)
		mems.Open(memsulator.fcrPort)

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

func main() {
	scenario := scenarios.NewScenario()
	scenario.Load("scenarios/coldstart.csv")
	//scenario.ConvertReadMemsLogToMemsFCR("scenarios/response.data")
	//scenario.SaveCSVFile("scenarios/response.csv")

	memsulator := NewMemsulator()
	memsulator.scenario = scenario
	memsulator.startECU()
}
