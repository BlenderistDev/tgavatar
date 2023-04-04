package log

type log interface {
	Errorln(args ...interface{})
	Infoln(args ...interface{})
}

//go:generate mockgen -source=log.go -destination=./mock_log/log.go -package=mock_log

// Logger log wrap for project
type Logger struct {
	log log
}

// NewLogger logger constructor
func NewLogger(log log) Logger {
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
