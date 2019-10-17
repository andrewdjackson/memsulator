package main

import (
	"andrewj.com/memsulator/readmems"
	"andrewj.com/memsulator/rosco"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/tarm/serial"
	"log"
	"os"
	"os/exec"
	"time"
)

var home, _ = homedir.Dir()
var ecu_port = home + "/ttyecu"
var pcr_port = home + "/ttycodereader"

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

func createVirtualSerialPorts(cmd *exec.Cmd) {
	fmt.Println("Creating virtual ports")

	// socat -d -d pty,link=ttycodereader,raw,echo=0 pty,link=ttyecu,raw,echo=0"
	binary, lookErr := exec.LookPath("socat")
	if lookErr != nil {
		panic(lookErr)
	}

	args := []string{"-d", "-d", "pty,link=" + pcr_port + ",raw,echo=0", "pty,link=" + ecu_port + ",raw,echo=0"}
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

func waitForVirtualPorts() bool {
	max := 2000

	// wait for virtual serial port alias is created
	for count := 0; count <= max; count++ {
		if fileExists(ecu_port) {
			break
		}

		time.Sleep(100)
	}

	fmt.Println("Virtual serial ports ready")

	return true
}

func (mems *mems) serialLoop() {
	rosco.UseResponseFile = true

	// commands are 1 byte (sometime 2 bytes but let's not worry about that at the moment)
	cmd := make([]byte, 1)

	// wait for a command from the code reader
	n, _ := mems.s.Read(cmd)

	if n > 0 {
		log.Printf("read (%d): %x", n, cmd[:n])

		// find the command response
		response := rosco.Response(cmd)

		// write the response to the code reader
		mems.s.Write(response)
		log.Printf("write: %x", response)
	}

	// sleep for a few milliseconds to regulate the responses
	time.Sleep(200)
}

func main() {
	var readmemsConfig readmems.Config
	var cmd *exec.Cmd

	defer cmd.Wait()
	go createVirtualSerialPorts(cmd)
	waitForVirtualPorts()

	// use if the readmems config is supplied
	readmemsConfig = readmems.ReadConfig()

	// if argument is supplied then use that as the port id
	if len(os.Args) > 1 {
		readmemsConfig.Port = os.Args[1]
	}

	mems := mems{}

	// connect to the code reader
	config := &serial.Config{Name: readmemsConfig.Port, Baud: 9600}

	s, err := serial.OpenPort(config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Emulating ECU on ", ecu_port)
	fmt.Println("Listening for Code Reader on ", pcr_port)

	mems.s = s
	mems.s.Flush()

	// listen for commands and respond
	for {
		mems.serialLoop()
	}
}
