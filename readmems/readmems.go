package readmems

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Config readmems configuration
type Config struct {
	Port       string
	Command    string
	Output     string
	Loop       string
	Connection string
}

func init() {

}

// reads a whole file into memory and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

// ReadConfig reads readmems.cfg file
func ReadConfig() Config {
	var c Config

	// set up defaults
	c.Port = "ttycodereader"
	c.Command = "read"
	c.Loop = "inf"
	c.Output = "stdout"
	c.Connection = "wait"

	lines, err := readLines("readmems.cfg")

	if err == nil {
		for i := range lines {
			// ignore comment lines or lines that are not value pairs
			if !strings.HasPrefix(lines[i], "#") {
				if strings.Contains(lines[i], "=") {
					data := strings.Split(lines[i], "=")
					switch data[0] {
					case "port":
						c.Port = data[1]
					case "command":
						c.Command = data[1]
					case "loop":
						c.Loop = data[1]
					case "output":
						c.Output = data[1]
					case "connection":
						c.Connection = data[1]
					}
				}
			}
		}
	}

	fmt.Println("ReadMems Config", c)

	return c
}
