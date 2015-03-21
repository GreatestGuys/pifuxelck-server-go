package log

import (
	"log"
)

const (
	Fatal   = 0
	Warn    = iota
	Info    = iota
	Debug   = iota
	Verbose = iota
)

var logLevel = Info

var levelToSymbol = map[int]string{
	Fatal:   "FATAL  ",
	Warn:    "WARNING",
	Info:    "INFO   ",
	Debug:   "DEBUG  ",
	Verbose: "VERBOSE",
}

// SetLogLevel sets the verbosity of logging. The default is Info.
func SetLogLevel(level int) {
	logLevel = level
}

// Logf logs a format string method at the given log level.
func Logf(level int, format string, args ...interface{}) {
	if level <= logLevel {
		symbol := levelToSymbol[level]
		if symbol == "" {
			symbol = "?"
		}
		format = symbol + " " + format
		if level > 0 {
			log.Printf(format, args...)
		} else {
			log.Panicf(format, args...)
		}
	}
}

// Fatalf logs a format string at the Fatal log level. Logging at this level
// will cause the server to panic.
func Fatalf(format string, args ...interface{}) {
	Logf(Fatal, format, args...)
}

// Warnf logs a format string at the Warn log level.
func Warnf(format string, args ...interface{}) {
	Logf(Warn, format, args...)
}

// Infof logs a format string at the Info log level.
func Infof(format string, args ...interface{}) {
	Logf(Info, format, args...)
}

// Debugf logs a format string at the Info log level.
func Debugf(format string, args ...interface{}) {
	Logf(Debug, format, args...)
}

// Verbosef logs a format string at the Info log level.
func Verbosef(format string, args ...interface{}) {
	Logf(Verbose, format, args...)
}
