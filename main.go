package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"golang.org/x/sys/windows/svc/eventlog"
	"log"
	"time"
)

type WindowsLogger struct {
	ev   *eventlog.Log
	errs chan<- error
}

type Config struct {
	Name     string
	Keyword  string
	Hostname string
	Type     string
	Id       string
	Msg      string
}

var (
	name     = flag.String("name", "", "event name")
	keyword  = flag.String("keyword", "", "keyword")
	hostname = flag.String("host", "", "host name")
	types    = flag.String("type", "", "CRITICAL/HARMLESS")
	id       = flag.String("id", "", "msg id")
	msg      = flag.String("msg", "", "message")
)

var c Config

func main() {
	initConfig()

	log.Println("start")

	logger := New(c.Name)
	defer Close(c.Name)

	logger.Error(NewMsg())

	log.Println("finished")
}

func initConfig() {
	_, err := toml.DecodeFile("config.toml", &c)
	if err != nil {
		log.Fatal(err)
	}

	flag.Parse()

	if *keyword != "" {
		c.Keyword = *keyword
	}
	if *name != "" {
		c.Name = *name
	}
	if *hostname != "" {
		c.Hostname = *hostname
	}
	if *types != "" {
		c.Type = *types
	}
	if *id != "" {
		c.Id = *id
	}
	if *msg != "" {
		c.Msg = *msg
	}
}

func NewMsg() string {
	format := time.Now().Format("%s 2006 Jan 02 15:04:05 %s %s %s %s")
	return fmt.Sprintf(format, c.keyword, c.Hostname, c.Type, c.Id, c.Msg)
}

func New(name string) *WindowsLogger {
	err := eventlog.InstallAsEventCreate(name, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		log.Fatalf("InstallAsEventCreate() failed: name = %s %s", name, err)
	}

	el, err := eventlog.Open(name)
	if err != nil {
		log.Fatalf("Cannot open eventlog: %s", err)
	}

	errs := make(chan<- error)
	return &WindowsLogger{el, errs}
}

func Close(name string) {
	err := eventlog.Remove(name)
	if err != nil {
		log.Fatalf("RemoveEventLogSource() failed: %s", err)
	}
}

func (l WindowsLogger) send(err error) error {
	if err == nil {
		return nil
	}
	if l.errs != nil {
		l.errs <- err
	}
	return err
}

// Error logs an error message.
func (l WindowsLogger) Error(v ...interface{}) error {
	return l.send(l.ev.Error(3, fmt.Sprint(v...)))
}

// Warning logs an warning message.
func (l WindowsLogger) Warning(v ...interface{}) error {
	return l.send(l.ev.Warning(2, fmt.Sprint(v...)))
}

// Info logs an info message.
func (l WindowsLogger) Info(v ...interface{}) error {
	return l.send(l.ev.Info(1, fmt.Sprint(v...)))
}

// Errorf logs an error message.
func (l WindowsLogger) Errorf(format string, a ...interface{}) error {
	return l.send(l.ev.Error(3, fmt.Sprintf(format, a...)))
}

// Warningf logs an warning message.
func (l WindowsLogger) Warningf(format string, a ...interface{}) error {
	return l.send(l.ev.Warning(2, fmt.Sprintf(format, a...)))
}

// Infof logs an info message.
func (l WindowsLogger) Infof(format string, a ...interface{}) error {
	return l.send(l.ev.Info(1, fmt.Sprintf(format, a...)))
}

// NError logs an error message and an event ID.
func (l WindowsLogger) NError(eventID uint32, v ...interface{}) error {
	return l.send(l.ev.Error(eventID, fmt.Sprint(v...)))
}

// NWarning logs an warning message and an event ID.
func (l WindowsLogger) NWarning(eventID uint32, v ...interface{}) error {
	return l.send(l.ev.Warning(eventID, fmt.Sprint(v...)))
}

// NInfo logs an info message and an event ID.
func (l WindowsLogger) NInfo(eventID uint32, v ...interface{}) error {
	return l.send(l.ev.Info(eventID, fmt.Sprint(v...)))
}

// NErrorf logs an error message and an event ID.
func (l WindowsLogger) NErrorf(eventID uint32, format string, a ...interface{}) error {
	return l.send(l.ev.Error(eventID, fmt.Sprintf(format, a...)))
}

// NWarningf logs an warning message and an event ID.
func (l WindowsLogger) NWarningf(eventID uint32, format string, a ...interface{}) error {
	return l.send(l.ev.Warning(eventID, fmt.Sprintf(format, a...)))
}

// NInfof logs an info message and an event ID.
func (l WindowsLogger) NInfof(eventID uint32, format string, a ...interface{}) error {
	return l.send(l.ev.Info(eventID, fmt.Sprintf(format, a...)))
}
