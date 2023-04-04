package log

import "github.com/sirupsen/logrus"

type Logger struct {
	log *logrus.Logger
}

func NewLogger(log *logrus.Logger) Logger {
	return Logger{
		log: log,
	}
}

func (l Logger) Info(args ...interface{}) {
	l.log.Infoln(args...)
}

func (l Logger) Error(args ...interface{}) {
	l.log.Errorln(args...)
}
