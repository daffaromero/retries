package logger

import (
	"log"
	"os"
)

type Log struct {
	out *log.Logger
	err *log.Logger
}

type Options struct {
	IsPrintStack bool
	IsExit       bool
	ExitCode     int
}

func NewLog(prefix string) *Log {
	return &Log{
		out: log.New(os.Stdout, "[LOG]["+prefix+"]", log.Default().Flags()),
		err: log.New(os.Stderr, "[ERROR]["+prefix+"]", log.Default().Flags()),
	}
}
