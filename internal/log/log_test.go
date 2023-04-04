package log

import (
	"testing"

	"github.com/golang/mock/gomock"
	"tgavatar/internal/log/mock_log"
)

const logMessage = "message"

func TestLogger_Info(t *testing.T) {
	ctrl := gomock.NewController(t)
	log := mock_log.NewMocklog(ctrl)
	log.EXPECT().Infoln(logMessage)

	l := NewLogger(log)
	l.Info(logMessage)
}

func TestLogger_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	log := mock_log.NewMocklog(ctrl)
	log.EXPECT().Errorln(logMessage)

	l := NewLogger(log)
	l.Error(logMessage)
}
