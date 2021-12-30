package fakedns

import (
	"log"
)

type Level int

const (
	ERROR Level = iota
	INFO
)

var levelNames = []string{
	"ERROR",
	"INFO",
}

func (p Level) String() string {
	return levelNames[p]
}

type Logger interface {
	Printf(lvl Level, format string, v ...interface{})
}

type defaultLogger struct {
	lvl Level
	log *log.Logger
}

func NewDefaultLogger(lvl Level, log *log.Logger) Logger {
	return &defaultLogger{
		lvl: lvl,
		log: log,
	}
}

func (dl *defaultLogger) Printf(lvl Level, format string, v ...interface{}) {
	if dl.lvl >= lvl {
		dl.log.Printf(format, v...)
	}
}
