package utils

import (
	"log"
	"path/filepath"
	"runtime"
	"strings"
)

// ECUCommandTrace color coding for Commands sent to the ECU
const ECUCommandTrace = "\u001b[38;5;200mECU_TX>\u001b[0m"

// ECUResponseTrace color coding for Reponses returned from the ECU
const ECUResponseTrace = "\u001b[38;5;200mECU_RX<\u001b[0m"

// ReceiveFromWebTrace color coding for Messages sent to the Web
const ReceiveFromWebTrace = "\u001b[38;5;21mWEB_RX<\u001b[0m"

// SendToWebTrace color coding for Messages received from the Web
const SendToWebTrace = "\u001b[38;5;21mWEB_TX>\u001b[0m"

var (
	// LogE logs as an error
	LogE = log.New(LogWriter{}, "\u001b[38;5;160mERROR: ", 0)
	// LogW logs as a warning
	LogW = log.New(LogWriter{}, "\u001b[38;5;214mWARNING: ", 0)
	// LogI logs as an info, no prefix
	LogI = log.New(LogWriter{}, "\u001b[38;5;70m", 0)
)

// LogWriter is used to format the log message
type LogWriter struct{}

// Write the log entry
func (f LogWriter) Write(p []byte) (n int, err error) {
	pc, file, line, ok := runtime.Caller(4)
	if !ok {
		file = "?"
		line = 0
	}

	fn := runtime.FuncForPC(pc)
	var fnName string
	if fn == nil {
		fnName = "?()"
	} else {
		dotName := filepath.Ext(fn.Name())
		fnName = strings.TrimLeft(dotName, ".") + "()"
	}

	log.Printf("%s\r \u001b[38;5;38m↵ %s: %d %s\u001b[0m", p, filepath.Base(file), line, fnName)
	return len(p), nil
}
