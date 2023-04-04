package log

import "github.com/sirupsen/logrus"

// Logger log wrap for project
type Logger struct {
	log *logrus.Logger
}

// NewLogger logger constructor
func NewLogger(log *logrus.Logger) Logger {
	return Logger{
		log: log,
	}
}

// Info log with info level
func (l Logger) Info(args ...interface{}) {
	l.log.Infoln(args...)
}

// Error log with error level
func (l Logger) Error(args ...interface{}) {
	l.log.Errorln(args...)
}
