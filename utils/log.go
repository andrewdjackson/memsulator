package utils

import (
	"log"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	// LogE logs as an error
	LogE = log.New(LogWriter{}, "ERROR: ", 0)
	// LogW logs as a warning
	LogW = log.New(LogWriter{}, "WARN: ", 0)
	// LogI logs as an info
	LogI = log.New(LogWriter{}, "INFO: ", 0)
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

	log.Printf("%s:%d %s: %s", filepath.Base(file), line, fnName, p)
	return len(p), nil
}
