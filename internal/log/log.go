package log

type innerLog interface {
	Errorln(args ...interface{})
	Infoln(args ...interface{})
}

//go:generate mockgen -source=log.go -destination=./mock_log/log.go -package=mock_log

type logger struct {
	log innerLog
}

// Logger log wrap interface for project
type Logger interface {
	Error(args ...interface{})
	Info(args ...interface{})
}

// NewLogger logger constructor
func NewLogger(log innerLog) Logger {
	return logger{
		log: log,
	}
}

// Info innerLog with info level
func (l logger) Info(args ...interface{}) {
	l.log.Infoln(args...)
}

// Error innerLog with error level
func (l logger) Error(args ...interface{}) {
	l.log.Errorln(args...)
}
