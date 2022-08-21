package httplogger

import (
	"fmt"
	"io"
	"os"
)

type curState byte

const (
	start         curState = iota
	foundCr                // found "\r"
	foundCrLf              // found "\r\n"
	foundCrLfCr            // found "\r\n\r"
	foundCrLfCrLf          // found "\r\n\r\n"
)

// HTTPLogger is to parse incoming byte stream and print http request and response header
type HTTPLogger struct {
	logger io.Writer
	state  curState
}

// New creates one logger instance
func New(l io.Writer) *HTTPLogger {
	if l == nil {
		l = os.Stdout
	}
	return &HTTPLogger{logger: l}
}

// Reset resets the internal state and start new parse
func (l *HTTPLogger) Reset() {
	l.state = start
}

// Print parse incoming bytes and print
func (l *HTTPLogger) Print(buf []byte) {
	for _, b := range buf {
		if l.state == foundCrLfCrLf {
			return
		}
		if b != '\r' && b != '\n' {
			fmt.Fprint(l.logger, string(b))
			l.Reset()
			continue
		}
		switch b {
		case '\r':
			switch l.state {
			case start:
				l.state = foundCr
			case foundCrLf:
				l.state = foundCrLfCr
			default:
				l.Reset()
			}
		case '\n':
			switch l.state {
			case foundCr:
				l.state = foundCrLf
			case foundCrLfCr:
				l.state = foundCrLfCrLf
			default:
				l.Reset()
			}
		}
		fmt.Fprint(l.logger, string(b))
	}
}
