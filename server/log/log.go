// The log package provides a light weight wrapper around go's standard log
// package that provides logging at different verbosity levels.
//
// Rules of thumb for log levels are as follows:
//
//     FATAL - Should only be used in a totally unrecoverable scenario. Logging
//             at the fatal level will cause the app to panic.
//
//     ERROR - The error level should only be used to log errors where the
//             request cannot be completed, typically because of some erroneous
//             server side state.
//
//     WARNING - The warning level should be used to log messages that should
//               not happen often, but are not fatal to the request. Typically
//               these will be client side errors, for example improperly
//               encoded request bodies.
//
//     INFO - The info level should be used to provide insight in the typical
//            activity of the app. There should be approximately one info log
//            per request.
//
//     DEBUG - The debug level should be used for messages that provide insight
//            in to the line by line path through the code. There can be many
//            debug statements per request. The debug log level should not be
//            used to log large go data structures.
//
//     VERBOSE - The verbose level can be used to log complete or partial go
//               data structures, or any other information that would otherwise
//               clutter the log output.
package log

import (
	"log"
	"runtime"
	"strconv"
	"strings"
)

const (
	Fatal   = 0
	Error   = iota
	Warn    = iota
	Info    = iota
	Debug   = iota
	Verbose = iota
)

var logLevel = Info

var levelToSymbol = map[int]string{
	Fatal:   "FATAL  ",
	Error:   "ERROR  ",
	Warn:    "WARNING",
	Info:    "INFO   ",
	Debug:   "DEBUG  ",
	Verbose: "VERBOSE",
}

// Init sets up the initial internal state of the logging package. To ensure
// consistent log statements, this method should be called prior to logging any
// thing.
func Init() {
	log.SetFlags(log.Ldate | log.Ltime)
}

// SetLogLevel sets the verbosity of logging. The default is Info.
func SetLogLevel(level int) {
	logLevel = level
}

// chop removes the prefix of haystack up to and including the needle string. If
// the needle is not found, haystack is returned as is.
func chop(haystack, needle string) string {
	i := strings.LastIndex(haystack, needle)
	if i < 0 {
		return haystack
	}
	return haystack[i+len(needle):]
}

func compress(path string) string {
	parts := strings.Split(path, "/")
	for i := 0; i < len(parts); i++ {
		if i != len(parts)-1 {
			parts[i] = parts[i][0:1]
		}
	}
	return strings.Join(parts, "/")
}

func logf(level int, format string, args ...interface{}) {
	if level > logLevel {
		return
	}

	// Set context to be the relative path of the source file which all folders
	// truncated to their first letter, e.g. this file would be s/l/log.go.
	_, file, line, ok := runtime.Caller(2)
	var context string
	if !ok {
		context = strings.Repeat("?", 10)
	} else {
		path := compress(chop(file, "/pifuxelck-server-go/"))
		context = path + ":" + strconv.Itoa(line)
	}

	symbol := levelToSymbol[level]
	if symbol == "" {
		symbol = levelToSymbol[Verbose]
	}

	format = symbol + " [" + context + "] " + format
	if level > 0 {
		log.Printf(format, args...)
	} else {
		log.Panicf(format, args...)
	}
}

// Fatalf logs a format string at the Fatal log level. Logging at this level
// will cause the server to panic.
func Fatalf(format string, args ...interface{}) {
	logf(Fatal, format, args...)
}

// Errorf logs a format string at the Error log level.
func Errorf(format string, args ...interface{}) {
	logf(Error, format, args...)
}

// Warnf logs a format string at the Warn log level.
func Warnf(format string, args ...interface{}) {
	logf(Warn, format, args...)
}

// Infof logs a format string at the Info log level.
func Infof(format string, args ...interface{}) {
	logf(Info, format, args...)
}

// Debugf logs a format string at the Info log level.
func Debugf(format string, args ...interface{}) {
	logf(Debug, format, args...)
}

// Verbosef logs a format string at the Info log level.
func Verbosef(format string, args ...interface{}) {
	logf(Verbose, format, args...)
}
