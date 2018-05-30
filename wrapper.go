// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package log4go

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
)

var (
	Global Logger
)

var GlobalLevel Level = DEBUG

func init() {
	Global = NewDefaultLogger(FINE)
}

// Wrapper for (*Logger).LoadConfiguration
func LoadConfiguration(filename string, types ...string) {
	if len(types) > 0 && types[0] == "xml" {
		Global.LoadConfiguration(filename)
	} else {
		Global.LoadJsonConfiguration(filename)
	}
}

// Wrapper for (*Logger).AddFilter
func AddFilter(name string, lvl Level, writer LogWriter) {
	Global.AddFilter(name, lvl, writer)
}

// Wrapper for (*Logger).Close (closes and removes all logwriters)
func Close() {
	Global.Close()
}

func Crash(args ...interface{}) {
	if len(args) > 0 {
		Global.intLogf(CRITICAL, strings.Repeat(" %v", len(args))[1:], args...)
	}
	panic(args)
}

// Logs the given message and crashes the program
func Crashf(format string, args ...interface{}) {
	Global.intLogf(CRITICAL, format, args...)
	Global.Close() // so that hopefully the messages get logged
	panic(fmt.Sprintf(format, args...))
}

// Compatibility with `log`
func Exit(args ...interface{}) {
	if len(args) > 0 {
		Global.intLogf(ERROR, strings.Repeat(" %v", len(args))[1:], args...)
	}
	Global.Close() // so that hopefully the messages get logged
	os.Exit(0)
}

// Compatibility with `log`
func Exitf(format string, args ...interface{}) {
	Global.intLogf(ERROR, format, args...)
	Global.Close() // so that hopefully the messages get logged
	os.Exit(0)
}

// Compatibility with `log`
func Stderr(args ...interface{}) {
	if len(args) > 0 {
		Global.intLogf(ERROR, strings.Repeat(" %v", len(args))[1:], args...)
	}
}

func GetGlobalLevel() Level {
	return GlobalLevel
}

func SetGlobalLevel(level Level) {
	GlobalLevel = level
}

// Compatibility with `log`
func Stderrf(format string, args ...interface{}) {
	Global.intLogf(ERROR, format, args...)
}

// Compatibility with `log`
func Stdout(args ...interface{}) {
	if len(args) > 0 {
		Global.intLogf(INFO, strings.Repeat(" %v", len(args))[1:], args...)
	}
}

// Compatibility with `log`
func Stdoutf(format string, args ...interface{}) {
	Global.intLogf(INFO, format, args...)
}

// Send a log message manually
// Wrapper for (*Logger).Log
func Log(lvl Level, source, message string) {
	Global.Log(lvl, source, message)
}

// Send a formatted log message easily
// Wrapper for (*Logger).Logf
func Logf(lvl Level, format string, args ...interface{}) {
	Global.intLogf(lvl, format, args...)
}

// Send a closure log message
// Wrapper for (*Logger).Logc
func Logc(lvl Level, closure func() string) {
	Global.intLogc(lvl, closure)
}

// Utility for debug log messages
// When given a string as the first argument, this behaves like Logf but with the DEBUG log level (e.g. the first argument is interpreted as a format for the latter arguments)
// When given a closure of type func()string, this logs the string returned by the closure iff it will be logged.  The closure runs at most one time.
// When given anything else, the log message will be each of the arguments formatted with %v and separated by spaces (ala Sprint).
// Wrapper for (*Logger).Debug
func Debugf(arg0 interface{}, args ...interface{}) {
	const (
		lvl = DEBUG
	)
	if lvl < GlobalLevel {
		return
	}
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.intLogc(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

func Debug(args ...interface{}) {
	const (
		lvl = DEBUG
	)
	if lvl < GlobalLevel {
		return
	}
	pc, _, lineno, ok := runtime.Caller(1)
	src := ""
	if ok {
		src = fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), lineno)
	}
	msg := fmt.Sprint(args)
	Global.Log(lvl, src, msg[1:len(msg)-1])
}

// Utility for trace log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Trace
func Tracef(arg0 interface{}, args ...interface{}) {
	const (
		lvl = TRACE
	)
	if lvl < GlobalLevel {
		return
	}
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.intLogc(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

func Trace(args ...interface{}) {
	const (
		lvl = TRACE
	)
	if lvl < GlobalLevel {
		return
	}
	pc, _, lineno, ok := runtime.Caller(1)
	src := ""
	if ok {
		src = fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), lineno)
	}
	msg := fmt.Sprint(args)
	Global.Log(lvl, src, msg[1:len(msg)-1])
}

// Utility for info log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Info
func Infof(arg0 interface{}, args ...interface{}) {
	const (
		lvl = INFO
	)
	if lvl < GlobalLevel {
		return
	}
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.intLogc(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

func Info(args ...interface{}) {
	const (
		lvl = INFO
	)
	if lvl < GlobalLevel {
		return
	}
	pc, _, lineno, ok := runtime.Caller(1)
	src := ""
	if ok {
		src = fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), lineno)
	}
	msg := fmt.Sprint(args)
	Global.Log(lvl, src, msg[1:len(msg)-1])
}

// Utility for warn log messages (returns an error for easy function returns) (see Debug() for parameter explanation)
// These functions will execute a closure exactly once, to build the error message for the return
// Wrapper for (*Logger).Warn
func Warnf(arg0 interface{}, args ...interface{}) error {
	const (
		lvl = WARNING
	)
	if lvl < GlobalLevel {
		return nil
	}
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(lvl, first, args...)
		return errors.New(fmt.Sprintf(first, args...))
	case func() string:
		// Log the closure (no other arguments used)
		str := first()
		Global.intLogf(lvl, "%s", str)
		return errors.New(str)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(lvl, fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
		return errors.New(fmt.Sprint(first) + fmt.Sprintf(strings.Repeat(" %v", len(args)), args...))
	}
	return nil
}

func Warn(args ...interface{}) {
	const (
		lvl = WARNING
	)
	if lvl < GlobalLevel {
		return
	}
	pc, _, lineno, ok := runtime.Caller(1)
	src := ""
	if ok {
		src = fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), lineno)
	}
	msg := fmt.Sprint(args)
	Global.Log(lvl, src, msg[1:len(msg)-1])
}

// Utility for error log messages (returns an error for easy function returns) (see Debug() for parameter explanation)
// These functions will execute a closure exactly once, to build the error message for the return
// Wrapper for (*Logger).Error
func Errorf(arg0 interface{}, args ...interface{}) error {
	const (
		lvl = ERROR
	)
	if lvl < GlobalLevel {
		return nil
	}
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(lvl, first, args...)
		return errors.New(fmt.Sprintf(first, args...))
	case func() string:
		// Log the closure (no other arguments used)
		str := first()
		Global.intLogf(lvl, "%s", str)
		return errors.New(str)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(lvl, fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
		return errors.New(fmt.Sprint(first) + fmt.Sprintf(strings.Repeat(" %v", len(args)), args...))
	}
	return nil
}

func Error(args ...interface{}) {
	const (
		lvl = ERROR
	)
	if lvl < GlobalLevel {
		return
	}
	pc, _, lineno, ok := runtime.Caller(1)
	src := ""
	if ok {
		src = fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), lineno)
	}
	msg := fmt.Sprint(args)
	Global.Log(lvl, src, msg[1:len(msg)-1])
}

// Utility for critical log messages (returns an error for easy function returns) (see Debug() for parameter explanation)
// These functions will execute a closure exactly once, to build the error message for the return
// Wrapper for (*Logger).Critical
func Criticalf(arg0 interface{}, args ...interface{}) error {
	const (
		lvl = CRITICAL
	)
	if lvl < GlobalLevel {
		return nil
	}
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(lvl, first, args...)
		return errors.New(fmt.Sprintf(first, args...))
	case func() string:
		// Log the closure (no other arguments used)
		str := first()
		Global.intLogf(lvl, "%s", str)
		return errors.New(str)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(lvl, fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
		return errors.New(fmt.Sprint(first) + fmt.Sprintf(strings.Repeat(" %v", len(args)), args...))
	}
	return nil
}

func Critical(args ...interface{}) {
	const (
		lvl = CRITICAL
	)
	if lvl < GlobalLevel {
		return
	}
	pc, _, lineno, ok := runtime.Caller(1)
	src := ""
	if ok {
		src = fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), lineno)
	}
	msg := fmt.Sprint(args)
	Global.Log(lvl, src, msg[1:len(msg)-1])
}
