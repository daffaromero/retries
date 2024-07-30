package logger

import (
	"fmt"
	"os"
	"runtime/debug"
)

func (l *Log) Error(message interface{}, options ...*Options) {
	msg := []interface{}{message}

	if len(options) > 0 && options[0].IsPrintStack {
		msg = append(msg, fmt.Sprintf("\n%s", debug.Stack()))
	}

	l.err.Println(msg...)

	if len(options) > 0 && options[0].IsExit {
		exitCode := 1
		if options[0].ExitCode > 1 {
			exitCode = options[0].ExitCode
		}

		os.Exit(exitCode)
	}
}

func (l *Log) CustomError(title string, message interface{}, options ...*Options) {
	msg := []interface{}{}
	msg = append(msg, title)
	msg = append(msg, message)

	if len(options) > 0 && options[0].IsPrintStack {
		msg = append(msg, fmt.Sprintf("\n%s", debug.Stack()))
	}

	l.err.Println(msg...)

	if len(options) > 0 && options[0].IsExit {
		exitCode := 1
		if options[0].ExitCode > 1 {
			exitCode = options[0].ExitCode
		}

		os.Exit(exitCode)
	}
}
