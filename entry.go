package logrjack

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/Sirupsen/logrus"
)

type entry struct {
	entry *logrus.Entry
}

// Entry is the interface returned by NewEntry.
//
// Use AddField or AddFields (opposed to Logrus' WithField)
// to add new fields to an entry.
//
// To add source file and line information call AddCallstack.
//
// Call the Info, Warn, Error, or Fatal methods to write the log entry.
type Entry interface {
	// AddField adds a single field to the Entry.
	AddField(key string, value interface{})

	// AddFields adds a map of fields to the Entry.
	AddFields(fields map[string]interface{})

	// AddCallstack adds the current callstack to the Entry using the key "callstack".
	//
	// Excludes runtime/proc.go, http/server.go, and files ending in .s from the callstack.
	AddCallstack()

	// Info logs the Entry with a status of "info".
	Info(args ...interface{})
	// Infof logs the Entry with a status of "info".
	Infof(format string, args ...interface{})
	// Warn logs the Entry with a status of "warn".
	Warn(args ...interface{})
	// Warnf logs the Entry with a status of "warn".
	Warnf(format string, args ...interface{})
	// Error logs the Entry with a status of "error".
	Error(args ...interface{})
	// Errorf logs the Entry with a status of "error".
	Errorf(format string, args ...interface{})
	// Fatal logs the Entry with a status of "fatal" and calls os.Exit(1).
	Fatal(args ...interface{})
	// Fatalf logs the Entry with a status of "fatal" and calls os.Exit(1).
	Fatalf(format string, args ...interface{})

	// AddError adds a field "err" with the specified error.
	AddError(err error)

	String() string
}

var std *logrus.Logger

func init() {
	std = logrus.StandardLogger()
}

// NewEntry creates a new log Entry.
func NewEntry() Entry {
	e := &entry{}
	e.entry = logrus.NewEntry(std)
	return e
}

// Callstack creates a new log Entry with the callstack. If 'err' is not nil,
// an entry "err" is added with the specified error.
func Callstack(err error) Entry {
	e := &entry{}
	e.entry = logrus.NewEntry(std)
	e.addCallstack(3)
	if err != nil {
		e.AddError(err)
	}
	return e
}

// String returns the string representation from the reader and
// ultimately the formatter.
func (e entry) String() string {
	s, err := e.entry.String()
	if err != nil {
		return fmt.Sprintf("%s - <%s>", s, err)
	}
	return s
}

// AddError adds a field "err" with the specified error.
func (e *entry) AddError(err error) {
	e.AddField("err", err)
}

// AddField adds a single field to the Entry.
func (e *entry) AddField(key string, value interface{}) {
	e.entry = e.entry.WithField(key, value)
}

// AddFields adds a map of fields to the Entry.
func (e *entry) AddFields(fields map[string]interface{}) {
	logrusFields := logrus.Fields{}
	for k, v := range fields {
		logrusFields[k] = v
	}
	e.entry = e.entry.WithFields(logrusFields)
}

// AddCallstack adds the current callstack to the Entry using the key "callstack".
//
// Excludes runtime/proc.go, http/server.go, and files ending in .s from the callstack.
func (e *entry) AddCallstack() {
	e.addCallstack(3)
}

func (e *entry) addCallstack(skipFrames int) {
	var cs []string
	for i := 0; ; i++ {
		file, line := getCaller(skipFrames + i)
		if file == "???" {
			break
		}
		if strings.HasSuffix(file, ".s") ||
			strings.HasPrefix(file, "http/server.go") ||
			strings.HasPrefix(file, "runtime/proc.go") {
			continue
		}
		cs = append(cs, fmt.Sprintf("%s:%d", file, line))
	}
	e.AddField("callstack", strings.Join(cs, ", "))
}

// Info logs the Entry with a status of "info".
func (e *entry) Info(args ...interface{}) {
	e.entry.Info(args...)
}

// Infof logs the Entry with a status of "info".
func (e *entry) Infof(format string, args ...interface{}) {
	e.entry.Infof(format, args...)
}

// Warn logs the Entry with a status of "warn".
func (e *entry) Warn(args ...interface{}) {
	e.entry.Warn(args...)
}

// Warnf logs the Entry with a status of "warn".
func (e *entry) Warnf(format string, args ...interface{}) {
	e.entry.Warnf(format, args...)
}

// Error logs the Entry with a status of "error".
func (e *entry) Error(args ...interface{}) {
	e.entry.Error(args...)
}

// Errorf logs the Entry with a status of "error".
func (e *entry) Errorf(format string, args ...interface{}) {
	e.entry.Errorf(format, args...)
}

// Fatal logs the Entry with a status of "fatal" and calls os.Exit(1).
func (e *entry) Fatal(args ...interface{}) {
	e.entry.Fatal(args...)
}

// Fatalf logs the Entry with a status of "fatal" and calls os.Exit(1).
func (e *entry) Fatalf(format string, args ...interface{}) {
	e.entry.Fatalf(format, args...)
}

func getCaller(skip int) (file string, line int) {
	var ok bool
	_, file, line, ok = runtime.Caller(skip)
	if !ok {
		file = "???"
		line = 0
	} else {
		short := file
		count := 0
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				count++
				if count == 2 {
					short = file[i+1:]
					break
				}
			}
		}
		file = short
	}
	return
}
