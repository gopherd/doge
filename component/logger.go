package component

import (
	"log"
)

// Logger represents log interface of component
type Logger interface {
	Debug(format string, args ...interface{})
	Trace(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
}

// logger implements Logger interface
type logger string

func (l logger) log(level string, format string, args ...interface{}) {
	log.Printf("["+level+"] ["+string(l)+"] "+format, args...)
}

func (l logger) Debug(format string, args ...interface{}) { l.log("D", format, args...) }
func (l logger) Trace(format string, args ...interface{}) { l.log("T", format, args...) }
func (l logger) Info(format string, args ...interface{})  { l.log("I", format, args...) }
func (l logger) Warn(format string, args ...interface{})  { l.log("W", format, args...) }
func (l logger) Error(format string, args ...interface{}) { l.log("E", format, args...) }
func (l logger) Fatal(format string, args ...interface{}) { l.log("C", format, args...) }
