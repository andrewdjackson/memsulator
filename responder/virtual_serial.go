package responder

import (
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
)

type VirtualSerialPort struct {
	homefolder      string
	ECUPort         string
	FCRPort         string
	virtualPortChan chan bool
	Connected       bool
}

func NewVirtualSerialPort() *VirtualSerialPort {
	vserial := &VirtualSerialPort{}
	vserial.homefolder, _ = homedir.Dir()
	vserial.ECUPort = filepath.ToSlash(vserial.homefolder + "/ttyecu")
	vserial.FCRPort = filepath.ToSlash(vserial.homefolder + "/ttycodereader")
	vserial.virtualPortChan = make(chan bool)
	vserial.Connected = false

	// create the virtual serial ports
	// this runs as a go routine as it waits for the file to be
	// created and send a response over the channel to signal
	// creation complete
	go vserial.CreateVirtualPorts()

	// this blocks until the port is ready
	var ready bool
	ready = <-vserial.virtualPortChan

	if ready {
		vserial.Connected = true
		log.Infof("virtual serial port created")
	}

	return vserial
}

func (vserial *VirtualSerialPort) CreateVirtualPorts() {
	if err := vserial.createVirtualSerialPorts(); err == nil {
		vserial.waitForPort()
	}
}

func (vserial *VirtualSerialPort) createVirtualSerialPorts() error {
	var cmd *exec.Cmd
	log.Infof("creating virtual ports")

	// socat -d -d pty,link=ttycodereader,raw,echo=0 pty,link=ttyecu,raw,echo=0"
	binary, lookErr := exec.LookPath("socat")
	if lookErr != nil {
		log.Warnf("unable to find socat command, brew install socat? (%s)", lookErr)
		return lookErr
	}

	args := []string{"-d", "-d", "pty,link='" + vserial.FCRPort + "',cfmakeraw,ignbrk=1,igncr=1,ignpar=1", "pty,link='" + vserial.ECUPort + "',cfmakeraw,ignbrk=1,igncr=1,ignpar=1"}
	env := os.Environ()
	cmd = exec.Command(binary)
	cmd.Args = args
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()

	if err != nil {
		log.Errorf("cmd.Run() failed with %s", err)
	}

	log.Infof("created virtual serial ports (%s)", cmd.String())

	return err
}

func (vserial *VirtualSerialPort) waitForPort() {
	for {
		if vserial.fileExists(vserial.ECUPort) {
			log.Infof("virtual serial ports ready")
			vserial.virtualPortChan <- true
			break
		}
	}
}

func (vserial *VirtualSerialPort) fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
