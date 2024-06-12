package responder

import (
	"context"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const ECU_PORT = "/ttyecu"
const FCR_PORT = "/ttycodereader"
const SOCAT_EXE = "socat"

type VirtualSerialPort struct {
	homeFolder      string
	ecuPort         string
	fcrPort         string
	virtualPortChan chan bool
}

func NewVirtualSerialPort() *VirtualSerialPort {
	vs := &VirtualSerialPort{}
	vs.homeFolder, _ = homedir.Dir()
	vs.ecuPort = filepath.ToSlash(vs.homeFolder + ECU_PORT)
	vs.fcrPort = filepath.ToSlash(vs.homeFolder + FCR_PORT)
	vs.virtualPortChan = make(chan bool)

	return vs
}

func (vs *VirtualSerialPort) CreateVirtualPorts() error {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*80))
	defer cancel()

	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return vs.createVirtualSerialPorts()
	})

	group.Go(func() error {
		vs.waitForPort()
		return nil
	})

	if err = group.Wait(); err != nil {
		log.Errorf("error creating virtual serial ports (%s)", err)
	} else {
		log.Infof("virtual serial port created")
	}

	return err
}

func (vs *VirtualSerialPort) createVirtualSerialPorts() error {
	var err error
	var path string

	log.Infof("creating virtual ports")

	if path, err = findSocat(); err == nil {
		cmdline := vs.buildCommandLine(path)
		if err = cmdline.Start(); err != nil {
			log.Errorf("cmd.Run() failed with %s", err)
		}

		log.Infof("created virtual serial ports (%s)", cmdline.String())
	}

	return err
}

func findSocat() (string, error) {
	var err error
	var path string

	if path, err = exec.LookPath(SOCAT_EXE); err != nil {
		log.Errorf("unable to find socat command, brew install socat? (%s)", err)
	}

	return path, err
}

// create the socat command line as follows
// socat -d -d pty, link=ttycodereader, cfmakeraw, ignbrk=1, igncr=1, ignpar=1, pty, link=ttyecu, cfmakeraw, ignbrk=1, igncr=1, ignpar=1
func (vs *VirtualSerialPort) buildCommandLine(socatPath string) *exec.Cmd {
	args := []string{"-d", "-d", "pty,link='" + vs.fcrPort + "',cfmakeraw,ignbrk=1,igncr=1,ignpar=1", "pty,link='" + vs.ecuPort + "',cfmakeraw,ignbrk=1,igncr=1,ignpar=1"}
	env := os.Environ()
	cmdline := exec.Command(socatPath)

	cmdline.Args = args
	cmdline.Env = env
	cmdline.Stdout = os.Stdout
	cmdline.Stderr = os.Stderr

	return cmdline
}

func (vs *VirtualSerialPort) waitForPort() {
	for {
		if vs.fileExists(vs.ecuPort) {
			break
		}
	}

	log.Infof("virtual serial port verified")
}

func (vs *VirtualSerialPort) fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
